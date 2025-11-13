// Package variable implements environment variable management resources and data sources.
package variable

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/variable/models"
)

// Ensure VariableResource implements required interfaces.
var (
	_ resource.Resource                = &VariableResource{}
	_ VariableResourceInterface        = &VariableResource{}
	_ resource.ResourceWithConfigure   = &VariableResource{}
	_ resource.ResourceWithImportState = &VariableResource{}
)

// VariableResourceInterface defines the interface for VariableResource.
type VariableResourceInterface interface {
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

// VariableResource defines the resource implementation for n8n variables.
// Note: n8n API has limitations - POST returns 201 with no body, no GET by ID endpoint.
// We work around this by using the LIST endpoint and filtering.
type VariableResource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewVariableResource creates and returns a new VariableResource instance.
//
// Returns:
//   - resource.Resource: A new VariableResource instance
func NewVariableResource() *VariableResource {
	// Return result.
	return &VariableResource{}
}

// NewVariableResourceWrapper creates a new VariableResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - resource.Resource: the wrapped VariableResource instance
func NewVariableResourceWrapper() resource.Resource {
	// Return the wrapped resource instance.
	return NewVariableResource()
}

// Metadata returns the resource type name.
//
// Params:
//   - ctx: Context for the operation
//   - req: Metadata request from Terraform plugin framework
//   - resp: Metadata response to populate with type name
//
// Returns:
//
//	(none, modifies resp parameter in place)
func (r *VariableResource) Metadata(_ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable"
}

// Schema defines the schema for the resource.
//
// Params:
//   - ctx: Context for the operation
//   - req: Schema request from Terraform plugin framework
//   - resp: Schema response to populate with resource schema
//
// Returns:
//
//	(none, modifies resp parameter in place)
func (r *VariableResource) Schema(_ctx context.Context, _req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "n8n variable resource. Note: API limitations require workarounds for Read operations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Variable identifier",
				Computed:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Variable key",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Variable value",
				Required:            true,
				Sensitive:           true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Variable type",
				Optional:            true,
				Computed:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Project ID to associate this variable with",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
//
// Params:
//   - ctx: Context for the operation
//   - req: Configure request from Terraform plugin framework
//   - resp: Configure response to store any errors
//
// Returns:
//
//	(none, modifies resp parameter in place)
func (r *VariableResource) Configure(_ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
		// Return result.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	// Check condition.
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
// Workaround: API returns 201 with no body, so we must call LIST to get the ID.
//
// Params:
//   - ctx: Context for the operation
//   - req: Create request from Terraform plugin framework
//   - resp: Create response to set state and store any errors
//
// Returns:
//
//	(none, modifies resp parameter in place)
func (r *VariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for error.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute create logic
	if !r.executeCreateLogic(ctx, plan, resp) {
		// Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// executeCreateLogic contains the main logic for creating a variable.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - resp: Create response
//
// Returns:
//   - bool: True if creation succeeded, false otherwise
func (r *VariableResource) executeCreateLogic(ctx context.Context, plan *models.Resource, resp *resource.CreateResponse) bool {
	// Execute create operation
	if !r.executeVariableCreate(ctx, plan, resp) {
		// Return failure.
		return false
	}

	// Find created variable
	foundVariable := r.findCreatedVariable(ctx, plan, resp)
	// Return early if variable not found.
	if foundVariable == nil {
		// Return failure.
		return false
	}

	// Map variable to model
	mapVariableToResourceModel(foundVariable, plan)

	// Return success.
	return true
}

// executeVariableCreate performs the API call to create a variable.
//
// Params:
//   - ctx: context for operation
//   - plan: variable resource model
//   - resp: create response
//
// Returns:
//   - bool: true if successful, false if error occurred
func (r *VariableResource) executeVariableCreate(ctx context.Context, plan *models.Resource, resp *resource.CreateResponse) bool {
	variableRequest := buildVariableRequest(plan)

	httpResp, err := r.client.APIClient.VariablesAPI.VariablesPost(ctx).
		VariableCreate(variableRequest).
		Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating variable",
			fmt.Sprintf("Could not create variable: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return false
	}

	return true
}

// findCreatedVariable retrieves and finds the created variable from the list.
//
// Params:
//   - ctx: context for operation
//   - plan: variable resource model
//   - resp: create response
//
// Returns:
//   - *n8nsdk.Variable: found variable or nil if not found
func (r *VariableResource) findCreatedVariable(ctx context.Context, plan *models.Resource, resp *resource.CreateResponse) *n8nsdk.Variable {
	variableList, httpResp, err := r.client.APIClient.VariablesAPI.VariablesGet(ctx).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading variable after creation",
			fmt.Sprintf("Variable was created but could not retrieve ID: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return with error.
		return nil
	}

	// Find our variable by key
	var foundVariable *n8nsdk.Variable
	var found bool
	// Check for non-nil value.
	if variableList.Data != nil {
		foundVariable, found = findVariableByKey(variableList.Data, plan.Key.ValueString())
	}

	// Check condition.
	if !found {
		resp.Diagnostics.AddError(
			"Error finding created variable",
			fmt.Sprintf("Variable with key '%s' was created but not found in list", plan.Key.ValueString()),
		)
		// Return with error.
		return nil
	}

	// Return result.
	return foundVariable
}

// Read refreshes the Terraform state with the latest data.
// Workaround: No GET by ID endpoint, so we use LIST and filter.
//
// Params:
//   - ctx: Context for the operation
//   - req: Read request from Terraform plugin framework
//   - resp: Read response to set state and store any errors
//
// Returns:
//
//	(none, modifies resp parameter in place)
func (r *VariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *models.Resource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute read logic
	if !r.executeReadLogic(ctx, state, resp) {
		// Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// executeReadLogic contains the main logic for reading a variable.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Read response
//
// Returns:
//   - bool: True if read succeeded, false otherwise
func (r *VariableResource) executeReadLogic(ctx context.Context, state *models.Resource, resp *resource.ReadResponse) bool {
	// Fetch variable list from API
	variableList, httpResp, err := r.client.APIClient.VariablesAPI.VariablesGet(ctx).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading variable",
			fmt.Sprintf("Could not read variable ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		// Return failure.
		return false
	}

	// Find variable by ID
	foundVariable := r.findVariableByID(variableList, state.ID.ValueString())
	// Check condition.
	if foundVariable == nil {
		// Variable not found = deleted outside Terraform
		resp.State.RemoveResource(ctx)
		// Return failure.
		return false
	}

	// Update state with found variable data
	r.updateStateFromVariable(foundVariable, state)

	// Return success.
	return true
}

// findVariableByID finds a variable in the list by ID.
//
// Params:
//   - variableList: list of variables from API
//   - id: variable ID to find
//
// Returns:
//   - *n8nsdk.Variable: found variable or nil if not found
func (r *VariableResource) findVariableByID(variableList *n8nsdk.VariableList, id string) *n8nsdk.Variable {
	// Check for non-nil value.
	if variableList.Data == nil {
		// Return result.
		return nil
	}
	// Iterate over items.
	for _, v := range variableList.Data {
		// Check for non-nil value.
		if v.Id != nil && *v.Id == id {
			// Return matching variable.
			return &v
		}
	}
	// Return result.
	return nil
}

// updateStateFromVariable updates the state model with data from the found variable.
//
// Params:
//   - foundVariable: variable data from API
//   - state: models.Resource to update
func (r *VariableResource) updateStateFromVariable(foundVariable *n8nsdk.Variable, state *models.Resource) {
	state.Key = types.StringValue(foundVariable.Key)
	state.Value = types.StringValue(foundVariable.Value)
	// Check for non-nil value.
	if foundVariable.Type != nil {
		state.Type = types.StringPointerValue(foundVariable.Type)
	}
	// Check for non-nil value.
	if foundVariable.Project != nil && foundVariable.Project.Id != nil {
		state.ProjectID = types.StringPointerValue(foundVariable.Project.Id)
	}
}

// Update updates the resource and sets the updated Terraform state on success.
// Workaround: API returns 204 with no body, so we must call LIST to verify.
//
// Params:
//   - ctx: Context for the operation
//   - req: Update request from Terraform plugin framework
//   - resp: Update response to set state and store any errors
//
// Returns:
//
//	(none, modifies resp parameter in place)
func (r *VariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for error.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute update logic
	if !r.executeUpdateLogic(ctx, plan, resp) {
		// Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// executeUpdateLogic contains the main logic for updating a variable.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - resp: Update response
//
// Returns:
//   - bool: True if update succeeded, false otherwise
func (r *VariableResource) executeUpdateLogic(ctx context.Context, plan *models.Resource, resp *resource.UpdateResponse) bool {
	// Execute the update API call
	if !r.executeVariableUpdate(ctx, plan, resp) {
		// Return failure.
		return false
	}

	// Verify the update by listing variables
	foundVariable := r.findUpdatedVariable(ctx, plan, resp)
	// Return early if variable not found.
	if foundVariable == nil {
		// Return failure.
		return false
	}

	// Map variable to model (only updates non-ID fields)
	mapVariableToResourceModel(foundVariable, plan)

	// Return success.
	return true
}

// executeVariableUpdate performs the API call to update a variable.
//
// Params:
//   - ctx: context for operation
//   - plan: variable resource model
//   - resp: update response
//
// Returns:
//   - bool: true if successful, false if error occurred
func (r *VariableResource) executeVariableUpdate(ctx context.Context, plan *models.Resource, resp *resource.UpdateResponse) bool {
	variableRequest := buildVariableRequest(plan)

	httpResp, err := r.client.APIClient.VariablesAPI.VariablesIdPut(ctx, plan.ID.ValueString()).
		VariableCreate(variableRequest).
		Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating variable",
			fmt.Sprintf("Could not update variable ID %s: %s\nHTTP Response: %v", plan.ID.ValueString(), err.Error(), httpResp),
		)
		return false
	}

	return true
}

// findUpdatedVariable retrieves and finds the updated variable from the list.
//
// Params:
//   - ctx: context for operation
//   - plan: variable resource model
//   - resp: update response
//
// Returns:
//   - *n8nsdk.Variable: found variable or nil if not found
func (r *VariableResource) findUpdatedVariable(ctx context.Context, plan *models.Resource, resp *resource.UpdateResponse) *n8nsdk.Variable {
	variableList, httpResp, err := r.client.APIClient.VariablesAPI.VariablesGet(ctx).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading variable after update",
			fmt.Sprintf("Variable was updated but could not verify: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return with error.
		return nil
	}

	// Find our variable by ID
	var foundVariable *n8nsdk.Variable
	var found bool
	// Check for non-nil value.
	if variableList.Data != nil {
		foundVariable, found = findVariableByID(variableList.Data, plan.ID.ValueString())
	}

	// Check condition.
	if !found {
		resp.Diagnostics.AddError(
			"Error verifying updated variable",
			fmt.Sprintf("Variable with ID '%s' was updated but not found in list", plan.ID.ValueString()),
		)
		// Return with error.
		return nil
	}

	// Return the found variable.
	return foundVariable
}

// Delete deletes the resource and removes the Terraform state on success.
//
// Params:
//   - ctx: Context for the operation
//   - req: Delete request from Terraform plugin framework
//   - resp: Delete response to store any errors
//
// Returns:
//
//	(none, modifies resp parameter in place)
func (r *VariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *models.Resource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute delete logic
	r.executeDeleteLogic(ctx, state, resp)
}

// executeDeleteLogic contains the main logic for deleting a variable.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Delete response
//
// Returns:
//   - bool: True if delete succeeded, false otherwise
func (r *VariableResource) executeDeleteLogic(ctx context.Context, state *models.Resource, resp *resource.DeleteResponse) bool {
	// DELETE returns 204 with no body
	httpResp, err := r.client.APIClient.VariablesAPI.VariablesIdDelete(ctx, state.ID.ValueString()).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting variable",
			fmt.Sprintf("Could not delete variable ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
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
//   - req: ImportState request from Terraform plugin framework
//   - resp: ImportState response to store state and any errors
//
// Returns:
//
//	(none, modifies resp parameter in place)
func (r *VariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
