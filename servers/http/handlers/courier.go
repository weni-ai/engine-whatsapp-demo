package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/weni/whatsapp-router/logger"
	"github.com/weni/whatsapp-router/services"
)

type CourierHandler struct {
	WhatsappService services.WhatsappService
}

func (c *CourierHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	logger.Debug(bodyString)
	header, body, err := c.WhatsappService.SendMessage(bodyBytes)

	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, nil)
		return
	}

	for k, v := range header {
		w.Header().Set(k, strings.Join(v, ""))
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, body)
}
