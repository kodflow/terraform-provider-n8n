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
)

// Ensure VariableResource implements required interfaces.
var (
	_ resource.Resource                = &VariableResource{}
	_ resource.ResourceWithConfigure   = &VariableResource{}
	_ resource.ResourceWithImportState = &VariableResource{}
)

// VariableResource defines the resource implementation for n8n variables.
// Note: n8n API has limitations - POST returns 201 with no body, no GET by ID endpoint.
// We work around this by using the LIST endpoint and filtering.
type VariableResource struct {
	client *client.N8nClient
}

// NewVariableResource creates and returns a new VariableResource instance.
func NewVariableResource() resource.Resource {
	// Return result.
	return &VariableResource{}
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
func (r *VariableResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
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
func (r *VariableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
func (r *VariableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
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
	var plan *VariableResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for error.
	if resp.Diagnostics.HasError() {
		return
	}

	// Build variable request
	variableRequest := buildVariableRequest(plan)

	// POST returns 201 with no body
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
		return
	}

	// Workaround: List all variables to find the one we just created
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
		return
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
		return
	}

	// Map variable to model
	mapVariableToResourceModel(foundVariable, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
	var state *VariableResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Workaround: No GET by ID, use LIST and filter
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
		// Return result.
		return
	}

	// Find our variable by ID
	var foundVariable *n8nsdk.Variable
	// Check for non-nil value.
	if variableList.Data != nil {
		// Iterate over items.
		for _, v := range variableList.Data {
			// Check for non-nil value.
			if v.Id != nil && *v.Id == state.ID.ValueString() {
				foundVariable = &v
				break
			}
		}
	}

	// Check condition.
	if foundVariable == nil {
		// Variable not found = deleted outside Terraform
		resp.State.RemoveResource(ctx)
		// Return result.
		return
	}

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

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
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
	var plan *VariableResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for error.
	if resp.Diagnostics.HasError() {
		return
	}

	// Build variable request
	variableRequest := buildVariableRequest(plan)

	// PUT returns 204 with no body
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
		return
	}

	// Workaround: List all variables to verify the update
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
		return
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
		return
	}

	// Map variable to model (only updates non-ID fields)
	mapVariableToResourceModel(foundVariable, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
	var state *VariableResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

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
		// Return result.
		return
	}
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
