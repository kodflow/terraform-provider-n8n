// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package models defines shared data structures for the provider.
package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// N8nProviderModel represents the configuration schema for the n8n Terraform provider.
// It defines the required configuration parameters for authenticating with n8n API.
type N8nProviderModel struct {
	// APIKey is the n8n API key used for authentication
	APIKey types.String `tfsdk:"api_key"`

	// BaseURL is the base URL for the n8n instance (e.g., "https://n8n.example.com")
	BaseURL types.String `tfsdk:"base_url"`
}
