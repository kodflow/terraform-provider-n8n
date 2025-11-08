package models

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestDataSource(t *testing.T) {
	t.Run("create with all fields", func(t *testing.T) {
		datasource := DataSource{
			ID:        types.StringValue("tag-123"),
			Name:      types.StringValue("Test Tag"),
			CreatedAt: types.StringValue("2024-01-01T00:00:00Z"),
			UpdatedAt: types.StringValue("2024-01-01T00:00:00Z"),
		}

		assert.Equal(t, "tag-123", datasource.ID.ValueString())
		assert.Equal(t, "Test Tag", datasource.Name.ValueString())
	})

	t.Run("zero value struct", func(t *testing.T) {
		var datasource DataSource
		assert.True(t, datasource.ID.IsNull())
		assert.True(t, datasource.Name.IsNull())
	})
}
