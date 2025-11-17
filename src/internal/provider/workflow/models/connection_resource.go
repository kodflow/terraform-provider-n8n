// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models defines data structures for workflow resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConnectionResource describes the workflow connection resource data model.
// This resource exists only in Terraform state and does not make API calls.
// It generates a connection entry that can be used in n8n_workflow resources.
type ConnectionResource struct {
	// ID is a computed unique identifier for this connection.
	ID types.String `tfsdk:"id"`

	// SourceNode is the name of the source node.
	SourceNode types.String `tfsdk:"source_node"`

	// SourceOutput is the output type (usually "main").
	SourceOutput types.String `tfsdk:"source_output"`

	// SourceOutputIndex is the index of the source output.
	// Used for nodes with multiple outputs like Switch.
	SourceOutputIndex types.Int64 `tfsdk:"source_output_index"`

	// TargetNode is the name of the destination node.
	TargetNode types.String `tfsdk:"target_node"`

	// TargetInput is the input type (usually "main").
	TargetInput types.String `tfsdk:"target_input"`

	// TargetInputIndex is the index of the target input.
	TargetInputIndex types.Int64 `tfsdk:"target_input_index"`

	// ConnectionJSON is the computed JSON representation of this connection.
	ConnectionJSON types.String `tfsdk:"connection_json"`
}
