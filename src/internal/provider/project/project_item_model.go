package project

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProjectItemModel represents a single project in the list returned from the n8n API.
// It contains project metadata including name, type, timestamps, and descriptive information.
type ProjectItemModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
}
