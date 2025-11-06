package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
)

// parseWorkflowJSON parses the JSON fields from a workflow model.
//
// Params:
//   - plan: The workflow resource model containing JSON data
//   - diags: Diagnostics for error reporting
//
// Returns:
//   - []n8nsdk.Node: Parsed workflow nodes
//   - map[string]interface{}: Parsed workflow connections
//   - n8nsdk.WorkflowSettings: Parsed workflow settings
func parseWorkflowJSON(plan *WorkflowResourceModel, diags *diag.Diagnostics) ([]n8nsdk.Node, map[string]interface{}, n8nsdk.WorkflowSettings) {
	// Parse nodes
	var nodes []n8nsdk.Node
	// Check for non-nil value.
	if !plan.NodesJSON.IsNull() && !plan.NodesJSON.IsUnknown() {
		// Check for error.
		if err := json.Unmarshal([]byte(plan.NodesJSON.ValueString()), &nodes); err != nil {
			diags.AddError("Invalid nodes JSON", fmt.Sprintf("Could not parse nodes_json: %s", err.Error()))
			// Return failure status.
			return []n8nsdk.Node{}, map[string]interface{}{}, n8nsdk.WorkflowSettings{}
		}
	} else {
		// Return empty slice.
		nodes = []n8nsdk.Node{}
	}

	// Parse connections
	var connections map[string]interface{}
	// Check for non-nil value.
	if !plan.ConnectionsJSON.IsNull() && !plan.ConnectionsJSON.IsUnknown() {
		// Check for error.
		if err := json.Unmarshal([]byte(plan.ConnectionsJSON.ValueString()), &connections); err != nil {
			diags.AddError("Invalid connections JSON", fmt.Sprintf("Could not parse connections_json: %s", err.Error()))
			// Return failure status.
			return []n8nsdk.Node{}, map[string]interface{}{}, n8nsdk.WorkflowSettings{}
		}
	} else {
		// Return empty slice.
		connections = map[string]interface{}{}
	}

	// Parse settings
	var settings n8nsdk.WorkflowSettings
	// Check for non-nil value.
	if !plan.SettingsJSON.IsNull() && !plan.SettingsJSON.IsUnknown() {
		// Check for error.
		if err := json.Unmarshal([]byte(plan.SettingsJSON.ValueString()), &settings); err != nil {
			diags.AddError("Invalid settings JSON", fmt.Sprintf("Could not parse settings_json: %s", err.Error()))
			// Return failure status.
			return []n8nsdk.Node{}, map[string]interface{}{}, n8nsdk.WorkflowSettings{}
		}
	}

	// Return result.
	return nodes, connections, settings
}

// mapTagsFromWorkflow maps tags from the SDK workflow to Terraform types.
//
// Params:
//   - ctx: Context for the API call
//   - workflow: The workflow from SDK containing tags
//   - diags: Diagnostics for error reporting
//
// Returns:
//   - types.List: Terraform list of tag IDs
func mapTagsFromWorkflow(ctx context.Context, workflow *n8nsdk.Workflow, diags *diag.Diagnostics) types.List {
	// Check length.
	if len(workflow.Tags) > 0 {
		tagIDs := make([]types.String, 0, len(workflow.Tags))
		// Iterate over items.
		for _, tag := range workflow.Tags {
			// Check for non-nil value.
			if tag.Id != nil {
				tagIDs = append(tagIDs, types.StringValue(*tag.Id))
			}
		}
		tagList, tagDiags := types.ListValueFrom(ctx, types.StringType, tagIDs)
		diags.Append(tagDiags...)
		// Return result.
		return tagList
	}

	// Empty list if no tags
	emptyList, emptyDiags := types.ListValueFrom(ctx, types.StringType, []types.String{})
	diags.Append(emptyDiags...)
	// Return result.
	return emptyList
}

// mapWorkflowToModel maps a workflow from the SDK to the Terraform model.
// This updates computed fields like timestamps, version, metadata, etc.
//
// Params:
//   - ctx: Context for the API call
//   - workflow: The workflow from SDK to map
//   - plan: The Terraform model to update
//   - diags: Diagnostics for error reporting
//
// Returns:
//   - None: Updates plan in-place
func mapWorkflowToModel(ctx context.Context, workflow *n8nsdk.Workflow, plan *WorkflowResourceModel, diags *diag.Diagnostics) {
	// Basic fields
	plan.Name = types.StringValue(workflow.Name)
	// Check for non-nil value.
	if workflow.Active != nil {
		plan.Active = types.BoolPointerValue(workflow.Active)
	}

	// Tags
	plan.Tags = mapTagsFromWorkflow(ctx, workflow, diags)

	// Timestamps and metadata
	// Check for non-nil value.
	if workflow.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(workflow.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	// Check for non-nil value.
	if workflow.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(workflow.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	// Check for non-nil value.
	if workflow.VersionId != nil {
		plan.VersionID = types.StringPointerValue(workflow.VersionId)
	}
	// Check for non-nil value.
	if workflow.IsArchived != nil {
		plan.IsArchived = types.BoolPointerValue(workflow.IsArchived)
	}
	// Check for non-nil value.
	if workflow.TriggerCount != nil {
		plan.TriggerCount = types.Int64Value(int64(*workflow.TriggerCount))
	}

	// Map objects
	// Check for non-nil value.
	if workflow.Meta != nil {
		metaMap, metaDiags := types.MapValueFrom(ctx, types.StringType, workflow.Meta)
		diags.Append(metaDiags...)
		// Check condition.
		if !diags.HasError() {
			plan.Meta = metaMap
		}
	}
	// Check for non-nil value.
	if workflow.PinData != nil {
		pinDataMap, pinDiags := types.MapValueFrom(ctx, types.StringType, workflow.PinData)
		diags.Append(pinDiags...)
		// Check condition.
		if !diags.HasError() {
			plan.PinData = pinDataMap
		}
	}

	// Serialize JSON fields
	serializeWorkflowJSON(workflow, plan)
}

// serializeWorkflowJSON serializes workflow nodes, connections and settings back to JSON strings.
//
// Params:
//   - workflow: The workflow from SDK containing data to serialize
//   - plan: The Terraform model to update with serialized JSON
//
// Returns:
//   - None: Updates plan in-place
func serializeWorkflowJSON(workflow *n8nsdk.Workflow, plan *WorkflowResourceModel) {
	// Check for non-nil value.
	if workflow.Nodes != nil {
		// Check for error.
		if nodesJSON, err := json.Marshal(workflow.Nodes); err == nil {
			plan.NodesJSON = types.StringValue(string(nodesJSON))
		}
	}
	// Check for non-nil value.
	if workflow.Connections != nil {
		// Check for error.
		if connectionsJSON, err := json.Marshal(workflow.Connections); err == nil {
			plan.ConnectionsJSON = types.StringValue(string(connectionsJSON))
		}
	}
	// Check for error.
	if settingsJSON, err := json.Marshal(workflow.Settings); err == nil {
		plan.SettingsJSON = types.StringValue(string(settingsJSON))
	}
}

// convertTagIDsToTagIdsInner converts string tag IDs to SDK TagIdsInner format.
//
// Params:
//   - tagIDs: List of tag ID strings
//
// Returns:
//   - []n8nsdk.TagIdsInner: Converted tag IDs in SDK format
func convertTagIDsToTagIdsInner(tagIDs []string) []n8nsdk.TagIdsInner {
	tagIdsInner := make([]n8nsdk.TagIdsInner, 0, len(tagIDs))
	// Iterate over items.
	for _, tagID := range tagIDs {
		tagIdsInner = append(tagIdsInner, n8nsdk.TagIdsInner{Id: tagID})
	}
	// Return result.
	return tagIdsInner
}

// handleWorkflowActivation handles the activation/deactivation of a workflow.
//
// Params:
//   - ctx: Context for the API call
//   - plan: The desired workflow state
//   - state: The current workflow state
//   - diags: Diagnostics for error reporting
//
// Returns:
//   - None: Updates workflow via API
func (r *WorkflowResource) handleWorkflowActivation(ctx context.Context, plan, state *WorkflowResourceModel, diags *diag.Diagnostics) {
	activeChanged := !plan.Active.IsNull() && !state.Active.IsNull() &&
		plan.Active.ValueBool() != state.Active.ValueBool()

	// Check condition.
	if !activeChanged {
		// Return success status.
		return
	}

	var httpResp *http.Response
	var err error

	// Use dedicated activate/deactivate endpoints
	// Check condition.
	if plan.Active.ValueBool() {
		_, httpResp, err = r.client.APIClient.WorkflowAPI.WorkflowsIdActivatePost(ctx, plan.ID.ValueString()).Execute()
	} else {
		_, httpResp, err = r.client.APIClient.WorkflowAPI.WorkflowsIdDeactivatePost(ctx, plan.ID.ValueString()).Execute()
	}

	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		action := "activate"
		// Check condition.
		if !plan.Active.ValueBool() {
			action = "deactivate"
		}
		diags.AddError(
			fmt.Sprintf("Error changing workflow activation status to %s", action),
			fmt.Sprintf("Could not %s workflow ID %s: %s\nHTTP Response: %v", action, plan.ID.ValueString(), err.Error(), httpResp),
		)
	}
}

// updateWorkflowTags updates the tags for a workflow.
//
// Params:
//   - ctx: Context for the API call
//   - workflowID: The workflow ID to update tags for
//   - plan: The Terraform model containing desired tags
//   - workflow: The SDK workflow to update
//   - diags: Diagnostics for error reporting
//
// Returns:
//   - None: Updates workflow tags via API
func (r *WorkflowResource) updateWorkflowTags(ctx context.Context, workflowID string, plan *WorkflowResourceModel, workflow *n8nsdk.Workflow, diags *diag.Diagnostics) {
	// Check for null value.
	if plan.Tags.IsNull() || plan.Tags.IsUnknown() {
		// Return success status.
		return
	}

	var tagIDs []string
	diags.Append(plan.Tags.ElementsAs(ctx, &tagIDs, false)...)
	// Check condition.
	if diags.HasError() {
		// Return failure status.
		return
	}

	tagIdsInner := convertTagIDsToTagIdsInner(tagIDs)

	tags, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdTagsPut(ctx, workflowID).
		TagIdsInner(tagIdsInner).
		Execute()

	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		diags.AddError(
			"Error updating workflow tags",
			fmt.Sprintf("Could not update tags for workflow ID %s: %s\nHTTP Response: %v", workflowID, err.Error(), httpResp),
		)
		// Return failure status.
		return
	}

	workflow.Tags = tags
}
