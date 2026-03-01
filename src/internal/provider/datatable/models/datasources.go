// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models contains data models for the datatable domain.
package models

// DataSources maps Terraform schema attributes for data table list data.
type DataSources struct {
	DataTables []Item `tfsdk:"data_tables"`
}
