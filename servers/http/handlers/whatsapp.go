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

const confirmationMessage = "Token válido, Whatsapp demo está pronto para sua utilização"

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

	payload := &eventPayload{}
	if err = json.Unmarshal(incomingWebhookEvent, &payload); err != nil {
		logger.Error(fmt.Sprintf("unable to parse request body: %s", err))
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, err.Error())
		return
	}

	if len(payload.Messages) <= 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

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

	textMessage := ""
	if payload.Messages[0].Type == "text" {
		textMessage = payload.Messages[0].Text.Body
	}

	if textMessage != "" && strings.Contains(textMessage, tokenPrefix) {
		channelFromToken, err := h.ChannelService.FindChannelByToken(textMessage)
		if err != nil {
			logger.Debug(err.Error())
		}
		if channelFromToken != nil {
			incomingContact.Channel = channelFromToken.ID
			if contact != nil {
				lastContactChannel, err := h.ChannelService.FindChannelById(contact.Channel.Hex())
				if err != nil {
					logger.Error(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				contact.Channel = channelFromToken.ID
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

				contactActivatedMetricDec := metric.NewContactActivated(lastContactChannel.UUID)
				h.Metrics.DecContactActivated(contactActivatedMetricDec)
				contactActivatedMetricInc := metric.NewContactActivated(lastContactChannel.UUID)
				h.Metrics.IncContactActivated(contactActivatedMetricInc)
				contactActivation := metric.NewContactActivation(channelFromToken.UUID)
				h.Metrics.SaveContactActivation(contactActivation)

				return
			} else {
				_, err := h.ContactService.CreateContact(incomingContact)
				if err != nil {
					logger.Error(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				_, b, err := h.sendTokenConfirmation(incomingContact)
				if err != nil {
					logger.Error(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				body, _ := ioutil.ReadAll(b)
				b.Close()
				logger.Debug(string(body))
				w.WriteHeader(http.StatusOK)

				contactActivation := metric.NewContactActivation(channelFromToken.UUID)
				h.Metrics.SaveContactActivation(contactActivation)
				contactActivated := metric.NewContactActivated(channelFromToken.UUID)
				h.Metrics.IncContactActivated(contactActivated)
				return
			}
		}
	} else {
		if contact != nil {
			channelId := contact.Channel.Hex()
			channel, err := h.ChannelService.FindChannelById(channelId)
			if err != nil {
				logger.Debug(err.Error())
			}
			if channel != nil {
				channelUUID := channel.UUID
				status, err := h.CourierService.RedirectMessage(channelUUID, string(incomingWebhookEvent))
				if err != nil {
					logger.Debug(err.Error())
					w.WriteHeader(status)
					fmt.Fprint(w, err)
					return
				}
				cmm := metric.NewContactMessage(channelUUID)
				h.Metrics.SaveContactMessage(cmm)
			}
		}
	}

	//returning status ok to avoid retry send mechanisms if contact not exists or token is not valid
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, errors.New("contact not found and token not valid"))
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
