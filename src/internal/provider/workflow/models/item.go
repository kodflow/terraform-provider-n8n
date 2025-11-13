// Package models defines data structures for workflow resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Item maps individual workflow attributes within the Terraform schema.
// Each item represents a single workflow with its identifier, name, and activation status.
type Item struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Active types.Bool   `tfsdk:"active"`
}
