// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models contains data models for the datatable domain.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSourceColumnModel describes a column in a data table datasource.
// It holds the column identifier, position index, name, and data type.
type DataSourceColumnModel struct {
	ID    types.String `tfsdk:"id" dto:"out,query,priv" json:"id"`
	Index types.Int64  `tfsdk:"index" json:"index"`
	Name  types.String `tfsdk:"name" json:"name"`
	Type  types.String `tfsdk:"type" json:"type"`
}

// DataSource describes the single data table datasource model.
// It holds the identifier, name, project association, timestamps, and column definitions.
type DataSource struct {
	ID        types.String            `tfsdk:"id" dto:"out,query,priv" json:"id"`
	Name      types.String            `tfsdk:"name" json:"name"`
	Columns   []DataSourceColumnModel `tfsdk:"columns" json:"columns"`
	ProjectID types.String            `tfsdk:"project_id" json:"project_id"`
	CreatedAt types.String            `tfsdk:"created_at" json:"created_at"`
	UpdatedAt types.String            `tfsdk:"updated_at" json:"updated_at"`
}
