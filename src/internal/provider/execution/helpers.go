// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package execution implements execution management resources.
package execution

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// ElementsAsser is a minimal interface for extracting elements from a Terraform set value.
type ElementsAsser interface {
	ElementsAs(ctx context.Context, target any, allowUnhandled bool) diag.Diagnostics
}

// diagnostics is a minimal interface for appending and checking Terraform diagnostics.
type diagnostics interface {
	Append(in ...diag.Diagnostic)
	HasError() bool
}

// extractTagIDs extracts a slice of tag ID strings from a set value.
//
// Params:
//   - ctx: Context for the operation
//   - tagIDs: The Terraform set of tag IDs
//   - diags: Diagnostics to append errors to
//
// Returns:
//   - []string: The extracted tag IDs
func extractTagIDs(ctx context.Context, tagIDs ElementsAsser, diags diagnostics) (items []string) {
	var tagIDStrings []string
	diags.Append(tagIDs.ElementsAs(ctx, &tagIDStrings, false)...)
	//: Check for error.
	if diags.HasError() {
		//: Return empty slice on error.
		return nil
	}
	//: Check for nil to return empty slice instead of nil.
	if tagIDStrings == nil {
		//: Return nil result.
		return nil
	}
	//: Return result.
	return tagIDStrings
}
