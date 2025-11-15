//go:build acceptance

package variable_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider"
)

func TestAccVariableResource(t *testing.T) {
	// NOTE: Variables require n8n enterprise license with feat:variables enabled.
	// This test will be skipped if the license does not support variables.
	t.Skip("Skipping variable acceptance test - requires n8n enterprise license with feat:variables")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccVariableResourceConfig("TERRAFORM_ACC_TEST_VAR", "test-value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_variable.test", "key", "TERRAFORM_ACC_TEST_VAR"),
					resource.TestCheckResourceAttr("n8n_variable.test", "value", "test-value"),
					resource.TestCheckResourceAttrSet("n8n_variable.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "n8n_variable.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccVariableResourceConfig("TERRAFORM_ACC_TEST_VAR", "test-value-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_variable.test", "value", "test-value-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccVariableResourceConfig(key, value string) string {
	return fmt.Sprintf(`
provider "n8n" {}

resource "n8n_variable" "test" {
  key   = %[1]q
  value = %[2]q
}
`, key, value)
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

// testAccCheckVariableDestroy verifies the variable has been destroyed
func testAccCheckVariableDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "n8n_variable" {
			continue
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("variable %s still exists", rs.Primary.ID)
		}
	}

	return nil
}
