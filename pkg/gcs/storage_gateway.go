package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ListStorageGatewaysOptions contains options for listing storage gateways.
type ListStorageGatewaysOptions struct {
	Filter   string // Filter storage gateways by name
	PageSize int    // Number of results per page
	Marker   string // Pagination marker
}

// StorageGatewayList represents a paginated list of storage gateways.
type StorageGatewayList struct {
	Data         []StorageGateway `json:"data"`
	HasNextPage  bool             `json:"has_next_page"`
	Marker       string           `json:"marker,omitempty"`
	TotalResults int              `json:"total,omitempty"`
}

// ListStorageGateways retrieves a list of storage gateways on the endpoint.
func (c *Client) ListStorageGateways(ctx context.Context, opts *ListStorageGatewaysOptions) (*StorageGatewayList, error) {
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

	path := "storage_gateways"
	if len(query) > 0 {
		path = fmt.Sprintf("%s?%s", path, query.Encode())
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list storage gateways: %w", err)
	}

	var list StorageGatewayList
	if err := c.decodeResponse(resp, &list); err != nil {
		return nil, err
	}

	return &list, nil
}

// GetStorageGateway retrieves a specific storage gateway by ID.
func (c *Client) GetStorageGateway(ctx context.Context, gatewayID string) (*StorageGateway, error) {
	if gatewayID == "" {
		return nil, fmt.Errorf("storage gateway ID is required")
	}

	path := fmt.Sprintf("storage_gateways/%s", gatewayID)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get storage gateway: %w", err)
	}

	var gateway StorageGateway
	if err := c.decodeResponse(resp, &gateway); err != nil {
		return nil, err
	}

	return &gateway, nil
}

// CreateStorageGateway creates a new storage gateway.
func (c *Client) CreateStorageGateway(ctx context.Context, gateway *StorageGateway) (*StorageGateway, error) {
	if gateway == nil {
		return nil, fmt.Errorf("storage gateway is required")
	}

	// Marshal the gateway to JSON
	body, err := json.Marshal(gateway)
	if err != nil {
		return nil, fmt.Errorf("marshal storage gateway: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "storage_gateways", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create storage gateway: %w", err)
	}

	var created StorageGateway
	if err := c.decodeResponse(resp, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

// UpdateStorageGateway updates an existing storage gateway.
func (c *Client) UpdateStorageGateway(ctx context.Context, gatewayID string, gateway *StorageGateway) (*StorageGateway, error) {
	if gatewayID == "" {
		return nil, fmt.Errorf("storage gateway ID is required")
	}
	if gateway == nil {
		return nil, fmt.Errorf("storage gateway is required")
	}

	// Marshal the gateway to JSON
	body, err := json.Marshal(gateway)
	if err != nil {
		return nil, fmt.Errorf("marshal storage gateway: %w", err)
	}

	path := fmt.Sprintf("storage_gateways/%s", gatewayID)
	resp, err := c.doRequest(ctx, http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("update storage gateway: %w", err)
	}

	var updated StorageGateway
	if err := c.decodeResponse(resp, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

// DeleteStorageGateway deletes a storage gateway.
func (c *Client) DeleteStorageGateway(ctx context.Context, gatewayID string) error {
	if gatewayID == "" {
		return fmt.Errorf("storage gateway ID is required")
	}

	path := fmt.Sprintf("storage_gateways/%s", gatewayID)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("delete storage gateway: %w", err)
	}

	// Close response body for DELETE requests
	defer func() { _ = resp.Body.Close() }()

	return nil
}
