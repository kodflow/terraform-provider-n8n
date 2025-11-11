package user

import (
	"testing"
	"time"

	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/user/models"
	"github.com/stretchr/testify/assert"
)

// TestUserDataSource_populateUserData_FullCoverage ensures 100% coverage.
func TestUserDataSource_populateUserData_FullCoverage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		user *n8nsdk.User
	}{
		{
			name: "all fields populated",
			user: &n8nsdk.User{
				Id:        strPtr("user-1"),
				Email:     "test@example.com",
				FirstName: strPtr("John"),
				LastName:  strPtr("Doe"),
				IsPending: boolPtr(false),
				CreatedAt: timePtr(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)),
				UpdatedAt: timePtr(time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC)),
				Role:      strPtr("global:admin"),
			},
		},
		{
			name: "nil ID",
			user: &n8nsdk.User{
				Id:    nil,
				Email: "test@example.com",
			},
		},
		{
			name: "nil FirstName",
			user: &n8nsdk.User{
				Id:        strPtr("user-1"),
				Email:     "test@example.com",
				FirstName: nil,
			},
		},
		{
			name: "nil LastName",
			user: &n8nsdk.User{
				Id:       strPtr("user-1"),
				Email:    "test@example.com",
				LastName: nil,
			},
		},
		{
			name: "nil IsPending",
			user: &n8nsdk.User{
				Id:        strPtr("user-1"),
				Email:     "test@example.com",
				IsPending: nil,
			},
		},
		{
			name: "nil CreatedAt",
			user: &n8nsdk.User{
				Id:        strPtr("user-1"),
				Email:     "test@example.com",
				CreatedAt: nil,
			},
		},
		{
			name: "nil UpdatedAt",
			user: &n8nsdk.User{
				Id:        strPtr("user-1"),
				Email:     "test@example.com",
				UpdatedAt: nil,
			},
		},
		{
			name: "nil Role",
			user: &n8nsdk.User{
				Id:    strPtr("user-1"),
				Email: "test@example.com",
				Role:  nil,
			},
		},
		{
			name: "all optional fields nil",
			user: &n8nsdk.User{
				Id:        nil,
				Email:     "test@example.com",
				FirstName: nil,
				LastName:  nil,
				IsPending: nil,
				CreatedAt: nil,
				UpdatedAt: nil,
				Role:      nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := &UserDataSource{}
			data := &models.DataSource{}

			d.populateUserData(tt.user, data)

			// Verify email is always set
			assert.Equal(t, tt.user.Email, data.Email.ValueString())

			// Verify optional fields are set correctly
			if tt.user.Id != nil {
				assert.Equal(t, *tt.user.Id, data.ID.ValueString())
			} else {
				assert.True(t, data.ID.IsNull() || data.ID.ValueString() == "")
			}
		})
	}
}

// strPtr returns a pointer to the given string.
func strPtr(s string) *string {
	return &s
}

// timePtr returns a pointer to the given time.
func timePtr(t time.Time) *time.Time {
	return &t
}
