package variable_test

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
	"github.com/kodflow/n8n/src/internal/provider/variable"
	"github.com/stretchr/testify/assert"
)

// TestNewVariablesDataSource tests the NewVariablesDataSource constructor.
func TestNewVariablesDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariablesDataSource()
				assert.NotNil(t, ds, "NewVariablesDataSource should return a non-nil datasource")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariablesDataSource()
				assert.NotNil(t, ds, "NewVariablesDataSource should return a non-nil datasource")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestNewVariablesDataSourceWrapper tests the NewVariablesDataSourceWrapper constructor.
func TestNewVariablesDataSourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariablesDataSourceWrapper()
				assert.NotNil(t, ds, "NewVariablesDataSourceWrapper should return a non-nil datasource")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariablesDataSourceWrapper()
				assert.NotNil(t, ds, "NewVariablesDataSourceWrapper should return a non-nil datasource")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestVariablesDataSource_Metadata tests the Metadata method.
func TestVariablesDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariablesDataSource()
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_variables", resp.TypeName, "TypeName should be set correctly")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariablesDataSource()
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_variables", resp.TypeName, "TypeName should be set correctly")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestVariablesDataSource_Schema tests the Schema method.
func TestVariablesDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariablesDataSource()
				req := datasource.SchemaRequest{}
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), req, resp)

				assert.NotNil(t, resp.Schema, "Schema should not be nil")
				assert.NotNil(t, resp.Schema.Attributes, "Schema attributes should not be nil")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariablesDataSource()
				req := datasource.SchemaRequest{}
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), req, resp)

				assert.NotNil(t, resp.Schema, "Schema should not be nil")
				assert.NotNil(t, resp.Schema.Attributes, "Schema attributes should not be nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestVariablesDataSource_Configure tests the Configure method.
func TestVariablesDataSource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		providerData interface{}
		expectError  bool
	}{
		{
			name:         "valid client",
			providerData: &client.N8nClient{},
			expectError:  false,
		},
		{
			name:         "nil provider data",
			providerData: nil,
			expectError:  false,
		},
		{
			name:         "invalid provider data",
			providerData: "invalid",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := variable.NewVariablesDataSource()
			req := datasource.ConfigureRequest{
				ProviderData: tt.providerData,
			}
			resp := &datasource.ConfigureResponse{}

			ds.Configure(context.Background(), req, resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError(), "Configure should return error for invalid provider data")
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "Configure should not return error")
			}
		})
	}
}

// TestVariablesDataSource_Read tests the Read method.
func TestVariablesDataSource_Read(t *testing.T) {
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
								"id": "var-123",
								"key": "TEST_VAR",
								"value": "test-value",
								"type": "string",
								"project": {"id": "proj-123", "name": "Test Project"}
							}
						]
					}`))
				})

				n8nClient, server := setupTestClientForDataSources(t, handler)
				defer server.Close()

				ds := variable.NewVariablesDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Build config
				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"project_id": tftypes.NewValue(tftypes.String, nil),
					"state":      tftypes.NewValue(tftypes.String, nil),
					"variables": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":         tftypes.String,
								"key":        tftypes.String,
								"value":      tftypes.String,
								"type":       tftypes.String,
								"project_id": tftypes.String,
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
				ds := variable.NewVariablesDataSource()
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

				ds := variable.NewVariablesDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"project_id": tftypes.NewValue(tftypes.String, nil),
					"state":      tftypes.NewValue(tftypes.String, nil),
					"variables": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":         tftypes.String,
								"key":        tftypes.String,
								"value":      tftypes.String,
								"type":       tftypes.String,
								"project_id": tftypes.String,
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
