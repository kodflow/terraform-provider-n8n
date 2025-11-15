// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package models defines data structures for tag resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSource maps the Terraform schema to a single tag from the n8n API.
// It contains tag metadata including name and creation/update timestamps.
type DataSource struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}
