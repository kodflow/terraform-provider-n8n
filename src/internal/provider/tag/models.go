// Package tag contains resources and datasources for tag management
package tag

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TagResourceModel describes the resource data model.
// Maps n8n tag attributes to Terraform schema, storing tag metadata including name and timestamps.
type TagResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}
