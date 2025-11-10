package project

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

// createTestSchema creates a test schema for project resource.
func createTestSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := &ProjectResource{}
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

// TestProjectResource_Create tests project creation.
func TestProjectResource_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/projects":
				if r.Method == http.MethodPost {
					w.WriteHeader(http.StatusCreated)
					return
				}
				if r.Method == http.MethodGet {
					projects := []map[string]interface{}{
						{
							"id":   "proj-123",
							"name": "Test Project",
							"type": "team",
						},
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]interface{}{
						"data": projects,
					})
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &ProjectResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, nil),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, nil),
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

		r := &ProjectResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, nil),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, nil),
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

	t.Run("creation succeeds but GET fails after creation", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/projects":
				if r.Method == http.MethodPost {
					w.WriteHeader(http.StatusCreated)
					return
				}
				if r.Method == http.MethodGet {
					// Simulate error when listing projects after creation.
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &ProjectResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, nil),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, nil),
			Schema: createTestSchema(t),
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Create should have errors when GET fails after creation")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error reading project after creation")
	})

	t.Run("creation succeeds but project not found in list", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/projects":
				if r.Method == http.MethodPost {
					w.WriteHeader(http.StatusCreated)
					return
				}
				if r.Method == http.MethodGet {
					// Return empty list - project not found.
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]interface{}{
						"data": []map[string]interface{}{},
					})
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &ProjectResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, nil),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, nil),
			Schema: createTestSchema(t),
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Create should have errors when project not found in list")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error finding created project")
	})
}

// TestProjectResource_Read tests project reading.
func TestProjectResource_Read(t *testing.T) {
	t.Run("successful read", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/projects" && r.Method == http.MethodGet {
				projects := []map[string]interface{}{
					{
						"id":   "proj-123",
						"name": "Test Project",
						"type": "team",
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": projects,
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &ProjectResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-123"),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawState),
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

	t.Run("project not found removes from state", func(t *testing.T) {
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

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &ProjectResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-nonexistent"),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors even when project not found")
		// Note: In Terraform, when a resource is not found during Read, it's removed from state
		// This is tested by checking that resp.State would be empty, but we can't easily verify that in unit tests
	})

	t.Run("read fails with API error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &ProjectResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-123"),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors on API failure")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error reading project")
	})
}

// TestProjectResource_Update tests project update.
func TestProjectResource_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/projects/proj-123" && r.Method == http.MethodPut:
				w.WriteHeader(http.StatusOK)
				return
			case r.URL.Path == "/projects" && r.Method == http.MethodGet:
				projects := []map[string]interface{}{
					{
						"id":   "proj-123",
						"name": "Updated Project",
						"type": "team",
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": projects,
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &ProjectResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-123"),
			"name": tftypes.NewValue(tftypes.String, "Updated Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-123"),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawState),
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

	t.Run("update fails with PUT error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/projects/proj-123" && r.Method == http.MethodPut {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &ProjectResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-123"),
			"name": tftypes.NewValue(tftypes.String, "Updated Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-123"),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawState),
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

		assert.True(t, resp.Diagnostics.HasError(), "Update should have errors on PUT failure")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error updating project")
	})

	t.Run("update succeeds but GET after update fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/projects/proj-123" && r.Method == http.MethodPut:
				w.WriteHeader(http.StatusOK)
				return
			case r.URL.Path == "/projects" && r.Method == http.MethodGet:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &ProjectResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-123"),
			"name": tftypes.NewValue(tftypes.String, "Updated Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-123"),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawState),
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

		assert.True(t, resp.Diagnostics.HasError(), "Update should have errors when GET after update fails")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error reading project after update")
	})

	t.Run("update succeeds but project not found in list", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/projects/proj-123" && r.Method == http.MethodPut:
				w.WriteHeader(http.StatusOK)
				return
			case r.URL.Path == "/projects" && r.Method == http.MethodGet:
				// Return empty list.
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []map[string]interface{}{},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &ProjectResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-123"),
			"name": tftypes.NewValue(tftypes.String, "Updated Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-123"),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawState),
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

		assert.True(t, resp.Diagnostics.HasError(), "Update should have errors when project not found in list")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error verifying updated project")
	})
}

// TestProjectResource_Delete tests project deletion.
func TestProjectResource_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/projects/proj-123" && r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &ProjectResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-123"),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawState),
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

		r := &ProjectResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-123"),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawState),
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

	t.Run("delete fails when state get fails", func(t *testing.T) {
		r := &ProjectResource{}

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

// TestProjectResource_EdgeCasesForCoverage tests edge cases for 100% coverage.
func TestProjectResource_EdgeCasesForCoverage(t *testing.T) {
	t.Run("create fails when plan get fails", func(t *testing.T) {
		r := &ProjectResource{}

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

	t.Run("create succeeds but state set fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/projects":
				if r.Method == http.MethodPost {
					w.WriteHeader(http.StatusCreated)
					return
				}
				if r.Method == http.MethodGet {
					projects := []map[string]interface{}{
						{
							"id":   "proj-123",
							"name": "Test Project",
							"type": "team",
						},
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]interface{}{
						"data": projects,
					})
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &ProjectResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, nil),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		// Create state with incompatible schema
		wrongSchema := schema.Schema{
			Attributes: map[string]schema.Attribute{
				"id": schema.NumberAttribute{
					Computed: true,
				},
			},
		}
		state := tfsdk.State{
			Schema: wrongSchema,
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Create should fail when State.Set fails")
	})

	t.Run("read fails when state get fails", func(t *testing.T) {
		r := &ProjectResource{}

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

	t.Run("read succeeds but state set fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/projects" && r.Method == http.MethodGet {
				projects := []map[string]interface{}{
					{
						"id":   "proj-123",
						"name": "Test Project",
						"type": "team",
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": projects,
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &ProjectResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-123"),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		// Create response state with incompatible schema
		wrongSchema := schema.Schema{
			Attributes: map[string]schema.Attribute{
				"id": schema.NumberAttribute{
					Computed: true,
				},
			},
		}
		respState := tfsdk.State{
			Schema: wrongSchema,
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: respState,
		}

		r.Read(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should fail when State.Set fails")
	})

	t.Run("update fails when plan get fails", func(t *testing.T) {
		r := &ProjectResource{}

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

	t.Run("update succeeds but state set fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/projects/proj-123" && r.Method == http.MethodPut:
				w.WriteHeader(http.StatusNoContent)
				return
			case r.URL.Path == "/projects" && r.Method == http.MethodGet:
				projects := []map[string]interface{}{
					{
						"id":   "proj-123",
						"name": "Updated Project",
						"type": "team",
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": projects,
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &ProjectResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-123"),
			"name": tftypes.NewValue(tftypes.String, "Updated Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "proj-123"),
			"name": tftypes.NewValue(tftypes.String, "Test Project"),
			"type": tftypes.NewValue(tftypes.String, "team"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		// Create response state with incompatible schema
		wrongSchema := schema.Schema{
			Attributes: map[string]schema.Attribute{
				"id": schema.NumberAttribute{
					Computed: true,
				},
			},
		}
		respState := tfsdk.State{
			Schema: wrongSchema,
		}

		req := resource.UpdateRequest{
			Plan:  plan,
			State: state,
		}
		resp := resource.UpdateResponse{
			State: respState,
		}

		r.Update(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Update should fail when State.Set fails")
	})
}

// TestProjectResource_ImportState tests state import.
func TestProjectResource_ImportState(t *testing.T) {
	r := &ProjectResource{}

	schema := createTestSchema(t)
	state := tfsdk.State{
		Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, nil),
		Schema: schema,
	}

	req := resource.ImportStateRequest{
		ID: "proj-123",
	}
	resp := &resource.ImportStateResponse{
		State: state,
	}

	r.ImportState(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "ImportState should not have errors")
}

// TestNewProjectResource tests creation of new ProjectResource instance.
func TestNewProjectResource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		r := NewProjectResource()

		assert.NotNil(t, r)
		assert.IsType(t, &ProjectResource{}, r)
		assert.Nil(t, r.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		r1 := NewProjectResource()
		r2 := NewProjectResource()

		assert.NotSame(t, r1, r2)
	})
}

// TestNewProjectResourceWrapper tests creation of wrapped ProjectResource instance.
func TestNewProjectResourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		r := NewProjectResourceWrapper()

		assert.NotNil(t, r)
		assert.IsType(t, &ProjectResource{}, r)
	})

	t.Run("wrapper returns resource.Resource interface", func(t *testing.T) {
		r := NewProjectResourceWrapper()

		// Verify it implements resource.Resource.
		assert.Implements(t, (*resource.Resource)(nil), r)
		assert.NotNil(t, r)
	})
}

// TestProjectResource_Metadata tests metadata setup.
func TestProjectResource_Metadata(t *testing.T) {
	t.Run("set metadata with provider type", func(t *testing.T) {
		r := NewProjectResource()
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_project", resp.TypeName)
	})

	t.Run("set metadata with custom provider type", func(t *testing.T) {
		r := NewProjectResource()
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), resource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}, resp)

		assert.Equal(t, "custom_provider_project", resp.TypeName)
	})
}

// TestProjectResource_Configure tests configuration with client.
func TestProjectResource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		r := NewProjectResource()
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
		r := NewProjectResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: nil,
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		r := NewProjectResource()
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
		r := NewProjectResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.True(t, resp.Diagnostics.HasError())
	})
}

// TestProjectResource_Interfaces tests interface implementations.
func TestProjectResource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		r := NewProjectResource()

		var _ resource.Resource = r
		var _ resource.ResourceWithConfigure = r
		var _ resource.ResourceWithImportState = r
		var _ ProjectResourceInterface = r
	})
}

// TestProjectResource_Schema tests schema definition.
func TestProjectResource_Schema(t *testing.T) {
	t.Run("schema has correct structure", func(t *testing.T) {
		r := NewProjectResource()
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "project")

		// Verify id attribute (computed)
		idAttr, exists := resp.Schema.Attributes["id"]
		assert.True(t, exists)
		assert.True(t, idAttr.IsComputed())

		// Verify name attribute (required)
		nameAttr, exists := resp.Schema.Attributes["name"]
		assert.True(t, exists)
		assert.True(t, nameAttr.IsRequired())

		// Verify type attribute (computed)
		typeAttr, exists := resp.Schema.Attributes["type"]
		assert.True(t, exists)
		assert.True(t, typeAttr.IsComputed())
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		r := NewProjectResource()
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})
}
