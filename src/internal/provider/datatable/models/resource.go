// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models contains data models for the datatable domain.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Column describes a single column in a data table.
type Column struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

// Resource describes the data table resource data model.
// Maps n8n data table attributes to Terraform schema for managing workspace data tables.
type Resource struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Columns   []Column     `tfsdk:"columns"`
	ProjectID types.String `tfsdk:"project_id"`
}
