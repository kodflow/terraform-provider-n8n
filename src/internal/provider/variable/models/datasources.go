// Package models defines data structures for variable resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSources maps the Terraform schema attributes for the variables datasource.
// It represents the complete set of variables data returned by the n8n API with optional filtering.
type DataSources struct {
	ProjectID types.String `tfsdk:"project_id"`
	State     types.String `tfsdk:"state"`
	Variables []Item       `tfsdk:"variables"`
}
