// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models defines data structures for credential resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CredentialItem describes a single credential in the list datasource.
// It contains credential metadata including name, type, and timestamps.
type CredentialItem struct {
	ID        types.String `tfsdk:"id" json:"id" dto:"out,query,priv"`
	Name      types.String `tfsdk:"name" json:"name" dto:"out,query,pub"`
	Type      types.String `tfsdk:"type" json:"type" dto:"out,query,pub"`
	CreatedAt types.String `tfsdk:"created_at" json:"createdAt" dto:"out,query,pub"`
	UpdatedAt types.String `tfsdk:"updated_at" json:"updatedAt" dto:"out,query,pub"`
}

// DataSources describes the list credentials datasource model.
// It holds all credentials returned from the n8n API.
type DataSources struct {
	Credentials []CredentialItem `tfsdk:"credentials" json:"credentials" dto:"out,query,priv"`
}
