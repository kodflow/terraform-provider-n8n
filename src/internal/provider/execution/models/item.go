package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Item represents a single execution in the list returned from the n8n API.
// It contains the execution metadata including timestamps, status, and workflow reference.
type Item struct {
	ID         types.String `tfsdk:"id"`
	WorkflowID types.String `tfsdk:"workflow_id"`
	Finished   types.Bool   `tfsdk:"finished"`
	Mode       types.String `tfsdk:"mode"`
	StartedAt  types.String `tfsdk:"started_at"`
	StoppedAt  types.String `tfsdk:"stopped_at"`
	Status     types.String `tfsdk:"status"`
}
