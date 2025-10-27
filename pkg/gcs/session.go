package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetSession retrieves the current CLI authentication session.
func (c *Client) GetSession(ctx context.Context) (*Session, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "session", nil)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	var session Session
	if err := c.decodeResponse(resp, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// UpdateSession updates the current session settings.
func (c *Client) UpdateSession(ctx context.Context, session *Session) (*Session, error) {
	if session == nil {
		return nil, fmt.Errorf("session is required")
	}

	body, err := json.Marshal(session)
	if err != nil {
		return nil, fmt.Errorf("marshal session: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPatch, "session", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("update session: %w", err)
	}

	var updated Session
	if err := c.decodeResponse(resp, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

// UpdateSessionConsents updates the consents for the current session.
func (c *Client) UpdateSessionConsents(ctx context.Context, consents []string) (*Session, error) {
	if len(consents) == 0 {
		return nil, fmt.Errorf("at least one consent is required")
	}

	payload := map[string][]string{
		"consents": consents,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal consents: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "session/consent", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("update session consents: %w", err)
	}

	var updated Session
	if err := c.decodeResponse(resp, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}
