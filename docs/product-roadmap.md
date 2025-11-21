# 📅 产品路线图 - EchoMind

**当前阶段**：Phase 2: Intelligent Analysis & Insights (v0.3.0 In Progress)
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
    *   [ ] **每日日报 (Daily Digest)**：每天早上 8:00 生成简报邮件/推送。
    *   [ ] **微信公众号接入 (Beta)**：
        *   [ ] 实现微信服务器验签与基础消息回复。
        *   [ ] 实现账户绑定流程 (OpenID <-> UserID)。
*   **用户界面**
    *   [x] **智能仪表盘**：基于分类的过滤视图。
    *   [x] **任务看板**：展示从邮件中提取的待办事项。
*   **技术基建**
    *   [x] 完善 AI 输出结构化 (JSON Output)。
    *   [x] 引入消息队列 (Redis/Kafka) 异步处理 AI 任务 (已引入 Asynq)。

## 🗓️ Phase 3: 关系网络与深度整合 (Beta) - 🚧 In Progress
**目标**：增强粘性，提供深层次洞察，拓宽用户群体至销售/投资人。
**预计周期**：8 - 10 周
**版本**：v0.4.0 (Target)
