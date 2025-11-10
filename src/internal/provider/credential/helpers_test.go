package credential

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/credential/models"
	"github.com/stretchr/testify/assert"
)

func TestWorkflowBackup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create workflow backup", wantErr: false},
		{name: "nil original workflow", wantErr: false},
		{name: "empty backup ID", wantErr: false},
		{name: "zero value struct", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create workflow backup":
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

			case "nil original workflow":
				backup := models.WorkflowBackup{
					ID:       "backup-nil",
					Original: nil,
				}

				assert.Equal(t, "backup-nil", backup.ID)
				assert.Nil(t, backup.Original)

			case "empty backup ID":
				backup := models.WorkflowBackup{
					ID:       "",
					Original: &n8nsdk.Workflow{},
				}

				assert.Equal(t, "", backup.ID)
				assert.NotNil(t, backup.Original)

			case "zero value struct":
				var backup models.WorkflowBackup
				assert.Equal(t, "", backup.ID)
				assert.Nil(t, backup.Original)

			case "error case - validation checks":
				// Test with various edge cases
				backup1 := models.WorkflowBackup{ID: "test", Original: nil}
				assert.Equal(t, "test", backup1.ID)
				assert.Nil(t, backup1.Original)

				backup2 := models.WorkflowBackup{ID: "", Original: &n8nsdk.Workflow{}}
				assert.Equal(t, "", backup2.ID)
				assert.NotNil(t, backup2.Original)
			}
		})
	}
}

func TestTransferResource(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create transfer resource with all fields", wantErr: false},
		{name: "create transfer resource with null values", wantErr: false},
		{name: "create transfer resource with unknown values", wantErr: false},
		{name: "partial initialization", wantErr: false},
		{name: "zero value struct", wantErr: false},
		{name: "error case - empty vs null vs unknown", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create transfer resource with all fields":
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

			case "create transfer resource with null values":
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

			case "create transfer resource with unknown values":
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

			case "partial initialization":
				transfer := models.TransferResource{
					ID:           types.StringValue("partial-transfer"),
					CredentialID: types.StringValue("cred-partial"),
				}

				assert.Equal(t, "partial-transfer", transfer.ID.ValueString())
				assert.Equal(t, "cred-partial", transfer.CredentialID.ValueString())
				assert.True(t, transfer.DestinationProjectID.IsNull())
				assert.True(t, transfer.TransferredAt.IsNull())

			case "zero value struct":
				var transfer models.TransferResource
				assert.True(t, transfer.ID.IsNull())
				assert.True(t, transfer.CredentialID.IsNull())
				assert.True(t, transfer.DestinationProjectID.IsNull())
				assert.True(t, transfer.TransferredAt.IsNull())

			case "error case - empty vs null vs unknown":
				// Test empty string value
				transfer1 := models.TransferResource{
					ID:           types.StringValue(""),
					CredentialID: types.StringValue(""),
				}
				assert.Equal(t, "", transfer1.ID.ValueString())
				assert.Equal(t, "", transfer1.CredentialID.ValueString())

				// Test mixed null and unknown
				transfer2 := models.TransferResource{
					ID:                   types.StringNull(),
					CredentialID:         types.StringUnknown(),
					DestinationProjectID: types.StringValue("test"),
				}
				assert.True(t, transfer2.ID.IsNull())
				assert.True(t, transfer2.CredentialID.IsUnknown())
				assert.Equal(t, "test", transfer2.DestinationProjectID.ValueString())
			}
		})
	}
}
