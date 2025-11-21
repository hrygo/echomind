# 🛡️ EchoMind Project Context & Guidelines

> **Product Vision:** A SaaS-level Intelligent Email Decision System. Not just a client, but an AI cognitive layer converting unstructured email into structured insights and tasks.
> **Target:** Executives (Decision View), Managers (Task Loop) & Dealmakers (Radar View).

## 1. Development Standards (The "Golden Rules")

### 🏗️ Architecture & Structure
*   **Monorepo**: Root contains `.github`, `docs`, `scripts`, `backend`, `frontend`, `deploy`.
*   **Documentation**: All architecture/product docs reside in `docs/`. Root contains meta-docs (`README`, `CONTRIBUTING`).
*   **Backend (Go)**:
    *   Frameworks: **Gin**, **GORM**, **Asynq**, **Viper**, **Zap**.
    *   Layout: `cmd/` (entry), `internal/` (private), `pkg/` (public).
*   **Frontend (Next.js)**:
    *   Stack: **Next.js 16 (App Router)**, **TypeScript**, **Tailwind CSS**, **pnpm**.

### 🔄 Workflow & Process
1.  **TDD First**: Write tests before code. (Red -> Green -> Refactor).
2.  **Conventional Commits**: `feat:` (feature), `fix:` (bug), `docs:` (docs), `refactor:` (code restructure), `chore:` (maintain).
3.  **Frequent Delivery**: Commit often. Release features independently.
4.  **Semantic Versioning & Release**: 
    *   Update version numbers in `frontend/package.json`, `Makefile`, and `docs/` *before* tagging.
    *   Tag releases (vMajor.Minor.Patch).
5.  **Documentation**: Keep `docs/` updated. Read `GEMINI.md` for context.
6.  **Safety**: Secrets via Env Vars only.

### 🤖 AI Agent Instructions
*   **Context**: Read `GEMINI.md` first.
*   **Tool Use**: Use `run_shell_command` for file ops.
*   **Testing**: Always verify with `make test` (Backend) or `pnpm test` (Frontend).

---

# 🛠️ Debugging & Local Development Guidelines

## 1. Service Startup & Management (Makefile)
When working on features or debugging locally, **always use the provided `Makefile` commands**. This ensures consistent setup, environment variable loading, and centralized log management.

### 🚀 Quick Start
*   **`make install`**: Install dependencies (Backend & Frontend).
*   **`make dev`**: Start ALL services in background (Postgres, Redis, Backend, Worker, Frontend).
*   **`make restart`**: Restart all services (useful after config changes).
*   **`make stop`**: Stop all services.

### 🔍 Debugging & Logs
*   **`make status`**: Check if services are running (PIDs) and view the last few log lines.
*   **`make logs`**: View last 500 lines of all logs.
*   **`make logs-backend`**: View last 500 lines of backend logs only.
*   **`make logs-worker`**: View last 500 lines of worker logs only.
*   **`make logs-frontend`**: View last 500 lines of frontend logs only.

## 2. Development Workflow
1.  **Feature Branch**: Create a branch for your task (e.g., `feat/stripe-integration`).
2.  **Local Dev**: Use `make dev` to run the stack.
3.  **Iterate**: Edit code -> `make restart` (if Backend/Worker code changed) -> Verify.
4.  **Test**: Run `make test` (Backend) and `pnpm test` (Frontend) before committing.
5.  **Commit**: Follow Conventional Commits (e.g., `feat: add stripe webhook handler`).

## 3. Configuration Management
*   **`backend/configs/config.example.yaml`**: Template for configuration.
*   **`backend/configs/config.yaml`**: **Ignored by Git**. Copy from example and fill in your local secrets (API Keys, DB Creds).
*   **Environment Variables**: The application supports overriding config via env vars (e.g., `ECHOMIND_DATABASE_DSN`).

### ✅ Debugging Rule: Prioritize `make` commands
When debugging or verifying changes locally, **always prefer using `make` commands** (e.g., `make run-backend`, `make status`, `make logs`) over directly executing binaries or raw shell commands. This ensures consistency with the project's defined environment and logging practices.

# 📅 Current Sprint: Phase 4 - Deep Insight & Relationship Intelligence (v0.5.0)

**Focus**: Contact Analytics, Relationship Graph, AI Smart Reply.
**Status**: Design Ready.

*Refer to `docs/tech-design-phase4.md` for architecture and implementation details.*

- [ ] **Backend: Contact Intelligence**:
    - [ ] **Model**: Update `Contact` with `InteractionCount`, `AvgSentiment`.
    - [ ] **Pipeline**: Update `AnalyzeTask` to aggregate contact stats.
    - [ ] **Migration**: Backfill stats for existing data.
- [ ] **Backend: Insight API**:
    - [ ] **API**: `GET /insights/network` (Nodes/Links).
    - [ ] **API**: `POST /ai/draft` (Generate reply).
- [ ] **Frontend: Visuals**:
    - [ ] **Graph**: Implement Relationship Network visualization.
    - [ ] **Smart Reply**: Add AI Draft button in Email Detail.

---

# ✅ Completed Sprints

## Phase 3: Real-World Sync Integration (v0.4.0)
**Focus**: User Credentials Management, Dynamic IMAP Connection, Settings UI.
**Status**: Completed.

*Refer to `docs/tech-design-phase3.md` for detailed specifications, acceptance criteria, and implementation plan.*

- [x] **Backend: Account Management**:
    - [x] **Model**: Create `EmailAccount` model (UserID, Server, Port, Username, EncryptedPassword).
    - [x] **Security**: Implement AES encryption for storing IMAP passwords (at rest).
    - [x] **API**: Add endpoints to `POST /settings/account` (Connect) and `GET /settings/account` (Status).
- [x] **Backend: Dynamic Sync Engine**:
    - [x] **Refactor**: Update `SyncService` to load user credentials and create a fresh IMAP connection per request.
    - [x] **Error Handling**: Handle connection failures (Auth error, Timeout) and update account status.
- [x] **Frontend: Settings & Connect**:
    - [x] **UI**: Build "Connect Email" form in `Settings` page (Host, Port, User, Password).
    - [x] **Status**: Display connection status (Connected/Failed) and "Last Synced At".

## Phase 2: Intelligent Analysis & Insights (v0.3.0)
- [x] **Backend**: AI Analysis Engine (Summary, Category, Sentiment, Action Items).
- [x] **Backend**: Auto-Classifier Logic.
- [x] **Frontend**: Smart Dashboard (Filters, Insights Card, Visual Badges).

---

# 🚧 Backlog: Future

- [ ] **Phase 5: Commercialization**:
    - [ ] Stripe Integration.
    - [ ] Usage Limits.
    - [ ] Team Collaboration.
- [ ] **Phase 6: Mobile & Desktop Native**:
    - [ ] Flutter / Electron App.
