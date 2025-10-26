package gcs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetInfo(t *testing.T) {
	expectedInfo := &Info{
		APIVersion:     "1.0",
		EndpointID:     "test-endpoint-id",
		ManagerVersion: "5.4.0",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/info" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/api/info")
		}
		if r.Method != http.MethodGet {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodGet)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expectedInfo)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL + "/api/",
		httpClient: &http.Client{},
		userAgent:  "test-agent",
	}

	ctx := context.Background()
	info, err := client.GetInfo(ctx)
	if err != nil {
		t.Fatalf("GetInfo() error: %v", err)
	}

	if info.APIVersion != expectedInfo.APIVersion {
		t.Errorf("APIVersion = %q, want %q", info.APIVersion, expectedInfo.APIVersion)
	}
	if info.EndpointID != expectedInfo.EndpointID {
		t.Errorf("EndpointID = %q, want %q", info.EndpointID, expectedInfo.EndpointID)
	}
	if info.ManagerVersion != expectedInfo.ManagerVersion {
		t.Errorf("ManagerVersion = %q, want %q", info.ManagerVersion, expectedInfo.ManagerVersion)
	}
}

func TestGetEndpoint(t *testing.T) {
	expectedEndpoint := &Endpoint{
		ID:           "test-endpoint-id",
		DisplayName:  "Test Endpoint",
		Organization: "Test Org",
		ContactEmail: "test@example.org",
		Public:       true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/endpoint" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/api/endpoint")
		}
		if r.Method != http.MethodGet {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodGet)
		}

		// Verify authorization
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer test-token")
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expectedEndpoint)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()
	endpoint, err := client.GetEndpoint(ctx)
	if err != nil {
		t.Fatalf("GetEndpoint() error: %v", err)
	}

	if endpoint.ID != expectedEndpoint.ID {
		t.Errorf("ID = %q, want %q", endpoint.ID, expectedEndpoint.ID)
	}
	if endpoint.DisplayName != expectedEndpoint.DisplayName {
		t.Errorf("DisplayName = %q, want %q", endpoint.DisplayName, expectedEndpoint.DisplayName)
	}
	if endpoint.Organization != expectedEndpoint.Organization {
		t.Errorf("Organization = %q, want %q", endpoint.Organization, expectedEndpoint.Organization)
	}
}

func TestUpdateEndpoint(t *testing.T) {
	inputEndpoint := &Endpoint{
		DisplayName:  "Updated Endpoint",
		Organization: "Updated Org",
		ContactEmail: "updated@example.org",
	}

	updatedEndpoint := &Endpoint{
		ID:           "test-endpoint-id",
		DisplayName:  "Updated Endpoint",
		Organization: "Updated Org",
		ContactEmail: "updated@example.org",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/endpoint" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/api/endpoint")
		}
		if r.Method != http.MethodPatch {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodPatch)
		}

		// Verify content type
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("Content-Type = %q, want %q", ct, "application/json")
		}

		// Decode request body
		var received Endpoint
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode request body: %v", err)
		}

		// Verify request data
		if received.DisplayName != inputEndpoint.DisplayName {
			t.Errorf("request DisplayName = %q, want %q", received.DisplayName, inputEndpoint.DisplayName)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(updatedEndpoint)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()
	result, err := client.UpdateEndpoint(ctx, inputEndpoint)
	if err != nil {
		t.Fatalf("UpdateEndpoint() error: %v", err)
	}

	if result.ID != updatedEndpoint.ID {
		t.Errorf("ID = %q, want %q", result.ID, updatedEndpoint.ID)
	}
	if result.DisplayName != updatedEndpoint.DisplayName {
		t.Errorf("DisplayName = %q, want %q", result.DisplayName, updatedEndpoint.DisplayName)
	}
}
