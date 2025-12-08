.PHONY: init install run-backend run-worker run-frontend docker-up stop stop-apps stop-infra restart reload dev build clean test test-fe test-e2e lint lint-fe deploy help status logs logs-backend logs-worker logs-frontend watch-logs watch-backend watch-worker watch-frontend db-shell redis-shell test-coverage clean-logs ci-status build-fe migrate-db doctor health-check backup-db restore-db quick-test profile format security-scan

# =============================================================================
# EchoMind Makefile - Optimized Version v1.1.1
# =============================================================================

# Version & Configuration
VERSION := 1.1.3
REPO_OWNER ?= your-username
DB_PASSWORD ?= change-me-in-prod

# Environment Configuration
ENVIRONMENT ?= development
LOG_LEVEL ?= info
CONFIG_FILE ?= backend/configs/config.yaml

# Ports Configuration
BACKEND_PORT := 8080
FRONTEND_PORT := 3000
DB_PORT := 5432
REDIS_PORT := 6380

# Database Configuration
DB_NAME := echomind_db
DB_USER := user
DB_HOST := localhost

# Log Management
LOG_DIR := logs
BACKEND_LOG := $(LOG_DIR)/backend.log
WORKER_LOG := $(LOG_DIR)/worker.log
FRONTEND_LOG := $(LOG_DIR)/frontend.log
MIGRATION_LOG := $(LOG_DIR)/migration.log

# Build Configuration
BUILD_DIR := bin
COVERAGE_DIR := coverage
TIMESTAMP := $(shell date +%Y%m%d_%H%M%S)

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
NC := \033[0m # No Color

# =============================================================================
# Help System
# =============================================================================

.DEFAULT_GOAL := help

help:
	@echo "$(CYAN)EchoMind Development Environment$(NC)"
	@echo "$(YELLOW)Version: $(VERSION) | Environment: $(ENVIRONMENT)$(NC)"
	@echo ""
	@echo "$(BLUE)🚀 Quick Start:$(NC)"
	@echo "  make init          - Initialize project (install dependencies)"
	@echo "  make dev           - Start all services (Infrastructure + Apps)"
	@echo "  make doctor        - Check system requirements and health"
	@echo ""
	@echo "$(BLUE)📋 Development Lifecycle:$(NC)"
	@echo "  make reload        - Restart only Apps (Backend, Worker, Frontend)"
	@echo "  make restart       - Restart EVERYTHING (including Docker)"
	@echo "  make stop          - Stop EVERYTHING"
	@echo "  make stop-apps     - Stop only Apps"
	@echo "  make clean         - Clean build artifacts and logs"
	@echo ""
	@echo "$(BLUE)🧪 Quality Assurance:$(NC)"
	@echo "  make test          - Run backend tests"
	@echo "  make test-fe        - Run frontend tests"
	@echo "  make test-e2e       - Run frontend E2E tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make quick-test     - Run quick validation tests"
	@echo "  make lint          - Lint backend code"
	@echo "  make lint-fe       - Lint frontend code"
	@echo "  make format        - Format all code"
	@echo "  make security-scan  - Run security vulnerability scan"
	@echo "  make build-check   - Build both backend and frontend"
	@echo ""
	@echo "$(BLUE)🏗️  Services:$(NC)"
	@echo "  make docker-up     - Start Local Docker services"
	@echo "  make run-backend   - Build and Start Backend API"
	@echo "  make run-worker    - Build and Start Worker"
	@echo "  make run-frontend  - Start Frontend"
	@echo "  make reindex       - Reindex all emails (generate embeddings)"
	@echo ""
	@echo "$(BLUE)🗄️  Database:$(NC)"
	@echo "  make db-init      - Initialize database schema"
	@echo "  make migrate-db   - Migrate database for vector dimensions"
	@echo "  make backup-db    - Backup database to file"
	@echo "  make restore-db   - Restore database from backup"
	@echo "  make db-shell     - Open PostgreSQL shell"
	@echo "  make redis-shell  - Open Redis shell"
	@echo ""
	@echo "$(BLUE)📊 Observability:$(NC)"
	@echo "  make status        - Check service status"
	@echo "  make health-check  - Comprehensive health check"
	@echo "  make ci-status     - Check latest GitHub CI/CD pipeline"
	@echo "  make logs          - View recent logs"
	@echo "  make watch-logs    - Follow logs in real-time"
	@echo "  make profile       - Profile application performance"
	@echo ""
	@echo "$(BLUE)⚙️  Advanced:$(NC)"
	@echo "  make run-backend-prod   - Run backend in production mode"
	@echo ""
	@echo "$(BLUE)📖 Examples:$(NC)"
	@echo "  make dev ENVIRONMENT=staging     - Start in staging mode"
	@echo "  make test LOG_LEVEL=debug        - Run tests with debug logging"
	@echo "  make logs SERVICE=backend        - View only backend logs"
	@echo ""
	@echo "$(RED)⚠️  Migration Note:$(NC) Run 'make migrate-db' if you encounter vector dimension errors."

# =============================================================================
# Utility Functions
# =============================================================================

ensure-log-dir:
	@mkdir -p $(LOG_DIR) $(BUILD_DIR) $(COVERAGE_DIR)

print-section:
	@echo "$(CYAN)=============================================================================$(NC)"
	@echo "$(CYAN)$(1)$(NC)"
	@echo "$(CYAN)=============================================================================$(NC)"

print-success:
	@echo "$(GREEN)✅ $(1)$(NC)"

print-error:
	@echo "$(RED)❌ $(1)$(NC)"

print-warning:
	@echo "$(YELLOW)⚠️  $(1)$(NC)"

print-info:
	@echo "$(BLUE)ℹ️  $(1)$(NC)

# =============================================================================
# System Health and Diagnostics
# =============================================================================

doctor:
	@$(call print-section,System Health Check)
	@echo "$(BLUE)Checking system requirements...$(NC)"
	@echo ""
	@echo "$(PURPLE)Required Tools:$(NC)"
	@command -v go >/dev/null 2>&1 && echo "✅ Go $(shell go version | awk '{print $$3}')" || echo "❌ Go not found"
	@command -v node >/dev/null 2>&1 && echo "✅ Node.js $(shell node --version)" || echo "❌ Node.js not found"
	@command -v pnpm >/dev/null 2>&1 && echo "✅ pnpm $(shell pnpm --version)" || echo "❌ pnpm not found"
	@command -v docker >/dev/null 2>&1 && echo "✅ Docker $(shell docker --version | awk '{print $$3}' | sed 's/,//')" || echo "❌ Docker not found"
	@command -v docker compose >/dev/null 2>&1 && echo "✅ Docker Compose" || echo "❌ Docker Compose not found"
	@echo ""
	@echo "$(PURPLE)Services Status:$(NC)"
	@echo "Backend (8080): $$(lsof -i:8080 -t >/dev/null 2>&1 && echo "$(GREEN)🟢 Running$(NC)" || echo "$(RED)🔴 Stopped$(NC)")"
	@echo "Frontend (3000): $$(lsof -i:3000 -t >/dev/null 2>&1 && echo "$(GREEN)🟢 Running$(NC)" || echo "$(RED)🔴 Stopped$(NC)")"
	@echo "Postgres (5432): $$(nc -z localhost 5432 2>/dev/null && echo "$(GREEN)🟢 Running$(NC)" || echo "$(RED)🔴 Stopped$(NC)")"
	@echo "Redis (6380): $$(nc -z localhost 6380 2>/dev/null && echo "$(GREEN)🟢 Running$(NC)" || echo "$(RED)🔴 Stopped$(NC)")"
	@echo ""
	@echo "$(PURPLE)Configuration:$(NC)"
	@echo "Environment: $(ENVIRONMENT)"
	@echo "Config File: $(CONFIG_FILE)"
	@test -f $(CONFIG_FILE) && echo "$(GREEN)✅ Config file exists$(NC)" || echo "$(RED)❌ Config file not found$(NC)"
	@echo ""
	@echo "$(PURPLE)Directories:$(NC)"
	@test -d backend && echo "$(GREEN)✅ Backend directory$(NC)" || echo "$(RED)❌ Backend directory missing$(NC)"
	@test -d frontend && echo "$(GREEN)✅ Frontend directory$(NC)" || echo "$(RED)❌ Frontend directory missing$(NC)"
	@test -d deploy && echo "$(GREEN)✅ Deploy directory$(NC)" || echo "$(RED)❌ Deploy directory missing$(NC)"
	@echo ""
	@echo "$(PURPLE)Memory & Disk:$(NC)"
	@echo "Free Memory: $$(free -h 2>/dev/null | grep '^Mem:' | awk '{print $$7}' || echo 'N/A (not Linux)')"
	@echo "Disk Space: $$(df -h . 2>/dev/null | tail -1 | awk '{print $$4}' || echo 'N/A')"

health-check: doctor
	@echo ""
	@$(call print-section,Application Health Check)
	@echo "$(BLUE)Testing API endpoints...$(NC)"
	@curl -s http://localhost:$(BACKEND_PORT)/health >/dev/null 2>&1 && echo "$(GREEN)✅ Backend Health API$(NC)" || echo "$(RED)❌ Backend Health API$(NC)"
	@curl -s http://localhost:$(FRONTEND_PORT) >/dev/null 2>&1 && echo "$(GREEN)✅ Frontend$(NC)" || echo "$(RED)❌ Frontend$(NC)"
	@echo ""
	@echo "$(BLUE)Database connectivity...$(NC)"
	@cd backend && go run -c 'package main; import "database/sql"; import _ "github.com/lib/pq"; func main() { db, _ := sql.Open("postgres", "postgres://user:password@localhost/echomind_db?sslmode=disable"); defer db.Close(); if err := db.Ping(); err == nil { println("✅ Database connection") } else { println("❌ Database connection") } }' 2>/dev/null || echo "$(YELLOW)⚠️  Database check skipped$(NC)"

# =============================================================================
# Project Initialization
# =============================================================================

init:
	@$(call print-section,Project Initialization)
	@echo "$(BLUE)Installing backend dependencies...$(NC)"
	@cd backend && go mod download && go mod tidy
	@echo "$(BLUE)Installing frontend dependencies...$(NC)"
	@cd frontend && pnpm install
	@echo "$(BLUE)Creating necessary directories...$(NC)"
	@mkdir -p $(LOG_DIR) $(BUILD_DIR) $(COVERAGE_DIR) tmp
	@$(call print-success,Project initialized successfully!)

install: init

# =============================================================================
# Development Environment
# =============================================================================

dev: clean-logs doctor docker-up wait-for-db run-backend run-worker run-frontend
	@$(call print-section,Development Environment Ready)
	@$(call print-success,All services started!)
	@echo "$(GREEN)Backend:  http://localhost:$(BACKEND_PORT)$(NC)"
	@echo "$(GREEN)Frontend: http://localhost:$(FRONTEND_PORT)$(NC)"
	@echo "$(YELLOW)API Health: http://localhost:$(BACKEND_PORT)/health$(NC)"
	@echo ""
	@echo "$(BLUE)🔍 Checking startup logs (first 20 lines)...$(NC)"
	@sleep 3
	@if [ -f $(BACKEND_LOG) ]; then head -n 20 $(BACKEND_LOG); else echo "$(YELLOW)No backend logs yet$(NC)"; fi

reload: stop-apps build run-backend run-worker run-frontend
	@$(call print-success,Applications reloaded!)

restart: stop dev

# =============================================================================
# Service Management
# =============================================================================

stop-apps:
	@echo "$(BLUE)Stopping applications...$(NC)"
	@pkill -f "bin/server" 2>/dev/null || true
	@pkill -f "bin/worker" 2>/dev/null || true
	@pkill -f "next-server" 2>/dev/null || true
	@pkill -f "next dev" 2>/dev/null || true
	@pkill -f "pnpm dev" 2>/dev/null || true
	# Force kill if processes are still running
	@lsof -ti:$(BACKEND_PORT) | xargs kill -9 2>/dev/null || true
	@lsof -ti:$(FRONTEND_PORT) | xargs kill -9 2>/dev/null || true
	@$(call print-success,Applications stopped)

stop-infra:
	@echo "$(BLUE)Stopping infrastructure...$(NC)"
	@cd deploy && docker compose down
	@$(call print-success,Infrastructure stopped)

stop: stop-apps stop-infra
	@$(call print-success,All services stopped)

# =============================================================================
# Infrastructure
# =============================================================================

docker-up:
	@echo "$(BLUE)Starting infrastructure services...$(NC)"
	@cd deploy && docker compose up -d
	@$(call print-success,Infrastructure services started)

wait-for-db:
	@echo "$(BLUE)Waiting for Database (port 5432)...$(NC)"
	@for i in {1..30}; do \
		if nc -z localhost $(DB_PORT) 2>/dev/null; then \
			echo "$(GREEN)✅ Database is ready!$(NC)"; \
			exit 0; \
		fi; \
		sleep 1; \
		echo -n "."; \
	done; \
	echo "$(RED)❌ Database failed to start in 30s.$(NC)"; \
	exit 1

wait-for-redis:
	@echo "$(BLUE)Waiting for Redis connection...$(NC)"
	@for i in {1..15}; do \
		if nc -z localhost $(REDIS_PORT) 2>/dev/null; then \
			$(call print-success,Redis is ready!); \
			exit 0; \
		fi; \
		sleep 1; \
		echo -n "."; \
	done; \
	$(call print-error,Redis failed to start in 15 seconds); \
	exit 1

# =============================================================================
# Database Operations
# =============================================================================

db-init: docker-up wait-for-db
	@$(call print-section,Database Initialization)
	@echo "$(BLUE)Initializing database schema...$(NC)"
	@cd backend && go run cmd/db_init/main.go
	@$(call print-success,Database initialized)

migrate-db: docker-up wait-for-db
	@$(call print-section,Database Migration)
	@echo "$(YELLOW)⚠️  WARNING: This will delete existing email embeddings!$(NC)"
	@read -p "Continue? (y/N) " confirm && [ "$$confirm" = "y" ] || (echo "$(RED)Migration cancelled$(NC)" && exit 1)
	@echo "$(BLUE)Running migration...$(NC)"
	@cd deploy && docker compose exec -T db psql -U $(DB_USER) -d $(DB_NAME) < "../backend/migrations/fix_vector_dimensions.sql" 2>&1 | tee ../$(MIGRATION_LOG)
	@$(call print-success,Database migration completed. Please restart the backend service.)

backup-db: docker-up wait-for-db
	@$(call print-section,Database Backup)
	@mkdir -p backups
	@BACKUP_FILE="backup_$(TIMESTAMP).sql"; \
	echo "$(BLUE)Creating backup: $$BACKUP_FILE$(NC)"; \
	cd deploy && docker compose exec db pg_dump -U $(DB_USER) $(DB_NAME) > "../backups/$$BACKUP_FILE"; \
	$(call print-success,Backup created: backups/$$BACKUP_FILE)

restore-db: docker-up wait-for-db
	@$(call print-section,Database Restore)
	@if [ -z "$(BACKUP_FILE)" ]; then \
		echo "$(RED)Error: BACKUP_FILE environment variable not set$(NC)"; \
		echo "$(BLUE)Usage: make restore-db BACKUP_FILE=backup_20231125_120000.sql$(NC)"; \
		exit 1; \
	fi; \
	echo "$(YELLOW)⚠️  WARNING: This will overwrite current database!$(NC)"; \
	read -p "Continue? (y/N) " confirm && [ "$$confirm" = "y" ] || (echo "$(RED)Restore cancelled$(NC)" && exit 1); \
	echo "$(BLUE)Restoring from backup: $(BACKUP_FILE)$(NC)"; \
	if [ ! -f "backups/$(BACKUP_FILE)" ]; then \
		echo "$(RED)Error: Backup file backups/$(BACKUP_FILE) not found$(NC)"; \
		exit 1; \
	fi; \
	cd deploy && docker compose exec -T db psql -U $(DB_USER) $(DB_NAME) < "../backups/$(BACKUP_FILE)"; \
	$(call print-success,Database restored from backup)

db-shell:
	@echo "$(BLUE)Connecting to PostgreSQL...$(NC)"
	@cd deploy && docker compose exec db psql -U $(DB_USER) -d $(DB_NAME)

redis-shell:
	@echo "$(BLUE)Connecting to Redis...$(NC)"
	@cd deploy && docker compose exec redis redis-cli

# =============================================================================
# Application Services
# =============================================================================

run-backend: ensure-log-dir build
	@echo "$(BLUE)Starting backend server...$(NC)"
	@cd backend && LOG_FILE_PATH=../$(BACKEND_LOG) nohup ../$(BUILD_DIR)/server >> ../$(BACKEND_LOG) 2>&1 & \
	BACKEND_PID=$$!; \
	echo "$(GREEN)✅ Backend started (PID: $$BACKEND_PID) with log: $(BACKEND_LOG)$(NC)"

run-worker: ensure-log-dir build
	@echo "$(BLUE)Starting worker with dedicated log file...$(NC)"
	@cd backend && LOG_FILE_PATH=../$(WORKER_LOG) nohup ../$(BUILD_DIR)/worker >> ../$(WORKER_LOG) 2>&1 & \
	WORKER_PID=$$!; \
	echo "$(GREEN)✅ Worker started (PID: $$WORKER_PID) with log: $(WORKER_LOG)$(NC)"

run-frontend: ensure-log-dir
	@echo "$(BLUE)Starting frontend development server...$(NC)"
	@cd frontend && nohup pnpm dev >> ../$(FRONTEND_LOG) 2>&1 & \
	FRONTEND_PID=$$!; \
	echo "$(GREEN)✅ Frontend started (PID: $$FRONTEND_PID)$(NC)"

run-backend-prod: ensure-log-dir build
	@echo "$(BLUE)Starting backend in production mode...$(NC)"
	@cd backend && nohup ../$(BUILD_DIR)/server -production >> ../$(BACKEND_LOG) 2>&1 & \
	BACKEND_PID=$$!; \
	echo "$(GREEN)✅ Backend (production) started (PID: $$BACKEND_PID)$(NC)"

reindex:
	@$(call print-section,Email Reindexing)
	@echo "$(BLUE)Reindexing all emails (this may take a while)...$(NC)"
	@cd backend && go run cmd/reindex/main.go
	@$(call print-success,Email reindexing completed)

# =============================================================================
# Build System
# =============================================================================

build:
	@$(call print-section,Building Backend)
	@echo "$(BLUE)Building backend binaries...$(NC)"
	@cd backend && go build -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o ../$(BUILD_DIR)/server ./cmd/main.go
	@cd backend && go build -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o ../$(BUILD_DIR)/worker ./cmd/worker/main.go
	@$(call print-success,Build completed)

build-fe:
	@$(call print-section,Building Frontend)
	@echo "$(BLUE)Building frontend application...$(NC)"
	@cd frontend && pnpm build
	@$(call print-success,Frontend build completed)

build-check: build build-fe
	@$(call print-success,All builds completed successfully)

# =============================================================================
# Code Quality
# =============================================================================

test:
	@$(call print-section,Backend Tests)
	@echo "$(BLUE)Running backend unit tests...$(NC)"
	@cd backend && go test -v -race ./...

test-fe:
	@$(call print-section,Frontend Tests)
	@echo "$(BLUE)Running frontend tests...$(NC)"
	@cd frontend && pnpm test

test-e2e:
	@$(call print-section,Frontend E2E Tests)
	@echo "$(BLUE)Running frontend E2E tests...$(NC)"
	@bash scripts/frontend/run-tests.sh

test-coverage:
	@$(call print-section,Backend Tests with Coverage)
	@echo "$(BLUE)Running tests with coverage report...$(NC)"
	@cd backend && go test -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	@cd backend && go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@cd backend && go tool cover -func=$(COVERAGE_DIR)/coverage.out
	@echo "$(GREEN)Coverage report generated: $(COVERAGE_DIR)/coverage.html$(NC)"

quick-test:
	@$(call print-section,Quick Validation Tests)
	@echo "$(BLUE)Running quick validation...$(NC)"
	@cd backend && go test -short ./...
	@cd frontend && pnpm type-check || true
	@$(call print-success,Quick validation passed)

lint:
	@$(call print-section,Backend Linting)
	@echo "$(BLUE)Linting backend code...$(NC)"
	@cd backend && golangci-lint run ./... || echo "$(YELLOW)golangci-lint not installed or found issues$(NC)"
	@cd backend && go vet ./...
	@cd backend && go fmt ./...

lint-fe:
	@$(call print-section,Frontend Linting)
	@echo "$(BLUE)Linting frontend code...$(NC)"
	@cd frontend && pnpm lint
	@cd frontend && pnpm type-check

format:
	@$(call print-section,Code Formatting)
	@echo "$(BLUE)Formatting backend code...$(NC)"
	@cd backend && go fmt ./...
	@cd backend && goimports -w . 2>/dev/null || true
	@echo "$(BLUE)Formatting frontend code...$(NC)"
	@cd frontend && pnpm format

security-scan:
	@$(call print-section,Security Scan)
	@echo "$(BLUE)Scanning for security vulnerabilities...$(NC)"
	@cd backend && go list -json -m all | nancy sleuth 2>/dev/null || echo "$(YELLOW)Nancy security scanner not installed$(NC)"
	@cd frontend && pnpm audit

# =============================================================================
# Performance Profiling
# =============================================================================

profile:
	@$(call print-section,Performance Profiling)
	@echo "$(BLUE)Starting backend with profiling enabled...$(NC)"
	@cd backend && go run -cpuprofile=$(COVERAGE_DIR)/cpu.prof -memprofile=$(COVERAGE_DIR)/mem.prof ./cmd/main.go &
	@echo "$(GREEN)Profiling started. Stop server with Ctrl+C when done.$(NC)"
	@echo "$(BLUE)Analyze with: go tool pprof $(COVERAGE_DIR)/cpu.prof$(NC)"

# =============================================================================
# Logging and Monitoring
# =============================================================================

status:
	@$(call print-section,Service Status)
	@echo "$(PURPLE)Application Services:$(NC)"
	@echo "Backend ($(BACKEND_PORT)):  $$(lsof -i:$(BACKEND_PORT) -t >/dev/null 2>&1 && echo "$(GREEN)🟢 Running$(NC)" || echo "$(RED)🔴 Stopped$(NC)")"
	@echo "Frontend ($(FRONTEND_PORT)): $$(lsof -i:$(FRONTEND_PORT) -t >/dev/null 2>&1 && echo "$(GREEN)🟢 Running$(NC)" || echo "$(RED)🔴 Stopped$(NC)")"
	@echo "Worker:                $$(pgrep -f "bin/worker" >/dev/null 2>&1 && echo "$(GREEN)🟢 Running$(NC)" || echo "$(RED)🔴 Stopped$(NC)")"
	@echo ""
	@echo "$(PURPLE)Infrastructure:$(NC)"
	@echo "Postgres ($(DB_PORT)):  $$(nc -z localhost $(DB_PORT) 2>/dev/null && echo "$(GREEN)🟢 Running$(NC)" || echo "$(RED)🔴 Stopped$(NC)")"
	@echo "Redis ($(REDIS_PORT)):    $$(nc -z localhost $(REDIS_PORT) 2>/dev/null && echo "$(GREEN)🟢 Running$(NC)" || echo "$(RED)🔴 Stopped$(NC)")"
	@echo ""
	@echo "$(PURPLE)Process IDs:$(NC)"
	@echo "Backend:  $$(pgrep -f "bin/server" 2>/dev/null || echo "Not running")"
	@echo "Worker:   $$(pgrep -f "bin/worker" 2>/dev/null || echo "Not running")"
	@echo "Frontend: $$(pgrep -f "next dev" 2>/dev/null || echo "Not running")"

logs:
	@$(call print-section,Recent Logs)
	@if [ "$(SERVICE)" = "backend" ] || [ -z "$(SERVICE)" ]; then echo "$(BLUE)Backend logs:$(NC)"; tail -n 50 $(BACKEND_LOG) 2>/dev/null || echo "  No logs yet"; fi
	@if [ "$(SERVICE)" = "worker" ] || [ -z "$(SERVICE)" ]; then echo "$(BLUE)Worker logs:$(NC)"; tail -n 50 $(WORKER_LOG) 2>/dev/null || echo "  No logs yet"; fi
	@if [ "$(SERVICE)" = "frontend" ] || [ -z "$(SERVICE)" ]; then echo "$(BLUE)Frontend logs:$(NC)"; tail -n 50 $(FRONTEND_LOG) 2>/dev/null || echo "  No logs yet"; fi

watch-logs:
	@$(call print-section,Real-time Logs)
	@echo "$(BLUE)Following logs (Ctrl+C to exit)...$(NC)"
	@if [ "$(SERVICE)" = "backend" ]; then tail -f $(BACKEND_LOG); \
	elif [ "$(SERVICE)" = "worker" ]; then tail -f $(WORKER_LOG); \
	elif [ "$(SERVICE)" = "frontend" ]; then tail -f $(FRONTEND_LOG); \
	else tail -f $(BACKEND_LOG) $(WORKER_LOG) $(FRONTEND_LOG); fi

watch-backend:
	@echo "$(BLUE)Following backend logs...$(NC)"
	@tail -f $(BACKEND_LOG)

watch-worker:
	@echo "$(BLUE)Following worker logs...$(NC)"
	@tail -f $(WORKER_LOG)

watch-frontend:
	@echo "$(BLUE)Following frontend logs...$(NC)"
	@tail -f $(FRONTEND_LOG)

ci-status:
	@$(call print-section,CI/CD Status)
	@./scripts/check_ci.sh 2>/dev/null || echo "$(YELLOW)CI status script not available$(NC)"

# =============================================================================
# Cleanup and Maintenance
# =============================================================================

clean:
	@$(call print-section,Cleanup)
	@echo "$(BLUE)Cleaning build artifacts and logs...$(NC)"
	@rm -rf $(BUILD_DIR) $(LOG_DIR) $(COVERAGE_DIR)
	@rm -f backend/coverage.out backend/*.prof
	@rm -rf frontend/.next frontend/dist
	@$(call print-success,Cleanup completed)

clean-logs:
	@echo "$(BLUE)Cleaning old logs...$(NC)"
	@rm -f $(LOG_DIR)/*.log 2>/dev/null || true
	@$(call print-success,Logs cleaned)

# =============================================================================
# Deployment
# =============================================================================

deploy:
	@$(call print-section,Production Deployment)
	@echo "$(YELLOW)⚠️  WARNING: This will deploy to production!$(NC)"
	@read -p "Continue? (y/N) " confirm && [ "$$confirm" = "y" ] || (echo "$(RED)Deployment cancelled$(NC)" && exit 1)
	@echo "$(BLUE)Deploying to production...$(NC)"
	@cd deploy && ./deploy.sh $(REPO_OWNER) $(DB_PASSWORD)
	@$(call print-success,Deployment completed)

# =============================================================================
# Quick Commands
# =============================================================================

quick-start: init dev
	@$(call print-success,Quick start completed!)

quick-test-deploy: build-check test
	@$(call print-success,Quick test and deploy ready!)

# =============================================================================
# End of Makefile
# =============================================================================