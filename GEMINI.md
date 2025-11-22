# 🛡️ EchoMind Project Context

> **Vision**: Personal Neural Interface for Executive Work. (个人智能神经中枢)  
> **Stage**: v0.7.2 (Alpha) | **Active Sprint**: Phase 6.1 - Team Collaboration Foundation

---

## 1. Technology Stack

**Backend**: Go 1.22+ (Gin, GORM, Asynq) | Postgres + pgvector | Redis  
**Frontend**: Next.js 16 (TypeScript, Tailwind CSS, Zustand)

---

## 2. Roadmap Status

### ✅ Recent Completion (v0.7.0 → v0.7.2)
**Phase 6.1 - Team Collaboration Foundation**
- ✅ Multi-tenant Architecture (Org/Team Models)
- ✅ Migration Strategy & Execution
- ✅ Frontend Context Switcher & State Management

### 🚧 Current Sprint (Phase 6.2+)
**Phase 6.2 - Shared Resources & Team Polish**
- Shared Email Inboxes
- Team Member Management UI
- Advanced Permission System

### 🔮 Future (6-Month Plan)
- **Phase 7** (2026.03-04): Cross-Platform (Desktop, WeChat)
- **Phase 8** (2026.05+): Commercialization (Stripe, SSO)

---

## 3. The Golden Rules (Non-Negotiable)

### 🛡️ Quality & Standards (Test-Driven)
1. **CI Mandatory**: `make test` (Backend) AND `pnpm build` (Frontend) MUST pass before commit.
   - Frontend must run `pnpm type-check` to catch strict type errors.
2. **Mock First**: Use mocks for external dependencies (AI, DB) in unit tests to ensure speed and stability.

### 🚀 Delivery (Frequent & Versioned)
1. **Commit Often**: Don't hoard changes. Atomic commits.
2. **Tag Immediately**: Release often.
   - Minor features (v0.x.Y): Daily if tests pass.
   - Fixes (v0.x.y): Immediate.
3. **Convention**: `feat:` | `fix:` | `docs:` | `refactor:`
4. **Versioning**: Update `Makefile`, `package.json`, `backend/cmd/main.go`.

### 🏗️ Architecture & Code Standards
1. **Refactoring Protocol (Blast Radius Control)**:
   - **Incremental**: When changing core APIs (e.g., `apiClient`), keep the old export deprecated temporarily.
   - **Search First**: Use `grep` or global search to identify ALL usage points before modifying types or exports.
2. **Frontend Components**:
   - **Check Existence**: Never assume a UI component (e.g., `Dialog`) exists. Check `src/components/ui` first.
   - **Atomic UI**: New features needing new UI components must include the component code in the commit.
3. **Database Schema**:
   - **Type Safety**: Changing a model field (e.g., `UUID` to `*UUID`) breaks code. Compile backend immediately after model changes.
   - **Compatibility**: Avoid DB-specific defaults (e.g., `gen_random_uuid()`) in GORM tags if they break SQLite tests. Generate IDs in application logic.

### 🌐 Internationalization (i18n)
- **Bilingual UI**: All user-facing text MUST support both English (en) and Chinese (zh).
- **No Hardcoding**: Use `t('key')`.

### 🔧 Tooling Usage (AI Agent SOP)
- **Precise Replacement**: When using `replace`, ensure `old_string` is unique and minimal. Avoid including long context that might have changed.
- **Verify State**: If a tool fails, use `read_file` to verify the current file state before retrying.