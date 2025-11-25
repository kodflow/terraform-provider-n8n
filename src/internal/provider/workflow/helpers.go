// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package workflow implements workflow management resources and data sources.
package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/models"
)

// CALLER_POLICY_DEFAULT is the default value for the CallerPolicy workflow setting.
// The n8n API returns this value even when not explicitly set by the user.
const CALLER_POLICY_DEFAULT string = "workflowsFromSameOwner"

// parseWorkflowJSON parses the JSON fields from a workflow model.
//
// Params:
//   - plan: The workflow resource model containing JSON data
//   - diags: Diagnostics for error reporting
//
// Returns:
//   - []n8nsdk.Node: Parsed workflow nodes
//   - map[string]any: Parsed workflow connections
//   - n8nsdk.WorkflowSettings: Parsed workflow settings
func parseWorkflowJSON(plan *models.Resource, diags *diag.Diagnostics) ([]n8nsdk.Node, map[string]any, n8nsdk.WorkflowSettings) {
	// Parse nodes
	var nodes []n8nsdk.Node
	// Check for non-nil value.
	if !plan.NodesJSON.IsNull() && !plan.NodesJSON.IsUnknown() {
		// Check for error.
		if err := json.Unmarshal([]byte(plan.NodesJSON.ValueString()), &nodes); err != nil {
			diags.AddError("Invalid nodes JSON", fmt.Sprintf("Could not parse nodes_json: %s", err.Error()))
			// Return failure status.
			return []n8nsdk.Node{}, map[string]any{}, n8nsdk.WorkflowSettings{}
		}
	} else {
		// Return empty slice.
		nodes = []n8nsdk.Node{}
	}

	// Parse connections
	var connections map[string]any
	// Check for non-nil value.
	if !plan.ConnectionsJSON.IsNull() && !plan.ConnectionsJSON.IsUnknown() {
		// Check for error.
		if err := json.Unmarshal([]byte(plan.ConnectionsJSON.ValueString()), &connections); err != nil {
			diags.AddError("Invalid connections JSON", fmt.Sprintf("Could not parse connections_json: %s", err.Error()))
			// Return failure status.
			return []n8nsdk.Node{}, map[string]any{}, n8nsdk.WorkflowSettings{}
		}
	} else {
		// Return empty slice.
		connections = map[string]any{}
	}

	// Parse settings
	var settings n8nsdk.WorkflowSettings
	// Check for non-nil value.
	if !plan.SettingsJSON.IsNull() && !plan.SettingsJSON.IsUnknown() {
		// Check for error.
		if err := json.Unmarshal([]byte(plan.SettingsJSON.ValueString()), &settings); err != nil {
			diags.AddError("Invalid settings JSON", fmt.Sprintf("Could not parse settings_json: %s", err.Error()))
			// Return failure status.
			return []n8nsdk.Node{}, map[string]any{}, n8nsdk.WorkflowSettings{}
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
//   - types.Set: Terraform set of tag IDs
func mapTagsFromWorkflow(ctx context.Context, workflow *n8nsdk.Workflow, diags *diag.Diagnostics) types.Set {
	// Check length.
	if len(workflow.Tags) > 0 {
		// Collect tag IDs
		tagIDs := make([]types.String, 0, len(workflow.Tags))
		// Iterate over items.
		for _, tag := range workflow.Tags {
			// Check for non-nil value.
			if tag.Id != nil {
				tagIDs = append(tagIDs, types.StringValue(*tag.Id))
			}
		}

		tagSet, tagDiags := types.SetValueFrom(ctx, types.StringType, tagIDs)
		diags.Append(tagDiags...)
		// Return result.
		return tagSet
	}

	// Return null set if no tags to avoid inconsistent result errors.
	return types.SetNull(types.StringType)
}

// mapWorkflowBasicFields maps basic workflow fields to the model.
//
// Params:
//   - workflow: The n8n workflow to map from
//   - plan: The resource model to update
func mapWorkflowBasicFields(workflow *n8nsdk.Workflow, plan *models.Resource) {
	// Set active status if available.
	if workflow.Active != nil {
		plan.Active = types.BoolPointerValue(workflow.Active)
	}
	// Set version ID if available.
	if workflow.VersionId != nil {
		plan.VersionID = types.StringPointerValue(workflow.VersionId)
	}
	// Set archived status if available.
	if workflow.IsArchived != nil {
		plan.IsArchived = types.BoolPointerValue(workflow.IsArchived)
	}
	// Set trigger count if available.
	if workflow.TriggerCount != nil {
		plan.TriggerCount = types.Int64Value(int64(*workflow.TriggerCount))
	}
}

// mapWorkflowProjectID maps the workflow project ID from shared info to the model.
//
// Params:
//   - workflow: The n8n workflow to map from
//   - plan: The resource model to update
func mapWorkflowProjectID(workflow *n8nsdk.Workflow, plan *models.Resource) {
	// Extract projectId from workflow.Shared[0].ProjectId
	// Check if workflow has shared projects
	if len(workflow.Shared) > 0 {
		// Check if the first shared project has a project ID
		if workflow.Shared[0].ProjectId != nil {
			plan.ProjectID = types.StringPointerValue(workflow.Shared[0].ProjectId)
		} else {
			// Project ID is nil, set to null in state
			plan.ProjectID = types.StringNull()
		}
	} else {
		// No shared projects, workflow is in default location
		plan.ProjectID = types.StringNull()
	}
}

// mapWorkflowTimestamps maps workflow timestamps to the model.
//
// Params:
//   - workflow: The n8n workflow to map from
//   - plan: The resource model to update
func mapWorkflowTimestamps(workflow *n8nsdk.Workflow, plan *models.Resource) {
	// Set creation timestamp if available.
	if workflow.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(workflow.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	// Set update timestamp if available.
	if workflow.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(workflow.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
}

// mapWorkflowToModel maps a workflow from the SDK to the Terraform model.
// This updates computed fields like timestamps, version, metadata, etc.
//
// Params:
//   - ctx: Context for the API call
//   - workflow: The workflow from SDK to map
//   - plan: The Terraform model to update
//   - diags: Diagnostics for error reporting
func mapWorkflowToModel(ctx context.Context, workflow *n8nsdk.Workflow, plan *models.Resource, diags *diag.Diagnostics) {
	// Basic fields
	plan.Name = types.StringValue(workflow.Name)

	// Map simple fields
	mapWorkflowBasicFields(workflow, plan)

	// Tags
	plan.Tags = mapTagsFromWorkflow(ctx, workflow, diags)

	// Project ID from shared workflow info
	mapWorkflowProjectID(workflow, plan)

	// Map timestamps
	mapWorkflowTimestamps(workflow, plan)

	// Map objects
	// Check for non-nil value.
	if workflow.Meta != nil {
		metaMap, metaDiags := types.MapValueFrom(ctx, types.StringType, workflow.Meta)
		diags.Append(metaDiags...)
		// Check condition.
		if !diags.HasError() {
			plan.Meta = metaMap
		}
	} else {
		// Set null map when API returns nil to ensure attribute is known.
		plan.Meta = types.MapNull(types.StringType)
	}
	// Check for non-nil value.
	if workflow.PinData != nil {
		pinDataMap, pinDiags := types.MapValueFrom(ctx, types.StringType, workflow.PinData)
		diags.Append(pinDiags...)
		// Check condition.
		if !diags.HasError() {
			plan.PinData = pinDataMap
		}
	} else {
		// Set null map when API returns nil to ensure attribute is known.
		plan.PinData = types.MapNull(types.StringType)
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
func serializeWorkflowJSON(workflow *n8nsdk.Workflow, plan *models.Resource) {
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
	// Normalize settings before serialization to avoid unnecessary diffs.
	normalizedSettings := normalizeWorkflowSettings(workflow.Settings)
	// Check for error.
	if settingsJSON, err := json.Marshal(normalizedSettings); err == nil {
		plan.SettingsJSON = types.StringValue(string(settingsJSON))
	}
}

// normalizeWorkflowSettings removes default values from settings.
// The n8n API returns default values for certain settings even when not explicitly set.
// This function removes callerPolicy and availableInMCP defaults to match user config.
//
// Params:
//   - settings: Original workflow settings from API
//
// Returns:
//   - n8nsdk.WorkflowSettings: Normalized settings without default values
func normalizeWorkflowSettings(settings n8nsdk.WorkflowSettings) n8nsdk.WorkflowSettings {
	// Create a copy to avoid modifying the original.
	normalized := settings

	// Remove callerPolicy if it's the default value.
	if normalized.CallerPolicy != nil && *normalized.CallerPolicy == CALLER_POLICY_DEFAULT {
		normalized.CallerPolicy = nil
	}

	// Remove availableInMCP if it's the default value (false).
	if normalized.AvailableInMCP != nil && !*normalized.AvailableInMCP {
		normalized.AvailableInMCP = nil
	}

	// Return result.
	return normalized
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
func (r *WorkflowResource) handleWorkflowActivation(ctx context.Context, plan, state *models.Resource, diags *diag.Diagnostics) {
	activeChanged := isActivationChanged(plan, state)

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
		// Deactivate workflow when Active is false
		_, httpResp, err = r.client.APIClient.WorkflowAPI.WorkflowsIdDeactivatePost(ctx, plan.ID.ValueString()).Execute()
	}

	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		action := getActivationAction(plan)
		diags.AddError(
			fmt.Sprintf("Error changing workflow activation status to %s", action),
			fmt.Sprintf("Could not %s workflow ID %s: %s\nHTTP Response: %v", action, plan.ID.ValueString(), err.Error(), httpResp),
		)
	}
}

// isActivationChanged checks if the workflow activation status has changed.
//
// Params:
//   - plan: The desired workflow state
//   - state: The current workflow state
//
// Returns:
//   - bool: True if activation status has changed
func isActivationChanged(plan, state *models.Resource) bool {
	// Return boolean result.
	return !plan.Active.IsNull() && !plan.Active.IsUnknown() &&
		!state.Active.IsNull() && !state.Active.IsUnknown() &&
		plan.Active.ValueBool() != state.Active.ValueBool()
}

// getActivationAction returns the activation action string.
//
// Params:
//   - plan: The desired workflow state
//
// Returns:
//   - string: "activate" or "deactivate"
func getActivationAction(plan *models.Resource) string {
	// Check condition.
	if plan.Active.ValueBool() {
		// Return activate action.
		return "activate"
	}
	// Return deactivate action.
	return "deactivate"
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
func (r *WorkflowResource) updateWorkflowTags(ctx context.Context, workflowID string, plan *models.Resource, workflow *n8nsdk.Workflow, diags *diag.Diagnostics) {
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

// createWorkflowViaAPI creates a new workflow via the n8n API.
//
// Params:
//   - ctx: Context for the API call
//   - workflowRequest: The workflow to create
//   - diags: Diagnostics for error reporting
//
// Returns:
//   - *n8nsdk.Workflow: Created workflow if successful, nil otherwise
func (r *WorkflowResource) createWorkflowViaAPI(ctx context.Context, workflowRequest n8nsdk.Workflow, diags *diag.Diagnostics) *n8nsdk.Workflow {
	workflow, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsPost(ctx).
		Workflow(workflowRequest).
		Execute()

	// Check for non-nil HTTP response
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for API error
	if err != nil {
		diags.AddError(
			"Error creating workflow",
			fmt.Sprintf("Could not create workflow, unexpected error: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return nil
	}

	// Return the created workflow
	return workflow
}

// updateWorkflowViaAPI updates an existing workflow via the n8n API.
//
// Params:
//   - ctx: Context for the API call
//   - workflowID: The ID of the workflow to update
//   - workflowRequest: The workflow updates to apply
//   - diags: Diagnostics for error reporting
//
// Returns:
//   - *n8nsdk.Workflow: Updated workflow if successful, nil otherwise
func (r *WorkflowResource) updateWorkflowViaAPI(ctx context.Context, workflowID string, workflowRequest n8nsdk.Workflow, diags *diag.Diagnostics) *n8nsdk.Workflow {
	workflow, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdPut(ctx, workflowID).
		Workflow(workflowRequest).
		Execute()

	// Check for non-nil HTTP response
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for API error
	if err != nil {
		diags.AddError(
			"Error updating workflow",
			fmt.Sprintf("Could not update workflow ID %s: %s\nHTTP Response: %v", workflowID, err.Error(), httpResp),
		)
		return nil
	}

	// Return the updated workflow
	return workflow
}

// transferWorkflowToProject transfers a workflow to a specified project.
//
// Params:
//   - ctx: Context for the API call
//   - workflowID: The workflow ID to transfer
//   - projectID: The destination project ID
//   - diags: Diagnostics for error reporting
//
// Returns:
//   - bool: True if transfer succeeded, false otherwise
func (r *WorkflowResource) transferWorkflowToProject(ctx context.Context, workflowID, projectID string, diags *diag.Diagnostics) bool {
	transferRequest := n8nsdk.WorkflowsIdTransferPutRequest{
		DestinationProjectId: projectID,
	}

	httpResp, err := r.client.APIClient.WorkflowAPI.
		WorkflowsIdTransferPut(ctx, workflowID).
		WorkflowsIdTransferPutRequest(transferRequest).
		Execute()

	// Check for non-nil HTTP response.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		diags.AddError(
			"Error transferring workflow to project",
			fmt.Sprintf("Could not transfer workflow ID %s to project %s: %s\nHTTP Response: %v", workflowID, projectID, err.Error(), httpResp),
		)
		// Return failure.
		return false
	}

	// Return success.
	return true
}

// handlePostCreation handles post-creation workflow operations (ID, tags, project, activation).
// This helper reduces complexity in the create operation by consolidating post-creation steps.
//
// Params:
//   - ctx: Context for the API call
//   - workflow: The created workflow
//   - plan: The resource model
//   - diags: Diagnostics for error reporting
//
// Returns:
//   - *n8nsdk.Workflow: Updated workflow if successful, nil otherwise
func (r *WorkflowResource) handlePostCreation(ctx context.Context, workflow *n8nsdk.Workflow, plan *models.Resource, diags *diag.Diagnostics) *n8nsdk.Workflow {
	// Set ID from created workflow
	plan.ID = types.StringPointerValue(workflow.Id)

	// Update tags if provided
	workflow = r.applyPostCreationTagsAndProject(ctx, workflow, plan, diags)
	// Check for tag/project errors
	if workflow == nil {
		return nil
	}

	// Handle workflow activation after creation if requested
	if !r.handlePostCreationActivation(ctx, plan, workflow, diags) {
		return nil
	}

	// Return the processed workflow with all post-creation operations applied
	return workflow
}

// applyPostCreationTagsAndProject applies tags and project assignment to a newly created workflow.
//
// Params:
//   - ctx: Context for the API call
//   - workflow: The created workflow
//   - plan: The resource model
//   - diags: Diagnostics for error reporting
//
// Returns:
//   - *n8nsdk.Workflow: Updated workflow if successful, nil otherwise
func (r *WorkflowResource) applyPostCreationTagsAndProject(ctx context.Context, workflow *n8nsdk.Workflow, plan *models.Resource, diags *diag.Diagnostics) *n8nsdk.Workflow {
	// Update tags if provided
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() && workflow.Id != nil {
		r.updateWorkflowTags(ctx, *workflow.Id, plan, workflow, diags)
		// Check for tag update errors
		if diags.HasError() {
			return nil
		}
	}

	// Transfer workflow to project if project_id is specified
	if !plan.ProjectID.IsNull() && !plan.ProjectID.IsUnknown() && workflow.Id != nil {
		updatedWorkflow := r.handleProjectAssignment(ctx, *workflow.Id, plan.ProjectID.ValueString(), diags)
		// Check if project assignment succeeded
		if diags.HasError() {
			return nil
		}
		// Use updated workflow if available
		if updatedWorkflow != nil {
			workflow = updatedWorkflow
		}
	}

	// Return the workflow with tags and project applied
	return workflow
}

// handleProjectAssignment handles transferring a workflow to a project and re-fetching the workflow.
// This helper reduces complexity in create and update operations.
//
// Params:
//   - ctx: Context for the API call
//   - workflowID: The workflow ID to transfer
//   - projectID: The destination project ID
//   - diags: Diagnostics for error reporting
//
// Returns:
//   - *n8nsdk.Workflow: Updated workflow if successful, nil otherwise
func (r *WorkflowResource) handleProjectAssignment(ctx context.Context, workflowID, projectID string, diags *diag.Diagnostics) *n8nsdk.Workflow {
	// Transfer workflow to project
	if !r.transferWorkflowToProject(ctx, workflowID, projectID, diags) {
		return nil
	}

	// Re-fetch workflow to get updated project info
	workflow, httpResp, err := r.client.APIClient.WorkflowAPI.
		WorkflowsIdGet(ctx, workflowID).
		Execute()

	// Check for non-nil HTTP response.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Return nil on re-fetch failure; add diagnostic since transfer
	// succeeded but workflow couldn't be re-fetched
	if err != nil {
		diags.AddError(
			"Error re-fetching workflow after project transfer",
			fmt.Sprintf("Workflow was transferred to project successfully, but failed to re-fetch workflow ID %s: %s\nHTTP Response: %v", workflowID, err.Error(), httpResp),
		)
		return nil
	}

	// Return the updated workflow with project information
	return workflow
}
