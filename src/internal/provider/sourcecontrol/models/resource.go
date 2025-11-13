// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package models contains data models for the sourcecontrol domain.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Package sourcecontrol contains resources and datasources for sourcecontrol management

// PullResource describes the source control pull resource data model.
// Maps n8n source control pull operation attributes to Terraform schema for pulling changes from remote repositories.
type PullResource struct {
	ID            types.String `tfsdk:"id"`
	Force         types.Bool   `tfsdk:"force"`
	VariablesJSON types.String `tfsdk:"variables_json"`
	ResultJSON    types.String `tfsdk:"result_json"`
}
