// Package credential implements the n8n credential resource with rotation support.
package credential

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/credential/models"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)

// ROTATION_THROTTLE_MILLISECONDS is the delay between workflow updates during credential rotation.
const ROTATION_THROTTLE_MILLISECONDS int = 100

// Ensure CredentialResource implements required interfaces.
var (
	_ resource.Resource                = &CredentialResource{}
	_ CredentialResourceInterface      = &CredentialResource{}
	_ resource.ResourceWithConfigure   = &CredentialResource{}
	_ resource.ResourceWithImportState = &CredentialResource{}
)

// CredentialResourceInterface defines the interface for CredentialResource.
type CredentialResourceInterface interface {
	resource.Resource
	Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse)
	Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse)
	Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse)
	Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse)
	Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse)
	Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse)
	Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse)
	ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse)
}

// CredentialResource defines the resource implementation for n8n credentials.
// Note: Uses rotation strategy for updates (CREATE new, UPDATE workflows, DELETE old).
type CredentialResource struct {
	// client is the N8n API client used for credential operations.
	client *client.N8nClient
}

// NewCredentialResource creates a new CredentialResource instance.
//
// Returns:
//   - *CredentialResource: A new CredentialResource instance
func NewCredentialResource() *CredentialResource {
	// Return result.
	return &CredentialResource{}
}

// NewCredentialResourceWrapper creates a new CredentialResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - resource.Resource: the wrapped CredentialResource instance
func NewCredentialResourceWrapper() resource.Resource {
	// Return the wrapped resource instance.
	return NewCredentialResource()
}

// Metadata returns the resource type name.
//
// Params:
//   - ctx: Context for the request
//   - req: Metadata request
//   - resp: Metadata response
func (r *CredentialResource) Metadata(_ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential"
}

// Schema defines the schema for the resource.
//
// Params:
//   - ctx: Context for the request
//   - req: Schema request
//   - resp: Schema response
func (r *CredentialResource) Schema(_ctx context.Context, _req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "n8n credential resource with automatic rotation on update.\n\n" +
			"**Update Behavior**: When updated, the credential is rotated:\n" +
			"1. New credential is created\n" +
			"2. All workflows using the old credential are updated\n" +
			"3. Old credential is deleted\n" +
			"4. If any step fails, automatic rollback is performed\n\n" +
			"**Note**: The credential ID will change after an update, but this is handled automatically.",
		Attributes: r.schemaAttributes(),
	}
}

// schemaAttributes returns the attribute definitions for the credential resource schema.
//
// Returns:
//   - map[string]schema.Attribute: the resource attribute definitions
func (r *CredentialResource) schemaAttributes() map[string]schema.Attribute {
	// Return credential schema attributes.
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "Credential identifier",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "Credential name",
			Required:            true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "Credential type (e.g., httpHeaderAuth, httpBasicAuth)",
			Required:            true,
		},
		"data": schema.MapAttribute{
			MarkdownDescription: "Credential data (secrets, passwords, API keys, etc.)",
			ElementType:         types.StringType,
			Required:            true,
			Sensitive:           true,
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "Timestamp when the credential was created",
			Computed:            true,
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "Timestamp when the credential was last updated",
			Computed:            true,
		},
	}
}

// Configure adds the provider configured client to the resource.
//
// Params:
//   - ctx: Context for the request
//   - req: Configure request with provider data
//   - resp: Configure response
func (r *CredentialResource) Configure(_ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
		// Return result.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	// Check if provider data is the expected client type.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		// Return result.
		return
	}

	r.client = clientData
}

// Create creates the resource and sets the initial Terraform state.
//
// Params:
//   - ctx: Context for the request
//   - req: Create request with plan data
//   - resp: Create response
func (r *CredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check if plan read succeeded.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute create logic
	if !r.executeCreateLogic(ctx, plan, resp) {
		// Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// extractCredentialData extracts credential data from Terraform types.Map.
// This function isolates the ElementsAs conversion which requires Terraform framework context.
//
// Params:
//   - ctx: Context for the request
//   - data: The Terraform types.Map containing credential data
//
// Returns:
//   - map[string]any: The extracted credential data
//   - diag.Diagnostics: Any diagnostics from the conversion
func extractCredentialData(ctx context.Context, data types.Map) (map[string]any, diag.Diagnostics) {
	// First convert to map[string]string since schema ElementType is types.StringType
	var credDataString map[string]string
	var diags diag.Diagnostics
	diags.Append(data.ElementsAs(ctx, &credDataString, false)...)
	// Check if conversion succeeded.
	if diags.HasError() {
		// Return empty map with diagnostics on conversion error.
		return map[string]any{}, diags
	}

	// Convert map[string]string to map[string]any for API compatibility
	credData := make(map[string]any, len(credDataString))
	// Iterate over credential data to convert string values to any.
	for key, value := range credDataString {
		credData[key] = value
	}

	// Return extracted data and diagnostics.
	return credData, diags
}

// executeCreateLogic contains the main logic for creating a credential.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - resp: Create response
//
// Returns:
//   - bool: True if creation succeeded, false otherwise
func (r *CredentialResource) executeCreateLogic(ctx context.Context, plan *models.Resource, resp *resource.CreateResponse) bool {
	// Extract credential data from Terraform types
	credData, diags := extractCredentialData(ctx, plan.Data)
	resp.Diagnostics.Append(diags...)
	// Check if credential data extraction succeeded.
	if resp.Diagnostics.HasError() {
		// Return failure.
		return false
	}

	// Execute creation with extracted data
	return r.executeCreateLogicWithData(ctx, plan, credData, resp)
}

// executeCreateLogicWithData executes the creation logic with pre-extracted credential data.
// This function is fully testable as it takes credData as a parameter, bypassing ElementsAs.
//
// Params:
//   - ctx: Context for the request
//   - plan: The Terraform plan
//   - credData: The credential data (already extracted)
//   - resp: The create response to populate with diagnostics
//
// Returns:
//   - bool: True if creation succeeded, false otherwise
func (r *CredentialResource) executeCreateLogicWithData(ctx context.Context, plan *models.Resource, credData map[string]any, resp *resource.CreateResponse) bool {
	// Create credential via API
	createResp, httpResp, err := r.executeCreate(ctx, plan.Name.ValueString(), plan.Type.ValueString(), credData)
	// Close HTTP response body if present.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check if credential creation failed.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating credential",
			fmt.Sprintf("Could not create credential: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return failure.
		return false
	}

	// Map response to plan
	r.mapCreateResponseToPlan(plan, createResp)

	// Return success.
	return true
}

// executeCreate executes the API call to create a credential.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - name: Credential name
//   - credType: Credential type
//   - data: Credential data
//
// Returns:
//   - *n8nsdk.CreateCredentialResponse: The created credential
//   - *http.Response: The HTTP response
//   - error: Error if any
func (r *CredentialResource) executeCreate(ctx context.Context, name, credType string, data map[string]any) (*n8nsdk.CreateCredentialResponse, *http.Response, error) {
	credRequest := n8nsdk.Credential{
		Name: name,
		Type: credType,
		Data: data,
	}

	// Execute API call and return result.
	return r.client.APIClient.CredentialAPI.CredentialsPost(ctx).
		Credential(credRequest).
		Execute()
}

// mapCreateResponseToPlan maps the API response to the Terraform plan.
// This helper function is separated for testability.
//
// Params:
//   - plan: The Terraform plan to update
//   - createResp: The API response
func (r *CredentialResource) mapCreateResponseToPlan(plan *models.Resource, createResp *n8nsdk.CreateCredentialResponse) {
	plan.ID = types.StringValue(createResp.Id)
	plan.Name = types.StringValue(createResp.Name)
	plan.Type = types.StringValue(createResp.Type)

	// Note: Data is not returned by the API for security reasons (contains secrets).
	// We keep the data from the plan

	// Map timestamps
	plan.CreatedAt = types.StringValue(createResp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	plan.UpdatedAt = types.StringValue(createResp.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
}

// Read refreshes the Terraform state with the latest data.
// WORKAROUND: n8n API doesn't support GET /credentials/{id}.
// We keep the state as-is and warn about drift detection limitations.
//
// Params:
//   - ctx: Context for the request
//   - req: Read request with state
//   - resp: Read response
func (r *CredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *models.Resource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check if state read succeeded.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// WORKAROUND: No API call - n8n doesn't support GET /credentials/{id}.
	tflog.Debug(ctx, fmt.Sprintf(
		"Read credential %s (state-only, no API verification)",
		state.ID.ValueString(),
	))

	// Keep state as-is
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource using rotation strategy.
// WORKAROUND: n8n API doesn't support PUT /credentials/{id}.
// Instead: CREATE new -> UPDATE workflows -> DELETE old.
//
// Params:
//   - ctx: Context for the request
//   - req: Update request with plan and state
//   - resp: Update response
func (r *CredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check if plan and state read succeeded.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute update logic
	if !r.executeUpdateLogic(ctx, plan, state, resp) {
		// Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// executeUpdateLogic contains the main logic for updating a credential using rotation strategy.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - state: The current resource state
//   - resp: Update response
//
// Returns:
//   - bool: True if update succeeded, false otherwise
func (r *CredentialResource) executeUpdateLogic(ctx context.Context, plan, state *models.Resource, resp *resource.UpdateResponse) bool {
	// Extract credential data from Terraform types
	credData, diags := extractCredentialData(ctx, plan.Data)
	resp.Diagnostics.Append(diags...)
	// Check if credential data extraction succeeded.
	if resp.Diagnostics.HasError() {
		// Return failure.
		return false
	}

	// Execute update with extracted data
	return r.executeUpdateLogicWithData(ctx, plan, state, credData, resp)
}

// executeUpdateLogicWithData executes the update logic with pre-extracted credential data.
// This function is fully testable as it takes credData as a parameter, bypassing ElementsAs.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - state: The current resource state
//   - credData: The credential data (already extracted)
//   - resp: Update response
//
// Returns:
//   - bool: True if update succeeded, false otherwise
func (r *CredentialResource) executeUpdateLogicWithData(ctx context.Context, plan, state *models.Resource, credData map[string]any, resp *resource.UpdateResponse) bool {
	oldCredID := state.ID.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Starting credential rotation for %s", oldCredID))

	// STEP 1: Create new credential
	newCred := r.createNewCredential(ctx, plan.Name.ValueString(), plan.Type.ValueString(), credData, &resp.Diagnostics)
	// Check if new credential creation succeeded.
	if resp.Diagnostics.HasError() {
		// Return failure.
		return false
	}

	newCredID := newCred.Id
	tflog.Info(ctx, fmt.Sprintf("Created new credential %s", newCredID))

	// STEP 2: Scan workflows using old credential
	affectedWorkflows, success := r.scanAffectedWorkflows(ctx, oldCredID, newCredID, &resp.Diagnostics)
	// Check if workflow scan succeeded.
	if !success {
		// Return failure.
		return false
	}

	// STEP 3: Update each workflow
	updatedWorkflows, success := r.updateAffectedWorkflows(ctx, affectedWorkflows, oldCredID, newCredID, &resp.Diagnostics)
	// Check if all workflow updates succeeded.
	if !success {
		// Return failure.
		return false
	}

	// STEP 4: Delete old credential
	r.deleteOldCredential(ctx, oldCredID, newCredID)

	// STEP 5: Update plan with new values
	plan.ID = types.StringValue(newCredID)
	plan.Name = types.StringValue(newCred.Name)
	plan.Type = types.StringValue(newCred.Type)
	plan.CreatedAt = types.StringValue(newCred.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	plan.UpdatedAt = types.StringValue(newCred.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	tflog.Info(ctx, fmt.Sprintf(
		"Credential rotated successfully: %s → %s (%d workflows updated)",
		oldCredID, newCredID, len(updatedWorkflows),
	))

	// Return success.
	return true
}

// Delete deletes the resource and removes the Terraform state on success.
//
// Params:
//   - ctx: Context for the request
//   - req: Delete request with state
//   - resp: Delete response
func (r *CredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *models.Resource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check if state read succeeded.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	_, httpResp, err := r.client.APIClient.CredentialAPI.DeleteCredential(ctx, state.ID.ValueString()).Execute()
	// Close HTTP response body if present.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check HTTP status code first - 200/204 means success even if SDK parsing fails.
	if httpResp != nil && (httpResp.StatusCode == http.StatusOK || httpResp.StatusCode == http.StatusNoContent) {
		tflog.Info(ctx, fmt.Sprintf("Deleted credential %s", state.ID.ValueString()))
		// Return success.
		return
	}

	// Check if credential deletion failed.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting credential",
			fmt.Sprintf("Could not delete credential ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
	}
}

// ImportState imports the resource into Terraform state.
//
// Params:
//   - ctx: Context for the request
//   - req: Import state request
//   - resp: Import state response
func (r *CredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.AddWarning(
		"Import limitation",
		"The n8n API does not support reading credentials. "+
			"This import assumes the credential ID is valid. "+
			"If the credential does not exist, operations will fail.",
	)

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// rollbackRotation rolls back a failed rotation.
//
// Params:
//   - ctx: Context for the operation
//   - newCredID: ID of the new credential to delete
//   - affectedWorkflows: List of workflow backups for restoration
//   - updatedWorkflows: List of workflows that were updated
func (r *CredentialResource) rollbackRotation(
	ctx context.Context,
	newCredID string,
	affectedWorkflows []models.WorkflowBackup,
	updatedWorkflows []string,
) {
	tflog.Error(ctx, "Rolling back credential rotation")

	r.deleteNewCredential(ctx, newCredID)
	restoredCount := r.restoreWorkflows(ctx, affectedWorkflows, updatedWorkflows)

	tflog.Info(ctx, fmt.Sprintf("Rollback complete: restored %d/%d workflows", restoredCount, len(updatedWorkflows)))
}

// deleteNewCredential deletes the newly created credential during rollback.
//
// Params:
//   - ctx: Context for the operation
//   - newCredID: ID of the new credential to delete
func (r *CredentialResource) deleteNewCredential(ctx context.Context, newCredID string) {
	_, httpResp, err := r.client.APIClient.CredentialAPI.DeleteCredential(ctx, newCredID).Execute()
	// Close HTTP response body if present.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check if credential deletion failed during rollback.
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("CRITICAL: Failed to delete new credential %s during rollback: %s", newCredID, err.Error()))
		// Handle alternative case.
	} else {
		tflog.Info(ctx, fmt.Sprintf("Deleted new credential %s during rollback", newCredID))
	}
}

// restoreWorkflows restores workflows to their original state during rollback.
//
// Params:
//   - ctx: Context for the operation
//   - affectedWorkflows: List of workflow backups for restoration
//   - updatedWorkflows: List of workflows that were updated
//
// Returns:
//   - restoredCount: Number of workflows successfully restored
func (r *CredentialResource) restoreWorkflows(
	ctx context.Context,
	affectedWorkflows []models.WorkflowBackup,
	updatedWorkflows []string,
) int {
	restoredCount := 0
	// Iterate through each updated workflow to restore it.
	for _, workflowID := range updatedWorkflows {
		original := r.findWorkflowBackup(affectedWorkflows, workflowID)
		// Check if workflow backup was found.
		if original == nil {
			tflog.Error(ctx, fmt.Sprintf("Cannot find original for workflow %s", workflowID))
			continue
		}

		// Restore workflow to original state.
		if r.restoreWorkflow(ctx, workflowID, original) {
			restoredCount++
		}
	}
	// Return count of successfully restored workflows.
	return restoredCount
}

// findWorkflowBackup finds the original workflow from backups.
//
// Params:
//   - affectedWorkflows: List of workflow backups
//   - workflowID: ID of the workflow to find
//
// Returns:
//   - original: Pointer to the original workflow, nil if not found
func (r *CredentialResource) findWorkflowBackup(
	affectedWorkflows []models.WorkflowBackup,
	workflowID string,
) *n8nsdk.Workflow {
	// Iterate through backups to find matching workflow.
	for _, backup := range affectedWorkflows {
		// Check if this backup matches the workflow being restored.
		if backup.ID == workflowID {
			// Return the original workflow from backup.
			return backup.Original
		}
	}
	// Return nil if no backup found.
	return nil
}

// restoreWorkflow restores a single workflow to its original state.
//
// Params:
//   - ctx: Context for the operation
//   - workflowID: ID of the workflow to restore
//   - original: Original workflow state
//
// Returns:
//   - success: True if workflow was restored successfully
func (r *CredentialResource) restoreWorkflow(
	ctx context.Context,
	workflowID string,
	original *n8nsdk.Workflow,
) bool {
	_, httpRespRestore, errRestore := r.client.APIClient.WorkflowAPI.
		WorkflowsIdPut(ctx, workflowID).
		Workflow(*original).
		Execute()
	// Close HTTP response body if present.
	if httpRespRestore != nil && httpRespRestore.Body != nil {
		defer httpRespRestore.Body.Close()
	}

	// Check if workflow restoration failed.
	if errRestore != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to restore workflow %s: %s", workflowID, errRestore.Error()))
		// Return failure status.
		return false
	}

	tflog.Debug(ctx, fmt.Sprintf("Restored workflow %s", workflowID))
	// Return success status.
	return true
}

// usesCredential checks if a workflow uses a specific credential.
//
// Params:
//   - workflow: The workflow to check
//   - credentialID: The credential ID to search for
//
// Returns:
//   - bool: True if workflow uses the credential, false otherwise
func usesCredential(workflow *n8nsdk.Workflow, credentialID string) bool {
	// Check if workflow is nil first.
	if workflow == nil {
		return false
	}

	// Check if workflow nodes are available.
	if workflow.Nodes == nil {
		// Return result.
		return false
	}

	// Iterate through workflow nodes to check for credential usage.
	for _, node := range workflow.Nodes {
		// Check if node has credentials defined.
		if node.Credentials != nil {
			// node.Credentials is already map[string]any
			for _, credValue := range node.Credentials {
				// Check if credential value is a map.
				if credInfo, okMap := credValue.(map[string]any); okMap {
					// Check if credential ID matches the target credential.
					if id, okID := credInfo["id"].(string); okID && id == credentialID {
						// Return result.
						return true
					}
				}
			}
		}
	}

	// Return result.
	return false
}

// replaceCredentialInWorkflow replaces all references to oldCredID with newCredID.
//
// Params:
//   - workflow: The workflow to modify
//   - oldCredID: The old credential ID to replace
//   - newCredID: The new credential ID to use
//
// Returns:
//   - *n8nsdk.Workflow: The modified workflow
func replaceCredentialInWorkflow(workflow *n8nsdk.Workflow, oldCredID, newCredID string) *n8nsdk.Workflow {
	// Check if workflow is nil first.
	if workflow == nil {
		return nil
	}

	// Check if workflow nodes are available.
	if workflow.Nodes == nil {
		// Return result.
		return workflow
	}

	// Iterate through workflow nodes to replace credentials.
	for i := range workflow.Nodes {
		node := &workflow.Nodes[i]

		// Check if node has credentials defined.
		if node.Credentials != nil {
			// node.Credentials is already map[string]any
			for credType, credValue := range node.Credentials {
				// Check if credential value is a map.
				if credInfo, okMap := credValue.(map[string]any); okMap {
					// Check if credential ID matches the old credential.
					if id, okID := credInfo["id"].(string); okID && id == oldCredID {
						credInfo["id"] = newCredID
						node.Credentials[credType] = credInfo
					}
				}
			}
		}
	}

	// Return result.
	return workflow
}
