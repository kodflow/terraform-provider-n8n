# ============================================================================
# Build & Installation
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

.PHONY: clean
clean: ## Clean Bazel build artifacts and reset openapi.yaml
	@echo ""
	@echo "$(BOLD)Cleaning build artifacts...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Removing Bazel cache\n"
	@bazel clean
	@printf "  $(CYAN)→$(RESET) Resetting openapi.yaml to committed version\n"
	@git checkout -- sdk/n8nsdk/api/openapi.yaml 2>/dev/null || true
	@echo "$(GREEN)✓$(RESET) Clean completed"
	@echo ""
