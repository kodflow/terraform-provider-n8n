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

// createTestSchema creates a test schema for tag resource
func createTestSchema(t *testing.T) schema.Schema {
	r := &TagResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), req, resp)
	return resp.Schema
}

// setupTestClient creates a test N8nClient with httptest server
func setupTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

// TestTagResource_Create tests tag creation
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

// TestTagResource_Read tests tag reading
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

// TestTagResource_Update tests tag update
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

// TestTagResource_Delete tests tag deletion
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

// TestTagResource_ImportState tests state import
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
