// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models defines data structures for variable resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSource maps Terraform schema attributes for variable data.
// It represents a single variable with all related attributes from the n8n API.
type DataSource struct {
	ID        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	Type      types.String `tfsdk:"type"`
	ProjectID types.String `tfsdk:"project_id"`
}
