// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models defines data structures for execution resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Item describes a single execution in the list datasource.
// It holds the execution identifier and its runtime metadata.
type Item struct {
	ID         types.String `tfsdk:"id" json:"id"`
	Mode       types.String `tfsdk:"mode" json:"mode"`
	Status     types.String `tfsdk:"status" json:"status"`
	WorkflowID types.String `tfsdk:"workflow_id" json:"workflow_id"`
	Finished   types.Bool   `tfsdk:"finished" json:"finished"`
	CreatedAt  types.String `tfsdk:"created_at" json:"created_at"`
	StartedAt  types.String `tfsdk:"started_at" json:"started_at"`
	StoppedAt  types.String `tfsdk:"stopped_at" json:"stopped_at"`
}

// DataSources describes the list executions datasource model.
// It holds optional filters and the collection of execution items.
type DataSources struct {
	WorkflowID types.String `tfsdk:"workflow_id" json:"workflow_id"`
	Status     types.String `tfsdk:"status" json:"status"`
	Executions []Item       `tfsdk:"executions" json:"executions"`
}
