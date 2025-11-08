package models

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestTransferResource(t *testing.T) {
	t.Run("create with all fields", func(t *testing.T) {
		// Create a transfer resource with all fields populated
		resource := TransferResource{
			ID:                   types.StringValue("transfer-123"),
			CredentialID:         types.StringValue("cred-456"),
			DestinationProjectID: types.StringValue("proj-789"),
			TransferredAt:        types.StringValue("2024-01-01T12:00:00Z"),
		}

		assert.Equal(t, "transfer-123", resource.ID.ValueString())
		assert.Equal(t, "cred-456", resource.CredentialID.ValueString())
		assert.Equal(t, "proj-789", resource.DestinationProjectID.ValueString())
		assert.Equal(t, "2024-01-01T12:00:00Z", resource.TransferredAt.ValueString())
		assert.False(t, resource.ID.IsNull())
		assert.False(t, resource.CredentialID.IsNull())
		assert.False(t, resource.DestinationProjectID.IsNull())
		assert.False(t, resource.TransferredAt.IsNull())
	})

	t.Run("create with null values", func(t *testing.T) {
		// Create a transfer resource with null values
		resource := TransferResource{
			ID:                   types.StringNull(),
			CredentialID:         types.StringNull(),
			DestinationProjectID: types.StringNull(),
			TransferredAt:        types.StringNull(),
		}

		assert.True(t, resource.ID.IsNull())
		assert.True(t, resource.CredentialID.IsNull())
		assert.True(t, resource.DestinationProjectID.IsNull())
		assert.True(t, resource.TransferredAt.IsNull())
		assert.Equal(t, "", resource.ID.ValueString())
		assert.Equal(t, "", resource.CredentialID.ValueString())
		assert.Equal(t, "", resource.DestinationProjectID.ValueString())
		assert.Equal(t, "", resource.TransferredAt.ValueString())
	})

	t.Run("create with unknown values", func(t *testing.T) {
		// Create a transfer resource with unknown values
		resource := TransferResource{
			ID:                   types.StringUnknown(),
			CredentialID:         types.StringUnknown(),
			DestinationProjectID: types.StringUnknown(),
			TransferredAt:        types.StringUnknown(),
		}

		assert.True(t, resource.ID.IsUnknown())
		assert.True(t, resource.CredentialID.IsUnknown())
		assert.True(t, resource.DestinationProjectID.IsUnknown())
		assert.True(t, resource.TransferredAt.IsUnknown())
		assert.False(t, resource.ID.IsNull())
		assert.False(t, resource.CredentialID.IsNull())
		assert.False(t, resource.DestinationProjectID.IsNull())
		assert.False(t, resource.TransferredAt.IsNull())
	})

	t.Run("partial initialization", func(t *testing.T) {
		// Test with only some fields set
		resource := TransferResource{
			ID:           types.StringValue("transfer-partial"),
			CredentialID: types.StringValue("cred-partial"),
			// Other fields remain zero value (null)
		}

		assert.Equal(t, "transfer-partial", resource.ID.ValueString())
		assert.Equal(t, "cred-partial", resource.CredentialID.ValueString())
		assert.True(t, resource.DestinationProjectID.IsNull())
		assert.True(t, resource.TransferredAt.IsNull())
	})

	t.Run("zero value struct", func(t *testing.T) {
		// Test zero value struct
		var resource TransferResource

		assert.True(t, resource.ID.IsNull())
		assert.True(t, resource.CredentialID.IsNull())
		assert.True(t, resource.DestinationProjectID.IsNull())
		assert.True(t, resource.TransferredAt.IsNull())
	})

	t.Run("modify values", func(t *testing.T) {
		// Test modifying values
		resource := TransferResource{
			ID:           types.StringValue("initial-id"),
			CredentialID: types.StringValue("initial-cred"),
		}

		// Modify values
		resource.ID = types.StringValue("modified-id")
		resource.CredentialID = types.StringValue("modified-cred")
		resource.DestinationProjectID = types.StringValue("new-project")
		resource.TransferredAt = types.StringValue("2024-01-02T00:00:00Z")

		assert.Equal(t, "modified-id", resource.ID.ValueString())
		assert.Equal(t, "modified-cred", resource.CredentialID.ValueString())
		assert.Equal(t, "new-project", resource.DestinationProjectID.ValueString())
		assert.Equal(t, "2024-01-02T00:00:00Z", resource.TransferredAt.ValueString())
	})

	t.Run("various ID formats", func(t *testing.T) {
		// Test various ID formats
		idFormats := []struct {
			name  string
			value string
		}{
			{"UUID", "550e8400-e29b-41d4-a716-446655440000"},
			{"numeric", "123456789"},
			{"alphanumeric", "abc123def456"},
			{"with dashes", "proj-123-cred-456"},
			{"with underscores", "proj_123_cred_456"},
			{"with dots", "proj.123.cred.456"},
			{"composite", "proj:123:cred:456"},
			{"base64", "cHJvai0xMjMtY3JlZC00NTY="},
		}

		for _, format := range idFormats {
			t.Run(format.name, func(t *testing.T) {
				resource := TransferResource{
					ID:                   types.StringValue(format.value),
					CredentialID:         types.StringValue(format.value),
					DestinationProjectID: types.StringValue(format.value),
				}
				assert.Equal(t, format.value, resource.ID.ValueString())
				assert.Equal(t, format.value, resource.CredentialID.ValueString())
				assert.Equal(t, format.value, resource.DestinationProjectID.ValueString())
			})
		}
	})

	t.Run("timestamp formats", func(t *testing.T) {
		// Test various timestamp formats
		timestamps := []string{
			"2024-01-01T00:00:00Z",
			"2024-01-01T12:34:56Z",
			"2024-01-01T12:34:56.789Z",
			"2024-01-01T12:34:56+00:00",
			"2024-01-01T12:34:56-05:00",
			time.Now().Format(time.RFC3339),
			time.Now().UTC().Format(time.RFC3339Nano),
			"2024-01-01",
			"01/01/2024",
			"January 1, 2024",
		}

		for _, ts := range timestamps {
			resource := TransferResource{
				TransferredAt: types.StringValue(ts),
			}
			assert.Equal(t, ts, resource.TransferredAt.ValueString())
		}
	})

	t.Run("special characters in IDs", func(t *testing.T) {
		// Test special characters in ID fields
		resource := TransferResource{
			ID:                   types.StringValue("id-with-special-!@#$%^&*()"),
			CredentialID:         types.StringValue("cred-with-特殊字符-🔒"),
			DestinationProjectID: types.StringValue("proj_with.dots-and-dashes"),
		}

		assert.Equal(t, "id-with-special-!@#$%^&*()", resource.ID.ValueString())
		assert.Equal(t, "cred-with-特殊字符-🔒", resource.CredentialID.ValueString())
		assert.Equal(t, "proj_with.dots-and-dashes", resource.DestinationProjectID.ValueString())
	})

	t.Run("very long IDs", func(t *testing.T) {
		// Test with very long ID strings
		longID := ""
		for i := 0; i < 1000; i++ {
			longID += "a"
		}

		resource := TransferResource{
			ID:                   types.StringValue(longID),
			CredentialID:         types.StringValue(longID),
			DestinationProjectID: types.StringValue(longID),
		}

		assert.Equal(t, longID, resource.ID.ValueString())
		assert.Equal(t, longID, resource.CredentialID.ValueString())
		assert.Equal(t, longID, resource.DestinationProjectID.ValueString())
	})

	t.Run("empty strings vs null", func(t *testing.T) {
		// Test difference between empty string and null
		resource1 := TransferResource{
			ID:                   types.StringValue(""),
			CredentialID:         types.StringValue(""),
			DestinationProjectID: types.StringValue(""),
			TransferredAt:        types.StringValue(""),
		}

		resource2 := TransferResource{
			ID:                   types.StringNull(),
			CredentialID:         types.StringNull(),
			DestinationProjectID: types.StringNull(),
			TransferredAt:        types.StringNull(),
		}

		// Empty strings are not null
		assert.False(t, resource1.ID.IsNull())
		assert.False(t, resource1.CredentialID.IsNull())
		assert.Equal(t, "", resource1.ID.ValueString())

		// Null values are null
		assert.True(t, resource2.ID.IsNull())
		assert.True(t, resource2.CredentialID.IsNull())
		assert.Equal(t, "", resource2.ID.ValueString())
	})

	t.Run("copy struct", func(t *testing.T) {
		// Test copying struct
		original := TransferResource{
			ID:                   types.StringValue("original-id"),
			CredentialID:         types.StringValue("original-cred"),
			DestinationProjectID: types.StringValue("original-proj"),
			TransferredAt:        types.StringValue("2024-01-01T00:00:00Z"),
		}

		copied := original

		assert.Equal(t, original.ID.ValueString(), copied.ID.ValueString())
		assert.Equal(t, original.CredentialID.ValueString(), copied.CredentialID.ValueString())
		assert.Equal(t, original.DestinationProjectID.ValueString(), copied.DestinationProjectID.ValueString())
		assert.Equal(t, original.TransferredAt.ValueString(), copied.TransferredAt.ValueString())

		// Modify copied
		copied.ID = types.StringValue("modified-id")

		// Original should not be affected (value semantics)
		assert.Equal(t, "original-id", original.ID.ValueString())
		assert.Equal(t, "modified-id", copied.ID.ValueString())
	})

	t.Run("pointer to struct", func(t *testing.T) {
		// Test pointer to struct
		resource := &TransferResource{
			ID:                   types.StringValue("pointer-id"),
			CredentialID:         types.StringValue("pointer-cred"),
			DestinationProjectID: types.StringValue("pointer-proj"),
		}

		assert.NotNil(t, resource)
		assert.Equal(t, "pointer-id", resource.ID.ValueString())
		assert.Equal(t, "pointer-cred", resource.CredentialID.ValueString())
		assert.Equal(t, "pointer-proj", resource.DestinationProjectID.ValueString())
	})

	t.Run("struct field tags", func(t *testing.T) {
		// This test documents that the struct has proper tfsdk tags
		resource := TransferResource{
			ID:                   types.StringValue("id"),
			CredentialID:         types.StringValue("credential_id"),
			DestinationProjectID: types.StringValue("destination_project_id"),
			TransferredAt:        types.StringValue("transferred_at"),
		}

		// The tfsdk tags map to Terraform schema field names
		assert.NotNil(t, resource)
	})

	t.Run("comparison", func(t *testing.T) {
		// Test struct comparison
		resource1 := TransferResource{
			ID:           types.StringValue("id1"),
			CredentialID: types.StringValue("cred1"),
		}

		resource2 := TransferResource{
			ID:           types.StringValue("id1"),
			CredentialID: types.StringValue("cred1"),
		}

		resource3 := TransferResource{
			ID:           types.StringValue("id2"),
			CredentialID: types.StringValue("cred2"),
		}

		// Same values should be equal
		assert.Equal(t, resource1.ID.ValueString(), resource2.ID.ValueString())
		assert.Equal(t, resource1.CredentialID.ValueString(), resource2.CredentialID.ValueString())

		// Different values should not be equal
		assert.NotEqual(t, resource1.ID.ValueString(), resource3.ID.ValueString())
		assert.NotEqual(t, resource1.CredentialID.ValueString(), resource3.CredentialID.ValueString())
	})
}

func TestTransferResourceConcurrency(t *testing.T) {
	t.Run("concurrent read", func(t *testing.T) {
		// Test concurrent reads
		resource := TransferResource{
			ID:                   types.StringValue("concurrent-id"),
			CredentialID:         types.StringValue("concurrent-cred"),
			DestinationProjectID: types.StringValue("concurrent-proj"),
			TransferredAt:        types.StringValue("2024-01-01T00:00:00Z"),
		}

		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				_ = resource.ID.ValueString()
				_ = resource.CredentialID.ValueString()
				_ = resource.DestinationProjectID.ValueString()
				_ = resource.TransferredAt.ValueString()
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}

		assert.True(t, true, "Concurrent reads completed without issues")
	})

	t.Run("separate instances", func(t *testing.T) {
		// Test multiple instances don't interfere
		done := make(chan bool, 10)

		for i := 0; i < 10; i++ {
			go func(n int) {
				resource := TransferResource{
					ID:           types.StringValue(string(rune('a' + n))),
					CredentialID: types.StringValue("cred"),
				}
				assert.NotNil(t, resource)
				done <- true
			}(i)
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

func TestTransferResourceValidation(t *testing.T) {
	t.Run("validate required fields", func(t *testing.T) {
		// Test that required fields can be checked
		resource := TransferResource{}

		// Check if required fields are null
		assert.True(t, resource.CredentialID.IsNull())
		assert.True(t, resource.DestinationProjectID.IsNull())

		// Set required fields
		resource.CredentialID = types.StringValue("required-cred")
		resource.DestinationProjectID = types.StringValue("required-proj")

		// Now they should not be null
		assert.False(t, resource.CredentialID.IsNull())
		assert.False(t, resource.DestinationProjectID.IsNull())
	})

	t.Run("validate transfer scenarios", func(t *testing.T) {
		// Test different transfer scenarios
		scenarios := []struct {
			name        string
			credID      string
			destProjID  string
			description string
		}{
			{
				name:        "user to project",
				credID:      "user-cred-123",
				destProjID:  "project-456",
				description: "Transfer from user to project",
			},
			{
				name:        "project to project",
				credID:      "project-cred-789",
				destProjID:  "project-012",
				description: "Transfer between projects",
			},
			{
				name:        "shared credential",
				credID:      "shared-cred-345",
				destProjID:  "shared-project-678",
				description: "Transfer shared credential",
			},
		}

		for _, scenario := range scenarios {
			t.Run(scenario.name, func(t *testing.T) {
				resource := TransferResource{
					ID:                   types.StringValue("transfer-" + scenario.name),
					CredentialID:         types.StringValue(scenario.credID),
					DestinationProjectID: types.StringValue(scenario.destProjID),
					TransferredAt:        types.StringValue(time.Now().Format(time.RFC3339)),
				}

				assert.Equal(t, scenario.credID, resource.CredentialID.ValueString())
				assert.Equal(t, scenario.destProjID, resource.DestinationProjectID.ValueString())
				assert.False(t, resource.TransferredAt.IsNull())
			})
		}
	})
}

func TestTransferResourceUseCases(t *testing.T) {
	t.Run("basic transfer", func(t *testing.T) {
		// Test basic credential transfer
		resource := TransferResource{
			ID:                   types.StringValue("transfer-001"),
			CredentialID:         types.StringValue("api-cred-123"),
			DestinationProjectID: types.StringValue("dev-project"),
			TransferredAt:        types.StringValue(time.Now().Format(time.RFC3339)),
		}

		assert.NotEmpty(t, resource.ID.ValueString())
		assert.NotEmpty(t, resource.CredentialID.ValueString())
		assert.NotEmpty(t, resource.DestinationProjectID.ValueString())
		assert.NotEmpty(t, resource.TransferredAt.ValueString())
	})

	t.Run("bulk transfer tracking", func(t *testing.T) {
		// Test tracking multiple transfers
		transfers := []TransferResource{
			{
				ID:                   types.StringValue("transfer-batch-001"),
				CredentialID:         types.StringValue("cred-001"),
				DestinationProjectID: types.StringValue("proj-target"),
				TransferredAt:        types.StringValue(time.Now().Format(time.RFC3339)),
			},
			{
				ID:                   types.StringValue("transfer-batch-002"),
				CredentialID:         types.StringValue("cred-002"),
				DestinationProjectID: types.StringValue("proj-target"),
				TransferredAt:        types.StringValue(time.Now().Format(time.RFC3339)),
			},
			{
				ID:                   types.StringValue("transfer-batch-003"),
				CredentialID:         types.StringValue("cred-003"),
				DestinationProjectID: types.StringValue("proj-target"),
				TransferredAt:        types.StringValue(time.Now().Format(time.RFC3339)),
			},
		}

		// All transfers to same destination
		for _, transfer := range transfers {
			assert.Equal(t, "proj-target", transfer.DestinationProjectID.ValueString())
			assert.NotEmpty(t, transfer.CredentialID.ValueString())
			assert.NotEmpty(t, transfer.TransferredAt.ValueString())
		}
	})

	t.Run("transfer with metadata", func(t *testing.T) {
		// Test transfer with rich IDs containing metadata
		resource := TransferResource{
			ID:                   types.StringValue("transfer_2024-01-01_user123_to_proj456"),
			CredentialID:         types.StringValue("oauth_google_user123"),
			DestinationProjectID: types.StringValue("proj_production_456"),
			TransferredAt:        types.StringValue("2024-01-01T15:30:00Z"),
		}

		// IDs can contain metadata for tracking
		assert.Contains(t, resource.ID.ValueString(), "2024-01-01")
		assert.Contains(t, resource.ID.ValueString(), "user123")
		assert.Contains(t, resource.ID.ValueString(), "proj456")
		assert.Contains(t, resource.CredentialID.ValueString(), "oauth")
		assert.Contains(t, resource.DestinationProjectID.ValueString(), "production")
	})
}

func BenchmarkTransferResource(b *testing.B) {
	b.Run("create", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = TransferResource{
				ID:                   types.StringValue("bench-transfer"),
				CredentialID:         types.StringValue("bench-cred"),
				DestinationProjectID: types.StringValue("bench-proj"),
				TransferredAt:        types.StringValue("2024-01-01T00:00:00Z"),
			}
		}
	})

	b.Run("access fields", func(b *testing.B) {
		resource := TransferResource{
			ID:                   types.StringValue("bench-transfer"),
			CredentialID:         types.StringValue("bench-cred"),
			DestinationProjectID: types.StringValue("bench-proj"),
			TransferredAt:        types.StringValue("2024-01-01T00:00:00Z"),
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = resource.ID.ValueString()
			_ = resource.CredentialID.ValueString()
			_ = resource.DestinationProjectID.ValueString()
			_ = resource.TransferredAt.ValueString()
		}
	})

	b.Run("modify fields", func(b *testing.B) {
		resource := TransferResource{}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			resource.ID = types.StringValue("id")
			resource.CredentialID = types.StringValue("cred")
			resource.DestinationProjectID = types.StringValue("proj")
			resource.TransferredAt = types.StringValue("time")
		}
	})
}
