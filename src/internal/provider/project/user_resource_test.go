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

func TestNewProjectUserResource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		r := NewProjectUserResource()

		assert.NotNil(t, r)
		assert.IsType(t, &ProjectUserResource{}, r)
		assert.Nil(t, r.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		r1 := NewProjectUserResource()
		r2 := NewProjectUserResource()

		assert.NotSame(t, r1, r2)
	})
}

func TestNewProjectUserResourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		r := NewProjectUserResourceWrapper()

		assert.NotNil(t, r)
		assert.IsType(t, &ProjectUserResource{}, r)
	})

	t.Run("wrapper returns resource.Resource interface", func(t *testing.T) {
		r := NewProjectUserResourceWrapper()

		// Verify it implements resource.Resource.
		assert.Implements(t, (*resource.Resource)(nil), r)
		assert.NotNil(t, r)
	})
}

func TestProjectUserResource_Metadata(t *testing.T) {
	t.Run("set metadata with provider type", func(t *testing.T) {
		r := NewProjectUserResource()
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_project_user", resp.TypeName)
	})

	t.Run("set metadata with custom provider type", func(t *testing.T) {
		r := NewProjectUserResource()
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), resource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}, resp)

		assert.Equal(t, "custom_provider_project_user", resp.TypeName)
	})
}

func TestProjectUserResource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		r := NewProjectUserResource()
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "user membership")

		expectedAttrs := []string{
			"id",
			"project_id",
			"user_id",
			"role",
		}

		for _, attr := range expectedAttrs {
			_, exists := resp.Schema.Attributes[attr]
			assert.True(t, exists, "Attribute %s should exist", attr)
		}
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		r := NewProjectUserResource()
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})
}

func TestProjectUserResource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		r := NewProjectUserResource()
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
		r := NewProjectUserResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: nil,
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		r := NewProjectUserResource()
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
		r := NewProjectUserResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.True(t, resp.Diagnostics.HasError())
	})
}

func TestProjectUserResource_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/projects/proj-123/users" && r.Method == http.MethodPost {
				w.WriteHeader(http.StatusCreated)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestUserResourceClient(t, handler)
		defer server.Close()

		r := &ProjectUserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
			Schema: createTestUserResourceSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, nil),
			Schema: createTestUserResourceSchema(t),
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

	t.Run("creation fails with API error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestUserResourceClient(t, handler)
		defer server.Close()

		r := &ProjectUserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
			Schema: createTestUserResourceSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, nil),
			Schema: createTestUserResourceSchema(t),
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Create should have errors")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error adding user to project")
	})
}

func TestProjectUserResource_Read(t *testing.T) {
	t.Run("successful read", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				userID := "user-456"
				role := "project:editor"
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []map[string]interface{}{
						{
							"id":    userID,
							"email": "test@example.com",
							"role":  role,
						},
					},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestUserResourceClient(t, handler)
		defer server.Close()

		r := &ProjectUserResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
			Schema: createTestUserResourceSchema(t),
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

	t.Run("user not found removes from state", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users" && r.Method == http.MethodGet &&
				r.URL.Query().Get("projectId") == "proj-123" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []map[string]interface{}{},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestUserResourceClient(t, handler)
		defer server.Close()

		r := &ProjectUserResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
			Schema: createTestUserResourceSchema(t),
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

	t.Run("read fails with API error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestUserResourceClient(t, handler)
		defer server.Close()

		r := &ProjectUserResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
			Schema: createTestUserResourceSchema(t),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error reading project users")
	})

	t.Run("read with nil data", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users" && r.Method == http.MethodGet &&
				r.URL.Query().Get("projectId") == "proj-123" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestUserResourceClient(t, handler)
		defer server.Close()

		r := &ProjectUserResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
			Schema: createTestUserResourceSchema(t),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should handle nil data gracefully")
	})
}

func TestProjectUserResource_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/projects/proj-123/users/user-456" && r.Method == http.MethodPatch {
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestUserResourceClient(t, handler)
		defer server.Close()

		r := &ProjectUserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:admin"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
			Schema: createTestUserResourceSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
			Schema: createTestUserResourceSchema(t),
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

	t.Run("update fails when changing project_id", func(t *testing.T) {
		r := &ProjectUserResource{}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-789/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-789"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
			Schema: createTestUserResourceSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
			Schema: createTestUserResourceSchema(t),
		}

		req := resource.UpdateRequest{
			Plan:  plan,
			State: state,
		}
		resp := resource.UpdateResponse{
			State: state,
		}

		r.Update(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Update should have errors when changing project_id")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Project/User Change Not Supported")
	})

	t.Run("update fails when changing user_id", func(t *testing.T) {
		r := &ProjectUserResource{}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-789"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-789"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
			Schema: createTestUserResourceSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
			Schema: createTestUserResourceSchema(t),
		}

		req := resource.UpdateRequest{
			Plan:  plan,
			State: state,
		}
		resp := resource.UpdateResponse{
			State: state,
		}

		r.Update(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Update should have errors when changing user_id")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Project/User Change Not Supported")
	})

	t.Run("update fails with API error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestUserResourceClient(t, handler)
		defer server.Close()

		r := &ProjectUserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:admin"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
			Schema: createTestUserResourceSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
			Schema: createTestUserResourceSchema(t),
		}

		req := resource.UpdateRequest{
			Plan:  plan,
			State: state,
		}
		resp := resource.UpdateResponse{
			State: state,
		}

		r.Update(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Update should have errors")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error updating user role in project")
	})

	t.Run("update succeeds when role unchanged", func(t *testing.T) {
		r := &ProjectUserResource{}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
			Schema: createTestUserResourceSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
			Schema: createTestUserResourceSchema(t),
		}

		req := resource.UpdateRequest{
			Plan:  plan,
			State: state,
		}
		resp := resource.UpdateResponse{
			State: state,
		}

		r.Update(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Update should not have errors when role is unchanged")
	})
}

func TestProjectUserResource_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/projects/proj-123/users/user-456" && r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestUserResourceClient(t, handler)
		defer server.Close()

		r := &ProjectUserResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
			Schema: createTestUserResourceSchema(t),
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

	t.Run("delete fails with API error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestUserResourceClient(t, handler)
		defer server.Close()

		r := &ProjectUserResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
			Schema: createTestUserResourceSchema(t),
		}

		req := resource.DeleteRequest{
			State: state,
		}
		resp := resource.DeleteResponse{
			State: state,
		}

		r.Delete(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Delete should have errors")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error removing user from project")
	})

	t.Run("delete fails when state get fails", func(t *testing.T) {
		r := &ProjectUserResource{}

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
			State: tfsdk.State{Schema: createTestUserResourceSchema(t)},
		}

		r.Delete(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Delete should fail when State.Get fails")
	})
}

// TestProjectUserResource_EdgeCasesForCoverage tests edge cases for 100% coverage.
func TestProjectUserResource_EdgeCasesForCoverage(t *testing.T) {
	t.Run("create fails when plan get fails", func(t *testing.T) {
		r := &ProjectUserResource{}

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
			State: tfsdk.State{Schema: createTestUserResourceSchema(t)},
		}

		r.Create(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Create should fail when Plan.Get fails")
	})

	t.Run("create succeeds but state set fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/projects/proj-123/users/user-456" && r.Method == http.MethodPut {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestUserResourceClient(t, handler)
		defer server.Close()

		r := &ProjectUserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
			Schema: createTestUserResourceSchema(t),
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
		r := &ProjectUserResource{}

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
			State: tfsdk.State{Schema: createTestUserResourceSchema(t)},
		}

		r.Read(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should fail when State.Get fails")
	})

	t.Run("read succeeds but state set fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users" && r.Method == http.MethodGet && r.URL.Query().Get("projectId") == "proj-123" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []map[string]interface{}{
						{
							"id":    "user-456",
							"email": "test@example.com",
							"role":  "project:editor",
						},
					},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestUserResourceClient(t, handler)
		defer server.Close()

		r := &ProjectUserResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
			Schema: createTestUserResourceSchema(t),
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
		r := &ProjectUserResource{}

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
			Schema: createTestUserResourceSchema(t),
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
			if r.URL.Path == "/projects/proj-123/users/user-456" && r.Method == http.MethodPut {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestUserResourceClient(t, handler)
		defer server.Close()

		r := &ProjectUserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:admin"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
			Schema: createTestUserResourceSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
			"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
			"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
			"role":       tftypes.NewValue(tftypes.String, "project:editor"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
			Schema: createTestUserResourceSchema(t),
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

func TestProjectUserResource_ImportState(t *testing.T) {
	t.Run("successful import", func(t *testing.T) {
		r := &ProjectUserResource{}

		schema := createTestUserResourceSchema(t)
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, nil),
			Schema: schema,
		}

		req := resource.ImportStateRequest{
			ID: "proj-123/user-456",
		}
		resp := &resource.ImportStateResponse{
			State: state,
		}

		r.ImportState(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError(), "ImportState should not have errors")
	})

	t.Run("import fails with invalid ID format", func(t *testing.T) {
		r := &ProjectUserResource{}

		schema := createTestUserResourceSchema(t)
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, nil),
			Schema: schema,
		}

		req := resource.ImportStateRequest{
			ID: "invalid-id",
		}
		resp := &resource.ImportStateResponse{
			State: state,
		}

		r.ImportState(context.Background(), req, resp)

		assert.True(t, resp.Diagnostics.HasError(), "ImportState should have errors with invalid ID")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid Import ID")
	})

	t.Run("import fails with too many parts", func(t *testing.T) {
		r := &ProjectUserResource{}

		schema := createTestUserResourceSchema(t)
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, nil),
			Schema: schema,
		}

		req := resource.ImportStateRequest{
			ID: "proj-123/user-456/extra",
		}
		resp := &resource.ImportStateResponse{
			State: state,
		}

		r.ImportState(context.Background(), req, resp)

		assert.True(t, resp.Diagnostics.HasError(), "ImportState should have errors with too many parts")
	})
}

func TestProjectUserResource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		r := NewProjectUserResource()

		var _ resource.Resource = r
		var _ resource.ResourceWithConfigure = r
		var _ resource.ResourceWithImportState = r
		var _ ProjectUserResourceInterface = r
	})
}

// createTestUserResourceSchema creates a test schema for project user resource.
func createTestUserResourceSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := &ProjectUserResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), req, resp)
	return resp.Schema
}

// setupTestUserResourceClient creates a test N8nClient with httptest server.
func setupTestUserResourceClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
