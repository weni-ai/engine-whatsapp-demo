package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/weni/whatsapp-router/config"
)

func HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)

	log.Print(bodyString)
	PostToWhatsapp(bodyBytes)
}

func PostToWhatsapp(body []byte) {
	wconfig := config.GetConfig().Whatsapp

	httpClient := &http.Client{}
	reqPath := "/v1/messages"

	reqURL, _ := url.Parse(wconfig.BaseURL + reqPath)
	req := &http.Request{
		Method: "POST",
		URL:    reqURL,
		Header: map[string][]string{
			"Content-Type":  {"application/json; charset=UTF-8"},
			"Authorization": {"Bearer " + wconfig.AuthToken},
		},
		Body: ioutil.NopCloser(bytes.NewReader(body)),
	}

	res, err := httpClient.Do(req)

	if err != nil {
		log.Println(err.Error())
	}

	fmt.Println(res.Status)
}
