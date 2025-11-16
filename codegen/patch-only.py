#!/usr/bin/env python3
"""
N8N OpenAPI Patch Only
Apply patches to openapi.yaml
"""

import subprocess
import sys
import shlex
from pathlib import Path

def run(cmd, check=True):
    """Run command and optionally exit on error (secure version without shell=True)"""
    if isinstance(cmd, str):
        cmd = shlex.split(cmd)
    result = subprocess.run(cmd, shell=False, capture_output=True, text=True)
    if check and result.returncode != 0:
        print(f"❌ Command failed: {' '.join(cmd)}", file=sys.stderr)
        print(result.stderr, file=sys.stderr)
        sys.exit(1)
    return result.stdout.strip()

def main():
    print("🩹 N8N OpenAPI Patch Application\n")

    API_DIR = Path("sdk/n8nsdk/api")
    spec_file = API_DIR / "openapi.yaml"
    patch_file = API_DIR / "openapi.patch"

    if not spec_file.exists():
        print("❌ OpenAPI file not found: sdk/n8nsdk/api/openapi.yaml")
        print("   Run 'make openapi/download' first")
        sys.exit(1)

    # 1. Apply patch
    print("🩹 Applying openapi.patch...")
    if patch_file.exists():
        # Try with fuzzy matching to allow for line number differences
        # Use stdin parameter instead of shell redirection
        with open(patch_file, 'r') as patch_input:
            result = subprocess.run(
                ['patch', '-p0', '--fuzz=3'],
                stdin=patch_input,
                capture_output=True,
                text=True,
                check=False
            )
        if result.returncode != 0:
            print("   ⚠️  Patch failed!")
            if result.stderr:
                print(result.stderr)
            if result.stdout:
                print(result.stdout)
            sys.exit(1)
        else:
            print("   ✓ Patched\n")
    else:
        print("   ⚠️  No patch file found\n")

    # 2. Add additionalProperties: true to credential schema
    print("🔧 Adding additionalProperties to credential schema...")
    with open(spec_file, 'r', encoding='utf-8') as f:
        lines = f.readlines()

    modified = False
    output = []
    i = 0
    while i < len(lines):
        line = lines[i]
        output.append(line)

        # Look for credential schema definition
        if line.strip() == 'credential:':
            # Skip to 'type: object' line
            j = i + 1
            while j < len(lines):
                output.append(lines[j])
                if lines[j].strip() == 'type: object':
                    # Add additionalProperties after type: object
                    indent = len(lines[j]) - len(lines[j].lstrip())
                    output.append(' ' * indent + 'additionalProperties: true\n')
                    modified = True
                    i = j
                    break
                j += 1
        i += 1

    if modified:
        with open(spec_file, 'w', encoding='utf-8') as f:
            f.writelines(output)
        print("   ✓ Added additionalProperties\n")
    else:
        print("   ⚠️  Credential schema not found\n")

    print("✅ OpenAPI spec patched!\n")
    print(f"   📄 {API_DIR}/openapi.yaml")
    print("\nNext step: Run 'make sdk' to generate the Go SDK\n")

if __name__ == '__main__':
    main()
