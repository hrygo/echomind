#!/bin/bash

# scripts/dev.sh - One-click development runner for EchoMind

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}Starting EchoMind Development Environment...${NC}"

# Function to kill all child processes on exit
cleanup() {
    echo -e "\n${RED}Stopping all services...${NC}"
    kill $(jobs -p) 2>/dev/null
    wait
    echo -e "${GREEN}All services stopped.${NC}"
}

# Trap SIGINT (Ctrl+C) and SIGTERM
trap cleanup SIGINT SIGTERM

# 1. Start Infrastructure (Docker)
echo -e "${GREEN}[1/4] Starting Infrastructure (Postgres & Redis)...${NC}"
cd deploy && docker compose up -d
if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to start Docker services.${NC}"
    exit 1
fi
cd ..

# 2. Start Backend Server
echo -e "${GREEN}[2/4] Starting Backend Server...${NC}"
cd backend && go run ./cmd/main.go &
SERVER_PID=$!

# 3. Start Backend Worker
echo -e "${GREEN}[3/4] Starting Backend Worker...${NC}"
cd backend && go run ./cmd/worker/main.go &
WORKER_PID=$!

# 4. Start Frontend
echo -e "${GREEN}[4/4] Starting Frontend...${NC}"
cd frontend && pnpm dev &
FRONTEND_PID=$!

echo -e "${BLUE}EchoMind is running!${NC}"
echo -e "Backend API: http://localhost:8080"
echo -e "Frontend:    http://localhost:3000"
echo -e "${BLUE}Press Ctrl+C to stop everything.${NC}"

# Wait for all background processes
wait
