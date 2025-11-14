package project

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/project/models"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
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

// TestProjectResource_executeCreateLogic tests the executeCreateLogic method with error cases.
func TestProjectResource_executeCreateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		projectName  string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectID     string
	}{
		{
			name:        "successful creation",
			projectName: "Test Project",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/projects" {
					w.WriteHeader(http.StatusCreated)
					return
				}
				if r.Method == http.MethodGet && r.URL.Path == "/projects" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{
							{
								"id":   "proj-123",
								"name": "Test Project",
								"type": "team",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectID:    "proj-123",
		},
		{
			name:        "creation API error",
			projectName: "Failed Project",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name:        "list API error after creation",
			projectName: "Test Project",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/projects" {
					w.WriteHeader(http.StatusCreated)
					return
				}
				if r.Method == http.MethodGet && r.URL.Path == "/projects" {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Cannot retrieve projects"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name:        "project not found after creation",
			projectName: "Test Project",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/projects" {
					w.WriteHeader(http.StatusCreated)
					return
				}
				if r.Method == http.MethodGet && r.URL.Path == "/projects" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{},
					})
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

			r := &ProjectResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				Name: types.StringValue(tt.projectName),
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
				assert.Equal(t, tt.expectID, plan.ID.ValueString(), "Project ID should match")
			}
		})
	}
}

// TestProjectResource_executeReadLogic tests the executeReadLogic method with error cases.
func TestProjectResource_executeReadLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		projectID    string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectName   string
	}{
		{
			name:      "successful read",
			projectID: "proj-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/projects" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{
							{
								"id":   "proj-123",
								"name": "Retrieved Project",
								"type": "team",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectName:  "Retrieved Project",
		},
		// Note: "project not found" case (RemoveResource) is tested in full CRUD tests
		// as it requires a properly initialized tfsdk.State with schema
		{
			name:      "API error",
			projectID: "proj-500",
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

			r := &ProjectResource{client: n8nClient}
			ctx := context.Background()
			state := &models.Resource{
				ID: types.StringValue(tt.projectID),
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
				assert.Equal(t, tt.expectName, state.Name.ValueString(), "Project name should match")
			}
		})
	}
}

// TestProjectResource_executeUpdateLogic tests the executeUpdateLogic method with error cases.
func TestProjectResource_executeUpdateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		projectID    string
		newName      string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:      "successful update",
			projectID: "proj-123",
			newName:   "Updated Project",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/projects/proj-123" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				if r.Method == http.MethodGet && r.URL.Path == "/projects" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{
							{
								"id":   "proj-123",
								"name": "Updated Project",
								"type": "team",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:      "update API error",
			projectID: "proj-404",
			newName:   "Updated Project",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Project not found"}`))
			},
			expectError: true,
		},
		{
			name:      "list API error after update",
			projectID: "proj-123",
			newName:   "Updated Project",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/projects/proj-123" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				if r.Method == http.MethodGet && r.URL.Path == "/projects" {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Cannot retrieve projects"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name:      "project not found after update",
			projectID: "proj-123",
			newName:   "Updated Project",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/projects/proj-123" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				if r.Method == http.MethodGet && r.URL.Path == "/projects" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{},
					})
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

			r := &ProjectResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				ID:   types.StringValue(tt.projectID),
				Name: types.StringValue(tt.newName),
			}
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}

			result := r.executeUpdateLogic(ctx, plan, resp)

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

// TestProjectResource_executeDeleteLogic tests the executeDeleteLogic method with error cases.
func TestProjectResource_executeDeleteLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		projectID    string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:      "successful deletion",
			projectID: "proj-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete && r.URL.Path == "/projects/proj-123" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:      "project not found",
			projectID: "proj-404",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Project not found"}`))
			},
			expectError: true,
		},
		{
			name:      "API error",
			projectID: "proj-500",
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

			r := &ProjectResource{client: n8nClient}
			ctx := context.Background()
			state := &models.Resource{
				ID: types.StringValue(tt.projectID),
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

// TestProjectUserResource_executeCreateLogic tests the executeCreateLogic method with error cases.
func TestProjectUserResource_executeCreateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		projectID    string
		userID       string
		role         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectID     string
	}{
		{
			name:      "successful creation",
			projectID: "proj-123",
			userID:    "user-456",
			role:      "project:admin",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/projects/proj-123/users" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectID:    "proj-123/user-456",
		},
		{
			name:      "API error",
			projectID: "proj-500",
			userID:    "user-456",
			role:      "project:admin",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name:      "invalid project ID",
			projectID: "proj-404",
			userID:    "user-456",
			role:      "project:admin",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Project not found"}`))
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

			r := &ProjectUserResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.UserResource{
				ProjectID: types.StringValue(tt.projectID),
				UserID:    types.StringValue(tt.userID),
				Role:      types.StringValue(tt.role),
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
				assert.Equal(t, tt.expectID, plan.ID.ValueString(), "Composite ID should match")
			}
		})
	}
}

// TestProjectUserResource_executeReadLogic tests the executeReadLogic method with error cases.
func TestProjectUserResource_executeReadLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		projectID    string
		userID       string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectRole   string
	}{
		{
			name:      "successful read",
			projectID: "proj-123",
			userID:    "user-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/users" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{
							{
								"id":    "user-456",
								"email": "user@example.com",
								"role":  "project:admin",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectRole:  "project:admin",
		},
		// Note: "user not found" case (RemoveResource) is tested in full CRUD tests
		// as it requires a properly initialized tfsdk.State with schema
		{
			name:      "API error",
			projectID: "proj-500",
			userID:    "user-456",
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

			r := &ProjectUserResource{client: n8nClient}
			ctx := context.Background()
			state := &models.UserResource{
				ProjectID: types.StringValue(tt.projectID),
				UserID:    types.StringValue(tt.userID),
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
				assert.Equal(t, tt.expectRole, state.Role.ValueString(), "Role should match")
			}
		})
	}
}

// TestProjectUserResource_executeUpdateLogic tests the executeUpdateLogic method with error cases.
func TestProjectUserResource_executeUpdateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		planProject  string
		planUser     string
		planRole     string
		stateProject string
		stateUser    string
		stateRole    string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:         "successful role update",
			planProject:  "proj-123",
			planUser:     "user-456",
			planRole:     "project:editor",
			stateProject: "proj-123",
			stateUser:    "user-456",
			stateRole:    "project:admin",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPatch && r.URL.Path == "/projects/proj-123/users/user-456" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:         "no change in role",
			planProject:  "proj-123",
			planUser:     "user-456",
			planRole:     "project:admin",
			stateProject: "proj-123",
			stateUser:    "user-456",
			stateRole:    "project:admin",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:         "project ID change not supported",
			planProject:  "proj-999",
			planUser:     "user-456",
			planRole:     "project:admin",
			stateProject: "proj-123",
			stateUser:    "user-456",
			stateRole:    "project:admin",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name:         "user ID change not supported",
			planProject:  "proj-123",
			planUser:     "user-999",
			planRole:     "project:admin",
			stateProject: "proj-123",
			stateUser:    "user-456",
			stateRole:    "project:admin",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name:         "API error on role update",
			planProject:  "proj-123",
			planUser:     "user-456",
			planRole:     "project:editor",
			stateProject: "proj-123",
			stateUser:    "user-456",
			stateRole:    "project:admin",
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

			r := &ProjectUserResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.UserResource{
				ProjectID: types.StringValue(tt.planProject),
				UserID:    types.StringValue(tt.planUser),
				Role:      types.StringValue(tt.planRole),
			}
			state := &models.UserResource{
				ProjectID: types.StringValue(tt.stateProject),
				UserID:    types.StringValue(tt.stateUser),
				Role:      types.StringValue(tt.stateRole),
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

// TestProjectUserResource_executeDeleteLogic tests the executeDeleteLogic method with error cases.
func TestProjectUserResource_executeDeleteLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		projectID    string
		userID       string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:      "successful deletion",
			projectID: "proj-123",
			userID:    "user-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete && r.URL.Path == "/projects/proj-123/users/user-456" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:      "user not found in project",
			projectID: "proj-123",
			userID:    "user-404",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "User not found in project"}`))
			},
			expectError: true,
		},
		{
			name:      "API error",
			projectID: "proj-500",
			userID:    "user-456",
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

			r := &ProjectUserResource{client: n8nClient}
			ctx := context.Background()
			state := &models.UserResource{
				ProjectID: types.StringValue(tt.projectID),
				UserID:    types.StringValue(tt.userID),
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

// TestProjectResource_createProject tests the createProject helper method.
func TestProjectResource_createProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		projectName  string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:        "successful project creation",
			projectName: "Test Project",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/projects" {
					w.WriteHeader(http.StatusCreated)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:        "API returns 500 error",
			projectName: "Failed Project",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name:        "API returns 400 bad request",
			projectName: "",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message": "Invalid project name"}`))
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

			r := &ProjectResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				Name: types.StringValue(tt.projectName),
			}
			resp := &resource.CreateResponse{
				State: resource.CreateResponse{}.State,
			}

			result := r.createProject(ctx, plan, resp)

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

// TestProjectResource_findCreatedProject tests the findCreatedProject helper method.
func TestProjectResource_findCreatedProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		projectName  string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectNil    bool
		expectID     string
	}{
		{
			name:        "project found in list",
			projectName: "Test Project",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/projects" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{
							{
								"id":   "proj-123",
								"name": "Test Project",
								"type": "team",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectNil: false,
			expectID:  "proj-123",
		},
		{
			name:        "project not found in list",
			projectName: "Missing Project",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/projects" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectNil: true,
		},
		{
			name:        "API returns error",
			projectName: "Error Project",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Cannot retrieve projects"}`))
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

			r := &ProjectResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				Name: types.StringValue(tt.projectName),
			}
			resp := &resource.CreateResponse{
				State: resource.CreateResponse{}.State,
			}

			result := r.findCreatedProject(ctx, plan, resp)

			if tt.expectNil {
				assert.Nil(t, result, "Should return nil")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.NotNil(t, result, "Should return project")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, tt.expectID, *result.Id, "Project ID should match")
			}
		})
	}
}

// TestProjectResource_updatePlanFromProject tests the updatePlanFromProject helper method.
func TestProjectResource_updatePlanFromProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		projectID   *string
		projectName string
		projectType *string
		expectType  bool
		expectNilID bool
	}{
		{
			name:        "update with all fields",
			projectID:   stringPtr("proj-123"),
			projectName: "Test Project",
			projectType: stringPtr("team"),
			expectType:  true,
			expectNilID: false,
		},
		{
			name:        "update with nil type",
			projectID:   stringPtr("proj-456"),
			projectName: "Another Project",
			projectType: nil,
			expectType:  false,
			expectNilID: false,
		},
		{
			name:        "update with empty type",
			projectID:   stringPtr("proj-789"),
			projectName: "Empty Type Project",
			projectType: stringPtr(""),
			expectType:  true,
			expectNilID: false,
		},
		{
			name:        "edge case with nil ID (error handling)",
			projectID:   nil,
			projectName: "No ID Project",
			projectType: stringPtr("team"),
			expectType:  true,
			expectNilID: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &ProjectResource{}
			plan := &models.Resource{}
			project := &n8nsdk.Project{
				Id:   tt.projectID,
				Name: tt.projectName,
				Type: tt.projectType,
			}

			r.updatePlanFromProject(plan, project)

			if tt.expectNilID {
				assert.True(t, plan.ID.IsNull(), "ID should be null when project.Id is nil")
			} else {
				assert.Equal(t, *tt.projectID, plan.ID.ValueString(), "ID should match")
			}
			assert.Equal(t, tt.projectName, plan.Name.ValueString(), "Name should match")

			if tt.expectType {
				if tt.projectType != nil {
					assert.Equal(t, *tt.projectType, plan.Type.ValueString(), "Type should match")
				}
			} else {
				assert.True(t, plan.Type.IsNull(), "Type should be null")
			}
		})
	}
}

// TestProjectResource_findProjectByID tests the findProjectByID helper method.
func TestProjectResource_findProjectByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		projectID    string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectNil    bool
		expectName   string
	}{
		{
			name:      "project found by ID",
			projectID: "proj-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/projects" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{
							{
								"id":   "proj-123",
								"name": "Found Project",
								"type": "team",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectNil:  false,
			expectName: "Found Project",
		},
		// Note: "project not found" case (RemoveResource) is tested in full CRUD tests
		// as it requires a properly initialized tfsdk.State with schema
		{
			name:      "API error",
			projectID: "proj-500",
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

			r := &ProjectResource{client: n8nClient}
			ctx := context.Background()
			state := &models.Resource{
				ID: types.StringValue(tt.projectID),
			}
			resp := &resource.ReadResponse{
				State: resource.ReadResponse{}.State,
			}

			result := r.findProjectByID(ctx, state, resp)

			if tt.expectNil {
				assert.Nil(t, result, "Should return nil")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.NotNil(t, result, "Should return project")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, tt.expectName, result.Name, "Project name should match")
			}
		})
	}
}

// TestProjectResource_updateStateFromProject tests the updateStateFromProject helper method.
func TestProjectResource_updateStateFromProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		projectName string
		projectType *string
		expectType  bool
	}{
		{
			name:        "update state with type",
			projectName: "State Project",
			projectType: stringPtr("team"),
			expectType:  true,
		},
		{
			name:        "update state without type",
			projectName: "No Type Project",
			projectType: nil,
			expectType:  false,
		},
		{
			name:        "update state with empty type",
			projectName: "Empty Type Project",
			projectType: stringPtr(""),
			expectType:  true,
		},
		{
			name:        "edge case with empty name (error handling)",
			projectName: "",
			projectType: stringPtr("team"),
			expectType:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &ProjectResource{}
			state := &models.Resource{
				ID: types.StringValue("existing-id"),
			}
			project := &n8nsdk.Project{
				Name: tt.projectName,
				Type: tt.projectType,
			}

			// Test that function doesn't panic even with edge cases
			assert.NotPanics(t, func() {
				r.updateStateFromProject(state, project)
			}, "Function should not panic")

			assert.Equal(t, tt.projectName, state.Name.ValueString(), "Name should match")
			assert.Equal(t, "existing-id", state.ID.ValueString(), "ID should remain unchanged")

			if tt.expectType {
				if tt.projectType != nil {
					assert.Equal(t, *tt.projectType, state.Type.ValueString(), "Type should match")
				}
			} else {
				assert.True(t, state.Type.IsNull(), "Type should be null")
			}
		})
	}
}

// TestProjectResource_executeProjectUpdate tests the executeProjectUpdate helper method.
func TestProjectResource_executeProjectUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		projectID    string
		projectName  string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:        "successful update",
			projectID:   "proj-123",
			projectName: "Updated Project",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/projects/proj-123" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:        "project not found",
			projectID:   "proj-404",
			projectName: "Missing Project",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Project not found"}`))
			},
			expectError: true,
		},
		{
			name:        "API error",
			projectID:   "proj-500",
			projectName: "Error Project",
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

			r := &ProjectResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				ID:   types.StringValue(tt.projectID),
				Name: types.StringValue(tt.projectName),
			}
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}

			result := r.executeProjectUpdate(ctx, plan, resp)

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

// TestProjectResource_findProjectAfterUpdate tests the findProjectAfterUpdate helper method.
func TestProjectResource_findProjectAfterUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		projectID    string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectNil    bool
		expectName   string
	}{
		{
			name:      "project found after update",
			projectID: "proj-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/projects" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{
							{
								"id":   "proj-123",
								"name": "Updated Project",
								"type": "team",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectNil:  false,
			expectName: "Updated Project",
		},
		{
			name:      "project not found after update",
			projectID: "proj-404",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/projects" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectNil: true,
		},
		{
			name:      "API error on list",
			projectID: "proj-500",
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

			r := &ProjectResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				ID: types.StringValue(tt.projectID),
			}
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}

			result := r.findProjectAfterUpdate(ctx, plan, resp)

			if tt.expectNil {
				assert.Nil(t, result, "Should return nil")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.NotNil(t, result, "Should return project")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, tt.expectName, result.Name, "Project name should match")
			}
		})
	}
}

// TestProjectResource_updateModelFromProject tests the updateModelFromProject helper method.
func TestProjectResource_updateModelFromProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		projectName string
		projectType *string
		expectType  bool
	}{
		{
			name:        "update model with all fields",
			projectName: "Model Project",
			projectType: stringPtr("team"),
			expectType:  true,
		},
		{
			name:        "update model without type",
			projectName: "No Type Model",
			projectType: nil,
			expectType:  false,
		},
		{
			name:        "update model with empty type",
			projectName: "Empty Type Model",
			projectType: stringPtr(""),
			expectType:  true,
		},
		{
			name:        "edge case with empty name (error handling)",
			projectName: "",
			projectType: stringPtr("team"),
			expectType:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &ProjectResource{}
			model := &models.Resource{
				ID: types.StringValue("existing-id"),
			}
			project := &n8nsdk.Project{
				Name: tt.projectName,
				Type: tt.projectType,
			}

			// Test that function doesn't panic even with edge cases
			assert.NotPanics(t, func() {
				r.updateModelFromProject(project, model)
			}, "Function should not panic")

			assert.Equal(t, tt.projectName, model.Name.ValueString(), "Name should match")
			assert.Equal(t, "existing-id", model.ID.ValueString(), "ID should remain unchanged")

			if tt.expectType {
				if tt.projectType != nil {
					assert.Equal(t, *tt.projectType, model.Type.ValueString(), "Type should match")
				}
			} else {
				assert.True(t, model.Type.IsNull(), "Type should be null")
			}
		})
	}
}

// stringPtr is a helper function to create string pointers for tests.
func stringPtr(s string) *string {
	return &s
}
