package workflow

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWorkflowsDataSourceInterface is a mock implementation of WorkflowsDataSourceInterface.
type MockWorkflowsDataSourceInterface struct {
	mock.Mock
}

func (m *MockWorkflowsDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockWorkflowsDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockWorkflowsDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockWorkflowsDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

func TestNewWorkflowsDataSource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		ds := NewWorkflowsDataSource()

		assert.NotNil(t, ds)
		assert.IsType(t, &WorkflowsDataSource{}, ds)
		assert.Nil(t, ds.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		ds1 := NewWorkflowsDataSource()
		ds2 := NewWorkflowsDataSource()

		assert.NotSame(t, ds1, ds2)
	})
}

func TestNewWorkflowsDataSourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		ds := NewWorkflowsDataSourceWrapper()

		assert.NotNil(t, ds)
		assert.IsType(t, &WorkflowsDataSource{}, ds)
	})

	t.Run("wrapper returns datasource.DataSource interface", func(t *testing.T) {
		ds := NewWorkflowsDataSourceWrapper()

		// ds is already of type datasource.DataSource, no assertion needed
		assert.NotNil(t, ds)
	})
}

func TestWorkflowsDataSource_Metadata(t *testing.T) {
	t.Run("set metadata with provider type", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
		resp := &datasource.MetadataResponse{}

		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_workflows", resp.TypeName)
	})

	t.Run("set metadata with different provider type", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_workflows", resp.TypeName)
	})

	t.Run("set metadata with empty provider type", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "_workflows", resp.TypeName)
	})
}

func TestWorkflowsDataSource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "n8n workflows")

		// Verify filter attributes
		activeAttr, exists := resp.Schema.Attributes["active"]
		assert.True(t, exists)
		assert.True(t, activeAttr.IsOptional())

		// Verify computed workflows list
		workflowsAttr, exists := resp.Schema.Attributes["workflows"]
		assert.True(t, exists)
		assert.True(t, workflowsAttr.IsComputed())
	})

	t.Run("workflows list has correct structure", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		workflowsAttr, exists := resp.Schema.Attributes["workflows"]
		assert.True(t, exists)
		assert.True(t, workflowsAttr.IsComputed())
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})

	t.Run("schema has correct attribute count", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		// Should have active and workflows
		assert.Len(t, resp.Schema.Attributes, 2)
	})
}

func TestWorkflowsDataSource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
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
		ds := NewWorkflowsDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: nil,
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
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
		ds := NewWorkflowsDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("configure multiple times", func(t *testing.T) {
		ds := NewWorkflowsDataSource()

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

func TestWorkflowsDataSource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		ds := NewWorkflowsDataSource()

		// Test that WorkflowsDataSource implements datasource.DataSource
		var _ datasource.DataSource = ds

		// Test that WorkflowsDataSource implements datasource.DataSourceWithConfigure
		var _ datasource.DataSourceWithConfigure = ds

		// Test that WorkflowsDataSource implements WorkflowsDataSourceInterface
		var _ WorkflowsDataSourceInterface = ds
	})
}

func TestWorkflowsDataSourceConcurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		ds := NewWorkflowsDataSource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)
				assert.Equal(t, "n8n_workflows", resp.TypeName)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		ds := NewWorkflowsDataSource()

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

func BenchmarkWorkflowsDataSource_Schema(b *testing.B) {
	ds := NewWorkflowsDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkWorkflowsDataSource_Metadata(b *testing.B) {
	ds := NewWorkflowsDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{}, resp)
	}
}

func BenchmarkWorkflowsDataSource_Configure(b *testing.B) {
	ds := NewWorkflowsDataSource()
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

func TestWorkflowsDataSource_Read(t *testing.T) {
	t.Run("successful read all workflows", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/workflows" && r.Method == http.MethodGet {
				response := map[string]interface{}{
					"data": []interface{}{
						map[string]interface{}{
							"id":          "wf-1",
							"name":        "Workflow 1",
							"active":      true,
							"nodes":       []interface{}{},
							"connections": map[string]interface{}{},
							"settings":    map[string]interface{}{},
						},
						map[string]interface{}{
							"id":          "wf-2",
							"name":        "Workflow 2",
							"active":      false,
							"nodes":       []interface{}{},
							"connections": map[string]interface{}{},
							"settings":    map[string]interface{}{},
						},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupDataSourceTestClient(t, handler)
		defer server.Close()

		ds := &WorkflowsDataSource{client: n8nClient}

		dsSchema := createDataSourceTestSchema(t)

		rawConfig := map[string]tftypes.Value{
			"active":    tftypes.NewValue(tftypes.Bool, nil),
			"workflows": tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "active": tftypes.Bool}}}, nil),
		}
		config := tfsdk.Config{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"active": tftypes.Bool, "workflows": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "active": tftypes.Bool}}}}}, rawConfig),
			Schema: dsSchema,
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"active": tftypes.Bool, "workflows": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "active": tftypes.Bool}}}}}, nil),
			Schema: dsSchema,
		}

		req := datasource.ReadRequest{
			Config: config,
		}
		resp := datasource.ReadResponse{
			State: state,
		}

		ds.Read(context.Background(), req, &resp)

		if resp.Diagnostics.HasError() {
			for _, diag := range resp.Diagnostics.Errors() {
				t.Logf("Diagnostic Error: %s - %s", diag.Summary(), diag.Detail())
			}
		}
		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")
	})

	t.Run("successful read with active filter", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/workflows" && r.Method == http.MethodGet {
				// Verify active filter is present
				activeParam := r.URL.Query().Get("active")
				assert.Equal(t, "true", activeParam)

				response := map[string]interface{}{
					"data": []interface{}{
						map[string]interface{}{
							"id":          "wf-1",
							"name":        "Active Workflow",
							"active":      true,
							"nodes":       []interface{}{},
							"connections": map[string]interface{}{},
							"settings":    map[string]interface{}{},
						},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupDataSourceTestClient(t, handler)
		defer server.Close()

		ds := &WorkflowsDataSource{client: n8nClient}

		rawConfig := map[string]tftypes.Value{
			"active":    tftypes.NewValue(tftypes.Bool, true),
			"workflows": tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "active": tftypes.Bool}}}, nil),
		}
		config := tfsdk.Config{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"active": tftypes.Bool, "workflows": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "active": tftypes.Bool}}}}}, rawConfig),
			Schema: createDataSourceTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"active": tftypes.Bool, "workflows": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "active": tftypes.Bool}}}}}, nil),
			Schema: createDataSourceTestSchema(t),
		}

		req := datasource.ReadRequest{
			Config: config,
		}
		resp := datasource.ReadResponse{
			State: state,
		}

		ds.Read(context.Background(), req, &resp)

		if resp.Diagnostics.HasError() {
			for _, diag := range resp.Diagnostics.Errors() {
				t.Logf("Diagnostic Error: %s - %s", diag.Summary(), diag.Detail())
			}
		}
		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")
	})

	t.Run("read fails with API error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupDataSourceTestClient(t, handler)
		defer server.Close()

		ds := &WorkflowsDataSource{client: n8nClient}

		dsSchema := createDataSourceTestSchema(t)

		rawConfig := map[string]tftypes.Value{
			"active":    tftypes.NewValue(tftypes.Bool, nil),
			"workflows": tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "active": tftypes.Bool}}}, nil),
		}
		config := tfsdk.Config{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"active": tftypes.Bool, "workflows": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "active": tftypes.Bool}}}}}, rawConfig),
			Schema: dsSchema,
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"active": tftypes.Bool, "workflows": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "active": tftypes.Bool}}}}}, nil),
			Schema: dsSchema,
		}

		req := datasource.ReadRequest{
			Config: config,
		}
		resp := datasource.ReadResponse{
			State: state,
		}

		ds.Read(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error reading workflows")
	})

	t.Run("read with empty workflow list", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/workflows" && r.Method == http.MethodGet {
				response := map[string]interface{}{
					"data": []interface{}{},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupDataSourceTestClient(t, handler)
		defer server.Close()

		ds := &WorkflowsDataSource{client: n8nClient}

		dsSchema := createDataSourceTestSchema(t)

		rawConfig := map[string]tftypes.Value{
			"active":    tftypes.NewValue(tftypes.Bool, nil),
			"workflows": tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "active": tftypes.Bool}}}, nil),
		}
		config := tfsdk.Config{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"active": tftypes.Bool, "workflows": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "active": tftypes.Bool}}}}}, rawConfig),
			Schema: dsSchema,
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"active": tftypes.Bool, "workflows": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "active": tftypes.Bool}}}}}, nil),
			Schema: dsSchema,
		}

		req := datasource.ReadRequest{
			Config: config,
		}
		resp := datasource.ReadResponse{
			State: state,
		}

		ds.Read(context.Background(), req, &resp)

		if resp.Diagnostics.HasError() {
			for _, diag := range resp.Diagnostics.Errors() {
				t.Logf("Diagnostic Error: %s - %s", diag.Summary(), diag.Detail())
			}
		}
		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")
	})
}

// createDataSourceTestSchema creates a test schema for workflows data source.
func createDataSourceTestSchema(t *testing.T) dsschema.Schema {
	t.Helper()
	ds := NewWorkflowsDataSource()
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), req, resp)
	return resp.Schema
}

// setupDataSourceTestClient creates a test N8nClient with httptest server for data source tests.
func setupDataSourceTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
