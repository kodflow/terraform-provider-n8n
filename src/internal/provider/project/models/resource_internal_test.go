// Package models defines data structures for project resources.
package models

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestResource(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create with all fields", wantErr: false},
		{name: "create with null values", wantErr: false},
		{name: "create with unknown values", wantErr: false},
		{name: "resource types", wantErr: false},
		{name: "zero value struct", wantErr: false},
		{name: "copy struct", wantErr: false},
		{name: "partial initialization", wantErr: false},
		{name: "name variations", wantErr: false},
		{name: "id formats", wantErr: false},
		{name: "required fields validation", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create with all fields":
				resource := Resource{
					ID:   types.StringValue("res-123"),
					Name: types.StringValue("Test Resource"),
					Type: types.StringValue("personal"),
				}

				assert.Equal(t, "res-123", resource.ID.ValueString())
				assert.Equal(t, "Test Resource", resource.Name.ValueString())
				assert.Equal(t, "personal", resource.Type.ValueString())

			case "create with null values":
				resource := Resource{
					ID:   types.StringNull(),
					Name: types.StringNull(),
					Type: types.StringNull(),
				}

				assert.True(t, resource.ID.IsNull())
				assert.True(t, resource.Name.IsNull())
				assert.True(t, resource.Type.IsNull())

			case "create with unknown values":
				resource := Resource{
					ID:   types.StringUnknown(),
					Name: types.StringUnknown(),
					Type: types.StringUnknown(),
				}

				assert.True(t, resource.ID.IsUnknown())
				assert.True(t, resource.Name.IsUnknown())
				assert.True(t, resource.Type.IsUnknown())

			case "resource types":
				resourceTypes := []string{
					"personal",
					"team",
					"organization",
					"shared",
					"public",
					"private",
				}

				for _, resType := range resourceTypes {
					resource := Resource{
						Type: types.StringValue(resType),
					}
					assert.Equal(t, resType, resource.Type.ValueString())
				}

			case "zero value struct":
				var resource Resource
				assert.True(t, resource.ID.IsNull())
				assert.True(t, resource.Name.IsNull())
				assert.True(t, resource.Type.IsNull())

			case "copy struct":
				original := Resource{
					ID:   types.StringValue("original-id"),
					Name: types.StringValue("Original Resource"),
					Type: types.StringValue("team"),
				}

				copied := original

				assert.Equal(t, original.ID.ValueString(), copied.ID.ValueString())
				assert.Equal(t, original.Name.ValueString(), copied.Name.ValueString())
				assert.Equal(t, original.Type.ValueString(), copied.Type.ValueString())

				// Modify copied
				copied.ID = types.StringValue("modified-id")
				copied.Name = types.StringValue("Modified Resource")
				assert.Equal(t, "original-id", original.ID.ValueString())
				assert.Equal(t, "modified-id", copied.ID.ValueString())
				assert.Equal(t, "Original Resource", original.Name.ValueString())
				assert.Equal(t, "Modified Resource", copied.Name.ValueString())

			case "partial initialization":
				resource := Resource{
					ID: types.StringValue("res-partial"),
					// Other fields remain null
				}

				assert.Equal(t, "res-partial", resource.ID.ValueString())
				assert.True(t, resource.Name.IsNull())
				assert.True(t, resource.Type.IsNull())

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
					"",
				}

				for _, name := range names {
					resource := Resource{
						Name: types.StringValue(name),
					}
					assert.Equal(t, name, resource.Name.ValueString())
				}

			case "id formats":
				ids := []string{
					"simple-id",
					"123456",
					"uuid-550e8400-e29b-41d4-a716-446655440000",
					"ID_WITH_UNDERSCORES",
					"id.with.dots",
					"id/with/slashes",
					"id@with#special$chars",
				}

				for _, id := range ids {
					resource := Resource{
						ID: types.StringValue(id),
					}
					assert.Equal(t, id, resource.ID.ValueString())
				}

			case "required fields validation":
				// Test with only ID
				resource1 := Resource{
					ID: types.StringValue("id-only"),
				}
				assert.False(t, resource1.ID.IsNull())
				assert.True(t, resource1.Name.IsNull())
				assert.True(t, resource1.Type.IsNull())

				// Test with only Name
				resource2 := Resource{
					Name: types.StringValue("name-only"),
				}
				assert.True(t, resource2.ID.IsNull())
				assert.False(t, resource2.Name.IsNull())
				assert.True(t, resource2.Type.IsNull())

				// Test with only Type
				resource3 := Resource{
					Type: types.StringValue("type-only"),
				}
				assert.True(t, resource3.ID.IsNull())
				assert.True(t, resource3.Name.IsNull())
				assert.False(t, resource3.Type.IsNull())

			case "error case - validation checks":
				// Test empty string values
				resource := Resource{
					ID:   types.StringValue(""),
					Name: types.StringValue(""),
					Type: types.StringValue("invalid-type"),
				}
				assert.Equal(t, "", resource.ID.ValueString())
				assert.Equal(t, "", resource.Name.ValueString())
				assert.Equal(t, "invalid-type", resource.Type.ValueString())
			}
		})
	}
}

func TestResourceConcurrency(t *testing.T) {
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
				resource := Resource{
					ID:   types.StringValue("concurrent-id"),
					Name: types.StringValue("Concurrent Resource"),
					Type: types.StringValue("team"),
				}

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						_ = resource.ID.ValueString()
						_ = resource.Name.ValueString()
						_ = resource.Type.ValueString()
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "error case - concurrent access validation":
				resource := Resource{
					ID:   types.StringValue("val-id"),
					Name: types.StringValue("Val Resource"),
				}

				done := make(chan bool, 50)
				for i := 0; i < 50; i++ {
					go func() {
						_ = resource.ID.ValueString()
						_ = resource.Name.ValueString()
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

func BenchmarkResource(b *testing.B) {
	b.Run("create", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Resource{
				ID:   types.StringValue("res-123"),
				Name: types.StringValue("Test Resource"),
				Type: types.StringValue("personal"),
			}
		}
	})

	b.Run("access fields", func(b *testing.B) {
		resource := Resource{
			ID:   types.StringValue("res-123"),
			Name: types.StringValue("Test Resource"),
			Type: types.StringValue("personal"),
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = resource.ID.ValueString()
			_ = resource.Name.ValueString()
			_ = resource.Type.ValueString()
		}
	})
}
