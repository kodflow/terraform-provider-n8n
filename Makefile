.DEFAULT_GOAL := help

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

.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Run tests
	@bazel test --test_verbose_timeout_warnings //src/...

.PHONY: build
build: ## Build provider locally and install for Terraform
	@echo "🔨 Compiling provider with Bazel..."
	@bazel build //src:terraform-provider-n8n
	@echo "📦 Installing to $(PLUGIN_DIR)..."
	@mkdir -p $(PLUGIN_DIR)
	@cp -f bazel-bin/src/terraform-provider-n8n_/terraform-provider-n8n $(PLUGIN_DIR)/terraform-provider-n8n_v$(VERSION)
	@chmod +x $(PLUGIN_DIR)/terraform-provider-n8n_v$(VERSION)
	@echo "✅ Provider installed successfully!"
	@echo "📍 Location: $(PLUGIN_DIR)/terraform-provider-n8n_v$(VERSION)"

.PHONY: clean
clean: ## Clean Bazel artifacts
	@bazel clean

.PHONY: fmt
fmt: ## Format all files in the project
	@echo "🎨 Formatting all files..."
	@echo "  → Go files..."
	@go fmt ./...
	@echo "  → Bazel files..."
	@buildifier -r .
	@echo "  → Shell scripts..."
	@shfmt -w -i 2 -ci -bn .
	@echo "  → YAML, JSON, Markdown..."
	@prettier --write "**/*.{json,yaml,yml,md}"
	@echo "  → Terraform files..."
	@terraform fmt -recursive examples/ 2>/dev/null || true
	@echo "✅ All files formatted!"

.PHONY: lint
lint: ## Run golangci-lint with ktn-linter
	@echo "🔍 Running golangci-lint with ktn-linter..."
	@golangci-lint run ./...
