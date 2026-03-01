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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/datatable/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

func setupDatatableDSTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

func TestDataTableDataSource_validateIdentifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		data        *models.DataSource
		expectValid bool
	}{
		{
			name: "valid with ID",
			data: &models.DataSource{
				ID: types.StringValue("dt-123"),
			},
			expectValid: true,
		},
		{
			name: "valid with name",
			data: &models.DataSource{
				Name: types.StringValue("My Table"),
			},
			expectValid: true,
		},
		{
			name: "valid with both ID and name",
			data: &models.DataSource{
				ID:   types.StringValue("dt-123"),
				Name: types.StringValue("My Table"),
			},
			expectValid: true,
		},
		{
			name: "invalid with both null",
			data: &models.DataSource{
				ID:   types.StringNull(),
				Name: types.StringNull(),
			},
			expectValid: false,
		},
		{
			name:        "invalid with uninitialized values",
			data:        &models.DataSource{},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := &DataTableDataSource{}
			resp := &datasource.ReadResponse{}

			valid := d.validateIdentifier(tt.data, resp)

			if tt.expectValid {
				assert.True(t, valid)
				assert.False(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, valid)
				assert.True(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestDataTableDataSource_fetchDataTableByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		tableID     string
		expectNil   bool
		expectError bool
	}{
		{
			name:    "data table found by ID",
			tableID: "dt-123",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/data-tables/dt-123" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"id":        "dt-123",
						"name":      "Test Table",
						"columns":   []interface{}{},
						"projectId": "proj-1",
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-01T00:00:00Z",
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   false,
			expectError: false,
		},
		{
			name:    "data table not found",
			tableID: "dt-999",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/data-tables/dt-999" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "Data table not found"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:    "API returns error",
			tableID: "dt-123",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/data-tables/dt-123" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupDatatableDSTestClient(t, tt.handler)
			defer server.Close()

			d := &DataTableDataSource{client: n8nClient}
			data := &models.DataSource{}
			data.ID = types.StringValue(tt.tableID)
			resp := &datasource.ReadResponse{}

			table := d.fetchDataTableByID(context.Background(), data, resp)

			if tt.expectNil {
				assert.Nil(t, table)
			} else {
				assert.NotNil(t, table)
			}

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestDataTableDataSource_fetchDataTableByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		tableName   string
		expectNil   bool
		expectError bool
	}{
		{
			name:      "data table found by name",
			tableName: "My Table",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/data-tables" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":        "dt-123",
								"name":      "My Table",
								"columns":   []interface{}{},
								"projectId": "proj-1",
								"createdAt": "2024-01-01T00:00:00Z",
								"updatedAt": "2024-01-01T00:00:00Z",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   false,
			expectError: false,
		},
		{
			name:      "data table not found in list",
			tableName: "Missing Table",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/data-tables" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":        "dt-456",
								"name":      "Other Table",
								"columns":   []interface{}{},
								"projectId": "proj-1",
								"createdAt": "2024-01-01T00:00:00Z",
								"updatedAt": "2024-01-01T00:00:00Z",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:      "API returns error",
			tableName: "My Table",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/data-tables" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:      "empty data table list",
			tableName: "My Table",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/data-tables" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupDatatableDSTestClient(t, tt.handler)
			defer server.Close()

			d := &DataTableDataSource{client: n8nClient}
			data := &models.DataSource{}
			data.Name = types.StringValue(tt.tableName)
			resp := &datasource.ReadResponse{}

			table := d.fetchDataTableByName(context.Background(), data, resp)

			if tt.expectNil {
				assert.Nil(t, table)
			} else {
				assert.NotNil(t, table)
			}

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}
