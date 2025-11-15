// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package models defines data structures for tag resources.
package models

// DataSources maps Terraform schema attributes for tag list data.
// It represents the complete data structure returned from the n8n tags API.
type DataSources struct {
	Tags []Item `tfsdk:"tags"`
}
