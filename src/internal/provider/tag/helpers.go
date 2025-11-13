// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package tag contains helper functions for tag operations.
package tag

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/tag/models"
)

// mapTagToDataSourceModel maps an SDK tag to the datasource model.
//
// Params:
//   - tag: The n8n SDK tag to map
//   - data: The models.DataSource to populate
//
// Returns:
//   - None
func mapTagToDataSourceModel(tag *n8nsdk.Tag, data *models.DataSource) {
	// Check for non-nil value.
	if tag.Id != nil {
		data.ID = types.StringValue(*tag.Id)
	}
	data.Name = types.StringValue(tag.Name)
	// Check for non-nil value.
	if tag.CreatedAt != nil {
		data.CreatedAt = types.StringValue(tag.CreatedAt.String())
	}
	// Check for non-nil value.
	if tag.UpdatedAt != nil {
		data.UpdatedAt = types.StringValue(tag.UpdatedAt.String())
	}
}

// findTagByName searches for a tag by name in a tag list.
//
// Params:
//   - tags: Slice of n8n SDK tags to search through
//   - name: The tag name to search for
//
// Returns:
//   - *n8nsdk.Tag: Pointer to the found tag, or nil if not found
//   - bool: True if tag was found, false otherwise
func findTagByName(tags []n8nsdk.Tag, name string) (*n8nsdk.Tag, bool) {
	// Iterate over items.
	for _, tag := range tags {
		// Check condition.
		if tag.Name == name {
			// Return success status.
			return &tag, true
		}
	}
	// Return nil to indicate failure.
	return nil, false
}
