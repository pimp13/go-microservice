# Variables
DOCKER_COMPOSE = docker-compose
PROJECT_NAME = echo-microservices

# Colors for output
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[1;33m
NC = \033[0m # No Color

.PHONY: help build up down restart logs clean test

help: ## نمایش راهنما
	@echo "$(GREEN)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'

build: ## Build تمام سرویس‌ها
	@echo "$(GREEN)Building all services...$(NC)"
	$(DOCKER_COMPOSE) build --no-cache

up: ## اجرای تمام سرویس‌ها
	@echo "$(GREEN)Starting all services...$(NC)"
	$(DOCKER_COMPOSE) up -d

up-build: ## Build و اجرای تمام سرویس‌ها
	@echo "$(GREEN)Building and starting all services...$(NC)"
	$(DOCKER_COMPOSE) up -d --build

down: ## متوقف کردن تمام سرویس‌ها
	@echo "$(RED)Stopping all services...$(NC)"
	$(DOCKER_COMPOSE) down

restart: ## Restart تمام سرویس‌ها
	@echo "$(YELLOW)Restarting all services...$(NC)"
	$(DOCKER_COMPOSE) restart

logs: ## نمایش logs تمام سرویس‌ها
	$(DOCKER_COMPOSE) logs -f

logs-service: ## نمایش logs سرویس خاص (استفاده: make logs-service SERVICE=user-service)
	$(DOCKER_COMPOSE) logs -f $(SERVICE)

status: ## نمایش وضعیت سرویس‌ها
	$(DOCKER_COMPOSE) ps

clean: ## پاک کردن containers، networks و volumes
	@echo "$(RED)Cleaning up...$(NC)"
	$(DOCKER_COMPOSE) down -v --remove-orphans
	docker system prune -f

clean-all: ## پاک کردن کامل (شامل images)
	@echo "$(RED)Complete cleanup...$(NC)"
	$(DOCKER_COMPOSE) down -v --remove-orphans --rmi all
	docker system prune -af

test: ## اجرای تست‌ها
	@echo "$(GREEN)Running tests...$(NC)"
	go test ./...

test-api: ## تست API endpoints
	@echo "$(GREEN)Testing API endpoints...$(NC)"
	curl -f http://localhost:8080/health || exit 1
	curl -f http://localhost:8081/health || exit 1
	curl -f http://localhost:8082/health || exit 1
	curl -f http://localhost:8083/health || exit 1

monitoring-up: ## اجرای سرویس‌های monitoring
	@echo "$(GREEN)Starting monitoring services...$(NC)"
	$(DOCKER_COMPOSE) --profile monitoring up -d

monitoring-down: ## متوقف کردن سرویس‌های monitoring
	@echo "$(RED)Stopping monitoring services...$(NC)"
	$(DOCKER_COMPOSE) --profile monitoring down

dev: ## اجرای محیط development
	@echo "$(GREEN)Starting development environment...$(NC)"
	$(DOCKER_COMPOSE) -f docker-compose.yml -f docker-compose.dev.yml up -d

prod: ## اجرای محیط production
	@echo "$(GREEN)Starting production environment...$(NC)"
	$(DOCKER_COMPOSE) -f docker-compose.yml -f docker-compose.prod.yml up -d

backup-db: ## Backup دیتابیس
	@echo "$(GREEN)Creating database backup...$(NC)"
	docker exec microservices_postgres pg_dumpall -U postgres > backup_$(shell date +%Y%m%d_%H%M%S).sql

restore-db: ## Restore دیتابیس (استفاده: make restore-db FILE=backup.sql)
	@echo "$(GREEN)Restoring database...$(NC)"
	docker exec -i microservices_postgres psql -U postgres < $(FILE)