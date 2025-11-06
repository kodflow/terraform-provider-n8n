// Package project contains resources and datasources for project management
package project

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProjectResourceModel describes the resource data model.
// Maps n8n project attributes to Terraform schema, storing project metadata and configuration.
type ProjectResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

