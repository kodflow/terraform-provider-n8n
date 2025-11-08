package workflow

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

// MockWorkflowDataSourceInterface is a mock implementation of WorkflowDataSourceInterface
type MockWorkflowDataSourceInterface struct {
	mock.Mock
}

func (m *MockWorkflowDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockWorkflowDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockWorkflowDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockWorkflowDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

func TestNewWorkflowDataSource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		ds := NewWorkflowDataSource()

		assert.NotNil(t, ds)
		assert.IsType(t, &WorkflowDataSource{}, ds)
		assert.Nil(t, ds.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		ds1 := NewWorkflowDataSource()
		ds2 := NewWorkflowDataSource()

		assert.NotSame(t, ds1, ds2)
	})
}

func TestNewWorkflowDataSourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		ds := NewWorkflowDataSourceWrapper()

		assert.NotNil(t, ds)
		assert.IsType(t, &WorkflowDataSource{}, ds)
	})

	t.Run("wrapper returns datasource.DataSource interface", func(t *testing.T) {
		ds := NewWorkflowDataSourceWrapper()

		_, ok := ds.(datasource.DataSource)
		assert.True(t, ok)
	})
}

func TestWorkflowDataSource_Metadata(t *testing.T) {
	t.Run("set metadata with provider type", func(t *testing.T) {
		ds := NewWorkflowDataSource()
		resp := &datasource.MetadataResponse{}

		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_workflow", resp.TypeName)
	})

	t.Run("set metadata with different provider type", func(t *testing.T) {
		ds := NewWorkflowDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_workflow", resp.TypeName)
	})

	t.Run("set metadata with empty provider type", func(t *testing.T) {
		ds := NewWorkflowDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "_workflow", resp.TypeName)
	})
}

func TestWorkflowDataSource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		ds := NewWorkflowDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "n8n workflow")

		// Verify required attributes
		idAttr, exists := resp.Schema.Attributes["id"]
		assert.True(t, exists)
		assert.True(t, idAttr.IsRequired())

		// Verify computed attributes
		nameAttr, exists := resp.Schema.Attributes["name"]
		assert.True(t, exists)
		assert.True(t, nameAttr.IsComputed())

		activeAttr, exists := resp.Schema.Attributes["active"]
		assert.True(t, exists)
		assert.True(t, activeAttr.IsComputed())
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		ds := NewWorkflowDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})

	t.Run("schema has correct attribute count", func(t *testing.T) {
		ds := NewWorkflowDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		// Should have id, name, active
		assert.Len(t, resp.Schema.Attributes, 3)
	})
}

func TestWorkflowDataSource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		ds := NewWorkflowDataSource()
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
		ds := NewWorkflowDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: nil,
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		ds := NewWorkflowDataSource()
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
		ds := NewWorkflowDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("configure multiple times", func(t *testing.T) {
		ds := NewWorkflowDataSource()

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

func TestWorkflowDataSource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		ds := NewWorkflowDataSource()

		// Test that WorkflowDataSource implements datasource.DataSource
		var _ datasource.DataSource = ds

		// Test that WorkflowDataSource implements datasource.DataSourceWithConfigure
		var _ datasource.DataSourceWithConfigure = ds

		// Test that WorkflowDataSource implements WorkflowDataSourceInterface
		var _ WorkflowDataSourceInterface = ds
	})
}

func TestWorkflowDataSourceConcurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		ds := NewWorkflowDataSource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)
				assert.Equal(t, "n8n_workflow", resp.TypeName)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		ds := NewWorkflowDataSource()

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

func BenchmarkWorkflowDataSource_Schema(b *testing.B) {
	ds := NewWorkflowDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkWorkflowDataSource_Metadata(b *testing.B) {
	ds := NewWorkflowDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{}, resp)
	}
}

func BenchmarkWorkflowDataSource_Configure(b *testing.B) {
	ds := NewWorkflowDataSource()
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

func TestWorkflowDataSource_Read(t *testing.T) {
	t.Run("successful read", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/workflows/wf-123" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":          "wf-123",
					"name":        "Test Workflow",
					"active":      true,
					"nodes":       []interface{}{},
					"connections": map[string]interface{}{},
					"settings":    map[string]interface{}{},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestDataSourceClient(t, handler)
		defer server.Close()

		ds := &WorkflowDataSource{client: n8nClient}

		rawConfig := map[string]tftypes.Value{
			"id":     tftypes.NewValue(tftypes.String, "wf-123"),
			"name":   tftypes.NewValue(tftypes.String, nil),
			"active": tftypes.NewValue(tftypes.Bool, nil),
		}
		attrTypes := map[string]tftypes.Type{
			"id":     tftypes.String,
			"name":   tftypes.String,
			"active": tftypes.Bool,
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

	t.Run("workflow not found", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "Workflow not found"}`))
		})

		n8nClient, server := setupTestDataSourceClient(t, handler)
		defer server.Close()

		ds := &WorkflowDataSource{client: n8nClient}

		rawConfig := map[string]tftypes.Value{
			"id":     tftypes.NewValue(tftypes.String, "wf-nonexistent"),
			"name":   tftypes.NewValue(tftypes.String, nil),
			"active": tftypes.NewValue(tftypes.Bool, nil),
		}
		attrTypes := map[string]tftypes.Type{
			"id":     tftypes.String,
			"name":   tftypes.String,
			"active": tftypes.Bool,
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

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors when workflow not found")
	})
}

// createTestDataSourceSchema creates a test schema for workflow datasource
func createTestDataSourceSchema(t *testing.T) schema.Schema {
	ds := &WorkflowDataSource{}
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), req, resp)
	return resp.Schema
}

// setupTestDataSourceClient creates a test N8nClient with httptest server
func setupTestDataSourceClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
