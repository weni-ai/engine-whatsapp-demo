package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/services"
	"github.com/weni/whatsapp-router/utils"
)

type IntegrationsHandler struct {
	ChannelService services.ChannelService
}

func (h *IntegrationsHandler) HandleCreateChannel(w http.ResponseWriter, r *http.Request) {
	ch := &models.Channel{}
	err := json.NewDecoder(r.Body).Decode(ch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ch.Token = utils.GenToken()
	_, err = h.ChannelService.CreateChannelDefault(ch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"token":"%s"}`, ch.Token)))
}
