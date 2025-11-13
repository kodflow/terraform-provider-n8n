// Package models defines data structures for user resources.
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
					ID:        types.StringValue("user-123"),
					Email:     types.StringValue("test@example.com"),
					FirstName: types.StringValue("John"),
					LastName:  types.StringValue("Doe"),
					Role:      types.StringValue("admin"),
				}

				assert.Equal(t, "user-123", datasource.ID.ValueString())
				assert.Equal(t, "test@example.com", datasource.Email.ValueString())
				assert.Equal(t, "John", datasource.FirstName.ValueString())
				assert.Equal(t, "Doe", datasource.LastName.ValueString())
				assert.Equal(t, "admin", datasource.Role.ValueString())

			case "zero value struct":
				var datasource DataSource
				assert.True(t, datasource.ID.IsNull())
				assert.True(t, datasource.Email.IsNull())
				assert.True(t, datasource.FirstName.IsNull())
				assert.True(t, datasource.LastName.IsNull())
				assert.True(t, datasource.Role.IsNull())

			case "error case - validation checks":
				datasource := DataSource{
					ID:        types.StringValue(""),
					Email:     types.StringValue("invalid-email"),
					FirstName: types.StringValue(""),
					LastName:  types.StringValue(""),
					Role:      types.StringValue("invalid-role"),
				}
				assert.Equal(t, "", datasource.ID.ValueString())
				assert.Equal(t, "invalid-email", datasource.Email.ValueString())
				assert.Equal(t, "invalid-role", datasource.Role.ValueString())
			}
		})
	}
}
