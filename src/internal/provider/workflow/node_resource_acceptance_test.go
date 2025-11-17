//go:build acceptance

package workflow_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
