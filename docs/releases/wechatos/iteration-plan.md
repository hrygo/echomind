# 🗓️ WeChat OS 迭代计划

> **总周期**: 2026.01 - 2026.02 (8 Weeks) | **负责人**: TBD

## Phase 7.1: 基础接入与绑定 (Weeks 1-2)
**目标**: 打通微信与 EchoMind 的连接，实现账号绑定和基础文本对话。

*   **Week 1: 基础设施搭建**
    *   [ ] 申请微信测试号/服务号。
    *   [ ] 配置公网回调域名 (Ngrok/Cloudflare)。
    *   [*] **SDK 集成**: 引入 `github.com/silenceper/wechat/v2`。
        *   配置 `RedisCache` 以复用现有 Redis 实例存储 `access_token`。
        *   初始化 `OfficialAccount` 实例。
    *   [ ] 实现 `WeChat Gateway` 基础接收与验签逻辑。
    *   [ ] 数据库迁移: 添加 `users` 表微信字段。
*   **Week 2: 账号绑定流程**
    *   [ ] 后端: 实现生成带参二维码接口。
    *   [ ] 后端: 实现扫码回调逻辑 (Event: `SCAN` / `SUBSCRIBE`)。
    *   [ ] 前端: 开发“设置 -> 微信绑定”页面。
    *   [ ] 测试: 验证绑定、解绑流程。

## Phase 7.2: 语音指挥官与意图识别 (Weeks 3-5)
**目标**: 实现语音转文字，并能通过自然语言查询邮件。

*   **Week 3: 语音处理管道**
    *   [ ] 集成微信素材下载接口。
    *   [ ] 集成 Whisper API (或使用微信自带识别)。
    *   [ ] 实现 `Voice -> Text` 转换模块。
*   **Week 4: 意图识别 (LLM Integration)**
    *   [ ] 设计 Prompt Template 和 Tools Definition。
    *   [ ] 实现 `Intent Analyzer` 服务。
    *   [ ] 开发 `EmailService` 的 Tool 适配层 (Search, GetSummary)。
*   **Week 5: 多轮对话与 FSM**
    *   [ ] 引入 Redis FSM 存储会话状态。
    *   [ ] 实现上下文管理 (Context Management)。
    *   [ ] 联调: 语音查询邮件 -> LLM 搜索 -> 返回结果。

## Phase 7.3: 主动智能与推送 (Weeks 6-7)
**目标**: 实现晨间简报和一键决策推送。

*   **Week 6: 晨间简报 (Morning Briefing)**
    *   [ ] 开发数据聚合服务 (Calendar + Email + Task)。
    *   [ ] 开发简报生成 Prompt。
    *   [ ] 实现 Cron Job 定时推送逻辑。
*   **Week 7: 一键决策 (One-Touch Decision)**
    *   [ ] 优化 `EmailIngestor`，识别决策类邮件。
    *   [ ] 设计微信模板消息 UI。
    *   [ ] 实现点击回调处理 (Approve/Reject Action)。

## Phase 7.4: 验收与发布 (Week 8)
**目标**: 整体测试、Bug 修复与文档完善。

*   **Week 8: Polish & Launch**
    *   [ ] **压力测试**: 模拟高并发消息处理。
    *   [ ] **安全审计**: 检查 OpenID 绑定安全性和数据脱敏。
    *   [ ] **用户文档**: 编写微信使用指南。
    *   [ ] **发布**: 部署到生产环境。

## 里程碑 (Milestones)

| 里程碑               | 时间点     | 交付物                             |
| :------------------- | :--------- | :--------------------------------- |
| **M1: Connectivity** | Week 2 End | 微信公众号可关注，扫码可绑定账号。 |
| **M2: Voice Query**  | Week 5 End | 可以发语音查邮件，系统能准确回复。 |
| **M3: Proactive AI** | Week 7 End | 每天早上收到简报，能处理审批推送。 |
| **M4: GA Release**   | Week 8 End | 全量发布 v0.9.5+。                 |
