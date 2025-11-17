# Complete Modular Workflow Example
# This demonstrates a REAL workflow using nodes from each category:
# - Trigger: Webhook
# - Core: Code, If, Set
# - Integration: HTTP Request

terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = "~> 1.0"
    }
  }
}

provider "n8n" {
  base_url = var.n8n_base_url
  api_key  = var.n8n_api_key
}

# ============================================================================
# TRIGGER NODE - Webhook receives data
# ============================================================================
resource "n8n_workflow_node" "api_webhook" {
  name     = "API Webhook"
  type     = "n8n-nodes-base.webhook"
  position = [250, 400]

  parameters = jsonencode({
    path         = "data-processor"
    httpMethod   = "POST"
    responseMode = "lastNode"
    responseData = "allEntries"
    options = {
      rawBody = false
    }
  })

  webhook_id = "data-processor-webhook"
}

# ============================================================================
# CORE NODE 1 - Code validates and enriches data
# ============================================================================
resource "n8n_workflow_node" "validate_data" {
  name     = "Validate Input"
  type     = "n8n-nodes-base.code"
  position = [450, 400]

  parameters = jsonencode({
    mode   = "runOnceForAllItems"
    jsCode = <<-EOT
      // Validate incoming webhook data
      const items = $input.all();

      return items.map(item => {
        const data = item.json;

        // Validation
        const isValid = data.email && data.email.includes('@');
        const score = data.score || 0;

        return {
          json: {
            ...data,
            isValid: isValid,
            score: parseInt(score),
            processedAt: new Date().toISOString(),
            category: score > 50 ? 'high' : 'low'
          }
        };
      });
    EOT
  })
}

# ============================================================================
# CORE NODE 2 - If routes based on validation
# ============================================================================
resource "n8n_workflow_node" "check_validity" {
  name     = "Check if Valid"
  type     = "n8n-nodes-base.if"
  position = [650, 400]

  parameters = jsonencode({
    conditions = {
      boolean = [
        {
          value1 = "={{ $json.isValid }}"
          value2 = true
        }
      ]
    }
  })
}

# ============================================================================
# CORE NODE 3 - Set transforms valid data
# ============================================================================
resource "n8n_workflow_node" "prepare_valid_data" {
  name     = "Prepare Valid Data"
  type     = "n8n-nodes-base.set"
  position = [850, 300]

  parameters = jsonencode({
    mode = "manual"
    fields = {
      values = [
        {
          name  = "email"
          type  = "string"
          value = "={{ $json.email }}"
        },
        {
          name  = "score"
          type  = "number"
          value = "={{ $json.score }}"
        },
        {
          name  = "category"
          type  = "string"
          value = "={{ $json.category }}"
        },
        {
          name  = "status"
          type  = "string"
          value = "success"
        },
        {
          name  = "timestamp"
          type  = "string"
          value = "={{ $json.processedAt }}"
        }
      ]
    }
    options = {}
  })
}

# ============================================================================
# CORE NODE 4 - Set prepares error response
# ============================================================================
resource "n8n_workflow_node" "prepare_error_data" {
  name     = "Prepare Error Data"
  type     = "n8n-nodes-base.set"
  position = [850, 500]

  parameters = jsonencode({
    mode = "manual"
    fields = {
      values = [
        {
          name  = "status"
          type  = "string"
          value = "error"
        },
        {
          name  = "message"
          type  = "string"
          value = "Invalid email format"
        },
        {
          name  = "receivedData"
          type  = "object"
          value = "={{ $json }}"
        }
      ]
    }
  })
}

# ============================================================================
# INTEGRATION NODE - HTTP Request sends to external API
# ============================================================================
resource "n8n_workflow_node" "send_to_api" {
  name     = "Send to External API"
  type     = "n8n-nodes-base.httpRequest"
  position = [1050, 300]

  parameters = jsonencode({
    method         = "POST"
    url            = "https://httpbin.org/post"
    authentication = "none"
    sendBody       = true
    bodyParameters = {
      parameters = [
        {
          name  = "data"
          value = "={{ $json }}"
        }
      ]
    }
    options = {
      redirect = {
        redirect = {
          maxRedirects = 3
        }
      }
      response = {
        response = {
          fullResponse       = false
          neverError         = false
          responseFormat     = "autodetect"
          outputPropertyName = "data"
        }
      }
    }
  })
}

# ============================================================================
# CORE NODE 5 - Code formats final response
# ============================================================================
resource "n8n_workflow_node" "format_response" {
  name     = "Format Response"
  type     = "n8n-nodes-base.code"
  position = [1250, 400]

  parameters = jsonencode({
    mode   = "runOnceForAllItems"
    jsCode = <<-EOT
      const items = $input.all();

      return items.map(item => ({
        json: {
          success: true,
          message: "Data processed successfully",
          result: item.json,
          workflow: "complete-modular-workflow"
        }
      }));
    EOT
  })
}

# ============================================================================
# CONNECTIONS - Define data flow between nodes
# ============================================================================

# Webhook -> Validate
resource "n8n_workflow_connection" "webhook_to_validate" {
  source_node         = n8n_workflow_node.api_webhook.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.validate_data.name
  target_input        = "main"
  target_input_index  = 0
}

# Validate -> Check
resource "n8n_workflow_connection" "validate_to_check" {
  source_node         = n8n_workflow_node.validate_data.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.check_validity.name
  target_input        = "main"
  target_input_index  = 0
}

# Check -> Prepare Valid (true branch)
resource "n8n_workflow_connection" "check_to_valid" {
  source_node         = n8n_workflow_node.check_validity.name
  source_output       = "main"
  source_output_index = 0 # true output
  target_node         = n8n_workflow_node.prepare_valid_data.name
  target_input        = "main"
  target_input_index  = 0
}

# Check -> Prepare Error (false branch)
resource "n8n_workflow_connection" "check_to_error" {
  source_node         = n8n_workflow_node.check_validity.name
  source_output       = "main"
  source_output_index = 1 # false output
  target_node         = n8n_workflow_node.prepare_error_data.name
  target_input        = "main"
  target_input_index  = 0
}

# Prepare Valid -> Send to API
resource "n8n_workflow_connection" "valid_to_api" {
  source_node         = n8n_workflow_node.prepare_valid_data.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.send_to_api.name
  target_input        = "main"
  target_input_index  = 0
}

# Send to API -> Format Response
resource "n8n_workflow_connection" "api_to_format" {
  source_node         = n8n_workflow_node.send_to_api.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.format_response.name
  target_input        = "main"
  target_input_index  = 0
}

# Prepare Error -> Format Response (merge error path)
resource "n8n_workflow_connection" "error_to_format" {
  source_node         = n8n_workflow_node.prepare_error_data.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.format_response.name
  target_input        = "main"
  target_input_index  = 0
}

# ============================================================================
# WORKFLOW - Assemble all nodes and connections
# ============================================================================
resource "n8n_workflow" "complete_modular" {
  name   = "ci-${var.run_id}-Complete Modular Workflow"
  active = false # Set to true after testing

  # Assemble all nodes
  nodes_json = jsonencode([
    jsondecode(n8n_workflow_node.api_webhook.node_json),
    jsondecode(n8n_workflow_node.validate_data.node_json),
    jsondecode(n8n_workflow_node.check_validity.node_json),
    jsondecode(n8n_workflow_node.prepare_valid_data.node_json),
    jsondecode(n8n_workflow_node.prepare_error_data.node_json),
    jsondecode(n8n_workflow_node.send_to_api.node_json),
    jsondecode(n8n_workflow_node.format_response.node_json),
  ])

  # Build connections structure for n8n
  connections_json = jsonencode({
    (n8n_workflow_node.api_webhook.name) = {
      main = [[{
        node  = n8n_workflow_node.validate_data.name
        type  = "main"
        index = 0
      }]]
    }
    (n8n_workflow_node.validate_data.name) = {
      main = [[{
        node  = n8n_workflow_node.check_validity.name
        type  = "main"
        index = 0
      }]]
    }
    (n8n_workflow_node.check_validity.name) = {
      main = [
        # True branch (index 0)
        [{
          node  = n8n_workflow_node.prepare_valid_data.name
          type  = "main"
          index = 0
        }],
        # False branch (index 1)
        [{
          node  = n8n_workflow_node.prepare_error_data.name
          type  = "main"
          index = 0
        }]
      ]
    }
    (n8n_workflow_node.prepare_valid_data.name) = {
      main = [[{
        node  = n8n_workflow_node.send_to_api.name
        type  = "main"
        index = 0
      }]]
    }
    (n8n_workflow_node.send_to_api.name) = {
      main = [[{
        node  = n8n_workflow_node.format_response.name
        type  = "main"
        index = 0
      }]]
    }
    (n8n_workflow_node.prepare_error_data.name) = {
      main = [[{
        node  = n8n_workflow_node.format_response.name
        type  = "main"
        index = 0
      }]]
    }
  })
}

# ============================================================================
# OUTPUTS
# ============================================================================
output "workflow_id" {
  value       = n8n_workflow.complete_modular.id
  description = "ID of the complete modular workflow"
}

output "webhook_url" {
  value       = "${var.n8n_base_url}/webhook/data-processor"
  description = "URL to trigger the workflow"
}

output "workflow_structure" {
  value = {
    total_nodes       = 7
    trigger_nodes     = 1
    core_nodes        = 5
    integration_nodes = 1
    connections       = 7
  }
  description = "Workflow structure summary"
}

output "test_curl_command" {
  value       = <<-EOT
    # Test with valid data:
    curl -X POST ${var.n8n_base_url}/webhook/data-processor \
      -H "Content-Type: application/json" \
      -d '{"email": "test@example.com", "score": 75}'

    # Test with invalid data:
    curl -X POST ${var.n8n_base_url}/webhook/data-processor \
      -H "Content-Type: application/json" \
      -d '{"email": "invalid-email", "score": 25}'
  EOT
  description = "Example curl commands to test the workflow"
}
