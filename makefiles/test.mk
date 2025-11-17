# ============================================================================
# Testing - Unit, Acceptance & Integration Tests
# ============================================================================

.PHONY: test
test: test/unit test/acceptance ## Run ALL tests (unit + acceptance) - Top level

.PHONY: test/unit
test/unit: ## Run all unit tests
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
test/acceptance: ## Run all E2E acceptance tests with real n8n instance
	@echo ""
	@echo "$(BOLD)Running E2E acceptance tests...$(RESET)"
	@if [ ! -f .env ]; then \
		printf "  $(YELLOW)⚠$(RESET)  .env file not found - skipping acceptance tests\n"; \
		printf "  $(CYAN)ℹ$(RESET)  Create .env with N8N_API_URL and N8N_API_KEY to run E2E tests\n"; \
		echo ""; \
		exit 0; \
	fi
	@printf "  $(CYAN)→$(RESET) Loading credentials from .env\n"
	@export $$(cat .env | xargs) && go test -v -tags=acceptance -timeout 30m ./src/internal/provider/credential/... ./src/internal/provider/tag/... ./src/internal/provider/variable/... ./src/internal/provider/workflow/... && echo "$(GREEN)✓$(RESET) E2E tests completed" || (printf "  $(YELLOW)⚠$(RESET)  E2E tests failed\n" && printf "  $(CYAN)ℹ$(RESET)  Verify N8N_API_URL is accessible and N8N_API_KEY is valid\n" && exit 1)
	@echo ""

.PHONY: test/acceptance/ci
test/acceptance/ci: ## Run E2E tests in CI mode (expects env vars)
	@echo ""
	@echo "$(BOLD)Running E2E acceptance tests...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Using credentials from environment variables\n"
	@go test -v -tags=acceptance -timeout 30m ./src/internal/provider/credential/... ./src/internal/provider/tag/... ./src/internal/provider/variable/... ./src/internal/provider/workflow/... && echo "$(GREEN)✓$(RESET) E2E tests completed" || (printf "  $(YELLOW)⚠$(RESET)  E2E tests failed\n" && printf "  $(CYAN)ℹ$(RESET)  Verify N8N_API_URL is accessible and N8N_API_KEY is valid\n" && exit 1)
	@echo ""

.PHONY: test/terraform
test/terraform: build ## Run ALL Terraform examples (plan/apply/destroy) with local provider
	@echo ""
	@echo "$(BOLD)Testing all Terraform examples...$(RESET)"
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
	@export $$(cat .env | xargs) && ./scripts/test-examples.sh
	@echo ""

.PHONY: test/tf
test/tf: test/terraform ## Alias for test/terraform (backward compatibility)

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
	TF_VAR_n8n_api_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform init -no-color -upgrade -plugin-dir=$(HOME)/.terraform.d/plugins > /dev/null && \
	printf "  $(CYAN)→$(RESET) Planning deployment\n" && \
	TF_VAR_n8n_api_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform plan -no-color -out=tfplan && \
	printf "  $(CYAN)→$(RESET) Applying changes\n" && \
	TF_VAR_n8n_api_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform apply -no-color -auto-approve tfplan && \
	echo "" && \
	echo "$(BOLD)$(GREEN)✅ Basic sample test completed successfully$(RESET)" && \
	echo "" && \
	echo "$(BOLD)Created Resources:$(RESET)" && \
	TF_VAR_n8n_api_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform output -no-color && \
	echo "" && \
	printf "  $(CYAN)→$(RESET) Destroying test resources\n" && \
	TF_VAR_n8n_api_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform destroy -no-color -auto-approve && \
	echo "$(GREEN)✓$(RESET) Test cleanup completed" && \
	echo "" || (echo "$(RED)✗$(RESET) Basic sample test failed" && exit 1)

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
	printf "  $(CYAN)→$(RESET) Destroying test resources\n" && \
	TF_VAR_n8n_api_url=$$N8N_API_URL TF_VAR_n8n_api_key=$$N8N_API_KEY terraform destroy -no-color -auto-approve && \
	echo "$(GREEN)✓$(RESET) Test cleanup completed" && \
	echo "" || (echo "$(RED)✗$(RESET) Community test failed" && exit 1)

.PHONY: test/validate-examples
test/validate-examples: build ## Validate all Terraform examples syntax with local provider (no .env needed)
	@echo ""
	@echo "$(BOLD)Validating Terraform examples syntax...$(RESET)"
	@./scripts/validate-examples.sh
	@echo ""
