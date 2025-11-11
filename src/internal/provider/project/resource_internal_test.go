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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/project/models"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

func TestProjectResource_createProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectError bool
	}{
		{
			name: "successful project creation",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodPost {
					w.WriteHeader(http.StatusCreated)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: false,
		},
		{
			name: "API returns error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodPost {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: true,
		},
		{
			name: "API returns bad request",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodPost {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"message": "Bad request"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: true,
		},
		{
			name: "network timeout simulation",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodPost {
					w.WriteHeader(http.StatusRequestTimeout)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClient(t, tt.handler)
			defer server.Close()

			r := &ProjectResource{client: n8nClient}
			plan := &models.Resource{}
			resp := &resource.CreateResponse{}

			success := r.createProject(context.Background(), plan, resp)

			if tt.expectError {
				assert.False(t, success)
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.True(t, success)
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestProjectResource_findCreatedProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		planName    string
		expectNil   bool
		expectError bool
	}{
		{
			name:     "project found by name",
			planName: "Test Project",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					id := "proj-123"
					projectType := "team"
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":   id,
								"name": "Test Project",
								"type": projectType,
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
			expectNil:   false,
			expectError: false,
		},
		{
			name:     "project not found in list",
			planName: "Missing Project",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":   "proj-456",
								"name": "Other Project",
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
			expectNil:   true,
			expectError: true,
		},
		{
			name:     "API returns error",
			planName: "Test Project",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:     "empty project list",
			planName: "Test Project",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:     "nil data in response",
			planName: "Test Project",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					response := map[string]interface{}{}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClient(t, tt.handler)
			defer server.Close()

			r := &ProjectResource{client: n8nClient}
			plan := &models.Resource{}
			plan.Name = types.StringValue(tt.planName)
			resp := &resource.CreateResponse{}

			foundProject := r.findCreatedProject(context.Background(), plan, resp)

			if tt.expectNil {
				assert.Nil(t, foundProject)
			} else {
				assert.NotNil(t, foundProject)
			}

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestProjectResource_updatePlanFromProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		project  *n8nsdk.Project
		validate func(*testing.T, *models.Resource)
	}{
		{
			name: "update with all fields",
			project: func() *n8nsdk.Project {
				id := "proj-123"
				projectType := "team"
				return &n8nsdk.Project{
					Id:   &id,
					Name: "Test Project",
					Type: &projectType,
				}
			}(),
			validate: func(t *testing.T, plan *models.Resource) {
				t.Helper()
				assert.Equal(t, "proj-123", plan.ID.ValueString())
				assert.Equal(t, "Test Project", plan.Name.ValueString())
				assert.Equal(t, "team", plan.Type.ValueString())
			},
		},
		{
			name: "update with nil ID",
			project: &n8nsdk.Project{
				Name: "Project Without ID",
			},
			validate: func(t *testing.T, plan *models.Resource) {
				t.Helper()
				assert.True(t, plan.ID.IsNull())
				assert.Equal(t, "Project Without ID", plan.Name.ValueString())
			},
		},
		{
			name: "update with nil type",
			project: func() *n8nsdk.Project {
				id := "proj-456"
				return &n8nsdk.Project{
					Id:   &id,
					Name: "Project Without Type",
				}
			}(),
			validate: func(t *testing.T, plan *models.Resource) {
				t.Helper()
				assert.Equal(t, "proj-456", plan.ID.ValueString())
				assert.Equal(t, "Project Without Type", plan.Name.ValueString())
				assert.True(t, plan.Type.IsNull())
			},
		},
		{
			name:    "error case - nil project",
			project: nil,
			validate: func(t *testing.T, plan *models.Resource) {
				t.Helper()
				// This would panic in real code, validating the test checks behavior
				assert.NotNil(t, plan)
			},
		},
		{
			name: "update empty name",
			project: func() *n8nsdk.Project {
				id := "proj-789"
				return &n8nsdk.Project{
					Id:   &id,
					Name: "",
				}
			}(),
			validate: func(t *testing.T, plan *models.Resource) {
				t.Helper()
				assert.Equal(t, "proj-789", plan.ID.ValueString())
				assert.Equal(t, "", plan.Name.ValueString())
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &ProjectResource{}
			plan := &models.Resource{}

			// Handle nil project case
			if tt.project == nil {
				assert.Panics(t, func() {
					r.updatePlanFromProject(plan, tt.project)
				})
				return
			}

			r.updatePlanFromProject(plan, tt.project)

			tt.validate(t, plan)
		})
	}
}

func TestProjectResource_findProjectByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		stateID     string
		expectNil   bool
		expectError bool
	}{
		{
			name:    "project found by ID",
			stateID: "proj-123",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					id := "proj-123"
					projectType := "team"
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":   id,
								"name": "Test Project",
								"type": projectType,
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
			expectNil:   false,
			expectError: false,
		},
		{
			name:    "project not found - removed from state",
			stateID: "proj-999",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":   "proj-123",
								"name": "Other Project",
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
			expectNil:   true,
			expectError: false,
		},
		{
			name:    "API returns error",
			stateID: "proj-123",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:    "empty project list",
			stateID: "proj-123",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: false,
		},
		{
			name:    "project with nil ID in list",
			stateID: "proj-123",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"name": "Project Without ID",
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
			expectNil:   true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClient(t, tt.handler)
			defer server.Close()

			r := &ProjectResource{client: n8nClient}
			state := &models.Resource{}
			state.ID = types.StringValue(tt.stateID)

			// Create a proper tfsdk.State with schema
			testSchema := schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id":   schema.StringAttribute{Computed: true},
					"name": schema.StringAttribute{Required: true},
					"type": schema.StringAttribute{Computed: true},
				},
			}
			rawState := map[string]tftypes.Value{
				"id":   tftypes.NewValue(tftypes.String, tt.stateID),
				"name": tftypes.NewValue(tftypes.String, "Test Project"),
				"type": tftypes.NewValue(tftypes.String, nil),
			}
			tfState := tfsdk.State{
				Raw: tftypes.NewValue(tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id":   tftypes.String,
						"name": tftypes.String,
						"type": tftypes.String,
					},
				}, rawState),
				Schema: testSchema,
			}
			resp := &resource.ReadResponse{
				State: tfState,
			}

			foundProject := r.findProjectByID(context.Background(), state, resp)

			if tt.expectNil {
				assert.Nil(t, foundProject)
			} else {
				assert.NotNil(t, foundProject)
			}

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestProjectResource_updateStateFromProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		project  *n8nsdk.Project
		validate func(*testing.T, *models.Resource)
	}{
		{
			name: "update state with all fields",
			project: func() *n8nsdk.Project {
				projectType := "team"
				return &n8nsdk.Project{
					Name: "Updated Project",
					Type: &projectType,
				}
			}(),
			validate: func(t *testing.T, state *models.Resource) {
				t.Helper()
				assert.Equal(t, "Updated Project", state.Name.ValueString())
				assert.Equal(t, "team", state.Type.ValueString())
			},
		},
		{
			name: "update state with nil type",
			project: &n8nsdk.Project{
				Name: "Project Without Type",
			},
			validate: func(t *testing.T, state *models.Resource) {
				t.Helper()
				assert.Equal(t, "Project Without Type", state.Name.ValueString())
				assert.True(t, state.Type.IsNull())
			},
		},
		{
			name: "update state with empty name",
			project: &n8nsdk.Project{
				Name: "",
			},
			validate: func(t *testing.T, state *models.Resource) {
				t.Helper()
				assert.Equal(t, "", state.Name.ValueString())
			},
		},
		{
			name: "update state preserves existing ID",
			project: func() *n8nsdk.Project {
				projectType := "personal"
				return &n8nsdk.Project{
					Name: "New Name",
					Type: &projectType,
				}
			}(),
			validate: func(t *testing.T, state *models.Resource) {
				t.Helper()
				assert.Equal(t, "New Name", state.Name.ValueString())
				assert.Equal(t, "personal", state.Type.ValueString())
			},
		},
		{
			name:    "error case - nil project",
			project: nil,
			validate: func(t *testing.T, state *models.Resource) {
				t.Helper()
				// This would panic in real code, validating the test checks behavior
				assert.NotNil(t, state)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &ProjectResource{}
			state := &models.Resource{}

			// Handle nil project case
			if tt.project == nil {
				assert.Panics(t, func() {
					r.updateStateFromProject(state, tt.project)
				})
				return
			}

			r.updateStateFromProject(state, tt.project)

			tt.validate(t, state)
		})
	}
}

func TestProjectResource_executeProjectUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectError bool
	}{
		{
			name: "successful update",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects/proj-123" && r.Method == http.MethodPut {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: false,
		},
		{
			name: "API returns error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects/proj-123" && r.Method == http.MethodPut {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: true,
		},
		{
			name: "not found error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects/proj-123" && r.Method == http.MethodPut {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "Project not found"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: true,
		},
		{
			name: "bad request error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects/proj-123" && r.Method == http.MethodPut {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"message": "Bad request"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClient(t, tt.handler)
			defer server.Close()

			r := &ProjectResource{client: n8nClient}
			plan := &models.Resource{}
			plan.ID = types.StringValue("proj-123")
			plan.Name = types.StringValue("Updated Project")
			resp := &resource.UpdateResponse{}

			success := r.executeProjectUpdate(context.Background(), plan, resp)

			if tt.expectError {
				assert.False(t, success)
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.True(t, success)
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestProjectResource_findProjectAfterUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		planID      string
		expectNil   bool
		expectError bool
	}{
		{
			name:   "project found after update",
			planID: "proj-123",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					id := "proj-123"
					projectType := "team"
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":   id,
								"name": "Updated Project",
								"type": projectType,
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
			expectNil:   false,
			expectError: false,
		},
		{
			name:   "project not found after update",
			planID: "proj-123",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":   "proj-456",
								"name": "Other Project",
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
			expectNil:   true,
			expectError: true,
		},
		{
			name:   "API returns error",
			planID: "proj-123",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:   "empty project list",
			planID: "proj-123",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:   "nil data in response",
			planID: "proj-123",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					response := map[string]interface{}{}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClient(t, tt.handler)
			defer server.Close()

			r := &ProjectResource{client: n8nClient}
			plan := &models.Resource{}
			plan.ID = types.StringValue(tt.planID)
			resp := &resource.UpdateResponse{}

			foundProject := r.findProjectAfterUpdate(context.Background(), plan, resp)

			if tt.expectNil {
				assert.Nil(t, foundProject)
			} else {
				assert.NotNil(t, foundProject)
			}

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestProjectResource_updateModelFromProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		project  *n8nsdk.Project
		validate func(*testing.T, *models.Resource)
	}{
		{
			name: "update model with all fields",
			project: func() *n8nsdk.Project {
				projectType := "team"
				return &n8nsdk.Project{
					Name: "Updated Project",
					Type: &projectType,
				}
			}(),
			validate: func(t *testing.T, model *models.Resource) {
				t.Helper()
				assert.Equal(t, "Updated Project", model.Name.ValueString())
				assert.Equal(t, "team", model.Type.ValueString())
			},
		},
		{
			name: "update model with nil type",
			project: &n8nsdk.Project{
				Name: "Project Without Type",
			},
			validate: func(t *testing.T, model *models.Resource) {
				t.Helper()
				assert.Equal(t, "Project Without Type", model.Name.ValueString())
				assert.True(t, model.Type.IsNull())
			},
		},
		{
			name: "update model with empty name",
			project: &n8nsdk.Project{
				Name: "",
			},
			validate: func(t *testing.T, model *models.Resource) {
				t.Helper()
				assert.Equal(t, "", model.Name.ValueString())
			},
		},
		{
			name: "update model multiple times",
			project: func() *n8nsdk.Project {
				projectType := "personal"
				return &n8nsdk.Project{
					Name: "Final Name",
					Type: &projectType,
				}
			}(),
			validate: func(t *testing.T, model *models.Resource) {
				t.Helper()
				assert.Equal(t, "Final Name", model.Name.ValueString())
				assert.Equal(t, "personal", model.Type.ValueString())
			},
		},
		{
			name:    "error case - nil project",
			project: nil,
			validate: func(t *testing.T, model *models.Resource) {
				t.Helper()
				// This would panic in real code, validating the test checks behavior
				assert.NotNil(t, model)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &ProjectResource{}
			model := &models.Resource{}

			// Handle nil project case
			if tt.project == nil {
				assert.Panics(t, func() {
					r.updateModelFromProject(tt.project, model)
				})
				return
			}

			r.updateModelFromProject(tt.project, model)

			tt.validate(t, model)
		})
	}
}

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

// TestProjectResource_Create_CRUD tests the Create method with full CRUD coverage.
func TestProjectResource_Create_CRUD(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "create with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodPost && r.URL.Path == "/projects" {
						w.WriteHeader(http.StatusCreated)
						return
					}
					if r.Method == http.MethodGet && r.URL.Path == "/projects" {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`{"data":[{"id":"proj-123","name":"test-project","type":"team"}]}`))
						return
					}
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &ProjectResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, nil),
					"name": tftypes.NewValue(tftypes.String, "test-project"),
					"type": tftypes.NewValue(tftypes.String, nil),
				})

				req := resource.CreateRequest{
					Plan: tfsdk.Plan{
						Raw:    planRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.CreateResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Create(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError(), "Expected no errors")
			},
		},
		{
			name: "error - create with invalid plan",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					t.Fatal("Should not reach API")
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &ProjectResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				planRaw := tftypes.NewValue(tftypes.String, "invalid")

				req := resource.CreateRequest{
					Plan: tfsdk.Plan{
						Raw:    planRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.CreateResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Create(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "Expected error with invalid plan")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestProjectResource_Read_CRUD tests the Read method with full CRUD coverage.
func TestProjectResource_Read_CRUD(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "read with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodGet && r.URL.Path == "/projects" {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`{"data":[{"id":"proj-123","name":"test-project","type":"team"}]}`))
					}
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &ProjectResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, "proj-123"),
					"name": tftypes.NewValue(tftypes.String, "test-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
				})

				req := resource.ReadRequest{
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.ReadResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Read(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError(), "Expected no errors")
			},
		},
		{
			name: "error - read with invalid state",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					t.Fatal("Should not reach API")
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &ProjectResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				stateRaw := tftypes.NewValue(tftypes.String, "invalid")

				req := resource.ReadRequest{
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.ReadResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "Expected error with invalid state")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestProjectResource_Update_CRUD tests the Update method with full CRUD coverage.
func TestProjectResource_Update_CRUD(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "update with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodPut && r.URL.Path == "/projects/proj-123" {
						w.WriteHeader(http.StatusNoContent)
						return
					}
					if r.Method == http.MethodGet && r.URL.Path == "/projects" {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`{"data":[{"id":"proj-123","name":"updated-project","type":"team"}]}`))
						return
					}
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &ProjectResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, "proj-123"),
					"name": tftypes.NewValue(tftypes.String, "test-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
				})

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, "proj-123"),
					"name": tftypes.NewValue(tftypes.String, "updated-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
				})

				req := resource.UpdateRequest{
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
					Plan: tfsdk.Plan{
						Raw:    planRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.UpdateResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Update(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError(), "Expected no errors")
			},
		},
		{
			name: "error - update with invalid plan",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					t.Fatal("Should not reach API")
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &ProjectResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, "proj-123"),
					"name": tftypes.NewValue(tftypes.String, "test-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
				})

				planRaw := tftypes.NewValue(tftypes.String, "invalid")

				req := resource.UpdateRequest{
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
					Plan: tfsdk.Plan{
						Raw:    planRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.UpdateResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Update(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "Expected error with invalid plan")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestProjectResource_Delete_CRUD tests the Delete method with full CRUD coverage.
func TestProjectResource_Delete_CRUD(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "delete with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodDelete && r.URL.Path == "/projects/proj-123" {
						w.WriteHeader(http.StatusNoContent)
					}
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &ProjectResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, "proj-123"),
					"name": tftypes.NewValue(tftypes.String, "test-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
				})

				req := resource.DeleteRequest{
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.DeleteResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Delete(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError(), "Expected no errors")
			},
		},
		{
			name: "error - delete with invalid state",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					t.Fatal("Should not reach API")
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &ProjectResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				stateRaw := tftypes.NewValue(tftypes.String, "invalid")

				req := resource.DeleteRequest{
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.DeleteResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Delete(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "Expected error with invalid state")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestProjectResource_ImportState_CRUD tests the ImportState method with full CRUD coverage.
func TestProjectResource_ImportState_CRUD(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "import state with valid ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &ProjectResource{}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				emptyValue := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, nil),
					"name": tftypes.NewValue(tftypes.String, nil),
					"type": tftypes.NewValue(tftypes.String, nil),
				})

				req := resource.ImportStateRequest{
					ID: "proj-123",
				}

				resp := &resource.ImportStateResponse{
					State: tfsdk.State{
						Raw:    emptyValue,
						Schema: schemaResp.Schema,
					},
				}

				r.ImportState(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError(), "Expected no errors")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}
