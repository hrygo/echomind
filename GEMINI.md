# 🛡️ EchoMind Project Specification

**Vision**: Personal Neural Interface | **Version**: v1.1.0 (Enterprise Release)
**Tech Stack**: Go(Gin/GORM/Asynq) + Next.js + Postgres(pgvector) + Redis

---

## 🚀 Version Release Specification

### Version Checklist
- `frontend/package.json` | `backend/pkg/logger/config.go` | `Makefile`
- `docs/openapi.yaml` | `backend/configs/logger*.yaml`
- `README*.md` | `CHANGELOG.md` | `docs/product-roadmap.md` | `docs/logger/README.md`

### Release Process
```bash
git add . && git commit -m "feat: v{version} - description" && git tag -a v{version}
```

### Version Strategy
- **Semantic**: `v{MAJOR}.{MINOR}.{PATCH}`
- **Enterprise**: v1.0+ marks production readiness
- **Sync**: All version references stay consistent

---

## ⚡ Core Development Rules

### Quality Assurance
- **Pre-commit**: `make test build test-fe build-fe`
- **Test First**: Mock external dependencies (AI, DB)
- **Build Verification**: Ensure compilation succeeds

### Architecture Principles
- **Database**: Compile verification after GORM model changes
- **Frontend**: Prioritize `src/components/ui` component reuse
- **Refactoring**: `grep` global search, preserve old APIs during transition
- **Internationalization**: Mandatory bilingual `t('key')`

### Tool Usage
```bash
make test build run-backend stop  # Core development commands
```

---

## 📋 AI Agent Operating Standards

### Working Environment
- **Directory**: Must be `~/aicoding/echomind`
- **Verification**: Confirm working directory before command execution

### Development Operations
- **File Operations**: Minimize context, prioritize state checks on failure
- **Commit Standards**: `feat:` `fix:` `docs:` `refactor:` prefixes
- **Atomic Commits**: Frequent, small-granularity commits

### Guiding Principles
- **Make Priority**: Use unified Make command interfaces
- **State Verification**: Use `read_file` on operation failure
- **Global Search**: Use `grep` to find all references before refactoring
- **Progressive**: Preserve old APIs, migrate gradually