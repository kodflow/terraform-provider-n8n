package variable

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// VariableDataSourceModel maps Terraform schema attributes for variable data.
// It represents a single variable with all related attributes from the n8n API.
type VariableDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	Type      types.String `tfsdk:"type"`
	ProjectID types.String `tfsdk:"project_id"`
}
