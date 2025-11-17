# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE in the project root for license information.

# Colors (inherit from main Makefile or define locally)
BLUE ?= \033[36m
GREEN ?= \033[32m
NC ?= \033[0m

.PHONY: nodes nodes/fetch nodes/parse nodes/diff nodes/generate nodes/workflows nodes/mega-workflow nodes/mega-workflow/test nodes/validate-coverage nodes/sync-report nodes/docs nodes/stats nodes/test nodes/clean

# Main nodes synchronization command
nodes: nodes/fetch nodes/parse nodes/sync-report nodes/diff nodes/generate nodes/test ## Synchronize n8n nodes from official repository
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

# Generate complete workflow for each node
nodes/workflows: ## Generate per-node workflow examples (296 examples)
	@echo "$(BLUE)[1mGenerating per-node workflow examples...$(NC)"
	@chmod +x scripts/nodes/generate-node-workflows.js
	@node scripts/nodes/generate-node-workflows.js
	@echo "$(GREEN)✓ Per-node workflows generated$(NC)"

# Generate MEGA workflow with ALL 296 nodes in a single workflow
nodes/mega-workflow: ## Generate single mega-workflow with ALL 296 nodes
	@echo "$(BLUE)[1mGenerating MEGA workflow (298 nodes in 1 workflow)...$(NC)"
	@chmod +x scripts/nodes/generate-mega-workflow.js
	@node scripts/nodes/generate-mega-workflow.js
	@echo "$(GREEN)✓ MEGA workflow generated in examples/mega-workflow/$(NC)"

# Test MEGA workflow (plan/apply/verify/destroy)
nodes/mega-workflow/test: nodes/mega-workflow ## Test MEGA workflow with real n8n instance
	@echo "$(BLUE)[1mTesting MEGA workflow...$(NC)"
	@echo "$(BLUE)ℹ  Using REAL n8n instance: $${N8N_API_URL}$(NC)"
	@echo ""
	@cd examples/mega-workflow && terraform init -upgrade >/dev/null 2>&1
	@cd examples/mega-workflow && terraform plan -var="n8n_base_url=$${N8N_API_URL}" -var="n8n_api_key=$${N8N_API_KEY}" -out=tfplan
	@echo ""
	@echo "$(BLUE)▶ Applying MEGA workflow (298 nodes)...$(NC)"
	@cd examples/mega-workflow && terraform apply -auto-approve tfplan
	@echo ""
	@echo "$(GREEN)✓ MEGA workflow created successfully$(NC)"
	@echo "$(BLUE)▶ Destroying MEGA workflow...$(NC)"
	@cd examples/mega-workflow && terraform destroy -var="n8n_base_url=$${N8N_API_URL}" -var="n8n_api_key=$${N8N_API_KEY}" -auto-approve
	@echo "$(GREEN)✓ MEGA workflow test completed$(NC)"

# Validate test coverage for all workflows
nodes/validate-coverage: ## Validate that all node tests cover ALL scenarios (inputs/outputs)
	@echo "$(BLUE)[1mValidating node workflow test coverage...$(NC)"
	@chmod +x scripts/nodes/validate-test-coverage.js
	@node scripts/nodes/validate-test-coverage.js

# Generate synchronization report for Claude agent
nodes/sync-report: ## Generate detailed sync report (NODES_SYNC.md)
	@echo "$(BLUE)[1mGenerating synchronization report...$(NC)"
	@chmod +x scripts/nodes/generate-sync-report.js
	@node scripts/nodes/generate-sync-report.js
	@echo "$(GREEN)✓ Sync report generated$(NC)"

# Generate comprehensive node documentation
nodes/docs: ## Generate SUPPORTED_NODES.md documentation
	@echo "$(BLUE)[1mGenerating node documentation...$(NC)"
	@chmod +x scripts/nodes/generate-nodes-documentation.js
	@node scripts/nodes/generate-nodes-documentation.js
	@echo "$(GREEN)✓ Node documentation generated$(NC)"

# Display node statistics
nodes/stats: ## Display node statistics
	@bash scripts/nodes/sync-n8n-nodes.sh stats

# NOTE: Use 'make test/unit' to run unit tests for workflow nodes/connections
# The paths //src/internal/provider/workflow/node/ and /connection/ no longer exist
# after flattening the structure. All tests are now in test/unit.

# Test all 296 workflow examples (VALIDATION ONLY - no real infrastructure)
nodes/test-workflows: ## Validate all 296 node examples (init/validate/plan with MOCK credentials)
	@echo "$(BLUE)[1mValidating all workflow examples (syntax only)...$(NC)"
	@echo "$(BLUE)ℹ  Using MOCK credentials for validation only$(NC)"
	@echo "$(BLUE)ℹ  For REAL infrastructure testing, use: make test/nodes$(NC)"
	@echo ""
	@chmod +x scripts/nodes/test-all-workflows.sh
	@bash scripts/nodes/test-all-workflows.sh
	@echo "$(GREEN)✓ Workflow validation completed$(NC)"

# Clean cache
nodes/clean: ## Clean nodes cache directory
	@echo "$(BLUE)[1mCleaning cache...$(NC)"
	@bash scripts/nodes/sync-n8n-nodes.sh clean
	@rm -f data/n8n-nodes-*.json data/n8n-nodes-*.md data/n8n-nodes-*.txt
	@echo "$(GREEN)✓ Cache cleaned$(NC)"
