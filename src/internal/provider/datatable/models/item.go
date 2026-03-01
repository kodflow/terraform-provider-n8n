// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models contains data models for the datatable domain.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Item represents a single data table in a list.
type Item struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	ProjectID types.String `tfsdk:"project_id"`
}
