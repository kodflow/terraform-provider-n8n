package execution

import (
	"testing"
	"time"

	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/stretchr/testify/assert"
)

func TestMapExecutionToItem(t *testing.T) {
	t.Run("map with all fields populated", func(t *testing.T) {
		id := float32(123)
		workflowID := float32(456)
		finished := true
		mode := "manual"
		startedAt := time.Now()
		stoppedAt := time.Now().Add(5 * time.Minute)
		status := "success"

		stoppedAtNullable := n8nsdk.NewNullableTime(&stoppedAt)

		execution := &n8nsdk.Execution{
			Id:         &id,
			WorkflowId: &workflowID,
			Finished:   &finished,
			Mode:       &mode,
			StartedAt:  &startedAt,
			StoppedAt:  *stoppedAtNullable,
			Status:     &status,
		}

		item := mapExecutionToItem(execution)

		assert.Equal(t, "123", item.ID.ValueString())
		assert.Equal(t, "456", item.WorkflowID.ValueString())
		assert.True(t, item.Finished.ValueBool())
		assert.Equal(t, "manual", item.Mode.ValueString())
		assert.Equal(t, startedAt.String(), item.StartedAt.ValueString())
		assert.Equal(t, stoppedAt.String(), item.StoppedAt.ValueString())
		assert.Equal(t, "success", item.Status.ValueString())
	})

	t.Run("map with nil fields", func(t *testing.T) {
		execution := &n8nsdk.Execution{}

		item := mapExecutionToItem(execution)

		assert.True(t, item.ID.IsNull())
		assert.True(t, item.WorkflowID.IsNull())
		assert.True(t, item.Finished.IsNull())
		assert.True(t, item.Mode.IsNull())
		assert.True(t, item.StartedAt.IsNull())
		assert.True(t, item.StoppedAt.IsNull())
		assert.True(t, item.Status.IsNull())
	})

	t.Run("map with partial fields", func(t *testing.T) {
		id := float32(789)
		workflowID := float32(999)
		status := "running"

		execution := &n8nsdk.Execution{
			Id:         &id,
			WorkflowId: &workflowID,
			Status:     &status,
			// Other fields nil
		}

		item := mapExecutionToItem(execution)

		assert.Equal(t, "789", item.ID.ValueString())
		assert.Equal(t, "999", item.WorkflowID.ValueString())
		assert.Equal(t, "running", item.Status.ValueString())
		assert.True(t, item.Finished.IsNull())
		assert.True(t, item.Mode.IsNull())
		assert.True(t, item.StartedAt.IsNull())
		assert.True(t, item.StoppedAt.IsNull())
	})

	t.Run("map with unset nullable stoppedAt", func(t *testing.T) {
		id := float32(999)
		workflowID := float32(888)
		finished := false
		status := "running"

		execution := &n8nsdk.Execution{
			Id:         &id,
			WorkflowId: &workflowID,
			Finished:   &finished,
			Status:     &status,
			StoppedAt:  n8nsdk.NullableTime{}, // Unset
		}

		item := mapExecutionToItem(execution)

		assert.Equal(t, "999", item.ID.ValueString())
		assert.Equal(t, "888", item.WorkflowID.ValueString())
		assert.False(t, item.Finished.ValueBool())
		assert.Equal(t, "running", item.Status.ValueString())
		assert.True(t, item.StoppedAt.IsNull())
	})

	t.Run("map different execution modes", func(t *testing.T) {
		modes := []string{"manual", "trigger", "schedule", "webhook", "api", "error", "retry"}

		for _, mode := range modes {
			modeCopy := mode
			execution := &n8nsdk.Execution{
				Mode: &modeCopy,
			}

			item := mapExecutionToItem(execution)

			assert.Equal(t, mode, item.Mode.ValueString())
		}
	})

	t.Run("map different status values", func(t *testing.T) {
		statuses := []string{"success", "error", "waiting", "running", "canceled", "crashed", "new"}

		for _, status := range statuses {
			statusCopy := status
			execution := &n8nsdk.Execution{
				Status: &statusCopy,
			}

			item := mapExecutionToItem(execution)

			assert.Equal(t, status, item.Status.ValueString())
		}
	})

	t.Run("map with zero value id", func(t *testing.T) {
		id := float32(0)
		execution := &n8nsdk.Execution{
			Id: &id,
		}

		item := mapExecutionToItem(execution)

		assert.Equal(t, "0", item.ID.ValueString())
	})

	t.Run("map with negative id", func(t *testing.T) {
		id := float32(-1)
		execution := &n8nsdk.Execution{
			Id: &id,
		}

		item := mapExecutionToItem(execution)

		assert.Equal(t, "-1", item.ID.ValueString())
	})

	t.Run("map with large id", func(t *testing.T) {
		id := float32(2147483647) // Max int32
		execution := &n8nsdk.Execution{
			Id: &id,
		}

		item := mapExecutionToItem(execution)

		// Note: float32 may lose precision for very large integers
		assert.Contains(t, item.ID.ValueString(), "2.1474")
	})

	t.Run("map finished states", func(t *testing.T) {
		// Test finished execution
		finished := true
		finishedExec := &n8nsdk.Execution{
			Finished: &finished,
		}

		finishedItem := mapExecutionToItem(finishedExec)
		assert.True(t, finishedItem.Finished.ValueBool())

		// Test running execution
		running := false
		runningExec := &n8nsdk.Execution{
			Finished: &running,
		}

		runningItem := mapExecutionToItem(runningExec)
		assert.False(t, runningItem.Finished.ValueBool())
	})

	t.Run("map with empty string values", func(t *testing.T) {
		mode := ""
		status := ""

		execution := &n8nsdk.Execution{
			Mode:   &mode,
			Status: &status,
		}

		item := mapExecutionToItem(execution)

		assert.Equal(t, "", item.Mode.ValueString())
		assert.Equal(t, "", item.Status.ValueString())
	})

	t.Run("map with special characters in string fields", func(t *testing.T) {
		mode := "mode-with-üñíçödé"
		status := "status_with_underscores"

		execution := &n8nsdk.Execution{
			Mode:   &mode,
			Status: &status,
		}

		item := mapExecutionToItem(execution)

		assert.Equal(t, mode, item.Mode.ValueString())
		assert.Equal(t, status, item.Status.ValueString())
	})

	t.Run("map timestamps at different times", func(t *testing.T) {
		startedAt := time.Now().Add(-1 * time.Hour)
		stoppedAt := time.Now()

		stoppedAtNullable := n8nsdk.NewNullableTime(&stoppedAt)

		execution := &n8nsdk.Execution{
			StartedAt: &startedAt,
			StoppedAt: *stoppedAtNullable,
		}

		item := mapExecutionToItem(execution)

		assert.Equal(t, startedAt.String(), item.StartedAt.ValueString())
		assert.Equal(t, stoppedAt.String(), item.StoppedAt.ValueString())
	})

	t.Run("map with fractional id values", func(t *testing.T) {
		id := float32(123.456)
		workflowID := float32(789.012)

		execution := &n8nsdk.Execution{
			Id:         &id,
			WorkflowId: &workflowID,
		}

		item := mapExecutionToItem(execution)

		assert.Equal(t, "123.456", item.ID.ValueString())
		assert.Equal(t, "789.012", item.WorkflowID.ValueString())
	})
}

func TestMapExecutionToItemConcurrency(t *testing.T) {
	t.Run("concurrent mapping", func(t *testing.T) {
		// Shared execution object to test thread safety
		id := float32(123)
		workflowID := float32(456)
		finished := true
		mode := "manual"
		status := "success"

		execution := &n8nsdk.Execution{
			Id:         &id,
			WorkflowId: &workflowID,
			Finished:   &finished,
			Mode:       &mode,
			Status:     &status,
		}

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				item := mapExecutionToItem(execution)
				assert.Equal(t, "123", item.ID.ValueString())
				assert.Equal(t, "456", item.WorkflowID.ValueString())
				assert.True(t, item.Finished.ValueBool())
				assert.Equal(t, "manual", item.Mode.ValueString())
				assert.Equal(t, "success", item.Status.ValueString())
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func BenchmarkMapExecutionToItem(b *testing.B) {
	id := float32(123)
	workflowID := float32(456)
	finished := true
	mode := "manual"
	startedAt := time.Now()
	stoppedAt := time.Now().Add(5 * time.Minute)
	status := "success"

	stoppedAtNullable := n8nsdk.NewNullableTime(&stoppedAt)

	execution := &n8nsdk.Execution{
		Id:         &id,
		WorkflowId: &workflowID,
		Finished:   &finished,
		Mode:       &mode,
		StartedAt:  &startedAt,
		StoppedAt:  *stoppedAtNullable,
		Status:     &status,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapExecutionToItem(execution)
	}
}
