# Advanced modular workflow example with Switch node (multiple branches)
# This demonstrates how one source node can connect to multiple target nodes

# Input webhook
resource "n8n_workflow_node" "api_webhook" {
  name     = "API Webhook"
  type     = "n8n-nodes-base.webhook"
  position = [250, 400]

  parameters = jsonencode({
    path         = "api-endpoint"
    httpMethod   = "POST"
    responseMode = "lastNode"
  })

  webhook_id = "api-endpoint-id"
}

# Switch node to route based on request type
resource "n8n_workflow_node" "route_switch" {
  name     = "Route by Type"
  type     = "n8n-nodes-base.switch"
  position = [450, 400]

  parameters = jsonencode({
    rules = {
      rules = [
        {
          conditions = {
            all = [
              {
                value1   = "={{ $json.type }}"
                value2   = "user"
                operator = "equals"
              }
            ]
          }
          outputKey = "users"
        },
        {
          conditions = {
            all = [
              {
                value1   = "={{ $json.type }}"
                value2   = "product"
                operator = "equals"
              }
            ]
          }
          outputKey = "products"
        },
        {
          conditions = {
            all = [
              {
                value1   = "={{ $json.type }}"
                value2   = "order"
                operator = "equals"
              }
            ]
          }
          outputKey = "orders"
        }
      ]
    }
  })
}

# Branch 1: Handle users
resource "n8n_workflow_node" "handle_users" {
  name     = "Handle Users"
  type     = "n8n-nodes-base.code"
  position = [650, 250]

  parameters = jsonencode({
    mode   = "runOnceForAllItems"
    jsCode = "return [{ json: { type: 'user', handled: true, data: $input.all() } }];"
  })
}

# Branch 2: Handle products
resource "n8n_workflow_node" "handle_products" {
  name     = "Handle Products"
  type     = "n8n-nodes-base.code"
  position = [650, 400]

  parameters = jsonencode({
    mode   = "runOnceForAllItems"
    jsCode = "return [{ json: { type: 'product', handled: true, data: $input.all() } }];"
  })
}

# Branch 3: Handle orders
resource "n8n_workflow_node" "handle_orders" {
  name     = "Handle Orders"
  type     = "n8n-nodes-base.code"
  position = [650, 550]

  parameters = jsonencode({
    mode   = "runOnceForAllItems"
    jsCode = "return [{ json: { type: 'order', handled: true, data: $input.all() } }];"
  })
}

# Connections: Webhook to Switch
resource "n8n_workflow_connection" "webhook_to_switch" {
  source_node         = n8n_workflow_node.api_webhook.name
  target_node         = n8n_workflow_node.route_switch.name
  source_output       = "main"
  source_output_index = 0
  target_input        = "main"
  target_input_index  = 0
}

# Switch to Branch 1 (output index 0 = first rule)
resource "n8n_workflow_connection" "switch_to_users" {
  source_node         = n8n_workflow_node.route_switch.name
  target_node         = n8n_workflow_node.handle_users.name
  source_output       = "main"
  source_output_index = 0 # First output
  target_input        = "main"
  target_input_index  = 0
}

# Switch to Branch 2 (output index 1 = second rule)
resource "n8n_workflow_connection" "switch_to_products" {
  source_node         = n8n_workflow_node.route_switch.name
  target_node         = n8n_workflow_node.handle_products.name
  source_output       = "main"
  source_output_index = 1 # Second output
  target_input        = "main"
  target_input_index  = 0
}

# Switch to Branch 3 (output index 2 = third rule)
resource "n8n_workflow_connection" "switch_to_orders" {
  source_node         = n8n_workflow_node.route_switch.name
  target_node         = n8n_workflow_node.handle_orders.name
  source_output       = "main"
  source_output_index = 2 # Third output
  target_input        = "main"
  target_input_index  = 0
}

# Create the workflow with switch routing
resource "n8n_workflow" "switch_example" {
  name   = "ci-${var.run_id}-Switch Routing Workflow"
  active = false

  nodes_json = jsonencode([
    jsondecode(n8n_workflow_node.api_webhook.node_json),
    jsondecode(n8n_workflow_node.route_switch.node_json),
    jsondecode(n8n_workflow_node.handle_users.node_json),
    jsondecode(n8n_workflow_node.handle_products.node_json),
    jsondecode(n8n_workflow_node.handle_orders.node_json),
  ])

  connections_json = jsonencode({
    (n8n_workflow_node.api_webhook.name) = {
      main = [[{
        node  = n8n_workflow_node.route_switch.name
        type  = "main"
        index = 0
      }]]
    }
    (n8n_workflow_node.route_switch.name) = {
      main = [
        # Output 0: Users branch
        [{
          node  = n8n_workflow_node.handle_users.name
          type  = "main"
          index = 0
        }],
        # Output 1: Products branch
        [{
          node  = n8n_workflow_node.handle_products.name
          type  = "main"
          index = 0
        }],
        # Output 2: Orders branch
        [{
          node  = n8n_workflow_node.handle_orders.name
          type  = "main"
          index = 0
        }]
      ]
    }
  })
}

output "switch_workflow_id" {
  value       = n8n_workflow.switch_example.id
  description = "ID of the switch routing workflow"
}
