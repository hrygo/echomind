# 🛡️ EchoMind Project Context & Guidelines

> **Product Vision:** A SaaS-level Intelligent Email Decision System. Not just a client, but an AI cognitive layer converting unstructured email into structured insights and tasks.
> **Target:** Executives (Decision View), Managers (Task Loop) & Dealmakers (Radar View).

## Development Standards (The "Golden Rules")

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

## 📅 Current Sprint: Phase 5 - Commercialization (v0.5.0)

**Focus**: Stripe Integration, Usage Limits, Team Collaboration.
**Status**: Backlog.

*Refer to `docs/tech-design-phase5.md` for architecture and implementation details (when available).*

- [ ] **Backend: Monetization**:
    - [ ] **Stripe Integration**: Implement Stripe webhook handler for subscription events.
    - [ ] **Usage Limits**: Track and enforce user usage limits for AI features.
    - [ ] **Team Accounts**: Support multiple users under a single organizational account.
- [ ] **Frontend: Billing & Teams**:
    - [ ] **Billing Page**: Create a dedicated billing and subscription management UI.
    - [ ] **Team Settings**: Develop UI for managing team members and permissions.

---

## ✅ Completed Sprints

## Phase 4: Deep Insight & Relationship Intelligence (v0.5.0)
**Focus**: Contact Analytics, Relationship Graph, AI Smart Reply.
**Status**: Completed.

*Refer to `docs/tech-design-phase4.md` for detailed specifications, acceptance criteria, and implementation plan.*

- [x] **Backend: Contact Intelligence**:
    - [x] **Model**: Update `Contact` with `InteractionCount`, `AvgSentiment`.
    - [x] **Pipeline**: Update `AnalyzeTask` to aggregate contact stats.
    - [x] **Migration**: Backfill stats for existing data.
- [x] **Backend: Insight API**:
    - [x] **API**: `GET /insights/network` (Nodes/Links).
    - [x] **API**: `POST /ai/draft` (Generate reply).
- [x] **Frontend: Visuals**:
    - [x] **Graph**: Implement Relationship Network visualization.
    - [x] **Smart Reply**: Add AI Draft button in Email Detail.

---

## 🚧 Backlog: Future


---

## 🚧 Backlog: Future
