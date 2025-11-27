package variable

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/variable/models"
	"github.com/stretchr/testify/assert"
)

// TestfindVariableByIDOrKey tests the findVariableByIDOrKey private function.
func Test_findVariableByIDOrKey(t *testing.T) {
	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "find variable by ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "key1", Value: "value1"},
				}
				foundVar, found := findVariableByIDOrKey(variables, types.StringValue("var-123"), types.StringNull())
				assert.True(t, found)
				assert.NotNil(t, foundVar)
				assert.Equal(t, "var-123", *foundVar.Id)
			},
		},
		{
			name: "find variable by key",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-456"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "key1", Value: "value1"},
				}
				foundVar, found := findVariableByIDOrKey(variables, types.StringNull(), types.StringValue("key1"))
				assert.True(t, found)
				assert.NotNil(t, foundVar)
				assert.Equal(t, "key1", foundVar.Key)
			},
		},
		{
			name: "find variable by ID when both ID and key provided",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-789"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "key1", Value: "value1"},
				}
				foundVar, found := findVariableByIDOrKey(variables, types.StringValue("var-789"), types.StringValue("key1"))
				assert.True(t, found)
				assert.NotNil(t, foundVar)
				assert.Equal(t, "var-789", *foundVar.Id)
			},
		},
		{
			name: "variable not found",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "key1", Value: "value1"},
				}
				foundVar, found := findVariableByIDOrKey(variables, types.StringValue("non-existent"), types.StringValue("non-existent"))
				assert.False(t, found)
				assert.Nil(t, foundVar)
			},
		},
		{
			name: "error case - empty variables list",
			testFunc: func(t *testing.T) {
				t.Helper()
				variables := []n8nsdk.Variable{}
				foundVar, found := findVariableByIDOrKey(variables, types.StringValue("var-123"), types.StringValue("key1"))
				assert.False(t, found)
				assert.Nil(t, foundVar)
			},
		},
		{
			name: "error case - nil ID and null key",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "key1", Value: "value1"},
				}
				foundVar, found := findVariableByIDOrKey(variables, types.StringNull(), types.StringNull())
				assert.False(t, found)
				assert.Nil(t, foundVar)
			},
		},
		{
			name: "error case - variable with nil ID field",
			testFunc: func(t *testing.T) {
				t.Helper()
				variables := []n8nsdk.Variable{
					{Id: nil, Key: "key1", Value: "value1"},
				}
				foundVar, found := findVariableByIDOrKey(variables, types.StringValue("var-123"), types.StringNull())
				assert.False(t, found)
				assert.Nil(t, foundVar)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestfindVariableByID tests the findVariableByID private function.
func Test_findVariableByID(t *testing.T) {
	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "find variable by ID success",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "key1", Value: "value1"},
				}
				foundVar, found := findVariableByID(variables, "var-123")
				assert.True(t, found)
				assert.NotNil(t, foundVar)
				assert.Equal(t, "var-123", *foundVar.Id)
			},
		},
		{
			name: "multiple variables find correct one",
			testFunc: func(t *testing.T) {
				t.Helper()
				id1 := "var-123"
				id2 := "var-456"
				variables := []n8nsdk.Variable{
					{Id: &id1, Key: "key1", Value: "value1"},
					{Id: &id2, Key: "key2", Value: "value2"},
				}
				foundVar, found := findVariableByID(variables, "var-456")
				assert.True(t, found)
				assert.NotNil(t, foundVar)
				assert.Equal(t, "var-456", *foundVar.Id)
			},
		},
		{
			name: "variable not found",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "key1", Value: "value1"},
				}
				foundVar, found := findVariableByID(variables, "non-existent")
				assert.False(t, found)
				assert.Nil(t, foundVar)
			},
		},
		{
			name: "error case - empty variables list",
			testFunc: func(t *testing.T) {
				t.Helper()
				variables := []n8nsdk.Variable{}
				foundVar, found := findVariableByID(variables, "var-123")
				assert.False(t, found)
				assert.Nil(t, foundVar)
			},
		},
		{
			name: "error case - variable with nil ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				variables := []n8nsdk.Variable{
					{Id: nil, Key: "key1", Value: "value1"},
				}
				foundVar, found := findVariableByID(variables, "var-123")
				assert.False(t, found)
				assert.Nil(t, foundVar)
			},
		},
		{
			name: "error case - empty search ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "key1", Value: "value1"},
				}
				foundVar, found := findVariableByID(variables, "")
				assert.False(t, found)
				assert.Nil(t, foundVar)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestfindVariableByKey tests the findVariableByKey private function.
func Test_findVariableByKey(t *testing.T) {
	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "find variable by key success",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "mykey", Value: "value1"},
				}
				foundVar, found := findVariableByKey(variables, "mykey")
				assert.True(t, found)
				assert.NotNil(t, foundVar)
				assert.Equal(t, "mykey", foundVar.Key)
			},
		},
		{
			name: "multiple variables find correct one",
			testFunc: func(t *testing.T) {
				t.Helper()
				id1 := "var-123"
				id2 := "var-456"
				variables := []n8nsdk.Variable{
					{Id: &id1, Key: "key1", Value: "value1"},
					{Id: &id2, Key: "key2", Value: "value2"},
				}
				foundVar, found := findVariableByKey(variables, "key2")
				assert.True(t, found)
				assert.NotNil(t, foundVar)
				assert.Equal(t, "key2", foundVar.Key)
			},
		},
		{
			name: "variable not found",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "key1", Value: "value1"},
				}
				foundVar, found := findVariableByKey(variables, "non-existent")
				assert.False(t, found)
				assert.Nil(t, foundVar)
			},
		},
		{
			name: "error case - empty variables list",
			testFunc: func(t *testing.T) {
				t.Helper()
				variables := []n8nsdk.Variable{}
				foundVar, found := findVariableByKey(variables, "key1")
				assert.False(t, found)
				assert.Nil(t, foundVar)
			},
		},
		{
			name: "error case - empty search key",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "key1", Value: "value1"},
				}
				foundVar, found := findVariableByKey(variables, "")
				assert.False(t, found)
				assert.Nil(t, foundVar)
			},
		},
		{
			name: "error case - case sensitive key search",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variables := []n8nsdk.Variable{
					{Id: &id, Key: "MyKey", Value: "value1"},
				}
				foundVar, found := findVariableByKey(variables, "mykey")
				assert.False(t, found)
				assert.Nil(t, foundVar)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestmapVariableToDataSourceModel tests the mapVariableToDataSourceModel private function.
func Test_mapVariableToDataSourceModel(t *testing.T) {
	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "map all fields success",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				varType := "string"
				projectID := "proj-456"
				variable := &n8nsdk.Variable{
					Id:    &id,
					Key:   "mykey",
					Value: "myvalue",
					Type:  &varType,
					Project: &n8nsdk.Project{
						Id: &projectID,
					},
				}
				data := &models.DataSource{}
				mapVariableToDataSourceModel(variable, data)
				assert.Equal(t, "var-123", data.ID.ValueString())
				assert.Equal(t, "mykey", data.Key.ValueString())
				assert.Equal(t, "myvalue", data.Value.ValueString())
				assert.Equal(t, "string", data.Type.ValueString())
				assert.Equal(t, "proj-456", data.ProjectID.ValueString())
			},
		},
		{
			name: "map with nil ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				variable := &n8nsdk.Variable{
					Id:    nil,
					Key:   "mykey",
					Value: "myvalue",
				}
				data := &models.DataSource{}
				mapVariableToDataSourceModel(variable, data)
				assert.True(t, data.ID.IsNull())
				assert.Equal(t, "mykey", data.Key.ValueString())
				assert.Equal(t, "myvalue", data.Value.ValueString())
			},
		},
		{
			name: "map with nil Type",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variable := &n8nsdk.Variable{
					Id:    &id,
					Key:   "mykey",
					Value: "myvalue",
					Type:  nil,
				}
				data := &models.DataSource{}
				mapVariableToDataSourceModel(variable, data)
				assert.True(t, data.Type.IsNull())
			},
		},
		{
			name: "map with nil Project",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variable := &n8nsdk.Variable{
					Id:      &id,
					Key:     "mykey",
					Value:   "myvalue",
					Project: nil,
				}
				data := &models.DataSource{}
				mapVariableToDataSourceModel(variable, data)
				assert.True(t, data.ProjectID.IsNull())
			},
		},
		{
			name: "error case - nil variable pointer causes panic",
			testFunc: func(t *testing.T) {
				t.Helper()
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic from nil variable, but did not panic")
					}
				}()
				data := &models.DataSource{}
				mapVariableToDataSourceModel(nil, data)
			},
		},
		{
			name: "error case - nil data pointer causes panic",
			testFunc: func(t *testing.T) {
				t.Helper()
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic from nil data, but did not panic")
					}
				}()
				id := "var-123"
				variable := &n8nsdk.Variable{
					Id:    &id,
					Key:   "mykey",
					Value: "myvalue",
				}
				mapVariableToDataSourceModel(variable, nil)
			},
		},
		{
			name: "error case - project with nil ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variable := &n8nsdk.Variable{
					Id:      &id,
					Key:     "mykey",
					Value:   "myvalue",
					Project: &n8nsdk.Project{Id: nil},
				}
				data := &models.DataSource{}
				mapVariableToDataSourceModel(variable, data)
				assert.True(t, data.ProjectID.IsNull())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestmapVariableToResourceModel tests the mapVariableToResourceModel private function.
func Test_mapVariableToResourceModel(t *testing.T) {
	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "map all fields success",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				varType := "string"
				projectID := "proj-456"
				variable := &n8nsdk.Variable{
					Id:    &id,
					Key:   "mykey",
					Value: "myvalue",
					Type:  &varType,
					Project: &n8nsdk.Project{
						Id: &projectID,
					},
				}
				data := &models.Resource{}
				mapVariableToResourceModel(variable, data)
				assert.Equal(t, "var-123", data.ID.ValueString())
				assert.Equal(t, "mykey", data.Key.ValueString())
				assert.Equal(t, "myvalue", data.Value.ValueString())
				assert.Equal(t, "string", data.Type.ValueString())
				assert.Equal(t, "proj-456", data.ProjectID.ValueString())
			},
		},
		{
			name: "map with nil ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				variable := &n8nsdk.Variable{
					Id:    nil,
					Key:   "mykey",
					Value: "myvalue",
				}
				data := &models.Resource{}
				mapVariableToResourceModel(variable, data)
				assert.True(t, data.ID.IsNull())
				assert.Equal(t, "mykey", data.Key.ValueString())
				assert.Equal(t, "myvalue", data.Value.ValueString())
			},
		},
		{
			name: "map with nil Type",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variable := &n8nsdk.Variable{
					Id:    &id,
					Key:   "mykey",
					Value: "myvalue",
					Type:  nil,
				}
				data := &models.Resource{}
				mapVariableToResourceModel(variable, data)
				assert.True(t, data.Type.IsNull())
			},
		},
		{
			name: "map with nil Project",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variable := &n8nsdk.Variable{
					Id:      &id,
					Key:     "mykey",
					Value:   "myvalue",
					Project: nil,
				}
				data := &models.Resource{}
				mapVariableToResourceModel(variable, data)
				assert.True(t, data.ProjectID.IsNull())
			},
		},
		{
			name: "error case - nil variable pointer causes panic",
			testFunc: func(t *testing.T) {
				t.Helper()
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic from nil variable, but did not panic")
					}
				}()
				data := &models.Resource{}
				mapVariableToResourceModel(nil, data)
			},
		},
		{
			name: "error case - nil data pointer causes panic",
			testFunc: func(t *testing.T) {
				t.Helper()
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic from nil data, but did not panic")
					}
				}()
				id := "var-123"
				variable := &n8nsdk.Variable{
					Id:    &id,
					Key:   "mykey",
					Value: "myvalue",
				}
				mapVariableToResourceModel(variable, nil)
			},
		},
		{
			name: "error case - project with nil ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := "var-123"
				variable := &n8nsdk.Variable{
					Id:      &id,
					Key:     "mykey",
					Value:   "myvalue",
					Project: &n8nsdk.Project{Id: nil},
				}
				data := &models.Resource{}
				mapVariableToResourceModel(variable, data)
				assert.True(t, data.ProjectID.IsNull())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestbuildVariableRequest tests the buildVariableRequest private function.
func Test_buildVariableRequest(t *testing.T) {
	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "build request with all fields",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					Key:       types.StringValue("mykey"),
					Value:     types.StringValue("myvalue"),
					Type:      types.StringValue("string"), // Type is computed/read-only, not sent to API
					ProjectID: types.StringValue("proj-123"),
				}
				request := buildVariableRequest(plan)
				assert.Equal(t, "mykey", request.Key)
				assert.Equal(t, "myvalue", request.Value)
				// Type is NOT sent to API - it's computed by n8n
				assert.Nil(t, request.Type)
				assert.True(t, request.ProjectId.IsSet())
				assert.Equal(t, "proj-123", *request.ProjectId.Get())
			},
		},
		{
			name: "build request with null Type",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					Key:   types.StringValue("mykey"),
					Value: types.StringValue("myvalue"),
					Type:  types.StringNull(),
				}
				request := buildVariableRequest(plan)
				assert.Equal(t, "mykey", request.Key)
				assert.Equal(t, "myvalue", request.Value)
				assert.Nil(t, request.Type)
			},
		},
		{
			name: "build request with unknown Type",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					Key:   types.StringValue("mykey"),
					Value: types.StringValue("myvalue"),
					Type:  types.StringUnknown(),
				}
				request := buildVariableRequest(plan)
				assert.Equal(t, "mykey", request.Key)
				assert.Equal(t, "myvalue", request.Value)
				assert.Nil(t, request.Type)
			},
		},
		{
			name: "build request with null ProjectID",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					Key:       types.StringValue("mykey"),
					Value:     types.StringValue("myvalue"),
					ProjectID: types.StringNull(),
				}
				request := buildVariableRequest(plan)
				assert.Equal(t, "mykey", request.Key)
				assert.Equal(t, "myvalue", request.Value)
				assert.False(t, request.ProjectId.IsSet())
			},
		},
		{
			name: "build request with unknown ProjectID",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					Key:       types.StringValue("mykey"),
					Value:     types.StringValue("myvalue"),
					ProjectID: types.StringUnknown(),
				}
				request := buildVariableRequest(plan)
				assert.Equal(t, "mykey", request.Key)
				assert.Equal(t, "myvalue", request.Value)
				assert.False(t, request.ProjectId.IsSet())
			},
		},
		{
			name: "error case - nil plan pointer causes panic",
			testFunc: func(t *testing.T) {
				t.Helper()
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic from nil plan, but did not panic")
					}
				}()
				buildVariableRequest(nil)
			},
		},
		{
			name: "error case - empty key and value",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					Key:   types.StringValue(""),
					Value: types.StringValue(""),
				}
				request := buildVariableRequest(plan)
				assert.Equal(t, "", request.Key)
				assert.Equal(t, "", request.Value)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}
