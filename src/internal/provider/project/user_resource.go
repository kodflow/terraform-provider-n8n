// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package project implements n8n project management resources and data sources.
package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/project/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
)

// CompositeIDParts is the expected number of parts in a composite ID (project_id/user_id).
const CompositeIDParts int = 2

// Ensure ProjectUserResource implements required interfaces.
var (
	_ resource.Resource                = &ProjectUserResource{}
	_ ProjectUserResourceInterface     = &ProjectUserResource{}
	_ resource.ResourceWithConfigure   = &ProjectUserResource{}
	_ resource.ResourceWithImportState = &ProjectUserResource{}
)

// ProjectUserResourceInterface defines the interface for ProjectUserResource.
type ProjectUserResourceInterface interface {
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

// ProjectUserResource defines the resource implementation for project-user relationships.
// Terraform resource that manages CRUD operations for user memberships in n8n projects via the n8n API.
type ProjectUserResource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewProjectUserResource creates a new ProjectUserResource instance.
// Returns:
//   - *ProjectUserResource: new ProjectUserResource instance
func NewProjectUserResource() (projectUserResource *ProjectUserResource) {
	//: Return new empty ProjectUserResource instance.
	return &ProjectUserResource{}
}

// NewProjectUserResourceWrapper creates a new ProjectUserResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - projectUserResource: the wrapped ProjectUserResource instance
func NewProjectUserResourceWrapper() (projectUserResource resource.Resource) {
	//: Return the wrapped resource instance.
	return NewProjectUserResource()
}

// Metadata returns the resource type name.
// Params:
//   - ctx: context
//   - req: metadata request
//   - resp: metadata response
//
// Returns:
//   - none
func (r *ProjectUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return when context is cancelled.
		return
	}
	resp.TypeName = req.ProviderTypeName + "_project_user"
}

// Schema defines the schema for the resource.
// Params:
//   - ctx: context
//   - req: schema request
//   - resp: schema response
//
// Returns:
//   - none
func (r *ProjectUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return when context is cancelled.
		return
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages user membership and roles within n8n projects. Allows adding users to projects, changing their roles, and removing them from projects.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Resource identifier in the format project_id/user_id",
				Computed:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "ID of the project",
				Required:            true,
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "ID of the user",
				Required:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "Role of the user in the project (e.g., 'project:admin', 'project:editor', 'project:viewer')",
				Required:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
// Params:
//   - ctx: context
//   - req: configure request
//   - resp: configure response
//
// Returns:
//   - none
func (r *ProjectUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create adds a user to a project with the specified role.
// Params:
//   - ctx: context
//   - req: create request
//   - resp: create response
//
// Returns:
//   - none
func (r *ProjectUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *models.UserResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	//: Check for plan parsing errors.
	if resp.Diagnostics.HasError() {
		//: Return after plan parsing failure.
		return
	}

	//: Execute create logic and return if it fails.
	if !r.executeCreateLogic(ctx, plan, resp) {
		//: Return after create logic failure.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// executeCreateLogic contains the main logic for adding a user to a project.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - resp: Create response
//
// Returns:
//   - ok: True if creation succeeded, false otherwise
func (r *ProjectUserResource) executeCreateLogic(ctx context.Context, plan *models.UserResource, resp *resource.CreateResponse) (ok bool) {
	//: Build request to add user to project.
	relation := n8nsdk.NewProjectsProjectIdUsersPostRequestRelationsInner(
		plan.UserID.ValueString(),
		plan.Role.ValueString(),
	)
	relations := []n8nsdk.ProjectsProjectIdUsersPostRequestRelationsInner{*relation}
	addUserReq := n8nsdk.NewProjectsProjectIdUsersPostRequest(relations)

	//: Add user to project via API.
	httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsProjectIdUsersPost(ctx, plan.ProjectID.ValueString()).
		ProjectsProjectIdUsersPostRequest(*addUserReq).Execute()
	//: Close HTTP response body if present to prevent resource leaks.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Silently discard close error on response body.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				resp.Diagnostics.AddWarning("Failed to close response body", closeErr.Error())
			}
		}()
	}
	//: Check for API error when adding user to project.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error adding user to project",
			fmt.Sprintf("Could not add user %s to project %s: %s\nHTTP Response: %v",
				plan.UserID.ValueString(), plan.ProjectID.ValueString(), err.Error(), httpResp),
		)
		//: Return failure when the API call fails.
		return false
	}

	//: Set composite ID combining project ID and user ID.
	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.ProjectID.ValueString(), plan.UserID.ValueString()))

	//: Return success after user was added to project.
	return true
}

// Read refreshes the resource state by checking if the user is in the project.
// Params:
//   - ctx: context
//   - req: read request
//   - resp: read response
//
// Returns:
//   - none
func (r *ProjectUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *models.UserResource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	//: Check for state parsing errors.
	if resp.Diagnostics.HasError() {
		//: Return after state parsing failure.
		return
	}

	//: Execute read logic and return if it fails.
	if !r.executeReadLogic(ctx, state, resp) {
		//: Return after read logic failure.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// executeReadLogic contains the main logic for reading a project user membership.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Read response
//
// Returns:
//   - ok: True if read succeeded, false otherwise
func (r *ProjectUserResource) executeReadLogic(ctx context.Context, state *models.UserResource, resp *resource.ReadResponse) (ok bool) {
	found := r.findUserInProject(ctx, state, resp)
	//: Check if user was found in the project.
	if !found {
		//: Return failure when user was not found.
		return false
	}

	//: Return success when user was found in project.
	return true
}

// findUserInProject checks if the user is in the project and updates the state.
//
// Params:
//   - ctx: context
//   - state: current resource state
//   - resp: read response
//
// Returns:
//   - ok: true if user was found in project
func (r *ProjectUserResource) findUserInProject(
	ctx context.Context,
	state *models.UserResource,
	resp *resource.ReadResponse,
) (ok bool) {
	userList, httpResp, err := r.client.APIClient.UserAPI.UsersGet(ctx).
		ProjectId(state.ProjectID.ValueString()).
		IncludeRole(true).
		Execute()
	//: Close HTTP response body if present to prevent resource leaks.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Silently discard close error on response body.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				resp.Diagnostics.AddWarning("Failed to close response body", closeErr.Error())
			}
		}()
	}
	//: Check for API error when listing project users.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project users",
			fmt.Sprintf("Could not read users for project %s: %s\nHTTP Response: %v",
				state.ProjectID.ValueString(), err.Error(), httpResp),
		)
		//: Return false to indicate failure.
		return false
	}

	//: Find the user in the project member list.
	found := r.searchUserInList(userList, state)

	//: Remove resource if user was not found (deleted outside Terraform).
	if !found {
		// User not found in project - resource has been deleted outside Terraform.
		resp.State.RemoveResource(ctx)
		//: Return false to indicate resource was removed.
		return false
	}

	//: Return true to indicate user was found.
	return true
}

// searchUserInList searches for the user in the user list and updates state.
//
// Params:
//   - userList: list of users from API
//   - state: current resource state to update
//
// Returns:
//   - ok: true if user was found in the list
func (r *ProjectUserResource) searchUserInList(
	userList *n8nsdk.UserList,
	state *models.UserResource,
) (ok bool) {
	//: Return false if no data is available in the user list.
	if userList.Data == nil {
		//: Return false if no data available.
		return false
	}

	//: Iterate over all users to find the matching one.
	for _, user := range userList.Data {
		//: Check if current user matches the state user ID.
		if user.Id != nil && *user.Id == state.UserID.ValueString() {
			//: Update role if available.
			if user.Role != nil {
				state.Role = types.StringPointerValue(user.Role)
			}
			//: Return true to indicate user was found.
			return true
		}
	}

	//: Return false if user was not found in the list.
	return false
}

// Update changes the user's role in the project.
// Params:
//   - ctx: context
//   - req: update request
//   - resp: update response
//
// Returns:
//   - none
func (r *ProjectUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *models.UserResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	//: Check for plan/state parsing errors.
	if resp.Diagnostics.HasError() {
		//: Return after plan/state parsing failure.
		return
	}

	//: Execute update logic and return if it fails.
	if !r.executeUpdateLogic(ctx, plan, state, resp) {
		//: Return after update logic failure.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// executeUpdateLogic contains the main logic for updating a project user membership.
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
func (r *ProjectUserResource) executeUpdateLogic(ctx context.Context, plan, state *models.UserResource, resp *resource.UpdateResponse) (ok bool) {
	//: Check if project or user changed - not supported.
	if !plan.ProjectID.Equal(state.ProjectID) || !plan.UserID.Equal(state.UserID) {
		resp.Diagnostics.AddError(
			"Project/User Change Not Supported",
			"Cannot change project_id or user_id. Please delete and recreate the resource.",
		)
		//: Return failure when unsupported change is requested.
		return false
	}

	//: Update role if it has changed.
	if !plan.Role.Equal(state.Role) {
		roleReq := n8nsdk.NewProjectsProjectIdUsersUserIdPatchRequest(plan.Role.ValueString())
		httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsProjectIdUsersUserIdPatch(
			ctx,
			plan.ProjectID.ValueString(),
			plan.UserID.ValueString(),
		).ProjectsProjectIdUsersUserIdPatchRequest(*roleReq).Execute()
		//: Close HTTP response body if present to prevent resource leaks.
		if httpResp != nil && httpResp.Body != nil {
			defer func() {
				//: Silently discard close error on response body.
				if closeErr := httpResp.Body.Close(); closeErr != nil {
					resp.Diagnostics.AddWarning("Failed to close response body", closeErr.Error())
				}
			}()
		}

		//: Check for API error when updating user role.
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating user role in project",
				fmt.Sprintf("Could not update role for user %s in project %s: %s\nHTTP Response: %v",
					plan.UserID.ValueString(), plan.ProjectID.ValueString(), err.Error(), httpResp),
			)
			//: Return failure when role update API call fails.
			return false
		}
	}

	//: Return success after role update or when no change was needed.
	return true
}

// Delete removes the user from the project.
// Params:
//   - ctx: context
//   - req: delete request
//   - resp: delete response
//
// Returns:
//   - none
func (r *ProjectUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *models.UserResource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	//: Check for state parsing errors.
	if resp.Diagnostics.HasError() {
		//: Return after state parsing failure.
		return
	}

	//: Execute delete logic.
	r.executeDeleteLogic(ctx, state, resp)
}

// executeDeleteLogic contains the main logic for removing a user from a project.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Delete response
//
// Returns:
//   - ok: True if delete succeeded, false otherwise
func (r *ProjectUserResource) executeDeleteLogic(ctx context.Context, state *models.UserResource, resp *resource.DeleteResponse) (ok bool) {
	//: Remove user from project via API.
	httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsProjectIdUsersUserIdDelete(
		ctx,
		state.ProjectID.ValueString(),
		state.UserID.ValueString(),
	).Execute()
	//: Close HTTP response body if present to prevent resource leaks.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Silently discard close error on response body.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				resp.Diagnostics.AddWarning("Failed to close response body", closeErr.Error())
			}
		}()
	}

	//: Check for API error when removing user from project.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error removing user from project",
			fmt.Sprintf("Could not remove user %s from project %s: %s\nHTTP Response: %v",
				state.UserID.ValueString(), state.ProjectID.ValueString(), err.Error(), httpResp),
		)
		//: Return failure when the delete API call fails.
		return false
	}

	//: Return success after user was removed from project.
	return true
}

// ImportState imports the resource using the composite ID format: project_id/user_id.
// Params:
//   - ctx: context
//   - req: import state request
//   - resp: import state response
//
// Returns:
//   - none
func (r *ProjectUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//: Split composite ID into project ID and user ID parts.
	parts := strings.Split(req.ID, "/")
	//: Validate that the import ID has exactly two parts.
	if len(parts) != CompositeIDParts {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format 'project_id/user_id', got: %s", req.ID),
		)
		//: Return after reporting invalid import ID format.
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), parts[1])...)
}
