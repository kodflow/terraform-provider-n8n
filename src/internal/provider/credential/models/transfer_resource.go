// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models contains data models for the credential domain.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Package credential contains resources and datasources for credential management

// TransferResource describes the credential transfer resource data model.
// Maps n8n credential transfer attributes to Terraform schema for transferring credentials between users or projects.
type TransferResource struct {
	ID                   types.String `tfsdk:"id"`
	CredentialID         types.String `tfsdk:"credential_id"`
	DestinationProjectID types.String `tfsdk:"destination_project_id"`
	TransferredAt        types.String `tfsdk:"transferred_at"`
}
