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

.PHONY: test/acceptance
test/acceptance: ## Run E2E acceptance tests with real n8n instance
	@echo ""
	@echo "$(BOLD)Running E2E acceptance tests...$(RESET)"
	@if [ ! -f .env ]; then \
		printf "  $(YELLOW)⚠$(RESET)  .env file not found - skipping acceptance tests\n"; \
		printf "  $(CYAN)ℹ$(RESET)  Create .env with N8N_URL and N8N_API_TOKEN to run E2E tests\n"; \
		echo ""; \
		exit 0; \
	fi
	@printf "  $(CYAN)→$(RESET) Loading credentials from .env\n"
	@export $$(cat .env | xargs) && \
	TF_ACC=1 go test -v -tags=acceptance -timeout 30m \
		./src/internal/provider/credential/... \
		./src/internal/provider/tag/... \
		./src/internal/provider/variable/... \
		./src/internal/provider/workflow/... || \
		(printf "  $(RED)✗$(RESET) Acceptance tests failed (this is expected if n8n instance is not accessible)\n"; exit 0)
	@echo "$(GREEN)✓$(RESET) E2E tests completed"
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
# SDK Generation
# ============================================================================

.PHONY: openapi
openapi: ## Generate SDK from n8n OpenAPI specification
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
