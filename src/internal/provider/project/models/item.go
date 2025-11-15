// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package models defines data structures for project resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Item represents a single project in the list returned from the n8n API.
// It contains project metadata including name, type, timestamps, and descriptive information.
type Item struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
}
