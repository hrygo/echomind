# EchoMind

**The Neural Interface for Executive Work.**
*(重塑决策直觉)*

EchoMind is a Personal Neural Interface designed for executives, managers, and dealmakers. It acts as an AI cognitive layer on top of your existing communication streams (Email, Calendar), transforming unstructured noise into structured insights, actionable tasks, and strategic intelligence.

> **Current Status**: v0.6.4 (Alpha) - Phase 5.3 RAG Polish

## 🚀 Core Capabilities

EchoMind doesn't just organize your email; it understands it.

### 🧠 The Cognitive Engine (RAG + AI)
*   **Semantic Search (RAG)**: Ask natural language questions like "What was the budget decision for Project Alpha?" and get answers grounded in your email history. Powered by `pgvector` and OpenAI embeddings.
*   **Smart Briefing**: Daily "God Mode" executive summary of risks, decisions, and high-priority items.
*   **Intent Radar**: Automatically identifies and visualizes business signals (Buying, Hiring, Partnership) and urgency levels.
*   **Multi-Model Intelligence**: seamlessly switches between **DeepSeek**, **OpenAI**, and **Gemini** based on task complexity (cost/performance optimized).

### ⚡ The Neural Interface
*   **AI Command Center**: A dashboard that prioritizes attention, not just chronology.
*   **Relationship Graph**: Visualizes your network strength and "sleeping" connections.
*   **Action Center**: Extracts implied tasks ("I'll send that by Tuesday") and tracks them to completion.
*   **Sentiment Analysis**: Detects conflict and urgency (e.g., "Angry Customer", "Urgent Request") before you even open the mail.

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

1.  **Secure Config**: Copy `backend/configs/config.example.yaml` to `backend/configs/config.yaml`.
2.  **Set Secrets**: Update `config.yaml` with your Database credentials and AI API Keys (OpenAI/DeepSeek).

### 4. Run Locally
Start the entire stack (DB, Backend, Worker, Frontend) with one command:

```bash
make dev
```

Visit `http://localhost:3000` to access the Neural Interface.

## 🧪 Quality Assurance

*   **Backend Tests**: `make test`
*   **Frontend Tests**: `cd frontend && pnpm test`

## 📂 Project Structure

```text
/echomind
├── backend/           # Go API, RAG Engine & Async Workers
│   ├── cmd/           # Entry points (server, worker, reindexer)
│   ├── internal/      # Core domain logic (Clean Architecture)
│   └── pkg/           # Shared libraries (AI Providers, IMAP)
├── frontend/          # Next.js 16 Web Application
├── deploy/            # Docker & Deployment configs
├── docs/              # Architecture & Product specs
└── scripts/           # Dev utilities
```