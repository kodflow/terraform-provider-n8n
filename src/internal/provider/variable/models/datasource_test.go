package models

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestDataSource(t *testing.T) {
	t.Run("create with all fields", func(t *testing.T) {
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
	})

	t.Run("zero value struct", func(t *testing.T) {
		var datasource DataSource
		assert.True(t, datasource.ID.IsNull())
		assert.True(t, datasource.Key.IsNull())
		assert.True(t, datasource.Value.IsNull())
		assert.True(t, datasource.Type.IsNull())
	})
}
