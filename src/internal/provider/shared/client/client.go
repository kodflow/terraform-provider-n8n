package client

import (
	"github.com/kodflow/n8n/sdk/n8nsdk"
)

// N8nClient wraps the generated SDK client with provider-specific configuration.
// It provides a convenient interface for resources to interact with the n8n API.
type N8nClient struct {
	// APIClient is the generated SDK client
	APIClient *n8nsdk.APIClient

	// BaseURL is the base URL of the n8n instance
	BaseURL string

	// APIKey is the API key for authentication
	APIKey string
}

// NewN8nClient creates a new N8nClient instance with the given configuration.
//
// Params:
//   - baseURL: the base URL of the n8n instance (e.g., "https://n8n.example.com")
//   - apiKey: the API key for authentication
//
// Returns:
//   - *N8nClient: configured client ready for API calls
func NewN8nClient(baseURL, apiKey string) *N8nClient {
	// Create SDK configuration
	cfg := n8nsdk.NewConfiguration()

	// Set the base URL (override the default /api/v1)
	cfg.Servers = n8nsdk.ServerConfigurations{
		n8nsdk.ServerConfiguration{
			URL:         baseURL + "/api/v1",
			Description: "n8n instance",
		},
	}

	// Add API key to default headers
	cfg.AddDefaultHeader("X-N8N-API-KEY", apiKey)

	// Create the API client
	apiClient := n8nsdk.NewAPIClient(cfg)

	// Return result.
	return &N8nClient{
		APIClient: apiClient,
		BaseURL:   baseURL,
		APIKey:    apiKey,
	}
}
