# 📅 产品路线图 - EchoMind (6-Month Plan)

**当前日期**: 2025-11-21
**规划周期**: 2025.12 - 2026.05

---

## 🗓️ Phase 5: 深度优化与 RAG (Core Intelligence)
**时间**: 2025.11 - 2025.12 (Month 1)
**版本**: v0.6.0 -> v0.6.x
**目标**: 构建私有知识库，实现自然语言搜索，让 Dashboard 数据“活”起来。

*   **RAG 引擎 (Vector & Memory)**
    *   [ ] **向量基础设施**: 集成 `pgvector`，实现邮件内容的向量化存储。
    *   [ ] **语义搜索 (Semantic Search)**: "帮我找上周关于预算的邮件"，取代关键词搜索。
    *   [ ] **自动记忆 (Auto-Memory)**: 自动提取并维护 "项目-联系人" 关系图谱 (Graph RAG 基础)。
*   **深度优化 (Optimization)**
    *   [x] **Spam Filtering**: 规则引擎过滤垃圾邮件。
    *   [ ] **本地缓存策略**: 优化前端缓存，减少网络请求。

## 🗓️ Phase 6: The Neural Interface (UI/UX Refactor)
**时间**: 2026.01 - 2026.02 (Month 2-3)
**版本**: v0.7.0
**目标**: 彻底重构 UI/UX，打造 AI Native 的交互体验 (参考 NotebookLM)。

*   **Neural Interface (Fluid Canvas)**
    *   [ ] **AI First Layout**: 移除传统邮件列表，首页变为 "Briefing" (智能简报) 和 "Studio" (问答画布)。
    *   [ ] **RAG Deep Integration**: 问答即交互核心，支持引用跳转 (Citations) 和流式输出。
    *   [ ] **Dynamic Widgets**: 根据意图自动渲染组件 (如：日历、关系图谱、审批卡片)。
    *   [ ] **Context Manager**: 左侧栏管理 "关注点" (Contexts) 而非文件夹。

## 🗓️ Phase 7: 团队协作 (Team Collaboration)
**时间**: 2026.02 - 2026.03 (Month 3-4)
**版本**: v0.8.0
**目标**: 完成多用户协作基础，商业化前置准备。

*   **团队协作 (Team Collaboration)**
    *   [ ] **组织架构 (Organization)**: 支持多用户归属同一组织，管理员权限。
    *   [ ] **共享标签/分类 (Shared Labels)**: 团队统一的邮件分类标准。
    *   [ ] **协作批注 (Comments)**: 团队成员在邮件内部进行讨论，不发送给外部。

## 🗓️ Phase 8: 多端体验 (Cross-Platform)
**时间**: 2026.03 - 2026.04 (Month 5-6)
**版本**: v0.9.0
**目标**: 覆盖桌面端与移动端，提供原生体验。

*   **桌面端 (Desktop)**
    *   [ ] **Electron/Flutter App**: 封装 Web 端，支持离线访问与系统通知。
    *   [ ] **快捷键增强**: 类似 Superhuman 的全键盘操作支持。
*   **移动端 (Mobile - WeChat First)**
    *   [ ] **微信公众号集成**: 唯一移动端入口，轻量级交互。
    *   [ ] **语音指令 (Voice AI)**: 发送语音 "帮我查一下未读的紧急邮件"，AI 自动回复摘要。
    *   [ ] **模板消息推送**: 实时推送高优先级邮件通知。

## 🗓️ Phase 9: 商业化 (Commercialization) - 🚀 Launch
**时间**: 2026.05+ (Month 7+)
**版本**: v1.0.0
**目标**: 正式推出付费服务。

*   **商业化 (Monetization)**
    *   [ ] **Stripe 集成**: 订阅支付。
    *   [ ] **多级套餐**: Free, Pro, Team。
    *   [ ] **企业级特性**: SSO, Audit Logs.
