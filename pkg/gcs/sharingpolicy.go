package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SharingPolicyList represents a list of sharing policies.
type SharingPolicyList struct {
	Data []SharingPolicy `json:"data"`
}

// ListSharingPolicies retrieves all sharing policies.
func (c *Client) ListSharingPolicies(ctx context.Context) (*SharingPolicyList, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "sharing-policies", nil)
	if err != nil {
		return nil, fmt.Errorf("list sharing policies: %w", err)
	}

	var list SharingPolicyList
	if err := c.decodeResponse(resp, &list); err != nil {
		return nil, err
	}

	return &list, nil
}

// GetSharingPolicy retrieves a specific sharing policy.
func (c *Client) GetSharingPolicy(ctx context.Context, policyID string) (*SharingPolicy, error) {
	if policyID == "" {
		return nil, fmt.Errorf("policy ID is required")
	}

	path := fmt.Sprintf("sharing-policies/%s", policyID)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get sharing policy: %w", err)
	}

	var policy SharingPolicy
	if err := c.decodeResponse(resp, &policy); err != nil {
		return nil, err
	}

	return &policy, nil
}

// CreateSharingPolicy creates a new sharing policy.
func (c *Client) CreateSharingPolicy(ctx context.Context, policy *SharingPolicy) (*SharingPolicy, error) {
	if policy == nil {
		return nil, fmt.Errorf("policy is required")
	}

	body, err := json.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("marshal policy: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "sharing-policies", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create sharing policy: %w", err)
	}

	var created SharingPolicy
	if err := c.decodeResponse(resp, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

// DeleteSharingPolicy deletes a sharing policy.
func (c *Client) DeleteSharingPolicy(ctx context.Context, policyID string) error {
	if policyID == "" {
		return fmt.Errorf("policy ID is required")
	}

	path := fmt.Sprintf("sharing-policies/%s", policyID)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("delete sharing policy: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}
