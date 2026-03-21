// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models contains data models for the datatable domain.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Item describes a single data table in the list datasource.
// It holds the identifier, name, project association, timestamps, and column definitions.
type Item struct {
	ID        types.String            `tfsdk:"id" dto:"out,query,priv" json:"id"`
	Name      types.String            `tfsdk:"name" json:"name"`
	Columns   []DataSourceColumnModel `tfsdk:"columns" json:"columns"`
	ProjectID types.String            `tfsdk:"project_id" json:"project_id"`
	CreatedAt types.String            `tfsdk:"created_at" json:"created_at"`
	UpdatedAt types.String            `tfsdk:"updated_at" json:"updated_at"`
}

// DataSources describes the list data tables datasource model.
// It holds the collection of data table items returned by the n8n API.
type DataSources struct {
	DataTables []Item `tfsdk:"data_tables" dto:"out,query,priv" json:"data_tables"`
}
