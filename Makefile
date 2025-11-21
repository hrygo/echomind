.PHONY: init install run-backend run-frontend docker-up run-worker test lint build clean deploy help status view-logs logs restart view-logs-backend view-logs-worker view-logs-frontend

# Version
VERSION := 0.4.0

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
	@echo "    make run-backend   - Start Backend API server (Background)"
	@echo "    make run-worker    - Start Backend Worker (Background)"
	@echo "    make run-frontend  - Start Frontend server (Background)"
	@echo "    make build         - Build backend binaries"
	@echo ""
	@echo "  Testing & QA:"
	@echo "    make test          - Run backend tests"
	@echo "    make lint          - Run linters"
	@echo "    make status        - Check status of running services and recent logs"
	@echo "    make view-logs     - View last 500 lines of all logs"
	@echo "    make view-logs-backend  - View last 500 lines of backend logs only"
	@echo "    make view-logs-worker   - View last 500 lines of worker logs only"
	@echo "    make view-logs-frontend - View last 500 lines of frontend logs only"
	@echo ""
	@echo "  Deployment:"
	@echo "    make deploy        - Deploy to production"

# Ensure log directory exists
logs:
	@mkdir -p $(LOG_DIR)

install: init

init:
	@echo "Initializing backend..."
	cd backend && go mod tidy
	@echo "Initializing frontend..."
	cd frontend && pnpm install

# Development Flow
dev: docker-up run-backend run-worker run-frontend
	@echo "----------------------------------------------------------------"
	@echo "🚀 All services started!"
	@echo "Backend:  http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@echo "Check status with 'make status' or logs with 'make view-logs'."
	@echo "----------------------------------------------------------------"

restart: stop dev

stop:
	@echo "Stopping services..."
	@pkill -f "./bin/server" || true
	@pkill -f "./bin/worker" || true
	@pkill -f "next dev" || true 
	@cd deploy && docker-compose down
	@echo "All services stopped."

# Infrastructure
docker-up:
	@echo "Bringing up Local Docker services..."
	cd deploy && docker-compose up -d

# Service Runners
run-backend: logs
	@echo "Starting backend server in background, logging to $(BACKEND_LOG)..."
	@if lsof -i:8080 >/dev/null; then \
		echo "⚠️  Port 8080 is already in use. Process might be running."; \
	else \
		nohup ./bin/server > $(BACKEND_LOG) 2>&1 & \
		echo "✅ Backend started (PID: $!)."; \
	fi

run-worker: logs
	@echo "Starting worker in background, logging to $(WORKER_LOG)..."
	@nohup ./bin/worker > $(WORKER_LOG) 2>&1 & \
	@echo "✅ Worker started (PID: $!).";

run-frontend: logs
	@echo "Starting frontend development server in background, logging to $(FRONTEND_LOG)..."
	@if lsof -i:3000 >/dev/null; then \
		echo "⚠️  Port 3000 is already in use. Process might be running."; \
	else \
		nohup pnpm --prefix frontend dev > $(FRONTEND_LOG) 2>&1 & \
		echo "✅ Frontend started (PID: $!)."; \
	fi

# Observability
status:
	@echo "--- Service Status ---"
	@echo "Backend (port 8080):"
	@lsof -i:8080 -t >/dev/null && echo "  🟢 Running (PID: $(lsof -i:8080 -t))" || echo "  🔴 Not running"
	@echo "Worker:"
	@pgrep -f "./bin/worker" >/dev/null && echo "  🟢 Running (PID: $(pgrep -f "./bin/worker"))" || echo "  🔴 Not running"
	@echo "Frontend (port 3000):"
	@lsof -i:3000 -t >/dev/null && echo "  🟢 Running (PID: $(lsof -i:3000 -t))" || echo "  🔴 Not running"
	@echo ""
	@echo "--- Recent Logs (Last 10 lines) ---"
	@echo "[Backend]:" && tail -n 10 $(BACKEND_LOG) 2>/dev/null || echo "  No logs."
	@echo "[Worker]:" && tail -n 10 $(WORKER_LOG) 2>/dev/null || echo "  No logs."
	@echo "[Frontend]:" && tail -n 10 $(FRONTEND_LOG) 2>/dev/null || echo "  No logs."

view-logs:
	@echo "Viewing last 500 lines of all logs..."
	@tail -n 500 $(BACKEND_LOG) $(WORKER_LOG) $(FRONTEND_LOG) 2>/dev/null || echo "  No logs yet or log files not found."

view-logs-backend:
	@echo "Viewing last 500 lines of backend logs..."
	@tail -n 500 $(BACKEND_LOG)

view-logs-worker:
	@echo "Viewing last 500 lines of worker logs..."
	@tail -n 500 $(WORKER_LOG)

view-logs-frontend:
	@echo "Viewing last 500 lines of frontend logs..."
	@tail -n 500 $(FRONTEND_LOG)

# Code Quality
test:
	@echo "Running backend tests..."
	cd backend && go test -v ./...

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

deploy:
	@echo "Deploying to Production..."
	@cd deploy && ./deploy.sh $(REPO_OWNER) $(DB_PASSWORD)