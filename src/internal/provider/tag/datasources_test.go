package tag

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTagsDataSourceInterface is a mock implementation of TagsDataSourceInterface.
type MockTagsDataSourceInterface struct {
	mock.Mock
}

func (m *MockTagsDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockTagsDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockTagsDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockTagsDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

func TestNewTagsDataSource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		ds := NewTagsDataSource()

		assert.NotNil(t, ds)
		assert.IsType(t, &TagsDataSource{}, ds)
		assert.Nil(t, ds.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		ds1 := NewTagsDataSource()
		ds2 := NewTagsDataSource()

		assert.NotSame(t, ds1, ds2)
	})
}

func TestNewTagsDataSourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		ds := NewTagsDataSourceWrapper()

		assert.NotNil(t, ds)
		assert.IsType(t, &TagsDataSource{}, ds)
	})

	t.Run("wrapper returns datasource.DataSource interface", func(t *testing.T) {
		ds := NewTagsDataSourceWrapper()

		// ds is already of type datasource.DataSource, no assertion needed
		assert.NotNil(t, ds)
	})
}

func TestTagsDataSource_Metadata(t *testing.T) {
	t.Run("set metadata with n8n provider", func(t *testing.T) {
		ds := NewTagsDataSource()
		resp := &datasource.MetadataResponse{}

		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_tags", resp.TypeName)
	})

	t.Run("set metadata with custom provider type", func(t *testing.T) {
		ds := NewTagsDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_tags", resp.TypeName)
	})

	t.Run("set metadata with empty provider type", func(t *testing.T) {
		ds := NewTagsDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "_tags", resp.TypeName)
	})
}

func TestTagsDataSource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		ds := NewTagsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "n8n tags")

		// Verify tags list exists
		tagsAttr, exists := resp.Schema.Attributes["tags"]
		assert.True(t, exists)
		assert.True(t, tagsAttr.IsComputed())
	})

	t.Run("tags attribute structure", func(t *testing.T) {
		ds := NewTagsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		tagsAttr, exists := resp.Schema.Attributes["tags"]
		assert.True(t, exists)
		assert.True(t, tagsAttr.IsComputed())
	})

	t.Run("schema has description", func(t *testing.T) {
		ds := NewTagsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		assert.NotEmpty(t, resp.Schema.MarkdownDescription)
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		ds := NewTagsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})
}

func TestTagsDataSource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		ds := NewTagsDataSource()
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
		ds := NewTagsDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: nil,
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		ds := NewTagsDataSource()
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
		ds := NewTagsDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("configure multiple times", func(t *testing.T) {
		ds := NewTagsDataSource()

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

	t.Run("reconfigure with nil after valid client", func(t *testing.T) {
		ds := NewTagsDataSource()

		// First configuration with valid client
		resp1 := &datasource.ConfigureResponse{}
		client1 := &client.N8nClient{}
		req1 := datasource.ConfigureRequest{
			ProviderData: client1,
		}
		ds.Configure(context.Background(), req1, resp1)
		assert.NotNil(t, ds.client)

		// Reconfigure with nil
		resp2 := &datasource.ConfigureResponse{}
		req2 := datasource.ConfigureRequest{
			ProviderData: nil,
		}
		ds.Configure(context.Background(), req2, resp2)
		assert.NotNil(t, ds.client) // Should remain unchanged
	})
}

func TestTagsDataSource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		ds := NewTagsDataSource()

		// Test that TagsDataSource implements datasource.DataSource
		var _ datasource.DataSource = ds

		// Test that TagsDataSource implements datasource.DataSourceWithConfigure
		var _ datasource.DataSourceWithConfigure = ds

		// Test that TagsDataSource implements TagsDataSourceInterface
		var _ TagsDataSourceInterface = ds
	})
}

func TestTagsDataSourceConcurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		ds := NewTagsDataSource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)
				assert.Equal(t, "n8n_tags", resp.TypeName)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		ds := NewTagsDataSource()

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
		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				ds := NewTagsDataSource()
				resp := &datasource.ConfigureResponse{}
				mockClient := &client.N8nClient{}
				req := datasource.ConfigureRequest{
					ProviderData: mockClient,
				}
				ds.Configure(context.Background(), req, resp)
				assert.NotNil(t, ds.client)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func BenchmarkTagsDataSource_Schema(b *testing.B) {
	ds := NewTagsDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkTagsDataSource_Metadata(b *testing.B) {
	ds := NewTagsDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)
	}
}

func BenchmarkTagsDataSource_Configure(b *testing.B) {
	ds := NewTagsDataSource()
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

// TestTagsDataSource_Read tests the Read method.
func TestTagsDataSource_Read(t *testing.T) {
	t.Run("successful read with tags", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/tags" && r.Method == http.MethodGet {
				tags := map[string]interface{}{
					"data": []map[string]interface{}{
						{
							"id":        "tag-123",
							"name":      "Test Tag 1",
							"createdAt": "2024-01-01T00:00:00Z",
							"updatedAt": "2024-01-02T00:00:00Z",
						},
						{
							"id":        "tag-456",
							"name":      "Test Tag 2",
							"createdAt": "2024-01-03T00:00:00Z",
							"updatedAt": "2024-01-04T00:00:00Z",
						},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(tags)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestDataSourcesClient(t, handler)
		defer server.Close()

		ds := &TagsDataSource{client: n8nClient}

		ctx := context.Background()
		schemaResp := datasource.SchemaResponse{}
		ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

		state := createEmptyTagsState(t, schemaResp.Schema)

		req := datasource.ReadRequest{}
		resp := &datasource.ReadResponse{
			State: state,
		}

		ds.Read(ctx, req, resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")
	})

	t.Run("successful read with empty list", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/tags" && r.Method == http.MethodGet {
				tags := map[string]interface{}{
					"data": []map[string]interface{}{},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(tags)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestDataSourcesClient(t, handler)
		defer server.Close()

		ds := &TagsDataSource{client: n8nClient}

		ctx := context.Background()
		schemaResp := datasource.SchemaResponse{}
		ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

		state := createEmptyTagsState(t, schemaResp.Schema)

		req := datasource.ReadRequest{}
		resp := &datasource.ReadResponse{
			State: state,
		}

		ds.Read(ctx, req, resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors with empty list")
	})

	t.Run("API error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestDataSourcesClient(t, handler)
		defer server.Close()

		ds := &TagsDataSource{client: n8nClient}

		ctx := context.Background()
		schemaResp := datasource.SchemaResponse{}
		ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

		state := createEmptyTagsState(t, schemaResp.Schema)

		req := datasource.ReadRequest{}
		resp := &datasource.ReadResponse{
			State: state,
		}

		ds.Read(ctx, req, resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors on API error")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error listing tags")
	})

	t.Run("read with nil data", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/tags" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestDataSourcesClient(t, handler)
		defer server.Close()

		ds := &TagsDataSource{client: n8nClient}

		ctx := context.Background()
		schemaResp := datasource.SchemaResponse{}
		ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

		state := createEmptyTagsState(t, schemaResp.Schema)

		req := datasource.ReadRequest{}
		resp := &datasource.ReadResponse{
			State: state,
		}

		ds.Read(ctx, req, resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors with nil data")
	})
}

// createEmptyTagsState creates an empty state for tags datasource tests.
func createEmptyTagsState(t *testing.T, schema schema.Schema) tfsdk.State {
	t.Helper()
	return tfsdk.State{
		Schema: schema,
	}
}

// setupTestDataSourcesClient creates a test N8nClient with httptest server for datasources.
func setupTestDataSourcesClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
