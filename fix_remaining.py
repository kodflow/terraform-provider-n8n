#!/usr/bin/env python3
"""Fix remaining KTN-FUNC-007 errors by adding Params/Returns to existing doc comments."""

import re
import sys

def fix_file(filepath):
    """Fix a single file."""
    with open(filepath, 'r', encoding='utf-8') as f:
        lines = f.readlines()

    modified = False
    i = 0
    while i < len(lines):
        line = lines[i]

        # Check if this is a function definition
        if line.strip().startswith('func '):
            # Look backwards for doc comment
            doc_start = i - 1
            while doc_start >= 0 and lines[doc_start].strip().startswith('//'):
                doc_start -= 1
            doc_start += 1

            # Check if doc comment exists
            if doc_start < i:
                # Check if it already has Params: and Returns:
                doc_block = ''.join(lines[doc_start:i])
                has_params = 'Params:' in doc_block or '//Params:' in doc_block
                has_returns = 'Returns:' in doc_block or '//Returns:' in doc_block

                if not has_params or not has_returns:
                    # Add Params and Returns
                    new_lines = []
                    new_lines.append('//\n')
                    new_lines.append('// Params:\n')
                    new_lines.append('//   - None\n')
                    new_lines.append('//\n')
                    new_lines.append('// Returns:\n')
                    new_lines.append('//   - None\n')

                    # Insert before function
                    lines = lines[:i] + new_lines + lines[i:]
                    modified = True
                    i += len(new_lines)

        i += 1

    if modified:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.writelines(lines)
        return True
    return False

# Files to fix
files = [
    '/workspace/src/internal/provider/datasources/datasources.go',
    '/workspace/src/internal/provider/datasources/execution.go',
    '/workspace/src/internal/provider/datasources/executions.go',
    '/workspace/src/internal/provider/datasources/project.go',
    '/workspace/src/internal/provider/datasources/projects.go',
    '/workspace/src/internal/provider/datasources/tag.go',
    '/workspace/src/internal/provider/datasources/tags.go',
    '/workspace/src/internal/provider/datasources/user.go',
    '/workspace/src/internal/provider/datasources/users.go',
    '/workspace/src/internal/provider/datasources/variable.go',
    '/workspace/src/internal/provider/datasources/variables.go',
    '/workspace/src/internal/provider/datasources/workflow.go',
    '/workspace/src/internal/provider/datasources/workflows.go',
    '/workspace/src/internal/provider/resources/resources.go',
    '/workspace/src/internal/provider/resources/credential.go',
    '/workspace/src/internal/provider/resources/credential_transfer.go',
    '/workspace/src/internal/provider/resources/execution_retry.go',
    '/workspace/src/internal/provider/resources/project.go',
    '/workspace/src/internal/provider/resources/project_user.go',
    '/workspace/src/internal/provider/resources/source_control_pull.go',
    '/workspace/src/internal/provider/resources/tag.go',
    '/workspace/src/internal/provider/resources/user.go',
    '/workspace/src/internal/provider/resources/variable.go',
    '/workspace/src/internal/provider/resources/workflow.go',
    '/workspace/src/internal/provider/resources/workflow_transfer.go',
]

fixed = 0
for filepath in files:
    try:
        if fix_file(filepath):
            print(f"Fixed {filepath}")
            fixed += 1
    except Exception as e:
        print(f"Error fixing {filepath}: {e}", file=sys.stderr)

print(f"\nFixed {fixed} files")
