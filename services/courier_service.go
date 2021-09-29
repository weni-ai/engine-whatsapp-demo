package services

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/weni/whatsapp-router/config"
	"github.com/weni/whatsapp-router/logger"
)

type CourierService interface {
	RedirectMessage(string, string) (int, error)
}

type DefaultCourierService struct {
}

func (cs DefaultCourierService) RedirectMessage(channelUUID string, msg string) (int, error) {
	courierBaseURL := config.GetConfig().Server.CourierBaseURL
	url := fmt.Sprintf("%v/%v/receive", courierBaseURL, channelUUID)
	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer([]byte(msg)))

	if err != nil {
		return resp.StatusCode, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	logger.Info(fmt.Sprintf("SENT: %v", string(body)))
	return http.StatusCreated, nil
}
