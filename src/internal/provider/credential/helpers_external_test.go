// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package credential_test provides black-box tests for the credential package helpers.
package credential_test

import (
	"testing"

	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/credential"
	"github.com/stretchr/testify/assert"
)

// TestFloat64BitSize verifies the exported constant value.
func TestFloat64BitSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected int
	}{
		{
			name:     "constant equals 64",
			expected: 64,
		},
		{
			name:     "error case - constant is positive",
			expected: 64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, credential.Float64BitSize)
		})
	}
}
