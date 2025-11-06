// Package project contains resources and datasources for project management
package project

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProjectUserResourceModel describes the resource data model.
// Maps n8n project-user relationship attributes to Terraform schema, storing user roles and project associations.
type ProjectUserResourceModel struct {
	ID        types.String `tfsdk:"id"`
	ProjectID types.String `tfsdk:"project_id"`
	UserID    types.String `tfsdk:"user_id"`
	Role      types.String `tfsdk:"role"`
}
