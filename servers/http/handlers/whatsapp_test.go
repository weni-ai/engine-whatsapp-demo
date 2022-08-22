package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/weni/whatsapp-router/metric"
	mocks "github.com/weni/whatsapp-router/mocks/services"
	"github.com/weni/whatsapp-router/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var testCases = []struct {
	Label  string
	Data   string
	Status int
}{
	{Label: "Receive Valid Text Message", Data: helloMsg, Status: 200},
	{Label: "Receive Valid Audio Message", Data: audioMsg, Status: 200},
	{Label: "Receive Valid Button Message", Data: buttonMsg, Status: 200},
	{Label: "Receive Valid Document Message", Data: documentMsg, Status: 200},
	{Label: "Receive Valid Image Message", Data: imageMsg, Status: 200},
	{Label: "Receive Valid Location Message", Data: locationMsg, Status: 200},
	{Label: "Receive Valid Video Message", Data: videoMsg, Status: 200},
	{Label: "Receive Valid Voice Message", Data: voiceMsg, Status: 200},
	{Label: "Receive Valid Contact Message", Data: contactMsg, Status: 200},
}

var channelID = primitive.NewObjectID()

var dummyChannel = &models.Channel{
	ID:    channelID,
	UUID:  "21ee95f6-3776-4b1e-aabc-742eb5dc9170",
	Name:  "local test",
	Token: "weni-demo-44a2m17t0x",
}
var dummyChannel2 = &models.Channel{
	ID:    primitive.NewObjectID(),
	UUID:  "21ee95f6-3776-4b1e-aabc-742eb5dc9170",
	Name:  "local test",
	Token: "weni-demo-1234567890",
}

var dummyContact = &models.Contact{
	URN:     "5582988887777",
	Name:    "Dummy",
	Channel: dummyChannel.ID,
}

var incomingDummyContact = &models.Contact{
	URN:  "5582988887777",
	Name: "Dummy",
}

func TestContactSendFlowsChoice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	payloadFlows := fmt.Sprintf(
		`{
			"to":"%s",
			"type":"interactive",
			"interative":{
				"type":"button",
				"body":"Escolha abaixo qual fluxo deseja iniciar.",
				"action": {
					"buttons": [
						{
							"type": "reply",
							"reply": {
							"id": "%s",
							"title": "%s" 
							}
						},
						{
							"type": "reply",
							"reply": {
							"id": "%s",
							"title": "%s" 
							}
						},
						{
							"type": "reply",
							"reply": {
							"id": "%s",
							"title": "%s" 
							}
						}
					] 
				}
			}
		}`,
		dummyContact.URN,
		fl.FlowsStarts[0].Name,
		fl.FlowsStarts[0].Name,
		fl.FlowsStarts[1].Name,
		fl.FlowsStarts[1].Name,
		fl.FlowsStarts[2].Name,
		fl.FlowsStarts[2].Name,
	)
	incomingRequest := `{"contacts":[{"profile":{"name":"Dummy"},"wa_id":"12341341234"}],"messages":[{"from":"5582988887777","id":"123456","text":{"body":"weni-demo-44a2m17t0x"},"timestamp":"623123123123","type":"text"}]}`

	metricService, err := metric.NewPrometheusService()
	assert.NoError(t, err)

	mockChannelService := mocks.NewMockChannelService(ctrl)
	mockContactService := mocks.NewMockContactService(ctrl)
	mockCourierService := mocks.NewMockCourierService(ctrl)
	mockWhatsappService := mocks.NewMockWhatsappService(ctrl)
	mockConfigService := mocks.NewMockConfigService(ctrl)
	mockFlowsService := mocks.NewMockFlowsService(ctrl)
	mockChannelService.EXPECT().FindChannelByToken(dummyChannel.Token).Return(dummyChannel, nil)
	mockContactService.EXPECT().FindContact(incomingDummyContact).Return(nil, errors.New("contact not found"))
	mockContactService.EXPECT().CreateContact(dummyContact).Return(dummyContact, nil)
	mockFlowsService.EXPECT().FindFlows(flows).Return(fl, nil)
	mockWhatsappService.EXPECT().SendMessage([]byte(payloadFlows)).Return(
		http.Header{"content-type": {"application/json"}},
		io.NopCloser(bytes.NewReader([]byte(`{"messages":{"id":"hBEGVYKZRIIyAgmiTgezkroUL2Q"}],"meta":{"api_status":"stable","version":"2.35.2"}}`))),
		nil,
	)

	wh := WhatsappHandler{mockContactService, mockChannelService, mockCourierService, mockWhatsappService, mockConfigService, metricService, mockFlowsService}
	router := chi.NewRouter()
	router.Post("/wr/receive/", wh.HandleIncomingRequests)
	request, _ := http.NewRequest(
		http.MethodPost,
		"/wr/receive/",
		strings.NewReader(incomingRequest))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)
}

func TestHandleIncomingRequest(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Label, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockChannelService := mocks.NewMockChannelService(ctrl)
			mockContactService := mocks.NewMockContactService(ctrl)
			mockCourierService := mocks.NewMockCourierService(ctrl)
			mockWhatsappService := mocks.NewMockWhatsappService(ctrl)
			mockConfigService := mocks.NewMockConfigService(ctrl)
			metricService, err := metric.NewPrometheusService()
			mockFlowsService := mocks.NewMockFlowsService(ctrl)
			assert.NoError(t, err)

			mockContactService.EXPECT().FindContact(incomingDummyContact).Return(dummyContact, nil)
			mockChannelService.EXPECT().FindChannelById(channelID.Hex()).Return(dummyChannel, nil)
			mockFlowsService.EXPECT().FindFlows(flows).Return(fl, nil)
			mockCourierService.EXPECT().RedirectMessage(dummyChannel.UUID, tc.Data).Return(tc.Status, nil)

			wh := WhatsappHandler{mockContactService, mockChannelService, mockCourierService, mockWhatsappService, mockConfigService, metricService, mockFlowsService}
			router := chi.NewRouter()
			router.Post("/wr/receive/", wh.HandleIncomingRequests)
			request, _ := http.NewRequest(
				http.MethodPost,
				"/wr/receive/",
				strings.NewReader(tc.Data),
			)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(t, response.Code, tc.Status)

			ctrl.Finish()
		})
	}
}

func TestContactTokenUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannelService := mocks.NewMockChannelService(ctrl)
	mockContactService := mocks.NewMockContactService(ctrl)
	mockCourierService := mocks.NewMockCourierService(ctrl)
	mockWhatsappService := mocks.NewMockWhatsappService(ctrl)
	mockConfigService := mocks.NewMockConfigService(ctrl)
	mockFlowsService := mocks.NewMockFlowsService(ctrl)
	metricService, err := metric.NewPrometheusService()
	assert.NoError(t, err)

	dummyUpdatedContact := &models.Contact{
		URN:     "5582988887777",
		Name:    "Dummy",
		Channel: dummyChannel2.ID,
	}

	payloadFlows := fmt.Sprintf(
		`{
			"to":"%s",
			"type":"interactive",
			"interative":{
				"type":"button",
				"body":"Escolha abaixo qual fluxo deseja iniciar.",
				"action": {
					"buttons": [
						{
							"type": "reply",
							"reply": {
							"id": "%s",
							"title": "%s" 
							}
						},
						{
							"type": "reply",
							"reply": {
							"id": "%s",
							"title": "%s" 
							}
						},
						{
							"type": "reply",
							"reply": {
							"id": "%s",
							"title": "%s" 
							}
						}
					] 
				}
			}
		}`,
		dummyUpdatedContact.URN,
		fl.FlowsStarts[0].Name,
		fl.FlowsStarts[0].Name,
		fl.FlowsStarts[1].Name,
		fl.FlowsStarts[1].Name,
		fl.FlowsStarts[2].Name,
		fl.FlowsStarts[2].Name,
	)

	incomingRequest := `{"contacts":[{"profile":{"name":"Dummy"},"wa_id":"12341341234"}],"messages":[{"from":"5582988887777","id":"123456","text":{"body":"weni-demo-1234567890"},"timestamp":"623123123123","type":"text"}]}`
	mockContactService.EXPECT().FindContact(incomingDummyContact).Return(dummyContact, nil)
	mockChannelService.EXPECT().FindChannelById(dummyContact.Channel.Hex()).Return(dummyChannel, nil)
	mockContactService.EXPECT().UpdateContact(dummyContact).Return(dummyUpdatedContact, nil)
	mockChannelService.EXPECT().FindChannelByToken(dummyChannel2.Token).Return(dummyChannel2, nil)
	mockFlowsService.EXPECT().FindFlows(flowsUpdate).Return(flUpdate, nil)
	mockWhatsappService.EXPECT().SendMessage([]byte(payloadFlows)).Return(
		http.Header{"content-type": {"application/json"}},
		io.NopCloser(bytes.NewReader([]byte(`{"messages":{"id":"hBEGVYKZRIIyAgmiTgezkroUL2Q"}],"meta":{"api_status":"stable","version":"2.35.2"}}`))),
		nil,
	)

	wh := WhatsappHandler{mockContactService, mockChannelService, mockCourierService, mockWhatsappService, mockConfigService, metricService, mockFlowsService}
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

func TestRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannelService := mocks.NewMockChannelService(ctrl)
	mockContactService := mocks.NewMockContactService(ctrl)
	mockCourierService := mocks.NewMockCourierService(ctrl)
	mockWhatsappService := mocks.NewMockWhatsappService(ctrl)
	mockConfigService := mocks.NewMockConfigService(ctrl)
	mockFlowsService := mocks.NewMockFlowsService(ctrl)
	metricService, err := metric.NewPrometheusService()
	assert.NoError(t, err)

	loginBody := `{"users":[{"token":"eyJhbGciOiAiSFMyNTYiLCAidHlwIjogIkpXVCJ9.eyJ1c2VyIjoiQWRtaW4iLCJpYXQiOjE2NDAxODIzMjIsImV4cCI6MTY0MDc4NzEyMiwid2E6cmFuZCI6ImVkMWU5OGU4ZjA4NmIxMDQzNDBlM2MxMGFjNGU3YzY3In0.2pEh32jyfBLUjxWNklEtgOrZqy7TgGj48y5pVTgl7FU","expires_after":"2021-12-29 14:12:02+00:00"}],"meta":{"version":"v2.37.1","api_status":"stable"}}`
	mockWhatsappService.EXPECT().Login().Return(
		&http.Response{
			Header:     http.Header{"content-type": {"application/json"}},
			Body:       io.NopCloser(bytes.NewReader([]byte(loginBody))),
			StatusCode: 200,
		},
		nil,
	)

	conf := &models.Config{
		Token: "eyJhbGciOiAiSFMyNTYiLCAidHlwIjogIkpXVCJ9.eyJ1c2VyIjoiQWRtaW4iLCJpYXQiOjE2NDAxODIzMjIsImV4cCI6MTY0MDc4NzEyMiwid2E6cmFuZCI6ImVkMWU5OGU4ZjA4NmIxMDQzNDBlM2MxMGFjNGU3YzY3In0.2pEh32jyfBLUjxWNklEtgOrZqy7TgGj48y5pVTgl7FU",
	}
	mockConfigService.EXPECT().CreateOrUpdate(
		conf,
	).Return(conf, nil)

	wh := WhatsappHandler{mockContactService, mockChannelService, mockCourierService, mockWhatsappService, mockConfigService, metricService, mockFlowsService}
	router := chi.NewRouter()
	testRoute := "/v1/users/login"
	router.Post(testRoute, wh.RefreshToken)
	request, _ := http.NewRequest(
		http.MethodPost,
		testRoute,
		nil,
	)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)
}

func TestHandleHealth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	healthBody := `{
		"health": {
			"192.168.11.206:wamaster-foobar-6bdc9bcc78-42vlm": {
				"gateway_status": "connected",
				"role": "primary_master"
			},
			"192.168.50.212:wacore-foobar-6d9d96959c-zzddg": {
				"gateway_status": "connected",
				"role": "coreapp"
			},
			"192.168.65.190:wacore-foobar-6d9d96959c-ddg5r": {
				"gateway_status": "connected",
				"role": "coreapp"
			}
		},
		"meta": {
			"version": "v2.37.1",
			"api_status": "stable"
		}
	}`

	mockWhatsappService := mocks.NewMockWhatsappService(ctrl)
	mockWhatsappService.EXPECT().Health().Return(
		&http.Response{
			Header:     http.Header{},
			Body:       io.NopCloser(bytes.NewReader([]byte(healthBody))),
			StatusCode: 200,
		},
		nil,
	)

	wh := WhatsappHandler{
		WhatsappService: mockWhatsappService,
	}
	router := chi.NewRouter()
	router.Get("/v1/health", wh.HandleHealth)
	request, _ := http.NewRequest(
		http.MethodGet,
		"/v1/health",
		nil,
	)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
}

func TestHandleGetMedia(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mediaID := "123456-qwerty-asdfgh-zxcvb"
	mockWhatsappService := mocks.NewMockWhatsappService(ctrl)
	mockWhatsappService.EXPECT().GetMedia(http.Header{}, mediaID).Return(
		&http.Response{
			Header:     http.Header{},
			Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			StatusCode: 200,
		},
		nil,
	)

	wh := WhatsappHandler{WhatsappService: mockWhatsappService}
	router := chi.NewRouter()
	router.Get("/v1/media/{mediaID}", wh.HandleGetMedia)
	request, err := http.NewRequest(
		http.MethodGet,
		"/v1/media/"+mediaID,
		nil,
	)
	assert.NoError(t, err)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
}

func TestHandlePostMedia(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWhatsappService := mocks.NewMockWhatsappService(ctrl)
	mockWhatsappService.EXPECT().PostMedia(
		http.Header{}, http.NoBody,
	).Return(
		&http.Response{
			Header:     http.Header{},
			Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			StatusCode: 201,
		},
		nil,
	)

	wh := WhatsappHandler{WhatsappService: mockWhatsappService}
	router := chi.NewRouter()
	router.Post("/v1/media", wh.HandlePostMedia)
	request, err := http.NewRequest(
		http.MethodPost,
		"/v1/media",
		bytes.NewReader([]byte("")),
	)
	assert.NoError(t, err)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, 201, response.Code)
}

var flows = &models.Flows{
	Channel: dummyContact.Channel,
}

var fl = &models.Flows{
	Channel: dummyChannel.ID,
	FlowsStarts: []models.Flow{
		{
			Name:    "flow1",
			UUID:    "b23662c3-935b-4748-b715-b62a44e9d228",
			Keyword: "flow1",
		},
		{
			Name:    "flow2",
			UUID:    "b23662c3-935b-4748-b725-b62a44e9d228",
			Keyword: "flow2",
		},
		{
			Name:    "flow3",
			UUID:    "b23662c3-935b-4748-b735-b62a44e9d228",
			Keyword: "flow3",
		},
	},
}

var flowsUpdate = &models.Flows{
	Channel: dummyChannel2.ID,
}

var flUpdate = &models.Flows{
	Channel: dummyChannel2.ID,
	FlowsStarts: []models.Flow{
		{
			Name:    "flow1",
			UUID:    "b23662c3-935b-4748-b715-b62a44e9d228",
			Keyword: "flow1",
		},
		{
			Name:    "flow2",
			UUID:    "b23662c3-935b-4748-b725-b62a44e9d228",
			Keyword: "flow2",
		},
		{
			Name:    "flow3",
			UUID:    "b23662c3-935b-4748-b735-b62a44e9d228",
			Keyword: "flow3",
		},
	},
}

var helloMsg = `{
	"contacts":[{
		"profile": {
			"name": "Dummy"
		},
		"wa_id": "5582988887777"
	}],
  "messages": [{
    "from": "5582988887777",
    "id": "41",
    "timestamp": "1454119029",
    "text": {
      "body": "hello world"
    },
    "type": "text"
   }]
}`

var audioMsg = `{
	"contacts":[{
		"profile": {
			"name": "Dummy"
		},
		"wa_id": "5582988887777"
	}],
	"messages": [{
		"from": "5582988887777",
		"id": "41",
		"timestamp": "1454119029",
		"type": "audio",
		"audio": {
			"file": "/path/to/v1/media/41",
			"id": "41",
			"link": "https://example.org/v1/media/41",
			"mime_type": "text/plain",
			"sha256": "the-sha-signature"
		}
	}]
}`

var buttonMsg = `{
	"contacts":[{
		"profile": {
			"name": "Dummy"
		},
		"wa_id": "5582988887777"
	}],
	"messages": [{
		"from": "5582988887777",
		"id": "41",
		"timestamp": "1454119029",
		"type": "button",
		"button": {
			"payload": null,
			"text": "BUTTON1"
		}
	}]
}`

var documentMsg = `{
	"contacts":[{
		"profile": {
			"name": "Dummy"
		},
		"wa_id": "5582988887777"
	}],
	"messages": [{
		"from": "5582988887777",
		"id": "41",
		"timestamp": "1454119029",
		"type": "document",
		"document": {
			"file": "/path/to/v1/media/41",
			"id": "41",
			"link": "https://example.org/v1/media/41",
			"mime_type": "text/plain",
			"sha256": "the-sha-signature",
			"caption": "the caption",
			"filename": "filename.type"
		}
	}]
}`

var imageMsg = `{
	"contacts":[{
		"profile": {
			"name": "Dummy"
		},
		"wa_id": "5582988887777"
	}],
	"messages": [{
		"from": "5582988887777",
		"id": "41",
		"timestamp": "1454119029",
		"type": "image",
		"image": {
			"file": "/path/to/v1/media/41",
			"id": "41",
			"link": "https://example.org/v1/media/41",
			"mime_type": "text/plain",
			"sha256": "the-sha-signature",
			"caption": "the caption"
		}
	}]
}`

var locationMsg = `{
	"contacts":[{
		"profile": {
			"name": "Dummy"
		},
		"wa_id": "5582988887777"
	}],
	"messages": [{
		"from": "5582988887777",
		"id": "41",
		"timestamp": "1454119029",
		"type": "location",
		"location": {
			"address": "some address",
			"latitude": 0.00,
			"longitude": 1.00,
			"name": "some name",
			"url": "https://example.org/"
		}
	}]
}`

var videoMsg = `{
	"contacts":[{
		"profile": {
			"name": "Dummy"
		},
		"wa_id": "5582988887777"
	}],
	"messages": [{
		"from": "5582988887777",
		"id": "41",
		"timestamp": "1454119029",
		"type": "video",
		"video": {
			"file": "/path/to/v1/media/41",
			"id": "41",
			"link": "https://example.org/v1/media/41",
			"mime_type": "text/plain",
			"sha256": "the-sha-signature"
		}
	}]
}`

var voiceMsg = `{
	"contacts":[{
		"profile": {
			"name": "Dummy"
		},
		"wa_id": "5582988887777"
	}],
	"messages": [{
		"from": "5582988887777",
		"id": "41",
		"timestamp": "1454119029",
		"type": "voice",
		"voice": {
			"file": "/path/to/v1/media/41",
			"id": "41",
			"link": "https://example.org/v1/media/41",
			"mime_type": "text/plain",
			"sha256": "the-sha-signature"
		}
	}]
}`

var contactMsg = `{
	"contacts":[{
		"profile": {
			"name": "Dummy"
		},
		"wa_id": "5582988887777"
	}],
	"messages": [{
		"from": "5582988887777",
		"id": "41",
		"timestamp": "1454119029",
		"type": "contacts",
		"contacts": [{
			"addresses": [],
			"emails": [],
			"ims": [],
			"name": {
				"first_name": "John Cruz",
				"formatted_name": "John Cruz"
			},
			"org": {},
			"phones": [
				{
					"phone": "+1 415-858-6273",
					"type": "CELL",
					"wa_id": "14158586273"
				},
				{
					"phone": "+1 415-858-6274",
					"type": "CELL",
					"wa_id": "14158586274"
				}
			],
			"urls": []
		}]
	}]
}`
