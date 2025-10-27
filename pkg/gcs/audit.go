package gcs

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// AuditLogList represents a list of audit log entries.
type AuditLogList struct {
	Data []AuditLog `json:"data"`
}

// AuditQueryParams represents query parameters for fetching audit logs.
type AuditQueryParams struct {
	StartTime  *time.Time
	EndTime    *time.Time
	EventType  string
	IdentityID string
	ResourceID string
	Action     string
	Result     string
	Limit      int
}

// GetAuditLogs retrieves audit logs from the GCS Manager API.
func (c *Client) GetAuditLogs(ctx context.Context, params *AuditQueryParams) (*AuditLogList, error) {
	// Build query parameters
	query := url.Values{}
	if params != nil {
		if params.StartTime != nil {
			query.Set("start_time", params.StartTime.Format(time.RFC3339))
		}
		if params.EndTime != nil {
			query.Set("end_time", params.EndTime.Format(time.RFC3339))
		}
		if params.EventType != "" {
			query.Set("event_type", params.EventType)
		}
		if params.IdentityID != "" {
			query.Set("identity_id", params.IdentityID)
		}
		if params.ResourceID != "" {
			query.Set("resource_id", params.ResourceID)
		}
		if params.Action != "" {
			query.Set("action", params.Action)
		}
		if params.Result != "" {
			query.Set("result", params.Result)
		}
		if params.Limit > 0 {
			query.Set("limit", fmt.Sprintf("%d", params.Limit))
		}
	}

	path := "audit-logs"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get audit logs: %w", err)
	}

	var list AuditLogList
	if err := c.decodeResponse(resp, &list); err != nil {
		return nil, err
	}

	return &list, nil
}
