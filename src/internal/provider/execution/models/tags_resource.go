// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models defines data structures for execution resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TagsResource describes the execution tags resource data model.
// It maps execution tag assignments between an execution and a set of tags.
type TagsResource struct {
	ExecutionID types.String `tfsdk:"execution_id"`
	TagIDs      types.Set    `tfsdk:"tag_ids"`
}
