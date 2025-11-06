package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure CredentialResource implements required interfaces.
var (
	_ resource.Resource                = &CredentialResource{}
	_ resource.ResourceWithConfigure   = &CredentialResource{}
	_ resource.ResourceWithImportState = &CredentialResource{}
)

// CredentialResource defines the resource implementation for n8n credentials.
// Note: Uses rotation strategy for updates (CREATE new, UPDATE workflows, DELETE old).
type CredentialResource struct {
	client *providertypes.N8nClient
}

// CredentialResourceModel describes the resource data model.
type CredentialResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Type      types.String `tfsdk:"type"`
	Data      types.Map    `tfsdk:"data"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// WorkflowBackup stores workflow state for rollback.
type WorkflowBackup struct {
	ID       string
	Original *n8nsdk.Workflow
}

// NewCredentialResource creates a new CredentialResource instance.
func NewCredentialResource() resource.Resource {
	// Return result.
	return &CredentialResource{}
}

// Metadata returns the resource type name.
func (r *CredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential"
}

// Schema defines the schema for the resource.
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
func (r *CredentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*providertypes.N8nClient)
	// Check condition.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providertypes.N8nClient, got: %T", req.ProviderData),
		)
		// Return result.
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *CredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CredentialResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare credential data
	var credData map[string]interface{}
	resp.Diagnostics.Append(plan.Data.ElementsAs(ctx, &credData, false)...)
	// Check condition.
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
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
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
func (r *CredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CredentialResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
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
func (r *CredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state CredentialResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	oldCredID := state.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Starting credential rotation for %s", oldCredID))

	// STEP 1: Create new credential
	// Prepare credential data
	var credData map[string]interface{}
	resp.Diagnostics.Append(plan.Data.ElementsAs(ctx, &credData, false)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	credRequest := n8nsdk.Credential{
		Name: plan.Name.ValueString(),
		Type: plan.Type.ValueString(),
		Data: credData,
	}

	newCred, httpResp, err := r.client.APIClient.CredentialAPI.CredentialsPost(ctx).
		Credential(credRequest).
		Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating new credential during rotation",
			fmt.Sprintf("Could not create new credential: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return result.
		return
	}

	newCredID := newCred.Id
	tflog.Info(ctx, fmt.Sprintf("Created new credential %s", newCredID))

	// STEP 2: Scan workflows using old credential
	workflowList, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsGet(ctx).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		// Rollback: delete new credential
		tflog.Error(ctx, "Failed to list workflows, rolling back")
		_, httpResp2, err := r.client.APIClient.CredentialAPI.DeleteCredential(ctx, newCredID).Execute()
		// Check for non-nil value.
		if httpResp2 != nil && httpResp2.Body != nil {
			defer httpResp2.Body.Close()
		}
		// Check for error.
		if err != nil {
			// Log error but continue - cleanup is best effort
			tflog.Error(ctx, fmt.Sprintf("Failed to delete credential during cleanup: %s", err.Error()))
		}

		resp.Diagnostics.AddError(
			"Error scanning workflows during rotation",
			fmt.Sprintf("Could not list workflows: %s", err.Error()),
		)
		// Return result.
		return
	}

	// Find affected workflows
	affectedWorkflows := []WorkflowBackup{}
	// Check for non-nil value.
	if workflowList.Data != nil {
		// Iterate over items.
		for _, workflow := range workflowList.Data {
			// Check condition.
			if usesCredential(&workflow, oldCredID) {
				affectedWorkflows = append(affectedWorkflows, WorkflowBackup{
					ID:       *workflow.Id,
					Original: &workflow,
				})
			}
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Found %d workflows using credential %s", len(affectedWorkflows), oldCredID))

	// STEP 3: Update each workflow
	updatedWorkflows := []string{}
	// Iterate over items.
	for i, backup := range affectedWorkflows {
		// Throttle to avoid rate limiting
		if i > 0 {
			time.Sleep(100 * time.Millisecond)
		}

		// Get fresh workflow data
		workflow, httpResp, err := r.client.APIClient.WorkflowAPI.
			WorkflowsIdGet(ctx, backup.ID).
			Execute()
		// Check for non-nil value.
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

		// Check for error.
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to get workflow %s, rolling back", backup.ID))
			r.rollbackRotation(ctx, newCredID, affectedWorkflows, updatedWorkflows)

			resp.Diagnostics.AddError(
				"Error reading workflow during rotation",
				fmt.Sprintf("Could not read workflow %s: %s\nRotation rolled back.", backup.ID, err.Error()),
			)
			// Return result.
			return
		}

		// Replace credential references
		updatedWorkflow := replaceCredentialInWorkflow(workflow, oldCredID, newCredID)

		// Update workflow
		_, httpResp, err = r.client.APIClient.WorkflowAPI.
			WorkflowsIdPut(ctx, backup.ID).
			Workflow(*updatedWorkflow).
			Execute()
		// Check for non-nil value.
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

		// Check for error.
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to update workflow %s, rolling back", backup.ID))
			r.rollbackRotation(ctx, newCredID, affectedWorkflows, updatedWorkflows)

			resp.Diagnostics.AddError(
				"Error updating workflow during rotation",
				fmt.Sprintf("Could not update workflow %s: %s\nRotation rolled back.", backup.ID, err.Error()),
			)
			// Return result.
			return
		}

		updatedWorkflows = append(updatedWorkflows, backup.ID)
		tflog.Debug(ctx, fmt.Sprintf("Updated workflow %s (%d/%d)", backup.ID, i+1, len(affectedWorkflows)))
	}

	// STEP 4: Delete old credential
	_, httpResp, err = r.client.APIClient.CredentialAPI.DeleteCredential(ctx, oldCredID).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		// Not critical - new credential works, old is just orphaned
		tflog.Warn(ctx, fmt.Sprintf(
			"Could not delete old credential %s: %s. New credential %s is active. Manual cleanup may be required.",
			oldCredID, err.Error(), newCredID,
		))
		// Handle alternative case.
	} else {
		tflog.Info(ctx, fmt.Sprintf("Deleted old credential %s", oldCredID))
	}

	// STEP 5: Update state
	plan.ID = types.StringValue(newCredID)
	plan.Name = types.StringValue(newCred.Name)
	plan.Type = types.StringValue(newCred.Type)

	// Map timestamps from new credential
	plan.CreatedAt = types.StringValue(newCred.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	plan.UpdatedAt = types.StringValue(newCred.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	tflog.Info(ctx, fmt.Sprintf(
		"Credential rotated successfully: %s → %s (%d workflows updated)",
		oldCredID, newCredID, len(updatedWorkflows),
	))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *CredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CredentialResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.APIClient.CredentialAPI.DeleteCredential(ctx, state.ID.ValueString()).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
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
func (r *CredentialResource) rollbackRotation(
	ctx context.Context,
	newCredID string,
	affectedWorkflows []WorkflowBackup,
	updatedWorkflows []string,
) {
	tflog.Error(ctx, "Rolling back credential rotation")

	// Delete new credential
	_, httpResp, err := r.client.APIClient.CredentialAPI.DeleteCredential(ctx, newCredID).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("CRITICAL: Failed to delete new credential %s during rollback: %s", newCredID, err.Error()))
		// Handle alternative case.
	} else {
		tflog.Info(ctx, fmt.Sprintf("Deleted new credential %s during rollback", newCredID))
	}

	// Restore updated workflows
	restoredCount := 0
	// Iterate over items.
	for _, workflowID := range updatedWorkflows {
		// Find original
		var original *n8nsdk.Workflow
		// Iterate over items.
		for _, backup := range affectedWorkflows {
			// Check condition.
			if backup.ID == workflowID {
				original = backup.Original
				break
			}
		}

		// Check for nil value.
		if original == nil {
			tflog.Error(ctx, fmt.Sprintf("Cannot find original for workflow %s", workflowID))
			continue
		}

		// Restore
		_, httpResp, err := r.client.APIClient.WorkflowAPI.
			WorkflowsIdPut(ctx, workflowID).
			Workflow(*original).
			Execute()
		// Check for non-nil value.
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

		// Check for error.
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to restore workflow %s: %s", workflowID, err.Error()))
			continue
		}

		restoredCount++
		tflog.Debug(ctx, fmt.Sprintf("Restored workflow %s", workflowID))
	}

	tflog.Info(ctx, fmt.Sprintf("Rollback complete: restored %d/%d workflows", restoredCount, len(updatedWorkflows)))
}

// usesCredential checks if a workflow uses a specific credential.
func usesCredential(workflow *n8nsdk.Workflow, credentialID string) bool {
	// Check for nil value.
	if workflow.Nodes == nil {
		// Return result.
		return false
	}

	// Iterate over items.
	for _, node := range workflow.Nodes {
		// Check for non-nil value.
		if node.Credentials != nil {
			// node.Credentials is already map[string]interface{}
			for _, credValue := range node.Credentials {
				// Check condition.
				if credInfo, ok := credValue.(map[string]interface{}); ok {
					// Check condition.
					if id, ok := credInfo["id"].(string); ok && id == credentialID {
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
func replaceCredentialInWorkflow(workflow *n8nsdk.Workflow, oldCredID, newCredID string) *n8nsdk.Workflow {
	// Check for nil value.
	if workflow.Nodes == nil {
		// Return result.
		return workflow
	}

	// Iterate over items.
	for i := range workflow.Nodes {
		node := &workflow.Nodes[i]

		// Check for non-nil value.
		if node.Credentials != nil {
			// node.Credentials is already map[string]interface{}
			for credType, credValue := range node.Credentials {
				// Check condition.
				if credInfo, ok := credValue.(map[string]interface{}); ok {
					// Check condition.
					if id, ok := credInfo["id"].(string); ok && id == oldCredID {
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
