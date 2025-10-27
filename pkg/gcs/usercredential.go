package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// UserCredentialList represents a list of user credentials.
type UserCredentialList struct {
	Data []UserCredential `json:"data"`
}

// ListUserCredentials retrieves all user credentials.
func (c *Client) ListUserCredentials(ctx context.Context) (*UserCredentialList, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "user-credentials", nil)
	if err != nil {
		return nil, fmt.Errorf("list user credentials: %w", err)
	}

	var list UserCredentialList
	if err := c.decodeResponse(resp, &list); err != nil {
		return nil, err
	}

	return &list, nil
}

// GetUserCredential retrieves a specific user credential.
func (c *Client) GetUserCredential(ctx context.Context, credentialID string) (*UserCredential, error) {
	if credentialID == "" {
		return nil, fmt.Errorf("credential ID is required")
	}

	path := fmt.Sprintf("user-credentials/%s", credentialID)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get user credential: %w", err)
	}

	var credential UserCredential
	if err := c.decodeResponse(resp, &credential); err != nil {
		return nil, err
	}

	return &credential, nil
}

// CreateActivescaleCredential creates an ActiveScale user credential.
func (c *Client) CreateActivescaleCredential(ctx context.Context, credential *UserCredential) (*UserCredential, error) {
	if credential == nil {
		return nil, fmt.Errorf("credential is required")
	}

	credential.Type = "activescale"
	body, err := json.Marshal(credential)
	if err != nil {
		return nil, fmt.Errorf("marshal credential: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "user-credentials", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create activescale credential: %w", err)
	}

	var created UserCredential
	if err := c.decodeResponse(resp, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

// CreateOAuthCredential creates an OAuth2 user credential.
func (c *Client) CreateOAuthCredential(ctx context.Context, credential *UserCredential) (*UserCredential, error) {
	if credential == nil {
		return nil, fmt.Errorf("credential is required")
	}

	credential.Type = "oauth"
	body, err := json.Marshal(credential)
	if err != nil {
		return nil, fmt.Errorf("marshal credential: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "user-credentials", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create oauth credential: %w", err)
	}

	var created UserCredential
	if err := c.decodeResponse(resp, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

// CreateS3Credential creates an S3 user credential.
func (c *Client) CreateS3Credential(ctx context.Context, credential *UserCredential) (*UserCredential, error) {
	if credential == nil {
		return nil, fmt.Errorf("credential is required")
	}

	credential.Type = "s3"
	body, err := json.Marshal(credential)
	if err != nil {
		return nil, fmt.Errorf("marshal credential: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "user-credentials", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create s3 credential: %w", err)
	}

	var created UserCredential
	if err := c.decodeResponse(resp, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

// AddS3Key adds an S3 IAM key to a credential.
func (c *Client) AddS3Key(ctx context.Context, credentialID string, key *S3Key) (*UserCredential, error) {
	if credentialID == "" {
		return nil, fmt.Errorf("credential ID is required")
	}
	if key == nil {
		return nil, fmt.Errorf("S3 key is required")
	}

	body, err := json.Marshal(key)
	if err != nil {
		return nil, fmt.Errorf("marshal key: %w", err)
	}

	path := fmt.Sprintf("user-credentials/%s/s3-keys", credentialID)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("add S3 key: %w", err)
	}

	var updated UserCredential
	if err := c.decodeResponse(resp, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

// UpdateS3Key updates an S3 IAM key.
func (c *Client) UpdateS3Key(ctx context.Context, credentialID, accessKeyID string, key *S3Key) (*UserCredential, error) {
	if credentialID == "" {
		return nil, fmt.Errorf("credential ID is required")
	}
	if accessKeyID == "" {
		return nil, fmt.Errorf("access key ID is required")
	}
	if key == nil {
		return nil, fmt.Errorf("S3 key is required")
	}

	body, err := json.Marshal(key)
	if err != nil {
		return nil, fmt.Errorf("marshal key: %w", err)
	}

	path := fmt.Sprintf("user-credentials/%s/s3-keys/%s", credentialID, accessKeyID)
	resp, err := c.doRequest(ctx, http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("update S3 key: %w", err)
	}

	var updated UserCredential
	if err := c.decodeResponse(resp, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

// DeleteS3Key deletes an S3 IAM key from a credential.
func (c *Client) DeleteS3Key(ctx context.Context, credentialID, accessKeyID string) error {
	if credentialID == "" {
		return fmt.Errorf("credential ID is required")
	}
	if accessKeyID == "" {
		return fmt.Errorf("access key ID is required")
	}

	path := fmt.Sprintf("user-credentials/%s/s3-keys/%s", credentialID, accessKeyID)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("delete S3 key: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}

// DeleteUserCredential deletes a user credential.
func (c *Client) DeleteUserCredential(ctx context.Context, credentialID string) error {
	if credentialID == "" {
		return fmt.Errorf("credential ID is required")
	}

	path := fmt.Sprintf("user-credentials/%s", credentialID)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("delete user credential: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}
