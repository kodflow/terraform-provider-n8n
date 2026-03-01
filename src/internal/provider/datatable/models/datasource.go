// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models contains data models for the datatable domain.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSource maps the Terraform schema to a single data table from the n8n API.
type DataSource struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	ProjectID types.String `tfsdk:"project_id"`
}
