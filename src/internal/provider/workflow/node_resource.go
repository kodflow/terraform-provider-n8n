// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// WorkflowNodeResource implements workflow node resources for modular workflow composition.
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/models"
)

const (
	// NODE_ATTRIBUTES_SIZE defines the initial capacity for node attributes map.
	NODE_ATTRIBUTES_SIZE int = 10
	// DEFAULT_TYPE_VERSION is the default node type version.
	DEFAULT_TYPE_VERSION int64 = 1
)

// Ensure WorkflowNodeResource implements required interfaces.
var (
	_ resource.Resource                = &WorkflowNodeResource{}
	_ resource.ResourceWithConfigure   = &WorkflowNodeResource{}
	_ resource.ResourceWithImportState = &WorkflowNodeResource{}
)

// WorkflowNodeResourceInterface defines the complete interface for workflow
// node resources.
type WorkflowNodeResourceInterface interface {
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

// WorkflowNodeResource defines a local-only resource for workflow nodes.
// This resource does not make API calls; it exists purely in Terraform state
// to generate node JSON for use in n8n_workflow resources.
type WorkflowNodeResource struct{}

// NewWorkflowNodeResource creates a new WorkflowNodeResource instance.
//
// Returns:
//   - *WorkflowNodeResource: A new WorkflowNodeResource instance.
func NewWorkflowNodeResource() *WorkflowNodeResource {
	// Return new instance.
	return &WorkflowNodeResource{}
}

// NewWorkflowNodeResourceWrapper creates a new WorkflowNodeResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - resource.Resource: the wrapped WorkflowNodeResource instance
func NewWorkflowNodeResourceWrapper() resource.Resource {
	// Return the wrapped resource instance.
	return NewWorkflowNodeResource()
}

// Metadata returns the resource type name.
//
// Params:
//   - _ctx: The context for the request (unused).
//   - req: The metadata request containing provider type name.
//   - resp: The metadata response to populate.
func (r *WorkflowNodeResource) Metadata(_ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_node"
}

// Schema defines the schema for the resource.
//
// Params:
//   - _ctx: The context for the request (unused).
//   - _req: The schema request (unused).
//   - resp: The schema response to populate.
func (r *WorkflowNodeResource) Schema(_ctx context.Context, _req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Defines a workflow node for modular workflow composition. " +
			"This resource exists only in Terraform state and generates JSON that can be " +
			"used in n8n_workflow resources. No API calls are made.",
		Attributes: r.getSchemaAttributes(),
	}
}

// getSchemaAttributes returns the schema attributes for the node resource.
//
// Returns:
//   - map[string]schema.Attribute: The schema attributes.
func (r *WorkflowNodeResource) getSchemaAttributes() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute, NODE_ATTRIBUTES_SIZE)

	// Add required attributes.
	r.addRequiredAttributes(attrs)

	// Add optional attributes.
	r.addOptionalAttributes(attrs)

	// Add computed JSON attribute.
	attrs["node_json"] = schema.StringAttribute{
		MarkdownDescription: "Computed JSON representation of this node for use in workflows",
		Computed:            true,
	}

	// Return complete attributes map.
	return attrs
}

// addRequiredAttributes adds required node attributes to the map.
//
// Params:
//   - attrs: The attributes map to update.
func (r *WorkflowNodeResource) addRequiredAttributes(attrs map[string]schema.Attribute) {
	attrs["id"] = schema.StringAttribute{
		MarkdownDescription: "Unique identifier for this node (computed)",
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	attrs["name"] = schema.StringAttribute{
		MarkdownDescription: "Display name of the node (used in connections)",
		Required:            true,
	}
	attrs["type"] = schema.StringAttribute{
		MarkdownDescription: "n8n node type (e.g., 'n8n-nodes-base.webhook')",
		Required:            true,
	}
	attrs["type_version"] = schema.Int64Attribute{
		MarkdownDescription: "Version of the node type",
		Optional:            true,
		Computed:            true,
		Default:             int64default.StaticInt64(DEFAULT_TYPE_VERSION),
	}
	attrs["position"] = schema.ListAttribute{
		MarkdownDescription: "Position [x, y] coordinates for UI display",
		ElementType:         types.Int64Type,
		Required:            true,
	}
}

// addOptionalAttributes adds optional node attributes to the map.
//
// Params:
//   - attrs: The attributes map to update.
func (r *WorkflowNodeResource) addOptionalAttributes(attrs map[string]schema.Attribute) {
	attrs["parameters"] = schema.StringAttribute{
		MarkdownDescription: "Node parameters as JSON string",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("{}"),
	}
	attrs["webhook_id"] = schema.StringAttribute{
		MarkdownDescription: "Webhook identifier for webhook nodes",
		Optional:            true,
	}
	attrs["disabled"] = schema.BoolAttribute{
		MarkdownDescription: "Whether the node is disabled",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	}
	attrs["notes"] = schema.StringAttribute{
		MarkdownDescription: "User notes about the node",
		Optional:            true,
	}
}

// Configure configures the resource (no-op for local resources).
//
// Params:
//   - _ctx: The context for the request (unused).
//   - _req: The configuration request (unused).
//   - _resp: The configuration response (unused).
func (r *WorkflowNodeResource) Configure(_ctx context.Context, _req resource.ConfigureRequest, _resp *resource.ConfigureResponse) {
	// No configuration needed for local-only resources.
}

// Create creates the resource in Terraform state.
//
// Params:
//   - ctx: The context for the request.
//   - req: The create request containing the planned resource data.
//   - resp: The create response to populate.
func (r *WorkflowNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *models.NodeResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check if there were errors retrieving the plan.
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate unique ID based on name and type.
	plan.ID = types.StringValue(r.generateNodeID(plan))

	// Generate node JSON.
	// Check if JSON generation failed.
	if !r.generateNodeJSON(ctx, plan, &resp.Diagnostics) {
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
func (r *WorkflowNodeResource) Read(_ctx context.Context, _req resource.ReadRequest, _resp *resource.ReadResponse) {
	// Local-only resource, state is always current.
}

// Update updates the resource in Terraform state.
//
// Params:
//   - ctx: The context for the request.
//   - req: The update request containing the planned resource data.
//   - resp: The update response to populate.
func (r *WorkflowNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *models.NodeResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check if there were errors retrieving the plan.
	if resp.Diagnostics.HasError() {
		return
	}

	// Regenerate node JSON with updated values.
	// Check if JSON generation failed.
	if !r.generateNodeJSON(ctx, plan, &resp.Diagnostics) {
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
func (r *WorkflowNodeResource) Delete(_ctx context.Context, _req resource.DeleteRequest, _resp *resource.DeleteResponse) {
	// Resource is removed from state automatically.
}

// ImportState imports the resource into Terraform state.
//
// Params:
//   - ctx: The context for the request.
//   - req: The import state request containing the resource ID.
//   - resp: The import state response to populate.
func (r *WorkflowNodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// generateNodeID creates a unique identifier for the node.
// Uses a hash of name + type to ensure consistency across applies.
//
// Params:
//   - plan: The resource data containing node information.
//
// Returns:
//   - string: A unique UUID based on the node name and type.
func (r *WorkflowNodeResource) generateNodeID(plan *models.NodeResource) string {
	// Use name + type to create a stable ID.
	input := fmt.Sprintf("%s-%s", plan.Name.ValueString(), plan.Type.ValueString())
	hash := sha256.Sum256([]byte(input))
	// Return UUID generated from node hash.
	return uuid.NewSHA1(uuid.NameSpaceOID, hash[:]).String()
}

// parseNodeParameters parses and adds parameters to the node.
//
// Params:
//   - plan: The resource data containing parameters.
//   - node: The node map to update with parameters.
//   - diags: The diagnostics to append errors to.
//
// Returns:
//   - bool: True if parsing succeeded, false otherwise.
func (r *WorkflowNodeResource) parseNodeParameters(plan *models.NodeResource, node map[string]any, diags *diag.Diagnostics) bool {
	// Check if parameters are provided.
	if !plan.Parameters.IsNull() && !plan.Parameters.IsUnknown() {
		var params map[string]any
		// Unmarshal JSON parameters. Check if parsing failed.
		if err := json.Unmarshal([]byte(plan.Parameters.ValueString()), &params); err != nil {
			diags.AddError(
				"Invalid parameters JSON",
				fmt.Sprintf("Could not parse parameters: %s", err.Error()),
			)
			return false
		}
		node["parameters"] = params
		// Handle null parameters by setting empty JSON.
	} else if plan.Parameters.IsNull() {
		// Set empty parameters if null to ensure it's known after apply.
		plan.Parameters = types.StringValue("{}")
		node["parameters"] = map[string]any{}
	}
	// Return success.
	return true
}

// addOptionalNodeFields adds optional fields to the node.
//
// Params:
//   - plan: The resource data containing optional fields.
//   - node: The node map to update with optional fields.
func (r *WorkflowNodeResource) addOptionalNodeFields(plan *models.NodeResource, node map[string]any) {
	// Check if webhook ID is provided.
	if !plan.WebhookID.IsNull() && !plan.WebhookID.IsUnknown() {
		node["webhookId"] = plan.WebhookID.ValueString()
	}

	// Check if node is disabled.
	if !plan.Disabled.IsNull() && !plan.Disabled.IsUnknown() && plan.Disabled.ValueBool() {
		node["disabled"] = true
	}

	// Check if notes are provided.
	if !plan.Notes.IsNull() && !plan.Notes.IsUnknown() {
		node["notes"] = plan.Notes.ValueString()
	}
}

// generateNodeJSON creates the JSON representation of the node.
//
// Params:
//   - ctx: The context for the request.
//   - plan: The resource data to convert to JSON.
//   - diags: The diagnostics to append errors to.
//
// Returns:
//   - bool: True if JSON generation succeeded, false otherwise.
func (r *WorkflowNodeResource) generateNodeJSON(ctx context.Context, plan *models.NodeResource, diags *diag.Diagnostics) bool {
	// Build node structure.
	node := map[string]any{
		"id":   plan.ID.ValueString(),
		"name": plan.Name.ValueString(),
		"type": plan.Type.ValueString(),
	}

	// Add type version.
	// Check if type version is provided.
	if !plan.TypeVersion.IsNull() && !plan.TypeVersion.IsUnknown() {
		node["typeVersion"] = plan.TypeVersion.ValueInt64()
	}

	// Add position.
	var position []int64
	diags.Append(plan.Position.ElementsAs(ctx, &position, false)...)
	// Check if position extraction failed.
	if diags.HasError() {
		return false
	}
	node["position"] = position

	// Add parameters.
	// Check if parameter parsing failed.
	if !r.parseNodeParameters(plan, node, diags) {
		return false
	}

	// Add optional fields.
	r.addOptionalNodeFields(plan, node)

	// Marshal to JSON.
	jsonBytes, err := json.Marshal(node)
	// Check if JSON marshalling failed.
	if err != nil {
		diags.AddError(
			"Failed to generate node JSON",
			fmt.Sprintf("Could not marshal node to JSON: %s", err.Error()),
		)
		return false
	}

	plan.NodeJSON = types.StringValue(string(jsonBytes))
	// Return success.
	return true
}
