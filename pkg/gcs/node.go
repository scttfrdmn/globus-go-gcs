package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ListNodesOptions contains options for listing nodes.
type ListNodesOptions struct {
	Filter   string // Filter nodes by name
	PageSize int    // Number of results per page
	Marker   string // Pagination marker
}

// NodeList represents a paginated list of nodes.
type NodeList struct {
	Data         []Node `json:"data"`
	HasNextPage  bool   `json:"has_next_page"`
	Marker       string `json:"marker,omitempty"`
	TotalResults int    `json:"total,omitempty"`
}

// ListNodes retrieves a list of nodes on the endpoint.
func (c *Client) ListNodes(ctx context.Context, opts *ListNodesOptions) (*NodeList, error) {
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

	path := "nodes"
	if len(query) > 0 {
		path = fmt.Sprintf("%s?%s", path, query.Encode())
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}

	var list NodeList
	if err := c.decodeResponse(resp, &list); err != nil {
		return nil, err
	}

	return &list, nil
}

// GetNode retrieves a specific node by ID.
func (c *Client) GetNode(ctx context.Context, nodeID string) (*Node, error) {
	if nodeID == "" {
		return nil, fmt.Errorf("node ID is required")
	}

	path := fmt.Sprintf("nodes/%s", nodeID)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get node: %w", err)
	}

	var node Node
	if err := c.decodeResponse(resp, &node); err != nil {
		return nil, err
	}

	return &node, nil
}

// CreateNode creates a new node.
func (c *Client) CreateNode(ctx context.Context, node *Node) (*Node, error) {
	if node == nil {
		return nil, fmt.Errorf("node is required")
	}

	// Marshal the node to JSON
	body, err := json.Marshal(node)
	if err != nil {
		return nil, fmt.Errorf("marshal node: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "nodes", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create node: %w", err)
	}

	var created Node
	if err := c.decodeResponse(resp, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

// UpdateNode updates an existing node.
func (c *Client) UpdateNode(ctx context.Context, nodeID string, node *Node) (*Node, error) {
	if nodeID == "" {
		return nil, fmt.Errorf("node ID is required")
	}
	if node == nil {
		return nil, fmt.Errorf("node is required")
	}

	// Marshal the node to JSON
	body, err := json.Marshal(node)
	if err != nil {
		return nil, fmt.Errorf("marshal node: %w", err)
	}

	path := fmt.Sprintf("nodes/%s", nodeID)
	resp, err := c.doRequest(ctx, http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("update node: %w", err)
	}

	var updated Node
	if err := c.decodeResponse(resp, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

// DeleteNode deletes a node.
func (c *Client) DeleteNode(ctx context.Context, nodeID string) error {
	if nodeID == "" {
		return fmt.Errorf("node ID is required")
	}

	path := fmt.Sprintf("nodes/%s", nodeID)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("delete node: %w", err)
	}

	// Close response body for DELETE requests
	defer func() { _ = resp.Body.Close() }()

	return nil
}

// SetupNode configures and initializes a new node.
func (c *Client) SetupNode(ctx context.Context, node *Node) (*Node, error) {
	if node == nil {
		return nil, fmt.Errorf("node configuration is required")
	}

	// Marshal the node to JSON
	body, err := json.Marshal(node)
	if err != nil {
		return nil, fmt.Errorf("marshal node: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "nodes/setup", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("setup node: %w", err)
	}

	var created Node
	if err := c.decodeResponse(resp, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

// CleanupNode removes a node and its configuration.
func (c *Client) CleanupNode(ctx context.Context, nodeID string) error {
	if nodeID == "" {
		return fmt.Errorf("node ID is required")
	}

	path := fmt.Sprintf("nodes/%s/cleanup", nodeID)
	resp, err := c.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return fmt.Errorf("cleanup node: %w", err)
	}

	// Close response body
	defer func() { _ = resp.Body.Close() }()

	return nil
}

// EnableNode activates a node for data transfers.
func (c *Client) EnableNode(ctx context.Context, nodeID string) error {
	if nodeID == "" {
		return fmt.Errorf("node ID is required")
	}

	path := fmt.Sprintf("nodes/%s/enable", nodeID)
	resp, err := c.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return fmt.Errorf("enable node: %w", err)
	}

	// Close response body
	defer func() { _ = resp.Body.Close() }()

	return nil
}

// DisableNode deactivates a node, preventing new transfers.
func (c *Client) DisableNode(ctx context.Context, nodeID string) error {
	if nodeID == "" {
		return fmt.Errorf("node ID is required")
	}

	path := fmt.Sprintf("nodes/%s/disable", nodeID)
	resp, err := c.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return fmt.Errorf("disable node: %w", err)
	}

	// Close response body
	defer func() { _ = resp.Body.Close() }()

	return nil
}

// NodeSecret represents a generated node authentication secret.
type NodeSecret struct {
	NodeID string `json:"node_id"`
	Secret string `json:"secret"`
}

// GenerateNodeSecret generates a new authentication secret for a node.
func (c *Client) GenerateNodeSecret(ctx context.Context, nodeID string) (*NodeSecret, error) {
	if nodeID == "" {
		return nil, fmt.Errorf("node ID is required")
	}

	path := fmt.Sprintf("nodes/%s/new-secret", nodeID)
	resp, err := c.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, fmt.Errorf("generate node secret: %w", err)
	}

	var secret NodeSecret
	if err := c.decodeResponse(resp, &secret); err != nil {
		return nil, err
	}

	return &secret, nil
}
