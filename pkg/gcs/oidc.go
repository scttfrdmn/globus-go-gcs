package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetOIDCServer retrieves the OIDC server configuration.
func (c *Client) GetOIDCServer(ctx context.Context) (*OIDCServer, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "oidc", nil)
	if err != nil {
		return nil, fmt.Errorf("get OIDC server: %w", err)
	}

	var server OIDCServer
	if err := c.decodeResponse(resp, &server); err != nil {
		return nil, err
	}

	return &server, nil
}

// CreateOIDCServer creates a new OIDC server configuration.
func (c *Client) CreateOIDCServer(ctx context.Context, server *OIDCServer) (*OIDCServer, error) {
	if server == nil {
		return nil, fmt.Errorf("OIDC server configuration is required")
	}

	body, err := json.Marshal(server)
	if err != nil {
		return nil, fmt.Errorf("marshal OIDC server: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "oidc", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create OIDC server: %w", err)
	}

	var created OIDCServer
	if err := c.decodeResponse(resp, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

// RegisterOIDCServer registers an existing OIDC server with the endpoint.
func (c *Client) RegisterOIDCServer(ctx context.Context, server *OIDCServer) (*OIDCServer, error) {
	if server == nil {
		return nil, fmt.Errorf("OIDC server configuration is required")
	}

	body, err := json.Marshal(server)
	if err != nil {
		return nil, fmt.Errorf("marshal OIDC server: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "oidc/register", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("register OIDC server: %w", err)
	}

	var registered OIDCServer
	if err := c.decodeResponse(resp, &registered); err != nil {
		return nil, err
	}

	return &registered, nil
}

// UpdateOIDCServer updates the OIDC server configuration.
func (c *Client) UpdateOIDCServer(ctx context.Context, server *OIDCServer) (*OIDCServer, error) {
	if server == nil {
		return nil, fmt.Errorf("OIDC server configuration is required")
	}

	body, err := json.Marshal(server)
	if err != nil {
		return nil, fmt.Errorf("marshal OIDC server: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPatch, "oidc", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("update OIDC server: %w", err)
	}

	var updated OIDCServer
	if err := c.decodeResponse(resp, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

// DeleteOIDCServer deletes the OIDC server configuration.
func (c *Client) DeleteOIDCServer(ctx context.Context) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, "oidc", nil)
	if err != nil {
		return fmt.Errorf("delete OIDC server: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}
