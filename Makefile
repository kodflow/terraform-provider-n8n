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
	@echo "$(BOLD)$(CYAN)Available commands:$(RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""

# ============================================================================
# Build & Test
# ============================================================================

.PHONY: test
test: ## Run test suite
	@printf "$(BOLD)Running test suite...$(RESET)\n"
	@bazel test --test_verbose_timeout_warnings //src/...
	@printf "$(GREEN)✓$(RESET) Tests completed\n"

.PHONY: build
build: ## Build and install provider
	@printf "$(BOLD)Building Terraform provider...$(RESET)\n"
	@printf "  $(CYAN)→$(RESET) Compiling with Bazel\n"
	@bazel build //src:terraform-provider-n8n
	@printf "  $(CYAN)→$(RESET) Installing to plugin directory\n"
	@mkdir -p $(PLUGIN_DIR)
	@cp -f bazel-bin/src/terraform-provider-n8n_/terraform-provider-n8n $(PLUGIN_DIR)/terraform-provider-n8n_v$(VERSION)
	@chmod +x $(PLUGIN_DIR)/terraform-provider-n8n_v$(VERSION)
	@printf "$(GREEN)✓$(RESET) Provider installed successfully\n"
	@printf "  $(CYAN)Location:$(RESET) $(PLUGIN_DIR)/terraform-provider-n8n_v$(VERSION)\n"

.PHONY: clean
clean: ## Remove build artifacts
	@printf "$(BOLD)Cleaning build artifacts...$(RESET)\n"
	@bazel clean
	@printf "$(GREEN)✓$(RESET) Cleanup completed\n"

# ============================================================================
# Code Quality
# ============================================================================

.PHONY: fmt
fmt: ## Format all source files
	@printf "$(BOLD)Formatting source files...$(RESET)\n"
	@printf "  $(CYAN)→$(RESET) Go files\n"
	@go fmt ./... > /dev/null
	@printf "  $(CYAN)→$(RESET) Bazel files\n"
	@buildifier -r . 2>&1 | grep -v "^$$" || true
	@printf "  $(CYAN)→$(RESET) Shell scripts\n"
	@find . -name "*.sh" ! -path "./bazel-*" ! -name "p10k.sh" -exec shfmt -w -i 2 -ci -bn {} \; 2>/dev/null
	@printf "  $(CYAN)→$(RESET) YAML, JSON, Markdown\n"
	@prettier --write "**/*.{json,yaml,yml,md}" --log-level silent
	@printf "  $(CYAN)→$(RESET) Terraform files\n"
	@terraform fmt -recursive examples/ > /dev/null 2>&1 || true
	@printf "$(GREEN)✓$(RESET) Formatting completed\n"

.PHONY: lint
lint: ## Run code linters
	@printf "$(BOLD)Running code analysis...$(RESET)\n"
	@printf "  $(CYAN)→$(RESET) golangci-lint\n"
	@golangci-lint run ./...
	@printf "  $(CYAN)→$(RESET) ktn-linter\n"
	@ktn-linter lint --simple ./... || true
	@printf "$(GREEN)✓$(RESET) Linting completed\n"

.PHONY: update
update: ## Update ktn-linter to latest version
	@printf "$(BOLD)Updating ktn-linter...$(RESET)\n"
	@KTN_ARCH=$$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/'); \
	KTN_VERSION=$$(curl -s https://api.github.com/repos/kodflow/ktn-linter/releases/latest | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/'); \
	printf "  $(CYAN)→$(RESET) Downloading version v$$KTN_VERSION for $$KTN_ARCH\n"; \
	mkdir -p $$HOME/.local/bin; \
	curl -fsSL "https://github.com/kodflow/ktn-linter/releases/download/v$${KTN_VERSION}/ktn-linter-linux-$${KTN_ARCH}" -o "$$HOME/.local/bin/ktn-linter" && \
	chmod +x "$$HOME/.local/bin/ktn-linter" && \
	printf "$(GREEN)✓$(RESET) ktn-linter updated to v$$KTN_VERSION\n"
