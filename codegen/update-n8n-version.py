#!/usr/bin/env python3
"""
Update N8N version to latest release
"""

import subprocess
import sys
import json
import re
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

def get_latest_version():
    """Get latest n8n version from GitHub API"""
    try:
        result = run("curl -s https://api.github.com/repos/n8n-io/n8n/releases/latest", check=False)
        data = json.loads(result)
        return data.get('tag_name', 'unknown').lstrip('n8n@')
    except:
        return None

def main():
    print("🔄 Updating to Latest n8n Version\n")

    # Get latest version and commit
    latest_version = get_latest_version()
    if not latest_version:
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
    script_path = Path("codegen/download-only.py")
    with open(script_path, 'r') as f:
        content = f.read()

    # Find current version
    current_match = re.search(r'N8N_COMMIT = "([a-f0-9]+)".*?# Frozen commit.*?\(n8n@([\d.]+)\)', content)
    if not current_match:
        print("❌ Could not find N8N_COMMIT in download-only.py")
        sys.exit(1)

    current_commit = current_match.group(1)
    current_version = current_match.group(2)

    print(f"📌 Current version: {current_version} (commit {current_commit[:8]})")
    print(f"🆕 Latest version:  {latest_version} (commit {latest_commit[:8]})\n")

    if current_version == latest_version:
        print("✅ Already up to date!")
        sys.exit(0)

    # Replace N8N_COMMIT line
    pattern = r'N8N_COMMIT = "[a-f0-9]+".*?# Frozen commit[^\n]*'
    replacement = f'N8N_COMMIT = "{latest_commit}"  # Frozen commit for API stability (n8n@{latest_version})'

    new_content = re.sub(pattern, replacement, content)

    # Write back
    with open(script_path, 'w') as f:
        f.write(new_content)

    print(f"✅ Updated codegen/download-only.py")
    print(f"   📌 New commit: {latest_commit[:8]} (n8n@{latest_version})")
    print()
    print("Next steps:")
    print("  1. Run 'make sdk/openapi/download' to download new spec")
    print("  2. Run 'make sdk/openapi/patch' to apply patches")
    print("  3. Run 'make sdk' to regenerate Go SDK")
    print("  4. Test the provider with the new version")
    print("  5. Commit changes if everything works")
    print()

if __name__ == '__main__':
    main()
