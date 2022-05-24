package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v11"
	"github.com/weni/whatsapp-router/config"
	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/services"
	"github.com/weni/whatsapp-router/utils"
)

var kkClient gocloak.GoCloak

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

func KeycloackAuth(next http.HandlerFunc) http.HandlerFunc {
	if kkClient == nil {
		kkClient = NewKeycloakClient()
	}
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		authSplit := strings.Split(authorization, " ")
		var token string
		token = authSplit[0]
		if len(authSplit) > 1 {
			token = authSplit[1]
		}

		ctx := context.Background()
		_, err := kkClient.GetUserInfo(ctx, token, config.GetConfig().OIDC.Realm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func NewKeycloakClient() gocloak.GoCloak {
	return gocloak.NewClient(config.GetConfig().OIDC.Host)
}
