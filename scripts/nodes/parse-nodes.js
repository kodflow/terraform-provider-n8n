#!/usr/bin/env node
/**
 * Copyright (c) 2024 Florent (Kodflow). All rights reserved.
 * Licensed under the Sustainable Use License 1.0
 * See LICENSE in the project root for license information.
 *
 * Parse n8n nodes from the repository and generate a registry JSON
 */

const fs = require('fs');
const path = require('path');

// Parse command line arguments
const [,, cacheDir, dataDir] = process.argv;

if (!cacheDir || !dataDir) {
    console.error('Usage: parse-nodes.js <cache-dir> <data-dir>');
    process.exit(1);
}

const nodesBaseDir = path.join(cacheDir, 'packages', 'nodes-base', 'nodes');
const registryFile = path.join(dataDir, 'n8n-nodes-registry.json');
const metadataFile = path.join(dataDir, 'n8n-nodes-metadata.json');

console.log('Parsing nodes from:', nodesBaseDir);

// Read all node directories
function discoverNodes() {
    const nodes = [];
    const categories = {};

    if (!fs.existsSync(nodesBaseDir)) {
        console.error('Nodes directory not found:', nodesBaseDir);
        process.exit(1);
    }

    const entries = fs.readdirSync(nodesBaseDir, { withFileTypes: true });

    for (const entry of entries) {
        if (!entry.isDirectory()) continue;

        const nodeName = entry.name;
        const nodeDir = path.join(nodesBaseDir, nodeName);

        // Look for .node.ts file
        const nodeFiles = fs.readdirSync(nodeDir).filter(f => f.endsWith('.node.ts'));

        if (nodeFiles.length === 0) continue;

        const nodeFile = path.join(nodeDir, nodeFiles[0]);
        const nodeJsonFile = path.join(nodeDir, nodeName + '.node.json');

        try {
            // Parse basic info from file content
            const content = fs.readFileSync(nodeFile, 'utf8');

            // Extract node type identifier (e.g., "n8n-nodes-base.webhook")
            const typeMatch = content.match(/name:\s*['"]([^'"]+)['"]/);
            const displayNameMatch = content.match(/displayName:\s*['"]([^'"]+)['"]/);
            const descriptionMatch = content.match(/description:\s*['"]([^'"]+)['"]/);
            const groupMatch = content.match(/group:\s*\[['"]([^'"]+)['"]\]/);
            const versionMatch = content.match(/version:\s*(\d+)/g);

            // Try to detect input/output types
            const inputsMatch = content.match(/inputs:\s*\[([^\]]+)\]/);
            const outputsMatch = content.match(/outputs:\s*\[([^\]]+)\]/);

            const nodeType = typeMatch ? typeMatch[1] : `n8n-nodes-base.${nodeName.toLowerCase()}`;
            const displayName = displayNameMatch ? displayNameMatch[1] : nodeName;
            const description = descriptionMatch ? descriptionMatch[1] : '';
            const group = groupMatch ? groupMatch[1] : 'action';
            const versions = versionMatch ? versionMatch.map(v => parseInt(v.match(/\d+/)[0])) : [1];

            // Parse inputs/outputs
            const inputs = inputsMatch ? parseConnectionArray(inputsMatch[1]) : ['main'];
            const outputs = outputsMatch ? parseConnectionArray(outputsMatch[1]) : ['main'];

            // Determine category from group or file location
            let category = 'Integration';
            if (group === 'trigger') category = 'Trigger';
            else if (['Code', 'Set', 'Merge', 'If', 'Switch'].includes(nodeName)) category = 'Core';
            else if (['PostgreSQL', 'MySQL', 'MongoDB'].some(db => nodeName.includes(db))) category = 'Database';

            // Try to read .node.json for additional metadata
            let resources = {};
            if (fs.existsSync(nodeJsonFile)) {
                try {
                    const nodeJson = JSON.parse(fs.readFileSync(nodeJsonFile, 'utf8'));
                    resources = nodeJson.resources || {};
                } catch (e) {
                    // Ignore JSON parse errors
                }
            }

            const nodeInfo = {
                name: displayName,
                type: nodeType,
                category,
                group,
                versions,
                latest_version: Math.max(...versions),
                description,
                inputs,
                outputs,
                file: path.relative(cacheDir, nodeFile),
                resources
            };

            nodes.push(nodeInfo);

            // Count by category
            categories[category] = (categories[category] || 0) + 1;

        } catch (error) {
            console.warn(`Warning: Failed to parse ${nodeName}:`, error.message);
        }
    }

    return { nodes, categories };
}

// Parse connection arrays like ['main'] or [NodeConnectionType.Main]
function parseConnectionArray(str) {
    const cleaned = str.replace(/NodeConnectionType\./g, '').replace(/['"]/g, '').trim();
    return cleaned.split(',').map(s => s.trim().toLowerCase()).filter(Boolean);
}

// Get n8n version
function getN8nVersion() {
    const versionFile = path.join(dataDir, 'n8n-nodes-version.txt');
    if (fs.existsSync(versionFile)) {
        return fs.readFileSync(versionFile, 'utf8').trim();
    }
    return 'unknown';
}

// Main execution
console.log('Starting node discovery...');
const { nodes, categories } = discoverNodes();

console.log(`Found ${nodes.length} nodes in ${Object.keys(categories).length} categories`);

// Sort nodes by name
nodes.sort((a, b) => a.name.localeCompare(b.name));

// Create registry
const registry = {
    version: getN8nVersion(),
    last_sync: new Date().toISOString(),
    total_nodes: nodes.length,
    nodes
};

// Create metadata
const metadata = {
    version: getN8nVersion(),
    last_sync: new Date().toISOString(),
    total_nodes: nodes.length,
    categories,
    groups: nodes.reduce((acc, node) => {
        acc[node.group] = (acc[node.group] || 0) + 1;
        return acc;
    }, {})
};

// Write files
fs.writeFileSync(registryFile, JSON.stringify(registry, null, 2));
fs.writeFileSync(metadataFile, JSON.stringify(metadata, null, 2));

console.log('✓ Registry written to:', registryFile);
console.log('✓ Metadata written to:', metadataFile);
console.log('');
console.log('Statistics:');
console.log(`  Total nodes: ${nodes.length}`);
console.log('  Categories:');
Object.entries(categories).forEach(([cat, count]) => {
    console.log(`    - ${cat}: ${count}`);
});
