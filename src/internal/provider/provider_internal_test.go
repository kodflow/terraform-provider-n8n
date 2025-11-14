// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_getEnvAPIKey tests the getEnvAPIKey function.
func Test_getEnvAPIKey(t *testing.T) {
	tests := []struct {
		name           string
		envAPIKey      string
		envAPIToken    string
		expectedResult string
	}{
		{
			name:           "prefers N8N_API_KEY over N8N_API_TOKEN",
			envAPIKey:      "preferred-key",
			envAPIToken:    "legacy-key",
			expectedResult: "preferred-key",
		},
		{
			name:           "falls back to N8N_API_TOKEN when N8N_API_KEY not set",
			envAPIKey:      "",
			envAPIToken:    "legacy-key",
			expectedResult: "legacy-key",
		},
		{
			name:           "returns empty string when no environment variable is set",
			envAPIKey:      "",
			envAPIToken:    "",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			// Set environment variables only if not empty
			if tt.envAPIKey != "" {
				t.Setenv("N8N_API_KEY", tt.envAPIKey)
			}
			// Set environment variables only if not empty
			if tt.envAPIToken != "" {
				t.Setenv("N8N_API_TOKEN", tt.envAPIToken)
			}

			result := getEnvAPIKey()

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

// Test_getEnvBaseURL tests the getEnvBaseURL function.
func Test_getEnvBaseURL(t *testing.T) {
	tests := []struct {
		name           string
		envAPIURL      string
		envURL         string
		expectedResult string
	}{
		{
			name:           "prefers N8N_API_URL over N8N_URL",
			envAPIURL:      "https://preferred.example.com",
			envURL:         "https://legacy.example.com",
			expectedResult: "https://preferred.example.com",
		},
		{
			name:           "falls back to N8N_URL when N8N_API_URL not set",
			envAPIURL:      "",
			envURL:         "https://legacy.example.com",
			expectedResult: "https://legacy.example.com",
		},
		{
			name:           "returns empty string when no environment variable is set",
			envAPIURL:      "",
			envURL:         "",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			// Set environment variables only if not empty
			if tt.envAPIURL != "" {
				t.Setenv("N8N_API_URL", tt.envAPIURL)
			}
			// Set environment variables only if not empty
			if tt.envURL != "" {
				t.Setenv("N8N_URL", tt.envURL)
			}

			result := getEnvBaseURL()

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
