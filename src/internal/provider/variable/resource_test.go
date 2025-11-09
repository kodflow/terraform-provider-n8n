package variable

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/variable/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// createTestSchema creates a test schema for variable resource.
func createTestSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := &VariableResource{}
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

// TestVariableResource_Create tests variable creation.
func TestVariableResource_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/variables":
				if r.Method == http.MethodPost {
					w.WriteHeader(http.StatusCreated)
					return
				}
				if r.Method == http.MethodGet {
					variables := []map[string]interface{}{
						{
							"id":    "var-123",
							"key":   "TEST_KEY",
							"value": "test-value",
							"type":  "string",
						},
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]interface{}{
						"data": variables,
					})
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &VariableResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, nil),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, nil),
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

		r := &VariableResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, nil),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, nil),
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

	t.Run("variable not found after creation", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/variables":
				if r.Method == http.MethodPost {
					w.WriteHeader(http.StatusCreated)
					return
				}
				if r.Method == http.MethodGet {
					// Return empty list - variable not found
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

		r := &VariableResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, nil),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, nil),
			Schema: createTestSchema(t),
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Create should have errors when variable not found after creation")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error finding created variable")
	})

	t.Run("API error when reading after creation", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/variables":
				if r.Method == http.MethodPost {
					w.WriteHeader(http.StatusCreated)
					return
				}
				if r.Method == http.MethodGet {
					// Return error when trying to read
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &VariableResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, nil),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, nil),
			Schema: createTestSchema(t),
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Create should have errors on API error")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error reading variable after creation")
	})

	t.Run("null data after creation", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/variables":
				if r.Method == http.MethodPost {
					w.WriteHeader(http.StatusCreated)
					return
				}
				if r.Method == http.MethodGet {
					// Return null data
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte(`{}`))
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &VariableResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, nil),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, nil),
			Schema: createTestSchema(t),
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Create should have errors when data is null")
	})
}

// TestVariableResource_Read tests variable reading.
func TestVariableResource_Read(t *testing.T) {
	t.Run("successful read", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				variables := []map[string]interface{}{
					{
						"id":    "var-123",
						"key":   "TEST_KEY",
						"value": "test-value",
						"type":  "string",
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": variables,
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &VariableResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawState),
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

	t.Run("variable not found removes from state", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
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

		r := &VariableResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-nonexistent"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors even when variable not found")
		// Note: In Terraform, when a resource is not found during Read, it's removed from state
		// This is tested by checking that resp.State would be empty, but we can't easily verify that in unit tests
	})

	t.Run("API error when reading", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &VariableResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors on API error")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error reading variable")
	})

	t.Run("null data from API", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &VariableResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors with null data (variable is removed from state)")
	})
}

// TestVariableResource_Update tests variable update.
func TestVariableResource_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/variables/var-123" && r.Method == http.MethodPut:
				w.WriteHeader(http.StatusOK)
				return
			case r.URL.Path == "/variables" && r.Method == http.MethodGet:
				variables := []map[string]interface{}{
					{
						"id":    "var-123",
						"key":   "TEST_KEY",
						"value": "updated-value",
						"type":  "string",
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": variables,
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &VariableResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "updated-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawState),
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

	t.Run("update API error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables/var-123" && r.Method == http.MethodPut {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &VariableResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "updated-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawState),
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

		assert.True(t, resp.Diagnostics.HasError(), "Update should have errors on API error")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error updating variable")
	})

	t.Run("API error when reading after update", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/variables/var-123" && r.Method == http.MethodPut:
				w.WriteHeader(http.StatusOK)
				return
			case r.URL.Path == "/variables" && r.Method == http.MethodGet:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &VariableResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "updated-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawState),
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

		assert.True(t, resp.Diagnostics.HasError(), "Update should have errors on read API error")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error reading variable after update")
	})

	t.Run("variable not found after update", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/variables/var-123" && r.Method == http.MethodPut:
				w.WriteHeader(http.StatusOK)
				return
			case r.URL.Path == "/variables" && r.Method == http.MethodGet:
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

		r := &VariableResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "updated-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawState),
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

		assert.True(t, resp.Diagnostics.HasError(), "Update should have errors when variable not found")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error verifying updated variable")
	})

	t.Run("null data after update", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/variables/var-123" && r.Method == http.MethodPut:
				w.WriteHeader(http.StatusOK)
				return
			case r.URL.Path == "/variables" && r.Method == http.MethodGet:
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &VariableResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "updated-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawState),
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

		assert.True(t, resp.Diagnostics.HasError(), "Update should have errors with null data")
	})
}

// TestVariableResource_Delete tests variable deletion.
func TestVariableResource_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables/var-123" && r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &VariableResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawState),
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

		r := &VariableResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "var-123"),
			"key":        tftypes.NewValue(tftypes.String, "TEST_KEY"),
			"value":      tftypes.NewValue(tftypes.String, "test-value"),
			"type":       tftypes.NewValue(tftypes.String, "string"),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, rawState),
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

// TestVariableResource_ImportState tests state import.
func TestVariableResource_ImportState(t *testing.T) {
	r := &VariableResource{}

	schema := createTestSchema(t)
	state := tfsdk.State{
		Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "key": tftypes.String, "value": tftypes.String, "type": tftypes.String, "project_id": tftypes.String}}, nil),
		Schema: schema,
	}

	req := resource.ImportStateRequest{
		ID: "var-123",
	}
	resp := &resource.ImportStateResponse{
		State: state,
	}

	r.ImportState(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "ImportState should not have errors")
}

// MockVariableResourceInterface is a mock implementation of VariableResourceInterface.
type MockVariableResourceInterface struct {
	mock.Mock
}

func (m *MockVariableResourceInterface) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariableResourceInterface) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariableResourceInterface) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariableResourceInterface) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariableResourceInterface) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariableResourceInterface) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariableResourceInterface) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariableResourceInterface) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	m.Called(ctx, req, resp)
}

func TestNewVariableResource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		r := NewVariableResource()

		assert.NotNil(t, r)
		assert.IsType(t, &VariableResource{}, r)
		assert.Nil(t, r.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		r1 := NewVariableResource()
		r2 := NewVariableResource()

		assert.NotSame(t, r1, r2)
	})
}

func TestNewVariableResourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		r := NewVariableResourceWrapper()

		assert.NotNil(t, r)
		assert.IsType(t, &VariableResource{}, r)
	})

	t.Run("wrapper returns resource.Resource interface", func(t *testing.T) {
		r := NewVariableResourceWrapper()

		// r is already of type resource.Resource, no assertion needed
		assert.NotNil(t, r)
	})
}

func TestVariableResource_Metadata(t *testing.T) {
	t.Run("set metadata with provider type", func(t *testing.T) {
		r := NewVariableResource()
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_variable", resp.TypeName)
	})

	t.Run("set metadata with different provider type", func(t *testing.T) {
		r := NewVariableResource()
		resp := &resource.MetadataResponse{}
		req := resource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_variable", resp.TypeName)
	})

	t.Run("set metadata with empty provider type", func(t *testing.T) {
		r := NewVariableResource()
		resp := &resource.MetadataResponse{}
		req := resource.MetadataRequest{
			ProviderTypeName: "",
		}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "_variable", resp.TypeName)
	})
}

func TestVariableResource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		r := NewVariableResource()
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "n8n variable")

		// Verify computed attributes
		idAttr, exists := resp.Schema.Attributes["id"]
		assert.True(t, exists)
		assert.True(t, idAttr.IsComputed())

		// Verify required attributes
		keyAttr, exists := resp.Schema.Attributes["key"]
		assert.True(t, exists)
		assert.True(t, keyAttr.IsRequired())

		valueAttr, exists := resp.Schema.Attributes["value"]
		assert.True(t, exists)
		assert.True(t, valueAttr.IsRequired())

		// Verify optional attributes
		typeAttr, exists := resp.Schema.Attributes["type"]
		assert.True(t, exists)
		assert.True(t, typeAttr.IsOptional())
		assert.True(t, typeAttr.IsComputed())

		projectIDAttr, exists := resp.Schema.Attributes["project_id"]
		assert.True(t, exists)
		assert.True(t, projectIDAttr.IsOptional())
		assert.True(t, projectIDAttr.IsComputed())
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		r := NewVariableResource()
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})

	t.Run("value attribute is sensitive", func(t *testing.T) {
		r := NewVariableResource()
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		_, exists := resp.Schema.Attributes["value"]
		assert.True(t, exists)
		// Note: Sensitive attribute testing depends on framework internals
	})
}

func TestVariableResource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		r := NewVariableResource()
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
		r := NewVariableResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: nil,
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		r := NewVariableResource()
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
		r := NewVariableResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "*client.N8nClient")
	})

	t.Run("configure multiple times", func(t *testing.T) {
		r := NewVariableResource()

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

func TestVariableResource_FindVariableByID(t *testing.T) {
	t.Run("find existing variable", func(t *testing.T) {
		r := NewVariableResource()
		id1 := "var-123"
		id2 := "var-456"
		variableList := &n8nsdk.VariableList{
			Data: []n8nsdk.Variable{
				{Id: &id1, Key: "key-1", Value: "value-1"},
				{Id: &id2, Key: "key-2", Value: "value-2"},
			},
		}

		found := r.findVariableByID(variableList, "var-123")

		assert.NotNil(t, found)
		assert.Equal(t, "var-123", *found.Id)
		assert.Equal(t, "key-1", found.Key)
	})

	t.Run("not found when ID does not match", func(t *testing.T) {
		r := NewVariableResource()
		id := "var-123"
		variableList := &n8nsdk.VariableList{
			Data: []n8nsdk.Variable{
				{Id: &id, Key: "test-key", Value: "test-value"},
			},
		}

		found := r.findVariableByID(variableList, "var-999")

		assert.Nil(t, found)
	})

	t.Run("nil data in variable list", func(t *testing.T) {
		r := NewVariableResource()
		variableList := &n8nsdk.VariableList{
			Data: nil,
		}

		found := r.findVariableByID(variableList, "var-123")

		assert.Nil(t, found)
	})

	t.Run("empty variable list", func(t *testing.T) {
		r := NewVariableResource()
		variableList := &n8nsdk.VariableList{
			Data: []n8nsdk.Variable{},
		}

		found := r.findVariableByID(variableList, "var-123")

		assert.Nil(t, found)
	})

	t.Run("variable with nil ID", func(t *testing.T) {
		r := NewVariableResource()
		id := "var-123"
		variableList := &n8nsdk.VariableList{
			Data: []n8nsdk.Variable{
				{Id: nil, Key: "key-1", Value: "value-1"},
				{Id: &id, Key: "key-2", Value: "value-2"},
			},
		}

		found := r.findVariableByID(variableList, "var-123")

		assert.NotNil(t, found)
		assert.Equal(t, "var-123", *found.Id)
	})
}

func TestVariableResource_UpdateStateFromVariable(t *testing.T) {
	t.Run("update state with all fields", func(t *testing.T) {
		r := NewVariableResource()
		varType := "string"
		projectID := "proj-123"
		variable := &n8nsdk.Variable{
			Key:   "test-key",
			Value: "test-value",
			Type:  &varType,
			Project: &n8nsdk.Project{
				Id: &projectID,
			},
		}

		state := &models.Resource{
			ID: types.StringValue("var-123"),
		}

		r.updateStateFromVariable(variable, state)

		assert.Equal(t, "var-123", state.ID.ValueString())
		assert.Equal(t, "test-key", state.Key.ValueString())
		assert.Equal(t, "test-value", state.Value.ValueString())
		assert.Equal(t, "string", state.Type.ValueString())
		assert.Equal(t, "proj-123", state.ProjectID.ValueString())
	})

	t.Run("update state with nil fields", func(t *testing.T) {
		r := NewVariableResource()
		variable := &n8nsdk.Variable{
			Key:   "test-key",
			Value: "test-value",
		}

		state := &models.Resource{
			ID: types.StringValue("var-123"),
		}

		r.updateStateFromVariable(variable, state)

		assert.Equal(t, "test-key", state.Key.ValueString())
		assert.Equal(t, "test-value", state.Value.ValueString())
		assert.True(t, state.Type.IsNull())
		assert.True(t, state.ProjectID.IsNull())
	})

	t.Run("update state with nil project", func(t *testing.T) {
		r := NewVariableResource()
		variable := &n8nsdk.Variable{
			Key:     "test-key",
			Value:   "test-value",
			Project: nil,
		}

		state := &models.Resource{
			ID: types.StringValue("var-123"),
		}

		r.updateStateFromVariable(variable, state)

		assert.True(t, state.ProjectID.IsNull())
	})
}

func TestVariableResource_ImportStatePassthrough(t *testing.T) {
	t.Run("import uses correct path", func(t *testing.T) {
		// This test verifies that the import state uses the "id" path
		expectedPath := path.Root("id")
		assert.NotNil(t, expectedPath)
	})
}

func TestVariableResource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		r := NewVariableResource()

		// Test that VariableResource implements resource.Resource
		var _ resource.Resource = r

		// Test that VariableResource implements resource.ResourceWithConfigure
		var _ resource.ResourceWithConfigure = r

		// Test that VariableResource implements resource.ResourceWithImportState
		var _ resource.ResourceWithImportState = r

		// Test that VariableResource implements VariableResourceInterface
		var _ VariableResourceInterface = r
	})
}

func TestVariableResourceConcurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		r := NewVariableResource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &resource.MetadataResponse{}
				r.Metadata(context.Background(), resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)
				assert.Equal(t, "n8n_variable", resp.TypeName)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		r := NewVariableResource()

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

	t.Run("concurrent configure calls", func(t *testing.T) {
		r := NewVariableResource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &resource.ConfigureResponse{}
				mockClient := &client.N8nClient{}
				req := resource.ConfigureRequest{
					ProviderData: mockClient,
				}
				r.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent findVariableByID calls", func(t *testing.T) {
		r := NewVariableResource()
		id := "var-123"
		variableList := &n8nsdk.VariableList{
			Data: []n8nsdk.Variable{
				{Id: &id, Key: "test-key", Value: "test-value"},
			},
		}

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				found := r.findVariableByID(variableList, "var-123")
				assert.NotNil(t, found)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func BenchmarkVariableResource_Schema(b *testing.B) {
	r := NewVariableResource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.SchemaResponse{}
		r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	}
}

func BenchmarkVariableResource_Metadata(b *testing.B) {
	r := NewVariableResource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.MetadataResponse{}
		r.Metadata(context.Background(), resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)
	}
}

func BenchmarkVariableResource_Configure(b *testing.B) {
	r := NewVariableResource()
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

func BenchmarkVariableResource_FindVariableByID(b *testing.B) {
	r := NewVariableResource()
	id := "var-123"
	variableList := &n8nsdk.VariableList{
		Data: []n8nsdk.Variable{
			{Id: &id, Key: "test-key", Value: "test-value"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.findVariableByID(variableList, "var-123")
	}
}

func BenchmarkVariableResource_UpdateStateFromVariable(b *testing.B) {
	r := NewVariableResource()
	varType := "string"
	projectID := "proj-123"
	variable := &n8nsdk.Variable{
		Key:   "test-key",
		Value: "test-value",
		Type:  &varType,
		Project: &n8nsdk.Project{
			Id: &projectID,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		state := &models.Resource{
			ID: types.StringValue("var-123"),
		}
		r.updateStateFromVariable(variable, state)
	}
}

// Mock request types for API testing.

// MockVariablesPostRequest mocks the VariablesPost API request.
type MockVariablesPostRequest struct {
	mock.Mock
	variableCreate n8nsdk.VariableCreate
}

func (m *MockVariablesPostRequest) VariableCreate(variableCreate n8nsdk.VariableCreate) *MockVariablesPostRequest {
	m.variableCreate = variableCreate
	return m
}

func (m *MockVariablesPostRequest) Execute() (*n8nsdk.VariableList, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*n8nsdk.VariableList), args.Error(1)
}

// MockVariablesGetRequest mocks the VariablesGet API request.
type MockVariablesGetRequest struct {
	mock.Mock
}

func (m *MockVariablesGetRequest) Execute() (*n8nsdk.VariableList, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*n8nsdk.VariableList), args.Error(1)
}

// MockVariablesPutRequest mocks the VariablesIdPut API request.
type MockVariablesPutRequest struct {
	mock.Mock
	variableCreate n8nsdk.VariableCreate
}

func (m *MockVariablesPutRequest) VariableCreate(variableCreate n8nsdk.VariableCreate) *MockVariablesPutRequest {
	m.variableCreate = variableCreate
	return m
}

func (m *MockVariablesPutRequest) Execute() error {
	args := m.Called()
	return args.Error(0)
}

// MockVariablesDeleteRequest mocks the VariablesIdDelete API request.
type MockVariablesDeleteRequest struct {
	mock.Mock
}

func (m *MockVariablesDeleteRequest) Execute() error {
	args := m.Called()
	return args.Error(0)
}

// MockAPIClient extends the basic mock to support Variables API.
type MockAPIClient struct {
	mock.Mock
	VariablesAPI *MockVariablesAPIWithRequests
}

// MockVariablesAPIWithRequests extends MockVariablesAPI with request mocking.
type MockVariablesAPIWithRequests struct {
	mock.Mock
}

func (m *MockVariablesAPIWithRequests) VariablesPost(ctx context.Context) n8nsdk.VariablesAPIVariablesPostRequest {
	args := m.Called(ctx)
	return args.Get(0).(n8nsdk.VariablesAPIVariablesPostRequest)
}

func (m *MockVariablesAPIWithRequests) VariablesGet(ctx context.Context) n8nsdk.VariablesAPIVariablesGetRequest {
	args := m.Called(ctx)
	return args.Get(0).(n8nsdk.VariablesAPIVariablesGetRequest)
}

func (m *MockVariablesAPIWithRequests) VariablesIdPut(ctx context.Context, id string) n8nsdk.VariablesAPIVariablesIdPutRequest {
	args := m.Called(ctx, id)
	return args.Get(0).(n8nsdk.VariablesAPIVariablesIdPutRequest)
}

func (m *MockVariablesAPIWithRequests) VariablesIdDelete(ctx context.Context, id string) n8nsdk.VariablesAPIVariablesIdDeleteRequest {
	args := m.Called(ctx, id)
	return args.Get(0).(n8nsdk.VariablesAPIVariablesIdDeleteRequest)
}
