.DEFAULT_GOAL := help

# ============================================================================
# Configuration
# ============================================================================

# Automatic OS and architecture detection
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Automatic dev version based on last git tag
# After v0.1.0 release -> 0.1.1-dev
# After v0.2.0 release -> 0.2.1-dev
# No tag yet -> 0.0.1-dev
LAST_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
LAST_VERSION := $(patsubst v%,%,$(LAST_TAG))
VERSION_PARTS := $(subst ., ,$(LAST_VERSION))
MAJOR := $(word 1,$(VERSION_PARTS))
MINOR := $(word 2,$(VERSION_PARTS))
PATCH := $(word 3,$(VERSION_PARTS))
NEXT_PATCH := $(shell echo $$(($(PATCH) + 1)))
VERSION := $(MAJOR).$(MINOR).$(NEXT_PATCH)-dev

PLUGIN_DIR := $(HOME)/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/$(VERSION)/$(GOOS)_$(GOARCH)

# Colors for output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m
BOLD := \033[1m

# ============================================================================
# Targets
# ============================================================================

.PHONY: help
help: ## Display available commands
	@echo ""
	@echo "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)"
	@echo "$(BOLD)  N8N Terraform Provider - Development Commands$(RESET)"
	@echo "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_/-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-22s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)"
	@echo ""

# ============================================================================
# Dependencies
# ============================================================================

.PHONY: deps
deps: ## Download Go module dependencies
	@echo ""
	@echo "$(BOLD)Downloading Go dependencies...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Running go mod download\n"
	@go mod download
	@echo "$(GREEN)✓$(RESET) Dependencies downloaded"
	@echo ""

.PHONY: tools
tools: ## Install development tools
	@echo ""
	@echo "$(BOLD)Installing development tools...$(RESET)"
	@printf "  $(CYAN)→$(RESET) goimports\n"
	@go install golang.org/x/tools/cmd/goimports@latest
	@printf "  $(CYAN)→$(RESET) buildifier\n"
	@go install github.com/bazelbuild/buildtools/buildifier@latest
	@printf "  $(CYAN)→$(RESET) shfmt\n"
	@go install mvdan.cc/sh/v3/cmd/shfmt@latest
	@printf "  $(CYAN)→$(RESET) prettier (requires npm)\n"
	@npm install -g prettier 2>/dev/null || echo "  $(YELLOW)⚠$(RESET)  npm not found, skipping prettier"
	@echo "$(GREEN)✓$(RESET) Tools installed"
	@echo ""

.PHONY: tools/lint
tools/lint: ## Install linting tools
	@echo ""
	@echo "$(BOLD)Installing linting tools...$(RESET)"
	@printf "  $(CYAN)→$(RESET) golangci-lint\n"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	@printf "  $(CYAN)→$(RESET) ktn-linter\n"
	@$(MAKE) update
	@echo "$(GREEN)✓$(RESET) Linting tools installed"
	@echo ""

.PHONY: tools/sdk
tools/sdk: ## Install SDK generation dependencies
	@echo ""
	@echo "$(BOLD)Installing SDK generation tools...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Python dependencies (pyyaml, requests)\n"
	@pip install -q pyyaml requests 2>/dev/null || echo "  $(YELLOW)⚠$(RESET)  pip not found, skipping Python deps"
	@printf "  $(CYAN)→$(RESET) OpenAPI Generator CLI\n"
	@if [ ! -f /tmp/openapi-generator-cli.jar ]; then \
		wget -q https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/7.10.0/openapi-generator-cli-7.10.0.jar -O /tmp/openapi-generator-cli.jar; \
		echo '#!/bin/bash' | sudo tee /usr/local/bin/openapi-generator > /dev/null; \
		echo 'java -jar /tmp/openapi-generator-cli.jar "$$@"' | sudo tee -a /usr/local/bin/openapi-generator > /dev/null; \
		sudo chmod +x /usr/local/bin/openapi-generator; \
	fi
	@echo "$(GREEN)✓$(RESET) SDK tools installed"
	@echo ""

# ============================================================================
# Build
# ============================================================================

.PHONY: build
build: ## Build and install provider
	@echo ""
	@echo "$(BOLD)Building Terraform provider...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Compiling with Bazel\n"
	@bazel build //src:terraform-provider-n8n
	@printf "  $(CYAN)→$(RESET) Installing to plugin directory\n"
	@mkdir -p $(PLUGIN_DIR)
	@cp -f bazel-bin/src/terraform-provider-n8n_/terraform-provider-n8n $(PLUGIN_DIR)/terraform-provider-n8n_v$(VERSION)
	@chmod +x $(PLUGIN_DIR)/terraform-provider-n8n_v$(VERSION)
	@echo "$(GREEN)✓$(RESET) Provider installed successfully"
	@echo "  $(CYAN)Location:$(RESET) $(PLUGIN_DIR)/terraform-provider-n8n_v$(VERSION)"
	@echo ""

# ============================================================================
# Tests
# ============================================================================

.PHONY: test
test: test/unit build test/acceptance ## Run all tests (unit + E2E)
	@echo ""
	@echo "$(BOLD)$(GREEN)✅ All tests passed$(RESET)"
	@echo ""

.PHONY: test/unit
test/unit: ## Run unit tests
	@echo ""
	@echo "$(BOLD)Running unit tests...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Executing Bazel tests\n"
	@bazel test --test_verbose_timeout_warnings //src/...
	@echo "$(GREEN)✓$(RESET) Unit tests completed"
	@echo ""

.PHONY: test/unit/ci
test/unit/ci: ## Run unit tests with CI-friendly output
	@echo ""
	@echo "$(BOLD)Running unit tests...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Executing Bazel tests\n"
	@bazel test --test_output=all --test_verbose_timeout_warnings //src/...
	@echo "$(GREEN)✓$(RESET) Unit tests completed"
	@echo ""

.PHONY: test/acceptance
test/acceptance: ## Run E2E acceptance tests with real n8n instance
	@echo ""
	@echo "$(BOLD)Running E2E acceptance tests...$(RESET)"
	@if [ ! -f .env ]; then \
		printf "  $(YELLOW)⚠$(RESET)  .env file not found - skipping acceptance tests\n"; \
		printf "  $(CYAN)ℹ$(RESET)  Create .env with N8N_API_URL and N8N_API_KEY to run E2E tests\n"; \
		echo ""; \
		exit 0; \
	fi
	@printf "  $(CYAN)→$(RESET) Loading credentials from .env\n"
	@export $$(cat .env | xargs) && \
	TF_ACC=1 go test -v -tags=acceptance -timeout 30m \
		./src/internal/provider/credential/... \
		./src/internal/provider/tag/... \
		./src/internal/provider/variable/... \
		./src/internal/provider/workflow/...
	@echo "$(GREEN)✓$(RESET) E2E tests completed"
	@echo ""

.PHONY: test/tf/community
test/tf/community: build ## Test community resources with Terraform (uses .env)
	@echo ""
	@echo "$(BOLD)Running Community Edition integration test...$(RESET)"
	@if [ ! -f .env ]; then \
		printf "  $(RED)✗$(RESET) .env file not found\n"; \
		printf "  $(CYAN)ℹ$(RESET)  Create .env with N8N_API_URL and N8N_API_KEY\n"; \
		printf "  $(CYAN)ℹ$(RESET)  Example:\n"; \
		printf "      N8N_API_URL=http://localhost:5678\n"; \
		printf "      N8N_API_KEY=your-api-key-here\n"; \
		echo ""; \
		exit 1; \
	fi
	@printf "  $(CYAN)→$(RESET) Loading credentials from .env\n"
	@export $$(cat .env | xargs) && \
	cd examples/community && \
	rm -rf .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup 2>/dev/null || true && \
	printf "  $(CYAN)→$(RESET) Initializing Terraform\n" && \
	TF_VAR_n8n_api_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform init -no-color -upgrade -plugin-dir=$(HOME)/.terraform.d/plugins > /dev/null && \
	printf "  $(CYAN)→$(RESET) Planning deployment\n" && \
	TF_VAR_n8n_api_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform plan -no-color -out=tfplan && \
	printf "  $(CYAN)→$(RESET) Applying changes\n" && \
	TF_VAR_n8n_api_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform apply -no-color -auto-approve tfplan && \
	echo "" && \
	echo "$(BOLD)$(GREEN)✅ Community test completed successfully$(RESET)" && \
	echo "" && \
	echo "$(BOLD)Created Resources:$(RESET)" && \
	TF_VAR_n8n_api_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform output -no-color && \
	echo "" && \
	printf "  $(CYAN)→$(RESET) Destroying resources\n" && \
	TF_VAR_n8n_api_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform destroy -no-color -auto-approve && \
	rm -rf .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup tfplan 2>/dev/null || true && \
	echo "$(GREEN)✓$(RESET) Resources cleaned up" && \
	echo ""

.PHONY: test/tf/basic-sample
test/tf/basic-sample: build ## Test basic-sample example with Terraform (uses .env)
	@echo ""
	@echo "$(BOLD)Running Basic Sample integration test...$(RESET)"
	@if [ ! -f .env ]; then \
		printf "  $(RED)✗$(RESET) .env file not found\n"; \
		printf "  $(CYAN)ℹ$(RESET)  Create .env with N8N_API_URL and N8N_API_KEY\n"; \
		printf "  $(CYAN)ℹ$(RESET)  Example:\n"; \
		printf "      N8N_API_URL=http://localhost:5678\n"; \
		printf "      N8N_API_KEY=your-api-key-here\n"; \
		echo ""; \
		exit 1; \
	fi
	@printf "  $(CYAN)→$(RESET) Loading credentials from .env\n"
	@export $$(cat .env | xargs) && \
	cd examples/basic-sample && \
	rm -rf .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup 2>/dev/null || true && \
	printf "  $(CYAN)→$(RESET) Initializing Terraform\n" && \
	TF_VAR_n8n_base_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform init -no-color -upgrade -plugin-dir=$(HOME)/.terraform.d/plugins > /dev/null && \
	printf "  $(CYAN)→$(RESET) Planning deployment\n" && \
	TF_VAR_n8n_base_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform plan -no-color -out=tfplan && \
	printf "  $(CYAN)→$(RESET) Applying changes\n" && \
	TF_VAR_n8n_base_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform apply -no-color -auto-approve tfplan && \
	echo "" && \
	echo "$(BOLD)$(GREEN)✅ Basic Sample test completed successfully$(RESET)" && \
	echo "" && \
	echo "$(BOLD)Created Resources:$(RESET)" && \
	TF_VAR_n8n_base_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform output -no-color && \
	echo "" && \
	printf "  $(CYAN)→$(RESET) Destroying resources\n" && \
	TF_VAR_n8n_base_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform destroy -no-color -auto-approve && \
	rm -rf .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup tfplan 2>/dev/null || true && \
	echo "$(GREEN)✓$(RESET) Resources cleaned up" && \
	echo ""

# ============================================================================
# Code Quality
# ============================================================================

.PHONY: fmt
fmt: ## Format all source files
	@echo ""
	@echo "$(BOLD)Formatting source files...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Go imports\n"
	@goimports -w $$(find . -type f -name "*.go" ! -path "./bazel-*" ! -path "./vendor/*") 2>/dev/null || true
	@printf "  $(CYAN)→$(RESET) Go files\n"
	@go fmt ./... > /dev/null
	@printf "  $(CYAN)→$(RESET) Bazel BUILD files (gazelle)\n"
	@bazel run //:gazelle 2>&1 | grep -E "^(ERROR|WARNING|INFO)" || true
	@printf "  $(CYAN)→$(RESET) Bazel files (buildifier)\n"
	@buildifier -r . 2>&1 | grep -v "^$$" || true
	@printf "  $(CYAN)→$(RESET) Shell scripts\n"
	@find . -name "*.sh" ! -path "./bazel-*" ! -name "p10k.sh" -exec shfmt -w -i 2 -ci -bn {} \; 2>/dev/null
	@printf "  $(CYAN)→$(RESET) YAML, JSON, Markdown\n"
	@prettier --write "**/*.{json,yaml,yml,md}" --ignore-path .prettierignore --log-level silent
	@printf "  $(CYAN)→$(RESET) Terraform files\n"
	@terraform fmt -recursive examples/ > /dev/null 2>&1 || true
	@echo "$(GREEN)✓$(RESET) Formatting completed"
	@echo ""

.PHONY: lint
lint: ## Run code linters
	@echo ""
	@echo "$(BOLD)Running code analysis...$(RESET)"
	@printf "  $(CYAN)→$(RESET) golangci-lint\n"
	@golangci-lint run ./...
	@printf "  $(CYAN)→$(RESET) ktn-linter\n"
	@ktn-linter lint ./... 2>&1 || true
	@echo "$(GREEN)✓$(RESET) Linting completed"
	@echo ""

.PHONY: update
update: ## Update ktn-linter to latest version
	@echo ""
	@echo "$(BOLD)Updating ktn-linter...$(RESET)"
	@KTN_ARCH=$$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/'); \
	KTN_VERSION=$$(curl -s https://api.github.com/repos/kodflow/ktn-linter/releases/latest | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/'); \
	printf "  $(CYAN)→$(RESET) Downloading version v$$KTN_VERSION for $$KTN_ARCH\n"; \
	mkdir -p $$HOME/.local/bin; \
	curl -fsSL "https://github.com/kodflow/ktn-linter/releases/download/v$${KTN_VERSION}/ktn-linter-linux-$${KTN_ARCH}" -o "$$HOME/.local/bin/ktn-linter" && \
	chmod +x "$$HOME/.local/bin/ktn-linter" && \
	printf "$(GREEN)✓$(RESET) ktn-linter updated to v$$KTN_VERSION\n"; \
	echo ""

# ============================================================================
# SDK Generation
# ============================================================================

.PHONY: openapi
openapi: ## Download and prepare OpenAPI spec from n8n
	@echo ""
	@python3 codegen/download-openapi.py
	@echo ""

.PHONY: openapi/version
openapi/version: ## Check n8n OpenAPI version info
	@python3 codegen/download-openapi.py --version

.PHONY: openapi/update
openapi/update: ## Update to latest n8n version
	@python3 codegen/download-openapi.py --update

.PHONY: openapi/patch/create
openapi/patch/create: ## Create patch from current git diff
	@python3 codegen/patch-openapi.py --create

.PHONY: openapi/patch/from-commit
openapi/patch/from-commit: ## Create patch from specific commit (usage: make openapi/patch/from-commit COMMIT=hash)
	@python3 codegen/patch-openapi.py --from-commit $(COMMIT)

.PHONY: sdk
sdk: ## Generate Go SDK from OpenAPI spec
	@echo ""
	@python3 codegen/generate-sdk.py
	@echo ""

.PHONY: clean
clean: ## Clean Bazel build artifacts
	@echo ""
	@echo "$(BOLD)Cleaning build artifacts...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Removing Bazel cache\n"
	@bazel clean
	@echo "$(GREEN)✓$(RESET) Clean completed"
	@echo ""

# ============================================================================
# Documentation
# ============================================================================

.PHONY: docs
docs: ## Generate all documentation (changelog + coverage)
	@echo ""
	@echo "$(BOLD)$(CYAN)📝 Generating documentation...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Generating CHANGELOG.md\n"
	@./scripts/generate-changelog.sh
	@printf "  $(CYAN)→$(RESET) Generating COVERAGE.md\n"
	@./scripts/generate-coverage.sh
	@echo "$(BOLD)$(GREEN)✅ All documentation generated$(RESET)"
	@echo ""
