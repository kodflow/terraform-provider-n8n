# ============================================================================
# SDK Generation - OpenAPI & Code Generation
# ============================================================================

.PHONY: sdk
sdk: sdk/openapi ## Generate complete Go SDK (download + patch + generate)
	@echo ""
	@python3 codegen/generate-sdk.py
	@echo ""
	@$(MAKE) fmt

.PHONY: sdk/openapi
sdk/openapi: sdk/openapi/download sdk/openapi/patch ## Download & patch OpenAPI spec

.PHONY: sdk/openapi/download
sdk/openapi/download: ## Download OpenAPI spec from n8n (no commit, no patch)
	@echo ""
	@python3 codegen/download-only.py
	@echo ""

.PHONY: sdk/openapi/patch
sdk/openapi/patch: ## Apply patches to OpenAPI spec
	@echo ""
	@python3 codegen/patch-only.py
	@echo ""

.PHONY: sdk/openapi/patch/create
sdk/openapi/patch/create: ## Create patch from current git diff
	@python3 codegen/patch-openapi.py --create

.PHONY: sdk/openapi/update
sdk/openapi/update: ## Update N8N_COMMIT to latest version
	@python3 codegen/update-n8n-version.py
