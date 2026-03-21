// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package sourcecontrol implements source control pull resource functionality.
package sourcecontrol

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/sourcecontrol/models"
)

// marshalImportedSlice serializes a slice to a JSON string.
// Returns types.StringNull() when the slice is nil or empty.
//
// Params:
//   - items: The slice to serialize
//
// Returns:
//   - types.String: JSON string or null when empty
func marshalImportedSlice(items any) (jsonStr types.String) {
	data, err := json.Marshal(items)
	//: Return null when serialization fails.
	if err != nil {
		//: Return null string on marshal error.
		return types.StringNull()
	}

	//: Store the string conversion once to avoid repeated allocations.
	s := string(data)
	//: Return null for empty array.
	if s == "null" || s == "[]" {
		//: Return null string for empty imports.
		return types.StringNull()
	}

	//: Return the serialized JSON string.
	return types.StringValue(s)
}

// mapImportResultToResource maps an SDK ImportResult to the resource model.
//
// Params:
//   - result: The n8n SDK ImportResult to map
//   - data: The models.Resource to populate
func mapImportResultToResource(result *n8nsdk.ImportResult, data *models.Resource) {
	//: Map workflows import result.
	if result.Workflows != nil {
		data.WorkflowsImported = marshalImportedSlice(result.Workflows)
	} else {
		//: Set null when no workflows were imported.
		data.WorkflowsImported = types.StringNull()
	}

	//: Map credentials import result.
	if result.Credentials != nil {
		data.CredentialsImported = marshalImportedSlice(result.Credentials)
	} else {
		//: Set null when no credentials were imported.
		data.CredentialsImported = types.StringNull()
	}

	//: Map tags import result.
	if result.Tags != nil {
		data.TagsImported = marshalImportedSlice(result.Tags)
	} else {
		//: Set null when no tags were imported.
		data.TagsImported = types.StringNull()
	}

	//: Map variables import result.
	if result.Variables != nil {
		data.VariablesImported = marshalImportedSlice(result.Variables)
	} else {
		//: Set null when no variables were imported.
		data.VariablesImported = types.StringNull()
	}
}
