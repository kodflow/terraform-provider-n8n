// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models contains data models for the tag domain.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Package tag contains resources and datasources for tag management

// Resource describes the tag resource data model.
// Maps n8n tag attributes to Terraform schema for managing workflow tags and organization.
type Resource struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}
