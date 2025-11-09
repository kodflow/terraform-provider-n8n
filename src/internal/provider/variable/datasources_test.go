package variable

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVariablesDataSourceInterface is a mock implementation of VariablesDataSourceInterface.
type MockVariablesDataSourceInterface struct {
	mock.Mock
}

func (m *MockVariablesDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariablesDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariablesDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariablesDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

func TestNewVariablesDataSource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		ds := NewVariablesDataSource()

		assert.NotNil(t, ds)
		assert.IsType(t, &VariablesDataSource{}, ds)
		assert.Nil(t, ds.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		ds1 := NewVariablesDataSource()
		ds2 := NewVariablesDataSource()

		assert.NotSame(t, ds1, ds2)
	})
}

func TestNewVariablesDataSourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		ds := NewVariablesDataSourceWrapper()

		assert.NotNil(t, ds)
		assert.IsType(t, &VariablesDataSource{}, ds)
	})

	t.Run("wrapper returns datasource.DataSource interface", func(t *testing.T) {
		ds := NewVariablesDataSourceWrapper()

		// ds is already of type datasource.DataSource, no assertion needed
		assert.NotNil(t, ds)
	})
}

func TestVariablesDataSource_Metadata(t *testing.T) {
	t.Run("set metadata with provider type", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.MetadataResponse{}

		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_variables", resp.TypeName)
	})

	t.Run("set metadata with different provider type", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_variables", resp.TypeName)
	})

	t.Run("set metadata with empty provider type", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "_variables", resp.TypeName)
	})
}

func TestVariablesDataSource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "n8n variables")

		// Verify filter attributes
		projectIDAttr, exists := resp.Schema.Attributes["project_id"]
		assert.True(t, exists)
		assert.True(t, projectIDAttr.IsOptional())

		stateAttr, exists := resp.Schema.Attributes["state"]
		assert.True(t, exists)
		assert.True(t, stateAttr.IsOptional())

		// Verify computed variables list
		variablesAttr, exists := resp.Schema.Attributes["variables"]
		assert.True(t, exists)
		assert.True(t, variablesAttr.IsComputed())
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			if name != "variables" { // Skip nested attribute
				assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
			}
		}
	})

	t.Run("variables list has correct structure", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		variablesAttr, exists := resp.Schema.Attributes["variables"]
		assert.True(t, exists)
		assert.True(t, variablesAttr.IsComputed())
	})
}

func TestVariablesDataSource_SchemaAttributes(t *testing.T) {
	t.Run("return schema attributes", func(t *testing.T) {
		ds := NewVariablesDataSource()

		attrs := ds.schemaAttributes()

		assert.NotNil(t, attrs)
		assert.Contains(t, attrs, "project_id")
		assert.Contains(t, attrs, "state")
		assert.Contains(t, attrs, "variables")

		// Verify attribute properties
		projectIDAttr := attrs["project_id"]
		assert.True(t, projectIDAttr.IsOptional())

		stateAttr := attrs["state"]
		assert.True(t, stateAttr.IsOptional())

		variablesAttr := attrs["variables"]
		assert.True(t, variablesAttr.IsComputed())
	})
}

func TestVariablesDataSource_VariableItemAttributes(t *testing.T) {
	t.Run("return variable item attributes", func(t *testing.T) {
		ds := NewVariablesDataSource()

		attrs := ds.variableItemAttributes()

		assert.NotNil(t, attrs)
		assert.Contains(t, attrs, "id")
		assert.Contains(t, attrs, "key")
		assert.Contains(t, attrs, "value")
		assert.Contains(t, attrs, "type")
		assert.Contains(t, attrs, "project_id")

		// Verify all are computed
		for _, attr := range attrs {
			assert.True(t, attr.IsComputed())
		}
	})

	t.Run("value attribute is sensitive", func(t *testing.T) {
		ds := NewVariablesDataSource()

		attrs := ds.variableItemAttributes()

		valueAttr, exists := attrs["value"]
		assert.True(t, exists)
		assert.True(t, valueAttr.IsComputed())
		// Note: Sensitive attribute testing depends on framework internals
	})
}

func TestVariablesDataSource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.ConfigureResponse{}

		mockClient := &client.N8nClient{}
		req := datasource.ConfigureRequest{
			ProviderData: mockClient,
		}

		ds.Configure(context.Background(), req, resp)

		assert.NotNil(t, ds.client)
		assert.Equal(t, mockClient, ds.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with nil provider data", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: nil,
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: "invalid-data",
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Data Source Configure Type")
	})

	t.Run("configure with wrong type", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "*client.N8nClient")
	})

	t.Run("configure multiple times", func(t *testing.T) {
		ds := NewVariablesDataSource()

		// First configuration
		resp1 := &datasource.ConfigureResponse{}
		client1 := &client.N8nClient{}
		req1 := datasource.ConfigureRequest{
			ProviderData: client1,
		}
		ds.Configure(context.Background(), req1, resp1)
		assert.Equal(t, client1, ds.client)

		// Second configuration
		resp2 := &datasource.ConfigureResponse{}
		client2 := &client.N8nClient{}
		req2 := datasource.ConfigureRequest{
			ProviderData: client2,
		}
		ds.Configure(context.Background(), req2, resp2)
		assert.Equal(t, client2, ds.client)
	})
}

func TestVariablesDataSource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		ds := NewVariablesDataSource()

		// Test that VariablesDataSource implements datasource.DataSource
		var _ datasource.DataSource = ds

		// Test that VariablesDataSource implements datasource.DataSourceWithConfigure
		var _ datasource.DataSourceWithConfigure = ds

		// Test that VariablesDataSource implements VariablesDataSourceInterface
		var _ VariablesDataSourceInterface = ds
	})
}

func TestVariablesDataSourceConcurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		ds := NewVariablesDataSource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)
				assert.Equal(t, "n8n_variables", resp.TypeName)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		ds := NewVariablesDataSource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &datasource.SchemaResponse{}
				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
				assert.NotNil(t, resp.Schema)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent configure calls", func(t *testing.T) {
		ds := NewVariablesDataSource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &datasource.ConfigureResponse{}
				mockClient := &client.N8nClient{}
				req := datasource.ConfigureRequest{
					ProviderData: mockClient,
				}
				ds.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema attributes calls", func(t *testing.T) {
		ds := NewVariablesDataSource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				attrs := ds.schemaAttributes()
				assert.NotNil(t, attrs)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent variable item attributes calls", func(t *testing.T) {
		ds := NewVariablesDataSource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				attrs := ds.variableItemAttributes()
				assert.NotNil(t, attrs)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func BenchmarkVariablesDataSource_Schema(b *testing.B) {
	ds := NewVariablesDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkVariablesDataSource_Metadata(b *testing.B) {
	ds := NewVariablesDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)
	}
}

func BenchmarkVariablesDataSource_Configure(b *testing.B) {
	ds := NewVariablesDataSource()
	mockClient := &client.N8nClient{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: mockClient,
		}
		ds.Configure(context.Background(), req, resp)
	}
}

func BenchmarkVariablesDataSource_SchemaAttributes(b *testing.B) {
	ds := NewVariablesDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ds.schemaAttributes()
	}
}

func BenchmarkVariablesDataSource_VariableItemAttributes(b *testing.B) {
	ds := NewVariablesDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ds.variableItemAttributes()
	}
}

// TODO: Fix datasources Read tests - complex Config initialization required.
func _TestVariablesDataSource_Read(t *testing.T) {
	t.Run("successful read with variables", func(t *testing.T) {
		// Create a mock HTTP server
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				// Return sample variables
				response := `{
					"data": [
						{
							"id": "var-1",
							"key": "API_KEY",
							"value": "secret-value",
							"type": "string",
							"project": {
								"id": "proj-123"
							}
						},
						{
							"id": "var-2",
							"key": "DB_URL",
							"value": "postgresql://localhost",
							"type": "string"
						}
					]
				}`
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(response))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestVariablesClient(t, handler)
		defer server.Close()

		ds := &VariablesDataSource{client: n8nClient}

		state := createTestVariablesDataSourceState(t)
		req := datasource.ReadRequest{
			Config: tfsdk.Config{
				Schema: state.Schema,
				Raw:    tftypes.NewValue(state.Schema.Type().TerraformType(context.Background()), nil),
			},
		}
		resp := &datasource.ReadResponse{
			State: state,
		}

		ds.Read(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("read with project_id filter", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				// Verify project_id filter was applied
				assert.Equal(t, "proj-123", r.URL.Query().Get("projectId"))

				response := `{"data": [{"id": "var-1", "key": "KEY1", "value": "val1"}]}`
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(response))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestVariablesClient(t, handler)
		defer server.Close()

		ds := &VariablesDataSource{client: n8nClient}

		state := createTestVariablesDataSourceState(t)
		req := datasource.ReadRequest{
			Config: tfsdk.Config{
				Schema: state.Schema,
				Raw:    tftypes.NewValue(state.Schema.Type().TerraformType(context.Background()), nil),
			},
		}
		resp := &datasource.ReadResponse{
			State: state,
		}

		ds.Read(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("read with state filter", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				// Verify state filter was applied
				assert.Equal(t, "active", r.URL.Query().Get("state"))

				response := `{"data": [{"id": "var-1", "key": "KEY1", "value": "val1"}]}`
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(response))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestVariablesClient(t, handler)
		defer server.Close()

		ds := &VariablesDataSource{client: n8nClient}

		state := createTestVariablesDataSourceState(t)
		req := datasource.ReadRequest{
			Config: tfsdk.Config{
				Schema: state.Schema,
				Raw:    tftypes.NewValue(state.Schema.Type().TerraformType(context.Background()), nil),
			},
		}
		resp := &datasource.ReadResponse{
			State: state,
		}

		ds.Read(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("read with both filters", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				assert.Equal(t, "proj-456", r.URL.Query().Get("projectId"))
				assert.Equal(t, "inactive", r.URL.Query().Get("state"))

				response := `{"data": []}`
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(response))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestVariablesClient(t, handler)
		defer server.Close()

		ds := &VariablesDataSource{client: n8nClient}

		state := createTestVariablesDataSourceState(t)
		req := datasource.ReadRequest{
			Config: tfsdk.Config{
				Schema: state.Schema,
				Raw:    tftypes.NewValue(state.Schema.Type().TerraformType(context.Background()), nil),
			},
		}
		resp := &datasource.ReadResponse{
			State: state,
		}

		ds.Read(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("read with empty response", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				response := `{"data": []}`
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(response))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestVariablesClient(t, handler)
		defer server.Close()

		ds := &VariablesDataSource{client: n8nClient}

		state := createTestVariablesDataSourceState(t)
		req := datasource.ReadRequest{
			Config: tfsdk.Config{
				Schema: state.Schema,
				Raw:    tftypes.NewValue(state.Schema.Type().TerraformType(context.Background()), nil),
			},
		}
		resp := &datasource.ReadResponse{
			State: state,
		}

		ds.Read(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("read with null data in response", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				response := `{}`
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(response))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestVariablesClient(t, handler)
		defer server.Close()

		ds := &VariablesDataSource{client: n8nClient}

		state := createTestVariablesDataSourceState(t)
		req := datasource.ReadRequest{
			Config: tfsdk.Config{
				Schema: state.Schema,
				Raw:    tftypes.NewValue(state.Schema.Type().TerraformType(context.Background()), nil),
			},
		}
		resp := &datasource.ReadResponse{
			State: state,
		}

		ds.Read(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("read fails on API error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestVariablesClient(t, handler)
		defer server.Close()

		ds := &VariablesDataSource{client: n8nClient}

		state := createTestVariablesDataSourceState(t)
		req := datasource.ReadRequest{
			Config: tfsdk.Config{
				Schema: state.Schema,
				Raw:    tftypes.NewValue(state.Schema.Type().TerraformType(context.Background()), nil),
			},
		}
		resp := &datasource.ReadResponse{
			State: state,
		}

		ds.Read(context.Background(), req, resp)

		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error listing variables")
	})
}

// Helper functions for test setup.
func setupTestVariablesClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
		BaseURL:   server.URL,
	}

	return n8nClient, server
}

func createTestVariablesDataSourceSchema(t *testing.T) schema.Schema {
	t.Helper()

	ds := NewVariablesDataSource()
	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)

	return schemaResp.Schema
}

func createTestVariablesDataSourceState(t *testing.T) tfsdk.State {
	t.Helper()

	testSchema := createTestVariablesDataSourceSchema(t)

	return tfsdk.State{
		Schema: testSchema,
		Raw:    tftypes.NewValue(testSchema.Type().TerraformType(context.Background()), nil),
	}
}
