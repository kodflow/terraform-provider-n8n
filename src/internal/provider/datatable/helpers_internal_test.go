package datatable

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/datatable/models"
	"github.com/stretchr/testify/assert"
)

// TestFindDataTableByName tests the findDataTableByName helper.
func TestFindDataTableByName(t *testing.T) {
	t.Parallel()

	tables := []n8nsdk.DataTable{
		{ID: "dt-1", Name: "Alpha"},
		{ID: "dt-2", Name: "Beta"},
	}

	tests := []struct {
		name      string
		search    string
		expectID  string
		expectHit bool
	}{
		{
			name:      "found first",
			search:    "Alpha",
			expectID:  "dt-1",
			expectHit: true,
		},
		{
			name:      "found second",
			search:    "Beta",
			expectID:  "dt-2",
			expectHit: true,
		},
		{
			name:      "not found",
			search:    "Gamma",
			expectHit: false,
		},
		{
			name:      "empty search",
			search:    "",
			expectHit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, found := findDataTableByName(tables, tt.search)

			assert.Equal(t, tt.expectHit, found)
			// Check condition.
			if tt.expectHit {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectID, result.ID)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

// TestMapDataTableColumnsToDataSource tests the mapDataTableColumnsToDataSource helper.
func TestMapDataTableColumnsToDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		columns     []n8nsdk.DataTableColumn
		expectCount int
		expectID    string
	}{
		{
			name: "columns with ID and index",
			columns: []n8nsdk.DataTableColumn{
				{ID: new("col-id-1"), Index: new(int64(0)), Name: "col1", Type: "string"},
			},
			expectCount: 1,
			expectID:    "col-id-1",
		},
		{
			name: "columns without ID and index",
			columns: []n8nsdk.DataTableColumn{
				{Name: "col2", Type: "number"},
			},
			expectCount: 1,
			expectID:    "",
		},
		{
			name:        "empty columns",
			columns:     []n8nsdk.DataTableColumn{},
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := mapDataTableColumnsToDataSource(tt.columns)

			assert.Len(t, result, tt.expectCount)
			// Check condition.
			if tt.expectCount > 0 && tt.expectID != "" {
				assert.Equal(t, tt.expectID, result[0].ID.ValueString())
			}
		})
	}
}

// TestMapDataTableToItem tests the mapDataTableToItem helper.
func TestMapDataTableToItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		dt         *n8nsdk.DataTable
		expectID   string
		expectName string
	}{
		{
			name: "basic table",
			dt: &n8nsdk.DataTable{
				ID:        "dt-1",
				Name:      "Test Table",
				Columns:   []n8nsdk.DataTableColumn{},
				ProjectID: "proj-1",
			},
			expectID:   "dt-1",
			expectName: "Test Table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			item := mapDataTableToItem(tt.dt)

			assert.Equal(t, tt.expectID, item.ID.ValueString())
			assert.Equal(t, tt.expectName, item.Name.ValueString())
		})
	}
}

// TestMapDataTableToDataSourceModel tests the mapDataTableToDataSourceModel helper.
func TestMapDataTableToDataSourceModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		dt              *n8nsdk.DataTable
		expectID        string
		expectName      string
		expectProjectID string
	}{
		{
			name: "maps all fields correctly",
			dt: &n8nsdk.DataTable{
				ID:        "dt-1",
				Name:      "Test Table",
				Columns:   []n8nsdk.DataTableColumn{},
				ProjectID: "proj-1",
			},
			expectID:        "dt-1",
			expectName:      "Test Table",
			expectProjectID: "proj-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := &models.DataSource{}
			mapDataTableToDataSourceModel(tt.dt, data)
			assert.Equal(t, tt.expectID, data.ID.ValueString())
			assert.Equal(t, tt.expectName, data.Name.ValueString())
			assert.Equal(t, tt.expectProjectID, data.ProjectID.ValueString())
		})
	}
}

// TestMapDataTableToResourceModel tests the mapDataTableToResourceModel helper.
func TestMapDataTableToResourceModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		dt              *n8nsdk.DataTable
		expectID        string
		expectName      string
		expectProjectID string
	}{
		{
			name: "maps all fields correctly",
			dt: &n8nsdk.DataTable{
				ID:        "dt-1",
				Name:      "Test Table",
				Columns:   []n8nsdk.DataTableColumn{},
				ProjectID: "proj-2",
			},
			expectID:        "dt-1",
			expectName:      "Test Table",
			expectProjectID: "proj-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := &models.Resource{}
			mapDataTableToResourceModel(tt.dt, data)
			assert.Equal(t, tt.expectID, data.ID.ValueString())
			assert.Equal(t, tt.expectName, data.Name.ValueString())
			assert.Equal(t, tt.expectProjectID, data.ProjectID.ValueString())
		})
	}
}

// TestParseExecutionID_ViaConstants checks Float32BitSize constant is usable.
func TestFindDataTableByName_EmptySlice(t *testing.T) {
	t.Parallel()

	result, found := findDataTableByName([]n8nsdk.DataTable{}, "test")

	assert.False(t, found)
	assert.Nil(t, result)
}

// TestMapDataTableColumnsToDataSource_NilID tests nil ID handling.
func TestMapDataTableColumnsToDataSource_NilID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		cols       []n8nsdk.DataTableColumn
		wantLen    int
		wantNilID  bool
		wantNilIdx bool
	}{
		{
			name:       "nil ID and index produce null values",
			cols:       []n8nsdk.DataTableColumn{{ID: nil, Index: nil, Name: "col", Type: "string"}},
			wantLen:    1,
			wantNilID:  true,
			wantNilIdx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := mapDataTableColumnsToDataSource(tt.cols)
			assert.Len(t, result, tt.wantLen)
			if tt.wantNilID {
				assert.Equal(t, types.StringNull(), result[0].ID)
			}
			if tt.wantNilIdx {
				assert.Equal(t, types.Int64Null(), result[0].Index)
			}
		})
	}
}
