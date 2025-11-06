#!/usr/bin/env python3
"""
Script to fix KTN-STRUCT errors in datasource files by adding:
1. Complete interface definitions (KTN-STRUCT-002)
2. Multi-line documentation for structs (KTN-STRUCT-004)
3. Interface implementation assertions (KTN-STRUCT-005)
"""

import re
import sys

def fix_datasource_file(filepath):
    """Fix a single datasource file."""
    with open(filepath, 'r') as f:
        content = f.read()

    # Extract the struct name from the file (e.g., "VariableDataSource")
    match = re.search(r'type (\w+DataSource(?:Model)?) struct', content)
    if not match:
        print(f"Could not find main struct in {filepath}")
        return False

    main_struct = match.group(1)
    base_name = main_struct.replace('DataSource', '')

    # Pattern 1: Add interface definition after var _ declarations
    interface_pattern = r'(var _ datasource\.DataSourceWithConfigure = &' + re.escape(main_struct) + r'\{\})\n\n(//.*?\ntype ' + re.escape(main_struct) + r' struct)'

    interface_code = f'''var _ {main_struct}Interface = &{main_struct}{{}}

// {main_struct}Interface defines the complete interface for {main_struct}.
// This interface includes all public methods required for the data source implementation.
type {main_struct}Interface interface {{
\tdatasource.DataSource
\tdatasource.DataSourceWithConfigure
}}

'''

    # Check if interface doesn't already exist
    if f'{main_struct}Interface' not in content:
        content = re.sub(interface_pattern, r'\1\n' + interface_code + r'\2', content)

    # Pattern 2: Enhance single-line struct documentation to multi-line
    doc_patterns = [
        (r'// ' + re.escape(main_struct) + r' defines the data source implementation.*?\n(type ' + re.escape(main_struct) + r' struct)',
         f'''// {main_struct} defines the data source implementation for {base_name.lower()}.
// This data source allows fetching detailed information about specific n8n {base_name.lower()}
// resources using unique identifiers or other query parameters.
type {main_struct} struct'''),
    ]

    for pattern, replacement in doc_patterns:
        if re.search(pattern, content):
            content = re.sub(pattern, replacement, content)

    # Pattern 3: Fix Model struct documentation
    model_struct = main_struct + 'Model'
    if model_struct in content:
        model_doc_pattern = r'// ' + re.escape(model_struct) + r' describes.*?\n(type ' + re.escape(model_struct) + r' struct)'
        model_doc_replacement = f'''// {model_struct} describes the data source data model.
// This model represents the schema for {base_name.lower()} data including all fields
// returned by the n8n API and required for Terraform state management.
type {model_struct} struct'''

        if re.search(model_doc_pattern, content):
            content = re.sub(model_doc_pattern, model_doc_replacement, content)

    # Write the updated content
    with open(filepath, 'w') as f:
        f.write(content)

    print(f"Fixed {filepath}")
    return True

# Fix all datasource files
files = [
    '/workspace/src/internal/provider/datasources/variable.go',
    '/workspace/src/internal/provider/datasources/variables.go',
    '/workspace/src/internal/provider/datasources/workflow.go',
    '/workspace/src/internal/provider/datasources/workflows.go',
]

for filepath in files:
    fix_datasource_file(filepath)

print("Done!")
