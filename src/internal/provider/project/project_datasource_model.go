package project

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProjectDataSourceModel maps the Terraform schema to a single project from the n8n API.
// It contains project metadata including timestamps, type, and descriptive information.
type ProjectDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
}
