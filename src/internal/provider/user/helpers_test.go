package user

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/user/models"
	"github.com/stretchr/testify/assert"
)

func TestMapUserToItem(t *testing.T) {
	t.Run("map with all fields populated", func(t *testing.T) {
		id := "user123"
		email := "test@example.com"
		firstName := "John"
		lastName := "Doe"
		isPending := false
		role := "global:admin"
		createdAt := time.Now()
		updatedAt := time.Now().Add(1 * time.Hour)

		user := &n8nsdk.User{
			Id:        &id,
			Email:     email,
			FirstName: &firstName,
			LastName:  &lastName,
			IsPending: &isPending,
			Role:      &role,
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
		}

		item := mapUserToItem(user)

		assert.Equal(t, "user123", item.ID.ValueString())
		assert.Equal(t, "test@example.com", item.Email.ValueString())
		assert.Equal(t, "John", item.FirstName.ValueString())
		assert.Equal(t, "Doe", item.LastName.ValueString())
		assert.False(t, item.IsPending.ValueBool())
		assert.Equal(t, "global:admin", item.Role.ValueString())
		assert.Equal(t, createdAt.String(), item.CreatedAt.ValueString())
		assert.Equal(t, updatedAt.String(), item.UpdatedAt.ValueString())
	})

	t.Run("map with nil optional fields", func(t *testing.T) {
		email := "test@example.com"

		user := &n8nsdk.User{
			Email: email,
		}

		item := mapUserToItem(user)

		assert.True(t, item.ID.IsNull())
		assert.Equal(t, "test@example.com", item.Email.ValueString())
		assert.True(t, item.FirstName.IsNull())
		assert.True(t, item.LastName.IsNull())
		assert.True(t, item.IsPending.IsNull())
		assert.True(t, item.Role.IsNull())
		assert.True(t, item.CreatedAt.IsNull())
		assert.True(t, item.UpdatedAt.IsNull())
	})

	t.Run("map with partial fields", func(t *testing.T) {
		id := "user456"
		email := "partial@example.com"
		role := "global:member"

		user := &n8nsdk.User{
			Id:    &id,
			Email: email,
			Role:  &role,
		}

		item := mapUserToItem(user)

		assert.Equal(t, "user456", item.ID.ValueString())
		assert.Equal(t, "partial@example.com", item.Email.ValueString())
		assert.Equal(t, "global:member", item.Role.ValueString())
		assert.True(t, item.FirstName.IsNull())
		assert.True(t, item.LastName.IsNull())
		assert.True(t, item.IsPending.IsNull())
	})

	t.Run("map with empty string values", func(t *testing.T) {
		email := ""
		firstName := ""
		lastName := ""

		user := &n8nsdk.User{
			Email:     email,
			FirstName: &firstName,
			LastName:  &lastName,
		}

		item := mapUserToItem(user)

		assert.Equal(t, "", item.Email.ValueString())
		assert.Equal(t, "", item.FirstName.ValueString())
		assert.Equal(t, "", item.LastName.ValueString())
	})

	t.Run("map different roles", func(t *testing.T) {
		roles := []string{"global:owner", "global:admin", "global:member"}

		for _, role := range roles {
			roleCopy := role
			user := &n8nsdk.User{
				Email: "test@example.com",
				Role:  &roleCopy,
			}

			item := mapUserToItem(user)

			assert.Equal(t, role, item.Role.ValueString())
		}
	})

	t.Run("map pending states", func(t *testing.T) {
		// Test pending user
		pending := true
		pendingUser := &n8nsdk.User{
			Email:     "pending@example.com",
			IsPending: &pending,
		}

		pendingItem := mapUserToItem(pendingUser)
		assert.True(t, pendingItem.IsPending.ValueBool())

		// Test active user
		active := false
		activeUser := &n8nsdk.User{
			Email:     "active@example.com",
			IsPending: &active,
		}

		activeItem := mapUserToItem(activeUser)
		assert.False(t, activeItem.IsPending.ValueBool())
	})

	t.Run("map with special characters in names", func(t *testing.T) {
		firstName := "José"
		lastName := "O'Brien"
		email := "jose.obrien@example.com"

		user := &n8nsdk.User{
			Email:     email,
			FirstName: &firstName,
			LastName:  &lastName,
		}

		item := mapUserToItem(user)

		assert.Equal(t, firstName, item.FirstName.ValueString())
		assert.Equal(t, lastName, item.LastName.ValueString())
	})

	t.Run("map timestamps at different times", func(t *testing.T) {
		createdAt := time.Now().Add(-24 * time.Hour)
		updatedAt := time.Now()

		user := &n8nsdk.User{
			Email:     "timestamps@example.com",
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
		}

		item := mapUserToItem(user)

		assert.Equal(t, createdAt.String(), item.CreatedAt.ValueString())
		assert.Equal(t, updatedAt.String(), item.UpdatedAt.ValueString())
	})
}

func TestMapUserToResourceModel(t *testing.T) {
	t.Run("map with all fields populated", func(t *testing.T) {
		id := "user789"
		email := "resource@example.com"
		firstName := "Alice"
		lastName := "Smith"
		isPending := false
		role := "global:owner"
		createdAt := time.Now()
		updatedAt := time.Now().Add(2 * time.Hour)

		user := &n8nsdk.User{
			Id:        &id,
			Email:     email,
			FirstName: &firstName,
			LastName:  &lastName,
			IsPending: &isPending,
			Role:      &role,
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
		}

		data := &models.Resource{}
		mapUserToResourceModel(user, data)

		assert.Equal(t, "user789", data.ID.ValueString())
		assert.Equal(t, "resource@example.com", data.Email.ValueString())
		assert.Equal(t, "Alice", data.FirstName.ValueString())
		assert.Equal(t, "Smith", data.LastName.ValueString())
		assert.False(t, data.IsPending.ValueBool())
		assert.Equal(t, "global:owner", data.Role.ValueString())
		assert.Equal(t, createdAt.String(), data.CreatedAt.ValueString())
		assert.Equal(t, updatedAt.String(), data.UpdatedAt.ValueString())
	})

	t.Run("map with nil optional fields", func(t *testing.T) {
		email := "minimal@example.com"

		user := &n8nsdk.User{
			Email: email,
		}

		data := &models.Resource{}
		mapUserToResourceModel(user, data)

		assert.True(t, data.ID.IsNull())
		assert.Equal(t, "minimal@example.com", data.Email.ValueString())
		assert.True(t, data.FirstName.IsNull())
		assert.True(t, data.LastName.IsNull())
		assert.True(t, data.IsPending.IsNull())
		assert.True(t, data.Role.IsNull())
		assert.True(t, data.CreatedAt.IsNull())
		assert.True(t, data.UpdatedAt.IsNull())
	})

	t.Run("map updates existing data model", func(t *testing.T) {
		id := "updated123"
		email := "updated@example.com"
		role := "global:admin"

		user := &n8nsdk.User{
			Id:    &id,
			Email: email,
			Role:  &role,
		}

		data := &models.Resource{
			Email: types.StringValue("old@example.com"),
			Role:  types.StringValue("global:member"),
		}

		mapUserToResourceModel(user, data)

		assert.Equal(t, "updated123", data.ID.ValueString())
		assert.Equal(t, "updated@example.com", data.Email.ValueString())
		assert.Equal(t, "global:admin", data.Role.ValueString())
	})

	t.Run("map preserves email when updating", func(t *testing.T) {
		id := "preserve123"
		email := "preserve@example.com"

		user := &n8nsdk.User{
			Id:    &id,
			Email: email,
		}

		data := &models.Resource{}
		mapUserToResourceModel(user, data)

		assert.Equal(t, "preserve@example.com", data.Email.ValueString())
	})

	t.Run("map with different role types", func(t *testing.T) {
		roles := []string{"global:owner", "global:admin", "global:member"}

		for _, role := range roles {
			roleCopy := role
			user := &n8nsdk.User{
				Email: "roles@example.com",
				Role:  &roleCopy,
			}

			data := &models.Resource{}
			mapUserToResourceModel(user, data)

			assert.Equal(t, role, data.Role.ValueString())
		}
	})

	t.Run("map pending and active states", func(t *testing.T) {
		// Test pending user
		pendingFlag := true
		pendingUser := &n8nsdk.User{
			Email:     "pending@example.com",
			IsPending: &pendingFlag,
		}

		pendingData := &models.Resource{}
		mapUserToResourceModel(pendingUser, pendingData)
		assert.True(t, pendingData.IsPending.ValueBool())

		// Test active user
		activeFlag := false
		activeUser := &n8nsdk.User{
			Email:     "active@example.com",
			IsPending: &activeFlag,
		}

		activeData := &models.Resource{}
		mapUserToResourceModel(activeUser, activeData)
		assert.False(t, activeData.IsPending.ValueBool())
	})

	t.Run("map overwrites previous values", func(t *testing.T) {
		id := "newid"
		email := "new@example.com"
		firstName := "New"
		lastName := "Name"

		user := &n8nsdk.User{
			Id:        &id,
			Email:     email,
			FirstName: &firstName,
			LastName:  &lastName,
		}

		data := &models.Resource{
			ID:        types.StringValue("oldid"),
			Email:     types.StringValue("old@example.com"),
			FirstName: types.StringValue("Old"),
			LastName:  types.StringValue("Name"),
		}

		mapUserToResourceModel(user, data)

		assert.Equal(t, "newid", data.ID.ValueString())
		assert.Equal(t, "new@example.com", data.Email.ValueString())
		assert.Equal(t, "New", data.FirstName.ValueString())
		assert.Equal(t, "Name", data.LastName.ValueString())
	})
}

func TestMapUserToItemConcurrency(t *testing.T) {
	t.Run("concurrent mapping", func(t *testing.T) {
		id := "concurrent123"
		email := "concurrent@example.com"
		firstName := "Concurrent"
		lastName := "Test"
		isPending := false
		role := "global:admin"

		user := &n8nsdk.User{
			Id:        &id,
			Email:     email,
			FirstName: &firstName,
			LastName:  &lastName,
			IsPending: &isPending,
			Role:      &role,
		}

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				item := mapUserToItem(user)
				assert.Equal(t, "concurrent123", item.ID.ValueString())
				assert.Equal(t, "concurrent@example.com", item.Email.ValueString())
				assert.Equal(t, "Concurrent", item.FirstName.ValueString())
				assert.Equal(t, "Test", item.LastName.ValueString())
				assert.False(t, item.IsPending.ValueBool())
				assert.Equal(t, "global:admin", item.Role.ValueString())
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func TestMapUserToResourceModelConcurrency(t *testing.T) {
	t.Run("concurrent resource mapping", func(t *testing.T) {
		id := "resourceconcurrent"
		email := "resourceconcurrent@example.com"
		role := "global:member"

		user := &n8nsdk.User{
			Id:    &id,
			Email: email,
			Role:  &role,
		}

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				data := &models.Resource{}
				mapUserToResourceModel(user, data)
				assert.Equal(t, "resourceconcurrent", data.ID.ValueString())
				assert.Equal(t, "resourceconcurrent@example.com", data.Email.ValueString())
				assert.Equal(t, "global:member", data.Role.ValueString())
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func BenchmarkMapUserToItem(b *testing.B) {
	id := "benchmark123"
	email := "benchmark@example.com"
	firstName := "Benchmark"
	lastName := "User"
	isPending := false
	role := "global:admin"
	createdAt := time.Now()
	updatedAt := time.Now().Add(1 * time.Hour)

	user := &n8nsdk.User{
		Id:        &id,
		Email:     email,
		FirstName: &firstName,
		LastName:  &lastName,
		IsPending: &isPending,
		Role:      &role,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapUserToItem(user)
	}
}

func BenchmarkMapUserToResourceModel(b *testing.B) {
	id := "benchmarkresource"
	email := "benchmarkresource@example.com"
	firstName := "Benchmark"
	lastName := "Resource"
	isPending := false
	role := "global:owner"
	createdAt := time.Now()
	updatedAt := time.Now().Add(2 * time.Hour)

	user := &n8nsdk.User{
		Id:        &id,
		Email:     email,
		FirstName: &firstName,
		LastName:  &lastName,
		IsPending: &isPending,
		Role:      &role,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}

	data := &models.Resource{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapUserToResourceModel(user, data)
	}
}

// Helper functions for tests
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func stringValue(s string) types.String {
	return types.StringValue(s)
}

func boolValue(b bool) types.Bool {
	return types.BoolValue(b)
}
