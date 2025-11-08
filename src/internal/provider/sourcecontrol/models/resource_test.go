package models

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestPullResource(t *testing.T) {
	t.Run("create with all fields", func(t *testing.T) {
		resource := PullResource{
			ID:            types.StringValue("pull-123"),
			Force:         types.BoolValue(true),
			VariablesJSON: types.StringValue(`{"key":"value"}`),
			ResultJSON:    types.StringValue(`{"status":"success"}`),
		}

		assert.Equal(t, "pull-123", resource.ID.ValueString())
		assert.True(t, resource.Force.ValueBool())
		assert.Equal(t, `{"key":"value"}`, resource.VariablesJSON.ValueString())
		assert.Equal(t, `{"status":"success"}`, resource.ResultJSON.ValueString())
	})

	t.Run("create with null values", func(t *testing.T) {
		resource := PullResource{
			ID:            types.StringNull(),
			Force:         types.BoolNull(),
			VariablesJSON: types.StringNull(),
			ResultJSON:    types.StringNull(),
		}

		assert.True(t, resource.ID.IsNull())
		assert.True(t, resource.Force.IsNull())
		assert.True(t, resource.VariablesJSON.IsNull())
		assert.True(t, resource.ResultJSON.IsNull())
	})

	t.Run("zero value struct", func(t *testing.T) {
		var resource PullResource
		assert.True(t, resource.ID.IsNull())
		assert.True(t, resource.Force.IsNull())
		assert.True(t, resource.VariablesJSON.IsNull())
		assert.True(t, resource.ResultJSON.IsNull())
	})

	t.Run("force flag variations", func(t *testing.T) {
		resourceForced := PullResource{
			Force: types.BoolValue(true),
		}
		assert.True(t, resourceForced.Force.ValueBool())

		resourceNotForced := PullResource{
			Force: types.BoolValue(false),
		}
		assert.False(t, resourceNotForced.Force.ValueBool())
	})

	t.Run("json field variations", func(t *testing.T) {
		jsonVariations := []string{
			`{}`,
			`{"simple":"value"}`,
			`{"nested":{"key":"value"}}`,
			`["array","of","values"]`,
			`null`,
			`"string value"`,
			`123`,
		}

		for _, json := range jsonVariations {
			resource := PullResource{
				VariablesJSON: types.StringValue(json),
				ResultJSON:    types.StringValue(json),
			}
			assert.Equal(t, json, resource.VariablesJSON.ValueString())
			assert.Equal(t, json, resource.ResultJSON.ValueString())
		}
	})
}
