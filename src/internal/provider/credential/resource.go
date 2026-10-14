// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package credential implements the n8n credential resource.
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
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/credential/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
)

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
		MarkdownDescription: "n8n credential resource.\n\n" +
			"**Update Behavior**: Updates are applied in place with `PATCH /credentials/{id}` " +
			"(requires n8n >= 2.3.0). The credential ID is preserved, so workflows referencing it " +
			"keep working without modification. The whole `data` object is replaced on each update.",
		Attributes: r.schemaAttributes(),
	}
}

// schemaAttributes returns the attribute definitions for the credential resource schema.
//
// Returns:
//   - map[string]schema.Attribute: the resource attribute definitions
func (r *CredentialResource) schemaAttributes() map[string]schema.Attribute {
	// Data attribute description with type conversion info.
	dataDesc := "Credential data (secrets, passwords, API keys, etc.). " +
		"String values are automatically converted to the correct type (number, boolean) " +
		"based on the credential schema."
	// Project ID attribute description.
	projectDesc := "Project ID to assign the credential to. " +
		"If not set, credential is created in personal space (General)."
	// Return credential schema attributes.
	return map[string]schema.Attribute{
		"id":         schema.StringAttribute{MarkdownDescription: "Credential identifier", Computed: true},
		"name":       schema.StringAttribute{MarkdownDescription: "Credential name", Required: true},
		"type":       schema.StringAttribute{MarkdownDescription: "Credential type (e.g., httpHeaderAuth)", Required: true},
		"data":       schema.MapAttribute{MarkdownDescription: dataDesc, ElementType: types.StringType, Required: true, Sensitive: true},
		"project_id": schema.StringAttribute{MarkdownDescription: projectDesc, Optional: true, Computed: true},
		"created_at": schema.StringAttribute{MarkdownDescription: "Timestamp when the credential was created", Computed: true},
		"updated_at": schema.StringAttribute{MarkdownDescription: "Timestamp when the credential was last updated", Computed: true},
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
	// Convert string values to proper types based on credential schema.
	convertedData := r.convertDataToSchemaTypes(ctx, plan.Type.ValueString(), credData)

	// Create credential via API
	createResp, httpResp, err := r.executeCreate(ctx, plan.Name.ValueString(), plan.Type.ValueString(), convertedData)
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

	// Handle project assignment if project_id is specified.
	if !plan.ProjectID.IsNull() && !plan.ProjectID.IsUnknown() {
		// Transfer credential to the specified project.
		if !r.transferCredentialToProject(ctx, createResp.Id, plan.ProjectID.ValueString(), &resp.Diagnostics) {
			// Return failure - project assignment failed.
			return false
		}
	}

	// Map project_id to plan (keep the requested value).
	mapCredentialProjectID(plan, plan.ProjectID)

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

// Update updates the credential in place via PATCH /credentials/{id}.
// The credential ID is preserved, so workflows referencing it keep working.
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

// executeUpdateLogic contains the main logic for updating a credential.
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
	credID := state.ID.ValueString()

	// Convert string values to proper types based on credential schema.
	convertedData := r.convertDataToSchemaTypes(ctx, plan.Type.ValueString(), credData)

	// Update credential in place via API
	updateResp, httpResp, err := r.executeUpdate(ctx, credID, plan.Name.ValueString(), plan.Type.ValueString(), convertedData)
	// Close HTTP response body if present.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check if credential update failed.
	if err != nil {
		detail := fmt.Sprintf("Could not update credential ID %s: %s\nHTTP Response: %v", credID, err.Error(), httpResp)
		// PATCH /credentials/{id} only exists since n8n 2.3.0.
		if httpResp != nil && httpResp.StatusCode == http.StatusMethodNotAllowed {
			detail += "\nUpdating credentials requires n8n >= 2.3.0."
		}
		resp.Diagnostics.AddError("Error updating credential", detail)
		// Return failure.
		return false
	}

	// Map response to plan (ID is unchanged)
	r.mapCreateResponseToPlan(plan, updateResp)

	// Move the credential if the requested project changed.
	projectID, ok := r.updateProjectAssignment(ctx, credID, plan.ProjectID, state.ProjectID, &resp.Diagnostics)
	// Check if project assignment succeeded.
	if !ok {
		// Return failure - project assignment failed.
		return false
	}

	// Map project_id to plan.
	mapCredentialProjectID(plan, projectID)

	tflog.Info(ctx, fmt.Sprintf("Updated credential %s", credID))

	// Return success.
	return true
}

// updateProjectAssignment transfers the credential when the requested project differs from the current one.
//
// Params:
//   - ctx: Context for the request
//   - credID: Credential ID
//   - planned: The planned project ID
//   - current: The project ID in state
//   - diags: Diagnostics collector for errors
//
// Returns:
//   - projectID: The project ID to store in state
//   - ok: True if the assignment succeeded, false otherwise
func (r *CredentialResource) updateProjectAssignment(ctx context.Context, credID string, planned, current types.String, diags *diag.Diagnostics) (projectID types.String, ok bool) {
	projectID = planned
	// project_id is unknown when removed from the config: the credential stays in its current project.
	if projectID.IsUnknown() {
		projectID = current
	}

	// Transfer only when the requested project differs from the current one.
	if !projectID.IsNull() && !projectID.Equal(current) {
		// Return transfer result.
		return projectID, r.transferCredentialToProject(ctx, credID, projectID.ValueString(), diags)
	}

	// Return unchanged project.
	return projectID, true
}

// executeUpdate executes the API call to update a credential.
// The whole data object is replaced (isPartialData is left to its false default).
//
// Params:
//   - ctx: Context for the request
//   - id: Credential ID
//   - name: Credential name
//   - credType: Credential type
//   - data: Credential data
//
// Returns:
//   - updated: The updated credential
//   - httpResp: The HTTP response
//   - err: Error if any
func (r *CredentialResource) executeUpdate(ctx context.Context, id, name, credType string, data map[string]any) (updated *n8nsdk.CreateCredentialResponse, httpResp *http.Response, err error) {
	updateRequest := n8nsdk.UpdateCredentialRequest{
		Name: &name,
		Type: &credType,
		Data: data,
	}

	// Execute API call and return result.
	return r.client.APIClient.CredentialAPI.UpdateCredential(ctx, id).
		UpdateCredentialRequest(updateRequest).
		Execute()
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
