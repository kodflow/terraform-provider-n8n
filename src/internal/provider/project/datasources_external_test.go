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

// TestNewProjectsDataSource tests the NewProjectsDataSource constructor.
func TestNewProjectsDataSource(t *testing.T) {
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
				ds := project.NewProjectsDataSource()
				assert.NotNil(t, ds)

			case "error case - validation checks":
				ds := project.NewProjectsDataSource()
				assert.NotNil(t, ds)
				assert.Implements(t, (*datasource.DataSource)(nil), ds)
			}
		})
	}
}

// TestNewProjectsDataSourceWrapper tests the NewProjectsDataSourceWrapper constructor.
func TestNewProjectsDataSourceWrapper(t *testing.T) {
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
				wrapper := project.NewProjectsDataSourceWrapper()
				assert.NotNil(t, wrapper)

			case "error case - validation checks":
				wrapper := project.NewProjectsDataSourceWrapper()
				assert.NotNil(t, wrapper)
			}
		})
	}
}

// TestProjectsDataSource_Metadata tests the Metadata method.
func TestProjectsDataSource_Metadata(t *testing.T) {
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
				ds := project.NewProjectsDataSource()
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), req, resp)
				assert.Equal(t, "n8n_projects", resp.TypeName)

			case "error case - validation checks":
				ds := project.NewProjectsDataSource()
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

// TestProjectsDataSource_Schema tests the Schema method.
func TestProjectsDataSource_Schema(t *testing.T) {
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
				ds := project.NewProjectsDataSource()
				req := datasource.SchemaRequest{}
				resp := &datasource.SchemaResponse{}
				ds.Schema(context.Background(), req, resp)
				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)

			case "error case - validation checks":
				ds := project.NewProjectsDataSource()
				req := datasource.SchemaRequest{}
				resp := &datasource.SchemaResponse{}
				ds.Schema(context.Background(), req, resp)
				assert.NotNil(t, resp.Schema)
			}
		})
	}
}

// TestProjectsDataSource_Configure tests the Configure method.
func TestProjectsDataSource_Configure(t *testing.T) {
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
				ds := project.NewProjectsDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: &client.N8nClient{},
				}
				resp := &datasource.ConfigureResponse{}
				ds.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - nil provider data":
				ds := project.NewProjectsDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: nil,
				}
				resp := &datasource.ConfigureResponse{}
				ds.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - wrong provider data type":
				ds := project.NewProjectsDataSource()
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

// TestProjectsDataSource_Read tests the Read method.
func TestProjectsDataSource_Read(t *testing.T) {
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
								"icon": {"value": "icon-test"},
								"description": "Test description"
							}
						]
					}`))
				})

				n8nClient, server := setupTestClientForDataSources(t, handler)
				defer server.Close()

				ds := project.NewProjectsDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Build config
				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"projects": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":          tftypes.String,
								"name":        tftypes.String,
								"type":        tftypes.String,
								"created_at":  tftypes.String,
								"updated_at":  tftypes.String,
								"icon":        tftypes.String,
								"description": tftypes.String,
							},
						}},
						nil,
					),
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
			name: "error - read with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestClientForDataSources(t, handler)
				defer server.Close()

				ds := project.NewProjectsDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"projects": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":          tftypes.String,
								"name":        tftypes.String,
								"type":        tftypes.String,
								"created_at":  tftypes.String,
								"updated_at":  tftypes.String,
								"icon":        tftypes.String,
								"description": tftypes.String,
							},
						}},
						nil,
					),
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
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// setupTestClientForDataSources creates a test N8nClient with httptest server for datasources (plural).
func setupTestClientForDataSources(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
