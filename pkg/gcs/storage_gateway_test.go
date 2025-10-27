package gcs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListStorageGateways(t *testing.T) {
	expectedList := &StorageGatewayList{
		Data: []StorageGateway{
			{
				ID:            "gateway-1",
				DisplayName:   "POSIX Storage",
				ConnectorID:   "posix",
				ConnectorName: "POSIX",
				Root:          "/data",
			},
			{
				ID:            "gateway-2",
				DisplayName:   "S3 Storage",
				ConnectorID:   "s3",
				ConnectorName: "Amazon S3",
				Root:          "s3://bucket",
			},
		},
		HasNextPage:  false,
		TotalResults: 2,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/storage_gateways" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/api/storage_gateways")
		}
		if r.Method != http.MethodGet {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodGet)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expectedList)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()

	t.Run("list all storage gateways", func(t *testing.T) {
		list, err := client.ListStorageGateways(ctx, nil)
		if err != nil {
			t.Fatalf("ListStorageGateways() error: %v", err)
		}

		if len(list.Data) != len(expectedList.Data) {
			t.Errorf("ListStorageGateways() returned %d gateways, want %d", len(list.Data), len(expectedList.Data))
		}

		if list.Data[0].ID != expectedList.Data[0].ID {
			t.Errorf("first gateway ID = %q, want %q", list.Data[0].ID, expectedList.Data[0].ID)
		}
	})

	t.Run("list with filter", func(t *testing.T) {
		opts := &ListStorageGatewaysOptions{
			Filter: "POSIX",
		}
		list, err := client.ListStorageGateways(ctx, opts)
		if err != nil {
			t.Fatalf("ListStorageGateways() error: %v", err)
		}

		if len(list.Data) != 2 {
			t.Errorf("ListStorageGateways() returned %d gateways, want 2", len(list.Data))
		}
	})
}

func TestGetStorageGateway(t *testing.T) {
	expectedGateway := &StorageGateway{
		ID:            "test-gateway-id",
		DisplayName:   "Test Gateway",
		ConnectorID:   "posix",
		ConnectorName: "POSIX",
		Root:          "/data/test",
		HighAssurance: true,
		RequireMFA:    false,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/storage_gateways/test-gateway-id"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q, want %q", r.URL.Path, expectedPath)
		}
		if r.Method != http.MethodGet {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodGet)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expectedGateway)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()

	t.Run("valid gateway ID", func(t *testing.T) {
		gateway, err := client.GetStorageGateway(ctx, "test-gateway-id")
		if err != nil {
			t.Fatalf("GetStorageGateway() error: %v", err)
		}

		if gateway.ID != expectedGateway.ID {
			t.Errorf("ID = %q, want %q", gateway.ID, expectedGateway.ID)
		}
		if gateway.DisplayName != expectedGateway.DisplayName {
			t.Errorf("DisplayName = %q, want %q", gateway.DisplayName, expectedGateway.DisplayName)
		}
	})

	t.Run("empty gateway ID", func(t *testing.T) {
		_, err := client.GetStorageGateway(ctx, "")
		if err == nil {
			t.Error("GetStorageGateway() expected error for empty ID, got nil")
		}
	})
}

func TestCreateStorageGateway(t *testing.T) {
	inputGateway := &StorageGateway{
		DisplayName: "New Gateway",
		ConnectorID: "posix",
		Root:        "/data/new",
	}

	createdGateway := &StorageGateway{
		ID:          "new-gateway-id",
		DisplayName: "New Gateway",
		ConnectorID: "posix",
		Root:        "/data/new",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/storage_gateways" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/api/storage_gateways")
		}
		if r.Method != http.MethodPost {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodPost)
		}

		// Verify content type
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("Content-Type = %q, want %q", ct, "application/json")
		}

		// Decode request body
		var received StorageGateway
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode request body: %v", err)
		}

		if received.DisplayName != inputGateway.DisplayName {
			t.Errorf("request DisplayName = %q, want %q", received.DisplayName, inputGateway.DisplayName)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(createdGateway)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()

	t.Run("create gateway", func(t *testing.T) {
		result, err := client.CreateStorageGateway(ctx, inputGateway)
		if err != nil {
			t.Fatalf("CreateStorageGateway() error: %v", err)
		}

		if result.ID != createdGateway.ID {
			t.Errorf("ID = %q, want %q", result.ID, createdGateway.ID)
		}
		if result.DisplayName != createdGateway.DisplayName {
			t.Errorf("DisplayName = %q, want %q", result.DisplayName, createdGateway.DisplayName)
		}
	})

	t.Run("nil gateway", func(t *testing.T) {
		_, err := client.CreateStorageGateway(ctx, nil)
		if err == nil {
			t.Error("CreateStorageGateway() expected error for nil gateway, got nil")
		}
	})
}

func TestUpdateStorageGateway(t *testing.T) {
	inputGateway := &StorageGateway{
		DisplayName: "Updated Gateway",
		Root:        "/data/updated",
	}

	updatedGateway := &StorageGateway{
		ID:          "test-gateway-id",
		DisplayName: "Updated Gateway",
		ConnectorID: "posix",
		Root:        "/data/updated",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/storage_gateways/test-gateway-id"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q, want %q", r.URL.Path, expectedPath)
		}
		if r.Method != http.MethodPatch {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodPatch)
		}

		// Decode request body
		var received StorageGateway
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode request body: %v", err)
		}

		if received.DisplayName != inputGateway.DisplayName {
			t.Errorf("request DisplayName = %q, want %q", received.DisplayName, inputGateway.DisplayName)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(updatedGateway)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()

	t.Run("update gateway", func(t *testing.T) {
		result, err := client.UpdateStorageGateway(ctx, "test-gateway-id", inputGateway)
		if err != nil {
			t.Fatalf("UpdateStorageGateway() error: %v", err)
		}

		if result.ID != updatedGateway.ID {
			t.Errorf("ID = %q, want %q", result.ID, updatedGateway.ID)
		}
		if result.DisplayName != updatedGateway.DisplayName {
			t.Errorf("DisplayName = %q, want %q", result.DisplayName, updatedGateway.DisplayName)
		}
	})

	t.Run("empty gateway ID", func(t *testing.T) {
		_, err := client.UpdateStorageGateway(ctx, "", inputGateway)
		if err == nil {
			t.Error("UpdateStorageGateway() expected error for empty ID, got nil")
		}
	})

	t.Run("nil gateway", func(t *testing.T) {
		_, err := client.UpdateStorageGateway(ctx, "test-id", nil)
		if err == nil {
			t.Error("UpdateStorageGateway() expected error for nil gateway, got nil")
		}
	})
}

func TestDeleteStorageGateway(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/storage_gateways/test-gateway-id"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q, want %q", r.URL.Path, expectedPath)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodDelete)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()

	t.Run("delete gateway", func(t *testing.T) {
		err := client.DeleteStorageGateway(ctx, "test-gateway-id")
		if err != nil {
			t.Errorf("DeleteStorageGateway() error: %v", err)
		}
	})

	t.Run("empty gateway ID", func(t *testing.T) {
		err := client.DeleteStorageGateway(ctx, "")
		if err == nil {
			t.Error("DeleteStorageGateway() expected error for empty ID, got nil")
		}
		if !strings.Contains(err.Error(), "required") {
			t.Errorf("DeleteStorageGateway() error message should contain 'required', got: %v", err)
		}
	})
}
