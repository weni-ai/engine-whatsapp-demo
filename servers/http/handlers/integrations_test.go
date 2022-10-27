package handlers

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Nerzal/gocloak/v11"
	"github.com/go-chi/chi"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weni/whatsapp-router/config"
	mocks "github.com/weni/whatsapp-router/mocks/services"
	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/servers/grpc/pb"
)

func TestHandleCreateChannel(t *testing.T) {
	dummyPayload := `{"uuid":"425b41f0-c554-4943-989c-5f88561a0cf5","name":"test-channel"}`

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFlowsService := mocks.NewMockFlowsService(ctrl)

	ih := IntegrationsHandler{mockChannelService{}, mockFlowsService}
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

func TestKeycloakAuth(t *testing.T) {
	cfg := GetConfig(t)
	kkClient = NewClientWithDebug(t)
	assert.NotNil(t, kkClient)

	SetUpTestUser(t, kkClient)

	token := GetUserToken(t, kkClient)

	config.GetConfig().OIDC.Realm = cfg.GoCloak.Realm

	log.Println("GGwp")

	router := chi.NewRouter()
	router.Post("/", KeycloackAuth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	request, err := http.NewRequest(
		http.MethodPost,
		"/",
		nil,
	)

	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	assert.NoError(t, err)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
}

func NewClientWithDebug(t testing.TB) gocloak.GoCloak {
	cfg := GetConfig(t)
	client := gocloak.NewClient(cfg.HostName)
	cond := func(resp *resty.Response, err error) bool {
		if resp != nil && resp.IsError() {
			if e, ok := resp.Error().(*gocloak.HTTPErrorResponse); ok {
				msg := e.String()
				return strings.Contains(msg, "Cached clientScope not found") || strings.Contains(msg, "unknown_error")
			}
		}
		return false
	}

	restyClient := client.RestyClient()

	restyClient.
		SetRetryCount(10).
		SetRetryWaitTime(2 * time.Second).
		AddRetryCondition(cond)

	return client
}

var (
	cfg        *Config
	configOnce sync.Once
	setupOnce  sync.Once
	testUserID string
)

func GetConfig(t testing.TB) *Config {
	configOnce.Do(func() {
		rand.Seed(time.Now().UTC().UnixNano())
		configFileName, ok := os.LookupEnv("GOCLOAK_TEST_CONFIG")
		if !ok {
			configFileName = filepath.Join("../../../testdata", "config.json")
		}
		configFile, err := os.Open(configFileName)
		require.NoError(t, err, "cannot open config.json")
		defer func() {
			err := configFile.Close()
			require.NoError(t, err, "cannot close config file")
		}()
		data, err := ioutil.ReadAll(configFile)
		require.NoError(t, err, "cannot read config.json")
		cfg = &Config{}
		err = json.Unmarshal(data, cfg)
		require.NoError(t, err, "cannot parse cfg.json")
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		if len(cfg.Proxy) != 0 {
			proxy, err := url.Parse(cfg.Proxy)
			require.NoError(t, err, "incorrect proxy url: "+cfg.Proxy)
			http.DefaultTransport.(*http.Transport).Proxy = http.ProxyURL(proxy)
		}
		if cfg.GoCloak.UserName == "" {
			cfg.GoCloak.UserName = "test_user"
		}
	})
	return cfg
}

func SetUpTestUser(t testing.TB, client gocloak.GoCloak) {
	setupOnce.Do(func() {
		cfg := GetConfig(t)
		token := GetAdminToken(t, client)

		user := gocloak.User{
			Username:      gocloak.StringP(cfg.GoCloak.UserName),
			Email:         gocloak.StringP(cfg.GoCloak.UserName + "@localhost.com"),
			EmailVerified: gocloak.BoolP(true),
			Enabled:       gocloak.BoolP(true),
		}

		createdUserID, err := client.CreateUser(
			context.Background(),
			token.AccessToken,
			cfg.GoCloak.Realm,
			user,
		)

		apiError, ok := err.(*gocloak.APIError)
		if ok && apiError.Code == http.StatusConflict {
			users, err := client.GetUsers(
				context.Background(),
				token.AccessToken,
				cfg.GoCloak.Realm,
				gocloak.GetUsersParams{
					Username: gocloak.StringP(cfg.GoCloak.UserName),
				})
			require.NoError(t, err, "GetUsers failed")
			for _, user := range users {
				if gocloak.PString(user.Username) == cfg.GoCloak.UserName {
					testUserID = gocloak.PString(user.ID)
					break
				}
			}
		} else {
			require.NoError(t, err, "CreateUser failed")
			testUserID = createdUserID
		}

		err = client.SetPassword(
			context.Background(),
			token.AccessToken,
			testUserID,
			cfg.GoCloak.Realm,
			cfg.GoCloak.Password,
			false)
		require.NoError(t, err, "SetPassword failed")
	})
}

func GetAdminToken(t testing.TB, client gocloak.GoCloak) *gocloak.JWT {
	cfg := GetConfig(t)
	token, err := client.LoginAdmin(
		context.Background(),
		cfg.Admin.UserName,
		cfg.Admin.Password,
		cfg.Admin.Realm)
	require.NoError(t, err, "Login Admin failed")
	return token
}

func GetUserToken(t *testing.T, client gocloak.GoCloak) *gocloak.JWT {
	SetUpTestUser(t, client)
	cfg := GetConfig(t)
	token, err := client.Login(
		context.Background(),
		cfg.GoCloak.ClientID,
		cfg.GoCloak.ClientSecret,
		cfg.GoCloak.Realm,
		cfg.GoCloak.UserName,
		cfg.GoCloak.Password)
	require.NoError(t, err, "Login failed")
	return token
}

type Config struct {
	HostName string        `json:"hostname"`
	Proxy    string        `json:"proxy,omitempty"`
	Admin    configAdmin   `json:"admin"`
	GoCloak  configGoCloak `json:"gocloak"`
}

type configAdmin struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Realm    string `json:"realm"`
}

type configGoCloak struct {
	UserName     string `json:"username"`
	Password     string `json:"password"`
	Realm        string `json:"realm"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

var DummyFl = &models.Flows{
	FlowsStarts: []models.Flow{
		{
			Name:    "test_flow1",
			UUID:    "507b6703-cc80-41fc-8a1b-cca573518dbb",
			Keyword: "hello1",
		},
		{
			Name:    "test_flow2",
			UUID:    "a76b3106-5e3d-462d-a0fc-4817c0d73ce7",
			Keyword: "hello2",
		},
		{
			Name:    "test_flow3",
			UUID:    "d7c97de5-bd06-4d7f-904f-63a7f8dd6b9d",
			Keyword: "hello3",
		},
	},
	Channel: "1d514cf1-a829-415d-8955-748563e173cf",
}

func TestHandleInitialProjectFlows(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dummyPayload := `{"channel_uuid":` + `"` + DummyFl.Channel + `"` + `,"flows_starts":[{"flow_name":"test_flow1","flow_uuid":"507b6703-cc80-41fc-8a1b-cca573518dbb","keyword":"hello1"},{"flow_name":"test_flow2","flow_uuid":"a76b3106-5e3d-462d-a0fc-4817c0d73ce7","keyword":"hello2"},{"flow_name":"test_flow3","flow_uuid":"d7c97de5-bd06-4d7f-904f-63a7f8dd6b9d","keyword":"hello3"}]}`

	mockFlowsService := mocks.NewMockFlowsService(ctrl)
	ih := IntegrationsHandler{mockChannelService{}, mockFlowsService}
	router := chi.NewRouter()
	router.Post("/v1/flows", ih.HandleInitialProjectFlows)
	mockFlowsService.EXPECT().FindFlows(DummyFl).Return(nil, nil)
	request, err := http.NewRequest(
		http.MethodPost,
		"/v1/flows",
		bytes.NewReader([]byte(dummyPayload)),
	)
	mockFlowsService.EXPECT().CreateFlows(DummyFl).Return(DummyFl, nil)
	assert.NoError(t, err)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, 201, response.Code)
}
