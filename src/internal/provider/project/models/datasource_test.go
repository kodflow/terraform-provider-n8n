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
			ID:          types.StringValue("proj-123"),
			Name:        types.StringValue("Test Project"),
			Type:        types.StringValue("personal"),
			CreatedAt:   types.StringValue(now),
			UpdatedAt:   types.StringValue(now),
			Icon:        types.StringValue("🚀"),
			Description: types.StringValue("A test project for unit testing"),
		}

		assert.Equal(t, "proj-123", datasource.ID.ValueString())
		assert.Equal(t, "Test Project", datasource.Name.ValueString())
		assert.Equal(t, "personal", datasource.Type.ValueString())
		assert.Equal(t, now, datasource.CreatedAt.ValueString())
		assert.Equal(t, now, datasource.UpdatedAt.ValueString())
		assert.Equal(t, "🚀", datasource.Icon.ValueString())
		assert.Equal(t, "A test project for unit testing", datasource.Description.ValueString())
	})

	t.Run("create with null values", func(t *testing.T) {
		datasource := DataSource{
			ID:          types.StringNull(),
			Name:        types.StringNull(),
			Type:        types.StringNull(),
			CreatedAt:   types.StringNull(),
			UpdatedAt:   types.StringNull(),
			Icon:        types.StringNull(),
			Description: types.StringNull(),
		}

		assert.True(t, datasource.ID.IsNull())
		assert.True(t, datasource.Name.IsNull())
		assert.True(t, datasource.Type.IsNull())
		assert.True(t, datasource.CreatedAt.IsNull())
		assert.True(t, datasource.UpdatedAt.IsNull())
		assert.True(t, datasource.Icon.IsNull())
		assert.True(t, datasource.Description.IsNull())
	})

	t.Run("create with unknown values", func(t *testing.T) {
		datasource := DataSource{
			ID:          types.StringUnknown(),
			Name:        types.StringUnknown(),
			Type:        types.StringUnknown(),
			CreatedAt:   types.StringUnknown(),
			UpdatedAt:   types.StringUnknown(),
			Icon:        types.StringUnknown(),
			Description: types.StringUnknown(),
		}

		assert.True(t, datasource.ID.IsUnknown())
		assert.True(t, datasource.Name.IsUnknown())
		assert.True(t, datasource.Type.IsUnknown())
		assert.True(t, datasource.CreatedAt.IsUnknown())
		assert.True(t, datasource.UpdatedAt.IsUnknown())
		assert.True(t, datasource.Icon.IsUnknown())
		assert.True(t, datasource.Description.IsUnknown())
	})

	t.Run("project types", func(t *testing.T) {
		projectTypes := []string{"personal", "team", "organization", "shared", "public"}

		for _, projectType := range projectTypes {
			datasource := DataSource{
				Type: types.StringValue(projectType),
			}
			assert.Equal(t, projectType, datasource.Type.ValueString())
		}
	})

	t.Run("zero value struct", func(t *testing.T) {
		var datasource DataSource
		assert.True(t, datasource.ID.IsNull())
		assert.True(t, datasource.Name.IsNull())
		assert.True(t, datasource.Type.IsNull())
		assert.True(t, datasource.CreatedAt.IsNull())
		assert.True(t, datasource.UpdatedAt.IsNull())
		assert.True(t, datasource.Icon.IsNull())
		assert.True(t, datasource.Description.IsNull())
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
			datasource := DataSource{
				CreatedAt: types.StringValue(ts),
				UpdatedAt: types.StringValue(ts),
			}
			assert.Equal(t, ts, datasource.CreatedAt.ValueString())
			assert.Equal(t, ts, datasource.UpdatedAt.ValueString())
		}
	})

	t.Run("copy struct", func(t *testing.T) {
		original := DataSource{
			ID:          types.StringValue("original-id"),
			Name:        types.StringValue("Original Project"),
			Type:        types.StringValue("team"),
			Icon:        types.StringValue("📁"),
			Description: types.StringValue("Original description"),
		}

		copied := original

		assert.Equal(t, original.ID.ValueString(), copied.ID.ValueString())
		assert.Equal(t, original.Name.ValueString(), copied.Name.ValueString())
		assert.Equal(t, original.Type.ValueString(), copied.Type.ValueString())
		assert.Equal(t, original.Icon.ValueString(), copied.Icon.ValueString())
		assert.Equal(t, original.Description.ValueString(), copied.Description.ValueString())

		// Modify copied
		copied.ID = types.StringValue("modified-id")
		copied.Name = types.StringValue("Modified Project")
		assert.Equal(t, "original-id", original.ID.ValueString())
		assert.Equal(t, "modified-id", copied.ID.ValueString())
		assert.Equal(t, "Original Project", original.Name.ValueString())
		assert.Equal(t, "Modified Project", copied.Name.ValueString())
	})

	t.Run("partial initialization", func(t *testing.T) {
		datasource := DataSource{
			ID:   types.StringValue("proj-partial"),
			Name: types.StringValue("Partial Project"),
			// Other fields remain null
		}

		assert.Equal(t, "proj-partial", datasource.ID.ValueString())
		assert.Equal(t, "Partial Project", datasource.Name.ValueString())
		assert.True(t, datasource.Type.IsNull())
		assert.True(t, datasource.CreatedAt.IsNull())
		assert.True(t, datasource.UpdatedAt.IsNull())
		assert.True(t, datasource.Icon.IsNull())
		assert.True(t, datasource.Description.IsNull())
	})

	t.Run("emoji icons", func(t *testing.T) {
		emojis := []string{"🚀", "📁", "💼", "🎯", "🔧", "⚡", "🌟", "📊", "🔐", "🌍"}

		for _, emoji := range emojis {
			datasource := DataSource{
				Icon: types.StringValue(emoji),
			}
			assert.Equal(t, emoji, datasource.Icon.ValueString())
		}
	})

	t.Run("description variations", func(t *testing.T) {
		descriptions := []string{
			"",
			"Short desc",
			"A longer description with multiple words",
			"Description with special characters: @#$%^&*()",
			"Multi-line\ndescription\nwith breaks",
			"Unicode description: 你好世界 مرحبا بالعالم",
		}

		for _, desc := range descriptions {
			datasource := DataSource{
				Description: types.StringValue(desc),
			}
			assert.Equal(t, desc, datasource.Description.ValueString())
		}
	})
}

func TestDataSourceConcurrency(t *testing.T) {
	t.Run("concurrent read", func(t *testing.T) {
		datasource := DataSource{
			ID:          types.StringValue("concurrent-id"),
			Name:        types.StringValue("Concurrent Project"),
			Type:        types.StringValue("team"),
			Icon:        types.StringValue("🔄"),
			Description: types.StringValue("Concurrent test"),
		}

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				_ = datasource.ID.ValueString()
				_ = datasource.Name.ValueString()
				_ = datasource.Type.ValueString()
				_ = datasource.Icon.ValueString()
				_ = datasource.Description.ValueString()
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
				ID:          types.StringValue("proj-123"),
				Name:        types.StringValue("Test Project"),
				Type:        types.StringValue("personal"),
				CreatedAt:   types.StringValue(now),
				UpdatedAt:   types.StringValue(now),
				Icon:        types.StringValue("🚀"),
				Description: types.StringValue("A test project"),
			}
		}
	})

	b.Run("access fields", func(b *testing.B) {
		datasource := DataSource{
			ID:          types.StringValue("proj-123"),
			Name:        types.StringValue("Test Project"),
			Type:        types.StringValue("personal"),
			Icon:        types.StringValue("🚀"),
			Description: types.StringValue("A test project"),
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = datasource.ID.ValueString()
			_ = datasource.Name.ValueString()
			_ = datasource.Type.ValueString()
			_ = datasource.Icon.ValueString()
			_ = datasource.Description.ValueString()
		}
	})
}
