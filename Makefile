.PHONY: init install run-backend run-worker run-frontend docker-up stop restart dev build clean test lint deploy help status logs logs-backend logs-worker logs-frontend watch-logs watch-backend watch-worker watch-frontend db-shell redis-shell test-coverage clean-logs

# Version
VERSION := 0.5.3

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
	@echo "    make install       - Alias for 'init'"
	@echo "    make dev           - Start all services (Infrastructure, Backend, Worker, Frontend) in background"
	@echo "    make restart       - Stop and restart all services"
	@echo "    make stop          - Stop all running services"
	@echo "    make clean         - Clean build artifacts and logs"
	@echo ""
	@echo "  Individual Services:"
	@echo "    make docker-up     - Start Local Docker services (Postgres & Redis)"
	@echo "    make run-backend   - Build and Start Backend API server (Background)"
	@echo "    make run-worker    - Build and Start Backend Worker (Background)"
	@echo "    make run-frontend  - Start Frontend server (Background)"
	@echo "    make build         - Build backend binaries"
	@echo ""
	@echo "  Infrastructure Interaction:"
	@echo "    make db-shell      - Connect to running Postgres database via psql"
	@echo "    make redis-shell   - Connect to running Redis instance via redis-cli"
	@echo ""
	@echo "  Testing & QA:"
	@echo "    make test          - Run backend tests"
	@echo "    make test-coverage - Run backend tests with coverage report"
	@echo "    make lint          - Run linters"
	@echo ""
	@echo "  Observability:"
	@echo "    make status        - Check status of running services and recent logs"
	@echo "    make logs          - View last 500 lines of all logs"
	@echo "    make watch-logs    - Follow (tail -f) all logs"
	@echo "    make watch-backend - Follow backend logs"
	@echo "    make watch-worker  - Follow worker logs"
	@echo "    make watch-frontend- Follow frontend logs"
	@echo ""
	@echo "  Deployment:"
	@echo "    make deploy        - Deploy to production"

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
dev: clean-logs docker-up run-backend run-worker run-frontend
	@echo "----------------------------------------------------------------"
	@echo "🚀 All services started!"
	@echo "Backend:  http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@echo "Check status with 'make status' or follow logs with 'make watch-logs'."
	@echo "----------------------------------------------------------------"
	@echo "⏳ Waiting 3 seconds for services to initialize..."
	@sleep 3
	@echo "🔍 Checking startup logs (head 100 lines)..."
	@echo "--- [Backend Log Head] ---"
	@head -n 100 $(BACKEND_LOG) 2>/dev/null || echo "Log file not created yet."
	@echo ""
	@echo "--- [Worker Log Head] ---"
	@head -n 100 $(WORKER_LOG) 2>/dev/null || echo "Log file not created yet."
	@echo ""
	@echo "--- [Frontend Log Head] ---"
	@head -n 100 $(FRONTEND_LOG) 2>/dev/null || echo "Log file not created yet."

restart: stop dev

stop:
	@echo "Stopping services..."
	@pkill -f "bin/server" || true
	@pkill -f "bin/worker" || true
	@pkill -f "next dev" || true 
	@cd deploy && docker-compose down
	@echo "All services stopped."

# Infrastructure
docker-up:
	@echo "Bringing up Local Docker services..."
	cd deploy && docker-compose up -d

db-shell:
	@echo "Connecting to Postgres..."
	@cd deploy && docker-compose exec postgres psql -U echomind -d echomind

redis-shell:
	@echo "Connecting to Redis..."
	@cd deploy && docker-compose exec redis redis-cli

# Service Runners
run-backend: ensure-log-dir build
	@echo "Starting backend server in background, logging to $(BACKEND_LOG)..."
	@if lsof -i:8080 -t >/dev/null; then \
		echo "⚠️  Port 8080 is already in use. Restarting backend..."; \
		pkill -f "bin/server" || true; \
		sleep 1; \
	fi
	@cd backend && nohup ../bin/server > ../$(BACKEND_LOG) 2>&1 & echo "✅ Backend started (PID: $$!)."

run-worker: ensure-log-dir build
	@echo "Starting worker in background, logging to $(WORKER_LOG)..."
	@pkill -f "bin/worker" || true
	@cd backend && nohup ../bin/worker > ../$(WORKER_LOG) 2>&1 & echo "✅ Worker started (PID: $$!)."

run-frontend: ensure-log-dir
	@echo "Starting frontend development server in background, logging to $(FRONTEND_LOG)..."
	@if lsof -i:3000 -t >/dev/null; then \
		echo "⚠️  Port 3000 is already in use. Process might be running."; \
	else \
		nohup pnpm --prefix frontend dev > $(FRONTEND_LOG) 2>&1 & \
		echo "✅ Frontend started."; \
	fi

# Observability
status:
	@echo "--- Service Status ---"
	@echo "Backend (port 8080):"
	@lsof -i:8080 -t >/dev/null && echo "  🟢 Running (PID: $$(lsof -i:8080 -t))" || echo "  🔴 Not running"
	@echo "Worker:"
	@pgrep -f "bin/worker" >/dev/null && echo "  🟢 Running (PID: $$(pgrep -f "bin/worker"))" || echo "  🔴 Not running"
	@echo "Frontend (port 3000):"
	@lsof -i:3000 -t >/dev/null && echo "  🟢 Running (PID: $$(lsof -i:3000 -t))" || echo "  🔴 Not running"
	@echo ""
	@echo "--- Recent Logs (Last 10 lines) ---"
	@echo "[Backend]:" && tail -n 10 $(BACKEND_LOG) 2>/dev/null || echo "  No logs."
	@echo "[Worker]:" && tail -n 10 $(WORKER_LOG) 2>/dev/null || echo "  No logs."
	@echo "[Frontend]:" && tail -n 10 $(FRONTEND_LOG) 2>/dev/null || echo "  No logs."

logs:
	@echo "Viewing last 500 lines of all logs..."
	@tail -n 500 $(BACKEND_LOG) $(WORKER_LOG) $(FRONTEND_LOG) 2>/dev/null || echo "  No logs yet or log files not found."

watch-logs:
	@echo "Following all logs (Press Ctrl+C to exit)..."
	@tail -f $(BACKEND_LOG) $(WORKER_LOG) $(FRONTEND_LOG)

watch-backend:
	@echo "Following backend logs..."
	@tail -f $(BACKEND_LOG)

watch-worker:
	@echo "Following worker logs..."
	@tail -f $(WORKER_LOG)

watch-frontend:
	@echo "Following frontend logs..."
	@tail -f $(FRONTEND_LOG)

logs-backend:
	@echo "Viewing last 500 lines of backend logs..."
	@tail -n 500 $(BACKEND_LOG)

logs-worker:
	@echo "Viewing last 500 lines of worker logs..."
	@tail -n 500 $(WORKER_LOG)

logs-frontend:
	@echo "Viewing last 500 lines of frontend logs..."
	@tail -n 500 $(FRONTEND_LOG)

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
