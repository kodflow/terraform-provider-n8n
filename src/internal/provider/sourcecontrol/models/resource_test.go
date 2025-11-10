package models

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestPullResource(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create with all fields", wantErr: false},
		{name: "create with null values", wantErr: false},
		{name: "zero value struct", wantErr: false},
		{name: "force flag variations", wantErr: false},
		{name: "json field variations", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create with all fields":
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

			case "create with null values":
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

			case "zero value struct":
				var resource PullResource
				assert.True(t, resource.ID.IsNull())
				assert.True(t, resource.Force.IsNull())
				assert.True(t, resource.VariablesJSON.IsNull())
				assert.True(t, resource.ResultJSON.IsNull())

			case "force flag variations":
				resourceForced := PullResource{
					Force: types.BoolValue(true),
				}
				assert.True(t, resourceForced.Force.ValueBool())

				resourceNotForced := PullResource{
					Force: types.BoolValue(false),
				}
				assert.False(t, resourceNotForced.Force.ValueBool())

			case "json field variations":
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

			case "error case - validation checks":
				resource := PullResource{
					ID:            types.StringValue(""),
					Force:         types.BoolValue(false),
					VariablesJSON: types.StringValue("invalid json"),
					ResultJSON:    types.StringValue(""),
				}
				assert.Equal(t, "", resource.ID.ValueString())
				assert.False(t, resource.Force.ValueBool())
				assert.Equal(t, "invalid json", resource.VariablesJSON.ValueString())
			}
		})
	}
}
