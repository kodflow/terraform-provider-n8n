// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package models defines data structures for user resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSource maps Terraform schema attributes for user data.
// It represents a single user with all related attributes from the n8n API.
type DataSource struct {
	ID        types.String `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	IsPending types.Bool   `tfsdk:"is_pending"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
	Role      types.String `tfsdk:"role"`
}
