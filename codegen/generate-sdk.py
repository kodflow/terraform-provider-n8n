#!/usr/bin/env python3
"""
N8N SDK Generation Pipeline
Generates Go SDK from OpenAPI spec using openapi-generator
"""

import subprocess
import sys
import shutil
import shlex
import tempfile
import os
import re
from pathlib import Path

def run(cmd, cwd=None):
    """Run command and exit on error (secure version without shell=True)"""
    if isinstance(cmd, str):
        cmd = shlex.split(cmd)
    # nosec B603 nosemgrep: python.lang.security.audit.dangerous-subprocess-use-audit
    result = subprocess.run(
        cmd, shell=False, cwd=cwd, capture_output=True, text=True, check=False
    )
    if result.returncode != 0:
        print(f"❌ Command failed: {' '.join(cmd)}", file=sys.stderr)
        print(result.stderr, file=sys.stderr)
        sys.exit(1)
    return result.stdout

def main():
    print("🚀 N8N SDK Generation Pipeline\n")

    # Config
    openapi_source = "sdk/n8nsdk/api/openapi.yaml"
    openapi_spec = "sdk/n8nsdk/api/openapi-generated.yaml"
    generator_jar = os.path.join(tempfile.gettempdir(), "openapi-generator-cli.jar")
    generator_version = "7.11.0"
    sdk_dir = Path("sdk/n8nsdk")

    # Check if OpenAPI source exists
    if not Path(openapi_source).exists():
        print(f"❌ Error: {openapi_source} not found", file=sys.stderr)
        print("Run 'make openapi' first to download and prepare the OpenAPI spec", file=sys.stderr)
        sys.exit(1)

    # Copy source to generated version
    print("📋 Copying OpenAPI spec for generation...")
    shutil.copy(openapi_source, openapi_spec)
    print(f"   ✓ Copied {openapi_source} → {openapi_spec}\n")

    # Backup openapi.yaml to restore it after generation
    openapi_backup = Path(openapi_source + ".backup")
    shutil.copy(openapi_source, openapi_backup)

    # 1. Enable additionalProperties: true on ALL object schemas
    # This allows the SDK to accept any unknown fields from the API
    print("🔧 Enabling additionalProperties: true on all schemas...")
    spec_content = Path(openapi_spec).read_text(encoding='utf-8')
    # Replace additionalProperties: false with additionalProperties: true
    spec_content = spec_content.replace(
        'additionalProperties: false',
        'additionalProperties: true'
    )
    Path(openapi_spec).write_text(spec_content, encoding='utf-8')
    count_changes = spec_content.count('additionalProperties: true')
    print(f"   ✓ Enabled additionalProperties: true ({count_changes} schemas)\n")

    # 2. Generate SDK from modified spec
    print("🔨 Generating SDK...\n")

    # Check Java
    if not shutil.which("java"):
        print("   ❌ Error: Java required", file=sys.stderr)
        sys.exit(1)

    # Download openapi-generator JAR if needed
    if not Path(generator_jar).exists():
        print("   → Downloading openapi-generator JAR...")
        jar_url = (
            f"https://repo1.maven.org/maven2/org/openapitools/"
            f"openapi-generator-cli/{generator_version}/"
            f"openapi-generator-cli-{generator_version}.jar"
        )
        run(f"wget -q {jar_url} -O {generator_jar}")
        print("   ✓ Downloaded\n")

    # Clean previous generation (keep api directory)
    print("   → Cleaning previous SDK files...")
    for item in sdk_dir.iterdir():
        if item.name != "api":
            if item.is_dir():
                shutil.rmtree(item)
            else:
                item.unlink()
    print("   ✓ Cleaned\n")

    # Generate SDK
    print("   → Running openapi-generator...")
    # nosec B603 nosemgrep
    result = subprocess.run(
        [
            'java', '-jar', generator_jar,
            'generate',
            '-i', openapi_spec,
            '-g', 'go',
            '-o', 'sdk/n8nsdk',
            '-c', 'codegen/openapi-generator-config.yaml',
            '--skip-validate-spec'
        ],
        shell=False,
        capture_output=True,
        text=True,
        check=False
    )
    if result.returncode != 0:
        print("   ❌ Generator failed", file=sys.stderr)
        print(result.stderr, file=sys.stderr)
        sys.exit(1)
    print("   ✓ Generated\n")

    # 3. Add typed fields to Workflow struct for known API fields
    # These fields are returned by the n8n API but not in the official spec
    # Adding them as typed fields provides better IDE support and type safety
    print("   → Adding typed fields to Workflow struct...")
    workflow_model = sdk_dir / "model_workflow.go"
    if workflow_model.exists():
        content = workflow_model.read_text(encoding='utf-8')
        # Add typed fields before AdditionalProperties
        workflow_fields = '''	Description *string `json:"description,omitempty"`
	VersionId *string `json:"versionId,omitempty"`
	IsArchived *bool `json:"isArchived,omitempty"`
	TriggerCount *float32 `json:"triggerCount,omitempty"`
	Meta map[string]interface{} `json:"meta,omitempty"`
	PinData map[string]interface{} `json:"pinData,omitempty"`
	VersionCounter *int32 `json:"versionCounter,omitempty"`
	AdditionalProperties map[string]interface{}'''
        content = content.replace(
            '\tAdditionalProperties map[string]interface{}',
            workflow_fields
        )
        workflow_model.write_text(content, encoding='utf-8')
        print("   ✓ Added typed fields to Workflow\n")
    else:
        print("   ⚠️  model_workflow.go not found\n")

    # 4. Remove DisallowUnknownFields() from all model files (if present)
    # This allows the API to return unknown fields without causing unmarshaling errors
    print("   → Removing DisallowUnknownFields() from model files...")
    count = 0
    for model_file in sdk_dir.glob("model_*.go"):
        content = model_file.read_text(encoding='utf-8')
        if 'DisallowUnknownFields()' in content:
            content = content.replace(
                '\tdecoder.DisallowUnknownFields()\n',
                ''
            )
            model_file.write_text(content, encoding='utf-8')
            count += 1
    print(f"   ✓ Fixed {count} model files\n")

    # 5. Fix module paths in .go files (OpenAPI Generator uses placeholder paths)
    print("   → Fixing module paths in .go files...")
    for go_file in sdk_dir.rglob("*.go"):
        content = go_file.read_text(encoding='utf-8')
        content = content.replace(
            "github.com/GIT_USER_ID/GIT_REPO_ID/n8nsdk",
            "github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
        )
        go_file.write_text(content, encoding='utf-8')
    print("   ✓ Fixed\n")

    # 6. Fix go.mod module declaration
    print("   → Fixing go.mod module path...")
    go_mod = sdk_dir / "go.mod"
    if go_mod.exists():
        content = go_mod.read_text(encoding='utf-8')
        content = content.replace(
            "module github.com/GIT_USER_ID/GIT_REPO_ID/n8nsdk",
            "module github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
        )
        go_mod.write_text(content, encoding='utf-8')
        print("   ✓ Fixed\n")
    else:
        print("   ⚠️  go.mod not found\n")

    # 7. Run go mod tidy
    print("   → Running go mod tidy...")
    # nosec B603 B607 nosemgrep
    subprocess.run(
        ['go', 'mod', 'tidy'],
        shell=False,
        cwd="sdk/n8nsdk",
        capture_output=True,
        check=False
    )
    print("   ✓ Done\n")

    # 8. Restore original openapi.yaml (generator may have reformatted it)
    print("   → Restoring original openapi.yaml...")
    shutil.copy(openapi_backup, openapi_source)
    openapi_backup.unlink()
    print("   ✓ Restored\n")

    # 9. Generate Bazel files
    print("🏗️  Generating BUILD files...")
    run("bazel run //:gazelle")
    print("   ✓ Done\n")

    print("✅ SDK generation complete!\n")

if __name__ == '__main__':
    main()
