package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ListCollectionsOptions contains options for listing collections.
type ListCollectionsOptions struct {
	Filter   string // Filter collections by name
	PageSize int    // Number of results per page
	Marker   string // Pagination marker
}

// CollectionList represents a paginated list of collections.
type CollectionList struct {
	Data         []Collection `json:"data"`
	HasNextPage  bool         `json:"has_next_page"`
	Marker       string       `json:"marker,omitempty"`
	TotalResults int          `json:"total,omitempty"`
}

// ListCollections retrieves a list of collections on the endpoint.
func (c *Client) ListCollections(ctx context.Context, opts *ListCollectionsOptions) (*CollectionList, error) {
	// Build query parameters
	query := url.Values{}
	if opts != nil {
		if opts.Filter != "" {
			query.Set("filter", opts.Filter)
		}
		if opts.PageSize > 0 {
			query.Set("page_size", fmt.Sprintf("%d", opts.PageSize))
		}
		if opts.Marker != "" {
			query.Set("marker", opts.Marker)
		}
	}

	path := "collections"
	if len(query) > 0 {
		path = fmt.Sprintf("%s?%s", path, query.Encode())
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list collections: %w", err)
	}

	var list CollectionList
	if err := c.decodeResponse(resp, &list); err != nil {
		return nil, err
	}

	return &list, nil
}

// GetCollection retrieves a specific collection by ID.
func (c *Client) GetCollection(ctx context.Context, collectionID string) (*Collection, error) {
	if collectionID == "" {
		return nil, fmt.Errorf("collection ID is required")
	}

	path := fmt.Sprintf("collections/%s", collectionID)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get collection: %w", err)
	}

	var collection Collection
	if err := c.decodeResponse(resp, &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}

// CreateCollection creates a new collection.
func (c *Client) CreateCollection(ctx context.Context, collection *Collection) (*Collection, error) {
	if collection == nil {
		return nil, fmt.Errorf("collection is required")
	}

	// Marshal the collection to JSON
	body, err := json.Marshal(collection)
	if err != nil {
		return nil, fmt.Errorf("marshal collection: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "collections", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create collection: %w", err)
	}

	var created Collection
	if err := c.decodeResponse(resp, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

// UpdateCollection updates an existing collection.
func (c *Client) UpdateCollection(ctx context.Context, collectionID string, collection *Collection) (*Collection, error) {
	if collectionID == "" {
		return nil, fmt.Errorf("collection ID is required")
	}
	if collection == nil {
		return nil, fmt.Errorf("collection is required")
	}

	// Marshal the collection to JSON
	body, err := json.Marshal(collection)
	if err != nil {
		return nil, fmt.Errorf("marshal collection: %w", err)
	}

	path := fmt.Sprintf("collections/%s", collectionID)
	resp, err := c.doRequest(ctx, http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("update collection: %w", err)
	}

	var updated Collection
	if err := c.decodeResponse(resp, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

// DeleteCollection deletes a collection.
func (c *Client) DeleteCollection(ctx context.Context, collectionID string) error {
	if collectionID == "" {
		return fmt.Errorf("collection ID is required")
	}

	path := fmt.Sprintf("collections/%s", collectionID)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}

	// Close response body for DELETE requests
	defer func() { _ = resp.Body.Close() }()

	return nil
}

// CollectionValidation represents the result of a collection validation check.
type CollectionValidation struct {
	CollectionID string             `json:"collection_id"`
	Valid        bool               `json:"valid"`
	Errors       []ValidationError  `json:"errors,omitempty"`
	Warnings     []ValidationError  `json:"warnings,omitempty"`
}

// ValidationError represents a validation error or warning.
type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// CheckCollection validates a collection's configuration.
func (c *Client) CheckCollection(ctx context.Context, collectionID string) (*CollectionValidation, error) {
	if collectionID == "" {
		return nil, fmt.Errorf("collection ID is required")
	}

	path := fmt.Sprintf("collections/%s/check", collectionID)
	resp, err := c.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, fmt.Errorf("check collection: %w", err)
	}

	var validation CollectionValidation
	if err := c.decodeResponse(resp, &validation); err != nil {
		return nil, err
	}

	return &validation, nil
}

// BatchDeleteResult represents the result of a batch delete operation.
type BatchDeleteResult struct {
	Deleted []string            `json:"deleted"`
	Failed  []BatchDeleteError  `json:"failed,omitempty"`
}

// BatchDeleteError represents a failure in batch delete.
type BatchDeleteError struct {
	CollectionID string `json:"collection_id"`
	Error        string `json:"error"`
}

// BatchDeleteCollections deletes multiple collections in a single operation.
func (c *Client) BatchDeleteCollections(ctx context.Context, collectionIDs []string) (*BatchDeleteResult, error) {
	if len(collectionIDs) == 0 {
		return nil, fmt.Errorf("at least one collection ID is required")
	}

	payload := map[string][]string{
		"collection_ids": collectionIDs,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "collections/batch-delete", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("batch delete collections: %w", err)
	}

	var result BatchDeleteResult
	if err := c.decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SetCollectionOwner designates the owner of a collection.
func (c *Client) SetCollectionOwner(ctx context.Context, collectionID, principalURN string) error {
	if collectionID == "" {
		return fmt.Errorf("collection ID is required")
	}
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

	path := fmt.Sprintf("collections/%s/owner", collectionID)
	resp, err := c.doRequest(ctx, http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("set collection owner: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}

// SetCollectionOwnerString sets a custom display name for the collection owner.
func (c *Client) SetCollectionOwnerString(ctx context.Context, collectionID, ownerString string) error {
	if collectionID == "" {
		return fmt.Errorf("collection ID is required")
	}
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

	path := fmt.Sprintf("collections/%s/owner-string", collectionID)
	resp, err := c.doRequest(ctx, http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("set collection owner string: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}

// ResetCollectionOwnerString resets the owner string to the default.
func (c *Client) ResetCollectionOwnerString(ctx context.Context, collectionID string) error {
	if collectionID == "" {
		return fmt.Errorf("collection ID is required")
	}

	path := fmt.Sprintf("collections/%s/owner-string", collectionID)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("reset collection owner string: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}

// SetSubscriptionAdminVerified sets the subscription admin verification status for a collection.
func (c *Client) SetSubscriptionAdminVerified(ctx context.Context, collectionID string, verified bool) error {
	if collectionID == "" {
		return fmt.Errorf("collection ID is required")
	}

	payload := map[string]bool{
		"subscription_admin_verified": verified,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	path := fmt.Sprintf("collections/%s/subscription-admin-verified", collectionID)
	resp, err := c.doRequest(ctx, http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("set subscription admin verified: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}

// SetupCollectionDomain configures a custom domain for a collection.
func (c *Client) SetupCollectionDomain(ctx context.Context, collectionID string, config *DomainConfig) error {
	if collectionID == "" {
		return fmt.Errorf("collection ID is required")
	}
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

	path := fmt.Sprintf("collections/%s/domain", collectionID)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("setup collection domain: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}

// GetCollectionDomain retrieves the custom domain configuration for a collection.
func (c *Client) GetCollectionDomain(ctx context.Context, collectionID string) (*DomainConfig, error) {
	if collectionID == "" {
		return nil, fmt.Errorf("collection ID is required")
	}

	path := fmt.Sprintf("collections/%s/domain", collectionID)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get collection domain: %w", err)
	}

	var domain DomainConfig
	if err := c.decodeResponse(resp, &domain); err != nil {
		return nil, err
	}

	return &domain, nil
}

// DeleteCollectionDomain removes the custom domain configuration from a collection.
func (c *Client) DeleteCollectionDomain(ctx context.Context, collectionID string) error {
	if collectionID == "" {
		return fmt.Errorf("collection ID is required")
	}

	path := fmt.Sprintf("collections/%s/domain", collectionID)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("delete collection domain: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}
