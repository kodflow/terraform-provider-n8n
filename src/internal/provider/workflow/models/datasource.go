package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSource maps the Terraform schema attributes for a single workflow datasource.
// It represents workflow metadata including its identifier, name, and activation status.
type DataSource struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Active types.Bool   `tfsdk:"active"`
}
