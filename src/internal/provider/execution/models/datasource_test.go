package models

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestDataSource(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create with all fields", wantErr: false},
		{name: "create with null values", wantErr: false},
		{name: "create with unknown values", wantErr: false},
		{name: "include data flag", wantErr: false},
		{name: "execution modes", wantErr: false},
		{name: "status values", wantErr: false},
		{name: "finished states", wantErr: false},
		{name: "zero value struct", wantErr: false},
		{name: "timestamp formats", wantErr: false},
		{name: "copy struct", wantErr: false},
		{name: "partial initialization", wantErr: false},
		{name: "execution states with data", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create with all fields":
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

			case "create with null values":
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

			case "create with unknown values":
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

			case "include data flag":
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

			case "execution modes":
				modes := []string{"manual", "trigger", "schedule", "webhook", "api", "error", "retry", "integrated"}

				for _, mode := range modes {
					datasource := DataSource{
						Mode: types.StringValue(mode),
					}
					assert.Equal(t, mode, datasource.Mode.ValueString())
				}

			case "status values":
				statuses := []string{"success", "error", "waiting", "running", "canceled", "crashed", "new"}

				for _, status := range statuses {
					datasource := DataSource{
						Status: types.StringValue(status),
					}
					assert.Equal(t, status, datasource.Status.ValueString())
				}

			case "finished states":
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

			case "zero value struct":
				var datasource DataSource
				assert.True(t, datasource.ID.IsNull())
				assert.True(t, datasource.WorkflowID.IsNull())
				assert.True(t, datasource.Finished.IsNull())
				assert.True(t, datasource.Mode.IsNull())
				assert.True(t, datasource.StartedAt.IsNull())
				assert.True(t, datasource.StoppedAt.IsNull())
				assert.True(t, datasource.Status.IsNull())
				assert.True(t, datasource.IncludeData.IsNull())

			case "timestamp formats":
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

			case "copy struct":
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

			case "partial initialization":
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

			case "execution states with data":
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

			case "error case - validation checks":
				// Test invalid/edge case values
				datasource1 := DataSource{
					ID:          types.StringValue(""),
					WorkflowID:  types.StringValue(""),
					IncludeData: types.BoolValue(true),
				}
				assert.Equal(t, "", datasource1.ID.ValueString())
				assert.True(t, datasource1.IncludeData.ValueBool())

				// Test mixed states
				datasource2 := DataSource{
					ID:       types.StringNull(),
					Finished: types.BoolValue(true),
					Status:   types.StringValue("unknown"),
				}
				assert.True(t, datasource2.ID.IsNull())
				assert.True(t, datasource2.Finished.ValueBool())
			}
		})
	}
}

func TestDataSourceConcurrency(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "concurrent read", wantErr: false},
		{name: "error case - concurrent access patterns", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "concurrent read":
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

			case "error case - concurrent access patterns":
				// Test concurrent access with various states
				datasources := []DataSource{
					{ID: types.StringValue("exec-1"), Finished: types.BoolValue(true)},
					{ID: types.StringValue("exec-2"), Finished: types.BoolValue(false)},
					{ID: types.StringNull(), Finished: types.BoolNull()},
				}

				done := make(chan bool, len(datasources)*10)
				for _, ds := range datasources {
					for i := 0; i < 10; i++ {
						go func(d DataSource) {
							_ = d.ID.IsNull()
							_ = d.Finished.IsNull()
							done <- true
						}(ds)
					}
				}

				for i := 0; i < len(datasources)*10; i++ {
					<-done
				}
			}
		})
	}
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
