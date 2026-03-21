// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package user contains helper functions for user operations.
package user

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/user/models"
)

// mapUserToItem maps an SDK user to the user item model.
//
// Params:
//   - user: The SDK user to map
//
// Returns:
//   - models.Item: The mapped user item model
func mapUserToItem(user *n8nsdk.User) (item models.Item) {
	//: Initialize model with pointer.
	result := &models.Item{
		Email: types.StringValue(user.Email),
	}

	//: Check for non-nil value.
	if user.Id != nil {
		result.ID = types.StringValue(*user.Id)
	}
	//: Check for non-nil value.
	if user.FirstName != nil {
		result.FirstName = types.StringPointerValue(user.FirstName)
	}
	//: Check for non-nil value.
	if user.LastName != nil {
		result.LastName = types.StringPointerValue(user.LastName)
	}
	//: Check for non-nil value.
	if user.IsPending != nil {
		result.IsPending = types.BoolPointerValue(user.IsPending)
	}
	//: Check for non-nil value.
	if user.CreatedAt != nil {
		result.CreatedAt = types.StringValue(user.CreatedAt.String())
	}
	//: Check for non-nil value.
	if user.UpdatedAt != nil {
		result.UpdatedAt = types.StringValue(user.UpdatedAt.String())
	}
	//: Check for non-nil value.
	if user.Role != nil {
		result.Role = types.StringPointerValue(user.Role)
	}

	//: Return result.
	return *result
}

// mapUserToResourceModel maps an SDK user to the user resource model.
// For pending users, firstName and lastName may be nil - we set them to Null
// to satisfy Terraform's requirement that all computed values be known after apply.
//
// Params:
//   - user: The SDK user to map
//   - data: The resource model to populate
func mapUserToResourceModel(user *n8nsdk.User, data *models.Resource) {
	//: Check for non-nil value.
	if user.Id != nil {
		data.ID = types.StringValue(*user.Id)
	}
	data.Email = types.StringValue(user.Email)
	//: FirstName may be nil for pending users.
	if user.FirstName != nil {
		data.FirstName = types.StringPointerValue(user.FirstName)
	} else {
		//: Set explicit Null instead of leaving unknown.
		data.FirstName = types.StringNull()
	}
	//: LastName may be nil for pending users.
	if user.LastName != nil {
		data.LastName = types.StringPointerValue(user.LastName)
	} else {
		//: Set explicit Null instead of leaving unknown.
		data.LastName = types.StringNull()
	}
	//: Check for non-nil value.
	if user.Role != nil {
		data.Role = types.StringPointerValue(user.Role)
	}
	//: Check for non-nil value.
	if user.IsPending != nil {
		data.IsPending = types.BoolPointerValue(user.IsPending)
	}
	//: Check for non-nil value.
	if user.CreatedAt != nil {
		data.CreatedAt = types.StringValue(user.CreatedAt.String())
	}
	//: Check for non-nil value.
	if user.UpdatedAt != nil {
		data.UpdatedAt = types.StringValue(user.UpdatedAt.String())
	}
}
