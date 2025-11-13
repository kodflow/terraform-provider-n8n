// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package models defines execution list data structures.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSources maps the Terraform schema to the datasource response.
// It represents the filtered execution list with workflow and execution details from the n8n API.
type DataSources struct {
	WorkflowID  types.String `tfsdk:"workflow_id"`
	ProjectID   types.String `tfsdk:"project_id"`
	Status      types.String `tfsdk:"status"`
	IncludeData types.Bool   `tfsdk:"include_data"`
	Executions  []Item       `tfsdk:"executions"`
}
