// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// WorkflowConnectionResource implements workflow connection resources for modular workflow composition.
package workflow

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/models"
)

const (
	// DEFAULT_OUTPUT_TYPE is the default output type for connections.
	DEFAULT_OUTPUT_TYPE string = "main"
	// DEFAULT_INPUT_TYPE is the default input type for connections.
	DEFAULT_INPUT_TYPE string = "main"
	// DEFAULT_OUTPUT_INDEX is the default output index.
	DEFAULT_OUTPUT_INDEX int64 = 0
	// DEFAULT_INPUT_INDEX is the default input index.
	DEFAULT_INPUT_INDEX int64 = 0
	// CONNECTION_ATTRIBUTES_SIZE defines the initial capacity for
	// connection attributes map.
	CONNECTION_ATTRIBUTES_SIZE int = 8
)

// Ensure WorkflowConnectionResource implements required interfaces.
var (
	_ resource.Resource                = &WorkflowConnectionResource{}
	_ resource.ResourceWithConfigure   = &WorkflowConnectionResource{}
	_ resource.ResourceWithImportState = &WorkflowConnectionResource{}
)

// WorkflowConnectionResourceInterface defines the complete interface for
// workflow connection resources.
type WorkflowConnectionResourceInterface interface {
	Metadata(context.Context, resource.MetadataRequest,
		*resource.MetadataResponse)
	Schema(context.Context, resource.SchemaRequest, *resource.SchemaResponse)
	Configure(context.Context, resource.ConfigureRequest,
		*resource.ConfigureResponse)
	Create(context.Context, resource.CreateRequest, *resource.CreateResponse)
	Read(context.Context, resource.ReadRequest, *resource.ReadResponse)
	Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse)
	Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse)
	ImportState(context.Context, resource.ImportStateRequest,
		*resource.ImportStateResponse)
}

// WorkflowConnectionResource defines a local-only resource for workflow
// connections. This resource does not make API calls; it exists purely in
// Terraform state to generate connection JSON for use in n8n_workflow resources.
type WorkflowConnectionResource struct{}

// NewWorkflowConnectionResource creates a new WorkflowConnectionResource
// instance.
//
// Returns:
//   - *WorkflowConnectionResource: A new WorkflowConnectionResource instance.
func NewWorkflowConnectionResource() *WorkflowConnectionResource {
	// Return new instance.
	return &WorkflowConnectionResource{}
}

// NewWorkflowConnectionResourceWrapper creates a new WorkflowConnectionResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - resource.Resource: the wrapped WorkflowConnectionResource instance
func NewWorkflowConnectionResourceWrapper() resource.Resource {
	// Return the wrapped resource instance.
	return NewWorkflowConnectionResource()
}

// Metadata returns the resource type name.
//
// Params:
//   - _ctx: The context for the request (unused).
//   - req: The metadata request containing provider type name.
//   - resp: The metadata response to populate.
func (r *WorkflowConnectionResource) Metadata(_ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_connection"
}

// Schema defines the schema for the resource.
//
// Params:
//   - _ctx: The context for the request (unused).
//   - _req: The schema request (unused).
//   - resp: The schema response to populate.
func (r *WorkflowConnectionResource) Schema(_ctx context.Context, _req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Defines a connection between workflow nodes for modular workflow composition. " +
			"This resource exists only in Terraform state and generates JSON that can be " +
			"used in n8n_workflow resources. No API calls are made.",
		Attributes: r.getSchemaAttributes(),
	}
}

// getSchemaAttributes returns the schema attributes for the connection resource.
//
// Returns:
//   - map[string]schema.Attribute: The schema attributes.
func (r *WorkflowConnectionResource) getSchemaAttributes() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute, CONNECTION_ATTRIBUTES_SIZE)

	// Add ID attribute.
	attrs["id"] = r.getIDAttribute()

	// Add source attributes.
	r.addSourceAttributes(attrs)

	// Add target attributes.
	r.addTargetAttributes(attrs)

	// Add computed JSON attribute.
	attrs["connection_json"] = schema.StringAttribute{
		MarkdownDescription: "Computed JSON representation of this connection for use in workflows",
		Computed:            true,
	}

	// Return complete attributes map.
	return attrs
}

// getIDAttribute returns the ID schema attribute.
//
// Returns:
//   - schema.StringAttribute: The ID attribute configuration.
func (r *WorkflowConnectionResource) getIDAttribute() schema.StringAttribute {
	// Return ID attribute with plan modifiers.
	return schema.StringAttribute{
		MarkdownDescription: "Unique identifier for this connection (computed)",
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// addSourceAttributes adds source node attributes to the map.
//
// Params:
//   - attrs: The attributes map to update.
func (r *WorkflowConnectionResource) addSourceAttributes(attrs map[string]schema.Attribute) {
	attrs["source_node"] = schema.StringAttribute{
		MarkdownDescription: "Name of the source node",
		Required:            true,
	}
	attrs["source_output"] = schema.StringAttribute{
		MarkdownDescription: "Output type from source node (default: 'main')",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(DEFAULT_OUTPUT_TYPE),
	}
	attrs["source_output_index"] = schema.Int64Attribute{
		MarkdownDescription: "Index of the source output (for nodes with multiple outputs like Switch)",
		Optional:            true,
		Computed:            true,
		Default:             int64default.StaticInt64(DEFAULT_OUTPUT_INDEX),
	}
}

// addTargetAttributes adds target node attributes to the map.
//
// Params:
//   - attrs: The attributes map to update.
func (r *WorkflowConnectionResource) addTargetAttributes(attrs map[string]schema.Attribute) {
	attrs["target_node"] = schema.StringAttribute{
		MarkdownDescription: "Name of the destination node",
		Required:            true,
	}
	attrs["target_input"] = schema.StringAttribute{
		MarkdownDescription: "Input type to target node (default: 'main')",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(DEFAULT_INPUT_TYPE),
	}
	attrs["target_input_index"] = schema.Int64Attribute{
		MarkdownDescription: "Index of the target input",
		Optional:            true,
		Computed:            true,
		Default:             int64default.StaticInt64(DEFAULT_INPUT_INDEX),
	}
}

// Configure configures the resource (no-op for local resources).
//
// Params:
//   - _ctx: The context for the request (unused).
//   - _req: The configuration request (unused).
//   - _resp: The configuration response (unused).
func (r *WorkflowConnectionResource) Configure(_ctx context.Context, _req resource.ConfigureRequest, _resp *resource.ConfigureResponse) {
	// No configuration needed for local-only resources.
}

// Create creates the resource in Terraform state.
//
// Params:
//   - ctx: The context for the request.
//   - req: The create request containing the planned resource data.
//   - resp: The create response to populate.
func (r *WorkflowConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *models.ConnectionResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check if there were errors retrieving the plan.
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate unique ID.
	plan.ID = types.StringValue(r.generateConnectionID(plan))

	// Generate connection JSON.
	// Check if JSON generation failed.
	if !r.generateConnectionJSON(plan, &resp.Diagnostics) {
		return
	}

	// Save to state.
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state (no-op for local resources).
//
// Params:
//   - _ctx: The context for the request (unused).
//   - _req: The read request (unused).
//   - _resp: The read response (unused).
func (r *WorkflowConnectionResource) Read(_ctx context.Context, _req resource.ReadRequest, _resp *resource.ReadResponse) {
	// Local-only resource, state is always current.
}

// Update updates the resource in Terraform state.
//
// Params:
//   - ctx: The context for the request.
//   - req: The update request containing the planned resource data.
//   - resp: The update response to populate.
func (r *WorkflowConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *models.ConnectionResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check if there were errors retrieving the plan.
	if resp.Diagnostics.HasError() {
		return
	}

	// Regenerate connection JSON.
	// Check if JSON generation failed.
	if !r.generateConnectionJSON(plan, &resp.Diagnostics) {
		return
	}

	// Save to state.
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete removes the resource from Terraform state.
//
// Params:
//   - _ctx: The context for the request (unused).
//   - _req: The delete request (unused).
//   - _resp: The delete response (unused).
func (r *WorkflowConnectionResource) Delete(_ctx context.Context, _req resource.DeleteRequest, _resp *resource.DeleteResponse) {
	// Resource is removed from state automatically.
}

// ImportState imports the resource into Terraform state.
//
// Params:
//   - ctx: The context for the request.
//   - req: The import state request containing the resource ID.
//   - resp: The import state response to populate.
func (r *WorkflowConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// generateConnectionID creates a unique identifier for the connection.
//
// Params:
//   - plan: The resource data containing connection information.
//
// Returns:
//   - string: A unique UUID based on the connection parameters.
func (r *WorkflowConnectionResource) generateConnectionID(plan *models.ConnectionResource) string {
	input := fmt.Sprintf("%s-%d-%s-%d",
		plan.SourceNode.ValueString(),
		plan.SourceOutputIndex.ValueInt64(),
		plan.TargetNode.ValueString(),
		plan.TargetInputIndex.ValueInt64(),
	)
	hash := sha256.Sum256([]byte(input))
	// Return UUID generated from connection hash.
	return uuid.NewSHA1(uuid.NameSpaceOID, hash[:]).String()
}

// generateConnectionJSON creates the JSON representation of the connection.
//
// Params:
//   - plan: The resource data to convert to JSON.
//   - diags: The diagnostics to append errors to.
//
// Returns:
//   - bool: True if JSON generation succeeded, false otherwise.
func (r *WorkflowConnectionResource) generateConnectionJSON(plan *models.ConnectionResource, diags *diag.Diagnostics) bool {
	// Build connection metadata structure.
	conn := map[string]any{
		"source_node":         plan.SourceNode.ValueString(),
		"source_output":       plan.SourceOutput.ValueString(),
		"source_output_index": plan.SourceOutputIndex.ValueInt64(),
		"target_node":         plan.TargetNode.ValueString(),
		"target_input":        plan.TargetInput.ValueString(),
		"target_input_index":  plan.TargetInputIndex.ValueInt64(),
	}

	// Marshal to JSON.
	jsonBytes, err := json.Marshal(conn)
	// Check if JSON marshalling failed.
	if err != nil {
		diags.AddError(
			"Failed to generate connection JSON",
			fmt.Sprintf("Could not marshal connection to JSON: %s", err.Error()),
		)
		return false
	}

	plan.ConnectionJSON = types.StringValue(string(jsonBytes))
	// Return success.
	return true
}
