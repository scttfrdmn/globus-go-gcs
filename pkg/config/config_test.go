package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	tests := []struct {
		name    string
		envVar  string
		wantErr bool
	}{
		{
			name:    "default home directory",
			envVar:  "",
			wantErr: false,
		},
		{
			name:    "custom from environment",
			envVar:  "/tmp/globus-test",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.envVar != "" {
				defer os.Unsetenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR")    //nolint:errcheck,gosec // Test cleanup
				os.Setenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR", tt.envVar) //nolint:errcheck,gosec // Test setup
			}

			got, err := GetConfigDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfigDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.envVar != "" && got != tt.envVar {
				t.Errorf("GetConfigDir() = %v, want %v", got, tt.envVar)
			}

			if tt.envVar == "" {
				home, _ := os.UserHomeDir()
				want := filepath.Join(home, DefaultConfigDir)
				if got != want {
					t.Errorf("GetConfigDir() = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestGetTokensDir(t *testing.T) {
	configDir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir() error = %v", err)
	}

	want := filepath.Join(configDir, "tokens")
	got, err := GetTokensDir()

	if err != nil {
		t.Errorf("GetTokensDir() error = %v", err)
	}

	if got != want {
		t.Errorf("GetTokensDir() = %v, want %v", got, want)
	}
}

func TestLoadClientConfig(t *testing.T) {
	tests := []struct {
		name             string
		clientIDEnv      string
		clientSecretEnv  string
		wantClientID     string
		wantClientSecret string
		wantProfile      string
	}{
		{
			name:             "default configuration",
			clientIDEnv:      "",
			clientSecretEnv:  "",
			wantClientID:     DefaultClientID,
			wantClientSecret: "",
			wantProfile:      DefaultProfile,
		},
		{
			name:             "custom client ID from environment",
			clientIDEnv:      "custom-client-id",
			clientSecretEnv:  "",
			wantClientID:     "custom-client-id",
			wantClientSecret: "",
			wantProfile:      DefaultProfile,
		},
		{
			name:             "custom client ID and secret from environment",
			clientIDEnv:      "custom-client-id",
			clientSecretEnv:  "custom-client-secret",
			wantClientID:     "custom-client-id",
			wantClientSecret: "custom-client-secret",
			wantProfile:      DefaultProfile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.clientIDEnv != "" {
				defer os.Unsetenv("GLOBUS_CLIENT_ID")         //nolint:errcheck,gosec // Test cleanup
				os.Setenv("GLOBUS_CLIENT_ID", tt.clientIDEnv) //nolint:errcheck,gosec // Test setup
			}
			if tt.clientSecretEnv != "" {
				defer os.Unsetenv("GLOBUS_CLIENT_SECRET")             //nolint:errcheck,gosec // Test cleanup
				os.Setenv("GLOBUS_CLIENT_SECRET", tt.clientSecretEnv) //nolint:errcheck,gosec // Test setup
			}

			cfg, err := LoadClientConfig()
			if err != nil {
				t.Fatalf("LoadClientConfig() error = %v", err)
			}

			if cfg.ClientID != tt.wantClientID {
				t.Errorf("LoadClientConfig().ClientID = %v, want %v", cfg.ClientID, tt.wantClientID)
			}

			if cfg.ClientSecret != tt.wantClientSecret {
				t.Errorf("LoadClientConfig().ClientSecret = %v, want %v", cfg.ClientSecret, tt.wantClientSecret)
			}

			if cfg.Profile != tt.wantProfile {
				t.Errorf("LoadClientConfig().Profile = %v, want %v", cfg.Profile, tt.wantProfile)
			}

			if cfg.ConfigDir == "" {
				t.Error("LoadClientConfig().ConfigDir is empty")
			}
		})
	}
}

func TestGetTokenFilePath(t *testing.T) {
	tests := []struct {
		name    string
		profile string
		wantErr bool
	}{
		{
			name:    "default profile",
			profile: "default",
			wantErr: false,
		},
		{
			name:    "custom profile",
			profile: "production",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTokenFilePath(tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTokenFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				tokensDir, _ := GetTokensDir()
				want := filepath.Join(tokensDir, tt.profile+".json")
				if got != want {
					t.Errorf("GetTokenFilePath() = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestEnsureConfigDir(t *testing.T) {
	// Use temporary directory for testing
	tmpDir := t.TempDir()
	testConfigDir := filepath.Join(tmpDir, "test-config")

	os.Setenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR", testConfigDir) //nolint:errcheck,gosec // Test setup
	defer os.Unsetenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR")        //nolint:errcheck,gosec // Test cleanup

	err := EnsureConfigDir()
	if err != nil {
		t.Fatalf("EnsureConfigDir() error = %v", err)
	}

	// Check directory was created
	info, err := os.Stat(testConfigDir)
	if err != nil {
		t.Fatalf("config directory not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("config path is not a directory")
	}

	// Check permissions (on Unix systems)
	if info.Mode().Perm() != 0700 {
		t.Errorf("config directory permissions = %o, want 0700", info.Mode().Perm())
	}
}

func TestEnsureTokensDir(t *testing.T) {
	// Use temporary directory for testing
	tmpDir := t.TempDir()
	testConfigDir := filepath.Join(tmpDir, "test-config")
	testTokensDir := filepath.Join(testConfigDir, "tokens")

	os.Setenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR", testConfigDir) //nolint:errcheck,gosec // Test setup
	defer os.Unsetenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR")        //nolint:errcheck,gosec // Test cleanup

	err := EnsureTokensDir()
	if err != nil {
		t.Fatalf("EnsureTokensDir() error = %v", err)
	}

	// Check directory was created
	info, err := os.Stat(testTokensDir)
	if err != nil {
		t.Fatalf("tokens directory not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("tokens path is not a directory")
	}

	// Check permissions (on Unix systems)
	if info.Mode().Perm() != 0700 {
		t.Errorf("tokens directory permissions = %o, want 0700", info.Mode().Perm())
	}
}
