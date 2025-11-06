package execution

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
)

// mapExecutionToItem maps an SDK execution to the execution item model.
//
// Params:
//   - execution: SDK execution to map
//
// Returns:
//   - ExecutionItemModel: Mapped execution item model
func mapExecutionToItem(execution *n8nsdk.Execution) ExecutionItemModel {
	item := ExecutionItemModel{}

	// Check for non-nil value.
	if execution.Id != nil {
		item.ID = types.StringValue(fmt.Sprintf("%v", *execution.Id))
	}
	// Check for non-nil value.
	if execution.WorkflowId != nil {
		item.WorkflowID = types.StringValue(fmt.Sprintf("%v", *execution.WorkflowId))
	}
	// Check for non-nil value.
	if execution.Finished != nil {
		item.Finished = types.BoolPointerValue(execution.Finished)
	}
	// Check for non-nil value.
	if execution.Mode != nil {
		item.Mode = types.StringPointerValue(execution.Mode)
	}
	// Check for non-nil value.
	if execution.StartedAt != nil {
		item.StartedAt = types.StringValue(execution.StartedAt.String())
	}
	// Check for non-nil value.
	if execution.StoppedAt.IsSet() {
		stoppedAt := execution.StoppedAt.Get()
		// Check for non-nil value.
		if stoppedAt != nil {
			item.StoppedAt = types.StringValue(stoppedAt.String())
		}
	}
	// Check for non-nil value.
	if execution.Status != nil {
		item.Status = types.StringPointerValue(execution.Status)
	}

	// Return result.
	return item
}
