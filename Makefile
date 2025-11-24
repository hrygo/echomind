.PHONY: init install run-backend run-worker run-frontend docker-up stop stop-apps stop-infra restart reload dev build clean test lint deploy help status logs logs-backend logs-worker logs-frontend watch-logs watch-backend watch-worker watch-frontend db-shell redis-shell test-coverage clean-logs ci-status

# Version
VERSION := 0.9.4

# Variables
REPO_OWNER ?= your-username
DB_PASSWORD ?= change-me-in-prod

# Log Management
LOG_DIR := logs
BACKEND_LOG := $(LOG_DIR)/backend.log
WORKER_LOG := $(LOG_DIR)/worker.log
FRONTEND_LOG := $(LOG_DIR)/frontend.log

# Default target
help:
	@echo "EchoMind Makefile Commands:"
	@echo "  Development Lifecycle:"
	@echo "    make init          - Initialize project (install dependencies)"
	@echo "    make dev           - Start all services (Infrastructure + Apps)"
	@echo "    make reload        - Restart only Apps (Backend, Worker, Frontend)"
	@echo "    make restart       - Restart EVERYTHING (including Docker)"
	@echo "    make stop          - Stop EVERYTHING"
	@echo "    make stop-apps     - Stop only Apps"
	@echo "    make clean         - Clean build artifacts and logs"
	@echo ""
	@echo "  Individual Services:"
	@echo "    make docker-up     - Start Local Docker services"
	@echo "    make run-backend   - Build and Start Backend API"
	@echo "    make run-worker    - Build and Start Worker"
	@echo "    make run-frontend  - Start Frontend"
	@echo "    make build         - Build backend binaries"
	@echo "    make reindex       - Reindex all emails (generate embeddings)"
	@echo ""
	@echo "  Observability:"
	@echo "    make status        - Check service status"
	@echo "    make ci-status     - Check latest GitHub CI/CD pipeline status"
	@echo "    make logs          - View logs"
	@echo "    make watch-logs    - Follow logs"
	@echo ""
	@echo "  Advanced:"
	@echo "    make run-backend-prod   - Run backend in production mode"
	@echo "    make build              - Build backend binaries"
	@echo "    make reindex            - Reindex all emails"
	@echo ""
	@echo "  CLI Parameters (for manual runs):"
	@echo "    ./bin/server -h              - Show backend CLI help"
	@echo "    ./bin/server -production     - Run in production mode"
	@echo "    ./bin/server -config=path    - Custom config file"

ensure-log-dir:
	@mkdir -p $(LOG_DIR)

install: init

init:
	@echo "Initializing backend..."
	cd backend && go mod tidy
	@echo "Initializing frontend..."
	cd frontend && pnpm install

clean-logs:
	@echo "Cleaning old logs..."
	@rm -f $(LOG_DIR)/*.log 2>/dev/null || true

# Development Flow
dev: clean-logs docker-up wait-for-db run-backend run-worker run-frontend
	@echo "----------------------------------------------------------------"
	@echo "🚀 All services started!"
	@echo "Backend:  http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@echo "----------------------------------------------------------------"
	@echo "🔍 Checking startup logs (head 20 lines)..."
	@sleep 2
	@head -n 20 $(BACKEND_LOG) 2>/dev/null || true

# Reload: Stop apps, rebuild, start apps (Keep DB running)
reload: stop-apps run-backend run-worker run-frontend
	@echo "♻️  Apps reloaded!"

# Restart: Stop everything, start everything
restart: stop dev

# Stop Apps only
stop-apps:
	@echo "Stopping applications..."
	@pkill -f "bin/server" || true
	@pkill -f "bin/worker" || true
	@pkill -f "next-server" || true
	@pkill -f "next dev" || true
	@# Kill processes on ports if pkill failed
	@lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	@lsof -ti:3000 | xargs kill -9 2>/dev/null || true

# Stop Infrastructure only
stop-infra:
	@echo "Stopping infrastructure..."
	@cd deploy && docker compose down

stop: stop-apps stop-infra
	@echo "🛑 All services stopped."

# Infrastructure
docker-up:
	@echo "Bringing up Local Docker services..."
	cd deploy && docker compose up -d

wait-for-db:
	@echo "⏳ Waiting for Database (port 5432)..."
	@for i in {1..30}; do \
		if nc -z localhost 5432 2>/dev/null; then \
			echo "✅ Database is ready!"; \
			exit 0; \
		fi; \
		sleep 1; \
	done; \
	echo "❌ Database failed to start in 30s."; \
	exit 1

db-shell:
	@echo "Connecting to Postgres..."
	@cd deploy && docker compose exec db psql -U echomind -d echomind

redis-shell:
	@echo "Connecting to Redis..."
	@cd deploy && docker compose exec redis redis-cli

# Service Runners
run-backend: ensure-log-dir build
	@echo "Starting backend server..."
	@cd backend && nohup ../bin/server > ../$(BACKEND_LOG) 2>&1 & echo "✅ Backend started (PID: $$!)."

run-worker: ensure-log-dir build
	@echo "Starting worker..."
	@cd backend && nohup ../bin/worker > ../$(WORKER_LOG) 2>&1 & echo "✅ Worker started (PID: $$!)."

run-frontend: ensure-log-dir
	@echo "Starting frontend..."
	@nohup pnpm --prefix frontend dev > $(FRONTEND_LOG) 2>&1 & echo "✅ Frontend started."

reindex:
	@echo "Reindexing emails..."
	@cd backend && go run cmd/reindex/main.go

db-init: docker-up wait-for-db
	@echo "Initializing database..."
	@cd backend && go run cmd/db_init/main.go

# Observability
status:
	@echo "--- Service Status ---"
	@echo "Backend (8080): $$(lsof -i:8080 -t >/dev/null && echo "🟢 Running" || echo "🔴 Stopped")"
	@echo "Frontend (3000): $$(lsof -i:3000 -t >/dev/null && echo "🟢 Running" || echo "🔴 Stopped")"
	@echo "Worker:         $$(pgrep -f "bin/worker" >/dev/null && echo "🟢 Running" || echo "🔴 Stopped")"
	@echo "Postgres (5432):$$(nc -z localhost 5432 2>/dev/null && echo "🟢 Running" || echo "🔴 Stopped")"
	@echo "Redis (6380):   $$(nc -z localhost 6380 2>/dev/null && echo "🟢 Running" || echo "🔴 Stopped")"

logs:
	@echo "Viewing last 500 lines of all logs..."
	@tail -n 500 $(BACKEND_LOG) $(WORKER_LOG) $(FRONTEND_LOG) 2>/dev/null || echo "  No logs yet."

watch-logs:
	@echo "Following all logs (Ctrl+C to exit)..."
	@tail -f $(BACKEND_LOG) $(WORKER_LOG) $(FRONTEND_LOG)

watch-backend:
	@tail -f $(BACKEND_LOG)

watch-worker:
	@tail -f $(WORKER_LOG)

watch-frontend:
	@tail -f $(FRONTEND_LOG)

ci-status:
	@./scripts/check_ci.sh

# Code Quality
test:
	@echo "Running backend tests..."
	cd backend && go test -v ./...

test-coverage:
	@echo "Running backend tests with coverage..."
	cd backend && go test -coverprofile=coverage.out ./...
	cd backend && go tool cover -func=coverage.out

lint:
	@echo "Linting backend..."
	cd backend && golangci-lint run ./... || echo "golangci-lint not installed or found issues"
	@echo "Linting frontend..."
	cd frontend && pnpm lint

build:
	@echo "Building backend..."
	cd backend && go build -o ../bin/server ./cmd/main.go
	cd backend && go build -o ../bin/worker ./cmd/worker/main.go

clean:
	@echo "Cleaning..."
	rm -rf bin/ $(LOG_DIR)/
	rm -f backend/coverage.out

deploy:
	@echo "Deploying to Production..."
	@cd deploy && ./deploy.sh $(REPO_OWNER) $(DB_PASSWORD)