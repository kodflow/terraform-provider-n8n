// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models contains data models for the datatable domain.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataTableColumnModel describes a column in a data table resource.
// It holds the column name and its data type for create/update operations.
type DataTableColumnModel struct {
	Name types.String `tfsdk:"name" dto:"inout,priv,pub" json:"name"`
	Type types.String `tfsdk:"type" json:"type"`
}

// Resource describes the data table resource data model.
// It holds the identifier, name, project association, timestamps, and column definitions.
type Resource struct {
	ID        types.String           `tfsdk:"id" dto:"inout,priv,priv" json:"id"`
	Name      types.String           `tfsdk:"name" json:"name"`
	Columns   []DataTableColumnModel `tfsdk:"columns" json:"columns"`
	ProjectID types.String           `tfsdk:"project_id" json:"project_id"`
	CreatedAt types.String           `tfsdk:"created_at" json:"created_at"`
	UpdatedAt types.String           `tfsdk:"updated_at" json:"updated_at"`
}
