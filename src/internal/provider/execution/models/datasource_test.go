package models

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestDataSource(t *testing.T) {
	t.Run("create with all fields", func(t *testing.T) {
		now := time.Now().Format(time.RFC3339)
		datasource := DataSource{
			ID:          types.StringValue("exec-789"),
			WorkflowID:  types.StringValue("wf-123"),
			Finished:    types.BoolValue(true),
			Mode:        types.StringValue("manual"),
			StartedAt:   types.StringValue(now),
			StoppedAt:   types.StringValue(now),
			Status:      types.StringValue("success"),
			IncludeData: types.BoolValue(true),
		}

		assert.Equal(t, "exec-789", datasource.ID.ValueString())
		assert.Equal(t, "wf-123", datasource.WorkflowID.ValueString())
		assert.True(t, datasource.Finished.ValueBool())
		assert.Equal(t, "manual", datasource.Mode.ValueString())
		assert.Equal(t, now, datasource.StartedAt.ValueString())
		assert.Equal(t, now, datasource.StoppedAt.ValueString())
		assert.Equal(t, "success", datasource.Status.ValueString())
		assert.True(t, datasource.IncludeData.ValueBool())
	})

	t.Run("create with null values", func(t *testing.T) {
		datasource := DataSource{
			ID:          types.StringNull(),
			WorkflowID:  types.StringNull(),
			Finished:    types.BoolNull(),
			Mode:        types.StringNull(),
			StartedAt:   types.StringNull(),
			StoppedAt:   types.StringNull(),
			Status:      types.StringNull(),
			IncludeData: types.BoolNull(),
		}

		assert.True(t, datasource.ID.IsNull())
		assert.True(t, datasource.WorkflowID.IsNull())
		assert.True(t, datasource.Finished.IsNull())
		assert.True(t, datasource.Mode.IsNull())
		assert.True(t, datasource.StartedAt.IsNull())
		assert.True(t, datasource.StoppedAt.IsNull())
		assert.True(t, datasource.Status.IsNull())
		assert.True(t, datasource.IncludeData.IsNull())
	})

	t.Run("create with unknown values", func(t *testing.T) {
		datasource := DataSource{
			ID:          types.StringUnknown(),
			WorkflowID:  types.StringUnknown(),
			Finished:    types.BoolUnknown(),
			Mode:        types.StringUnknown(),
			StartedAt:   types.StringUnknown(),
			StoppedAt:   types.StringUnknown(),
			Status:      types.StringUnknown(),
			IncludeData: types.BoolUnknown(),
		}

		assert.True(t, datasource.ID.IsUnknown())
		assert.True(t, datasource.WorkflowID.IsUnknown())
		assert.True(t, datasource.Finished.IsUnknown())
		assert.True(t, datasource.Mode.IsUnknown())
		assert.True(t, datasource.StartedAt.IsUnknown())
		assert.True(t, datasource.StoppedAt.IsUnknown())
		assert.True(t, datasource.Status.IsUnknown())
		assert.True(t, datasource.IncludeData.IsUnknown())
	})

	t.Run("include data flag", func(t *testing.T) {
		withData := DataSource{
			ID:          types.StringValue("exec-with-data"),
			IncludeData: types.BoolValue(true),
		}
		assert.True(t, withData.IncludeData.ValueBool())
		assert.Equal(t, "exec-with-data", withData.ID.ValueString())

		withoutData := DataSource{
			ID:          types.StringValue("exec-without-data"),
			IncludeData: types.BoolValue(false),
		}
		assert.False(t, withoutData.IncludeData.ValueBool())
		assert.Equal(t, "exec-without-data", withoutData.ID.ValueString())
	})

	t.Run("execution modes", func(t *testing.T) {
		modes := []string{"manual", "trigger", "schedule", "webhook", "api", "error", "retry", "integrated"}

		for _, mode := range modes {
			datasource := DataSource{
				Mode: types.StringValue(mode),
			}
			assert.Equal(t, mode, datasource.Mode.ValueString())
		}
	})

	t.Run("status values", func(t *testing.T) {
		statuses := []string{"success", "error", "waiting", "running", "canceled", "crashed", "new"}

		for _, status := range statuses {
			datasource := DataSource{
				Status: types.StringValue(status),
			}
			assert.Equal(t, status, datasource.Status.ValueString())
		}
	})

	t.Run("finished states", func(t *testing.T) {
		finishedDatasource := DataSource{
			Finished:  types.BoolValue(true),
			Status:    types.StringValue("success"),
			StoppedAt: types.StringValue("2024-01-01T00:00:00Z"),
		}
		assert.True(t, finishedDatasource.Finished.ValueBool())
		assert.Equal(t, "success", finishedDatasource.Status.ValueString())
		assert.NotEmpty(t, finishedDatasource.StoppedAt.ValueString())

		runningDatasource := DataSource{
			Finished:  types.BoolValue(false),
			Status:    types.StringValue("running"),
			StoppedAt: types.StringValue(""),
		}
		assert.False(t, runningDatasource.Finished.ValueBool())
		assert.Equal(t, "running", runningDatasource.Status.ValueString())
		assert.Empty(t, runningDatasource.StoppedAt.ValueString())
	})

	t.Run("zero value struct", func(t *testing.T) {
		var datasource DataSource
		assert.True(t, datasource.ID.IsNull())
		assert.True(t, datasource.WorkflowID.IsNull())
		assert.True(t, datasource.Finished.IsNull())
		assert.True(t, datasource.Mode.IsNull())
		assert.True(t, datasource.StartedAt.IsNull())
		assert.True(t, datasource.StoppedAt.IsNull())
		assert.True(t, datasource.Status.IsNull())
		assert.True(t, datasource.IncludeData.IsNull())
	})

	t.Run("timestamp formats", func(t *testing.T) {
		timestamps := []string{
			"2024-01-01T00:00:00Z",
			"2024-01-01T12:34:56Z",
			"2024-01-01T12:34:56.789Z",
			"2024-01-01T12:34:56+00:00",
			"2024-01-01T12:34:56-05:00",
			time.Now().Format(time.RFC3339),
			time.Now().UTC().Format(time.RFC3339Nano),
		}

		for _, ts := range timestamps {
			datasource := DataSource{
				StartedAt: types.StringValue(ts),
				StoppedAt: types.StringValue(ts),
			}
			assert.Equal(t, ts, datasource.StartedAt.ValueString())
			assert.Equal(t, ts, datasource.StoppedAt.ValueString())
		}
	})

	t.Run("copy struct", func(t *testing.T) {
		original := DataSource{
			ID:          types.StringValue("original-id"),
			WorkflowID:  types.StringValue("original-wf"),
			Finished:    types.BoolValue(true),
			Status:      types.StringValue("success"),
			IncludeData: types.BoolValue(false),
		}

		copied := original

		assert.Equal(t, original.ID.ValueString(), copied.ID.ValueString())
		assert.Equal(t, original.WorkflowID.ValueString(), copied.WorkflowID.ValueString())
		assert.Equal(t, original.Finished.ValueBool(), copied.Finished.ValueBool())
		assert.Equal(t, original.Status.ValueString(), copied.Status.ValueString())
		assert.Equal(t, original.IncludeData.ValueBool(), copied.IncludeData.ValueBool())

		// Modify copied
		copied.ID = types.StringValue("modified-id")
		copied.IncludeData = types.BoolValue(true)
		assert.Equal(t, "original-id", original.ID.ValueString())
		assert.Equal(t, "modified-id", copied.ID.ValueString())
		assert.False(t, original.IncludeData.ValueBool())
		assert.True(t, copied.IncludeData.ValueBool())
	})

	t.Run("partial initialization", func(t *testing.T) {
		datasource := DataSource{
			ID:          types.StringValue("exec-partial"),
			WorkflowID:  types.StringValue("wf-partial"),
			IncludeData: types.BoolValue(true),
			// Other fields remain null
		}

		assert.Equal(t, "exec-partial", datasource.ID.ValueString())
		assert.Equal(t, "wf-partial", datasource.WorkflowID.ValueString())
		assert.True(t, datasource.IncludeData.ValueBool())
		assert.True(t, datasource.Finished.IsNull())
		assert.True(t, datasource.Mode.IsNull())
		assert.True(t, datasource.StartedAt.IsNull())
		assert.True(t, datasource.StoppedAt.IsNull())
		assert.True(t, datasource.Status.IsNull())
	})

	t.Run("execution states with data", func(t *testing.T) {
		// Test running execution with data
		running := DataSource{
			ID:          types.StringValue("exec-running"),
			WorkflowID:  types.StringValue("wf-running"),
			Finished:    types.BoolValue(false),
			Status:      types.StringValue("running"),
			StartedAt:   types.StringValue("2024-01-01T00:00:00Z"),
			StoppedAt:   types.StringValue(""),
			IncludeData: types.BoolValue(true),
		}

		assert.False(t, running.Finished.ValueBool())
		assert.Equal(t, "running", running.Status.ValueString())
		assert.NotEmpty(t, running.StartedAt.ValueString())
		assert.Empty(t, running.StoppedAt.ValueString())
		assert.True(t, running.IncludeData.ValueBool())

		// Test completed execution without data
		completed := DataSource{
			ID:          types.StringValue("exec-completed"),
			WorkflowID:  types.StringValue("wf-completed"),
			Finished:    types.BoolValue(true),
			Status:      types.StringValue("success"),
			StartedAt:   types.StringValue("2024-01-01T00:00:00Z"),
			StoppedAt:   types.StringValue("2024-01-01T00:05:00Z"),
			IncludeData: types.BoolValue(false),
		}

		assert.True(t, completed.Finished.ValueBool())
		assert.Equal(t, "success", completed.Status.ValueString())
		assert.NotEmpty(t, completed.StartedAt.ValueString())
		assert.NotEmpty(t, completed.StoppedAt.ValueString())
		assert.False(t, completed.IncludeData.ValueBool())
	})
}

func TestDataSourceConcurrency(t *testing.T) {
	t.Run("concurrent read", func(t *testing.T) {
		datasource := DataSource{
			ID:          types.StringValue("exec-concurrent"),
			WorkflowID:  types.StringValue("wf-concurrent"),
			Finished:    types.BoolValue(true),
			Mode:        types.StringValue("manual"),
			Status:      types.StringValue("success"),
			IncludeData: types.BoolValue(true),
		}

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				_ = datasource.ID.ValueString()
				_ = datasource.WorkflowID.ValueString()
				_ = datasource.Finished.ValueBool()
				_ = datasource.Mode.ValueString()
				_ = datasource.Status.ValueString()
				_ = datasource.IncludeData.ValueBool()
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func BenchmarkDataSource(b *testing.B) {
	b.Run("create", func(b *testing.B) {
		now := time.Now().Format(time.RFC3339)
		for i := 0; i < b.N; i++ {
			_ = DataSource{
				ID:          types.StringValue("exec-789"),
				WorkflowID:  types.StringValue("wf-123"),
				Finished:    types.BoolValue(true),
				Mode:        types.StringValue("manual"),
				StartedAt:   types.StringValue(now),
				StoppedAt:   types.StringValue(now),
				Status:      types.StringValue("success"),
				IncludeData: types.BoolValue(true),
			}
		}
	})

	b.Run("access fields", func(b *testing.B) {
		datasource := DataSource{
			ID:          types.StringValue("exec-789"),
			WorkflowID:  types.StringValue("wf-123"),
			Finished:    types.BoolValue(true),
			Mode:        types.StringValue("manual"),
			Status:      types.StringValue("success"),
			IncludeData: types.BoolValue(true),
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = datasource.ID.ValueString()
			_ = datasource.WorkflowID.ValueString()
			_ = datasource.Finished.ValueBool()
			_ = datasource.Mode.ValueString()
			_ = datasource.Status.ValueString()
			_ = datasource.IncludeData.ValueBool()
		}
	})
}
