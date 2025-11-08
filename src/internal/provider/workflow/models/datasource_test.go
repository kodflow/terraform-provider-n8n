package models

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestDataSource(t *testing.T) {
	t.Run("create with all fields", func(t *testing.T) {
		datasource := DataSource{
			ID:     types.StringValue("wf-123"),
			Name:   types.StringValue("Test Workflow"),
			Active: types.BoolValue(true),
		}

		assert.Equal(t, "wf-123", datasource.ID.ValueString())
		assert.Equal(t, "Test Workflow", datasource.Name.ValueString())
		assert.True(t, datasource.Active.ValueBool())
	})

	t.Run("create with null values", func(t *testing.T) {
		datasource := DataSource{
			ID:     types.StringNull(),
			Name:   types.StringNull(),
			Active: types.BoolNull(),
		}

		assert.True(t, datasource.ID.IsNull())
		assert.True(t, datasource.Name.IsNull())
		assert.True(t, datasource.Active.IsNull())
	})

	t.Run("create with unknown values", func(t *testing.T) {
		datasource := DataSource{
			ID:     types.StringUnknown(),
			Name:   types.StringUnknown(),
			Active: types.BoolUnknown(),
		}

		assert.True(t, datasource.ID.IsUnknown())
		assert.True(t, datasource.Name.IsUnknown())
		assert.True(t, datasource.Active.IsUnknown())
	})

	t.Run("zero value struct", func(t *testing.T) {
		var datasource DataSource
		assert.True(t, datasource.ID.IsNull())
		assert.True(t, datasource.Name.IsNull())
		assert.True(t, datasource.Active.IsNull())
	})

	t.Run("active state variations", func(t *testing.T) {
		activeWorkflow := DataSource{
			ID:     types.StringValue("wf-active"),
			Name:   types.StringValue("Active Workflow"),
			Active: types.BoolValue(true),
		}
		assert.True(t, activeWorkflow.Active.ValueBool())

		inactiveWorkflow := DataSource{
			ID:     types.StringValue("wf-inactive"),
			Name:   types.StringValue("Inactive Workflow"),
			Active: types.BoolValue(false),
		}
		assert.False(t, inactiveWorkflow.Active.ValueBool())
	})

	t.Run("copy struct", func(t *testing.T) {
		original := DataSource{
			ID:     types.StringValue("original-id"),
			Name:   types.StringValue("Original Workflow"),
			Active: types.BoolValue(true),
		}

		copy := original

		assert.Equal(t, original.ID.ValueString(), copy.ID.ValueString())
		assert.Equal(t, original.Name.ValueString(), copy.Name.ValueString())
		assert.Equal(t, original.Active.ValueBool(), copy.Active.ValueBool())

		// Modify copy
		copy.ID = types.StringValue("modified-id")
		copy.Name = types.StringValue("Modified Workflow")
		copy.Active = types.BoolValue(false)

		assert.Equal(t, "original-id", original.ID.ValueString())
		assert.Equal(t, "modified-id", copy.ID.ValueString())
		assert.Equal(t, "Original Workflow", original.Name.ValueString())
		assert.Equal(t, "Modified Workflow", copy.Name.ValueString())
		assert.True(t, original.Active.ValueBool())
		assert.False(t, copy.Active.ValueBool())
	})

	t.Run("name variations", func(t *testing.T) {
		names := []string{
			"Simple Workflow",
			"workflow-with-dashes",
			"workflow_with_underscores",
			"Workflow With Numbers 123",
			"Workflow.With.Dots",
			"Workflow/With/Slashes",
			"Workflow@With#Special$Chars",
			"Unicode工作流程",
			"",
		}

		for _, name := range names {
			datasource := DataSource{
				Name: types.StringValue(name),
			}
			assert.Equal(t, name, datasource.Name.ValueString())
		}
	})

	t.Run("partial initialization", func(t *testing.T) {
		// Only ID set
		datasource1 := DataSource{
			ID: types.StringValue("wf-id-only"),
		}
		assert.False(t, datasource1.ID.IsNull())
		assert.True(t, datasource1.Name.IsNull())
		assert.True(t, datasource1.Active.IsNull())

		// Only Name set
		datasource2 := DataSource{
			Name: types.StringValue("Name Only"),
		}
		assert.True(t, datasource2.ID.IsNull())
		assert.False(t, datasource2.Name.IsNull())
		assert.True(t, datasource2.Active.IsNull())

		// Only Active set
		datasource3 := DataSource{
			Active: types.BoolValue(true),
		}
		assert.True(t, datasource3.ID.IsNull())
		assert.True(t, datasource3.Name.IsNull())
		assert.False(t, datasource3.Active.IsNull())
	})
}

func TestDataSourceConcurrency(t *testing.T) {
	t.Run("concurrent read", func(t *testing.T) {
		datasource := DataSource{
			ID:     types.StringValue("concurrent-id"),
			Name:   types.StringValue("Concurrent Workflow"),
			Active: types.BoolValue(true),
		}

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				_ = datasource.ID.ValueString()
				_ = datasource.Name.ValueString()
				_ = datasource.Active.ValueBool()
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
		for i := 0; i < b.N; i++ {
			_ = DataSource{
				ID:     types.StringValue("wf-123"),
				Name:   types.StringValue("Test Workflow"),
				Active: types.BoolValue(true),
			}
		}
	})

	b.Run("access fields", func(b *testing.B) {
		datasource := DataSource{
			ID:     types.StringValue("wf-123"),
			Name:   types.StringValue("Test Workflow"),
			Active: types.BoolValue(true),
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = datasource.ID.ValueString()
			_ = datasource.Name.ValueString()
			_ = datasource.Active.ValueBool()
		}
	})
}
