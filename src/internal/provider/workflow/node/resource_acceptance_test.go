//go:build acceptance

package node_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
)

// TestAccWorkflowNodeResource tests the n8n_workflow_node resource.
func TestAccWorkflowNodeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing - Basic node
			{
				Config: testAccWorkflowNodeResourceConfigBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_workflow_node.webhook", "name", "Test Webhook"),
					resource.TestCheckResourceAttr("n8n_workflow_node.webhook", "type", "n8n-nodes-base.webhook"),
					resource.TestCheckResourceAttr("n8n_workflow_node.webhook", "type_version", "1"),
					resource.TestCheckResourceAttr("n8n_workflow_node.webhook", "disabled", "false"),
					resource.TestCheckResourceAttrSet("n8n_workflow_node.webhook", "id"),
					resource.TestCheckResourceAttrSet("n8n_workflow_node.webhook", "node_json"),
					testAccCheckNodeJSONValid("n8n_workflow_node.webhook"),
				),
			},
			// Create and Read testing - Code node with parameters
			{
				Config: testAccWorkflowNodeResourceConfigCode(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_workflow_node.code", "name", "Process Data"),
					resource.TestCheckResourceAttr("n8n_workflow_node.code", "type", "n8n-nodes-base.code"),
					resource.TestCheckResourceAttrSet("n8n_workflow_node.code", "parameters"),
					testAccCheckNodeJSONValid("n8n_workflow_node.code"),
					testAccCheckNodeJSONContainsField("n8n_workflow_node.code", "parameters"),
				),
			},
		},
	})
}

// TestAccWorkflowNodeResource_InCompleteWorkflow tests nodes as part of a complete workflow.
func TestAccWorkflowNodeResource_InCompleteWorkflow(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowNodeResourceConfigComplete(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check nodes exist in state
					resource.TestCheckResourceAttrSet("n8n_workflow_node.start", "id"),
					resource.TestCheckResourceAttrSet("n8n_workflow_node.process", "id"),
					resource.TestCheckResourceAttrSet("n8n_workflow_node.end", "id"),
					// Check workflow exists in n8n
					testAccCheckWorkflowExists("n8n_workflow.test"),
					testAccCheckWorkflowHasNodes("n8n_workflow.test", 3),
				),
			},
		},
	})
}

// testAccWorkflowNodeResourceConfigBasic returns config for basic webhook node.
func testAccWorkflowNodeResourceConfigBasic() string {
	return `
provider "n8n" {}

resource "n8n_workflow_node" "webhook" {
  name     = "Test Webhook"
  type     = "n8n-nodes-base.webhook"
  position = [250, 300]

  parameters = jsonencode({
    path         = "test-webhook"
    httpMethod   = "POST"
    responseMode = "lastNode"
  })
}
`
}

// testAccWorkflowNodeResourceConfigCode returns config for code node.
func testAccWorkflowNodeResourceConfigCode() string {
	return `
provider "n8n" {}

resource "n8n_workflow_node" "code" {
  name     = "Process Data"
  type     = "n8n-nodes-base.code"
  position = [450, 300]

  parameters = jsonencode({
    mode   = "runOnceForAllItems"
    jsCode = "return $input.all();"
  })
}
`
}

// testAccWorkflowNodeResourceConfigComplete returns complete workflow with nodes.
func testAccWorkflowNodeResourceConfigComplete() string {
	return `
provider "n8n" {}

resource "n8n_workflow_node" "start" {
  name     = "Start"
  type     = "n8n-nodes-base.start"
  position = [250, 300]
}

resource "n8n_workflow_node" "process" {
  name     = "Process"
  type     = "n8n-nodes-base.code"
  position = [450, 300]

  parameters = jsonencode({
    mode   = "runOnceForAllItems"
    jsCode = "return [{json: {processed: true}}];"
  })
}

resource "n8n_workflow_node" "end" {
  name     = "End"
  type     = "n8n-nodes-base.set"
  position = [650, 300]

  parameters = jsonencode({
    mode = "manual"
    fields = {
      values = [{
        name  = "result"
        type  = "string"
        value = "success"
      }]
    }
  })
}

resource "n8n_workflow" "test" {
  name   = "terraform-acc-test-modular-workflow"
  active = false

  nodes_json = jsonencode([
    jsondecode(n8n_workflow_node.start.node_json),
    jsondecode(n8n_workflow_node.process.node_json),
    jsondecode(n8n_workflow_node.end.node_json),
  ])

  connections_json = jsonencode({
    (n8n_workflow_node.start.name) = {
      main = [[{
        node  = n8n_workflow_node.process.name
        type  = "main"
        index = 0
      }]]
    }
    (n8n_workflow_node.process.name) = {
      main = [[{
        node  = n8n_workflow_node.end.name
        type  = "main"
        index = 0
      }]]
    }
  })
}
`
}

func testAccPreCheck(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping acceptance test")
	}

	provider.TestAccPreCheckEnv(t)
	t.Logf("🔍 Running node acceptance tests against n8n instance: %s", os.Getenv("N8N_API_URL"))
}

var testAccProtoV6ProviderFactories = provider.TestAccProtoV6ProviderFactories

// testAccCheckNodeJSONValid verifies the node_json attribute is valid JSON.
func testAccCheckNodeJSONValid(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		nodeJSON := rs.Primary.Attributes["node_json"]
		if nodeJSON == "" {
			return fmt.Errorf("node_json is empty")
		}

		var node map[string]any
		if err := json.Unmarshal([]byte(nodeJSON), &node); err != nil {
			return fmt.Errorf("node_json is not valid JSON: %w", err)
		}

		// Verify required fields
		requiredFields := []string{"id", "name", "type"}
		for _, field := range requiredFields {
			if _, ok := node[field]; !ok {
				return fmt.Errorf("node_json missing required field: %s", field)
			}
		}

		return nil
	}
}

// testAccCheckNodeJSONContainsField verifies node_json contains specific field.
func testAccCheckNodeJSONContainsField(resourceName, field string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		nodeJSON := rs.Primary.Attributes["node_json"]
		var node map[string]any
		if err := json.Unmarshal([]byte(nodeJSON), &node); err != nil {
			return fmt.Errorf("failed to parse node_json: %w", err)
		}

		if _, ok := node[field]; !ok {
			return fmt.Errorf("node_json missing field: %s", field)
		}

		return nil
	}
}

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

		n8nClient := client.NewN8nClient(os.Getenv("N8N_API_URL"), os.Getenv("N8N_API_KEY"))

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

		n8nClient := client.NewN8nClient(os.Getenv("N8N_API_URL"), os.Getenv("N8N_API_KEY"))

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
	n8nClient := client.NewN8nClient(os.Getenv("N8N_API_URL"), os.Getenv("N8N_API_KEY"))

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "n8n_workflow" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, httpResp, err := n8nClient.APIClient.WorkflowAPI.WorkflowsIdGet(context.Background(), rs.Primary.ID).Execute()
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

		if err == nil {
			return fmt.Errorf("workflow %s still exists in n8n", rs.Primary.ID)
		}

		if httpResp != nil && httpResp.StatusCode != 404 {
			return fmt.Errorf("unexpected error checking workflow deletion: status %d", httpResp.StatusCode)
		}
	}

	return nil
}
