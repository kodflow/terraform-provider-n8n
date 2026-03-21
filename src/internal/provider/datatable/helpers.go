// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package datatable contains helper functions for data table operations.
package datatable

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/datatable/models"
)

// mapDataTableToResourceModel maps an SDK DataTable to the resource model.
//
// Params:
//   - dt: The n8n SDK DataTable to map
//   - data: The models.Resource to populate
func mapDataTableToResourceModel(dt *n8nsdk.DataTable, data *models.Resource) {
	data.ID = types.StringValue(dt.ID)
	data.Name = types.StringValue(dt.Name)
	data.ProjectID = types.StringValue(dt.ProjectID)
	data.CreatedAt = types.StringValue(dt.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.UpdatedAt = types.StringValue(dt.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
}

// mapDataTableToDataSourceModel maps an SDK DataTable to the datasource model.
//
// Params:
//   - dt: The n8n SDK DataTable to map
//   - data: The models.DataSource to populate
func mapDataTableToDataSourceModel(dt *n8nsdk.DataTable, data *models.DataSource) {
	data.ID = types.StringValue(dt.ID)
	data.Name = types.StringValue(dt.Name)
	data.ProjectID = types.StringValue(dt.ProjectID)
	data.CreatedAt = types.StringValue(dt.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.UpdatedAt = types.StringValue(dt.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.Columns = mapDataTableColumnsToDataSource(dt.Columns)
}

// mapDataTableToItem maps an SDK DataTable to an Item model for list datasources.
//
// Params:
//   - dt: The n8n SDK DataTable to map
//
// Returns:
//   - models.Item: The mapped item model
func mapDataTableToItem(dt *n8nsdk.DataTable) (item models.Item) {
	item = models.Item{
		ID:        types.StringValue(dt.ID),
		Name:      types.StringValue(dt.Name),
		ProjectID: types.StringValue(dt.ProjectID),
		CreatedAt: types.StringValue(dt.CreatedAt.Format("2006-01-02T15:04:05Z07:00")),
		UpdatedAt: types.StringValue(dt.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")),
		Columns:   mapDataTableColumnsToDataSource(dt.Columns),
	}
	//: Return the fully populated item.
	return item
}

// mapDataTableColumnsToDataSource maps SDK DataTableColumn slice to datasource column models.
//
// Params:
//   - columns: The SDK column definitions to map
//
// Returns:
//   - []models.DataSourceColumnModel: The mapped column models
func mapDataTableColumnsToDataSource(columns []n8nsdk.DataTableColumn) (items []models.DataSourceColumnModel) {
	result := make([]models.DataSourceColumnModel, 0, len(columns))
	//: Iterate over each column and build the column model.
	for _, col := range columns {
		colModel := models.DataSourceColumnModel{
			Name: types.StringValue(col.Name),
			Type: types.StringValue(col.Type),
		}
		//: Set the column ID from the pointer value or null when absent.
		if col.ID != nil {
			colModel.ID = types.StringValue(*col.ID)
			//: Use null when the column ID pointer is absent.
		} else {
			colModel.ID = types.StringNull()
		}
		//: Set the column index from the pointer value or null when absent.
		if col.Index != nil {
			colModel.Index = types.Int64Value(*col.Index)
			//: Use null when the column index pointer is absent.
		} else {
			colModel.Index = types.Int64Null()
		}
		result = append(result, colModel)
	}
	//: Return the fully populated column model slice.
	return result
}

// findDataTableByName searches for a data table by name in a data table list.
//
// Params:
//   - tables: Slice of n8n SDK data tables to search through
//   - name: The data table name to search for
//
// Returns:
//   - *n8nsdk.DataTable: Pointer to the found data table, or nil if not found
//   - bool: True if data table was found, false otherwise
func findDataTableByName(tables []n8nsdk.DataTable, name string) (dataTable *n8nsdk.DataTable, ok bool) {
	//: Iterate over each table to find a name match.
	for _, dt := range tables {
		//: Return the matching table when the name equals the search term.
		if dt.Name == name {
			//: Return the matched table pointer and success flag.
			return &dt, true
		}
	}
	//: Return nil to indicate no matching table was found.
	return nil, false
}
