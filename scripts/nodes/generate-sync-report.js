#!/usr/bin/env node

/**
 * Generate detailed synchronization report for n8n nodes
 *
 * This script compares the current node registry with the previous version
 * and generates a detailed sync.md report for the Claude agent to process.
 *
 * The report includes:
 * - Added nodes (new nodes that weren't in previous version)
 * - Removed nodes (nodes that existed before but are gone)
 * - Modified nodes (nodes with changed properties)
 * - Detailed property-level changes
 *
 * Usage: node scripts/nodes/generate-sync-report.js
 */

const fs = require('fs');
const path = require('path');

// Paths
const CURRENT_REGISTRY = path.join(__dirname, '../../data/n8n-nodes-registry.json');
const PREVIOUS_REGISTRY = path.join(__dirname, '../../data/n8n-nodes-registry.previous.json');
const SYNC_REPORT = path.join(__dirname, '../../NODES_SYNC.md');

/**
 * Deep comparison of two objects
 */
function deepDiff(obj1, obj2, path = '') {
  const changes = [];

  // Get all keys from both objects
  const allKeys = new Set([...Object.keys(obj1 || {}), ...Object.keys(obj2 || {})]);

  for (const key of allKeys) {
    const currentPath = path ? `${path}.${key}` : key;
    const val1 = obj1?.[key];
    const val2 = obj2?.[key];

    if (val1 === undefined && val2 !== undefined) {
      changes.push({
        type: 'added',
        path: currentPath,
        value: val2
      });
    } else if (val1 !== undefined && val2 === undefined) {
      changes.push({
        type: 'removed',
        path: currentPath,
        value: val1
      });
    } else if (typeof val1 === 'object' && typeof val2 === 'object') {
      if (Array.isArray(val1) && Array.isArray(val2)) {
        // Array comparison
        if (JSON.stringify(val1) !== JSON.stringify(val2)) {
          changes.push({
            type: 'modified',
            path: currentPath,
            oldValue: val1,
            newValue: val2
          });
        }
      } else if (val1 !== null && val2 !== null) {
        // Recursive object comparison
        changes.push(...deepDiff(val1, val2, currentPath));
      }
    } else if (val1 !== val2) {
      changes.push({
        type: 'modified',
        path: currentPath,
        oldValue: val1,
        newValue: val2
      });
    }
  }

  return changes;
}

/**
 * Compare two registries and detect changes
 */
function compareRegistries(current, previous) {
  const currentNodes = new Map(current.nodes.map(n => [n.type, n]));
  const previousNodes = new Map(previous.nodes.map(n => [n.type, n]));

  const added = [];
  const removed = [];
  const modified = [];

  // Find added nodes
  for (const [type, node] of currentNodes) {
    if (!previousNodes.has(type)) {
      added.push(node);
    }
  }

  // Find removed nodes
  for (const [type, node] of previousNodes) {
    if (!currentNodes.has(type)) {
      removed.push(node);
    }
  }

  // Find modified nodes
  for (const [type, currentNode] of currentNodes) {
    const previousNode = previousNodes.get(type);
    if (previousNode) {
      const changes = deepDiff(previousNode, currentNode);
      if (changes.length > 0) {
        modified.push({
          node: currentNode,
          changes
        });
      }
    }
  }

  return { added, removed, modified };
}

/**
 * Generate markdown report
 */
function generateReport(current, previous, diff) {
  const timestamp = new Date().toISOString();

  let report = `# N8N Nodes Synchronization Report

**Generated**: ${timestamp}
**Previous Version**: ${previous?.version || 'N/A'}
**Current Version**: ${current.version}
**Last Sync**: ${current.last_sync}

## Summary

- **Total Nodes**: ${current.nodes.length}
- **Added**: ${diff.added.length}
- **Removed**: ${diff.removed.length}
- **Modified**: ${diff.modified.length}
- **Unchanged**: ${current.nodes.length - diff.added.length - diff.removed.length - diff.modified.length}

---

## 🎯 Actions Required for Claude Agent

`;

  if (diff.added.length === 0 && diff.removed.length === 0 && diff.modified.length === 0) {
    report += `### ✅ No changes detected

All nodes are synchronized. No action required.

`;
  } else {
    report += `### Tasks to Execute

`;

    if (diff.added.length > 0) {
      report += `#### 1. Add Support for ${diff.added.length} New Node(s)

For each new node:
- [ ] Generate Terraform example in \`examples/nodes/{category}/{node-slug}/\`
- [ ] Create acceptance test
- [ ] Add to documentation
- [ ] Update README with supported nodes list

`;
    }

    if (diff.removed.length > 0) {
      report += `#### 2. Remove ${diff.removed.length} Deprecated Node(s)

For each removed node:
- [ ] Mark as deprecated in documentation
- [ ] Add migration guide if replacement exists
- [ ] Remove examples after grace period

`;
    }

    if (diff.modified.length > 0) {
      report += `#### 3. Update ${diff.modified.length} Modified Node(s)

For each modified node:
- [ ] Review property changes
- [ ] Update examples if needed
- [ ] Update tests to cover new properties
- [ ] Update documentation

`;
    }
  }

  report += `---

## 📋 Detailed Changes

`;

  // Added nodes
  if (diff.added.length > 0) {
    report += `### ✅ Added Nodes (${diff.added.length})

`;
    diff.added.forEach(node => {
      report += `#### ${node.name}

- **Type**: \`${node.type}\`
- **Category**: ${node.category}
- **Latest Version**: ${node.latest_version}
- **Description**: ${node.description || 'N/A'}
- **Inputs**: ${node.inputs.join(', ') || 'None'}
- **Outputs**: ${node.outputs.join(', ') || 'None'}
- **File**: \`${node.file}\`

**Full JSON**:
\`\`\`json
${JSON.stringify(node, null, 2)}
\`\`\`

`;
    });
  }

  // Removed nodes
  if (diff.removed.length > 0) {
    report += `### ❌ Removed Nodes (${diff.removed.length})

`;
    diff.removed.forEach(node => {
      report += `#### ${node.name}

- **Type**: \`${node.type}\`
- **Category**: ${node.category}
- **Was in**: \`${node.file}\`

`;
    });
  }

  // Modified nodes
  if (diff.modified.length > 0) {
    report += `### 🔄 Modified Nodes (${diff.modified.length})

`;
    diff.modified.forEach(({ node, changes }) => {
      report += `#### ${node.name} (\`${node.type}\`)

**Changes detected**: ${changes.length}

`;
      changes.forEach(change => {
        if (change.type === 'added') {
          report += `- ✅ **Added** \`${change.path}\`: \`${JSON.stringify(change.value)}\`\n`;
        } else if (change.type === 'removed') {
          report += `- ❌ **Removed** \`${change.path}\`: was \`${JSON.stringify(change.value)}\`\n`;
        } else if (change.type === 'modified') {
          report += `- 🔄 **Modified** \`${change.path}\`:\n`;
          report += `  - **Old**: \`${JSON.stringify(change.oldValue)}\`\n`;
          report += `  - **New**: \`${JSON.stringify(change.newValue)}\`\n`;
        }
      });

      report += `\n`;
    });
  }

  report += `---

## 📊 Node Distribution

`;

  const categories = {};
  current.nodes.forEach(node => {
    categories[node.category] = (categories[node.category] || 0) + 1;
  });

  report += `| Category | Count |\n`;
  report += `|----------|-------|\n`;
  Object.entries(categories).sort((a, b) => b[1] - a[1]).forEach(([cat, count]) => {
    report += `| ${cat} | ${count} |\n`;
  });

  report += `\n---

## 🔍 Testing Requirements

`;

  if (diff.added.length > 0 || diff.modified.length > 0) {
    report += `### Nodes Requiring Testing

`;

    const nodesToTest = [
      ...diff.added.map(n => ({ ...n, reason: 'New node' })),
      ...diff.modified.map(({ node }) => ({ ...node, reason: 'Modified properties' }))
    ];

    nodesToTest.forEach(node => {
      const requiresCredentials = node.type.includes('api') ||
                                  node.type.includes('slack') ||
                                  node.type.includes('github') ||
                                  node.type.includes('oauth');

      report += `#### ${node.name}

- **Reason**: ${node.reason}
- **Requires Credentials**: ${requiresCredentials ? '⚠️ YES' : '✅ NO'}
- **Test Priority**: ${requiresCredentials ? 'Medium (mock credentials)' : 'High (test immediately)'}

`;
    });
  } else {
    report += `✅ No new testing required - all nodes unchanged.

`;
  }

  report += `---

## 📝 Documentation Updates

`;

  if (diff.added.length > 0 || diff.removed.length > 0 || diff.modified.length > 0) {
    report += `### Files to Update

- [ ] \`README.md\` - Add new nodes to supported list
- [ ] \`NODES_REFERENCE.md\` - Update node catalog
- [ ] \`CHANGELOG.md\` - Document changes
`;

    if (diff.added.length > 0) {
      report += `- [ ] Generate ${diff.added.length} new example(s) in \`examples/nodes/\`\n`;
    }

    if (diff.modified.length > 0) {
      report += `- [ ] Update ${diff.modified.length} existing example(s)\n`;
    }
  }

  report += `\n---

## 🤖 Automation Notes

This report is auto-generated and should not be committed to the repository.
It serves as input for the Claude agent to process node synchronization changes.

**Next Steps**:
1. Review this report
2. Execute required actions (code gen, examples, tests)
3. Run full test suite
4. Update documentation
5. Commit changes

`;

  return report;
}

/**
 * Main execution
 */
function main() {
  console.log('📊 Generating synchronization report...\n');

  // Load current registry
  let current;
  try {
    current = JSON.parse(fs.readFileSync(CURRENT_REGISTRY, 'utf8'));
    console.log(`✅ Loaded current registry: ${current.nodes.length} nodes`);
  } catch (error) {
    console.error(`❌ Failed to load current registry: ${error.message}`);
    process.exit(1);
  }

  // Load previous registry (if exists)
  let previous = null;
  if (fs.existsSync(PREVIOUS_REGISTRY)) {
    try {
      previous = JSON.parse(fs.readFileSync(PREVIOUS_REGISTRY, 'utf8'));
      console.log(`✅ Loaded previous registry: ${previous.nodes.length} nodes`);
    } catch (error) {
      console.warn(`⚠️  Failed to load previous registry: ${error.message}`);
      console.warn(`   Treating all nodes as new`);
    }
  } else {
    console.warn(`⚠️  No previous registry found - first sync`);
    console.warn(`   Treating all nodes as new`);
  }

  // Compare registries
  const diff = previous
    ? compareRegistries(current, previous)
    : { added: current.nodes, removed: [], modified: [] };

  console.log(`\n📊 Changes detected:`);
  console.log(`   Added: ${diff.added.length}`);
  console.log(`   Removed: ${diff.removed.length}`);
  console.log(`   Modified: ${diff.modified.length}`);

  // Generate report
  const report = generateReport(current, previous, diff);

  // Write report
  fs.writeFileSync(SYNC_REPORT, report);
  console.log(`\n✅ Report generated: ${SYNC_REPORT}`);

  // Save current as previous for next run
  fs.copyFileSync(CURRENT_REGISTRY, PREVIOUS_REGISTRY);
  console.log(`✅ Saved registry backup for next comparison`);

  // Summary
  console.log(`\n📋 Summary:`);
  if (diff.added.length === 0 && diff.removed.length === 0 && diff.modified.length === 0) {
    console.log(`   🎉 No changes - all nodes synchronized!`);
  } else {
    console.log(`   ⚠️  ${diff.added.length + diff.removed.length + diff.modified.length} total changes detected`);
    console.log(`   📄 Review ${SYNC_REPORT} for details`);
  }
}

// Run
main();
