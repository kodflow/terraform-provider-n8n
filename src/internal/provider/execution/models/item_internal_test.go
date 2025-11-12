// Package models defines data structures for execution resources.
package models

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestItem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create with all fields", wantErr: false},
		{name: "create with null values", wantErr: false},
		{name: "create with unknown values", wantErr: false},
		{name: "execution modes", wantErr: false},
		{name: "status values", wantErr: false},
		{name: "finished states", wantErr: false},
		{name: "zero value struct", wantErr: false},
		{name: "timestamp formats", wantErr: false},
		{name: "copy struct", wantErr: false},
		{name: "partial initialization", wantErr: false},
		{name: "execution states", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create with all fields":
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

			case "create with null values":
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

			case "create with unknown values":
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

			case "execution modes":
				modes := []string{"manual", "trigger", "schedule", "webhook", "api", "error", "retry", "internal"}

				for _, mode := range modes {
					item := Item{
						Mode: types.StringValue(mode),
					}
					assert.Equal(t, mode, item.Mode.ValueString())
				}

			case "status values":
				statuses := []string{"success", "error", "waiting", "running", "canceled", "crashed", "unknown"}

				for _, status := range statuses {
					item := Item{
						Status: types.StringValue(status),
					}
					assert.Equal(t, status, item.Status.ValueString())
				}

			case "finished states":
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

			case "zero value struct":
				var item Item
				assert.True(t, item.ID.IsNull())
				assert.True(t, item.WorkflowID.IsNull())
				assert.True(t, item.Finished.IsNull())
				assert.True(t, item.Mode.IsNull())
				assert.True(t, item.StartedAt.IsNull())
				assert.True(t, item.StoppedAt.IsNull())
				assert.True(t, item.Status.IsNull())

			case "timestamp formats":
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

			case "copy struct":
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

			case "partial initialization":
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

			case "execution states":
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

			case "error case - validation checks":
				// Test empty string values
				item := Item{
					ID:         types.StringValue(""),
					WorkflowID: types.StringValue(""),
					Status:     types.StringValue("invalid-status"),
				}
				assert.Equal(t, "", item.ID.ValueString())
				assert.Equal(t, "", item.WorkflowID.ValueString())
				assert.Equal(t, "invalid-status", item.Status.ValueString())
			}
		})
	}
}

func TestItemConcurrency(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "concurrent read", wantErr: false},
		{name: "error case - concurrent access validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// NO t.Parallel() here - concurrent goroutines don't work well with t.Parallel()
			switch tt.name {
			case "concurrent read":
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

			case "error case - concurrent access validation":
				item := Item{
					ID:         types.StringValue("exec-val"),
					WorkflowID: types.StringValue("wf-val"),
				}

				done := make(chan bool, 50)
				for i := 0; i < 50; i++ {
					go func() {
						_ = item.ID.ValueString()
						_ = item.WorkflowID.ValueString()
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
