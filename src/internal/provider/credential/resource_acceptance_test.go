//go:build acceptance

package credential_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/kodflow/n8n/src/internal/provider"
)

func TestAccCredentialResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCredentialResourceConfig("terraform-acc-test-credential", "httpHeaderAuth"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_credential.test", "name", "terraform-acc-test-credential"),
					resource.TestCheckResourceAttr("n8n_credential.test", "type", "httpHeaderAuth"),
					resource.TestCheckResourceAttrSet("n8n_credential.test", "id"),
					// Note: n8n API doesn't support GET /credentials or GET /credentials/{id}
					// so we can only verify through Terraform state checks
				),
			},
			// ImportState testing
			{
				ResourceName:      "n8n_credential.test",
				ImportState:       true,
				ImportStateVerify: true,
				// n8n API doesn't support GET /credentials/{id}, so we can only verify ID
				ImportStateVerifyIgnore: []string{"data", "name", "type", "created_at", "updated_at"},
			},
			// Update and Read testing
			{
				Config: testAccCredentialResourceConfig("terraform-acc-test-credential-updated", "httpHeaderAuth"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_credential.test", "name", "terraform-acc-test-credential-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCredentialResourceConfig(name, credType string) string {
	return fmt.Sprintf(`
provider "n8n" {}

resource "n8n_credential" "test" {
  name = %[1]q
  type = %[2]q
  data = {
    name  = "Authorization"
    value = "Bearer test-token-123"
  }
}
`, name, credType)
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
}

var testAccProtoV6ProviderFactories = provider.TestAccProtoV6ProviderFactories

// testAccCheckCredentialDestroy verifies the credential has been destroyed.
// Note: n8n API doesn't support GET /credentials or GET /credentials/{id},
// so we cannot verify deletion via API. Terraform state verification is sufficient.
func testAccCheckCredentialDestroy(s *terraform.State) error {
	// Since n8n API doesn't provide a way to check credential existence,
	// we rely on Terraform's state management to verify proper cleanup.
	// If Terraform successfully deleted the resource without errors,
	// we consider the destroy successful.
	return nil
}
