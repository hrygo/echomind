.PHONY: init run-backend run-frontend docker-up run-worker test lint build clean deploy help

# Version
VERSION := 0.2.0

# Variables
REPO_OWNER ?= your-username
DB_PASSWORD ?= change-me-in-prod

# Default target
help:
	@echo "EchoMind Makefile Commands:"
	@echo "  make init         - Initialize project (install dependencies)"
	@echo "  make docker-up    - Start Local Docker services (Postgres & Redis)"
	@echo "  make run-backend  - Run the backend API server (Local)"
	@echo "  make run-worker   - Run the backend async worker (Local)"
	@echo "  make run-frontend - Run the frontend development server (Local)"
	@echo "  make test         - Run backend tests"
	@echo "  make lint         - Run linters (if installed)"
	@echo "  make build        - Build backend binaries"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make deploy       - Deploy using production docker-compose (Requires REPO_OWNER & DB_PASSWORD)"

init:
	@echo "Initializing backend..."
	cd backend && go mod tidy
	@echo "Initializing frontend..."
	cd frontend && pnpm install

docker-up:
	@echo "Bringing up Local Docker services..."
	cd deploy && docker-compose up -d

run-backend:
	@echo "Running backend server..."
	cd backend && go run ./cmd/main.go

run-worker:
	@echo "Running worker..."
	cd backend && go run ./cmd/worker/main.go

run-frontend:
	@echo "Running frontend..."
	cd frontend && pnpm dev

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
	rm -rf bin/

deploy:
	@echo "Deploying to Production..."
	@cd deploy && ./deploy.sh $(REPO_OWNER) $(DB_PASSWORD)