# ============================================================================
# Local n8n Instance - Start, stop, and configure for acceptance tests
# ============================================================================

COMPOSE_FILE := .devcontainer/docker-compose.yml

.PHONY: n8n/start
n8n/start: ## Start the local n8n service (http://localhost:5678)
	@echo ""
	@echo "$(BOLD)Starting local n8n instance...$(RESET)"
	@docker compose -f $(COMPOSE_FILE) up -d n8n
	@printf "  $(CYAN)→$(RESET) n8n starting at http://localhost:5678\n"
	@printf "  $(CYAN)ℹ$(RESET)  Run $(BOLD)make n8n/setup$(RESET) once to create credentials\n"
	@echo ""

.PHONY: n8n/stop
n8n/stop: ## Stop the local n8n service
	@echo ""
	@echo "$(BOLD)Stopping local n8n instance...$(RESET)"
	@docker compose -f $(COMPOSE_FILE) stop n8n
	@echo "$(GREEN)✓$(RESET) n8n stopped"
	@echo ""

.PHONY: n8n/logs
n8n/logs: ## Tail logs from the local n8n service
	@docker compose -f $(COMPOSE_FILE) logs -f n8n

.PHONY: n8n/setup
n8n/setup: ## Configure local n8n: create owner + API key, write .env
	@echo ""
	@printf "  $(CYAN)→$(RESET) Running n8n setup script\n"
	@chmod +x scripts/setup-local-n8n.sh
	@bash scripts/setup-local-n8n.sh

.PHONY: n8n/reset
n8n/reset: ## Delete n8n data volume and restart fresh (destroys all data)
	@echo ""
	@printf "  $(YELLOW)⚠$(RESET)  This will destroy all n8n data and workflows\n"
	@docker compose -f $(COMPOSE_FILE) stop n8n
	@docker volume rm $$(docker compose -f $(COMPOSE_FILE) config --volumes | grep n8n-data || true) 2>/dev/null || true
	@docker compose -f $(COMPOSE_FILE) up -d n8n
	@printf "  $(CYAN)ℹ$(RESET)  Run $(BOLD)make n8n/setup$(RESET) to reconfigure\n"
	@echo ""
