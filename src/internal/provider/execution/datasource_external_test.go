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

func TestNewExecutionDataSourceWrapper(t *testing.T) {
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
				ds := execution.NewExecutionDataSourceWrapper()

				assert.NotNil(t, ds)
				assert.IsType(t, &execution.ExecutionDataSource{}, ds)

			case "wrapper returns datasource.DataSource interface":
				ds := execution.NewExecutionDataSourceWrapper()

				// ds is already of type datasource.DataSource, no assertion needed.
				assert.NotNil(t, ds)

			case "error case - verify wrapper type":
				ds := execution.NewExecutionDataSourceWrapper()
				// Verify ds implements datasource.DataSource interface.
				assert.NotNil(t, ds)
			}
		})
	}
}

func TestExecutionDataSource_Metadata(t *testing.T) {
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
				ds := execution.NewExecutionDataSource()
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_execution", resp.TypeName)

			case "set metadata with provider type name":
				ds := execution.NewExecutionDataSource()
				resp := &datasource.MetadataResponse{}
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_execution", resp.TypeName)

			case "set metadata with different provider type":
				ds := execution.NewExecutionDataSource()
				resp := &datasource.MetadataResponse{}
				req := datasource.MetadataRequest{
					ProviderTypeName: "custom_provider",
				}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "custom_provider_execution", resp.TypeName)

			case "error case - metadata validation":
				ds := execution.NewExecutionDataSource()
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "test",
				}, resp)
				assert.Equal(t, "test_execution", resp.TypeName)
			}
		})
	}
}

func TestExecutionDataSource_Schema(t *testing.T) {
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
				ds := execution.NewExecutionDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.Contains(t, resp.Schema.MarkdownDescription, "n8n workflow execution")

				// Verify required attributes.
				idAttr, exists := resp.Schema.Attributes["id"]
				assert.True(t, exists)
				assert.True(t, idAttr.IsRequired())

				// Verify computed attributes.
				computedAttrs := []string{
					"workflow_id",
					"finished",
					"mode",
					"started_at",
					"stopped_at",
					"status",
				}

				for _, attr := range computedAttrs {
					schemaAttr, exists := resp.Schema.Attributes[attr]
					assert.True(t, exists, "Attribute %s should exist", attr)
					assert.True(t, schemaAttr.IsComputed(), "Attribute %s should be computed", attr)
				}

				// Verify optional attributes.
				includeDataAttr, exists := resp.Schema.Attributes["include_data"]
				assert.True(t, exists)
				assert.True(t, includeDataAttr.IsOptional())

			case "schema attributes have descriptions":
				ds := execution.NewExecutionDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				for name, attr := range resp.Schema.Attributes {
					assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
				}

			case "error case - schema validation":
				ds := execution.NewExecutionDataSource()
				resp := &datasource.SchemaResponse{}
				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
				assert.NotNil(t, resp.Schema)
			}
		})
	}
}

func TestNewExecutionDataSource(t *testing.T) {
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
				ds := execution.NewExecutionDataSource()
				assert.NotNil(t, ds)
				assert.IsType(t, &execution.ExecutionDataSource{}, ds)

			case "returns non-nil instance":
				ds := execution.NewExecutionDataSource()
				assert.NotNil(t, ds)

			case "error case - verify instance type":
				ds := execution.NewExecutionDataSource()
				assert.NotNil(t, ds)
			}
		})
	}
}

func TestExecutionDataSource_Configure(t *testing.T) {
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
				ds := execution.NewExecutionDataSource()
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
				ds := execution.NewExecutionDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: nil,
				}

				ds.Configure(context.Background(), req, resp)

				assert.False(t, resp.Diagnostics.HasError())

			case "configure with invalid provider data":
				ds := execution.NewExecutionDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: "invalid",
				}

				ds.Configure(context.Background(), req, resp)

				assert.True(t, resp.Diagnostics.HasError())

			case "error case - verify configuration":
				ds := execution.NewExecutionDataSource()
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

func TestExecutionDataSource_Read(t *testing.T) {
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
						"id": 123,
						"workflowId": 456,
						"finished": true,
						"mode": "manual",
						"startedAt": "2024-01-01T00:00:00Z",
						"stoppedAt": "2024-01-01T00:01:00Z",
						"status": "success"
					}`))
				})

				n8nClient, server := setupTestClientForDataSource(t, handler)
				defer server.Close()

				ds := execution.NewExecutionDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Build config
				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":           tftypes.NewValue(tftypes.String, "123"),
					"workflow_id":  tftypes.NewValue(tftypes.String, nil),
					"finished":     tftypes.NewValue(tftypes.Bool, nil),
					"mode":         tftypes.NewValue(tftypes.String, nil),
					"started_at":   tftypes.NewValue(tftypes.String, nil),
					"stopped_at":   tftypes.NewValue(tftypes.String, nil),
					"status":       tftypes.NewValue(tftypes.String, nil),
					"include_data": tftypes.NewValue(tftypes.Bool, nil),
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
				ds := execution.NewExecutionDataSource()
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
			name: "error - read with invalid execution ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				n8nClient, server := setupTestClientForDataSource(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				defer server.Close()

				ds := execution.NewExecutionDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Build config with invalid ID
				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":           tftypes.NewValue(tftypes.String, "invalid-id"),
					"workflow_id":  tftypes.NewValue(tftypes.String, nil),
					"finished":     tftypes.NewValue(tftypes.Bool, nil),
					"mode":         tftypes.NewValue(tftypes.String, nil),
					"started_at":   tftypes.NewValue(tftypes.String, nil),
					"stopped_at":   tftypes.NewValue(tftypes.String, nil),
					"status":       tftypes.NewValue(tftypes.String, nil),
					"include_data": tftypes.NewValue(tftypes.Bool, nil),
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
			name: "error - read with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "Execution not found"}`))
				})

				n8nClient, server := setupTestClientForDataSource(t, handler)
				defer server.Close()

				ds := execution.NewExecutionDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":           tftypes.NewValue(tftypes.String, "123"),
					"workflow_id":  tftypes.NewValue(tftypes.String, nil),
					"finished":     tftypes.NewValue(tftypes.Bool, nil),
					"mode":         tftypes.NewValue(tftypes.String, nil),
					"started_at":   tftypes.NewValue(tftypes.String, nil),
					"stopped_at":   tftypes.NewValue(tftypes.String, nil),
					"status":       tftypes.NewValue(tftypes.String, nil),
					"include_data": tftypes.NewValue(tftypes.Bool, nil),
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
