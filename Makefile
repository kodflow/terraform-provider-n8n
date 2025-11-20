.DEFAULT_GOAL := help

# ============================================================================
# Configuration
# ============================================================================

# Automatic OS and architecture detection
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Dynamic version detection from git tags
LAST_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v1.0.0")
VERSION := $(subst v,,$(LAST_TAG))

PLUGIN_DIR := $(HOME)/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/$(VERSION)/$(GOOS)_$(GOARCH)

# Colors for output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m
BOLD := \033[1m

# ============================================================================
# Include Modular Makefiles
# ============================================================================

include makefiles/sdk.mk
include makefiles/test.mk
include makefiles/nodes.mk
include makefiles/quality.mk
include makefiles/build.mk
include makefiles/tools.mk
include makefiles/terraform.mk

# ============================================================================
# Global Update Target
# ============================================================================

.PHONY: update
update: ## Update ALL dependencies (Go deps + n8n SDK + ktn-linter + README badge)
	@echo ""
	@echo "$(BOLD)$(CYAN)🔄 Updating all dependencies...$(RESET)"
	@echo ""
	@printf "$(BOLD)1/5 Updating Go dependencies to latest versions...$(RESET)\n"
	@./scripts/update-go-dependencies.sh
	@echo ""
	@printf "$(BOLD)2/5 Updating n8n commit to latest version...$(RESET)\n"
	@$(MAKE) sdk/openapi/update
	@echo ""
	@printf "$(BOLD)3/5 Updating ktn-linter to latest version...$(RESET)\n"
	@$(MAKE) tools/update
	@echo ""
	@printf "$(BOLD)4/5 Updating n8n version badge in README...$(RESET)\n"
	@./scripts/update-n8n-badge.sh
	@echo ""
	@printf "$(BOLD)5/5 Regenerating SDK and documentation...$(RESET)\n"
	@$(MAKE) sdk
	@$(MAKE) docs
	@echo ""
	@printf "$(BOLD)Formatting all source files...$(RESET)\n"
	@$(MAKE) fmt
	@echo ""
	@echo "$(BOLD)$(GREEN)✅ All updates completed$(RESET)"
	@echo ""
	@echo "$(YELLOW)ℹ$(RESET)  Next steps:"
	@echo "  1. Review changes: git diff"
	@echo "  2. Test changes:   make test"
	@echo "  3. Commit changes: git add . && git commit"
	@echo ""

# ============================================================================
# Help Target
# ============================================================================

.PHONY: help
help: ## Display available commands
	@echo ""
	@echo "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)"
	@echo "$(BOLD)  N8N Terraform Provider - Development Commands$(RESET)"
	@echo "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)"
	@echo ""
	@echo "$(BOLD)Main Commands:$(RESET)"
	@printf "  $(CYAN)%-28s$(RESET) %s\n" "update" "Update ALL dependencies (n8n SDK + ktn-linter + README badge)"
	@echo ""
	@echo "$(BOLD)Build & Clean:$(RESET)"
	@grep -h -E '^[a-zA-Z_/-]+:.*?## .*$$' makefiles/build.mk | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-28s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BOLD)Testing:$(RESET)"
	@grep -h -E '^[a-zA-Z_/-]+:.*?## .*$$' makefiles/test.mk | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-28s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BOLD)Terraform Operations:$(RESET)"
	@grep -h -E '^[a-zA-Z_/-]+:.*?## .*$$' makefiles/terraform.mk | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-28s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BOLD)Nodes Operations:$(RESET)"
	@grep -h -E '^[a-zA-Z_/-]+:.*?## .*$$' makefiles/nodes.mk | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-28s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BOLD)Code Quality:$(RESET)"
	@grep -h -E '^[a-zA-Z_/-]+:.*?## .*$$' makefiles/quality.mk | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-28s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BOLD)SDK Generation:$(RESET)"
	@grep -h -E '^[a-zA-Z_/-]+:.*?## .*$$' makefiles/sdk.mk | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-28s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BOLD)Tools & Dependencies:$(RESET)"
	@grep -h -E '^[a-zA-Z_/-]+:.*?## .*$$' makefiles/tools.mk | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-28s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)"
	@echo ""
