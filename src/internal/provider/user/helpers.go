// Package user contains helper functions for user operations.
package user

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/user/models"
)

// mapUserToItem maps an SDK user to the user item model.
//
// Params:
//   - user: The SDK user to map
//
// Returns:
//   - models.Item: The mapped user item model
func mapUserToItem(user *n8nsdk.User) models.Item {
	// Initialize model with pointer.
	item := &models.Item{
		Email: types.StringValue(user.Email),
	}

	// Check for non-nil value.
	if user.Id != nil {
		item.ID = types.StringValue(*user.Id)
	}
	// Check for non-nil value.
	if user.FirstName != nil {
		item.FirstName = types.StringPointerValue(user.FirstName)
	}
	// Check for non-nil value.
	if user.LastName != nil {
		item.LastName = types.StringPointerValue(user.LastName)
	}
	// Check for non-nil value.
	if user.IsPending != nil {
		item.IsPending = types.BoolPointerValue(user.IsPending)
	}
	// Check for non-nil value.
	if user.CreatedAt != nil {
		item.CreatedAt = types.StringValue(user.CreatedAt.String())
	}
	// Check for non-nil value.
	if user.UpdatedAt != nil {
		item.UpdatedAt = types.StringValue(user.UpdatedAt.String())
	}
	// Check for non-nil value.
	if user.Role != nil {
		item.Role = types.StringPointerValue(user.Role)
	}

	// Return result.
	return *item
}

// mapUserToResourceModel maps an SDK user to the user resource model.
//
// Params:
//   - user: The SDK user to map
//   - data: The resource model to populate
//
// Returns:
//   - error: Error if operation fails
func mapUserToResourceModel(user *n8nsdk.User, data *models.Resource) {
	// Check for non-nil value.
	if user.Id != nil {
		data.ID = types.StringValue(*user.Id)
	}
	data.Email = types.StringValue(user.Email)
	// Check for non-nil value.
	if user.FirstName != nil {
		data.FirstName = types.StringPointerValue(user.FirstName)
	}
	// Check for non-nil value.
	if user.LastName != nil {
		data.LastName = types.StringPointerValue(user.LastName)
	}
	// Check for non-nil value.
	if user.Role != nil {
		data.Role = types.StringPointerValue(user.Role)
	}
	// Check for non-nil value.
	if user.IsPending != nil {
		data.IsPending = types.BoolPointerValue(user.IsPending)
	}
	// Check for non-nil value.
	if user.CreatedAt != nil {
		data.CreatedAt = types.StringValue(user.CreatedAt.String())
	}
	// Check for non-nil value.
	if user.UpdatedAt != nil {
		data.UpdatedAt = types.StringValue(user.UpdatedAt.String())
	}
}
