// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

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
		expectedResult string
	}{
		{
			name:           "returns N8N_API_KEY when set",
			envAPIKey:      "test-api-key",
			expectedResult: "test-api-key",
		},
		{
			name:           "error case - returns empty string when N8N_API_KEY not set",
			envAPIKey:      "",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			// Always use t.Setenv to ensure clean environment
			// This properly clears any existing value when testing empty case
			t.Setenv("N8N_API_KEY", tt.envAPIKey)

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
		expectedResult string
	}{
		{
			name:           "returns N8N_API_URL when set",
			envAPIURL:      "https://test.example.com",
			expectedResult: "https://test.example.com",
		},
		{
			name:           "error case - returns empty string when N8N_API_URL not set",
			envAPIURL:      "",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			// Always use t.Setenv to ensure clean environment
			// This properly clears any existing value when testing empty case
			t.Setenv("N8N_API_URL", tt.envAPIURL)

			result := getEnvBaseURL()

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
