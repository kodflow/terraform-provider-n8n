#!/usr/bin/env node

/**
 * Validate that each node workflow test covers ALL scenarios
 * - Checks if all outputs are tested
 * - Checks if all inputs are tested
 * - Generates coverage report
 */

const fs = require('fs');
const path = require('path');

// Known nodes with special connection requirements
const NODE_REQUIREMENTS = {
  'n8n-nodes-base.if': {
    outputs: 2,
    outputNames: ['True (output[0])', 'False (output[1])'],
    description: 'IF node must test both true and false branches'
  },
  'n8n-nodes-base.switch': {
    outputs: 4,
    outputNames: ['Output 1', 'Output 2', 'Output 3', 'Fallback'],
    description: 'Switch node must test all routing cases'
  },
  'n8n-nodes-base.filter': {
    outputs: 2,
    outputNames: ['Pass (output[0])', 'Fail (output[1])'],
    description: 'Filter node must test both pass and fail outputs'
  },
  'n8n-nodes-base.splitInBatches': {
    outputs: 2,
    outputNames: ['Batch (output[0])', 'Done (output[1])'],
    description: 'Split In Batches must test both batch and done outputs'
  },
  'n8n-nodes-base.compareDatasets': {
    outputs: 3,
    outputNames: ['Match', 'Mismatch', 'No Match'],
    description: 'Compare Datasets must test all three outputs'
  },
  'n8n-nodes-base.merge': {
    inputs: 2,
    inputNames: ['Input 1', 'Input 2'],
    description: 'Merge node must test multiple inputs'
  }
};

// Find all node workflow directories
function findAllNodeWorkflows(baseDir) {
  const workflows = [];

  function walk(dir) {
    const entries = fs.readdirSync(dir, { withFileTypes: true });

    for (const entry of entries) {
      const fullPath = path.join(dir, entry.name);

      if (entry.isDirectory()) {
        // Check if this directory has a main.tf (it's a node workflow)
        const mainTf = path.join(fullPath, 'main.tf');
        if (fs.existsSync(mainTf)) {
          workflows.push({
            path: fullPath,
            category: path.basename(path.dirname(fullPath)),
            name: entry.name,
            mainTf
          });
        } else {
          // Recurse into subdirectories
          walk(fullPath);
        }
      }
    }
  }

  walk(baseDir);
  return workflows;
}

// Parse main.tf to extract node type and connections
function analyzeWorkflow(mainTfPath) {
  const content = fs.readFileSync(mainTfPath, 'utf-8');

  // Extract node type from comment
  const typeMatch = content.match(/# Type: ([^\n]+)/);
  const nodeType = typeMatch ? typeMatch[1].trim() : null;

  if (!nodeType) {
    return { nodeType: null, error: 'Could not determine node type' };
  }

  // Count output nodes
  const outputNodeMatches = content.match(/resource "n8n_workflow_node" "output_\d+"/g);
  const outputCount = outputNodeMatches ? outputNodeMatches.length : 0;

  // Count display_result (single output)
  const hasDisplayResult = content.includes('resource "n8n_workflow_node" "display_result"');
  const totalOutputs = outputCount || (hasDisplayResult ? 1 : 0);

  // Count input nodes
  const inputNodeMatches = content.match(/resource "n8n_workflow_node" "input_\d+"/g);
  const inputCount = inputNodeMatches ? inputNodeMatches.length : 0;

  // Count manual_trigger (single input)
  const hasManualTrigger = content.includes('resource "n8n_workflow_node" "manual_trigger"');
  const totalInputs = inputCount || (hasManualTrigger ? 1 : 0);

  // Extract output connections
  const outputConnections = [];
  const outputConnRegex = /source_output_index\s*=\s*(\d+)/g;
  let match;
  while ((match = outputConnRegex.exec(content)) !== null) {
    const index = parseInt(match[1]);
    if (!outputConnections.includes(index)) {
      outputConnections.push(index);
    }
  }
  outputConnections.sort((a, b) => a - b);

  return {
    nodeType,
    totalOutputs,
    totalInputs,
    outputConnections,
    testedOutputs: outputConnections.length
  };
}

// Validate coverage
function validateCoverage(workflow, analysis) {
  const nodeType = analysis.nodeType;
  const requirements = NODE_REQUIREMENTS[nodeType];

  if (!requirements) {
    // Node without special requirements
    return {
      status: 'ok',
      message: 'Standard node (single input/output)'
    };
  }

  const issues = [];

  // Check outputs
  if (requirements.outputs) {
    if (analysis.totalOutputs < requirements.outputs) {
      issues.push(`Missing outputs: has ${analysis.totalOutputs}, needs ${requirements.outputs}`);
    }

    if (analysis.testedOutputs < requirements.outputs) {
      issues.push(`Not all outputs connected: ${analysis.testedOutputs}/${requirements.outputs} tested`);
    }

    // Check if all output indices are covered
    const expectedIndices = Array.from({ length: requirements.outputs }, (_, i) => i);
    const missingIndices = expectedIndices.filter(i => !analysis.outputConnections.includes(i));

    if (missingIndices.length > 0) {
      issues.push(`Missing output indices: ${missingIndices.map(i => `output[${i}]`).join(', ')}`);
    }
  }

  // Check inputs
  if (requirements.inputs) {
    if (analysis.totalInputs < requirements.inputs) {
      issues.push(`Missing inputs: has ${analysis.totalInputs}, needs ${requirements.inputs}`);
    }
  }

  if (issues.length > 0) {
    return {
      status: 'incomplete',
      message: issues.join('; '),
      expected: requirements.description
    };
  }

  return {
    status: 'complete',
    message: `✓ All ${requirements.outputs || requirements.inputs} ${requirements.outputs ? 'outputs' : 'inputs'} tested`
  };
}

// Main execution
function main() {
  console.log('🔍 Validating node workflow test coverage\n');
  console.log('━'.repeat(80));

  const nodesDir = path.join(__dirname, '../../examples/nodes');

  if (!fs.existsSync(nodesDir)) {
    console.error('❌ Nodes directory not found:', nodesDir);
    process.exit(1);
  }

  const workflows = findAllNodeWorkflows(nodesDir);
  console.log(`\n📦 Found ${workflows.length} node workflows\n`);

  let totalCount = 0;
  let completeCount = 0;
  let incompleteCount = 0;
  let standardCount = 0;

  const incompleteNodes = [];

  for (const workflow of workflows) {
    totalCount++;
    const analysis = analyzeWorkflow(workflow.mainTf);

    if (analysis.error) {
      console.log(`⚠️  ${workflow.category}/${workflow.name}: ${analysis.error}`);
      continue;
    }

    const validation = validateCoverage(workflow, analysis);

    if (validation.status === 'incomplete') {
      incompleteCount++;
      incompleteNodes.push({
        name: `${workflow.category}/${workflow.name}`,
        nodeType: analysis.nodeType,
        ...validation
      });
      console.log(`❌ ${workflow.category}/${workflow.name} (${analysis.nodeType})`);
      console.log(`   ${validation.message}`);
      console.log(`   Expected: ${validation.expected}`);
    } else if (validation.status === 'complete') {
      completeCount++;
      console.log(`✅ ${workflow.category}/${workflow.name} (${analysis.nodeType})`);
      console.log(`   ${validation.message}`);
    } else {
      standardCount++;
      // Don't print standard nodes to reduce noise
    }
  }

  // Summary
  console.log('\n' + '━'.repeat(80));
  console.log('\n📊 Coverage Summary\n');
  console.log(`Total workflows:         ${totalCount}`);
  console.log(`Standard nodes:          ${standardCount} (no special requirements)`);
  console.log(`Complete coverage:       ${completeCount} (all scenarios tested)`);
  console.log(`Incomplete coverage:     ${incompleteCount} (missing scenarios)`);

  const coveragePercent = ((completeCount + standardCount) / totalCount * 100).toFixed(1);
  console.log(`\nOverall coverage:        ${coveragePercent}%`);

  // List incomplete nodes
  if (incompleteNodes.length > 0) {
    console.log('\n' + '━'.repeat(80));
    console.log('\n❌ Nodes with incomplete test coverage:\n');

    for (const node of incompleteNodes) {
      console.log(`  • ${node.name}`);
      console.log(`    Type: ${node.nodeType}`);
      console.log(`    Issue: ${node.message}`);
      console.log(`    Fix: ${node.expected}`);
      console.log('');
    }

    console.log('💡 Run "make nodes/workflows" to regenerate all workflows with complete coverage\n');
    process.exit(1);
  } else {
    console.log('\n✅ All nodes have complete test coverage!\n');
    process.exit(0);
  }
}

main();
