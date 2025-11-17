# ============================================================================
# Terraform Operations - Interactive Plan/Apply/Destroy
# ============================================================================

# Terraform context management
TF_CONTEXT_FILE := .terraform-context
TF_EXAMPLE ?= $(shell cat $(TF_CONTEXT_FILE) 2>/dev/null || echo "")

# Helper function to validate and save context
define save_tf_context
	@if [ -z "$(1)" ]; then \
		if [ -z "$(TF_EXAMPLE)" ]; then \
			printf "  $(RED)✗$(RESET) No Terraform example specified and no context saved\n"; \
			printf "  $(CYAN)ℹ$(RESET)  Usage: make $(2) examples/basic-sample\n"; \
			printf "  $(CYAN)ℹ$(RESET)  Available examples:\n"; \
			find examples -name "main.tf" -exec dirname {} \; | sort | sed 's/^/      - /'; \
			echo ""; \
			exit 1; \
		fi; \
	else \
		if [ ! -f "$(1)/main.tf" ]; then \
			printf "  $(RED)✗$(RESET) Invalid example path: $(1)\n"; \
			printf "  $(CYAN)ℹ$(RESET)  main.tf not found in $(1)\n"; \
			printf "  $(CYAN)ℹ$(RESET)  Available examples:\n"; \
			find examples -name "main.tf" -exec dirname {} \; | sort | sed 's/^/      - /'; \
			echo ""; \
			exit 1; \
		fi; \
		echo "$(1)" > $(TF_CONTEXT_FILE); \
	fi
endef

# Helper function to get current context
define get_tf_context
	$(if $(1),$(1),$(TF_EXAMPLE))
endef

# Helper function to load env vars
define load_tf_env
	@if [ ! -f .env ]; then \
		printf "  $(RED)✗$(RESET) .env file not found\n"; \
		printf "  $(CYAN)ℹ$(RESET)  Create .env with N8N_BASE_URL and N8N_API_KEY\n"; \
		printf "  $(CYAN)ℹ$(RESET)  Example:\n"; \
		printf "      N8N_BASE_URL=http://localhost:5678\n"; \
		printf "      N8N_API_KEY=your-api-key-here\n"; \
		echo ""; \
		exit 1; \
	fi
endef

# Helper function to export env vars
define export_tf_vars
	export $$(cat .env | xargs)
endef

.PHONY: tf/context
tf/context: ## Show current Terraform context
	@echo ""
	@echo "$(BOLD)Current Terraform Context:$(RESET)"
	@if [ -f "$(TF_CONTEXT_FILE)" ] && [ -s "$(TF_CONTEXT_FILE)" ]; then \
		printf "  $(GREEN)→$(RESET) $(BOLD)$$(cat $(TF_CONTEXT_FILE))$(RESET)\n"; \
	else \
		printf "  $(YELLOW)⚠$(RESET)  No context set\n"; \
		printf "  $(CYAN)ℹ$(RESET)  Set context with: make plan examples/basic-sample\n"; \
	fi
	@echo ""

.PHONY: tf/init
tf/init: build ## Initialize Terraform for current context [EXAMPLE]
	@$(call save_tf_context,$(EXAMPLE),tf/init)
	@$(call load_tf_env)
	@echo ""
	@echo "$(BOLD)Initializing Terraform...$(RESET)"
	@EXAMPLE_PATH=$(call get_tf_context,$(EXAMPLE)); \
	printf "  $(CYAN)→$(RESET) Context: $(BOLD)$$EXAMPLE_PATH$(RESET)\n"; \
	cd $$EXAMPLE_PATH && \
	printf "  $(CYAN)→$(RESET) Cleaning previous state\n" && \
	rm -rf .terraform .terraform.lock.hcl 2>/dev/null || true && \
	printf "  $(CYAN)→$(RESET) Running terraform init\n" && \
	$(call export_tf_vars) && \
	TF_VAR_n8n_base_url=$$N8N_BASE_URL TF_VAR_n8n_api_key=$$N8N_API_KEY \
	terraform init -upgrade -plugin-dir=$(HOME)/.terraform.d/plugins && \
	echo "$(GREEN)✓$(RESET) Initialization completed"
	@echo ""

.PHONY: tf/plan
tf/plan: build ## Plan Terraform changes for current context [EXAMPLE]
	@$(call save_tf_context,$(EXAMPLE),tf/plan)
	@$(call load_tf_env)
	@echo ""
	@echo "$(BOLD)Planning Terraform changes...$(RESET)"
	@EXAMPLE_PATH=$(call get_tf_context,$(EXAMPLE)); \
	printf "  $(CYAN)→$(RESET) Context: $(BOLD)$$EXAMPLE_PATH$(RESET)\n"; \
	cd $$EXAMPLE_PATH && \
	if [ ! -d ".terraform" ]; then \
		printf "  $(CYAN)→$(RESET) Initializing Terraform\n"; \
		$(call export_tf_vars) && \
		TF_VAR_n8n_base_url=$$N8N_BASE_URL TF_VAR_n8n_api_key=$$N8N_API_KEY \
		terraform init -upgrade -plugin-dir=$(HOME)/.terraform.d/plugins > /dev/null; \
	fi && \
	printf "  $(CYAN)→$(RESET) Running terraform plan\n" && \
	$(call export_tf_vars) && \
	TF_VAR_n8n_base_url=$$N8N_BASE_URL TF_VAR_n8n_api_key=$$N8N_API_KEY \
	terraform plan -out=tfplan && \
	echo "$(GREEN)✓$(RESET) Plan completed"
	@echo ""

.PHONY: tf/apply
tf/apply: ## Apply Terraform changes for current context [EXAMPLE]
	@$(call save_tf_context,$(EXAMPLE),tf/apply)
	@$(call load_tf_env)
	@echo ""
	@echo "$(BOLD)Applying Terraform changes...$(RESET)"
	@EXAMPLE_PATH=$(call get_tf_context,$(EXAMPLE)); \
	printf "  $(CYAN)→$(RESET) Context: $(BOLD)$$EXAMPLE_PATH$(RESET)\n"; \
	cd $$EXAMPLE_PATH && \
	if [ ! -f "tfplan" ]; then \
		printf "  $(YELLOW)⚠$(RESET)  No plan found, running plan first\n"; \
		if [ ! -d ".terraform" ]; then \
			printf "  $(CYAN)→$(RESET) Initializing Terraform\n"; \
			$(call export_tf_vars) && \
			TF_VAR_n8n_base_url=$$N8N_BASE_URL TF_VAR_n8n_api_key=$$N8N_API_KEY \
			terraform init -upgrade -plugin-dir=$(HOME)/.terraform.d/plugins > /dev/null; \
		fi && \
		printf "  $(CYAN)→$(RESET) Running terraform plan\n"; \
		$(call export_tf_vars) && \
		TF_VAR_n8n_base_url=$$N8N_BASE_URL TF_VAR_n8n_api_key=$$N8N_API_KEY \
		terraform plan -out=tfplan; \
	fi && \
	printf "  $(CYAN)→$(RESET) Running terraform apply\n" && \
	$(call export_tf_vars) && \
	TF_VAR_n8n_base_url=$$N8N_BASE_URL TF_VAR_n8n_api_key=$$N8N_API_KEY \
	terraform apply tfplan && \
	rm -f tfplan && \
	echo "" && \
	echo "$(BOLD)$(GREEN)✅ Apply completed successfully$(RESET)" && \
	echo "" && \
	echo "$(BOLD)Outputs:$(RESET)" && \
	$(call export_tf_vars) && \
	TF_VAR_n8n_base_url=$$N8N_BASE_URL TF_VAR_n8n_api_key=$$N8N_API_KEY \
	terraform output && \
	echo "$(GREEN)✓$(RESET) Apply completed"
	@echo ""

.PHONY: tf/destroy
tf/destroy: ## Destroy Terraform resources for current context [EXAMPLE]
	@$(call save_tf_context,$(EXAMPLE),tf/destroy)
	@$(call load_tf_env)
	@echo ""
	@echo "$(BOLD)$(YELLOW)Destroying Terraform resources...$(RESET)"
	@EXAMPLE_PATH=$(call get_tf_context,$(EXAMPLE)); \
	printf "  $(CYAN)→$(RESET) Context: $(BOLD)$$EXAMPLE_PATH$(RESET)\n"; \
	cd $$EXAMPLE_PATH && \
	if [ ! -d ".terraform" ]; then \
		printf "  $(CYAN)→$(RESET) Initializing Terraform\n"; \
		$(call export_tf_vars) && \
		TF_VAR_n8n_base_url=$$N8N_BASE_URL TF_VAR_n8n_api_key=$$N8N_API_KEY \
		terraform init -upgrade -plugin-dir=$(HOME)/.terraform.d/plugins > /dev/null; \
	fi && \
	printf "  $(CYAN)→$(RESET) Running terraform destroy\n" && \
	$(call export_tf_vars) && \
	TF_VAR_n8n_base_url=$$N8N_BASE_URL TF_VAR_n8n_api_key=$$N8N_API_KEY \
	terraform destroy -auto-approve && \
	echo "$(GREEN)✓$(RESET) Destroy completed"
	@echo ""

.PHONY: tf/output
tf/output: ## Show Terraform outputs for current context [EXAMPLE]
	@$(call save_tf_context,$(EXAMPLE),tf/output)
	@$(call load_tf_env)
	@echo ""
	@echo "$(BOLD)Terraform Outputs:$(RESET)"
	@EXAMPLE_PATH=$(call get_tf_context,$(EXAMPLE)); \
	printf "  $(CYAN)→$(RESET) Context: $(BOLD)$$EXAMPLE_PATH$(RESET)\n"; \
	echo ""; \
	cd $$EXAMPLE_PATH && \
	if [ ! -d ".terraform" ]; then \
		printf "  $(RED)✗$(RESET) Terraform not initialized\n"; \
		printf "  $(CYAN)ℹ$(RESET)  Run: make tf/init\n"; \
		echo ""; \
		exit 1; \
	fi && \
	$(call export_tf_vars) && \
	TF_VAR_n8n_base_url=$$N8N_BASE_URL TF_VAR_n8n_api_key=$$N8N_API_KEY \
	terraform output
	@echo ""

.PHONY: tf/clean
tf/clean: ## Clean all Terraform state and cache files [EXAMPLE]
	@$(call save_tf_context,$(EXAMPLE),tf/clean)
	@echo ""
	@echo "$(BOLD)Cleaning Terraform state...$(RESET)"
	@EXAMPLE_PATH=$(call get_tf_context,$(EXAMPLE)); \
	printf "  $(CYAN)→$(RESET) Context: $(BOLD)$$EXAMPLE_PATH$(RESET)\n"; \
	cd $$EXAMPLE_PATH && \
	printf "  $(CYAN)→$(RESET) Removing .terraform directory\n" && \
	rm -rf .terraform && \
	printf "  $(CYAN)→$(RESET) Removing lock file\n" && \
	rm -f .terraform.lock.hcl && \
	printf "  $(CYAN)→$(RESET) Removing state files\n" && \
	rm -f terraform.tfstate terraform.tfstate.backup && \
	printf "  $(CYAN)→$(RESET) Removing plan files\n" && \
	rm -f tfplan && \
	echo "$(GREEN)✓$(RESET) Clean completed"
	@echo ""

.PHONY: tf/list
tf/list: ## List all available Terraform examples
	@echo ""
	@echo "$(BOLD)Available Terraform Examples:$(RESET)"
	@find examples -name "main.tf" -exec dirname {} \; | sort | while read example; do \
		if [ -f "$(TF_CONTEXT_FILE)" ] && [ "$$(cat $(TF_CONTEXT_FILE))" = "$$example" ]; then \
			printf "  $(GREEN)●$(RESET) $(BOLD)$$example$(RESET) $(GREEN)(current)$(RESET)\n"; \
		else \
			printf "  $(CYAN)○$(RESET) $$example\n"; \
		fi; \
	done
	@echo ""
	@printf "  $(CYAN)ℹ$(RESET)  Use: make tf/plan examples/basic-sample\n"
	@echo ""

# ============================================================================
# Convenience Aliases (same as tf/plan, tf/apply, tf/destroy)
# ============================================================================
# These are shortcuts for the most common Terraform operations.
# They work exactly the same as their tf/* counterparts.
#
# Examples:
#   make plan examples/basic-sample   →  same as: make tf/plan examples/basic-sample
#   make apply                        →  same as: make tf/apply (uses saved context)
#   make destroy                      →  same as: make tf/destroy (uses saved context)
# ============================================================================

.PHONY: plan
plan: tf/plan ## Terraform plan (alias for tf/plan) [EXAMPLE]

.PHONY: apply
apply: tf/apply ## Terraform apply (alias for tf/apply) [EXAMPLE]

.PHONY: destroy
destroy: tf/destroy ## Terraform destroy (alias for tf/destroy) [EXAMPLE]
