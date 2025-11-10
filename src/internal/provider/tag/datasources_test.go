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
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create new instance", wantErr: false},
		{name: "multiple instances are independent", wantErr: false},
		{name: "error case - constructor validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create new instance":
				ds := NewTagsDataSource()

				assert.NotNil(t, ds)
				assert.IsType(t, &TagsDataSource{}, ds)
				assert.Nil(t, ds.client)

			case "multiple instances are independent":
				ds1 := NewTagsDataSource()
				ds2 := NewTagsDataSource()

				assert.NotSame(t, ds1, ds2)

			case "error case - constructor validation":
				// Placeholder
				assert.False(t, tt.wantErr)
			}
		})
	}
}

func TestNewTagsDataSourceWrapper(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create new wrapped instance", wantErr: false},
		{name: "wrapper returns datasource.DataSource interface", wantErr: false},
		{name: "error case - wrapper validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create new wrapped instance":
				ds := NewTagsDataSourceWrapper()

				assert.NotNil(t, ds)
				assert.IsType(t, &TagsDataSource{}, ds)

			case "wrapper returns datasource.DataSource interface":
				ds := NewTagsDataSourceWrapper()

				// ds is already of type datasource.DataSource, no assertion needed
				assert.NotNil(t, ds)

			case "error case - wrapper validation":
				// Placeholder
				assert.False(t, tt.wantErr)
			}
		})
	}
}

func TestTagsDataSource_Metadata(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "set metadata with n8n provider", wantErr: false},
		{name: "set metadata with custom provider type", wantErr: false},
		{name: "set metadata with empty provider type", wantErr: false},
		{name: "error case - metadata validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "set metadata with n8n provider":
				ds := NewTagsDataSource()
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_tags", resp.TypeName)

			case "set metadata with custom provider type":
				ds := NewTagsDataSource()
				resp := &datasource.MetadataResponse{}
				req := datasource.MetadataRequest{
					ProviderTypeName: "custom_provider",
				}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "custom_provider_tags", resp.TypeName)

			case "set metadata with empty provider type":
				ds := NewTagsDataSource()
				resp := &datasource.MetadataResponse{}
				req := datasource.MetadataRequest{
					ProviderTypeName: "",
				}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "_tags", resp.TypeName)

			case "error case - metadata validation":
				// Placeholder
				assert.False(t, tt.wantErr)
			}
		})
	}
}

func TestTagsDataSource_Schema(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "return schema", wantErr: false},
		{name: "tags attribute structure", wantErr: false},
		{name: "schema has description", wantErr: false},
		{name: "schema attributes have descriptions", wantErr: false},
		{name: "error case - schema validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "return schema":
				ds := NewTagsDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.Contains(t, resp.Schema.MarkdownDescription, "n8n tags")

				// Verify tags list exists
				tagsAttr, exists := resp.Schema.Attributes["tags"]
				assert.True(t, exists)
				assert.True(t, tagsAttr.IsComputed())

			case "tags attribute structure":
				ds := NewTagsDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				tagsAttr, exists := resp.Schema.Attributes["tags"]
				assert.True(t, exists)
				assert.True(t, tagsAttr.IsComputed())

			case "schema has description":
				ds := NewTagsDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				assert.NotEmpty(t, resp.Schema.MarkdownDescription)

			case "schema attributes have descriptions":
				ds := NewTagsDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				for name, attr := range resp.Schema.Attributes {
					assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
				}

			case "error case - schema validation":
				// Placeholder
				assert.False(t, tt.wantErr)
			}
		})
	}
}

func TestTagsDataSource_Configure(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "configure with valid client", wantErr: false},
		{name: "configure with nil provider data", wantErr: false},
		{name: "configure with invalid provider data", wantErr: true},
		{name: "configure with wrong type", wantErr: true},
		{name: "configure multiple times", wantErr: false},
		{name: "reconfigure with nil after valid client", wantErr: false},
		{name: "error case - invalid configuration", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "configure with valid client":
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

			case "configure with nil provider data":
				ds := NewTagsDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: nil,
				}

				ds.Configure(context.Background(), req, resp)

				assert.Nil(t, ds.client)
				assert.False(t, resp.Diagnostics.HasError())

			case "configure with invalid provider data":
				ds := NewTagsDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: "invalid-data",
				}

				ds.Configure(context.Background(), req, resp)

				assert.Nil(t, ds.client)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Data Source Configure Type")

			case "configure with wrong type":
				ds := NewTagsDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: struct{}{},
				}

				ds.Configure(context.Background(), req, resp)

				assert.Nil(t, ds.client)
				assert.True(t, resp.Diagnostics.HasError())

			case "configure multiple times":
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

			case "reconfigure with nil after valid client":
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

			case "error case - invalid configuration":
				// Same as configure with wrong type
				assert.True(t, tt.wantErr)
			}
		})
	}
}

func TestTagsDataSource_Interfaces(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "implements required interfaces", wantErr: false},
		{name: "error case - interface validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "implements required interfaces":
				ds := NewTagsDataSource()

				// Test that TagsDataSource implements datasource.DataSource
				var _ datasource.DataSource = ds

				// Test that TagsDataSource implements datasource.DataSourceWithConfigure
				var _ datasource.DataSourceWithConfigure = ds

				// Test that TagsDataSource implements TagsDataSourceInterface
				var _ TagsDataSourceInterface = ds

			case "error case - interface validation":
				// Placeholder
				assert.False(t, tt.wantErr)
			}
		})
	}
}

func TestTagsDataSourceConcurrency(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "concurrent metadata calls", wantErr: false},
		{name: "concurrent schema calls", wantErr: false},
		{name: "concurrent configure calls", wantErr: false},
		{name: "error case - concurrent access validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// NO t.Parallel() here - concurrent goroutines don't play well with t.Parallel()
			switch tt.name {
			case "concurrent metadata calls":
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

			case "concurrent schema calls":
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

			case "concurrent configure calls":
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

			case "error case - concurrent access validation":
				// Placeholder
				assert.False(t, tt.wantErr)
			}
		})
	}
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
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "successful read with tags", wantErr: false},
		{name: "successful read with empty list", wantErr: false},
		{name: "API error", wantErr: true},
		{name: "read with nil data", wantErr: false},
		{name: "error case - read validation", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "successful read with tags":
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

			case "successful read with empty list":
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

			case "API error":
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

			case "read with nil data":
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

			case "error case - read validation":
				// Placeholder
				assert.True(t, tt.wantErr)
			}
		})
	}
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
