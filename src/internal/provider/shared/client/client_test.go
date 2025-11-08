package client

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewN8nClient(t *testing.T) {
	t.Run("creates client with base URL and API key", func(t *testing.T) {
		baseURL := "https://n8n.example.com"
		apiKey := "test-api-key-123"

		client := NewN8nClient(baseURL, apiKey)

		assert.NotNil(t, client, "Client should not be nil")
		assert.Equal(t, baseURL, client.BaseURL, "BaseURL should be set")
		assert.Equal(t, apiKey, client.APIKey, "APIKey should be set")
		assert.NotNil(t, client.APIClient, "APIClient should be initialized")
	})

	t.Run("creates client with empty base URL", func(t *testing.T) {
		client := NewN8nClient("", "test-key")

		assert.NotNil(t, client, "Client should be created")
		assert.Equal(t, "", client.BaseURL, "Should handle empty base URL")
		assert.NotNil(t, client.APIClient, "APIClient should be initialized")
	})

	t.Run("creates client with empty API key", func(t *testing.T) {
		client := NewN8nClient("https://n8n.example.com", "")

		assert.NotNil(t, client, "Client should be created")
		assert.Equal(t, "", client.APIKey, "Should handle empty API key")
		assert.NotNil(t, client.APIClient, "APIClient should be initialized")
	})

	t.Run("creates client with both empty values", func(t *testing.T) {
		client := NewN8nClient("", "")

		assert.NotNil(t, client, "Client should be created")
		assert.Equal(t, "", client.BaseURL, "BaseURL should be empty")
		assert.Equal(t, "", client.APIKey, "APIKey should be empty")
		assert.NotNil(t, client.APIClient, "APIClient should be initialized")
	})

	t.Run("API client configuration includes base URL", func(t *testing.T) {
		baseURL := "https://n8n.example.com"
		client := NewN8nClient(baseURL, "test-key")

		require.NotNil(t, client.APIClient)
		require.NotNil(t, client.APIClient.GetConfig())
		require.NotEmpty(t, client.APIClient.GetConfig().Servers)

		// Verify the server URL is set correctly (should append /api/v1)
		serverURL := client.APIClient.GetConfig().Servers[0].URL
		assert.Equal(t, baseURL+"/api/v1", serverURL, "Server URL should include /api/v1 suffix")
	})

	t.Run("API client configuration includes API key header", func(t *testing.T) {
		apiKey := "test-api-key-123"
		client := NewN8nClient("https://n8n.example.com", apiKey)

		require.NotNil(t, client.APIClient)
		require.NotNil(t, client.APIClient.GetConfig())

		// Verify the API key header is set
		defaultHeaders := client.APIClient.GetConfig().DefaultHeader
		assert.Contains(t, defaultHeaders, "X-N8N-API-KEY", "Should have X-N8N-API-KEY header")
		assert.Equal(t, apiKey, defaultHeaders["X-N8N-API-KEY"], "API key should be in header")
	})

	t.Run("creates multiple independent clients", func(t *testing.T) {
		client1 := NewN8nClient("https://n8n1.example.com", "key1")
		client2 := NewN8nClient("https://n8n2.example.com", "key2")

		assert.NotEqual(t, client1.BaseURL, client2.BaseURL, "Clients should have different base URLs")
		assert.NotEqual(t, client1.APIKey, client2.APIKey, "Clients should have different API keys")
		assert.NotSame(t, client1, client2, "Should be different instances")
		assert.NotSame(t, client1.APIClient, client2.APIClient, "Should have different API clients")
	})

	t.Run("handles various base URL formats", func(t *testing.T) {
		testCases := []struct {
			name     string
			baseURL  string
			expected string
		}{
			{"with https", "https://n8n.example.com", "https://n8n.example.com/api/v1"},
			{"with http", "http://n8n.example.com", "http://n8n.example.com/api/v1"},
			{"with port", "https://n8n.example.com:8080", "https://n8n.example.com:8080/api/v1"},
			{"with trailing slash", "https://n8n.example.com/", "https://n8n.example.com//api/v1"},
			{"localhost", "http://localhost:5678", "http://localhost:5678/api/v1"},
			{"IP address", "http://192.168.1.1", "http://192.168.1.1/api/v1"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				client := NewN8nClient(tc.baseURL, "test-key")
				require.NotNil(t, client.APIClient)
				require.NotNil(t, client.APIClient.GetConfig())
				require.NotEmpty(t, client.APIClient.GetConfig().Servers)

				serverURL := client.APIClient.GetConfig().Servers[0].URL
				assert.Equal(t, tc.expected, serverURL, "Server URL should be formatted correctly")
			})
		}
	})

	t.Run("handles various API key formats", func(t *testing.T) {
		testCases := []struct {
			name   string
			apiKey string
		}{
			{"alphanumeric", "abc123xyz789"},
			{"with special chars", "key-with-dashes-and_underscores"},
			{"very long key", strings.Repeat("a", 500)},
			{"single char", "k"},
			{"UUID format", "550e8400-e29b-41d4-a716-446655440000"},
			{"base64-like", "dGVzdC1hcGkta2V5LTEyMw=="},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				client := NewN8nClient("https://n8n.example.com", tc.apiKey)
				assert.Equal(t, tc.apiKey, client.APIKey, "Should handle API key: %s", tc.name)

				defaultHeaders := client.APIClient.GetConfig().DefaultHeader
				assert.Equal(t, tc.apiKey, defaultHeaders["X-N8N-API-KEY"], "Header should contain API key")
			})
		}
	})

	t.Run("server description is set", func(t *testing.T) {
		client := NewN8nClient("https://n8n.example.com", "test-key")

		require.NotNil(t, client.APIClient)
		require.NotNil(t, client.APIClient.GetConfig())
		require.NotEmpty(t, client.APIClient.GetConfig().Servers)

		description := client.APIClient.GetConfig().Servers[0].Description
		assert.Equal(t, "n8n instance", description, "Server should have description")
	})

	t.Run("configuration is properly initialized", func(t *testing.T) {
		client := NewN8nClient("https://n8n.example.com", "test-key")

		require.NotNil(t, client.APIClient)
		cfg := client.APIClient.GetConfig()
		require.NotNil(t, cfg)

		// Verify servers are configured
		assert.NotEmpty(t, cfg.Servers, "Servers should be configured")
		assert.Len(t, cfg.Servers, 1, "Should have exactly one server")

		// Verify headers are configured
		assert.NotNil(t, cfg.DefaultHeader, "Default headers should be set")
	})
}

func TestN8nClientStructure(t *testing.T) {
	t.Run("client has all required fields", func(t *testing.T) {
		client := NewN8nClient("https://n8n.example.com", "test-key")

		assert.NotNil(t, client.APIClient, "APIClient field should be set")
		assert.NotEmpty(t, client.BaseURL, "BaseURL field should be set")
		assert.NotEmpty(t, client.APIKey, "APIKey field should be set")
	})

	t.Run("client fields are accessible", func(t *testing.T) {
		baseURL := "https://n8n.example.com"
		apiKey := "test-key-123"
		client := NewN8nClient(baseURL, apiKey)

		// Verify fields are exported and accessible
		assert.Equal(t, baseURL, client.BaseURL)
		assert.Equal(t, apiKey, client.APIKey)
		assert.NotNil(t, client.APIClient)
	})
}

func TestN8nClientEdgeCases(t *testing.T) {
	t.Run("handles nil after creation", func(t *testing.T) {
		client := NewN8nClient("https://n8n.example.com", "test-key")
		assert.NotNil(t, client, "Client should never be nil")
	})

	t.Run("handles very long base URL", func(t *testing.T) {
		longURL := "https://" + strings.Repeat("a", 1000) + ".com"
		client := NewN8nClient(longURL, "test-key")

		assert.NotNil(t, client)
		assert.Equal(t, longURL, client.BaseURL)
	})

	t.Run("handles special characters in base URL", func(t *testing.T) {
		// Note: This might not be a valid URL, but the function should handle it
		specialURL := "https://n8n-test.example.com:8080/path"
		client := NewN8nClient(specialURL, "test-key")

		assert.NotNil(t, client)
		assert.Equal(t, specialURL, client.BaseURL)
	})

	t.Run("handles whitespace in inputs", func(t *testing.T) {
		client := NewN8nClient("  https://n8n.example.com  ", "  test-key  ")

		assert.NotNil(t, client)
		// The function doesn't trim, so whitespace is preserved
		assert.Equal(t, "  https://n8n.example.com  ", client.BaseURL)
		assert.Equal(t, "  test-key  ", client.APIKey)
	})
}

func TestN8nClientServerConfiguration(t *testing.T) {
	t.Run("server URL appends /api/v1", func(t *testing.T) {
		baseURL := "https://n8n.example.com"
		client := NewN8nClient(baseURL, "test-key")

		serverURL := client.APIClient.GetConfig().Servers[0].URL
		assert.True(t, strings.HasSuffix(serverURL, "/api/v1"), "Server URL should end with /api/v1")
		assert.True(t, strings.HasPrefix(serverURL, baseURL), "Server URL should start with base URL")
	})

	t.Run("handles base URL with existing path", func(t *testing.T) {
		baseURL := "https://n8n.example.com/custom/path"
		client := NewN8nClient(baseURL, "test-key")

		serverURL := client.APIClient.GetConfig().Servers[0].URL
		assert.Equal(t, baseURL+"/api/v1", serverURL, "Should append /api/v1 to existing path")
	})

	t.Run("handles base URL with trailing slash", func(t *testing.T) {
		baseURL := "https://n8n.example.com/"
		client := NewN8nClient(baseURL, "test-key")

		serverURL := client.APIClient.GetConfig().Servers[0].URL
		// Will result in double slash, which is the actual behavior
		assert.Equal(t, "https://n8n.example.com//api/v1", serverURL)
	})
}

func TestN8nClientHeaderConfiguration(t *testing.T) {
	t.Run("X-N8N-API-KEY header is set", func(t *testing.T) {
		apiKey := "my-secret-key"
		client := NewN8nClient("https://n8n.example.com", apiKey)

		headers := client.APIClient.GetConfig().DefaultHeader
		assert.Contains(t, headers, "X-N8N-API-KEY")
		assert.Equal(t, apiKey, headers["X-N8N-API-KEY"])
	})

	t.Run("header name is correct case", func(t *testing.T) {
		client := NewN8nClient("https://n8n.example.com", "test-key")

		headers := client.APIClient.GetConfig().DefaultHeader
		// Verify exact header name (case-sensitive)
		_, exists := headers["X-N8N-API-KEY"]
		assert.True(t, exists, "Header should be exactly 'X-N8N-API-KEY'")
	})

	t.Run("empty API key is still set in header", func(t *testing.T) {
		client := NewN8nClient("https://n8n.example.com", "")

		headers := client.APIClient.GetConfig().DefaultHeader
		assert.Contains(t, headers, "X-N8N-API-KEY")
		assert.Equal(t, "", headers["X-N8N-API-KEY"])
	})
}

func TestN8nClientIsolation(t *testing.T) {
	t.Run("multiple clients are independent", func(t *testing.T) {
		client1 := NewN8nClient("https://n8n1.example.com", "key1")
		client2 := NewN8nClient("https://n8n2.example.com", "key2")

		// Modifying one client shouldn't affect the other
		client1.BaseURL = "https://modified.example.com"
		assert.NotEqual(t, client1.BaseURL, client2.BaseURL)
		assert.Equal(t, "https://n8n2.example.com", client2.BaseURL)
	})

	t.Run("API clients are separate instances", func(t *testing.T) {
		client1 := NewN8nClient("https://n8n1.example.com", "key1")
		client2 := NewN8nClient("https://n8n2.example.com", "key2")

		assert.NotSame(t, client1.APIClient, client2.APIClient)

		// Verify configurations are independent
		config1 := client1.APIClient.GetConfig()
		config2 := client2.APIClient.GetConfig()
		assert.NotSame(t, config1, config2)
	})
}

func TestN8nClientUseCases(t *testing.T) {
	t.Run("typical production use case", func(t *testing.T) {
		client := NewN8nClient(
			"https://n8n.mycompany.com",
			"prod-api-key-12345",
		)

		assert.NotNil(t, client)
		assert.Equal(t, "https://n8n.mycompany.com", client.BaseURL)
		assert.Equal(t, "prod-api-key-12345", client.APIKey)
		assert.NotNil(t, client.APIClient)
	})

	t.Run("local development use case", func(t *testing.T) {
		client := NewN8nClient(
			"http://localhost:5678",
			"dev-api-key",
		)

		assert.NotNil(t, client)
		assert.Equal(t, "http://localhost:5678", client.BaseURL)
		assert.Equal(t, "dev-api-key", client.APIKey)

		serverURL := client.APIClient.GetConfig().Servers[0].URL
		assert.Equal(t, "http://localhost:5678/api/v1", serverURL)
	})

	t.Run("custom port use case", func(t *testing.T) {
		client := NewN8nClient(
			"https://n8n.example.com:8443",
			"api-key",
		)

		assert.NotNil(t, client)
		serverURL := client.APIClient.GetConfig().Servers[0].URL
		assert.Contains(t, serverURL, ":8443")
	})
}

func TestN8nClientReturnsPointer(t *testing.T) {
	t.Run("returns pointer to N8nClient", func(t *testing.T) {
		client := NewN8nClient("https://n8n.example.com", "test-key")

		// Verify it's a pointer
		assert.IsType(t, &N8nClient{}, client)

		// Verify we can modify it
		client.BaseURL = "https://modified.example.com"
		assert.Equal(t, "https://modified.example.com", client.BaseURL)
	})
}

func TestN8nClientConfigurationObject(t *testing.T) {
	t.Run("configuration object is valid", func(t *testing.T) {
		client := NewN8nClient("https://n8n.example.com", "test-key")

		cfg := client.APIClient.GetConfig()
		assert.NotNil(t, cfg, "Configuration should not be nil")
		assert.NotNil(t, cfg.Servers, "Servers should not be nil")
		assert.NotNil(t, cfg.DefaultHeader, "DefaultHeader should not be nil")
	})

	t.Run("configuration has one server", func(t *testing.T) {
		client := NewN8nClient("https://n8n.example.com", "test-key")

		cfg := client.APIClient.GetConfig()
		assert.Len(t, cfg.Servers, 1, "Should have exactly one server configured")
	})

	t.Run("server configuration has required fields", func(t *testing.T) {
		client := NewN8nClient("https://n8n.example.com", "test-key")

		server := client.APIClient.GetConfig().Servers[0]
		assert.NotEmpty(t, server.URL, "Server URL should not be empty")
		assert.NotEmpty(t, server.Description, "Server description should not be empty")
	})
}

func TestN8nClientConcurrentCreation(t *testing.T) {
	t.Run("can create clients concurrently", func(t *testing.T) {
		done := make(chan *N8nClient, 10)

		for i := 0; i < 10; i++ {
			go func(id int) {
				client := NewN8nClient(
					"https://n8n.example.com",
					"test-key",
				)
				done <- client
			}(i)
		}

		// Collect all clients
		clients := make([]*N8nClient, 10)
		for i := 0; i < 10; i++ {
			clients[i] = <-done
			assert.NotNil(t, clients[i])
		}

		// Verify all clients were created
		for i, client := range clients {
			assert.NotNil(t, client, "Client %d should not be nil", i)
			assert.NotNil(t, client.APIClient, "Client %d should have API client", i)
		}
	})
}
