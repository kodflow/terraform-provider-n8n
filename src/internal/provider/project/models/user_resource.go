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
