package gcs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListRoles(t *testing.T) {
	expectedList := &RoleList{
		Data: []Role{
			{
				ID:         "role-1",
				Collection: "collection-1",
				Principal:  "user@example.org",
				Role:       "administrator",
			},
			{
				ID:         "role-2",
				Collection: "collection-1",
				Principal:  "group@example.org",
				Role:       "access_manager",
			},
		},
		HasNextPage:  false,
		TotalResults: 2,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/roles" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/api/roles")
		}
		if r.Method != http.MethodGet {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodGet)
		}

		// Check query parameters
		query := r.URL.Query()
		if collection := query.Get("collection"); collection != "" {
			if collection != "collection-1" {
				t.Errorf("collection = %q, want %q", collection, "collection-1")
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

	t.Run("list all roles", func(t *testing.T) {
		list, err := client.ListRoles(ctx, nil)
		if err != nil {
			t.Fatalf("ListRoles() error: %v", err)
		}

		if len(list.Data) != len(expectedList.Data) {
			t.Errorf("ListRoles() returned %d roles, want %d", len(list.Data), len(expectedList.Data))
		}

		if list.Data[0].ID != expectedList.Data[0].ID {
			t.Errorf("first role ID = %q, want %q", list.Data[0].ID, expectedList.Data[0].ID)
		}
	})

	t.Run("list with collection filter", func(t *testing.T) {
		opts := &ListRolesOptions{
			Collection: "collection-1",
		}
		list, err := client.ListRoles(ctx, opts)
		if err != nil {
			t.Fatalf("ListRoles() error: %v", err)
		}

		if len(list.Data) != 2 {
			t.Errorf("ListRoles() returned %d roles, want 2", len(list.Data))
		}
	})
}

func TestGetRole(t *testing.T) {
	expectedRole := &Role{
		ID:         "test-role-id",
		Collection: "test-collection",
		Principal:  "user@example.org",
		Role:       "administrator",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/roles/test-role-id"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q, want %q", r.URL.Path, expectedPath)
		}
		if r.Method != http.MethodGet {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodGet)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expectedRole)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()

	t.Run("valid role ID", func(t *testing.T) {
		role, err := client.GetRole(ctx, "test-role-id")
		if err != nil {
			t.Fatalf("GetRole() error: %v", err)
		}

		if role.ID != expectedRole.ID {
			t.Errorf("ID = %q, want %q", role.ID, expectedRole.ID)
		}
		if role.Collection != expectedRole.Collection {
			t.Errorf("Collection = %q, want %q", role.Collection, expectedRole.Collection)
		}
		if role.Principal != expectedRole.Principal {
			t.Errorf("Principal = %q, want %q", role.Principal, expectedRole.Principal)
		}
		if role.Role != expectedRole.Role {
			t.Errorf("Role = %q, want %q", role.Role, expectedRole.Role)
		}
	})

	t.Run("empty role ID", func(t *testing.T) {
		_, err := client.GetRole(ctx, "")
		if err == nil {
			t.Error("GetRole() expected error for empty ID, got nil")
		}
	})
}

func TestCreateRole(t *testing.T) {
	inputRole := &Role{
		Collection: "collection-1",
		Principal:  "user@example.org",
		Role:       "administrator",
	}

	createdRole := &Role{
		ID:         "new-role-id",
		Collection: "collection-1",
		Principal:  "user@example.org",
		Role:       "administrator",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/roles" {
			t.Errorf("request path = %q, want %q", r.URL.Path, "/api/roles")
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
		var received Role
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode request body: %v", err)
		}

		if received.Collection != inputRole.Collection {
			t.Errorf("request Collection = %q, want %q", received.Collection, inputRole.Collection)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(createdRole)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()

	t.Run("create role", func(t *testing.T) {
		result, err := client.CreateRole(ctx, inputRole)
		if err != nil {
			t.Fatalf("CreateRole() error: %v", err)
		}

		if result.ID != createdRole.ID {
			t.Errorf("ID = %q, want %q", result.ID, createdRole.ID)
		}
		if result.Collection != createdRole.Collection {
			t.Errorf("Collection = %q, want %q", result.Collection, createdRole.Collection)
		}
	})

	t.Run("nil role", func(t *testing.T) {
		_, err := client.CreateRole(ctx, nil)
		if err == nil {
			t.Error("CreateRole() expected error for nil role, got nil")
		}
	})
}

func TestUpdateRole(t *testing.T) {
	inputRole := &Role{
		Role: "access_manager",
	}

	updatedRole := &Role{
		ID:         "test-role-id",
		Collection: "collection-1",
		Principal:  "user@example.org",
		Role:       "access_manager",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/roles/test-role-id"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q, want %q", r.URL.Path, expectedPath)
		}
		if r.Method != http.MethodPatch {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodPatch)
		}

		// Decode request body
		var received Role
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode request body: %v", err)
		}

		if received.Role != inputRole.Role {
			t.Errorf("request Role = %q, want %q", received.Role, inputRole.Role)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(updatedRole)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL + "/api/",
		httpClient:  &http.Client{},
		accessToken: "test-token",
		userAgent:   "test-agent",
	}

	ctx := context.Background()

	t.Run("update role", func(t *testing.T) {
		result, err := client.UpdateRole(ctx, "test-role-id", inputRole)
		if err != nil {
			t.Fatalf("UpdateRole() error: %v", err)
		}

		if result.ID != updatedRole.ID {
			t.Errorf("ID = %q, want %q", result.ID, updatedRole.ID)
		}
		if result.Role != updatedRole.Role {
			t.Errorf("Role = %q, want %q", result.Role, updatedRole.Role)
		}
	})

	t.Run("empty role ID", func(t *testing.T) {
		_, err := client.UpdateRole(ctx, "", inputRole)
		if err == nil {
			t.Error("UpdateRole() expected error for empty ID, got nil")
		}
	})

	t.Run("nil role", func(t *testing.T) {
		_, err := client.UpdateRole(ctx, "test-id", nil)
		if err == nil {
			t.Error("UpdateRole() expected error for nil role, got nil")
		}
	})
}

func TestDeleteRole(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/roles/test-role-id"
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

	t.Run("delete role", func(t *testing.T) {
		err := client.DeleteRole(ctx, "test-role-id")
		if err != nil {
			t.Errorf("DeleteRole() error: %v", err)
		}
	})

	t.Run("empty role ID", func(t *testing.T) {
		err := client.DeleteRole(ctx, "")
		if err == nil {
			t.Error("DeleteRole() expected error for empty ID, got nil")
		}
		if !strings.Contains(err.Error(), "required") {
			t.Errorf("DeleteRole() error message should contain 'required', got: %v", err)
		}
	})
}
