// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package models defines data structures for workflow resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSources maps the Terraform schema attributes for the workflows datasource.
// It represents the complete set of workflows data returned by the n8n API with optional active status filtering.
type DataSources struct {
	Workflows []Item     `tfsdk:"workflows"`
	Active    types.Bool `tfsdk:"active"`
}
