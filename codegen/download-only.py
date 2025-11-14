#!/usr/bin/env python3
"""
N8N OpenAPI Download Only
Downloads OpenAPI spec from n8n repo, bundles it, and fixes aliases
NO PATCHING - NO COMMITS
"""

import subprocess
import sys
import shutil
import json
from pathlib import Path

def run(cmd, cwd=None, check=True):
    """Run command and optionally exit on error"""
    result = subprocess.run(cmd, shell=True, cwd=cwd, capture_output=True, text=True)
    if check and result.returncode != 0:
        print(f"❌ Command failed: {cmd}", file=sys.stderr)
        print(result.stderr, file=sys.stderr)
        sys.exit(1)
    return result.stdout.strip()

def get_n8n_version(temp_dir, commit):
    """Extract n8n version from package.json at specific commit"""
    try:
        package_json = Path(temp_dir) / "package.json"
        if package_json.exists():
            with open(package_json, 'r') as f:
                data = json.load(f)
                return data.get('version', 'unknown')
    except:
        pass
    return 'unknown'

def get_latest_version():
    """Get latest n8n version from GitHub API"""
    try:
        result = run("curl -s https://api.github.com/repos/n8n-io/n8n/releases/latest", check=False)
        data = json.loads(result)
        return data.get('tag_name', 'unknown').lstrip('n8n@')
    except:
        return 'unknown'

def main():
    print("🚀 N8N OpenAPI Download (No Commit)\n")

    # Config
    N8N_COMMIT = "4bf741ae67124724eb94e582de94daf0d70f9bd0"  # Frozen commit for API stability (n8n@1.119.2)
    TEMP_DIR = "/tmp/n8n-openapi-download"
    API_DIR = Path("sdk/n8nsdk/api")

    # 1. Download from GitHub
    print("📥 Downloading OpenAPI from GitHub...")
    print(f"   🔒 Using frozen commit: {N8N_COMMIT[:8]}\n")

    if Path(TEMP_DIR).exists():
        shutil.rmtree(TEMP_DIR)

    run(f"git clone --depth 1 --filter=blob:none --sparse https://github.com/n8n-io/n8n.git {TEMP_DIR}")
    run("git sparse-checkout set packages/cli/src/public-api", cwd=TEMP_DIR)
    run(f"git checkout {N8N_COMMIT}", cwd=TEMP_DIR)

    # Get package.json separately to extract version
    run(f"git checkout {N8N_COMMIT} -- package.json", cwd=TEMP_DIR, check=False)

    # Get version from the frozen commit
    frozen_version = get_n8n_version(TEMP_DIR, N8N_COMMIT)
    print(f"   📌 Frozen n8n version: {frozen_version}")

    # Get latest version from GitHub
    latest_version = get_latest_version()
    print(f"   🆕 Latest n8n version: {latest_version}")

    if frozen_version != latest_version and frozen_version != 'unknown' and latest_version != 'unknown':
        print(f"   ⚠️  Version mismatch detected!")
        print(f"      Consider updating to latest by changing N8N_COMMIT\n")
    else:
        print(f"   ✓ Versions in sync\n")

    source_path = API_DIR / "openapi-source"
    if source_path.exists():
        # Force remove with chmod to handle permission issues
        run(f"chmod -R u+w {source_path}")
        shutil.rmtree(source_path)
    shutil.copytree(f"{TEMP_DIR}/packages/cli/src/public-api", source_path)
    shutil.rmtree(TEMP_DIR)
    print("   ✓ Downloaded\n")

    # 2. Bundle YAML
    print("📦 Bundling YAML files...")
    run(f"npx --yes @redocly/cli@latest bundle {source_path}/v1/openapi.yml -o {API_DIR}/openapi.yaml")
    print("   ✓ Bundled\n")

    # 3. Fix schema aliases
    import yaml
    import re
    spec_file = API_DIR / "openapi.yaml"

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

    # 4. Add version info to OpenAPI spec
    print("📝 Adding version information...")
    with open(spec_file, 'r') as f:
        openapi_content = f.read()

    # Add x-n8n-version-info extension to info section
    import yaml
    spec = yaml.safe_load(openapi_content)

    if 'info' not in spec:
        spec['info'] = {}

    spec['info']['x-n8n-version-info'] = {
        'frozenVersion': frozen_version,
        'frozenCommit': N8N_COMMIT,
        'latestVersion': latest_version,
        'inSync': frozen_version == latest_version,
        'note': 'This OpenAPI spec is generated from a frozen commit for API stability. Check latestVersion to see if updates are available.'
    }

    # Write back with proper indentation (2 spaces for nested content, 2 spaces offset for sequences)
    with open(spec_file, 'w') as f:
        yaml.dump(spec, f, default_flow_style=False, sort_keys=False, allow_unicode=True, indent=2, width=float("inf"))

    print(f"   ✓ Added version info\n")

    print("✅ OpenAPI spec downloaded and prepared!\n")
    print(f"   📄 {API_DIR}/openapi.yaml (bundled + aliases fixed + versioned)")
    print(f"   📌 Based on n8n {frozen_version} (commit {N8N_COMMIT[:8]})")
    if frozen_version != latest_version and latest_version != 'unknown':
        print(f"   🆕 Latest available: {latest_version}")
    print("\nNext step: Run 'make openapi/patch' to apply custom patches\n")

if __name__ == '__main__':
    main()
