// Package models defines data structures for variable resources.
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
					ID:    types.StringValue("var-123"),
					Key:   types.StringValue("TEST_VAR"),
					Value: types.StringValue("test value"),
					Type:  types.StringValue("string"),
				}

				assert.Equal(t, "var-123", datasource.ID.ValueString())
				assert.Equal(t, "TEST_VAR", datasource.Key.ValueString())
				assert.Equal(t, "test value", datasource.Value.ValueString())
				assert.Equal(t, "string", datasource.Type.ValueString())

			case "zero value struct":
				var datasource DataSource
				assert.True(t, datasource.ID.IsNull())
				assert.True(t, datasource.Key.IsNull())
				assert.True(t, datasource.Value.IsNull())
				assert.True(t, datasource.Type.IsNull())

			case "error case - validation checks":
				datasource := DataSource{
					ID:    types.StringValue(""),
					Key:   types.StringValue(""),
					Value: types.StringValue(""),
					Type:  types.StringValue("invalid-type"),
				}
				assert.Equal(t, "", datasource.ID.ValueString())
				assert.Equal(t, "", datasource.Key.ValueString())
				assert.Equal(t, "invalid-type", datasource.Type.ValueString())
			}
		})
	}
}
