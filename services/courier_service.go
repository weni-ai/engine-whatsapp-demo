package services

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/weni/whatsapp-router/config"
	"github.com/weni/whatsapp-router/logger"
)

type CourierService interface {
	RedirectMessage(string, string) (int, error)
	RedirectMessageWac(string, *http.Request, string) (int, error)
}

type DefaultCourierService struct {
}

func (cs DefaultCourierService) RedirectMessage(channelUUID string, msg string) (int, error) {
	courierBaseURL := config.GetConfig().App.CourierBaseURL
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

func (cs DefaultCourierService) RedirectMessageWac(channelUUID string, r *http.Request, token string) (int, error) {
	courierBaseURL := config.GetConfig().App.CloudURL
	// url := fmt.Sprintf("%v/%v/receive", courierBaseURL, channelUUID)
	url := fmt.Sprintf("%v/receive", courierBaseURL)

	fmt.Println("URL", url)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}

	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	req.Header.Add("X-Router-Token", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resp.StatusCode, err
	}

	body, _ = io.ReadAll(resp.Body)
	log.Println("Status Code:", resp.StatusCode)
	log.Println("Body:", string(body))
	log.Printf("Response: %v, Body: %v, Response Body: %v", resp, string(body), body)
	logger.Error(fmt.Sprintf("Response: %v, Body: %v, Response Body: %v", resp, string(body), body))
	resp.Body.Close()

	return resp.StatusCode, nil
}

func NewCourierService() DefaultCourierService {
	return DefaultCourierService{}
}
