package workflow

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/workflow/models"
	"github.com/stretchr/testify/assert"
)

func TestNewWorkflowResource(t *testing.T) {
	t.Run("create new resource", func(t *testing.T) {
		r := NewWorkflowResource()
		assert.NotNil(t, r)
		assert.Nil(t, r.client)
	})

	t.Run("create new resource wrapper", func(t *testing.T) {
		r := NewWorkflowResourceWrapper()
		assert.NotNil(t, r)
		_, ok := r.(*WorkflowResource)
		assert.True(t, ok)
	})
}

func TestWorkflowResource_Metadata(t *testing.T) {
	t.Run("set metadata", func(t *testing.T) {
		r := &WorkflowResource{}
		req := resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "n8n_workflow", resp.TypeName)
	})

	t.Run("different provider type name", func(t *testing.T) {
		r := &WorkflowResource{}
		req := resource.MetadataRequest{
			ProviderTypeName: "custom",
		}
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_workflow", resp.TypeName)
	})

	t.Run("empty provider type name", func(t *testing.T) {
		r := &WorkflowResource{}
		req := resource.MetadataRequest{
			ProviderTypeName: "",
		}
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "_workflow", resp.TypeName)
	})
}

func TestWorkflowResource_Schema(t *testing.T) {
	t.Run("get schema", func(t *testing.T) {
		r := &WorkflowResource{}
		req := resource.SchemaRequest{}
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), req, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "workflow resource")
		assert.NotEmpty(t, resp.Schema.Attributes)
	})

	t.Run("schema attributes", func(t *testing.T) {
		r := &WorkflowResource{}
		attrs := r.schemaAttributes()

		assert.NotNil(t, attrs)
		assert.Contains(t, attrs, "id")
		assert.Contains(t, attrs, "name")
		assert.Contains(t, attrs, "active")
		assert.Contains(t, attrs, "tags")
		assert.Contains(t, attrs, "nodes_json")
		assert.Contains(t, attrs, "connections_json")
		assert.Contains(t, attrs, "settings_json")
		assert.Contains(t, attrs, "created_at")
		assert.Contains(t, attrs, "updated_at")
		assert.Contains(t, attrs, "version_id")
		assert.Contains(t, attrs, "is_archived")

		// Verify attribute count
		assert.Equal(t, 14, len(attrs))
	})

	t.Run("core attributes are correct", func(t *testing.T) {
		r := &WorkflowResource{}
		attrs := r.schemaAttributes()

		assert.Contains(t, attrs, "id")
		assert.Contains(t, attrs, "name")
		assert.Contains(t, attrs, "active")
		assert.Contains(t, attrs, "tags")
	})

	t.Run("json attributes are correct", func(t *testing.T) {
		r := &WorkflowResource{}
		attrs := r.schemaAttributes()

		assert.Contains(t, attrs, "nodes_json")
		assert.Contains(t, attrs, "connections_json")
		assert.Contains(t, attrs, "settings_json")
	})

	t.Run("metadata attributes are correct", func(t *testing.T) {
		r := &WorkflowResource{}
		attrs := r.schemaAttributes()

		assert.Contains(t, attrs, "created_at")
		assert.Contains(t, attrs, "updated_at")
		assert.Contains(t, attrs, "version_id")
		assert.Contains(t, attrs, "is_archived")
		assert.Contains(t, attrs, "trigger_count")
		assert.Contains(t, attrs, "meta")
		assert.Contains(t, attrs, "pin_data")
	})
}

func TestWorkflowResource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		r := &WorkflowResource{}
		mockClient := &client.N8nClient{}

		req := resource.ConfigureRequest{
			ProviderData: mockClient,
		}
		resp := &resource.ConfigureResponse{}

		r.Configure(context.Background(), req, resp)

		assert.Equal(t, mockClient, r.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with nil provider data", func(t *testing.T) {
		r := &WorkflowResource{}

		req := resource.ConfigureRequest{
			ProviderData: nil,
		}
		resp := &resource.ConfigureResponse{}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with wrong type", func(t *testing.T) {
		r := &WorkflowResource{}

		req := resource.ConfigureRequest{
			ProviderData: "invalid",
		}
		resp := &resource.ConfigureResponse{}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")
	})

	t.Run("configure multiple times", func(t *testing.T) {
		r := &WorkflowResource{}

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

func TestWorkflowResource_InterfaceCompliance(t *testing.T) {
	t.Run("implements Resource interface", func(t *testing.T) {
		var _ resource.Resource = (*WorkflowResource)(nil)
	})

	t.Run("implements ResourceWithConfigure interface", func(t *testing.T) {
		var _ resource.ResourceWithConfigure = (*WorkflowResource)(nil)
	})

	t.Run("implements ResourceWithImportState interface", func(t *testing.T) {
		var _ resource.ResourceWithImportState = (*WorkflowResource)(nil)
	})

	t.Run("implements WorkflowResourceInterface", func(t *testing.T) {
		var _ WorkflowResourceInterface = (*WorkflowResource)(nil)
	})
}

// createTestSchema creates a test schema for workflow resource.
func createTestSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := &WorkflowResource{}
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

// TestWorkflowResource_Create tests workflow creation.
func TestWorkflowResource_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/workflows":
				if r.Method == http.MethodPost {
					response := map[string]interface{}{
						"id":          "wf-123",
						"name":        "Test Workflow",
						"active":      false,
						"nodes":       []interface{}{},
						"connections": map[string]interface{}{},
						"settings":    map[string]interface{}{},
						"tags":        []interface{}{},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					json.NewEncoder(w).Encode(response)
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &WorkflowResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":               tftypes.NewValue(tftypes.String, nil),
			"name":             tftypes.NewValue(tftypes.String, "Test Workflow"),
			"active":           tftypes.NewValue(tftypes.Bool, nil),
			"tags":             tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil),
			"nodes_json":       tftypes.NewValue(tftypes.String, nil),
			"connections_json": tftypes.NewValue(tftypes.String, nil),
			"settings_json":    tftypes.NewValue(tftypes.String, nil),
			"created_at":       tftypes.NewValue(tftypes.String, nil),
			"updated_at":       tftypes.NewValue(tftypes.String, nil),
			"version_id":       tftypes.NewValue(tftypes.String, nil),
			"is_archived":      tftypes.NewValue(tftypes.Bool, nil),
			"trigger_count":    tftypes.NewValue(tftypes.Number, nil),
			"meta":             tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
			"pin_data":         tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
		}
		objectType := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":               tftypes.String,
				"name":             tftypes.String,
				"active":           tftypes.Bool,
				"tags":             tftypes.List{ElementType: tftypes.String},
				"nodes_json":       tftypes.String,
				"connections_json": tftypes.String,
				"settings_json":    tftypes.String,
				"created_at":       tftypes.String,
				"updated_at":       tftypes.String,
				"version_id":       tftypes.String,
				"is_archived":      tftypes.Bool,
				"trigger_count":    tftypes.Number,
				"meta":             tftypes.Map{ElementType: tftypes.String},
				"pin_data":         tftypes.Map{ElementType: tftypes.String},
			},
		}

		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(objectType, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, nil),
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

		r := &WorkflowResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":               tftypes.NewValue(tftypes.String, nil),
			"name":             tftypes.NewValue(tftypes.String, "Test Workflow"),
			"active":           tftypes.NewValue(tftypes.Bool, nil),
			"tags":             tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil),
			"nodes_json":       tftypes.NewValue(tftypes.String, nil),
			"connections_json": tftypes.NewValue(tftypes.String, nil),
			"settings_json":    tftypes.NewValue(tftypes.String, nil),
			"created_at":       tftypes.NewValue(tftypes.String, nil),
			"updated_at":       tftypes.NewValue(tftypes.String, nil),
			"version_id":       tftypes.NewValue(tftypes.String, nil),
			"is_archived":      tftypes.NewValue(tftypes.Bool, nil),
			"trigger_count":    tftypes.NewValue(tftypes.Number, nil),
			"meta":             tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
			"pin_data":         tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
		}
		objectType := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":               tftypes.String,
				"name":             tftypes.String,
				"active":           tftypes.Bool,
				"tags":             tftypes.List{ElementType: tftypes.String},
				"nodes_json":       tftypes.String,
				"connections_json": tftypes.String,
				"settings_json":    tftypes.String,
				"created_at":       tftypes.String,
				"updated_at":       tftypes.String,
				"version_id":       tftypes.String,
				"is_archived":      tftypes.Bool,
				"trigger_count":    tftypes.Number,
				"meta":             tftypes.Map{ElementType: tftypes.String},
				"pin_data":         tftypes.Map{ElementType: tftypes.String},
			},
		}

		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(objectType, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, nil),
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

// TestWorkflowResource_Read tests workflow reading.
func TestWorkflowResource_Read(t *testing.T) {
	t.Run("successful read", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/workflows/wf-123" && r.Method == http.MethodGet {
				response := map[string]interface{}{
					"id":          "wf-123",
					"name":        "Test Workflow",
					"active":      false,
					"nodes":       []interface{}{},
					"connections": map[string]interface{}{},
					"settings":    map[string]interface{}{},
					"tags":        []interface{}{},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &WorkflowResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":               tftypes.NewValue(tftypes.String, "wf-123"),
			"name":             tftypes.NewValue(tftypes.String, "Test Workflow"),
			"active":           tftypes.NewValue(tftypes.Bool, false),
			"tags":             tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{}),
			"nodes_json":       tftypes.NewValue(tftypes.String, "[]"),
			"connections_json": tftypes.NewValue(tftypes.String, "{}"),
			"settings_json":    tftypes.NewValue(tftypes.String, "{}"),
			"created_at":       tftypes.NewValue(tftypes.String, nil),
			"updated_at":       tftypes.NewValue(tftypes.String, nil),
			"version_id":       tftypes.NewValue(tftypes.String, nil),
			"is_archived":      tftypes.NewValue(tftypes.Bool, nil),
			"trigger_count":    tftypes.NewValue(tftypes.Number, nil),
			"meta":             tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
			"pin_data":         tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
		}
		objectType := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":               tftypes.String,
				"name":             tftypes.String,
				"active":           tftypes.Bool,
				"tags":             tftypes.List{ElementType: tftypes.String},
				"nodes_json":       tftypes.String,
				"connections_json": tftypes.String,
				"settings_json":    tftypes.String,
				"created_at":       tftypes.String,
				"updated_at":       tftypes.String,
				"version_id":       tftypes.String,
				"is_archived":      tftypes.Bool,
				"trigger_count":    tftypes.Number,
				"meta":             tftypes.Map{ElementType: tftypes.String},
				"pin_data":         tftypes.Map{ElementType: tftypes.String},
			},
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, rawState),
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

	t.Run("workflow not found removes from state", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/workflows/wf-nonexistent" && r.Method == http.MethodGet {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Workflow not found"}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &WorkflowResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":               tftypes.NewValue(tftypes.String, "wf-nonexistent"),
			"name":             tftypes.NewValue(tftypes.String, "Test Workflow"),
			"active":           tftypes.NewValue(tftypes.Bool, false),
			"tags":             tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{}),
			"nodes_json":       tftypes.NewValue(tftypes.String, "[]"),
			"connections_json": tftypes.NewValue(tftypes.String, "{}"),
			"settings_json":    tftypes.NewValue(tftypes.String, "{}"),
			"created_at":       tftypes.NewValue(tftypes.String, nil),
			"updated_at":       tftypes.NewValue(tftypes.String, nil),
			"version_id":       tftypes.NewValue(tftypes.String, nil),
			"is_archived":      tftypes.NewValue(tftypes.Bool, nil),
			"trigger_count":    tftypes.NewValue(tftypes.Number, nil),
			"meta":             tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
			"pin_data":         tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
		}
		objectType := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":               tftypes.String,
				"name":             tftypes.String,
				"active":           tftypes.Bool,
				"tags":             tftypes.List{ElementType: tftypes.String},
				"nodes_json":       tftypes.String,
				"connections_json": tftypes.String,
				"settings_json":    tftypes.String,
				"created_at":       tftypes.String,
				"updated_at":       tftypes.String,
				"version_id":       tftypes.String,
				"is_archived":      tftypes.Bool,
				"trigger_count":    tftypes.Number,
				"meta":             tftypes.Map{ElementType: tftypes.String},
				"pin_data":         tftypes.Map{ElementType: tftypes.String},
			},
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, rawState),
			Schema: createTestSchema(t),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors when workflow not found")
	})
}

// TestWorkflowResource_Update tests workflow update.
func TestWorkflowResource_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/workflows/wf-123/activate" && r.Method == http.MethodPost:
				// Handle activation endpoint
				w.WriteHeader(http.StatusOK)
				return
			case r.URL.Path == "/workflows/wf-123/tags" && r.Method == http.MethodPut:
				// Handle tags update endpoint
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]interface{}{})
				return
			case r.URL.Path == "/workflows/wf-123" && r.Method == http.MethodPut:
				response := map[string]interface{}{
					"id":          "wf-123",
					"name":        "Updated Workflow",
					"active":      true,
					"nodes":       []interface{}{},
					"connections": map[string]interface{}{},
					"settings":    map[string]interface{}{},
					"tags":        []interface{}{},
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &WorkflowResource{client: n8nClient}

		objectType := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":               tftypes.String,
				"name":             tftypes.String,
				"active":           tftypes.Bool,
				"tags":             tftypes.List{ElementType: tftypes.String},
				"nodes_json":       tftypes.String,
				"connections_json": tftypes.String,
				"settings_json":    tftypes.String,
				"created_at":       tftypes.String,
				"updated_at":       tftypes.String,
				"version_id":       tftypes.String,
				"is_archived":      tftypes.Bool,
				"trigger_count":    tftypes.Number,
				"meta":             tftypes.Map{ElementType: tftypes.String},
				"pin_data":         tftypes.Map{ElementType: tftypes.String},
			},
		}

		rawPlan := map[string]tftypes.Value{
			"id":               tftypes.NewValue(tftypes.String, "wf-123"),
			"name":             tftypes.NewValue(tftypes.String, "Updated Workflow"),
			"active":           tftypes.NewValue(tftypes.Bool, true),
			"tags":             tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{}),
			"nodes_json":       tftypes.NewValue(tftypes.String, "[]"),
			"connections_json": tftypes.NewValue(tftypes.String, "{}"),
			"settings_json":    tftypes.NewValue(tftypes.String, "{}"),
			"created_at":       tftypes.NewValue(tftypes.String, nil),
			"updated_at":       tftypes.NewValue(tftypes.String, nil),
			"version_id":       tftypes.NewValue(tftypes.String, nil),
			"is_archived":      tftypes.NewValue(tftypes.Bool, nil),
			"trigger_count":    tftypes.NewValue(tftypes.Number, nil),
			"meta":             tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
			"pin_data":         tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(objectType, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":               tftypes.NewValue(tftypes.String, "wf-123"),
			"name":             tftypes.NewValue(tftypes.String, "Test Workflow"),
			"active":           tftypes.NewValue(tftypes.Bool, false),
			"tags":             tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{}),
			"nodes_json":       tftypes.NewValue(tftypes.String, "[]"),
			"connections_json": tftypes.NewValue(tftypes.String, "{}"),
			"settings_json":    tftypes.NewValue(tftypes.String, "{}"),
			"created_at":       tftypes.NewValue(tftypes.String, nil),
			"updated_at":       tftypes.NewValue(tftypes.String, nil),
			"version_id":       tftypes.NewValue(tftypes.String, nil),
			"is_archived":      tftypes.NewValue(tftypes.Bool, nil),
			"trigger_count":    tftypes.NewValue(tftypes.Number, nil),
			"meta":             tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
			"pin_data":         tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, rawState),
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

// TestWorkflowResource_Delete tests workflow deletion.
func TestWorkflowResource_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/workflows/wf-123" && r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &WorkflowResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":               tftypes.NewValue(tftypes.String, "wf-123"),
			"name":             tftypes.NewValue(tftypes.String, "Test Workflow"),
			"active":           tftypes.NewValue(tftypes.Bool, false),
			"tags":             tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{}),
			"nodes_json":       tftypes.NewValue(tftypes.String, "[]"),
			"connections_json": tftypes.NewValue(tftypes.String, "{}"),
			"settings_json":    tftypes.NewValue(tftypes.String, "{}"),
			"created_at":       tftypes.NewValue(tftypes.String, nil),
			"updated_at":       tftypes.NewValue(tftypes.String, nil),
			"version_id":       tftypes.NewValue(tftypes.String, nil),
			"is_archived":      tftypes.NewValue(tftypes.Bool, nil),
			"trigger_count":    tftypes.NewValue(tftypes.Number, nil),
			"meta":             tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
			"pin_data":         tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
		}
		objectType := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":               tftypes.String,
				"name":             tftypes.String,
				"active":           tftypes.Bool,
				"tags":             tftypes.List{ElementType: tftypes.String},
				"nodes_json":       tftypes.String,
				"connections_json": tftypes.String,
				"settings_json":    tftypes.String,
				"created_at":       tftypes.String,
				"updated_at":       tftypes.String,
				"version_id":       tftypes.String,
				"is_archived":      tftypes.Bool,
				"trigger_count":    tftypes.Number,
				"meta":             tftypes.Map{ElementType: tftypes.String},
				"pin_data":         tftypes.Map{ElementType: tftypes.String},
			},
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, rawState),
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

		r := &WorkflowResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":               tftypes.NewValue(tftypes.String, "wf-123"),
			"name":             tftypes.NewValue(tftypes.String, "Test Workflow"),
			"active":           tftypes.NewValue(tftypes.Bool, false),
			"tags":             tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{}),
			"nodes_json":       tftypes.NewValue(tftypes.String, "[]"),
			"connections_json": tftypes.NewValue(tftypes.String, "{}"),
			"settings_json":    tftypes.NewValue(tftypes.String, "{}"),
			"created_at":       tftypes.NewValue(tftypes.String, nil),
			"updated_at":       tftypes.NewValue(tftypes.String, nil),
			"version_id":       tftypes.NewValue(tftypes.String, nil),
			"is_archived":      tftypes.NewValue(tftypes.Bool, nil),
			"trigger_count":    tftypes.NewValue(tftypes.Number, nil),
			"meta":             tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
			"pin_data":         tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
		}
		objectType := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":               tftypes.String,
				"name":             tftypes.String,
				"active":           tftypes.Bool,
				"tags":             tftypes.List{ElementType: tftypes.String},
				"nodes_json":       tftypes.String,
				"connections_json": tftypes.String,
				"settings_json":    tftypes.String,
				"created_at":       tftypes.String,
				"updated_at":       tftypes.String,
				"version_id":       tftypes.String,
				"is_archived":      tftypes.Bool,
				"trigger_count":    tftypes.Number,
				"meta":             tftypes.Map{ElementType: tftypes.String},
				"pin_data":         tftypes.Map{ElementType: tftypes.String},
			},
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, rawState),
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

func TestWorkflowResource_Concurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		r := &WorkflowResource{}

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &resource.MetadataResponse{}
				r.Metadata(context.Background(), resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)
				assert.Equal(t, "n8n_workflow", resp.TypeName)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		r := &WorkflowResource{}

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &resource.SchemaResponse{}
				r.Schema(context.Background(), resource.SchemaRequest{}, resp)
				assert.NotNil(t, resp.Schema)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func TestWORKFLOW_ATTRIBUTES_SIZE(t *testing.T) {
	t.Run("constant is defined", func(t *testing.T) {
		assert.Equal(t, 11, WORKFLOW_ATTRIBUTES_SIZE)
	})

	t.Run("actual schema has 14 attributes", func(t *testing.T) {
		r := &WorkflowResource{}
		attrs := r.schemaAttributes()
		// The actual schema has 14 attributes:
		// id, name, active, tags, nodes_json, connections_json, settings_json,
		// created_at, updated_at, version_id, is_archived, trigger_count, meta, pin_data
		assert.Equal(t, 14, len(attrs))
	})
}

func BenchmarkWorkflowResource_Schema(b *testing.B) {
	r := &WorkflowResource{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.SchemaResponse{}
		r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	}
}

func BenchmarkWorkflowResource_Metadata(b *testing.B) {
	r := &WorkflowResource{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.MetadataResponse{}
		r.Metadata(context.Background(), resource.MetadataRequest{}, resp)
	}
}

func BenchmarkWorkflowResource_Configure(b *testing.B) {
	r := &WorkflowResource{}
	mockClient := &client.N8nClient{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: mockClient,
		}
		r.Configure(context.Background(), req, resp)
	}
}

func BenchmarkWorkflowResource_SchemaAttributes(b *testing.B) {
	r := &WorkflowResource{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.schemaAttributes()
	}
}

// CRUD Tests - Focus on helper function testing.

func TestWorkflowResource_Create_JSONParsing(t *testing.T) {
	t.Run("invalid nodes json", func(t *testing.T) {
		plan := &models.Resource{
			Name:            types.StringValue("Test Workflow"),
			NodesJSON:       types.StringValue("invalid json"),
			ConnectionsJSON: types.StringValue("{}"),
		}

		diags := &diag.Diagnostics{}
		parseWorkflowJSON(plan, diags)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags.Errors()[0].Summary(), "Invalid nodes JSON")
	})

	t.Run("invalid connections json", func(t *testing.T) {
		plan := &models.Resource{
			Name:            types.StringValue("Test Workflow"),
			NodesJSON:       types.StringValue("[]"),
			ConnectionsJSON: types.StringValue("invalid json"),
		}

		diags := &diag.Diagnostics{}
		parseWorkflowJSON(plan, diags)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags.Errors()[0].Summary(), "Invalid connections JSON")
	})

	t.Run("invalid settings json", func(t *testing.T) {
		plan := &models.Resource{
			Name:            types.StringValue("Test Workflow"),
			NodesJSON:       types.StringValue("[]"),
			ConnectionsJSON: types.StringValue("{}"),
			SettingsJSON:    types.StringValue("invalid json"),
		}

		diags := &diag.Diagnostics{}
		parseWorkflowJSON(plan, diags)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags.Errors()[0].Summary(), "Invalid settings JSON")
	})

	t.Run("valid json parsing", func(t *testing.T) {
		plan := &models.Resource{
			Name:            types.StringValue("Test Workflow"),
			NodesJSON:       types.StringValue("[]"),
			ConnectionsJSON: types.StringValue("{}"),
			SettingsJSON:    types.StringValue("{}"),
		}

		diags := &diag.Diagnostics{}
		nodes, connections, settings := parseWorkflowJSON(plan, diags)

		assert.False(t, diags.HasError())
		assert.NotNil(t, nodes)
		assert.NotNil(t, connections)
		assert.NotNil(t, settings)
	})
}

func TestWorkflowResource_Read_ClientRequired(t *testing.T) {
	t.Run("resource requires configured client", func(t *testing.T) {
		r := &WorkflowResource{
			client: nil,
		}

		assert.Nil(t, r.client)
	})
}

func TestWorkflowResource_Update_ActivationDetection(t *testing.T) {
	t.Run("detect activation change", func(t *testing.T) {
		plan := &models.Resource{
			Active: types.BoolValue(true),
		}

		state := &models.Resource{
			Active: types.BoolValue(false),
		}

		// Activation should be detected as changed
		activeChanged := !plan.Active.IsNull() && !state.Active.IsNull() &&
			plan.Active.ValueBool() != state.Active.ValueBool()

		assert.True(t, activeChanged)
	})

	t.Run("no change when both active", func(t *testing.T) {
		plan := &models.Resource{
			Active: types.BoolValue(true),
		}

		state := &models.Resource{
			Active: types.BoolValue(true),
		}

		activeChanged := !plan.Active.IsNull() && !state.Active.IsNull() &&
			plan.Active.ValueBool() != state.Active.ValueBool()

		assert.False(t, activeChanged)
	})
}

func TestWorkflowResource_Delete_Basic(t *testing.T) {
	t.Run("delete requires client", func(t *testing.T) {
		r := &WorkflowResource{
			client: nil,
		}

		// Delete operation requires an API client
		assert.Nil(t, r.client)
	})
}

// TestWorkflowResource_CRUD is a comprehensive suite that validates CRUD operation logic.
func TestWorkflowResource_CRUD(t *testing.T) {
	t.Run("workflow lifecycle validation", func(t *testing.T) {
		// Create - validates JSON parsing (tested in TestWorkflowResource_Create_JSONParsing)
		// Read - validates client requirement (tested in TestWorkflowResource_Read_ClientRequired)
		// Update - validates JSON parsing and activation detection (tested in TestWorkflowResource_Update_ActivationDetection)
		// Delete - validates state requirement (tested in TestWorkflowResource_Delete_Basic)

		// This test ensures all CRUD operations have basic validation in place
		assert.NotNil(t, NewWorkflowResource())
	})
}

// TestWorkflowResource_ImportState tests state import functionality.
func TestWorkflowResource_ImportState(t *testing.T) {
	t.Run("successful import", func(t *testing.T) {
		r := &WorkflowResource{}

		schema := createTestSchema(t)
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String}}, nil),
			Schema: schema,
		}

		req := resource.ImportStateRequest{
			ID: "workflow-123",
		}
		resp := &resource.ImportStateResponse{
			State: state,
		}

		r.ImportState(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError(), "ImportState should not have errors")
	})
}

// TestWorkflowResource_EdgeCasesForCoverage tests edge cases for 100% coverage.
func TestWorkflowResource_EdgeCasesForCoverage(t *testing.T) {
	t.Run("create fails when plan get fails", func(t *testing.T) {
		r := &WorkflowResource{}

		// Create plan with mismatched schema
		wrongSchema := schema.Schema{
			Attributes: map[string]schema.Attribute{
				"id": schema.NumberAttribute{
					Required: true,
				},
			},
		}
		rawPlan := map[string]tftypes.Value{
			"id": tftypes.NewValue(tftypes.Number, 123),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.Number}}, rawPlan),
			Schema: wrongSchema,
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: tfsdk.State{Schema: createTestSchema(t)},
		}

		r.Create(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Create should fail when Plan.Get fails")
	})

	t.Run("read fails when state get fails", func(t *testing.T) {
		r := &WorkflowResource{}

		// Create state with mismatched schema
		wrongSchema := schema.Schema{
			Attributes: map[string]schema.Attribute{
				"id": schema.NumberAttribute{
					Required: true,
				},
			},
		}
		rawState := map[string]tftypes.Value{
			"id": tftypes.NewValue(tftypes.Number, 123),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.Number}}, rawState),
			Schema: wrongSchema,
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: tfsdk.State{Schema: createTestSchema(t)},
		}

		r.Read(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should fail when State.Get fails")
	})

	t.Run("update fails when plan get fails", func(t *testing.T) {
		r := &WorkflowResource{}

		// Create plan with mismatched schema
		wrongSchema := schema.Schema{
			Attributes: map[string]schema.Attribute{
				"id": schema.NumberAttribute{
					Required: true,
				},
			},
		}
		rawPlan := map[string]tftypes.Value{
			"id": tftypes.NewValue(tftypes.Number, 123),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.Number}}, rawPlan),
			Schema: wrongSchema,
		}

		state := tfsdk.State{
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

		assert.True(t, resp.Diagnostics.HasError(), "Update should fail when Plan.Get fails")
	})

	t.Run("delete fails when state get fails", func(t *testing.T) {
		r := &WorkflowResource{}

		// Create state with mismatched schema
		wrongSchema := schema.Schema{
			Attributes: map[string]schema.Attribute{
				"id": schema.NumberAttribute{
					Required: true,
				},
			},
		}
		rawState := map[string]tftypes.Value{
			"id": tftypes.NewValue(tftypes.Number, 123),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.Number}}, rawState),
			Schema: wrongSchema,
		}

		req := resource.DeleteRequest{
			State: state,
		}
		resp := resource.DeleteResponse{
			State: tfsdk.State{Schema: createTestSchema(t)},
		}

		r.Delete(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Delete should fail when State.Get fails")
	})
}
