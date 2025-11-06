// Package user contains resources and datasources for user management
package user

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UserResourceModel describes the resource data model.
// Maps n8n user attributes to Terraform schema, storing user identity, role, and account information.
type UserResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	Role      types.String `tfsdk:"role"`
	IsPending types.Bool   `tfsdk:"is_pending"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

