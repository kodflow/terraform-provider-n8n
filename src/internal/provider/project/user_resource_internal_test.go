package project

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/project/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

func TestProjectUserResource_findUserInProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		handler      http.HandlerFunc
		projectID    string
		userID       string
		expectFound  bool
		expectError  bool
		expectedRole string
	}{
		{
			name:         "user found in project",
			projectID:    "proj-123",
			userID:       "user-456",
			expectedRole: "project:admin",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/users" && r.Method == http.MethodGet {
					userID := "user-456"
					role := "project:admin"
					response := map[string]any{
						"data": []any{
							map[string]any{
								"id":    userID,
								"email": "user@example.com",
								"role":  role,
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectFound: true,
			expectError: false,
		},
		{
			name:      "user not found in project",
			projectID: "proj-123",
			userID:    "user-999",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/users" && r.Method == http.MethodGet {
					response := map[string]any{
						"data": []any{
							map[string]any{
								"id":    "user-456",
								"email": "other@example.com",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectFound: false,
			expectError: false,
		},
		{
			name:      "API returns error",
			projectID: "proj-123",
			userID:    "user-456",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/users" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectFound: false,
			expectError: true,
		},
		{
			name:      "empty user list",
			projectID: "proj-123",
			userID:    "user-456",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/users" && r.Method == http.MethodGet {
					response := map[string]any{
						"data": []any{},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectFound: false,
			expectError: false,
		},
		{
			name:      "nil data in response",
			projectID: "proj-123",
			userID:    "user-456",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/users" && r.Method == http.MethodGet {
					response := map[string]any{}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectFound: false,
			expectError: false,
		},
		{
			name:         "user found with nil role",
			projectID:    "proj-123",
			userID:       "user-456",
			expectedRole: "",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/users" && r.Method == http.MethodGet {
					userID := "user-456"
					response := map[string]any{
						"data": []any{
							map[string]any{
								"id":    userID,
								"email": "user@example.com",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectFound: true,
			expectError: false,
		},
		{
			name:      "user with nil ID in list",
			projectID: "proj-123",
			userID:    "user-456",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/users" && r.Method == http.MethodGet {
					response := map[string]any{
						"data": []any{
							map[string]any{
								"email": "user@example.com",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectFound: false,
			expectError: false,
		},
		{
			name:         "multiple users in list",
			projectID:    "proj-123",
			userID:       "user-789",
			expectedRole: "project:editor",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/users" && r.Method == http.MethodGet {
					response := map[string]any{
						"data": []any{
							map[string]any{
								"id":    "user-456",
								"email": "user1@example.com",
								"role":  "project:admin",
							},
							map[string]any{
								"id":    "user-789",
								"email": "user2@example.com",
								"role":  "project:editor",
							},
							map[string]any{
								"id":    "user-101",
								"email": "user3@example.com",
								"role":  "project:viewer",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectFound: true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupUserTestClient(t, tt.handler)
			defer server.Close()

			r := &ProjectUserResource{client: n8nClient}
			state := &models.UserResource{}
			state.ProjectID = types.StringValue(tt.projectID)
			state.UserID = types.StringValue(tt.userID)

			// Create a proper tfsdk.State with schema
			testSchema := schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id":         schema.StringAttribute{Computed: true},
					"project_id": schema.StringAttribute{Required: true},
					"user_id":    schema.StringAttribute{Required: true},
					"role":       schema.StringAttribute{Required: true},
				},
			}
			rawState := map[string]tftypes.Value{
				"id":         tftypes.NewValue(tftypes.String, nil),
				"project_id": tftypes.NewValue(tftypes.String, tt.projectID),
				"user_id":    tftypes.NewValue(tftypes.String, tt.userID),
				"role":       tftypes.NewValue(tftypes.String, nil),
			}
			tfState := tfsdk.State{
				Raw: tftypes.NewValue(tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id":         tftypes.String,
						"project_id": tftypes.String,
						"user_id":    tftypes.String,
						"role":       tftypes.String,
					},
				}, rawState),
				Schema: testSchema,
			}
			resp := &resource.ReadResponse{
				State: tfState,
			}

			found := r.findUserInProject(t.Context(), state, resp)

			if tt.expectFound {
				assert.True(t, found)
				if tt.expectedRole != "" {
					assert.Equal(t, tt.expectedRole, state.Role.ValueString())
				}
			} else {
				assert.False(t, found)
			}

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestProjectUserResource_searchUserInList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userList     *n8nsdk.UserList
		stateUserID  string
		expectFound  bool
		expectedRole string
	}{
		{
			name: "user found with role",
			userList: func() *n8nsdk.UserList {
				return &n8nsdk.UserList{
					Data: []n8nsdk.User{
						{
							Id:   new("user-123"),
							Role: new("project:admin"),
						},
					},
				}
			}(),
			stateUserID:  "user-123",
			expectFound:  true,
			expectedRole: "project:admin",
		},
		{
			name: "user found without role",
			userList: func() *n8nsdk.UserList {
				return &n8nsdk.UserList{
					Data: []n8nsdk.User{
						{
							Id: new("user-456"),
						},
					},
				}
			}(),
			stateUserID: "user-456",
			expectFound: true,
		},
		{
			name: "user not found",
			userList: func() *n8nsdk.UserList {
				return &n8nsdk.UserList{
					Data: []n8nsdk.User{
						{
							Id: new("user-123"),
						},
					},
				}
			}(),
			stateUserID: "user-999",
			expectFound: false,
		},
		{
			name: "empty user list",
			userList: &n8nsdk.UserList{
				Data: []n8nsdk.User{},
			},
			stateUserID: "user-123",
			expectFound: false,
		},
		{
			name:        "nil data in user list",
			userList:    &n8nsdk.UserList{},
			stateUserID: "user-123",
			expectFound: false,
		},
		{
			name: "user with nil ID",
			userList: &n8nsdk.UserList{
				Data: []n8nsdk.User{
					{
						Id: nil,
					},
				},
			},
			stateUserID: "user-123",
			expectFound: false,
		},
		{
			name: "multiple users in list",
			userList: func() *n8nsdk.UserList {
				return &n8nsdk.UserList{
					Data: []n8nsdk.User{
						{Id: new("user-123"), Role: new("project:admin")},
						{Id: new("user-456"), Role: new("project:editor")},
						{Id: new("user-789"), Role: new("project:viewer")},
					},
				}
			}(),
			stateUserID:  "user-456",
			expectFound:  true,
			expectedRole: "project:editor",
		},
		{
			name: "user found is first in list",
			userList: func() *n8nsdk.UserList {
				return &n8nsdk.UserList{
					Data: []n8nsdk.User{
						{Id: new("user-123"), Role: new("project:admin")},
						{Id: new("user-456")},
					},
				}
			}(),
			stateUserID:  "user-123",
			expectFound:  true,
			expectedRole: "project:admin",
		},
		{
			name: "user found is last in list",
			userList: func() *n8nsdk.UserList {
				return &n8nsdk.UserList{
					Data: []n8nsdk.User{
						{Id: new("user-123")},
						{Id: new("user-456"), Role: new("project:viewer")},
					},
				}
			}(),
			stateUserID:  "user-456",
			expectFound:  true,
			expectedRole: "project:viewer",
		},
		{
			name: "empty string user ID",
			userList: func() *n8nsdk.UserList {
				return &n8nsdk.UserList{
					Data: []n8nsdk.User{
						{Id: new("")},
					},
				}
			}(),
			stateUserID: "",
			expectFound: true,
		},
		{
			name: "special characters in user ID",
			userList: func() *n8nsdk.UserList {
				return &n8nsdk.UserList{
					Data: []n8nsdk.User{
						{Id: new("user-!@#$%^&*()"), Role: new("project:admin")},
					},
				}
			}(),
			stateUserID:  "user-!@#$%^&*()",
			expectFound:  true,
			expectedRole: "project:admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &ProjectUserResource{}
			state := &models.UserResource{}
			state.UserID = types.StringValue(tt.stateUserID)

			found := r.searchUserInList(tt.userList, state)

			if tt.expectFound {
				assert.True(t, found)
				if tt.expectedRole != "" {
					assert.Equal(t, tt.expectedRole, state.Role.ValueString())
				}
			} else {
				assert.False(t, found)
			}
		})
	}
}

func setupUserTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
