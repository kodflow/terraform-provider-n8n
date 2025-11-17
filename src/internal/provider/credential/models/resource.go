// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models contains data models for the credential domain.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Resource describes the resource data model.
// Maps n8n credential attributes to Terraform schema, storing credential metadata and sensitive data.
type Resource struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Type      types.String `tfsdk:"type"`
	Data      types.Map    `tfsdk:"data"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}
