package project

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

// MockProjectDataSourceInterface is a mock implementation of ProjectDataSourceInterface.
type MockProjectDataSourceInterface struct {
	mock.Mock
}

func (m *MockProjectDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockProjectDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockProjectDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockProjectDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

func TestNewProjectDataSource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		ds := NewProjectDataSource()

		assert.NotNil(t, ds)
		assert.IsType(t, &ProjectDataSource{}, ds)
		assert.Nil(t, ds.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		ds1 := NewProjectDataSource()
		ds2 := NewProjectDataSource()

		assert.NotSame(t, ds1, ds2)
		// Each instance is a different pointer, even if they have the same initial state
	})
}

func TestNewProjectDataSourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		ds := NewProjectDataSourceWrapper()

		assert.NotNil(t, ds)
		assert.IsType(t, &ProjectDataSource{}, ds)
	})

	t.Run("wrapper returns datasource.DataSource interface", func(t *testing.T) {
		ds := NewProjectDataSourceWrapper()

		// ds is already of type datasource.DataSource, no assertion needed
		assert.NotNil(t, ds)
	})
}

func TestProjectDataSource_Metadata(t *testing.T) {
	t.Run("set metadata without provider type", func(t *testing.T) {
		ds := NewProjectDataSource()
		resp := &datasource.MetadataResponse{}

		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_project", resp.TypeName)
	})

	t.Run("set metadata with provider type name", func(t *testing.T) {
		ds := NewProjectDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "n8n_project", resp.TypeName)
	})

	t.Run("set metadata with different provider type", func(t *testing.T) {
		ds := NewProjectDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_project", resp.TypeName)
	})
}

func TestProjectDataSource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		ds := NewProjectDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "n8n project")

		// Verify attributes
		expectedAttrs := []string{
			"id",
			"name",
			"type",
			"created_at",
			"updated_at",
			"icon",
			"description",
		}

		for _, attr := range expectedAttrs {
			_, exists := resp.Schema.Attributes[attr]
			assert.True(t, exists, "Attribute %s should exist", attr)
		}
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		ds := NewProjectDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})
}

func TestProjectDataSource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		ds := NewProjectDataSource()
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
		ds := NewProjectDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: nil,
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		ds := NewProjectDataSource()
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
		ds := NewProjectDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.True(t, resp.Diagnostics.HasError())
	})
}

func TestProjectDataSource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		ds := NewProjectDataSource()

		// Test that ProjectDataSource implements datasource.DataSource
		var _ datasource.DataSource = ds

		// Test that ProjectDataSource implements datasource.DataSourceWithConfigure
		var _ datasource.DataSourceWithConfigure = ds

		// Test that ProjectDataSource implements ProjectDataSourceInterface
		var _ ProjectDataSourceInterface = ds
	})
}

func TestProjectDataSourceConcurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		ds := NewProjectDataSource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)
				assert.Equal(t, "n8n_project", resp.TypeName)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		ds := NewProjectDataSource()

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

func BenchmarkProjectDataSource_Schema(b *testing.B) {
	ds := NewProjectDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkProjectDataSource_Metadata(b *testing.B) {
	ds := NewProjectDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)
	}
}

func BenchmarkProjectDataSource_Configure(b *testing.B) {
	ds := NewProjectDataSource()
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

func TestProjectDataSource_Read(t *testing.T) {
	t.Run("successful read by ID", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/projects" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []map[string]interface{}{
						{
							"id":          "proj-123",
							"name":        "Test Project",
							"type":        "team",
							"createdAt":   "2024-01-01T00:00:00Z",
							"updatedAt":   "2024-01-02T00:00:00Z",
							"icon":        "home",
							"description": "Test project description",
						},
					},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestDataSourceClient(t, handler)
		defer server.Close()

		ds := &ProjectDataSource{client: n8nClient}

		rawConfig := map[string]tftypes.Value{
			"id":          tftypes.NewValue(tftypes.String, "proj-123"),
			"name":        tftypes.NewValue(tftypes.String, nil),
			"type":        tftypes.NewValue(tftypes.String, nil),
			"created_at":  tftypes.NewValue(tftypes.String, nil),
			"updated_at":  tftypes.NewValue(tftypes.String, nil),
			"icon":        tftypes.NewValue(tftypes.String, nil),
			"description": tftypes.NewValue(tftypes.String, nil),
		}
		attrTypes := map[string]tftypes.Type{
			"id":          tftypes.String,
			"name":        tftypes.String,
			"type":        tftypes.String,
			"created_at":  tftypes.String,
			"updated_at":  tftypes.String,
			"icon":        tftypes.String,
			"description": tftypes.String,
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

	t.Run("successful read by name", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/projects" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []map[string]interface{}{
						{
							"id":          "proj-456",
							"name":        "Test Project By Name",
							"type":        "personal",
							"createdAt":   "2024-01-01T00:00:00Z",
							"updatedAt":   "2024-01-02T00:00:00Z",
							"icon":        "star",
							"description": "Another project",
						},
					},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestDataSourceClient(t, handler)
		defer server.Close()

		ds := &ProjectDataSource{client: n8nClient}

		rawConfig := map[string]tftypes.Value{
			"id":          tftypes.NewValue(tftypes.String, nil),
			"name":        tftypes.NewValue(tftypes.String, "Test Project By Name"),
			"type":        tftypes.NewValue(tftypes.String, nil),
			"created_at":  tftypes.NewValue(tftypes.String, nil),
			"updated_at":  tftypes.NewValue(tftypes.String, nil),
			"icon":        tftypes.NewValue(tftypes.String, nil),
			"description": tftypes.NewValue(tftypes.String, nil),
		}
		attrTypes := map[string]tftypes.Type{
			"id":          tftypes.String,
			"name":        tftypes.String,
			"type":        tftypes.String,
			"created_at":  tftypes.String,
			"updated_at":  tftypes.String,
			"icon":        tftypes.String,
			"description": tftypes.String,
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

	t.Run("project not found", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/projects" && r.Method == http.MethodGet {
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

		ds := &ProjectDataSource{client: n8nClient}

		rawConfig := map[string]tftypes.Value{
			"id":          tftypes.NewValue(tftypes.String, "proj-nonexistent"),
			"name":        tftypes.NewValue(tftypes.String, nil),
			"type":        tftypes.NewValue(tftypes.String, nil),
			"created_at":  tftypes.NewValue(tftypes.String, nil),
			"updated_at":  tftypes.NewValue(tftypes.String, nil),
			"icon":        tftypes.NewValue(tftypes.String, nil),
			"description": tftypes.NewValue(tftypes.String, nil),
		}
		attrTypes := map[string]tftypes.Type{
			"id":          tftypes.String,
			"name":        tftypes.String,
			"type":        tftypes.String,
			"created_at":  tftypes.String,
			"updated_at":  tftypes.String,
			"icon":        tftypes.String,
			"description": tftypes.String,
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

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors when project not found")
	})
}

// createTestDataSourceSchema creates a test schema for project datasource.
func createTestDataSourceSchema(t *testing.T) schema.Schema {
	t.Helper()
	ds := &ProjectDataSource{}
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
