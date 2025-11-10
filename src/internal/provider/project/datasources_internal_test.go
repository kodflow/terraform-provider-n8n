package project

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

func TestNewProjectsDataSource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		ds := NewProjectsDataSource()

		assert.NotNil(t, ds)
		assert.IsType(t, &ProjectsDataSource{}, ds)
		assert.Nil(t, ds.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		ds1 := NewProjectsDataSource()
		ds2 := NewProjectsDataSource()

		assert.NotSame(t, ds1, ds2)
	})
}

func TestNewProjectsDataSourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		ds := NewProjectsDataSourceWrapper()

		assert.NotNil(t, ds)
		assert.IsType(t, &ProjectsDataSource{}, ds)
	})

	t.Run("wrapper returns datasource.DataSource interface", func(t *testing.T) {
		ds := NewProjectsDataSourceWrapper()

		// Verify it implements datasource.DataSource.
		assert.Implements(t, (*datasource.DataSource)(nil), ds)
		assert.NotNil(t, ds)
	})
}

func TestProjectsDataSource_Metadata(t *testing.T) {
	t.Run("set metadata with provider type", func(t *testing.T) {
		ds := NewProjectsDataSource()
		resp := &datasource.MetadataResponse{}

		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_projects", resp.TypeName)
	})

	t.Run("set metadata with custom provider type", func(t *testing.T) {
		ds := NewProjectsDataSource()
		resp := &datasource.MetadataResponse{}

		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}, resp)

		assert.Equal(t, "custom_provider_projects", resp.TypeName)
	})
}

func TestProjectsDataSource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		ds := NewProjectsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "n8n projects")

		// Verify "projects" attribute exists.
		_, exists := resp.Schema.Attributes["projects"]
		assert.True(t, exists, "Attribute 'projects' should exist")
	})

	t.Run("schema has description", func(t *testing.T) {
		ds := NewProjectsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		assert.NotEmpty(t, resp.Schema.MarkdownDescription)
	})
}

func TestProjectsDataSource_schemaAttributes(t *testing.T) {
	t.Run("returns projects attribute", func(t *testing.T) {
		ds := NewProjectsDataSource()
		attrs := ds.schemaAttributes()

		assert.NotNil(t, attrs)
		assert.Contains(t, attrs, "projects")
	})

	t.Run("projects attribute is list nested", func(t *testing.T) {
		ds := NewProjectsDataSource()
		attrs := ds.schemaAttributes()

		projectsAttr := attrs["projects"]
		assert.NotNil(t, projectsAttr)
		assert.NotEmpty(t, projectsAttr.GetMarkdownDescription())
	})
}

func TestProjectsDataSource_projectAttributes(t *testing.T) {
	t.Run("returns all project attributes", func(t *testing.T) {
		ds := NewProjectsDataSource()
		attrs := ds.projectAttributes()

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
			_, exists := attrs[attr]
			assert.True(t, exists, "Attribute %s should exist", attr)
		}
	})

	t.Run("all attributes have descriptions", func(t *testing.T) {
		ds := NewProjectsDataSource()
		attrs := ds.projectAttributes()

		for name, attr := range attrs {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})
}

func TestProjectsDataSource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		ds := NewProjectsDataSource()
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
		ds := NewProjectsDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: nil,
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		ds := NewProjectsDataSource()
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
		ds := NewProjectsDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.True(t, resp.Diagnostics.HasError())
	})
}

func TestProjectsDataSource_Read(t *testing.T) {
	t.Run("successful read with projects", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/projects" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []map[string]interface{}{
						{
							"id":          "proj-123",
							"name":        "Test Project 1",
							"type":        "team",
							"createdAt":   "2024-01-01T00:00:00Z",
							"updatedAt":   "2024-01-02T00:00:00Z",
							"icon":        "home",
							"description": "First test project",
						},
						{
							"id":          "proj-456",
							"name":        "Test Project 2",
							"type":        "personal",
							"createdAt":   "2024-01-03T00:00:00Z",
							"updatedAt":   "2024-01-04T00:00:00Z",
							"icon":        "star",
							"description": "Second test project",
						},
					},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestProjectsClient(t, handler)
		defer server.Close()

		ds := &ProjectsDataSource{client: n8nClient}

		req := datasource.ReadRequest{}
		resp := &datasource.ReadResponse{
			State: createTestProjectsDataSourceState(t),
		}

		ds.Read(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")
	})

	t.Run("successful read with empty project list", func(t *testing.T) {
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

		n8nClient, server := setupTestProjectsClient(t, handler)
		defer server.Close()

		ds := &ProjectsDataSource{client: n8nClient}

		req := datasource.ReadRequest{}
		resp := &datasource.ReadResponse{
			State: createTestProjectsDataSourceState(t),
		}

		ds.Read(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors with empty list")
	})

	t.Run("successful read with nil data", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/projects" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestProjectsClient(t, handler)
		defer server.Close()

		ds := &ProjectsDataSource{client: n8nClient}

		req := datasource.ReadRequest{}
		resp := &datasource.ReadResponse{
			State: createTestProjectsDataSourceState(t),
		}

		ds.Read(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should handle nil data gracefully")
	})

	t.Run("read fails with API error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestProjectsClient(t, handler)
		defer server.Close()

		ds := &ProjectsDataSource{client: n8nClient}

		req := datasource.ReadRequest{}
		resp := &datasource.ReadResponse{
			State: createTestProjectsDataSourceState(t),
		}

		ds.Read(context.Background(), req, resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors on API failure")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error listing projects")
	})

	t.Run("read fails with unauthorized error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message": "Unauthorized"}`))
		})

		n8nClient, server := setupTestProjectsClient(t, handler)
		defer server.Close()

		ds := &ProjectsDataSource{client: n8nClient}

		req := datasource.ReadRequest{}
		resp := &datasource.ReadResponse{
			State: createTestProjectsDataSourceState(t),
		}

		ds.Read(context.Background(), req, resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors on unauthorized")
	})
}

func TestProjectsDataSource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		ds := NewProjectsDataSource()

		var _ datasource.DataSource = ds
		var _ datasource.DataSourceWithConfigure = ds
		var _ ProjectsDataSourceInterface = ds
	})
}

// createTestProjectsDataSourceState creates a test state for projects datasource.
func createTestProjectsDataSourceState(t *testing.T) tfsdk.State {
	t.Helper()
	ds := &ProjectsDataSource{}
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), req, resp)
	return tfsdk.State{
		Schema: resp.Schema,
	}
}

// setupTestProjectsClient creates a test N8nClient with httptest server.
func setupTestProjectsClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
