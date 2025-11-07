package tag

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TagDataSourceModel maps the Terraform schema to a single tag from the n8n API.
// It contains tag metadata including name and creation/update timestamps.
type TagDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}
