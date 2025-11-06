package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure UserResource implements required interfaces.
var (
	_ resource.Resource                = &UserResource{}
	_ resource.ResourceWithConfigure   = &UserResource{}
	_ resource.ResourceWithImportState = &UserResource{}
)

// UserResource defines the resource implementation for user management.
type UserResource struct {
	client *providertypes.N8nClient
}

// UserResourceModel describes the resource data model.
type UserResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	Role      types.String `tfsdk:"role"`
	IsPending types.Bool   `tfsdk:"is_pending"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// NewUserResource creates a new UserResource instance.
func NewUserResource() resource.Resource {
 // Return result.
	return &UserResource{}
}

// Metadata returns the resource type name.
func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the resource.
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
func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*providertypes.N8nClient)
	// Check condition.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providertypes.N8nClient, got: %T", req.ProviderData),
		)
		// Return result.
		return
	}

	r.client = client
}

// Create creates a new user.
func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
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

	// Populate response data with full user details
	if user.Id != nil {
		plan.ID = types.StringValue(*user.Id)
	}
	plan.Email = types.StringValue(user.Email)
	// Check for non-nil value.
	if user.FirstName != nil {
		plan.FirstName = types.StringPointerValue(user.FirstName)
	}
	// Check for non-nil value.
	if user.LastName != nil {
		plan.LastName = types.StringPointerValue(user.LastName)
	}
	// Check for non-nil value.
	if user.Role != nil {
		plan.Role = types.StringPointerValue(user.Role)
	}
	// Check for non-nil value.
	if user.IsPending != nil {
		plan.IsPending = types.BoolPointerValue(user.IsPending)
	}
	// Check for non-nil value.
	if user.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(user.CreatedAt.String())
	}
	// Check for non-nil value.
	if user.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(user.UpdatedAt.String())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the resource state with the latest data.
func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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

	// Update state with latest data
	if user.Id != nil {
		state.ID = types.StringValue(*user.Id)
	}
	state.Email = types.StringValue(user.Email)
	// Check for non-nil value.
	if user.FirstName != nil {
		state.FirstName = types.StringPointerValue(user.FirstName)
	}
	// Check for non-nil value.
	if user.LastName != nil {
		state.LastName = types.StringPointerValue(user.LastName)
	}
	// Check for non-nil value.
	if user.Role != nil {
		state.Role = types.StringPointerValue(user.Role)
	}
	// Check for non-nil value.
	if user.IsPending != nil {
		state.IsPending = types.BoolPointerValue(user.IsPending)
	}
	// Check for non-nil value.
	if user.CreatedAt != nil {
		state.CreatedAt = types.StringValue(user.CreatedAt.String())
	}
	// Check for non-nil value.
	if user.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(user.UpdatedAt.String())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the user. Only role changes are supported by the API.
func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state UserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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
			// Return result.
			return
		}
	}

	// Refresh user data after update
	user, httpResp, err := r.client.APIClient.UserAPI.UsersIdGet(ctx, state.ID.ValueString()).IncludeRole(true).Execute()
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

	// Update state with latest data
	if user.Id != nil {
		plan.ID = types.StringValue(*user.Id)
	}
	plan.Email = types.StringValue(user.Email)
	// Check for non-nil value.
	if user.FirstName != nil {
		plan.FirstName = types.StringPointerValue(user.FirstName)
	}
	// Check for non-nil value.
	if user.LastName != nil {
		plan.LastName = types.StringPointerValue(user.LastName)
	}
	// Check for non-nil value.
	if user.Role != nil {
		plan.Role = types.StringPointerValue(user.Role)
	}
	// Check for non-nil value.
	if user.IsPending != nil {
		plan.IsPending = types.BoolPointerValue(user.IsPending)
	}
	// Check for non-nil value.
	if user.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(user.CreatedAt.String())
	}
	// Check for non-nil value.
	if user.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(user.UpdatedAt.String())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the user.
func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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
func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
