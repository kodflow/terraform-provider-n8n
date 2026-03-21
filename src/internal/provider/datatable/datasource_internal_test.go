package datatable

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/datatable/models"
	"github.com/stretchr/testify/assert"
)

// TestDataTableDataSource_executeReadByID tests the executeReadByID method.
func TestDataTableDataSource_executeReadByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		tableID      string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectName   string
	}{
		{
			name:    "successful read by ID",
			tableID: "dt-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/data-tables/dt-123" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":        "dt-123",
						"name":      "My Table",
						"columns":   []any{},
						"projectId": "proj-1",
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-01T00:00:00Z",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectName:  "My Table",
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClient(t, http.HandlerFunc(tt.setupHandler))
			defer server.Close()

			d := &DataTableDataSource{client: n8nClient}
			ctx := t.Context()
			data := &models.DataSource{
				ID: types.StringValue(tt.tableID),
			}
			resp := &datasource.ReadResponse{
				State: datasource.ReadResponse{}.State,
			}

			result := d.executeReadByID(ctx, data, resp)

			//: Verify the result matches the expected error state.
			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, tt.expectName, data.Name.ValueString(), "Table name should match")
			}
		})
	}
}

// TestDataTableDataSource_executeReadByName tests the executeReadByName method.
func TestDataTableDataSource_executeReadByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		tableName    string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectID     string
	}{
		{
			name:      "successful read by name",
			tableName: "My Table",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/data-tables" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []any{
							map[string]any{
								"id":        "dt-123",
								"name":      "My Table",
								"columns":   []any{},
								"projectId": "proj-1",
								"createdAt": "2024-01-01T00:00:00Z",
								"updatedAt": "2024-01-01T00:00:00Z",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectID:    "dt-123",
		},
		{
			name:      "table not found by name",
			tableName: "Missing Table",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/data-tables" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name:      "API error",
			tableName: "Any Table",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClient(t, http.HandlerFunc(tt.setupHandler))
			defer server.Close()

			d := &DataTableDataSource{client: n8nClient}
			ctx := t.Context()
			data := &models.DataSource{
				Name: types.StringValue(tt.tableName),
			}
			resp := &datasource.ReadResponse{
				State: datasource.ReadResponse{}.State,
			}

			result := d.executeReadByName(ctx, data, resp)

			//: Verify the result matches the expected error state.
			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, tt.expectID, data.ID.ValueString(), "Table ID should match")
			}
		})
	}
}

// TestDataTableColumnSchema tests the dataTableColumnSchema helper.
func TestDataTableColumnSchema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		expectAttrs []string
	}{
		{
			name:        "returns all column attributes",
			expectAttrs: []string{"id", "index", "name", "type"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			nestedAttr := dataTableColumnSchema()

			//: Verify the nested attribute contains all expected column attributes.
			assert.NotNil(t, nestedAttr)
			for _, attr := range tt.expectAttrs {
				assert.Contains(t, nestedAttr.Attributes, attr)
			}
		})
	}
}

// TestResolveNamedDataTable tests the resolveNamedDataTable helper.
func TestResolveNamedDataTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		dtList      *n8nsdk.DataTableList
		searchName  string
		expectFound bool
	}{
		{
			name:       "table found by name",
			searchName: "My Table",
			dtList: &n8nsdk.DataTableList{
				Data: []n8nsdk.DataTable{
					{ID: "dt-1", Name: "My Table"},
				},
			},
			expectFound: true,
		},
		{
			name:       "table not found in list",
			searchName: "Missing",
			dtList: &n8nsdk.DataTableList{
				Data: []n8nsdk.DataTable{},
			},
			expectFound: false,
		},
		{
			name:        "nil list",
			searchName:  "Any",
			dtList:      nil,
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data := &models.DataSource{}
			data.Name = types.StringValue(tt.searchName)

			resp := &datasource.ReadResponse{
				State: datasource.ReadResponse{}.State,
			}

			result := resolveNamedDataTable(tt.dtList, data, resp)

			//: Verify the result matches the expected state.
			if tt.expectFound {
				assert.True(t, result, "Should return true when table is found")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have error when found")
			} else {
				assert.False(t, result, "Should return false when table is not found")
				assert.True(t, resp.Diagnostics.HasError(), "Should have error when not found")
			}
		})
	}
}
