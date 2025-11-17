// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models defines data structures for workflow resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Transfer describes the workflow transfer resource data model.
// Captures the workflow ID and destination project ID for transfer operations.
type Transfer struct {
	ID                   types.String `tfsdk:"id"`
	WorkflowID           types.String `tfsdk:"workflow_id"`
	DestinationProjectID types.String `tfsdk:"destination_project_id"`
	TransferredAt        types.String `tfsdk:"transferred_at"`
}
