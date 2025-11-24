# 🛡️ EchoMind Project Context

> **Vision**: Personal Neural Interface for Executive Work. (个人智能神经中枢)  
> **Stage**: v0.9.4 (Beta) | **Active Sprint**: Phase 7 - WeChat Integration

---

## 1. Technology Stack

**Backend**: Go 1.22+ (Gin, GORM, Asynq) | Postgres + pgvector | Redis  
**Frontend**: Next.js 16 (TypeScript, Tailwind CSS, Zustand)

---

## 2. Roadmap Status

### ✅ Recent Completion (v0.9.2 → v0.9.4)
**Phase 6.3 - The Neural Nexus (Smart Copilot & Generative Widgets)**
- ✅ **Context Bridge** (v0.9.2): Seamless Search → Chat integration with context passing.
- ✅ **Omni-Bar** (v0.9.3): Unified search/chat interface with smart routing and RAG.
- ✅ **Generative Widgets** (v0.9.4): Dynamic UI rendering in chat stream (TaskList, EmailDraft, CalendarEvent, SearchResult widgets).
- ✅ **Widget Framework**: Complete widget system with 5 widget types and E2E tests.
- ✅ **Stream Integration**: Real-time widget parsing from SSE chat stream.

### 🚧 Current Sprint (v0.9.5+ Target)
**Phase 7 - WeChat Integration (Conversational OS)**
- **Voice Commander**: WeChat voice → Whisper transcription → Intent execution.
- **One-Touch Decision**: Push approval/decision cards to WeChat for instant action.
- **Calendar Gatekeeper**: Auto-detect conflicts and generate smart replies.
- **Morning Briefing**: Daily digest push (tasks + key emails).

### 🔮 Future (6-Month Plan)
- **Phase 8** (2026.03+): Commercialization (Stripe, SSO, Multi-tier Plans)
- **Phase 9** (TBD): Team Collaboration (基础设施已就绪，功能开发暂停)

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
