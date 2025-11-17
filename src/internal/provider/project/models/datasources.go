// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models defines data structures for project resources.
package models

// DataSources maps the Terraform schema to the datasource response.
// It represents a list of projects retrieved from the n8n API with all project details.
type DataSources struct {
	Projects []Item `tfsdk:"projects"`
}
