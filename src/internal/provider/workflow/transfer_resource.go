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

// WorkflowTransferResourceInterface defines the complete interface for workflow transfer resource.
// It combines all standard Terraform resource interfaces required for managing workflow transfers.
type WorkflowTransferResourceInterface interface {
	resource.Resource
	resource.ResourceWithConfigure
	resource.ResourceWithImportState
	Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse)
	Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse)
	Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse)
	Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse)
	Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse)
	Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse)
	Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse)
	ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse)
}

// Ensure WorkflowTransferResource implements required interfaces.
var (
	_ WorkflowTransferResourceInterface = &WorkflowTransferResource{}
	_ resource.Resource                 = &WorkflowTransferResource{}
	_ resource.ResourceWithConfigure    = &WorkflowTransferResource{}
	_ resource.ResourceWithImportState  = &WorkflowTransferResource{}
)

// WorkflowTransferResource defines the resource implementation for transferring a workflow to a project.
// This resource manages one-time workflow transfer operations in n8n, moving workflows between projects
// via the n8n API. The transfer operation is triggered on resource creation.
type WorkflowTransferResource struct {
	client *client.N8nClient
}

// NewWorkflowTransferResource creates a new WorkflowTransferResource instance.
//
// Returns:
//   - resource.Resource: a new workflow transfer resource ready for use
func NewWorkflowTransferResource() resource.Resource {
	// Initialize and return new WorkflowTransferResource instance
	return &WorkflowTransferResource{}
}

// Metadata returns the resource type name.
//
// Params:
//   - ctx: context for the operation
//   - req: metadata request containing provider type name
//   - resp: metadata response to populate with resource type name
//
// Returns:
//   - void: modifies resp in place, sets resp.TypeName
func (r *WorkflowTransferResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_transfer"
}

// Schema defines the schema for the resource.
//
// Params:
//   - ctx: context for the operation
//   - req: schema request
//   - resp: schema response to populate with resource schema definition
//
// Returns:
//   - void: modifies resp in place, sets resp.Schema with workflow transfer attributes
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
//
// Params:
//   - ctx: context for the operation
//   - req: configure request containing provider data (N8nClient)
//   - resp: configure response to populate with diagnostics on error
//
// Returns:
//   - void: modifies r.client in place, populates resp.Diagnostics on error
func (r *WorkflowTransferResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Early return if no provider data is configured
	if req.ProviderData == nil {
		return
	}

	// Type assert provider data to N8nClient
	clientData, ok := req.ProviderData.(*client.N8nClient)
	// Return error if provider data is not N8nClient type
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		return
	}

	// Store client for use in CRUD operations
	r.client = clientData
}

// Create triggers the workflow transfer.
//
// Params:
//   - ctx: context for the operation
//   - req: create request containing planned workflow transfer configuration
//   - resp: create response to populate with state and diagnostics
//
// Returns:
//   - void: modifies resp in place, populates resp.State or resp.Diagnostics on error
func (r *WorkflowTransferResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Initialize model from plan
	plan := &WorkflowTransferResourceModel{}

	// Parse plan into model
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	// Return early if plan parsing failed
	if resp.Diagnostics.HasError() {
		return
	}

	// Build transfer request with destination project ID
	transferReq := n8nsdk.NewWorkflowsIdTransferPutRequest(plan.DestinationProjectID.ValueString())

	// Execute workflow transfer operation via n8n API
	httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdTransferPut(ctx, plan.WorkflowID.ValueString()).
		WorkflowsIdTransferPutRequest(*transferReq).Execute()
	// Clean up HTTP response body if present
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Return error if transfer operation failed
	if err != nil {
		resp.Diagnostics.AddError(
			"Error transferring workflow",
			fmt.Sprintf("Could not transfer workflow %s to project %s: %s\nHTTP Response: %v",
				plan.WorkflowID.ValueString(), plan.DestinationProjectID.ValueString(), err.Error(), httpResp),
		)
		return
	}

	// Update model with computed values from transfer operation
	plan.ID = types.StringValue(fmt.Sprintf("%s-to-%s", plan.WorkflowID.ValueString(), plan.DestinationProjectID.ValueString()))
	plan.TransferredAt = types.StringValue(fmt.Sprintf("transfer-response-%d", httpResp.StatusCode))

	// Persist updated model to state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the resource state. For transfer operations, we just keep the current state.
//
// Params:
//   - ctx: context for the operation
//   - req: read request containing current state
//   - resp: read response to populate with refreshed state
//
// Returns:
//   - void: modifies resp in place, populates resp.State or resp.Diagnostics on error
func (r *WorkflowTransferResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Initialize state model
	state := &WorkflowTransferResourceModel{}

	// Read state from Terraform
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	// Return early if state read failed
	if resp.Diagnostics.HasError() {
		return
	}

	// Transfer operations are one-time actions, so maintain the existing state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update is not supported for transfer operations.
//
// Params:
//   - ctx: context for the operation
//   - req: update request containing plan and state
//   - resp: update response to populate with diagnostics
//
// Returns:
//   - void: modifies resp in place, populates resp.Diagnostics with error
func (r *WorkflowTransferResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Workflow transfer resources cannot be updated. To transfer again, create a new resource.",
	)
}

// Delete removes the resource from state but doesn't perform any API operation.
//
// Params:
//   - ctx: context for the operation
//   - req: delete request containing state
//   - resp: delete response (unused for transfer operations)
//
// Returns:
//   - void: no operation performed, resource removed from state
func (r *WorkflowTransferResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Transfer operations cannot be undone, so we just remove from state
}

// ImportState imports the resource into Terraform state.
//
// Params:
//   - ctx: context for the operation
//   - req: import request containing resource ID
//   - resp: import response to populate with state
//
// Returns:
//   - void: modifies resp in place, populates resp.State with imported ID
func (r *WorkflowTransferResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
