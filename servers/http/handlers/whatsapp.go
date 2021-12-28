package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/weni/whatsapp-router/config"
	"github.com/weni/whatsapp-router/logger"
	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/services"
)

const confirmationMessage = "Token válido, Whatsapp demo está pronto para sua utilização"

const tokenPrefix = "weni-demo"

type WhatsappHandler struct {
	ContactService  services.ContactService
	ChannelService  services.ChannelService
	CourierService  services.CourierService
	WhatsappService services.WhatsappService
}

func (h *WhatsappHandler) HandleIncomingRequests(w http.ResponseWriter, r *http.Request) {
	incomingMsg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error("unexpected server error - " + err.Error())
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, err.Error())
		return
	}

	incomingContact := parseToContact(string(incomingMsg))
	if incomingContact == nil {
		err := errors.New("request without being from a contact")
		logger.Debug(fmt.Sprintf("%v: %v", err.Error(), string(incomingMsg)))
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, err.Error())
		return
	}

	contact, err := h.ContactService.FindContact(incomingContact)
	if err != nil {
		logger.Debug(err.Error())
	}

	textMessage := extractTextMessage(string(incomingMsg))

	if textMessage != "" {
		if strings.Contains(textMessage, tokenPrefix) {
			channelFromToken, err := h.ChannelService.FindChannelByToken(textMessage)
			if err != nil {
				logger.Debug(err.Error())
			}
			if channelFromToken != nil {
				incomingContact.Channel = channelFromToken.ID
				if contact != nil {
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
					} else {
						body, _ := ioutil.ReadAll(b)
						logger.Debug(string(body))
						w.WriteHeader(http.StatusOK)
						return
					}
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
					} else {
						body, _ := ioutil.ReadAll(b)
						logger.Debug(string(body))
						w.WriteHeader(http.StatusOK)
						return
					}
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
					status, err := h.CourierService.RedirectMessage(channelUUID, string(incomingMsg))
					if err != nil {
						logger.Debug(err.Error())
						w.WriteHeader(status)
						fmt.Fprint(w, err)
					}
				}
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
	log.Println(bdString)

	if err := json.Unmarshal(bdBytes, &login); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		logger.Error(err.Error())
		return
	}

	newToken := login.Users[0].Token

	config.UpdateAuthToken(newToken)
	logger.Info("Whatsapp token update")
	w.WriteHeader(http.StatusOK)
	for k, v := range res.Header {
		w.Header().Set(k, strings.Join(v, ""))
	}
	fmt.Fprint(w, bdString)
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

func parseToContact(m string) *models.Contact {
	name := extractName(m)
	number := extractNumber(m)
	if name != "" && number != "" {
		return &models.Contact{
			URN:  number,
			Name: name,
		}
	}
	return nil
}

func extractName(m string) string {
	var result map[string][]map[string]map[string]interface{}
	json.Unmarshal([]byte(m), &result)
	if result["contacts"] != nil {
		return result["contacts"][0]["profile"]["name"].(string)
	}
	return ""
}

func extractNumber(m string) string {
	var result map[string][]map[string]interface{}
	json.Unmarshal([]byte(m), &result)
	if result["messages"] != nil {
		return result["messages"][0]["from"].(string)
	}
	return ""
}

func extractTextMessage(m string) string {
	var result map[string][]map[string]map[string]interface{}
	json.Unmarshal([]byte(m), &result)
	if result["messages"] != nil {
		return result["messages"][0]["text"]["body"].(string)
	}
	return ""
}
