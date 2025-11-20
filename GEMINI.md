# 🛡️ EchoMind Project Context & Guidelines

> **Product Vision:** A SaaS-level Intelligent Email Decision System. Not just a client, but an AI cognitive layer converting unstructured email into structured insights and tasks.
> **Target:** Executives (Decision View) & Managers (Task Loop).

## 1. Development Standards (The "Golden Rules")

### 🏗️ Architecture & Structure
*   **Monorepo**: Root contains `.github`, `docs`, `scripts`, `backend`, `frontend`, `deploy`.
*   **Backend (Go)**:
    *   Frameworks: **Gin** (Web), **GORM** (ORM/Postgres), **Asynq** (Queue/Redis), **Viper** (Config), **Zap** (Log).
    *   Layout: `cmd/` (entry), `internal/` (private logic), `pkg/` (public libs).
    *   Style: **snake_case** for DB columns/JSON keys. **CamelCase** for Go structs.
*   **Frontend (Next.js)**:
    *   Stack: **Next.js 16 (App Router)**, **TypeScript**, **Tailwind CSS**, **pnpm**.
    *   Style: Functional components, Hooks-based state.
*   **Data Layer**:
    *   **PostgreSQL**: `snake_case` tables. IDs are UUID/Snowflake.
    *   **Redis**: Keys prefixed with `echomind:`.

### 🔄 Workflow & Process
1.  **TDD First**: Write the test *before* the implementation. (Red -> Green -> Refactor).
2.  **Conventional Commits**: `feat:`, `fix:`, `docs:`, `refactor:`, `chore:`.
3.  **Frequent Delivery**: Commit often (atomic commits). Release independent feature modules immediately upon completion.
4.  **Semantic Versioning**: Follow SemVer (Major.Minor.Patch). Tag releases (e.g., `v1.0.1`) upon module completion.
5.  **Documentation**: Update `GEMINI.md` (Task Status) and `tech-architecture.md` (System Changes) *before* marking a feature complete.
6.  **Safety**: Never commit secrets. Use env vars.

### 🤖 AI Agent Instructions (Self-Correction)
*   **Context**: Always read `GEMINI.md` first to understand the current Sprint status.
*   **Tool Use**: Prefer `run_shell_command` for file ops, `search_file_content` for code lookup.
*   **Testing**: Always run `go test ./...` (Backend) or `pnpm test` (Frontend) after changes.

---

# 📅 Current Sprint: Verification before Phase 4 (Completed)

**Focus**: Ensure current codebase stability by running all tests.

- [x] **Run all backend tests** (`make test`)
- [x] **Run all frontend tests** (`cd frontend && pnpm test`)

# 📅 Next Up: Phase 4 - Commercialization & Scaling

**Focus**: Multi-tenancy, User Authentication, and Payment Integration.

- [ ] **Backend: User System**:
    - [ ] Model: `User` (Email, PasswordHash, StripeCustomerID).
    - [ ] Auth: JWT-based authentication middleware.
- [ ] **Backend: Multi-tenancy**:
    - [ ] Update: `Email` and `Contact` models to include `UserID`.
    - [ ] Middleware: Enforce data isolation.
- [ ] **Commercialization**:
    - [ ] Stripe Integration (Subscription management).
    - [ ] Usage Limits (AI quotas).