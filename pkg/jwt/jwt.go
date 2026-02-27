package jwt

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Signer generates JWT tokens for authenticating with external services (e.g. Flows API).
type Signer struct {
	privateKey *rsa.PrivateKey
	expiration time.Duration
}

// NewSigner creates a new JWT signer from a PEM-encoded RSA private key.
// expirationMins is the token lifetime in minutes.
func NewSigner(privateKeyPEM string, expirationMins int64) (*Signer, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
	}

	return &Signer{
		privateKey: key,
		expiration: time.Duration(expirationMins) * time.Minute,
	}, nil
}

// GenerateToken creates a signed JWT with the given claims and standard exp/iat.
// Typically used with channel_uuid for Flows API authentication.
func (s *Signer) GenerateToken(claims map[string]interface{}) (string, error) {
	now := time.Now()
	mapClaims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(s.expiration).Unix(),
	}
	for k, v := range claims {
		mapClaims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, mapClaims)
	return token.SignedString(s.privateKey)
}
