package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/weni/whatsapp-router/config"
	"github.com/weni/whatsapp-router/utils"
)

const (
	messagePath = "/v1/messages"
	loginPath   = "/v1/users/login"
	healthPath  = "/v1/health"
	mediaPath   = "/v1/media/"
)

type WhatsappService interface {
	SendMessage([]byte) (http.Header, io.ReadCloser, error)
	Login() (*http.Response, error)
	Health() (*http.Response, error)
	GetMedia(http.Header, string) (*http.Response, error)
	PostMedia(http.Header, io.ReadCloser) (*http.Response, error)
}

type DefaultWhatsappService struct {
}

func NewWhatsappService() DefaultWhatsappService {
	return DefaultWhatsappService{}
}

func (ws DefaultWhatsappService) SendMessage(body []byte) (http.Header, io.ReadCloser, error) {
	wconfig := config.GetConfig().Whatsapp

	httpClient := utils.GetHTTPClient()

	reqURL, _ := url.Parse(wconfig.BaseURL + messagePath)
	req := &http.Request{
		Method: "POST",
		URL:    reqURL,
		Header: map[string][]string{
			"Content-Type":  {"application/json"},
			"Accept":        {"application/json"},
			"Authorization": {"Bearer " + config.GetAuthToken()},
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

func (ws DefaultWhatsappService) Login() (*http.Response, error) {
	wconfig := config.GetConfig().Whatsapp
	httpClient := utils.GetHTTPClient()
	reqURL, _ := url.Parse(wconfig.BaseURL + loginPath)

	req := &http.Request{
		Method: "POST",
		URL:    reqURL,
		Header: map[string][]string{},
		Body:   nil,
	}

	req.SetBasicAuth(wconfig.Username, wconfig.Password)
	return httpClient.Do(req)
}

func (ws DefaultWhatsappService) Health() (*http.Response, error) {
	wconfig := config.GetConfig().Whatsapp
	httpClient := utils.GetHTTPClient()
	reqURL, _ := url.Parse(wconfig.BaseURL + healthPath)

	req := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Content-Type":  {"application/json"},
			"Accept":        {"application/json"},
			"Authorization": {"Bearer " + config.GetAuthToken()},
		},
		Body: nil,
	}
	return httpClient.Do(req)
}

func (ws DefaultWhatsappService) GetMedia(header http.Header, mediaID string) (*http.Response, error) {
	wconfig := config.GetConfig().Whatsapp
	httpClient := utils.GetHTTPClient()
	req, err := http.NewRequest(
		"GET",
		wconfig.BaseURL+mediaPath+mediaID,
		nil,
	)
	if err != nil {
		return nil, err
	}
	utils.CopyHeader(req.Header, header)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetAuthToken()))
	return httpClient.Do(req)
}

func (ws DefaultWhatsappService) PostMedia(header http.Header, body io.ReadCloser) (*http.Response, error) {
	wconfig := config.GetConfig().Whatsapp
	httpClient := utils.GetHTTPClient()
	req, err := http.NewRequest(
		"POST",
		wconfig.BaseURL+mediaPath,
		body,
	)
	defer body.Close()
	if err != nil {
		return nil, err
	}
	utils.CopyHeader(req.Header, header)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetAuthToken()))
	return httpClient.Do(req)
}

type LoginWhatsapp struct {
	Users []struct {
		Token        string
		ExpiresAfter string
	}
	Meta struct {
		Version   string
		ApiStatus string
	}
}
