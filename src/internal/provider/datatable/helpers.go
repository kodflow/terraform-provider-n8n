// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package datatable implements data table management resources.
package datatable

import (
	"math"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/datatable/models"
)

// mapDataTableToDataSourceModel maps an SDK data table to the datasource model.
func mapDataTableToDataSourceModel(table *n8nsdk.DataTable, data *models.DataSource) {
	data.ID = types.StringValue(table.Id)
	data.Name = types.StringValue(table.Name)
	data.ProjectID = types.StringValue(table.ProjectId)
}

// findDataTableByName searches for a data table by name in a list.
func findDataTableByName(tables []n8nsdk.DataTable, name string) (*n8nsdk.DataTable, bool) {
	for _, table := range tables {
		if table.Name == name {
			return &table, true
		}
	}
	return nil, false
}

// mapColumnsToModel maps SDK column slice to model column slice.
// Columns are ordered by DataTableColumnsInner.Index so that parameter order is stable.
func mapColumnsToModel(cols []n8nsdk.DataTableColumnsInner) []models.Column {
	if len(cols) == 0 {
		return []models.Column{}
	}

	sorted := make([]n8nsdk.DataTableColumnsInner, len(cols))
	copy(sorted, cols)
	var defaultIdx int32 = math.MaxInt32
	if n := len(cols); n <= math.MaxInt32 {
		defaultIdx = int32(n)
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		idxI := defaultIdx
		if sorted[i].Index != nil {
			idxI = *sorted[i].Index
		}
		idxJ := defaultIdx
		if sorted[j].Index != nil {
			idxJ = *sorted[j].Index
		}
		return idxI < idxJ
	})

	out := make([]models.Column, 0, len(sorted))
	for _, c := range sorted {
		name := ""
		if c.Name != nil {
			name = *c.Name
		}
		typ := ""
		if c.Type != nil {
			typ = *c.Type
		}
		out = append(out, models.Column{
			Name: types.StringValue(name),
			Type: types.StringValue(typ),
		})
	}
	return out
}
