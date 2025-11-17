// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models defines data structures for variable resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSources maps the Terraform schema attributes for the variables datasource.
// It represents the complete set of variables data returned by the n8n API with optional filtering.
type DataSources struct {
	ProjectID types.String `tfsdk:"project_id"`
	State     types.String `tfsdk:"state"`
	Variables []Item       `tfsdk:"variables"`
}
