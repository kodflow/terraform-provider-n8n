package variable

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// VariablesDataSourceModel maps the Terraform schema attributes for the variables datasource.
// It represents the complete set of variables data returned by the n8n API with optional filtering.
type VariablesDataSourceModel struct {
	ProjectID types.String        `tfsdk:"project_id"`
	State     types.String        `tfsdk:"state"`
	Variables []VariableItemModel `tfsdk:"variables"`
}
