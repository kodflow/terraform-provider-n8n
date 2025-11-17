#!/usr/bin/env node
/**
 * Copyright (c) 2024 Florent (Kodflow). All rights reserved.
 * Licensed under the Sustainable Use License 1.0
 * See LICENSE in the project root for license information.
 *
 * Generate Terraform examples for each node category
 */

const fs = require('fs');
const path = require('path');

const [,, dataDir, examplesDir] = process.argv;

if (!dataDir || !examplesDir) {
    console.error('Usage: generate-examples.js <data-dir> <examples-dir>');
    process.exit(1);
}

const registryFile = path.join(dataDir, 'n8n-nodes-registry.json');
const registry = JSON.parse(fs.readFileSync(registryFile, 'utf8'));

console.log('Generating Terraform examples...');

// Group nodes by category
const byCategory = registry.nodes.reduce((acc, node) => {
    if (!acc[node.category]) acc[node.category] = [];
    acc[node.category].push(node);
    return acc;
}, {});

// Generate example for each category
for (const [category, nodes] of Object.entries(byCategory)) {
    const categoryDir = path.join(examplesDir, category.toLowerCase());
    fs.mkdirSync(categoryDir, { recursive: true });

    // Generate showcase file with examples of each node
    const tfLines = [];
    tfLines.push('# ' + category + ' Nodes Showcase');
    tfLines.push('# Auto-generated from n8n repository');
    tfLines.push('');
    tfLines.push('terraform {');
    tfLines.push('  required_providers {');
    tfLines.push('    n8n = {');
    tfLines.push('      source  = "kodflow/n8n"');
    tfLines.push('      version = "~> 1.0"');
    tfLines.push('    }');
    tfLines.push('  }');
    tfLines.push('}');
    tfLines.push('');
    tfLines.push('provider "n8n" {');
    tfLines.push('  base_url = var.n8n_base_url');
    tfLines.push('  api_key  = var.n8n_api_key');
    tfLines.push('}');
    tfLines.push('');

    // Generate example for first 5 nodes of category
    const exampleNodes = nodes.slice(0, 5);
    let position = [250, 300];

    for (let i = 0; i < exampleNodes.length; i++) {
        const node = exampleNodes[i];
        const resourceName = node.name.toLowerCase().replace(/[^a-z0-9]/g, '_');

        tfLines.push(`# ${node.name}`);
        tfLines.push(`# ${node.description}`);
        tfLines.push(`resource "n8n_workflow_node" "${resourceName}" {`);
        tfLines.push(`  name     = "${node.name}"`);
        tfLines.push(`  type     = "${node.type}"`);
        tfLines.push(`  position = [${position[0]}, ${position[1]}]`);
        tfLines.push('');
        tfLines.push('  parameters = jsonencode({');
        tfLines.push('    # Add node-specific parameters here');
        tfLines.push('  })');
        tfLines.push('}');
        tfLines.push('');

        position[0] += 200;
    }

    // Add workflow that uses these nodes
    tfLines.push('# Example workflow combining the nodes above');
    tfLines.push(`resource "n8n_workflow" "${category.toLowerCase()}_showcase" {`);
    tfLines.push(`  name   = "ci-\${var.run_id}-${category} Showcase"`);
    tfLines.push('  active = false');
    tfLines.push('');
    tfLines.push('  nodes_json = jsonencode([');
    for (let i = 0; i < exampleNodes.length; i++) {
        const node = exampleNodes[i];
        const resourceName = node.name.toLowerCase().replace(/[^a-z0-9]/g, '_');
        const comma = i < exampleNodes.length - 1 ? ',' : '';
        tfLines.push(`    jsondecode(n8n_workflow_node.${resourceName}.node_json)${comma}`);
    }
    tfLines.push('  ])');
    tfLines.push('');
    tfLines.push('  connections_json = jsonencode({})');
    tfLines.push('}');

    fs.writeFileSync(path.join(categoryDir, 'main.tf'), tfLines.join('\n'));

    // Generate variables.tf
    const varsLines = [];
    varsLines.push(`# Variables for ${category} examples`);
    varsLines.push('variable "n8n_base_url" {');
    varsLines.push('  description = "Base URL of the n8n instance"');
    varsLines.push('  type        = string');
    varsLines.push('  default     = "http://localhost:5678"');
    varsLines.push('}');
    varsLines.push('');
    varsLines.push('variable "n8n_api_key" {');
    varsLines.push('  description = "API key for n8n authentication"');
    varsLines.push('  type        = string');
    varsLines.push('  sensitive   = true');
    varsLines.push('}');
    varsLines.push('');
    varsLines.push('variable "run_id" {');
    varsLines.push('  description = "Unique run identifier for CI/CD"');
    varsLines.push('  type        = string');
    varsLines.push('  default     = "local"');
    varsLines.push('}');

    fs.writeFileSync(path.join(categoryDir, 'variables.tf'), varsLines.join('\n'));

    // Generate README
    const readmeLines = [];
    readmeLines.push(`# ${category} Nodes Examples`);
    readmeLines.push('');
    readmeLines.push(`This directory contains examples for **${category}** category nodes.`);
    readmeLines.push('');
    readmeLines.push(`Total nodes in this category: **${nodes.length}**`);
    readmeLines.push('');
    readmeLines.push('## Nodes in this Category');
    readmeLines.push('');
    for (const node of nodes) {
        readmeLines.push(`- **${node.name}** (\`${node.type}\`) - ${node.description}`);
    }
    readmeLines.push('');
    readmeLines.push('## Usage');
    readmeLines.push('');
    readmeLines.push('```bash');
    readmeLines.push('cd ' + path.relative(process.cwd(), categoryDir));
    readmeLines.push('terraform init');
    readmeLines.push('terraform plan');
    readmeLines.push('terraform apply');
    readmeLines.push('```');

    fs.writeFileSync(path.join(categoryDir, 'README.md'), readmeLines.join('\n'));

    console.log(`✓ Generated examples for ${category} (${nodes.length} nodes)`);
}

// Generate comprehensive index
const indexLines = [];
indexLines.push('# N8N Nodes Examples Index');
indexLines.push('');
indexLines.push(`Generated from n8n ${registry.version} on ${new Date().toISOString()}`);
indexLines.push('');
indexLines.push(`Total nodes: **${registry.total_nodes}**`);
indexLines.push('');
indexLines.push('## Categories');
indexLines.push('');
for (const [category, nodes] of Object.entries(byCategory)) {
    indexLines.push(`- [${category}](./${category.toLowerCase()}/README.md) - ${nodes.length} nodes`);
}

fs.writeFileSync(path.join(examplesDir, 'INDEX.md'), indexLines.join('\n'));

console.log('');
console.log('✓ Examples generated successfully');
console.log(`  Categories: ${Object.keys(byCategory).length}`);
console.log(`  Total nodes: ${registry.total_nodes}`);
