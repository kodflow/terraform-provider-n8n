# ============================================================================
# Code Quality - Formatting, Linting & Documentation
# ============================================================================

.PHONY: quality
quality: fmt lint docs ## Run ALL quality checks (format + lint + docs) - Top level

.PHONY: fmt
fmt: ## Format all source files
	@echo ""
	@echo "$(BOLD)Formatting source files...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Go imports\n"
	@goimports -w $$(find . -type f -name "*.go" ! -path "./bazel-*" ! -path "./vendor/*") 2>/dev/null || true
	@printf "  $(CYAN)→$(RESET) Go files\n"
	@go fmt ./... > /dev/null
	@printf "  $(CYAN)→$(RESET) Bazel BUILD files (gazelle)\n"
	@bazel run //:gazelle 2>&1 | grep -E "^(ERROR|WARNING|INFO)" || true
	@printf "  $(CYAN)→$(RESET) Bazel files (buildifier)\n"
	@buildifier -r . 2>&1 | grep -v "^$$" || true
	@printf "  $(CYAN)→$(RESET) Shell scripts\n"
	@find . -name "*.sh" ! -path "./bazel-*" ! -name "p10k.sh" -exec shfmt -w -i 2 -ci -bn {} \; 2>/dev/null
	@printf "  $(CYAN)→$(RESET) YAML, JSON, Markdown\n"
	@prettier --write "**/*.{json,yaml,yml,md}" --ignore-path .prettierignore --log-level silent
	@printf "  $(CYAN)→$(RESET) Terraform files\n"
	@terraform fmt -recursive examples/ > /dev/null 2>&1 || true
	@echo "$(GREEN)✓$(RESET) Formatting completed"
	@echo ""

.PHONY: lint
lint: build ## Run code linters + validate Terraform examples
	@echo ""
	@echo "$(BOLD)Running code analysis...$(RESET)"
	@printf "  $(CYAN)→$(RESET) golangci-lint\n"
	@golangci-lint run ./...
	@printf "  $(CYAN)→$(RESET) ktn-linter\n"
	@ktn-linter lint ./... 2>&1 || true
	@printf "  $(CYAN)→$(RESET) Terraform examples validation\n"
	@./scripts/validate-examples.sh
	@echo "$(GREEN)✓$(RESET) Linting completed"
	@echo ""

.PHONY: docs
docs: ## Generate ALL documentation (Terraform docs + COVERAGE.MD + nodes README.md)
	@echo ""
	@echo "$(BOLD)$(CYAN)📝 Generating documentation...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Cleaning previous documentation\n"
	@rm -rf docs
	@printf "  $(CYAN)→$(RESET) Generating Terraform provider documentation\n"
	@go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir src --provider-name n8n --rendered-website-dir ../docs
	@printf "  $(CYAN)→$(RESET) Generating examples/nodes/README.md\n"
	@chmod +x scripts/nodes/generate-nodes-documentation.js
	@node scripts/nodes/generate-nodes-documentation.js
	@printf "  $(CYAN)→$(RESET) Generating COVERAGE.MD\n"
	@./scripts/generate-coverage.sh
	@printf "  $(CYAN)→$(RESET) Formatting generated documentation\n"
	@prettier --write "**/*.md" --ignore-path .prettierignore --log-level silent 2>/dev/null || true
	@echo "$(BOLD)$(GREEN)✅ Documentation generated$(RESET)"
	@echo ""
