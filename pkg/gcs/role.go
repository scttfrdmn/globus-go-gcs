package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ListRolesOptions contains options for listing roles.
type ListRolesOptions struct {
	Collection string // Filter roles by collection ID
	Principal  string // Filter roles by principal (identity)
	PageSize   int    // Number of results per page
	Marker     string // Pagination marker
}

// RoleList represents a paginated list of roles.
type RoleList struct {
	Data         []Role `json:"data"`
	HasNextPage  bool   `json:"has_next_page"`
	Marker       string `json:"marker,omitempty"`
	TotalResults int    `json:"total,omitempty"`
}

// ListRoles retrieves a list of roles on the endpoint.
func (c *Client) ListRoles(ctx context.Context, opts *ListRolesOptions) (*RoleList, error) {
	// Build query parameters
	query := url.Values{}
	if opts != nil {
		if opts.Collection != "" {
			query.Set("collection", opts.Collection)
		}
		if opts.Principal != "" {
			query.Set("principal", opts.Principal)
		}
		if opts.PageSize > 0 {
			query.Set("page_size", fmt.Sprintf("%d", opts.PageSize))
		}
		if opts.Marker != "" {
			query.Set("marker", opts.Marker)
		}
	}

	path := "roles"
	if len(query) > 0 {
		path = fmt.Sprintf("%s?%s", path, query.Encode())
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}

	var list RoleList
	if err := c.decodeResponse(resp, &list); err != nil {
		return nil, err
	}

	return &list, nil
}

// GetRole retrieves a specific role by ID.
func (c *Client) GetRole(ctx context.Context, roleID string) (*Role, error) {
	if roleID == "" {
		return nil, fmt.Errorf("role ID is required")
	}

	path := fmt.Sprintf("roles/%s", roleID)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get role: %w", err)
	}

	var role Role
	if err := c.decodeResponse(resp, &role); err != nil {
		return nil, err
	}

	return &role, nil
}

// CreateRole creates a new role assignment.
func (c *Client) CreateRole(ctx context.Context, role *Role) (*Role, error) {
	if role == nil {
		return nil, fmt.Errorf("role is required")
	}

	// Marshal the role to JSON
	body, err := json.Marshal(role)
	if err != nil {
		return nil, fmt.Errorf("marshal role: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "roles", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create role: %w", err)
	}

	var created Role
	if err := c.decodeResponse(resp, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

// UpdateRole updates an existing role assignment.
func (c *Client) UpdateRole(ctx context.Context, roleID string, role *Role) (*Role, error) {
	if roleID == "" {
		return nil, fmt.Errorf("role ID is required")
	}
	if role == nil {
		return nil, fmt.Errorf("role is required")
	}

	// Marshal the role to JSON
	body, err := json.Marshal(role)
	if err != nil {
		return nil, fmt.Errorf("marshal role: %w", err)
	}

	path := fmt.Sprintf("roles/%s", roleID)
	resp, err := c.doRequest(ctx, http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("update role: %w", err)
	}

	var updated Role
	if err := c.decodeResponse(resp, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

// DeleteRole deletes a role assignment.
func (c *Client) DeleteRole(ctx context.Context, roleID string) error {
	if roleID == "" {
		return fmt.Errorf("role ID is required")
	}

	path := fmt.Sprintf("roles/%s", roleID)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("delete role: %w", err)
	}

	// Close response body for DELETE requests
	defer func() { _ = resp.Body.Close() }()

	return nil
}
