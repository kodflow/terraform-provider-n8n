package project_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/project"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// TestNewProjectDataSource tests the NewProjectDataSource constructor.
func TestNewProjectDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "creates valid datasource",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "creates valid datasource":
				ds := project.NewProjectDataSource()
				assert.NotNil(t, ds)

			case "error case - validation checks":
				ds := project.NewProjectDataSource()
				assert.NotNil(t, ds)
				assert.Implements(t, (*datasource.DataSource)(nil), ds)
			}
		})
	}
}

// TestNewProjectDataSourceWrapper tests the NewProjectDataSourceWrapper constructor.
func TestNewProjectDataSourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "creates valid datasource wrapper",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "creates valid datasource wrapper":
				wrapper := project.NewProjectDataSourceWrapper()
				assert.NotNil(t, wrapper)

			case "error case - validation checks":
				wrapper := project.NewProjectDataSourceWrapper()
				assert.NotNil(t, wrapper)
			}
		})
	}
}

// TestProjectDataSource_Metadata tests the Metadata method.
func TestProjectDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "sets correct type name",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "sets correct type name":
				ds := project.NewProjectDataSource()
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), req, resp)
				assert.Equal(t, "n8n_project", resp.TypeName)

			case "error case - validation checks":
				ds := project.NewProjectDataSource()
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), req, resp)
				assert.NotEmpty(t, resp.TypeName)
			}
		})
	}
}

// TestProjectDataSource_Schema tests the Schema method.
func TestProjectDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "returns valid schema",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "returns valid schema":
				ds := project.NewProjectDataSource()
				req := datasource.SchemaRequest{}
				resp := &datasource.SchemaResponse{}
				ds.Schema(context.Background(), req, resp)
				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)

			case "error case - validation checks":
				ds := project.NewProjectDataSource()
				req := datasource.SchemaRequest{}
				resp := &datasource.SchemaResponse{}
				ds.Schema(context.Background(), req, resp)
				assert.NotNil(t, resp.Schema)
			}
		})
	}
}

// TestProjectDataSource_Configure tests the Configure method.
func TestProjectDataSource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "configures with valid client",
			wantErr: false,
		},
		{
			name:    "error case - nil provider data",
			wantErr: false,
		},
		{
			name:    "error case - wrong provider data type",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "configures with valid client":
				ds := project.NewProjectDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: &client.N8nClient{},
				}
				resp := &datasource.ConfigureResponse{}
				ds.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - nil provider data":
				ds := project.NewProjectDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: nil,
				}
				resp := &datasource.ConfigureResponse{}
				ds.Configure(context.Background(), req, resp)
				// Should not error on nil provider data
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - wrong provider data type":
				ds := project.NewProjectDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: "wrong type",
				}
				resp := &datasource.ConfigureResponse{}
				ds.Configure(context.Background(), req, resp)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Data Source Configure Type")
			}
		})
	}
}

// TestProjectDataSource_Read tests the Read method.
func TestProjectDataSource_Read(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "read with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
						"data": [
							{
								"id": "project-123",
								"name": "Test Project",
								"type": "team",
								"createdAt": "2024-01-01T00:00:00Z",
								"updatedAt": "2024-01-02T00:00:00Z",
								"icon": "icon-test",
								"description": "Test description"
							}
						]
					}`))
				})

				n8nClient, server := setupTestClientForDataSource(t, handler)
				defer server.Close()

				ds := project.NewProjectDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Build config
				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, "project-123"),
					"name":        tftypes.NewValue(tftypes.String, nil),
					"type":        tftypes.NewValue(tftypes.String, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"updated_at":  tftypes.NewValue(tftypes.String, nil),
					"icon":        tftypes.NewValue(tftypes.String, nil),
					"description": tftypes.NewValue(tftypes.String, nil),
				})

				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := &datasource.ReadResponse{
					State: state,
				}

				// Call Read
				ds.Read(ctx, req, resp)

				// Verify success
				if resp.Diagnostics.HasError() {
					for _, diag := range resp.Diagnostics.Errors() {
						t.Logf("Error: %s - %s", diag.Summary(), diag.Detail())
					}
				}
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "read with successful API call using name",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
						"data": [
							{
								"id": "project-123",
								"name": "Test Project",
								"type": "team",
								"createdAt": "2024-01-01T00:00:00Z",
								"updatedAt": "2024-01-02T00:00:00Z",
								"icon": "icon-test",
								"description": "Test description"
							}
						]
					}`))
				})

				n8nClient, server := setupTestClientForDataSource(t, handler)
				defer server.Close()

				ds := project.NewProjectDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Build config with name instead of id
				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, nil),
					"name":        tftypes.NewValue(tftypes.String, "Test Project"),
					"type":        tftypes.NewValue(tftypes.String, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"updated_at":  tftypes.NewValue(tftypes.String, nil),
					"icon":        tftypes.NewValue(tftypes.String, nil),
					"description": tftypes.NewValue(tftypes.String, nil),
				})

				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := &datasource.ReadResponse{
					State: state,
				}

				// Call Read
				ds.Read(ctx, req, resp)

				// Verify success
				if resp.Diagnostics.HasError() {
					for _, diag := range resp.Diagnostics.Errors() {
						t.Logf("Error: %s - %s", diag.Summary(), diag.Detail())
					}
				}
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - read with invalid config",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := project.NewProjectDataSource()
				ctx := context.Background()

				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Create invalid config
				configRaw := tftypes.NewValue(tftypes.String, "invalid")

				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := &datasource.ReadResponse{
					State: state,
				}

				ds.Read(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - read with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "Project not found"}`))
				})

				n8nClient, server := setupTestClientForDataSource(t, handler)
				defer server.Close()

				ds := project.NewProjectDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, "project-123"),
					"name":        tftypes.NewValue(tftypes.String, nil),
					"type":        tftypes.NewValue(tftypes.String, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"updated_at":  tftypes.NewValue(tftypes.String, nil),
					"icon":        tftypes.NewValue(tftypes.String, nil),
					"description": tftypes.NewValue(tftypes.String, nil),
				})

				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := &datasource.ReadResponse{
					State: state,
				}

				ds.Read(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - missing required identifiers",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := project.NewProjectDataSource()
				ctx := context.Background()

				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Create config with both id and name as null
				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, nil),
					"name":        tftypes.NewValue(tftypes.String, nil),
					"type":        tftypes.NewValue(tftypes.String, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"updated_at":  tftypes.NewValue(tftypes.String, nil),
					"icon":        tftypes.NewValue(tftypes.String, nil),
					"description": tftypes.NewValue(tftypes.String, nil),
				})

				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := &datasource.ReadResponse{
					State: state,
				}

				ds.Read(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing Required Attribute")
			},
		},
		{
			name: "error - project not found in list",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
						"data": [
							{
								"id": "other-project-456",
								"name": "Other Project",
								"type": "team"
							}
						]
					}`))
				})

				n8nClient, server := setupTestClientForDataSource(t, handler)
				defer server.Close()

				ds := project.NewProjectDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, "project-123"),
					"name":        tftypes.NewValue(tftypes.String, nil),
					"type":        tftypes.NewValue(tftypes.String, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"updated_at":  tftypes.NewValue(tftypes.String, nil),
					"icon":        tftypes.NewValue(tftypes.String, nil),
					"description": tftypes.NewValue(tftypes.String, nil),
				})

				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := &datasource.ReadResponse{
					State: state,
				}

				ds.Read(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Project Not Found")
			},
		},
		{
			name: "error - project not found by name",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
						"data": [
							{
								"id": "other-project-456",
								"name": "Other Project",
								"type": "team"
							}
						]
					}`))
				})

				n8nClient, server := setupTestClientForDataSource(t, handler)
				defer server.Close()

				ds := project.NewProjectDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, nil),
					"name":        tftypes.NewValue(tftypes.String, "NonExistent Project"),
					"type":        tftypes.NewValue(tftypes.String, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"updated_at":  tftypes.NewValue(tftypes.String, nil),
					"icon":        tftypes.NewValue(tftypes.String, nil),
					"description": tftypes.NewValue(tftypes.String, nil),
				})

				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := &datasource.ReadResponse{
					State: state,
				}

				ds.Read(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Project Not Found")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "NonExistent Project")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// setupTestClientForDataSource creates a test N8nClient with httptest server for datasources.
func setupTestClientForDataSource(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
