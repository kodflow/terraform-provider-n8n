package execution

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/execution/models"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/shared/constants"
)

// Ensure ExecutionRetryResource implements required interfaces.
var (
	_ resource.Resource                = &ExecutionRetryResource{}
	_ ExecutionRetryResourceInterface  = &ExecutionRetryResource{}
	_ resource.ResourceWithConfigure   = &ExecutionRetryResource{}
	_ resource.ResourceWithImportState = &ExecutionRetryResource{}
)

// NewExecutionRetryResource creates a new ExecutionRetryResource instance.
// Returns:
//   - *ExecutionRetryResource: A new ExecutionRetryResource instance
func NewExecutionRetryResource() *ExecutionRetryResource {
	// Return result.
	return &ExecutionRetryResource{}
}

// NewExecutionRetryResourceWrapper creates a new ExecutionRetryResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - resource.Resource: the wrapped ExecutionRetryResource instance
func NewExecutionRetryResourceWrapper() resource.Resource {
	// Return the wrapped resource instance.
	return NewExecutionRetryResource()
}

// ExecutionRetryResourceInterface defines the interface for ExecutionRetryResource.
type ExecutionRetryResourceInterface interface {
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

// ExecutionRetryResource defines the resource implementation for retrying an execution.
// Terraform resource that manages CRUD operations for execution retries with the n8n API.
type ExecutionRetryResource struct {
	// client is the N8n API client used for execution retry operations.
	client *client.N8nClient
}

// Metadata returns the resource type name.
// Params:
//   - ctx: Context for request handling
//   - req: Metadata request from Terraform
//   - resp: Metadata response to Terraform
//
// Returns:
//   - void
func (r *ExecutionRetryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_execution_retry"
}

// Schema defines the schema for the resource.
// Params:
//   - ctx: Context for request handling
//   - req: Schema request from Terraform
//   - resp: Schema response to Terraform
//
// Returns:
//   - void
func (r *ExecutionRetryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retries a failed n8n workflow execution. This resource triggers a retry operation and captures the resulting execution details.",
		Attributes:          r.schemaAttributes(),
	}
}

// schemaAttributes returns the attribute definitions for the execution retry resource schema.
//
// Returns:
//   - map[string]schema.Attribute: the resource attribute definitions
func (r *ExecutionRetryResource) schemaAttributes() map[string]schema.Attribute {
	// Return execution retry schema attributes.
	return map[string]schema.Attribute{
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
	}
}

// Configure adds the provider configured client to the resource.
// Params:
//   - ctx: Context for request handling
//   - req: Configure request from Terraform
//   - resp: Configure response to Terraform
//
// Returns:
//   - void
func (r *ExecutionRetryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create triggers the execution retry.
// Params:
//   - ctx: Context for request handling
//   - req: Create request from Terraform
//   - resp: Create response to Terraform
//
// Returns:
//   - void
func (r *ExecutionRetryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *models.RetryResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Convert execution ID to float32 as required by the API
	executionID, err := strconv.ParseFloat(plan.ExecutionID.ValueString(), constants.FLOAT32_BIT_SIZE)
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
	populateRetryExecutionData(execution, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the resource state. For retry operations, we just keep the current state.
// Params:
//   - ctx: Context for request handling
//   - req: Read request from Terraform
//   - resp: Read response to Terraform
//
// Returns:
//   - void
func (r *ExecutionRetryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *models.RetryResource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Retry operations are one-time actions, so we just maintain the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is not supported for retry operations.
// Params:
//   - ctx: Context for request handling
//   - req: Update request from Terraform
//   - resp: Update response to Terraform
//
// Returns:
//   - void
func (r *ExecutionRetryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Execution retry resources cannot be updated. To retry again, create a new resource.",
	)
}

// Delete removes the resource from state but doesn't perform any API operation.
// Params:
//   - ctx: Context for request handling
//   - req: Delete request from Terraform
//   - resp: Delete response to Terraform
//
// Returns:
//   - void
func (r *ExecutionRetryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retry operations cannot be undone, so we just remove from state
}

// ImportState imports the resource into Terraform state.
// Params:
//   - ctx: Context for request handling
//   - req: ImportState request from Terraform
//   - resp: ImportState response to Terraform
//
// Returns:
//   - void
func (r *ExecutionRetryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("execution_id"), req, resp)
}

// populateRetryExecutionData populates the retry execution resource model from the API response.
//
// Params:
//   - execution: the execution response from the API
//   - model: the resource model to populate
func populateRetryExecutionData(execution *n8nsdk.Execution, model *models.RetryResource) {
	// Check if execution ID is available.
	if execution.Id != nil {
		model.NewExecutionID = types.StringValue(fmt.Sprintf("%v", *execution.Id))
	}
	// Check if workflow ID is available.
	if execution.WorkflowId != nil {
		model.WorkflowID = types.StringValue(fmt.Sprintf("%v", *execution.WorkflowId))
	}
	// Check if finished status is available.
	if execution.Finished != nil {
		model.Finished = types.BoolPointerValue(execution.Finished)
	}
	// Check if execution mode is available.
	if execution.Mode != nil {
		model.Mode = types.StringPointerValue(execution.Mode)
	}
	// Check if started timestamp is available.
	if execution.StartedAt != nil {
		model.StartedAt = types.StringValue(execution.StartedAt.String())
	}
	// Check if stopped timestamp is available.
	if execution.StoppedAt.IsSet() {
		stoppedAt := execution.StoppedAt.Get()
		// Check if stopped timestamp value is not nil.
		if stoppedAt != nil {
			model.StoppedAt = types.StringValue(stoppedAt.String())
		}
	}
	// Check if execution status is available.
	if execution.Status != nil {
		model.Status = types.StringPointerValue(execution.Status)
	}
}
