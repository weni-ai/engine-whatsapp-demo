package flows

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/weni/whatsapp-router/pkg/jwt"
)

// Client calls the Flows API (e.g. project language).
type Client struct {
	BaseURL    string
	JWTSigner  *jwt.Signer
	HTTPClient *http.Client
}

// NewClient creates a new Flows API client. jwtSigner may be nil; then requests are sent without Authorization.
func NewClient(baseURL string, jwtSigner *jwt.Signer) *Client {
	return &Client{
		BaseURL:    baseURL,
		JWTSigner:  jwtSigner,
		HTTPClient: http.DefaultClient,
	}
}

// ProjectLanguageResponse is the response from GET /api/v2/projects/project_language.
type ProjectLanguageResponse struct {
	Language string `json:"language"`
}

// GetProjectLanguage returns the project language for the given channel UUID.
// GET /api/v2/projects/project_language?channel_uuid=<channel_uuid>
func (c *Client) GetProjectLanguage(channelUUID string) (string, error) {
	rawURL := fmt.Sprintf("%s/api/v2/projects/project_language", c.BaseURL)
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid flows base URL: %w", err)
	}
	q := u.Query()
	q.Set("channel_uuid", channelUUID)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	if err := c.addJWTAuthHeader(req, channelUUID); err != nil {
		return "", fmt.Errorf("add JWT auth: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("project_language returned status %d: %s", resp.StatusCode, string(body))
	}

	var out ProjectLanguageResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return out.Language, nil
}

func (c *Client) addJWTAuthHeader(req *http.Request, channelUUID string) error {
	if c.JWTSigner == nil {
		return nil
	}

	token, err := c.JWTSigner.GenerateToken(map[string]interface{}{
		"channel_uuid": channelUUID,
	})
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}
