package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)

// Ensure UserResource implements required interfaces.
var (
	_ resource.Resource                = &UserResource{}
	_ resource.ResourceWithConfigure   = &UserResource{}
	_ resource.ResourceWithImportState = &UserResource{}
)

// UserResource defines the resource implementation for user management.
// Terraform resource that manages CRUD operations for n8n users via the n8n API.
type UserResource struct {
	client *client.N8nClient
}


// NewUserResource creates a new UserResource instance.
//
// Returns:
//   - resource.Resource: A new UserResource configured for Terraform
func NewUserResource() resource.Resource {
	// Return result.
	return &UserResource{}
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
func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
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
func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages n8n users. Only available for the instance owner. Note: The API only supports updating the user's role, not other fields.",

		Attributes: map[string]schema.Attribute{
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
func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
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
	plan := &UserResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

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
		// Return result.
		return
	}

	// API returns UsersPost200Response with limited user data
	if result.User == nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			"API did not return user data",
		)
		// Return result.
		return
	}

	// Get user ID from creation response
	if result.User.Id == nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			"API did not return user ID",
		)
		// Return result.
		return
	}

	userID := *result.User.Id

	// Fetch full user details using the ID
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
		// Return result.
		return
	}

	// Map user to model
	mapUserToResourceModel(user, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
	state := &UserResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

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
		// Return result.
		return
	}

	// Map user to model
	mapUserToResourceModel(user, state)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
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
	plan := &UserResourceModel{}
	state := &UserResourceModel{}
	var err error
	var httpResp *http.Response

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if email changed - not supported
	if !plan.Email.Equal(state.Email) {
		resp.Diagnostics.AddError(
			"Email Change Not Supported",
			"The n8n API does not support changing a user's email address. Please delete and recreate the user if needed.",
		)
		// Return result.
		return
	}

	// Check if role changed
	roleChanged := !plan.Role.IsNull() && !state.Role.IsNull() &&
		!plan.Role.Equal(state.Role)

	// Check condition.
	if roleChanged {
		// Update role
		roleReq := n8nsdk.NewUsersIdRolePatchRequest(plan.Role.ValueString())
		httpResp, err = r.client.APIClient.UserAPI.UsersIdRolePatch(ctx, state.ID.ValueString()).
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
			// Return result.
			return
		}
	}

	// Refresh user data after update
	var user *n8nsdk.User
	user, httpResp, err = r.client.APIClient.UserAPI.UsersIdGet(ctx, state.ID.ValueString()).IncludeRole(true).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user after update",
			fmt.Sprintf("Could not read user %s after update: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		// Return result.
		return
	}

	// Map user to model
	mapUserToResourceModel(user, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
	state := &UserResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

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
		// Return result.
		return
	}
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
