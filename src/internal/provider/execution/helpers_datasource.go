// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package execution implements execution management resources.
package execution

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/execution/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/constants"
)

// nullableTime is the minimal interface for nullable time values from the SDK.
// Using an interface instead of the concrete type keeps the helper testable and decoupled.
type nullableTime interface {
	//: Get returns the underlying time pointer.
	Get() *time.Time
	//: IsSet reports whether the value is set.
	IsSet() bool
}

// stringValue is the minimal interface for Terraform string attributes.
// Decouples isValueSet from the concrete basetypes.StringValue type.
type stringValue interface {
	//: IsNull reports whether the attribute is null.
	IsNull() bool
	//: IsUnknown reports whether the attribute is unknown.
	IsUnknown() bool
}

// isValueSet reports whether a string attribute is set (not null and not unknown).
// Used to guard filter application without repeating the double nil check.
//
// Params:
//   - v: Terraform string attribute to test
//
// Returns:
//   - set: True when the value should be used as a filter
func isValueSet(v stringValue) (set bool) {
	//: Return true only when the value carries a usable string.
	return !v.IsNull() && !v.IsUnknown()
}

// buildExecutionItems converts the API execution list into the Terraform model slice.
// Returns an empty slice when the list is nil or contains no data.
//
// Params:
//   - execList: The n8n SDK ExecutionList response
//
// Returns:
//   - []models.Item: The mapped item slice (never nil)
func buildExecutionItems(execList *n8nsdk.ExecutionList) (items []models.Item) {
	items = make([]models.Item, 0, constants.DefaultListCapacity)
	//: Populate only when the response contains data.
	if execList != nil && execList.Data != nil {
		//: Iterate over each execution and append it to the slice.
		for _, exec := range execList.Data {
			items = append(items, mapExecutionToItem(&exec))
		}
	}
	//: Return the populated or empty slice.
	return items
}

// formatExecutionID converts a float32 execution ID pointer to a Terraform string value.
// Returns types.StringNull() when the pointer is nil.
//
// Params:
//   - id: Optional float32 execution ID pointer
//
// Returns:
//   - types.String: The formatted ID string or null
func formatExecutionID(id *float32) (idStr types.String) {
	//: Return null when the ID pointer is absent.
	if id == nil {
		//: Return null string result.
		return types.StringNull()
	}
	//: Format the float32 ID as an integer string for readability.
	return types.StringValue(fmt.Sprintf("%.0f", *id))
}

// formatNullableTime converts a nullableTime to a Terraform string value.
// Returns types.StringNull() when the time is not set.
//
// Params:
//   - t: nullableTime from the SDK
//
// Returns:
//   - types.String: The formatted ISO 8601 time string or null
func formatNullableTime(t nullableTime) (timeStr types.String) {
	//: Return null when the time is not set.
	if !t.IsSet() || t.Get() == nil {
		//: Return null string result.
		return types.StringNull()
	}
	//: Format the time in ISO 8601 format.
	return types.StringValue(t.Get().Format("2006-01-02T15:04:05Z07:00"))
}

// formatNullableBool converts a bool pointer to a Terraform bool value.
// Returns types.BoolNull() when the pointer is nil.
//
// Params:
//   - b: Optional bool pointer
//
// Returns:
//   - types.Bool: The bool value or null
func formatNullableBool(b *bool) (boolVal types.Bool) {
	//: Return null when the bool pointer is absent.
	if b == nil {
		//: Return null bool result.
		return types.BoolNull()
	}
	//: Return the bool value.
	return types.BoolValue(*b)
}

// formatNullableString converts a string pointer to a Terraform string value.
// Returns types.StringNull() when the pointer is nil.
//
// Params:
//   - s: Optional string pointer
//
// Returns:
//   - types.String: The string value or null
func formatNullableString(s *string) (strVal types.String) {
	//: Return null when the string pointer is absent.
	if s == nil {
		//: Return null string result.
		return types.StringNull()
	}
	//: Return the string value.
	return types.StringValue(*s)
}

// mapExecutionToDataSource maps an SDK Execution to a DataSource model.
//
// Params:
//   - exec: The n8n SDK Execution to map
//   - data: The models.DataSource to populate
func mapExecutionToDataSource(exec *n8nsdk.Execution, data *models.DataSource) {
	data.ID = formatExecutionID(exec.Id)
	data.Mode = formatNullableString(exec.Mode)
	data.Status = formatNullableString(exec.Status)
	data.WorkflowID = formatExecutionID(exec.WorkflowId)
	data.Finished = formatNullableBool(exec.Finished)
	data.CreatedAt = formatNullableTime(exec.CreatedAt)
	data.StartedAt = formatNullableTime(exec.StartedAt)
	data.StoppedAt = formatNullableTime(exec.StoppedAt)
}

// mapExecutionToItem maps an SDK Execution to an Item model for list datasources.
//
// Params:
//   - exec: The n8n SDK Execution to map
//
// Returns:
//   - models.Item: The mapped item model
func mapExecutionToItem(exec *n8nsdk.Execution) (item models.Item) {
	item = models.Item{
		ID:         formatExecutionID(exec.Id),
		Mode:       formatNullableString(exec.Mode),
		Status:     formatNullableString(exec.Status),
		WorkflowID: formatExecutionID(exec.WorkflowId),
		Finished:   formatNullableBool(exec.Finished),
		CreatedAt:  formatNullableTime(exec.CreatedAt),
		StartedAt:  formatNullableTime(exec.StartedAt),
		StoppedAt:  formatNullableTime(exec.StoppedAt),
	}
	//: Return the fully populated item.
	return item
}
