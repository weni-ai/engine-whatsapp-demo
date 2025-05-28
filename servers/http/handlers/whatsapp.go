package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/weni/whatsapp-router/config"
	"github.com/weni/whatsapp-router/logger"
	"github.com/weni/whatsapp-router/metric"
	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/services"
	"github.com/weni/whatsapp-router/utils"
)

var confirmationMessage = config.GetConfig().Whatsapp.WelcomeMessage

const tokenPrefix = "weni-demo"

type WhatsappHandler struct {
	ContactService  services.ContactService
	ChannelService  services.ChannelService
	CourierService  services.CourierService
	WhatsappService services.WhatsappService
	ConfigService   services.ConfigService
	Metrics         *metric.Service
}

func (h *WhatsappHandler) HandleIncomingRequests(w http.ResponseWriter, r *http.Request) {
	incomingWebhookEvent, err := ioutil.ReadAll(io.LimitReader(r.Body, 1000000))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(incomingWebhookEvent))
	defer r.Body.Close()
	if err != nil {
		logger.Error(fmt.Sprintf("unable to read request body: %s", err))
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, err.Error())
		return
	}

	payload, err := h.parseEventPayload(incomingWebhookEvent)
	if err != nil {
		logger.Error(fmt.Sprintf("unable to parse request body: %s", err))
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, err.Error())
		return
	}

	if len(payload.Messages) <= 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	contact := h.getOrCreateContact(payload)
	textMessage := h.getMessageText(payload)

	// Handle token-based channel registration/update
	if textMessage != "" && strings.Contains(textMessage, tokenPrefix) {
		h.handleTokenMessage(w, contact, textMessage, payload)
		return
	}

	// Handle regular message routing
	if contact != nil {
		h.routeMessage(w, contact, incomingWebhookEvent)
		return
	}

	//returning status ok to avoid retry send mechanisms if contact not exists or token is not valid
	logger.Debug("contact not found and token not valid")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, errors.New("contact not found and token not valid"))
}

func (h *WhatsappHandler) HandleIncomingRequestsWac(w http.ResponseWriter, r *http.Request) {
	incomingWebhookEvent, err := ioutil.ReadAll(io.LimitReader(r.Body, 1000000))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(incomingWebhookEvent))
	defer r.Body.Close()
	if err != nil {
		logger.Error(fmt.Sprintf("unable to read request body: %s", err))
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, err.Error())
		return
	}

	payload, err := h.parseEventPayloadWac(incomingWebhookEvent)
	if err != nil {
		logger.Error(fmt.Sprintf("unable to parse request body: %s", err))
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, err.Error())
		return
	}

	if len(payload.Entry[0].Changes[0].Value.Messages) <= 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	contact := h.getOrCreateContactWac(payload)

	textMessage := h.getMessageTextWac(payload)

	// Handle token-based channel registration/update
	if textMessage != "" && strings.Contains(textMessage, tokenPrefix) {
		h.handleTokenMessageWac(w, contact, textMessage, payload)
		return
	}

	// Handle regular message routing
	if contact != nil {
		h.routeMessageWac(w, r, contact, incomingWebhookEvent)
		return
	}

	//returning status ok to avoid retry send mechanisms if contact not exists or token is not valid
	logger.Debug("contact not found and token not valid")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, errors.New("contact not found and token not valid"))
}

// Helper functions for HandleIncomingRequests

func (h *WhatsappHandler) parseEventPayload(data []byte) (*eventPayload, error) {
	payload := &eventPayload{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func (h *WhatsappHandler) parseEventPayloadWac(data []byte) (*eventPayloadWac, error) {
	payload := &eventPayloadWac{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func (h *WhatsappHandler) getOrCreateContact(payload *eventPayload) *models.Contact {
	cName := ""
	if len(payload.Contacts) > 0 {
		cName = payload.Contacts[0].Profile.Name
	}
	incomingContact := &models.Contact{
		URN:  payload.Messages[0].From,
		Name: cName,
	}

	contact, err := h.ContactService.FindContact(incomingContact)
	if err != nil {
		logger.Debug(err.Error())
	}

	return contact
}

func (h *WhatsappHandler) getOrCreateContactWac(payload *eventPayloadWac) *models.Contact {
	cName := ""
	if len(payload.Entry[0].Changes[0].Value.Contacts) > 0 {
		cName = payload.Entry[0].Changes[0].Value.Contacts[0].Profile.Name
	}
	incomingContact := &models.Contact{
		URN:  payload.Entry[0].Changes[0].Value.Contacts[0].WaID,
		Name: cName,
	}

	contact, err := h.ContactService.FindContact(incomingContact)
	if err != nil {
		logger.Debug(err.Error())
	}

	return contact
}

func (h *WhatsappHandler) getMessageText(payload *eventPayload) string {
	if payload.Messages[0].Type == "text" {
		return payload.Messages[0].Text.Body
	}
	return ""
}

func (h *WhatsappHandler) getMessageTextWac(payload *eventPayloadWac) string {
	if payload.Entry[0].Changes[0].Value.Messages[0].Type == "text" {
		return payload.Entry[0].Changes[0].Value.Messages[0].Text.Body
	}
	return ""
}

func (h *WhatsappHandler) handleTokenMessage(w http.ResponseWriter, contact *models.Contact, token string, payload *eventPayload) {
	channelFromToken, err := h.ChannelService.FindChannelByToken(token)
	if err != nil {
		logger.Debug(err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}

	if channelFromToken == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Create a new contact or update existing one
	if contact != nil {
		h.updateExistingContact(w, contact, channelFromToken)
	} else {
		// Create contact from payload
		cName := ""
		if len(payload.Contacts) > 0 {
			cName = payload.Contacts[0].Profile.Name
		}

		incomingContact := &models.Contact{
			URN:     payload.Messages[0].From,
			Name:    cName,
			Channel: channelFromToken.ID,
		}

		h.createNewContact(w, incomingContact, channelFromToken)
	}
}

func (h *WhatsappHandler) handleTokenMessageWac(w http.ResponseWriter, contact *models.Contact, token string, payload *eventPayloadWac) {
	channelFromToken, err := h.ChannelService.FindChannelByToken(token)
	if err != nil {
		logger.Debug(err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}

	if channelFromToken == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Create a new contact or update existing one
	if contact != nil {
		h.updateExistingContactWac(w, contact, channelFromToken)
	} else {
		// Create contact from payload
		cName := ""
		if len(payload.Entry[0].Changes[0].Value.Contacts) > 0 {
			cName = payload.Entry[0].Changes[0].Value.Contacts[0].Profile.Name
		}

		incomingContact := &models.Contact{
			URN:     payload.Entry[0].Changes[0].Value.Messages[0].From,
			Name:    cName,
			Channel: channelFromToken.ID,
		}

		h.createNewContactWac(w, incomingContact, channelFromToken)
	}
}

func (h *WhatsappHandler) updateExistingContact(w http.ResponseWriter, contact *models.Contact, newChannel *models.Channel) {
	lastContactChannel, err := h.ChannelService.FindChannelById(contact.Channel.Hex())
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contact.Channel = newChannel.ID
	_, err = h.ContactService.UpdateContact(contact)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, b, err := h.sendTokenConfirmation(contact)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ := ioutil.ReadAll(b)
	b.Close()
	logger.Debug(string(body))
	w.WriteHeader(http.StatusOK)

	// Update metrics
	contactActivatedMetricDec := metric.NewContactActivated(lastContactChannel.UUID)
	h.Metrics.DecContactActivated(contactActivatedMetricDec)
	contactActivatedMetricInc := metric.NewContactActivated(newChannel.UUID)
	h.Metrics.IncContactActivated(contactActivatedMetricInc)
	contactActivation := metric.NewContactActivation(newChannel.UUID)
	h.Metrics.SaveContactActivation(contactActivation)
}

func (h *WhatsappHandler) updateExistingContactWac(w http.ResponseWriter, contact *models.Contact, newChannel *models.Channel) {
	lastContactChannel, err := h.ChannelService.FindChannelById(contact.Channel.Hex())
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contact.Channel = newChannel.ID
	_, err = h.ContactService.UpdateContact(contact)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, b, err := h.sendTokenConfirmationWac(contact)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ := ioutil.ReadAll(b)
	b.Close()
	logger.Debug(string(body))
	w.WriteHeader(http.StatusOK)

	// Update metrics
	contactActivatedMetricDec := metric.NewContactActivated(lastContactChannel.UUID)
	h.Metrics.DecContactActivated(contactActivatedMetricDec)
	contactActivatedMetricInc := metric.NewContactActivated(newChannel.UUID)
	h.Metrics.IncContactActivated(contactActivatedMetricInc)
	contactActivation := metric.NewContactActivation(newChannel.UUID)
	h.Metrics.SaveContactActivation(contactActivation)
}

func (h *WhatsappHandler) createNewContact(w http.ResponseWriter, contact *models.Contact, channel *models.Channel) {
	_, err := h.ContactService.CreateContact(contact)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, b, err := h.sendTokenConfirmation(contact)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, _ := ioutil.ReadAll(b)
	b.Close()
	logger.Debug(string(body))
	w.WriteHeader(http.StatusOK)

	// Update metrics
	contactActivation := metric.NewContactActivation(channel.UUID)
	h.Metrics.SaveContactActivation(contactActivation)
	contactActivated := metric.NewContactActivated(channel.UUID)
	h.Metrics.IncContactActivated(contactActivated)
}

func (h *WhatsappHandler) createNewContactWac(w http.ResponseWriter, contact *models.Contact, channel *models.Channel) {
	_, err := h.ContactService.CreateContact(contact)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, b, err := h.sendTokenConfirmationWac(contact)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, _ := ioutil.ReadAll(b)
	b.Close()
	logger.Debug(string(body))
	w.WriteHeader(http.StatusOK)

	// Update metrics
	contactActivation := metric.NewContactActivation(channel.UUID)
	h.Metrics.SaveContactActivation(contactActivation)
	contactActivated := metric.NewContactActivated(channel.UUID)
	h.Metrics.IncContactActivated(contactActivated)
}

func (h *WhatsappHandler) routeMessage(w http.ResponseWriter, contact *models.Contact, message []byte) {
	channelId := contact.Channel.Hex()
	channel, err := h.ChannelService.FindChannelById(channelId)
	if err != nil {
		logger.Debug(err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}

	if channel == nil {
		logger.Debug("channel not found")
		w.WriteHeader(http.StatusOK)
		return
	}

	channelUUID := channel.UUID
	status, err := h.CourierService.RedirectMessage(channelUUID, string(message))
	if err != nil {
		logger.Debug(err.Error())
		w.WriteHeader(status)
		fmt.Fprint(w, err)
		return
	}

	if status >= 400 {
		logger.Debug(fmt.Sprintf("message redirect with status %d for channel %s", status, channelUUID))
		w.WriteHeader(http.StatusOK)
		return
	}

	cmm := metric.NewContactMessage(channelUUID)
	h.Metrics.SaveContactMessage(cmm)
	w.WriteHeader(http.StatusOK)
}

func (h *WhatsappHandler) routeMessageWac(w http.ResponseWriter, r *http.Request, contact *models.Contact, message []byte) {
	channelId := contact.Channel.Hex()
	channel, err := h.ChannelService.FindChannelById(channelId)
	if err != nil {
		logger.Debug(err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}

	if channel == nil {
		logger.Debug("channel not found")
		w.WriteHeader(http.StatusOK)
		return
	}

	channelUUID := channel.UUID
	status, err := h.CourierService.RedirectMessageWac(channelUUID, r)
	if err != nil {
		logger.Debug(err.Error())
		w.WriteHeader(status)
		fmt.Fprint(w, err)
		return
	}

	if status >= 400 {
		logger.Debug(fmt.Sprintf("message redirect with status %d for channel %s", status, channelUUID))
		w.WriteHeader(http.StatusOK)
		return
	}

	cmm := metric.NewContactMessage(channelUUID)
	h.Metrics.SaveContactMessage(cmm)
	w.WriteHeader(http.StatusOK)
}

func (h *WhatsappHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	res, err := h.WhatsappService.Login()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		logger.Error(err.Error())
		return
	}

	var login services.LoginWhatsapp

	bdBytes, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	bdString := string(bdBytes)

	if res.StatusCode != 200 {
		w.WriteHeader(res.StatusCode)
		w.Write(bdBytes)
		logger.Error(fmt.Sprintf("Couldn't update token: %s", bdString))
		return
	}

	if err := json.Unmarshal(bdBytes, &login); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		logger.Error(err.Error())
		return
	}

	newToken := login.Users[0].Token

	config.UpdateAuthToken(newToken)

	h.ConfigService.CreateOrUpdate(&models.Config{Token: newToken})

	utils.CopyHeader(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)
	w.Write(bdBytes)
}

func (h *WhatsappHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	res, err := h.WhatsappService.Health()
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.CopyHeader(w.Header(), res.Header)
	io.Copy(w, res.Body)
	res.Body.Close()
}

func (h *WhatsappHandler) HandleGetMedia(w http.ResponseWriter, r *http.Request) {
	mediaID := chi.URLParam(r, "mediaID")
	res, err := h.WhatsappService.GetMedia(r.Header, mediaID)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.CopyHeader(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
	res.Body.Close()
}

func (h *WhatsappHandler) HandlePostMedia(w http.ResponseWriter, r *http.Request) {
	res, err := h.WhatsappService.PostMedia(r.Header, r.Body)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.CopyHeader(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
	res.Body.Close()
}

func (h *WhatsappHandler) sendTokenConfirmation(contact *models.Contact) (http.Header, io.ReadCloser, error) {
	urn := contact.URN
	payload := fmt.Sprintf(
		`{"to":"%s","type":"text","text":{"body":"%s"}}`,
		urn,
		confirmationMessage,
	)
	payloadBytes := []byte(payload)

	return h.WhatsappService.SendMessage(payloadBytes)
}

func (h *WhatsappHandler) sendTokenConfirmationWac(contact *models.Contact) (http.Header, io.ReadCloser, error) {
	urn := contact.URN
	payload := fmt.Sprintf(
		`{"messaging_product":"whatsapp","recipient_type":"individual","to":"%s","type":"text","text":{"body":"%s"}}`,
		urn,
		confirmationMessage,
	)
	payloadBytes := []byte(payload)

	return h.WhatsappService.SendMessageWac(payloadBytes)
}

type eventPayload struct {
	Contacts []struct {
		Profile struct {
			Name string `json:"name"`
		} `json:"profile"`
		WaID string `json:"wa_id"`
	} `json:"contacts"`
	Messages []struct {
		From      string `json:"from"      validate:"required"`
		ID        string `json:"id"        validate:"required"`
		Timestamp string `json:"timestamp" validate:"required"`
		Type      string `json:"type"      validate:"required"`
		Text      struct {
			Body string `json:"body"`
		} `json:"text"`
	}
}

type eventPayloadWac struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					From      string `json:"from"      validate:"required"`
					ID        string `json:"id"        validate:"required"`
					Timestamp string `json:"timestamp" validate:"required"`
					Type      string `json:"type"      validate:"required"`
					Text      struct {
						Body string `json:"body"`
					} `json:"text"`
				} `json:"messages"`
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
}
