package workflow

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WorkflowItemModel maps individual workflow attributes within the Terraform schema.
// Each item represents a single workflow with its identifier, name, and activation status.
type WorkflowItemModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Active types.Bool   `tfsdk:"active"`
}
