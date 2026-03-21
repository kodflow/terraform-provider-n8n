// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package workflow implements workflow management resources and data sources.
package workflow

import "github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"

// workflowContent holds the parsed workflow content fields for API update operations.
type workflowContent struct {
	// Name is the workflow display name.
	Name string
	// Nodes are the workflow nodes.
	Nodes []n8nsdk.Node
	// Connections holds workflow connections.
	Connections map[string]any
	// Settings holds workflow settings.
	Settings n8nsdk.WorkflowSettings
}
