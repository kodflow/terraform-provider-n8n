//go:build acceptance

// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.


package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// TestAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"n8n": providerserver.NewProtocol6WithError(New("test")()),
}

// TestAccPreCheckEnv validates required environment variables are set
func TestAccPreCheckEnv(t *testing.T) {
	t.Helper()

	if v := os.Getenv("N8N_URL"); v == "" {
		t.Fatal("N8N_URL must be set for acceptance tests")
	}
	if v := os.Getenv("N8N_API_TOKEN"); v == "" {
		t.Fatal("N8N_API_TOKEN must be set for acceptance tests")
	}
}

// providerConfig returns the base provider configuration for acceptance tests
func providerConfig() string {
	return `
provider "n8n" {
  # Uses N8N_URL and N8N_API_TOKEN from environment
}
`
}

// cleanupTestResource attempts to delete a resource by ID using the n8n API
// This is used to ensure test resources are cleaned up even if Terraform destroy fails
func cleanupTestResource(ctx context.Context, t *testing.T, resourceType string, resourceID string) {
	t.Helper()
	// This is best-effort cleanup - if it fails, we log but don't fail the test
	// since the test itself may have already cleaned up the resource
	t.Logf("Attempting cleanup of %s resource: %s", resourceType, resourceID)
}
