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

# 📅 Current Sprint: Phase 2 - Intelligent Analysis & Insights (v0.3.0)

**Focus**: AI Classification, Action Item Extraction, Smart Dashboard.

- [ ] **Backend: AI Analysis Engine**:
    - [ ] **Refine Summary**: Update `Summarize` prompt to return structured JSON (Category, Sentiment, Action Items).
    - [ ] **Model Update**: Add `Category`, `ActionItems` (JSONB) to `Email` model.
    - [ ] **Classifier**: Implement logic to auto-tag emails (Work, Personal, Newsletter).
- [ ] **Frontend: Smart Dashboard**:
    - [ ] **Filters**: Add Sidebar/Tabs for `Category` filtering.
    - [ ] **Action Items**: Display extracted tasks in Email Detail view.
    - [ ] **Visuals**: Add badges for `Sentiment` and `Urgency`.

---

# 🚧 Backlog: Advanced Features & Commercialization

- [ ] **Phase 3: Deep Integration**:
    - [ ] Relationship Graph.
    - [ ] Calendar Integration.
- [ ] **Phase 4: Commercialization**:
    - [ ] Stripe Integration.
    - [ ] Usage Limits.
