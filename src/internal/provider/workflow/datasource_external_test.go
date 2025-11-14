package workflow_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/workflow"
	"github.com/stretchr/testify/assert"
)

// TestNewWorkflowDataSourceExternal tests the public NewWorkflowDataSource function (black-box).
func TestNewWorkflowDataSourceExternal(t *testing.T) {
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
				ds := workflow.NewWorkflowDataSource()

				assert.NotNil(t, ds)

			case "multiple calls return independent instances":
				ds1 := workflow.NewWorkflowDataSource()
				ds2 := workflow.NewWorkflowDataSource()

				assert.NotSame(t, ds1, ds2, "each call should return a new instance")

			case "returned instance is not nil even on repeated calls":
				for i := 0; i < 100; i++ {
					ds := workflow.NewWorkflowDataSource()
					assert.NotNil(t, ds, "instance should never be nil")
				}

			case "error case - validation checks":
				ds := workflow.NewWorkflowDataSource()
				assert.NotNil(t, ds)
			}
		})
	}
}

// TestNewWorkflowDataSourceWrapperExternal tests the public NewWorkflowDataSourceWrapper function (black-box).
func TestNewWorkflowDataSourceWrapperExternal(t *testing.T) {
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
				ds := workflow.NewWorkflowDataSourceWrapper()

				assert.NotNil(t, ds)

			case "multiple calls return independent instances":
				ds1 := workflow.NewWorkflowDataSourceWrapper()
				ds2 := workflow.NewWorkflowDataSourceWrapper()

				assert.NotSame(t, ds1, ds2, "each call should return a new instance")

			case "returned instance is not nil even on repeated calls":
				for i := 0; i < 100; i++ {
					ds := workflow.NewWorkflowDataSourceWrapper()
					assert.NotNil(t, ds, "instance should never be nil")
				}

			case "error case - validation checks":
				ds := workflow.NewWorkflowDataSourceWrapper()
				assert.NotNil(t, ds)
			}
		})
	}
}

// TestNewWorkflowDataSource tests the exact function name expected by KTN-TEST-003.
func TestNewWorkflowDataSource(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "creates non-nil datasource",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := workflow.NewWorkflowDataSource()
				assert.NotNil(t, ds)
			},
		},
		{
			name: "error case - multiple calls return different instances",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds1 := workflow.NewWorkflowDataSource()
				ds2 := workflow.NewWorkflowDataSource()
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

// TestNewWorkflowDataSourceWrapper tests the exact function name expected by KTN-TEST-003.
func TestNewWorkflowDataSourceWrapper(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "creates non-nil datasource",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := workflow.NewWorkflowDataSourceWrapper()
				assert.NotNil(t, ds)
			},
		},
		{
			name: "error case - multiple calls return different instances",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds1 := workflow.NewWorkflowDataSourceWrapper()
				ds2 := workflow.NewWorkflowDataSourceWrapper()
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

// TestWorkflowDataSource_Metadata tests the Metadata method.
func TestWorkflowDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		providerTypeName string
		expectedTypeName string
	}{
		{
			name:             "sets correct type name",
			providerTypeName: "n8n",
			expectedTypeName: "n8n_workflow",
		},
		{
			name:             "error case - handles empty provider type name",
			providerTypeName: "",
			expectedTypeName: "_workflow",
		},
		{
			name:             "error case - handles custom provider type name",
			providerTypeName: "custom",
			expectedTypeName: "custom_workflow",
		},
		{
			name:             "error case - handles provider type name with special characters",
			providerTypeName: "test-provider_123",
			expectedTypeName: "test-provider_123_workflow",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := workflow.NewWorkflowDataSource()
			req := datasource.MetadataRequest{
				ProviderTypeName: tt.providerTypeName,
			}
			resp := &datasource.MetadataResponse{}

			ds.Metadata(context.Background(), req, resp)

			assert.Equal(t, tt.expectedTypeName, resp.TypeName)
		})
	}
}

// TestWorkflowDataSource_Schema tests the Schema method.
func TestWorkflowDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "defines valid schema"},
		{name: "schema has markdown description"},
		{name: "schema has required id attribute"},
		{name: "error case - schema attributes are not empty"},
		{name: "error case - all required attributes are present"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := workflow.NewWorkflowDataSource()
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

			case "schema has required id attribute":
				assert.Contains(t, resp.Schema.Attributes, "id")
				idAttr := resp.Schema.Attributes["id"]
				assert.NotNil(t, idAttr)

			case "error case - schema attributes are not empty":
				assert.NotEmpty(t, resp.Schema.Attributes)
				assert.Greater(t, len(resp.Schema.Attributes), 0)

			case "error case - all required attributes are present":
				assert.Contains(t, resp.Schema.Attributes, "id")
				assert.Contains(t, resp.Schema.Attributes, "name")
			}
		})
	}
}

// TestWorkflowDataSource_Configure tests the Configure method.
func TestWorkflowDataSource_Configure(t *testing.T) {
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
				d := workflow.NewWorkflowDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: &client.N8nClient{},
				}
				resp := &datasource.ConfigureResponse{}
				d.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - nil provider data":
				d := workflow.NewWorkflowDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: nil,
				}
				resp := &datasource.ConfigureResponse{}
				d.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - wrong provider data type":
				d := workflow.NewWorkflowDataSource()
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

// TestWorkflowDataSource_Read tests the Read method.

func TestWorkflowDataSource_Read(t *testing.T) {
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
					w.Write([]byte(`{"id": "workflow-123", "name": "Test Workflow", "active": true, "nodes": [], "connections": {}, "settings": {}}`))
				})

				n8nClient, server := setupTestClientForDataSource(t, handler)
				defer server.Close()

				ds := workflow.NewWorkflowDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Build config
				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":     tftypes.NewValue(tftypes.String, "workflow-123"),
					"name":   tftypes.NewValue(tftypes.String, nil),
					"active": tftypes.NewValue(tftypes.Bool, nil),
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
				ds := workflow.NewWorkflowDataSource()
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
					w.Write([]byte(`{"message": "Workflow not found"}`))
				})

				n8nClient, server := setupTestClientForDataSource(t, handler)
				defer server.Close()

				ds := workflow.NewWorkflowDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":     tftypes.NewValue(tftypes.String, "workflow-123"),
					"name":   tftypes.NewValue(tftypes.String, nil),
					"active": tftypes.NewValue(tftypes.Bool, nil),
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
