# EchoMind

**The Neural Interface for Executive Work.**
*(重塑决策直觉)*

EchoMind is a Personal Neural Interface designed for executives, managers, and dealmakers. It acts as an AI cognitive layer on top of your existing communication streams (Email, Calendar), transforming unstructured noise into structured insights, actionable tasks, and strategic intelligence.

> **Current Status**: v0.9.0 (Beta) - Phase 6.2: Smart Contexts & Actionable Intelligence

## 🚀 Core Capabilities

EchoMind doesn't just organize your email; it understands it.

### 🧠 The Cognitive Engine (RAG + AI)
*   **AI Chat Copilot**: A conversational assistant that answers questions about your work context ("What did Alice say about the budget?"), powered by RAG and streaming responses.
*   **Semantic Search**: Ask natural language questions and get answers grounded in your email history.
*   **Smart Actions**: Automatically detects actionable items (meetings, tasks) in emails and provides one-click buttons to add them to your calendar or todo list.
*   **Intent Radar**: Visualizes business signals (Buying, Hiring, Partnership) and urgency levels.
*   **Multi-Model Intelligence**: Switches between DeepSeek, OpenAI, and Gemini for cost/performance optimization.

### ⚡ The Neural Interface
*   **Mobile First**: A fully responsive design with a "Collapse-to-Expand" search bar and swipeable drawers, optimized for executives on the go.
*   **AI Command Center**: A dashboard that prioritizes attention, not just chronology.
*   **Relationship Graph**: Visualizes your network strength and "sleeping" connections.
*   **Action Center**: Extracts implied tasks and tracks them to completion.

## 📚 Documentation

*   **[Product Vision (PRD)](docs/prd.md)**: The "Why" and "What" - User personas and strategic value.
*   **[System Architecture](docs/tech-architecture.md)**: The "How" - Database schema, RAG pipeline, and API design.
*   **[Design System](docs/product-design.md)**: UI/UX principles and component breakdown.
*   **[Roadmap](docs/product-roadmap.md)**: Future plans including Team Collaboration and Mobile apps.

## 🛠️ Tech Stack

Built for performance, privacy, and scalability.

*   **Backend**: Go 1.22+ (Gin, GORM, Asynq, Viper)
*   **Frontend**: Next.js 16 (App Router, Zustand, Tailwind CSS)
*   **Database**: PostgreSQL + `pgvector` (Vector Search)
*   **Infrastructure**: Docker Compose, Redis (Queue/Cache)
*   **AI/ML**: OpenAI / DeepSeek / Gemini SDKs, LangChain concepts

## 🏁 Getting Started

### 1. Prerequisites
*   Go (1.22+)
*   Node.js (20+) & pnpm
*   Docker & Docker Compose

### 2. Quick Start

```bash
# Clone the repo
git clone https://github.com/your-username/echomind.git
cd echomind

# Initialize environment
make init
```

### 3. Configuration

1.  **Copy Config Template**: 
    ```bash
    cp backend/configs/config.example.yaml backend/configs/config.yaml
    ```
    
2.  **Update Secrets**: Edit `backend/configs/config.yaml`:
    *   Database credentials (`database.dsn`)
    *   JWT secret (`server.jwt.secret`)
    *   Encryption key (`security.encryption_key`)
    *   AI API Keys (`ai.providers.*.settings.api_key`)

### 4. Run Locally

Start the entire stack (DB, Backend, Worker, Frontend):

```bash
make dev
```

Visit `http://localhost:3000` to access the Neural Interface.

### 5. Advanced Usage

**CLI Parameters** (Introduced in v0.9.0):
```bash
# Production mode
./bin/server -production=true

# Custom config file
CONFIG_PATH=/path/to/config.yaml ./bin/server

# View all options
./bin/server -h

# Graceful shutdown
# Press Ctrl+C or send SIGTERM to gracefully stop services
```

## 🧪 Quality Assurance

*   **Backend Tests**: `make test`
*   **Frontend Tests**: `cd frontend && pnpm test`
*   **Type Checking**: `cd frontend && pnpm type-check`

## 📂 Project Structure

```text
/echomind
├── backend/           # Go API, RAG Engine & Async Workers
│   ├── cmd/           # Entry points (server, worker, reindex, backfill_contexts)
│   ├── internal/      # Core domain logic (Clean Architecture)
│   │   ├── app/       # Dependency injection & CLI config (NEW in v0.9.0)
│   │   ├── router/    # Route & middleware management (NEW in v0.9.0)
│   │   └── ...
│   └── pkg/           # Shared libraries (AI Providers, IMAP)
├── frontend/          # Next.js 16 Web Application
├── deploy/            # Docker & Deployment configs
├── docs/              # Architecture & Product specs
└── scripts/           # Dev utilities
```

## 🔄 Recent Updates (v0.9.0)

**Backend Optimizations:**
- ✅ CLI parameter support (`-config`, `-production`)
- ✅ Graceful shutdown (SIGINT/SIGTERM handling)
- ✅ Dependency injection container
- ✅ Router modularization
- ✅ Configuration centralization

**Feature Enhancements:**
- ✅ Smart Contexts (project/topic organization)
- ✅ Task management system
- ✅ Enhanced logging with request ID tracing

---

**Built with ❤️ for executives who value clarity over chaos.**