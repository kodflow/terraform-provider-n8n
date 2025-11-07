// Package execution contains resources and datasources for execution management
package execution

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ExecutionRetryResourceModel describes the resource data model.
// Maps n8n execution retry attributes to Terraform schema, storing execution details and retry results.
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
