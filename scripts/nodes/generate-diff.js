#!/usr/bin/env node
/**
 * Copyright (c) 2024 Florent (Kodflow). All rights reserved.
 * Licensed under the Sustainable Use License 1.0
 * See LICENSE in the project root for license information.
 *
 * Generate a changelog by comparing current and previous node registries
 */

const fs = require('fs');
const path = require('path');

const [,, dataDir] = process.argv;

if (!dataDir) {
    console.error('Usage: generate-diff.js <data-dir>');
    process.exit(1);
}

const registryFile = path.join(dataDir, 'n8n-nodes-registry.json');
const registryBackup = path.join(dataDir, 'n8n-nodes-registry.json.backup');
const changelogFile = path.join(dataDir, 'n8n-nodes-changelog.md');

// Backup current registry for next diff
if (fs.existsSync(registryFile) && !fs.existsSync(registryBackup)) {
    fs.copyFileSync(registryFile, registryBackup);
    console.log('No previous backup found. Current registry saved as baseline.');
    process.exit(0);
}

if (!fs.existsSync(registryBackup)) {
    console.log('No previous registry to compare with.');
    process.exit(0);
}

// Load registries
const current = JSON.parse(fs.readFileSync(registryFile, 'utf8'));
const previous = JSON.parse(fs.readFileSync(registryBackup, 'utf8'));

// Compare
const changes = {
    added: [],
    removed: [],
    modified: []
};

// Create maps for efficient lookup
const currentMap = new Map(current.nodes.map(n => [n.type, n]));
const previousMap = new Map(previous.nodes.map(n => [n.type, n]));

// Find added and modified
for (const node of current.nodes) {
    if (!previousMap.has(node.type)) {
        changes.added.push(node);
    } else {
        const prev = previousMap.get(node.type);
        if (JSON.stringify(node) !== JSON.stringify(prev)) {
            changes.modified.push({
                type: node.type,
                name: node.name,
                changes: detectChanges(prev, node)
            });
        }
    }
}

// Find removed
for (const node of previous.nodes) {
    if (!currentMap.has(node.type)) {
        changes.removed.push(node);
    }
}

// Detect specific changes
function detectChanges(prev, curr) {
    const diffs = [];

    if (prev.latest_version !== curr.latest_version) {
        diffs.push(`Version updated: ${prev.latest_version} → ${curr.latest_version}`);
    }

    if (prev.description !== curr.description) {
        diffs.push('Description changed');
    }

    if (JSON.stringify(prev.inputs) !== JSON.stringify(curr.inputs)) {
        diffs.push(`Inputs changed: ${JSON.stringify(prev.inputs)} → ${JSON.stringify(curr.inputs)}`);
    }

    if (JSON.stringify(prev.outputs) !== JSON.stringify(curr.outputs)) {
        diffs.push(`Outputs changed: ${JSON.stringify(prev.outputs)} → ${JSON.stringify(curr.outputs)}`);
    }

    return diffs;
}

// Generate changelog
const lines = [];
lines.push('# N8N Nodes Changelog');
lines.push('');
lines.push(`Generated: ${new Date().toISOString()}`);
lines.push('');
lines.push(`Comparing ${previous.version} → ${current.version}`);
lines.push('');

if (changes.added.length === 0 && changes.removed.length === 0 && changes.modified.length === 0) {
    lines.push('## No Changes Detected');
    lines.push('');
    lines.push('All nodes remain unchanged.');
} else {
    lines.push('## Summary');
    lines.push('');
    lines.push(`- **Added:** ${changes.added.length} nodes`);
    lines.push(`- **Removed:** ${changes.removed.length} nodes`);
    lines.push(`- **Modified:** ${changes.modified.length} nodes`);
    lines.push('');

    if (changes.added.length > 0) {
        lines.push('## Added Nodes');
        lines.push('');
        for (const node of changes.added) {
            lines.push(`### ${node.name}`);
            lines.push(`- **Type:** \`${node.type}\``);
            lines.push(`- **Category:** ${node.category}`);
            lines.push(`- **Description:** ${node.description}`);
            lines.push('');
        }
    }

    if (changes.removed.length > 0) {
        lines.push('## Removed Nodes');
        lines.push('');
        for (const node of changes.removed) {
            lines.push(`- **${node.name}** (\`${node.type}\`)`);
        }
        lines.push('');
    }

    if (changes.modified.length > 0) {
        lines.push('## Modified Nodes');
        lines.push('');
        for (const mod of changes.modified) {
            lines.push(`### ${mod.name}`);
            lines.push(`- **Type:** \`${mod.type}\``);
            lines.push('- **Changes:**');
            for (const change of mod.changes) {
                lines.push(`  - ${change}`);
            }
            lines.push('');
        }
    }
}

// Write changelog
fs.writeFileSync(changelogFile, lines.join('\n'));

// Update backup
fs.copyFileSync(registryFile, registryBackup);

console.log('✓ Changelog generated:', changelogFile);
console.log('');
console.log('Changes detected:');
console.log(`  Added: ${changes.added.length}`);
console.log(`  Removed: ${changes.removed.length}`);
console.log(`  Modified: ${changes.modified.length}`);
