#!/usr/bin/env python3
"""
Script to fix FUNC linter errors in Go files.
Handles FUNC-001 (>35 lines), FUNC-003 (magic numbers), FUNC-005 (complexity >10), FUNC-007 (missing docs).
"""

import re
import subprocess
import sys
from pathlib import Path

def add_constant_for_magic_number(content, number, context_line):
    """Add a constant for a magic number."""
    # Generate constant name from context
    const_name = f"const{number}Limit"
    if "32" in str(number) and "Parse" in context_line:
        const_name = "float32BitSize"
    elif "100" in str(number):
        const_name = "defaultPageSize"
    elif number == 2:
        const_name = "minRequiredFields"

    # Add constant at top of file after package and imports
    const_decl = f'\nconst (\n\t// {const_name} constant value.\n\t{const_name} = {number}\n)\n'

    # Find where to insert (after imports)
    import_end = content.rfind(')')
    if import_end != -1:
        # Find next line after imports
        next_newline = content.find('\n', import_end)
        if next_newline != -1:
            content = content[:next_newline+1] + const_decl + content[next_newline+1:]

    # Replace the number with the constant name
    content = content.replace(f'strconv.ParseFloat(data.ID.ValueString(), {number})',
                              f'strconv.ParseFloat(data.ID.ValueString(), {const_name})')
    content = content.replace(f'strconv.ParseFloat(id, {number})',
                              f'strconv.ParseFloat(id, {const_name})')

    return content

def extract_schema_attributes(content, func_start, func_end):
    """Extract attributes map to separate function."""
    # Find the Attributes: map[string]schema.Attribute{ section
    attr_match = re.search(r'Attributes:\s+map\[string\]schema\.Attribute\{', content[func_start:func_end])
    if not attr_match:
        return content

    attr_start = func_start + attr_match.end()
    # Find matching closing brace
    brace_count = 1
    i = attr_start
    while i < func_end and brace_count > 0:
        if content[i] == '{':
            brace_count += 1
        elif content[i] == '}':
            brace_count -= 1
        i += 1
    attr_end = i - 1

    attributes_body = content[attr_start:attr_end]

    # Create helper function name based on context
    func_name_match = re.search(r'func \(.*?\) (\w+)\(', content[max(0, func_start-200):func_start])
    helper_name = "schemaAttributes"
    if func_name_match:
        helper_name = func_name_match.group(1).lower() + "SchemaAttributes"

    # Create helper function
    helper_func = f'''
// {helper_name} returns the schema attributes.
//
// Returns:
//   - map of schema attributes
func {helper_name}() map[string]schema.Attribute {{
\treturn map[string]schema.Attribute{{{attributes_body}}}
}}
'''

    # Replace in original function
    new_schema = f'Attributes:          {helper_name}(),'
    content = content[:func_start + attr_match.start()] + new_schema + content[attr_end+1:]

    # Add helper function before the Schema function
    content = content[:func_start] + helper_func + '\n' + content[func_start:]

    return content

def main():
    # Get all Go files with FUNC errors
    result = subprocess.run(
        ['ktn-linter', 'lint', './...'],
        capture_output=True,
        text=True,
        cwd='/workspace/src'
    )

    # Parse errors
    lines = result.stdout.split('\n') + result.stderr.split('\n')
    files_to_fix = set()

    for line in lines:
        if 'KTN-FUNC-' in line and ('/workspace/' in line or '.go:' in line):
            # Extract filename
            match = re.search(r'(/workspace/[^\s:]+\.go)', line)
            if match:
                files_to_fix.add(match.group(1))

    print(f"Found {len(files_to_fix)} files to fix")

    for filepath in sorted(files_to_fix):
        print(f"Processing {filepath}...")
        try:
            with open(filepath, 'r') as f:
                content = f.read()

            # TODO: Apply fixes based on error type
            # For now, just report
            print(f"  Would fix: {filepath}")
        except Exception as e:
            print(f"  Error processing {filepath}: {e}")

if __name__ == '__main__':
    main()
