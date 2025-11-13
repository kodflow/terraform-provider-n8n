package execution

import (
	"testing"
	"time"

	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/stretchr/testify/assert"
)

func Test_mapExecutionToItem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "map with all fields populated", wantErr: false},
		{name: "map with nil fields", wantErr: false},
		{name: "map with partial fields", wantErr: false},
		{name: "map with unset nullable stoppedAt", wantErr: false},
		{name: "map different execution modes", wantErr: false},
		{name: "map different status values", wantErr: false},
		{name: "map with zero value id", wantErr: false},
		{name: "map with negative id", wantErr: false},
		{name: "map with large id", wantErr: false},
		{name: "map finished states", wantErr: false},
		{name: "map with empty string values", wantErr: false},
		{name: "map with special characters in string fields", wantErr: false},
		{name: "map timestamps at different times", wantErr: false},
		{name: "map with fractional id values", wantErr: false},
		{name: "error case - nil execution", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "map with all fields populated":
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
					StartedAt:  *n8nsdk.NewNullableTime(&startedAt),
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

			case "map with nil fields":
				execution := &n8nsdk.Execution{}

				item := mapExecutionToItem(execution)

				assert.True(t, item.ID.IsNull())
				assert.True(t, item.WorkflowID.IsNull())
				assert.True(t, item.Finished.IsNull())
				assert.True(t, item.Mode.IsNull())
				assert.True(t, item.StartedAt.IsNull())
				assert.True(t, item.StoppedAt.IsNull())
				assert.True(t, item.Status.IsNull())

			case "map with partial fields":
				id := float32(789)
				workflowID := float32(999)
				status := "running"

				execution := &n8nsdk.Execution{
					Id:         &id,
					WorkflowId: &workflowID,
					Status:     &status,
				}

				item := mapExecutionToItem(execution)

				assert.Equal(t, "789", item.ID.ValueString())
				assert.Equal(t, "999", item.WorkflowID.ValueString())
				assert.Equal(t, "running", item.Status.ValueString())
				assert.True(t, item.Finished.IsNull())
				assert.True(t, item.Mode.IsNull())
				assert.True(t, item.StartedAt.IsNull())
				assert.True(t, item.StoppedAt.IsNull())

			case "map with unset nullable stoppedAt":
				id := float32(999)
				workflowID := float32(888)
				finished := false
				status := "running"

				execution := &n8nsdk.Execution{
					Id:         &id,
					WorkflowId: &workflowID,
					Finished:   &finished,
					Status:     &status,
					StoppedAt:  n8nsdk.NullableTime{},
				}

				item := mapExecutionToItem(execution)

				assert.Equal(t, "999", item.ID.ValueString())
				assert.Equal(t, "888", item.WorkflowID.ValueString())
				assert.False(t, item.Finished.ValueBool())
				assert.Equal(t, "running", item.Status.ValueString())
				assert.True(t, item.StoppedAt.IsNull())

			case "map different execution modes":
				modes := []string{"manual", "trigger", "schedule", "webhook", "api", "error", "retry"}

				for _, mode := range modes {
					modeCopy := mode
					execution := &n8nsdk.Execution{
						Mode: &modeCopy,
					}

					item := mapExecutionToItem(execution)

					assert.Equal(t, mode, item.Mode.ValueString())
				}

			case "map different status values":
				statuses := []string{"success", "error", "waiting", "running", "canceled", "crashed", "new"}

				for _, status := range statuses {
					statusCopy := status
					execution := &n8nsdk.Execution{
						Status: &statusCopy,
					}

					item := mapExecutionToItem(execution)

					assert.Equal(t, status, item.Status.ValueString())
				}

			case "map with zero value id":
				id := float32(0)
				execution := &n8nsdk.Execution{
					Id: &id,
				}

				item := mapExecutionToItem(execution)

				assert.Equal(t, "0", item.ID.ValueString())

			case "map with negative id":
				id := float32(-1)
				execution := &n8nsdk.Execution{
					Id: &id,
				}

				item := mapExecutionToItem(execution)

				assert.Equal(t, "-1", item.ID.ValueString())

			case "map with large id":
				id := float32(2147483647)
				execution := &n8nsdk.Execution{
					Id: &id,
				}

				item := mapExecutionToItem(execution)

				assert.Contains(t, item.ID.ValueString(), "2.1474")

			case "map finished states":
				finished := true
				finishedExec := &n8nsdk.Execution{
					Finished: &finished,
				}

				finishedItem := mapExecutionToItem(finishedExec)
				assert.True(t, finishedItem.Finished.ValueBool())

				running := false
				runningExec := &n8nsdk.Execution{
					Finished: &running,
				}

				runningItem := mapExecutionToItem(runningExec)
				assert.False(t, runningItem.Finished.ValueBool())

			case "map with empty string values":
				mode := ""
				status := ""

				execution := &n8nsdk.Execution{
					Mode:   &mode,
					Status: &status,
				}

				item := mapExecutionToItem(execution)

				assert.Equal(t, "", item.Mode.ValueString())
				assert.Equal(t, "", item.Status.ValueString())

			case "map with special characters in string fields":
				mode := "mode-with-üñíçödé"
				status := "status_with_underscores"

				execution := &n8nsdk.Execution{
					Mode:   &mode,
					Status: &status,
				}

				item := mapExecutionToItem(execution)

				assert.Equal(t, mode, item.Mode.ValueString())
				assert.Equal(t, status, item.Status.ValueString())

			case "map timestamps at different times":
				startedAt := time.Now().Add(-1 * time.Hour)
				stoppedAt := time.Now()

				stoppedAtNullable := n8nsdk.NewNullableTime(&stoppedAt)

				execution := &n8nsdk.Execution{
					StartedAt: *n8nsdk.NewNullableTime(&startedAt),
					StoppedAt: *stoppedAtNullable,
				}

				item := mapExecutionToItem(execution)

				assert.Equal(t, startedAt.String(), item.StartedAt.ValueString())
				assert.Equal(t, stoppedAt.String(), item.StoppedAt.ValueString())

			case "map with fractional id values":
				id := float32(123.456)
				workflowID := float32(789.012)

				execution := &n8nsdk.Execution{
					Id:         &id,
					WorkflowId: &workflowID,
				}

				item := mapExecutionToItem(execution)

				assert.Equal(t, "123.456", item.ID.ValueString())
				assert.Equal(t, "789.012", item.WorkflowID.ValueString())

			case "error case - nil execution":
				execution := &n8nsdk.Execution{}
				item := mapExecutionToItem(execution)
				assert.NotNil(t, item)
			}
		})
	}
}

func Test_mapExecutionToItemConcurrency(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "concurrent mapping", wantErr: false},
		{name: "error case - concurrent access validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// NO t.Parallel() here - concurrent goroutines don't work well with t.Parallel()
			switch tt.name {
			case "concurrent mapping":
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

			case "error case - concurrent access validation":
				id := float32(999)
				execution := &n8nsdk.Execution{
					Id: &id,
				}

				done := make(chan bool, 50)
				for i := 0; i < 50; i++ {
					go func() {
						item := mapExecutionToItem(execution)
						assert.Equal(t, "999", item.ID.ValueString())
						done <- true
					}()
				}

				for i := 0; i < 50; i++ {
					<-done
				}
			}
		})
	}
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
		StartedAt:  *n8nsdk.NewNullableTime(&startedAt),
		StoppedAt:  *stoppedAtNullable,
		Status:     &status,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapExecutionToItem(execution)
	}
}
