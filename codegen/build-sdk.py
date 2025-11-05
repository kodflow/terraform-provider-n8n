#!/usr/bin/env python3
"""
N8N SDK Build Pipeline
Simple orchestration script for SDK generation
"""

import subprocess
import sys
import shutil
from pathlib import Path

def run(cmd, cwd=None):
    """Run command and exit on error"""
    result = subprocess.run(cmd, shell=True, cwd=cwd, capture_output=True, text=True)
    if result.returncode != 0:
        print(f"❌ Command failed: {cmd}", file=sys.stderr)
        print(result.stderr, file=sys.stderr)
        sys.exit(1)
    return result.stdout

def main():
    print("🚀 N8N SDK Build Pipeline\n")

    # Config
    N8N_COMMIT = "0ccf47044a2ba5b94140bfdd2ba36b868091288d"
    TEMP_DIR = "/tmp/n8n-openapi-download"
    API_DIR = Path("sdk/n8nsdk/api")

    # 1. Download from GitHub
    print("📥 Downloading OpenAPI from GitHub...")
    if Path(TEMP_DIR).exists():
        shutil.rmtree(TEMP_DIR)

    run(f"git clone --depth 1 --filter=blob:none --sparse https://github.com/n8n-io/n8n.git {TEMP_DIR}")
    run("git sparse-checkout set packages/cli/src/public-api", cwd=TEMP_DIR)
    run(f"git checkout {N8N_COMMIT}", cwd=TEMP_DIR)

    source_path = API_DIR / "openapi-source"
    if source_path.exists():
        shutil.rmtree(source_path)
    shutil.copytree(f"{TEMP_DIR}/packages/cli/src/public-api", source_path)
    shutil.rmtree(TEMP_DIR)
    print("   ✓ Downloaded\n")

    # 2. Bundle YAML
    print("📦 Bundling YAML files...")
    run(f"npx --yes @redocly/cli@latest bundle {source_path}/v1/openapi.yml -o {API_DIR}/openapi.yaml")
    print("   ✓ Bundled\n")

    # 3. Apply patch
    import yaml
    import re
    spec_file = API_DIR / "openapi.yaml"

    print("🩹 Applying openapi.patch...")
    patch_file = API_DIR / "openapi.patch"
    if patch_file.exists():
        result = subprocess.run(
            f"patch -p0 < {patch_file}",
            shell=True,
            capture_output=True,
            text=True
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

    # 4. Fix schema aliases
    print("🔧 Fixing schema aliases...")

    # Read YAML to find aliases
    with open(spec_file, 'r') as f:
        spec = yaml.safe_load(f)

    schemas = spec.get('components', {}).get('schemas', {})
    aliases = [(name, def_['$ref'].split('/')[-1])
               for name, def_ in schemas.items()
               if isinstance(def_, dict) and len(def_) == 1 and '$ref' in def_]

    # Read as text to preserve formatting
    with open(spec_file, 'r') as f:
        content = f.read()

    # Replace each alias with its target schema content
    for alias_name, target_name in aliases:
        if target_name in schemas:
            # Find the target schema in text using DOTALL and negative lookahead
            # Matches from "    schemaname:\n" until the next schema at same indentation level
            target_pattern = rf"(\n    {re.escape(target_name)}:\n)((?:(?!\n    [a-zA-Z]).)*)"
            target_match = re.search(target_pattern, content, re.DOTALL)

            if target_match:
                target_content = target_match.group(2)

                # Replace the alias $ref with the target content
                # Note: target_content already starts with proper indentation
                alias_pattern = rf"(\n    {re.escape(alias_name)}:\n      \$ref: '#/components/schemas/{re.escape(target_name)}')"
                content = re.sub(alias_pattern, f"\n    {alias_name}:\n{target_content}", content, count=1)

    with open(spec_file, 'w') as f:
        f.write(content)
    print(f"   ✓ Fixed {len(aliases)} aliases\n")

    # 5. Generate SDK
    print("🔨 Generating SDK...")
    run("bash codegen/generate-sdk.sh")
    print("   ✓ Generated\n")

    # 6. Generate Bazel files
    print("🏗️  Generating BUILD files...")
    run("bazel run //:gazelle")
    print("   ✓ Done\n")

    print("✅ SDK build complete!\n")

if __name__ == '__main__':
    main()
