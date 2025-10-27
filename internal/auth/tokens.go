// Package auth provides authentication and token management for the GCS CLI.
//
// This package handles:
//   - OAuth2 authentication flows using globus-go-sdk
//   - Token storage and retrieval (per-profile)
//   - Token validation and refresh
//   - Token encryption at rest using AES-256-GCM
//
// Token files are stored at ~/.globus-connect-server/tokens/{profile}.json
// with 0600 permissions (user-only read/write).
//
// As of v2.0, tokens are encrypted at rest using AES-256-GCM with keys
// stored in the system keyring. This provides HIPAA/PHI compliance.
// Plaintext tokens from v1.x are automatically migrated on first load.
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

// EncryptedTokenFile represents the on-disk format of encrypted token files.
//
// This wraps the encrypted data with metadata to identify it as encrypted
// and support versioning for future changes to the encryption scheme.
type EncryptedTokenFile struct {
	// Format identifies this as an encrypted token file
	Format string `json:"format"`

	// EncryptedData contains the encrypted TokenInfo
	EncryptedData *EncryptedData `json:"encrypted_data"`
}

const (
	// EncryptedTokenFormat is the format identifier for encrypted tokens
	EncryptedTokenFormat = "encrypted-v1"
)

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
// As of v2.0, tokens are encrypted at rest using AES-256-GCM with keys
// stored in the system keyring. This provides HIPAA/PHI compliance.
//
// The token file is created with 0600 permissions (user-only read/write)
// for additional security.
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

	// Marshal token to JSON (plaintext, will be encrypted)
	plaintext, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("marshal token: %w", err)
	}

	// Encrypt the token data
	encryptedData, err := Encrypt(plaintext)
	if err != nil {
		return fmt.Errorf("encrypt token: %w", err)
	}

	// Wrap in encrypted file format
	encryptedFile := &EncryptedTokenFile{
		Format:        EncryptedTokenFormat,
		EncryptedData: encryptedData,
	}

	// Marshal encrypted file to JSON
	fileData, err := json.MarshalIndent(encryptedFile, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal encrypted token file: %w", err)
	}

	// Write with user-only permissions
	if err := os.WriteFile(tokenPath, fileData, 0600); err != nil {
		return fmt.Errorf("write token file: %w", err)
	}

	return nil
}

// LoadToken loads token information from disk for a given profile.
//
// As of v2.0, tokens are stored encrypted. This function supports automatic
// migration of plaintext tokens from v1.x - they will be detected, loaded,
// and automatically re-saved in encrypted format.
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

	// Try to parse as encrypted token file first
	var encryptedFile EncryptedTokenFile
	if err := json.Unmarshal(data, &encryptedFile); err == nil && encryptedFile.Format == EncryptedTokenFormat {
		// This is an encrypted token - decrypt it
		plaintext, err := Decrypt(encryptedFile.EncryptedData)
		if err != nil {
			return nil, fmt.Errorf("decrypt token: %w", err)
		}

		// Parse decrypted token
		var token TokenInfo
		if err := json.Unmarshal(plaintext, &token); err != nil {
			return nil, fmt.Errorf("parse decrypted token: %w", err)
		}

		return &token, nil
	}

	// Not encrypted - try to parse as plaintext token (v1.x format)
	var token TokenInfo
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("parse token file: %w (file may be corrupted)", err)
	}

	// Automatically migrate plaintext token to encrypted format
	// This is transparent to the user and happens on first load
	if err := SaveToken(profile, &token); err != nil {
		// Migration failed - log warning but still return the token
		// This allows the CLI to continue working even if keyring is unavailable
		fmt.Fprintf(os.Stderr, "Warning: Could not migrate token to encrypted format: %v\n", err)
		fmt.Fprintf(os.Stderr, "Your token is still stored in plaintext. To enable encryption, ensure your system keyring is available.\n")
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
