# ============================================================================
# Tools & Dependencies Installation
# ============================================================================

.PHONY: tools
tools: tools/dev tools/lint tools/sdk ## Install ALL tools (dev + lint + sdk) - Top level

.PHONY: tools/dev
tools/dev: ## Install development tools
	@echo ""
	@echo "$(BOLD)Installing development tools...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Go module dependencies\n"
	@go mod download
	@printf "  $(CYAN)→$(RESET) goimports\n"
	@go install golang.org/x/tools/cmd/goimports@latest
	@printf "  $(CYAN)→$(RESET) buildifier\n"
	@go install github.com/bazelbuild/buildtools/buildifier@latest
	@printf "  $(CYAN)→$(RESET) shfmt\n"
	@go install mvdan.cc/sh/v3/cmd/shfmt@latest
	@printf "  $(CYAN)→$(RESET) prettier (requires npm)\n"
	@npm install -g prettier 2>/dev/null || echo "  $(YELLOW)⚠$(RESET)  npm not found, skipping prettier"
	@echo "$(GREEN)✓$(RESET) Development tools installed"
	@echo ""

.PHONY: tools/lint
tools/lint: ## Install linting tools
	@echo ""
	@echo "$(BOLD)Installing linting tools...$(RESET)"
	@printf "  $(CYAN)→$(RESET) golangci-lint\n"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	@printf "  $(CYAN)→$(RESET) ktn-linter (latest version)\n"
	@KTN_ARCH=$$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/'); \
	if [ -n "$$GITHUB_TOKEN" ]; then \
		KTN_VERSION=$$(curl -s -H "Authorization: token $$GITHUB_TOKEN" https://api.github.com/repos/kodflow/ktn-linter/releases/latest | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/'); \
	else \
		KTN_VERSION=$$(curl -s https://api.github.com/repos/kodflow/ktn-linter/releases/latest | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/'); \
	fi; \
	if [ -z "$$KTN_VERSION" ]; then echo "  $(RED)✗$(RESET) Failed to get ktn-linter version from GitHub API"; exit 1; fi; \
	printf "  $(CYAN)  →$(RESET) Downloading version v$$KTN_VERSION for $$KTN_ARCH\n"; \
	mkdir -p $$HOME/.local/bin; \
	curl -fsSL "https://github.com/kodflow/ktn-linter/releases/download/v$${KTN_VERSION}/ktn-linter-linux-$${KTN_ARCH}" -o "$$HOME/.local/bin/ktn-linter" && \
	chmod +x "$$HOME/.local/bin/ktn-linter" && \
	printf "  $(GREEN)  ✓$(RESET) ktn-linter v$$KTN_VERSION installed\n"
	@echo "$(GREEN)✓$(RESET) Linting tools installed"
	@echo ""

.PHONY: tools/sdk
tools/sdk: ## Install SDK generation dependencies
	@echo ""
	@echo "$(BOLD)Installing SDK generation tools...$(RESET)"
	@printf "  $(CYAN)→$(RESET) Python dependencies (pyyaml, requests)\n"
	@pip install -q pyyaml requests 2>/dev/null || echo "  $(YELLOW)⚠$(RESET)  pip not found, skipping Python deps"
	@printf "  $(CYAN)→$(RESET) OpenAPI Generator CLI\n"
	@if [ ! -f /tmp/openapi-generator-cli.jar ]; then \
		wget -q https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/7.10.0/openapi-generator-cli-7.10.0.jar -O /tmp/openapi-generator-cli.jar; \
		echo '#!/bin/bash' | sudo tee /usr/local/bin/openapi-generator > /dev/null; \
		echo 'java -jar /tmp/openapi-generator-cli.jar "$$@"' | sudo tee -a /usr/local/bin/openapi-generator > /dev/null; \
		sudo chmod +x /usr/local/bin/openapi-generator; \
	fi
	@echo "$(GREEN)✓$(RESET) SDK tools installed"
	@echo ""

.PHONY: tools/update
tools/update: ## Update ktn-linter to latest version
	@echo ""
	@echo "$(BOLD)Updating ktn-linter...$(RESET)"
	@KTN_ARCH=$$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/'); \
	if [ -n "$$GITHUB_TOKEN" ]; then \
		KTN_VERSION=$$(curl -s -H "Authorization: token $$GITHUB_TOKEN" https://api.github.com/repos/kodflow/ktn-linter/releases/latest | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/'); \
	else \
		KTN_VERSION=$$(curl -s https://api.github.com/repos/kodflow/ktn-linter/releases/latest | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/'); \
	fi; \
	if [ -z "$$KTN_VERSION" ]; then echo "  $(RED)✗$(RESET) Failed to get ktn-linter version from GitHub API"; exit 1; fi; \
	printf "  $(CYAN)→$(RESET) Downloading version v$$KTN_VERSION for $$KTN_ARCH\n"; \
	mkdir -p $$HOME/.local/bin; \
	curl -fsSL "https://github.com/kodflow/ktn-linter/releases/download/v$${KTN_VERSION}/ktn-linter-linux-$${KTN_ARCH}" -o "$$HOME/.local/bin/ktn-linter" && \
	chmod +x "$$HOME/.local/bin/ktn-linter" && \
	printf "$(GREEN)✓$(RESET) ktn-linter updated to v$$KTN_VERSION\n"; \
	echo ""
