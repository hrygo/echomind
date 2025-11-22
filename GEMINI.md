# 🛡️ EchoMind Project Context

> **Vision**: Personal Neural Interface for Executive Work. (个人智能神经中枢)  
> **Stage**: v0.6.4 (Alpha) | **Active Sprint**: Phase 5.3 - RAG Polish & Phase 6 Prep

---

## 1. Technology Stack

**Backend**: Go 1.22+ (Gin, GORM, Asynq) | Postgres + pgvector | Redis  
**Frontend**: Next.js 16 (TypeScript, Tailwind CSS, Zustand)

---

## 2. Roadmap Status

### ✅ Recent Completion (v0.6.5 → v0.7.0)
**Phase 5.3 - RAG Polish & Phase 6 Prep**
- ✅ Performance: < 500ms search
- ✅ Design: Team Collaboration architecture

### 🚧 Current Sprint (v0.7.1+)
**Phase 6.1 - Team Collaboration Foundation** | [Plan](docs/sprints/week3_team_collab/sprint-plan.md)
- Database Models (Org, Team, Member)
- Multi-tenant Migration
- Organization API & UI Switcher

### 🔮 Future (6-Month Plan)
- **Phase 6** (2026.01-02): Team Collaboration
- **Phase 7** (2026.03-04): Cross-Platform (Desktop, WeChat)
- **Phase 8** (2026.05+): Commercialization (Stripe, SSO)

---

## 3. The Golden Rules (Non-Negotiable)

### Quality (Test-Driven)
`make test` MUST pass before commit. Red → Green → Refactor.

### Delivery (Frequent & Versioned)
1. **Commit Often**: Don't hoard changes
2. **Tag Immediately**: After quality verification, release
   - Minor features (v0.x.Y): Daily if tests pass
   - Major features (v0.X.0): Milestone complete
   - Fixes (v0.x.y): Immediate
3. **Convention**: `feat:` | `fix:` | `docs:`
4. **Versioning**: Update `Makefile`, `package.json`, `backend/cmd/main.go`
5. **Principle**: "Done" = "Released & Tagged"

### Structure (Convention over Configuration)
- Follow `internal/`, `pkg/` directory structure
- Respect `backend/` vs `frontend/` boundary
- Maintain naming conventions

### Internationalization (i18n)
- **Bilingual UI**: All user-facing text MUST support both English (en) and Chinese (zh).
- **No Hardcoding**: Never hardcode strings in components. Use `t('key')` and update locale files.
