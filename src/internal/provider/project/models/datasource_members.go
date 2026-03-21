// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models defines data structures for project resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MemberItem describes a single project member in the datasource.
// It maps the member's profile and role information from the n8n API.
type MemberItem struct {
	UserID    types.String `tfsdk:"user_id"    json:"user_id"    dto:"out,query,priv"`
	Role      types.String `tfsdk:"role"       json:"role"`
	Email     types.String `tfsdk:"email"      json:"email"`
	FirstName types.String `tfsdk:"first_name" json:"first_name"`
	LastName  types.String `tfsdk:"last_name"  json:"last_name"`
}

// DataSourceMembers describes the project members datasource model.
// It holds all members for a given project from the n8n API.
type DataSourceMembers struct {
	ProjectID types.String `tfsdk:"project_id" json:"project_id" dto:"out,query,priv"`
	Members   []MemberItem `tfsdk:"members"    json:"members"`
}
