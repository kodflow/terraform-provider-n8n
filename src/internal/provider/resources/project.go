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

// Ensure ProjectResource implements required interfaces.
var (
	_ resource.Resource                = &ProjectResource{}
	_ resource.ResourceWithConfigure   = &ProjectResource{}
	_ resource.ResourceWithImportState = &ProjectResource{}
)

// ProjectResource defines the resource implementation for n8n projects.
// Note: n8n API has limitations - POST returns 201 with no body, no GET by ID endpoint.
// We work around this by using the LIST endpoint and filtering.
type ProjectResource struct {
	client *providertypes.N8nClient
}

// ProjectResourceModel describes the resource data model.
type ProjectResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

// NewProjectResource creates a new ProjectResource instance.
func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

// Metadata returns the resource type name.
func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the resource.
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
func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*providertypes.N8nClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providertypes.N8nClient, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
// Workaround: API returns 201 with no body, so we must call LIST to get the ID.
func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProjectResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectRequest := n8nsdk.Project{
		Name: plan.Name.ValueString(),
	}

	// POST returns 201 with no body
	httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsPost(ctx).
		Project(projectRequest).
		Execute()
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			fmt.Sprintf("Could not create project: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Workaround: List all projects to find the one we just created
	projectList, httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsGet(ctx).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project after creation",
			fmt.Sprintf("Project was created but could not retrieve ID: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Find our project by name
	var foundProject *n8nsdk.Project
	if projectList.Data != nil {
		for _, p := range projectList.Data {
			if p.Name == plan.Name.ValueString() {
				foundProject = &p
				break
			}
		}
	}

	if foundProject == nil {
		resp.Diagnostics.AddError(
			"Error finding created project",
			fmt.Sprintf("Project with name '%s' was created but not found in list", plan.Name.ValueString()),
		)
		return
	}

	plan.ID = types.StringPointerValue(foundProject.Id)
	plan.Name = types.StringValue(foundProject.Name)
	if foundProject.Type != nil {
		plan.Type = types.StringPointerValue(foundProject.Type)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
// Workaround: No GET by ID endpoint, so we use LIST and filter.
func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ProjectResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Workaround: No GET by ID, use LIST and filter
	projectList, httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsGet(ctx).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project",
			fmt.Sprintf("Could not read project ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		return
	}

	// Find our project by ID
	var foundProject *n8nsdk.Project
	if projectList.Data != nil {
		for _, p := range projectList.Data {
			if p.Id != nil && *p.Id == state.ID.ValueString() {
				foundProject = &p
				break
			}
		}
	}

	if foundProject == nil {
		// Project not found = deleted outside Terraform
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(foundProject.Name)
	if foundProject.Type != nil {
		state.Type = types.StringPointerValue(foundProject.Type)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
// Workaround: API returns 204 with no body, so we must call LIST to verify.
func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProjectResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectRequest := n8nsdk.Project{
		Name: plan.Name.ValueString(),
	}

	// PUT returns 204 with no body
	httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsProjectIdPut(ctx, plan.ID.ValueString()).
		Project(projectRequest).
		Execute()
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating project",
			fmt.Sprintf("Could not update project ID %s: %s\nHTTP Response: %v", plan.ID.ValueString(), err.Error(), httpResp),
		)
		return
	}

	// Workaround: List all projects to verify the update
	projectList, httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsGet(ctx).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project after update",
			fmt.Sprintf("Project was updated but could not verify: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Find our project by ID
	var foundProject *n8nsdk.Project
	if projectList.Data != nil {
		for _, p := range projectList.Data {
			if p.Id != nil && *p.Id == plan.ID.ValueString() {
				foundProject = &p
				break
			}
		}
	}

	if foundProject == nil {
		resp.Diagnostics.AddError(
			"Error verifying updated project",
			fmt.Sprintf("Project with ID '%s' was updated but not found in list", plan.ID.ValueString()),
		)
		return
	}

	plan.Name = types.StringValue(foundProject.Name)
	if foundProject.Type != nil {
		plan.Type = types.StringPointerValue(foundProject.Type)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ProjectResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// DELETE returns 204 with no body
	httpResp, err := r.client.APIClient.ProjectsAPI.ProjectsProjectIdDelete(ctx, state.ID.ValueString()).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting project",
			fmt.Sprintf("Could not delete project ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		return
	}
}

// ImportState imports the resource into Terraform state.
func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
