package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
		log.Print(err)
		return
	}

	incomingContact := incomingMsg.ToContact()
	if incomingContact == nil {
		return
	}

	contact, err := h.ContactService.FindContact(incomingContact)
	if err != nil {
		log.Print(err)
	}

	if contact != nil {
		channelId := contact.Channel.Hex()
		channel, err2 := h.ChannelService.FindChannelById(channelId)
		if err2 != nil {
			log.Print(err)
		}
		if channel != nil {
			//TODO redirect message to channel handler
			jsonMsg, _ := json.Marshal(incomingMsg)

			// channelUUID := "624476b1-032f-46d8-becd-d0f14e23bfbb"
			channelUUID := channel.UUID

			RedirectRequest(r, channelUUID, string(jsonMsg))
		} else {
			//TODO nothing to do
		}

	} else {
		possibleToken := incomingMsg.Messages[0].Text.Body
		ch, err := h.ChannelService.FindChannelByToken(possibleToken)
		if err != nil {
			log.Print(err)
		}
		if ch != nil {
			incomingContact.Channel = ch.ID
			h.ContactService.CreateContact(incomingContact)

		}
	}
}

const courierBaseURL = "https://f3e9-179-235-152-98.ngrok.io/c/wa"

//TODO finish this
func RedirectRequest(r *http.Request, channelUUID string, msg string) {
	resp, err := http.Post(
		fmt.Sprintf("%v/%v/receive", courierBaseURL, channelUUID),
		"application/json",
		bytes.NewBuffer([]byte(msg)))
	if err != nil {
		fmt.Printf("err %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Printf("Body: %s", body)
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

func (m *MessagePayload) ToContact() *models.Contact {
	if len(m.Messages) > 0 && len(m.Contacts) > 0 {
		return &models.Contact{
			URN:  m.Messages[0].From,
			Name: m.Contacts[0].Profile.Name,
		}
	}
	return nil
}
