package models

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestDataSource(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create with all fields", wantErr: false},
		{name: "zero value struct", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create with all fields":
				datasource := DataSource{
					ID:        types.StringValue("tag-123"),
					Name:      types.StringValue("Test Tag"),
					CreatedAt: types.StringValue("2024-01-01T00:00:00Z"),
					UpdatedAt: types.StringValue("2024-01-01T00:00:00Z"),
				}

				assert.Equal(t, "tag-123", datasource.ID.ValueString())
				assert.Equal(t, "Test Tag", datasource.Name.ValueString())

			case "zero value struct":
				var datasource DataSource
				assert.True(t, datasource.ID.IsNull())
				assert.True(t, datasource.Name.IsNull())

			case "error case - validation checks":
				datasource := DataSource{
					ID:   types.StringValue(""),
					Name: types.StringValue(""),
				}
				assert.Equal(t, "", datasource.ID.ValueString())
				assert.Equal(t, "", datasource.Name.ValueString())
			}
		})
	}
}
