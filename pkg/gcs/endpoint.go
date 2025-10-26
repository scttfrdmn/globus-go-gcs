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
