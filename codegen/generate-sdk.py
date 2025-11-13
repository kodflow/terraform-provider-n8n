#!/usr/bin/env python3
"""
N8N SDK Generation Pipeline
Generates Go SDK from OpenAPI spec using openapi-generator
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
    print("🚀 N8N SDK Generation Pipeline\n")

    # Config
    OPENAPI_SOURCE = "sdk/n8nsdk/api/openapi.yaml"
    OPENAPI_SPEC = "sdk/n8nsdk/api/openapi-generated.yaml"
    GENERATOR_JAR = "/tmp/openapi-generator-cli.jar"
    GENERATOR_VERSION = "7.11.0"
    sdk_dir = Path("sdk/n8nsdk")

    # Check if OpenAPI source exists
    if not Path(OPENAPI_SOURCE).exists():
        print(f"❌ Error: {OPENAPI_SOURCE} not found", file=sys.stderr)
        print("Run 'make openapi' first to download and prepare the OpenAPI spec", file=sys.stderr)
        sys.exit(1)

    # Copy source to generated version
    print("📋 Copying OpenAPI spec for generation...")
    shutil.copy(OPENAPI_SOURCE, OPENAPI_SPEC)
    print(f"   ✓ Copied {OPENAPI_SOURCE} → {OPENAPI_SPEC}\n")

    # 1. Generate SDK
    print("🔨 Generating SDK...\n")

    # Check Java
    if not shutil.which("java"):
        print("   ❌ Error: Java required", file=sys.stderr)
        sys.exit(1)

    # Download openapi-generator JAR if needed
    if not Path(GENERATOR_JAR).exists():
        print("   → Downloading openapi-generator JAR...")
        run(f"wget -q https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/{GENERATOR_VERSION}/openapi-generator-cli-{GENERATOR_VERSION}.jar -O {GENERATOR_JAR}")
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
    result = subprocess.run(
        f"java -jar {GENERATOR_JAR} generate -i {OPENAPI_SPEC} -g go -o sdk/n8nsdk -c codegen/openapi-generator-config.yaml --skip-validate-spec",
        shell=True,
        capture_output=True,
        text=True
    )
    if result.returncode != 0:
        print(f"   ❌ Generator failed", file=sys.stderr)
        print(result.stderr, file=sys.stderr)
        sys.exit(1)
    print("   ✓ Generated\n")

    # 2. Fix module paths
    print("   → Fixing module paths...")
    for go_file in sdk_dir.rglob("*.go"):
        content = go_file.read_text()
        content = content.replace(
            "github.com/GIT_USER_ID/GIT_REPO_ID/n8nsdk",
            "github.com/kodflow/n8n/sdk/n8nsdk"
        )
        go_file.write_text(content)
    print("   ✓ Fixed\n")

    # 3. Run go mod tidy
    print("   → Running go mod tidy...")
    subprocess.run("go mod tidy", shell=True, cwd="sdk/n8nsdk", capture_output=True)
    print("   ✓ Done\n")

    # 4. Generate Bazel files
    print("🏗️  Generating BUILD files...")
    run("bazel run //:gazelle")
    print("   ✓ Done\n")

    print("✅ SDK generation complete!\n")

if __name__ == '__main__':
    main()
