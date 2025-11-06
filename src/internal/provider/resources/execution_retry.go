package resources

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure ExecutionRetryResource implements required interfaces.
var (
	_ resource.Resource                = &ExecutionRetryResource{}
	_ resource.ResourceWithConfigure   = &ExecutionRetryResource{}
	_ resource.ResourceWithImportState = &ExecutionRetryResource{}
)

// ExecutionRetryResource defines the resource implementation for retrying an execution.
type ExecutionRetryResource struct {
	client *providertypes.N8nClient
}

// ExecutionRetryResourceModel describes the resource data model.
type ExecutionRetryResourceModel struct {
	ExecutionID    types.String `tfsdk:"execution_id"`
	NewExecutionID types.String `tfsdk:"new_execution_id"`
	WorkflowID     types.String `tfsdk:"workflow_id"`
	Finished       types.Bool   `tfsdk:"finished"`
	Mode           types.String `tfsdk:"mode"`
	StartedAt      types.String `tfsdk:"started_at"`
	StoppedAt      types.String `tfsdk:"stopped_at"`
	Status         types.String `tfsdk:"status"`
}

// NewExecutionRetryResource creates a new ExecutionRetryResource instance.
func NewExecutionRetryResource() resource.Resource {
	// Return result.
	return &ExecutionRetryResource{}
}

// Metadata returns the resource type name.
func (r *ExecutionRetryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_execution_retry"
}

// Schema defines the schema for the resource.
func (r *ExecutionRetryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retries a failed n8n workflow execution. This resource triggers a retry operation and captures the resulting execution details.",

		Attributes: map[string]schema.Attribute{
			"execution_id": schema.StringAttribute{
				MarkdownDescription: "ID of the execution to retry",
				Required:            true,
			},
			"new_execution_id": schema.StringAttribute{
				MarkdownDescription: "ID of the new execution created by the retry",
				Computed:            true,
			},
			"workflow_id": schema.StringAttribute{
				MarkdownDescription: "ID of the workflow that was executed",
				Computed:            true,
			},
			"finished": schema.BoolAttribute{
				MarkdownDescription: "Whether the retried execution finished",
				Computed:            true,
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "Execution mode",
				Computed:            true,
			},
			"started_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the retried execution started",
				Computed:            true,
			},
			"stopped_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the retried execution stopped",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Status of the retried execution",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *ExecutionRetryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create triggers the execution retry.
func (r *ExecutionRetryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ExecutionRetryResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert execution ID to float32 as required by the API
	executionID, err := strconv.ParseFloat(plan.ExecutionID.ValueString(), 32)
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Execution ID",
			fmt.Sprintf("Could not parse execution ID %s as number: %s", plan.ExecutionID.ValueString(), err.Error()),
		)
		// Return result.
		return
	}

	// Retry the execution
	execution, httpResp, err := r.client.APIClient.ExecutionAPI.ExecutionsIdRetryPost(ctx, float32(executionID)).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrying execution",
			fmt.Sprintf("Could not retry execution ID %s: %s\nHTTP Response: %v", plan.ExecutionID.ValueString(), err.Error(), httpResp),
		)
		// Return result.
		return
	}

	// Populate response data
	if execution.Id != nil {
		plan.NewExecutionID = types.StringValue(fmt.Sprintf("%v", *execution.Id))
	}
	// Check for non-nil value.
	if execution.WorkflowId != nil {
		plan.WorkflowID = types.StringValue(fmt.Sprintf("%v", *execution.WorkflowId))
	}
	// Check for non-nil value.
	if execution.Finished != nil {
		plan.Finished = types.BoolPointerValue(execution.Finished)
	}
	// Check for non-nil value.
	if execution.Mode != nil {
		plan.Mode = types.StringPointerValue(execution.Mode)
	}
	// Check for non-nil value.
	if execution.StartedAt != nil {
		plan.StartedAt = types.StringValue(execution.StartedAt.String())
	}
	// Check condition.
	if execution.StoppedAt.IsSet() {
		stoppedAt := execution.StoppedAt.Get()
		// Check for non-nil value.
		if stoppedAt != nil {
			plan.StoppedAt = types.StringValue(stoppedAt.String())
		}
	}
	// Check for non-nil value.
	if execution.Status != nil {
		plan.Status = types.StringPointerValue(execution.Status)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the resource state. For retry operations, we just keep the current state.
func (r *ExecutionRetryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ExecutionRetryResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Retry operations are one-time actions, so we just maintain the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is not supported for retry operations.
func (r *ExecutionRetryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Execution retry resources cannot be updated. To retry again, create a new resource.",
	)
}

// Delete removes the resource from state but doesn't perform any API operation.
func (r *ExecutionRetryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retry operations cannot be undone, so we just remove from state
}

// ImportState imports the resource into Terraform state.
func (r *ExecutionRetryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("execution_id"), req, resp)
}
