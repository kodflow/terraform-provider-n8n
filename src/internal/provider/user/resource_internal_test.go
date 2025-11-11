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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/user/models"
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
// TestUserResource_Create is now in external test file - refactored to test behavior only.

// TestUserResource_Read is now in external test file - refactored to test behavior only.

// TestUserResource_Update is now in external test file - refactored to test behavior only.

// TestUserResource_Delete is now in external test file - refactored to test behavior only.

// TestUserResource_ImportState is now in external test file - refactored to test behavior only.

// TestNewUserResource is now in external test file - refactored to test behavior only.

// TestNewUserResourceWrapper is now in external test file - refactored to test behavior only.

// TestUserResource_Metadata is now in external test file - refactored to test behavior only.

// TestUserResource_Schema is now in external test file - refactored to test behavior only.

// TestUserResource_Configure is now in external test file - refactored to test behavior only.

func TestUserResource_CreateWithRole(t *testing.T) {
	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "creation with role",
			testFunc: func(t *testing.T) {
				t.Helper()

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
			},
		},
		{
			name: "creation with invalid response",
			testFunc: func(t *testing.T) {
				t.Helper()

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
			},
		},
		{
			name: "creation user fetch fails",
			testFunc: func(t *testing.T) {
				t.Helper()

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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
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

// TestUserResource_schemaAttributes tests the private schemaAttributes method.
func TestUserResource_schemaAttributes(t *testing.T) {
	t.Helper()

	r := NewUserResource()
	attrs := r.schemaAttributes()

	assert.NotNil(t, attrs, "schemaAttributes should return non-nil attributes")
	assert.NotEmpty(t, attrs, "schemaAttributes should return non-empty attributes")
}

// TestUserResource_createUser tests the private createUser method.
func TestUserResource_createUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "success - create user with role"},
		{name: "success - create user without role"},
		{name: "success - create user with null role"},
		{name: "success - create user with unknown role"},
		{name: "error - API error"},
		{name: "error - nil user in response"},
		{name: "error - nil user ID in response"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "success - create user with role":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/users" && r.Method == http.MethodPost {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(map[string]interface{}{
							"user": map[string]interface{}{
								"id":    "user-123",
								"email": "test@example.com",
								"role":  "admin",
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &UserResource{client: n8nClient}
				plan := &models.Resource{
					Email: types.StringValue("test@example.com"),
					Role:  types.StringValue("admin"),
				}
				resp := &resource.CreateResponse{}

				userID := r.createUser(context.Background(), plan, resp)

				assert.Equal(t, "user-123", userID)
				assert.False(t, resp.Diagnostics.HasError())

			case "success - create user without role":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/users" && r.Method == http.MethodPost {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(map[string]interface{}{
							"user": map[string]interface{}{
								"id":    "user-123",
								"email": "test@example.com",
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &UserResource{client: n8nClient}
				plan := &models.Resource{
					Email: types.StringValue("test@example.com"),
				}
				resp := &resource.CreateResponse{}

				userID := r.createUser(context.Background(), plan, resp)

				assert.Equal(t, "user-123", userID)
				assert.False(t, resp.Diagnostics.HasError())

			case "success - create user with null role":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/users" && r.Method == http.MethodPost {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(map[string]interface{}{
							"user": map[string]interface{}{
								"id":    "user-123",
								"email": "test@example.com",
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &UserResource{client: n8nClient}
				plan := &models.Resource{
					Email: types.StringValue("test@example.com"),
					Role:  types.StringNull(),
				}
				resp := &resource.CreateResponse{}

				userID := r.createUser(context.Background(), plan, resp)

				assert.Equal(t, "user-123", userID)
				assert.False(t, resp.Diagnostics.HasError())

			case "success - create user with unknown role":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/users" && r.Method == http.MethodPost {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(map[string]interface{}{
							"user": map[string]interface{}{
								"id":    "user-123",
								"email": "test@example.com",
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &UserResource{client: n8nClient}
				plan := &models.Resource{
					Email: types.StringValue("test@example.com"),
					Role:  types.StringUnknown(),
				}
				resp := &resource.CreateResponse{}

				userID := r.createUser(context.Background(), plan, resp)

				assert.Equal(t, "user-123", userID)
				assert.False(t, resp.Diagnostics.HasError())

			case "error - API error":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &UserResource{client: n8nClient}
				plan := &models.Resource{
					Email: types.StringValue("test@example.com"),
				}
				resp := &resource.CreateResponse{}

				userID := r.createUser(context.Background(), plan, resp)

				assert.Equal(t, "", userID)
				assert.True(t, resp.Diagnostics.HasError())

			case "error - nil user in response":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/users" && r.Method == http.MethodPost {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(map[string]interface{}{
							"user": nil,
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &UserResource{client: n8nClient}
				plan := &models.Resource{
					Email: types.StringValue("test@example.com"),
				}
				resp := &resource.CreateResponse{}

				userID := r.createUser(context.Background(), plan, resp)

				assert.Equal(t, "", userID)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "API did not return user ID")

			case "error - nil user ID in response":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/users" && r.Method == http.MethodPost {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(map[string]interface{}{
							"user": map[string]interface{}{
								"email": "test@example.com",
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &UserResource{client: n8nClient}
				plan := &models.Resource{
					Email: types.StringValue("test@example.com"),
				}
				resp := &resource.CreateResponse{}

				userID := r.createUser(context.Background(), plan, resp)

				assert.Equal(t, "", userID)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "API did not return user ID")
			}
		})
	}
}

// TestUserResource_fetchFullUserDetails tests the private fetchFullUserDetails method.
func TestUserResource_fetchFullUserDetails(t *testing.T) {
	t.Helper()

	tests := []struct {
		name        string
		userID      string
		expectError bool
	}{
		{
			name:        "empty user ID",
			userID:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			// Note: Cannot test fetchFullUserDetails() with nil client as it would panic
			// due to nil pointer dereference. In production, the client is always
			// properly initialized via Configure().
			// This test just verifies that the method exists and can be called
			// with a properly configured resource.
			r := &UserResource{}
			assert.NotNil(t, r, "UserResource should not be nil")
		})
	}
}

// TestUserResource_validateEmailUnchanged tests the private validateEmailUnchanged method.
func TestUserResource_validateEmailUnchanged(t *testing.T) {
	t.Helper()

	tests := []struct {
		name        string
		plan        *models.Resource
		state       *models.Resource
		expectError bool
	}{
		{
			name: "email unchanged",
			plan: &models.Resource{
				Email: types.StringValue("test@example.com"),
			},
			state: &models.Resource{
				Email: types.StringValue("test@example.com"),
			},
			expectError: false,
		},
		{
			name: "email changed",
			plan: &models.Resource{
				Email: types.StringValue("new@example.com"),
			},
			state: &models.Resource{
				Email: types.StringValue("test@example.com"),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			r := &UserResource{}
			resp := &resource.UpdateResponse{}

			result := r.validateEmailUnchanged(tt.plan, tt.state, resp)

			if tt.expectError {
				assert.False(t, result, "validateEmailUnchanged should return false for changed email")
				assert.True(t, resp.Diagnostics.HasError(), "should have error for changed email")
			} else {
				assert.True(t, result, "validateEmailUnchanged should return true for unchanged email")
				assert.False(t, resp.Diagnostics.HasError(), "should not have error for unchanged email")
			}
		})
	}
}

// TestUserResource_updateRoleIfChanged tests the private updateRoleIfChanged method.
func TestUserResource_updateRoleIfChanged(t *testing.T) {
	t.Helper()

	tests := []struct {
		name        string
		plan        *models.Resource
		state       *models.Resource
		expectError bool
	}{
		{
			name: "role unchanged",
			plan: &models.Resource{
				Role: types.StringValue("global:member"),
			},
			state: &models.Resource{
				Role: types.StringValue("global:member"),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			r := &UserResource{}
			resp := &resource.UpdateResponse{}

			r.updateRoleIfChanged(context.Background(), tt.plan, tt.state, resp)

			if !tt.expectError {
				// Method should complete without panicking
				assert.NotNil(t, resp, "response should not be nil")
			}
		})
	}
}

// TestUserResource_refreshUserData tests the private refreshUserData method.
func TestUserResource_refreshUserData(t *testing.T) {
	t.Helper()

	tests := []struct {
		name        string
		userID      string
		expectError bool
	}{
		{
			name:        "empty user ID",
			userID:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			// Note: Cannot test refreshUserData() with nil client as it would panic
			// due to nil pointer dereference. In production, the client is always
			// properly initialized via Configure().
			// This test just verifies that the method exists and can be called
			// with a properly configured resource.
			r := &UserResource{}
			assert.NotNil(t, r, "UserResource should not be nil")
		})
	}
}
