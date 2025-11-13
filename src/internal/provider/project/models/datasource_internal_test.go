// Package models defines data structures for project resources.
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
		{name: "project types", wantErr: false},
		{name: "zero value struct", wantErr: false},
		{name: "timestamp formats", wantErr: false},
		{name: "copy struct", wantErr: false},
		{name: "partial initialization", wantErr: false},
		{name: "emoji icons", wantErr: false},
		{name: "description variations", wantErr: false},
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

			case "create with null values":
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

			case "create with unknown values":
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

			case "project types":
				projectTypes := []string{"personal", "team", "organization", "shared", "public"}

				for _, projectType := range projectTypes {
					datasource := DataSource{
						Type: types.StringValue(projectType),
					}
					assert.Equal(t, projectType, datasource.Type.ValueString())
				}

			case "zero value struct":
				var datasource DataSource
				assert.True(t, datasource.ID.IsNull())
				assert.True(t, datasource.Name.IsNull())
				assert.True(t, datasource.Type.IsNull())
				assert.True(t, datasource.CreatedAt.IsNull())
				assert.True(t, datasource.UpdatedAt.IsNull())
				assert.True(t, datasource.Icon.IsNull())
				assert.True(t, datasource.Description.IsNull())

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
					datasource := DataSource{
						CreatedAt: types.StringValue(ts),
						UpdatedAt: types.StringValue(ts),
					}
					assert.Equal(t, ts, datasource.CreatedAt.ValueString())
					assert.Equal(t, ts, datasource.UpdatedAt.ValueString())
				}

			case "copy struct":
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

			case "partial initialization":
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

			case "emoji icons":
				emojis := []string{"🚀", "📁", "💼", "🎯", "🔧", "⚡", "🌟", "📊", "🔐", "🌍"}

				for _, emoji := range emojis {
					datasource := DataSource{
						Icon: types.StringValue(emoji),
					}
					assert.Equal(t, emoji, datasource.Icon.ValueString())
				}

			case "description variations":
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

			case "error case - validation checks":
				datasource := DataSource{
					ID:          types.StringValue(""),
					Name:        types.StringValue(""),
					Type:        types.StringValue("invalid-type"),
					Description: types.StringValue(""),
				}
				assert.Equal(t, "", datasource.ID.ValueString())
				assert.Equal(t, "", datasource.Name.ValueString())
				assert.Equal(t, "invalid-type", datasource.Type.ValueString())
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
		{name: "error case - concurrent access validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// NO t.Parallel() here - concurrent goroutines don't work well with t.Parallel()
			switch tt.name {
			case "concurrent read":
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

			case "error case - concurrent access validation":
				datasource := DataSource{
					ID:   types.StringValue("val-id"),
					Name: types.StringValue("Val Project"),
				}

				done := make(chan bool, 50)
				for i := 0; i < 50; i++ {
					go func() {
						_ = datasource.ID.ValueString()
						_ = datasource.Name.ValueString()
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
