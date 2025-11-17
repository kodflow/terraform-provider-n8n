// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package node implements workflow node resources for modular workflow composition.
package node

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
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/node/models"
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

// WorkflowNodeResource defines a local-only resource for workflow nodes.
// This resource does not make API calls; it exists purely in Terraform state
// to generate node JSON for use in n8n_workflow resources.
type WorkflowNodeResource struct{}

// NewWorkflowNodeResource creates a new WorkflowNodeResource instance.
//
// Returns:
//   - resource.Resource: A new WorkflowNodeResource instance
func NewWorkflowNodeResource() resource.Resource {
	return &WorkflowNodeResource{}
}

// Metadata returns the resource type name.
func (r *WorkflowNodeResource) Metadata(_ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_node"
}

// Schema defines the schema for the resource.
func (r *WorkflowNodeResource) Schema(_ctx context.Context, _req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Defines a workflow node for modular workflow composition. " +
			"This resource exists only in Terraform state and generates JSON that can be " +
			"used in n8n_workflow resources. No API calls are made.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for this node (computed)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Display name of the node (used in connections)",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "n8n node type (e.g., 'n8n-nodes-base.webhook')",
				Required:            true,
			},
			"type_version": schema.Int64Attribute{
				MarkdownDescription: "Version of the node type",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(DEFAULT_TYPE_VERSION),
			},
			"position": schema.ListAttribute{
				MarkdownDescription: "Position [x, y] coordinates for UI display",
				ElementType:         types.Int64Type,
				Required:            true,
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: "Node parameters as JSON string",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("{}"),
			},
			"webhook_id": schema.StringAttribute{
				MarkdownDescription: "Webhook identifier for webhook nodes",
				Optional:            true,
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the node is disabled",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "User notes about the node",
				Optional:            true,
			},
			"node_json": schema.StringAttribute{
				MarkdownDescription: "Computed JSON representation of this node for use in workflows",
				Computed:            true,
			},
		},
	}
}

// Configure configures the resource (no-op for local resources).
func (r *WorkflowNodeResource) Configure(_ctx context.Context, _req resource.ConfigureRequest, _resp *resource.ConfigureResponse) {
	// No configuration needed for local-only resources.
}

// Create creates the resource in Terraform state.
func (r *WorkflowNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate unique ID based on name and type.
	plan.ID = types.StringValue(r.generateNodeID(plan))

	// Generate node JSON.
	if !r.generateNodeJSON(ctx, plan, &resp.Diagnostics) {
		return
	}

	// Save to state.
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state (no-op for local resources).
func (r *WorkflowNodeResource) Read(_ctx context.Context, _req resource.ReadRequest, _resp *resource.ReadResponse) {
	// Local-only resource, state is always current.
}

// Update updates the resource in Terraform state.
func (r *WorkflowNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Regenerate node JSON with updated values.
	if !r.generateNodeJSON(ctx, plan, &resp.Diagnostics) {
		return
	}

	// Save to state.
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete removes the resource from Terraform state.
func (r *WorkflowNodeResource) Delete(_ctx context.Context, _req resource.DeleteRequest, _resp *resource.DeleteResponse) {
	// Resource is removed from state automatically.
}

// ImportState imports the resource into Terraform state.
func (r *WorkflowNodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// generateNodeID creates a unique identifier for the node.
// Uses a hash of name + type to ensure consistency across applies.
func (r *WorkflowNodeResource) generateNodeID(plan *models.Resource) string {
	// Use name + type to create a stable ID.
	input := fmt.Sprintf("%s-%s", plan.Name.ValueString(), plan.Type.ValueString())
	hash := sha256.Sum256([]byte(input))
	return uuid.NewSHA1(uuid.NameSpaceOID, hash[:]).String()
}

// generateNodeJSON creates the JSON representation of the node.
func (r *WorkflowNodeResource) generateNodeJSON(ctx context.Context, plan *models.Resource, diags *diag.Diagnostics) bool {
	// Build node structure.
	node := map[string]any{
		"id":   plan.ID.ValueString(),
		"name": plan.Name.ValueString(),
		"type": plan.Type.ValueString(),
	}

	// Add type version.
	if !plan.TypeVersion.IsNull() && !plan.TypeVersion.IsUnknown() {
		node["typeVersion"] = plan.TypeVersion.ValueInt64()
	}

	// Add position.
	var position []int64
	diags.Append(plan.Position.ElementsAs(ctx, &position, false)...)
	if diags.HasError() {
		return false
	}
	node["position"] = position

	// Add parameters.
	if !plan.Parameters.IsNull() && !plan.Parameters.IsUnknown() {
		var params map[string]any
		if err := json.Unmarshal([]byte(plan.Parameters.ValueString()), &params); err != nil {
			diags.AddError(
				"Invalid parameters JSON",
				fmt.Sprintf("Could not parse parameters: %s", err.Error()),
			)
			return false
		}
		node["parameters"] = params
	} else if plan.Parameters.IsNull() {
		// Set empty parameters if null to ensure it's known after apply
		plan.Parameters = types.StringValue("{}")
		node["parameters"] = map[string]any{}
	}

	// Add optional fields.
	if !plan.WebhookID.IsNull() && !plan.WebhookID.IsUnknown() {
		node["webhookId"] = plan.WebhookID.ValueString()
	}

	if !plan.Disabled.IsNull() && !plan.Disabled.IsUnknown() && plan.Disabled.ValueBool() {
		node["disabled"] = true
	}

	if !plan.Notes.IsNull() && !plan.Notes.IsUnknown() {
		node["notes"] = plan.Notes.ValueString()
	}

	// Marshal to JSON.
	jsonBytes, err := json.Marshal(node)
	if err != nil {
		diags.AddError(
			"Failed to generate node JSON",
			fmt.Sprintf("Could not marshal node to JSON: %s", err.Error()),
		)
		return false
	}

	plan.NodeJSON = types.StringValue(string(jsonBytes))
	return true
}
