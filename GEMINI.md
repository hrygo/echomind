# 🛡️ EchoMind Project Context

> **Vision**: AI-powered Email Decision System. "External Brain for Executive Communication."
> **Stage**: v0.5.1 (Stable/Beta).
> **Active Sprint**: Phase 5 - Commercialization.

## 1. The Golden Rules (Non-Negotiable)

1.  **Frequent Delivery**: Commit often. Don't hoard changes.
2.  **Semantic Versioning**: `vMajor.Minor.Patch`. Update `Makefile` & `package.json` before tagging.
3.  **Test-Driven**: Red -> Green -> Refactor. `make test` MUST pass before commit.
4.  **Convention over Configuration**: Follow existing directory structure (`internal/`, `pkg/`) and naming.
5.  **Monorepo Discipline**: Respect the `backend/` vs `frontend/` boundary.

## 2. Technology Stack

*   **Backend**: Go 1.22+
    *   *Web*: Gin
    *   *DB*: GORM (Postgres)
    *   *Async*: Asynq (Redis)
    *   *Config*: Viper
*   **Frontend**: Next.js 16 (TypeScript)
    *   *UI*: Tailwind CSS + Radix UI
    *   *State*: Zustand
    *   *Fetch*: Axios

## 3. Development Workflow

### 🚀 Start
*   `make dev`: Starts Postgres, Redis, Backend, Worker, and Frontend.
*   `make logs`: View all logs.

### 🧪 Verify
*   `make test`: Run backend unit/integration tests.
*   `cd frontend && pnpm test`: Run frontend component tests.

### 📦 Release
1.  Finish Feature/Fix.
2.  Verify Tests (`make test`).
3.  Bump Version (`Makefile`, `frontend/package.json`, `backend/cmd/main.go`).
4.  Commit: `feat: ...` or `fix: ...`.
5.  Tag: `git tag vX.Y.Z`.

## 4. Roadmap Status

### ✅ Completed (v0.1.0 - v0.5.1)
*   **Core**: IMAP Sync, Email Parsing, Auth (JWT).
*   **AI**: Summary, Sentiment, Classification, Contact Intelligence, Smart Reply.
*   **UI**: Dashboard, Insights Graph, Account Settings.

### 🚧 In Progress (Phase 5: v0.6.0 Target)
*   **Monetization**: Stripe Integration, Webhooks.
*   **Quotas**: Usage limits per user/org.
*   **Teams**: Organization support.

### 🔮 Future
*   **Desktop App**: Flutter/Electron.
*   **Advanced RAG**: Search across email history.