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
	// Test with value set.
	t.Run("returns N8N_API_KEY when set", func(t *testing.T) {
		t.Helper()
		t.Setenv("N8N_API_KEY", "test-api-key")

		result := getEnvAPIKey()

		assert.Equal(t, "test-api-key", result)
	})

	// Test empty case only when not running acceptance tests.
	// Acceptance tests set these env vars globally and t.Setenv cannot override them.
	t.Run("error case - returns empty string when N8N_API_KEY not set", func(t *testing.T) {
		t.Helper()
		// When TF_ACC is set, env vars are preset and cannot be cleared with t.Setenv.
		if os.Getenv("TF_ACC") != "" {
			return
		}
		t.Setenv("N8N_API_KEY", "")

		result := getEnvAPIKey()

		assert.Equal(t, "", result)
	})
}

// Test_getEnvBaseURL tests the getEnvBaseURL function.
func Test_getEnvBaseURL(t *testing.T) {
	// Test with value set.
	t.Run("returns N8N_API_URL when set", func(t *testing.T) {
		t.Helper()
		t.Setenv("N8N_API_URL", "https://test.example.com")

		result := getEnvBaseURL()

		assert.Equal(t, "https://test.example.com", result)
	})

	// Test empty case only when not running acceptance tests.
	// Acceptance tests set these env vars globally and t.Setenv cannot override them.
	t.Run("error case - returns empty string when N8N_API_URL not set", func(t *testing.T) {
		t.Helper()
		// When TF_ACC is set, env vars are preset and cannot be cleared with t.Setenv.
		if os.Getenv("TF_ACC") != "" {
			return
		}
		t.Setenv("N8N_API_URL", "")

		result := getEnvBaseURL()

		assert.Equal(t, "", result)
	})
}
