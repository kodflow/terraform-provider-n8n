# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE in the project root for license information.

# Colors (inherit from main Makefile or define locally)
BLUE ?= \033[36m
GREEN ?= \033[32m
NC ?= \033[0m

.PHONY: nodes nodes/fetch nodes/parse nodes/diff nodes/generate nodes/stats nodes/test nodes/clean

# Main nodes synchronization command
nodes: nodes/fetch nodes/parse nodes/diff nodes/generate nodes/test ## Synchronize n8n nodes from official repository
	@echo ""
	@echo "$(GREEN)[1m✓ Nodes synchronization completed$(NC)"

# Fetch n8n repository
nodes/fetch: ## Fetch n8n official repository
	@echo "$(BLUE)[1mFetching n8n repository...$(NC)"
	@bash scripts/nodes/sync-n8n-nodes.sh fetch

# Parse nodes and generate registry
nodes/parse: ## Parse nodes and generate registry JSON
	@echo "$(BLUE)[1mParsing nodes...$(NC)"
	@bash scripts/nodes/sync-n8n-nodes.sh parse

# Generate changelog from differences
nodes/diff: ## Generate changelog from differences
	@echo "$(BLUE)[1mGenerating changelog...$(NC)"
	@bash scripts/nodes/sync-n8n-nodes.sh diff

# Generate Go code and Terraform examples
nodes/generate: ## Generate Go code and Terraform examples
	@echo "$(BLUE)[1mGenerating code and examples...$(NC)"
	@bash scripts/nodes/sync-n8n-nodes.sh generate

# Display node statistics
nodes/stats: ## Display node statistics
	@bash scripts/nodes/sync-n8n-nodes.sh stats

# Run node-related tests
nodes/test: ## Run node-related tests
	@echo "$(BLUE)[1mRunning node tests...$(NC)"
	@bazel test //src/internal/provider/workflow/node/...
	@bazel test //src/internal/provider/workflow/connection/...

# Clean cache
nodes/clean: ## Clean nodes cache directory
	@echo "$(BLUE)[1mCleaning cache...$(NC)"
	@bash scripts/nodes/sync-n8n-nodes.sh clean
	@rm -f data/n8n-nodes-*.json data/n8n-nodes-*.md data/n8n-nodes-*.txt
	@echo "$(GREEN)✓ Cache cleaned$(NC)"
