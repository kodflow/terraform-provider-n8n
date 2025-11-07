package workflow

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WorkflowDataSourceModel maps the Terraform schema attributes for a single workflow datasource.
// It represents workflow metadata including its identifier, name, and activation status.
type WorkflowDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Active types.Bool   `tfsdk:"active"`
}
