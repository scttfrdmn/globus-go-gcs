package gcs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		endpointFQDN string
		opts        []ClientOption
		wantErr     bool
		wantBaseURL string
	}{
		{
			name:        "valid FQDN",
			endpointFQDN: "test.example.org",
			opts:        nil,
			wantErr:     false,
			wantBaseURL: "https://test.example.org/api/",
		},
		{
			name:        "empty FQDN",
			endpointFQDN: "",
			opts:        nil,
			wantErr:     true,
		},
		{
			name:        "with custom options",
			endpointFQDN: "gcs.example.org",
			opts: []ClientOption{
				WithAccessToken("test-token"),
				WithUserAgent("custom-agent"),
			},
			wantErr:     false,
			wantBaseURL: "https://gcs.example.org/api/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.endpointFQDN, tt.opts...)

			if tt.wantErr {
				if err == nil {
					t.Error("NewClient() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewClient() unexpected error: %v", err)
				return
			}

			if client == nil {
				t.Fatal("NewClient() returned nil client")
			}

			if client.baseURL != tt.wantBaseURL {
				t.Errorf("NewClient() baseURL = %q, want %q", client.baseURL, tt.wantBaseURL)
			}
		})
	}
}

func TestClientOptions(t *testing.T) {
	t.Run("WithAccessToken", func(t *testing.T) {
		token := "test-access-token"
		client, err := NewClient("test.example.org", WithAccessToken(token))
		if err != nil {
			t.Fatalf("NewClient() error: %v", err)
		}

		if client.accessToken != token {
			t.Errorf("accessToken = %q, want %q", client.accessToken, token)
		}
	})

	t.Run("WithUserAgent", func(t *testing.T) {
		agent := "custom-user-agent/1.0"
		client, err := NewClient("test.example.org", WithUserAgent(agent))
		if err != nil {
			t.Fatalf("NewClient() error: %v", err)
		}

		if client.userAgent != agent {
			t.Errorf("userAgent = %q, want %q", client.userAgent, agent)
		}
	})

	t.Run("WithHTTPClient", func(t *testing.T) {
		customClient := &http.Client{Timeout: 10 * time.Second}
		client, err := NewClient("test.example.org", WithHTTPClient(customClient))
		if err != nil {
			t.Fatalf("NewClient() error: %v", err)
		}

		if client.httpClient != customClient {
			t.Error("httpClient not set correctly")
		}
	})

	t.Run("WithTimeout", func(t *testing.T) {
		timeout := 45 * time.Second
		client, err := NewClient("test.example.org", WithTimeout(timeout))
		if err != nil {
			t.Fatalf("NewClient() error: %v", err)
		}

		if client.httpClient.Timeout != timeout {
			t.Errorf("httpClient.Timeout = %v, want %v", client.httpClient.Timeout, timeout)
		}
	})
}

func TestSetAccessToken(t *testing.T) {
	client, err := NewClient("test.example.org")
	if err != nil {
		t.Fatalf("NewClient() error: %v", err)
	}

	newToken := "new-token"
	client.SetAccessToken(newToken)

	if client.accessToken != newToken {
		t.Errorf("SetAccessToken() token = %q, want %q", client.accessToken, newToken)
	}
}

func TestDoRequest(t *testing.T) {
	t.Run("authenticated request", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify authorization header
			auth := r.Header.Get("Authorization")
			if auth != "Bearer test-token" {
				t.Errorf("Authorization header = %q, want %q", auth, "Bearer test-token")
			}

			// Verify user agent
			ua := r.Header.Get("User-Agent")
			if ua == "" {
				t.Error("User-Agent header not set")
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		}))
		defer server.Close()

		// Create client pointing to test server
		client := &Client{
			baseURL:     server.URL + "/",
			httpClient:  &http.Client{},
			accessToken: "test-token",
			userAgent:   "test-agent",
		}

		// Make request
		ctx := context.Background()
		resp, err := client.doRequest(ctx, http.MethodGet, "test", nil)
		if err != nil {
			t.Errorf("doRequest() error: %v", err)
		}
		if resp != nil {
			_ = resp.Body.Close()
		}
	})

	t.Run("error response", func(t *testing.T) {
		// Create test server that returns error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"not found"}`))
		}))
		defer server.Close()

		client := &Client{
			baseURL:    server.URL + "/",
			httpClient: &http.Client{},
			userAgent:  "test-agent",
		}

		ctx := context.Background()
		_, err := client.doRequest(ctx, http.MethodGet, "test", nil)
		if err == nil {
			t.Error("doRequest() expected error for 404, got nil")
		}
	})
}
