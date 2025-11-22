# 📅 产品路线图 - EchoMind (6-Month Plan)

**当前日期**: 2025-11-22
**当前版本**: v0.7.2 (Alpha)
**规划周期**: 2025.12 - 2026.05

---

## ✅ Phase 5: 深度优化与 RAG (Core Intelligence) - COMPLETED
**时间**: 2025.11 - 2025.12 (Month 1)
**版本**: v0.6.0 -> v0.7.2
**状态**: ✅ 已完成 (2025-11-22)
**目标**: 构建私有知识库，实现自然语言搜索，多租户基础架构。

*   **RAG 引擎**: 向量基础设施，语义搜索。
*   **架构升级**: 多租户 (Org/Team) 数据库模型与迁移完成。
*   **工程质量**: 完善 E2E 测试与 CI 流程。

## 🚧 Phase 6.0: Personal Intelligence Deep-Dive (Current Sprint)
**时间**: 2025.12 (Month 2)
**版本**: v0.8.0
**目标**: 深化个人 AI 体验，聚焦 Web 端交互与微信生态。

*   **AI Copilot (Web)**
    *   [ ] **Chat Sidebar**: 右侧边栏对话框，支持自然语言查信。
    *   [ ] **Smart Actions**: 基于 AI 分析的快捷操作 ("添加到待办", "草拟回复")。
*   **移动端策略 (Mobile Strategy)**
    *   [ ] **Web Mobile Polish**: 优化手机浏览器上的 Dashboard 阅读体验。
    *   [ ] **WeChat Integration**: 唯一的官方移动端入口（微信公众号）。
    *   *注：原生 App / 小程序计划待定 (TBD)，优先资源投入 Web 与微信。*

## 🗓️ Phase 7: WeChat Connect (Conversational OS)
**时间**: 2026.01 - 2026.02 (Month 3-4)
**版本**: v0.9.0
**目标**: 将微信打造为 EchoMind 的核心移动交互系统。

*   **Conversational Core**
    *   [ ] **Voice Commander**: 微信语音 -> Whisper 转录 -> 意图执行 (回复邮件/查询)。
    *   [ ] **One-Touch Decision**: 推送 "审批/决策" 卡片，微信内直接点击 [批准]/[驳回]。
    *   [ ] **Thought Catcher**: 随时语音记录灵感/待办，自动同步到 Dashboard。
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

- **v0.7.2** (2025-11-22): ✅ Team Foundation & Fixes
- **v0.7.0-beta** (2025-11-22): ✅ RAG Polish
- **v0.6.4** (2025-11-22): ✅ RAG & Semantic Search
- **v0.6.0** (2025-11-22): pgvector Infrastructure