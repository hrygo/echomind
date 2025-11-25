# 🛡️ EchoMind Project Context

- **Vision**: Personal Neural Interface (个人智能神经中枢)
- **Status**: `v0.9.7` | **Current Sprint**: Dashboard API Integration Phase 1
- **Stack**:
    - **Backend**: Go 1.22+ (Gin, GORM, Asynq) | Postgres+pgvector | Redis
    - **Frontend**: Next.js 16 (TypeScript, Tailwind, Zustand)

---
## Roadmap

- ✅ **v0.9.2-4 (Neural Nexus)**: Context Bridge, Omni-Bar, Generative Widget Framework.
- 🚧 **v0.9.5+ (WeChat OS)**: Voice Commander, One-Touch Decisions, Calendar Gatekeeper, Morning Briefing.

---
## The Golden Rules (Non-Negotiable)

- 🛡️ **Quality (TDD)**
    - **CI**: `make test` (BE) & `make test-fe` (FE) must pass before commit.
    - **Tests**: Mock external dependencies (AI, DB) for speed & stability.
    - **Build**: Use `make build` (BE) & `make build-fe` (FE) for compilation verification.

- 🚀 **Delivery (Frequent & Versioned)**
    - **Commits**: Atomic, frequent, use conventional prefixes (`feat:`, `fix:`).
    - **Versioning**: Tag releases often. Update version in `Makefile`, `package.json`, `backend/cmd/main.go`, `README.md`, `README-zh.md`, `docs/openapi.yaml`, `docs/*.md`.

- 🏗️ **Architecture & Code**
    - **Refactor**: Use `grep` to find all usages. Keep old APIs temporarily for core changes.
    - **Frontend**: Check for existing components (`src/components/ui`) before creating new ones.
    - **Database**: Compile BE after model changes. Avoid DB-specific defaults (e.g., `gen_random_uuid()`) in GORM tags.

- 🌐 **Internationalization (i18n)**
    - All UI text must be bilingual (en/zh) via `t('key')`. No hardcoded strings.

- 🔧 **Tooling (AI Agent SOP)**
    - **Working Directory**: Always ensure `~/aicoding/echomind` as the base directory before executing any commands.
    - **Preferred Commands**: Prioritize `make` commands over direct tool calls:
      - Use `make test` instead of `go test ./...`
      - Use `make build` instead of `go build ./cmd/main.go`
      - Use `make test-fe` instead of `pnpm build && pnpm type-check`
      - Use `make build-fe` instead of `pnpm build`
    - **`replace`**: Use minimal, unique context for `old_string`.
    - **`verify`**: On tool failure, use `read_file` to check state before retrying.
    - **Directory Awareness**: Before any command execution, verify working directory with `pwd` and navigate to project root if needed.
