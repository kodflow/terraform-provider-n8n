//go:build acceptance

package credential_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider"
)

func TestAccCredentialResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCredentialResourceConfig("terraform-acc-test-credential", "httpHeaderAuth"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_credential.test", "name", "terraform-acc-test-credential"),
					resource.TestCheckResourceAttr("n8n_credential.test", "type", "httpHeaderAuth"),
					resource.TestCheckResourceAttrSet("n8n_credential.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "n8n_credential.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data"}, // data is sensitive and may not match exactly
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
}

var testAccProtoV6ProviderFactories = provider.TestAccProtoV6ProviderFactories

// testAccCheckCredentialDestroy verifies the credential has been destroyed
func testAccCheckCredentialDestroy(s *terraform.State) error {
	// Verify all n8n_credential resources have been destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "n8n_credential" {
			continue
		}

		// If we reach here, the resource still exists in state, which means destroy failed
		if rs.Primary.ID != "" {
			return fmt.Errorf("credential %s still exists", rs.Primary.ID)
		}
	}

	return nil
}
