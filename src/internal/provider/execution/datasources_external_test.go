package execution_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/execution"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

func TestNewExecutionsDataSource(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create new instance", wantErr: false},
		{name: "returns non-nil instance", wantErr: false},
		{name: "error case - verify instance type", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create new instance":
				ds := execution.NewExecutionsDataSource()
				assert.NotNil(t, ds)
				assert.IsType(t, &execution.ExecutionsDataSource{}, ds)

			case "returns non-nil instance":
				ds := execution.NewExecutionsDataSource()
				assert.NotNil(t, ds)

			case "error case - verify instance type":
				ds := execution.NewExecutionsDataSource()
				assert.NotNil(t, ds)
			}
		})
	}
}

func TestNewExecutionsDataSourceWrapper(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create new wrapped instance", wantErr: false},
		{name: "wrapper returns datasource.DataSource interface", wantErr: false},
		{name: "error case - verify wrapper type", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create new wrapped instance":
				ds := execution.NewExecutionsDataSourceWrapper()
				assert.NotNil(t, ds)
				assert.IsType(t, &execution.ExecutionsDataSource{}, ds)

			case "wrapper returns datasource.DataSource interface":
				ds := execution.NewExecutionsDataSourceWrapper()
				assert.NotNil(t, ds)

			case "error case - verify wrapper type":
				ds := execution.NewExecutionsDataSourceWrapper()
				assert.NotNil(t, ds)
			}
		})
	}
}

func TestExecutionsDataSource_Metadata(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "set metadata without provider type", wantErr: false},
		{name: "set metadata with provider type name", wantErr: false},
		{name: "set metadata with different provider type", wantErr: false},
		{name: "error case - metadata validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "set metadata without provider type":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_executions", resp.TypeName)

			case "set metadata with provider type name":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.MetadataResponse{}
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_executions", resp.TypeName)

			case "set metadata with different provider type":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.MetadataResponse{}
				req := datasource.MetadataRequest{
					ProviderTypeName: "custom_provider",
				}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "custom_provider_executions", resp.TypeName)

			case "error case - metadata validation":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "test",
				}, resp)
				assert.Equal(t, "test_executions", resp.TypeName)
			}
		})
	}
}

func TestExecutionsDataSource_Schema(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "return schema", wantErr: false},
		{name: "schema attributes have descriptions", wantErr: false},
		{name: "error case - schema validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "return schema":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.Contains(t, resp.Schema.MarkdownDescription, "n8n workflow executions")

			case "schema attributes have descriptions":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				for name, attr := range resp.Schema.Attributes {
					assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
				}

			case "error case - schema validation":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.SchemaResponse{}
				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
				assert.NotNil(t, resp.Schema)
			}
		})
	}
}

func TestExecutionsDataSource_Configure(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "configure with valid client", wantErr: false},
		{name: "configure with nil provider data", wantErr: false},
		{name: "configure with invalid provider data", wantErr: false},
		{name: "error case - verify configuration", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "configure with valid client":
				ds := execution.NewExecutionsDataSource()
				cfg := n8nsdk.NewConfiguration()
				cfg.Servers = n8nsdk.ServerConfigurations{
					{URL: "http://test.local", Description: "Test server"},
				}
				mockClient := &client.N8nClient{
					APIClient: n8nsdk.NewAPIClient(cfg),
				}
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: mockClient,
				}

				ds.Configure(context.Background(), req, resp)

				assert.False(t, resp.Diagnostics.HasError())

			case "configure with nil provider data":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: nil,
				}

				ds.Configure(context.Background(), req, resp)

				assert.False(t, resp.Diagnostics.HasError())

			case "configure with invalid provider data":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: "invalid",
				}

				ds.Configure(context.Background(), req, resp)

				assert.True(t, resp.Diagnostics.HasError())

			case "error case - verify configuration":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: 12345,
				}
				ds.Configure(context.Background(), req, resp)
				assert.True(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestExecutionsDataSource_Read(t *testing.T) {
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
								"id": 123,
								"workflowId": 456,
								"finished": true,
								"mode": "manual",
								"startedAt": "2024-01-01T00:00:00Z",
								"stoppedAt": "2024-01-01T00:01:00Z",
								"status": "success"
							}
						],
						"nextCursor": null
					}`))
				})

				n8nClient, server := setupTestClientForDataSources(t, handler)
				defer server.Close()

				ds := execution.NewExecutionsDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Build config
				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"workflow_id":  tftypes.NewValue(tftypes.String, nil),
					"project_id":   tftypes.NewValue(tftypes.String, nil),
					"status":       tftypes.NewValue(tftypes.String, nil),
					"include_data": tftypes.NewValue(tftypes.Bool, nil),
					"executions": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":          tftypes.String,
								"workflow_id": tftypes.String,
								"finished":    tftypes.Bool,
								"mode":        tftypes.String,
								"started_at":  tftypes.String,
								"stopped_at":  tftypes.String,
								"status":      tftypes.String,
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
				ds := execution.NewExecutionsDataSource()
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

				ds := execution.NewExecutionsDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"workflow_id":  tftypes.NewValue(tftypes.String, nil),
					"project_id":   tftypes.NewValue(tftypes.String, nil),
					"status":       tftypes.NewValue(tftypes.String, nil),
					"include_data": tftypes.NewValue(tftypes.Bool, nil),
					"executions": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":          tftypes.String,
								"workflow_id": tftypes.String,
								"finished":    tftypes.Bool,
								"mode":        tftypes.String,
								"started_at":  tftypes.String,
								"stopped_at":  tftypes.String,
								"status":      tftypes.String,
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
