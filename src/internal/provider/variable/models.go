// Package variable contains resources and datasources for variable management
package variable

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// VariableResourceModel describes the resource data model.
// Maps n8n variable attributes to Terraform schema, storing variable configuration and project association.
type VariableResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	Type      types.String `tfsdk:"type"`
	ProjectID types.String `tfsdk:"project_id"`
}
