// Package auth provides authentication and token management for the GCS CLI.
//
// This package handles:
//   - OAuth2 authentication flows using globus-go-sdk
//   - Token storage and retrieval (per-profile)
//   - Token validation and refresh
//
// Token files are stored at ~/.globus-connect-server/tokens/{profile}.json
// with 0600 permissions (user-only read/write).
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/scttfrdmn/globus-go-gcs/pkg/config"
	"github.com/scttfrdmn/globus-go-sdk/v3/pkg/services/auth"
)

const (
	// Token refresh buffer - refresh tokens this much before expiry
	tokenRefreshBuffer = 5 * time.Minute
)

// TokenInfo represents stored authentication tokens for a profile.
type TokenInfo struct {
	// AccessToken is the Bearer token for API requests.
	AccessToken string `json:"access_token"`

	// RefreshToken is used to obtain new access tokens (optional).
	// Only present when using offline_access scope.
	RefreshToken string `json:"refresh_token,omitempty"`

	// ExpiresAt is when the access token expires.
	ExpiresAt time.Time `json:"expires_at"`

	// Scopes are the OAuth2 scopes granted for this token.
	Scopes []string `json:"scopes,omitempty"`

	// ResourceServer is the resource server this token is for.
	ResourceServer string `json:"resource_server,omitempty"`
}

// IsValid returns true if the token is valid (not expired with buffer).
//
// Uses a 5-minute buffer to prevent edge cases where the token
// expires during an API request.
func (t *TokenInfo) IsValid() bool {
	if t == nil {
		return false
	}

	// Check if token will expire within the buffer window
	return time.Now().Add(tokenRefreshBuffer).Before(t.ExpiresAt)
}

// CanRefresh returns true if the token can be refreshed.
func (t *TokenInfo) CanRefresh() bool {
	return t != nil && t.RefreshToken != ""
}

// SaveToken saves token information to disk for a given profile.
//
// The token file is created with 0600 permissions (user-only read/write)
// for security.
func SaveToken(profile string, token *TokenInfo) error {
	// Ensure tokens directory exists
	if err := config.EnsureTokensDir(); err != nil {
		return fmt.Errorf("ensure tokens directory: %w", err)
	}

	// Get token file path
	tokenPath, err := config.GetTokenFilePath(profile)
	if err != nil {
		return fmt.Errorf("get token file path: %w", err)
	}

	// Marshal token to JSON
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal token: %w", err)
	}

	// Write with user-only permissions
	if err := os.WriteFile(tokenPath, data, 0600); err != nil {
		return fmt.Errorf("write token file: %w", err)
	}

	return nil
}

// LoadToken loads token information from disk for a given profile.
//
// Returns an error if the token file doesn't exist or is invalid.
func LoadToken(profile string) (*TokenInfo, error) {
	tokenPath, err := config.GetTokenFilePath(profile)
	if err != nil {
		return nil, fmt.Errorf("get token file path: %w", err)
	}

	// Read token file
	data, err := os.ReadFile(tokenPath) //nolint:gosec // Intentional file read from config directory
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not logged in (no token found for profile %q)", profile)
		}
		return nil, fmt.Errorf("read token file: %w", err)
	}

	// Unmarshal token
	var token TokenInfo
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("parse token file: %w", err)
	}

	return &token, nil
}

// DeleteToken deletes the token file for a given profile.
func DeleteToken(profile string) error {
	tokenPath, err := config.GetTokenFilePath(profile)
	if err != nil {
		return fmt.Errorf("get token file path: %w", err)
	}

	if err := os.Remove(tokenPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete token file: %w", err)
	}

	return nil
}

// RefreshTokenIfNeeded refreshes the token if it's expired or will expire soon.
//
// Uses the globus-go-sdk auth client to refresh the token.
// Returns true if the token was refreshed, false otherwise.
func RefreshTokenIfNeeded(ctx context.Context, profile string, authClient *auth.Client) (bool, error) {
	token, err := LoadToken(profile)
	if err != nil {
		return false, err
	}

	// Token is still valid
	if token.IsValid() {
		return false, nil
	}

	// Cannot refresh without refresh token
	if !token.CanRefresh() {
		return false, fmt.Errorf("token expired and cannot be refreshed (no refresh token)")
	}

	// Refresh the token
	tokenResp, err := authClient.RefreshToken(ctx, token.RefreshToken)
	if err != nil {
		return false, fmt.Errorf("refresh token: %w", err)
	}

	// Update token info
	newToken := &TokenInfo{
		AccessToken:    tokenResp.AccessToken,
		RefreshToken:   tokenResp.RefreshToken,
		ExpiresAt:      time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		Scopes:         token.Scopes,
		ResourceServer: tokenResp.ResourceServer,
	}

	// If refresh token not returned, keep the old one
	if newToken.RefreshToken == "" {
		newToken.RefreshToken = token.RefreshToken
	}

	// Save updated token
	if err := SaveToken(profile, newToken); err != nil {
		return false, fmt.Errorf("save refreshed token: %w", err)
	}

	return true, nil
}

// TokenFromAuthResponse converts an auth.TokenResponse to TokenInfo.
func TokenFromAuthResponse(resp *auth.TokenResponse) *TokenInfo {
	return &TokenInfo{
		AccessToken:    resp.AccessToken,
		RefreshToken:   resp.RefreshToken,
		ExpiresAt:      time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second),
		Scopes:         []string{resp.Scope},
		ResourceServer: resp.ResourceServer,
	}
}
