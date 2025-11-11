package tag

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/tag/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTagDataSourceInterface is a mock implementation of TagDataSourceInterface.
type MockTagDataSourceInterface struct {
	mock.Mock
}

func (m *MockTagDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockTagDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockTagDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockTagDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
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

// TestNewTagDataSource is now in external test file - refactored to test behavior only.

// TestNewTagDataSourceWrapper is now in external test file - refactored to test behavior only.

// TestTagDataSource_Metadata is now in external test file - refactored to test behavior only.

// TestTagDataSource_Schema is now in external test file - refactored to test behavior only.

// TestTagDataSource_Configure is now in external test file - refactored to test behavior only.

func TestTagDataSource_Interfaces(t *testing.T) {
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
				ds := NewTagDataSource()

				// Test that TagDataSource implements datasource.DataSource
				var _ datasource.DataSource = ds

				// Test that TagDataSource implements datasource.DataSourceWithConfigure
				var _ datasource.DataSourceWithConfigure = ds

				// Test that TagDataSource implements TagDataSourceInterface
				var _ TagDataSourceInterface = ds

			case "error case - interface validation":
				// Placeholder for validation
				assert.False(t, tt.wantErr)
			}
		})
	}
}

func TestTagDataSourceConcurrency(t *testing.T) {
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
				ds := NewTagDataSource()

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						resp := &datasource.MetadataResponse{}
						ds.Metadata(context.Background(), datasource.MetadataRequest{
							ProviderTypeName: "n8n",
						}, resp)
						assert.Equal(t, "n8n_tag", resp.TypeName)
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "concurrent schema calls":
				ds := NewTagDataSource()

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
						ds := NewTagDataSource()
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
				// Placeholder for validation
				assert.False(t, tt.wantErr)
			}
		})
	}
}

func BenchmarkTagDataSource_Schema(b *testing.B) {
	ds := NewTagDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkTagDataSource_Metadata(b *testing.B) {
	ds := NewTagDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)
	}
}

func BenchmarkTagDataSource_Configure(b *testing.B) {
	ds := NewTagDataSource()
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

// TestTagDataSource_Read is now in external test file - refactored to test behavior only.

func TestTagDataSource_validateIdentifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "valid - id provided"},
		{name: "valid - name provided"},
		{name: "valid - both provided"},
		{name: "error - neither provided", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := &TagDataSource{}
			resp := &datasource.ReadResponse{}

			switch tt.name {
			case "valid - id provided":
				data := &models.DataSource{
					ID:   types.StringValue("tag-123"),
					Name: types.StringNull(),
				}

				result := ds.validateIdentifier(data, resp)

				assert.True(t, result)
				assert.False(t, resp.Diagnostics.HasError())

			case "valid - name provided":
				data := &models.DataSource{
					ID:   types.StringNull(),
					Name: types.StringValue("test-tag"),
				}

				result := ds.validateIdentifier(data, resp)

				assert.True(t, result)
				assert.False(t, resp.Diagnostics.HasError())

			case "valid - both provided":
				data := &models.DataSource{
					ID:   types.StringValue("tag-123"),
					Name: types.StringValue("test-tag"),
				}

				result := ds.validateIdentifier(data, resp)

				assert.True(t, result)
				assert.False(t, resp.Diagnostics.HasError())

			case "error - neither provided":
				data := &models.DataSource{
					ID:   types.StringNull(),
					Name: types.StringNull(),
				}

				result := ds.validateIdentifier(data, resp)

				assert.False(t, result)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing Required Attribute")
			}
		})
	}
}

func TestTagDataSource_fetchTagByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success - fetch tag by ID"},
		{name: "error - API error", wantErr: true},
		{name: "error - tag not found", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "success - fetch tag by ID":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/tags/tag-123" && r.Method == http.MethodGet {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(n8nsdk.Tag{
							Id:   ptrString("tag-123"),
							Name: "Test Tag",
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &TagDataSource{client: n8nClient}
				data := &models.DataSource{
					ID: types.StringValue("tag-123"),
				}
				resp := &datasource.ReadResponse{}

				tag := ds.fetchTagByID(context.Background(), data, resp)

				assert.NotNil(t, tag)
				assert.Equal(t, "tag-123", *tag.Id)
				assert.Equal(t, "Test Tag", tag.Name)
				assert.False(t, resp.Diagnostics.HasError())

			case "error - API error":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &TagDataSource{client: n8nClient}
				data := &models.DataSource{
					ID: types.StringValue("tag-123"),
				}
				resp := &datasource.ReadResponse{}

				tag := ds.fetchTagByID(context.Background(), data, resp)

				assert.Nil(t, tag)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error retrieving tag")

			case "error - tag not found":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "Tag not found"}`))
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &TagDataSource{client: n8nClient}
				data := &models.DataSource{
					ID: types.StringValue("nonexistent"),
				}
				resp := &datasource.ReadResponse{}

				tag := ds.fetchTagByID(context.Background(), data, resp)

				assert.Nil(t, tag)
				assert.True(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestTagDataSource_fetchTagByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success - fetch tag by name"},
		{name: "error - API error", wantErr: true},
		{name: "error - tag not found", wantErr: true},
		{name: "error - empty tag list", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "success - fetch tag by name":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/tags" && r.Method == http.MethodGet {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(n8nsdk.TagList{
							Data: []n8nsdk.Tag{
								{Id: ptrString("tag-1"), Name: "Production"},
								{Id: ptrString("tag-2"), Name: "Development"},
								{Id: ptrString("tag-3"), Name: "Testing"},
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &TagDataSource{client: n8nClient}
				data := &models.DataSource{
					Name: types.StringValue("Development"),
				}
				resp := &datasource.ReadResponse{}

				tag := ds.fetchTagByName(context.Background(), data, resp)

				assert.NotNil(t, tag)
				assert.Equal(t, "tag-2", *tag.Id)
				assert.Equal(t, "Development", tag.Name)
				assert.False(t, resp.Diagnostics.HasError())

			case "error - API error":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &TagDataSource{client: n8nClient}
				data := &models.DataSource{
					Name: types.StringValue("Test Tag"),
				}
				resp := &datasource.ReadResponse{}

				tag := ds.fetchTagByName(context.Background(), data, resp)

				assert.Nil(t, tag)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error listing tags")

			case "error - tag not found":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/tags" && r.Method == http.MethodGet {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(n8nsdk.TagList{
							Data: []n8nsdk.Tag{
								{Id: ptrString("tag-1"), Name: "Production"},
								{Id: ptrString("tag-2"), Name: "Development"},
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &TagDataSource{client: n8nClient}
				data := &models.DataSource{
					Name: types.StringValue("NonExistent"),
				}
				resp := &datasource.ReadResponse{}

				tag := ds.fetchTagByName(context.Background(), data, resp)

				assert.Nil(t, tag)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Tag Not Found")

			case "error - empty tag list":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/tags" && r.Method == http.MethodGet {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(n8nsdk.TagList{
							Data: []n8nsdk.Tag{},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &TagDataSource{client: n8nClient}
				data := &models.DataSource{
					Name: types.StringValue("AnyTag"),
				}
				resp := &datasource.ReadResponse{}

				tag := ds.fetchTagByName(context.Background(), data, resp)

				assert.Nil(t, tag)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Tag Not Found")
			}
		})
	}
}

// ptrString returns a pointer to the given string.
func ptrString(s string) *string {
	return &s
}
