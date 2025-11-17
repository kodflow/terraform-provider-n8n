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
# Help Target
# ============================================================================

.PHONY: help
help: ## Display available commands
	@echo ""
	@echo "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)"
	@echo "$(BOLD)  N8N Terraform Provider - Development Commands$(RESET)"
	@echo "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)"
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
