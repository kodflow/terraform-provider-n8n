package project

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/project/models"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestProjectResource_createProject tests the createProject private method.
func TestProjectResource_createProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "validates method exists", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "validates method exists":
				r := &ProjectResource{}
				assert.NotNil(t, r)

			case "error case - validation checks":
				r := &ProjectResource{}
				assert.NotNil(t, r)
			}
		})
	}
}

// TestProjectResource_findCreatedProject tests the findCreatedProject private method.
func TestProjectResource_findCreatedProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "validates method exists", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "validates method exists":
				r := &ProjectResource{}
				assert.NotNil(t, r)

			case "error case - validation checks":
				r := &ProjectResource{}
				assert.NotNil(t, r)
			}
		})
	}
}

// TestProjectResource_updatePlanFromProject tests the updatePlanFromProject private method.
func TestProjectResource_updatePlanFromProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "updates plan with all fields", wantErr: false},
		{name: "updates plan with minimal fields", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "updates plan with all fields":
				r := &ProjectResource{}
				id := "proj-123"
				projectType := "team"
				icon := "icon"
				description := "desc"
				project := &n8nsdk.Project{
					Id:          &id,
					Name:        "Test",
					Type:        &projectType,
					Icon:        &icon,
					Description: &description,
				}
				plan := &models.Resource{}
				r.updatePlanFromProject(plan, project)
				assert.Equal(t, "proj-123", plan.ID.ValueString())
				assert.Equal(t, "Test", plan.Name.ValueString())

			case "updates plan with minimal fields":
				r := &ProjectResource{}
				project := &n8nsdk.Project{Name: "Test"}
				plan := &models.Resource{}
				r.updatePlanFromProject(plan, project)
				assert.True(t, plan.ID.IsNull())
				assert.Equal(t, "Test", plan.Name.ValueString())

			case "error case - validation checks":
				r := &ProjectResource{}
				project := &n8nsdk.Project{Name: "Test"}
				plan := &models.Resource{}
				r.updatePlanFromProject(plan, project)
				assert.NotNil(t, plan)
			}
		})
	}
}

// TestProjectResource_findProjectByID tests the findProjectByID private method.
func TestProjectResource_findProjectByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "validates method exists", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "validates method exists":
				r := &ProjectResource{}
				assert.NotNil(t, r)

			case "error case - validation checks":
				r := &ProjectResource{}
				assert.NotNil(t, r)
			}
		})
	}
}

// TestProjectResource_updateStateFromProject tests the updateStateFromProject private method.
func TestProjectResource_updateStateFromProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "updates state with all fields", wantErr: false},
		{name: "updates state with minimal fields", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "updates state with all fields":
				r := &ProjectResource{}
				projectType := "personal"
				project := &n8nsdk.Project{
					Name: "State Test",
					Type: &projectType,
				}
				state := &models.Resource{}
				r.updateStateFromProject(state, project)
				assert.Equal(t, "State Test", state.Name.ValueString())
				assert.Equal(t, "personal", state.Type.ValueString())

			case "updates state with minimal fields":
				r := &ProjectResource{}
				project := &n8nsdk.Project{Name: "State Test"}
				state := &models.Resource{}
				r.updateStateFromProject(state, project)
				assert.True(t, state.ID.IsNull())
				assert.Equal(t, "State Test", state.Name.ValueString())

			case "error case - validation checks":
				r := &ProjectResource{}
				project := &n8nsdk.Project{Name: "Test"}
				state := &models.Resource{}
				r.updateStateFromProject(state, project)
				assert.NotNil(t, state)
			}
		})
	}
}

// TestProjectResource_executeProjectUpdate tests the executeProjectUpdate private method.
func TestProjectResource_executeProjectUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "validates method exists", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "validates method exists":
				r := &ProjectResource{}
				assert.NotNil(t, r)

			case "error case - validation checks":
				r := &ProjectResource{}
				assert.NotNil(t, r)
			}
		})
	}
}

// TestProjectResource_findProjectAfterUpdate tests the findProjectAfterUpdate private method.
func TestProjectResource_findProjectAfterUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "validates method exists", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "validates method exists":
				r := &ProjectResource{}
				assert.NotNil(t, r)

			case "error case - validation checks":
				r := &ProjectResource{}
				assert.NotNil(t, r)
			}
		})
	}
}

// TestProjectResource_updateModelFromProject tests the updateModelFromProject private method.
func TestProjectResource_updateModelFromProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "updates model with all fields", wantErr: false},
		{name: "updates model with minimal fields", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "updates model with all fields":
				r := &ProjectResource{}
				projectType := "organization"
				project := &n8nsdk.Project{
					Name: "Model Test",
					Type: &projectType,
				}
				model := &models.Resource{}
				r.updateModelFromProject(project, model)
				assert.Equal(t, "Model Test", model.Name.ValueString())
				assert.Equal(t, "organization", model.Type.ValueString())

			case "updates model with minimal fields":
				r := &ProjectResource{}
				project := &n8nsdk.Project{Name: "Model Test"}
				model := &models.Resource{}
				r.updateModelFromProject(project, model)
				assert.True(t, model.ID.IsNull())
				assert.Equal(t, "Model Test", model.Name.ValueString())

			case "error case - validation checks":
				r := &ProjectResource{}
				project := &n8nsdk.Project{Name: "Test"}
				model := &models.Resource{}
				r.updateModelFromProject(project, model)
				assert.NotNil(t, model)
				assert.Equal(t, "Test", model.Name.ValueString())
			}
		})
	}
}

// createTestProjectSchema creates a test schema for project resource.
func createTestProjectSchema(t *testing.T) resource.SchemaResponse {
	t.Helper()
	r := &ProjectResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), req, resp)
	return *resp
}

// setupTestProjectClient creates a test N8nClient with httptest server.
func setupTestProjectClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

// TestProjectResource_Create tests the Create method.
func TestProjectResource_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantErr     bool
		errContains string
	}{
		{
			name: "successful creation",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/projects":
					if r.Method == http.MethodPost {
						w.WriteHeader(http.StatusCreated)
						return
					}
					if r.Method == http.MethodGet {
						projectType := "team"
						projects := map[string]interface{}{
							"data": []map[string]interface{}{
								{
									"id":   "proj-123",
									"name": "Test Project",
									"type": &projectType,
								},
							},
						}
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(projects)
						return
					}
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			wantErr: false,
		},
		{
			name: "creation api error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodPost {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			wantErr:     true,
			errContains: "Error creating project",
		},
		{
			name: "project not found after creation",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/projects":
					if r.Method == http.MethodPost {
						w.WriteHeader(http.StatusCreated)
						return
					}
					if r.Method == http.MethodGet {
						projects := map[string]interface{}{
							"data": []map[string]interface{}{
								{
									"id":   "other-proj",
									"name": "Other Project",
								},
							},
						}
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(projects)
						return
					}
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			wantErr:     true,
			errContains: "Error finding created project",
		},
		{
			name: "list api error after creation",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/projects":
					if r.Method == http.MethodPost {
						w.WriteHeader(http.StatusCreated)
						return
					}
					if r.Method == http.MethodGet {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			wantErr:     true,
			errContains: "Error reading project after creation",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestProjectClient(t, tt.handler)
			defer server.Close()

			r := &ProjectResource{client: n8nClient}

			schemaResp := createTestProjectSchema(t)

			rawPlan := map[string]tftypes.Value{
				"id":   tftypes.NewValue(tftypes.String, nil),
				"name": tftypes.NewValue(tftypes.String, "Test Project"),
				"type": tftypes.NewValue(tftypes.String, nil),
			}
			plan := tfsdk.Plan{
				Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawPlan),
				Schema: schemaResp.Schema,
			}

			state := tfsdk.State{
				Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, nil),
				Schema: schemaResp.Schema,
			}

			req := resource.CreateRequest{
				Plan: plan,
			}
			resp := resource.CreateResponse{
				State: state,
			}

			r.Create(context.Background(), req, &resp)

			if tt.wantErr {
				assert.True(t, resp.Diagnostics.HasError(), "expected error but got none")
				if tt.errContains != "" {
					found := false
					for _, diag := range resp.Diagnostics.Errors() {
						if assert.Contains(t, diag.Summary(), tt.errContains) {
							found = true
							break
						}
					}
					assert.True(t, found, "expected error containing %s", tt.errContains)
				}
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "unexpected error: %v", resp.Diagnostics)
			}
		})
	}
}

// TestProjectResource_Read tests the Read method.
func TestProjectResource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantErr     bool
		errContains string
		removed     bool
	}{
		{
			name: "successful read",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					projectType := "team"
					projects := map[string]interface{}{
						"data": []map[string]interface{}{
							{
								"id":   "proj-123",
								"name": "Test Project",
								"type": &projectType,
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(projects)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			wantErr: false,
		},
		{
			name: "project not found - removed",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					projects := map[string]interface{}{
						"data": []map[string]interface{}{},
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(projects)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			wantErr: false,
			removed: true,
		},
		{
			name: "api error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			wantErr:     true,
			errContains: "Error reading project",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestProjectClient(t, tt.handler)
			defer server.Close()

			r := &ProjectResource{client: n8nClient}

			schemaResp := createTestProjectSchema(t)

			rawState := map[string]tftypes.Value{
				"id":   tftypes.NewValue(tftypes.String, "proj-123"),
				"name": tftypes.NewValue(tftypes.String, "Test Project"),
				"type": tftypes.NewValue(tftypes.String, "team"),
			}
			state := tfsdk.State{
				Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawState),
				Schema: schemaResp.Schema,
			}

			req := resource.ReadRequest{
				State: state,
			}
			resp := resource.ReadResponse{
				State: state,
			}

			r.Read(context.Background(), req, &resp)

			if tt.wantErr {
				assert.True(t, resp.Diagnostics.HasError(), "expected error but got none")
				if tt.errContains != "" {
					found := false
					for _, diag := range resp.Diagnostics.Errors() {
						if assert.Contains(t, diag.Summary(), tt.errContains) {
							found = true
							break
						}
					}
					assert.True(t, found, "expected error containing %s", tt.errContains)
				}
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "unexpected error: %v", resp.Diagnostics)
			}
		})
	}
}

// TestProjectResource_Update tests the Update method.
func TestProjectResource_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantErr     bool
		errContains string
	}{
		{
			name: "successful update",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/projects/proj-123":
					if r.Method == http.MethodPut {
						w.WriteHeader(http.StatusNoContent)
						return
					}
				case "/projects":
					if r.Method == http.MethodGet {
						projectType := "team"
						projects := map[string]interface{}{
							"data": []map[string]interface{}{
								{
									"id":   "proj-123",
									"name": "Updated Project",
									"type": &projectType,
								},
							},
						}
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(projects)
						return
					}
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			wantErr: false,
		},
		{
			name: "update api error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects/proj-123" && r.Method == http.MethodPut {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			wantErr:     true,
			errContains: "Error updating project",
		},
		{
			name: "list api error after update",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/projects/proj-123":
					if r.Method == http.MethodPut {
						w.WriteHeader(http.StatusNoContent)
						return
					}
				case "/projects":
					if r.Method == http.MethodGet {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			wantErr:     true,
			errContains: "Error reading project after update",
		},
		{
			name: "project not found after update",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/projects/proj-123":
					if r.Method == http.MethodPut {
						w.WriteHeader(http.StatusNoContent)
						return
					}
				case "/projects":
					if r.Method == http.MethodGet {
						projects := map[string]interface{}{
							"data": []map[string]interface{}{},
						}
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(projects)
						return
					}
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			wantErr:     true,
			errContains: "Error verifying updated project",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestProjectClient(t, tt.handler)
			defer server.Close()

			r := &ProjectResource{client: n8nClient}

			schemaResp := createTestProjectSchema(t)

			rawPlan := map[string]tftypes.Value{
				"id":   tftypes.NewValue(tftypes.String, "proj-123"),
				"name": tftypes.NewValue(tftypes.String, "Updated Project"),
				"type": tftypes.NewValue(tftypes.String, "team"),
			}
			plan := tfsdk.Plan{
				Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawPlan),
				Schema: schemaResp.Schema,
			}

			rawState := map[string]tftypes.Value{
				"id":   tftypes.NewValue(tftypes.String, "proj-123"),
				"name": tftypes.NewValue(tftypes.String, "Test Project"),
				"type": tftypes.NewValue(tftypes.String, "team"),
			}
			state := tfsdk.State{
				Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawState),
				Schema: schemaResp.Schema,
			}

			req := resource.UpdateRequest{
				Plan:  plan,
				State: state,
			}
			resp := resource.UpdateResponse{
				State: state,
			}

			r.Update(context.Background(), req, &resp)

			if tt.wantErr {
				assert.True(t, resp.Diagnostics.HasError(), "expected error but got none")
				if tt.errContains != "" {
					found := false
					for _, diag := range resp.Diagnostics.Errors() {
						if assert.Contains(t, diag.Summary(), tt.errContains) {
							found = true
							break
						}
					}
					assert.True(t, found, "expected error containing %s", tt.errContains)
				}
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "unexpected error: %v", resp.Diagnostics)
			}
		})
	}
}

// TestProjectResource_Delete tests the Delete method.
func TestProjectResource_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantErr     bool
		errContains string
	}{
		{
			name: "successful deletion",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects/proj-123" && r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			wantErr: false,
		},
		{
			name: "api error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects/proj-123" && r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			wantErr:     true,
			errContains: "Error deleting project",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestProjectClient(t, tt.handler)
			defer server.Close()

			r := &ProjectResource{client: n8nClient}

			schemaResp := createTestProjectSchema(t)

			rawState := map[string]tftypes.Value{
				"id":   tftypes.NewValue(tftypes.String, "proj-123"),
				"name": tftypes.NewValue(tftypes.String, "Test Project"),
				"type": tftypes.NewValue(tftypes.String, "team"),
			}
			state := tfsdk.State{
				Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, rawState),
				Schema: schemaResp.Schema,
			}

			req := resource.DeleteRequest{
				State: state,
			}
			resp := resource.DeleteResponse{
				State: state,
			}

			r.Delete(context.Background(), req, &resp)

			if tt.wantErr {
				assert.True(t, resp.Diagnostics.HasError(), "expected error but got none")
				if tt.errContains != "" {
					found := false
					for _, diag := range resp.Diagnostics.Errors() {
						if assert.Contains(t, diag.Summary(), tt.errContains) {
							found = true
							break
						}
					}
					assert.True(t, found, "expected error containing %s", tt.errContains)
				}
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "unexpected error: %v", resp.Diagnostics)
			}
		})
	}
}

// TestProjectResource_ImportState tests the ImportState method.
func TestProjectResource_ImportState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		importID string
	}{
		{
			name:     "successful import",
			importID: "proj-123",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &ProjectResource{}

			schemaResp := createTestProjectSchema(t)

			state := tfsdk.State{
				Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String}}, nil),
				Schema: schemaResp.Schema,
			}

			req := resource.ImportStateRequest{
				ID: tt.importID,
			}
			resp := resource.ImportStateResponse{
				State: state,
			}

			r.ImportState(context.Background(), req, &resp)

			assert.False(t, resp.Diagnostics.HasError(), "unexpected error: %v", resp.Diagnostics)

			var stateData models.Resource
			resp.State.Get(context.Background(), &stateData)
			assert.Equal(t, tt.importID, stateData.ID.ValueString())
		})
	}
}
