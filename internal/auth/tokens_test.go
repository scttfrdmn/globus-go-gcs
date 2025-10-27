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

	os.Setenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR", testConfigDir) //nolint:errcheck,gosec // Test setup
	defer os.Unsetenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR")        //nolint:errcheck,gosec // Test cleanup

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

	os.Setenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR", testConfigDir) //nolint:errcheck,gosec // Test setup
	defer os.Unsetenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR")        //nolint:errcheck,gosec // Test cleanup

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

	os.Setenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR", testConfigDir) //nolint:errcheck,gosec // Test setup
	defer os.Unsetenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR")        //nolint:errcheck,gosec // Test cleanup

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

	os.Setenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR", testConfigDir) //nolint:errcheck,gosec // Test setup
	defer os.Unsetenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR")        //nolint:errcheck,gosec // Test cleanup

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

	os.Setenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR", testConfigDir) //nolint:errcheck,gosec // Test setup
	defer os.Unsetenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR")        //nolint:errcheck,gosec // Test cleanup

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

func TestSaveToken_IsEncrypted(t *testing.T) {
	// Use temporary directory for testing
	tmpDir := t.TempDir()
	testConfigDir := filepath.Join(tmpDir, "test-config")

	os.Setenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR", testConfigDir) //nolint:errcheck,gosec // Test setup
	defer os.Unsetenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR")        //nolint:errcheck,gosec // Test cleanup

	profile := "test-profile"
	token := &TokenInfo{
		AccessToken:    "secret_access_token_12345",
		RefreshToken:   "secret_refresh_token_67890",
		ExpiresAt:      time.Now().Add(1 * time.Hour),
		Scopes:         []string{"scope1", "scope2"},
		ResourceServer: "test.api.globus.org",
	}

	// Save token
	if err := SaveToken(profile, token); err != nil {
		t.Fatalf("SaveToken() error = %v", err)
	}

	// Read raw file content
	tokenPath := filepath.Join(testConfigDir, "tokens", profile+".json")
	data, err := os.ReadFile(tokenPath) //nolint:gosec // Intentional test file read from controlled test directory
	if err != nil {
		t.Fatalf("read token file: %v", err)
	}

	fileContent := string(data)

	// Debug: print file content
	t.Logf("Token file content:\n%s", fileContent)

	// Verify file contains encryption format marker (with flexible whitespace)
	if !contains(fileContent, "\"format\"") || !contains(fileContent, "encrypted-v1") {
		t.Errorf("Token file is missing encryption format marker. File content:\n%s", fileContent)
	}

	// Verify file contains encrypted_data field
	if !contains(fileContent, "encrypted_data") {
		t.Error("Token file is missing encrypted_data field")
	}

	// Verify plaintext tokens are NOT visible in file
	if contains(fileContent, token.AccessToken) {
		t.Error("Plaintext access token found in encrypted file - encryption failed!")
	}
	if contains(fileContent, token.RefreshToken) {
		t.Error("Plaintext refresh token found in encrypted file - encryption failed!")
	}

	// Verify we can still load and decrypt the token
	loaded, err := LoadToken(profile)
	if err != nil {
		t.Fatalf("LoadToken() error = %v", err)
	}

	if loaded.AccessToken != token.AccessToken {
		t.Errorf("Decrypted token doesn't match original")
	}
}

func TestLoadToken_AutoMigratePlaintext(t *testing.T) {
	// Use temporary directory for testing
	tmpDir := t.TempDir()
	testConfigDir := filepath.Join(tmpDir, "test-config")

	os.Setenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR", testConfigDir) //nolint:errcheck,gosec // Test setup
	defer os.Unsetenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR")        //nolint:errcheck,gosec // Test cleanup

	profile := "test-profile"

	// Create tokens directory
	tokensDir := filepath.Join(testConfigDir, "tokens")
	if err := os.MkdirAll(tokensDir, 0700); err != nil {
		t.Fatalf("create tokens directory: %v", err)
	}

	// Write a plaintext token file (v1.x format)
	tokenPath := filepath.Join(tokensDir, profile+".json")
	//nolint:gosec // Test fixture with intentional hardcoded test credentials
	plaintextToken := `{
  "access_token": "plaintext_access_token",
  "refresh_token": "plaintext_refresh_token",
  "expires_at": "2025-12-31T23:59:59Z",
  "scopes": ["scope1"],
  "resource_server": "test.api.globus.org"
}`
	if err := os.WriteFile(tokenPath, []byte(plaintextToken), 0600); err != nil {
		t.Fatalf("write plaintext token: %v", err)
	}

	// Verify it's plaintext
	data, _ := os.ReadFile(tokenPath) //nolint:gosec // Intentional test file read from controlled test directory
	if !contains(string(data), "plaintext_access_token") {
		t.Fatal("Test setup failed - plaintext token not written correctly")
	}

	// Load token (should auto-migrate to encrypted)
	loaded, err := LoadToken(profile)
	if err != nil {
		t.Fatalf("LoadToken() error = %v", err)
	}

	// Verify token was loaded correctly
	if loaded.AccessToken != "plaintext_access_token" {
		t.Errorf("AccessToken = %v, want plaintext_access_token", loaded.AccessToken)
	}
	if loaded.RefreshToken != "plaintext_refresh_token" {
		t.Errorf("RefreshToken = %v, want plaintext_refresh_token", loaded.RefreshToken)
	}

	// Re-read file to verify it's now encrypted
	data, err = os.ReadFile(tokenPath) //nolint:gosec // Intentional test file read from controlled test directory
	if err != nil {
		t.Fatalf("read migrated token file: %v", err)
	}

	fileContent := string(data)

	// Debug: print migrated file content
	t.Logf("Migrated token file content:\n%s", fileContent)

	// Verify file is now encrypted (with flexible whitespace)
	if !contains(fileContent, "\"format\"") || !contains(fileContent, "encrypted-v1") {
		t.Errorf("Token was not migrated to encrypted format. File content:\n%s", fileContent)
	}

	// Verify plaintext is no longer visible
	if contains(fileContent, "plaintext_access_token") {
		t.Error("Plaintext token still visible after migration - encryption failed!")
	}

	// Verify we can still load the migrated token
	loaded2, err := LoadToken(profile)
	if err != nil {
		t.Fatalf("LoadToken() after migration error = %v", err)
	}

	if loaded2.AccessToken != "plaintext_access_token" {
		t.Errorf("Migrated token AccessToken = %v, want plaintext_access_token", loaded2.AccessToken)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
