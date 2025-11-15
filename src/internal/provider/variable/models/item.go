// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package models defines data structures for variable resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Item maps individual variable attributes within the Terraform schema.
// Each item represents a single variable with its ID, key, value, type, and associated project.
type Item struct {
	ID        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	Type      types.String `tfsdk:"type"`
	ProjectID types.String `tfsdk:"project_id"`
}
