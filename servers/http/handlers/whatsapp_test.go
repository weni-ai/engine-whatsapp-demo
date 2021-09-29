package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	mocks "github.com/weni/whatsapp-router/mocks/services"
	"github.com/weni/whatsapp-router/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestContactTokenConfirmation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannelService := mocks.NewMockChannelService(ctrl)
	dummyChannel := models.Channel{
		primitive.NilObjectID,
		"21ee95f6-3776-4b1e-aabc-742eb5dc9170",
		"local test",
		"MK3avjS7g",
	}
	//TODO change token value to new type
	mockChannelService.EXPECT().FindChannelByToken("MK3avjS7g").Return(
		dummyChannel,
		nil,
	)

	mockContactService := mocks.NewMockContactService(ctrl)
	dummyContact := models.Contact{
		primitive.NilObjectID,
		"5582988887777",
		"Dummy",
		primitive.NilObjectID,
	}
	mockContactService.EXPECT().FindContact(dummyContact).Return(nil, errors.New("contact not found"))
	mockContactService.EXPECT().CreateContact(dummyContact).Return(dummyContact, nil)

	wh := WhatsappHandler{mockContactService, mockChannelService}

	router := chi.NewRouter()
	router.Post("/wr/receive/", wh.HandleIncomingRequests)

	request, _ := http.NewRequest(
		http.MethodPost,
		"/wr/receive/",
		strings.NewReader(
			`{"contacts":[{"profile":{"name":"user_name"},"wa_id":"12341341234"}],"messages":[{"from":"558299990000","id":"123456","text":{"body":"hi dude."},"timestamp":"623123123123","type":"text"}]}`))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
}
