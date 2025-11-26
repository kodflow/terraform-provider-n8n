// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package workflow implements workflow management resources and data sources.
package workflow

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/models"
)

// WORKFLOW_ATTRIBUTES_SIZE defines the initial capacity for workflow attributes map.
const WORKFLOW_ATTRIBUTES_SIZE int = 15

// Ensure WorkflowResource implements required interfaces.
var (
	_ resource.Resource                = &WorkflowResource{}
	_ WorkflowResourceInterface        = &WorkflowResource{}
	_ resource.ResourceWithConfigure   = &WorkflowResource{}
	_ resource.ResourceWithImportState = &WorkflowResource{}
)

// WorkflowResource defines the resource implementation for n8n workflows.
// Terraform resource that manages CRUD operations for n8n workflows via the n8n API.
//
// Params:
//   - client: The n8n API client for communicating with n8n server
//
// Returns:
//   - None: This is a struct, not a function
//
// WorkflowResourceInterface defines the interface for WorkflowResource.
type WorkflowResourceInterface interface {
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

// WorkflowResource defines the resource implementation for workflows.
// Terraform resource that manages CRUD operations for n8n workflows via the n8n API.
// It handles workflow lifecycle including creation, updates, deletion, and import operations.
type WorkflowResource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewWorkflowResource creates a new WorkflowResource instance.
//
// Returns:
//   - resource.Resource: A new WorkflowResource instance
func NewWorkflowResource() *WorkflowResource {
	// Return result.
	return &WorkflowResource{}
}

// NewWorkflowResourceWrapper creates a new WorkflowResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - resource.Resource: the wrapped WorkflowResource instance
func NewWorkflowResourceWrapper() resource.Resource {
	// Return the wrapped resource instance.
	return NewWorkflowResource()
}

// Metadata returns the resource type name.
//
// Params:
//   - ctx: Context for the operation
//   - req: Metadata request containing provider type name
//   - resp: Metadata response to populate with type name
//
// Returns:
//   - None: Updates resp in-place
func (r *WorkflowResource) Metadata(_ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

// Schema defines the schema for the resource.
//
// Params:
//   - ctx: Context for the operation
//   - req: Schema request
//   - resp: Schema response to populate with resource schema
//
// Returns:
//   - None: Updates resp in-place
func (r *WorkflowResource) Schema(_ctx context.Context, _req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "n8n workflow resource using generated SDK",
		Attributes:          r.schemaAttributes(),
	}
}

// schemaAttributes returns the attribute definitions for the workflow resource schema.
//
// Returns:
//   - map[string]schema.Attribute: the resource attribute definitions
func (r *WorkflowResource) schemaAttributes() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute, WORKFLOW_ATTRIBUTES_SIZE)

	r.addCoreAttributes(attrs)
	r.addJSONAttributes(attrs)
	r.addMetadataAttributes(attrs)

	// Return schema attributes.
	return attrs
}

// addCoreAttributes adds the core workflow attributes to the schema.
//
// Params:
//   - attrs: attribute map to populate
func (r *WorkflowResource) addCoreAttributes(attrs map[string]schema.Attribute) {
	attrs["id"] = schema.StringAttribute{
		MarkdownDescription: "Workflow identifier",
		Computed:            true,
	}
	attrs["name"] = schema.StringAttribute{
		MarkdownDescription: "Workflow name",
		Required:            true,
	}
	attrs["active"] = schema.BoolAttribute{
		MarkdownDescription: "Whether the workflow is active",
		Optional:            true,
		Computed:            true,
	}
	attrs["tags"] = schema.SetAttribute{
		MarkdownDescription: "Set of tag IDs associated with this workflow",
		ElementType:         types.StringType,
		Optional:            true,
	}
	attrs["project_id"] = schema.StringAttribute{
		MarkdownDescription: "Project ID where the workflow should be created. If not specified, workflow is created in the default 'Overview' location. The workflow can be transferred to a different project by updating this value. Note: Once assigned to a project, a workflow cannot be moved back to the Overview location due to n8n API limitations.",
		Optional:            true,
		Computed:            true,
	}
}

// addJSONAttributes adds the JSON-based workflow attributes to the schema.
//
// Params:
//   - attrs: attribute map to populate
func (r *WorkflowResource) addJSONAttributes(attrs map[string]schema.Attribute) {
	attrs["nodes_json"] = schema.StringAttribute{
		MarkdownDescription: "Workflow nodes as JSON string. Must be valid JSON array of node objects.",
		Optional:            true,
		Computed:            true,
	}
	attrs["connections_json"] = schema.StringAttribute{
		MarkdownDescription: "Workflow connections as JSON string. Must be valid JSON object mapping node connections.",
		Optional:            true,
		Computed:            true,
	}
	attrs["settings_json"] = schema.StringAttribute{
		MarkdownDescription: "Workflow settings as JSON string. Must be valid JSON object.",
		Optional:            true,
		Computed:            true,
	}
}

// addMetadataAttributes adds the metadata workflow attributes to the schema.
//
// Params:
//   - attrs: attribute map to populate
func (r *WorkflowResource) addMetadataAttributes(attrs map[string]schema.Attribute) {
	attrs["created_at"] = schema.StringAttribute{
		MarkdownDescription: "Timestamp when the workflow was created",
		Computed:            true,
	}
	attrs["updated_at"] = schema.StringAttribute{
		MarkdownDescription: "Timestamp when the workflow was last updated",
		Computed:            true,
	}
	attrs["version_id"] = schema.StringAttribute{
		MarkdownDescription: "Version identifier of the workflow",
		Computed:            true,
	}
	attrs["is_archived"] = schema.BoolAttribute{
		MarkdownDescription: "Whether the workflow is archived",
		Computed:            true,
	}
	attrs["trigger_count"] = schema.Int64Attribute{
		MarkdownDescription: "Number of triggers in the workflow",
		Computed:            true,
	}
	attrs["meta"] = schema.MapAttribute{
		MarkdownDescription: "Workflow metadata",
		ElementType:         types.StringType,
		Computed:            true,
	}
	attrs["pin_data"] = schema.MapAttribute{
		MarkdownDescription: "Pinned test data for the workflow",
		ElementType:         types.StringType,
		Computed:            true,
	}
}

// Configure adds the provider configured client to the resource.
//
// Params:
//   - ctx: Context for the operation
//   - req: Configure request containing provider data
//   - resp: Configure response for error handling
//
// Returns:
//   - None: Updates resource in-place, populates resp with errors if any
func (r *WorkflowResource) Configure(_ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		// Return with error.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	// Check provider data type.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		// Return result.
		return
	}

	r.client = clientData
}

// Create creates the resource and sets the initial Terraform state.
//
// Params:
//   - ctx: Context for the operation
//   - req: Create request containing plan data
//   - resp: Create response for state and error handling
//
// Returns:
//   - None: Updates resp with state and error diagnostics
func (r *WorkflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for plan parsing errors.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute create logic
	if !r.executeCreateLogic(ctx, plan, resp) {
		// Return with error.
		return
	}

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// executeCreateLogic contains the main logic for creating a workflow.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - resp: Create response
//
// Returns:
//   - bool: True if creation succeeded, false otherwise
func (r *WorkflowResource) executeCreateLogic(ctx context.Context, plan *models.Resource, resp *resource.CreateResponse) bool {
	// Parse JSON fields.
	nodes, connections, settings := parseWorkflowJSON(plan, &resp.Diagnostics)
	// Check for JSON parsing errors.
	if resp.Diagnostics.HasError() {
		// Return failure.
		return false
	}

	// Create workflow (Note: active field is read-only during creation)
	workflowRequest := n8nsdk.Workflow{
		Name:        plan.Name.ValueString(),
		Nodes:       nodes,
		Connections: connections,
		Settings:    settings,
	}

	workflow := r.createWorkflowViaAPI(ctx, workflowRequest, &resp.Diagnostics)
	// Check for API error
	if resp.Diagnostics.HasError() {
		return false
	}

	// Handle post-creation operations (ID, tags, project, activation)
	workflow = r.handlePostCreation(ctx, workflow, plan, &resp.Diagnostics)
	// Check for post-creation errors
	if resp.Diagnostics.HasError() || workflow == nil {
		return false
	}

	// Map workflow state to model.
	// Note: plan.Active retains value from plan or activation result.
	mapWorkflowToModel(ctx, workflow, plan, &resp.Diagnostics)

	// Return success.
	return true
}

// handlePostCreationActivation handles workflow activation after creation.
//
// Params:
//   - ctx: context for the operation
//   - plan: planned resource state
//   - workflow: created workflow from API
//   - diags: diagnostics for error reporting
//
// Returns:
//   - bool: true if successful or activation not needed, false on error
func (r *WorkflowResource) handlePostCreationActivation(ctx context.Context, plan *models.Resource, workflow *n8nsdk.Workflow, diags *diag.Diagnostics) bool {
	// Check if user wants workflow activated
	wantsActive := !plan.Active.IsNull() && !plan.Active.IsUnknown() && plan.Active.ValueBool()
	// Activation not requested
	if !wantsActive {
		// Return success
		return true
	}

	// Check workflow has nodes before activation
	hasNodes := len(workflow.Nodes) > 0
	// Cannot activate workflow without nodes
	if !hasNodes {
		diags.AddError(
			"Cannot activate empty workflow",
			"Cannot create workflow with active=true because it has no nodes. Either set active=false or add at least one node (trigger, poller, or webhook) to the workflow before activating.",
		)
		// Return failure
		return false
	}

	// Prepare state for activation
	currentState := &models.Resource{
		ID:     plan.ID,
		Active: types.BoolPointerValue(workflow.Active),
	}

	// Activate the workflow
	r.handleWorkflowActivation(ctx, plan, currentState, diags)

	// Update workflow object to reflect successful activation
	// This ensures mapWorkflowToModel will use the correct active status
	if !diags.HasError() {
		activated := true
		workflow.Active = &activated
	}

	// Return success if no errors occurred
	return !diags.HasError()
}

// Read refreshes the Terraform state with the latest data.
//
// Params:
//   - ctx: Context for the operation
//   - req: Read request containing state data
//   - resp: Read response for state and error handling
//
// Returns:
//   - None: Updates resp with refreshed state and error diagnostics
func (r *WorkflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *models.Resource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check for state parsing errors.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute read logic
	if !r.executeReadLogic(ctx, state, resp) {
		// Return with error.
		return
	}

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// executeReadLogic contains the main logic for reading a workflow.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Read response
//
// Returns:
//   - bool: True if read succeeded, false otherwise
func (r *WorkflowResource) executeReadLogic(ctx context.Context, state *models.Resource, resp *resource.ReadResponse) bool {
	// Get workflow from SDK.
	workflow, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdGet(ctx, state.ID.ValueString()).Execute()

	// Check for non-nil HTTP response.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for API error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading workflow",
			fmt.Sprintf("Could not read workflow ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		// Return failure.
		return false
	}

	// Map response to state.
	mapWorkflowToModel(ctx, workflow, state, &resp.Diagnostics)

	// Return success.
	return true
}

// Update updates the resource and sets the updated Terraform state on success.
//
// Params:
//   - ctx: Context for the operation
//   - req: Update request containing plan and state data
//   - resp: Update response for state and error handling
//
// Returns:
//   - None: Updates resp with updated state and error diagnostics
func (r *WorkflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *models.Resource

	// Read Terraform plan and state data.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check for plan/state parsing errors.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute update logic
	if !r.executeUpdateLogic(ctx, plan, state, resp) {
		// Return with error.
		return
	}

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// executeUpdateLogic contains the main logic for updating a workflow.
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
func (r *WorkflowResource) executeUpdateLogic(ctx context.Context, plan, state *models.Resource, resp *resource.UpdateResponse) bool {
	// Parse JSON fields.
	nodes, connections, settings := parseWorkflowJSON(plan, &resp.Diagnostics)
	// Check for JSON parsing errors.
	if resp.Diagnostics.HasError() {
		// Return failure.
		return false
	}

	// Use state.ID for the workflow ID since plan.ID may be Unknown
	// for Computed attributes. The ID cannot change during an update,
	// so state.ID is always the correct value.
	workflowID := state.ID.ValueString()

	// Copy ID from state to plan for consistency.
	plan.ID = state.ID

	// Handle activation change.
	r.handleWorkflowActivation(ctx, plan, state, &resp.Diagnostics)
	// Check for activation change errors.
	if resp.Diagnostics.HasError() {
		// Return failure.
		return false
	}

	// Update workflow content.
	workflowRequest := n8nsdk.Workflow{
		Name:        plan.Name.ValueString(),
		Nodes:       nodes,
		Connections: connections,
		Settings:    settings,
	}

	workflow := r.updateWorkflowViaAPI(ctx, workflowID, workflowRequest, &resp.Diagnostics)
	// Check for API error
	if resp.Diagnostics.HasError() {
		return false
	}

	// Update tags.
	r.updateWorkflowTags(ctx, workflowID, plan, workflow, &resp.Diagnostics)
	// Check for tag update errors.
	if resp.Diagnostics.HasError() {
		// Return failure.
		return false
	}

	// Handle project transfer if project_id changed.
	if !plan.ProjectID.Equal(state.ProjectID) {
		// Transfer to new project if project_id is set
		// Note: Removing a workflow from a project (changing from a
		// value to null) is not supported by the n8n API. The workflow
		// will remain in its current project if project_id changes to null.
		if !plan.ProjectID.IsNull() && !plan.ProjectID.IsUnknown() {
			updatedWorkflow := r.handleProjectAssignment(ctx, workflowID, plan.ProjectID.ValueString(), &resp.Diagnostics)
			// Check if project assignment succeeded
			if resp.Diagnostics.HasError() {
				// Return failure.
				return false
			}
			// Use updated workflow if available.
			if updatedWorkflow != nil {
				workflow = updatedWorkflow
			}
		}
	}

	// Map workflow response to state.
	mapWorkflowToModel(ctx, workflow, plan, &resp.Diagnostics)

	// Ensure computed timestamp fields have known values.
	// The PUT API may not return these fields, so fall back to state values.
	if plan.CreatedAt.IsUnknown() || plan.CreatedAt.IsNull() {
		plan.CreatedAt = state.CreatedAt
	}
	if plan.UpdatedAt.IsUnknown() || plan.UpdatedAt.IsNull() {
		plan.UpdatedAt = state.UpdatedAt
	}

	// Return success.
	return true
}

// Delete deletes the resource and removes the Terraform state on success.
//
// Params:
//   - ctx: Context for the operation
//   - req: Delete request containing state data
//   - resp: Delete response for error handling
//
// Returns:
//   - None: Updates resp with error diagnostics if any
func (r *WorkflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *models.Resource

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check for state parsing errors.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute delete logic
	r.executeDeleteLogic(ctx, state, resp)
}

// executeDeleteLogic contains the main logic for deleting a workflow.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Delete response
//
// Returns:
//   - bool: True if delete succeeded, false otherwise
func (r *WorkflowResource) executeDeleteLogic(ctx context.Context, state *models.Resource, resp *resource.DeleteResponse) bool {
	// Delete workflow using SDK.
	_, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdDelete(ctx, state.ID.ValueString()).Execute()
	// Check for non-nil HTTP response.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for API error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting workflow",
			fmt.Sprintf("Could not delete workflow ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		// Return failure.
		return false
	}

	// Return success.
	return true
}

// ImportState imports the resource into Terraform state.
//
// Params:
//   - ctx: Context for the operation
//   - req: Import state request containing resource ID
//   - resp: Import state response for state and error handling
//
// Returns:
//   - None: Updates resp with imported state and error diagnostics
func (r *WorkflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
