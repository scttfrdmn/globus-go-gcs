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
