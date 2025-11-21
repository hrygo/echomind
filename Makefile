.PHONY: init run-backend run-frontend docker-up run-worker test lint build clean deploy help status tail-logs logs

# Version
VERSION := 0.3.0

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
	@echo "  make init         - Initialize project (install dependencies)"
	@echo "  make docker-up    - Start Local Docker services (Postgres & Redis)"
	@echo "  make run-backend  - Run the backend API server (Local, Background)"
	@echo "  make run-worker   - Run the backend async worker (Local, Background)"
	@echo "  make run-frontend - Run the frontend development server (Local, Background)"
	@echo "  make test         - Run backend tests"
	@echo "  make lint         - Run linters (if installed)"
	@echo "  make build        - Build backend binaries"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make dev          - Start all services (Infrastructure, Backend, Worker, Frontend in background)"
	@echo "  make stop         - Stop all services"
	@echo "  make deploy       - Deploy using production docker-compose (Requires REPO_OWNER & DB_PASSWORD)"
	@echo "  make status       - Check status of running services and last log entries"
	@echo "  make tail-logs    - Tail all service logs in real-time"

# Ensure log directory exists
logs:
	@mkdir -p $(LOG_DIR)

dev: docker-up run-backend run-worker run-frontend
	@echo "All services started. Check status with 'make status' or logs with 'make tail-logs'."

stop:
	@echo "Stopping services..."
	@pkill -f "./bin/server" || true
	@pkill -f "./bin/worker" || true
	@pkill -f "next dev" || true 
	@cd deploy && docker-compose down
	@echo "All services stopped."

init:
	@echo "Initializing backend..."
	cd backend && go mod tidy
	@echo "Initializing frontend..."
	cd frontend && pnpm install

docker-up:
	@echo "Bringing up Local Docker services..."
	cd deploy && docker-compose up -d

run-backend: logs
	@echo "Starting backend server in background, logging to $(BACKEND_LOG)..."
	@if lsof -i:8080 >/dev/null; then \
		echo "Port 8080 is already in use. Please stop the existing process or use 'make stop'."; \
	else \
		nohup ./bin/server > $(BACKEND_LOG) 2>&1 & \
		echo "Backend started (PID: $!). Check logs with 'make tail-logs'."; \
	fi

run-worker: logs
	@echo "Starting worker in background, logging to $(WORKER_LOG)..."
	@nohup ./bin/worker > $(WORKER_LOG) 2>&1 & \
	echo "Worker started (PID: $!). Check logs with 'make tail-logs'.";

run-frontend: logs
	@echo "Starting frontend development server in background, logging to $(FRONTEND_LOG)..."
	@if lsof -i:3000 >/dev/null; then \
		echo "Port 3000 is already in use. Please stop the existing process or use 'make stop'."; \
	else \
		nohup pnpm --prefix frontend dev > $(FRONTEND_LOG) 2>&1 & \
		echo "Frontend started (PID: $!). Check logs with 'make tail-logs'."; \
	fi

status:
	@echo "--- Service Status ---"
	@echo "Backend (port 8080):"
	@lsof -i:8080 -t >/dev/null && echo "  Running (PID: $(lsof -i:8080 -t))" || echo "  Not running."
	@echo "Worker:"
	@pgrep -f "./bin/worker" >/dev/null && echo "  Running (PID: $(pgrep -f "./bin/worker"))" || echo "  Not running."
	@echo "Frontend (port 3000):"
	@lsof -i:3000 -t >/dev/null && echo "  Running (PID: $(lsof -i:3000 -t))" || echo "  Not running."
	@echo "--- Last 5 Log Entries ---"
	@tail -n 5 $(BACKEND_LOG) $(WORKER_LOG) $(FRONTEND_LOG) 2>/dev/null || echo "  No logs yet or log files not found."

tail-logs:
	@echo "Tailing service logs. Press Ctrl+C to stop."
	@tail -f $(BACKEND_LOG) $(WORKER_LOG) $(FRONTEND_LOG) 2>/dev/null

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