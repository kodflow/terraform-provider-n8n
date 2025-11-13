//go:build acceptance

package workflow_test

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

func TestAccWorkflowResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWorkflowDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkflowResourceConfig("terraform-acc-test-workflow", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_workflow.test", "name", "terraform-acc-test-workflow"),
					resource.TestCheckResourceAttr("n8n_workflow.test", "active", "false"),
					resource.TestCheckResourceAttrSet("n8n_workflow.test", "id"),
					// ✅ Verify workflow exists in n8n via API
					testAccCheckWorkflowExists("n8n_workflow.test"),
					// ✅ Verify workflow has correct number of nodes
					testAccCheckWorkflowHasNodes("n8n_workflow.test", 2),
				),
			},
			// ImportState testing
			{
				ResourceName:      "n8n_workflow.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// NOTE: Update testing skipped - n8n API does not support PUT /workflows/{id} (405 Method Not Allowed)
			// NOTE: Activation testing skipped - n8n API does not support PATCH /workflows/{id} (405 Method Not Allowed)
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

  nodes_json = jsonencode([
    {
      id         = "start"
      name       = "Start"
      type       = "n8n-nodes-base.start"
      position   = [250, 300]
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
  ])

  connections_json = jsonencode({
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
  })

  settings_json = jsonencode({})
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

	// Log test environment info
	t.Logf("🔍 Running acceptance tests against n8n instance: %s", os.Getenv("N8N_URL"))
}

var testAccProtoV6ProviderFactories = provider.TestAccProtoV6ProviderFactories

// testAccCheckWorkflowExists verifies the workflow exists in n8n via API call.
func testAccCheckWorkflowExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("workflow ID is not set")
		}

		// Create n8n client
		n8nClient := client.NewN8nClient(os.Getenv("N8N_URL"), os.Getenv("N8N_API_TOKEN"))

		// Call n8n API to verify workflow exists
		workflow, httpResp, err := n8nClient.APIClient.WorkflowAPI.WorkflowsIdGet(context.Background(), rs.Primary.ID).Execute()
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

		if err != nil {
			return fmt.Errorf("workflow %s does not exist in n8n: %w", rs.Primary.ID, err)
		}

		if workflow.Id == nil || *workflow.Id != rs.Primary.ID {
			return fmt.Errorf("workflow ID mismatch: expected %s, got %v", rs.Primary.ID, workflow.Id)
		}

		return nil
	}
}

// testAccCheckWorkflowHasNodes verifies workflow has expected number of nodes via API.
func testAccCheckWorkflowHasNodes(resourceName string, expectedNodeCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		// Create n8n client
		n8nClient := client.NewN8nClient(os.Getenv("N8N_URL"), os.Getenv("N8N_API_TOKEN"))

		// Get workflow from API
		workflow, httpResp, err := n8nClient.APIClient.WorkflowAPI.WorkflowsIdGet(context.Background(), rs.Primary.ID).Execute()
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

		if err != nil {
			return fmt.Errorf("failed to get workflow: %w", err)
		}

		if len(workflow.Nodes) != expectedNodeCount {
			return fmt.Errorf("workflow has %d nodes, expected %d", len(workflow.Nodes), expectedNodeCount)
		}

		return nil
	}
}

// testAccCheckWorkflowDestroy verifies the workflow has been destroyed via API call.
func testAccCheckWorkflowDestroy(s *terraform.State) error {
	// Create n8n client
	n8nClient := client.NewN8nClient(os.Getenv("N8N_URL"), os.Getenv("N8N_API_TOKEN"))

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "n8n_workflow" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		// Try to get workflow from n8n API
		_, httpResp, err := n8nClient.APIClient.WorkflowAPI.WorkflowsIdGet(context.Background(), rs.Primary.ID).Execute()
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

		// If workflow still exists (no 404), return error
		if err == nil {
			return fmt.Errorf("workflow %s still exists in n8n", rs.Primary.ID)
		}

		// Check if error is 404 Not Found (expected after delete)
		if httpResp != nil && httpResp.StatusCode != 404 {
			return fmt.Errorf("unexpected error checking workflow deletion: status %d", httpResp.StatusCode)
		}
	}

	return nil
}
