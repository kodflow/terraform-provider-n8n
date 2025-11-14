# ============================================================================
# GPG Key Management for Terraform Registry
# ============================================================================

# Default values from git config or devcontainer
GPG_NAME ?= $(shell git config --global user.name || echo "Kodflow")
GPG_EMAIL ?= $(shell git config --global user.email || echo "133899878+kodflow@users.noreply.github.com")
GPG_KEY_ID ?= $(shell gpg --list-secret-keys --with-colons "$(GPG_EMAIL)" 2>/dev/null | awk -F: '/sec:/ {print $$5; exit}')

.PHONY: gpg/generate
gpg/generate: ## Generate GPG key for signing releases
	@bash scripts/gpg-generate.sh "$(GPG_NAME)" "$(GPG_EMAIL)"

.PHONY: gpg/configure
gpg/configure: ## Configure git to use GPG signing
	@echo "$(CYAN)⚙️  Configuring git GPG signing...$(RESET)"
	@if [ -z "$(GPG_KEY_ID)" ]; then \
		echo "$(RED)❌ No GPG key found. Run 'make gpg/generate' first$(RESET)"; \
		exit 1; \
	fi
	@git config --global user.name "$(GPG_NAME)"
	@git config --global user.email "$(GPG_EMAIL)"
	@git config --global user.signingkey "$(GPG_KEY_ID)"
	@git config --global gpg.program gpg
	@git config --global commit.gpgsign true
	@git config --global tag.gpgsign true
	@echo "$(GREEN)✅ Git configured to sign commits and tags with key $(GPG_KEY_ID)$(RESET)"

.PHONY: gpg/export
gpg/export: ## Export GPG keys for GitHub Secrets
	@echo "$(CYAN)📤 Exporting GPG keys...$(RESET)"
	@if [ -z "$(GPG_KEY_ID)" ]; then \
		echo "$(RED)❌ No GPG key found. Run 'make gpg/generate' first$(RESET)"; \
		exit 1; \
	fi
	@echo ""
	@mkdir -p .gpg-export
	@gpg --armor --export "$(GPG_KEY_ID)" > .gpg-export/public-key.asc
	@gpg --armor --export-secret-keys "$(GPG_KEY_ID)" > .gpg-export/private-key.asc
	@echo "$(GREEN)✅ Keys exported to .gpg-export/$(RESET)"
	@echo ""
	@echo "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)"
	@echo "$(BOLD)  GitHub Secrets Configuration$(RESET)"
	@echo "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)"
	@echo ""
	@echo "$(BOLD)1. Go to:$(RESET) https://github.com/kodflow/n8n/settings/secrets/actions"
	@echo ""
	@echo "$(BOLD)2. Add these secrets:$(RESET)"
	@echo ""
	@echo "   $(CYAN)GPG_PRIVATE_KEY:$(RESET)"
	@cat .gpg-export/private-key.asc
	@echo ""
	@echo "   $(CYAN)GPG_PASSPHRASE:$(RESET) (leave empty if no passphrase)"
	@echo ""
	@echo "$(BOLD)3. For Terraform Registry, upload public key:$(RESET)"
	@echo "   File: .gpg-export/public-key.asc"
	@echo "   URL:  https://registry.terraform.io/settings/gpg-keys"
	@echo ""
	@echo "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)"
	@echo ""
	@echo "$(YELLOW)⚠️  IMPORTANT: The private key is SENSITIVE. Delete .gpg-export/ after adding to GitHub!$(RESET)"
	@echo "$(YELLOW)   Run: rm -rf .gpg-export/$(RESET)"
	@echo ""

.PHONY: gpg/export/public
gpg/export/public: ## Export only public key for Terraform Registry
	@if [ -z "$(GPG_KEY_ID)" ]; then \
		echo "$(RED)❌ No GPG key found. Run 'make gpg/generate' first$(RESET)"; \
		exit 1; \
	fi
	@mkdir -p .gpg-export
	@gpg --armor --export "$(GPG_KEY_ID)" > .gpg-export/public-key.asc
	@echo "$(GREEN)✅ Public key exported to .gpg-export/public-key.asc$(RESET)"
	@echo ""
	@echo "Upload this file to: $(CYAN)https://registry.terraform.io/settings/gpg-keys$(RESET)"

.PHONY: gpg/info
gpg/info: ## Show current GPG key information
	@echo "$(CYAN)🔍 GPG Configuration:$(RESET)"
	@echo ""
	@echo "$(BOLD)Name:$(RESET)  $(GPG_NAME)"
	@echo "$(BOLD)Email:$(RESET) $(GPG_EMAIL)"
	@if [ -n "$(GPG_KEY_ID)" ]; then \
		echo "$(BOLD)Key ID:$(RESET) $(GREEN)$(GPG_KEY_ID)$(RESET)"; \
		echo ""; \
		echo "$(BOLD)Key Details:$(RESET)"; \
		gpg --list-secret-keys "$(GPG_EMAIL)"; \
	else \
		echo "$(BOLD)Key ID:$(RESET) $(RED)Not found$(RESET)"; \
		echo ""; \
		echo "$(YELLOW)Run 'make gpg/generate' to create a new GPG key$(RESET)"; \
	fi
	@echo ""
	@echo "$(BOLD)Git GPG Signing:$(RESET)"
	@if [ "$$(git config --global commit.gpgsign)" = "true" ]; then \
		echo "  Commits: $(GREEN)✓ Enabled$(RESET)"; \
	else \
		echo "  Commits: $(RED)✗ Disabled$(RESET)"; \
	fi
	@if [ "$$(git config --global tag.gpgsign)" = "true" ]; then \
		echo "  Tags:    $(GREEN)✓ Enabled$(RESET)"; \
	else \
		echo "  Tags:    $(RED)✗ Disabled$(RESET)"; \
	fi

.PHONY: gpg/delete
gpg/delete: ## Delete GPG key (WARNING: irreversible!)
	@echo "$(RED)⚠️  WARNING: This will permanently delete your GPG key!$(RESET)"
	@echo ""
	@if [ -z "$(GPG_KEY_ID)" ]; then \
		echo "$(YELLOW)No GPG key found to delete$(RESET)"; \
		exit 0; \
	fi
	@echo "$(BOLD)Key to delete:$(RESET) $(GPG_KEY_ID)"
	@echo ""
	@read -p "Are you sure? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		gpg --batch --yes --delete-secret-keys "$(GPG_KEY_ID)" 2>/dev/null || true; \
		gpg --batch --yes --delete-keys "$(GPG_KEY_ID)" 2>/dev/null || true; \
		echo "$(GREEN)✅ GPG key deleted$(RESET)"; \
	else \
		echo "$(YELLOW)Cancelled$(RESET)"; \
	fi

.PHONY: gpg/test
gpg/test: ## Test GPG signing with a test file
	@echo "$(CYAN)🧪 Testing GPG signing...$(RESET)"
	@if [ -z "$(GPG_KEY_ID)" ]; then \
		echo "$(RED)❌ No GPG key found. Run 'make gpg/generate' first$(RESET)"; \
		exit 1; \
	fi
	@echo "Test message" > /tmp/gpg-test.txt
	@gpg --default-key "$(GPG_KEY_ID)" --armor --detach-sign /tmp/gpg-test.txt
	@if [ -f /tmp/gpg-test.txt.asc ]; then \
		echo "$(GREEN)✅ GPG signing works!$(RESET)"; \
		gpg --verify /tmp/gpg-test.txt.asc /tmp/gpg-test.txt; \
		rm -f /tmp/gpg-test.txt /tmp/gpg-test.txt.asc; \
	else \
		echo "$(RED)❌ GPG signing failed$(RESET)"; \
		exit 1; \
	fi

.PHONY: gpg/fingerprint
gpg/fingerprint: ## Show GPG key fingerprint for GoReleaser
	@if [ -z "$(GPG_KEY_ID)" ]; then \
		echo "$(RED)❌ No GPG key found. Run 'make gpg/generate' first$(RESET)"; \
		exit 1; \
	fi
	@echo "$(BOLD)GPG Fingerprint:$(RESET)"
	@gpg --fingerprint "$(GPG_EMAIL)" | grep -A 1 "Key fingerprint" || gpg --fingerprint "$(GPG_KEY_ID)"
	@echo ""
	@echo "$(BOLD)For GoReleaser, set this environment variable:$(RESET)"
	@echo "  $(CYAN)GPG_FINGERPRINT=$(RESET)$$(gpg --fingerprint "$(GPG_EMAIL)" | grep -A 1 "Key fingerprint" | tail -n 1 | tr -d ' ')"

.PHONY: gpg/setup
gpg/setup: gpg/generate gpg/configure gpg/export ## Complete GPG setup (generate + configure + export)
	@echo ""
	@echo "$(GREEN)$(BOLD)✅ GPG setup complete!$(RESET)"
	@echo ""
	@echo "$(BOLD)Next steps:$(RESET)"
	@echo "  1. Add GitHub Secrets (see instructions above)"
	@echo "  2. Upload public key to Terraform Registry"
	@echo "  3. Delete .gpg-export/ directory: $(CYAN)rm -rf .gpg-export/$(RESET)"
	@echo "  4. Test signing: $(CYAN)make gpg/test$(RESET)"
	@echo ""

.PHONY: gpg/clean
gpg/clean: ## Remove exported GPG keys directory
	@if [ -d .gpg-export ]; then \
		echo "$(YELLOW)🗑️  Removing .gpg-export/ directory...$(RESET)"; \
		rm -rf .gpg-export; \
		echo "$(GREEN)✅ Cleaned$(RESET)"; \
	else \
		echo "$(YELLOW)Nothing to clean$(RESET)"; \
	fi
