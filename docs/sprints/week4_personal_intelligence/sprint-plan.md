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
