//go:build acceptance

package workflow_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
)

// TestAccWorkflowConnectionResource tests the n8n_workflow_connection resource.
func TestAccWorkflowConnectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing - Basic connection
			{
				Config: testAccWorkflowConnectionResourceConfigBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_workflow_connection.basic", "source_node", "Start"),
					resource.TestCheckResourceAttr("n8n_workflow_connection.basic", "target_node", "End"),
					resource.TestCheckResourceAttr("n8n_workflow_connection.basic", "source_output", "main"),
					resource.TestCheckResourceAttr("n8n_workflow_connection.basic", "target_input", "main"),
					resource.TestCheckResourceAttr("n8n_workflow_connection.basic", "source_output_index", "0"),
					resource.TestCheckResourceAttr("n8n_workflow_connection.basic", "target_input_index", "0"),
					resource.TestCheckResourceAttrSet("n8n_workflow_connection.basic", "id"),
					resource.TestCheckResourceAttrSet("n8n_workflow_connection.basic", "connection_json"),
					testAccCheckConnectionJSONValid("n8n_workflow_connection.basic"),
				),
			},
			// Test multi-output connection (If node)
			{
				Config: testAccWorkflowConnectionResourceConfigMultiOutput(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_workflow_connection.true_branch", "source_output_index", "0"),
					resource.TestCheckResourceAttr("n8n_workflow_connection.false_branch", "source_output_index", "1"),
					testAccCheckConnectionJSONValid("n8n_workflow_connection.true_branch"),
					testAccCheckConnectionJSONValid("n8n_workflow_connection.false_branch"),
				),
			},
		},
	})
}

// TestAccWorkflowConnectionResource_InCompleteWorkflow tests connections in workflow.
func TestAccWorkflowConnectionResource_InCompleteWorkflow(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowConnectionResourceConfigComplete(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check connections exist in state
					resource.TestCheckResourceAttrSet("n8n_workflow_connection.start_to_if", "id"),
					resource.TestCheckResourceAttrSet("n8n_workflow_connection.if_true", "id"),
					resource.TestCheckResourceAttrSet("n8n_workflow_connection.if_false", "id"),
					// Check workflow has correct structure
					testAccCheckWorkflowExists("n8n_workflow.test"),
					testAccCheckWorkflowHasNodes("n8n_workflow.test", 4),
					testAccCheckWorkflowHasConnections("n8n_workflow.test", 3),
				),
			},
		},
	})
}

// testAccWorkflowConnectionResourceConfigBasic returns basic connection config.
func testAccWorkflowConnectionResourceConfigBasic() string {
	return `
provider "n8n" {}

resource "n8n_workflow_connection" "basic" {
  source_node         = "Start"
  source_output       = "main"
  source_output_index = 0
  target_node         = "End"
  target_input        = "main"
  target_input_index  = 0
}
`
}

// testAccWorkflowConnectionResourceConfigMultiOutput returns multi-output config.
func testAccWorkflowConnectionResourceConfigMultiOutput() string {
	return `
provider "n8n" {}

resource "n8n_workflow_connection" "true_branch" {
  source_node         = "Check Data"
  source_output       = "main"
  source_output_index = 0
  target_node         = "Process Valid"
  target_input        = "main"
  target_input_index  = 0
}

resource "n8n_workflow_connection" "false_branch" {
  source_node         = "Check Data"
  source_output       = "main"
  source_output_index = 1
  target_node         = "Process Invalid"
  target_input        = "main"
  target_input_index  = 0
}
`
}

// testAccWorkflowConnectionResourceConfigComplete returns complete workflow.
func testAccWorkflowConnectionResourceConfigComplete() string {
	return `
provider "n8n" {}

resource "n8n_workflow_node" "start" {
  name     = "Start"
  type     = "n8n-nodes-base.start"
  position = [250, 300]
}

resource "n8n_workflow_node" "check" {
  name     = "Check"
  type     = "n8n-nodes-base.if"
  position = [450, 300]

  parameters = jsonencode({
    conditions = {
      boolean = [{
        value1 = "={{ $json.valid }}"
        value2 = true
      }]
    }
  })
}

resource "n8n_workflow_node" "valid" {
  name     = "Valid Path"
  type     = "n8n-nodes-base.set"
  position = [650, 200]

  parameters = jsonencode({
    mode = "manual"
    fields = {
      values = [{
        name  = "status"
        type  = "string"
        value = "valid"
      }]
    }
  })
}

resource "n8n_workflow_node" "invalid" {
  name     = "Invalid Path"
  type     = "n8n-nodes-base.set"
  position = [650, 400]

  parameters = jsonencode({
    mode = "manual"
    fields = {
      values = [{
        name  = "status"
        type  = "string"
        value = "invalid"
      }]
    }
  })
}

resource "n8n_workflow_connection" "start_to_if" {
  source_node         = n8n_workflow_node.start.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.check.name
  target_input        = "main"
  target_input_index  = 0
}

resource "n8n_workflow_connection" "if_true" {
  source_node         = n8n_workflow_node.check.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.valid.name
  target_input        = "main"
  target_input_index  = 0
}

resource "n8n_workflow_connection" "if_false" {
  source_node         = n8n_workflow_node.check.name
  source_output       = "main"
  source_output_index = 1
  target_node         = n8n_workflow_node.invalid.name
  target_input        = "main"
  target_input_index  = 0
}

resource "n8n_workflow" "test" {
  name   = "terraform-acc-test-branching-workflow"
  active = false

  nodes_json = jsonencode([
    jsondecode(n8n_workflow_node.start.node_json),
    jsondecode(n8n_workflow_node.check.node_json),
    jsondecode(n8n_workflow_node.valid.node_json),
    jsondecode(n8n_workflow_node.invalid.node_json),
  ])

  connections_json = jsonencode({
    (n8n_workflow_node.start.name) = {
      main = [[{
        node  = n8n_workflow_node.check.name
        type  = "main"
        index = 0
      }]]
    }
    (n8n_workflow_node.check.name) = {
      main = [
        [{
          node  = n8n_workflow_node.valid.name
          type  = "main"
          index = 0
        }],
        [{
          node  = n8n_workflow_node.invalid.name
          type  = "main"
          index = 0
        }]
      ]
    }
  })
}
`
}

// testAccCheckConnectionJSONValid verifies connection_json is valid JSON.
func testAccCheckConnectionJSONValid(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		connectionJSON := rs.Primary.Attributes["connection_json"]
		if connectionJSON == "" {
			return fmt.Errorf("connection_json is empty")
		}

		var conn map[string]any
		if err := json.Unmarshal([]byte(connectionJSON), &conn); err != nil {
			return fmt.Errorf("connection_json is not valid JSON: %w", err)
		}

		// Verify required fields
		requiredFields := []string{"source_node", "target_node", "source_output", "target_input"}
		for _, field := range requiredFields {
			if _, ok := conn[field]; !ok {
				return fmt.Errorf("connection_json missing required field: %s", field)
			}
		}

		return nil
	}
}

// testAccCheckWorkflowHasConnections verifies workflow has expected connections.
func testAccCheckWorkflowHasConnections(resourceName string, expectedConnectionCount int) resource.TestCheckFunc {
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

		// Count total connections
		connectionCount := 0
		if workflow.Connections != nil {
			// workflow.Connections is already map[string]interface{}
			for _, nodeConns := range workflow.Connections {
				nodeConnsMap, ok := nodeConns.(map[string]any)
				if !ok {
					continue
				}

				for _, outputConns := range nodeConnsMap {
					outputConnsList, ok := outputConns.([]any)
					if !ok {
						continue
					}

					for _, connGroup := range outputConnsList {
						connGroupList, ok := connGroup.([]any)
						if !ok {
							continue
						}
						connectionCount += len(connGroupList)
					}
				}
			}
		}

		if connectionCount != expectedConnectionCount {
			return fmt.Errorf("workflow has %d connections, expected %d", connectionCount, expectedConnectionCount)
		}

		return nil
	}
}
