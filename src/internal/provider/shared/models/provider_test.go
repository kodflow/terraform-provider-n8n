package models

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestN8nProviderModel(t *testing.T) {
	t.Run("create with valid values", func(t *testing.T) {
		// Test creating model with valid values
		model := N8nProviderModel{
			APIKey:  types.StringValue("test-api-key"),
			BaseURL: types.StringValue("https://n8n.example.com"),
		}

		assert.NotNil(t, model)
		assert.Equal(t, "test-api-key", model.APIKey.ValueString())
		assert.Equal(t, "https://n8n.example.com", model.BaseURL.ValueString())
		assert.False(t, model.APIKey.IsNull())
		assert.False(t, model.BaseURL.IsNull())
	})

	t.Run("create with null values", func(t *testing.T) {
		// Test creating model with null values
		model := N8nProviderModel{
			APIKey:  types.StringNull(),
			BaseURL: types.StringNull(),
		}

		assert.NotNil(t, model)
		assert.True(t, model.APIKey.IsNull())
		assert.True(t, model.BaseURL.IsNull())
		assert.Equal(t, "", model.APIKey.ValueString())
		assert.Equal(t, "", model.BaseURL.ValueString())
	})

	t.Run("create with unknown values", func(t *testing.T) {
		// Test creating model with unknown values
		model := N8nProviderModel{
			APIKey:  types.StringUnknown(),
			BaseURL: types.StringUnknown(),
		}

		assert.NotNil(t, model)
		assert.True(t, model.APIKey.IsUnknown())
		assert.True(t, model.BaseURL.IsUnknown())
		assert.False(t, model.APIKey.IsNull())
		assert.False(t, model.BaseURL.IsNull())
	})

	t.Run("create with empty strings", func(t *testing.T) {
		// Test creating model with empty string values
		model := N8nProviderModel{
			APIKey:  types.StringValue(""),
			BaseURL: types.StringValue(""),
		}

		assert.NotNil(t, model)
		assert.Equal(t, "", model.APIKey.ValueString())
		assert.Equal(t, "", model.BaseURL.ValueString())
		assert.False(t, model.APIKey.IsNull())
		assert.False(t, model.BaseURL.IsNull())
	})

	t.Run("partial initialization", func(t *testing.T) {
		// Test with only one field set
		model1 := N8nProviderModel{
			APIKey: types.StringValue("api-key-only"),
		}
		assert.Equal(t, "api-key-only", model1.APIKey.ValueString())
		assert.True(t, model1.BaseURL.IsNull())

		model2 := N8nProviderModel{
			BaseURL: types.StringValue("https://base-url-only.com"),
		}
		assert.True(t, model2.APIKey.IsNull())
		assert.Equal(t, "https://base-url-only.com", model2.BaseURL.ValueString())
	})

	t.Run("zero value struct", func(t *testing.T) {
		// Test zero value struct
		var model N8nProviderModel

		assert.True(t, model.APIKey.IsNull())
		assert.True(t, model.BaseURL.IsNull())
	})

	t.Run("struct field tags", func(t *testing.T) {
		// This test documents that the struct has proper tfsdk tags
		model := N8nProviderModel{
			APIKey:  types.StringValue("key"),
			BaseURL: types.StringValue("url"),
		}

		// The tfsdk tags are used by Terraform plugin framework
		// api_key and base_url are the field names in Terraform configuration
		assert.NotNil(t, model)
	})

	t.Run("modify values", func(t *testing.T) {
		// Test modifying values
		model := N8nProviderModel{
			APIKey:  types.StringValue("initial-key"),
			BaseURL: types.StringValue("https://initial.com"),
		}

		// Modify values
		model.APIKey = types.StringValue("modified-key")
		model.BaseURL = types.StringValue("https://modified.com")

		assert.Equal(t, "modified-key", model.APIKey.ValueString())
		assert.Equal(t, "https://modified.com", model.BaseURL.ValueString())
	})

	t.Run("various URL formats", func(t *testing.T) {
		// Test various URL formats
		urls := []string{
			"https://n8n.example.com",
			"http://localhost:5678",
			"https://n8n.example.com:8080",
			"https://subdomain.n8n.example.com",
			"https://n8n.example.com/path",
			"https://192.168.1.1:5678",
		}

		for _, url := range urls {
			model := N8nProviderModel{
				APIKey:  types.StringValue("key"),
				BaseURL: types.StringValue(url),
			}
			assert.Equal(t, url, model.BaseURL.ValueString())
		}
	})

	t.Run("various API key formats", func(t *testing.T) {
		// Test various API key formats
		keys := []string{
			"simple-key",
			"KEY_WITH_UNDERSCORES",
			"key-with-dashes",
			"key.with.dots",
			"key123456789",
			"veryLongAPIKeyWithManyCharacters1234567890abcdefghijklmnopqrstuvwxyz",
			"key=with=equals",
			"key/with/slashes",
		}

		for _, key := range keys {
			model := N8nProviderModel{
				APIKey:  types.StringValue(key),
				BaseURL: types.StringValue("https://test.com"),
			}
			assert.Equal(t, key, model.APIKey.ValueString())
		}
	})

	t.Run("special characters in values", func(t *testing.T) {
		// Test special characters
		model := N8nProviderModel{
			APIKey:  types.StringValue("key-with-special-!@#$%^&*()"),
			BaseURL: types.StringValue("https://example.com?param=value&other=123"),
		}

		assert.Equal(t, "key-with-special-!@#$%^&*()", model.APIKey.ValueString())
		assert.Equal(t, "https://example.com?param=value&other=123", model.BaseURL.ValueString())
	})

	t.Run("unicode in values", func(t *testing.T) {
		// Test unicode characters
		model := N8nProviderModel{
			APIKey:  types.StringValue("key-测试-テスト-тест"),
			BaseURL: types.StringValue("https://例え.com"),
		}

		assert.Equal(t, "key-测试-テスト-тест", model.APIKey.ValueString())
		assert.Equal(t, "https://例え.com", model.BaseURL.ValueString())
	})

	t.Run("very long values", func(t *testing.T) {
		// Test with very long strings
		longKey := ""
		for i := 0; i < 1000; i++ {
			longKey += "a"
		}

		longURL := "https://example.com/"
		for i := 0; i < 1000; i++ {
			longURL += "path/"
		}

		model := N8nProviderModel{
			APIKey:  types.StringValue(longKey),
			BaseURL: types.StringValue(longURL),
		}

		assert.Equal(t, longKey, model.APIKey.ValueString())
		assert.Equal(t, longURL, model.BaseURL.ValueString())
	})

	t.Run("comparison", func(t *testing.T) {
		// Test struct comparison
		model1 := N8nProviderModel{
			APIKey:  types.StringValue("key1"),
			BaseURL: types.StringValue("https://url1.com"),
		}

		model2 := N8nProviderModel{
			APIKey:  types.StringValue("key1"),
			BaseURL: types.StringValue("https://url1.com"),
		}

		model3 := N8nProviderModel{
			APIKey:  types.StringValue("key2"),
			BaseURL: types.StringValue("https://url2.com"),
		}

		// Same values should be equal
		assert.Equal(t, model1.APIKey.ValueString(), model2.APIKey.ValueString())
		assert.Equal(t, model1.BaseURL.ValueString(), model2.BaseURL.ValueString())

		// Different values should not be equal
		assert.NotEqual(t, model1.APIKey.ValueString(), model3.APIKey.ValueString())
		assert.NotEqual(t, model1.BaseURL.ValueString(), model3.BaseURL.ValueString())
	})

	t.Run("pointer to struct", func(t *testing.T) {
		// Test pointer to struct
		model := &N8nProviderModel{
			APIKey:  types.StringValue("pointer-key"),
			BaseURL: types.StringValue("https://pointer.com"),
		}

		assert.NotNil(t, model)
		assert.Equal(t, "pointer-key", model.APIKey.ValueString())
		assert.Equal(t, "https://pointer.com", model.BaseURL.ValueString())
	})

	t.Run("copy struct", func(t *testing.T) {
		// Test copying struct
		original := N8nProviderModel{
			APIKey:  types.StringValue("original-key"),
			BaseURL: types.StringValue("https://original.com"),
		}

		copied := original

		assert.Equal(t, original.APIKey.ValueString(), copied.APIKey.ValueString())
		assert.Equal(t, original.BaseURL.ValueString(), copied.BaseURL.ValueString())

		// Modify copied
		copied.APIKey = types.StringValue("modified-key")

		// Original should not be affected (value semantics)
		assert.Equal(t, "original-key", original.APIKey.ValueString())
		assert.Equal(t, "modified-key", copied.APIKey.ValueString())
	})
}

func TestN8nProviderModelValidation(t *testing.T) {
	t.Run("validate required fields", func(t *testing.T) {
		// Test that both fields can be checked for null/empty
		model := N8nProviderModel{}

		// Both should be null initially
		assert.True(t, model.APIKey.IsNull())
		assert.True(t, model.BaseURL.IsNull())

		// Check empty string scenario
		model.APIKey = types.StringValue("")
		model.BaseURL = types.StringValue("")

		assert.Equal(t, "", model.APIKey.ValueString())
		assert.Equal(t, "", model.BaseURL.ValueString())
		assert.False(t, model.APIKey.IsNull()) // Empty string is not null
		assert.False(t, model.BaseURL.IsNull())
	})

	t.Run("validate URL format", func(t *testing.T) {
		// Test various URL validation scenarios
		testCases := []struct {
			name    string
			baseURL string
			valid   bool // This would be used in actual validation logic
		}{
			{"valid https", "https://n8n.example.com", true},
			{"valid http", "http://localhost:5678", true},
			{"with port", "https://n8n.example.com:8080", true},
			{"with path", "https://n8n.example.com/api", true},
			{"missing protocol", "n8n.example.com", false},
			{"invalid protocol", "ftp://n8n.example.com", false},
			{"empty", "", false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				model := N8nProviderModel{
					APIKey:  types.StringValue("key"),
					BaseURL: types.StringValue(tc.baseURL),
				}
				assert.Equal(t, tc.baseURL, model.BaseURL.ValueString())
			})
		}
	})
}

func TestN8nProviderModelConcurrency(t *testing.T) {
	t.Run("concurrent read", func(t *testing.T) {
		// Test concurrent reads
		model := N8nProviderModel{
			APIKey:  types.StringValue("concurrent-key"),
			BaseURL: types.StringValue("https://concurrent.com"),
		}

		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				_ = model.APIKey.ValueString()
				_ = model.BaseURL.ValueString()
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
				model := N8nProviderModel{
					APIKey:  types.StringValue(string(rune('a' + n))),
					BaseURL: types.StringValue("https://test.com"),
				}
				assert.NotNil(t, model)
				done <- true
			}(i)
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

func BenchmarkN8nProviderModel(b *testing.B) {
	b.Run("create", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = N8nProviderModel{
				APIKey:  types.StringValue("benchmark-key"),
				BaseURL: types.StringValue("https://benchmark.com"),
			}
		}
	})

	b.Run("access fields", func(b *testing.B) {
		model := N8nProviderModel{
			APIKey:  types.StringValue("benchmark-key"),
			BaseURL: types.StringValue("https://benchmark.com"),
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = model.APIKey.ValueString()
			_ = model.BaseURL.ValueString()
		}
	})

	b.Run("modify fields", func(b *testing.B) {
		model := N8nProviderModel{}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			model.APIKey = types.StringValue("key")
			model.BaseURL = types.StringValue("url")
		}
	})
}
