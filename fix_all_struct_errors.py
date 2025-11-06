#!/usr/bin/env python3
"""
Comprehensive script to fix ALL KTN-STRUCT errors in provider files.
Handles datasources and resources.
"""

import re
import os
import glob

def extract_struct_info(content):
    """Extract main struct names and their types."""
    structs = []

    # Find all exported struct definitions
    pattern = r'type\s+(\w+(?:DataSource|Resource)(?:Model)?)\s+struct'
    matches = re.finditer(pattern, content)

    for match in matches:
        structs.append(match.group(1))

    return structs

def get_resource_type(struct_name):
    """Determine if this is a DataSource or Resource."""
    if 'DataSource' in struct_name:
        return 'datasource'
    elif 'Resource' in struct_name:
        return 'resource'
    return None

def generate_interface(struct_name, resource_type):
    """Generate interface definition for a struct."""
    if resource_type == 'datasource':
        return f'''var _ {struct_name}Interface = &{struct_name}{{}}

// {struct_name}Interface defines the complete interface for {struct_name}.
// This interface includes all public methods required for the data source implementation.
type {struct_name}Interface interface {{
\tdatasource.DataSource
\tdatasource.DataSourceWithConfigure
}}

'''
    elif resource_type == 'resource':
        return f'''var _ {struct_name}Interface = &{struct_name}{{}}

// {struct_name}Interface defines the complete interface for {struct_name}.
// This interface includes all public methods required for the resource implementation.
type {struct_name}Interface interface {{
\tresource.Resource
\tresource.ResourceWithConfigure
\tresource.ResourceWithImportState
}}

'''
    return ''

def enhance_struct_doc(content, struct_name):
    """Enhance single-line struct documentation to multi-line."""
    # Find current documentation
    pattern = r'(//\s*' + re.escape(struct_name) + r'.*?\n)(type\s+' + re.escape(struct_name) + r'\s+struct)'

    # Determine struct type
    if 'Model' in struct_name:
        doc = f'''// {struct_name} describes the data model for this resource.
// This model represents the complete schema including all fields required
// for managing the resource through Terraform and the n8n API.
'''
    elif 'DataSource' in struct_name:
        base = struct_name.replace('DataSource', '')
        doc = f'''// {struct_name} defines the data source implementation.
// This data source allows fetching information about n8n {base.lower()}
// resources using the n8n API.
'''
    elif 'Resource' in struct_name:
        base = struct_name.replace('Resource', '')
        doc = f'''// {struct_name} defines the resource implementation.
// This resource manages n8n {base.lower()} lifecycle including
// creation, updates, deletion, and state management.
'''
    else:
        return content

    # Replace single-line doc with multi-line
    return re.sub(pattern, doc + r'\2', content)

def fix_file(filepath):
    """Fix all KTN-STRUCT errors in a single file."""
    print(f"Processing {filepath}...")

    with open(filepath, 'r') as f:
        content = f.read()

    # Extract all structs
    structs = extract_struct_info(content)
    if not structs:
        print(f"  No structs found in {filepath}")
        return False

    main_struct = structs[0]  # The main struct (DataSource or Resource)
    resource_type = get_resource_type(main_struct)

    if not resource_type:
        print(f"  Could not determine type for {filepath}")
        return False

    # Step 1: Add interface if missing
    if f'{main_struct}Interface' not in content:
        # Find the position after var _ declarations
        if resource_type == 'datasource':
            pattern = r'(var\s+_\s+datasource\.DataSourceWithConfigure\s+=\s+&' + re.escape(main_struct) + r'\{\})'
        else:
            pattern = r'(var\s+_\s+resource\.ResourceWithImportState\s+=\s+&' + re.escape(main_struct) + r'\{\})'

        interface_code = generate_interface(main_struct, resource_type)
        content = re.sub(pattern + r'\n', r'\1\n' + interface_code, content)
        print(f"  Added interface for {main_struct}")

    # Step 2: Enhance documentation for all structs
    for struct in structs:
        old_content = content
        content = enhance_struct_doc(content, struct)
        if content != old_content:
            print(f"  Enhanced documentation for {struct}")

    # Write back
    with open(filepath, 'w') as f:
        f.write(content)

    return True

def main():
    # Find all datasource and resource files
    datasource_files = glob.glob('/workspace/src/internal/provider/datasources/*.go')
    resource_files = glob.glob('/workspace/src/internal/provider/resources/*.go')

    all_files = datasource_files + resource_files

    # Exclude datasources.go and resources.go (they are factory files)
    all_files = [f for f in all_files if not f.endswith(('datasources.go', 'resources.go'))]

    print(f"Found {len(all_files)} files to process")

    fixed = 0
    for filepath in all_files:
        if fix_file(filepath):
            fixed += 1

    print(f"\nProcessed {fixed}/{len(all_files)} files")

if __name__ == '__main__':
    main()
