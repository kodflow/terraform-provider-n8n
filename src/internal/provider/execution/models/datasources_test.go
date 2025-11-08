package models

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestDataSources(t *testing.T) {
	t.Run("create with executions list", func(t *testing.T) {
		datasources := DataSources{
			WorkflowID:  types.StringValue("wf-main"),
			ProjectID:   types.StringValue("proj-123"),
			Status:      types.StringValue("success"),
			IncludeData: types.BoolValue(true),
			Executions: []Item{
				{
					ID:         types.StringValue("exec-1"),
					WorkflowID: types.StringValue("wf-1"),
					Finished:   types.BoolValue(true),
					Mode:       types.StringValue("manual"),
					StartedAt:  types.StringValue("2024-01-01T00:00:00Z"),
					StoppedAt:  types.StringValue("2024-01-01T00:01:00Z"),
					Status:     types.StringValue("success"),
				},
				{
					ID:         types.StringValue("exec-2"),
					WorkflowID: types.StringValue("wf-2"),
					Finished:   types.BoolValue(false),
					Mode:       types.StringValue("trigger"),
					StartedAt:  types.StringValue("2024-01-01T00:02:00Z"),
					StoppedAt:  types.StringValue(""),
					Status:     types.StringValue("running"),
				},
			},
		}

		assert.Equal(t, "wf-main", datasources.WorkflowID.ValueString())
		assert.Equal(t, "proj-123", datasources.ProjectID.ValueString())
		assert.Equal(t, "success", datasources.Status.ValueString())
		assert.True(t, datasources.IncludeData.ValueBool())
		assert.Len(t, datasources.Executions, 2)
		assert.Equal(t, "exec-1", datasources.Executions[0].ID.ValueString())
		assert.Equal(t, "exec-2", datasources.Executions[1].ID.ValueString())
	})

	t.Run("create with empty executions", func(t *testing.T) {
		datasources := DataSources{
			WorkflowID:  types.StringValue("wf-empty"),
			ProjectID:   types.StringValue("proj-empty"),
			Status:      types.StringValue(""),
			IncludeData: types.BoolValue(false),
			Executions:  []Item{},
		}

		assert.Equal(t, "wf-empty", datasources.WorkflowID.ValueString())
		assert.Equal(t, "proj-empty", datasources.ProjectID.ValueString())
		assert.Equal(t, "", datasources.Status.ValueString())
		assert.False(t, datasources.IncludeData.ValueBool())
		assert.Len(t, datasources.Executions, 0)
	})

	t.Run("create with nil executions", func(t *testing.T) {
		datasources := DataSources{
			WorkflowID:  types.StringValue("wf-nil"),
			ProjectID:   types.StringValue("proj-nil"),
			Status:      types.StringValue("all"),
			IncludeData: types.BoolValue(false),
			Executions:  nil,
		}

		assert.Equal(t, "wf-nil", datasources.WorkflowID.ValueString())
		assert.Nil(t, datasources.Executions)
	})

	t.Run("create with null values", func(t *testing.T) {
		datasources := DataSources{
			WorkflowID:  types.StringNull(),
			ProjectID:   types.StringNull(),
			Status:      types.StringNull(),
			IncludeData: types.BoolNull(),
			Executions:  nil,
		}

		assert.True(t, datasources.WorkflowID.IsNull())
		assert.True(t, datasources.ProjectID.IsNull())
		assert.True(t, datasources.Status.IsNull())
		assert.True(t, datasources.IncludeData.IsNull())
		assert.Nil(t, datasources.Executions)
	})

	t.Run("create with unknown values", func(t *testing.T) {
		datasources := DataSources{
			WorkflowID:  types.StringUnknown(),
			ProjectID:   types.StringUnknown(),
			Status:      types.StringUnknown(),
			IncludeData: types.BoolUnknown(),
			Executions:  []Item{},
		}

		assert.True(t, datasources.WorkflowID.IsUnknown())
		assert.True(t, datasources.ProjectID.IsUnknown())
		assert.True(t, datasources.Status.IsUnknown())
		assert.True(t, datasources.IncludeData.IsUnknown())
		assert.NotNil(t, datasources.Executions)
		assert.Len(t, datasources.Executions, 0)
	})

	t.Run("filter by status", func(t *testing.T) {
		statuses := []string{"success", "error", "waiting", "running", "canceled"}

		for _, status := range statuses {
			datasources := DataSources{
				Status: types.StringValue(status),
			}
			assert.Equal(t, status, datasources.Status.ValueString())
		}
	})

	t.Run("include data flag", func(t *testing.T) {
		withData := DataSources{
			IncludeData: types.BoolValue(true),
		}
		assert.True(t, withData.IncludeData.ValueBool())

		withoutData := DataSources{
			IncludeData: types.BoolValue(false),
		}
		assert.False(t, withoutData.IncludeData.ValueBool())
	})

	t.Run("copy struct", func(t *testing.T) {
		original := DataSources{
			WorkflowID: types.StringValue("original-wf"),
			ProjectID:  types.StringValue("original-proj"),
			Executions: []Item{
				{
					ID: types.StringValue("exec-original"),
				},
			},
		}

		copied := original

		assert.Equal(t, original.WorkflowID.ValueString(), copied.WorkflowID.ValueString())
		assert.Equal(t, original.ProjectID.ValueString(), copied.ProjectID.ValueString())

		// Modify copied
		copied.WorkflowID = types.StringValue("modified-wf")
		assert.Equal(t, "original-wf", original.WorkflowID.ValueString())
		assert.Equal(t, "modified-wf", copied.WorkflowID.ValueString())

		// Note: slice is shared between original and copied
		copied.Executions[0].ID = types.StringValue("exec-modified")
		assert.Equal(t, "exec-modified", original.Executions[0].ID.ValueString())
	})

	t.Run("zero value struct", func(t *testing.T) {
		var datasources DataSources
		assert.True(t, datasources.WorkflowID.IsNull())
		assert.True(t, datasources.ProjectID.IsNull())
		assert.True(t, datasources.Status.IsNull())
		assert.True(t, datasources.IncludeData.IsNull())
		assert.Nil(t, datasources.Executions)
	})

	t.Run("multiple executions with different states", func(t *testing.T) {
		datasources := DataSources{
			WorkflowID: types.StringValue("wf-multi"),
			Executions: []Item{
				{
					ID:       types.StringValue("exec-success"),
					Finished: types.BoolValue(true),
					Status:   types.StringValue("success"),
				},
				{
					ID:       types.StringValue("exec-error"),
					Finished: types.BoolValue(true),
					Status:   types.StringValue("error"),
				},
				{
					ID:       types.StringValue("exec-running"),
					Finished: types.BoolValue(false),
					Status:   types.StringValue("running"),
				},
			},
		}

		assert.Len(t, datasources.Executions, 3)
		assert.True(t, datasources.Executions[0].Finished.ValueBool())
		assert.True(t, datasources.Executions[1].Finished.ValueBool())
		assert.False(t, datasources.Executions[2].Finished.ValueBool())
	})
}

func TestDataSourcesConcurrency(t *testing.T) {
	t.Run("concurrent read", func(t *testing.T) {
		datasources := DataSources{
			WorkflowID:  types.StringValue("wf-concurrent"),
			ProjectID:   types.StringValue("proj-concurrent"),
			Status:      types.StringValue("running"),
			IncludeData: types.BoolValue(true),
		}

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				_ = datasources.WorkflowID.ValueString()
				_ = datasources.ProjectID.ValueString()
				_ = datasources.Status.ValueString()
				_ = datasources.IncludeData.ValueBool()
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func BenchmarkDataSources(b *testing.B) {
	b.Run("create", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = DataSources{
				WorkflowID:  types.StringValue("wf-123"),
				ProjectID:   types.StringValue("proj-123"),
				Status:      types.StringValue("success"),
				IncludeData: types.BoolValue(true),
				Executions:  []Item{},
			}
		}
	})

	b.Run("access fields", func(b *testing.B) {
		datasources := DataSources{
			WorkflowID:  types.StringValue("wf-123"),
			ProjectID:   types.StringValue("proj-123"),
			Status:      types.StringValue("success"),
			IncludeData: types.BoolValue(true),
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = datasources.WorkflowID.ValueString()
			_ = datasources.ProjectID.ValueString()
			_ = datasources.Status.ValueString()
			_ = datasources.IncludeData.ValueBool()
		}
	})
}
