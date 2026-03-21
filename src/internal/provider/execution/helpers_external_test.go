// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package execution_test

import (
	"testing"

	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/execution"
)

// TestNewExecutionTagsResourceNotNil verifies that NewExecutionTagsResource
// returns a non-nil resource instance ready for Terraform framework registration.
func TestNewExecutionTagsResourceNotNil(t *testing.T) {
	t.Parallel()

	r := execution.NewExecutionTagsResource()
	//: Verify resource was created successfully.
	if r == nil {
		t.Error("expected non-nil ExecutionTagsResource")
	}
}

// TestNewExecutionTagsResourceWrapperNotNil verifies that the Terraform
// provider wrapper returns a non-nil resource.Resource interface implementation.
func TestNewExecutionTagsResourceWrapperNotNil(t *testing.T) {
	t.Parallel()

	r := execution.NewExecutionTagsResourceWrapper()
	//: Verify wrapper returns valid resource.
	if r == nil {
		t.Error("expected non-nil resource.Resource from wrapper")
	}
}
