// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models defines data structures for source control resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Resource describes the source control pull resource model.
// Triggering a pull imports workflows, credentials, variables, and tags from the
// configured git repository into the n8n instance.
type Resource struct {
	// Force triggers the pull even when no remote changes are detected.
	Force types.Bool `tfsdk:"force" json:"force"`
	// WorkflowsImported contains a JSON summary of imported workflow changes.
	WorkflowsImported types.String `tfsdk:"workflows_imported" json:"workflows_imported"`
	// CredentialsImported contains a JSON summary of imported credential changes.
	CredentialsImported types.String `tfsdk:"credentials_imported" json:"credentials_imported"`
	// TagsImported contains a JSON summary of imported tag changes.
	TagsImported types.String `tfsdk:"tags_imported" json:"tags_imported"`
	// VariablesImported contains a JSON summary of imported variable changes.
	VariablesImported types.String `tfsdk:"variables_imported" json:"variables_imported"`
}
