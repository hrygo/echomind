# 🏗️ 技术架构详细规范 (Technical Architecture Spec)

## 1. 工程结构 (Monorepo)

采用 Monorepo 模式，严格遵循 GitHub 开源最佳实践。

```text
/echomind
├── .github/                 # [Best Practice] GitHub 自动化与协作配置
│   └── workflows/           # GitHub Actions CI/CD (CI: Test/Lint, CD: Docker Build/Push)
├── docs/                    # [Best Practice] 项目文档中心
│   ├── prd.md               # 产品需求文档
│   ├── tech-architecture.md # 技术架构文档
│   └── ...
├── backend/                 # Go 后端服务
│   ├── cmd/                 # 应用程序入口
│   │   ├── main.go          # HTTP Server
│   │   └── worker/          # Async Task Worker
│   ├── configs/             # 配置文件模板 (config.example.yaml)
│   ├── internal/            # 私有业务逻辑
│   │   ├── handler/         # HTTP Handlers (Gin)
│   │   ├── model/           # Database Models (GORM)
│   │   ├── service/         # Business Logic & Factory
│   │   └── tasks/           # Asynq Task Handlers
│   ├── pkg/                 # 公共库
│   │   ├── ai/              # AI Providers (OpenAI, Gemini, DeepSeek)
│   │   └── imap/            # IMAP Connector & Body Parser
│   └── go.mod
├── frontend/                # Next.js 前端应用
│   ├── src/
│   │   ├── app/             # Next.js App Router
│   │   ├── components/      # UI 组件库
│   │   └── hooks/           # Custom React Hooks
│   └── package.json
├── deploy/                  # 部署与基础设施
│   ├── deploy.sh            # 生产部署脚本
│   ├── docker-compose.yml   # 本地开发编排
│   └── docker-compose.prod.yml # 生产环境编排
├── scripts/                 # 工具脚本
├── .gitignore               # 全局忽略文件
├── Makefile                 # [Best Practice] 统一任务入口
├── README.md                # 项目主页
└── CONTRIBUTING.md          # [Best Practice] 贡献与开发规约
```

## 2. 后端技术栈 (Go Ecosystem)

*   **Web Framework**: `Gin`
*   **ORM**: `GORM` (PostgreSQL)
*   **Config**: `Viper` (Supports YAML & Environment Variables)
*   **Async Queue**: `Asynq` (Redis-based) - Used for background email analysis tasks.
*   **WeChat Gateway**: Handles WeChat XML callbacks, signature verification, and voice processing via **OpenAI Whisper** (High-accuracy STT).
*   **Spam Filter**: Rule-based filter (`internal/spam`) to pre-screen emails before AI processing.
*   **AI Engine**: 
    *   **Architecture**: Adapter Pattern & Factory Pattern.
    *   **Interface**: `pkg/ai/AIProvider` (Methods: `Summarize`, `Classify`, `AnalyzeSentiment`).
    *   **Implementations**: 
        *   `openai`: Uses `go-openai` SDK.
        *   `gemini`: Uses `generative-ai-go` SDK.
        *   `deepseek`: Adapts `openai` implementation with custom BaseURL.
    *   **RAG Support (v0.6.0+)**:
        *   **Embeddings**: OpenAI `text-embedding-3-small` or compatible.
        *   **Vector DB**: **pgvector** (Postgres extension) for storing email embeddings (No external vector DB required).
    *   **Configuration**: Prompts are externalized in `config.yaml`.
*   **Logging**: `Zap` (Structured Logging)

## 3. 自动化工作流 (CI/CD)

*   **CI (GitHub Actions)**:
    *   **Backend**: Go Mod Tidy, Test (`go test ./...`), Lint.
    *   **Frontend**: Pnpm Lint, Test (`jest`).
*   **CD (GitHub Actions)**:
    *   **Docker**: Multi-stage builds for Backend and Frontend.
    *   **Registry**: Push images to GitHub Container Registry (GHCR) on `main` branch push.
    *   **Deploy**: Shell script triggering Docker Compose update.

## 4. 开发规范

*   **Workflow**: TDD (Test-Driven Development).
*   **Version Control**:
    *   **Commits**: Frequent, Atomic, Conventional Commits (`feat:`, `fix:`, `chore:`).
    *   **Releases**: Independent feature releases tagged with Semantic Versioning (`vX.Y.Z`).

## 5. 数据层规范 (Data Layer)

### PostgreSQL Schema
*   **Naming**: snake_case.
*   **ID**: UUID/Snowflake.
*   **Core Entities**:
    *   `emails`: Stores email metadata, content snippet, and AI insights (`Summary`, `Sentiment`, `Urgency`).
    *   `contacts`: Stores sender info (`Name`, `Email`) and interaction stats (`InteractionCount`, `LastInteractedAt`).
    *   `users` (Planned Phase 4): Auth & Tenant isolation.

### Redis Keys
*   `asynq:{queue}`: Background task queues.
*   `echomind:cache:{key}`: General caching.

## 6. 可观测性 (Observability)

系统集成了结构化日志与基础监控指标，确保生产环境的可见性。

*   **Structured Logging**: 使用 Uber `Zap` 库。
    *   **Request IDs**: 每个 HTTP 请求分配唯一 `X-Request-ID`，贯穿处理链路。
    *   **Levels**: `Info` (常规操作), `Warn` (业务异常), `Error` (系统故障/Panic).
    *   **Fields**: Log entries include `userID`, `duration`, `query` (for search), etc.
*   **Health Checks**:
    *   `GET /api/v1/health`: 检查数据库 (Postgres + pgvector) 连接状态。
*   **Metrics (Logs)**:
    *   Search Latency: 记录每次搜索的耗时。
    *   Embedding Latency: 记录向量生成的耗时（外部 API 调用）。

## 7. API 接口 (API)

详细 API 文档请参考: [docs/api.md](api.md)

## 8. Multi-Tenancy Architecture (v0.7.0+)

EchoMind 支持多租户架构，允许用户创建和管理组织（Organization）和团队（Team）。

*   **Models**:
    *   `Organization`: 最高层级单元，拥有资源和成员。
    *   `OrganizationMember`: 关联用户与组织，包含角色（Owner, Admin, Member）。
    *   `Team`: 组织内的子组，可拥有特定的资源（如共享邮箱）。
    *   `TeamMember`: 关联用户与团队。
*   **Resource Ownership**:
    *   资源（如 `EmailAccount`, `Contact`）现在支持三种所有权模式：
        1.  **Personal**: `UserID` 不为空，`TeamID`/`OrgID` 为空。
        2.  **Organization**: `OrganizationID` 不为空，`TeamID`/`UserID` 为空。
        3.  **Team**: `TeamID` 不为空，`UserID` 为空。
*   **Context Switching**:
    *   API 请求通过 Header `X-Organization-ID` 传递当前上下文。
    *   前端使用 Zustand store 管理当前选中的组织。

