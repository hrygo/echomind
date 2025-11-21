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

# 📅 Current Sprint: Phase 3 - Real-World Sync Integration (v0.4.0)

**Focus**: User Credentials Management, Dynamic IMAP Connection, Settings UI.
**Status**: Design Complete. Ready for Implementation.

*Refer to `docs/tech-design-phase3.md` for detailed specifications, acceptance criteria, and implementation plan.*

- [ ] **Backend: Account Management**:
    - [ ] **Model**: Create `EmailAccount` model (UserID, Server, Port, Username, EncryptedPassword).
    - [ ] **Security**: Implement AES encryption for storing IMAP passwords (at rest).
    - [ ] **API**: Add endpoints to `POST /settings/account` (Connect) and `GET /settings/account` (Status).
- [ ] **Backend: Dynamic Sync Engine**:
    - [ ] **Refactor**: Update `SyncService` to load user credentials and create a fresh IMAP connection per request.
    - [ ] **Error Handling**: Handle connection failures (Auth error, Timeout) and update account status.
- [ ] **Frontend: Settings & Connect**:
    - [ ] **UI**: Build "Connect Email" form in `Settings` page (Host, Port, User, Password).
    - [ ] **Status**: Display connection status (Connected/Failed) and "Last Synced At".

---

# ✅ Completed Sprints

## Phase 2: Intelligent Analysis & Insights (v0.3.0)
- [x] **Backend**: AI Analysis Engine (Summary, Category, Sentiment, Action Items).
- [x] **Backend**: Auto-Classifier Logic.
- [x] **Frontend**: Smart Dashboard (Filters, Insights Card, Visual Badges).

---

# 🚧 Backlog: Advanced Features & Commercialization

- [ ] **Phase 3: Deep Integration**:
    - [ ] Relationship Graph.
    - [ ] Calendar Integration.
- [ ] **Phase 4: Commercialization**:
    - [ ] Stripe Integration.
    - [ ] Usage Limits.
