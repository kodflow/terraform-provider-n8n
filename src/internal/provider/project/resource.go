// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package project implements n8n project management resources and data sources.
package project

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/project/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
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
func NewProjectResource() (projectResource *ProjectResource) {
	//: Return new empty ProjectResource instance.
	return &ProjectResource{}
}

// NewProjectResourceWrapper creates a new ProjectResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - projectResource: the wrapped ProjectResource instance
func NewProjectResourceWrapper() (projectResource resource.Resource) {
	//: Return the wrapped resource instance.
	return NewProjectResource()
}

// Metadata returns the resource type name.
//
// Params:
//   - ctx: context for request cancellation
//   - req: metadata request
//   - resp: metadata response
func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return when context is cancelled.
		return
	}
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the resource.
//
// Params:
//   - ctx: context for request cancellation
//   - req: schema request
//   - resp: schema response
func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return when context is cancelled.
		return
	}
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
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return when context is cancelled.
		return
	}
	//: Skip configuration when provider data is not yet available.
	if req.ProviderData == nil {
		//: Return early when no provider data is set.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	//: Validate the provider data type is the expected N8nClient.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		//: Return after reporting the type mismatch error.
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
	//: Check for plan parsing errors.
	if resp.Diagnostics.HasError() {
		//: Return after plan parsing failure.
		return
	}

	//: Execute create logic and return if it fails.
	if !r.executeCreateLogic(ctx, &plan, resp) {
		//: Return after create logic failure.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// executeCreateLogic contains the main logic for creating a project.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - resp: Create response
//
// Returns:
//   - ok: True if creation succeeded, false otherwise
func (r *ProjectResource) executeCreateLogic(ctx context.Context, plan *models.Resource, resp *resource.CreateResponse) (ok bool) {
	//: Create project via API and return failure if it errors.
	if !r.createProject(ctx, plan, resp) {
		//: Return failure when project creation API call fails.
		return false
	}

	//: Retrieve created project details to get the server-assigned ID.
	foundProject := r.findCreatedProject(ctx, plan, resp)
	//: Check if project was found after creation.
	if foundProject == nil {
		//: Return failure when created project cannot be found.
		return false
	}

	r.updatePlanFromProject(plan, foundProject)

	//: Return success after project creation and state update.
	return true
}

// createProject sends the project creation request to the API.
//
// Params:
//   - ctx: context for request cancellation
//   - plan: planned project configuration
//   - resp: create response
//
// Returns:
//   - ok: true if project was created successfully
func (r *ProjectResource) createProject(
	ctx context.Context,
	plan *models.Resource,
	resp *resource.CreateResponse,
) (ok bool) {
	projectRequest := n8nsdk.Project{
		Name: plan.Name.ValueString(),
	}

	httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsPost(ctx).
		Project(projectRequest).
		Execute()
	//: Close HTTP response body if present to prevent resource leaks.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Silently discard close error on response body.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				_ = closeErr
			}
		}()
	}

	//: Check for error from the create API call.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			fmt.Sprintf("Could not create project: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		//: Return failure when the API call fails.
		return false
	}

	//: Return success after successful project creation.
	return true
}

// listProjects fetches all projects from the n8n API.
//
// Params:
//   - ctx: context for request cancellation
//
// Returns:
//   - data: list of projects or nil on error
//   - err: error from the API call
func (r *ProjectResource) listProjects(ctx context.Context) (data []n8nsdk.Project, err error) {
	//: Execute the projects list API call.
	projectList, httpResp, apiErr := r.client.APIClient.ProjectsAPI.ProjectsGet(ctx).Execute()
	//: Close HTTP response body if present to prevent resource leaks.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Silently discard close error on response body.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				_ = closeErr
			}
		}()
	}
	//: Return error if API call failed.
	if apiErr != nil {
		//: Return nil list and error when API fails.
		return nil, apiErr
	}
	//: Return projects data if available, or empty slice.
	if projectList.Data != nil {
		//: Return the projects list from the API response.
		return projectList.Data, nil
	}
	//: Return empty list when no data is present.
	return nil, nil
}

// findCreatedProject finds the newly created project by listing all projects.
//
// Params:
//   - ctx: context for request cancellation
//   - plan: planned project configuration
//   - resp: create response
//
// Returns:
//   - project: pointer to the created project, nil if not found
func (r *ProjectResource) findCreatedProject(
	ctx context.Context,
	plan *models.Resource,
	resp *resource.CreateResponse,
) (project *n8nsdk.Project) {
	projects, err := r.listProjects(ctx)
	//: Check for error from the list API call.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project after creation",
			fmt.Sprintf("Project was created but could not retrieve ID: %s", err.Error()),
		)
		//: Return nil to indicate failure.
		return nil
	}

	//: Iterate over all projects to find by name.
	for _, p := range projects {
		//: Return immediately on name match.
		if p.Name == plan.Name.ValueString() {
			//: Return pointer to matched project.
			return &p
		}
	}

	//: Report error if the created project was not found in the list.
	resp.Diagnostics.AddError(
		"Error finding created project",
		fmt.Sprintf("Project with name '%s' was created but not found in list", plan.Name.ValueString()),
	)
	//: Return nil to indicate failure.
	return nil
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
	//: Set type if project has one.
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
	//: Check for state parsing errors.
	if resp.Diagnostics.HasError() {
		//: Return after state parsing failure.
		return
	}

	//: Execute read logic and return if it fails.
	if !r.executeReadLogic(ctx, &state, resp) {
		//: Return after read logic failure.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// executeReadLogic contains the main logic for reading a project.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Read response
//
// Returns:
//   - ok: True if read succeeded, false otherwise
func (r *ProjectResource) executeReadLogic(ctx context.Context, state *models.Resource, resp *resource.ReadResponse) (ok bool) {
	foundProject := r.findProjectByID(ctx, state, resp)
	//: Check if project was found in the list.
	if foundProject == nil {
		//: Return failure when project is not found.
		return false
	}

	r.updateStateFromProject(state, foundProject)

	//: Return success after state update.
	return true
}

// findProjectByID finds a project by ID using the list endpoint.
//
// Params:
//   - ctx: context for request cancellation
//   - state: current state containing project ID
//   - resp: read response
//
// Returns:
//   - project: pointer to the found project, nil if not found or error
func (r *ProjectResource) findProjectByID(
	ctx context.Context,
	state *models.Resource,
	resp *resource.ReadResponse,
) (project *n8nsdk.Project) {
	projects, err := r.listProjects(ctx)
	//: Check for error from the list API call.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project",
			fmt.Sprintf("Could not read project ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		//: Return nil to indicate failure.
		return nil
	}

	//: Iterate over all projects to find by ID.
	for _, p := range projects {
		//: Return immediately when matching ID is found.
		if p.Id != nil && *p.Id == state.ID.ValueString() {
			//: Return pointer to the matched project.
			return &p
		}
	}

	//: Remove resource from state since it was deleted outside Terraform.
	resp.State.RemoveResource(ctx)
	//: Return nil to indicate resource was removed.
	return nil
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
	//: Set type if project has one.
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
	var plan, state models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	//: Check for plan/state parsing errors.
	if resp.Diagnostics.HasError() {
		//: Return after plan/state parsing failure.
		return
	}

	//: Execute update logic and return if it fails.
	if !r.executeUpdateLogic(ctx, &plan, &state, resp) {
		//: Return after update logic failure.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// executeUpdateLogic contains the main logic for updating a project.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - state: The current resource state
//   - resp: Update response
//
// Returns:
//   - ok: True if update succeeded, false otherwise
func (r *ProjectResource) executeUpdateLogic(ctx context.Context, plan, state *models.Resource, resp *resource.UpdateResponse) (ok bool) {
	// Use state.ID for the project ID since plan.ID may be Unknown
	// for Computed attributes.
	projectID := state.ID.ValueString()

	//: Copy ID from state to plan for consistency.
	plan.ID = state.ID

	//: Execute the update API call and return failure if it errors.
	if !r.executeProjectUpdate(ctx, projectID, plan, resp) {
		//: Return failure when update API call fails.
		return false
	}

	//: Verify the update by listing projects.
	foundProject := r.findProjectAfterUpdate(ctx, projectID, resp)
	//: Check if updated project was found.
	if foundProject == nil {
		//: Return failure when updated project cannot be found.
		return false
	}

	//: Update model with found project data.
	r.updateModelFromProject(foundProject, plan)

	//: Return success after update and state synchronization.
	return true
}

// executeProjectUpdate performs the API call to update a project.
//
// Params:
//   - ctx: context for request cancellation
//   - projectID: the project ID to update
//   - plan: planned project model
//   - resp: update response
//
// Returns:
//   - ok: true if successful, false if error occurred
func (r *ProjectResource) executeProjectUpdate(ctx context.Context, projectID string, plan *models.Resource, resp *resource.UpdateResponse) (ok bool) {
	projectRequest := n8nsdk.Project{
		Name: plan.Name.ValueString(),
	}

	httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsProjectIdPut(ctx, projectID).
		Project(projectRequest).
		Execute()
	//: Close HTTP response body if present to prevent resource leaks.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Silently discard close error on response body.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				_ = closeErr
			}
		}()
	}

	//: Check for error from the update API call.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating project",
			fmt.Sprintf("Could not update project ID %s: %s\nHTTP Response: %v", projectID, err.Error(), httpResp),
		)
		//: Return failure when the update API call fails.
		return false
	}

	//: Return success after successful project update.
	return true
}

// findProjectAfterUpdate retrieves and finds the updated project from the list.
//
// Params:
//   - ctx: context for request cancellation
//   - projectID: the project ID to find
//   - resp: update response
//
// Returns:
//   - project: found project or nil if not found
func (r *ProjectResource) findProjectAfterUpdate(ctx context.Context, projectID string, resp *resource.UpdateResponse) (project *n8nsdk.Project) {
	projects, err := r.listProjects(ctx)
	//: Check for error from the list API call.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project after update",
			fmt.Sprintf("Project was updated but could not verify: %s", err.Error()),
		)
		//: Return nil after reporting the API error.
		return nil
	}

	//: Iterate over all projects to find by ID.
	for _, p := range projects {
		//: Return immediately when matching ID is found.
		if p.Id != nil && *p.Id == projectID {
			//: Return pointer to the matched project.
			return &p
		}
	}

	//: Report error if updated project was not found in the list.
	resp.Diagnostics.AddError(
		"Error verifying updated project",
		fmt.Sprintf("Project with ID '%s' was updated but not found in list", projectID),
	)
	//: Return nil after reporting the project not found error.
	return nil
}

// updateModelFromProject updates the model with data from the project.
//
// Params:
//   - project: source project
//   - model: target project model
func (r *ProjectResource) updateModelFromProject(project *n8nsdk.Project, model *models.Resource) {
	model.Name = types.StringValue(project.Name)
	//: Set type if project has one.
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
	//: Check for state parsing errors.
	if resp.Diagnostics.HasError() {
		//: Return after state parsing failure.
		return
	}

	//: Execute delete logic.
	r.executeDeleteLogic(ctx, &state, resp)
}

// executeDeleteLogic contains the main logic for deleting a project.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Delete response
//
// Returns:
//   - ok: True if delete succeeded, false otherwise
func (r *ProjectResource) executeDeleteLogic(ctx context.Context, state *models.Resource, resp *resource.DeleteResponse) (ok bool) {
	// DELETE returns 204 with no body.
	httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsProjectIdDelete(ctx, state.ID.ValueString()).Execute()
	//: Close HTTP response body if present to prevent resource leaks.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Silently discard close error on response body.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				_ = closeErr
			}
		}()
	}

	//: Check for error from the delete API call.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting project",
			fmt.Sprintf("Could not delete project ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		//: Return failure when the delete API call fails.
		return false
	}

	//: Return success after successful project deletion.
	return true
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
