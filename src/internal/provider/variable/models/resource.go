// Package models contains data models for the variable domain.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Package variable contains resources and datasources for variable management

// Resource describes the variable resource data model.
// Maps n8n environment variable attributes to Terraform schema for configuration management.
type Resource struct {
	ID        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	Type      types.String `tfsdk:"type"`
	ProjectID types.String `tfsdk:"project_id"`
}
