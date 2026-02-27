package flows

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weni/whatsapp-router/pkg/jwt"
)

func generateTestRSAPEM(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	require.NoError(t, err)
	bytes := x509.MarshalPKCS1PrivateKey(key)
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: bytes}
	return string(pem.EncodeToMemory(block))
}

func TestGetProjectLanguage_WithoutJWT(t *testing.T) {
	var gotURL string
	var authHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotURL = r.URL.String()
		authHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ProjectLanguageResponse{Language: "pt-br"})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	lang, err := client.GetProjectLanguage("channel-uuid-123")
	require.NoError(t, err)
	assert.Equal(t, "pt-br", lang)
	assert.Contains(t, gotURL, "channel_uuid=channel-uuid-123")
	assert.Empty(t, authHeader, "Authorization header should be empty when no JWT signer")
}

func TestGetProjectLanguage_WithJWT(t *testing.T) {
	var authHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ProjectLanguageResponse{Language: "es"})
	}))
	defer server.Close()

	pemStr := generateTestRSAPEM(t)
	signer, err := jwt.NewSigner(pemStr, 60)
	require.NoError(t, err)
	client := NewClient(server.URL, signer)

	lang, err := client.GetProjectLanguage("channel-456")
	require.NoError(t, err)
	assert.Equal(t, "es", lang)
	assert.NotEmpty(t, authHeader, "Authorization header should be set when JWT signer is present")
	assert.Contains(t, authHeader, "Bearer ", "Authorization should be Bearer token")
}

func TestGetProjectLanguage_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not found"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	_, err := client.GetProjectLanguage("channel-789")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

func TestGetProjectLanguage_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`not json`))
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	_, err := client.GetProjectLanguage("channel-789")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode")
}

func TestGetProjectLanguage_URLHasChannelUUID(t *testing.T) {
	var gotURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotURL = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ProjectLanguageResponse{Language: "en-us"})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	_, err := client.GetProjectLanguage("my-channel-uuid")
	require.NoError(t, err)
	assert.Equal(t, "channel_uuid=my-channel-uuid", gotURL)
}
