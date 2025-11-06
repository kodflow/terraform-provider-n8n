#!/usr/bin/env python3
"""Script to automatically fix KTN-FUNC-007 errors by adding Params and Returns sections."""

import re
import subprocess
import sys

def get_func007_errors():
    """Get all KTN-FUNC-007 errors from ktn-linter."""
    result = subprocess.run(
        ['ktn-linter', 'lint', './...'],
        capture_output=True,
        text=True,
        cwd='/workspace/src'
    )

    # Parse errors: file:line:col
    errors = []
    lines = result.stdout.split('\n') + result.stderr.split('\n')
    for line in lines:
        if 'KTN-FUNC-007' in line:
            # Extract file path and line number
            match = re.search(r'/workspace/src/([^:]+):(\d+):\d+', line)
            if match:
                filepath = f'/workspace/src/{match.group(1)}'
                lineno = int(match.group(2))
                errors.append((filepath, lineno))

    return errors

def read_file(filepath):
    """Read file content."""
    with open(filepath, 'r', encoding='utf-8') as f:
        return f.readlines()

def write_file(filepath, lines):
    """Write file content."""
    with open(filepath, 'w', encoding='utf-8') as f:
        f.writelines(lines)

def extract_function_signature(lines, start_idx):
    """Extract function signature starting at start_idx."""
    # Find the function declaration
    func_line = lines[start_idx]

    # Check if it's a function
    if 'func ' not in func_line:
        return None, None, None

    # Extract receiver if present
    receiver_match = re.search(r'func\s+\(([^)]+)\)', func_line)
    receiver = receiver_match.group(1) if receiver_match else None

    # Extract function name and signature
    if receiver:
        sig_match = re.search(r'func\s+\([^)]+\)\s+(\w+)\s*\(([^)]*)\)\s*(.*)', func_line)
    else:
        sig_match = re.search(r'func\s+(\w+)\s*\(([^)]*)\)\s*(.*)', func_line)

    if not sig_match:
        return None, None, None

    func_name = sig_match.group(1)
    params_str = sig_match.group(2)
    returns_str = sig_match.group(3).strip()

    # Handle multi-line signatures
    if not returns_str.endswith('{'):
        idx = start_idx + 1
        while idx < len(lines) and '{' not in lines[idx]:
            returns_str += ' ' + lines[idx].strip()
            idx += 1
        if idx < len(lines):
            returns_str += ' ' + lines[idx].split('{')[0].strip()

    # Clean up returns
    returns_str = returns_str.replace('{', '').strip()

    return func_name, params_str, returns_str

def parse_params(params_str):
    """Parse parameter string into list of (name, type) tuples."""
    if not params_str.strip():
        return []

    params = []
    # Split by comma, but handle nested types
    parts = []
    level = 0
    current = []
    for char in params_str + ',':
        if char in '([{':
            level += 1
            current.append(char)
        elif char in ')]}':
            level -= 1
            current.append(char)
        elif char == ',' and level == 0:
            parts.append(''.join(current).strip())
            current = []
        else:
            current.append(char)

    for part in parts:
        if not part:
            continue
        # Handle "name type" or just "type"
        tokens = part.strip().split()
        if len(tokens) >= 2:
            # Last token is type, everything before is name
            name = ' '.join(tokens[:-1])
            typ = tokens[-1]
            params.append((name, typ))
        else:
            params.append(('', part))

    return params

def parse_returns(returns_str):
    """Parse return types."""
    if not returns_str:
        return []

    returns_str = returns_str.strip()

    # Handle (type1, type2) format
    if returns_str.startswith('(') and returns_str.endswith(')'):
        inner = returns_str[1:-1]
        # Split by comma
        types = []
        level = 0
        current = []
        for char in inner + ',':
            if char in '([{':
                level += 1
                current.append(char)
            elif char in ')]}':
                level -= 1
                current.append(char)
            elif char == ',' and level == 0:
                types.append(''.join(current).strip())
                current = []
            else:
                current.append(char)
        return [t for t in types if t]
    else:
        return [returns_str] if returns_str else []

def generate_doc_comment(func_name, params, returns):
    """Generate Params and Returns documentation."""
    doc_lines = []
    doc_lines.append('//\n')
    doc_lines.append('// Params:\n')

    if not params or all(not name for name, _ in params):
        doc_lines.append('//   - None\n')
    else:
        for name, typ in params:
            if name:
                # Generate description based on common patterns
                if 'ctx' in name.lower() or typ == 'context.Context':
                    desc = 'The context for the request'
                elif 'req' in name.lower():
                    desc = 'The request object'
                elif 'resp' in name.lower():
                    desc = 'The response object'
                elif 'data' in name.lower():
                    desc = 'Data model for the operation'
                elif 'client' in name.lower():
                    desc = 'API client instance'
                elif name == 'cmd':
                    desc = 'The cobra command being executed'
                elif name == 'args':
                    desc = 'Command line arguments'
                elif name == 'v':
                    desc = 'The version string to set'
                else:
                    desc = f'Parameter of type {typ}'
                doc_lines.append(f'//   - {name}: {desc}\n')

    doc_lines.append('//\n')
    doc_lines.append('// Returns:\n')

    if not returns:
        doc_lines.append('//   - None\n')
    else:
        for ret_type in returns:
            if ret_type == 'error':
                desc = 'Error if operation fails'
            elif 'datasource.DataSource' in ret_type:
                desc = 'New data source instance'
            elif 'resource.Resource' in ret_type:
                desc = 'New resource instance'
            elif '[]' in ret_type and 'datasource.DataSource' in ret_type:
                desc = 'Slice of data source factory functions'
            elif '[]' in ret_type and 'resource.Resource' in ret_type:
                desc = 'Slice of resource factory functions'
            else:
                desc = f'Value of type {ret_type}'
            doc_lines.append(f'//   - {ret_type}: {desc}\n')

    return doc_lines

def fix_function_doc(lines, func_line_idx):
    """Fix documentation for function at func_line_idx."""
    # Find existing doc comment
    doc_start_idx = func_line_idx - 1
    while doc_start_idx >= 0 and lines[doc_start_idx].strip().startswith('//'):
        doc_start_idx -= 1
    doc_start_idx += 1

    # Extract function signature
    func_name, params_str, returns_str = extract_function_signature(lines, func_line_idx)
    if not func_name:
        return lines

    # Parse params and returns
    params = parse_params(params_str)
    returns = parse_returns(returns_str)

    # Generate new doc comment
    new_doc = generate_doc_comment(func_name, params, returns)

    # Check if Params/Returns already exist
    has_params = False
    has_returns = False
    for i in range(doc_start_idx, func_line_idx):
        line = lines[i]
        if '// Params:' in line or '//Params:' in line:
            has_params = True
        if '// Returns:' in line or '//Returns:' in line:
            has_returns = True

    if has_params and has_returns:
        # Already has both sections
        return lines

    # Find where to insert (after first line of doc comment)
    insert_idx = doc_start_idx
    # Skip first doc line
    if insert_idx < func_line_idx and lines[insert_idx].strip().startswith('//'):
        insert_idx += 1

    # Insert new doc
    return lines[:insert_idx] + new_doc + lines[insert_idx:]

def main():
    """Main function."""
    print("Fetching KTN-FUNC-007 errors...")
    errors = get_func007_errors()

    if not errors:
        print("No KTN-FUNC-007 errors found!")
        return 0

    print(f"Found {len(errors)} KTN-FUNC-007 errors")

    # Group by file
    files_to_fix = {}
    for filepath, lineno in errors:
        if filepath not in files_to_fix:
            files_to_fix[filepath] = []
        files_to_fix[filepath].append(lineno)

    print(f"Fixing {len(files_to_fix)} files...")

    for filepath, line_numbers in files_to_fix.items():
        print(f"Fixing {filepath}...")
        lines = read_file(filepath)

        # Sort line numbers in reverse to fix from bottom to top
        line_numbers.sort(reverse=True)

        for lineno in line_numbers:
            # Convert to 0-based index
            idx = lineno - 1
            if idx < 0 or idx >= len(lines):
                continue

            lines = fix_function_doc(lines, idx)

        write_file(filepath, lines)

    print("Done! Re-running linter to verify...")
    result = subprocess.run(
        ['ktn-linter', 'lint', './...'],
        capture_output=True,
        text=True,
        cwd='/workspace/src'
    )

    # Count remaining errors
    remaining = result.stdout.count('KTN-FUNC-007') + result.stderr.count('KTN-FUNC-007')
    print(f"Remaining KTN-FUNC-007 errors: {remaining}")

    return 0 if remaining == 0 else 1

if __name__ == '__main__':
    sys.exit(main())
