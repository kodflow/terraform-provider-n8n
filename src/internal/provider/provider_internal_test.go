// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package provider

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_getEnvAPIKey tests the getEnvAPIKey function.
func Test_getEnvAPIKey(t *testing.T) {
	tests := []struct {
		name          string
		envValue      string
		want          string
		skipWhenTFACC bool
	}{
		{
			name:          "returns N8N_API_KEY when set",
			envValue:      "test-api-key",
			want:          "test-api-key",
			skipWhenTFACC: false,
		},
		{
			name:          "returns N8N_API_KEY with special characters",
			envValue:      "key-with-special-chars!@#$%",
			want:          "key-with-special-chars!@#$%",
			skipWhenTFACC: false,
		},
		{
			name:          "returns N8N_API_KEY with JWT format",
			envValue:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			want:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			skipWhenTFACC: false,
		},
		{
			name:          "error case - returns empty string when N8N_API_KEY not set",
			envValue:      "",
			want:          "",
			skipWhenTFACC: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip tests that require empty env vars when TF_ACC is set.
			// Acceptance tests set these env vars globally and t.Setenv cannot override them.
			if tt.skipWhenTFACC && os.Getenv("TF_ACC") != "" {
				return
			}

			t.Setenv("N8N_API_KEY", tt.envValue)

			result := getEnvAPIKey()

			assert.Equal(t, tt.want, result)
		})
	}
}

// Test_getEnvBaseURL tests the getEnvBaseURL function.
func Test_getEnvBaseURL(t *testing.T) {
	tests := []struct {
		name          string
		envValue      string
		want          string
		skipWhenTFACC bool
	}{
		{
			name:          "returns N8N_API_URL when set",
			envValue:      "https://test.example.com",
			want:          "https://test.example.com",
			skipWhenTFACC: false,
		},
		{
			name:          "returns N8N_API_URL with port",
			envValue:      "http://localhost:5678",
			want:          "http://localhost:5678",
			skipWhenTFACC: false,
		},
		{
			name:          "returns N8N_API_URL with path",
			envValue:      "https://n8n.example.com/api/v1",
			want:          "https://n8n.example.com/api/v1",
			skipWhenTFACC: false,
		},
		{
			name:          "error case - returns empty string when N8N_API_URL not set",
			envValue:      "",
			want:          "",
			skipWhenTFACC: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip tests that require empty env vars when TF_ACC is set.
			// Acceptance tests set these env vars globally and t.Setenv cannot override them.
			if tt.skipWhenTFACC && os.Getenv("TF_ACC") != "" {
				return
			}

			t.Setenv("N8N_API_URL", tt.envValue)

			result := getEnvBaseURL()

			assert.Equal(t, tt.want, result)
		})
	}
}
