package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// generateTestRSAPEM returns a PEM-encoded RSA private key for tests.
func generateTestRSAPEM(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	require.NoError(t, err)
	bytes := x509.MarshalPKCS1PrivateKey(key)
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: bytes}
	return string(pem.EncodeToMemory(block))
}

func TestNewSigner_InvalidPEM(t *testing.T) {
	_, err := NewSigner("not a valid pem", 60)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse RSA private key")
}

func TestNewSigner_ValidPEM(t *testing.T) {
	pemStr := generateTestRSAPEM(t)
	signer, err := NewSigner(pemStr, 60)
	require.NoError(t, err)
	require.NotNil(t, signer)
}

func TestSigner_GenerateToken_WithClaims(t *testing.T) {
	pemStr := generateTestRSAPEM(t)
	signer, err := NewSigner(pemStr, 60)
	require.NoError(t, err)

	claims := map[string]interface{}{
		"channel_uuid": "test-channel-123",
	}
	tokenStr, err := signer.GenerateToken(claims)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	// Parse and verify token (use same key's public part)
	key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(pemStr))
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return &key.PublicKey, nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claimsOut := token.Claims.(jwt.MapClaims)
	assert.Equal(t, "test-channel-123", claimsOut["channel_uuid"])
	assert.Contains(t, claimsOut, "iat")
	assert.Contains(t, claimsOut, "exp")

	exp := time.Unix(int64(claimsOut["exp"].(float64)), 0)
	iat := time.Unix(int64(claimsOut["iat"].(float64)), 0)
	assert.True(t, exp.After(iat), "exp should be after iat")
}

func TestSigner_GenerateToken_Expiration(t *testing.T) {
	pemStr := generateTestRSAPEM(t)
	// 30 minutes expiration
	signer, err := NewSigner(pemStr, 30)
	require.NoError(t, err)

	tokenStr, err := signer.GenerateToken(map[string]interface{}{"channel_uuid": "ch"})
	require.NoError(t, err)

	key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(pemStr))
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return &key.PublicKey, nil
	})
	require.NoError(t, err)
	claims := token.Claims.(jwt.MapClaims)
	exp := time.Unix(int64(claims["exp"].(float64)), 0)
	iat := time.Unix(int64(claims["iat"].(float64)), 0)
	diff := exp.Sub(iat)
	assert.Equal(t, 30*time.Minute, diff.Round(time.Minute))
}
