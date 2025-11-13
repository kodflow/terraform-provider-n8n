// Package models defines data structures for execution resources.
package models

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestRetryResource(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create with all fields", wantErr: false},
		{name: "create with null values", wantErr: false},
		{name: "create with unknown values", wantErr: false},
		{name: "retry modes", wantErr: false},
		{name: "status values", wantErr: false},
		{name: "finished states", wantErr: false},
		{name: "zero value struct", wantErr: false},
		{name: "timestamp formats", wantErr: false},
		{name: "copy struct", wantErr: false},
		{name: "partial initialization", wantErr: false},
		{name: "execution ID relationships", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create with all fields":
				now := time.Now().Format(time.RFC3339)
				resource := RetryResource{
					ExecutionID:    types.StringValue("exec-456"),
					NewExecutionID: types.StringValue("exec-789"),
					WorkflowID:     types.StringValue("wf-789"),
					Finished:       types.BoolValue(true),
					Mode:           types.StringValue("manual"),
					StartedAt:      types.StringValue(now),
					StoppedAt:      types.StringValue(now),
					Status:         types.StringValue("success"),
				}

				assert.Equal(t, "exec-456", resource.ExecutionID.ValueString())
				assert.Equal(t, "exec-789", resource.NewExecutionID.ValueString())
				assert.Equal(t, "wf-789", resource.WorkflowID.ValueString())
				assert.True(t, resource.Finished.ValueBool())
				assert.Equal(t, "manual", resource.Mode.ValueString())
				assert.Equal(t, now, resource.StartedAt.ValueString())
				assert.Equal(t, now, resource.StoppedAt.ValueString())
				assert.Equal(t, "success", resource.Status.ValueString())

			case "create with null values":
				resource := RetryResource{
					ExecutionID:    types.StringNull(),
					NewExecutionID: types.StringNull(),
					WorkflowID:     types.StringNull(),
					Finished:       types.BoolNull(),
					Mode:           types.StringNull(),
					StartedAt:      types.StringNull(),
					StoppedAt:      types.StringNull(),
					Status:         types.StringNull(),
				}

				assert.True(t, resource.ExecutionID.IsNull())
				assert.True(t, resource.NewExecutionID.IsNull())
				assert.True(t, resource.WorkflowID.IsNull())
				assert.True(t, resource.Finished.IsNull())
				assert.True(t, resource.Mode.IsNull())
				assert.True(t, resource.StartedAt.IsNull())
				assert.True(t, resource.StoppedAt.IsNull())
				assert.True(t, resource.Status.IsNull())

			case "create with unknown values":
				resource := RetryResource{
					ExecutionID:    types.StringUnknown(),
					NewExecutionID: types.StringUnknown(),
					WorkflowID:     types.StringUnknown(),
					Finished:       types.BoolUnknown(),
					Mode:           types.StringUnknown(),
					StartedAt:      types.StringUnknown(),
					StoppedAt:      types.StringUnknown(),
					Status:         types.StringUnknown(),
				}

				assert.True(t, resource.ExecutionID.IsUnknown())
				assert.True(t, resource.NewExecutionID.IsUnknown())
				assert.True(t, resource.WorkflowID.IsUnknown())
				assert.True(t, resource.Finished.IsUnknown())
				assert.True(t, resource.Mode.IsUnknown())
				assert.True(t, resource.StartedAt.IsUnknown())
				assert.True(t, resource.StoppedAt.IsUnknown())
				assert.True(t, resource.Status.IsUnknown())

			case "retry modes":
				modes := []string{"manual", "automatic", "scheduled", "trigger", "api"}

				for _, mode := range modes {
					resource := RetryResource{
						Mode: types.StringValue(mode),
					}
					assert.Equal(t, mode, resource.Mode.ValueString())
				}

			case "status values":
				statuses := []string{"success", "failed", "running", "waiting", "canceled", "error"}

				for _, status := range statuses {
					resource := RetryResource{
						Status: types.StringValue(status),
					}
					assert.Equal(t, status, resource.Status.ValueString())
				}

			case "finished states":
				finishedResource := RetryResource{
					Finished: types.BoolValue(true),
				}
				assert.True(t, finishedResource.Finished.ValueBool())

				runningResource := RetryResource{
					Finished: types.BoolValue(false),
				}
				assert.False(t, runningResource.Finished.ValueBool())

			case "zero value struct":
				var resource RetryResource
				assert.True(t, resource.ExecutionID.IsNull())
				assert.True(t, resource.NewExecutionID.IsNull())
				assert.True(t, resource.WorkflowID.IsNull())
				assert.True(t, resource.Finished.IsNull())
				assert.True(t, resource.Mode.IsNull())
				assert.True(t, resource.StartedAt.IsNull())
				assert.True(t, resource.StoppedAt.IsNull())
				assert.True(t, resource.Status.IsNull())

			case "timestamp formats":
				timestamps := []string{
					"2024-01-01T00:00:00Z",
					"2024-01-01T12:34:56Z",
					"2024-01-01T12:34:56.789Z",
					time.Now().Format(time.RFC3339),
				}

				for _, ts := range timestamps {
					resource := RetryResource{
						StartedAt: types.StringValue(ts),
						StoppedAt: types.StringValue(ts),
					}
					assert.Equal(t, ts, resource.StartedAt.ValueString())
					assert.Equal(t, ts, resource.StoppedAt.ValueString())
				}

			case "copy struct":
				original := RetryResource{
					ExecutionID:    types.StringValue("exec-original"),
					NewExecutionID: types.StringValue("exec-new-original"),
					WorkflowID:     types.StringValue("wf-original"),
					Finished:       types.BoolValue(true),
				}

				copied := original

				assert.Equal(t, original.ExecutionID.ValueString(), copied.ExecutionID.ValueString())
				assert.Equal(t, original.NewExecutionID.ValueString(), copied.NewExecutionID.ValueString())

				// Modify copied
				copied.ExecutionID = types.StringValue("exec-modified")
				assert.Equal(t, "exec-original", original.ExecutionID.ValueString())
				assert.Equal(t, "exec-modified", copied.ExecutionID.ValueString())

			case "partial initialization":
				// Test with only some fields set
				resource := RetryResource{
					ExecutionID: types.StringValue("exec-partial"),
					WorkflowID:  types.StringValue("wf-partial"),
					// Other fields remain null
				}

				assert.Equal(t, "exec-partial", resource.ExecutionID.ValueString())
				assert.Equal(t, "wf-partial", resource.WorkflowID.ValueString())
				assert.True(t, resource.NewExecutionID.IsNull())
				assert.True(t, resource.Finished.IsNull())
				assert.True(t, resource.Mode.IsNull())
				assert.True(t, resource.StartedAt.IsNull())
				assert.True(t, resource.StoppedAt.IsNull())
				assert.True(t, resource.Status.IsNull())

			case "execution ID relationships":
				// Test the relationship between ExecutionID and NewExecutionID
				resource := RetryResource{
					ExecutionID:    types.StringValue("original-exec"),
					NewExecutionID: types.StringValue("retry-exec"),
				}

				// In a retry scenario, ExecutionID is the original and NewExecutionID is the retry
				assert.NotEqual(t, resource.ExecutionID.ValueString(), resource.NewExecutionID.ValueString())

			case "error case - validation checks":
				// Test empty string values
				resource := RetryResource{
					ExecutionID:    types.StringValue(""),
					NewExecutionID: types.StringValue(""),
					Status:         types.StringValue("invalid-status"),
				}
				assert.Equal(t, "", resource.ExecutionID.ValueString())
				assert.Equal(t, "", resource.NewExecutionID.ValueString())
				assert.Equal(t, "invalid-status", resource.Status.ValueString())
			}
		})
	}
}

func TestRetryResourceConcurrency(t *testing.T) {
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
				resource := RetryResource{
					ExecutionID:    types.StringValue("exec-concurrent"),
					NewExecutionID: types.StringValue("exec-new-concurrent"),
					WorkflowID:     types.StringValue("wf-concurrent"),
					Finished:       types.BoolValue(true),
				}

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						_ = resource.ExecutionID.ValueString()
						_ = resource.NewExecutionID.ValueString()
						_ = resource.WorkflowID.ValueString()
						_ = resource.Finished.ValueBool()
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "error case - concurrent access validation":
				resource := RetryResource{
					ExecutionID:    types.StringValue("exec-val"),
					NewExecutionID: types.StringValue("exec-new-val"),
				}

				done := make(chan bool, 50)
				for i := 0; i < 50; i++ {
					go func() {
						_ = resource.ExecutionID.ValueString()
						_ = resource.NewExecutionID.ValueString()
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

func BenchmarkRetryResource(b *testing.B) {
	b.Run("create", func(b *testing.B) {
		now := time.Now().Format(time.RFC3339)
		for i := 0; i < b.N; i++ {
			_ = RetryResource{
				ExecutionID:    types.StringValue("exec-456"),
				NewExecutionID: types.StringValue("exec-789"),
				WorkflowID:     types.StringValue("wf-789"),
				Finished:       types.BoolValue(true),
				Mode:           types.StringValue("manual"),
				StartedAt:      types.StringValue(now),
				StoppedAt:      types.StringValue(now),
				Status:         types.StringValue("success"),
			}
		}
	})

	b.Run("access fields", func(b *testing.B) {
		resource := RetryResource{
			ExecutionID:    types.StringValue("exec-456"),
			NewExecutionID: types.StringValue("exec-789"),
			WorkflowID:     types.StringValue("wf-789"),
			Finished:       types.BoolValue(true),
			Mode:           types.StringValue("manual"),
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = resource.ExecutionID.ValueString()
			_ = resource.NewExecutionID.ValueString()
			_ = resource.WorkflowID.ValueString()
			_ = resource.Finished.ValueBool()
			_ = resource.Mode.ValueString()
		}
	})
}
