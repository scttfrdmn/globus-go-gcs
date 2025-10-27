package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AuthPolicyList represents a list of authentication policies.
type AuthPolicyList struct {
	Data []AuthPolicy `json:"data"`
}

// ListAuthPolicies retrieves all authentication policies.
func (c *Client) ListAuthPolicies(ctx context.Context) (*AuthPolicyList, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "auth-policies", nil)
	if err != nil {
		return nil, fmt.Errorf("list auth policies: %w", err)
	}

	var list AuthPolicyList
	if err := c.decodeResponse(resp, &list); err != nil {
		return nil, err
	}

	return &list, nil
}

// GetAuthPolicy retrieves a specific authentication policy by ID.
func (c *Client) GetAuthPolicy(ctx context.Context, policyID string) (*AuthPolicy, error) {
	if policyID == "" {
		return nil, fmt.Errorf("policy ID is required")
	}

	path := fmt.Sprintf("auth-policies/%s", policyID)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get auth policy: %w", err)
	}

	var policy AuthPolicy
	if err := c.decodeResponse(resp, &policy); err != nil {
		return nil, err
	}

	return &policy, nil
}

// CreateAuthPolicy creates a new authentication policy.
func (c *Client) CreateAuthPolicy(ctx context.Context, policy *AuthPolicy) (*AuthPolicy, error) {
	if policy == nil {
		return nil, fmt.Errorf("policy is required")
	}

	body, err := json.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("marshal policy: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "auth-policies", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create auth policy: %w", err)
	}

	var created AuthPolicy
	if err := c.decodeResponse(resp, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

// UpdateAuthPolicy updates an existing authentication policy.
func (c *Client) UpdateAuthPolicy(ctx context.Context, policyID string, policy *AuthPolicy) (*AuthPolicy, error) {
	if policyID == "" {
		return nil, fmt.Errorf("policy ID is required")
	}
	if policy == nil {
		return nil, fmt.Errorf("policy is required")
	}

	body, err := json.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("marshal policy: %w", err)
	}

	path := fmt.Sprintf("auth-policies/%s", policyID)
	resp, err := c.doRequest(ctx, http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("update auth policy: %w", err)
	}

	var updated AuthPolicy
	if err := c.decodeResponse(resp, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

// DeleteAuthPolicy deletes an authentication policy.
func (c *Client) DeleteAuthPolicy(ctx context.Context, policyID string) error {
	if policyID == "" {
		return fmt.Errorf("policy ID is required")
	}

	path := fmt.Sprintf("auth-policies/%s", policyID)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("delete auth policy: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}
