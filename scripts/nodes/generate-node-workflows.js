#!/usr/bin/env node

/**
 * Generate complete workflow examples for each n8n node
 *
 * This script creates 1 folder per node with:
 * - main.tf: Complete workflow testing the node
 * - variables.tf: Provider configuration
 * - README.md: Node documentation
 *
 * Usage: node scripts/nodes/generate-node-workflows.js
 */

const fs = require('fs');
const path = require('path');

// Paths
const REGISTRY_FILE = path.join(__dirname, '../../data/n8n-nodes-registry.json');
const EXAMPLES_DIR = path.join(__dirname, '../../examples/nodes');

// Load registry
let registry;
try {
  const registryContent = fs.readFileSync(REGISTRY_FILE, 'utf8');
  registry = JSON.parse(registryContent);
  console.log(`📦 Loaded registry with ${registry.nodes.length} nodes`);
} catch (error) {
  console.error(`❌ Failed to load registry: ${error.message}`);
  process.exit(1);
}

/**
 * Convert node name to slug (filesystem-safe)
 */
function toSlug(name) {
  return name
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '');
}

/**
 * Get full node type with n8n-nodes-base prefix if needed
 */
function getFullNodeType(type) {
  // If type already has a prefix (n8n-nodes-*, @n8n/*), return as-is
  if (type.startsWith('n8n-nodes-') || type.startsWith('@n8n/')) {
    return type;
  }
  // Otherwise add n8n-nodes-base. prefix
  return `n8n-nodes-base.${type}`;
}

/**
 * Get example parameters for a node
 */
function getExampleParameters(node) {
  const fullType = getFullNodeType(node.type);

  // Node-specific parameter examples
  const examples = {
    'n8n-nodes-base.webhook': {
      path: `test-${toSlug(node.name)}`,
      httpMethod: 'POST',
      responseMode: 'onReceived'
    },
    'n8n-nodes-base.code': {
      mode: 'runOnceForAllItems',
      jsCode: `// Process data\nconst items = $input.all();\nreturn items.map(item => ({\n  json: {\n    ...item.json,\n    processed: true,\n    timestamp: new Date().toISOString()\n  }\n}));`
    },
    'n8n-nodes-base.if': {
      conditions: {
        boolean: [{
          value1: '={{ $json.isValid }}',
          value2: true
        }]
      }
    },
    'n8n-nodes-base.switch': {
      mode: 'rules',
      rules: {
        values: [
          { value: '={{ $json.type === "A" }}' },
          { value: '={{ $json.type === "B" }}' }
        ]
      }
    },
    'n8n-nodes-base.set': {
      mode: 'manual',
      fields: {
        values: [{
          name: 'output',
          type: 'string',
          value: '={{ $json }}'
        }]
      }
    },
    'n8n-nodes-base.merge': {
      mode: 'combine',
      mergeByFields: {
        values: [{ field1: 'id', field2: 'id' }]
      }
    },
    'n8n-nodes-base.httpRequest': {
      method: 'GET',
      url: 'https://httpbin.org/get',
      authentication: 'none',
      options: {}
    },
    'n8n-nodes-base.schedule': {
      rule: {
        interval: [{
          field: 'hours',
          hoursInterval: 1
        }]
      }
    },
    'n8n-nodes-base.manualTrigger': {}
  };

  return examples[fullType] || {
    // Generic fallback
    note: `Configure ${node.name} parameters here`
  };
}

/**
 * Generate main.tf for a node
 */
function generateMainTf(node) {
  const isTrigger = node.category === 'Trigger';
  const needsInput = !isTrigger;
  const nodeSlug = toSlug(node.name);
  const fullNodeType = getFullNodeType(node.type);

  const params = getExampleParameters(node);
  const paramsJson = JSON.stringify(params, null, 4).split('\n').map(line => '    ' + line).join('\n');

  let content = `# Test workflow for ${node.name}
# Category: ${node.category}
# Type: ${fullNodeType}

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

`;

  // Add input node if needed
  if (needsInput) {
    content += `# INPUT: Manual trigger to start the workflow
resource "n8n_workflow_node" "manual_trigger" {
  name     = "Manual Trigger"
  type     = "n8n-nodes-base.manualTrigger"
  position = [250, 300]
}

`;
  }

  // Add the tested node
  content += `# TESTED NODE: ${node.name}
resource "n8n_workflow_node" "test_node" {
  name     = "${node.name}"
  type     = "${fullNodeType}"
  position = [${needsInput ? 450 : 250}, 300]

  parameters = jsonencode(
${paramsJson}
  )
}

`;

  // Add output node
  content += `# OUTPUT: Display result
resource "n8n_workflow_node" "display_result" {
  name     = "Display Result"
  type     = "n8n-nodes-base.set"
  position = [${needsInput ? 650 : 450}, 300]

  parameters = jsonencode({
    mode = "manual"
    fields = {
      values = [{
        name  = "result"
        type  = "string"
        value = "={{ $json }}"
      }]
    }
  })
}

`;

  // Add connections
  content += `# CONNECTIONS\n`;

  if (needsInput) {
    content += `resource "n8n_workflow_connection" "input_to_test" {
  source_node         = n8n_workflow_node.manual_trigger.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.test_node.name
  target_input        = "main"
  target_input_index  = 0
}

`;
  }

  content += `resource "n8n_workflow_connection" "test_to_output" {
  source_node         = n8n_workflow_node.test_node.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.display_result.name
  target_input        = "main"
  target_input_index  = 0
}

`;

  // Build workflow
  const nodes = needsInput
    ? ['manual_trigger', 'test_node', 'display_result']
    : ['test_node', 'display_result'];

  content += `# WORKFLOW
resource "n8n_workflow" "test_${nodeSlug}" {
  name   = "Test: ${node.name}"
  active = false

  nodes_json = jsonencode([
    ${nodes.map(n => `jsondecode(n8n_workflow_node.${n}.node_json)`).join(',\n    ')}
  ])

  connections_json = jsonencode({
`;

  if (needsInput) {
    content += `    (n8n_workflow_node.manual_trigger.name) = {
      main = [[{
        node  = n8n_workflow_node.test_node.name
        type  = "main"
        index = 0
      }]]
    }
`;
  }

  content += `    (n8n_workflow_node.test_node.name) = {
      main = [[{
        node  = n8n_workflow_node.display_result.name
        type  = "main"
        index = 0
      }]]
    }
  })
}

# OUTPUTS
output "workflow_id" {
  value       = n8n_workflow.test_${nodeSlug}.id
  description = "ID of the test workflow"
}

output "workflow_name" {
  value       = n8n_workflow.test_${nodeSlug}.name
  description = "Name of the test workflow"
}
`;

  return content;
}

/**
 * Generate variables.tf
 */
function generateVariablesTf() {
  return `# Variables for node test workflow

variable "n8n_base_url" {
  description = "Base URL of the n8n instance"
  type        = string
  default     = "http://localhost:5678"
}

variable "n8n_api_key" {
  description = "API key for n8n authentication"
  type        = string
  sensitive   = true
}
`;
}

/**
 * Generate README.md for a node
 */
function generateReadme(node) {
  const isTrigger = node.category === 'Trigger';
  const fullNodeType = getFullNodeType(node.type);

  return `# ${node.name} Node Test

**Category**: ${node.category}
**Type**: \`${fullNodeType}\`
**Latest Version**: ${node.latest_version}

## Description

${node.description || `Test workflow for the ${node.name} node.`}

## Node Information

- **Inputs**: ${node.inputs.length > 0 ? node.inputs.join(', ') : 'None (Trigger node)'}
- **Outputs**: ${node.outputs.join(', ')}
- **File**: \`${node.file}\`

## Workflow Structure

\`\`\`
${isTrigger ? '' : 'Manual Trigger → '}${node.name} → Display Result
\`\`\`

## Usage

### 1. Initialize Terraform

\`\`\`bash
terraform init
\`\`\`

### 2. Plan

\`\`\`bash
terraform plan -var="n8n_api_key=YOUR_API_KEY"
\`\`\`

### 3. Apply

\`\`\`bash
terraform apply -var="n8n_api_key=YOUR_API_KEY"
\`\`\`

### 4. Test in n8n

1. Open the workflow in n8n UI
2. ${isTrigger ? 'Trigger the workflow (webhook, schedule, etc.)' : 'Click "Execute Workflow"'}
3. Check the result in the Display Result node

## Customization

Edit the \`parameters\` in \`main.tf\` to customize the ${node.name} node behavior.

See [n8n ${node.name} documentation](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.${toSlug(node.name)}/) for available parameters.

## Notes

${isTrigger
  ? `- This is a trigger node, so it doesn't need an input\n- The workflow will be triggered by ${node.name.toLowerCase()} events`
  : `- This workflow uses a manual trigger for testing\n- In production, replace with appropriate trigger node`
}
`;
}

/**
 * Generate example for a single node
 */
function generateNodeExample(node) {
  const categorySlug = toSlug(node.category);
  const nodeSlug = toSlug(node.name);

  const nodeDir = path.join(EXAMPLES_DIR, categorySlug, nodeSlug);

  // Create directory
  fs.mkdirSync(nodeDir, { recursive: true });

  // Generate files
  fs.writeFileSync(
    path.join(nodeDir, 'main.tf'),
    generateMainTf(node)
  );

  fs.writeFileSync(
    path.join(nodeDir, 'variables.tf'),
    generateVariablesTf()
  );

  fs.writeFileSync(
    path.join(nodeDir, 'README.md'),
    generateReadme(node)
  );

  console.log(`✅ Generated: examples/nodes/${categorySlug}/${nodeSlug}/`);
}

/**
 * Main execution
 */
function main() {
  console.log('🚀 Generating per-node workflow examples...\n');

  let generated = 0;
  let skipped = 0;

  for (const node of registry.nodes) {
    try {
      generateNodeExample(node);
      generated++;
    } catch (error) {
      console.error(`❌ Failed to generate ${node.name}: ${error.message}`);
      skipped++;
    }
  }

  console.log(`\n✅ Generation complete!`);
  console.log(`   Generated: ${generated} node examples`);
  console.log(`   Skipped: ${skipped}`);
  console.log(`\n📁 Examples location: ${EXAMPLES_DIR}`);
}

// Run
main();
