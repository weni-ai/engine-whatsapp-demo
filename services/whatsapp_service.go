package services

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/weni/whatsapp-router/config"
)

const messagePath = "/v1/messages"

type WhatsappService interface {
	SendMessage([]byte) (http.Header, io.ReadCloser, error)
}

type DefaultWhatsappService struct {
}

func NewWhatsappService() DefaultWhatsappService {
	return DefaultWhatsappService{}
}

func (ws DefaultWhatsappService) SendMessage(body []byte) (http.Header, io.ReadCloser, error) {
	wconfig := config.GetConfig().Whatsapp

	httpClient := &http.Client{}

	reqURL, _ := url.Parse(wconfig.BaseURL + messagePath)
	req := &http.Request{
		Method: "POST",
		URL:    reqURL,
		Header: map[string][]string{
			"Content-Type":  {"application/json"},
			"Accept":        {"application/json"},
			"Authorization": {"Bearer " + wconfig.AuthToken},
		},
		Body: ioutil.NopCloser(bytes.NewReader(body)),
	}

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, nil, err
	}
	if res.StatusCode == 401 {
		return nil, nil, errors.New(res.Status)
	}
	return res.Header, res.Body, nil
}
