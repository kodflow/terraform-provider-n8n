package credential

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/credential/models"
	"github.com/stretchr/testify/assert"
)

func TestWorkflowBackup(t *testing.T) {
	t.Run("create workflow backup", func(t *testing.T) {
		active := true
		id := "123"
		originalWorkflow := &n8nsdk.Workflow{
			Id:     &id,
			Name:   "Test Workflow",
			Active: &active,
		}

		backup := models.WorkflowBackup{
			ID:       "backup-123",
			Original: originalWorkflow,
		}

		assert.Equal(t, "backup-123", backup.ID)
		assert.NotNil(t, backup.Original)
		assert.Equal(t, "123", *backup.Original.Id)
		assert.Equal(t, "Test Workflow", backup.Original.Name)
		assert.True(t, *backup.Original.Active)
	})

	t.Run("nil original workflow", func(t *testing.T) {
		backup := models.WorkflowBackup{
			ID:       "backup-nil",
			Original: nil,
		}

		assert.Equal(t, "backup-nil", backup.ID)
		assert.Nil(t, backup.Original)
	})

	t.Run("empty backup ID", func(t *testing.T) {
		backup := models.WorkflowBackup{
			ID:       "",
			Original: &n8nsdk.Workflow{},
		}

		assert.Equal(t, "", backup.ID)
		assert.NotNil(t, backup.Original)
	})

	t.Run("zero value struct", func(t *testing.T) {
		var backup models.WorkflowBackup
		assert.Equal(t, "", backup.ID)
		assert.Nil(t, backup.Original)
	})
}

func TestTransferResource(t *testing.T) {
	t.Run("create transfer resource with all fields", func(t *testing.T) {
		transfer := models.TransferResource{
			ID:                   types.StringValue("transfer-123"),
			CredentialID:         types.StringValue("cred-456"),
			DestinationProjectID: types.StringValue("proj-789"),
			TransferredAt:        types.StringValue("2024-01-01T00:00:00Z"),
		}

		assert.Equal(t, "transfer-123", transfer.ID.ValueString())
		assert.Equal(t, "cred-456", transfer.CredentialID.ValueString())
		assert.Equal(t, "proj-789", transfer.DestinationProjectID.ValueString())
		assert.Equal(t, "2024-01-01T00:00:00Z", transfer.TransferredAt.ValueString())
	})

	t.Run("create transfer resource with null values", func(t *testing.T) {
		transfer := models.TransferResource{
			ID:                   types.StringNull(),
			CredentialID:         types.StringNull(),
			DestinationProjectID: types.StringNull(),
			TransferredAt:        types.StringNull(),
		}

		assert.True(t, transfer.ID.IsNull())
		assert.True(t, transfer.CredentialID.IsNull())
		assert.True(t, transfer.DestinationProjectID.IsNull())
		assert.True(t, transfer.TransferredAt.IsNull())
	})

	t.Run("create transfer resource with unknown values", func(t *testing.T) {
		transfer := models.TransferResource{
			ID:                   types.StringUnknown(),
			CredentialID:         types.StringUnknown(),
			DestinationProjectID: types.StringUnknown(),
			TransferredAt:        types.StringUnknown(),
		}

		assert.True(t, transfer.ID.IsUnknown())
		assert.True(t, transfer.CredentialID.IsUnknown())
		assert.True(t, transfer.DestinationProjectID.IsUnknown())
		assert.True(t, transfer.TransferredAt.IsUnknown())
	})

	t.Run("partial initialization", func(t *testing.T) {
		transfer := models.TransferResource{
			ID:           types.StringValue("partial-transfer"),
			CredentialID: types.StringValue("cred-partial"),
		}

		assert.Equal(t, "partial-transfer", transfer.ID.ValueString())
		assert.Equal(t, "cred-partial", transfer.CredentialID.ValueString())
		assert.True(t, transfer.DestinationProjectID.IsNull())
		assert.True(t, transfer.TransferredAt.IsNull())
	})

	t.Run("zero value struct", func(t *testing.T) {
		var transfer models.TransferResource
		assert.True(t, transfer.ID.IsNull())
		assert.True(t, transfer.CredentialID.IsNull())
		assert.True(t, transfer.DestinationProjectID.IsNull())
		assert.True(t, transfer.TransferredAt.IsNull())
	})
}
