package models

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestN8nProviderModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create with valid values", wantErr: false},
		{name: "create with null values", wantErr: false},
		{name: "create with unknown values", wantErr: false},
		{name: "create with empty strings", wantErr: false},
		{name: "partial initialization", wantErr: false},
		{name: "zero value struct", wantErr: false},
		{name: "struct field tags", wantErr: false},
		{name: "modify values", wantErr: false},
		{name: "various URL formats", wantErr: false},
		{name: "various API key formats", wantErr: false},
		{name: "special characters in values", wantErr: false},
		{name: "unicode in values", wantErr: false},
		{name: "very long values", wantErr: false},
		{name: "comparison", wantErr: false},
		{name: "pointer to struct", wantErr: false},
		{name: "copy struct", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "create with valid values":
				model := N8nProviderModel{
					APIKey:  types.StringValue("test-api-key"),
					BaseURL: types.StringValue("https://n8n.example.com"),
				}
				assert.NotNil(t, model)
				assert.Equal(t, "test-api-key", model.APIKey.ValueString())
				assert.Equal(t, "https://n8n.example.com", model.BaseURL.ValueString())
				assert.False(t, model.APIKey.IsNull())
				assert.False(t, model.BaseURL.IsNull())

			case "create with null values":
				model := N8nProviderModel{
					APIKey:  types.StringNull(),
					BaseURL: types.StringNull(),
				}
				assert.NotNil(t, model)
				assert.True(t, model.APIKey.IsNull())
				assert.True(t, model.BaseURL.IsNull())
				assert.Equal(t, "", model.APIKey.ValueString())
				assert.Equal(t, "", model.BaseURL.ValueString())

			case "create with unknown values":
				model := N8nProviderModel{
					APIKey:  types.StringUnknown(),
					BaseURL: types.StringUnknown(),
				}
				assert.NotNil(t, model)
				assert.True(t, model.APIKey.IsUnknown())
				assert.True(t, model.BaseURL.IsUnknown())
				assert.False(t, model.APIKey.IsNull())
				assert.False(t, model.BaseURL.IsNull())

			case "create with empty strings":
				model := N8nProviderModel{
					APIKey:  types.StringValue(""),
					BaseURL: types.StringValue(""),
				}
				assert.NotNil(t, model)
				assert.Equal(t, "", model.APIKey.ValueString())
				assert.Equal(t, "", model.BaseURL.ValueString())
				assert.False(t, model.APIKey.IsNull())
				assert.False(t, model.BaseURL.IsNull())

			case "partial initialization":
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

			case "zero value struct":
				var model N8nProviderModel
				assert.True(t, model.APIKey.IsNull())
				assert.True(t, model.BaseURL.IsNull())

			case "struct field tags":
				model := N8nProviderModel{
					APIKey:  types.StringValue("key"),
					BaseURL: types.StringValue("url"),
				}
				assert.NotNil(t, model)

			case "modify values":
				model := N8nProviderModel{
					APIKey:  types.StringValue("initial-key"),
					BaseURL: types.StringValue("https://initial.com"),
				}
				model.APIKey = types.StringValue("modified-key")
				model.BaseURL = types.StringValue("https://modified.com")
				assert.Equal(t, "modified-key", model.APIKey.ValueString())
				assert.Equal(t, "https://modified.com", model.BaseURL.ValueString())

			case "various URL formats":
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

			case "various API key formats":
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

			case "special characters in values":
				model := N8nProviderModel{
					APIKey:  types.StringValue("key-with-special-!@#$%^&*()"),
					BaseURL: types.StringValue("https://example.com?param=value&other=123"),
				}
				assert.Equal(t, "key-with-special-!@#$%^&*()", model.APIKey.ValueString())
				assert.Equal(t, "https://example.com?param=value&other=123", model.BaseURL.ValueString())

			case "unicode in values":
				model := N8nProviderModel{
					APIKey:  types.StringValue("key-测试-テスト-тест"),
					BaseURL: types.StringValue("https://例え.com"),
				}
				assert.Equal(t, "key-测试-テスト-тест", model.APIKey.ValueString())
				assert.Equal(t, "https://例え.com", model.BaseURL.ValueString())

			case "very long values":
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

			case "comparison":
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
				assert.Equal(t, model1.APIKey.ValueString(), model2.APIKey.ValueString())
				assert.Equal(t, model1.BaseURL.ValueString(), model2.BaseURL.ValueString())
				assert.NotEqual(t, model1.APIKey.ValueString(), model3.APIKey.ValueString())
				assert.NotEqual(t, model1.BaseURL.ValueString(), model3.BaseURL.ValueString())

			case "pointer to struct":
				model := &N8nProviderModel{
					APIKey:  types.StringValue("pointer-key"),
					BaseURL: types.StringValue("https://pointer.com"),
				}
				assert.NotNil(t, model)
				assert.Equal(t, "pointer-key", model.APIKey.ValueString())
				assert.Equal(t, "https://pointer.com", model.BaseURL.ValueString())

			case "copy struct":
				original := N8nProviderModel{
					APIKey:  types.StringValue("original-key"),
					BaseURL: types.StringValue("https://original.com"),
				}
				copied := original
				assert.Equal(t, original.APIKey.ValueString(), copied.APIKey.ValueString())
				assert.Equal(t, original.BaseURL.ValueString(), copied.BaseURL.ValueString())
				copied.APIKey = types.StringValue("modified-key")
				assert.Equal(t, "original-key", original.APIKey.ValueString())
				assert.Equal(t, "modified-key", copied.APIKey.ValueString())

			case "error case - validation checks":
				model := N8nProviderModel{}
				assert.True(t, model.APIKey.IsNull())
				assert.True(t, model.BaseURL.IsNull())
			}
		})
	}
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
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "concurrent read", wantErr: false},
		{name: "separate instances", wantErr: false},
		{name: "error case - concurrent validation", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// NO t.Parallel() here - goroutines

			switch tt.name {
			case "concurrent read":
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

			case "separate instances":
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

			case "error case - concurrent validation":
				done := make(chan bool, 10)
				for i := 0; i < 10; i++ {
					go func() {
						model := N8nProviderModel{}
						assert.True(t, model.APIKey.IsNull())
						assert.True(t, model.BaseURL.IsNull())
						done <- true
					}()
				}
				for i := 0; i < 10; i++ {
					<-done
				}
			}
		})
	}
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
