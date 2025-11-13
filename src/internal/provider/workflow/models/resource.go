// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package models contains data models for the workflow domain.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Resource describes the workflow resource data model.
// Maps n8n workflow attributes to Terraform schema, including nodes, connections, and settings.
type Resource struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Active          types.Bool   `tfsdk:"active"`
	Tags            types.List   `tfsdk:"tags"`
	NodesJSON       types.String `tfsdk:"nodes_json"`
	ConnectionsJSON types.String `tfsdk:"connections_json"`
	SettingsJSON    types.String `tfsdk:"settings_json"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
	VersionID       types.String `tfsdk:"version_id"`
	IsArchived      types.Bool   `tfsdk:"is_archived"`
	TriggerCount    types.Int64  `tfsdk:"trigger_count"`
	Meta            types.Map    `tfsdk:"meta"`
	PinData         types.Map    `tfsdk:"pin_data"`
}
