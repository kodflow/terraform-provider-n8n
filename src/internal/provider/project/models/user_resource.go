// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package models defines data structures for project resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UserResource describes the resource data model.
// Maps n8n project-user relationship attributes to Terraform schema, storing user roles and project associations.
type UserResource struct {
	ID        types.String `tfsdk:"id"`
	ProjectID types.String `tfsdk:"project_id"`
	UserID    types.String `tfsdk:"user_id"`
	Role      types.String `tfsdk:"role"`
}
