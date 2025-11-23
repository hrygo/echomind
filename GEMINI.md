# 🛡️ EchoMind Project Context

> **Vision**: Personal Neural Interface for Executive Work. (个人智能神经中枢)  
> **Stage**: v0.9.3 (Beta) | **Active Sprint**: Phase 6.3 - Stage 3: Generative Widgets

---

## 1. Technology Stack

**Backend**: Go 1.22+ (Gin, GORM, Asynq) | Postgres + pgvector | Redis  
**Frontend**: Next.js 16 (TypeScript, Tailwind CSS, Zustand)

---

## 2. Roadmap Status

### ✅ Recent Completion (v0.9.2 → v0.9.3)
**Phase 6.3 - Stage 2: Smart Copilot (Omni-Bar)**
- ✅ **Omni-Bar**: Unified Header Search and Chat into a single context-aware input (`CopilotWidget`).
- ✅ **RAG Integration**: `ChatService` prioritizes explicit context (`context_ref_ids`) from search results before falling back to auto-search.
- ✅ **Streaming**: Implemented real-time token streaming for Chat using SSE (`/api/v1/chat/completions`).
- ✅ **UX Polish**: Seamless mode switching between "Instant Search" and "AI Chat".
- ✅ **State Management**: Centralized `useCopilotStore` for managing Search/Chat modes and data.

### 🚧 Current Sprint (v0.9.4 Target)
**Phase 6.3 - Stage 3: Generative Widgets**
- **Widget Protocol**: Standardize the JSON schema for widgets embedded in AI response streams.
- **Frontend Renderer**: Implement a `WidgetRenderer` component to dynamically display `TaskCard`, `CalendarSlot`, and `SearchResultCard` within the chat stream.
- **Backend Support**: Enhance `ChatService` to detect intent and inject widget data structures into the response stream.

### 🔮 Future (6-Month Plan)
- **Phase 7** (2026.03-04): WeChat Integration (Official Account)
- **Phase 8** (2026.05+): Commercialization (Stripe, SSO)
- **Phase 9** (TBD): Team Collaboration (重新评估优先级)

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
- **Precise Replacement**: When using `replace`, ensure `old_string` 是唯一且最小的，避免包含可能已更改的长上下文。
- **Verify State**: 如果工具失败，使用 `read_file` 验证当前文件状态，然后再重试。
