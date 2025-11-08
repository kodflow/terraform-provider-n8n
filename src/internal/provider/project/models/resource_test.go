package models

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestResource(t *testing.T) {
	t.Run("create with all fields", func(t *testing.T) {
		resource := Resource{
			ID:   types.StringValue("res-123"),
			Name: types.StringValue("Test Resource"),
			Type: types.StringValue("personal"),
		}

		assert.Equal(t, "res-123", resource.ID.ValueString())
		assert.Equal(t, "Test Resource", resource.Name.ValueString())
		assert.Equal(t, "personal", resource.Type.ValueString())
	})

	t.Run("create with null values", func(t *testing.T) {
		resource := Resource{
			ID:   types.StringNull(),
			Name: types.StringNull(),
			Type: types.StringNull(),
		}

		assert.True(t, resource.ID.IsNull())
		assert.True(t, resource.Name.IsNull())
		assert.True(t, resource.Type.IsNull())
	})

	t.Run("create with unknown values", func(t *testing.T) {
		resource := Resource{
			ID:   types.StringUnknown(),
			Name: types.StringUnknown(),
			Type: types.StringUnknown(),
		}

		assert.True(t, resource.ID.IsUnknown())
		assert.True(t, resource.Name.IsUnknown())
		assert.True(t, resource.Type.IsUnknown())
	})

	t.Run("resource types", func(t *testing.T) {
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
	})

	t.Run("zero value struct", func(t *testing.T) {
		var resource Resource
		assert.True(t, resource.ID.IsNull())
		assert.True(t, resource.Name.IsNull())
		assert.True(t, resource.Type.IsNull())
	})

	t.Run("copy struct", func(t *testing.T) {
		original := Resource{
			ID:   types.StringValue("original-id"),
			Name: types.StringValue("Original Resource"),
			Type: types.StringValue("team"),
		}

		copy := original

		assert.Equal(t, original.ID.ValueString(), copy.ID.ValueString())
		assert.Equal(t, original.Name.ValueString(), copy.Name.ValueString())
		assert.Equal(t, original.Type.ValueString(), copy.Type.ValueString())

		// Modify copy
		copy.ID = types.StringValue("modified-id")
		copy.Name = types.StringValue("Modified Resource")
		assert.Equal(t, "original-id", original.ID.ValueString())
		assert.Equal(t, "modified-id", copy.ID.ValueString())
		assert.Equal(t, "Original Resource", original.Name.ValueString())
		assert.Equal(t, "Modified Resource", copy.Name.ValueString())
	})

	t.Run("partial initialization", func(t *testing.T) {
		resource := Resource{
			ID: types.StringValue("res-partial"),
			// Other fields remain null
		}

		assert.Equal(t, "res-partial", resource.ID.ValueString())
		assert.True(t, resource.Name.IsNull())
		assert.True(t, resource.Type.IsNull())
	})

	t.Run("name variations", func(t *testing.T) {
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
	})

	t.Run("id formats", func(t *testing.T) {
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
	})

	t.Run("required fields validation", func(t *testing.T) {
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
	})
}

func TestResourceConcurrency(t *testing.T) {
	t.Run("concurrent read", func(t *testing.T) {
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
	})
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
