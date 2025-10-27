package gcs

import (
	"crypto/tls"
	"net/http"
	"time"

	globusauth "github.com/scttfrdmn/globus-go-sdk/v3/pkg/services/auth"
)

// ClientOption is a function that configures a GCS Client.
type ClientOption func(*clientOptions)

// clientOptions holds the configuration for a GCS Client.
type clientOptions struct {
	httpClient   *http.Client
	authClient   *globusauth.Client
	accessToken  string
	timeout      time.Duration
	userAgent    string
	tlsConfig    *tls.Config
}

// defaultOptions returns the default client options.
//
// As of v2.0, uses SecureHTTPClient() which enforces TLS 1.2+ with
// NIST-approved cipher suites for HIPAA/PHI compliance (NIST 800-53 SC-8, SC-13).
func defaultOptions() *clientOptions {
	timeout := 30 * time.Second
	return &clientOptions{
		httpClient: SecureHTTPClient(timeout),
		timeout:    timeout,
		userAgent:  "globus-go-gcs/dev",
		tlsConfig:  SecureTLSConfig(),
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(opts *clientOptions) {
		opts.httpClient = client
	}
}

// WithAuthClient sets the Globus Auth client for authentication.
func WithAuthClient(client *globusauth.Client) ClientOption {
	return func(opts *clientOptions) {
		opts.authClient = client
	}
}

// WithAccessToken sets the access token for authentication.
func WithAccessToken(token string) ClientOption {
	return func(opts *clientOptions) {
		opts.accessToken = token
	}
}

// WithTimeout sets the HTTP request timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(opts *clientOptions) {
		opts.timeout = timeout
		if opts.httpClient != nil {
			opts.httpClient.Timeout = timeout
		}
	}
}

// WithUserAgent sets the User-Agent header.
func WithUserAgent(userAgent string) ClientOption {
	return func(opts *clientOptions) {
		opts.userAgent = userAgent
	}
}

// WithTLSConfig sets a custom TLS configuration.
//
// Use this to customize TLS settings beyond the secure defaults.
// For example, to add custom root CAs or adjust minimum TLS version.
//
// Example:
//
//	tlsConfig := gcs.CustomTLSConfig(
//	    gcs.WithMinTLSVersion(tls.VersionTLS13),
//	    gcs.WithRootCAs(customCertPool),
//	)
//	client, _ := gcs.NewClient(endpoint, gcs.WithTLSConfig(tlsConfig))
func WithTLSConfig(config *tls.Config) ClientOption {
	return func(opts *clientOptions) {
		opts.tlsConfig = config
		// Update HTTP client transport with new TLS config
		if opts.httpClient != nil {
			if transport, ok := opts.httpClient.Transport.(*http.Transport); ok {
				transport.TLSClientConfig = config
			} else {
				// Create new transport with TLS config
				opts.httpClient.Transport = &http.Transport{
					TLSClientConfig:     config,
					MaxIdleConns:        100,
					MaxIdleConnsPerHost: 10,
					IdleConnTimeout:     90 * time.Second,
					ForceAttemptHTTP2:   true,
				}
			}
		}
	}
}

// WithInsecureSkipVerify disables TLS certificate verification.
//
// WARNING: This makes connections vulnerable to man-in-the-middle attacks.
// Only use this for testing against local/development servers.
// Never use in production!
func WithInsecureSkipVerify() ClientOption {
	return func(opts *clientOptions) {
		if opts.tlsConfig == nil {
			opts.tlsConfig = SecureTLSConfig()
		}
		opts.tlsConfig.InsecureSkipVerify = true

		// Update HTTP client transport
		if opts.httpClient != nil {
			if transport, ok := opts.httpClient.Transport.(*http.Transport); ok {
				transport.TLSClientConfig = opts.tlsConfig
			}
		}
	}
}

// WithMinTLSVersion sets the minimum TLS version (default is TLS 1.2).
//
// Use this to enforce TLS 1.3 if your environment requires it:
//
//	client, _ := gcs.NewClient(endpoint, gcs.WithMinTLSVersion(tls.VersionTLS13))
func WithMinTLSVersion(version uint16) ClientOption {
	return func(opts *clientOptions) {
		if opts.tlsConfig == nil {
			opts.tlsConfig = SecureTLSConfig()
		}
		opts.tlsConfig.MinVersion = version

		// Update HTTP client transport
		if opts.httpClient != nil {
			if transport, ok := opts.httpClient.Transport.(*http.Transport); ok {
				transport.TLSClientConfig = opts.tlsConfig
			}
		}
	}
}
