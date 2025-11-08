package models

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestDataSource(t *testing.T) {
	t.Run("create with all fields", func(t *testing.T) {
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
	})

	t.Run("zero value struct", func(t *testing.T) {
		var datasource DataSource
		assert.True(t, datasource.ID.IsNull())
		assert.True(t, datasource.Email.IsNull())
		assert.True(t, datasource.FirstName.IsNull())
		assert.True(t, datasource.LastName.IsNull())
		assert.True(t, datasource.Role.IsNull())
	})
}
