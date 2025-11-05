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

// Ensure WorkflowResource implements required interfaces
var (
	_ resource.Resource                = &WorkflowResource{}
	_ resource.ResourceWithConfigure   = &WorkflowResource{}
	_ resource.ResourceWithImportState = &WorkflowResource{}
)

// WorkflowResource defines the resource implementation for n8n workflows.
type WorkflowResource struct {
	client *providertypes.N8nClient
}

// WorkflowResourceModel describes the resource data model.
type WorkflowResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Active       types.Bool   `tfsdk:"active"`
	Tags         types.List   `tfsdk:"tags"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
	VersionId    types.String `tfsdk:"version_id"`
	IsArchived   types.Bool   `tfsdk:"is_archived"`
	TriggerCount types.Int64  `tfsdk:"trigger_count"`
	Meta         types.Map    `tfsdk:"meta"`
	PinData      types.Map    `tfsdk:"pin_data"`
}

// NewWorkflowResource creates a new WorkflowResource instance.
func NewWorkflowResource() resource.Resource {
	return &WorkflowResource{}
}

// Metadata returns the resource type name.
func (r *WorkflowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

// Schema defines the schema for the resource.
func (r *WorkflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "n8n workflow resource using generated SDK",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Workflow identifier",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Workflow name",
				Required:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the workflow is active",
				Optional:            true,
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "List of tag IDs associated with this workflow",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the workflow was created",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the workflow was last updated",
				Computed:            true,
			},
			"version_id": schema.StringAttribute{
				MarkdownDescription: "Version identifier of the workflow",
				Computed:            true,
			},
			"is_archived": schema.BoolAttribute{
				MarkdownDescription: "Whether the workflow is archived",
				Computed:            true,
			},
			"trigger_count": schema.Int64Attribute{
				MarkdownDescription: "Number of triggers in the workflow",
				Computed:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "Workflow metadata",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"pin_data": schema.MapAttribute{
				MarkdownDescription: "Pinned test data for the workflow",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *WorkflowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*providertypes.N8nClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providertypes.N8nClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *WorkflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WorkflowResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create workflow using SDK
	// Note: 'active' and 'tags' fields are read-only during creation
	// Note: 'nodes' and 'connections' are required by the API
	workflowRequest := n8nsdk.Workflow{
		Name:        plan.Name.ValueString(),
		Nodes:       []n8nsdk.Node{},            // Empty nodes array required by API
		Connections: map[string]interface{}{}, // Empty connections object required by API
	}

	// Call SDK API to create workflow
	workflow, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsPost(ctx).
		Workflow(workflowRequest).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating workflow",
			fmt.Sprintf("Could not create workflow, unexpected error: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// If tags were provided, update the workflow with tags using the dedicated tags endpoint
	// (tags field is read-only in POST/PUT /workflows but can be set via PUT /workflows/{id}/tags)
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() && workflow.Id != nil {
		var tagIDs []string
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tagIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		tagIdsInner := make([]n8nsdk.TagIdsInner, len(tagIDs))
		for i, tagID := range tagIDs {
			tagIdsInner[i] = n8nsdk.TagIdsInner{Id: tagID}
		}

		tags, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdTagsPut(ctx, *workflow.Id).
			TagIdsInner(tagIdsInner).
			Execute()

		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating workflow with tags",
				fmt.Sprintf("Workflow created but could not set tags: %s\nHTTP Response: %v", err.Error(), httpResp),
			)
			return
		}

		// Update workflow tags from response
		workflow.Tags = tags
	}

	// Map response to state
	plan.ID = types.StringPointerValue(workflow.Id)
	plan.Name = types.StringValue(workflow.Name)
	if workflow.Active != nil {
		plan.Active = types.BoolPointerValue(workflow.Active)
	}

	// Map tags from response
	if workflow.Tags != nil && len(workflow.Tags) > 0 {
		tagIDs := make([]types.String, len(workflow.Tags))
		for i, tag := range workflow.Tags {
			if tag.Id != nil {
				tagIDs[i] = types.StringValue(*tag.Id)
			}
		}
		tagList, diags := types.ListValueFrom(ctx, types.StringType, tagIDs)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			plan.Tags = tagList
		}
	}

	// Map computed fields from response
	if workflow.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(workflow.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	if workflow.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(workflow.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	if workflow.VersionId != nil {
		plan.VersionId = types.StringPointerValue(workflow.VersionId)
	}
	if workflow.IsArchived != nil {
		plan.IsArchived = types.BoolPointerValue(workflow.IsArchived)
	}
	if workflow.TriggerCount != nil {
		plan.TriggerCount = types.Int64Value(int64(*workflow.TriggerCount))
	}
	if workflow.Meta != nil {
		metaMap, diags := types.MapValueFrom(ctx, types.StringType, workflow.Meta)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			plan.Meta = metaMap
		}
	}
	if workflow.PinData != nil {
		pinDataMap, diags := types.MapValueFrom(ctx, types.StringType, workflow.PinData)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			plan.PinData = pinDataMap
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *WorkflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkflowResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get workflow from SDK
	workflow, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdGet(ctx, state.ID.ValueString()).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading workflow",
			fmt.Sprintf("Could not read workflow ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		return
	}

	// Update state
	state.Name = types.StringValue(workflow.Name)
	if workflow.Active != nil {
		state.Active = types.BoolPointerValue(workflow.Active)
	}

	// Map tags from response
	if workflow.Tags != nil && len(workflow.Tags) > 0 {
		tagIDs := make([]types.String, len(workflow.Tags))
		for i, tag := range workflow.Tags {
			if tag.Id != nil {
				tagIDs[i] = types.StringValue(*tag.Id)
			}
		}
		tagList, diags := types.ListValueFrom(ctx, types.StringType, tagIDs)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			state.Tags = tagList
		}
	} else {
		// Set to empty list if no tags
		emptyList, diags := types.ListValueFrom(ctx, types.StringType, []types.String{})
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			state.Tags = emptyList
		}
	}

	// Map computed fields from response
	if workflow.CreatedAt != nil {
		state.CreatedAt = types.StringValue(workflow.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	if workflow.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(workflow.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	if workflow.VersionId != nil {
		state.VersionId = types.StringPointerValue(workflow.VersionId)
	}
	if workflow.IsArchived != nil {
		state.IsArchived = types.BoolPointerValue(workflow.IsArchived)
	}
	if workflow.TriggerCount != nil {
		state.TriggerCount = types.Int64Value(int64(*workflow.TriggerCount))
	}
	if workflow.Meta != nil {
		metaMap, diags := types.MapValueFrom(ctx, types.StringType, workflow.Meta)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			state.Meta = metaMap
		}
	}
	if workflow.PinData != nil {
		pinDataMap, diags := types.MapValueFrom(ctx, types.StringType, workflow.PinData)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			state.PinData = pinDataMap
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *WorkflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WorkflowResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update workflow using SDK
	workflowRequest := n8nsdk.Workflow{
		Name: plan.Name.ValueString(),
	}

	if !plan.Active.IsNull() {
		active := plan.Active.ValueBool()
		workflowRequest.Active = &active
	}

	// Call SDK API to update workflow (name, active, etc.)
	workflow, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdPut(ctx, plan.ID.ValueString()).
		Workflow(workflowRequest).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating workflow",
			fmt.Sprintf("Could not update workflow ID %s: %s\nHTTP Response: %v", plan.ID.ValueString(), err.Error(), httpResp),
		)
		return
	}

	// Update tags separately using the dedicated tags endpoint
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		var tagIDs []string
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tagIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		tagIdsInner := make([]n8nsdk.TagIdsInner, len(tagIDs))
		for i, tagID := range tagIDs {
			tagIdsInner[i] = n8nsdk.TagIdsInner{Id: tagID}
		}

		tags, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdTagsPut(ctx, plan.ID.ValueString()).
			TagIdsInner(tagIdsInner).
			Execute()

		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating workflow tags",
				fmt.Sprintf("Could not update tags for workflow ID %s: %s\nHTTP Response: %v", plan.ID.ValueString(), err.Error(), httpResp),
			)
			return
		}

		// Update workflow tags from response
		workflow.Tags = tags
	}

	// Update state with response
	plan.Name = types.StringValue(workflow.Name)
	if workflow.Active != nil {
		plan.Active = types.BoolPointerValue(workflow.Active)
	}

	// Map tags from response
	if workflow.Tags != nil && len(workflow.Tags) > 0 {
		tagIDs := make([]types.String, len(workflow.Tags))
		for i, tag := range workflow.Tags {
			if tag.Id != nil {
				tagIDs[i] = types.StringValue(*tag.Id)
			}
		}
		tagList, diags := types.ListValueFrom(ctx, types.StringType, tagIDs)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			plan.Tags = tagList
		}
	} else {
		// Set to empty list if no tags
		emptyList, diags := types.ListValueFrom(ctx, types.StringType, []types.String{})
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			plan.Tags = emptyList
		}
	}

	// Map computed fields from response
	if workflow.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(workflow.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	if workflow.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(workflow.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	if workflow.VersionId != nil {
		plan.VersionId = types.StringPointerValue(workflow.VersionId)
	}
	if workflow.IsArchived != nil {
		plan.IsArchived = types.BoolPointerValue(workflow.IsArchived)
	}
	if workflow.TriggerCount != nil {
		plan.TriggerCount = types.Int64Value(int64(*workflow.TriggerCount))
	}
	if workflow.Meta != nil {
		metaMap, diags := types.MapValueFrom(ctx, types.StringType, workflow.Meta)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			plan.Meta = metaMap
		}
	}
	if workflow.PinData != nil {
		pinDataMap, diags := types.MapValueFrom(ctx, types.StringType, workflow.PinData)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			plan.PinData = pinDataMap
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *WorkflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkflowResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete workflow using SDK
	_, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdDelete(ctx, state.ID.ValueString()).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting workflow",
			fmt.Sprintf("Could not delete workflow ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		return
	}
}

// ImportState imports the resource into Terraform state.
func (r *WorkflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
