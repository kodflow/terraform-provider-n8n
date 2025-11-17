// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models defines data structures for workflow resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NodeResource describes the workflow node resource data model.
// This resource exists only in Terraform state and does not make API calls.
// It generates a node JSON structure that can be used in n8n_workflow resources.
type NodeResource struct {
	// ID is a computed unique identifier for this node (auto-generated).
	ID types.String `tfsdk:"id"`

	// Name is the display name of the node (used in connections).
	Name types.String `tfsdk:"name"`

	// Type is the n8n node type (e.g., "n8n-nodes-base.webhook").
	Type types.String `tfsdk:"type"`

	// TypeVersion is the version of the node type (optional).
	TypeVersion types.Int64 `tfsdk:"type_version"`

	// Position contains the [x, y] coordinates for UI display.
	Position types.List `tfsdk:"position"`

	// Parameters is a JSON string containing node-specific configuration.
	Parameters types.String `tfsdk:"parameters"`

	// WebhookID is an optional webhook identifier for webhook nodes.
	WebhookID types.String `tfsdk:"webhook_id"`

	// Disabled indicates if the node is disabled in the workflow.
	Disabled types.Bool `tfsdk:"disabled"`

	// Notes contains optional user notes about the node.
	Notes types.String `tfsdk:"notes"`

	// NodeJSON is the computed JSON representation of this node.
	// This is what gets inserted into the workflow's nodes array.
	NodeJSON types.String `tfsdk:"node_json"`
}
