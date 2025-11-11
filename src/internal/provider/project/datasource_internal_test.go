package project

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// createTestDataSourceSchema creates a test schema for project datasource.
func createTestDataSourceSchema(t *testing.T) datasource.SchemaResponse {
	t.Helper()
	d := &ProjectDataSource{}
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), req, resp)
	return *resp
}

// setupTestDataSourceClient creates a test N8nClient with httptest server.
func setupTestDataSourceClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

// TestProjectDataSource_Read tests the Read method.
func TestProjectDataSource_Read(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		idValue     interface{}
		nameValue   interface{}
		wantErr     bool
		errContains string
	}{
		{
			name: "successful read by id",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					projectType := "team"
					icon := "icon-test"
					description := "test description"
					projects := map[string]interface{}{
						"data": []map[string]interface{}{
							{
								"id":          "proj-123",
								"name":        "Test Project",
								"type":        projectType,
								"createdAt":   now.Format(time.RFC3339),
								"updatedAt":   now.Format(time.RFC3339),
								"icon":        icon,
								"description": description,
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(projects)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			idValue:   "proj-123",
			nameValue: nil,
			wantErr:   false,
		},
		{
			name: "successful read by name",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					projectType := "team"
					projects := map[string]interface{}{
						"data": []map[string]interface{}{
							{
								"id":   "proj-123",
								"name": "Test Project",
								"type": projectType,
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(projects)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			idValue:   nil,
			nameValue: "Test Project",
			wantErr:   false,
		},
		{
			name: "missing required attribute",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}),
			idValue:     nil,
			nameValue:   nil,
			wantErr:     true,
			errContains: "Missing Required Attribute",
		},
		{
			name: "api error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			idValue:     "proj-123",
			nameValue:   nil,
			wantErr:     true,
			errContains: "Error listing projects",
		},
		{
			name: "project not found",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					projects := map[string]interface{}{
						"data": []map[string]interface{}{},
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(projects)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			idValue:     "proj-123",
			nameValue:   nil,
			wantErr:     true,
			errContains: "Project Not Found",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestDataSourceClient(t, tt.handler)
			defer server.Close()

			d := &ProjectDataSource{client: n8nClient}

			schemaResp := createTestDataSourceSchema(t)

			rawConfig := map[string]tftypes.Value{
				"id":          tftypes.NewValue(tftypes.String, tt.idValue),
				"name":        tftypes.NewValue(tftypes.String, tt.nameValue),
				"type":        tftypes.NewValue(tftypes.String, nil),
				"created_at":  tftypes.NewValue(tftypes.String, nil),
				"updated_at":  tftypes.NewValue(tftypes.String, nil),
				"icon":        tftypes.NewValue(tftypes.String, nil),
				"description": tftypes.NewValue(tftypes.String, nil),
			}
			config := tfsdk.Config{
				Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String, "icon": tftypes.String, "description": tftypes.String}}, rawConfig),
				Schema: schemaResp.Schema,
			}

			state := tfsdk.State{
				Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String, "icon": tftypes.String, "description": tftypes.String}}, nil),
				Schema: schemaResp.Schema,
			}

			req := datasource.ReadRequest{
				Config: config,
			}
			resp := datasource.ReadResponse{
				State: state,
			}

			d.Read(context.Background(), req, &resp)

			if tt.wantErr {
				assert.True(t, resp.Diagnostics.HasError(), "expected error but got none")
				if tt.errContains != "" {
					found := false
					for _, diag := range resp.Diagnostics.Errors() {
						if assert.Contains(t, diag.Summary(), tt.errContains) {
							found = true
							break
						}
					}
					assert.True(t, found, "expected error containing %s", tt.errContains)
				}
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "unexpected error: %v", resp.Diagnostics)
			}
		})
	}
}

// TestProjectDataSource_Configure tests the Configure method.
func TestProjectDataSource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		providerData interface{}
		wantErr      bool
		errContains  string
	}{
		{
			name:         "nil provider data",
			providerData: nil,
			wantErr:      false,
		},
		{
			name:         "valid provider data",
			providerData: &client.N8nClient{},
			wantErr:      false,
		},
		{
			name:         "invalid provider data type",
			providerData: "invalid",
			wantErr:      true,
			errContains:  "Unexpected Data Source Configure Type",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := &ProjectDataSource{}

			req := datasource.ConfigureRequest{
				ProviderData: tt.providerData,
			}
			resp := &datasource.ConfigureResponse{}

			d.Configure(context.Background(), req, resp)

			if tt.wantErr {
				assert.True(t, resp.Diagnostics.HasError(), "expected error but got none")
				if tt.errContains != "" {
					found := false
					for _, diag := range resp.Diagnostics.Errors() {
						if assert.Contains(t, diag.Summary(), tt.errContains) {
							found = true
							break
						}
					}
					assert.True(t, found, "expected error containing %s", tt.errContains)
				}
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "unexpected error: %v", resp.Diagnostics)
			}
		})
	}
}
