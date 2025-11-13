//go:build acceptance

package tag_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider"
)

func TestAccTagResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTagResourceConfig("terraform-acc-test-tag"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_tag.test", "name", "terraform-acc-test-tag"),
					resource.TestCheckResourceAttrSet("n8n_tag.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "n8n_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccTagResourceConfig("terraform-acc-test-tag-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_tag.test", "name", "terraform-acc-test-tag-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTagResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "n8n" {}

resource "n8n_tag" "test" {
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
}

var testAccProtoV6ProviderFactories = provider.TestAccProtoV6ProviderFactories

// testAccCheckTagDestroy verifies the tag has been destroyed
func testAccCheckTagDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "n8n_tag" {
			continue
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("tag %s still exists", rs.Primary.ID)
		}
	}

	return nil
}
