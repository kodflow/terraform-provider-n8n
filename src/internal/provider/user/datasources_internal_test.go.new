package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUsersDataSourceInterface is a mock implementation of UsersDataSourceInterface.
type MockUsersDataSourceInterface struct {
	mock.Mock
}

func (m *MockUsersDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockUsersDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockUsersDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockUsersDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

func TestNewUsersDataSource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		ds := NewUsersDataSource()

		assert.NotNil(t, ds)
		assert.IsType(t, &UsersDataSource{}, ds)
		assert.Nil(t, ds.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		ds1 := NewUsersDataSource()
		ds2 := NewUsersDataSource()

		assert.NotSame(t, ds1, ds2)
	})
}

func TestNewUsersDataSourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		ds := NewUsersDataSourceWrapper()

		assert.NotNil(t, ds)
		assert.IsType(t, &UsersDataSource{}, ds)
	})

	t.Run("wrapper returns datasource.DataSource interface", func(t *testing.T) {
		ds := NewUsersDataSourceWrapper()

		// ds is already of type datasource.DataSource, no assertion needed
		assert.NotNil(t, ds)
	})
}

func TestUsersDataSource_Metadata(t *testing.T) {
	t.Run("set metadata with provider type", func(t *testing.T) {
		ds := NewUsersDataSource()
		resp := &datasource.MetadataResponse{}

		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_users", resp.TypeName)
	})

	t.Run("set metadata with provider type name", func(t *testing.T) {
		ds := NewUsersDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "n8n_users", resp.TypeName)
	})

	t.Run("set metadata with different provider type", func(t *testing.T) {
		ds := NewUsersDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_users", resp.TypeName)
	})
}

func TestUsersDataSource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		ds := NewUsersDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "n8n users")

		// Verify users list attribute
		usersAttr, exists := resp.Schema.Attributes["users"]
		assert.True(t, exists)
		assert.True(t, usersAttr.IsComputed())
	})

	t.Run("users list has correct structure", func(t *testing.T) {
		ds := NewUsersDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		usersAttr, exists := resp.Schema.Attributes["users"]
		assert.True(t, exists)
		assert.True(t, usersAttr.IsComputed())
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		ds := NewUsersDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			if name == "users" {
				assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
			}
		}
	})
}

func TestUsersDataSource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		ds := NewUsersDataSource()
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
		ds := NewUsersDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: nil,
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		ds := NewUsersDataSource()
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
		ds := NewUsersDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("configure multiple times", func(t *testing.T) {
		ds := NewUsersDataSource()

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

func TestUsersDataSource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		ds := NewUsersDataSource()

		// Test that UsersDataSource implements datasource.DataSource
		var _ datasource.DataSource = ds

		// Test that UsersDataSource implements datasource.DataSourceWithConfigure
		var _ datasource.DataSourceWithConfigure = ds

		// Test that UsersDataSource implements UsersDataSourceInterface
		var _ UsersDataSourceInterface = ds
	})
}

func TestUsersDataSourceConcurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		ds := NewUsersDataSource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)
				assert.Equal(t, "n8n_users", resp.TypeName)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		ds := NewUsersDataSource()

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

func BenchmarkUsersDataSource_Schema(b *testing.B) {
	ds := NewUsersDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkUsersDataSource_Metadata(b *testing.B) {
	ds := NewUsersDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{}, resp)
	}
}

func BenchmarkUsersDataSource_Configure(b *testing.B) {
	ds := NewUsersDataSource()
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

// TestUsersDataSource_Read tests the read functionality.
func TestUsersDataSource_Read(t *testing.T) {
	t.Run("successful read with users", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []map[string]interface{}{
						{
							"id":        "user-1",
							"email":     "user1@example.com",
							"firstName": "User",
							"lastName":  "One",
							"isPending": false,
							"createdAt": "2024-01-01T00:00:00Z",
							"updatedAt": "2024-01-02T00:00:00Z",
							"role":      "global:member",
						},
						{
							"id":        "user-2",
							"email":     "user2@example.com",
							"firstName": "User",
							"lastName":  "Two",
							"isPending": true,
							"createdAt": "2024-01-01T00:00:00Z",
							"updatedAt": "2024-01-02T00:00:00Z",
							"role":      "global:admin",
						},
					},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestUsersDataSourceClient(t, handler)
		defer server.Close()

		ds := &UsersDataSource{client: n8nClient}

		resp := datasource.ReadResponse{
			State: tfsdk.State{
				Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"users": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String, "role": tftypes.String}}}}}, nil),
				Schema: createTestUsersDataSourceSchema(t),
			},
		}

		ds.Read(context.Background(), datasource.ReadRequest{}, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")
	})

	t.Run("successful read with empty users", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []map[string]interface{}{},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestUsersDataSourceClient(t, handler)
		defer server.Close()

		ds := &UsersDataSource{client: n8nClient}

		resp := datasource.ReadResponse{
			State: tfsdk.State{
				Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"users": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String, "role": tftypes.String}}}}}, nil),
				Schema: createTestUsersDataSourceSchema(t),
			},
		}

		ds.Read(context.Background(), datasource.ReadRequest{}, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")
	})

	t.Run("successful read with nil data", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestUsersDataSourceClient(t, handler)
		defer server.Close()

		ds := &UsersDataSource{client: n8nClient}

		resp := datasource.ReadResponse{
			State: tfsdk.State{
				Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"users": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String, "role": tftypes.String}}}}}, nil),
				Schema: createTestUsersDataSourceSchema(t),
			},
		}

		ds.Read(context.Background(), datasource.ReadRequest{}, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")
	})

	t.Run("read fails with error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestUsersDataSourceClient(t, handler)
		defer server.Close()

		ds := &UsersDataSource{client: n8nClient}

		resp := datasource.ReadResponse{
			State: tfsdk.State{
				Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"users": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String, "role": tftypes.String}}}}}, nil),
				Schema: createTestUsersDataSourceSchema(t),
			},
		}

		ds.Read(context.Background(), datasource.ReadRequest{}, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors when API fails")
	})
}

// createTestUsersDataSourceSchema creates a test schema for users datasource.
func createTestUsersDataSourceSchema(t *testing.T) datasource_schema.Schema {
	t.Helper()
	ds := &UsersDataSource{}
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), req, resp)
	return resp.Schema
}

// setupTestUsersDataSourceClient creates a test N8nClient with httptest server.
func setupTestUsersDataSourceClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
