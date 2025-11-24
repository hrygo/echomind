# 📅 产品路线图 - EchoMind (6-Month Plan)

**当前日期**: 2025-11-24
**当前版本**: v0.9.4 (Beta)
**规划周期**: 2025.11 - 2026.05

---

## ✅ Phase 6.2: Actionable Intelligence - COMPLETED
**时间**: 2025.11 (Month 1)
**版本**: v0.9.0 -> v0.9.1
**状态**: ✅ 已完成 (2025-11-23)
**目标**: 从“被动阅读”到“主动决策”。实现可交互的 Dashboard 和任务闭环。

*   **Smart Contexts**: 智能情境管理与自动聚合。
*   **Actionable Dashboard**: "Smart Feed" 与一键行动 (Approve/Snooze/Dismiss)。
*   **Backend Infrastructure**: 引入 `bootstrap` 包统一初始化流程。

## ✅ Phase 6.3: The Neural Nexus (智能中枢整合) - COMPLETED
**时间**: 2025.11 - 2025.12 (Month 2)
**版本**: v0.9.2 -> v0.9.4
**状态**: ✅ 已完成 (2025-11-23)
**目标**: 打破“搜索”与“对话”的边界，让用户在同一个心流中完成查找信息与分析决策。

### Stage 1: Context Bridge (上下文桥梁) [v0.9.2]
*   **目标**: 实现 Search 到 Chat 的无缝“投递”。
*   **Key Features**:
    *   [x] **Integration**: Chat Store 支持 `activeContext` 状态。
    *   [x] **UI**: 搜索结果增加 "Ask Copilot" 按钮。
    *   [x] **Chat**: 聊天组件自动读取上下文并生成回答。

### Stage 2: Omni-Bar (全能入口) [v0.9.3]
*   **目标**: 搜索框具备“路由”能力，自动识别搜索或对话意图。
*   **Key Features**:
    *   [x] **Smart Routing**: 识别问句并自动跳转 Chat。
    *   [x] **Mixed Mode**: 先搜索后分析的混合交互。

### Stage 3: Generative Widgets (生成式组件) [v0.9.4]
*   **目标**: Chat 中渲染 UI 组件。
*   **Key Features**:
    *   [x] **Widget Rendering**: Chat 流中显示 TaskList 或 SearchResultCard。

## 🗓️ Phase 7: WeChat Connect (Conversational OS)
**时间**: 2026.01 - 2026.02 (Month 3-4)
**版本**: v0.9.5+
**目标**: 将微信打造为 EchoMind 的核心移动交互系统。

*   **Conversational Core**
    *   [ ] **Voice Commander**: 微信语音 -> Whisper 转录 -> 意图执行 (回复邮件/查询)。
    *   [ ] **One-Touch Decision**: 推送 "审批/决策" 卡片，微信内直接点击 [批准]/[驳回]。
*   **Intelligent Features**
    *   [ ] **Calendar Gatekeeper**: "下周二下午有空吗？" -> 自动检测冲突并生成建议回复。
    *   [ ] **Morning Briefing**: 每日晨报推送 (今日待办 + 关键邮件)。

## 🗓️ Phase 8: 商业化 (Commercialization) - 🚀 Launch
**时间**: 2026.03+ (Month 5+)
**版本**: v1.0.0
**目标**: 正式推出付费服务。

*   **商业化 (Monetization)**
    *   [ ] **Stripe 集成**: 订阅支付。
    *   [ ] **多级套餐**: Free, Pro, Team。
    *   [ ] **企业级特性**: SSO, Audit Logs.

## ⏸️ Phase 9: 团队协作 (Team Collaboration) - ON HOLD
**状态**: 基础设施已就绪，功能开发暂停。
**目标**: 待个人用户基数稳定后，再开启团队/组织功能。

---

## 📝 版本历史

- **v0.9.4** (2025-11-23): ✅ Generative Widgets (Chat UI)
- **v0.9.3** (2025-11-23): ✅ Omni-Bar (Smart Routing & Mixed Mode)
- **v0.9.2** (2025-11-23): ✅ Context Bridge (Search -> Chat Integration)
- **v0.9.1** (2025-11-23): ✅ Smart Contexts, Actionable Dashboard, Action Service (Approve/Snooze), i18n
- **v0.9.0** (2025-11-23): ✅ Task Engine, Backend Refactor (Bootstrap)
- **v0.8.0** (2025-11-23): ✅ Smart Actions, AI Copilot, Mobile-First UI
- **v0.7.4** (2025-11-22): ✅ AI Chat Interface (Copilot) with Streaming Response
- **v0.7.0-beta** (2025-11-22): ✅ RAG Polish & Semantic Search
