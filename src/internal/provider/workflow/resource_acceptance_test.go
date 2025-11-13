//go:build acceptance

package workflow_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider"
)

func TestAccWorkflowResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkflowResourceConfig("terraform-acc-test-workflow", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_workflow.test", "name", "terraform-acc-test-workflow"),
					resource.TestCheckResourceAttr("n8n_workflow.test", "active", "false"),
					resource.TestCheckResourceAttrSet("n8n_workflow.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "n8n_workflow.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccWorkflowResourceConfig("terraform-acc-test-workflow-updated", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_workflow.test", "name", "terraform-acc-test-workflow-updated"),
					resource.TestCheckResourceAttr("n8n_workflow.test", "active", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWorkflowResourceConfig(name string, active bool) string {
	return fmt.Sprintf(`
provider "n8n" {}

resource "n8n_workflow" "test" {
  name   = %[1]q
  active = %[2]t

  nodes = [
    {
      id   = "start"
      name = "Start"
      type = "n8n-nodes-base.start"
      position = [250, 300]
      parameters = {}
    },
    {
      id   = "http"
      name = "HTTP Request"
      type = "n8n-nodes-base.httpRequest"
      position = [450, 300]
      parameters = {
        url            = "https://api.github.com/repos/n8n-io/n8n"
        method         = "GET"
        responseFormat = "json"
      }
    }
  ]

  connections = {
    "Start" = {
      main = [
        [
          {
            node  = "HTTP Request"
            type  = "main"
            index = 0
          }
        ]
      ]
    }
  }

  settings = {}
}
`, name, active)
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

// testAccCheckWorkflowDestroy verifies the workflow has been destroyed
func testAccCheckWorkflowDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "n8n_workflow" {
			continue
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("workflow %s still exists", rs.Primary.ID)
		}
	}

	return nil
}
