package services

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/weni/whatsapp-router/config"
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

	return resp.StatusCode, nil
}

func NewCourierService() DefaultCourierService {
	return DefaultCourierService{}
}
