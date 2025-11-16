#!/usr/bin/env python3
"""
OpenAPI Patch Management
Create and apply patches to openapi.yaml
"""

import sys
import shutil
import subprocess
import shlex
import yaml
from pathlib import Path

def run(cmd):
    """Run command and return output (secure version without shell=True)"""
    if isinstance(cmd, str):
        cmd = shlex.split(cmd)
    result = subprocess.run(cmd, shell=False, capture_output=True, text=True, check=False)  # nosec B603
    return result.returncode, result.stdout, result.stderr

def create_patch():
    """Create a patch from current git diff"""
    print("\n📝 Creating OpenAPI Patch\n")

    openapi_file = Path("sdk/n8nsdk/api/openapi.yaml")
    patch_file = Path("sdk/n8nsdk/api/openapi.patch")

    if not openapi_file.exists():
        print("❌ openapi.yaml not found")
        print("   Run 'make openapi' first")
        sys.exit(1)

    # Check if there are changes
    rc, stdout, stderr = run("git diff sdk/n8nsdk/api/openapi.yaml")
    if rc != 0:
        print(f"❌ Git error: {stderr}")
        sys.exit(1)

    if not stdout.strip():
        print("⚠️  No changes detected in openapi.yaml")
        print("   Make your modifications first")
        sys.exit(1)

    # Extract only the diff part (remove git headers)
    lines = stdout.split('\n')
    diff_start = None
    for i, line in enumerate(lines):
        if line.startswith('---'):
            diff_start = i
            break

    if diff_start is None:
        print("❌ Could not parse git diff output")
        sys.exit(1)

    # Fix the paths for patch -p0 format
    patch_content = []
    for line in lines[diff_start:]:
        if line.startswith('--- a/'):
            line = '--- ' + line[6:]
        elif line.startswith('+++ b/'):
            line = '+++ ' + line[6:]
        patch_content.append(line)

    # Write patch file
    with open(patch_file, 'w', encoding='utf-8') as f:
        f.write('\n'.join(patch_content))

    print(f"✅ Patch created: {patch_file}")
    print(f"   {len([l for l in patch_content if l.startswith('+') and not l.startswith('+++')])} additions")
    print(f"   {len([l for l in patch_content if l.startswith('-') and not l.startswith('---')])} deletions")
    print()

def apply_git_commit_patch(commit_hash):
    """Apply patch from a specific git commit"""
    print(f"\n🔄 Applying Patch from Commit {commit_hash[:8]}\n")

    openapi_file = Path("sdk/n8nsdk/api/openapi.yaml")
    patch_file = Path("sdk/n8nsdk/api/openapi.patch")

    if not openapi_file.exists():
        print("❌ openapi.yaml not found")
        print("   Run 'make openapi' first")
        sys.exit(1)

    # Get the diff from commit
    print("📥 Extracting changes from commit...")
    rc, stdout, stderr = run(f"git show {commit_hash} sdk/n8nsdk/api/openapi.yaml")
    if rc != 0:
        print(f"❌ Git error: {stderr}")
        sys.exit(1)

    # Extract only the diff part
    lines = stdout.split('\n')
    diff_start = None
    for i, line in enumerate(lines):
        if line.startswith('---'):
            diff_start = i
            break

    if diff_start is None:
        print("❌ Could not parse git show output")
        sys.exit(1)

    # Fix the paths for patch -p0 format
    patch_content = []
    for line in lines[diff_start:]:
        if line.startswith('--- a/'):
            line = '--- ' + line[6:]
        elif line.startswith('+++ b/'):
            line = '+++ ' + line[6:]
        patch_content.append(line)

    # Write patch file
    with open(patch_file, 'w', encoding='utf-8') as f:
        f.write('\n'.join(patch_content))

    print("   ✓ Extracted changes")
    print(f"   📝 {len([l for l in patch_content if l.startswith('+') and not l.startswith('+++')])} additions")
    print(f"   📝 {len([l for l in patch_content if l.startswith('-') and not l.startswith('---')])} deletions")
    print()

    # Apply the patch to current file
    print("🩹 Applying patch to openapi.yaml...")
    backup_file = openapi_file.with_suffix('.yaml.backup')
    shutil.copy(openapi_file, backup_file)

    # Try with fuzzy matching (allows line number differences)
    # Use stdin parameter instead of shell redirection
    with open(patch_file, 'r', encoding='utf-8') as patch_input:
        result = subprocess.run(  # nosec B603 B607
            ['patch', '-p0', '--fuzz=3'],
            stdin=patch_input,
            capture_output=True,
            text=True,
            check=False
        )
    if result.returncode:
        print("❌ Patch failed!")
        print(result.stderr)
        # Restore backup
        shutil.copy(backup_file, openapi_file)
        backup_file.unlink()
        sys.exit(1)

    backup_file.unlink()
    print("   ✓ Patch applied successfully")
    print()

    # Validate YAML
    print("✅ Validating YAML...")
    try:
        with open(openapi_file, 'r', encoding='utf-8') as f:
            yaml.safe_load(f)
        print("   ✓ YAML is valid")
    except (yaml.YAMLError, IOError) as e:
        print(f"❌ YAML validation failed: {e}")
        sys.exit(1)

    print()
    print("✅ Patch from commit applied successfully!")
    print(f"   📄 Patch saved to: {patch_file}")
    print()
    print("Next steps:")
    print("  1. Test with 'make openapi' to ensure patch applies cleanly")
    print("  2. Commit the patch file: git add sdk/n8nsdk/api/openapi.patch")
    print()

if __name__ == '__main__':
    import argparse

    parser = argparse.ArgumentParser(description='OpenAPI Patch Management')
    parser.add_argument('--create', action='store_true', help='Create patch from git diff')
    parser.add_argument('--from-commit', metavar='HASH', help='Create patch from specific commit')
    args = parser.parse_args()

    if args.create:
        create_patch()
    elif args.from_commit:
        apply_git_commit_patch(args.from_commit)
    else:
        parser.print_help()
