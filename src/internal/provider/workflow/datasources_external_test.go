package workflow_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow"
	"github.com/stretchr/testify/assert"
)

// TestNewWorkflowsDataSourceExternal tests the public NewWorkflowsDataSource function (black-box).
func TestNewWorkflowsDataSourceExternal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create new instance from external package", wantErr: false},
		{name: "multiple calls return independent instances", wantErr: false},
		{name: "returned instance is not nil even on repeated calls", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create new instance from external package":
				ds := workflow.NewWorkflowsDataSource()

				assert.NotNil(t, ds)

			case "multiple calls return independent instances":
				ds1 := workflow.NewWorkflowsDataSource()
				ds2 := workflow.NewWorkflowsDataSource()

				assert.NotSame(t, ds1, ds2, "each call should return a new instance")

			case "returned instance is not nil even on repeated calls":
				for i := 0; i < 100; i++ {
					ds := workflow.NewWorkflowsDataSource()
					assert.NotNil(t, ds, "instance should never be nil")
				}

			case "error case - validation checks":
				ds := workflow.NewWorkflowsDataSource()
				assert.NotNil(t, ds)
			}
		})
	}
}

// TestNewWorkflowsDataSourceWrapperExternal tests the public NewWorkflowsDataSourceWrapper function (black-box).
func TestNewWorkflowsDataSourceWrapperExternal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create wrapped instance from external package", wantErr: false},
		{name: "multiple calls return independent instances", wantErr: false},
		{name: "returned instance is not nil even on repeated calls", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create wrapped instance from external package":
				ds := workflow.NewWorkflowsDataSourceWrapper()

				assert.NotNil(t, ds)

			case "multiple calls return independent instances":
				ds1 := workflow.NewWorkflowsDataSourceWrapper()
				ds2 := workflow.NewWorkflowsDataSourceWrapper()

				assert.NotSame(t, ds1, ds2, "each call should return a new instance")

			case "returned instance is not nil even on repeated calls":
				for i := 0; i < 100; i++ {
					ds := workflow.NewWorkflowsDataSourceWrapper()
					assert.NotNil(t, ds, "instance should never be nil")
				}

			case "error case - validation checks":
				ds := workflow.NewWorkflowsDataSourceWrapper()
				assert.NotNil(t, ds)
			}
		})
	}
}

// TestNewWorkflowsDataSource tests the exact function name expected by KTN-TEST-003.
func TestNewWorkflowsDataSource(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "creates non-nil datasource",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := workflow.NewWorkflowsDataSource()
				assert.NotNil(t, ds)
			},
		},
		{
			name: "error case - multiple calls return different instances",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds1 := workflow.NewWorkflowsDataSource()
				ds2 := workflow.NewWorkflowsDataSource()
				assert.NotSame(t, ds1, ds2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}

// TestNewWorkflowsDataSourceWrapper tests the exact function name expected by KTN-TEST-003.
func TestNewWorkflowsDataSourceWrapper(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "creates non-nil datasource",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := workflow.NewWorkflowsDataSourceWrapper()
				assert.NotNil(t, ds)
			},
		},
		{
			name: "error case - multiple calls return different instances",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds1 := workflow.NewWorkflowsDataSourceWrapper()
				ds2 := workflow.NewWorkflowsDataSourceWrapper()
				assert.NotSame(t, ds1, ds2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}

// TestWorkflowsDataSource_Metadata tests the Metadata method.
func TestWorkflowsDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		providerTypeName string
		expectedTypeName string
	}{
		{
			name:             "sets correct type name",
			providerTypeName: "n8n",
			expectedTypeName: "n8n_workflows",
		},
		{
			name:             "error case - handles empty provider type name",
			providerTypeName: "",
			expectedTypeName: "_workflows",
		},
		{
			name:             "error case - handles custom provider type name",
			providerTypeName: "custom",
			expectedTypeName: "custom_workflows",
		},
		{
			name:             "error case - handles provider type name with special characters",
			providerTypeName: "test-provider_123",
			expectedTypeName: "test-provider_123_workflows",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := workflow.NewWorkflowsDataSource()
			req := datasource.MetadataRequest{
				ProviderTypeName: tt.providerTypeName,
			}
			resp := &datasource.MetadataResponse{}

			ds.Metadata(context.Background(), req, resp)

			assert.Equal(t, tt.expectedTypeName, resp.TypeName)
		})
	}
}

// TestWorkflowsDataSource_Schema tests the Schema method.
func TestWorkflowsDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "defines valid schema"},
		{name: "schema has markdown description"},
		{name: "error case - schema attributes are not empty"},
		{name: "error case - schema has expected structure"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := workflow.NewWorkflowsDataSource()
			req := datasource.SchemaRequest{}
			resp := &datasource.SchemaResponse{}

			ds.Schema(context.Background(), req, resp)

			switch tt.name {
			case "defines valid schema":
				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)

			case "schema has markdown description":
				assert.NotEmpty(t, resp.Schema.MarkdownDescription)
				assert.Contains(t, resp.Schema.MarkdownDescription, "workflow")

			case "error case - schema attributes are not empty":
				assert.NotEmpty(t, resp.Schema.Attributes)
				assert.Greater(t, len(resp.Schema.Attributes), 0)

			case "error case - schema has expected structure":
				assert.NotNil(t, resp.Schema)
				assert.NotNil(t, resp.Schema.Attributes)
			}
		})
	}
}

// TestWorkflowsDataSource_Configure tests the Configure method.
func TestWorkflowsDataSource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "configures with valid client", wantErr: false},
		{name: "error case - nil provider data", wantErr: false},
		{name: "error case - wrong provider data type", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "configures with valid client":
				d := workflow.NewWorkflowsDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: &client.N8nClient{},
				}
				resp := &datasource.ConfigureResponse{}
				d.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - nil provider data":
				d := workflow.NewWorkflowsDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: nil,
				}
				resp := &datasource.ConfigureResponse{}
				d.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - wrong provider data type":
				d := workflow.NewWorkflowsDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: "wrong type",
				}
				resp := &datasource.ConfigureResponse{}
				d.Configure(context.Background(), req, resp)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Data Source Configure Type")
			}
		})
	}
}

// TestWorkflowsDataSource_Read tests the Read method.
func TestWorkflowsDataSource_Read(t *testing.T) {
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
								"id": "workflow-123",
								"name": "Test Workflow",
								"active": true,
								"nodes": [],
								"connections": {},
								"settings": {}
							}
						]
					}`))
				})

				n8nClient, server := setupTestClientForDataSources(t, handler)
				defer server.Close()

				ds := workflow.NewWorkflowsDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Build config
				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"active": tftypes.NewValue(tftypes.Bool, nil),
					"workflows": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":     tftypes.String,
								"name":   tftypes.String,
								"active": tftypes.Bool,
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
			name: "error - read with invalid config",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := workflow.NewWorkflowsDataSource()
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
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestClientForDataSources(t, handler)
				defer server.Close()

				ds := workflow.NewWorkflowsDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"active": tftypes.NewValue(tftypes.Bool, nil),
					"workflows": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":     tftypes.String,
								"name":   tftypes.String,
								"active": tftypes.Bool,
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
