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

    # 1. Generate SDK
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

    # 2. Fix model_workflow.go (add missing fields that OpenAPI Generator ignores)
    print("   → Fixing model_workflow.go...")
    workflow_model = sdk_dir / "model_workflow.go"
    if workflow_model.exists():
        content = workflow_model.read_text(encoding='utf-8')

        # Check if fields are already there
        if 'VersionId' not in content:
            # Add fields to struct
            content = content.replace(
                '\tShared []SharedWorkflow `json:"shared,omitempty"`\n}',
                '\tShared []SharedWorkflow `json:"shared,omitempty"`\n' +
                '\tVersionId *string `json:"versionId,omitempty"`\n' +
                '\tIsArchived *bool `json:"isArchived,omitempty"`\n' +
                '\tTriggerCount *float32 `json:"triggerCount,omitempty"`\n' +
                '\tMeta map[string]interface{} `json:"meta,omitempty"`\n' +
                '\tPinData map[string]interface{} `json:"pinData,omitempty"`\n' +
                '}'
            )

            # Add ToMap serialization
            content = content.replace(
                '\tif !IsNil(o.Shared) {\n\t\ttoSerialize["shared"] = o.Shared\n\t}\n\treturn toSerialize, nil\n}',
                '\tif !IsNil(o.Shared) {\n\t\ttoSerialize["shared"] = o.Shared\n\t}\n' +
                '\tif !IsNil(o.VersionId) {\n\t\ttoSerialize["versionId"] = o.VersionId\n\t}\n' +
                '\tif !IsNil(o.IsArchived) {\n\t\ttoSerialize["isArchived"] = o.IsArchived\n\t}\n' +
                '\tif !IsNil(o.TriggerCount) {\n\t\ttoSerialize["triggerCount"] = o.TriggerCount\n\t}\n' +
                '\tif !IsNil(o.Meta) {\n\t\ttoSerialize["meta"] = o.Meta\n\t}\n' +
                '\tif !IsNil(o.PinData) {\n\t\ttoSerialize["pinData"] = o.PinData\n\t}\n' +
                '\treturn toSerialize, nil\n}'
            )

            # Add getter/setter methods (simplified - just the struct is enough for now)
            # Full methods can be added later if needed

            workflow_model.write_text(content, encoding='utf-8')
            print("   ✓ Added missing workflow fields\n")
        else:
            print("   ✓ Fields already present\n")

        # Remove DisallowUnknownFields() to allow API to return extra fields
        content = workflow_model.read_text(encoding='utf-8')
        if 'decoder.DisallowUnknownFields()' in content:
            lines = content.split('\n')
            filtered_lines = [line for line in lines if 'DisallowUnknownFields()' not in line]
            content = '\n'.join(filtered_lines)
            workflow_model.write_text(content, encoding='utf-8')
            print("   ✓ Removed DisallowUnknownFields() to allow extra API fields\n")
    else:
        print("   ⚠️  model_workflow.go not found\n")

    # 3. Fix module paths in .go files
    print("   → Fixing module paths in .go files...")
    for go_file in sdk_dir.rglob("*.go"):
        content = go_file.read_text(encoding='utf-8')
        content = content.replace(
            "github.com/GIT_USER_ID/GIT_REPO_ID/n8nsdk",
            "github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
        )
        go_file.write_text(content, encoding='utf-8')
    print("   ✓ Fixed\n")

    # 4. Fix go.mod module declaration
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

    # 5. Run go mod tidy
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

    # 6. Restore original openapi.yaml (generator may have reformatted it)
    print("   → Restoring original openapi.yaml...")
    shutil.copy(openapi_backup, openapi_source)
    openapi_backup.unlink()
    print("   ✓ Restored\n")

    # 7. Generate Bazel files
    print("🏗️  Generating BUILD files...")
    run("bazel run //:gazelle")
    print("   ✓ Done\n")

    print("✅ SDK generation complete!\n")

if __name__ == '__main__':
    main()
