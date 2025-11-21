# 📅 产品路线图 - EchoMind

**当前阶段**：Phase 5: Commercialization (v0.6.0 Planned)
**发布周期**：双周迭代 (Bi-weekly Sprint)

---

## 🗓️ Phase 1: 基础连接与数据概览 (MVP) - ✅ Completed
**目标**：跑通 "连接 -> 同步 -> 展示" 核心链路，验证技术可行性。
**版本**：v0.1.0 (Basic Sync), v0.2.0 (Auth & Multi-tenancy)

*   **后端 (Backend)**
    *   [x] 实现通用 IMAP 连接器 (Go)，支持 SSL/TLS 连接。
    *   [x] 实现邮件元数据同步 (Subject, Sender, Date) 存入 PostgreSQL。
    *   [x] 基础邮件正文解析 (HTML -> Text)。
    *   [x] **用户系统**：注册/登录 (JWT)，多租户数据隔离。
*   **前端 (Frontend)**
    *   [x] 登录/注册页 & 邮箱绑定向导。
    *   [x] 统一收件箱列表视图 (Unified Inbox)。
    *   [x] 邮件详情页 (基础阅读体验)。
*   **AI (Prelim)**
    *   [x] 接入 DeepSeek API。
    *   [x] 实现简单的 "一键摘要" 功能 (On-demand summary)。

## 🗓️ Phase 2: 智能分析与双视图构建 (Alpha) - ✅ Completed
**目标**：实现核心价值功能，AI 驱动的分类与任务提取。
**版本**：v0.3.0

*   **核心功能**
    *   [x] **智能分类器**：自动区分 Newsletter, Notification, Personal, Work。
    *   [x] **任务提取引擎**：从邮件中识别 Action Items 并存入待办库。
    *   [x] **智能仪表盘**：基于分类的过滤视图。
    *   [x] **任务看板**：展示从邮件中提取的待办事项。
*   **技术基建**
    *   [x] 完善 AI 输出结构化 (JSON Output)。
    *   [x] 引入消息队列 (Asynq) 异步处理 AI 任务。

## 🗓️ Phase 3: 真实环境同步与设置 (Beta) - ✅ Completed
**目标**：支持用户自定义 IMAP 配置，完善账户管理。
**版本**：v0.4.0

*   [x] **账户管理**：加密存储 IMAP 凭证。
*   [x] **动态同步引擎**：基于用户配置的动态 IMAP 连接。
*   [x] **设置界面**：前端连接向导与状态展示。

## 🗓️ Phase 4: 深度洞察与关系智能 (Beta 2) - ✅ Completed
**目标**：从单封邮件分析升级为关系网络分析，提供 AI 辅助回复。
**版本**：v0.5.0

*   [x] **联系人智能**：自动聚合联系人互动统计 (次数, 情感均值)。
*   [x] **关系图谱**：可视化展示用户社交/工作网络。
*   [x] **智能回复 (Smart Reply)**：基于上下文生成回复草稿。

## 🗓️ Phase 5: 商业化与团队协作 (RC) - 🚧 In Progress
**目标**：构建 SaaS 基础能力，支持付费订阅与团队管理。
**版本**：v0.6.0 (Target)

*   **商业化 (Monetization)**
    *   [ ] **Stripe 集成**：订阅支付与 Webhook 处理。
    *   [ ] **用量限制**：基于订阅等级的 AI 调用限额。
*   **团队协作 (Team)**
    *   [ ] **组织架构**：多用户归属同一组织。
    *   [ ] **共享收件箱**：团队成员查看特定分类邮件。
