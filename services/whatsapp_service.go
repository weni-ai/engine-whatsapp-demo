package services

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/weni/whatsapp-router/config"
)

type WhatsappService interface {
	SendMessage([]byte) (http.Header, io.ReadCloser, error)
}

type DefaultWhatsappService struct {
}

func (ws DefaultWhatsappService) SendMessage(body []byte) (http.Header, io.ReadCloser, error) {
	wconfig := config.GetConfig().Whatsapp

	httpClient := &http.Client{}
	reqPath := "/v1/messages"

	reqURL, _ := url.Parse(wconfig.BaseURL + reqPath)
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
	return res.Header, res.Body, nil
}

func NewWhatsappService() DefaultWhatsappService {
	return DefaultWhatsappService{}
}
