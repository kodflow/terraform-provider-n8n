package models

import (
	"github.com/kodflow/n8n/sdk/n8nsdk"
)

// WorkflowBackup stores workflow state for rollback during credential rotation.
// Captures original workflow data during credential rotation to enable recovery if the operation fails.
type WorkflowBackup struct {
	ID       string
	Original *n8nsdk.Workflow
}
