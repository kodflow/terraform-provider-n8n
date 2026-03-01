// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package datatable

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/datatable/models"
	"github.com/stretchr/testify/assert"
)

func Test_mapDataTableToDataSourceModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "map with all fields populated"},
		{name: "map with empty strings"},
		{name: "map preserves existing data fields"},
		{name: "map multiple times"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "map with all fields populated":
				table := &n8nsdk.DataTable{
					Id:        "dt-123",
					Name:      "Test Table",
					ProjectId: "proj-1",
					Columns:   []n8nsdk.DataTableColumnsInner{},
				}

				data := &models.DataSource{}
				mapDataTableToDataSourceModel(table, data)

				assert.Equal(t, "dt-123", data.ID.ValueString())
				assert.Equal(t, "Test Table", data.Name.ValueString())
				assert.Equal(t, "proj-1", data.ProjectID.ValueString())

			case "map with empty strings":
				table := &n8nsdk.DataTable{
					Id:        "",
					Name:      "",
					ProjectId: "",
					Columns:   nil,
				}

				data := &models.DataSource{}
				mapDataTableToDataSourceModel(table, data)

				assert.Equal(t, "", data.ID.ValueString())
				assert.Equal(t, "", data.Name.ValueString())
				assert.Equal(t, "", data.ProjectID.ValueString())

			case "map preserves existing data fields":
				table := &n8nsdk.DataTable{
					Id:        "dt-new",
					Name:      "New Name",
					ProjectId: "proj-new",
					Columns:   []n8nsdk.DataTableColumnsInner{},
				}

				data := &models.DataSource{
					ID:   types.StringValue("old-id"),
					Name: types.StringValue("Old Name"),
				}

				mapDataTableToDataSourceModel(table, data)

				assert.Equal(t, "dt-new", data.ID.ValueString())
				assert.Equal(t, "New Name", data.Name.ValueString())
				assert.Equal(t, "proj-new", data.ProjectID.ValueString())

			case "map multiple times":
				table1 := &n8nsdk.DataTable{
					Id:   "dt-1",
					Name: "First Table",
				}

				data := &models.DataSource{}
				mapDataTableToDataSourceModel(table1, data)
				assert.Equal(t, "First Table", data.Name.ValueString())

				table2 := &n8nsdk.DataTable{
					Id:   "dt-2",
					Name: "Second Table",
				}
				mapDataTableToDataSourceModel(table2, data)
				assert.Equal(t, "Second Table", data.Name.ValueString())
			}
		})
	}
}

func Test_findDataTableByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "find existing table"},
		{name: "table not found"},
		{name: "empty table list"},
		{name: "nil table list"},
		{name: "find first table"},
		{name: "find last table"},
		{name: "case sensitive search"},
		{name: "duplicate names returns first"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "find existing table":
				tables := []n8nsdk.DataTable{
					{Id: "dt-1", Name: "Production"},
					{Id: "dt-2", Name: "Development"},
					{Id: "dt-3", Name: "Testing"},
				}

				table, found := findDataTableByName(tables, "Development")

				assert.True(t, found)
				assert.NotNil(t, table)
				assert.Equal(t, "Development", table.Name)
				assert.Equal(t, "dt-2", table.Id)

			case "table not found":
				tables := []n8nsdk.DataTable{
					{Id: "dt-1", Name: "Production"},
				}

				table, found := findDataTableByName(tables, "NonExistent")

				assert.False(t, found)
				assert.Nil(t, table)

			case "empty table list":
				tables := []n8nsdk.DataTable{}

				table, found := findDataTableByName(tables, "AnyTable")

				assert.False(t, found)
				assert.Nil(t, table)

			case "nil table list":
				var tables []n8nsdk.DataTable

				table, found := findDataTableByName(tables, "AnyTable")

				assert.False(t, found)
				assert.Nil(t, table)

			case "find first table":
				tables := []n8nsdk.DataTable{
					{Id: "dt-1", Name: "First"},
					{Id: "dt-2", Name: "Second"},
				}

				table, found := findDataTableByName(tables, "First")

				assert.True(t, found)
				assert.NotNil(t, table)
				assert.Equal(t, "First", table.Name)
				assert.Equal(t, "dt-1", table.Id)

			case "find last table":
				tables := []n8nsdk.DataTable{
					{Id: "dt-1", Name: "First"},
					{Id: "dt-2", Name: "Last"},
				}

				table, found := findDataTableByName(tables, "Last")

				assert.True(t, found)
				assert.NotNil(t, table)
				assert.Equal(t, "Last", table.Name)
				assert.Equal(t, "dt-2", table.Id)

			case "case sensitive search":
				tables := []n8nsdk.DataTable{
					{Id: "dt-1", Name: "Production"},
				}

				table, found := findDataTableByName(tables, "production")

				assert.False(t, found)
				assert.Nil(t, table)

			case "duplicate names returns first":
				tables := []n8nsdk.DataTable{
					{Id: "dt-1", Name: "Duplicate"},
					{Id: "dt-2", Name: "Duplicate"},
				}

				table, found := findDataTableByName(tables, "Duplicate")

				assert.True(t, found)
				assert.NotNil(t, table)
				assert.Equal(t, "dt-1", table.Id)
			}
		})
	}
}

func Test_mapColumnsToModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		cols     []n8nsdk.DataTableColumnsInner
		wantLen  int
		wantNil  bool
		checkVal func(t *testing.T, out []models.Column)
	}{
		{
			name:    "nil slice returns empty slice",
			cols:    nil,
			wantLen: 0,
		},
		{
			name:    "empty slice returns empty slice",
			cols:    []n8nsdk.DataTableColumnsInner{},
			wantLen: 0,
		},
		{
			name: "single column",
			cols: func() []n8nsdk.DataTableColumnsInner {
				n := "col1"
				ty := "string"
				return []n8nsdk.DataTableColumnsInner{
					{Name: &n, Type: &ty},
				}
			}(),
			wantLen: 1,
			checkVal: func(t *testing.T, out []models.Column) {
				t.Helper()
				assert.Equal(t, "col1", out[0].Name.ValueString())
				assert.Equal(t, "string", out[0].Type.ValueString())
			},
		},
		{
			name: "multiple columns",
			cols: func() []n8nsdk.DataTableColumnsInner {
				n1, t1 := "a", "string"
				n2, t2 := "b", "number"
				n3, t3 := "c", "boolean"
				return []n8nsdk.DataTableColumnsInner{
					{Name: &n1, Type: &t1},
					{Name: &n2, Type: &t2},
					{Name: &n3, Type: &t3},
				}
			}(),
			wantLen: 3,
			checkVal: func(t *testing.T, out []models.Column) {
				t.Helper()
				assert.Equal(t, "a", out[0].Name.ValueString())
				assert.Equal(t, "number", out[1].Type.ValueString())
				assert.Equal(t, "c", out[2].Name.ValueString())
			},
		},
		{
			name: "column with nil name and type",
			cols: []n8nsdk.DataTableColumnsInner{
				{Name: nil, Type: nil},
			},
			wantLen: 1,
			checkVal: func(t *testing.T, out []models.Column) {
				t.Helper()
				assert.Equal(t, "", out[0].Name.ValueString())
				assert.Equal(t, "", out[0].Type.ValueString())
			},
		},
		{
			name: "columns ordered by index when input order differs",
			cols: func() []n8nsdk.DataTableColumnsInner {
				n1, t1 := "first", "string"
				n2, t2 := "second", "number"
				n3, t3 := "third", "boolean"
				idx0 := int32(0)
				idx1 := int32(1)
				idx2 := int32(2)
				return []n8nsdk.DataTableColumnsInner{
					{Name: &n3, Type: &t3, Index: &idx2},
					{Name: &n1, Type: &t1, Index: &idx0},
					{Name: &n2, Type: &t2, Index: &idx1},
				}
			}(),
			wantLen: 3,
			checkVal: func(t *testing.T, out []models.Column) {
				t.Helper()
				assert.Equal(t, "first", out[0].Name.ValueString())
				assert.Equal(t, "second", out[1].Name.ValueString())
				assert.Equal(t, "third", out[2].Name.ValueString())
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			out := mapColumnsToModel(tt.cols)

			if tt.wantNil {
				assert.Nil(t, out)
				return
			}
			assert.Len(t, out, tt.wantLen)
			if tt.checkVal != nil {
				tt.checkVal(t, out)
			}
		})
	}
}
