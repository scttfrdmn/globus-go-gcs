package gcs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client is a client for the Globus Connect Server Manager API.
// Unlike other Globus service clients, the GCS Manager API is hosted on
// individual GCS endpoint hosts rather than a centralized service.
type Client struct {
	baseURL     string
	httpClient  *http.Client
	accessToken string
	userAgent   string
}

// NewClient creates a new GCS Manager API client.
// The endpointFQDN is the fully qualified domain name of the GCS endpoint
// (e.g., "abc.def.data.globus.org").
func NewClient(endpointFQDN string, opts ...ClientOption) (*Client, error) {
	if endpointFQDN == "" {
		return nil, fmt.Errorf("endpoint FQDN is required")
	}

	// Apply default options
	options := defaultOptions()

	// Apply user options
	for _, opt := range opts {
		opt(options)
	}

	// Construct base URL
	baseURL := fmt.Sprintf("https://%s/api/", endpointFQDN)

	client := &Client{
		baseURL:     baseURL,
		httpClient:  options.httpClient,
		accessToken: options.accessToken,
		userAgent:   options.userAgent,
	}

	return client, nil
}

// SetAccessToken sets the access token for authentication.
// This can be used to update the token after the client is created.
func (c *Client) SetAccessToken(token string) {
	c.accessToken = token
}

// doRequest performs an HTTP request with authentication.
func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	// Construct full URL
	url := c.baseURL + strings.TrimPrefix(path, "/")

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", c.userAgent)
	if c.accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		defer func() { _ = resp.Body.Close() }()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

// decodeResponse decodes a JSON response into the target struct.
func (c *Client) decodeResponse(resp *http.Response, target interface{}) error {
	defer func() { _ = resp.Body.Close() }()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}
