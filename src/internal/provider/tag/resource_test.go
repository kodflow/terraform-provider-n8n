package tag

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// createTestSchema creates a test schema for tag resource.
func createTestSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := &TagResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), req, resp)
	return resp.Schema
}

// setupTestClient creates a test N8nClient with httptest server.
func setupTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)

	cfg := n8nsdk.NewConfiguration()
	// Parse server URL to extract host and port
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

// TestTagResource_Create tests tag creation.
func TestTagResource_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/tags":
				if r.Method == http.MethodPost {
					tag := map[string]interface{}{
						"id":   "tag-123",
						"name": "Test Tag",
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					json.NewEncoder(w).Encode(tag)
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &TagResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
			Schema: createTestSchema(t),
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Create should not have errors")
	})

	t.Run("creation fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &TagResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
			Schema: createTestSchema(t),
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Create should have errors")
	})
}

// TestTagResource_Read tests tag reading.
func TestTagResource_Read(t *testing.T) {
	t.Run("successful read", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/tags/tag-123" && r.Method == http.MethodGet {
				tag := map[string]interface{}{
					"id":   "tag-123",
					"name": "Test Tag",
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(tag)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &TagResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "tag-123"),
			"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")
	})

	t.Run("tag not found removes from state", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/tags/tag-nonexistent" && r.Method == http.MethodGet {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Tag not found"}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &TagResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "tag-nonexistent"),
			"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors when tag not found")
		// Note: In Terraform, when a resource is not found during Read, it's removed from state
		// This is tested by checking that resp.State would be empty, but we can't easily verify that in unit tests
	})
}

// TestTagResource_Update tests tag update.
func TestTagResource_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/tags/tag-123" && r.Method == http.MethodPut:
				tag := map[string]interface{}{
					"id":   "tag-123",
					"name": "Updated Tag",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(tag)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &TagResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "tag-123"),
			"name":       tftypes.NewValue(tftypes.String, "Updated Tag"),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "tag-123"),
			"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		req := resource.UpdateRequest{
			Plan:  plan,
			State: state,
		}
		resp := resource.UpdateResponse{
			State: state,
		}

		r.Update(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Update should not have errors")
	})
}

// TestTagResource_Delete tests tag deletion.
func TestTagResource_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/tags/tag-123" && r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &TagResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "tag-123"),
			"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		req := resource.DeleteRequest{
			State: state,
		}
		resp := resource.DeleteResponse{
			State: state,
		}

		r.Delete(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Delete should not have errors")
	})

	t.Run("delete fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &TagResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "tag-123"),
			"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		req := resource.DeleteRequest{
			State: state,
		}
		resp := resource.DeleteResponse{
			State: state,
		}

		r.Delete(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Delete should have errors")
	})
}

// TestTagResource_ImportState tests state import.
func TestTagResource_ImportState(t *testing.T) {
	r := &TagResource{}

	schema := createTestSchema(t)
	state := tfsdk.State{
		Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
		Schema: schema,
	}

	req := resource.ImportStateRequest{
		ID: "tag-123",
	}
	resp := &resource.ImportStateResponse{
		State: state,
	}

	r.ImportState(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "ImportState should not have errors")
}

// TestNewTagResource tests NewTagResource constructor.
func TestNewTagResource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		r := NewTagResource()

		assert.NotNil(t, r)
		assert.IsType(t, &TagResource{}, r)
		assert.Nil(t, r.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		r1 := NewTagResource()
		r2 := NewTagResource()

		assert.NotSame(t, r1, r2)
	})
}

// TestNewTagResourceWrapper tests NewTagResourceWrapper constructor.
func TestNewTagResourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		r := NewTagResourceWrapper()

		assert.NotNil(t, r)
		assert.IsType(t, &TagResource{}, r)
	})

	t.Run("wrapper returns resource.Resource interface", func(t *testing.T) {
		r := NewTagResourceWrapper()
		assert.NotNil(t, r)
	})
}

// TestTagResource_Metadata tests Metadata method.
func TestTagResource_Metadata(t *testing.T) {
	t.Run("set metadata with n8n provider", func(t *testing.T) {
		r := NewTagResource()
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_tag", resp.TypeName)
	})

	t.Run("set metadata with custom provider type", func(t *testing.T) {
		r := NewTagResource()
		resp := &resource.MetadataResponse{}
		req := resource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_tag", resp.TypeName)
	})

	t.Run("set metadata with empty provider type", func(t *testing.T) {
		r := NewTagResource()
		resp := &resource.MetadataResponse{}
		req := resource.MetadataRequest{
			ProviderTypeName: "",
		}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "_tag", resp.TypeName)
	})
}

// TestTagResource_Configure tests Configure method.
func TestTagResource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		r := NewTagResource()
		resp := &resource.ConfigureResponse{}

		mockClient := &client.N8nClient{}
		req := resource.ConfigureRequest{
			ProviderData: mockClient,
		}

		r.Configure(context.Background(), req, resp)

		assert.NotNil(t, r.client)
		assert.Equal(t, mockClient, r.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with nil provider data", func(t *testing.T) {
		r := NewTagResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: nil,
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		r := NewTagResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: "invalid-data",
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")
	})

	t.Run("configure with wrong type", func(t *testing.T) {
		r := NewTagResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("configure multiple times", func(t *testing.T) {
		r := NewTagResource()

		// First configuration
		resp1 := &resource.ConfigureResponse{}
		client1 := &client.N8nClient{}
		req1 := resource.ConfigureRequest{
			ProviderData: client1,
		}
		r.Configure(context.Background(), req1, resp1)
		assert.Equal(t, client1, r.client)

		// Second configuration
		resp2 := &resource.ConfigureResponse{}
		client2 := &client.N8nClient{}
		req2 := resource.ConfigureRequest{
			ProviderData: client2,
		}
		r.Configure(context.Background(), req2, resp2)
		assert.Equal(t, client2, r.client)
	})
}

// TestTagResource_Create_WithTimestamps tests create with timestamps.
func TestTagResource_Create_WithTimestamps(t *testing.T) {
	t.Run("successful creation with timestamps", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/tags":
				if r.Method == http.MethodPost {
					tag := map[string]interface{}{
						"id":        "tag-123",
						"name":      "Test Tag",
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-02T00:00:00Z",
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					json.NewEncoder(w).Encode(tag)
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &TagResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
			Schema: createTestSchema(t),
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Create should not have errors")
	})

	t.Run("creation with invalid plan", func(t *testing.T) {
		r := &TagResource{}

		ctx := context.Background()
		schema := createTestSchema(t)

		// Create invalid state
		state := tfsdk.State{
			Schema: schema,
			Raw:    tftypes.NewValue(tftypes.String, "invalid"),
		}

		plan := tfsdk.Plan{
			Schema: schema,
			Raw:    tftypes.NewValue(tftypes.String, "invalid"),
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(ctx, req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Create should have errors with invalid plan")
	})
}

// TestTagResource_Read_WithTimestamps tests read with timestamps.
func TestTagResource_Read_WithTimestamps(t *testing.T) {
	t.Run("successful read with timestamps", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/tags/tag-123" && r.Method == http.MethodGet {
				tag := map[string]interface{}{
					"id":        "tag-123",
					"name":      "Test Tag",
					"createdAt": "2024-01-01T00:00:00Z",
					"updatedAt": "2024-01-02T00:00:00Z",
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(tag)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &TagResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "tag-123"),
			"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")
	})

	t.Run("read with invalid state", func(t *testing.T) {
		r := &TagResource{}

		ctx := context.Background()
		schema := createTestSchema(t)

		// Create invalid state
		state := tfsdk.State{
			Schema: schema,
			Raw:    tftypes.NewValue(tftypes.String, "invalid"),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(ctx, req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors with invalid state")
	})
}

// TestTagResource_Update_WithTimestamps tests update with timestamps.
func TestTagResource_Update_WithTimestamps(t *testing.T) {
	t.Run("successful update with timestamps", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/tags/tag-123" && r.Method == http.MethodPut:
				tag := map[string]interface{}{
					"id":        "tag-123",
					"name":      "Updated Tag",
					"createdAt": "2024-01-01T00:00:00Z",
					"updatedAt": "2024-01-03T00:00:00Z",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(tag)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &TagResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "tag-123"),
			"name":       tftypes.NewValue(tftypes.String, "Updated Tag"),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "tag-123"),
			"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		req := resource.UpdateRequest{
			Plan:  plan,
			State: state,
		}
		resp := resource.UpdateResponse{
			State: state,
		}

		r.Update(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Update should not have errors")
	})

	t.Run("update with invalid plan", func(t *testing.T) {
		r := &TagResource{}

		ctx := context.Background()
		schema := createTestSchema(t)

		// Create invalid state
		state := tfsdk.State{
			Schema: schema,
			Raw:    tftypes.NewValue(tftypes.String, "invalid"),
		}

		plan := tfsdk.Plan{
			Schema: schema,
			Raw:    tftypes.NewValue(tftypes.String, "invalid"),
		}

		req := resource.UpdateRequest{
			Plan:  plan,
			State: state,
		}
		resp := resource.UpdateResponse{
			State: state,
		}

		r.Update(ctx, req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Update should have errors with invalid plan")
	})
}

// TestTagResource_Delete_WithState tests delete with state.
func TestTagResource_Delete_WithState(t *testing.T) {
	t.Run("delete with invalid state", func(t *testing.T) {
		r := &TagResource{}

		ctx := context.Background()
		schema := createTestSchema(t)

		// Create invalid state
		state := tfsdk.State{
			Schema: schema,
			Raw:    tftypes.NewValue(tftypes.String, "invalid"),
		}

		req := resource.DeleteRequest{
			State: state,
		}
		resp := resource.DeleteResponse{
			State: state,
		}

		r.Delete(ctx, req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Delete should have errors with invalid state")
	})
}

// TestTagResource_Interfaces tests that resource implements required interfaces.
func TestTagResource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		r := NewTagResource()

		// Test that TagResource implements resource.Resource
		var _ resource.Resource = r

		// Test that TagResource implements resource.ResourceWithConfigure
		var _ resource.ResourceWithConfigure = r

		// Test that TagResource implements resource.ResourceWithImportState
		var _ resource.ResourceWithImportState = r

		// Test that TagResource implements TagResourceInterface
		var _ TagResourceInterface = r
	})
}
