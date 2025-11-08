package execution

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

// MockExecutionDataSourceInterface is a mock implementation of ExecutionDataSourceInterface.
type MockExecutionDataSourceInterface struct {
	mock.Mock
}

func (m *MockExecutionDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

func TestNewExecutionDataSource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		ds := NewExecutionDataSource()

		assert.NotNil(t, ds)
		assert.IsType(t, &ExecutionDataSource{}, ds)
		assert.Nil(t, ds.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		ds1 := NewExecutionDataSource()
		ds2 := NewExecutionDataSource()

		assert.NotSame(t, ds1, ds2)
		// Each instance is a different pointer, even if they have the same initial state
	})
}

func TestNewExecutionDataSourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		ds := NewExecutionDataSourceWrapper()

		assert.NotNil(t, ds)
		assert.IsType(t, &ExecutionDataSource{}, ds)
	})

	t.Run("wrapper returns datasource.DataSource interface", func(t *testing.T) {
		ds := NewExecutionDataSourceWrapper()

		// ds is already of type datasource.DataSource, no assertion needed
		assert.NotNil(t, ds)
	})
}

func TestExecutionDataSource_Metadata(t *testing.T) {
	t.Run("set metadata without provider type", func(t *testing.T) {
		ds := NewExecutionDataSource()
		resp := &datasource.MetadataResponse{}

		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_execution", resp.TypeName)
	})

	t.Run("set metadata with provider type name", func(t *testing.T) {
		ds := NewExecutionDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "n8n_execution", resp.TypeName)
	})

	t.Run("set metadata with different provider type", func(t *testing.T) {
		ds := NewExecutionDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_execution", resp.TypeName)
	})
}

func TestExecutionDataSource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		ds := NewExecutionDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "n8n workflow execution")

		// Verify required attributes
		idAttr, exists := resp.Schema.Attributes["id"]
		assert.True(t, exists)
		assert.True(t, idAttr.IsRequired())

		// Verify computed attributes
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

		// Verify optional attributes
		includeDataAttr, exists := resp.Schema.Attributes["include_data"]
		assert.True(t, exists)
		assert.True(t, includeDataAttr.IsOptional())
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		ds := NewExecutionDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})
}

func TestExecutionDataSource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		ds := NewExecutionDataSource()
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
		ds := NewExecutionDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: nil,
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		ds := NewExecutionDataSource()
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
		ds := NewExecutionDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("configure multiple times", func(t *testing.T) {
		ds := NewExecutionDataSource()

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

func TestExecutionDataSource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		ds := NewExecutionDataSource()

		// Test that ExecutionDataSource implements datasource.DataSource
		var _ datasource.DataSource = ds

		// Test that ExecutionDataSource implements datasource.DataSourceWithConfigure
		var _ datasource.DataSourceWithConfigure = ds

		// Test that ExecutionDataSource implements ExecutionDataSourceInterface
		var _ ExecutionDataSourceInterface = ds
	})
}

func TestExecutionDataSourceConcurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		ds := NewExecutionDataSource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)
				assert.Equal(t, "n8n_execution", resp.TypeName)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		ds := NewExecutionDataSource()

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
}

func BenchmarkExecutionDataSource_Schema(b *testing.B) {
	ds := NewExecutionDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkExecutionDataSource_Metadata(b *testing.B) {
	ds := NewExecutionDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{}, resp)
	}
}

func BenchmarkExecutionDataSource_Configure(b *testing.B) {
	ds := NewExecutionDataSource()
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

func TestExecutionDataSource_Read(t *testing.T) {
	t.Run("successful read", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/executions/123" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":         123,
					"workflowId": 456,
					"finished":   true,
					"mode":       "manual",
					"startedAt":  "2024-01-01T00:00:00Z",
					"stoppedAt":  "2024-01-01T00:01:00Z",
					"status":     "success",
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestDataSourceClient(t, handler)
		defer server.Close()

		ds := &ExecutionDataSource{client: n8nClient}

		rawConfig := map[string]tftypes.Value{
			"id":           tftypes.NewValue(tftypes.String, "123"),
			"workflow_id":  tftypes.NewValue(tftypes.String, nil),
			"finished":     tftypes.NewValue(tftypes.Bool, nil),
			"mode":         tftypes.NewValue(tftypes.String, nil),
			"started_at":   tftypes.NewValue(tftypes.String, nil),
			"stopped_at":   tftypes.NewValue(tftypes.String, nil),
			"status":       tftypes.NewValue(tftypes.String, nil),
			"include_data": tftypes.NewValue(tftypes.Bool, nil),
		}
		attrTypes := map[string]tftypes.Type{
			"id":           tftypes.String,
			"workflow_id":  tftypes.String,
			"finished":     tftypes.Bool,
			"mode":         tftypes.String,
			"started_at":   tftypes.String,
			"stopped_at":   tftypes.String,
			"status":       tftypes.String,
			"include_data": tftypes.Bool,
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

	t.Run("execution not found", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "Execution not found"}`))
		})

		n8nClient, server := setupTestDataSourceClient(t, handler)
		defer server.Close()

		ds := &ExecutionDataSource{client: n8nClient}

		rawConfig := map[string]tftypes.Value{
			"id":           tftypes.NewValue(tftypes.String, "99999"),
			"workflow_id":  tftypes.NewValue(tftypes.String, nil),
			"finished":     tftypes.NewValue(tftypes.Bool, nil),
			"mode":         tftypes.NewValue(tftypes.String, nil),
			"started_at":   tftypes.NewValue(tftypes.String, nil),
			"stopped_at":   tftypes.NewValue(tftypes.String, nil),
			"status":       tftypes.NewValue(tftypes.String, nil),
			"include_data": tftypes.NewValue(tftypes.Bool, nil),
		}
		attrTypes := map[string]tftypes.Type{
			"id":           tftypes.String,
			"workflow_id":  tftypes.String,
			"finished":     tftypes.Bool,
			"mode":         tftypes.String,
			"started_at":   tftypes.String,
			"stopped_at":   tftypes.String,
			"status":       tftypes.String,
			"include_data": tftypes.Bool,
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

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors when execution not found")
	})
}

// createTestDataSourceSchema creates a test schema for execution datasource.
func createTestDataSourceSchema(t *testing.T) schema.Schema {
	t.Helper()
	ds := &ExecutionDataSource{}
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
