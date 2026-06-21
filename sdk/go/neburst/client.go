package neburst

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Client is the Neburst OpenAPI client.
type Client struct {
	baseURL    string
	keyID      string
	secret     string
	httpClient *http.Client
}

// NewClient creates a new Neburst API client.
//
// Supports two calling styles:
//
//	NewClient("https://api.neburst.com", "nb_key_...", "nb_secret_...")   // key ID + secret
//	NewClient("https://api.neburst.com", "eyJrZXlfaWQ...base64...")      // combined base64 key
//
// The combined key is a base64-encoded JSON: {"key_id":"nb_key_...","secret":"nb_secret_..."}
func NewClient(baseURL string, auth ...string) *Client {
	var keyID, secret string
	switch len(auth) {
	case 1:
		keyID, secret = parseCombinedKey(auth[0])
	case 2:
		keyID, secret = auth[0], auth[1]
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		keyID:   keyID,
		secret:  secret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func parseCombinedKey(combined string) (string, string) {
	data, err := base64.StdEncoding.DecodeString(combined)
	if err != nil {
		return combined, ""
	}
	var parsed struct {
		KeyID  string `json:"key_id"`
		Secret string `json:"secret"`
	}
	if json.Unmarshal(data, &parsed) == nil && parsed.KeyID != "" && parsed.Secret != "" {
		return parsed.KeyID, parsed.Secret
	}
	return combined, ""
}

// SetHTTPClient replaces the default http.Client used for requests.
func (c *Client) SetHTTPClient(hc *http.Client) {
	c.httpClient = hc
}

// doRequest executes a signed API request.
//
// method: HTTP method (GET, POST, etc.)
// path: API path (e.g. "/open/v1/compute/instances")
// query: optional query parameters (may be nil)
// body: optional request body; will be JSON-marshaled if non-nil
// result: pointer to the type T that will receive the response Data field
func (c *Client) doRequest(method, path string, query url.Values, body any, result any) error {
	// Marshal body
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("neburst: marshal body: %w", err)
		}
	}

	// Build full URL
	fullURL := c.baseURL + path
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	// Sign
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := uuidV4()
	sortedQ := buildSortedQueryString(query)
	bodyHash := sha256Hex(bodyBytes) // handles nil -> SHA256("")
	stringToSign := timestamp + "\n" + method + "\n" + path + "\n" + sortedQ + "\n" + bodyHash
	signature := hmacSHA256Hex(c.secret, stringToSign)

	// Build HTTP request
	var bodyReader io.Reader
	if bodyBytes != nil {
		bodyReader = bytes.NewReader(bodyBytes)
	}
	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("neburst: create request: %w", err)
	}

	req.Header.Set("User-Agent", "Neburst OpenAPI SDK/1.0")
	req.Header.Set("X-Nb-Key", c.keyID)
	req.Header.Set("X-Nb-Timestamp", timestamp)
	req.Header.Set("X-Nb-Nonce", nonce)
	req.Header.Set("X-Nb-Signature", signature)
	if bodyBytes != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("neburst: http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("neburst: read response: %w", err)
	}

	// Decode the envelope
	if result != nil {
		var envelope apiResponse[json.RawMessage]
		if err := json.Unmarshal(respBody, &envelope); err != nil {
			return fmt.Errorf("neburst: decode response: %w (status %d, body: %s)", err, resp.StatusCode, string(respBody))
		}
		if envelope.Code != 0 {
			return &APIError{Code: envelope.Code, Message: envelope.Msg}
		}
		if err := json.Unmarshal(envelope.Data, result); err != nil {
			return fmt.Errorf("neburst: decode data: %w", err)
		}
	} else {
		// No result expected, but still check for API errors
		var envelope apiResponse[json.RawMessage]
		if err := json.Unmarshal(respBody, &envelope); err != nil {
			return fmt.Errorf("neburst: decode response: %w (status %d, body: %s)", err, resp.StatusCode, string(respBody))
		}
		if envelope.Code != 0 {
			return &APIError{Code: envelope.Code, Message: envelope.Msg}
		}
	}

	return nil
}

// sha256Hex returns the hex-encoded SHA-256 hash of data.
// For nil or empty data, this returns SHA256("") which is the well-known constant.
func sha256Hex(data []byte) string {
	if data == nil {
		data = []byte{}
	}
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// hmacSHA256Hex computes HMAC-SHA256(secret, data) and returns the hex encoding.
func hmacSHA256Hex(secret, data string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

// buildSortedQueryString builds a sorted, URL-encoded query string from url.Values.
// Keys are sorted alphabetically. Each key-value pair is URL-encoded and joined with "&".
// If query is nil or empty, returns empty string "".
func buildSortedQueryString(query url.Values) string {
	if len(query) == 0 {
		return ""
	}

	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		parts = append(parts, url.QueryEscape(k)+"="+url.QueryEscape(query.Get(k)))
	}
	return strings.Join(parts, "&")
}

// uuidV4 generates a random UUID v4 string using crypto/rand.
func uuidV4() string {
	var u [16]byte
	_, _ = rand.Read(u[:])
	u[6] = (u[6] & 0x0f) | 0x40 // version 4
	u[8] = (u[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		u[0:4], u[4:6], u[6:8], u[8:10], u[10:16])
}
