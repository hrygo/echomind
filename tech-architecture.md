# 🏗️ 技术架构详细规范 (Technical Architecture Spec) - GitHub Optimized

## 1. 工程结构 (Monorepo)

采用 Monorepo 模式，严格遵循 GitHub 开源最佳实践。

```text
/echomind-root
├── .github/                 # [Best Practice] GitHub 自动化与协作配置
│   ├── ISSUE_TEMPLATE/      # Issue 规范模板 (Bug report, Feature request)
│   ├── workflows/           # GitHub Actions CI/CD (Go Test, Lint, Build)
│   └── PULL_REQUEST_TEMPLATE.md
├── docs/                    # [Best Practice] 项目文档中心
│   ├── architecture/        # 架构设计文档 (md + images)
│   └── api/                 # API 定义 (OpenAPI/Swagger)
├── backend/                 # Go 后端服务 (遵循 golang-standards)
│   ├── api/                 # API 协议定义 (Proto/OpenAPI)
│   ├── cmd/                 # 应用程序入口 (Main applications)
│   │   ├── server/          # -> main.go (HTTP Server)
│   │   └── worker/          # -> main.go (Async Task Worker)
│   ├── configs/             # 配置文件模板 (config.example.yaml)
│   ├── internal/            # 私有业务逻辑 (Private application code)
│   │   ├── handler/         # HTTP Handlers (Gin) -> sync.go (/api/v1/sync)
│   │   ├── model/           # Database Models (GORM) -> email.go
│   │   ├── service/         # Business Logic -> sync.go (SyncEmails)
│   │   ├── repository/      # Data Access Layer
│   │   └── middleware/      # HTTP Middlewares
│   ├── pkg/                 # 公共库 (可被外部引用的代码，如 Utils, SDK)
│   │   └── imap/            # IMAP Connector, Fetcher & Body Parser (connector.go, fetch.go, body.go)
│   └── go.mod
├── frontend/                # Next.js 前端应用
│   ├── src/
│   │   ├── app/             # Next.js App Router
│   │   ├── components/      # UI 组件库
│   │   └── hooks/           # Custom React Hooks
│   └── package.json
├── deploy/                  # 部署与基础设施
│   ├── docker/              # Dockerfiles
│   └── docker-compose.yml   # 本地开发编排
├── scripts/                 # 构建与维护脚本 (Shell/Python)
├── .editorconfig            # [Best Practice] 跨编辑器代码风格统一
├── .gitignore               # 全局忽略文件
├── Makefile                 # [Best Practice] 统一任务入口 (make run, make build)
├── LICENSE                  # 开源协议
├── README.md                # 项目主页 (Badges, Quick Start)
└── CONTRIBUTING.md          # [Best Practice] 贡献指南
```

## 2. 后端技术栈 (Go Ecosystem)

*   **Web Framework**: `Gin`
*   **ORM**: `GORM` (PostgreSQL)
*   **Config**: `Viper`
*   **Async Queue**: `Asynq` (Redis-based) - Used for background email analysis.
*   **AI Engine**: Strategy Pattern for LLM Providers (DeepSeek, OpenAI, etc.), Config-driven.
    *   Interface: `pkg/ai/provider.go`
    *   Implementations: `pkg/ai/deepseek`, `pkg/ai/openai`
*   **Logging**: `Zap` (Structured Logging)
*   **Linting**: `golangci-lint` (集成在 CI 中)

## 3. 自动化工作流 (CI/CD)

*   **Pre-commit**: 本地检查格式 (gofmt, prettier)。
*   **CI (GitHub Actions)**:
    *   `go-test`: 每次 Push 自动运行 Go 单元测试。
    *   `go-lint`: 运行 golangci-lint 检查代码质量。
    *   `frontend-build`: 检查 Next.js 构建是否通过。

## 4. 开发规范

*   **Commit Message**: 遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范 (e.g., `feat: add email sync`, `fix: task status update`).
*   **Branching**: Feature Branch Workflow (`main` <- `feature/xyz`).

## 5. 数据层规范 (Data Layer)

### PostgreSQL Schema
*   **命名**: snake_case。
*   **ID**: UUID/Snowflake。
*   **Entities**:
    *   `emails`: Stores email content, metadata, and AI analysis results (summary, sentiment, urgency).
    *   `contacts`: Stores sender information and interaction stats (count, last_interacted).
*   **Migration**: 使用 GORM AutoMigrate (开发阶段) 或 Golang-Migrate (生产阶段)。

### Redis Keys
*   `echomind:sess:{token}`
*   `echomind:queue:{task_id}`