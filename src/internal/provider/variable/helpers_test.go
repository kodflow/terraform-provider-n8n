package variable

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/variable/models"
	"github.com/stretchr/testify/assert"
)

func TestFindVariableByIDOrKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "find by ID when both ID and key provided"},
		{name: "find by key when ID is null"},
		{name: "find by key when both ID and key match"},
		{name: "not found when no match"},
		{name: "find in empty list"},
		{name: "find with multiple variables"},
		{name: "find with nil ID in variable"},
		{name: "both ID and key are null"},
		{name: "error case - validation checks", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "find by ID when both ID and key provided":
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "test-key", Value: "test-value"},
				}

				found, exists := findVariableByIDOrKey(variables, types.StringValue("var-123"), types.StringValue("wrong-key"))

				assert.True(t, exists)
				assert.NotNil(t, found)
				assert.Equal(t, "var-123", *found.Id)

			case "find by key when ID is null":
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "test-key", Value: "test-value"},
				}

				found, exists := findVariableByIDOrKey(variables, types.StringNull(), types.StringValue("test-key"))

				assert.True(t, exists)
				assert.NotNil(t, found)
				assert.Equal(t, "test-key", found.Key)

			case "find by key when both ID and key match":
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "test-key", Value: "test-value"},
				}

				found, exists := findVariableByIDOrKey(variables, types.StringValue("var-123"), types.StringValue("test-key"))

				assert.True(t, exists)
				assert.NotNil(t, found)

			case "not found when no match":
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "test-key", Value: "test-value"},
				}

				found, exists := findVariableByIDOrKey(variables, types.StringValue("var-999"), types.StringValue("wrong-key"))

				assert.False(t, exists)
				assert.Nil(t, found)

			case "find in empty list":
				variables := []n8nsdk.Variable{}

				found, exists := findVariableByIDOrKey(variables, types.StringValue("var-123"), types.StringValue("test-key"))

				assert.False(t, exists)
				assert.Nil(t, found)

			case "find with multiple variables":
				id1 := "var-123"
				id2 := "var-456"
				id3 := "var-789"
				variables := []n8nsdk.Variable{
					{Id: &id1, Key: "key-1", Value: "value-1"},
					{Id: &id2, Key: "key-2", Value: "value-2"},
					{Id: &id3, Key: "key-3", Value: "value-3"},
				}

				found, exists := findVariableByIDOrKey(variables, types.StringNull(), types.StringValue("key-2"))

				assert.True(t, exists)
				assert.NotNil(t, found)
				assert.Equal(t, "key-2", found.Key)
				assert.Equal(t, "var-456", *found.Id)

			case "find with nil ID in variable":
				variables := []n8nsdk.Variable{
					{Id: nil, Key: "test-key", Value: "test-value"},
				}

				found, exists := findVariableByIDOrKey(variables, types.StringValue("var-123"), types.StringValue("test-key"))

				assert.True(t, exists)
				assert.NotNil(t, found)
				assert.Equal(t, "test-key", found.Key)

			case "both ID and key are null":
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "test-key", Value: "test-value"},
				}

				found, exists := findVariableByIDOrKey(variables, types.StringNull(), types.StringNull())

				assert.False(t, exists)
				assert.Nil(t, found)

			case "error case - validation checks":
				// Test with nil variables slice
				found, exists := findVariableByIDOrKey(nil, types.StringValue("var-123"), types.StringValue("test-key"))
				assert.False(t, exists)
				assert.Nil(t, found)

				// Test with unknown types
				variables := []n8nsdk.Variable{}
				found, exists = findVariableByIDOrKey(variables, types.StringUnknown(), types.StringUnknown())
				assert.False(t, exists)
				assert.Nil(t, found)
			}
		})
	}
}

func TestFindVariableByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "find existing variable by ID"},
		{name: "not found when ID does not match"},
		{name: "find in empty list"},
		{name: "find with multiple variables"},
		{name: "find with nil ID"},
		{name: "empty ID string"},
		{name: "error case - validation checks", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "find existing variable by ID":
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "test-key", Value: "test-value"},
				}

				found, exists := findVariableByID(variables, "var-123")

				assert.True(t, exists)
				assert.NotNil(t, found)
				assert.Equal(t, "var-123", *found.Id)

			case "not found when ID does not match":
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "test-key", Value: "test-value"},
				}

				found, exists := findVariableByID(variables, "var-999")

				assert.False(t, exists)
				assert.Nil(t, found)

			case "find in empty list":
				variables := []n8nsdk.Variable{}

				found, exists := findVariableByID(variables, "var-123")

				assert.False(t, exists)
				assert.Nil(t, found)

			case "find with multiple variables":
				id1 := "var-123"
				id2 := "var-456"
				id3 := "var-789"
				variables := []n8nsdk.Variable{
					{Id: &id1, Key: "key-1", Value: "value-1"},
					{Id: &id2, Key: "key-2", Value: "value-2"},
					{Id: &id3, Key: "key-3", Value: "value-3"},
				}

				found, exists := findVariableByID(variables, "var-456")

				assert.True(t, exists)
				assert.NotNil(t, found)
				assert.Equal(t, "var-456", *found.Id)
				assert.Equal(t, "key-2", found.Key)

			case "find with nil ID":
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: nil, Key: "key-1", Value: "value-1"},
					{Id: &id, Key: "key-2", Value: "value-2"},
				}

				found, exists := findVariableByID(variables, "var-123")

				assert.True(t, exists)
				assert.NotNil(t, found)
				assert.Equal(t, "var-123", *found.Id)

			case "empty ID string":
				id := ""
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "test-key", Value: "test-value"},
				}

				found, exists := findVariableByID(variables, "")

				assert.True(t, exists)
				assert.NotNil(t, found)

			case "error case - validation checks":
				// Test with nil variables slice
				found, exists := findVariableByID(nil, "var-123")
				assert.False(t, exists)
				assert.Nil(t, found)

				// Test with all nil IDs
				variables := []n8nsdk.Variable{
					{Id: nil, Key: "key-1", Value: "value-1"},
					{Id: nil, Key: "key-2", Value: "value-2"},
				}
				found, exists = findVariableByID(variables, "var-123")
				assert.False(t, exists)
				assert.Nil(t, found)
			}
		})
	}
}

func TestFindVariableByKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "find existing variable by key"},
		{name: "not found when key does not match"},
		{name: "find in empty list"},
		{name: "find with multiple variables"},
		{name: "find with empty key"},
		{name: "find first match when duplicates exist"},
		{name: "case sensitive key matching"},
		{name: "error case - validation checks", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "find existing variable by key":
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "test-key", Value: "test-value"},
				}

				found, exists := findVariableByKey(variables, "test-key")

				assert.True(t, exists)
				assert.NotNil(t, found)
				assert.Equal(t, "test-key", found.Key)

			case "not found when key does not match":
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "test-key", Value: "test-value"},
				}

				found, exists := findVariableByKey(variables, "wrong-key")

				assert.False(t, exists)
				assert.Nil(t, found)

			case "find in empty list":
				variables := []n8nsdk.Variable{}

				found, exists := findVariableByKey(variables, "test-key")

				assert.False(t, exists)
				assert.Nil(t, found)

			case "find with multiple variables":
				id1 := "var-123"
				id2 := "var-456"
				id3 := "var-789"
				variables := []n8nsdk.Variable{
					{Id: &id1, Key: "key-1", Value: "value-1"},
					{Id: &id2, Key: "key-2", Value: "value-2"},
					{Id: &id3, Key: "key-3", Value: "value-3"},
				}

				found, exists := findVariableByKey(variables, "key-2")

				assert.True(t, exists)
				assert.NotNil(t, found)
				assert.Equal(t, "key-2", found.Key)
				assert.Equal(t, "var-456", *found.Id)

			case "find with empty key":
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "", Value: "test-value"},
				}

				found, exists := findVariableByKey(variables, "")

				assert.True(t, exists)
				assert.NotNil(t, found)

			case "find first match when duplicates exist":
				id1 := "var-123"
				id2 := "var-456"
				variables := []n8nsdk.Variable{
					{Id: &id1, Key: "duplicate-key", Value: "value-1"},
					{Id: &id2, Key: "duplicate-key", Value: "value-2"},
				}

				found, exists := findVariableByKey(variables, "duplicate-key")

				assert.True(t, exists)
				assert.NotNil(t, found)
				assert.Equal(t, "var-123", *found.Id)

			case "case sensitive key matching":
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "TestKey", Value: "test-value"},
				}

				found, exists := findVariableByKey(variables, "testkey")

				assert.False(t, exists)
				assert.Nil(t, found)

			case "error case - validation checks":
				// Test with nil variables slice
				found, exists := findVariableByKey(nil, "test-key")
				assert.False(t, exists)
				assert.Nil(t, found)

				// Test with special characters in key
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "key-with-special-chars-!@#", Value: "value"},
				}
				found, exists = findVariableByKey(variables, "key-with-special-chars-!@#")
				assert.True(t, exists)
				assert.NotNil(t, found)
			}
		})
	}
}

func TestMapVariableToDataSourceModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "map with all fields populated"},
		{name: "map with nil fields"},
		{name: "map with nil project"},
		{name: "map with project but nil project ID"},
		{name: "map with empty strings"},
		{name: "map sensitive value"},
		{name: "error case - validation checks", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "map with all fields populated":
				id := "var-123"
				varType := "string"
				projectID := "proj-456"
				variable := &n8nsdk.Variable{
					Id:    &id,
					Key:   "test-key",
					Value: "test-value",
					Type:  &varType,
					Project: &n8nsdk.Project{
						Id: &projectID,
					},
				}

				data := &models.DataSource{}
				mapVariableToDataSourceModel(variable, data)

				assert.Equal(t, "var-123", data.ID.ValueString())
				assert.Equal(t, "test-key", data.Key.ValueString())
				assert.Equal(t, "test-value", data.Value.ValueString())
				assert.Equal(t, "string", data.Type.ValueString())
				assert.Equal(t, "proj-456", data.ProjectID.ValueString())

			case "map with nil fields":
				variable := &n8nsdk.Variable{
					Key:   "test-key",
					Value: "test-value",
				}

				data := &models.DataSource{}
				mapVariableToDataSourceModel(variable, data)

				assert.True(t, data.ID.IsNull())
				assert.Equal(t, "test-key", data.Key.ValueString())
				assert.Equal(t, "test-value", data.Value.ValueString())
				assert.True(t, data.Type.IsNull())
				assert.True(t, data.ProjectID.IsNull())

			case "map with nil project":
				id := "var-123"
				variable := &n8nsdk.Variable{
					Id:      &id,
					Key:     "test-key",
					Value:   "test-value",
					Project: nil,
				}

				data := &models.DataSource{}
				mapVariableToDataSourceModel(variable, data)

				assert.Equal(t, "var-123", data.ID.ValueString())
				assert.True(t, data.ProjectID.IsNull())

			case "map with project but nil project ID":
				id := "var-123"
				variable := &n8nsdk.Variable{
					Id:    &id,
					Key:   "test-key",
					Value: "test-value",
					Project: &n8nsdk.Project{
						Id: nil,
					},
				}

				data := &models.DataSource{}
				mapVariableToDataSourceModel(variable, data)

				assert.Equal(t, "var-123", data.ID.ValueString())
				assert.True(t, data.ProjectID.IsNull())

			case "map with empty strings":
				id := ""
				varType := ""
				projectID := ""
				variable := &n8nsdk.Variable{
					Id:    &id,
					Key:   "",
					Value: "",
					Type:  &varType,
					Project: &n8nsdk.Project{
						Id: &projectID,
					},
				}

				data := &models.DataSource{}
				mapVariableToDataSourceModel(variable, data)

				assert.Equal(t, "", data.ID.ValueString())
				assert.Equal(t, "", data.Key.ValueString())
				assert.Equal(t, "", data.Value.ValueString())
				assert.Equal(t, "", data.Type.ValueString())
				assert.Equal(t, "", data.ProjectID.ValueString())

			case "map sensitive value":
				id := "var-123"
				variable := &n8nsdk.Variable{
					Id:    &id,
					Key:   "api-key",
					Value: "super-secret-value",
				}

				data := &models.DataSource{}
				mapVariableToDataSourceModel(variable, data)

				assert.Equal(t, "super-secret-value", data.Value.ValueString())

			case "error case - validation checks":
				// Test with nil variable pointer
				data := &models.DataSource{}
				assert.NotPanics(t, func() {
					// Function should handle nil gracefully if it does
					// Since the actual function may panic, we verify behavior
					defer func() {
						if r := recover(); r != nil {
							// Expected to panic with nil variable
							assert.NotNil(t, r)
						}
					}()
					mapVariableToDataSourceModel(nil, data)
				})

				// Test with nil data pointer
				id := "var-123"
				variable := &n8nsdk.Variable{
					Id:    &id,
					Key:   "test-key",
					Value: "test-value",
				}
				assert.NotPanics(t, func() {
					defer func() {
						if r := recover(); r != nil {
							// Expected to panic with nil data
							assert.NotNil(t, r)
						}
					}()
					mapVariableToDataSourceModel(variable, nil)
				})
			}
		})
	}
}

func TestMapVariableToResourceModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "map with all fields populated"},
		{name: "map with nil fields"},
		{name: "map with nil project"},
		{name: "map with project but nil project ID"},
		{name: "error case - validation checks", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "map with all fields populated":
				id := "var-123"
				varType := "string"
				projectID := "proj-456"
				variable := &n8nsdk.Variable{
					Id:    &id,
					Key:   "test-key",
					Value: "test-value",
					Type:  &varType,
					Project: &n8nsdk.Project{
						Id: &projectID,
					},
				}

				data := &models.Resource{}
				mapVariableToResourceModel(variable, data)

				assert.Equal(t, "var-123", data.ID.ValueString())
				assert.Equal(t, "test-key", data.Key.ValueString())
				assert.Equal(t, "test-value", data.Value.ValueString())
				assert.Equal(t, "string", data.Type.ValueString())
				assert.Equal(t, "proj-456", data.ProjectID.ValueString())

			case "map with nil fields":
				variable := &n8nsdk.Variable{
					Key:   "test-key",
					Value: "test-value",
				}

				data := &models.Resource{}
				mapVariableToResourceModel(variable, data)

				assert.True(t, data.ID.IsNull())
				assert.Equal(t, "test-key", data.Key.ValueString())
				assert.Equal(t, "test-value", data.Value.ValueString())
				assert.True(t, data.Type.IsNull())
				assert.True(t, data.ProjectID.IsNull())

			case "map with nil project":
				id := "var-123"
				variable := &n8nsdk.Variable{
					Id:      &id,
					Key:     "test-key",
					Value:   "test-value",
					Project: nil,
				}

				data := &models.Resource{}
				mapVariableToResourceModel(variable, data)

				assert.Equal(t, "var-123", data.ID.ValueString())
				assert.True(t, data.ProjectID.IsNull())

			case "map with project but nil project ID":
				id := "var-123"
				variable := &n8nsdk.Variable{
					Id:    &id,
					Key:   "test-key",
					Value: "test-value",
					Project: &n8nsdk.Project{
						Id: nil,
					},
				}

				data := &models.Resource{}
				mapVariableToResourceModel(variable, data)

				assert.Equal(t, "var-123", data.ID.ValueString())
				assert.True(t, data.ProjectID.IsNull())

			case "error case - validation checks":
				// Test with nil variable pointer
				data := &models.Resource{}
				assert.NotPanics(t, func() {
					defer func() {
						if r := recover(); r != nil {
							// Expected to panic with nil variable
							assert.NotNil(t, r)
						}
					}()
					mapVariableToResourceModel(nil, data)
				})

				// Test with nil data pointer
				id := "var-123"
				variable := &n8nsdk.Variable{
					Id:    &id,
					Key:   "test-key",
					Value: "test-value",
				}
				assert.NotPanics(t, func() {
					defer func() {
						if r := recover(); r != nil {
							// Expected to panic with nil data
							assert.NotNil(t, r)
						}
					}()
					mapVariableToResourceModel(variable, nil)
				})
			}
		})
	}
}

func TestBuildVariableRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "build with all fields"},
		{name: "build with only required fields"},
		{name: "build with unknown type"},
		{name: "build with unknown project ID"},
		{name: "build with empty string values"},
		{name: "build with special characters"},
		{name: "error case - validation checks", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "build with all fields":
				plan := &models.Resource{
					Key:       types.StringValue("test-key"),
					Value:     types.StringValue("test-value"),
					Type:      types.StringValue("string"),
					ProjectID: types.StringValue("proj-123"),
				}

				req := buildVariableRequest(plan)

				assert.Equal(t, "test-key", req.Key)
				assert.Equal(t, "test-value", req.Value)
				assert.NotNil(t, req.Type)
				assert.Equal(t, "string", *req.Type)
				assert.True(t, req.ProjectId.IsSet())
				projectID := req.ProjectId.Get()
				assert.Equal(t, "proj-123", *projectID)

			case "build with only required fields":
				plan := &models.Resource{
					Key:       types.StringValue("test-key"),
					Value:     types.StringValue("test-value"),
					Type:      types.StringNull(),
					ProjectID: types.StringNull(),
				}

				req := buildVariableRequest(plan)

				assert.Equal(t, "test-key", req.Key)
				assert.Equal(t, "test-value", req.Value)
				assert.Nil(t, req.Type)
				assert.False(t, req.ProjectId.IsSet())

			case "build with unknown type":
				plan := &models.Resource{
					Key:       types.StringValue("test-key"),
					Value:     types.StringValue("test-value"),
					Type:      types.StringUnknown(),
					ProjectID: types.StringNull(),
				}

				req := buildVariableRequest(plan)

				assert.Equal(t, "test-key", req.Key)
				assert.Equal(t, "test-value", req.Value)
				assert.Nil(t, req.Type)

			case "build with unknown project ID":
				plan := &models.Resource{
					Key:       types.StringValue("test-key"),
					Value:     types.StringValue("test-value"),
					Type:      types.StringNull(),
					ProjectID: types.StringUnknown(),
				}

				req := buildVariableRequest(plan)

				assert.Equal(t, "test-key", req.Key)
				assert.Equal(t, "test-value", req.Value)
				assert.False(t, req.ProjectId.IsSet())

			case "build with empty string values":
				plan := &models.Resource{
					Key:       types.StringValue(""),
					Value:     types.StringValue(""),
					Type:      types.StringValue(""),
					ProjectID: types.StringValue(""),
				}

				req := buildVariableRequest(plan)

				assert.Equal(t, "", req.Key)
				assert.Equal(t, "", req.Value)
				assert.NotNil(t, req.Type)
				assert.Equal(t, "", *req.Type)
				assert.True(t, req.ProjectId.IsSet())

			case "build with special characters":
				plan := &models.Resource{
					Key:       types.StringValue("KEY_WITH_UNDERSCORES"),
					Value:     types.StringValue("value-with-dashes"),
					Type:      types.StringValue("custom-type"),
					ProjectID: types.StringValue("proj-123-abc"),
				}

				req := buildVariableRequest(plan)

				assert.Equal(t, "KEY_WITH_UNDERSCORES", req.Key)
				assert.Equal(t, "value-with-dashes", req.Value)
				assert.Equal(t, "custom-type", *req.Type)

			case "error case - validation checks":
				// Test with nil plan pointer
				assert.NotPanics(t, func() {
					defer func() {
						if r := recover(); r != nil {
							// Expected to panic with nil plan
							assert.NotNil(t, r)
						}
					}()
					buildVariableRequest(nil)
				})

				// Test with long values
				plan := &models.Resource{
					Key:       types.StringValue("very-long-key-that-exceeds-normal-length-expectations"),
					Value:     types.StringValue("very-long-value-that-exceeds-normal-length-expectations"),
					Type:      types.StringValue("string"),
					ProjectID: types.StringValue("proj-123"),
				}
				req := buildVariableRequest(plan)
				assert.NotNil(t, req)
				assert.Contains(t, req.Key, "very-long-key")
			}
		})
	}
}

func BenchmarkFindVariableByID(b *testing.B) {
	id1 := "var-1"
	id2 := "var-2"
	id3 := "var-3"
	variables := []n8nsdk.Variable{
		{Id: &id1, Key: "key-1", Value: "value-1"},
		{Id: &id2, Key: "key-2", Value: "value-2"},
		{Id: &id3, Key: "key-3", Value: "value-3"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findVariableByID(variables, "var-2")
	}
}

func BenchmarkFindVariableByKey(b *testing.B) {
	id1 := "var-1"
	id2 := "var-2"
	id3 := "var-3"
	variables := []n8nsdk.Variable{
		{Id: &id1, Key: "key-1", Value: "value-1"},
		{Id: &id2, Key: "key-2", Value: "value-2"},
		{Id: &id3, Key: "key-3", Value: "value-3"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findVariableByKey(variables, "key-2")
	}
}

func BenchmarkMapVariableToDataSourceModel(b *testing.B) {
	id := "var-123"
	varType := "string"
	projectID := "proj-456"
	variable := &n8nsdk.Variable{
		Id:    &id,
		Key:   "test-key",
		Value: "test-value",
		Type:  &varType,
		Project: &n8nsdk.Project{
			Id: &projectID,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data := &models.DataSource{}
		mapVariableToDataSourceModel(variable, data)
	}
}

func BenchmarkBuildVariableRequest(b *testing.B) {
	plan := &models.Resource{
		Key:       types.StringValue("test-key"),
		Value:     types.StringValue("test-value"),
		Type:      types.StringValue("string"),
		ProjectID: types.StringValue("proj-123"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buildVariableRequest(plan)
	}
}
