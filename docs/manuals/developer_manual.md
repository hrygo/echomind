# 🛠️ EchoMind 研发手册 (Developer Manual)

> **架构**: Cloud-Native Microservices (Monolithic Repo)
> **技术栈**: Go + Next.js + Postgres + Redis + LLM

## 1. 技术架构 (Architecture)

### 1.1 后端 (Backend)
基于 **Go 1.22+** 构建，遵循 Clean Architecture（整洁架构）原则。
*   **目录结构**:
    *   `cmd/`: 应用程序入口 (`server`, `worker`, `tools`)。均使用 `internal/bootstrap` 进行统一初始化。
    *   `internal/model/`: GORM 数据模型定义。
    *   `internal/service/`: 核心业务逻辑 (Email, Context, Action, Search)。
    *   `internal/handler/`: HTTP 路由处理层 (Gin)。
    *   `internal/tasks/`: Asynq 异步任务定义。
    *   `pkg/`: 通用工具库 (AI Provider, IMAP, Utils)。
*   **关键组件**:
    *   **Web Server (Gin)**: 提供 RESTful API。
    *   **Worker (Asynq)**: 处理邮件同步、LLM 分析、向量生成等耗时任务。
    *   **Database (Postgres)**: 使用 `pgvector` 扩展存储文本向量。

### 1.2 前端 (Frontend)
基于 **Next.js 16 (App Router)** 和 **TypeScript**。
*   **状态管理**: Zustand (`useEmailStore`, `useActionStore`, `useContextStore`, `useCopilotStore` - v0.9.3 新增)。
*   **UI 框架**: Tailwind CSS + Lucide React Icons。
*   **国际化**: React Context (`LanguageContext`) 加载 JSON 字典。
*   **核心组件**: `CopilotWidget` (v0.9.3 新增) 统一了搜索和 AI 对话入口。
*   **交互模式**:
    *   **Optimistic UI**: 用户点击操作（如批准）时，界面立即响应，后台异步请求。
    *   **Undo Mechanism**: 关键操作提供 Toast 撤销功能。

## 2. 数据模型设计 (ERD 核心)
*   **Users**: 用户基础信息。
*   **Emails**: 核心存储。
    *   包含 `Summary` (AI摘要), `Sentiment` (情感), `Urgency` (紧急度)。
    *   `SnoozedUntil`: 标记小睡截止时间。
*   **EmailEmbeddings**: 存储邮件正文的 Vector Embedding (1536维)，用于 RAG 搜索。
*   **Contexts**: 智能情境定义（Keywords, Stakeholders JSON）。
*   **EmailContexts**: 多对多关联表，记录邮件命中了哪些情境。
*   **Tasks**: 从邮件中提取的待办事项。

## 3. 核心业务流程

### 3.1 邮件同步与分析流水线 (The Pipeline)
1.  **Sync**: 用户触发或定时触发 `SyncService`，通过 IMAP 拉取新邮件。
2.  **Queue**: 新邮件 ID 被推送到 Redis 队列 `email:analyze`。
3.  **Worker Processing**:
    *   **Spam Check**: 规则过滤垃圾邮件。
    *   **LLM Analysis**: 调用 OpenAI/Gemini 生成摘要、情感评分、提取 Action Items。
    *   **Context Matching**: 遍历用户定义的 Context 规则，进行关键词和发件人匹配，打标签。
    *   **Embedding**: 调用 Embedding API 生成向量并存入 `pgvector`。

### 3.2 搜索与 AI 对话 (RAG Search & AI Chat)
1.  用户通过 **智能副驾（Omni-Bar）**输入查询或提问。
2.  **搜索模式**: 如果是关键词搜索，后端将查询语句转化为 Vector，数据库进行余弦相似度搜索 (`<=>` 运算符)，找出语义最接近的邮件。结合传统 SQL 过滤（如时间范围、Context ID）返回结果。
3.  **AI 对话模式**: 如果是提问，前端会将对话历史和当前显示的搜索结果（作为 `context_ref_ids`）发送给后端。
4.  **后端处理**: `ChatService` 优先从 `context_ref_ids` 获取邮件内容注入系统 Prompt，然后调用 AI 模型进行对话，并通过 SSE 流式返回响应。

### 3.3 极速行动 (Actions)
*   **API**: `POST /api/v1/actions/{type}`
*   **Approve**: 执行软删除 (Soft Delete)，从收件箱视图移除。
*   **Snooze**: 设置 `snoozed_until` 字段。`ListEmails` 接口默认过滤掉 `snoozed_until > NOW()` 的记录。
*   **Dismiss**: 将邮件的 `Urgency` 字段强制降级为 `Low`。

## 4. 开发与部署指南

### 4.1 环境依赖
*   Go 1.22+
*   Node.js 20+ (pnpm)
*   Docker Compose (Postgres 15+ with pgvector, Redis 7)

### 4.2 常用命令 (Makefile)
EchoMind 使用 Makefile 管理全生命周期：
*   `make init`: 初始化 Go mod 和 pnpm 依赖。
*   `make docker-up`: 启动数据库和 Redis。
*   `make dev`: 一键启动所有服务（Backend, Worker, Frontend, DB）。
*   `make build`: 编译后端二进制文件。
*   `make test`: 运行后端单元测试。
*   `make reindex`: 手动重新生成所有邮件的向量索引。

### 4.3 配置管理
配置文件位于 `backend/configs/config.yaml`。
*   **敏感信息**: 生产环境建议通过环境变量覆盖，例如 `ECHOMIND_AI_OPENAI_API_KEY`。
*   **引导**: `backend/internal/bootstrap` 包负责加载配置并初始化全局单例。

### 4.4 贡献代码
1.  **Frontend**: 修改组件后运行 `pnpm type-check` 确保类型安全。
2.  **Backend**: 修改接口需同步更新 `internal/handler` 和 `internal/router/routes.go`。
3.  **Tests**: 新增 Service 逻辑必须编写对应的 `_test.go`。
