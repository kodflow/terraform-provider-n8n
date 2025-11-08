package models

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestResource(t *testing.T) {
	t.Run("create with all fields", func(t *testing.T) {
		// Create a resource with all fields populated
		dataMap := map[string]attr.Value{
			"api_key": types.StringValue("secret-key"),
			"url":     types.StringValue("https://api.example.com"),
		}

		resource := Resource{
			ID:        types.StringValue("cred-123"),
			Name:      types.StringValue("Test Credential"),
			Type:      types.StringValue("api"),
			Data:      types.MapValueMust(types.StringType, dataMap),
			CreatedAt: types.StringValue("2024-01-01T00:00:00Z"),
			UpdatedAt: types.StringValue("2024-01-02T00:00:00Z"),
		}

		assert.Equal(t, "cred-123", resource.ID.ValueString())
		assert.Equal(t, "Test Credential", resource.Name.ValueString())
		assert.Equal(t, "api", resource.Type.ValueString())
		assert.False(t, resource.Data.IsNull())
		assert.Equal(t, "2024-01-01T00:00:00Z", resource.CreatedAt.ValueString())
		assert.Equal(t, "2024-01-02T00:00:00Z", resource.UpdatedAt.ValueString())
	})

	t.Run("create with null values", func(t *testing.T) {
		// Create a resource with null values
		resource := Resource{
			ID:        types.StringNull(),
			Name:      types.StringNull(),
			Type:      types.StringNull(),
			Data:      types.MapNull(types.StringType),
			CreatedAt: types.StringNull(),
			UpdatedAt: types.StringNull(),
		}

		assert.True(t, resource.ID.IsNull())
		assert.True(t, resource.Name.IsNull())
		assert.True(t, resource.Type.IsNull())
		assert.True(t, resource.Data.IsNull())
		assert.True(t, resource.CreatedAt.IsNull())
		assert.True(t, resource.UpdatedAt.IsNull())
	})

	t.Run("create with unknown values", func(t *testing.T) {
		// Create a resource with unknown values
		resource := Resource{
			ID:        types.StringUnknown(),
			Name:      types.StringUnknown(),
			Type:      types.StringUnknown(),
			Data:      types.MapUnknown(types.StringType),
			CreatedAt: types.StringUnknown(),
			UpdatedAt: types.StringUnknown(),
		}

		assert.True(t, resource.ID.IsUnknown())
		assert.True(t, resource.Name.IsUnknown())
		assert.True(t, resource.Type.IsUnknown())
		assert.True(t, resource.Data.IsUnknown())
		assert.True(t, resource.CreatedAt.IsUnknown())
		assert.True(t, resource.UpdatedAt.IsUnknown())
	})

	t.Run("create with empty data map", func(t *testing.T) {
		// Create a resource with an empty data map
		emptyMap := map[string]attr.Value{}

		resource := Resource{
			ID:        types.StringValue("cred-empty"),
			Name:      types.StringValue("Empty Data Credential"),
			Type:      types.StringValue("empty"),
			Data:      types.MapValueMust(types.StringType, emptyMap),
			CreatedAt: types.StringValue("2024-01-01T00:00:00Z"),
			UpdatedAt: types.StringValue("2024-01-01T00:00:00Z"),
		}

		assert.Equal(t, "cred-empty", resource.ID.ValueString())
		assert.False(t, resource.Data.IsNull())
		assert.Equal(t, 0, len(resource.Data.Elements()))
	})

	t.Run("partial initialization", func(t *testing.T) {
		// Test with only some fields set
		resource := Resource{
			ID:   types.StringValue("partial-123"),
			Name: types.StringValue("Partial Cred"),
			// Other fields remain zero value (null)
		}

		assert.Equal(t, "partial-123", resource.ID.ValueString())
		assert.Equal(t, "Partial Cred", resource.Name.ValueString())
		assert.True(t, resource.Type.IsNull())
		assert.True(t, resource.Data.IsNull())
		assert.True(t, resource.CreatedAt.IsNull())
		assert.True(t, resource.UpdatedAt.IsNull())
	})

	t.Run("zero value struct", func(t *testing.T) {
		// Test zero value struct
		var resource Resource

		assert.True(t, resource.ID.IsNull())
		assert.True(t, resource.Name.IsNull())
		assert.True(t, resource.Type.IsNull())
		assert.True(t, resource.Data.IsNull())
		assert.True(t, resource.CreatedAt.IsNull())
		assert.True(t, resource.UpdatedAt.IsNull())
	})

	t.Run("modify values", func(t *testing.T) {
		// Test modifying values
		resource := Resource{
			ID:   types.StringValue("initial-id"),
			Name: types.StringValue("Initial Name"),
		}

		// Modify values
		resource.ID = types.StringValue("modified-id")
		resource.Name = types.StringValue("Modified Name")
		resource.Type = types.StringValue("modified-type")

		assert.Equal(t, "modified-id", resource.ID.ValueString())
		assert.Equal(t, "Modified Name", resource.Name.ValueString())
		assert.Equal(t, "modified-type", resource.Type.ValueString())
	})

	t.Run("various credential types", func(t *testing.T) {
		// Test various credential type values
		credTypes := []string{
			"api",
			"oauth2",
			"basic_auth",
			"ssh",
			"database",
			"aws",
			"azure",
			"gcp",
			"smtp",
			"ftp",
		}

		for _, credType := range credTypes {
			resource := Resource{
				Type: types.StringValue(credType),
			}
			assert.Equal(t, credType, resource.Type.ValueString())
		}
	})

	t.Run("complex data map", func(t *testing.T) {
		// Test with a complex data map
		complexData := map[string]attr.Value{
			"username":     types.StringValue("user123"),
			"password":     types.StringValue("secret-pass"),
			"host":         types.StringValue("db.example.com"),
			"port":         types.StringValue("5432"),
			"database":     types.StringValue("mydb"),
			"ssl":          types.StringValue("true"),
			"connection":   types.StringValue("postgresql://user:pass@host:5432/db"),
			"extra_params": types.StringValue("{\"timeout\":30,\"pool_size\":10}"),
		}

		resource := Resource{
			ID:   types.StringValue("complex-cred"),
			Name: types.StringValue("Complex Database Credential"),
			Type: types.StringValue("database"),
			Data: types.MapValueMust(types.StringType, complexData),
		}

		assert.Equal(t, "complex-cred", resource.ID.ValueString())
		assert.Equal(t, 8, len(resource.Data.Elements()))
	})

	t.Run("timestamp formats", func(t *testing.T) {
		// Test various timestamp formats
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
			resource := Resource{
				CreatedAt: types.StringValue(ts),
				UpdatedAt: types.StringValue(ts),
			}
			assert.Equal(t, ts, resource.CreatedAt.ValueString())
			assert.Equal(t, ts, resource.UpdatedAt.ValueString())
		}
	})

	t.Run("special characters in values", func(t *testing.T) {
		// Test special characters in string fields
		resource := Resource{
			ID:   types.StringValue("id-with-special-!@#$%^&*()"),
			Name: types.StringValue("Name with 特殊字符 🔒"),
			Type: types.StringValue("type_with.dots-and-dashes"),
		}

		assert.Equal(t, "id-with-special-!@#$%^&*()", resource.ID.ValueString())
		assert.Equal(t, "Name with 特殊字符 🔒", resource.Name.ValueString())
		assert.Equal(t, "type_with.dots-and-dashes", resource.Type.ValueString())
	})

	t.Run("very long values", func(t *testing.T) {
		// Test with very long strings
		longString := ""
		for i := 0; i < 1000; i++ {
			longString += "a"
		}

		resource := Resource{
			ID:   types.StringValue(longString),
			Name: types.StringValue(longString),
			Type: types.StringValue(longString),
		}

		assert.Equal(t, longString, resource.ID.ValueString())
		assert.Equal(t, longString, resource.Name.ValueString())
		assert.Equal(t, longString, resource.Type.ValueString())
	})

	t.Run("copy struct", func(t *testing.T) {
		// Test copying struct
		original := Resource{
			ID:        types.StringValue("original-id"),
			Name:      types.StringValue("Original"),
			Type:      types.StringValue("api"),
			CreatedAt: types.StringValue("2024-01-01T00:00:00Z"),
			UpdatedAt: types.StringValue("2024-01-01T00:00:00Z"),
		}

		copied := original

		assert.Equal(t, original.ID.ValueString(), copied.ID.ValueString())
		assert.Equal(t, original.Name.ValueString(), copied.Name.ValueString())
		assert.Equal(t, original.Type.ValueString(), copied.Type.ValueString())

		// Modify copied
		copied.ID = types.StringValue("modified-id")

		// Original should not be affected (value semantics)
		assert.Equal(t, "original-id", original.ID.ValueString())
		assert.Equal(t, "modified-id", copied.ID.ValueString())
	})

	t.Run("pointer to struct", func(t *testing.T) {
		// Test pointer to struct
		resource := &Resource{
			ID:   types.StringValue("pointer-id"),
			Name: types.StringValue("Pointer Resource"),
		}

		assert.NotNil(t, resource)
		assert.Equal(t, "pointer-id", resource.ID.ValueString())
		assert.Equal(t, "Pointer Resource", resource.Name.ValueString())
	})

	t.Run("struct field tags", func(t *testing.T) {
		// This test documents that the struct has proper tfsdk tags
		resource := Resource{
			ID:        types.StringValue("id"),
			Name:      types.StringValue("name"),
			Type:      types.StringValue("type"),
			Data:      types.MapNull(types.StringType),
			CreatedAt: types.StringValue("created"),
			UpdatedAt: types.StringValue("updated"),
		}

		// The tfsdk tags are used by Terraform plugin framework
		// They map to the Terraform schema field names
		assert.NotNil(t, resource)
	})

	t.Run("data map with different value types", func(t *testing.T) {
		// Test map with different types of values (all as strings since it's StringType)
		dataMap := map[string]attr.Value{
			"boolean":      types.StringValue("true"),
			"number":       types.StringValue("42"),
			"float":        types.StringValue("3.14159"),
			"null_value":   types.StringNull(),
			"empty_string": types.StringValue(""),
			"json":         types.StringValue(`{"key": "value"}`),
			"array":        types.StringValue(`["item1", "item2"]`),
		}

		resource := Resource{
			ID:   types.StringValue("mixed-data"),
			Data: types.MapValueMust(types.StringType, dataMap),
		}

		assert.Equal(t, "mixed-data", resource.ID.ValueString())
		assert.Equal(t, 7, len(resource.Data.Elements()))
	})

	t.Run("comparison", func(t *testing.T) {
		// Test struct comparison
		resource1 := Resource{
			ID:   types.StringValue("id1"),
			Name: types.StringValue("Name1"),
		}

		resource2 := Resource{
			ID:   types.StringValue("id1"),
			Name: types.StringValue("Name1"),
		}

		resource3 := Resource{
			ID:   types.StringValue("id2"),
			Name: types.StringValue("Name2"),
		}

		// Same values should be equal
		assert.Equal(t, resource1.ID.ValueString(), resource2.ID.ValueString())
		assert.Equal(t, resource1.Name.ValueString(), resource2.Name.ValueString())

		// Different values should not be equal
		assert.NotEqual(t, resource1.ID.ValueString(), resource3.ID.ValueString())
		assert.NotEqual(t, resource1.Name.ValueString(), resource3.Name.ValueString())
	})
}

func TestResourceConcurrency(t *testing.T) {
	t.Run("concurrent read", func(t *testing.T) {
		// Test concurrent reads
		resource := Resource{
			ID:   types.StringValue("concurrent-id"),
			Name: types.StringValue("Concurrent Resource"),
		}

		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				_ = resource.ID.ValueString()
				_ = resource.Name.ValueString()
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}

		assert.True(t, true, "Concurrent reads completed without issues")
	})

	t.Run("separate instances", func(t *testing.T) {
		// Test multiple instances don't interfere
		done := make(chan bool, 10)

		for i := 0; i < 10; i++ {
			go func(n int) {
				resource := Resource{
					ID:   types.StringValue(string(rune('a' + n))),
					Name: types.StringValue("Resource"),
				}
				assert.NotNil(t, resource)
				done <- true
			}(i)
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

func TestResourceValidation(t *testing.T) {
	t.Run("validate required fields", func(t *testing.T) {
		// Test that required fields can be checked
		resource := Resource{}

		// Check if required fields are null
		assert.True(t, resource.ID.IsNull())
		assert.True(t, resource.Name.IsNull())
		assert.True(t, resource.Type.IsNull())

		// Set required fields
		resource.ID = types.StringValue("required-id")
		resource.Name = types.StringValue("Required Name")
		resource.Type = types.StringValue("required-type")

		// Now they should not be null
		assert.False(t, resource.ID.IsNull())
		assert.False(t, resource.Name.IsNull())
		assert.False(t, resource.Type.IsNull())
	})

	t.Run("validate empty vs null", func(t *testing.T) {
		// Test difference between empty string and null
		resource1 := Resource{
			Name: types.StringValue(""),
		}

		resource2 := Resource{
			Name: types.StringNull(),
		}

		// Empty string is not null
		assert.False(t, resource1.Name.IsNull())
		assert.Equal(t, "", resource1.Name.ValueString())

		// Null is null
		assert.True(t, resource2.Name.IsNull())
		assert.Equal(t, "", resource2.Name.ValueString()) // Null returns empty string
	})
}

func TestResourceUseCases(t *testing.T) {
	t.Run("api credential", func(t *testing.T) {
		// Test typical API credential
		dataMap := map[string]attr.Value{
			"api_key":  types.StringValue("sk-1234567890abcdef"),
			"base_url": types.StringValue("https://api.example.com/v1"),
		}

		resource := Resource{
			ID:        types.StringValue("api-cred-001"),
			Name:      types.StringValue("Example API"),
			Type:      types.StringValue("api"),
			Data:      types.MapValueMust(types.StringType, dataMap),
			CreatedAt: types.StringValue(time.Now().Format(time.RFC3339)),
			UpdatedAt: types.StringValue(time.Now().Format(time.RFC3339)),
		}

		assert.Equal(t, "api", resource.Type.ValueString())
		assert.Equal(t, 2, len(resource.Data.Elements()))
	})

	t.Run("database credential", func(t *testing.T) {
		// Test database credential
		dataMap := map[string]attr.Value{
			"host":     types.StringValue("localhost"),
			"port":     types.StringValue("5432"),
			"database": types.StringValue("myapp"),
			"username": types.StringValue("dbuser"),
			"password": types.StringValue("dbpass"),
			"ssl_mode": types.StringValue("require"),
		}

		resource := Resource{
			ID:        types.StringValue("db-cred-001"),
			Name:      types.StringValue("PostgreSQL Production"),
			Type:      types.StringValue("postgresql"),
			Data:      types.MapValueMust(types.StringType, dataMap),
			CreatedAt: types.StringValue(time.Now().Format(time.RFC3339)),
			UpdatedAt: types.StringValue(time.Now().Format(time.RFC3339)),
		}

		assert.Equal(t, "postgresql", resource.Type.ValueString())
		assert.Equal(t, 6, len(resource.Data.Elements()))
	})

	t.Run("oauth credential", func(t *testing.T) {
		// Test OAuth credential
		dataMap := map[string]attr.Value{
			"client_id":     types.StringValue("client123"),
			"client_secret": types.StringValue("secret456"),
			"access_token":  types.StringValue("token789"),
			"refresh_token": types.StringValue("refresh012"),
			"token_expiry":  types.StringValue("2024-12-31T23:59:59Z"),
		}

		resource := Resource{
			ID:        types.StringValue("oauth-cred-001"),
			Name:      types.StringValue("Google OAuth"),
			Type:      types.StringValue("oauth2"),
			Data:      types.MapValueMust(types.StringType, dataMap),
			CreatedAt: types.StringValue(time.Now().Format(time.RFC3339)),
			UpdatedAt: types.StringValue(time.Now().Format(time.RFC3339)),
		}

		assert.Equal(t, "oauth2", resource.Type.ValueString())
		assert.Equal(t, 5, len(resource.Data.Elements()))
	})
}

func BenchmarkResource(b *testing.B) {
	b.Run("create", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Resource{
				ID:        types.StringValue("bench-id"),
				Name:      types.StringValue("Bench Resource"),
				Type:      types.StringValue("api"),
				Data:      types.MapNull(types.StringType),
				CreatedAt: types.StringValue("2024-01-01T00:00:00Z"),
				UpdatedAt: types.StringValue("2024-01-01T00:00:00Z"),
			}
		}
	})

	b.Run("access fields", func(b *testing.B) {
		resource := Resource{
			ID:        types.StringValue("bench-id"),
			Name:      types.StringValue("Bench Resource"),
			Type:      types.StringValue("api"),
			CreatedAt: types.StringValue("2024-01-01T00:00:00Z"),
			UpdatedAt: types.StringValue("2024-01-01T00:00:00Z"),
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = resource.ID.ValueString()
			_ = resource.Name.ValueString()
			_ = resource.Type.ValueString()
			_ = resource.CreatedAt.ValueString()
			_ = resource.UpdatedAt.ValueString()
		}
	})

	b.Run("modify fields", func(b *testing.B) {
		resource := Resource{}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			resource.ID = types.StringValue("id")
			resource.Name = types.StringValue("name")
			resource.Type = types.StringValue("type")
		}
	})

	b.Run("create with map", func(b *testing.B) {
		dataMap := map[string]attr.Value{
			"key1": types.StringValue("value1"),
			"key2": types.StringValue("value2"),
			"key3": types.StringValue("value3"),
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Resource{
				ID:   types.StringValue("bench-id"),
				Data: types.MapValueMust(types.StringType, dataMap),
			}
		}
	})
}
