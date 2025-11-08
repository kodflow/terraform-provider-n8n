package models

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestUserResource(t *testing.T) {
	t.Run("create with all fields", func(t *testing.T) {
		userResource := UserResource{
			ID:        types.StringValue("user-res-123"),
			ProjectID: types.StringValue("proj-456"),
			UserID:    types.StringValue("user-789"),
			Role:      types.StringValue("admin"),
		}

		assert.Equal(t, "user-res-123", userResource.ID.ValueString())
		assert.Equal(t, "proj-456", userResource.ProjectID.ValueString())
		assert.Equal(t, "user-789", userResource.UserID.ValueString())
		assert.Equal(t, "admin", userResource.Role.ValueString())
	})

	t.Run("create with null values", func(t *testing.T) {
		userResource := UserResource{
			ID:        types.StringNull(),
			ProjectID: types.StringNull(),
			UserID:    types.StringNull(),
			Role:      types.StringNull(),
		}

		assert.True(t, userResource.ID.IsNull())
		assert.True(t, userResource.ProjectID.IsNull())
		assert.True(t, userResource.UserID.IsNull())
		assert.True(t, userResource.Role.IsNull())
	})

	t.Run("create with unknown values", func(t *testing.T) {
		userResource := UserResource{
			ID:        types.StringUnknown(),
			ProjectID: types.StringUnknown(),
			UserID:    types.StringUnknown(),
			Role:      types.StringUnknown(),
		}

		assert.True(t, userResource.ID.IsUnknown())
		assert.True(t, userResource.ProjectID.IsUnknown())
		assert.True(t, userResource.UserID.IsUnknown())
		assert.True(t, userResource.Role.IsUnknown())
	})

	t.Run("user roles", func(t *testing.T) {
		roles := []string{
			"admin",
			"editor",
			"viewer",
			"owner",
			"member",
			"contributor",
			"guest",
		}

		for _, role := range roles {
			userResource := UserResource{
				Role: types.StringValue(role),
			}
			assert.Equal(t, role, userResource.Role.ValueString())
		}
	})

	t.Run("zero value struct", func(t *testing.T) {
		var userResource UserResource
		assert.True(t, userResource.ID.IsNull())
		assert.True(t, userResource.ProjectID.IsNull())
		assert.True(t, userResource.UserID.IsNull())
		assert.True(t, userResource.Role.IsNull())
	})

	t.Run("copy struct", func(t *testing.T) {
		original := UserResource{
			ID:        types.StringValue("original-id"),
			ProjectID: types.StringValue("original-proj"),
			UserID:    types.StringValue("original-user"),
			Role:      types.StringValue("admin"),
		}

		copied := original

		assert.Equal(t, original.ID.ValueString(), copied.ID.ValueString())
		assert.Equal(t, original.ProjectID.ValueString(), copied.ProjectID.ValueString())
		assert.Equal(t, original.UserID.ValueString(), copied.UserID.ValueString())
		assert.Equal(t, original.Role.ValueString(), copied.Role.ValueString())

		// Modify copied
		copied.ID = types.StringValue("modified-id")
		copied.Role = types.StringValue("viewer")
		assert.Equal(t, "original-id", original.ID.ValueString())
		assert.Equal(t, "modified-id", copied.ID.ValueString())
		assert.Equal(t, "admin", original.Role.ValueString())
		assert.Equal(t, "viewer", copied.Role.ValueString())
	})

	t.Run("partial initialization", func(t *testing.T) {
		userResource := UserResource{
			ProjectID: types.StringValue("proj-partial"),
			UserID:    types.StringValue("user-partial"),
			// Other fields remain null
		}

		assert.True(t, userResource.ID.IsNull())
		assert.Equal(t, "proj-partial", userResource.ProjectID.ValueString())
		assert.Equal(t, "user-partial", userResource.UserID.ValueString())
		assert.True(t, userResource.Role.IsNull())
	})

	t.Run("id formats", func(t *testing.T) {
		ids := []string{
			"simple-id",
			"123456",
			"uuid-550e8400-e29b-41d4-a716-446655440000",
			"ID_WITH_UNDERSCORES",
			"id.with.dots",
			"id/with/slashes",
			"composite:proj-123:user-456",
		}

		for _, id := range ids {
			userResource := UserResource{
				ID: types.StringValue(id),
			}
			assert.Equal(t, id, userResource.ID.ValueString())
		}
	})

	t.Run("project and user id relationships", func(t *testing.T) {
		userResource := UserResource{
			ID:        types.StringValue("rel-123"),
			ProjectID: types.StringValue("proj-abc"),
			UserID:    types.StringValue("user-xyz"),
			Role:      types.StringValue("member"),
		}

		// Verify IDs are independent
		assert.NotEqual(t, userResource.ID.ValueString(), userResource.ProjectID.ValueString())
		assert.NotEqual(t, userResource.ID.ValueString(), userResource.UserID.ValueString())
		assert.NotEqual(t, userResource.ProjectID.ValueString(), userResource.UserID.ValueString())
	})

	t.Run("multiple users in same project", func(t *testing.T) {
		projectID := "shared-project"

		user1 := UserResource{
			ID:        types.StringValue("rel-1"),
			ProjectID: types.StringValue(projectID),
			UserID:    types.StringValue("user-1"),
			Role:      types.StringValue("admin"),
		}

		user2 := UserResource{
			ID:        types.StringValue("rel-2"),
			ProjectID: types.StringValue(projectID),
			UserID:    types.StringValue("user-2"),
			Role:      types.StringValue("editor"),
		}

		user3 := UserResource{
			ID:        types.StringValue("rel-3"),
			ProjectID: types.StringValue(projectID),
			UserID:    types.StringValue("user-3"),
			Role:      types.StringValue("viewer"),
		}

		// All share same project
		assert.Equal(t, projectID, user1.ProjectID.ValueString())
		assert.Equal(t, projectID, user2.ProjectID.ValueString())
		assert.Equal(t, projectID, user3.ProjectID.ValueString())

		// But have different user IDs and roles
		assert.NotEqual(t, user1.UserID.ValueString(), user2.UserID.ValueString())
		assert.NotEqual(t, user1.UserID.ValueString(), user3.UserID.ValueString())
		assert.NotEqual(t, user2.UserID.ValueString(), user3.UserID.ValueString())

		assert.NotEqual(t, user1.Role.ValueString(), user2.Role.ValueString())
		assert.NotEqual(t, user1.Role.ValueString(), user3.Role.ValueString())
		assert.NotEqual(t, user2.Role.ValueString(), user3.Role.ValueString())
	})

	t.Run("same user in multiple projects", func(t *testing.T) {
		userID := "shared-user"

		proj1 := UserResource{
			ID:        types.StringValue("rel-1"),
			ProjectID: types.StringValue("proj-1"),
			UserID:    types.StringValue(userID),
			Role:      types.StringValue("admin"),
		}

		proj2 := UserResource{
			ID:        types.StringValue("rel-2"),
			ProjectID: types.StringValue("proj-2"),
			UserID:    types.StringValue(userID),
			Role:      types.StringValue("editor"),
		}

		proj3 := UserResource{
			ID:        types.StringValue("rel-3"),
			ProjectID: types.StringValue("proj-3"),
			UserID:    types.StringValue(userID),
			Role:      types.StringValue("viewer"),
		}

		// All share same user
		assert.Equal(t, userID, proj1.UserID.ValueString())
		assert.Equal(t, userID, proj2.UserID.ValueString())
		assert.Equal(t, userID, proj3.UserID.ValueString())

		// But have different projects and potentially different roles
		assert.NotEqual(t, proj1.ProjectID.ValueString(), proj2.ProjectID.ValueString())
		assert.NotEqual(t, proj1.ProjectID.ValueString(), proj3.ProjectID.ValueString())
		assert.NotEqual(t, proj2.ProjectID.ValueString(), proj3.ProjectID.ValueString())
	})
}

func TestUserResourceConcurrency(t *testing.T) {
	t.Run("concurrent read", func(t *testing.T) {
		userResource := UserResource{
			ID:        types.StringValue("concurrent-id"),
			ProjectID: types.StringValue("concurrent-proj"),
			UserID:    types.StringValue("concurrent-user"),
			Role:      types.StringValue("admin"),
		}

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				_ = userResource.ID.ValueString()
				_ = userResource.ProjectID.ValueString()
				_ = userResource.UserID.ValueString()
				_ = userResource.Role.ValueString()
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func BenchmarkUserResource(b *testing.B) {
	b.Run("create", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = UserResource{
				ID:        types.StringValue("user-res-123"),
				ProjectID: types.StringValue("proj-456"),
				UserID:    types.StringValue("user-789"),
				Role:      types.StringValue("admin"),
			}
		}
	})

	b.Run("access fields", func(b *testing.B) {
		userResource := UserResource{
			ID:        types.StringValue("user-res-123"),
			ProjectID: types.StringValue("proj-456"),
			UserID:    types.StringValue("user-789"),
			Role:      types.StringValue("admin"),
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = userResource.ID.ValueString()
			_ = userResource.ProjectID.ValueString()
			_ = userResource.UserID.ValueString()
			_ = userResource.Role.ValueString()
		}
	})
}
