// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package user implements user management resources and data sources.
package user

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/user/models"
)

// Ensure UserResource implements required interfaces.
var (
	_ resource.Resource                = &UserResource{}
	_ UserResourceInterface            = &UserResource{}
	_ resource.ResourceWithConfigure   = &UserResource{}
	_ resource.ResourceWithImportState = &UserResource{}
)

// UserResourceInterface defines the interface for UserResource.
type UserResourceInterface interface {
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

// UserResource defines the resource implementation for user management.
// Terraform resource that manages CRUD operations for n8n users via the n8n API.
type UserResource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewUserResource creates a new UserResource instance.
//
// Returns:
//   - resource.Resource: A new UserResource configured for Terraform
func NewUserResource() *UserResource {
	// Return result.
	return &UserResource{}
}

// NewUserResourceWrapper creates a new UserResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - resource.Resource: the wrapped UserResource instance
func NewUserResourceWrapper() resource.Resource {
	// Return the wrapped resource instance.
	return NewUserResource()
}

// Metadata returns the resource type name.
//
// Params:
//   - ctx: Context for operation cancellation and deadlines
//   - req: Metadata request from Terraform
//   - resp: Metadata response to populate with resource type name
//
// Returns:
//   - (none)
func (r *UserResource) Metadata(_ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// Set resource type name.
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the resource.
//
// Params:
//   - ctx: Context for operation cancellation and deadlines
//   - req: Schema request from Terraform
//   - resp: Schema response to populate with resource schema definition
//
// Returns:
//   - (none)
func (r *UserResource) Schema(_ctx context.Context, _req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages n8n users. Only available for the instance owner. Note: The API only supports updating the user's role, not other fields.",
		Attributes:          r.schemaAttributes(),
	}
}

// schemaAttributes returns the attribute definitions for the user resource schema.
//
// Returns:
//   - map[string]schema.Attribute: the resource attribute definitions
func (r *UserResource) schemaAttributes() map[string]schema.Attribute {
	// Return schema attributes.
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "User identifier",
			Computed:            true,
		},
		"email": schema.StringAttribute{
			MarkdownDescription: "User's email address",
			Required:            true,
		},
		"first_name": schema.StringAttribute{
			MarkdownDescription: "User's first name",
			Computed:            true,
		},
		"last_name": schema.StringAttribute{
			MarkdownDescription: "User's last name",
			Computed:            true,
		},
		"role": schema.StringAttribute{
			MarkdownDescription: "User's global role (e.g., 'global:admin', 'global:member')",
			Optional:            true,
			Computed:            true,
		},
		"is_pending": schema.BoolAttribute{
			MarkdownDescription: "Whether the user has finished setting up their account",
			Computed:            true,
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "Timestamp when the user was created",
			Computed:            true,
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "Timestamp when the user was last updated",
			Computed:            true,
		},
	}
}

// Configure adds the provider configured client to the resource.
//
// Params:
//   - ctx: Context for operation cancellation and deadlines
//   - req: Configure request from Terraform containing provider data
//   - resp: Configure response to populate with diagnostics
//
// Returns:
//   - (none)
func (r *UserResource) Configure(_ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new user.
//
// Params:
//   - ctx: Context for operation cancellation and deadlines
//   - req: Create request from Terraform containing planned resource state
//   - resp: Create response to populate with created resource state
//
// Returns:
//   - (none)
func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &models.Resource{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute create logic
	if !r.executeCreateLogic(ctx, plan, resp) {
		// Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// executeCreateLogic contains the main logic for creating a user.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - resp: Create response
//
// Returns:
//   - bool: True if creation succeeded, false otherwise
func (r *UserResource) executeCreateLogic(ctx context.Context, plan *models.Resource, resp *resource.CreateResponse) bool {
	// Create user and get ID
	userID := r.createUser(ctx, plan, resp)
	// Check condition.
	if userID == "" {
		// Return failure.
		return false
	}

	// Fetch full user details
	user := r.fetchFullUserDetails(ctx, userID, resp)
	// Check for nil value.
	if user == nil {
		// Return failure.
		return false
	}

	// Map user to model
	mapUserToResourceModel(user, plan)

	// Return success.
	return true
}

// createUser creates a user and returns the user ID.
//
// Params:
//   - ctx: context for operation
//   - plan: user resource model
//   - resp: create response
//
// Returns:
//   - string: user ID or empty string if error occurred
func (r *UserResource) createUser(ctx context.Context, plan *models.Resource, resp *resource.CreateResponse) string {
	// Build create request (API accepts array of users, we're creating one)
	userReq := n8nsdk.NewUsersPostRequestInner(plan.Email.ValueString())
	// Check condition.
	if !plan.Role.IsNull() && !plan.Role.IsUnknown() {
		role := plan.Role.ValueString()
		userReq.SetRole(role)
	}

	usersArray := []n8nsdk.UsersPostRequestInner{*userReq}

	// Create user
	result, httpResp, err := r.client.APIClient.UserAPI.UsersPost(ctx).UsersPostRequestInner(usersArray).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			fmt.Sprintf("Could not create user %s: %s\nHTTP Response: %v", plan.Email.ValueString(), err.Error(), httpResp),
		)
		// Return empty string on error.
		return ""
	}

	// Validate response
	if result.User == nil || result.User.Id == nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			"API did not return user ID",
		)
		// Return empty string on validation error.
		return ""
	}

	// Return the user ID.
	return *result.User.Id
}

// fetchFullUserDetails retrieves full user details by ID.
//
// Params:
//   - ctx: context for operation
//   - userID: user identifier
//   - resp: create response
//
// Returns:
//   - *n8nsdk.User: user object or nil if error occurred
func (r *UserResource) fetchFullUserDetails(ctx context.Context, userID string, resp *resource.CreateResponse) *n8nsdk.User {
	user, httpResp, err := r.client.APIClient.UserAPI.UsersIdGet(ctx, userID).IncludeRole(true).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading created user",
			fmt.Sprintf("User was created but could not read full details: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return nil on error.
		return nil
	}
	// Return the user details.
	return user
}

// Read refreshes the resource state with the latest data.
//
// Params:
//   - ctx: Context for operation cancellation and deadlines
//   - req: Read request from Terraform containing current resource state
//   - resp: Read response to populate with refreshed resource state
//
// Returns:
//   - (none)
func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &models.Resource{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute read logic
	if !r.executeReadLogic(ctx, state, resp) {
		// Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// executeReadLogic contains the main logic for reading a user.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Read response
//
// Returns:
//   - bool: True if read succeeded, false otherwise
func (r *UserResource) executeReadLogic(ctx context.Context, state *models.Resource, resp *resource.ReadResponse) bool {
	// Fetch user by ID
	user, httpResp, err := r.client.APIClient.UserAPI.UsersIdGet(ctx, state.ID.ValueString()).IncludeRole(true).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user",
			fmt.Sprintf("Could not read user %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		// Return failure.
		return false
	}

	// Map user to model
	mapUserToResourceModel(user, state)

	// Return success.
	return true
}

// Update updates the user. Only role changes are supported by the API.
//
// Params:
//   - ctx: Context for operation cancellation and deadlines
//   - req: Update request from Terraform containing planned and current resource state
//   - resp: Update response to populate with updated resource state
//
// Returns:
//   - (none)
func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &models.Resource{}
	state := &models.Resource{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute update logic
	if !r.executeUpdateLogic(ctx, plan, state, resp) {
		// Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// executeUpdateLogic contains the main logic for updating a user.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - state: The current resource state
//   - resp: Update response
//
// Returns:
//   - bool: True if update succeeded, false otherwise
func (r *UserResource) executeUpdateLogic(ctx context.Context, plan, state *models.Resource, resp *resource.UpdateResponse) bool {
	// Validate email hasn't changed
	if !r.validateEmailUnchanged(plan, state, resp) {
		// Return failure.
		return false
	}

	// Update role if changed
	r.updateRoleIfChanged(ctx, plan, state, resp)
	// Check for errors in diagnostics.
	if resp.Diagnostics.HasError() {
		// Return failure.
		return false
	}

	// Refresh user data
	user := r.refreshUserData(ctx, state.ID.ValueString(), resp)
	// Check for nil value.
	if user == nil {
		// Return failure.
		return false
	}

	// Map user to model
	mapUserToResourceModel(user, plan)

	// Return success.
	return true
}

// validateEmailUnchanged checks if email has changed and reports error.
//
// Params:
//   - plan: planned state
//   - state: current state
//   - resp: update response
//
// Returns:
//   - bool: true if email unchanged, false otherwise
func (r *UserResource) validateEmailUnchanged(plan, state *models.Resource, resp *resource.UpdateResponse) bool {
	// Check condition.
	if !plan.Email.Equal(state.Email) {
		resp.Diagnostics.AddError(
			"Email Change Not Supported",
			"The n8n API does not support changing a user's email address. Please delete and recreate the user if needed.",
		)
		return false
	}
	return true
}

// updateRoleIfChanged updates the user role if it has changed.
//
// Params:
//   - ctx: context for operation
//   - plan: planned state
//   - state: current state
//   - resp: update response
func (r *UserResource) updateRoleIfChanged(ctx context.Context, plan, state *models.Resource, resp *resource.UpdateResponse) {
	// Check if role changed
	roleChanged := !plan.Role.IsNull() && !state.Role.IsNull() &&
		!plan.Role.Equal(state.Role)

	// Check condition.
	if !roleChanged {
		// Return result.
		return
	}

	// Update role
	roleReq := n8nsdk.NewUsersIdRolePatchRequest(plan.Role.ValueString())
	httpResp, err := r.client.APIClient.UserAPI.UsersIdRolePatch(ctx, state.ID.ValueString()).
		UsersIdRolePatchRequest(*roleReq).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating user role",
			fmt.Sprintf("Could not update role for user %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
	}
}

// refreshUserData retrieves updated user data from the API.
//
// Params:
//   - ctx: context for operation
//   - userID: user identifier
//   - resp: update response
//
// Returns:
//   - *n8nsdk.User: user object or nil if error occurred
func (r *UserResource) refreshUserData(ctx context.Context, userID string, resp *resource.UpdateResponse) *n8nsdk.User {
	user, httpResp, err := r.client.APIClient.UserAPI.UsersIdGet(ctx, userID).IncludeRole(true).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user after update",
			fmt.Sprintf("Could not read user %s after update: %s\nHTTP Response: %v", userID, err.Error(), httpResp),
		)
		// Return nil on error.
		return nil
	}
	// Return the user details.
	return user
}

// Delete deletes the user.
//
// Params:
//   - ctx: Context for operation cancellation and deadlines
//   - req: Delete request from Terraform containing resource state to delete
//   - resp: Delete response to populate with diagnostics
//
// Returns:
//   - (none)
func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &models.Resource{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute delete logic
	r.executeDeleteLogic(ctx, state, resp)
}

// executeDeleteLogic contains the main logic for deleting a user.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Delete response
//
// Returns:
//   - bool: True if delete succeeded, false otherwise
func (r *UserResource) executeDeleteLogic(ctx context.Context, state *models.Resource, resp *resource.DeleteResponse) bool {
	// Delete user
	httpResp, err := r.client.APIClient.UserAPI.UsersIdDelete(ctx, state.ID.ValueString()).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting user",
			fmt.Sprintf("Could not delete user %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		// Return failure.
		return false
	}

	// Return success.
	return true
}

// ImportState imports the resource into Terraform state.
//
// Params:
//   - ctx: Context for operation cancellation and deadlines
//   - req: Import state request from Terraform containing resource ID to import
//   - resp: Import state response to populate with imported resource state
//
// Returns:
//   - (none)
func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
