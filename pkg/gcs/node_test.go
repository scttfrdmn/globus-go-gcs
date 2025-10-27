package gcs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListNodes(t *testing.T) {
	expectedList := &NodeList{
		Data: []Node{
			{
				ID:       "node-1",
				Name:     "Primary Node",
				Incoming: true,
				Outgoing: true,
			},
			{
				ID:       "node-2",
				Name:     "Secondary Node",
				Incoming: true,
				Outgoing: false,
			},
		},
		HasNextPage:  false,
		TotalResults: 2,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/nodes" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/api/nodes")
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

	t.Run("list all nodes", func(t *testing.T) {
		list, err := client.ListNodes(ctx, nil)
		if err != nil {
			t.Fatalf("ListNodes() error: %v", err)
		}

		if len(list.Data) != len(expectedList.Data) {
			t.Errorf("ListNodes() returned %d nodes, want %d", len(list.Data), len(expectedList.Data))
		}

		if list.Data[0].ID != expectedList.Data[0].ID {
			t.Errorf("first node ID = %q, want %q", list.Data[0].ID, expectedList.Data[0].ID)
		}
	})

	t.Run("list with filter", func(t *testing.T) {
		opts := &ListNodesOptions{
			Filter: "Primary",
		}
		list, err := client.ListNodes(ctx, opts)
		if err != nil {
			t.Fatalf("ListNodes() error: %v", err)
		}

		if len(list.Data) != 2 {
			t.Errorf("ListNodes() returned %d nodes, want 2", len(list.Data))
		}
	})
}

func TestGetNode(t *testing.T) {
	expectedNode := &Node{
		ID:       "test-node-id",
		Name:     "Test Node",
		Incoming: true,
		Outgoing: true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/nodes/test-node-id"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q, want %q", r.URL.Path, expectedPath)
		}
		if r.Method != http.MethodGet {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodGet)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expectedNode)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()

	t.Run("valid node ID", func(t *testing.T) {
		node, err := client.GetNode(ctx, "test-node-id")
		if err != nil {
			t.Fatalf("GetNode() error: %v", err)
		}

		if node.ID != expectedNode.ID {
			t.Errorf("ID = %q, want %q", node.ID, expectedNode.ID)
		}
		if node.Name != expectedNode.Name {
			t.Errorf("Name = %q, want %q", node.Name, expectedNode.Name)
		}
		if node.Incoming != expectedNode.Incoming {
			t.Errorf("Incoming = %t, want %t", node.Incoming, expectedNode.Incoming)
		}
		if node.Outgoing != expectedNode.Outgoing {
			t.Errorf("Outgoing = %t, want %t", node.Outgoing, expectedNode.Outgoing)
		}
	})

	t.Run("empty node ID", func(t *testing.T) {
		_, err := client.GetNode(ctx, "")
		if err == nil {
			t.Error("GetNode() expected error for empty ID, got nil")
		}
	})
}

func TestCreateNode(t *testing.T) {
	inputNode := &Node{
		Name:     "New Node",
		Incoming: true,
		Outgoing: false,
	}

	createdNode := &Node{
		ID:       "new-node-id",
		Name:     "New Node",
		Incoming: true,
		Outgoing: false,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/nodes" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/api/nodes")
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
		var received Node
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode request body: %v", err)
		}

		if received.Name != inputNode.Name {
			t.Errorf("request Name = %q, want %q", received.Name, inputNode.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(createdNode)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()

	t.Run("create node", func(t *testing.T) {
		result, err := client.CreateNode(ctx, inputNode)
		if err != nil {
			t.Fatalf("CreateNode() error: %v", err)
		}

		if result.ID != createdNode.ID {
			t.Errorf("ID = %q, want %q", result.ID, createdNode.ID)
		}
		if result.Name != createdNode.Name {
			t.Errorf("Name = %q, want %q", result.Name, createdNode.Name)
		}
	})

	t.Run("nil node", func(t *testing.T) {
		_, err := client.CreateNode(ctx, nil)
		if err == nil {
			t.Error("CreateNode() expected error for nil node, got nil")
		}
	})
}

func TestUpdateNode(t *testing.T) {
	inputNode := &Node{
		Name:     "Updated Node",
		Incoming: true,
		Outgoing: true,
	}

	updatedNode := &Node{
		ID:       "test-node-id",
		Name:     "Updated Node",
		Incoming: true,
		Outgoing: true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/nodes/test-node-id"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q, want %q", r.URL.Path, expectedPath)
		}
		if r.Method != http.MethodPatch {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodPatch)
		}

		// Decode request body
		var received Node
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode request body: %v", err)
		}

		if received.Name != inputNode.Name {
			t.Errorf("request Name = %q, want %q", received.Name, inputNode.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(updatedNode)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()

	t.Run("update node", func(t *testing.T) {
		result, err := client.UpdateNode(ctx, "test-node-id", inputNode)
		if err != nil {
			t.Fatalf("UpdateNode() error: %v", err)
		}

		if result.ID != updatedNode.ID {
			t.Errorf("ID = %q, want %q", result.ID, updatedNode.ID)
		}
		if result.Name != updatedNode.Name {
			t.Errorf("Name = %q, want %q", result.Name, updatedNode.Name)
		}
	})

	t.Run("empty node ID", func(t *testing.T) {
		_, err := client.UpdateNode(ctx, "", inputNode)
		if err == nil {
			t.Error("UpdateNode() expected error for empty ID, got nil")
		}
	})

	t.Run("nil node", func(t *testing.T) {
		_, err := client.UpdateNode(ctx, "test-id", nil)
		if err == nil {
			t.Error("UpdateNode() expected error for nil node, got nil")
		}
	})
}

func TestDeleteNode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/nodes/test-node-id"
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

	t.Run("delete node", func(t *testing.T) {
		err := client.DeleteNode(ctx, "test-node-id")
		if err != nil {
			t.Errorf("DeleteNode() error: %v", err)
		}
	})

	t.Run("empty node ID", func(t *testing.T) {
		err := client.DeleteNode(ctx, "")
		if err == nil {
			t.Error("DeleteNode() expected error for empty ID, got nil")
		}
		if !strings.Contains(err.Error(), "required") {
			t.Errorf("DeleteNode() error message should contain 'required', got: %v", err)
		}
	})
}
