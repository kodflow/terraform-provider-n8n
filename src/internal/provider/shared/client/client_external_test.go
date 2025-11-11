package client_test

import (
	"strings"
	"testing"

	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewN8nClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "creates client with base URL and API key", wantErr: false},
		{name: "creates client with empty base URL", wantErr: false},
		{name: "creates client with empty API key", wantErr: false},
		{name: "creates client with both empty values", wantErr: false},
		{name: "API client configuration includes base URL", wantErr: false},
		{name: "API client configuration includes API key header", wantErr: false},
		{name: "creates multiple independent clients", wantErr: false},
		{name: "handles various base URL formats", wantErr: false},
		{name: "handles various API key formats", wantErr: false},
		{name: "server description is set", wantErr: false},
		{name: "configuration is properly initialized", wantErr: false},
		{name: "error case - nil checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "creates client with base URL and API key":
				baseURL := "https://n8n.example.com"
				apiKey := "test-api-key-123"

				c := client.NewN8nClient(baseURL, apiKey)

				assert.NotNil(t, c, "Client should not be nil")
				assert.Equal(t, baseURL, c.BaseURL, "BaseURL should be set")
				assert.Equal(t, apiKey, c.APIKey, "APIKey should be set")
				assert.NotNil(t, c.APIClient, "APIClient should be initialized")

			case "creates client with empty base URL":
				c := client.NewN8nClient("", "test-key")

				assert.NotNil(t, c, "Client should be created")
				assert.Equal(t, "", c.BaseURL, "Should handle empty base URL")
				assert.NotNil(t, c.APIClient, "APIClient should be initialized")

			case "creates client with empty API key":
				c := client.NewN8nClient("https://n8n.example.com", "")

				assert.NotNil(t, c, "Client should be created")
				assert.Equal(t, "", c.APIKey, "Should handle empty API key")
				assert.NotNil(t, c.APIClient, "APIClient should be initialized")

			case "creates client with both empty values":
				c := client.NewN8nClient("", "")

				assert.NotNil(t, c, "Client should be created")
				assert.Equal(t, "", c.BaseURL, "BaseURL should be empty")
				assert.Equal(t, "", c.APIKey, "APIKey should be empty")
				assert.NotNil(t, c.APIClient, "APIClient should be initialized")

			case "API client configuration includes base URL":
				baseURL := "https://n8n.example.com"
				c := client.NewN8nClient(baseURL, "test-key")

				require.NotNil(t, c.APIClient)
				require.NotNil(t, c.APIClient.GetConfig())
				require.NotEmpty(t, c.APIClient.GetConfig().Servers)

				serverURL := c.APIClient.GetConfig().Servers[0].URL
				assert.Equal(t, baseURL+"/api/v1", serverURL, "Server URL should include /api/v1 suffix")

			case "API client configuration includes API key header":
				apiKey := "test-api-key-123"
				c := client.NewN8nClient("https://n8n.example.com", apiKey)

				require.NotNil(t, c.APIClient)
				require.NotNil(t, c.APIClient.GetConfig())

				defaultHeaders := c.APIClient.GetConfig().DefaultHeader
				assert.Contains(t, defaultHeaders, "X-N8N-API-KEY", "Should have X-N8N-API-KEY header")
				assert.Equal(t, apiKey, defaultHeaders["X-N8N-API-KEY"], "API key should be in header")

			case "creates multiple independent clients":
				c1 := client.NewN8nClient("https://n8n1.example.com", "key1")
				c2 := client.NewN8nClient("https://n8n2.example.com", "key2")

				assert.NotEqual(t, c1.BaseURL, c2.BaseURL, "Clients should have different base URLs")
				assert.NotEqual(t, c1.APIKey, c2.APIKey, "Clients should have different API keys")
				assert.NotSame(t, c1, c2, "Should be different instances")
				assert.NotSame(t, c1.APIClient, c2.APIClient, "Should have different API clients")

			case "handles various base URL formats":
				testCases := []struct {
					baseURL  string
					expected string
				}{
					{"https://n8n.example.com", "https://n8n.example.com/api/v1"},
					{"http://n8n.example.com", "http://n8n.example.com/api/v1"},
					{"https://n8n.example.com:8080", "https://n8n.example.com:8080/api/v1"},
					{"http://localhost:5678", "http://localhost:5678/api/v1"},
				}

				for _, tc := range testCases {
					c := client.NewN8nClient(tc.baseURL, "test-key")
					require.NotNil(t, c.APIClient)
					require.NotNil(t, c.APIClient.GetConfig())
					require.NotEmpty(t, c.APIClient.GetConfig().Servers)

					serverURL := c.APIClient.GetConfig().Servers[0].URL
					assert.Equal(t, tc.expected, serverURL)
				}

			case "handles various API key formats":
				testKeys := []string{
					"abc123xyz789",
					"key-with-dashes",
					strings.Repeat("a", 100),
					"k",
				}

				for _, apiKey := range testKeys {
					c := client.NewN8nClient("https://n8n.example.com", apiKey)
					assert.Equal(t, apiKey, c.APIKey)

					defaultHeaders := c.APIClient.GetConfig().DefaultHeader
					assert.Equal(t, apiKey, defaultHeaders["X-N8N-API-KEY"])
				}

			case "server description is set":
				c := client.NewN8nClient("https://n8n.example.com", "test-key")

				require.NotNil(t, c.APIClient)
				require.NotNil(t, c.APIClient.GetConfig())
				require.NotEmpty(t, c.APIClient.GetConfig().Servers)

				description := c.APIClient.GetConfig().Servers[0].Description
				assert.Equal(t, "n8n instance", description)

			case "configuration is properly initialized":
				c := client.NewN8nClient("https://n8n.example.com", "test-key")

				require.NotNil(t, c.APIClient)
				cfg := c.APIClient.GetConfig()
				require.NotNil(t, cfg)

				assert.NotEmpty(t, cfg.Servers, "Servers should be configured")
				assert.Len(t, cfg.Servers, 1, "Should have exactly one server")
				assert.NotNil(t, cfg.DefaultHeader, "Default headers should be set")

			case "error case - nil checks":
				c := client.NewN8nClient("https://n8n.example.com", "test-key")
				assert.NotNil(t, c)
				assert.NotNil(t, c.APIClient)
				assert.NotNil(t, c.APIClient.GetConfig())
			}
		})
	}
}

func TestN8nClientStructure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "client has all required fields", wantErr: false},
		{name: "client fields are accessible", wantErr: false},
		{name: "error case - field validation", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "client has all required fields":
				c := client.NewN8nClient("https://n8n.example.com", "test-key")

				assert.NotNil(t, c.APIClient)
				assert.NotEmpty(t, c.BaseURL)
				assert.NotEmpty(t, c.APIKey)

			case "client fields are accessible":
				baseURL := "https://n8n.example.com"
				apiKey := "test-key-123"
				c := client.NewN8nClient(baseURL, apiKey)

				assert.Equal(t, baseURL, c.BaseURL)
				assert.Equal(t, apiKey, c.APIKey)
				assert.NotNil(t, c.APIClient)

			case "error case - field validation":
				c := client.NewN8nClient("", "")
				assert.NotNil(t, c)
				assert.Equal(t, "", c.BaseURL)
				assert.Equal(t, "", c.APIKey)
			}
		})
	}
}

func TestN8nClientEdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "handles nil after creation", wantErr: false},
		{name: "handles very long base URL", wantErr: false},
		{name: "handles special characters", wantErr: false},
		{name: "error case - whitespace handling", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "handles nil after creation":
				c := client.NewN8nClient("https://n8n.example.com", "test-key")
				assert.NotNil(t, c)

			case "handles very long base URL":
				longURL := "https://" + strings.Repeat("a", 500) + ".com"
				c := client.NewN8nClient(longURL, "test-key")

				assert.NotNil(t, c)
				assert.Equal(t, longURL, c.BaseURL)

			case "handles special characters":
				specialURL := "https://n8n-test.example.com:8080"
				c := client.NewN8nClient(specialURL, "test-key")

				assert.NotNil(t, c)
				assert.Equal(t, specialURL, c.BaseURL)

			case "error case - whitespace handling":
				c := client.NewN8nClient("  https://n8n.example.com  ", "  key  ")
				assert.NotNil(t, c)
				assert.Equal(t, "  https://n8n.example.com  ", c.BaseURL)
				assert.Equal(t, "  key  ", c.APIKey)
			}
		})
	}
}
