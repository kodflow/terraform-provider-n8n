// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package credential provides helper functions for credential operations including rotation and workflow updates.
package credential

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/credential/models"
)

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
			// Check if workflow uses the old credential.
			if usesCredential(&workflow, oldCredID) {
				affectedWorkflows = append(affectedWorkflows, models.WorkflowBackup{
					ID:       *workflow.Id,
					Original: &workflow,
				})
			}
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Found %d workflows using credential %s", len(affectedWorkflows), oldCredID))
	// Return the list of affected workflows and success status.
	return affectedWorkflows, true
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
		// Close response body if present.
		if httpRespGet != nil && httpRespGet.Body != nil {
			defer httpRespGet.Body.Close()
		}

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
		// Close response body if present.
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

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
