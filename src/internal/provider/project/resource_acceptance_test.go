//go:build acceptance

// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package project_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider"
)

func TestAccProjectResource(t *testing.T) {
	// NOTE: Projects require n8n enterprise license.
	// This test requires an enterprise n8n instance.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfig("terraform-acc-test-project"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_project.test", "name", "terraform-acc-test-project"),
					resource.TestCheckResourceAttrSet("n8n_project.test", "id"),
					resource.TestCheckResourceAttrSet("n8n_project.test", "type"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "n8n_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProjectResourceConfig("terraform-acc-test-project-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_project.test", "name", "terraform-acc-test-project-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProjectResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "n8n" {}

resource "n8n_project" "test" {
  name = %[1]q
}
`, name)
}

func testAccPreCheck(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping acceptance test")
	}

	// Verify required environment variables
	provider.TestAccPreCheckEnv(t)

	// Log test environment info
	t.Logf("🔍 Running acceptance tests against n8n instance: %s", os.Getenv("N8N_API_URL"))
	t.Logf("📋 Testing project feature (requires enterprise license)")
}

var testAccProtoV6ProviderFactories = provider.TestAccProtoV6ProviderFactories

// testAccCheckProjectDestroy verifies the project has been destroyed.
func testAccCheckProjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "n8n_project" {
			continue
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("project %s still exists", rs.Primary.ID)
		}
	}

	return nil
}
