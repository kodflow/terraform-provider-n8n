package user

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

// MockUserDataSourceInterface is a mock implementation of UserDataSourceInterface
type MockUserDataSourceInterface struct {
	mock.Mock
}

func (m *MockUserDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockUserDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockUserDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockUserDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

func TestNewUserDataSource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		ds := NewUserDataSource()

		assert.NotNil(t, ds)
		assert.IsType(t, &UserDataSource{}, ds)
		assert.Nil(t, ds.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		ds1 := NewUserDataSource()
		ds2 := NewUserDataSource()

		assert.NotSame(t, ds1, ds2)
	})
}

func TestNewUserDataSourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		ds := NewUserDataSourceWrapper()

		assert.NotNil(t, ds)
		assert.IsType(t, &UserDataSource{}, ds)
	})

	t.Run("wrapper returns datasource.DataSource interface", func(t *testing.T) {
		ds := NewUserDataSourceWrapper()

		_, ok := ds.(datasource.DataSource)
		assert.True(t, ok)
	})
}

func TestUserDataSource_Metadata(t *testing.T) {
	t.Run("set metadata with provider type", func(t *testing.T) {
		ds := NewUserDataSource()
		resp := &datasource.MetadataResponse{}

		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_user", resp.TypeName)
	})

	t.Run("set metadata with provider type name", func(t *testing.T) {
		ds := NewUserDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "n8n_user", resp.TypeName)
	})

	t.Run("set metadata with different provider type", func(t *testing.T) {
		ds := NewUserDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_user", resp.TypeName)
	})
}

func TestUserDataSource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		ds := NewUserDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "n8n user")

		// Verify optional attributes that can be used as identifiers
		idAttr, exists := resp.Schema.Attributes["id"]
		assert.True(t, exists)
		assert.True(t, idAttr.IsOptional())
		assert.True(t, idAttr.IsComputed())

		emailAttr, exists := resp.Schema.Attributes["email"]
		assert.True(t, exists)
		assert.True(t, emailAttr.IsOptional())
		assert.True(t, emailAttr.IsComputed())

		// Verify computed attributes
		computedAttrs := []string{
			"first_name",
			"last_name",
			"is_pending",
			"created_at",
			"updated_at",
			"role",
		}

		for _, attr := range computedAttrs {
			schemaAttr, exists := resp.Schema.Attributes[attr]
			assert.True(t, exists, "Attribute %s should exist", attr)
			assert.True(t, schemaAttr.IsComputed(), "Attribute %s should be computed", attr)
		}
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		ds := NewUserDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})
}

func TestUserDataSource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		ds := NewUserDataSource()
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
		ds := NewUserDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: nil,
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		ds := NewUserDataSource()
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
		ds := NewUserDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("configure multiple times", func(t *testing.T) {
		ds := NewUserDataSource()

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

func TestUserDataSource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		ds := NewUserDataSource()

		// Test that UserDataSource implements datasource.DataSource
		var _ datasource.DataSource = ds

		// Test that UserDataSource implements datasource.DataSourceWithConfigure
		var _ datasource.DataSourceWithConfigure = ds

		// Test that UserDataSource implements UserDataSourceInterface
		var _ UserDataSourceInterface = ds
	})
}

func TestUserDataSourceConcurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		ds := NewUserDataSource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)
				assert.Equal(t, "n8n_user", resp.TypeName)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		ds := NewUserDataSource()

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

func BenchmarkUserDataSource_Schema(b *testing.B) {
	ds := NewUserDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkUserDataSource_Metadata(b *testing.B) {
	ds := NewUserDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{}, resp)
	}
}

func BenchmarkUserDataSource_Configure(b *testing.B) {
	ds := NewUserDataSource()
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

func TestUserDataSource_Read(t *testing.T) {
	t.Run("successful read by ID", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users/user-123" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":        "user-123",
					"email":     "test@example.com",
					"firstName": "Test",
					"lastName":  "User",
					"isPending": false,
					"createdAt": "2024-01-01T00:00:00Z",
					"updatedAt": "2024-01-02T00:00:00Z",
					"role":      "global:member",
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestDataSourceClient(t, handler)
		defer server.Close()

		ds := &UserDataSource{client: n8nClient}

		rawConfig := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-123"),
			"email":      tftypes.NewValue(tftypes.String, nil),
			"first_name": tftypes.NewValue(tftypes.String, nil),
			"last_name":  tftypes.NewValue(tftypes.String, nil),
			"is_pending": tftypes.NewValue(tftypes.Bool, nil),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
			"role":       tftypes.NewValue(tftypes.String, nil),
		}
		attrTypes := map[string]tftypes.Type{
			"id":         tftypes.String,
			"email":      tftypes.String,
			"first_name": tftypes.String,
			"last_name":  tftypes.String,
			"is_pending": tftypes.Bool,
			"created_at": tftypes.String,
			"updated_at": tftypes.String,
			"role":       tftypes.String,
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

	t.Run("successful read by email", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users/test@example.com" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":        "user-456",
					"email":     "test@example.com",
					"firstName": "Test",
					"lastName":  "User",
					"isPending": true,
					"createdAt": "2024-01-01T00:00:00Z",
					"updatedAt": "2024-01-02T00:00:00Z",
					"role":      "global:admin",
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestDataSourceClient(t, handler)
		defer server.Close()

		ds := &UserDataSource{client: n8nClient}

		rawConfig := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, nil),
			"last_name":  tftypes.NewValue(tftypes.String, nil),
			"is_pending": tftypes.NewValue(tftypes.Bool, nil),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
			"role":       tftypes.NewValue(tftypes.String, nil),
		}
		attrTypes := map[string]tftypes.Type{
			"id":         tftypes.String,
			"email":      tftypes.String,
			"first_name": tftypes.String,
			"last_name":  tftypes.String,
			"is_pending": tftypes.Bool,
			"created_at": tftypes.String,
			"updated_at": tftypes.String,
			"role":       tftypes.String,
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

	t.Run("user not found", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "User not found"}`))
		})

		n8nClient, server := setupTestDataSourceClient(t, handler)
		defer server.Close()

		ds := &UserDataSource{client: n8nClient}

		rawConfig := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-nonexistent"),
			"email":      tftypes.NewValue(tftypes.String, nil),
			"first_name": tftypes.NewValue(tftypes.String, nil),
			"last_name":  tftypes.NewValue(tftypes.String, nil),
			"is_pending": tftypes.NewValue(tftypes.Bool, nil),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
			"role":       tftypes.NewValue(tftypes.String, nil),
		}
		attrTypes := map[string]tftypes.Type{
			"id":         tftypes.String,
			"email":      tftypes.String,
			"first_name": tftypes.String,
			"last_name":  tftypes.String,
			"is_pending": tftypes.Bool,
			"created_at": tftypes.String,
			"updated_at": tftypes.String,
			"role":       tftypes.String,
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

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors when user not found")
	})
}

// createTestDataSourceSchema creates a test schema for user datasource
func createTestDataSourceSchema(t *testing.T) schema.Schema {
	ds := &UserDataSource{}
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
