package gcs

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"
)

// SecureTLSConfig returns a TLS configuration that enforces:
//   - TLS 1.2+ (NIST 800-53 SC-8, SC-13)
//   - NIST-approved cipher suites
//   - Strong key exchange
//
// This configuration is required for HIPAA/PHI compliance and meets
// NIST Special Publication 800-52 Rev. 2 requirements.
func SecureTLSConfig() *tls.Config {
	return &tls.Config{
		// Enforce TLS 1.2 minimum (NIST 800-52 Rev. 2)
		MinVersion: tls.VersionTLS12,

		// Prefer server cipher suite ordering
		PreferServerCipherSuites: true,

		// NIST-approved cipher suites (FIPS 140-2 compatible)
		// Prioritized by security strength and performance
		CipherSuites: []uint16{
			// TLS 1.3 cipher suites (strongest, hardware-accelerated)
			// Note: TLS 1.3 suites are automatically enabled when using TLS 1.3
			// and cannot be explicitly configured in Go's TLS 1.3 implementation

			// TLS 1.2 cipher suites (NIST-approved, FIPS 140-2)
			// ECDHE with AES-GCM (strongest for TLS 1.2)
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,

			// RSA with AES-GCM (fallback, still secure)
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256, //nolint:gosec // G402: RSA-GCM is secure and NIST-approved
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		},

		// Curve preferences (NIST-approved)
		CurvePreferences: []tls.CurveID{
			tls.X25519,    // Modern, fast, secure
			tls.CurveP256, // NIST P-256 (required for FIPS)
			tls.CurveP384, // NIST P-384 (higher security)
		},

		// Enable session resumption for performance
		// (reduces handshake overhead)
		SessionTicketsDisabled: false,

		// Verify server certificates by default
		InsecureSkipVerify: false,
	}
}

// SecureHTTPClient creates an HTTP client with secure TLS configuration.
//
// This is the recommended way to create HTTP clients for production use,
// as it enforces TLS 1.2+ and NIST-approved cipher suites.
func SecureHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig:     SecureTLSConfig(),
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			// Enable HTTP/2 support
			ForceAttemptHTTP2: true,
		},
	}
}

// TLSConfigOption is a function that modifies a TLS configuration.
type TLSConfigOption func(*tls.Config)

// WithTLSMinVersion sets the minimum TLS version for CustomTLSConfig.
// Use this to enforce TLS 1.3 if needed (default is TLS 1.2).
func WithTLSMinVersion(version uint16) TLSConfigOption {
	return func(cfg *tls.Config) {
		cfg.MinVersion = version
	}
}

// WithTLSInsecureSkipVerify disables certificate verification for CustomTLSConfig.
//
// WARNING: Only use this for testing! This makes connections vulnerable
// to man-in-the-middle attacks. Never use in production.
func WithTLSInsecureSkipVerify() TLSConfigOption {
	return func(cfg *tls.Config) {
		cfg.InsecureSkipVerify = true
	}
}

// WithRootCAs sets custom root CA certificates for verification.
// Use this when connecting to servers with self-signed certificates
// or internal CAs.
func WithRootCAs(certPool *x509.CertPool) TLSConfigOption {
	return func(cfg *tls.Config) {
		cfg.RootCAs = certPool
	}
}

// WithServerName sets the server name for SNI (Server Name Indication).
// Use this when the server hostname doesn't match the certificate.
func WithServerName(serverName string) TLSConfigOption {
	return func(cfg *tls.Config) {
		cfg.ServerName = serverName
	}
}

// CustomTLSConfig creates a TLS configuration with custom options.
// Starts with secure defaults and applies the provided options.
func CustomTLSConfig(opts ...TLSConfigOption) *tls.Config {
	cfg := SecureTLSConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// ValidateTLSConfig validates that a TLS configuration meets security requirements.
//
// Returns an error if the configuration is insecure:
//   - TLS version < 1.2
//   - InsecureSkipVerify enabled (unless explicitly allowed)
//   - Weak cipher suites
func ValidateTLSConfig(cfg *tls.Config, allowInsecure bool) error {
	if cfg == nil {
		return fmt.Errorf("TLS config is nil")
	}

	// Check minimum TLS version
	if cfg.MinVersion != 0 && cfg.MinVersion < tls.VersionTLS12 {
		return fmt.Errorf("TLS version < 1.2 is not allowed (found: 0x%04x)", cfg.MinVersion)
	}

	// Check for insecure skip verify
	if cfg.InsecureSkipVerify && !allowInsecure {
		return fmt.Errorf("InsecureSkipVerify is enabled (certificate verification disabled)")
	}

	// Check for weak cipher suites (if any are explicitly configured)
	if len(cfg.CipherSuites) > 0 {
		weakCiphers := []uint16{
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		}

		for _, cipher := range cfg.CipherSuites {
			for _, weak := range weakCiphers {
				if cipher == weak {
					return fmt.Errorf("weak cipher suite detected: 0x%04x", cipher)
				}
			}
		}
	}

	return nil
}

// GetTLSVersion returns a human-readable string for a TLS version constant.
func GetTLSVersion(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("Unknown (0x%04x)", version)
	}
}

// GetCipherSuiteName returns a human-readable string for a cipher suite constant.
func GetCipherSuiteName(cipher uint16) string {
	// This returns the cipher suite name for logging/debugging
	return tls.CipherSuiteName(cipher)
}
