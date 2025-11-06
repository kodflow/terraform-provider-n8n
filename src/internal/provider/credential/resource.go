package credential

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)

// RotationThrottleMilliseconds is the delay between workflow updates during credential rotation.
const RotationThrottleMilliseconds int = 100

// Ensure CredentialResource implements required interfaces.
var (
	_ resource.Resource                = &CredentialResource{}
	_ resource.ResourceWithConfigure   = &CredentialResource{}
	_ resource.ResourceWithImportState = &CredentialResource{}
)

// CredentialResource defines the resource implementation for n8n credentials.
// Note: Uses rotation strategy for updates (CREATE new, UPDATE workflows, DELETE old).
type CredentialResource struct {
	client *client.N8nClient
}


// CredentialWorkflowBackup stores workflow state for rollback during credential rotation.
// Captures original workflow data during credential rotation to enable recovery if the operation fails.
type CredentialWorkflowBackup struct {
	ID       string
	Original *n8nsdk.Workflow
}

// NewCredentialResource creates a new CredentialResource instance.
//
// Returns:
//   - resource.Resource: A new CredentialResource instance
func NewCredentialResource() resource.Resource {
	// Return result.
	return &CredentialResource{}
}

// Metadata returns the resource type name.
//
// Params:
//   - ctx: Context for the request
//   - req: Metadata request
//   - resp: Metadata response
func (r *CredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential"
}

// Schema defines the schema for the resource.
//
// Params:
//   - ctx: Context for the request
//   - req: Schema request
//   - resp: Schema response
func (r *CredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "n8n credential resource with automatic rotation on update.\n\n" +
			"**Update Behavior**: When updated, the credential is rotated:\n" +
			"1. New credential is created\n" +
			"2. All workflows using the old credential are updated\n" +
			"3. Old credential is deleted\n" +
			"4. If any step fails, automatic rollback is performed\n\n" +
			"**Note**: The credential ID will change after an update, but this is handled automatically.",

		Attributes: map[string]schema.Attribute{
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
		},
	}
}

// Configure adds the provider configured client to the resource.
//
// Params:
//   - ctx: Context for the request
//   - req: Configure request with provider data
//   - resp: Configure response
func (r *CredentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
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
	var plan *CredentialResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check if plan read succeeded.
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare credential data
	var credData map[string]interface{}
	resp.Diagnostics.Append(plan.Data.ElementsAs(ctx, &credData, false)...)
	// Check if credential data extraction succeeded.
	if resp.Diagnostics.HasError() {
		return
	}

	credRequest := n8nsdk.Credential{
		Name: plan.Name.ValueString(),
		Type: plan.Type.ValueString(),
		Data: credData,
	}

	createResp, httpResp, err := r.client.APIClient.CredentialAPI.CredentialsPost(ctx).
		Credential(credRequest).
		Execute()
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
		// Return result.
		return
	}

	plan.ID = types.StringValue(createResp.Id)
	plan.Name = types.StringValue(createResp.Name)
	plan.Type = types.StringValue(createResp.Type)

	// Note: Data is not returned by the API for security reasons (contains secrets).
	// We keep the data from the plan

	// Map timestamps
	plan.CreatedAt = types.StringValue(createResp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	plan.UpdatedAt = types.StringValue(createResp.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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
	var state *CredentialResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check if state read succeeded.
	if resp.Diagnostics.HasError() {
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
	var plan, state *CredentialResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check if plan and state read succeeded.
	if resp.Diagnostics.HasError() {
		return
	}

	oldCredID := state.ID.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Starting credential rotation for %s", oldCredID))

	// Prepare credential data
	var credData map[string]interface{}
	resp.Diagnostics.Append(plan.Data.ElementsAs(ctx, &credData, false)...)
	// Check if credential data extraction succeeded.
	if resp.Diagnostics.HasError() {
		return
	}

	// STEP 1: Create new credential
	newCred := r.createNewCredential(ctx, plan.Name.ValueString(), plan.Type.ValueString(), credData, &resp.Diagnostics)
	// Check if new credential creation succeeded.
	if resp.Diagnostics.HasError() {
		return
	}

	newCredID := newCred.Id
	tflog.Info(ctx, fmt.Sprintf("Created new credential %s", newCredID))

	// STEP 2: Scan workflows using old credential
	affectedWorkflows, success := r.scanAffectedWorkflows(ctx, oldCredID, newCredID, &resp.Diagnostics)
	// Check if workflow scan succeeded.
	if !success {
		return
	}

	// STEP 3: Update each workflow
	updatedWorkflows, success := r.updateAffectedWorkflows(ctx, affectedWorkflows, oldCredID, newCredID, &resp.Diagnostics)
	// Check if all workflow updates succeeded.
	if !success {
		return
	}

	// STEP 4: Delete old credential
	r.deleteOldCredential(ctx, oldCredID, newCredID)

	// STEP 5: Update state
	plan.ID = types.StringValue(newCredID)
	plan.Name = types.StringValue(newCred.Name)
	plan.Type = types.StringValue(newCred.Type)
	plan.CreatedAt = types.StringValue(newCred.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	plan.UpdatedAt = types.StringValue(newCred.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	tflog.Info(ctx, fmt.Sprintf(
		"Credential rotated successfully: %s → %s (%d workflows updated)",
		oldCredID, newCredID, len(updatedWorkflows),
	))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
//
// Params:
//   - ctx: Context for the request
//   - req: Delete request with state
//   - resp: Delete response
func (r *CredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *CredentialResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check if state read succeeded.
	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.APIClient.CredentialAPI.DeleteCredential(ctx, state.ID.ValueString()).Execute()
	// Close HTTP response body if present.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check if credential deletion failed.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting credential",
			fmt.Sprintf("Could not delete credential ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		// Return result.
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleted credential %s", state.ID.ValueString()))
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
	affectedWorkflows []CredentialWorkflowBackup,
	updatedWorkflows []string,
) {
	tflog.Error(ctx, "Rolling back credential rotation")

	// Delete new credential
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

	// Restore updated workflows
	restoredCount := 0
	// Iterate through each updated workflow to restore it.
	for _, workflowID := range updatedWorkflows {
		// Find original workflow backup.
		var original *n8nsdk.Workflow
		// Iterate through backups to find matching workflow.
		for _, backup := range affectedWorkflows {
			// Check if this backup matches the workflow being restored.
			if backup.ID == workflowID {
				original = backup.Original
				break
			}
		}

		// Check if workflow backup was found.
		if original == nil {
			tflog.Error(ctx, fmt.Sprintf("Cannot find original for workflow %s", workflowID))
			continue
		}

		// Restore workflow to original state.
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
			continue
		}

		restoredCount++
		tflog.Debug(ctx, fmt.Sprintf("Restored workflow %s", workflowID))
	}

	tflog.Info(ctx, fmt.Sprintf("Rollback complete: restored %d/%d workflows", restoredCount, len(updatedWorkflows)))
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
	// Check if workflow nodes are available.
	if workflow.Nodes == nil {
		// Return result.
		return false
	}

	// Iterate through workflow nodes to check for credential usage.
	for _, node := range workflow.Nodes {
		// Check if node has credentials defined.
		if node.Credentials != nil {
			// node.Credentials is already map[string]interface{}
			for _, credValue := range node.Credentials {
				// Check if credential value is a map.
				if credInfo, okMap := credValue.(map[string]interface{}); okMap {
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
			// node.Credentials is already map[string]interface{}
			for credType, credValue := range node.Credentials {
				// Check if credential value is a map.
				if credInfo, okMap := credValue.(map[string]interface{}); okMap {
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
