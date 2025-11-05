#!/bin/bash
# Generate n8n SDK from OpenAPI specification
set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RESET='\033[0m'

OPENAPI_SPEC="sdk/n8nsdk/api/openapi.yaml"
GENERATOR_JAR="/tmp/openapi-generator-cli.jar"
GENERATOR_VERSION="7.11.0"

# Check Java
if ! command -v java &> /dev/null; then
    echo -e "${RED}Error: Java required${RESET}"
    exit 1
fi

# Download openapi-generator JAR if needed
if [ ! -f "$GENERATOR_JAR" ]; then
    echo "  → Downloading openapi-generator JAR..."
    wget -q "https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/${GENERATOR_VERSION}/openapi-generator-cli-${GENERATOR_VERSION}.jar" -O "$GENERATOR_JAR"
fi

# Clean previous generation (keep api directory)
find sdk/n8nsdk -mindepth 1 -maxdepth 1 ! -name 'api' -exec rm -rf {} +

# Generate SDK
echo "  → Running openapi-generator..."
java -jar "$GENERATOR_JAR" generate \
  -i "$OPENAPI_SPEC" \
  -g go \
  -o sdk/n8nsdk \
  -c codegen/openapi-generator-config.yaml \
  --skip-validate-spec \
  > /dev/null 2>&1

# Fix module paths
echo "  → Fixing module paths..."
find sdk/n8nsdk -name "*.go" -type f -exec sed -i 's|github.com/GIT_USER_ID/GIT_REPO_ID/n8nsdk|github.com/kodflow/n8n/sdk/n8nsdk|g' {} +

# Run go mod tidy
cd sdk/n8nsdk
go mod tidy 2>&1 | grep -v "finding module" | grep -v "found github" || true
cd ../..
