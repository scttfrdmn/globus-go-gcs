package gcs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListCollections(t *testing.T) {
	expectedList := &CollectionList{
		Data: []Collection{
			{
				ID:               "collection-1",
				DisplayName:      "Test Collection 1",
				CollectionType:   "mapped",
				StorageGatewayID: "gateway-1",
			},
			{
				ID:               "collection-2",
				DisplayName:      "Test Collection 2",
				CollectionType:   "guest",
				StorageGatewayID: "gateway-1",
			},
		},
		HasNextPage:  false,
		TotalResults: 2,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/collections" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/api/collections")
		}
		if r.Method != http.MethodGet {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodGet)
		}

		// Check query parameters
		query := r.URL.Query()
		if filter := query.Get("filter"); filter != "" {
			if filter != "test" {
				t.Errorf("filter = %q, want %q", filter, "test")
			}
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

	t.Run("list all collections", func(t *testing.T) {
		list, err := client.ListCollections(ctx, nil)
		if err != nil {
			t.Fatalf("ListCollections() error: %v", err)
		}

		if len(list.Data) != len(expectedList.Data) {
			t.Errorf("ListCollections() returned %d collections, want %d", len(list.Data), len(expectedList.Data))
		}

		if list.Data[0].ID != expectedList.Data[0].ID {
			t.Errorf("first collection ID = %q, want %q", list.Data[0].ID, expectedList.Data[0].ID)
		}
	})

	t.Run("list with filter", func(t *testing.T) {
		opts := &ListCollectionsOptions{
			Filter: "test",
		}
		list, err := client.ListCollections(ctx, opts)
		if err != nil {
			t.Fatalf("ListCollections() error: %v", err)
		}

		if len(list.Data) != 2 {
			t.Errorf("ListCollections() returned %d collections, want 2", len(list.Data))
		}
	})
}

func TestGetCollection(t *testing.T) {
	expectedCollection := &Collection{
		ID:               "test-collection-id",
		DisplayName:      "Test Collection",
		CollectionType:   "mapped",
		StorageGatewayID: "gateway-1",
		Public:           true,
		Organization:     "Test Org",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/collections/test-collection-id"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q, want %q", r.URL.Path, expectedPath)
		}
		if r.Method != http.MethodGet {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodGet)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expectedCollection)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()

	t.Run("valid collection ID", func(t *testing.T) {
		collection, err := client.GetCollection(ctx, "test-collection-id")
		if err != nil {
			t.Fatalf("GetCollection() error: %v", err)
		}

		if collection.ID != expectedCollection.ID {
			t.Errorf("ID = %q, want %q", collection.ID, expectedCollection.ID)
		}
		if collection.DisplayName != expectedCollection.DisplayName {
			t.Errorf("DisplayName = %q, want %q", collection.DisplayName, expectedCollection.DisplayName)
		}
	})

	t.Run("empty collection ID", func(t *testing.T) {
		_, err := client.GetCollection(ctx, "")
		if err == nil {
			t.Error("GetCollection() expected error for empty ID, got nil")
		}
	})
}

func TestCreateCollection(t *testing.T) {
	inputCollection := &Collection{
		DisplayName:      "New Collection",
		CollectionType:   "mapped",
		StorageGatewayID: "gateway-1",
		Public:           false,
	}

	createdCollection := &Collection{
		ID:               "new-collection-id",
		DisplayName:      "New Collection",
		CollectionType:   "mapped",
		StorageGatewayID: "gateway-1",
		Public:           false,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/collections" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/api/collections")
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
		var received Collection
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode request body: %v", err)
		}

		if received.DisplayName != inputCollection.DisplayName {
			t.Errorf("request DisplayName = %q, want %q", received.DisplayName, inputCollection.DisplayName)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(createdCollection)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()

	t.Run("create collection", func(t *testing.T) {
		result, err := client.CreateCollection(ctx, inputCollection)
		if err != nil {
			t.Fatalf("CreateCollection() error: %v", err)
		}

		if result.ID != createdCollection.ID {
			t.Errorf("ID = %q, want %q", result.ID, createdCollection.ID)
		}
		if result.DisplayName != createdCollection.DisplayName {
			t.Errorf("DisplayName = %q, want %q", result.DisplayName, createdCollection.DisplayName)
		}
	})

	t.Run("nil collection", func(t *testing.T) {
		_, err := client.CreateCollection(ctx, nil)
		if err == nil {
			t.Error("CreateCollection() expected error for nil collection, got nil")
		}
	})
}

func TestUpdateCollection(t *testing.T) {
	inputCollection := &Collection{
		DisplayName: "Updated Collection",
		Public:      true,
	}

	updatedCollection := &Collection{
		ID:               "test-collection-id",
		DisplayName:      "Updated Collection",
		CollectionType:   "mapped",
		StorageGatewayID: "gateway-1",
		Public:           true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/collections/test-collection-id"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q, want %q", r.URL.Path, expectedPath)
		}
		if r.Method != http.MethodPatch {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodPatch)
		}

		// Decode request body
		var received Collection
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode request body: %v", err)
		}

		if received.DisplayName != inputCollection.DisplayName {
			t.Errorf("request DisplayName = %q, want %q", received.DisplayName, inputCollection.DisplayName)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(updatedCollection)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()

	t.Run("update collection", func(t *testing.T) {
		result, err := client.UpdateCollection(ctx, "test-collection-id", inputCollection)
		if err != nil {
			t.Fatalf("UpdateCollection() error: %v", err)
		}

		if result.ID != updatedCollection.ID {
			t.Errorf("ID = %q, want %q", result.ID, updatedCollection.ID)
		}
		if result.DisplayName != updatedCollection.DisplayName {
			t.Errorf("DisplayName = %q, want %q", result.DisplayName, updatedCollection.DisplayName)
		}
	})

	t.Run("empty collection ID", func(t *testing.T) {
		_, err := client.UpdateCollection(ctx, "", inputCollection)
		if err == nil {
			t.Error("UpdateCollection() expected error for empty ID, got nil")
		}
	})

	t.Run("nil collection", func(t *testing.T) {
		_, err := client.UpdateCollection(ctx, "test-id", nil)
		if err == nil {
			t.Error("UpdateCollection() expected error for nil collection, got nil")
		}
	})
}

func TestDeleteCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/collections/test-collection-id"
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

	t.Run("delete collection", func(t *testing.T) {
		err := client.DeleteCollection(ctx, "test-collection-id")
		if err != nil {
			t.Errorf("DeleteCollection() error: %v", err)
		}
	})

	t.Run("empty collection ID", func(t *testing.T) {
		err := client.DeleteCollection(ctx, "")
		if err == nil {
			t.Error("DeleteCollection() expected error for empty ID, got nil")
		}
		if !strings.Contains(err.Error(), "required") {
			t.Errorf("DeleteCollection() error message should contain 'required', got: %v", err)
		}
	})
}
