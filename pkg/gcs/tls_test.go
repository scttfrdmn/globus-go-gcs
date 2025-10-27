package gcs

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"testing"
	"time"
)

func TestSecureTLSConfig(t *testing.T) {
	cfg := SecureTLSConfig()

	// Verify minimum TLS version
	if cfg.MinVersion != tls.VersionTLS12 {
		t.Errorf("MinVersion = 0x%04x, want TLS 1.2 (0x%04x)", cfg.MinVersion, tls.VersionTLS12)
	}

	// Verify cipher suites are configured
	if len(cfg.CipherSuites) == 0 {
		t.Error("No cipher suites configured")
	}

	// Verify at least one ECDHE-GCM suite (strongest)
	hasECDHEGCM := false
	for _, cipher := range cfg.CipherSuites {
		if cipher == tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 ||
			cipher == tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384 ||
			cipher == tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256 ||
			cipher == tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384 {
			hasECDHEGCM = true
			break
		}
	}
	if !hasECDHEGCM {
		t.Error("No ECDHE-GCM cipher suites configured")
	}

	// Verify curve preferences
	if len(cfg.CurvePreferences) == 0 {
		t.Error("No curve preferences configured")
	}

	// Verify certificate verification is enabled by default
	if cfg.InsecureSkipVerify {
		t.Error("InsecureSkipVerify should be false by default")
	}
}

func TestSecureHTTPClient(t *testing.T) {
	timeout := 10 * time.Second
	client := SecureHTTPClient(timeout)

	// Verify timeout
	if client.Timeout != timeout {
		t.Errorf("Timeout = %v, want %v", client.Timeout, timeout)
	}

	// Verify transport is configured
	if client.Transport == nil {
		t.Fatal("Transport is nil")
	}

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("Transport is not *http.Transport")
	}

	// Verify TLS config exists
	if transport.TLSClientConfig == nil {
		t.Fatal("TLSClientConfig is nil")
	}

	// Verify minimum TLS version
	if transport.TLSClientConfig.MinVersion != tls.VersionTLS12 {
		t.Errorf("MinVersion = 0x%04x, want TLS 1.2", transport.TLSClientConfig.MinVersion)
	}

	// Verify connection pooling
	if transport.MaxIdleConns <= 0 {
		t.Error("MaxIdleConns not configured")
	}
	if transport.MaxIdleConnsPerHost <= 0 {
		t.Error("MaxIdleConnsPerHost not configured")
	}

	// Verify HTTP/2 support
	if !transport.ForceAttemptHTTP2 {
		t.Error("HTTP/2 support not enabled")
	}
}

func TestCustomTLSConfig(t *testing.T) {
	tests := []struct {
		name    string
		opts    []TLSConfigOption
		check   func(*testing.T, *tls.Config)
	}{
		{
			name: "with min TLS 1.3",
			opts: []TLSConfigOption{
				WithTLSMinVersion(tls.VersionTLS13),
			},
			check: func(t *testing.T, cfg *tls.Config) {
				if cfg.MinVersion != tls.VersionTLS13 {
					t.Errorf("MinVersion = 0x%04x, want TLS 1.3", cfg.MinVersion)
				}
			},
		},
		{
			name: "with insecure skip verify",
			opts: []TLSConfigOption{
				WithTLSInsecureSkipVerify(),
			},
			check: func(t *testing.T, cfg *tls.Config) {
				if !cfg.InsecureSkipVerify {
					t.Error("InsecureSkipVerify should be true")
				}
			},
		},
		{
			name: "with server name",
			opts: []TLSConfigOption{
				WithServerName("example.com"),
			},
			check: func(t *testing.T, cfg *tls.Config) {
				if cfg.ServerName != "example.com" {
					t.Errorf("ServerName = %q, want example.com", cfg.ServerName)
				}
			},
		},
		{
			name: "with custom root CAs",
			opts: []TLSConfigOption{
				WithRootCAs(x509.NewCertPool()),
			},
			check: func(t *testing.T, cfg *tls.Config) {
				if cfg.RootCAs == nil {
					t.Error("RootCAs should be set")
				}
			},
		},
		{
			name: "multiple options",
			opts: []TLSConfigOption{
				WithTLSMinVersion(tls.VersionTLS13),
				WithServerName("test.example.com"),
			},
			check: func(t *testing.T, cfg *tls.Config) {
				if cfg.MinVersion != tls.VersionTLS13 {
					t.Error("MinVersion not set correctly")
				}
				if cfg.ServerName != "test.example.com" {
					t.Error("ServerName not set correctly")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := CustomTLSConfig(tt.opts...)
			tt.check(t, cfg)
		})
	}
}

func TestValidateTLSConfig(t *testing.T) {
	tests := []struct {
		name          string
		config        *tls.Config
		allowInsecure bool
		wantErr       bool
		errContains   string
	}{
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "secure config",
			config: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			wantErr: false,
		},
		{
			name: "TLS 1.0 not allowed",
			//nolint:gosec // G402: Intentionally testing rejection of TLS 1.0
			config: &tls.Config{
				MinVersion: tls.VersionTLS10,
			},
			wantErr:     true,
			errContains: "TLS version < 1.2",
		},
		{
			name: "TLS 1.1 not allowed",
			//nolint:gosec // G402: Intentionally testing rejection of TLS 1.1
			config: &tls.Config{
				MinVersion: tls.VersionTLS11,
			},
			wantErr:     true,
			errContains: "TLS version < 1.2",
		},
		{
			name: "insecure skip verify blocked",
			//nolint:gosec // G402: Intentionally testing rejection of InsecureSkipVerify
			config: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				InsecureSkipVerify: true,
			},
			allowInsecure: false,
			wantErr:       true,
			errContains:   "InsecureSkipVerify",
		},
		{
			name: "insecure skip verify allowed",
			//nolint:gosec // G402: Intentionally testing acceptance when explicitly allowed
			config: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				InsecureSkipVerify: true,
			},
			allowInsecure: true,
			wantErr:       false,
		},
		{
			name: "weak cipher suite RC4",
			//nolint:gosec // G402: Intentionally testing rejection of RC4
			config: &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_RSA_WITH_RC4_128_SHA,
				},
			},
			wantErr:     true,
			errContains: "weak cipher suite",
		},
		{
			name: "weak cipher suite 3DES",
			//nolint:gosec // G402: Intentionally testing rejection of 3DES
			config: &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
				},
			},
			wantErr:     true,
			errContains: "weak cipher suite",
		},
		{
			name: "weak cipher suite CBC",
			//nolint:gosec // G402: Intentionally testing rejection of CBC
			config: &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				},
			},
			wantErr:     true,
			errContains: "weak cipher suite",
		},
		{
			name: "strong cipher suites only",
			config: &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTLSConfig(tt.config, tt.allowInsecure)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTLSConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.errContains != "" {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.errContains)
				} else if !contains(err.Error(), tt.errContains) {
					t.Errorf("Error %q does not contain %q", err.Error(), tt.errContains)
				}
			}
		})
	}
}

func TestGetTLSVersion(t *testing.T) {
	tests := []struct {
		version uint16
		want    string
	}{
		{tls.VersionTLS10, "TLS 1.0"},
		{tls.VersionTLS11, "TLS 1.1"},
		{tls.VersionTLS12, "TLS 1.2"},
		{tls.VersionTLS13, "TLS 1.3"},
		{0x9999, "Unknown (0x9999)"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := GetTLSVersion(tt.version)
			if got != tt.want {
				t.Errorf("GetTLSVersion(0x%04x) = %q, want %q", tt.version, got, tt.want)
			}
		})
	}
}

func TestGetCipherSuiteName(t *testing.T) {
	// Test a few known cipher suites
	tests := []struct {
		cipher uint16
		want   string
	}{
		{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"},
		{tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384, "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"},
		{tls.TLS_RSA_WITH_AES_128_GCM_SHA256, "TLS_RSA_WITH_AES_128_GCM_SHA256"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := GetCipherSuiteName(tt.cipher)
			if got != tt.want {
				t.Errorf("GetCipherSuiteName(0x%04x) = %q, want %q", tt.cipher, got, tt.want)
			}
		})
	}
}

func TestSecureTLSConfig_NoCBCCiphers(t *testing.T) {
	cfg := SecureTLSConfig()

	// Verify no CBC cipher suites (vulnerable to BEAST, Lucky 13)
	cbcCiphers := []uint16{
		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		// Note: TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384 not defined in all Go versions
	}

	for _, cbcCipher := range cbcCiphers {
		for _, configuredCipher := range cfg.CipherSuites {
			if configuredCipher == cbcCipher {
				t.Errorf("Insecure CBC cipher suite found: 0x%04x (%s)",
					cbcCipher, GetCipherSuiteName(cbcCipher))
			}
		}
	}
}

func TestSecureTLSConfig_NoRC4(t *testing.T) {
	cfg := SecureTLSConfig()

	// Verify no RC4 cipher suites (broken, RFC 7465)
	rc4Ciphers := []uint16{
		tls.TLS_RSA_WITH_RC4_128_SHA,
		tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
	}

	for _, rc4Cipher := range rc4Ciphers {
		for _, configuredCipher := range cfg.CipherSuites {
			if configuredCipher == rc4Cipher {
				t.Errorf("Insecure RC4 cipher suite found: 0x%04x", rc4Cipher)
			}
		}
	}
}

func TestSecureTLSConfig_No3DES(t *testing.T) {
	cfg := SecureTLSConfig()

	// Verify no 3DES cipher suites (deprecated, 64-bit block size)
	desCiphers := []uint16{
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	}

	for _, desCipher := range desCiphers {
		for _, configuredCipher := range cfg.CipherSuites {
			if configuredCipher == desCipher {
				t.Errorf("Insecure 3DES cipher suite found: 0x%04x", desCipher)
			}
		}
	}
}

func TestSecureTLSConfig_OnlyGCM(t *testing.T) {
	cfg := SecureTLSConfig()

	// Verify all cipher suites use GCM (authenticated encryption)
	for _, cipher := range cfg.CipherSuites {
		name := GetCipherSuiteName(cipher)
		if !contains(name, "GCM") {
			t.Errorf("Non-GCM cipher suite found: 0x%04x (%s)", cipher, name)
		}
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
