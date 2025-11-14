#!/usr/bin/env python3
"""
N8N OpenAPI Download & Patch Pipeline
Downloads OpenAPI spec from n8n repo, bundles it, fixes aliases, and applies patches
"""

import subprocess
import sys
import shutil
import json
import argparse
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

def check_version():
    """Display version information from OpenAPI spec"""
    import yaml

    API_DIR = Path("sdk/n8nsdk/api")
    spec_file = API_DIR / "openapi.yaml"

    if not spec_file.exists():
        print("❌ OpenAPI file not found: sdk/n8nsdk/api/openapi.yaml")
        print("   Run 'make openapi' to generate it")
        sys.exit(1)

    try:
        with open(spec_file, 'r') as f:
            spec = yaml.safe_load(f)

        info = spec.get('info', {})
        version_info = info.get('x-n8n-version-info', {})

        if not version_info:
            print("⚠️  No version info found in OpenAPI spec")
            print("   Run 'make openapi' to regenerate with version tracking")
            sys.exit(1)

        frozen_version = version_info.get('frozenVersion', 'unknown')
        frozen_commit = version_info.get('frozenCommit', 'unknown')
        latest_version = version_info.get('latestVersion', 'unknown')
        in_sync = version_info.get('inSync', False)

        print("\n🔍 N8N Version Check\n")
        print(f"📌 Frozen Version:  {frozen_version}")
        print(f"🔒 Frozen Commit:   {frozen_commit[:8] if len(frozen_commit) > 8 else frozen_commit}")
        print(f"🆕 Latest Version:  {latest_version}")
        print()

        if in_sync:
            print("✅ Versions are in sync!")
        else:
            print("⚠️  Version mismatch detected!")
            print(f"   Consider updating by changing N8N_COMMIT in codegen/download-openapi.py")
            print(f"   From: {frozen_version} ({frozen_commit[:8]})")
            print(f"   To:   {latest_version}")
        print()

    except Exception as e:
        print(f"❌ Error: {e}")
        sys.exit(1)

def update_to_latest():
    """Update N8N_COMMIT to latest version"""
    print("\n🔄 Updating to Latest n8n Version\n")

    # Get latest version and commit
    latest_version = get_latest_version()
    if latest_version == 'unknown':
        print("❌ Failed to fetch latest version from GitHub")
        sys.exit(1)

    print(f"📥 Fetching commit hash for n8n@{latest_version}...")

    try:
        result = run(f"curl -s https://api.github.com/repos/n8n-io/n8n/git/refs/tags/n8n@{latest_version}", check=False)
        data = json.loads(result)
        latest_commit = data.get('object', {}).get('sha', None)

        if not latest_commit:
            print("❌ Failed to get commit hash for latest version")
            sys.exit(1)

        print(f"   ✓ Found commit: {latest_commit[:8]}\n")

    except Exception as e:
        print(f"❌ Error fetching commit: {e}")
        sys.exit(1)

    # Read current script
    script_path = Path(__file__)
    with open(script_path, 'r') as f:
        content = f.read()

    # Find and replace N8N_COMMIT line
    import re
    pattern = r'N8N_COMMIT = "[a-f0-9]+".*?# Frozen commit[^\n]*'
    replacement = f'N8N_COMMIT = "{latest_commit}"  # Frozen commit for API stability (n8n@{latest_version})'

    if not re.search(pattern, content):
        print("❌ Could not find N8N_COMMIT in script")
        sys.exit(1)

    new_content = re.sub(pattern, replacement, content)

    # Write back
    with open(script_path, 'w') as f:
        f.write(new_content)

    print(f"✅ Updated codegen/download-openapi.py")
    print(f"   📌 New commit: {latest_commit[:8]} (n8n@{latest_version})")
    print()
    print("Next steps:")
    print("  1. Run 'make openapi' to regenerate OpenAPI spec")
    print("  2. Run 'make sdk' to regenerate Go SDK")
    print("  3. Test the provider with the new version")
    print("  4. Commit changes if everything works")
    print()

def main():
    print("🚀 N8N OpenAPI Download Pipeline\n")

    # Config
    N8N_COMMIT = "53fd5c94c8798292cc981508a937b09532bbcf64"  # Frozen commit for API stability (n8n@1.119.1)
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

    # 4.5. Fix sharedWorkflow additionalProperties (before patch)
    print("🔧 Fixing sharedWorkflow schema...")
    with open(spec_file, 'r') as f:
        content = f.read()

    # Fix sharedWorkflow: additionalProperties: false -> true
    content = content.replace(
        '    sharedWorkflow:\n      type: object\n      additionalProperties: false',
        '    sharedWorkflow:\n      type: object\n      additionalProperties: true'
    )

    with open(spec_file, 'w') as f:
        f.write(content)

    print(f"   ✓ Fixed sharedWorkflow\n")

    # 5. Git commit the clean OpenAPI spec (before patch) - using temporary index
    print("💾 Committing clean OpenAPI spec...")
    commit_message = f"chore(sdk): bump n8n OpenAPI spec to {frozen_version}"

    # Save current index state, commit only openapi.yaml, restore index
    # This ensures we only commit openapi.yaml regardless of what's staged
    commit_script = f'''
    # Save what's currently staged
    git diff --cached > /tmp/staged_changes.patch 2>/dev/null || true
    # Reset index
    git reset --quiet HEAD 2>/dev/null || true
    # Add only openapi.yaml
    git add sdk/n8nsdk/api/openapi.yaml
    # Commit it (bypass hooks with --no-verify to avoid changelog/coverage regeneration)
    git commit --no-verify -m "{commit_message}" 2>&1
    commit_status=$?
    # Restore previously staged files (if any)
    if [ -s /tmp/staged_changes.patch ]; then
        git apply --cached /tmp/staged_changes.patch 2>/dev/null || true
    fi
    rm -f /tmp/staged_changes.patch
    exit $commit_status
    '''

    result = subprocess.run(
        commit_script,
        shell=True,
        capture_output=True,
        text=True,
        executable='/bin/bash'
    )
    if result.returncode == 0:
        print(f"   ✓ Committed: {commit_message}\n")
    else:
        # Check if there's nothing to commit
        if "nothing to commit" in result.stdout or "nothing to commit" in result.stderr:
            print("   ℹ No changes to commit (already up to date)\n")
        else:
            print("   ⚠️  Commit failed!")
            if result.stderr:
                print(f"      {result.stderr}")
            if result.stdout:
                print(f"      {result.stdout}\n")

    # 6. Apply patch
    print("🩹 Applying openapi.patch...")
    patch_file = API_DIR / "openapi.patch"
    if patch_file.exists():
        # Try with fuzzy matching to allow for line number differences
        result = subprocess.run(
            f"patch -p0 --fuzz=3 < {patch_file}",
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

    # 7. Add additionalProperties: true to credential schema
    print("🔧 Adding additionalProperties to credential schema...")
    with open(spec_file, 'r') as f:
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
        with open(spec_file, 'w') as f:
            f.writelines(output)
        print("   ✓ Added additionalProperties\n")
    else:
        print("   ⚠️  Credential schema not found\n")

    print("✅ OpenAPI spec ready!\n")
    print(f"   📄 {API_DIR}/openapi.yaml (bundled + aliases fixed + patched + versioned)")
    print(f"   📌 Based on n8n {frozen_version} (commit {N8N_COMMIT[:8]})")
    if frozen_version != latest_version and latest_version != 'unknown':
        print(f"   🆕 Latest available: {latest_version}")
    print("\nNext step: Run 'make sdk' to generate the Go SDK\n")

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='N8N OpenAPI Download & Version Check')
    parser.add_argument('--version', action='store_true', help='Check version information')
    parser.add_argument('--update', action='store_true', help='Update N8N_COMMIT to latest version')
    args = parser.parse_args()

    if args.version:
        check_version()
    elif args.update:
        update_to_latest()
    else:
        main()
