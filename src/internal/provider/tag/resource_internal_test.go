package tag

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/tag/models"
	"github.com/stretchr/testify/assert"
)

// setupTestClient creates a test N8nClient with httptest server.
func setupTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)

	cfg := n8nsdk.NewConfiguration()
	cfg.Servers = n8nsdk.ServerConfigurations{
		{
			URL:         server.URL,
			Description: "Test server",
		},
	}
	cfg.HTTPClient = server.Client()
	cfg.AddDefaultHeader("X-N8N-API-KEY", "test-key")

	apiClient := n8nsdk.NewAPIClient(cfg)
	n8nClient := &client.N8nClient{
		APIClient: apiClient,
	}

	return n8nClient, server
}

// TestTagResource_executeCreateLogic tests the executeCreateLogic method with error cases.
func TestTagResource_executeCreateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		tagName      string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectID     string
	}{
		{
			name:    "successful creation",
			tagName: "Test Tag",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/tags" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					json.NewEncoder(w).Encode(map[string]any{
						"id":        "tag-123",
						"name":      "Test Tag",
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-01T00:00:00Z",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectID:    "tag-123",
		},
		{
			name:    "API error",
			tagName: "Failed Tag",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name:    "empty tag name",
			tagName: "",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message": "Tag name is required"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &TagResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				Name: types.StringValue(tt.tagName),
			}
			resp := &resource.CreateResponse{
				State: resource.CreateResponse{}.State,
			}

			result := r.executeCreateLogic(ctx, plan, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, tt.expectID, plan.ID.ValueString(), "Tag ID should match")
			}
		})
	}
}

// TestTagResource_executeReadLogic tests the executeReadLogic method with error cases.
func TestTagResource_executeReadLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		tagID        string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectName   string
	}{
		{
			name:  "successful read",
			tagID: "tag-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/tags/tag-123" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":        "tag-123",
						"name":      "Retrieved Tag",
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-01T00:00:00Z",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectName:  "Retrieved Tag",
		},
		{
			name:  "tag not found",
			tagID: "tag-404",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Tag not found"}`))
			},
			expectError: true,
		},
		{
			name:  "API error",
			tagID: "tag-500",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &TagResource{client: n8nClient}
			ctx := context.Background()
			state := &models.Resource{
				ID: types.StringValue(tt.tagID),
			}
			resp := &resource.ReadResponse{
				State: resource.ReadResponse{}.State,
			}

			result := r.executeReadLogic(ctx, state, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, tt.expectName, state.Name.ValueString(), "Tag name should match")
			}
		})
	}
}

// TestTagResource_executeUpdateLogic tests the executeUpdateLogic method with error cases.
func TestTagResource_executeUpdateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		tagID        string
		newName      string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:    "successful update",
			tagID:   "tag-123",
			newName: "Updated Tag",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/tags/tag-123" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":        "tag-123",
						"name":      "Updated Tag",
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-02T00:00:00Z",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:    "tag not found",
			tagID:   "tag-404",
			newName: "Updated Tag",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Tag not found"}`))
			},
			expectError: true,
		},
		{
			name:    "API error",
			tagID:   "tag-500",
			newName: "Updated Tag",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &TagResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				ID:   types.StringValue(tt.tagID),
				Name: types.StringValue(tt.newName),
			}
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}

			result := r.executeUpdateLogic(ctx, plan, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			}
		})
	}
}

// TestTagResource_executeDeleteLogic tests the executeDeleteLogic method with error cases.
func TestTagResource_executeDeleteLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		tagID        string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:  "successful deletion",
			tagID: "tag-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete && r.URL.Path == "/tags/tag-123" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:  "tag not found",
			tagID: "tag-404",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Tag not found"}`))
			},
			expectError: true,
		},
		{
			name:  "API error",
			tagID: "tag-500",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &TagResource{client: n8nClient}
			ctx := context.Background()
			state := &models.Resource{
				ID: types.StringValue(tt.tagID),
			}
			resp := &resource.DeleteResponse{
				State: resource.DeleteResponse{}.State,
			}

			result := r.executeDeleteLogic(ctx, state, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			}
		})
	}
}
