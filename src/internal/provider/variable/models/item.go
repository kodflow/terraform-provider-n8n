package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Item maps individual variable attributes within the Terraform schema.
// Each item represents a single variable with its ID, key, value, type, and associated project.
type Item struct {
	ID        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	Type      types.String `tfsdk:"type"`
	ProjectID types.String `tfsdk:"project_id"`
}
