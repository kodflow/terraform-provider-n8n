// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package connection implements workflow connection resources for modular workflow composition.
package connection

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
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/connection/models"
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
)

// Ensure WorkflowConnectionResource implements required interfaces.
var (
	_ resource.Resource                = &WorkflowConnectionResource{}
	_ resource.ResourceWithConfigure   = &WorkflowConnectionResource{}
	_ resource.ResourceWithImportState = &WorkflowConnectionResource{}
)

// WorkflowConnectionResource defines a local-only resource for workflow connections.
// This resource does not make API calls; it exists purely in Terraform state
// to generate connection JSON for use in n8n_workflow resources.
type WorkflowConnectionResource struct{}

// NewWorkflowConnectionResource creates a new WorkflowConnectionResource instance.
//
// Returns:
//   - resource.Resource: A new WorkflowConnectionResource instance
func NewWorkflowConnectionResource() resource.Resource {
	return &WorkflowConnectionResource{}
}

// Metadata returns the resource type name.
func (r *WorkflowConnectionResource) Metadata(_ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_connection"
}

// Schema defines the schema for the resource.
func (r *WorkflowConnectionResource) Schema(_ctx context.Context, _req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Defines a connection between workflow nodes for modular workflow composition. " +
			"This resource exists only in Terraform state and generates JSON that can be " +
			"used in n8n_workflow resources. No API calls are made.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for this connection (computed)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_node": schema.StringAttribute{
				MarkdownDescription: "Name of the source node",
				Required:            true,
			},
			"source_output": schema.StringAttribute{
				MarkdownDescription: "Output type from source node (default: 'main')",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(DEFAULT_OUTPUT_TYPE),
			},
			"source_output_index": schema.Int64Attribute{
				MarkdownDescription: "Index of the source output (for nodes with multiple outputs like Switch)",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(DEFAULT_OUTPUT_INDEX),
			},
			"target_node": schema.StringAttribute{
				MarkdownDescription: "Name of the destination node",
				Required:            true,
			},
			"target_input": schema.StringAttribute{
				MarkdownDescription: "Input type to target node (default: 'main')",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(DEFAULT_INPUT_TYPE),
			},
			"target_input_index": schema.Int64Attribute{
				MarkdownDescription: "Index of the target input",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(DEFAULT_INPUT_INDEX),
			},
			"connection_json": schema.StringAttribute{
				MarkdownDescription: "Computed JSON representation of this connection for use in workflows",
				Computed:            true,
			},
		},
	}
}

// Configure configures the resource (no-op for local resources).
func (r *WorkflowConnectionResource) Configure(_ctx context.Context, _req resource.ConfigureRequest, _resp *resource.ConfigureResponse) {
	// No configuration needed for local-only resources.
}

// Create creates the resource in Terraform state.
func (r *WorkflowConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate unique ID.
	plan.ID = types.StringValue(r.generateConnectionID(plan))

	// Generate connection JSON.
	if !r.generateConnectionJSON(plan, &resp.Diagnostics) {
		return
	}

	// Save to state.
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state (no-op for local resources).
func (r *WorkflowConnectionResource) Read(_ctx context.Context, _req resource.ReadRequest, _resp *resource.ReadResponse) {
	// Local-only resource, state is always current.
}

// Update updates the resource in Terraform state.
func (r *WorkflowConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Regenerate connection JSON.
	if !r.generateConnectionJSON(plan, &resp.Diagnostics) {
		return
	}

	// Save to state.
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete removes the resource from Terraform state.
func (r *WorkflowConnectionResource) Delete(_ctx context.Context, _req resource.DeleteRequest, _resp *resource.DeleteResponse) {
	// Resource is removed from state automatically.
}

// ImportState imports the resource into Terraform state.
func (r *WorkflowConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// generateConnectionID creates a unique identifier for the connection.
func (r *WorkflowConnectionResource) generateConnectionID(plan *models.Resource) string {
	input := fmt.Sprintf("%s-%d-%s-%d",
		plan.SourceNode.ValueString(),
		plan.SourceOutputIndex.ValueInt64(),
		plan.TargetNode.ValueString(),
		plan.TargetInputIndex.ValueInt64(),
	)
	hash := sha256.Sum256([]byte(input))
	return uuid.NewSHA1(uuid.NameSpaceOID, hash[:]).String()
}

// generateConnectionJSON creates the JSON representation of the connection.
func (r *WorkflowConnectionResource) generateConnectionJSON(plan *models.Resource, diags *diag.Diagnostics) bool {
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
	if err != nil {
		diags.AddError(
			"Failed to generate connection JSON",
			fmt.Sprintf("Could not marshal connection to JSON: %s", err.Error()),
		)
		return false
	}

	plan.ConnectionJSON = types.StringValue(string(jsonBytes))
	return true
}
