// Package models defines data structures for execution data sources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSource describes the data source data model.
// It maps the Terraform schema attributes for reading a single execution,
// including workflow ID, execution status, timestamps, and optional data inclusion.
type DataSource struct {
	ID          types.String `tfsdk:"id"`
	WorkflowID  types.String `tfsdk:"workflow_id"`
	Finished    types.Bool   `tfsdk:"finished"`
	Mode        types.String `tfsdk:"mode"`
	StartedAt   types.String `tfsdk:"started_at"`
	StoppedAt   types.String `tfsdk:"stopped_at"`
	Status      types.String `tfsdk:"status"`
	IncludeData types.Bool   `tfsdk:"include_data"`
}
