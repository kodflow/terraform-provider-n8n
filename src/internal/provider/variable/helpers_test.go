package variable

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/variable/models"
	"github.com/stretchr/testify/assert"
)

func TestFindVariableByIDOrKey(t *testing.T) {
	t.Run("find by ID when both ID and key provided", func(t *testing.T) {
		id := "var-123"
		variables := []n8nsdk.Variable{
			{Id: &id, Key: "test-key", Value: "test-value"},
		}

		found, exists := findVariableByIDOrKey(variables, types.StringValue("var-123"), types.StringValue("wrong-key"))

		assert.True(t, exists)
		assert.NotNil(t, found)
		assert.Equal(t, "var-123", *found.Id)
	})

	t.Run("find by key when ID is null", func(t *testing.T) {
		id := "var-123"
		variables := []n8nsdk.Variable{
			{Id: &id, Key: "test-key", Value: "test-value"},
		}

		found, exists := findVariableByIDOrKey(variables, types.StringNull(), types.StringValue("test-key"))

		assert.True(t, exists)
		assert.NotNil(t, found)
		assert.Equal(t, "test-key", found.Key)
	})

	t.Run("find by key when both ID and key match", func(t *testing.T) {
		id := "var-123"
		variables := []n8nsdk.Variable{
			{Id: &id, Key: "test-key", Value: "test-value"},
		}

		found, exists := findVariableByIDOrKey(variables, types.StringValue("var-123"), types.StringValue("test-key"))

		assert.True(t, exists)
		assert.NotNil(t, found)
	})

	t.Run("not found when no match", func(t *testing.T) {
		id := "var-123"
		variables := []n8nsdk.Variable{
			{Id: &id, Key: "test-key", Value: "test-value"},
		}

		found, exists := findVariableByIDOrKey(variables, types.StringValue("var-999"), types.StringValue("wrong-key"))

		assert.False(t, exists)
		assert.Nil(t, found)
	})

	t.Run("find in empty list", func(t *testing.T) {
		variables := []n8nsdk.Variable{}

		found, exists := findVariableByIDOrKey(variables, types.StringValue("var-123"), types.StringValue("test-key"))

		assert.False(t, exists)
		assert.Nil(t, found)
	})

	t.Run("find with multiple variables", func(t *testing.T) {
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
	})

	t.Run("find with nil ID in variable", func(t *testing.T) {
		variables := []n8nsdk.Variable{
			{Id: nil, Key: "test-key", Value: "test-value"},
		}

		found, exists := findVariableByIDOrKey(variables, types.StringValue("var-123"), types.StringValue("test-key"))

		assert.True(t, exists)
		assert.NotNil(t, found)
		assert.Equal(t, "test-key", found.Key)
	})

	t.Run("both ID and key are null", func(t *testing.T) {
		id := "var-123"
		variables := []n8nsdk.Variable{
			{Id: &id, Key: "test-key", Value: "test-value"},
		}

		found, exists := findVariableByIDOrKey(variables, types.StringNull(), types.StringNull())

		assert.False(t, exists)
		assert.Nil(t, found)
	})
}

func TestFindVariableByID(t *testing.T) {
	t.Run("find existing variable by ID", func(t *testing.T) {
		id := "var-123"
		variables := []n8nsdk.Variable{
			{Id: &id, Key: "test-key", Value: "test-value"},
		}

		found, exists := findVariableByID(variables, "var-123")

		assert.True(t, exists)
		assert.NotNil(t, found)
		assert.Equal(t, "var-123", *found.Id)
	})

	t.Run("not found when ID does not match", func(t *testing.T) {
		id := "var-123"
		variables := []n8nsdk.Variable{
			{Id: &id, Key: "test-key", Value: "test-value"},
		}

		found, exists := findVariableByID(variables, "var-999")

		assert.False(t, exists)
		assert.Nil(t, found)
	})

	t.Run("find in empty list", func(t *testing.T) {
		variables := []n8nsdk.Variable{}

		found, exists := findVariableByID(variables, "var-123")

		assert.False(t, exists)
		assert.Nil(t, found)
	})

	t.Run("find with multiple variables", func(t *testing.T) {
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
	})

	t.Run("find with nil ID", func(t *testing.T) {
		id := "var-123"
		variables := []n8nsdk.Variable{
			{Id: nil, Key: "key-1", Value: "value-1"},
			{Id: &id, Key: "key-2", Value: "value-2"},
		}

		found, exists := findVariableByID(variables, "var-123")

		assert.True(t, exists)
		assert.NotNil(t, found)
		assert.Equal(t, "var-123", *found.Id)
	})

	t.Run("empty ID string", func(t *testing.T) {
		id := ""
		variables := []n8nsdk.Variable{
			{Id: &id, Key: "test-key", Value: "test-value"},
		}

		found, exists := findVariableByID(variables, "")

		assert.True(t, exists)
		assert.NotNil(t, found)
	})
}

func TestFindVariableByKey(t *testing.T) {
	t.Run("find existing variable by key", func(t *testing.T) {
		id := "var-123"
		variables := []n8nsdk.Variable{
			{Id: &id, Key: "test-key", Value: "test-value"},
		}

		found, exists := findVariableByKey(variables, "test-key")

		assert.True(t, exists)
		assert.NotNil(t, found)
		assert.Equal(t, "test-key", found.Key)
	})

	t.Run("not found when key does not match", func(t *testing.T) {
		id := "var-123"
		variables := []n8nsdk.Variable{
			{Id: &id, Key: "test-key", Value: "test-value"},
		}

		found, exists := findVariableByKey(variables, "wrong-key")

		assert.False(t, exists)
		assert.Nil(t, found)
	})

	t.Run("find in empty list", func(t *testing.T) {
		variables := []n8nsdk.Variable{}

		found, exists := findVariableByKey(variables, "test-key")

		assert.False(t, exists)
		assert.Nil(t, found)
	})

	t.Run("find with multiple variables", func(t *testing.T) {
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
	})

	t.Run("find with empty key", func(t *testing.T) {
		id := "var-123"
		variables := []n8nsdk.Variable{
			{Id: &id, Key: "", Value: "test-value"},
		}

		found, exists := findVariableByKey(variables, "")

		assert.True(t, exists)
		assert.NotNil(t, found)
	})

	t.Run("find first match when duplicates exist", func(t *testing.T) {
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
	})

	t.Run("case sensitive key matching", func(t *testing.T) {
		id := "var-123"
		variables := []n8nsdk.Variable{
			{Id: &id, Key: "TestKey", Value: "test-value"},
		}

		found, exists := findVariableByKey(variables, "testkey")

		assert.False(t, exists)
		assert.Nil(t, found)
	})
}

func TestMapVariableToDataSourceModel(t *testing.T) {
	t.Run("map with all fields populated", func(t *testing.T) {
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
	})

	t.Run("map with nil fields", func(t *testing.T) {
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
	})

	t.Run("map with nil project", func(t *testing.T) {
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
	})

	t.Run("map with project but nil project ID", func(t *testing.T) {
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
	})

	t.Run("map with empty strings", func(t *testing.T) {
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
	})

	t.Run("map sensitive value", func(t *testing.T) {
		id := "var-123"
		variable := &n8nsdk.Variable{
			Id:    &id,
			Key:   "api-key",
			Value: "super-secret-value",
		}

		data := &models.DataSource{}
		mapVariableToDataSourceModel(variable, data)

		assert.Equal(t, "super-secret-value", data.Value.ValueString())
	})
}

func TestMapVariableToResourceModel(t *testing.T) {
	t.Run("map with all fields populated", func(t *testing.T) {
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
	})

	t.Run("map with nil fields", func(t *testing.T) {
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
	})

	t.Run("map with nil project", func(t *testing.T) {
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
	})

	t.Run("map with project but nil project ID", func(t *testing.T) {
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
	})
}

func TestBuildVariableRequest(t *testing.T) {
	t.Run("build with all fields", func(t *testing.T) {
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
	})

	t.Run("build with only required fields", func(t *testing.T) {
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
	})

	t.Run("build with unknown type", func(t *testing.T) {
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
	})

	t.Run("build with unknown project ID", func(t *testing.T) {
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
	})

	t.Run("build with empty string values", func(t *testing.T) {
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
	})

	t.Run("build with special characters", func(t *testing.T) {
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
	})
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
