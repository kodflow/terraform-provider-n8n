// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package models defines data structures for credential resources.
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
