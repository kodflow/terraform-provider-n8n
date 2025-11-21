//go:build acceptance

// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package user_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider"
)

func TestAccUserResource(t *testing.T) {
	// NOTE: Users require n8n enterprise license.
	// This test requires an enterprise n8n instance with user management enabled.
	// The email must be unique for each test run.
	testEmail := fmt.Sprintf("terraform-acc-test-%d@example.com", time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceConfig(testEmail, "global:member"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_user.test", "email", testEmail),
					resource.TestCheckResourceAttr("n8n_user.test", "role", "global:member"),
					resource.TestCheckResourceAttrSet("n8n_user.test", "id"),
					resource.TestCheckResourceAttrSet("n8n_user.test", "is_pending"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "n8n_user.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing (role change)
			{
				Config: testAccUserResourceConfig(testEmail, "global:admin"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_user.test", "email", testEmail),
					resource.TestCheckResourceAttr("n8n_user.test", "role", "global:admin"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccUserResourceConfig(email, role string) string {
	return fmt.Sprintf(`
provider "n8n" {}

resource "n8n_user" "test" {
  email = %[1]q
  role  = %[2]q
}
`, email, role)
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
	t.Logf("📋 Testing user feature (requires enterprise license)")
}

var testAccProtoV6ProviderFactories = provider.TestAccProtoV6ProviderFactories

// testAccCheckUserDestroy verifies the user has been destroyed.
func testAccCheckUserDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "n8n_user" {
			continue
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("user %s still exists", rs.Primary.ID)
		}
	}

	return nil
}
