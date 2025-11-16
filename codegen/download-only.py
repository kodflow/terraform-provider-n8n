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
import shlex
import tempfile
import re
import yaml
from pathlib import Path

def run(cmd, cwd=None, check=True):
    """Run command and optionally exit on error (secure version without shell=True)"""
    if isinstance(cmd, str):
        cmd = shlex.split(cmd)
    result = subprocess.run(cmd, shell=False, cwd=cwd, capture_output=True, text=True, check=False)  # nosec B603
    if check and result.returncode != 0:
        print(f"❌ Command failed: {' '.join(cmd)}", file=sys.stderr)
        print(result.stderr, file=sys.stderr)
        sys.exit(1)
    return result.stdout.strip()

def get_n8n_version(temp_dir, commit):
    """Extract n8n version from package.json at specific commit"""
    try:
        package_json = Path(temp_dir) / "package.json"
        if package_json.exists():
            with open(package_json, 'r', encoding='utf-8') as f:
                data = json.load(f)
                return data.get('version', 'unknown')
    except (json.JSONDecodeError, IOError, KeyError):
        pass
    return 'unknown'

def get_latest_version():
    """Get latest n8n version from GitHub API"""
    try:
        result = run("curl -s https://api.github.com/repos/n8n-io/n8n/releases/latest", check=False)
        data = json.loads(result)
        return data.get('tag_name', 'unknown').lstrip('n8n@')
    except (json.JSONDecodeError, KeyError, TypeError):
        return 'unknown'

def download_from_github(n8n_commit, temp_dir, api_dir):
    """Download OpenAPI spec from GitHub and return version info"""
    print("📥 Downloading OpenAPI from GitHub...")
    print(f"   🔒 Using frozen commit: {n8n_commit[:8]}\n")

    if Path(temp_dir).exists():
        shutil.rmtree(temp_dir)

    # Long line: git clone command
    run(
        f"git clone --depth 1 --filter=blob:none --sparse "
        f"https://github.com/n8n-io/n8n.git {temp_dir}"
    )
    run("git sparse-checkout set packages/cli/src/public-api", cwd=temp_dir)
    run(f"git checkout {n8n_commit}", cwd=temp_dir)

    # Get package.json separately to extract version
    run(f"git checkout {n8n_commit} -- package.json", cwd=temp_dir, check=False)

    # Get version from the frozen commit
    frozen_version = get_n8n_version(temp_dir, n8n_commit)
    print(f"   📌 Frozen n8n version: {frozen_version}")

    # Get latest version from GitHub
    latest_version = get_latest_version()
    print(f"   🆕 Latest n8n version: {latest_version}")

    if frozen_version != latest_version and frozen_version != 'unknown' and latest_version != 'unknown':
        print("   ⚠️  Version mismatch detected!")
        print("      Consider updating to latest by changing N8N_COMMIT\n")
    else:
        print("   ✓ Versions in sync\n")

    source_path = api_dir / "openapi-source"
    if source_path.exists():
        # Force remove with chmod to handle permission issues
        run(f"chmod -R u+w {source_path}")
        shutil.rmtree(source_path)
    shutil.copytree(f"{temp_dir}/packages/cli/src/public-api", source_path)
    shutil.rmtree(temp_dir)
    print("   ✓ Downloaded\n")

    return frozen_version, latest_version

def bundle_yaml_spec(source_path, api_dir):
    """Bundle YAML files into single OpenAPI spec"""
    print("📦 Bundling YAML files...")
    # Long line: npx bundling command
    run(
        f"npx --yes @redocly/cli@latest bundle {source_path}/v1/openapi.yml "
        f"-o {api_dir}/openapi.yaml"
    )
    print("   ✓ Bundled\n")

def fix_schema_aliases(spec_file):
    """Fix schema aliases in OpenAPI spec"""
    print("🔧 Fixing schema aliases...")

    # Read YAML to find aliases
    with open(spec_file, 'r', encoding='utf-8') as f:
        spec = yaml.safe_load(f)

    schemas = spec.get('components', {}).get('schemas', {})
    aliases = [(name, def_['$ref'].split('/')[-1])
               for name, def_ in schemas.items()
               if isinstance(def_, dict) and len(def_) == 1 and '$ref' in def_]

    # Read as text to preserve formatting
    with open(spec_file, 'r', encoding='utf-8') as f:
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

    with open(spec_file, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f"   ✓ Fixed {len(aliases)} aliases\n")

def add_version_info(spec_file, frozen_version, n8n_commit, latest_version):
    """Add version information to OpenAPI spec"""
    print("📝 Adding version information...")
    with open(spec_file, 'r', encoding='utf-8') as f:
        openapi_content = f.read()

    # Add x-n8n-version-info extension to info section
    spec = yaml.safe_load(openapi_content)

    if 'info' not in spec:
        spec['info'] = {}

    spec['info']['x-n8n-version-info'] = {
        'frozenVersion': frozen_version,
        'frozenCommit': n8n_commit,
        'latestVersion': latest_version,
        'inSync': frozen_version == latest_version,
        'note': 'This OpenAPI spec is generated from a frozen commit for API stability. Check latestVersion to see if updates are available.'
    }

    # Write back with proper indentation (2 spaces for nested content, 2 spaces offset for sequences)
    with open(spec_file, 'w', encoding='utf-8') as f:
        yaml.dump(spec, f, default_flow_style=False, sort_keys=False, allow_unicode=True, indent=2, width=float("inf"))

    print("   ✓ Added version info\n")

def main():
    print("🚀 N8N OpenAPI Download (No Commit)\n")

    # Config
    n8n_commit = "4bf741ae67124724eb94e582de94daf0d70f9bd0"  # Frozen commit for API stability (n8n@1.119.2)
    temp_dir = tempfile.mkdtemp(prefix="n8n-openapi-")
    api_dir = Path("sdk/n8nsdk/api")

    try:
        # 1. Download from GitHub
        frozen_version, latest_version = download_from_github(n8n_commit, temp_dir, api_dir)

        # 2. Bundle YAML
        bundle_yaml_spec(api_dir / "openapi-source", api_dir)

        # 3. Fix schema aliases
        spec_file = api_dir / "openapi.yaml"
        fix_schema_aliases(spec_file)

        # 4. Add version info to OpenAPI spec
        add_version_info(spec_file, frozen_version, n8n_commit, latest_version)

        # Final output
        print("✅ OpenAPI spec downloaded and prepared!\n")
        print(f"   📄 {api_dir}/openapi.yaml (bundled + aliases fixed + versioned)")
        print(f"   📌 Based on n8n {frozen_version} (commit {n8n_commit[:8]})")
        if frozen_version != latest_version and latest_version != 'unknown':
            print(f"   🆕 Latest available: {latest_version}")
        print("\nNext step: Run 'make openapi/patch' to apply custom patches\n")
    finally:
        # Clean up temporary directory
        if Path(temp_dir).exists():
            shutil.rmtree(temp_dir)

if __name__ == '__main__':
    main()
