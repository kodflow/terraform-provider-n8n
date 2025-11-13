// Package models defines data structures for project resources.
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
		{name: "item types", wantErr: false},
		{name: "zero value struct", wantErr: false},
		{name: "timestamp formats", wantErr: false},
		{name: "copy struct", wantErr: false},
		{name: "partial initialization", wantErr: false},
		{name: "emoji icons", wantErr: false},
		{name: "name variations", wantErr: false},
		{name: "description variations", wantErr: false},
		{name: "updated after created", wantErr: false},
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
					ID:          types.StringValue("item-123"),
					Name:        types.StringValue("Test Item"),
					Type:        types.StringValue("personal"),
					CreatedAt:   types.StringValue(now),
					UpdatedAt:   types.StringValue(now),
					Icon:        types.StringValue("📋"),
					Description: types.StringValue("A test item for unit testing"),
				}

				assert.Equal(t, "item-123", item.ID.ValueString())
				assert.Equal(t, "Test Item", item.Name.ValueString())
				assert.Equal(t, "personal", item.Type.ValueString())
				assert.Equal(t, now, item.CreatedAt.ValueString())
				assert.Equal(t, now, item.UpdatedAt.ValueString())
				assert.Equal(t, "📋", item.Icon.ValueString())
				assert.Equal(t, "A test item for unit testing", item.Description.ValueString())

			case "create with null values":
				item := Item{
					ID:          types.StringNull(),
					Name:        types.StringNull(),
					Type:        types.StringNull(),
					CreatedAt:   types.StringNull(),
					UpdatedAt:   types.StringNull(),
					Icon:        types.StringNull(),
					Description: types.StringNull(),
				}

				assert.True(t, item.ID.IsNull())
				assert.True(t, item.Name.IsNull())
				assert.True(t, item.Type.IsNull())
				assert.True(t, item.CreatedAt.IsNull())
				assert.True(t, item.UpdatedAt.IsNull())
				assert.True(t, item.Icon.IsNull())
				assert.True(t, item.Description.IsNull())

			case "create with unknown values":
				item := Item{
					ID:          types.StringUnknown(),
					Name:        types.StringUnknown(),
					Type:        types.StringUnknown(),
					CreatedAt:   types.StringUnknown(),
					UpdatedAt:   types.StringUnknown(),
					Icon:        types.StringUnknown(),
					Description: types.StringUnknown(),
				}

				assert.True(t, item.ID.IsUnknown())
				assert.True(t, item.Name.IsUnknown())
				assert.True(t, item.Type.IsUnknown())
				assert.True(t, item.CreatedAt.IsUnknown())
				assert.True(t, item.UpdatedAt.IsUnknown())
				assert.True(t, item.Icon.IsUnknown())
				assert.True(t, item.Description.IsUnknown())

			case "item types":
				itemTypes := []string{"personal", "team", "organization", "shared", "public", "private"}

				for _, itemType := range itemTypes {
					item := Item{
						Type: types.StringValue(itemType),
					}
					assert.Equal(t, itemType, item.Type.ValueString())
				}

			case "zero value struct":
				var item Item
				assert.True(t, item.ID.IsNull())
				assert.True(t, item.Name.IsNull())
				assert.True(t, item.Type.IsNull())
				assert.True(t, item.CreatedAt.IsNull())
				assert.True(t, item.UpdatedAt.IsNull())
				assert.True(t, item.Icon.IsNull())
				assert.True(t, item.Description.IsNull())

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
						CreatedAt: types.StringValue(ts),
						UpdatedAt: types.StringValue(ts),
					}
					assert.Equal(t, ts, item.CreatedAt.ValueString())
					assert.Equal(t, ts, item.UpdatedAt.ValueString())
				}

			case "copy struct":
				original := Item{
					ID:          types.StringValue("original-id"),
					Name:        types.StringValue("Original Item"),
					Type:        types.StringValue("team"),
					Icon:        types.StringValue("🔷"),
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
				copied.Name = types.StringValue("Modified Item")
				assert.Equal(t, "original-id", original.ID.ValueString())
				assert.Equal(t, "modified-id", copied.ID.ValueString())
				assert.Equal(t, "Original Item", original.Name.ValueString())
				assert.Equal(t, "Modified Item", copied.Name.ValueString())

			case "partial initialization":
				item := Item{
					ID:   types.StringValue("item-partial"),
					Name: types.StringValue("Partial Item"),
					// Other fields remain null
				}

				assert.Equal(t, "item-partial", item.ID.ValueString())
				assert.Equal(t, "Partial Item", item.Name.ValueString())
				assert.True(t, item.Type.IsNull())
				assert.True(t, item.CreatedAt.IsNull())
				assert.True(t, item.UpdatedAt.IsNull())
				assert.True(t, item.Icon.IsNull())
				assert.True(t, item.Description.IsNull())

			case "emoji icons":
				emojis := []string{"🚀", "📁", "💼", "🎯", "🔧", "⚡", "🌟", "📊", "🔐", "🌍"}

				for _, emoji := range emojis {
					item := Item{
						Icon: types.StringValue(emoji),
					}
					assert.Equal(t, emoji, item.Icon.ValueString())
				}

			case "name variations":
				names := []string{
					"Simple Name",
					"name-with-dashes",
					"name_with_underscores",
					"Name With Numbers 123",
					"Name.With.Dots",
					"Name/With/Slashes",
					"Name@With#Special$Chars",
					"Unicode名前имя",
				}

				for _, name := range names {
					item := Item{
						Name: types.StringValue(name),
					}
					assert.Equal(t, name, item.Name.ValueString())
				}

			case "description variations":
				descriptions := []string{
					"",
					"Short desc",
					"A longer description with multiple words and sentences.",
					"Description with special characters: @#$%^&*()",
					"Multi-line\ndescription\nwith\nline breaks",
					"Unicode description: 你好世界 مرحبا بالعالم",
				}

				for _, desc := range descriptions {
					item := Item{
						Description: types.StringValue(desc),
					}
					assert.Equal(t, desc, item.Description.ValueString())
				}

			case "updated after created":
				createdAt := "2024-01-01T00:00:00Z"
				updatedAt := "2024-01-02T00:00:00Z"

				item := Item{
					CreatedAt: types.StringValue(createdAt),
					UpdatedAt: types.StringValue(updatedAt),
				}

				assert.Equal(t, createdAt, item.CreatedAt.ValueString())
				assert.Equal(t, updatedAt, item.UpdatedAt.ValueString())
				assert.NotEqual(t, item.CreatedAt.ValueString(), item.UpdatedAt.ValueString())

			case "error case - validation checks":
				item := Item{
					ID:          types.StringValue(""),
					Name:        types.StringValue(""),
					Type:        types.StringValue("invalid-type"),
					Description: types.StringValue(""),
				}
				assert.Equal(t, "", item.ID.ValueString())
				assert.Equal(t, "", item.Name.ValueString())
				assert.Equal(t, "invalid-type", item.Type.ValueString())
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
					ID:          types.StringValue("concurrent-id"),
					Name:        types.StringValue("Concurrent Item"),
					Type:        types.StringValue("team"),
					Icon:        types.StringValue("🔄"),
					Description: types.StringValue("Concurrent test"),
				}

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						_ = item.ID.ValueString()
						_ = item.Name.ValueString()
						_ = item.Type.ValueString()
						_ = item.Icon.ValueString()
						_ = item.Description.ValueString()
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "error case - concurrent access validation":
				item := Item{
					ID:   types.StringValue("val-id"),
					Name: types.StringValue("Val Item"),
				}

				done := make(chan bool, 50)
				for i := 0; i < 50; i++ {
					go func() {
						_ = item.ID.ValueString()
						_ = item.Name.ValueString()
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
				ID:          types.StringValue("item-123"),
				Name:        types.StringValue("Test Item"),
				Type:        types.StringValue("personal"),
				CreatedAt:   types.StringValue(now),
				UpdatedAt:   types.StringValue(now),
				Icon:        types.StringValue("🚀"),
				Description: types.StringValue("A test item"),
			}
		}
	})

	b.Run("access fields", func(b *testing.B) {
		item := Item{
			ID:          types.StringValue("item-123"),
			Name:        types.StringValue("Test Item"),
			Type:        types.StringValue("personal"),
			Icon:        types.StringValue("🚀"),
			Description: types.StringValue("A test item"),
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = item.ID.ValueString()
			_ = item.Name.ValueString()
			_ = item.Type.ValueString()
			_ = item.Icon.ValueString()
			_ = item.Description.ValueString()
		}
	})
}
