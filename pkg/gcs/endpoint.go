package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetInfo retrieves the GCS Manager service information.
// This endpoint does not require authentication.
func (c *Client) GetInfo(ctx context.Context) (*Info, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "info", nil)
	if err != nil {
		return nil, fmt.Errorf("get info: %w", err)
	}

	var info Info
	if err := c.decodeResponse(resp, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

// GetEndpoint retrieves the endpoint configuration.
func (c *Client) GetEndpoint(ctx context.Context) (*Endpoint, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "endpoint", nil)
	if err != nil {
		return nil, fmt.Errorf("get endpoint: %w", err)
	}

	var endpoint Endpoint
	if err := c.decodeResponse(resp, &endpoint); err != nil {
		return nil, err
	}

	return &endpoint, nil
}

// UpdateEndpoint updates the endpoint configuration.
func (c *Client) UpdateEndpoint(ctx context.Context, endpoint *Endpoint) (*Endpoint, error) {
	// Marshal the endpoint to JSON
	body, err := json.Marshal(endpoint)
	if err != nil {
		return nil, fmt.Errorf("marshal endpoint: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPatch, "endpoint", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("update endpoint: %w", err)
	}

	var updated Endpoint
	if err := c.decodeResponse(resp, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

// SetupEndpoint creates and initializes a new GCS endpoint.
func (c *Client) SetupEndpoint(ctx context.Context, endpoint *Endpoint) (*Endpoint, error) {
	if endpoint == nil {
		return nil, fmt.Errorf("endpoint configuration is required")
	}

	// Marshal the endpoint to JSON
	body, err := json.Marshal(endpoint)
	if err != nil {
		return nil, fmt.Errorf("marshal endpoint: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "endpoint", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("setup endpoint: %w", err)
	}

	var created Endpoint
	if err := c.decodeResponse(resp, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

// CleanupEndpoint permanently removes the endpoint configuration.
func (c *Client) CleanupEndpoint(ctx context.Context) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, "endpoint", nil)
	if err != nil {
		return fmt.Errorf("cleanup endpoint: %w", err)
	}

	// Close response body for DELETE requests
	defer func() { _ = resp.Body.Close() }()

	return nil
}

// DeploymentKeyResult represents the result of a key conversion operation.
type DeploymentKeyResult struct {
	OldKey string `json:"old_key,omitempty"`
	NewKey string `json:"new_key"`
}

// ConvertDeploymentKey converts an old deployment key to a new one.
func (c *Client) ConvertDeploymentKey(ctx context.Context, oldKey string) (*DeploymentKeyResult, error) {
	if oldKey == "" {
		return nil, fmt.Errorf("old deployment key is required")
	}

	payload := map[string]string{
		"old_key": oldKey,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "endpoint/key-convert", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("convert deployment key: %w", err)
	}

	var result DeploymentKeyResult
	if err := c.decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
