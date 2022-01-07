package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/weni/whatsapp-router/logger"
	"github.com/weni/whatsapp-router/services"
	"github.com/weni/whatsapp-router/utils"
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

	utils.CopyHeader(w.Header(), header)
	b, _ := ioutil.ReadAll(body)
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}
