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
		log.Print("unexpected server error - " + err.Error())
		return
	}

	incomingContact := incomingMsg.ToContact()
	if incomingContact == nil {
		log.Println("bad request for logical error")
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
			log.Println(err.Error())
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
		fmt.Println(err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err.Error())
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
