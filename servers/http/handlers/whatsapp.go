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

var welcomeMessage = config.GetConfig().Whatsapp.WelcomeMessage

const tokenPrefix = "weni-demo"

type WhatsappHandler struct {
	ContactService  services.ContactService
	ChannelService  services.ChannelService
	CourierService  services.CourierService
	WhatsappService services.WhatsappService
	ConfigService   services.ConfigService
	Metrics         *metric.Service
	FlowsService    services.FlowsService
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
	} else if payload.Messages[0].Type == "interactive" {
		textMessage = payload.Messages[0].Interactive.ButtonReply.Title
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
				flows := &models.Flows{
					Channel: channelFromToken.UUID,
				}

				fl, err := h.FlowsService.FindFlows(flows)
				if err != nil {
					logger.Debug(err.Error())
				}

				var b io.ReadCloser
				if fl != nil {
					_, b, err = h.sendFlowsChoice(channelFromToken, contact, fl)
				} else {
					_, b, err = h.sendTokenConfirmation(contact)
				}

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
				contact, err := h.ContactService.CreateContact(incomingContact)
				if err != nil {
					logger.Error(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				flows := &models.Flows{
					Channel: channelFromToken.UUID,
				}

				fl, err := h.FlowsService.FindFlows(flows)
				if err != nil {
					logger.Debug(err.Error())
				}

				var b io.ReadCloser
				if fl != nil {
					_, b, err = h.sendFlowsChoice(channelFromToken, contact, fl)
				} else {
					_, b, err = h.sendTokenConfirmation(contact)
				}

				if err != nil {
					logger.Error(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				body, _ := ioutil.ReadAll(b)
				b.Close()
				logger.Debug(string(body))

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
			hasKeyword := false
			if channel != nil {

				flows := &models.Flows{
					Channel: channel.UUID,
				}

				fls, err := h.FlowsService.FindFlows(flows)
				if err != nil {
					logger.Debug(err.Error())
				}
				var keyword string
				for _, fl := range fls.FlowsStarts {
					if textMessage == fl.Name {
						keyword = fl.Keyword
						hasKeyword = true
						break
					}
				}
				channelUUID := channel.UUID
				var status int
				if hasKeyword {
					payloadInteractive := &eventPayload{
						Contacts: payload.Contacts,
					}

					payloadInteractive.Messages = append(payloadInteractive.Messages, struct {
						From      string "json:\"from\"      validate:\"required\""
						ID        string "json:\"id\"        validate:\"required\""
						Timestamp string "json:\"timestamp\" validate:\"required\""
						Type      string "json:\"type\"      validate:\"required\""
						Text      struct {
							Body string "json:\"body\""
						} "json:\"text\""
						Interactive struct {
							ButtonReply struct {
								ID    string "json:\"id\""
								Title string "json:\"title\""
							} "json:\"button_reply\""
							Type string "json:\"type\""
						} "json:\"interactive,omitempty\""
					}{From: payload.Messages[0].From, ID: payload.Messages[0].ID, Text: struct {
						Body string "json:\"body\""
					}{keyword}, Timestamp: payload.Messages[0].Timestamp, Type: "text"})

					payloadBytes, err := json.Marshal(payloadInteractive)
					if err != nil {
						logger.Debug(err.Error())
					}
					status, err = h.CourierService.RedirectMessage(channelUUID, string(payloadBytes))
					if err != nil {
						logger.Debug(err.Error())
						w.WriteHeader(status)
						fmt.Fprint(w, err)
						return
					}
				} else {
					status, err = h.CourierService.RedirectMessage(channelUUID, string(incomingWebhookEvent))
					if err != nil {
						logger.Debug(err.Error())
						w.WriteHeader(status)
						fmt.Fprint(w, err)
						return
					}
				}

				if status >= 400 {
					logger.Debug(fmt.Sprintf("message redirect with status %d for channel %s", status, channelUUID))
					return
				}
				cmm := metric.NewContactMessage(channelUUID)
				h.Metrics.SaveContactMessage(cmm)
				w.WriteHeader(http.StatusOK)
				return
			}
			logger.Debug("channel not found")
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	//returning status ok to avoid retry send mechanisms if contact not exists or token is not valid
	logger.Debug("contact not found and token not valid")
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

func (h *WhatsappHandler) sendFlowsChoice(channel *models.Channel, contact *models.Contact, fl *models.Flows) (http.Header, io.ReadCloser, error) {
	welcomeMessageFlows := "Olá, bem vindo ao WhatsApp Demo, escolha um dos fluxos abaixo para iniciar."
	urn := contact.URN

	payload := fmt.Sprintf(
		`{"to":"%s","type":"interactive","interactive":{"type":"button","body":{"text": "%s"},`,
		urn,
		welcomeMessageFlows,
	)

	if len(fl.FlowsStarts) > 0 {
		payload = payload + `"action":{"buttons":[`

		for i, f := range fl.FlowsStarts {
			payload = payload + fmt.Sprintf(`{"type": "reply","reply": {"id": "%s","title": "%s"}}`, f.Name, f.Name)
			if i != len(fl.FlowsStarts)-1 {
				payload = payload + `,`
			}
		}
		payload = payload + `]}`
	}

	payload = payload + `}}`
	payloadBytes := []byte(payload)

	return h.WhatsappService.SendMessage(payloadBytes)
}

func (h *WhatsappHandler) sendTokenConfirmation(contact *models.Contact) (http.Header, io.ReadCloser, error) {
	urn := contact.URN
	payload := fmt.Sprintf(
		`{"to":"%s","type":"text","text":{"body":"%s"}}`,
		urn,
		welcomeMessage,
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
		Interactive struct {
			ButtonReply struct {
				ID    string `json:"id"`
				Title string `json:"title"`
			} `json:"button_reply"`
			Type string `json:"type"`
		} `json:"interactive,omitempty"`
	}
}
