// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package audit implements audit report datasource functionality.
package audit

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/audit/models"
)

// marshalRiskReport serializes a map[string]any risk report to a JSON string.
// Returns types.StringNull() when the map is nil or empty.
//
// Params:
//   - report: The risk report map to serialize
//
// Returns:
//   - types.String: JSON string or null when empty
func marshalRiskReport(report map[string]any) (jsonStr types.String) {
	//: Return null when the report is absent.
	if len(report) == 0 {
		//: Return null string for absent reports.
		return types.StringNull()
	}

	data, err := json.Marshal(report)
	//: Return null when serialization fails.
	if err != nil {
		//: Return null string on marshal error.
		return types.StringNull()
	}

	//: Return the serialized JSON string.
	return types.StringValue(string(data))
}

// mapAuditToDataSource maps an SDK Audit response to the datasource model.
//
// Params:
//   - audit: The n8n SDK Audit to map
//   - data: The models.DataSource to populate
func mapAuditToDataSource(audit *n8nsdk.Audit, data *models.DataSource) {
	data.CredentialsRiskReport = marshalRiskReport(audit.CredentialsRiskReport)
	data.DatabaseRiskReport = marshalRiskReport(audit.DatabaseRiskReport)
	data.FilesystemRiskReport = marshalRiskReport(audit.FilesystemRiskReport)
	data.NodesRiskReport = marshalRiskReport(audit.NodesRiskReport)
	data.InstanceRiskReport = marshalRiskReport(audit.InstanceRiskReport)
}
