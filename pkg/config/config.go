// Package config provides configuration management for the Globus Connect Server CLI.
//
// Configuration is loaded from multiple sources with the following precedence:
//  1. Environment variables (GLOBUS_CLIENT_ID, GLOBUS_CLIENT_SECRET, etc.)
//  2. Configuration file (~/.globus-connect-server/config.yaml)
//  3. Default values
//
// The configuration directory structure follows Python CLI compatibility:
//
//	~/.globus-connect-server/
//	├── config.yaml           # CLI configuration
//	├── tokens/               # Token storage (per profile)
//	│   └── default.json      # Default profile tokens
//	└── deployment-key.json   # Optional: endpoint deployment key
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// DefaultClientID is the public client ID for the GCS CLI.
	// This matches the Python CLI's client ID for compatibility.
	DefaultClientID = "e6c75d97-532a-4c88-b031-f5a3014430e3"

	// DefaultConfigDir is the directory where CLI configuration is stored.
	DefaultConfigDir = ".globus-connect-server"

	// DefaultProfile is the name of the default profile.
	DefaultProfile = "default"
)

// Config represents the CLI configuration.
type Config struct {
	// ClientID is the Globus Auth OAuth2 client ID.
	ClientID string `json:"client_id" yaml:"client_id"`

	// ClientSecret is the Globus Auth OAuth2 client secret (optional).
	// If not provided, uses public client flow.
	ClientSecret string `json:"client_secret,omitempty" yaml:"client_secret,omitempty"`

	// ConfigDir is the directory where configuration files are stored.
	// Defaults to ~/.globus-connect-server
	ConfigDir string `json:"-" yaml:"-"`

	// Profile is the active profile name.
	// Defaults to "default".
	Profile string `json:"-" yaml:"-"`
}

// GetConfigDir returns the configuration directory path.
//
// Priority:
//  1. GLOBUS_CONNECT_SERVER_CONFIG_DIR environment variable
//  2. $HOME/.globus-connect-server
func GetConfigDir() (string, error) {
	// Check environment variable first
	if dir := os.Getenv("GLOBUS_CONNECT_SERVER_CONFIG_DIR"); dir != "" {
		return dir, nil
	}

	// Use default in home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}

	return filepath.Join(home, DefaultConfigDir), nil
}

// GetTokensDir returns the tokens directory path.
func GetTokensDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "tokens"), nil
}

// EnsureConfigDir creates the configuration directory if it doesn't exist.
func EnsureConfigDir() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Create directory with user-only permissions
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	return nil
}

// EnsureTokensDir creates the tokens directory if it doesn't exist.
func EnsureTokensDir() error {
	tokensDir, err := GetTokensDir()
	if err != nil {
		return err
	}

	// Create directory with user-only permissions
	if err := os.MkdirAll(tokensDir, 0700); err != nil {
		return fmt.Errorf("create tokens directory: %w", err)
	}

	return nil
}

// LoadClientConfig loads client configuration from environment variables and defaults.
//
// Priority:
//  1. GLOBUS_CLIENT_ID, GLOBUS_CLIENT_SECRET environment variables
//  2. Default client ID (public client)
func LoadClientConfig() (*Config, error) {
	cfg := &Config{
		ClientID: DefaultClientID,
		Profile:  DefaultProfile,
	}

	// Override from environment if set
	if clientID := os.Getenv("GLOBUS_CLIENT_ID"); clientID != "" {
		cfg.ClientID = clientID
	}

	if clientSecret := os.Getenv("GLOBUS_CLIENT_SECRET"); clientSecret != "" {
		cfg.ClientSecret = clientSecret
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}
	cfg.ConfigDir = configDir

	return cfg, nil
}

// GetTokenFilePath returns the path to the token file for a given profile.
func GetTokenFilePath(profile string) (string, error) {
	tokensDir, err := GetTokensDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(tokensDir, profile+".json"), nil
}
