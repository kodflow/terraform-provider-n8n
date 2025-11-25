// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package credential provides helper functions for credential operations including rotation and workflow updates.
package credential

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/credential/models"
)

// float64BitSize is the bit size for float64 parsing.
const float64BitSize = 64

// createNewCredential creates a new credential via the API.
// Returns the new credential response or adds an error to diagnostics.
//
// Params:
//   - ctx: Context for the API call
//   - name: Name of the credential
//   - credType: Type of the credential
//   - credData: Credential data as key-value pairs
//   - diags: Diagnostics collector for errors
//
// Returns:
//   - *n8nsdk.CreateCredentialResponse: The created credential or nil on error
func (r *CredentialResource) createNewCredential(ctx context.Context, name, credType string, credData map[string]any, diags *diag.Diagnostics) *n8nsdk.CreateCredentialResponse {
	credRequest := n8nsdk.Credential{
		Name: name,
		Type: credType,
		Data: credData,
	}

	newCred, httpResp, err := r.client.APIClient.CredentialAPI.CredentialsPost(ctx).
		Credential(credRequest).
		Execute()
	// Close response body if present.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error during credential creation.
	if err != nil {
		diags.AddError(
			"Error creating new credential during rotation",
			fmt.Sprintf("Could not create new credential: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return nil to indicate failure.
		return nil
	}

	// Return the created credential.
	return newCred
}

// scanAffectedWorkflows lists all workflows and finds those using the old credential.
// Returns a list of workflow backups and whether the scan succeeded.
//
// Params:
//   - ctx: Context for the API call
//   - oldCredID: ID of the old credential to search for
//   - newCredID: ID of the new credential (for rollback if needed)
//   - diags: Diagnostics collector for errors
//
// Returns:
//   - []models.WorkflowBackup: List of workflows using the old credential
//   - bool: True if scan succeeded, false otherwise
func (r *CredentialResource) scanAffectedWorkflows(ctx context.Context, oldCredID, newCredID string, diags *diag.Diagnostics) ([]models.WorkflowBackup, bool) {
	workflowList, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsGet(ctx).Execute()
	// Close response body if present.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error listing workflows.
	if err != nil {
		// Rollback: delete new credential
		tflog.Error(ctx, "Failed to list workflows, rolling back")
		r.deleteCredentialBestEffort(ctx, newCredID)

		diags.AddError(
			"Error scanning workflows during rotation",
			fmt.Sprintf("Could not list workflows: %s", err.Error()),
		)
		// Return empty slice and failure status.
		return []models.WorkflowBackup{}, false
	}

	// Find affected workflows
	affectedWorkflows := []models.WorkflowBackup{}
	// Check if workflow data is available.
	if workflowList.Data != nil {
		// Iterate through workflows to find those using the credential.
		for _, workflow := range workflowList.Data {
			// Create a copy to avoid loop pointer aliasing
			wf := workflow
			// Check if workflow uses the old credential.
			if usesCredential(&wf, oldCredID) {
				affectedWorkflows = append(affectedWorkflows, models.WorkflowBackup{
					ID:       *wf.Id,
					Original: &wf,
				})
			}
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Found %d workflows using credential %s", len(affectedWorkflows), oldCredID))
	// Return the list of affected workflows and success status.
	return affectedWorkflows, true
}

// closeResponseBody closes an HTTP response body and logs any errors.
// This is a helper to avoid defer in loops and handle close errors properly.
//
// Params:
//   - ctx: Context for logging
//   - resp: HTTP response to close (can be nil)
func closeResponseBody(ctx context.Context, resp *http.Response) {
	// Only attempt to close if response and body are not nil.
	if resp != nil && resp.Body != nil {
		// Log warning if response body fails to close.
		if closeErr := resp.Body.Close(); closeErr != nil {
			tflog.Warn(ctx, fmt.Sprintf("Failed to close response body: %s", closeErr.Error()))
		}
	}
}

// updateAffectedWorkflows updates all workflows to use the new credential.
// Returns a list of successfully updated workflow IDs and whether all updates succeeded.
//
// Params:
//   - ctx: Context for the API call
//   - affectedWorkflows: List of workflow backups to update
//   - oldCredID: ID of the old credential to replace
//   - newCredID: ID of the new credential to use
//   - diags: Diagnostics collector for errors
//
// Returns:
//   - []string: List of successfully updated workflow IDs
//   - bool: True if all updates succeeded, false otherwise
func (r *CredentialResource) updateAffectedWorkflows(ctx context.Context, affectedWorkflows []models.WorkflowBackup, oldCredID, newCredID string, diags *diag.Diagnostics) ([]string, bool) {
	updatedWorkflows := []string{}

	// Iterate through all affected workflows.
	for i, backup := range affectedWorkflows {
		// Throttle to avoid rate limiting
		// Skip sleep for first iteration.
		if i > 0 {
			time.Sleep(time.Duration(ROTATION_THROTTLE_MILLISECONDS) * time.Millisecond)
		}

		// Get fresh workflow data
		workflow, httpRespGet, errGet := r.client.APIClient.WorkflowAPI.
			WorkflowsIdGet(ctx, backup.ID).
			Execute()
		closeResponseBody(ctx, httpRespGet)

		// Check for error reading workflow.
		if errGet != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to get workflow %s, rolling back", backup.ID))
			r.rollbackRotation(ctx, newCredID, affectedWorkflows, updatedWorkflows)

			diags.AddError(
				"Error reading workflow during rotation",
				fmt.Sprintf("Could not read workflow %s: %s\nRotation rolled back.", backup.ID, errGet.Error()),
			)
			// Return partial results and failure status.
			return updatedWorkflows, false
		}

		// Replace credential references
		updatedWorkflow := replaceCredentialInWorkflow(workflow, oldCredID, newCredID)

		// Update workflow
		_, httpResp, err := r.client.APIClient.WorkflowAPI.
			WorkflowsIdPut(ctx, backup.ID).
			Workflow(*updatedWorkflow).
			Execute()
		closeResponseBody(ctx, httpResp)

		// Check for error updating workflow.
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to update workflow %s, rolling back", backup.ID))
			r.rollbackRotation(ctx, newCredID, affectedWorkflows, updatedWorkflows)

			diags.AddError(
				"Error updating workflow during rotation",
				fmt.Sprintf("Could not update workflow %s: %s\nRotation rolled back.", backup.ID, err.Error()),
			)
			// Return partial results and failure status.
			return updatedWorkflows, false
		}

		updatedWorkflows = append(updatedWorkflows, backup.ID)
		tflog.Debug(ctx, fmt.Sprintf("Updated workflow %s (%d/%d)", backup.ID, i+1, len(affectedWorkflows)))
	}

	// Return all updated workflows and success status.
	return updatedWorkflows, true
}

// deleteCredentialBestEffort attempts to delete a credential but does not fail if unsuccessful.
// Used for cleanup operations where failure is not critical.
//
// Params:
//   - ctx: Context for the API call
//   - credID: ID of the credential to delete
func (r *CredentialResource) deleteCredentialBestEffort(ctx context.Context, credID string) {
	_, httpResp, err := r.client.APIClient.CredentialAPI.DeleteCredential(ctx, credID).Execute()
	// Close response body if present.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Log error if deletion failed, but don't fail the operation.
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to delete credential during cleanup: %s", err.Error()))
	}
}

// deleteOldCredential deletes the old credential after rotation.
// Logs a warning if deletion fails but does not fail the rotation.
//
// Params:
//   - ctx: Context for the API call
//   - oldCredID: ID of the old credential to delete
//   - newCredID: ID of the new credential (for logging)
func (r *CredentialResource) deleteOldCredential(ctx context.Context, oldCredID, newCredID string) {
	_, httpResp, err := r.client.APIClient.CredentialAPI.DeleteCredential(ctx, oldCredID).Execute()
	// Close response body if present.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check if deletion failed.
	if err != nil {
		// Not critical - new credential works, old is just orphaned
		tflog.Warn(ctx, fmt.Sprintf(
			"Could not delete old credential %s: %s. New credential %s is active. Manual cleanup may be required.",
			oldCredID, err.Error(), newCredID,
		))
	} else {
		// Log success.
		tflog.Info(ctx, fmt.Sprintf("Deleted old credential %s", oldCredID))
	}
}

// convertDataToSchemaTypes converts string values in credential data to their proper types
// based on the credential schema from n8n API.
// This enables support for credential types that require number or boolean fields.
//
// Params:
//   - ctx: Context for the API call
//   - credType: The credential type name (e.g., "imap", "mongoDb")
//   - data: The credential data with string values
//
// Returns:
//   - map[string]any: The credential data with converted types
func (r *CredentialResource) convertDataToSchemaTypes(ctx context.Context, credType string, data map[string]any) map[string]any {
	// Fetch credential schema from n8n API.
	schema, httpResp, err := r.client.APIClient.CredentialAPI.
		CredentialsSchemaCredentialTypeNameGet(ctx, credType).
		Execute()
	// Close response body if present.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// If schema fetch fails, return original data with warning.
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf(
			"Could not fetch credential schema for type '%s': %s. Using string values as-is.",
			credType, err.Error(),
		))
		// Return original data without type conversion.
		return data
	}

	// Convert data based on schema properties.
	return r.applySchemaTypeConversions(ctx, schema, data)
}

// applySchemaTypeConversions applies type conversions based on the schema.
//
// Params:
//   - ctx: Context for logging
//   - schema: The credential schema from n8n API
//   - data: The credential data with string values
//
// Returns:
//   - map[string]any: The credential data with converted types
func (r *CredentialResource) applySchemaTypeConversions(ctx context.Context, schema map[string]any, data map[string]any) map[string]any {
	// Extract properties from schema.
	properties, ok := schema["properties"].(map[string]any)
	// If no properties found, return original data.
	if !ok {
		tflog.Debug(ctx, "No properties found in credential schema, using original data")
		// Return original data.
		return data
	}

	// Create result map with converted values.
	result := make(map[string]any, len(data))
	// Iterate over data keys to convert each value.
	for key, value := range data {
		// Only convert string values.
		strValue, isString := value.(string)
		// If not a string, keep original value.
		if !isString {
			result[key] = value
			continue
		}

		// Get property schema for this key.
		propSchema, propExists := properties[key].(map[string]any)
		// If no property schema, keep original string value.
		if !propExists {
			result[key] = value
			continue
		}

		// Convert based on property type.
		result[key] = convertValueByType(ctx, key, strValue, propSchema)
	}

	// Return converted data.
	return result
}

// convertValueByType converts a string value to the appropriate type based on schema.
//
// Params:
//   - ctx: Context for logging
//   - key: The property key (for logging)
//   - value: The string value to convert
//   - propSchema: The property schema containing type information
//
// Returns:
//   - any: The converted value (or original string if conversion fails)
func convertValueByType(ctx context.Context, key, value string, propSchema map[string]any) any {
	// Get the type from property schema.
	propType, hasType := propSchema["type"].(string)
	// If no type specified, return original value.
	if !hasType {
		// Return original string value.
		return value
	}

	// Convert based on type.
	switch propType {
	// Handle numeric types (number, integer) by converting string to float64.
	case "number", "integer":
		// Try to parse as float64 (handles both int and float).
		if f, err := strconv.ParseFloat(value, float64BitSize); err == nil {
			tflog.Debug(ctx, fmt.Sprintf("Converted '%s' from string to number: %v", key, f))
			// Return converted number.
			return f
		}
		tflog.Debug(ctx, fmt.Sprintf("Could not convert '%s' value '%s' to number, keeping as string", key, value))
	// Handle boolean type by converting string to bool.
	case "boolean":
		// Try to parse as boolean.
		if b, err := strconv.ParseBool(value); err == nil {
			tflog.Debug(ctx, fmt.Sprintf("Converted '%s' from string to boolean: %v", key, b))
			// Return converted boolean.
			return b
		}
		tflog.Debug(ctx, fmt.Sprintf("Could not convert '%s' value '%s' to boolean, keeping as string", key, value))
	// Handle other types (string, array, object, etc.) - keep as string.
	default:
	}

	// Return original string value for other types or conversion failures.
	return value
}

// transferCredentialToProject transfers a credential to a specified project.
//
// Params:
//   - ctx: Context for the API call
//   - credentialID: The credential ID to transfer
//   - projectID: The destination project ID
//   - diags: Diagnostics collector for errors
//
// Returns:
//   - bool: True if transfer succeeded, false otherwise
func (r *CredentialResource) transferCredentialToProject(ctx context.Context, credentialID, projectID string, diags *diag.Diagnostics) bool {
	transferRequest := n8nsdk.CredentialsIdTransferPutRequest{
		DestinationProjectId: projectID,
	}

	httpResp, err := r.client.APIClient.CredentialAPI.
		CredentialsIdTransferPut(ctx, credentialID).
		CredentialsIdTransferPutRequest(transferRequest).
		Execute()

	// Close response body if present.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		diags.AddError(
			"Error transferring credential to project",
			fmt.Sprintf("Could not transfer credential ID %s to project %s: %s\nHTTP Response: %v", credentialID, projectID, err.Error(), httpResp),
		)
		// Return failure.
		return false
	}

	tflog.Info(ctx, fmt.Sprintf("Transferred credential %s to project %s", credentialID, projectID))
	// Return success.
	return true
}

// mapCredentialProjectID maps the project_id to the model.
// Since n8n API doesn't return project info in credential responses,
// we keep the value from the plan.
//
// Params:
//   - plan: The resource model to update
//   - requestedProjectID: The project ID that was requested
func mapCredentialProjectID(plan *models.Resource, requestedProjectID types.String) {
	// If project_id was set in plan, keep it.
	if !requestedProjectID.IsNull() && !requestedProjectID.IsUnknown() {
		plan.ProjectID = requestedProjectID
	} else {
		// Set to null if not specified.
		plan.ProjectID = types.StringNull()
	}
}
