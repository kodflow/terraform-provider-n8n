#!/usr/bin/env bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE.md in the project root for license information.

# Test all Terraform examples with plan/apply/destroy
# This script tests both unitary examples (one resource each) and complex integration examples

set -e # Exit on first error

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKSPACE_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 Testing all Terraform examples"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check required environment variables
if [ -z "$N8N_API_URL" ]; then
  echo "❌ ERROR: N8N_API_URL environment variable is not set"
  echo "   Please set: export N8N_API_URL=https://your-n8n-instance.com"
  exit 1
fi

if [ -z "$N8N_API_KEY" ]; then
  echo "❌ ERROR: N8N_API_KEY environment variable is not set"
  echo "   Please set: export N8N_API_KEY=your-api-key"
  exit 1
fi

echo "✅ Environment variables configured"
echo ""

# Setup local provider override using Terraform CLI configuration
echo "📦 Setting up local provider override..."

# Find provider binary (installed by 'make build')
PROVIDER_BINARY=""
if [ -f "$WORKSPACE_DIR/bazel-bin/src/terraform-provider-n8n_/terraform-provider-n8n" ]; then
  PROVIDER_BINARY="$WORKSPACE_DIR/bazel-bin/src/terraform-provider-n8n_/terraform-provider-n8n"
elif [ -f "$HOME/go/bin/terraform-provider-n8n" ]; then
  PROVIDER_BINARY="$HOME/go/bin/terraform-provider-n8n"
else
  echo "❌ ERROR: Provider binary not found"
  echo "   Please build the provider first: make build"
  echo ""
  echo "   Searched locations:"
  echo "     - $WORKSPACE_DIR/bazel-bin/src/terraform-provider-n8n_/terraform-provider-n8n"
  echo "     - $HOME/go/bin/terraform-provider-n8n"
  exit 1
fi

echo "✅ Found provider binary at: $PROVIDER_BINARY"

# Get the directory containing the provider binary
PROVIDER_DIR=$(dirname "$PROVIDER_BINARY")

# Create Terraform CLI configuration file with dev_overrides
# This tells Terraform to use our local provider binary instead of downloading from registry
TERRAFORMRC_PATH="$HOME/.terraformrc"
echo "   Creating Terraform CLI configuration with dev_overrides..."

cat >"$TERRAFORMRC_PATH" <<EOF
provider_installation {
  dev_overrides {
    "kodflow/n8n" = "$PROVIDER_DIR"
  }

  # For all other providers, use the default registry
  direct {}
}
EOF

echo "✅ Terraform configured to use local provider"
echo "   Provider directory: $PROVIDER_DIR"
echo "   Config file: $TERRAFORMRC_PATH"
echo ""

# Export Terraform variables from environment
# This prevents Terraform from prompting for input in CI
export TF_VAR_n8n_api_key="$N8N_API_KEY"
export TF_VAR_n8n_base_url="$N8N_API_URL"

echo "✅ Terraform variables configured from environment"
echo ""

# Track overall status
FAILED_EXAMPLES=()
PASSED_EXAMPLES=()

# Function to test an example
test_example() {
  local example_dir="$1"
  local example_name
  local parent_dir

  example_name=$(basename "$example_dir")
  parent_dir=$(basename "$(dirname "$example_dir")")

  if [ "$parent_dir" != "examples" ]; then
    example_name="${parent_dir}/${example_name}"
  fi

  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "📝 Testing: $example_name"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo ""

  cd "$example_dir"

  # Initialize Terraform
  echo "🔧 Initializing Terraform..."
  if ! terraform init -no-color -input=false 2>&1 | tee /tmp/tf-init.log; then
    echo "❌ FAILED: terraform init for $example_name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    cat /tmp/tf-init.log
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    FAILED_EXAMPLES+=("$example_name (init)")
    return 1
  fi
  echo "✅ Init successful"
  echo ""

  # Validate configuration
  echo "🔍 Validating configuration..."
  if ! terraform validate -no-color 2>&1 | tee /tmp/tf-validate.log; then
    echo "❌ FAILED: terraform validate for $example_name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    cat /tmp/tf-validate.log
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    FAILED_EXAMPLES+=("$example_name (validate)")
    return 1
  fi
  echo "✅ Validation successful"
  echo ""

  # Plan
  echo "📋 Running terraform plan..."
  if ! terraform plan -no-color -input=false -out=tfplan 2>&1 | tee /tmp/tf-plan.log; then
    echo "❌ FAILED: terraform plan for $example_name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    cat /tmp/tf-plan.log
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    FAILED_EXAMPLES+=("$example_name (plan)")
    return 1
  fi
  echo "✅ Plan successful"
  echo ""

  # Apply
  echo "🚀 Running terraform apply..."
  if ! terraform apply -no-color -auto-approve tfplan 2>&1 | tee /tmp/tf-apply.log; then
    echo "❌ FAILED: terraform apply for $example_name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    cat /tmp/tf-apply.log
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    FAILED_EXAMPLES+=("$example_name (apply)")

    # Attempt cleanup even if apply failed
    echo "🧹 Attempting cleanup..."
    terraform destroy -no-color -input=false -auto-approve 2>&1 || true
    return 1
  fi
  echo "✅ Apply successful"
  echo ""

  # Show outputs
  echo "📤 Outputs:"
  terraform output -no-color 2>&1 || true
  echo ""

  # Destroy
  echo "🧹 Running terraform destroy..."
  if ! terraform destroy -no-color -input=false -auto-approve 2>&1 | tee /tmp/tf-destroy.log; then
    echo "❌ FAILED: terraform destroy for $example_name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    cat /tmp/tf-destroy.log
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    FAILED_EXAMPLES+=("$example_name (destroy)")
    return 1
  fi
  echo "✅ Destroy successful"
  echo ""

  PASSED_EXAMPLES+=("$example_name")
  echo "✅ SUCCESS: $example_name"
  echo ""

  return 0
}

# Test unitary examples (one resource type each)
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📦 PHASE 1: Testing unitary examples (one resource type each)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Find all example directories in community/
cd "$WORKSPACE_DIR"
UNITARY_EXAMPLES=$(find examples/community -mindepth 2 -maxdepth 2 -type d -name "*" ! -name "executions" 2>/dev/null | sort)

for example_dir in $UNITARY_EXAMPLES; do
  # Skip if no .tf files exist
  if ! ls "$example_dir"/*.tf >/dev/null 2>&1; then
    continue
  fi

  test_example "$WORKSPACE_DIR/$example_dir" || true # Continue even if this example fails

  # Return to workspace root
  cd "$WORKSPACE_DIR"
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏗️  PHASE 2: Testing complex integration example"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Test basic-sample (complex example with multiple resources)
test_example "$WORKSPACE_DIR/examples/basic-sample" || true

cd "$WORKSPACE_DIR"

# Print summary
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Test Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "✅ Passed: ${#PASSED_EXAMPLES[@]} examples"
for example in "${PASSED_EXAMPLES[@]}"; do
  echo "   ✓ $example"
done
echo ""

if [ ${#FAILED_EXAMPLES[@]} -gt 0 ]; then
  echo "❌ Failed: ${#FAILED_EXAMPLES[@]} examples"
  for example in "${FAILED_EXAMPLES[@]}"; do
    echo "   ✗ $example"
  done
  echo ""
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "❌ EXAMPLES TEST FAILED"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  exit 1
else
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "✅ ALL EXAMPLES PASSED"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
fi
