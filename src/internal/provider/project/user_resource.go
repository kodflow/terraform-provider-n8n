package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)

// CompositeIDParts is the expected number of parts in a composite ID (project_id/user_id).
const CompositeIDParts int = 2

// Ensure ProjectUserResource implements required interfaces.
var (
	_ resource.Resource                = &ProjectUserResource{}
	_ resource.ResourceWithConfigure   = &ProjectUserResource{}
	_ resource.ResourceWithImportState = &ProjectUserResource{}
)

// ProjectUserResource defines the resource implementation for project-user relationships.
// Terraform resource that manages CRUD operations for user memberships in n8n projects via the n8n API.
type ProjectUserResource struct {
	client *client.N8nClient
}


// NewProjectUserResource creates a new ProjectUserResource instance.
// Returns:
//   - resource.Resource: new ProjectUserResource instance
func NewProjectUserResource() resource.Resource {
	// Return result.
	return &ProjectUserResource{}
}

// Metadata returns the resource type name.
// Params:
//   - ctx: context
//   - req: metadata request
//   - resp: metadata response
// Returns:
//   - none
func (r *ProjectUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_user"
}

// Schema defines the schema for the resource.
// Params:
//   - ctx: context
//   - req: schema request
//   - resp: schema response
// Returns:
//   - none
func (r *ProjectUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
// Returns:
//   - none
func (r *ProjectUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create adds a user to a project with the specified role.
// Params:
//   - ctx: context
//   - req: create request
//   - resp: create response
// Returns:
//   - none
func (r *ProjectUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *ProjectUserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request to add user to project.
	relation := n8nsdk.NewProjectsProjectIdUsersPostRequestRelationsInner(
		plan.UserID.ValueString(),
		plan.Role.ValueString(),
	)
	relations := []n8nsdk.ProjectsProjectIdUsersPostRequestRelationsInner{*relation}
	addUserReq := n8nsdk.NewProjectsProjectIdUsersPostRequest(relations)

	// Add user to project.
	httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsProjectIdUsersPost(ctx, plan.ProjectID.ValueString()).
		ProjectsProjectIdUsersPostRequest(*addUserReq).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error adding user to project",
			fmt.Sprintf("Could not add user %s to project %s: %s\nHTTP Response: %v",
				plan.UserID.ValueString(), plan.ProjectID.ValueString(), err.Error(), httpResp),
		)
		// Return result.
		return
	}

	// Set composite ID.
	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.ProjectID.ValueString(), plan.UserID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the resource state by checking if the user is in the project.
// Params:
//   - ctx: context
//   - req: read request
//   - resp: read response
// Returns:
//   - none
func (r *ProjectUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *ProjectUserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Use users API with project filter to check if user is in project.
	userList, httpResp, err := r.client.APIClient.UserAPI.UsersGet(ctx).
		ProjectId(state.ProjectID.ValueString()).
		IncludeRole(true).
		Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project users",
			fmt.Sprintf("Could not read users for project %s: %s\nHTTP Response: %v",
				state.ProjectID.ValueString(), err.Error(), httpResp),
		)
		// Return result.
		return
	}

	// Find the user in the project.
	found := false
	// Check for non-nil value.
	if userList.Data != nil {
		// Iterate over items.
		for _, user := range userList.Data {
			// Check for non-nil value.
			if user.Id != nil && *user.Id == state.UserID.ValueString() {
				found = true
				// Update role if available.
				if user.Role != nil {
					state.Role = types.StringPointerValue(user.Role)
				}
				break
			}
		}
	}

	// Check condition.
	if !found {
		// User not found in project - resource has been deleted outside Terraform.
		resp.State.RemoveResource(ctx)
		// Return result.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update changes the user's role in the project.
// Params:
//   - ctx: context
//   - req: update request
//   - resp: update response
// Returns:
//   - none
func (r *ProjectUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *ProjectUserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if project or user changed - not supported.
	if !plan.ProjectID.Equal(state.ProjectID) || !plan.UserID.Equal(state.UserID) {
		resp.Diagnostics.AddError(
			"Project/User Change Not Supported",
			"Cannot change project_id or user_id. Please delete and recreate the resource.",
		)
		// Return result.
		return
	}

	// Update role if changed.
	if !plan.Role.Equal(state.Role) {
		roleReq := n8nsdk.NewProjectsProjectIdUsersUserIdPatchRequest(plan.Role.ValueString())
		httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsProjectIdUsersUserIdPatch(
			ctx,
			plan.ProjectID.ValueString(),
			plan.UserID.ValueString(),
		).ProjectsProjectIdUsersUserIdPatchRequest(*roleReq).Execute()
		// Check for non-nil value.
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

		// Check for error.
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating user role in project",
				fmt.Sprintf("Could not update role for user %s in project %s: %s\nHTTP Response: %v",
					plan.UserID.ValueString(), plan.ProjectID.ValueString(), err.Error(), httpResp),
			)
			// Return result.
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete removes the user from the project.
// Params:
//   - ctx: context
//   - req: delete request
//   - resp: delete response
// Returns:
//   - none
func (r *ProjectUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *ProjectUserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove user from project.
	httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsProjectIdUsersUserIdDelete(
		ctx,
		state.ProjectID.ValueString(),
		state.UserID.ValueString(),
	).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error removing user from project",
			fmt.Sprintf("Could not remove user %s from project %s: %s\nHTTP Response: %v",
				state.UserID.ValueString(), state.ProjectID.ValueString(), err.Error(), httpResp),
		)
		// Return result.
		return
	}
}

// ImportState imports the resource using the composite ID format: project_id/user_id.
// Params:
//   - ctx: context
//   - req: import state request
//   - resp: import state response
// Returns:
//   - none
func (r *ProjectUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Split composite ID.
	parts := strings.Split(req.ID, "/")
	// Check condition.
	if len(parts) != CompositeIDParts {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format 'project_id/user_id', got: %s", req.ID),
		)
		// Return result.
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), parts[1])...)
}
