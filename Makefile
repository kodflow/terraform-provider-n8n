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
# Build & Test
# ============================================================================

.PHONY: test
test: ## Run test suite
	@echo ""
	@echo "$(BOLD)Running test suite...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Executing Bazel tests\n"
	@bazel test --test_verbose_timeout_warnings //src/...
	@echo "$(GREEN)✓$(RESET) Tests completed"
	@echo ""

.PHONY: test/n8n
test/n8n: build test/n8n/plan test/n8n/apply test/n8n/destroy ## Test provider with real n8n instance (plan → apply → destroy)
	@echo ""
	@echo "$(GREEN)✓$(RESET) Full test cycle completed successfully"
	@echo ""

.PHONY: test/n8n/plan
test/n8n/plan: build ## Run terraform plan with n8n provider
	@echo ""
	@echo "$(BOLD)Running terraform plan...$(RESET)"
	@if [ ! -f .env ]; then \
		printf "  $(RED)✗$(RESET) .env file not found. Please create it with N8N_URL and N8N_API_TOKEN\n"; \
		exit 1; \
	fi
	@printf "  $(CYAN)→$(RESET) Validating configuration\n"
	@export $$(cat .env | xargs) && cd examples/basic-sample && \
	terraform plan \
		-var="n8n_api_key=$$N8N_API_TOKEN" \
		-var="n8n_base_url=$$N8N_URL"
	@echo "$(GREEN)✓$(RESET) Plan completed successfully"
	@echo ""

.PHONY: test/n8n/apply
test/n8n/apply: build ## Run terraform apply with n8n provider
	@echo ""
	@echo "$(BOLD)Running terraform apply...$(RESET)"
	@if [ ! -f .env ]; then \
		printf "  $(RED)✗$(RESET) .env file not found. Please create it with N8N_URL and N8N_API_TOKEN\n"; \
		exit 1; \
	fi
	@printf "  $(CYAN)→$(RESET) Applying configuration\n"
	@export $$(cat .env | xargs) && cd examples/basic-sample && \
	terraform apply -auto-approve \
		-var="n8n_api_key=$$N8N_API_TOKEN" \
		-var="n8n_base_url=$$N8N_URL"
	@echo "$(GREEN)✓$(RESET) Apply completed successfully"
	@echo ""

.PHONY: test/n8n/destroy
test/n8n/destroy: ## Run terraform destroy with n8n provider
	@echo ""
	@echo "$(BOLD)Running terraform destroy...$(RESET)"
	@if [ ! -f .env ]; then \
		printf "  $(RED)✗$(RESET) .env file not found. Please create it with N8N_URL and N8N_API_TOKEN\n"; \
		exit 1; \
	fi
	@printf "  $(CYAN)→$(RESET) Destroying resources\n"
	@export $$(cat .env | xargs) && cd examples/basic-sample && \
	terraform destroy -auto-approve \
		-var="n8n_api_key=$$N8N_API_TOKEN" \
		-var="n8n_base_url=$$N8N_URL"
	@echo "$(GREEN)✓$(RESET) Destroy completed successfully"
	@echo ""

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
	@prettier --write "**/*.{json,yaml,yml,md}" --log-level silent
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
# API Tools
# ============================================================================

.PHONY: openapi
openapi: ## Download n8n OpenAPI from GitHub, patch, and generate SDK - Complete pipeline
	@echo ""
	@python3 codegen/build-sdk.py
	@echo ""
	@$(MAKE) fmt

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

# ============================================================================
# AI Context Generation
# ============================================================================

.PHONY: repomix
repomix: ## Generate compressed project context for AI (70% token reduction)
	@echo ""
	@echo "$(BOLD)$(CYAN)🤖 Generating AI context with Repomix...$(RESET)"
	@if ! command -v repomix >/dev/null 2>&1; then \
		printf "  $(RED)✗$(RESET) Repomix not installed. Run: npm install -g repomix\n"; \
		exit 1; \
	fi
	@printf "  $(CYAN)→$(RESET) Compressing codebase (excludes tests)\n"
	@repomix
	@echo "$(GREEN)✓$(RESET) Context generated: repomix-output.md"
	@echo ""
	@echo "$(BOLD)Token Statistics:$(RESET)"
	@grep -A 3 "Token Statistics" repomix-output.md || true
	@echo ""

.PHONY: repomix/full
repomix/full: ## Generate full project context including tests
	@echo ""
	@echo "$(BOLD)$(CYAN)🤖 Generating full AI context...$(RESET)"
	@if ! command -v repomix >/dev/null 2>&1; then \
		printf "  $(RED)✗$(RESET) Repomix not installed. Run: npm install -g repomix\n"; \
		exit 1; \
	fi
	@printf "  $(CYAN)→$(RESET) Compressing full codebase (with tests)\n"
	@repomix --include "**/*.go" --include "**/*.sh" --compress
	@echo "$(GREEN)✓$(RESET) Full context generated: repomix-output.md"
	@echo ""

.PHONY: repomix/install
repomix/install: ## Install Repomix globally
	@echo ""
	@echo "$(BOLD)$(CYAN)📦 Installing Repomix...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Installing via npm\n"
	@npm install -g repomix
	@echo "$(GREEN)✓$(RESET) Repomix installed successfully"
	@echo ""
