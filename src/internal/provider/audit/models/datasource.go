// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models defines data structures for audit resources.
package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSource describes the audit datasource model.
// Each risk report is serialized as a JSON string since the reports contain
// dynamic key-value structures that vary per n8n instance.
type DataSource struct {
	CredentialsRiskReport types.String `tfsdk:"credentials_risk_report" json:"credentials_risk_report"`
	DatabaseRiskReport    types.String `tfsdk:"database_risk_report" json:"database_risk_report"`
	FilesystemRiskReport  types.String `tfsdk:"filesystem_risk_report" json:"filesystem_risk_report"`
	NodesRiskReport       types.String `tfsdk:"nodes_risk_report" json:"nodes_risk_report"`
	InstanceRiskReport    types.String `tfsdk:"instance_risk_report" json:"instance_risk_report"`
}
