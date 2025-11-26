# 📅 Sprint Plans Archive


# [Source: v0.6.0_rag/weekly-plan.md]
# 📅 Weekly Sprint Plan: RAG & Semantic Search (v0.6.0)

> **Goal**: Implement End-to-End Semantic Search for Emails.
> **Timeline**: 5 Days (Mon-Fri).

## 🗓️ Schedule

### Day 1: Infrastructure & Schema (Foundation)
**Objective**: Enable Vector Storage capability.

*   **Morning**: Docker & Environment
    *   [ ] Update `docker-compose.yml` to use `pgvector/pgvector:pg16`.
    *   [ ] Verify extension activation (`CREATE EXTENSION vector;`).
*   **Afternoon**: GORM & Migration
    *   [ ] Install `github.com/pgvector/pgvector-go`.
    *   [ ] Create `internal/model/embedding.go`.
    *   [ ] Implement AutoMigrate logic in `main.go`.
    *   [ ] Create `HNSW` index via raw SQL migration/init script.

### Day 2: Embedding Service (The "Eyes")
**Objective**: Connect to AI models to generate vectors.

*   **Morning**: Provider Interface
    *   [ ] Define `EmbeddingProvider` interface in `pkg/ai`.
    *   [ ] Implement `OpenAI` provider (using `text-embedding-3-small`).
    *   [ ] (Optional) Implement `DeepSeek` provider if compatible/available.
*   **Afternoon**: Text Processing
    *   [ ] Create `pkg/utils/chunker.go`.
    *   [ ] Implement HTML-to-Text stripping.
    *   [ ] Implement "Sliding Window" or "Paragraph-based" chunking strategy.

### Day 3: Ingestion Pipeline (The "Memory")
**Objective**: Automate the flow from Email -> Vector DB.

*   **Morning**: Worker Integration
    *   [ ] Update `internal/tasks/analyze.go`.
    *   [ ] Workflow: `Parse Email` -> `Summary` -> `Chunk` -> `Embed` -> `Save`.
*   **Afternoon**: Backfill & Testing
    *   [ ] Create a CLI command `make reindex` to process existing emails.
    *   [ ] Verify data in `email_embeddings` table.

### Day 4: Search Service (The "Brain")
**Objective**: Query the vector database.

*   **Morning**: Search Logic
    *   [ ] Create `internal/service/search.go`.
    *   [ ] Implement Cosine Distance search using GORM/SQL.
    *   [ ] Join results with `emails` table to get metadata (Subject, Sender).
*   **Afternoon**: API Exposure
    *   [ ] Create `internal/handler/search.go`.
    *   [ ] Register `GET /api/v1/search` endpoint.
    *   [ ] Define Request/Response DTOs (including `relevance_score`).

### Day 5: UI Integration & Polish (The Experience)
**Objective**: Deliver the feature to the user.

*   **Morning**: Frontend
    *   [ ] Update `Header.tsx` Search Bar to handle `Enter` key.
    *   [ ] Create `SearchResults` component (Dropdown or Dedicated Page).
    *   [ ] Show "Semantic Match" highlights (optional).
*   **Afternoon**: QA & Release
    *   [ ] Test queries: "Budget from last week", "Meeting requests from Alice".
    *   [ ] Performance tuning (adjust chunk size or limit).
    *   [ ] Tag v0.6.0-beta.


# [Source: v0.9.0_actionable_intelligence/plan.md]
# v0.9.0 Detailed Design: Actionable Intelligence

> **Phase**: 6.2
> **Version Target**: v0.9.0
> **Duration**: 4 Weeks
> **Theme**: From Insight to Action

---

## 1. Core Philosophy: The "Active" Dashboard

Moving from *Reading* to *Doing*. The Dashboard becomes a command center where decisions are made instantly, without navigating away.

---

## 2. Feature Specifications

### 2.1 Actionable Dashboard Cards

**Goal**: Enable "One-Click Decisions" directly from the Briefing view.

#### UI Components
1.  **Pending Decision Card**:
    *   **Content**: AI Summary of the request + Key Entities (Sender, Deadline).
    *   **Actions**:
        *   `[Approve]`: Sends a pre-generated "Approved" reply, archives email.
        *   `[Reply...]`: Opens a mini-editor with AI draft options.
        *   `[Snooze]`: Hides card for 4h/Tomorrow.
    *   **Interaction**:
        *   Click Action -> Card shows loading spinner -> Card fades out (Optimistic UI).
        *   "Undo" toast appears for 5 seconds.

2.  **Risk Warning Card**:
    *   **Actions**:
        *   `[Dismiss]`: Marks risk as handled.
        *   `[Investigate]`: Opens Chat Copilot with context pre-loaded ("Why is this high risk?").

#### API Endpoints
*   `POST /api/v1/actions/approve`: `{ email_id: "uuid" }`
*   `POST /api/v1/actions/snooze`: `{ email_id: "uuid", duration: "4h" }`

---

### 2.2 Smart Contexts (The "Project" View)

**Goal**: Slice the massive inbox into manageable "Attention Scopes".

#### Logic & Rules
A `Context` is a dynamic filter defined by:
1.  **Keywords**: "Project Alpha", "Budget", "Q4".
2.  **Key Stakeholders**: List of email addresses (client@corp.com).
3.  **Timeframe**: "Last 30 days" (rolling) or "Oct 1 - Dec 31" (fixed).

#### Interaction Flow
1.  **Creation**: User clicks "+" in Sidebar -> "Create Smart Context".
    *   Input: Name, Keywords, Key People.
    *   Preview: "Found 42 matching emails".
2.  **Activation**: Clicking a Context in Sidebar (`/dashboard?context=project_alpha`).
    *   **Dashboard**: Re-calculates stats (Risks, Tasks) *only* for this context.
    *   **Search/Chat**: RAG scope is limited to documents/emails in this context.
3.  **Auto-Tagging**:
    *   Backend `AnalyzeTask` checks new emails against active Context rules.
    *   If match: Adds `context_ids` to Email metadata (for fast filtering).

#### API Endpoints
*   `POST /api/v1/contexts`: Create definition.
*   `GET /api/v1/dashboard/stats?context_id=...`: Scoped stats.

---

### 2.3 Task Hub (Internal Task System)

**Goal**: Centralize "To-Dos" from emails and manual entry, preparing for WeChat push.

#### Data Structure (`tasks` table)
*   `source`: "email" | "manual" | "ai_inference"
*   `notify_wechat`: boolean (Default true for High priority).
*   `status`: "todo" | "in_progress" | "done" | "archived"

#### Integration Points
*   **Smart Action**: Clicking "Create Task" in Email Detail -> `POST /api/v1/tasks`.
*   **Dashboard Widget**: A dedicated "Action Items" list replacing the current static list.
    *   Support: Checkbox (Complete), Edit (Rename/Reschedule).

---

## 3. Technical Architecture & Schema

### 3.1 Database Schema (Postgres)

```sql
CREATE TABLE contexts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    color VARCHAR(20) DEFAULT 'blue',
    keywords TEXT[], -- Array of strings
    stakeholders TEXT[], -- Array of email addresses
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    source_email_id UUID REFERENCES emails(id), -- Optional link to email
    context_id UUID REFERENCES contexts(id),    -- Optional link to context
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'todo', -- todo, done
    priority VARCHAR(20) DEFAULT 'medium', -- high, medium, low
    due_date TIMESTAMPTZ,
    notify_wechat BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Many-to-Many for Email-Context caching (Optimization)
CREATE TABLE email_contexts (
    email_id UUID REFERENCES emails(id),
    context_id UUID REFERENCES contexts(id),
    PRIMARY KEY (email_id, context_id)
);
```

### 3.2 RAG Pipeline Update
*   **Ingestion**: When indexing an email, run it against `Context` rules. If match, add `context_id` to Vector Metadata.
*   **Query**: `SearchService.Search(query, contextID)` -> Adds filter `metadata['context_id'] == contextID`.

---

## 4. Implementation Roadmap (Weekly)

### Week 1: The Task Engine
*   [BE] `tasks` Table Migration & Model.
*   [BE] `TaskService` (CRUD).
*   [FE] `TaskWidget` component (List, Checkbox, optimistic updates).
*   [Integration] Connect Email Detail "Create Task" button to API.

### Week 2: The Context Brain
*   [BE] `contexts` Table & `ContextService`.
*   [BE] Rule Matcher Logic (Regex/Keyword matching).
*   [FE] Context Sidebar UI & Creator Modal.
*   [BE] Background job: Backfill contexts for existing emails.

### Week 3: Actionable Dashboard
*   [FE] Refactor Dashboard Cards to support "Actions".
*   [BE] `ActionService`: Handle `Approve`, `Snooze` (simple implementation: tag update + archive).
*   [FE] "Undo" Toast mechanism.

### Week 4: Chat & Polish
*   [AI] Prompt Engineering: "Extract tasks as JSON for Task Hub".
*   [FE] Chat Widget: Render `TaskCard` in chat stream.
*   [QA] End-to-end testing of the "Email -> Decision -> Archive" loop.

# [Source: week2_rag_polish/sprint-plan.md]
# Week 2 Sprint Plan: RAG Polish & Phase 6 Preparation

> **Sprint**: Phase 5.3 - RAG Polish + Phase 6.0 Kickoff  
> **Version Target**: v0.6.5 - v0.7.0  
> **Duration**: 5 days  
> **Dependencies**: Completed Phase 5.2 (v0.6.4)

---

## Sprint Goals

### Primary Objectives
1. **Polish RAG Features**: Performance optimization and testing
2. **Prepare Phase 6**: Design Team Collaboration architecture
3. **Production Readiness**: Monitoring, error handling, and UX improvements

### Success Criteria
- ✅ Search performance < 500ms for 10k emails
- ✅ Integration tests for search API
- ✅ E2E tests for search UI
- ✅ Team Collaboration design approved
- ✅ All tests passing, no regressions

---

## Day 1: Performance & Monitoring (v0.6.5)

### Objective
Optimize search performance and add observability.

### Morning: Performance Optimization
- [x] **Benchmark Current Performance**
    - [x] Create `backend/internal/service/search_bench_test.go`
    - [x] Test with 1k, 10k, 100k emails
    - [x] Document baseline metrics
- [x] **Optimize Chunking Strategy**
    - [x] Experiment with chunk sizes (500, 1000, 2000 chars)
    - [x] Measure embedding generation time
    - [x] Update `TextChunker` defaults if needed
- [x] **Database Query Optimization**
    - [x] Add `EXPLAIN ANALYZE` for search queries
    - [x] Tune HNSW index parameters (`m`, `ef_construction`)
    - [x] Consider materialized views for metadata

### Afternoon: Monitoring & Logging
- [x] **Add Metrics**
    - [x] Search latency histogram
    - [x] Embedding generation duration
    - [x] Vector DB query time
    - [x] Cache hit rate (if implemented)
- [x] **Structured Logging**
    - [x] Update search handlers with request IDs
    - [x] Log search queries for analytics
    - [x] Add error categorization
- [x] **Health Checks**
    - [x] Add `/health` endpoint for search service
    - [x] Verify pgvector extension status

**Deliverable**: v0.6.5 with performance improvements and monitoring

---

## Day 2: Testing & Quality (v0.6.6)

### Objective
Comprehensive testing coverage for RAG features.

### Morning: Backend Integration Tests
- [x] **Search API Tests**
    - [x] Create `backend/internal/handler/search_test.go`
    - [x] Test authentication (401 errors)
    - [x] Test query validation (empty, too long)
    - [x] Test limit parameter edge cases
- [x] **Service Layer Tests**
    - [x] Create `backend/internal/service/search_test.go`
    - [x] Mock embedding provider
    - [x] Mock database queries
    - [x] Test ranking and scoring logic
- [x] **Worker Tests**
    - [x] Update `analyze_test.go` for embedding scenarios
    - [x] Test embedding failure handling
    - [x] Test reindex command

### Afternoon: Frontend E2E Tests
- [x] **Setup E2E Framework**
    - [x] Install Playwright or Cypress
    - [x] Configure test database
- [x] **Search Flow Tests**
    - [x] Test search input interaction
    - [x] Test results display
    - [x] Test result navigation
    - [x] Test loading and error states
- [ ] **Accessibility Tests**
    - [ ] Keyboard navigation
    - [ ] Screen reader compatibility

**Deliverable**: v0.6.6 with 80%+ test coverage

---

## Day 3: UX Improvements (v0.6.7)

### Objective
Enhance user experience based on sprint feedback.

### Morning: Search UX
- [x] **Search History**
    - [x] Store recent searches in localStorage
    - [x] Show suggestions dropdown
    - [x] Clear history option
- [x] **Result Highlighting**
    - [x] Highlight matching terms in snippets
    - [x] Add relevance score explanation tooltip
- [x] **Filters**
    - [x] Add date range filter UI
    - [x] Add sender filter
    - [x] Update API to support filters

### Afternoon: Error Handling & Feedback
- [x] **Better Error Messages**
    - [x] User-friendly API error messages
    - [x] Retry mechanism for failed searches
    - [x] Offline state handling (via Error Boundary/Mock)
- [x] **Loading States**
    - [x] Skeleton loaders for results
    - [x] Progress indicator for reindex (N/A for search UI)
- [x] **Empty States**
    - [x] Onboarding tips for first search
    - [x] Suggestions when no results

**Deliverable**: v0.6.7 with enhanced UX

---

## Day 4: Phase 6 Design - Team Collaboration (Planning)

### Objective
Design architecture for Team Collaboration features.

### Morning: Requirements Analysis
- [x] **Feature Scope**
    - [x] Shared email labels/tags
    - [x] Team inboxes (shared mailboxes)
    - [x] Permission system (Owner, Admin, Member)
    - [x] Activity feed
- [x] **Data Model Design**
    - [x] `organizations` table
    - [x] `organization_members` table
    - [x] `shared_labels` table
    - [x] `team_inboxes` table
- [x] **API Design**
    - [x] Organization CRUD endpoints
    - [x] Member management endpoints
    - [x] Label sharing endpoints

### Afternoon: Implementation Plan
- [x] **Create Design Doc**
    - [x] Database schema with migrations
    - [x] API specifications
    - [x] Frontend component breakdown
    - [x] Security considerations
    - [x] Link: [Phase 6 Design](./technical_designs.md)
- [x] **Create Task Breakdown**
    - [x] Estimate effort for each feature
    - [x] Identify dependencies
    - [x] Define MVP scope

**Deliverable**: Phase 6 Design Document for review

---

## Day 5: Polish & Release Prep (v0.7.0-beta)

### Objective
Finalize RAG polish sprint and prepare for Phase 6.

### Morning: Documentation & Cleanup
- [x] **Update Documentation**
    - [x] Add search performance guide (N/A - covered in bench test comments)
    - [x] Document monitoring metrics
    - [x] Update API documentation (Verified)
- [x] **Code Cleanup**
    - [x] Remove debug logs
    - [x] Fix linting issues
    - [x] Refactor duplicate code
- [x] **Update GEMINI.md**
    - [x] Link: [Phase 6 Tasks](./summaries.md)
    - [x] Update Active Sprint to Phase 6

### Afternoon: Release Planning
- [x] **Version Management**
    - [x] Bump to v0.7.0-beta
    - [x] Tag release
    - [x] Update changelog
- [x] **QA Checklist**
    - [x] Run full test suite
    - [x] Manual testing of critical paths
    - [x] Performance benchmarks
- [x] **Deployment Prep**
    - [x] Create deployment guide
    - [x] Document environment variables
    - [x] Prepare rollback plan

**Deliverable**: v0.7.0-beta release ready

---

## Technical Details

### Performance Targets
- **Search Latency**: < 500ms (p95)
- **Embedding Generation**: < 2s per email
- **Reindex Throughput**: 100 emails/min
- **Database Query**: < 100ms

### Testing Coverage Goals
- **Backend**: 80% line coverage
- **Frontend**: 70% component coverage
- **E2E**: Critical user flows covered

### Phase 6 Preparation
- **Database Schema**: Designed and reviewed
- **API Contracts**: Defined in OpenAPI spec
- **Frontend Mockups**: Low-fidelity wireframes

---

## Risks & Mitigations

### Risk 1: Performance Regression
- **Mitigation**: Continuous benchmarking, rollback plan
- **Owner**: Backend team

### Risk 2: Test Flakiness
- **Mitigation**: Use deterministic test data, retry logic
- **Owner**: QA

### Risk 3: Scope Creep on Phase 6
- **Mitigation**: Strict MVP definition, timebox design phase
- **Owner**: Product

---

## Dependencies

### External
- ✅ pgvector extension running
- ✅ OpenAI API access
- ⚠️ Monitoring infrastructure (optional)

### Internal
- ✅ Completed Phase 5.2 (v0.6.4)
- ⚠️ Frontend E2E framework setup
- ⚠️ Performance testing infrastructure

---

## Success Metrics

### Quantitative
- Search response time: < 500ms
- Test coverage: > 75%
- Zero critical bugs in production
- 100% uptime during sprint

### Qualitative
- User feedback on search UX
- Team confidence in Phase 6 design
- Code review quality scores

---

## Next Steps After This Sprint

### Short-term (Week 3)
- Start Phase 6.1: Organization Management
- Implement basic team features
- Continue performance monitoring

### Medium-term (Month 2)
- Advanced team collaboration features
- Multi-tenant isolation
- Audit logging

### Long-term (Q1 2025+)
- Cross-platform support (Phase 7)
- Commercialization features (Phase 8)

---

## Notes

- **Focus**: Balance between polish and innovation
- **Quality**: No shortcuts on testing
- **Communication**: Daily standups, async updates
- **Flexibility**: Adjust scope based on velocity

---

**Created**: November 22, 2025  
**Sprint Lead**: TBD  
**Stakeholders**: Product, Engineering, QA


# [Source: week4_personal_intelligence/sprint-plan.md]
# Week 4 Sprint Plan: Phase 6.0 - Personal Intelligence Deep-Dive

> **Sprint**: Phase 6.0 - Personal Intelligence Deep-Dive  
> **Version Target**: v0.8.0  
> **Duration**: 1 Week  
> **Dependencies**: Completed Phase 5.3 (RAG Polish & Frontend Fixes)

---

## 1. 核心目标 (Objectives)

1.  **Chat Interface (Copilot)**:
    *   实现右侧边栏的 **AI 助手对话框**。
    *   支持自然语言指令：“帮我总结今天来自 CEO 的邮件”、“查找上周关于合同的附件”。
    *   基于 RAG 的多轮对话。
2.  **移动端适配 (Mobile First Polish)**:
    *   优化现有 Web UI 的响应式布局，确保在手机浏览器上体验流畅 (PWA 准备)。
    *   解决 Sidebar 在移动端的交互问题。
3.  **智能自动化 (Smart Actions)**:
    *   **一键动作**：基于 AI 分析结果，提供“添加到待办”、“创建日历事件”的快捷操作。

### 成功标准 (Success Criteria)
- ✅ AI Chatbot 能够进行多轮对话并回答邮件相关问题。
- ✅ 核心 Dashboard 页面在移动端浏览器上具有良好可读性和交互性。
- ✅ 邮件详情页提供至少一种 AI 驱动的快捷操作按钮。

---

## 2. 任务分解 (Task Breakdown)

#### Day 1: Chat UI & Infrastructure
- [ ] **Frontend (UI)**:
    - [ ] 设计并实现右侧 `Copilot` 抽屉/边栏组件。
    - [ ] 实现输入框、消息展示区域、加载状态。
- [ ] **Backend (API)**:
    - [ ] 创建 `POST /api/v1/chat` 接口，接收用户消息。
    - [ ] 支持流式响应 (Server-Sent Events - SSE)。

#### Day 2: Chat RAG & Context Integration
- [ ] **Backend (RAG)**:
    - [ ] 将 `POST /api/v1/chat` 接口与现有的 `SearchService` (RAG) 对接，实现邮件内容检索。
    - [ ] 实现简单的对话上下文管理（例如，保存最近 N 条对话）。
    - [ ] 对接 AI Provider (如 Gemini)。

#### Day 3: Mobile Strategy (Web & WeChat)
- [ ] **Frontend (Web Mobile Polish)**:
    - [ ] 优化 `Sidebar` 和 `Header` 在手机浏览器上的显示（折叠/隐藏）。
    - [ ] 确保核心流程（查信、看详情）在移动端 Web 可用。
    - [ ] *Note: 不做原生 App 适配，仅保证 Web 响应式。*
- [ ] **Backend (WeChat Prep)**:
    - [ ] 调研微信公众号接口 (WeChat Official Account API)。
    - [ ] 设计 `WeChatGateway` 基础结构（接收 XML 回调）。

#### Day 4: Smart Actions (Actionable AI)
- [ ] **Backend (Prompt Engineering)**:
    - [ ] 优化 AI Prompt，让 AI 输出结构化的 `suggested_actions` (例如，JSON 格式的动作列表)。
- [ ] **Frontend (Integration)**:
    - [ ] 在邮件详情页 (EmailDetailPage) 渲染 Action Buttons (e.g., "Add to To-Do", "Create Calendar Event")。
    - [ ] 实现这些动作的简单点击逻辑（模拟或调用基础 API）。

#### Day 5: Polish & Release Prep (v0.8.0)
- [ ] **E2E Tests**: 
    - [ ] 编写 Playwright E2E 测试，覆盖核心对话流程 (AI Chat)。
    - [ ] 测试移动端布局的关键元素。
- [ ] **Performance**: 检查新引入功能的性能。
- [ ] **Documentation**: 
    - [ ] 更新用户指南，解释新的 AI Copilot 和智能动作功能。
    - [ ] 更新 `CHANGELOG.md`。
- [ ] **Version Management**: 
    - [ ] 升级版本号到 `v0.8.0` (`Makefile`, `package.json`, `backend/cmd/main.go`)。

---

## 3. 风险与缓解 (Risks & Mitigations)

*   **风险**: AI Chat 质量不及预期。
    *   **缓解**: 持续优化 Prompt Engineering，探索更先进的 RAG 策略，或考虑使用更强大的 LLM。
*   **风险**: 移动端适配工作量大，影响迭代速度。
    *   **缓解**: 优先关注核心 Dashboard 视图，对非核心页面进行渐进式优化。

---

## 4. 后续计划 (Next Steps)

*   **Phase 7**: 跨平台支持 (桌面端，微信小程序)。
*   **Phase 8**: 商业化功能 (Stripe, SSO)。
