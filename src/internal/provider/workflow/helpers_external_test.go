// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package workflow_test provides black-box tests for the workflow package helpers.
package workflow_test

import (
	"testing"

	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow"
	"github.com/stretchr/testify/assert"
)

// TestCallerPolicyDefault verifies the exported CallerPolicyDefault constant value.
func TestCallerPolicyDefault(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "constant has correct value",
			expected: "workflowsFromSameOwner",
		},
		{
			name:     "error case - constant is non-empty",
			expected: "workflowsFromSameOwner",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, workflow.CallerPolicyDefault)
			assert.NotEmpty(t, workflow.CallerPolicyDefault)
		})
	}
}
