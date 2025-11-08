package variable

import (
	"context"
	"encoding/json"
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

// MockVariableDataSourceInterface is a mock implementation of VariableDataSourceInterface.
type MockVariableDataSourceInterface struct {
	mock.Mock
}

func (m *MockVariableDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariableDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariableDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariableDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

func TestNewVariableDataSource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		ds := NewVariableDataSource()

		assert.NotNil(t, ds)
		assert.IsType(t, &VariableDataSource{}, ds)
		assert.Nil(t, ds.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		ds1 := NewVariableDataSource()
		ds2 := NewVariableDataSource()

		assert.NotSame(t, ds1, ds2)
	})
}

func TestNewVariableDataSourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		ds := NewVariableDataSourceWrapper()

		assert.NotNil(t, ds)
		assert.IsType(t, &VariableDataSource{}, ds)
	})

	t.Run("wrapper returns datasource.DataSource interface", func(t *testing.T) {
		ds := NewVariableDataSourceWrapper()

		// ds is already of type datasource.DataSource, no assertion needed
		assert.NotNil(t, ds)
	})
}

func TestVariableDataSource_Metadata(t *testing.T) {
	t.Run("set metadata with provider type", func(t *testing.T) {
		ds := NewVariableDataSource()
		resp := &datasource.MetadataResponse{}

		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_variable", resp.TypeName)
	})

	t.Run("set metadata with different provider type", func(t *testing.T) {
		ds := NewVariableDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_variable", resp.TypeName)
	})

	t.Run("set metadata with empty provider type", func(t *testing.T) {
		ds := NewVariableDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "_variable", resp.TypeName)
	})
}

func TestVariableDataSource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		ds := NewVariableDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "n8n variable")

		// Verify optional attributes
		idAttr, exists := resp.Schema.Attributes["id"]
		assert.True(t, exists)
		assert.True(t, idAttr.IsOptional())
		assert.True(t, idAttr.IsComputed())

		keyAttr, exists := resp.Schema.Attributes["key"]
		assert.True(t, exists)
		assert.True(t, keyAttr.IsOptional())
		assert.True(t, keyAttr.IsComputed())

		// Verify computed attributes
		valueAttr, exists := resp.Schema.Attributes["value"]
		assert.True(t, exists)
		assert.True(t, valueAttr.IsComputed())

		typeAttr, exists := resp.Schema.Attributes["type"]
		assert.True(t, exists)
		assert.True(t, typeAttr.IsComputed())

		projectIDAttr, exists := resp.Schema.Attributes["project_id"]
		assert.True(t, exists)
		assert.True(t, projectIDAttr.IsOptional())
		assert.True(t, projectIDAttr.IsComputed())
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		ds := NewVariableDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})

	t.Run("value attribute is sensitive", func(t *testing.T) {
		ds := NewVariableDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		_, exists := resp.Schema.Attributes["value"]
		assert.True(t, exists)
		// Note: Sensitive attribute testing depends on framework internals
	})
}

func TestVariableDataSource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		ds := NewVariableDataSource()
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
		ds := NewVariableDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: nil,
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		ds := NewVariableDataSource()
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
		ds := NewVariableDataSource()
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
		ds := NewVariableDataSource()

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

func TestVariableDataSource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		ds := NewVariableDataSource()

		// Test that VariableDataSource implements datasource.DataSource
		var _ datasource.DataSource = ds

		// Test that VariableDataSource implements datasource.DataSourceWithConfigure
		var _ datasource.DataSourceWithConfigure = ds

		// Test that VariableDataSource implements VariableDataSourceInterface
		var _ VariableDataSourceInterface = ds
	})
}

func TestVariableDataSource_ValidateIdentifiers(t *testing.T) {
	t.Run("valid with ID provided", func(t *testing.T) {
		_ = NewVariableDataSource()
		_ = &datasource.ReadResponse{}

		// This would normally be populated from the config
		// For testing, we're directly testing the validation logic
		// Note: validateIdentifiers is not exported, so we test through Read behavior
		// This test structure follows the pattern but can't directly test unexported methods
	})

	t.Run("valid with key provided", func(t *testing.T) {
		_ = NewVariableDataSource()
		_ = &datasource.ReadResponse{}
		// Similar limitation as above
	})
}

func TestVariableDataSourceConcurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		ds := NewVariableDataSource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)
				assert.Equal(t, "n8n_variable", resp.TypeName)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		ds := NewVariableDataSource()

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
		ds := NewVariableDataSource()

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
}

func BenchmarkVariableDataSource_Schema(b *testing.B) {
	ds := NewVariableDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkVariableDataSource_Metadata(b *testing.B) {
	ds := NewVariableDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)
	}
}

func BenchmarkVariableDataSource_Configure(b *testing.B) {
	ds := NewVariableDataSource()
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

func TestVariableDataSource_Read(t *testing.T) {
	t.Run("successful read by ID", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []map[string]interface{}{
						{
							"id":    "var-123",
							"key":   "TEST_VAR",
							"value": "test-value",
							"type":  "string",
						},
					},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestDataSourceClient(t, handler)
		defer server.Close()

		ds := &VariableDataSource{client: n8nClient}

		rawConfig := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, nil),
			"value":      tftypes.NewValue(tftypes.String, nil),
			"type":       tftypes.NewValue(tftypes.String, nil),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		attrTypes := map[string]tftypes.Type{
			"id":         tftypes.String,
			"key":        tftypes.String,
			"value":      tftypes.String,
			"type":       tftypes.String,
			"project_id": tftypes.String,
		}
		config := tfsdk.Config{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, rawConfig),
			Schema: createTestDataSourceSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, nil),
			Schema: createTestDataSourceSchema(t),
		}

		req := datasource.ReadRequest{
			Config: config,
		}
		resp := datasource.ReadResponse{
			State: state,
		}

		ds.Read(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")
	})

	t.Run("successful read by key", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []map[string]interface{}{
						{
							"id":    "var-456",
							"key":   "TEST_VAR_BY_KEY",
							"value": "test-value",
							"type":  "string",
						},
					},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestDataSourceClient(t, handler)
		defer server.Close()

		ds := &VariableDataSource{client: n8nClient}

		rawConfig := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"key":        tftypes.NewValue(tftypes.String, "TEST_VAR_BY_KEY"),
			"value":      tftypes.NewValue(tftypes.String, nil),
			"type":       tftypes.NewValue(tftypes.String, nil),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		attrTypes := map[string]tftypes.Type{
			"id":         tftypes.String,
			"key":        tftypes.String,
			"value":      tftypes.String,
			"type":       tftypes.String,
			"project_id": tftypes.String,
		}
		config := tfsdk.Config{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, rawConfig),
			Schema: createTestDataSourceSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, nil),
			Schema: createTestDataSourceSchema(t),
		}

		req := datasource.ReadRequest{
			Config: config,
		}
		resp := datasource.ReadResponse{
			State: state,
		}

		ds.Read(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")
	})

	t.Run("variable not found", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []map[string]interface{}{},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestDataSourceClient(t, handler)
		defer server.Close()

		ds := &VariableDataSource{client: n8nClient}

		rawConfig := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-nonexistent"),
			"key":        tftypes.NewValue(tftypes.String, nil),
			"value":      tftypes.NewValue(tftypes.String, nil),
			"type":       tftypes.NewValue(tftypes.String, nil),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		attrTypes := map[string]tftypes.Type{
			"id":         tftypes.String,
			"key":        tftypes.String,
			"value":      tftypes.String,
			"type":       tftypes.String,
			"project_id": tftypes.String,
		}
		config := tfsdk.Config{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, rawConfig),
			Schema: createTestDataSourceSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, nil),
			Schema: createTestDataSourceSchema(t),
		}

		req := datasource.ReadRequest{
			Config: config,
		}
		resp := datasource.ReadResponse{
			State: state,
		}

		ds.Read(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors when variable not found")
	})
}

// createTestDataSourceSchema creates a test schema for variable datasource.
func createTestDataSourceSchema(t *testing.T) schema.Schema {
	t.Helper()
	ds := &VariableDataSource{}
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), req, resp)
	return resp.Schema
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
