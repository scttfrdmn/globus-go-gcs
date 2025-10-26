package gcs

import (
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
}

// defaultOptions returns the default client options.
func defaultOptions() *clientOptions {
	return &clientOptions{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		timeout:   30 * time.Second,
		userAgent: "globus-go-gcs/dev",
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
