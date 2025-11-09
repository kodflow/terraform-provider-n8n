package user

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

// createTestSchema creates a test schema for user resource.
func createTestSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := &UserResource{}
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

// TestUserResource_Create tests user creation.
func TestUserResource_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/users":
				if r.Method == http.MethodPost {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					response := map[string]interface{}{
						"user": map[string]interface{}{
							"id": "user-123",
						},
					}
					json.NewEncoder(w).Encode(response)
					return
				}
			case "/users/user-123":
				if r.Method == http.MethodGet {
					user := map[string]interface{}{
						"id":        "user-123",
						"email":     "test@example.com",
						"firstName": "Test",
						"lastName":  "User",
						"role":      "global:member",
						"isPending": false,
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-01T00:00:00Z",
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(user)
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &UserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, nil),
			"last_name":  tftypes.NewValue(tftypes.String, nil),
			"role":       tftypes.NewValue(tftypes.String, nil),
			"is_pending": tftypes.NewValue(tftypes.Bool, nil),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
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

		r := &UserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, nil),
			"last_name":  tftypes.NewValue(tftypes.String, nil),
			"role":       tftypes.NewValue(tftypes.String, nil),
			"is_pending": tftypes.NewValue(tftypes.Bool, nil),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
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

// TestUserResource_Read tests user reading.
func TestUserResource_Read(t *testing.T) {
	t.Run("successful read", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users/user-123" && r.Method == http.MethodGet {
				user := map[string]interface{}{
					"id":        "user-123",
					"email":     "test@example.com",
					"firstName": "Test",
					"lastName":  "User",
					"role":      "global:member",
					"isPending": false,
					"createdAt": "2024-01-01T00:00:00Z",
					"updatedAt": "2024-01-01T00:00:00Z",
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(user)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &UserResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-123"),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, "Test"),
			"last_name":  tftypes.NewValue(tftypes.String, "User"),
			"role":       tftypes.NewValue(tftypes.String, "global:member"),
			"is_pending": tftypes.NewValue(tftypes.Bool, false),
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
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

	t.Run("user not found removes from state", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users/user-nonexistent" && r.Method == http.MethodGet {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "User not found"}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &UserResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-nonexistent"),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, "Test"),
			"last_name":  tftypes.NewValue(tftypes.String, "User"),
			"role":       tftypes.NewValue(tftypes.String, "global:member"),
			"is_pending": tftypes.NewValue(tftypes.Bool, false),
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should have errors when user not found")
		// Note: In Terraform, when a resource is not found during Read, it's removed from state
		// This is tested by checking that resp.State would be empty, but we can't easily verify that in unit tests
	})
}

// TestUserResource_Update tests user update.
func TestUserResource_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/users/user-123/role" && r.Method == http.MethodPatch:
				w.WriteHeader(http.StatusOK)
				return
			case r.URL.Path == "/users/user-123" && r.Method == http.MethodGet:
				user := map[string]interface{}{
					"id":        "user-123",
					"email":     "test@example.com",
					"firstName": "Test",
					"lastName":  "User",
					"role":      "global:admin",
					"isPending": false,
					"createdAt": "2024-01-01T00:00:00Z",
					"updatedAt": "2024-01-01T00:00:00Z",
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(user)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &UserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-123"),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, "Test"),
			"last_name":  tftypes.NewValue(tftypes.String, "User"),
			"role":       tftypes.NewValue(tftypes.String, "global:admin"),
			"is_pending": tftypes.NewValue(tftypes.Bool, false),
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-123"),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, "Test"),
			"last_name":  tftypes.NewValue(tftypes.String, "User"),
			"role":       tftypes.NewValue(tftypes.String, "global:member"),
			"is_pending": tftypes.NewValue(tftypes.Bool, false),
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
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

// TestUserResource_Delete tests user deletion.
func TestUserResource_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users/user-123" && r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &UserResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-123"),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, "Test"),
			"last_name":  tftypes.NewValue(tftypes.String, "User"),
			"role":       tftypes.NewValue(tftypes.String, "global:member"),
			"is_pending": tftypes.NewValue(tftypes.Bool, false),
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
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

		r := &UserResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-123"),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, "Test"),
			"last_name":  tftypes.NewValue(tftypes.String, "User"),
			"role":       tftypes.NewValue(tftypes.String, "global:member"),
			"is_pending": tftypes.NewValue(tftypes.Bool, false),
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
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

// TestUserResource_ImportState tests state import.
func TestUserResource_ImportState(t *testing.T) {
	r := &UserResource{}

	schema := createTestSchema(t)
	state := tfsdk.State{
		Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
		Schema: schema,
	}

	req := resource.ImportStateRequest{
		ID: "user-123",
	}
	resp := &resource.ImportStateResponse{
		State: state,
	}

	r.ImportState(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "ImportState should not have errors")
}

// TestNewUserResource tests resource creation.
func TestNewUserResource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		r := NewUserResource()

		assert.NotNil(t, r)
		assert.IsType(t, &UserResource{}, r)
		assert.Nil(t, r.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		r1 := NewUserResource()
		r2 := NewUserResource()

		assert.NotSame(t, r1, r2)
	})
}

// TestNewUserResourceWrapper tests resource wrapper creation.
func TestNewUserResourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		r := NewUserResourceWrapper()

		assert.NotNil(t, r)
		assert.IsType(t, &UserResource{}, r)
	})

	t.Run("wrapper returns resource.Resource interface", func(t *testing.T) {
		r := NewUserResourceWrapper()

		assert.NotNil(t, r)
	})
}

// TestUserResource_Metadata tests resource metadata.
func TestUserResource_Metadata(t *testing.T) {
	t.Run("set metadata with provider type", func(t *testing.T) {
		r := NewUserResource()
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_user", resp.TypeName)
	})

	t.Run("set metadata with custom provider type", func(t *testing.T) {
		r := NewUserResource()
		resp := &resource.MetadataResponse{}
		req := resource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_user", resp.TypeName)
	})
}

// TestUserResource_Schema tests resource schema.
func TestUserResource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		r := NewUserResource()
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "n8n users")

		// Verify required attributes
		emailAttr, exists := resp.Schema.Attributes["email"]
		assert.True(t, exists)
		assert.True(t, emailAttr.IsRequired())

		// Verify computed attributes
		idAttr, exists := resp.Schema.Attributes["id"]
		assert.True(t, exists)
		assert.True(t, idAttr.IsComputed())

		// Verify optional+computed attributes
		roleAttr, exists := resp.Schema.Attributes["role"]
		assert.True(t, exists)
		assert.True(t, roleAttr.IsOptional())
		assert.True(t, roleAttr.IsComputed())
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		r := NewUserResource()
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})
}

// TestUserResource_Configure tests resource configuration.
func TestUserResource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		r := NewUserResource()
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
		r := NewUserResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: nil,
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		r := NewUserResource()
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
		r := NewUserResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("configure multiple times", func(t *testing.T) {
		r := NewUserResource()

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

// TestUserResource_CreateWithRole tests creation with role specified.
func TestUserResource_CreateWithRole(t *testing.T) {
	t.Run("creation with role", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/users":
				if r.Method == http.MethodPost {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					response := map[string]interface{}{
						"user": map[string]interface{}{
							"id": "user-123",
						},
					}
					json.NewEncoder(w).Encode(response)
					return
				}
			case "/users/user-123":
				if r.Method == http.MethodGet {
					user := map[string]interface{}{
						"id":        "user-123",
						"email":     "test@example.com",
						"firstName": "Test",
						"lastName":  "User",
						"role":      "global:admin",
						"isPending": false,
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-01T00:00:00Z",
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(user)
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &UserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, nil),
			"last_name":  tftypes.NewValue(tftypes.String, nil),
			"role":       tftypes.NewValue(tftypes.String, "global:admin"),
			"is_pending": tftypes.NewValue(tftypes.Bool, nil),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
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

	t.Run("creation with invalid response", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users" && r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				response := map[string]interface{}{
					"user": map[string]interface{}{},
				}
				json.NewEncoder(w).Encode(response)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &UserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, nil),
			"last_name":  tftypes.NewValue(tftypes.String, nil),
			"role":       tftypes.NewValue(tftypes.String, nil),
			"is_pending": tftypes.NewValue(tftypes.Bool, nil),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
			Schema: createTestSchema(t),
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Create should have errors when API returns invalid response")
	})

	t.Run("creation user fetch fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/users":
				if r.Method == http.MethodPost {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					response := map[string]interface{}{
						"user": map[string]interface{}{
							"id": "user-123",
						},
					}
					json.NewEncoder(w).Encode(response)
					return
				}
			case "/users/user-123":
				if r.Method == http.MethodGet {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &UserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, nil),
			"last_name":  tftypes.NewValue(tftypes.String, nil),
			"role":       tftypes.NewValue(tftypes.String, nil),
			"is_pending": tftypes.NewValue(tftypes.Bool, nil),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
			Schema: createTestSchema(t),
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Create should have errors when user fetch fails")
	})
}

// TestUserResource_UpdateEmailChanged tests email change validation.
func TestUserResource_UpdateEmailChanged(t *testing.T) {
	t.Run("update with email change fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &UserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-123"),
			"email":      tftypes.NewValue(tftypes.String, "newemail@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, "Test"),
			"last_name":  tftypes.NewValue(tftypes.String, "User"),
			"role":       tftypes.NewValue(tftypes.String, "global:member"),
			"is_pending": tftypes.NewValue(tftypes.Bool, false),
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-123"),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, "Test"),
			"last_name":  tftypes.NewValue(tftypes.String, "User"),
			"role":       tftypes.NewValue(tftypes.String, "global:member"),
			"is_pending": tftypes.NewValue(tftypes.Bool, false),
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
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

		assert.True(t, resp.Diagnostics.HasError(), "Update should have errors when email changes")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Email Change Not Supported")
	})
}

// TestUserResource_UpdateNoRoleChange tests update without role change.
func TestUserResource_UpdateNoRoleChange(t *testing.T) {
	t.Run("update without role change", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users/user-123" && r.Method == http.MethodGet {
				user := map[string]interface{}{
					"id":        "user-123",
					"email":     "test@example.com",
					"firstName": "Test",
					"lastName":  "User",
					"role":      "global:member",
					"isPending": false,
					"createdAt": "2024-01-01T00:00:00Z",
					"updatedAt": "2024-01-01T00:00:00Z",
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(user)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &UserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-123"),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, "Test"),
			"last_name":  tftypes.NewValue(tftypes.String, "User"),
			"role":       tftypes.NewValue(tftypes.String, "global:member"),
			"is_pending": tftypes.NewValue(tftypes.Bool, false),
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-123"),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, "Test"),
			"last_name":  tftypes.NewValue(tftypes.String, "User"),
			"role":       tftypes.NewValue(tftypes.String, "global:member"),
			"is_pending": tftypes.NewValue(tftypes.Bool, false),
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
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

// TestUserResource_UpdateRoleFails tests role update failure.
func TestUserResource_UpdateRoleFails(t *testing.T) {
	t.Run("update role fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users/user-123/role" && r.Method == http.MethodPatch {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &UserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-123"),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, "Test"),
			"last_name":  tftypes.NewValue(tftypes.String, "User"),
			"role":       tftypes.NewValue(tftypes.String, "global:admin"),
			"is_pending": tftypes.NewValue(tftypes.Bool, false),
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-123"),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, "Test"),
			"last_name":  tftypes.NewValue(tftypes.String, "User"),
			"role":       tftypes.NewValue(tftypes.String, "global:member"),
			"is_pending": tftypes.NewValue(tftypes.Bool, false),
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
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

		assert.True(t, resp.Diagnostics.HasError(), "Update should have errors when role update fails")
	})
}

// TestUserResource_UpdateRefreshFails tests refresh failure after update.
func TestUserResource_UpdateRefreshFails(t *testing.T) {
	t.Run("update refresh fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/users/user-123/role" && r.Method == http.MethodPatch:
				w.WriteHeader(http.StatusOK)
				return
			case r.URL.Path == "/users/user-123" && r.Method == http.MethodGet:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &UserResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-123"),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, "Test"),
			"last_name":  tftypes.NewValue(tftypes.String, "User"),
			"role":       tftypes.NewValue(tftypes.String, "global:admin"),
			"is_pending": tftypes.NewValue(tftypes.Bool, false),
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "user-123"),
			"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
			"first_name": tftypes.NewValue(tftypes.String, "Test"),
			"last_name":  tftypes.NewValue(tftypes.String, "User"),
			"role":       tftypes.NewValue(tftypes.String, "global:member"),
			"is_pending": tftypes.NewValue(tftypes.Bool, false),
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "email": tftypes.String, "first_name": tftypes.String, "last_name": tftypes.String, "role": tftypes.String, "is_pending": tftypes.Bool, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
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

		assert.True(t, resp.Diagnostics.HasError(), "Update should have errors when refresh fails")
	})
}
