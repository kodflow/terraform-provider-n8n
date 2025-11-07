package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Item represents a single tag in the list.
// It maps individual tag attributes from the n8n API to Terraform schema.
type Item struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}
