package workflow

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WorkflowTransferResourceModel describes the workflow transfer resource data model.
// Captures the workflow ID and destination project ID for transfer operations.
type WorkflowTransferResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	WorkflowID           types.String `tfsdk:"workflow_id"`
	DestinationProjectID types.String `tfsdk:"destination_project_id"`
	TransferredAt        types.String `tfsdk:"transferred_at"`
}
