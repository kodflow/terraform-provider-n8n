// Package sourcecontrol contains resources and datasources for sourcecontrol management
package sourcecontrol

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SourceControlPullResourceModel describes the resource data model.
// Maps n8n source control pull attributes to Terraform schema, storing pull configuration and import results.
type SourceControlPullResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Force         types.Bool   `tfsdk:"force"`
	VariablesJSON types.String `tfsdk:"variables_json"`
	ResultJSON    types.String `tfsdk:"result_json"`
}
