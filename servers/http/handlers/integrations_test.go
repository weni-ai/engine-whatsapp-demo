package handlers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/servers/grpc/pb"
)

func TestHandleCreateChannel(t *testing.T) {
	dummyPayload := `{"uuid":"425b41f0-c554-4943-989c-5f88561a0cf5","name":"test-channel"}`

	ih := IntegrationsHandler{mockChannelService{}}
	router := chi.NewRouter()
	router.Post("/v1/channels", ih.HandleCreateChannel)
	request, err := http.NewRequest(
		http.MethodPost,
		"/v1/channels",
		bytes.NewReader([]byte(dummyPayload)),
	)
	assert.NoError(t, err)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, 201, response.Code)

	bodyString, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	assert.Equal(t, true, strings.Contains(string(bodyString), "weni-demo"))
}

var DummyCh = &models.Channel{
	UUID: "425b41f0-c554-4943-989c-5f88561a0cf5",
	Name: "test-channel",
}

type mockChannelService struct {
}

func (cs mockChannelService) CreateChannelDefault(*models.Channel) (*models.Channel, error) {
	return DummyCh, nil
}

func (cs mockChannelService) CreateChannel(ctx context.Context, channel *pb.ChannelRequest) (*pb.ChannelResponse, error) {
	return nil, nil
}

func (cs mockChannelService) FindChannel(channel *models.Channel) (*models.Channel, error) {
	return nil, nil
}

func (cs mockChannelService) FindChannelById(id string) (*models.Channel, error) {
	return nil, nil
}

func (cs mockChannelService) FindChannelByToken(token string) (*models.Channel, error) {
	return nil, nil
}
