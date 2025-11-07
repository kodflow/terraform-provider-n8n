package execution

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ExecutionsDataSourceModel maps the Terraform schema to the datasource response.
// It represents the filtered execution list with workflow and execution details from the n8n API.
type ExecutionsDataSourceModel struct {
	WorkflowID  types.String         `tfsdk:"workflow_id"`
	ProjectID   types.String         `tfsdk:"project_id"`
	Status      types.String         `tfsdk:"status"`
	IncludeData types.Bool           `tfsdk:"include_data"`
	Executions  []ExecutionItemModel `tfsdk:"executions"`
}
