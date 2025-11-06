package workflow

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)

// Ensure WorkflowResource implements required interfaces.
var (
	_ resource.Resource                = &WorkflowResource{}
	_ resource.ResourceWithConfigure   = &WorkflowResource{}
	_ resource.ResourceWithImportState = &WorkflowResource{}
)

// WorkflowResource defines the resource implementation for n8n workflows.
// Terraform resource that manages CRUD operations for n8n workflows via the n8n API.
//
// Params:
//   - client: The n8n API client for communicating with n8n server
//
// Returns:
//   - None: This is a struct, not a function
type WorkflowResource struct {
	client *client.N8nClient
}

// NewWorkflowResource creates a new WorkflowResource instance.
//
// Returns:
//   - resource.Resource: A new WorkflowResource instance
func NewWorkflowResource() resource.Resource {
	// Return result.
	return &WorkflowResource{}
}

// Metadata returns the resource type name.
//
// Params:
//   - ctx: Context for the operation
//   - req: Metadata request containing provider type name
//   - resp: Metadata response to populate with type name
//
// Returns:
//   - None: Updates resp in-place
func (r *WorkflowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

// Schema defines the schema for the resource.
//
// Params:
//   - ctx: Context for the operation
//   - req: Schema request
//   - resp: Schema response to populate with resource schema
//
// Returns:
//   - None: Updates resp in-place
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
			"nodes_json": schema.StringAttribute{
				MarkdownDescription: "Workflow nodes as JSON string. Must be valid JSON array of node objects.",
				Optional:            true,
				Computed:            true,
			},
			"connections_json": schema.StringAttribute{
				MarkdownDescription: "Workflow connections as JSON string. Must be valid JSON object mapping node connections.",
				Optional:            true,
				Computed:            true,
			},
			"settings_json": schema.StringAttribute{
				MarkdownDescription: "Workflow settings as JSON string. Must be valid JSON object.",
				Optional:            true,
				Computed:            true,
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
//
// Params:
//   - ctx: Context for the operation
//   - req: Configure request containing provider data
//   - resp: Configure response for error handling
//
// Returns:
//   - None: Updates resource in-place, populates resp with errors if any
func (r *WorkflowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	// Check provider data type.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		// Return result.
		return
	}

	r.client = clientData
}

// Create creates the resource and sets the initial Terraform state.
//
// Params:
//   - ctx: Context for the operation
//   - req: Create request containing plan data
//   - resp: Create response for state and error handling
//
// Returns:
//   - None: Updates resp with state and error diagnostics
func (r *WorkflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *WorkflowResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for plan parsing errors.
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse JSON fields.
	nodes, connections, settings := parseWorkflowJSON(plan, &resp.Diagnostics)
	// Check for JSON parsing errors.
	if resp.Diagnostics.HasError() {
		return
	}

	// Create workflow
	workflowRequest := n8nsdk.Workflow{
		Name:        plan.Name.ValueString(),
		Nodes:       nodes,
		Connections: connections,
		Settings:    settings,
	}

	workflow, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsPost(ctx).
		Workflow(workflowRequest).
		Execute()

	// Check for non-nil HTTP response.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for API error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating workflow",
			fmt.Sprintf("Could not create workflow, unexpected error: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Update tags if provided.
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() && workflow.Id != nil {
		r.updateWorkflowTags(ctx, *workflow.Id, plan, workflow, &resp.Diagnostics)
		// Check for tag update errors.
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Set ID and map response to state.
	plan.ID = types.StringPointerValue(workflow.Id)
	mapWorkflowToModel(ctx, workflow, plan, &resp.Diagnostics)

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}
// Read refreshes the Terraform state with the latest data.
//
// Params:
//   - ctx: Context for the operation
//   - req: Read request containing state data
//   - resp: Read response for state and error handling
//
// Returns:
//   - None: Updates resp with refreshed state and error diagnostics
func (r *WorkflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *WorkflowResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check for state parsing errors.
	if resp.Diagnostics.HasError() {
		return
	}

	// Get workflow from SDK.
	workflow, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdGet(ctx, state.ID.ValueString()).Execute()

	// Check for non-nil HTTP response.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for API error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading workflow",
			fmt.Sprintf("Could not read workflow ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		return
	}

	// Map response to state.
	mapWorkflowToModel(ctx, workflow, state, &resp.Diagnostics)

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
// Update updates the resource and sets the updated Terraform state on success.
//
// Params:
//   - ctx: Context for the operation
//   - req: Update request containing plan and state data
//   - resp: Update response for state and error handling
//
// Returns:
//   - None: Updates resp with updated state and error diagnostics
func (r *WorkflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *WorkflowResourceModel

	// Read Terraform plan and state data.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check for plan/state parsing errors.
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse JSON fields.
	nodes, connections, settings := parseWorkflowJSON(plan, &resp.Diagnostics)
	// Check for JSON parsing errors.
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle activation change.
	r.handleWorkflowActivation(ctx, plan, state, &resp.Diagnostics)
	// Check for activation change errors.
	if resp.Diagnostics.HasError() {
		return
	}

	// Update workflow content.
	workflowRequest := n8nsdk.Workflow{
		Name:        plan.Name.ValueString(),
		Nodes:       nodes,
		Connections: connections,
		Settings:    settings,
	}

	workflow, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdPut(ctx, plan.ID.ValueString()).
		Workflow(workflowRequest).
		Execute()

	// Check for non-nil HTTP response.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for API error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating workflow",
			fmt.Sprintf("Could not update workflow ID %s: %s\nHTTP Response: %v", plan.ID.ValueString(), err.Error(), httpResp),
		)
		return
	}

	// Update tags.
	r.updateWorkflowTags(ctx, plan.ID.ValueString(), plan, workflow, &resp.Diagnostics)
	// Check for tag update errors.
	if resp.Diagnostics.HasError() {
		return
	}

	// Map workflow response to state.
	mapWorkflowToModel(ctx, workflow, plan, &resp.Diagnostics)

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}
// Delete deletes the resource and removes the Terraform state on success.
//
// Params:
//   - ctx: Context for the operation
//   - req: Delete request containing state data
//   - resp: Delete response for error handling
//
// Returns:
//   - None: Updates resp with error diagnostics if any
func (r *WorkflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *WorkflowResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check for state parsing errors.
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete workflow using SDK.
	_, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdDelete(ctx, state.ID.ValueString()).Execute()
	// Check for non-nil HTTP response.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for API error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting workflow",
			fmt.Sprintf("Could not delete workflow ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		// Return error status.
		return
	}
}

// ImportState imports the resource into Terraform state.
//
// Params:
//   - ctx: Context for the operation
//   - req: Import state request containing resource ID
//   - resp: Import state response for state and error handling
//
// Returns:
//   - None: Updates resp with imported state and error diagnostics
func (r *WorkflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
