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

// SetEndpointOwner assigns the endpoint owner role to a specified principal.
func (c *Client) SetEndpointOwner(ctx context.Context, principalURN string) error {
	if principalURN == "" {
		return fmt.Errorf("principal URN is required")
	}

	payload := map[string]string{
		"principal": principalURN,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPut, "endpoint/owner", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("set endpoint owner: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}

// SetEndpointOwnerString sets a custom display name for the endpoint owner.
func (c *Client) SetEndpointOwnerString(ctx context.Context, ownerString string) error {
	if ownerString == "" {
		return fmt.Errorf("owner string is required")
	}

	payload := map[string]string{
		"owner_string": ownerString,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPut, "endpoint/owner-string", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("set endpoint owner string: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}

// ResetEndpointOwnerString resets the owner string to the default (ClientID).
func (c *Client) ResetEndpointOwnerString(ctx context.Context) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, "endpoint/owner-string", nil)
	if err != nil {
		return fmt.Errorf("reset endpoint owner string: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}

// SetSubscriptionID updates the subscription assignment for the endpoint.
func (c *Client) SetSubscriptionID(ctx context.Context, subscriptionID string) error {
	if subscriptionID == "" {
		return fmt.Errorf("subscription ID is required")
	}

	payload := map[string]string{
		"subscription_id": subscriptionID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPut, "endpoint/subscription", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("set subscription ID: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}

// SetupEndpointDomain configures a custom domain for the endpoint.
func (c *Client) SetupEndpointDomain(ctx context.Context, config *DomainConfig) error {
	if config == nil {
		return fmt.Errorf("domain configuration is required")
	}
	if config.Domain == "" {
		return fmt.Errorf("domain is required")
	}

	body, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "endpoint/domain", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("setup endpoint domain: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}

// GetEndpointDomain retrieves the custom domain configuration for the endpoint.
func (c *Client) GetEndpointDomain(ctx context.Context) (*DomainConfig, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "endpoint/domain", nil)
	if err != nil {
		return nil, fmt.Errorf("get endpoint domain: %w", err)
	}

	var domain DomainConfig
	if err := c.decodeResponse(resp, &domain); err != nil {
		return nil, err
	}

	return &domain, nil
}

// DeleteEndpointDomain removes the custom domain configuration from the endpoint.
func (c *Client) DeleteEndpointDomain(ctx context.Context) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, "endpoint/domain", nil)
	if err != nil {
		return fmt.Errorf("delete endpoint domain: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}
