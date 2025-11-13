//go:build acceptance

package tag_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/kodflow/n8n/src/internal/provider"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)

func TestAccTagResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTagResourceConfig("terraform-acc-test-tag"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_tag.test", "name", "terraform-acc-test-tag"),
					resource.TestCheckResourceAttrSet("n8n_tag.test", "id"),
					// ✅ Verify tag exists in n8n via API
					testAccCheckTagExists("n8n_tag.test"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "n8n_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// NOTE: Update testing skipped - n8n API does not support PUT /tags (405 Method Not Allowed)
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

	// Log test environment info
	t.Logf("🔍 Running acceptance tests against n8n instance: %s", os.Getenv("N8N_URL"))
}

var testAccProtoV6ProviderFactories = provider.TestAccProtoV6ProviderFactories

// testAccCheckTagExists verifies the tag exists in n8n via API call.
func testAccCheckTagExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("tag ID is not set")
		}

		// Create n8n client
		n8nClient := client.NewN8nClient(os.Getenv("N8N_URL"), os.Getenv("N8N_API_TOKEN"))

		// Get tag from n8n API
		tag, httpResp, err := n8nClient.APIClient.TagsAPI.TagsIdGet(context.Background(), rs.Primary.ID).Execute()
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

		if err != nil {
			return fmt.Errorf("tag %s does not exist in n8n: %w", rs.Primary.ID, err)
		}

		if tag.Id == nil || *tag.Id != rs.Primary.ID {
			return fmt.Errorf("tag ID mismatch: expected %s, got %v", rs.Primary.ID, tag.Id)
		}

		return nil
	}
}

// testAccCheckTagDestroy verifies the tag has been destroyed via API call.
func testAccCheckTagDestroy(s *terraform.State) error {
	// Create n8n client
	n8nClient := client.NewN8nClient(os.Getenv("N8N_URL"), os.Getenv("N8N_API_TOKEN"))

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "n8n_tag" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		// Try to get tag from n8n API
		_, httpResp, err := n8nClient.APIClient.TagsAPI.TagsIdGet(context.Background(), rs.Primary.ID).Execute()
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

		// If tag still exists (no 404), return error
		if err == nil {
			return fmt.Errorf("tag %s still exists in n8n", rs.Primary.ID)
		}

		// Check if error is 404 Not Found (expected after delete)
		if httpResp != nil && httpResp.StatusCode != 404 {
			return fmt.Errorf("unexpected error checking tag deletion: status %d", httpResp.StatusCode)
		}
	}

	return nil
}
