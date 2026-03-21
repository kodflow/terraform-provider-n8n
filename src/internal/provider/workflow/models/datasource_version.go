// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models defines data structures for workflow resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSourceVersion describes the workflow version datasource model.
// It maps a specific workflow version from the n8n API to Terraform schema.
type DataSourceVersion struct {
	WorkflowID      types.String `tfsdk:"workflow_id"`
	VersionID       types.String `tfsdk:"version_id"`
	Name            types.String `tfsdk:"name"`
	Authors         types.String `tfsdk:"authors"`
	NodesJSON       types.String `tfsdk:"nodes_json"`
	ConnectionsJSON types.String `tfsdk:"connections_json"`
	CreatedAt       types.String `tfsdk:"created_at"`
}
