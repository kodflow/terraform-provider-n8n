package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Resource describes the resource data model.
// Maps n8n project attributes to Terraform schema, storing project metadata and configuration.
type Resource struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}
