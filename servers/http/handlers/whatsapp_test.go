package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "github.com/weni/whatsapp-router/mocks/services"
	"github.com/weni/whatsapp-router/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestContactTokenConfirmation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannelService := mocks.NewMockChannelService(ctrl)
	mockContactService := mocks.NewMockContactService(ctrl)
	mockCourierService := mocks.NewMockCourierService(ctrl)
	mockWhatsappService := mocks.NewMockWhatsappService(ctrl)

	channelID := primitive.NewObjectID()

	dummyChannel := &models.Channel{
		ID:    channelID,
		UUID:  "21ee95f6-3776-4b1e-aabc-742eb5dc9170",
		Name:  "local test",
		Token: "localtest-whatsapp-demo-44a2m17t0x",
	}
	mockChannelService.EXPECT().FindChannelByToken("localtest-whatsapp-demo-44a2m17t0x").Return(
		dummyChannel,
		nil,
	)
	dummyContact := &models.Contact{
		URN:     "5582988887777",
		Name:    "Dummy",
		Channel: primitive.NilObjectID,
	}

	newDummyContact := &models.Contact{
		URN:     "5582988887777",
		Name:    "Dummy",
		Channel: channelID,
	}

	urn := newDummyContact.URN
	payload := fmt.Sprintf(
		`{"to":"%s","type":"text","text":{"body":"%s"}}`,
		urn,
		confirmationMessage,
	)
	dummyPayloadBytes := []byte(payload)

	mockContactService.EXPECT().FindContact(dummyContact).Return(nil, errors.New("contact not found"))
	mockContactService.EXPECT().CreateContact(newDummyContact).Return(newDummyContact, nil)
	mockWhatsappService.EXPECT().SendMessage(dummyPayloadBytes).Return(
		http.Header{
			"content-type": {"application/json"},
		},
		ioutil.NopCloser(bytes.NewReader([]byte(`{"messages":{"id":"gBEGVYKZRIIyAgmiTgezkroUL2Q"}],"meta":{"api_status":"stable","version":"2.35.2"}}`))),
		nil,
	)

	wh := WhatsappHandler{mockContactService, mockChannelService, mockCourierService, mockWhatsappService}

	router := chi.NewRouter()
	router.Post("/wr/receive/", wh.HandleIncomingRequests)

	request, _ := http.NewRequest(
		http.MethodPost,
		"/wr/receive/",
		strings.NewReader(
			`{"contacts":[{"profile":{"name":"Dummy"},"wa_id":"12341341234"}],"messages":[{"from":"5582988887777","id":"123456","text":{"body":"localtest-whatsapp-demo-44a2m17t0x"},"timestamp":"623123123123","type":"text"}]}`))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, response.Code, 201)
}

func TestSendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannelService := mocks.NewMockChannelService(ctrl)
	mockContactService := mocks.NewMockContactService(ctrl)
	mockCourierService := mocks.NewMockCourierService(ctrl)
	mockWhatsappService := mocks.NewMockWhatsappService(ctrl)

	channelID := primitive.NewObjectID()

	dummyChannel := &models.Channel{
		ID:    channelID,
		UUID:  "21ee95f6-3776-4b1e-aabc-742eb5dc9170",
		Name:  "local test",
		Token: "localtest-whatsapp-demo-44a2m17t0x",
	}

	incomingDummyContact := &models.Contact{
		URN:  "5582988887777",
		Name: "Dummy",
	}

	dummyContact := &models.Contact{
		URN:     "5582988887777",
		Name:    "Dummy",
		Channel: dummyChannel.ID,
	}

	incomingRequest := `{"contacts":[{"profile":{"name":"Dummy"},"wa_id":"12341341234"}],"messages":[{"from":"5582988887777","id":"123456","text":{"body":"hello"},"timestamp":"623123123123","type":"text"}]}`
	mockContactService.EXPECT().FindContact(incomingDummyContact).Return(dummyContact, nil)
	mockChannelService.EXPECT().FindChannelById(channelID.Hex()).Return(dummyChannel, nil)
	mockChannelService.EXPECT().FindChannelByToken(extractTextMessage(incomingRequest)).Return(nil, nil)
	mockCourierService.EXPECT().RedirectMessage(dummyChannel.UUID, incomingRequest).Return(http.StatusOK, nil)

	wh := WhatsappHandler{mockContactService, mockChannelService, mockCourierService, mockWhatsappService}

	router := chi.NewRouter()
	router.Post("/wr/receive/", wh.HandleIncomingRequests)

	request, _ := http.NewRequest(
		http.MethodPost,
		"/wr/receive/",
		strings.NewReader(incomingRequest),
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, response.Code, 200)
}

func TestContactTokenUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannelService := mocks.NewMockChannelService(ctrl)
	mockContactService := mocks.NewMockContactService(ctrl)
	mockCourierService := mocks.NewMockCourierService(ctrl)
	mockWhatsappService := mocks.NewMockWhatsappService(ctrl)

	dummyChannel := &models.Channel{
		ID:    primitive.NewObjectID(),
		UUID:  "21ee95f6-3776-4b1e-aabc-742eb5dc9170",
		Name:  "local test",
		Token: "localtest-whatsapp-demo-44a2m17t0x",
	}

	dummyChannel2 := &models.Channel{
		ID:    primitive.NewObjectID(),
		UUID:  "21ee95f6-3776-4b1e-aabc-742eb5dc9170",
		Name:  "local test",
		Token: "localtest-whatsapp-demo-1234567890",
	}

	incomingDummyContact := &models.Contact{
		URN:  "5582988887777",
		Name: "Dummy",
	}

	dummyContact := &models.Contact{
		URN:     "5582988887777",
		Name:    "Dummy",
		Channel: dummyChannel.ID,
	}

	dummyUpdatedContact := &models.Contact{
		URN:     "5582988887777",
		Name:    "Dummy",
		Channel: dummyChannel2.ID,
	}

	urn := dummyContact.URN
	payload := fmt.Sprintf(
		`{"to":"%s","type":"text","text":{"body":"%s"}}`,
		urn,
		confirmationMessage,
	)
	dummyPayloadBytes := []byte(payload)

	incomingRequest := `{"contacts":[{"profile":{"name":"Dummy"},"wa_id":"12341341234"}],"messages":[{"from":"5582988887777","id":"123456","text":{"body":"localtest-whatsapp-demo-1234567890"},"timestamp":"623123123123","type":"text"}]}`
	mockContactService.EXPECT().FindContact(incomingDummyContact).Return(dummyContact, nil)
	mockContactService.EXPECT().UpdateContact(dummyContact).Return(dummyUpdatedContact, nil)
	mockChannelService.EXPECT().FindChannelByToken(extractTextMessage(incomingRequest)).Return(dummyChannel2, nil)
	mockWhatsappService.EXPECT().SendMessage(dummyPayloadBytes).Return(
		http.Header{
			"content-type": {"application/json"},
		},
		ioutil.NopCloser(bytes.NewReader([]byte(`{"messages":{"id":"gBEGVYKZRIIyAgmiTgezkroUL2Q"}],"meta":{"api_status":"stable","version":"2.35.2"}}`))),
		nil,
	)

	wh := WhatsappHandler{mockContactService, mockChannelService, mockCourierService, mockWhatsappService}

	router := chi.NewRouter()
	router.Post("/wr/receive/", wh.HandleIncomingRequests)

	request, _ := http.NewRequest(
		http.MethodPost,
		"/wr/receive/",
		strings.NewReader(incomingRequest),
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, response.Code, 200)
}
