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

// Ensure WorkflowTransferResource implements required interfaces.
var (
	_ resource.Resource                = &WorkflowTransferResource{}
	_ resource.ResourceWithConfigure   = &WorkflowTransferResource{}
	_ resource.ResourceWithImportState = &WorkflowTransferResource{}
)

// WorkflowTransferResource defines the resource implementation for transferring a workflow to a project.
type WorkflowTransferResource struct {
	client *providertypes.N8nClient
}

// WorkflowTransferResourceModel describes the resource data model.
type WorkflowTransferResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	WorkflowID           types.String `tfsdk:"workflow_id"`
	DestinationProjectID types.String `tfsdk:"destination_project_id"`
	TransferredAt        types.String `tfsdk:"transferred_at"`
}

// NewWorkflowTransferResource creates a new WorkflowTransferResource instance.
func NewWorkflowTransferResource() resource.Resource {
	// Return result.
	return &WorkflowTransferResource{}
}

// Metadata returns the resource type name.
func (r *WorkflowTransferResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_transfer"
}

// Schema defines the schema for the resource.
func (r *WorkflowTransferResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Transfers a workflow to another project. This is a one-time operation resource that triggers the transfer when created.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Resource identifier (generated)",
				Computed:            true,
			},
			"workflow_id": schema.StringAttribute{
				MarkdownDescription: "ID of the workflow to transfer",
				Required:            true,
			},
			"destination_project_id": schema.StringAttribute{
				MarkdownDescription: "ID of the destination project",
				Required:            true,
			},
			"transferred_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the transfer occurred",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *WorkflowTransferResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create triggers the workflow transfer.
func (r *WorkflowTransferResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WorkflowTransferResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Build transfer request
	transferReq := n8nsdk.NewWorkflowsIdTransferPutRequest(plan.DestinationProjectID.ValueString())

	// Execute transfer
	httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdTransferPut(ctx, plan.WorkflowID.ValueString()).
		WorkflowsIdTransferPutRequest(*transferReq).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error transferring workflow",
			fmt.Sprintf("Could not transfer workflow %s to project %s: %s\nHTTP Response: %v",
				plan.WorkflowID.ValueString(), plan.DestinationProjectID.ValueString(), err.Error(), httpResp),
		)
		// Return result.
		return
	}

	// Set computed fields
	plan.ID = types.StringValue(fmt.Sprintf("%s-to-%s", plan.WorkflowID.ValueString(), plan.DestinationProjectID.ValueString()))
	plan.TransferredAt = types.StringValue(fmt.Sprintf("transfer-response-%d", httpResp.StatusCode))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the resource state. For transfer operations, we just keep the current state.
func (r *WorkflowTransferResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkflowTransferResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Transfer operations are one-time actions, so we just maintain the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is not supported for transfer operations.
func (r *WorkflowTransferResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Workflow transfer resources cannot be updated. To transfer again, create a new resource.",
	)
}

// Delete removes the resource from state but doesn't perform any API operation.
func (r *WorkflowTransferResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Transfer operations cannot be undone, so we just remove from state
}

// ImportState imports the resource into Terraform state.
func (r *WorkflowTransferResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
