package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/user/models"
	"github.com/stretchr/testify/assert"
)

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

// TestUserResource_executeCreateLogic tests the executeCreateLogic method with error cases.
func TestUserResource_executeCreateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		email        string
		role         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectID     string
	}{
		{
			name:  "successful creation",
			email: "test@example.com",
			role:  "global:member",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/users" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					// n8n API returns an array response
					json.NewEncoder(w).Encode([]map[string]any{
						{
							"user": map[string]any{
								"id": "user-123",
							},
							"error": "",
						},
					})
					return
				}
				if r.Method == http.MethodGet && r.URL.Path == "/users/user-123" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":        "user-123",
						"email":     "test@example.com",
						"firstName": "Test",
						"lastName":  "User",
						"role":      "global:member",
						"isPending": false,
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-01T00:00:00Z",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectID:    "user-123",
		},
		{
			name:  "API error on creation",
			email: "failed@example.com",
			role:  "global:member",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name:  "missing user ID in response",
			email: "noid@example.com",
			role:  "global:member",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/users" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					// n8n API returns an array response with empty user
					json.NewEncoder(w).Encode([]map[string]any{
						{
							"user":  map[string]any{},
							"error": "",
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name:  "error fetching full details",
			email: "fetchfail@example.com",
			role:  "global:member",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/users" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					// n8n API returns an array response
					json.NewEncoder(w).Encode([]map[string]any{
						{
							"user": map[string]any{
								"id": "user-456",
							},
							"error": "",
						},
					})
					return
				}
				if r.Method == http.MethodGet && r.URL.Path == "/users/user-456" {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "User not found"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &UserResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				Email: types.StringValue(tt.email),
			}
			if tt.role != "" {
				plan.Role = types.StringValue(tt.role)
			}
			resp := &resource.CreateResponse{
				State: resource.CreateResponse{}.State,
			}

			result := r.executeCreateLogic(ctx, plan, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, tt.expectID, plan.ID.ValueString(), "User ID should match")
			}
		})
	}
}

// TestUserResource_executeReadLogic tests the executeReadLogic method with error cases.
func TestUserResource_executeReadLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectEmail  string
	}{
		{
			name:   "successful read",
			userID: "user-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/users/user-123" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":        "user-123",
						"email":     "retrieved@example.com",
						"firstName": "Retrieved",
						"lastName":  "User",
						"role":      "global:admin",
						"isPending": false,
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-01T00:00:00Z",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectEmail: "retrieved@example.com",
		},
		{
			name:   "user not found",
			userID: "user-404",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "User not found"}`))
			},
			expectError: true,
		},
		{
			name:   "API error",
			userID: "user-500",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &UserResource{client: n8nClient}
			ctx := context.Background()
			state := &models.Resource{
				ID: types.StringValue(tt.userID),
			}
			resp := &resource.ReadResponse{
				State: resource.ReadResponse{}.State,
			}

			result := r.executeReadLogic(ctx, state, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, tt.expectEmail, state.Email.ValueString(), "User email should match")
			}
		})
	}
}

// TestUserResource_executeUpdateLogic tests the executeUpdateLogic method with error cases.
func TestUserResource_executeUpdateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       string
		oldEmail     string
		newEmail     string
		oldRole      string
		newRole      string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:     "successful update with role change",
			userID:   "user-123",
			oldEmail: "test@example.com",
			newEmail: "test@example.com",
			oldRole:  "global:member",
			newRole:  "global:admin",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPatch && r.URL.Path == "/users/user-123/role" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				if r.Method == http.MethodGet && r.URL.Path == "/users/user-123" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":        "user-123",
						"email":     "test@example.com",
						"firstName": "Test",
						"lastName":  "User",
						"role":      "global:admin",
						"isPending": false,
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-02T00:00:00Z",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:     "email change not supported",
			userID:   "user-123",
			oldEmail: "old@example.com",
			newEmail: "new@example.com",
			oldRole:  "global:member",
			newRole:  "global:member",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name:     "role update fails",
			userID:   "user-404",
			oldEmail: "test@example.com",
			newEmail: "test@example.com",
			oldRole:  "global:member",
			newRole:  "global:admin",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPatch && r.URL.Path == "/users/user-404/role" {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "User not found"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name:     "refresh after update fails",
			userID:   "user-789",
			oldEmail: "test@example.com",
			newEmail: "test@example.com",
			oldRole:  "global:member",
			newRole:  "global:admin",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPatch && r.URL.Path == "/users/user-789/role" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				if r.Method == http.MethodGet && r.URL.Path == "/users/user-789" {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name:     "no role change",
			userID:   "user-999",
			oldEmail: "test@example.com",
			newEmail: "test@example.com",
			oldRole:  "global:member",
			newRole:  "global:member",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/users/user-999" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":        "user-999",
						"email":     "test@example.com",
						"firstName": "Test",
						"lastName":  "User",
						"role":      "global:member",
						"isPending": false,
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-01T00:00:00Z",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &UserResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				ID:    types.StringValue(tt.userID),
				Email: types.StringValue(tt.newEmail),
				Role:  types.StringValue(tt.newRole),
			}
			state := &models.Resource{
				ID:    types.StringValue(tt.userID),
				Email: types.StringValue(tt.oldEmail),
				Role:  types.StringValue(tt.oldRole),
			}
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}

			result := r.executeUpdateLogic(ctx, plan, state, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			}
		})
	}
}

// TestUserResource_executeDeleteLogic tests the executeDeleteLogic method with error cases.
func TestUserResource_executeDeleteLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:   "successful deletion",
			userID: "user-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete && r.URL.Path == "/users/user-123" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:   "user not found",
			userID: "user-404",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "User not found"}`))
			},
			expectError: true,
		},
		{
			name:   "API error",
			userID: "user-500",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &UserResource{client: n8nClient}
			ctx := context.Background()
			state := &models.Resource{
				ID: types.StringValue(tt.userID),
			}
			resp := &resource.DeleteResponse{
				State: resource.DeleteResponse{}.State,
			}

			result := r.executeDeleteLogic(ctx, state, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			}
		})
	}
}

// TestUserResource_schemaAttributes tests the schemaAttributes method.
func TestUserResource_schemaAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		expectedFields []string
	}{
		{
			name: "returns all expected schema attributes",
			expectedFields: []string{
				"id",
				"email",
				"first_name",
				"last_name",
				"role",
				"is_pending",
				"created_at",
				"updated_at",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &UserResource{}
			attrs := r.schemaAttributes()

			assert.NotNil(t, attrs, "Schema attributes should not be nil")
			assert.Len(t, attrs, len(tt.expectedFields), "Should have correct number of fields")

			for _, field := range tt.expectedFields {
				assert.Contains(t, attrs, field, "Should contain field: %s", field)
			}

			// Verify specific attribute properties
			assert.True(t, attrs["id"].IsComputed(), "id should be computed")
			assert.True(t, attrs["email"].IsRequired(), "email should be required")
			assert.True(t, attrs["role"].IsOptional(), "role should be optional")
			assert.True(t, attrs["role"].IsComputed(), "role should be computed")
		})
	}
}

// TestUserResource_createUser tests the createUser method with error cases.
func TestUserResource_createUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		email        string
		role         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectUserID string
		expectError  bool
	}{
		{
			name:  "successful user creation with role",
			email: "test@example.com",
			role:  "global:admin",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				// n8n API returns an array response
				json.NewEncoder(w).Encode([]map[string]any{
					{
						"user": map[string]any{
							"id": "user-123",
						},
						"error": "",
					},
				})
			},
			expectUserID: "user-123",
			expectError:  false,
		},
		{
			name:  "successful user creation without role",
			email: "norole@example.com",
			role:  "",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				// n8n API returns an array response
				json.NewEncoder(w).Encode([]map[string]any{
					{
						"user": map[string]any{
							"id": "user-456",
						},
						"error": "",
					},
				})
			},
			expectUserID: "user-456",
			expectError:  false,
		},
		{
			name:  "API error on creation",
			email: "error@example.com",
			role:  "global:member",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message": "Invalid request"}`))
			},
			expectUserID: "",
			expectError:  true,
		},
		{
			name:  "missing user ID in response",
			email: "noid@example.com",
			role:  "global:member",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				// n8n API returns an array response with empty user
				json.NewEncoder(w).Encode([]map[string]any{
					{
						"user":  map[string]any{},
						"error": "",
					},
				})
			},
			expectUserID: "",
			expectError:  true,
		},
		{
			name:  "empty array response",
			email: "niluser@example.com",
			role:  "global:member",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				// n8n API returns empty array
				json.NewEncoder(w).Encode([]map[string]any{})
			},
			expectUserID: "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &UserResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				Email: types.StringValue(tt.email),
			}
			if tt.role != "" {
				plan.Role = types.StringValue(tt.role)
			}
			resp := &resource.CreateResponse{
				State: resource.CreateResponse{}.State,
			}

			userID := r.createUser(ctx, plan, resp)

			if tt.expectError {
				assert.Empty(t, userID, "Should return empty string on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.Equal(t, tt.expectUserID, userID, "User ID should match")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			}
		})
	}
}

// TestUserResource_fetchFullUserDetails tests the fetchFullUserDetails method with error cases.
func TestUserResource_fetchFullUserDetails(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectNil    bool
		expectEmail  string
	}{
		{
			name:   "successful fetch",
			userID: "user-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"id":        "user-123",
					"email":     "fetched@example.com",
					"firstName": "Fetched",
					"lastName":  "User",
					"role":      "global:member",
					"isPending": false,
					"createdAt": "2024-01-01T00:00:00Z",
					"updatedAt": "2024-01-01T00:00:00Z",
				})
			},
			expectNil:   false,
			expectEmail: "fetched@example.com",
		},
		{
			name:   "user not found",
			userID: "user-404",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "User not found"}`))
			},
			expectNil: true,
		},
		{
			name:   "API error",
			userID: "user-500",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectNil: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &UserResource{client: n8nClient}
			ctx := context.Background()
			resp := &resource.CreateResponse{
				State: resource.CreateResponse{}.State,
			}

			user := r.fetchFullUserDetails(ctx, tt.userID, resp)

			if tt.expectNil {
				assert.Nil(t, user, "Should return nil on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.NotNil(t, user, "Should return user object")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				if user != nil {
					assert.Equal(t, tt.expectEmail, user.Email, "Email should match")
				}
			}
		})
	}
}

// TestUserResource_validateEmailUnchanged tests the validateEmailUnchanged method with error cases.
func TestUserResource_validateEmailUnchanged(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		planEmail   string
		stateEmail  string
		expectValid bool
		expectError bool
	}{
		{
			name:        "email unchanged",
			planEmail:   "test@example.com",
			stateEmail:  "test@example.com",
			expectValid: true,
			expectError: false,
		},
		{
			name:        "email changed",
			planEmail:   "new@example.com",
			stateEmail:  "old@example.com",
			expectValid: false,
			expectError: true,
		},
		{
			name:        "email unchanged case sensitive",
			planEmail:   "Test@Example.com",
			stateEmail:  "Test@Example.com",
			expectValid: true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &UserResource{}
			plan := &models.Resource{
				Email: types.StringValue(tt.planEmail),
			}
			state := &models.Resource{
				Email: types.StringValue(tt.stateEmail),
			}
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}

			result := r.validateEmailUnchanged(plan, state, resp)

			assert.Equal(t, tt.expectValid, result, "Validation result should match expected")
			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error for email change")
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			}
		})
	}
}

// TestUserResource_updateRoleIfChanged tests the updateRoleIfChanged method with error cases.
func TestUserResource_updateRoleIfChanged(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		planRole     string
		stateRole    string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:      "role changed successfully",
			planRole:  "global:admin",
			stateRole: "global:member",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPatch && r.URL.Path == "/users/user-123/role" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:      "no role change",
			planRole:  "global:member",
			stateRole: "global:member",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:      "API error on role update",
			planRole:  "global:admin",
			stateRole: "global:member",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message": "Invalid role"}`))
			},
			expectError: true,
		},
		{
			name:      "user not found during role update",
			planRole:  "global:admin",
			stateRole: "global:member",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "User not found"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &UserResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				ID:   types.StringValue("user-123"),
				Role: types.StringValue(tt.planRole),
			}
			state := &models.Resource{
				ID:   types.StringValue("user-123"),
				Role: types.StringValue(tt.stateRole),
			}
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}

			r.updateRoleIfChanged(ctx, plan, state, resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			}
		})
	}
}

// TestUserResource_refreshUserData tests the refreshUserData method with error cases.
func TestUserResource_refreshUserData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectNil    bool
		expectRole   string
	}{
		{
			name:   "successful refresh",
			userID: "user-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"id":        "user-123",
					"email":     "refreshed@example.com",
					"firstName": "Refreshed",
					"lastName":  "User",
					"role":      "global:admin",
					"isPending": false,
					"createdAt": "2024-01-01T00:00:00Z",
					"updatedAt": "2024-01-02T00:00:00Z",
				})
			},
			expectNil:  false,
			expectRole: "global:admin",
		},
		{
			name:   "user not found",
			userID: "user-404",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "User not found"}`))
			},
			expectNil: true,
		},
		{
			name:   "API error on refresh",
			userID: "user-500",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectNil: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &UserResource{client: n8nClient}
			ctx := context.Background()
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}

			user := r.refreshUserData(ctx, tt.userID, resp)

			if tt.expectNil {
				assert.Nil(t, user, "Should return nil on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.NotNil(t, user, "Should return user object")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				if user != nil && user.Role != nil {
					assert.Equal(t, tt.expectRole, *user.Role, "Role should match")
				}
			}
		})
	}
}

// TestUserResource_createUser_ErrorPaths tests createUser error handling paths.
func TestUserResource_createUser_ErrorPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		email        string
		role         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectUserID string
		expectError  bool
	}{
		{
			name:  "invalid JSON response",
			email: "test@example.com",
			role:  "global:member",
			setupHandler: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				// Return invalid JSON
				_, _ = w.Write([]byte(`{invalid json`))
			},
			expectUserID: "",
			expectError:  true,
		},
		{
			name:  "API error in response body",
			email: "existing@example.com",
			role:  "global:member",
			setupHandler: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				// n8n API returns an array response with error field populated
				_ = json.NewEncoder(w).Encode([]map[string]any{
					{
						"user": map[string]any{
							"id": "",
						},
						"error": "Email already exists",
					},
				})
			},
			expectUserID: "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &UserResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				Email: types.StringValue(tt.email),
				Role:  types.StringValue(tt.role),
			}
			resp := &resource.CreateResponse{
				State: resource.CreateResponse{}.State,
			}

			userID := r.createUser(ctx, plan, resp)

			if tt.expectError {
				assert.Empty(t, userID, "Should return empty string on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.Equal(t, tt.expectUserID, userID, "User ID should match")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			}
		})
	}
}

// TestUserResource_createUser_NilHTTPClient tests createUser with nil HTTP client uses default.
func TestUserResource_createUser_NilHTTPClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		email        string
		role         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectUserID string
		expectError  bool
	}{
		{
			name:  "nil HTTP client uses default",
			email: "test@example.com",
			role:  "global:member",
			setupHandler: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode([]map[string]any{
					{
						"user": map[string]any{
							"id": "user-from-default-client",
						},
						"error": "",
					},
				})
			},
			expectUserID: "user-from-default-client",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(tt.setupHandler))
			defer server.Close()

			// Create config with nil HTTPClient
			cfg := n8nsdk.NewConfiguration()
			cfg.Servers = n8nsdk.ServerConfigurations{
				{
					URL:         server.URL,
					Description: "Test server",
				},
			}
			cfg.HTTPClient = nil // Explicitly set to nil
			cfg.AddDefaultHeader("X-N8N-API-KEY", "test-key")

			apiClient := n8nsdk.NewAPIClient(cfg)
			n8nClient := &client.N8nClient{
				APIClient: apiClient,
			}

			r := &UserResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				Email: types.StringValue(tt.email),
				Role:  types.StringValue(tt.role),
			}
			resp := &resource.CreateResponse{
				State: resource.CreateResponse{}.State,
			}

			userID := r.createUser(ctx, plan, resp)

			if tt.expectError {
				assert.Empty(t, userID, "Should return empty string on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.Equal(t, tt.expectUserID, userID, "Should return user ID using default client")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			}
		})
	}
}

// TestUserResource_buildUserCreateRequest tests the buildUserCreateRequest helper function.
func TestUserResource_buildUserCreateRequest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		email       string
		role        string
		expectBody  bool
		expectError bool
	}{
		{
			name:        "success - with role",
			email:       "test@example.com",
			role:        "global:member",
			expectBody:  true,
			expectError: false,
		},
		{
			name:        "success - without role",
			email:       "test@example.com",
			role:        "",
			expectBody:  true,
			expectError: false,
		},
		{
			name:        "success - with admin role",
			email:       "admin@example.com",
			role:        "global:admin",
			expectBody:  true,
			expectError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &UserResource{}
			plan := &models.Resource{
				Email: types.StringValue(tt.email),
			}
			// Check if role is provided.
			if tt.role != "" {
				plan.Role = types.StringValue(tt.role)
			} else {
				plan.Role = types.StringNull()
			}
			resp := &resource.CreateResponse{}
			body := r.buildUserCreateRequest(plan, resp)
			// Check body expectations.
			if tt.expectBody {
				assert.NotEmpty(t, body, "Should return non-empty body")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			}
		})
	}
}

// TestUserResource_executeUserHTTPRequest tests the executeUserHTTPRequest helper function.
func TestUserResource_executeUserHTTPRequest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectBody   bool
		expectError  bool
		expectStatus int
	}{
		{
			name: "success - valid response",
			setupHandler: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`[{"user":{"id":"user-123"},"error":""}]`))
			},
			expectBody:   true,
			expectError:  false,
			expectStatus: http.StatusCreated,
		},
		{
			name: "error - server error",
			setupHandler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal server error"}`))
			},
			expectBody:   true,
			expectError:  false,
			expectStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(tt.setupHandler))
			defer server.Close()
			cfg := n8nsdk.NewConfiguration()
			cfg.Servers = n8nsdk.ServerConfigurations{{URL: server.URL, Description: "Test server"}}
			cfg.HTTPClient = server.Client()
			cfg.AddDefaultHeader("X-N8N-API-KEY", "test-key")
			apiClient := n8nsdk.NewAPIClient(cfg)
			r := &UserResource{
				client: &client.N8nClient{APIClient: apiClient},
			}
			ctx := context.Background()
			jsonBody := []byte(`[{"email":"test@example.com"}]`)
			resp := &resource.CreateResponse{}
			body, statusCode := r.executeUserHTTPRequest(ctx, jsonBody, resp)
			// Check body expectations.
			if tt.expectBody {
				assert.NotEmpty(t, body, "Should return non-empty body")
			}
			assert.Equal(t, tt.expectStatus, statusCode, "Should return expected status code")
		})
	}
}

// TestUserResource_parseUserCreateResponse tests the parseUserCreateResponse helper function.
func TestUserResource_parseUserCreateResponse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		body         []byte
		statusCode   int
		expectUserID string
		expectError  bool
	}{
		{
			name:         "success - valid response",
			body:         []byte(`[{"user":{"id":"user-123","email":"test@example.com","role":"global:member"},"error":""}]`),
			statusCode:   http.StatusCreated,
			expectUserID: "user-123",
			expectError:  false,
		},
		{
			name:         "error - status code out of range",
			body:         []byte(`{"error":"not found"}`),
			statusCode:   http.StatusNotFound,
			expectUserID: "",
			expectError:  true,
		},
		{
			name:         "error - invalid JSON",
			body:         []byte(`{invalid json`),
			statusCode:   http.StatusCreated,
			expectUserID: "",
			expectError:  true,
		},
		{
			name:         "error - empty response array",
			body:         []byte(`[]`),
			statusCode:   http.StatusCreated,
			expectUserID: "",
			expectError:  true,
		},
		{
			name:         "error - API error in response",
			body:         []byte(`[{"user":{"id":""},"error":"Email already exists"}]`),
			statusCode:   http.StatusCreated,
			expectUserID: "",
			expectError:  true,
		},
		{
			name:         "error - empty user ID",
			body:         []byte(`[{"user":{"id":"","email":"test@example.com"},"error":""}]`),
			statusCode:   http.StatusCreated,
			expectUserID: "",
			expectError:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &UserResource{}
			resp := &resource.CreateResponse{}
			userID := r.parseUserCreateResponse(tt.body, tt.statusCode, resp)
			// Check expectations.
			if tt.expectError {
				assert.Empty(t, userID, "Should return empty string on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.Equal(t, tt.expectUserID, userID, "Should return expected user ID")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			}
		})
	}
}
