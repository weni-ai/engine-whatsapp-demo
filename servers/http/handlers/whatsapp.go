package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/weni/whatsapp-router/config"
	"github.com/weni/whatsapp-router/logger"
	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/services"
)

type WhatsappHandler struct {
	ContactService services.ContactService
	ChannelService services.ChannelService
}

func (h *WhatsappHandler) HandleIncomingRequests(w http.ResponseWriter, r *http.Request) {
	incomingMsg := MessagePayload{}

	if err := json.NewDecoder(r.Body).Decode(&incomingMsg); err != nil {
		logger.Error("unexpected server error - " + err.Error())
		return
	}

	incomingContact := incomingMsg.ToContact()
	if incomingContact == nil {
		logger.Error("bad request for logical error")
		return
	}

	contact, err := h.ContactService.FindContact(incomingContact)
	if err != nil {
		logger.Error(err.Error())
	}

	if contact != nil {
		channelId := contact.Channel.Hex()
		channel, err2 := h.ChannelService.FindChannelById(channelId)
		if err2 != nil {
			logger.Error(err.Error())
		}
		if channel != nil {
			jsonMsg, _ := json.Marshal(incomingMsg)
			channelUUID := channel.UUID
			RedirectRequest(r, channelUUID, string(jsonMsg))
		}

	} else {
		possibleToken := incomingMsg.Messages[0].Text.Body
		ch, err := h.ChannelService.FindChannelByToken(possibleToken)
		if err != nil {
			logger.Error(err.Error())
		}
		if ch != nil {
			incomingContact.Channel = ch.ID
			h.ContactService.CreateContact(incomingContact)

		}
	}
}

func RedirectRequest(r *http.Request, channelUUID string, msg string) {
	courierBaseURL := config.GetConfig().Server.CourierBaseURL
	url := fmt.Sprintf("%v/%v/receive", courierBaseURL, channelUUID)
	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer([]byte(msg)))

	if err != nil {
		logger.Error(err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Info(fmt.Sprintf("SENT: %v", string(body)))
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	wconfig := config.GetConfig().Whatsapp
	httpClient := &http.Client{}
	reqPath := "/v1/users/login"

	reqURL, _ := url.Parse(wconfig.BaseURL + reqPath)

	req := &http.Request{
		Method: "POST",
		URL:    reqURL,
		Header: map[string][]string{},
		Body:   r.Body,
	}

	req.SetBasicAuth(config.AppConf.Whatsapp.Username, config.AppConf.Whatsapp.Password)

	res, err := httpClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		logger.Error(err.Error())
		return
	}

	var login LoginPayload

	if err := json.NewDecoder(res.Body).Decode(&login); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		logger.Error(err.Error())
		return
	}

	newToken := login.Users[0].Token

	config.UpdateToken(newToken)
	logger.Info("Whatsapp token update")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(login)
	for k, v := range res.Header {
		w.Header().Set(k, strings.Join(v, ""))
	}
	fmt.Fprint(w, string(b))
}

type MessagePayload struct {
	Contacts []struct {
		Profile struct {
			Name string `json:"name"`
		} `json:"profile"`
		WaID string `json:"wa_id"`
	} `json:"contacts"`
	Messages []struct {
		From string `json:"from"`
		ID   string `json:"id"`
		Text struct {
			Body string `json:"body"`
		} `json:"text"`
		Timestamp string `json:"timestamp"`
		Type      string `json:"type"`
	} `json:"messages"`
}

type LoginPayload struct {
	Users []struct {
		Token        string
		ExpiresAfter string
	}
	Meta struct {
		Version   string
		ApiStatus string
	}
}

func (m *MessagePayload) ToContact() *models.Contact {
	if len(m.Messages) > 0 && len(m.Contacts) > 0 {
		return &models.Contact{
			URN:  m.Messages[0].From,
			Name: m.Contacts[0].Profile.Name,
		}
	}
	return nil
}
