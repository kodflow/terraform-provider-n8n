.DEFAULT_GOAL := help

# Détection automatique de l'OS et de l'architecture
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
VERSION := 0.0.1
PLUGIN_DIR := $(HOME)/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/$(VERSION)/$(GOOS)_$(GOARCH)

.PHONY: help
help: ## Affiche cette aide
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Lance les tests
	@bazel test --test_verbose_timeout_warnings //src/...

.PHONY: build
build: ## Build le provider localement et l'installe pour Terraform
	@echo "🔨 Compilation du provider avec Bazel..."
	@bazel build //src:terraform-provider-n8n
	@echo "📦 Installation dans $(PLUGIN_DIR)..."
	@mkdir -p $(PLUGIN_DIR)
	@cp -f bazel-bin/src/terraform-provider-n8n_/terraform-provider-n8n $(PLUGIN_DIR)/terraform-provider-n8n_v$(VERSION)
	@chmod +x $(PLUGIN_DIR)/terraform-provider-n8n_v$(VERSION)
	@echo "✅ Provider installé avec succès !"
	@echo "📍 Location: $(PLUGIN_DIR)/terraform-provider-n8n_v$(VERSION)"

.PHONY: clean
clean: ## Nettoie les artifacts Bazel
	@bazel clean
