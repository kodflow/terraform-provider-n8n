// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models contains data models for the user domain.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Package user contains resources and datasources for user management

// Resource describes the user resource data model.
// Maps n8n user attributes to Terraform schema for managing user accounts and roles.
type Resource struct {
	ID        types.String `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	Role      types.String `tfsdk:"role"`
	IsPending types.Bool   `tfsdk:"is_pending"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}
