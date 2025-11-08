package models

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestItem(t *testing.T) {
	t.Run("create with all fields", func(t *testing.T) {
		now := time.Now().Format(time.RFC3339)
		item := Item{
			ID:         types.StringValue("exec-123"),
			WorkflowID: types.StringValue("wf-456"),
			Finished:   types.BoolValue(true),
			Mode:       types.StringValue("manual"),
			StartedAt:  types.StringValue(now),
			StoppedAt:  types.StringValue(now),
			Status:     types.StringValue("success"),
		}

		assert.Equal(t, "exec-123", item.ID.ValueString())
		assert.Equal(t, "wf-456", item.WorkflowID.ValueString())
		assert.True(t, item.Finished.ValueBool())
		assert.Equal(t, "manual", item.Mode.ValueString())
		assert.Equal(t, now, item.StartedAt.ValueString())
		assert.Equal(t, now, item.StoppedAt.ValueString())
		assert.Equal(t, "success", item.Status.ValueString())
	})

	t.Run("create with null values", func(t *testing.T) {
		item := Item{
			ID:         types.StringNull(),
			WorkflowID: types.StringNull(),
			Finished:   types.BoolNull(),
			Mode:       types.StringNull(),
			StartedAt:  types.StringNull(),
			StoppedAt:  types.StringNull(),
			Status:     types.StringNull(),
		}

		assert.True(t, item.ID.IsNull())
		assert.True(t, item.WorkflowID.IsNull())
		assert.True(t, item.Finished.IsNull())
		assert.True(t, item.Mode.IsNull())
		assert.True(t, item.StartedAt.IsNull())
		assert.True(t, item.StoppedAt.IsNull())
		assert.True(t, item.Status.IsNull())
	})

	t.Run("create with unknown values", func(t *testing.T) {
		item := Item{
			ID:         types.StringUnknown(),
			WorkflowID: types.StringUnknown(),
			Finished:   types.BoolUnknown(),
			Mode:       types.StringUnknown(),
			StartedAt:  types.StringUnknown(),
			StoppedAt:  types.StringUnknown(),
			Status:     types.StringUnknown(),
		}

		assert.True(t, item.ID.IsUnknown())
		assert.True(t, item.WorkflowID.IsUnknown())
		assert.True(t, item.Finished.IsUnknown())
		assert.True(t, item.Mode.IsUnknown())
		assert.True(t, item.StartedAt.IsUnknown())
		assert.True(t, item.StoppedAt.IsUnknown())
		assert.True(t, item.Status.IsUnknown())
	})

	t.Run("execution modes", func(t *testing.T) {
		modes := []string{"manual", "trigger", "schedule", "webhook", "api", "error", "retry", "internal"}

		for _, mode := range modes {
			item := Item{
				Mode: types.StringValue(mode),
			}
			assert.Equal(t, mode, item.Mode.ValueString())
		}
	})

	t.Run("status values", func(t *testing.T) {
		statuses := []string{"success", "error", "waiting", "running", "canceled", "crashed", "unknown"}

		for _, status := range statuses {
			item := Item{
				Status: types.StringValue(status),
			}
			assert.Equal(t, status, item.Status.ValueString())
		}
	})

	t.Run("finished states", func(t *testing.T) {
		finishedItem := Item{
			Finished:  types.BoolValue(true),
			StoppedAt: types.StringValue("2024-01-01T00:00:00Z"),
		}
		assert.True(t, finishedItem.Finished.ValueBool())
		assert.NotEmpty(t, finishedItem.StoppedAt.ValueString())

		runningItem := Item{
			Finished:  types.BoolValue(false),
			StoppedAt: types.StringValue(""),
		}
		assert.False(t, runningItem.Finished.ValueBool())
		assert.Empty(t, runningItem.StoppedAt.ValueString())
	})

	t.Run("zero value struct", func(t *testing.T) {
		var item Item
		assert.True(t, item.ID.IsNull())
		assert.True(t, item.WorkflowID.IsNull())
		assert.True(t, item.Finished.IsNull())
		assert.True(t, item.Mode.IsNull())
		assert.True(t, item.StartedAt.IsNull())
		assert.True(t, item.StoppedAt.IsNull())
		assert.True(t, item.Status.IsNull())
	})

	t.Run("timestamp formats", func(t *testing.T) {
		timestamps := []string{
			"2024-01-01T00:00:00Z",
			"2024-01-01T12:34:56Z",
			"2024-01-01T12:34:56.789Z",
			"2024-01-01T12:34:56+00:00",
			"2024-01-01T12:34:56-05:00",
			time.Now().Format(time.RFC3339),
		}

		for _, ts := range timestamps {
			item := Item{
				StartedAt: types.StringValue(ts),
				StoppedAt: types.StringValue(ts),
			}
			assert.Equal(t, ts, item.StartedAt.ValueString())
			assert.Equal(t, ts, item.StoppedAt.ValueString())
		}
	})

	t.Run("copy struct", func(t *testing.T) {
		original := Item{
			ID:         types.StringValue("original-id"),
			WorkflowID: types.StringValue("original-wf"),
			Finished:   types.BoolValue(true),
			Status:     types.StringValue("success"),
		}

		copied := original

		assert.Equal(t, original.ID.ValueString(), copied.ID.ValueString())
		assert.Equal(t, original.WorkflowID.ValueString(), copied.WorkflowID.ValueString())
		assert.Equal(t, original.Finished.ValueBool(), copied.Finished.ValueBool())
		assert.Equal(t, original.Status.ValueString(), copied.Status.ValueString())

		// Modify copied
		copied.ID = types.StringValue("modified-id")
		assert.Equal(t, "original-id", original.ID.ValueString())
		assert.Equal(t, "modified-id", copied.ID.ValueString())
	})

	t.Run("partial initialization", func(t *testing.T) {
		item := Item{
			ID:     types.StringValue("exec-partial"),
			Status: types.StringValue("running"),
			// Other fields remain null
		}

		assert.Equal(t, "exec-partial", item.ID.ValueString())
		assert.Equal(t, "running", item.Status.ValueString())
		assert.True(t, item.WorkflowID.IsNull())
		assert.True(t, item.Finished.IsNull())
		assert.True(t, item.Mode.IsNull())
		assert.True(t, item.StartedAt.IsNull())
		assert.True(t, item.StoppedAt.IsNull())
	})

	t.Run("execution states", func(t *testing.T) {
		// Test running execution
		running := Item{
			ID:        types.StringValue("exec-running"),
			Finished:  types.BoolValue(false),
			Status:    types.StringValue("running"),
			StartedAt: types.StringValue("2024-01-01T00:00:00Z"),
			StoppedAt: types.StringValue(""),
		}

		assert.False(t, running.Finished.ValueBool())
		assert.Equal(t, "running", running.Status.ValueString())
		assert.NotEmpty(t, running.StartedAt.ValueString())
		assert.Empty(t, running.StoppedAt.ValueString())

		// Test completed execution
		completed := Item{
			ID:        types.StringValue("exec-completed"),
			Finished:  types.BoolValue(true),
			Status:    types.StringValue("success"),
			StartedAt: types.StringValue("2024-01-01T00:00:00Z"),
			StoppedAt: types.StringValue("2024-01-01T00:05:00Z"),
		}

		assert.True(t, completed.Finished.ValueBool())
		assert.Equal(t, "success", completed.Status.ValueString())
		assert.NotEmpty(t, completed.StartedAt.ValueString())
		assert.NotEmpty(t, completed.StoppedAt.ValueString())

		// Test failed execution
		failed := Item{
			ID:        types.StringValue("exec-failed"),
			Finished:  types.BoolValue(true),
			Status:    types.StringValue("error"),
			StartedAt: types.StringValue("2024-01-01T00:00:00Z"),
			StoppedAt: types.StringValue("2024-01-01T00:00:30Z"),
		}

		assert.True(t, failed.Finished.ValueBool())
		assert.Equal(t, "error", failed.Status.ValueString())
		assert.NotEmpty(t, failed.StartedAt.ValueString())
		assert.NotEmpty(t, failed.StoppedAt.ValueString())
	})
}

func TestItemConcurrency(t *testing.T) {
	t.Run("concurrent read", func(t *testing.T) {
		item := Item{
			ID:         types.StringValue("exec-concurrent"),
			WorkflowID: types.StringValue("wf-concurrent"),
			Finished:   types.BoolValue(true),
			Mode:       types.StringValue("manual"),
			Status:     types.StringValue("success"),
		}

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				_ = item.ID.ValueString()
				_ = item.WorkflowID.ValueString()
				_ = item.Finished.ValueBool()
				_ = item.Mode.ValueString()
				_ = item.Status.ValueString()
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func BenchmarkItem(b *testing.B) {
	b.Run("create", func(b *testing.B) {
		now := time.Now().Format(time.RFC3339)
		for i := 0; i < b.N; i++ {
			_ = Item{
				ID:         types.StringValue("exec-123"),
				WorkflowID: types.StringValue("wf-456"),
				Finished:   types.BoolValue(true),
				Mode:       types.StringValue("manual"),
				StartedAt:  types.StringValue(now),
				StoppedAt:  types.StringValue(now),
				Status:     types.StringValue("success"),
			}
		}
	})

	b.Run("access fields", func(b *testing.B) {
		item := Item{
			ID:         types.StringValue("exec-123"),
			WorkflowID: types.StringValue("wf-456"),
			Finished:   types.BoolValue(true),
			Mode:       types.StringValue("manual"),
			Status:     types.StringValue("success"),
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = item.ID.ValueString()
			_ = item.WorkflowID.ValueString()
			_ = item.Finished.ValueBool()
			_ = item.Mode.ValueString()
			_ = item.Status.ValueString()
		}
	})
}
