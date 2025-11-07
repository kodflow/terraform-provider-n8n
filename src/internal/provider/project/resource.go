package project

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/project/models"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)

// Ensure ProjectResource implements required interfaces.
var (
	_ resource.Resource                = &ProjectResource{}
	_ ProjectResourceInterface         = &ProjectResource{}
	_ resource.ResourceWithConfigure   = &ProjectResource{}
	_ resource.ResourceWithImportState = &ProjectResource{}
)

// ProjectResourceInterface defines the interface for ProjectResource.
type ProjectResourceInterface interface {
	resource.Resource
	Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse)
	Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse)
	Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse)
	Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse)
	Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse)
	Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse)
	Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse)
	ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse)
}

// ProjectResource defines the resource implementation for n8n projects.
// Note: n8n API has limitations - POST returns 201 with no body, no GET by ID endpoint.
// We work around this by using the LIST endpoint and filtering.
type ProjectResource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewProjectResource creates a new ProjectResource instance.
//
// Returns:
//   - resource.Resource: new ProjectResource instance
func NewProjectResource() *ProjectResource {
	// Return result.
	return &ProjectResource{}
}

// NewProjectResourceWrapper creates a new ProjectResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - resource.Resource: the wrapped ProjectResource instance
func NewProjectResourceWrapper() resource.Resource {
	// Return the wrapped resource instance.
	return NewProjectResource()
}

// Metadata returns the resource type name.
//
// Params:
//   - ctx: context for request cancellation
//   - req: metadata request
//   - resp: metadata response
func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the resource.
//
// Params:
//   - ctx: context for request cancellation
//   - req: schema request
//   - resp: schema response
func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "n8n project resource. Note: API limitations require workarounds for Read operations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Project identifier",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Project name",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Project type",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
//
// Params:
//   - ctx: context for request cancellation
//   - req: configure request
//   - resp: configure response
func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
		// Return result.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	// Check condition.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		// Return result.
		return
	}

	r.client = clientData
}

// Create creates the resource and sets the initial Terraform state.
// Workaround: API returns 201 with no body, so we must call LIST to get the ID.
//
// Params:
//   - ctx: context for request cancellation
//   - req: create request
//   - resp: create response
func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Create project via API
	if !r.createProject(ctx, &plan, resp) {
		// Return result.
		return
	}

	// Retrieve created project details
	foundProject := r.findCreatedProject(ctx, &plan, resp)
	// Check for nil value.
	if foundProject == nil {
		// Return result.
		return
	}

	r.updatePlanFromProject(&plan, foundProject)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// createProject sends the project creation request to the API.
//
// Params:
//   - ctx: context for request cancellation
//   - plan: planned project configuration
//   - resp: create response
//
// Returns:
//   - success: true if project was created successfully
func (r *ProjectResource) createProject(
	ctx context.Context,
	plan *models.Resource,
	resp *resource.CreateResponse,
) bool {
	projectRequest := n8nsdk.Project{
		Name: plan.Name.ValueString(),
	}

	httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsPost(ctx).
		Project(projectRequest).
		Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			fmt.Sprintf("Could not create project: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return failure status.
		return false
	}

	// Return success status.
	return true
}

// findCreatedProject finds the newly created project by listing all projects.
//
// Params:
//   - ctx: context for request cancellation
//   - plan: planned project configuration
//   - resp: create response
//
// Returns:
//   - foundProject: pointer to the created project, nil if not found
func (r *ProjectResource) findCreatedProject(
	ctx context.Context,
	plan *models.Resource,
	resp *resource.CreateResponse,
) *n8nsdk.Project {
	projectList, httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsGet(ctx).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project after creation",
			fmt.Sprintf("Project was created but could not retrieve ID: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return nil to indicate failure.
		return nil
	}

	// Find our project by name
	var foundProject *n8nsdk.Project
	// Check for non-nil value.
	if projectList.Data != nil {
		// Iterate over items.
		for _, p := range projectList.Data {
			// Check condition.
			if p.Name == plan.Name.ValueString() {
				foundProject = &p
				break
			}
		}
	}

	// Check for nil value.
	if foundProject == nil {
		resp.Diagnostics.AddError(
			"Error finding created project",
			fmt.Sprintf("Project with name '%s' was created but not found in list", plan.Name.ValueString()),
		)
		// Return nil to indicate failure.
		return nil
	}

	// Return found project.
	return foundProject
}

// updatePlanFromProject updates the plan model with data from the created project.
//
// Params:
//   - plan: planned project configuration to update
//   - project: created project data from API
func (r *ProjectResource) updatePlanFromProject(
	plan *models.Resource,
	project *n8nsdk.Project,
) {
	plan.ID = types.StringPointerValue(project.Id)
	plan.Name = types.StringValue(project.Name)
	// Check for non-nil value.
	if project.Type != nil {
		plan.Type = types.StringPointerValue(project.Type)
	}
}

// Read refreshes the Terraform state with the latest data.
// Workaround: No GET by ID endpoint, so we use LIST and filter.
//
// Params:
//   - ctx: context for request cancellation
//   - req: read request
//   - resp: read response
func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.Resource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	foundProject := r.findProjectByID(ctx, &state, resp)
	// Check for nil value.
	if foundProject == nil {
		// Return result.
		return
	}

	r.updateStateFromProject(&state, foundProject)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// findProjectByID finds a project by ID using the list endpoint.
//
// Params:
//   - ctx: context for request cancellation
//   - state: current state containing project ID
//   - resp: read response
//
// Returns:
//   - foundProject: pointer to the found project, nil if not found or error
func (r *ProjectResource) findProjectByID(
	ctx context.Context,
	state *models.Resource,
	resp *resource.ReadResponse,
) *n8nsdk.Project {
	projectList, httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsGet(ctx).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project",
			fmt.Sprintf("Could not read project ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		// Return nil to indicate failure.
		return nil
	}

	// Find our project by ID
	var foundProject *n8nsdk.Project
	// Check for non-nil value.
	if projectList.Data != nil {
		// Iterate over items.
		for _, p := range projectList.Data {
			// Check for non-nil value.
			if p.Id != nil && *p.Id == state.ID.ValueString() {
				foundProject = &p
				break
			}
		}
	}

	// Check for nil value.
	if foundProject == nil {
		// Project not found = deleted outside Terraform
		resp.State.RemoveResource(ctx)
		// Return nil to indicate resource was removed.
		return nil
	}

	// Return found project.
	return foundProject
}

// updateStateFromProject updates the state model with data from the found project.
//
// Params:
//   - state: current state to update
//   - project: project data from API
func (r *ProjectResource) updateStateFromProject(
	state *models.Resource,
	project *n8nsdk.Project,
) {
	state.Name = types.StringValue(project.Name)
	// Check for non-nil value.
	if project.Type != nil {
		state.Type = types.StringPointerValue(project.Type)
	}
}

// Update updates the resource and sets the updated Terraform state on success.
// Workaround: API returns 204 with no body, so we must call LIST to verify.
//
// Params:
//   - ctx: context for request cancellation
//   - req: update request
//   - resp: update response
func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute the update API call
	if !r.executeProjectUpdate(ctx, &plan, resp) {
		// Return result.
		return
	}

	// Verify the update by listing projects
	foundProject := r.findProjectAfterUpdate(ctx, &plan, resp)
	// Check for nil value.
	if foundProject == nil {
		// Return result.
		return
	}

	// Update model with found project data
	r.updateModelFromProject(foundProject, &plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// executeProjectUpdate performs the API call to update a project.
//
// Params:
//   - ctx: context for request cancellation
//   - plan: planned project model
//   - resp: update response
//
// Returns:
//   - bool: true if successful, false if error occurred
func (r *ProjectResource) executeProjectUpdate(ctx context.Context, plan *models.Resource, resp *resource.UpdateResponse) bool {
	projectRequest := n8nsdk.Project{
		Name: plan.Name.ValueString(),
	}

	httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsProjectIdPut(ctx, plan.ID.ValueString()).
		Project(projectRequest).
		Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating project",
			fmt.Sprintf("Could not update project ID %s: %s\nHTTP Response: %v", plan.ID.ValueString(), err.Error(), httpResp),
		)
		return false
	}

	return true
}

// findProjectAfterUpdate retrieves and finds the updated project from the list.
//
// Params:
//   - ctx: context for request cancellation
//   - plan: planned project model
//   - resp: update response
//
// Returns:
//   - *n8nsdk.Project: found project or nil if not found
func (r *ProjectResource) findProjectAfterUpdate(ctx context.Context, plan *models.Resource, resp *resource.UpdateResponse) *n8nsdk.Project {
	projectList, httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsGet(ctx).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project after update",
			fmt.Sprintf("Project was updated but could not verify: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return with error.
		return nil
	}

	// Find our project by ID
	var foundProject *n8nsdk.Project
	// Check for non-nil value.
	if projectList.Data != nil {
		// Iterate over items.
		for _, p := range projectList.Data {
			// Check for non-nil value.
			if p.Id != nil && *p.Id == plan.ID.ValueString() {
				foundProject = &p
				break
			}
		}
	}

	// Check for nil value.
	if foundProject == nil {
		resp.Diagnostics.AddError(
			"Error verifying updated project",
			fmt.Sprintf("Project with ID '%s' was updated but not found in list", plan.ID.ValueString()),
		)
		// Return with error.
		return nil
	}

	// Return result.
	return foundProject
}

// updateModelFromProject updates the model with data from the project.
//
// Params:
//   - project: source project
//   - model: target project model
func (r *ProjectResource) updateModelFromProject(project *n8nsdk.Project, model *models.Resource) {
	model.Name = types.StringValue(project.Name)
	// Check for non-nil value.
	if project.Type != nil {
		model.Type = types.StringPointerValue(project.Type)
	}
}

// Delete deletes the resource and removes the Terraform state on success.
//
// Params:
//   - ctx: context for request cancellation
//   - req: delete request
//   - resp: delete response
func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.Resource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// DELETE returns 204 with no body
	httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsProjectIdDelete(ctx, state.ID.ValueString()).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting project",
			fmt.Sprintf("Could not delete project ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		// Return result.
		return
	}
}

// ImportState imports the resource into Terraform state.
//
// Params:
//   - ctx: context for request cancellation
//   - req: import state request
//   - resp: import state response
func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
