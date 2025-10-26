package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTokenInfo_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		token     *TokenInfo
		wantValid bool
	}{
		{
			name:      "nil token",
			token:     nil,
			wantValid: false,
		},
		{
			name: "valid token (expires in 1 hour)",
			token: &TokenInfo{
				AccessToken: "valid-token",
				ExpiresAt:   time.Now().Add(1 * time.Hour),
			},
			wantValid: true,
		},
		{
			name: "expired token",
			token: &TokenInfo{
				AccessToken: "expired-token",
				ExpiresAt:   time.Now().Add(-1 * time.Hour),
			},
			wantValid: false,
		},
		{
			name: "token expires within buffer (4 minutes)",
			token: &TokenInfo{
				AccessToken: "soon-expired-token",
				ExpiresAt:   time.Now().Add(4 * time.Minute),
			},
			wantValid: false,
		},
		{
			name: "token expires just after buffer (6 minutes)",
			token: &TokenInfo{
				AccessToken: "valid-token",
				ExpiresAt:   time.Now().Add(6 * time.Minute),
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.token.IsValid(); got != tt.wantValid {
				t.Errorf("TokenInfo.IsValid() = %v, want %v", got, tt.wantValid)
			}
		})
	}
}

func TestTokenInfo_CanRefresh(t *testing.T) {
	tests := []struct {
		name    string
		token   *TokenInfo
		wantCan bool
	}{
		{
			name:    "nil token",
			token:   nil,
			wantCan: false,
		},
		{
			name: "token with refresh token",
			token: &TokenInfo{
				AccessToken:  "access",
				RefreshToken: "refresh",
			},
			wantCan: true,
		},
		{
			name: "token without refresh token",
			token: &TokenInfo{
				AccessToken: "access",
			},
			wantCan: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.token.CanRefresh(); got != tt.wantCan {
				t.Errorf("TokenInfo.CanRefresh() = %v, want %v", got, tt.wantCan)
			}
		})
	}
}

func TestSaveAndLoadToken(t *testing.T) {
	// Use temporary directory for testing
	tmpDir := t.TempDir()
	testConfigDir := filepath.Join(tmpDir, "test-config")

	os.Setenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR", testConfigDir)
	defer os.Unsetenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR")

	profile := "test-profile"
	token := &TokenInfo{
		AccessToken:    "test-access-token",
		RefreshToken:   "test-refresh-token",
		ExpiresAt:      time.Now().Add(1 * time.Hour),
		Scopes:         []string{"scope1", "scope2"},
		ResourceServer: "test.api.globus.org",
	}

	// Save token
	err := SaveToken(profile, token)
	if err != nil {
		t.Fatalf("SaveToken() error = %v", err)
	}

	// Load token
	loaded, err := LoadToken(profile)
	if err != nil {
		t.Fatalf("LoadToken() error = %v", err)
	}

	// Compare tokens
	if loaded.AccessToken != token.AccessToken {
		t.Errorf("AccessToken = %v, want %v", loaded.AccessToken, token.AccessToken)
	}
	if loaded.RefreshToken != token.RefreshToken {
		t.Errorf("RefreshToken = %v, want %v", loaded.RefreshToken, token.RefreshToken)
	}
	if loaded.ResourceServer != token.ResourceServer {
		t.Errorf("ResourceServer = %v, want %v", loaded.ResourceServer, token.ResourceServer)
	}

	// ExpiresAt may have slight differences due to JSON serialization
	diff := loaded.ExpiresAt.Sub(token.ExpiresAt)
	if diff > time.Second || diff < -time.Second {
		t.Errorf("ExpiresAt difference too large: %v", diff)
	}

	// Check scopes
	if len(loaded.Scopes) != len(token.Scopes) {
		t.Errorf("Scopes length = %v, want %v", len(loaded.Scopes), len(token.Scopes))
	}
}

func TestLoadToken_NotExists(t *testing.T) {
	// Use temporary directory for testing
	tmpDir := t.TempDir()
	testConfigDir := filepath.Join(tmpDir, "test-config")

	os.Setenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR", testConfigDir)
	defer os.Unsetenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR")

	profile := "nonexistent-profile"

	_, err := LoadToken(profile)
	if err == nil {
		t.Error("LoadToken() expected error for nonexistent profile, got nil")
	}
}

func TestDeleteToken(t *testing.T) {
	// Use temporary directory for testing
	tmpDir := t.TempDir()
	testConfigDir := filepath.Join(tmpDir, "test-config")

	os.Setenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR", testConfigDir)
	defer os.Unsetenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR")

	profile := "test-profile"
	token := &TokenInfo{
		AccessToken: "test-token",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}

	// Save token
	if err := SaveToken(profile, token); err != nil {
		t.Fatalf("SaveToken() error = %v", err)
	}

	// Delete token
	if err := DeleteToken(profile); err != nil {
		t.Fatalf("DeleteToken() error = %v", err)
	}

	// Verify token is deleted
	_, err := LoadToken(profile)
	if err == nil {
		t.Error("LoadToken() expected error after DeleteToken(), got nil")
	}
}

func TestDeleteToken_NotExists(t *testing.T) {
	// Use temporary directory for testing
	tmpDir := t.TempDir()
	testConfigDir := filepath.Join(tmpDir, "test-config")

	os.Setenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR", testConfigDir)
	defer os.Unsetenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR")

	profile := "nonexistent-profile"

	// Should not error when deleting nonexistent token
	err := DeleteToken(profile)
	if err != nil {
		t.Errorf("DeleteToken() unexpected error for nonexistent token: %v", err)
	}
}

func TestSaveToken_FilePermissions(t *testing.T) {
	// Use temporary directory for testing
	tmpDir := t.TempDir()
	testConfigDir := filepath.Join(tmpDir, "test-config")

	os.Setenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR", testConfigDir)
	defer os.Unsetenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR")

	profile := "test-profile"
	token := &TokenInfo{
		AccessToken: "test-token",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}

	// Save token
	if err := SaveToken(profile, token); err != nil {
		t.Fatalf("SaveToken() error = %v", err)
	}

	// Check file permissions
	tokenPath := filepath.Join(testConfigDir, "tokens", profile+".json")
	info, err := os.Stat(tokenPath)
	if err != nil {
		t.Fatalf("stat token file: %v", err)
	}

	// Check permissions are 0600 (user-only read/write)
	if info.Mode().Perm() != 0600 {
		t.Errorf("token file permissions = %o, want 0600", info.Mode().Perm())
	}
}
