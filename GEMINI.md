# 🛡️ EchoMind Project Specification

**Vision**: Personal Neural Interface | **Version**: v1.1.0 (Enterprise Release)
**Tech Stack**: Go(Gin/GORM/Asynq) + Next.js + Postgres(pgvector) + Redis

---

## 🚀 Version Release

**Checklist**: `frontend/package.json` | `backend/pkg/logger/config.go` | `Makefile` | `docs/openapi.yaml` | `backend/configs/logger*.yaml` | `README*.md` | `CHANGELOG.md` | `docs/product-roadmap.md` | `docs/logger/README.md`

**Release**: `git add . && git commit -m "feat: v{version} - description" && git tag -a v{version}`

**Strategy**: Semantic v{MAJOR}.{MINOR}.{PATCH} | Enterprise v1.0+ | Sync all references

---

## ⚡ Development Rules

**Quality**: Pre-commit `make test build test-fe build-fe` | Test-first (mock AI/DB) | Build verification

**Architecture**: DB compile verification | Frontend `src/components/ui` reuse | `grep` refactoring | Bilingual `t('key')`

**Workflow**: TDD + Make-First (Tests → Implementation → Makefile → Verification)

**Commands**: `make test/test-fe/build/build-fe/run-backend/run-worker/stop/db-init/lint`

**Pre-commit**: `make test && make build && make test-fe && make build-fe`

---

## 📋 AI Agent Standards

**Environment**: Directory `~/aicoding/echomind` | Make-First priority | Verify working directory

**Operations**: TDD cycle (red-green-refactor) | Minimal context | Commit prefixes `feat: fix: docs: refactor:` | Atomic commits

**Principles**: Exclusive Make commands | Full test coverage | `read_file` verification | `grep` before refactoring | Gradual API migration