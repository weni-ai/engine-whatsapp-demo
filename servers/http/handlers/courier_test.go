package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "github.com/weni/whatsapp-router/mocks/services"
)

func TestHandleMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	payload := `{"to":"5582988887777","type":"text","text":{"body":"hello"}}`
	payloadBytes := []byte(payload)

	mockWhatsappService := mocks.NewMockWhatsappService(ctrl)
	mockWhatsappService.EXPECT().SendMessage(payloadBytes).Return(
		http.Header{
			"content-type": {"application/json"},
		},
		ioutil.NopCloser(bytes.NewReader([]byte(`{"messages":{"id":"gBEGVYKZRIIyAgmiTgezkroUL2Q"}],"meta":{"api_status":"stable","version":"2.35.2"}}`))),
		nil,
	)

	ch := CourierHandler{mockWhatsappService}

	router := chi.NewRouter()
	router.Post("/v1/messages", ch.HandleSendMessage)

	requestBody := strings.NewReader(payload)
	request, _ := http.NewRequest(
		http.MethodPost,
		"/v1/messages",
		requestBody,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, response.Code, 201)
}
