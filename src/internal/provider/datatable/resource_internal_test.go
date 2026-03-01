// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package datatable

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/datatable/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
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

// TestDataTableResource_executeCreateLogic tests the executeCreateLogic method with error cases.
func TestDataTableResource_executeCreateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		tableName    string
		columns      []models.Column
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectID     string
	}{
		{
			name:      "successful creation",
			tableName: "My Table",
			columns: []models.Column{
				{Name: types.StringValue("col1"), Type: types.StringValue("string")},
				{Name: types.StringValue("col2"), Type: types.StringValue("number")},
			},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/data-tables" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					json.NewEncoder(w).Encode(map[string]any{
						"id":   "dt-123",
						"name": "My Table",
						"columns": []any{
							map[string]any{"id": "c1", "name": "col1", "type": "string"},
							map[string]any{"id": "c2", "name": "col2", "type": "number"},
						},
						"projectId": "proj-1",
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-01T00:00:00Z",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectID:    "dt-123",
		},
		{
			name:      "API error",
			tableName: "Failed Table",
			columns:   []models.Column{{Name: types.StringValue("a"), Type: types.StringValue("string")}},
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

			r := &DataTableResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				Name:    types.StringValue(tt.tableName),
				Columns: tt.columns,
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
				assert.Equal(t, tt.expectID, plan.ID.ValueString(), "Data table ID should match")
			}
		})
	}
}

// TestDataTableResource_executeReadLogic tests the executeReadLogic method with error cases.
func TestDataTableResource_executeReadLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		tableID      string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectName   string
	}{
		{
			name:    "successful read",
			tableID: "dt-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/data-tables/dt-123" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":   "dt-123",
						"name": "Retrieved Table",
						"columns": []any{
							map[string]any{"id": "c1", "name": "col1", "type": "string"},
						},
						"projectId": "proj-1",
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-01T00:00:00Z",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectName:  "Retrieved Table",
		},
		{
			name:    "table not found",
			tableID: "dt-404",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Data table not found"}`))
			},
			expectError: true,
		},
		{
			name:    "API error",
			tableID: "dt-500",
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

			r := &DataTableResource{client: n8nClient}
			ctx := context.Background()
			state := &models.Resource{
				ID: types.StringValue(tt.tableID),
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
				assert.Equal(t, tt.expectName, state.Name.ValueString(), "Data table name should match")
			}
		})
	}
}

// TestDataTableResource_executeUpdateLogic tests the executeUpdateLogic method with error cases.
func TestDataTableResource_executeUpdateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		tableID      string
		newName      string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:    "successful update",
			tableID: "dt-123",
			newName: "Updated Table",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPatch && r.URL.Path == "/data-tables/dt-123" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":   "dt-123",
						"name": "Updated Table",
						"columns": []any{
							map[string]any{"id": "c1", "name": "col1", "type": "string"},
						},
						"projectId": "proj-1",
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
			name:    "table not found",
			tableID: "dt-404",
			newName: "Updated Table",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Data table not found"}`))
			},
			expectError: true,
		},
		{
			name:    "API error",
			tableID: "dt-500",
			newName: "Updated Table",
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

			r := &DataTableResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				Name: types.StringValue(tt.newName),
			}
			state := &models.Resource{
				ID:   types.StringValue(tt.tableID),
				Name: types.StringValue("Old Name"),
			}
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}

			result := r.executeUpdateLogic(ctx, plan, state, resp)

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

// TestDataTableResource_executeDeleteLogic tests the executeDeleteLogic method with error cases.
func TestDataTableResource_executeDeleteLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		tableID      string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:    "successful deletion",
			tableID: "dt-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete && r.URL.Path == "/data-tables/dt-123" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:    "table not found",
			tableID: "dt-404",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Data table not found"}`))
			},
			expectError: true,
		},
		{
			name:    "API error",
			tableID: "dt-500",
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

			r := &DataTableResource{client: n8nClient}
			ctx := context.Background()
			state := &models.Resource{
				ID: types.StringValue(tt.tableID),
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
