# 🛡️ EchoMind Project Context

> **Vision**: AI-powered Email Decision System. "External Brain for Executive Communication."  
> **Stage**: v0.6.4 (Alpha)  
> **Active Sprint**: Phase 5.3 - RAG Polish & Phase 6 Prep

---

## 1. Technology Stack

### Backend
- **Language**: Go 1.22+
- **Web Framework**: Gin
- **Database**: GORM (Postgres + pgvector)
- **Async Jobs**: Asynq (Redis)
- **Configuration**: Viper

### Frontend
- **Framework**: Next.js 16 (TypeScript)
- **UI Library**: Tailwind CSS + Radix UI
- **State Management**: Zustand
- **HTTP Client**: Axios

---

## 2. Development Workflow

### 🚀 Start Development
```bash
make dev    # Starts Postgres, Redis, Backend, Worker, and Frontend
make logs   # View all logs
```

### 🧪 Verify Quality
```bash
make test                  # Backend unit/integration tests
cd frontend && pnpm test   # Frontend component tests
```

### 📦 Release Process
1. Finish Feature/Fix
2. Verify Tests: `make test`
3. Bump Version: `Makefile`, `frontend/package.json`, `backend/cmd/main.go`
4. Commit: `feat: ...` or `fix: ...`
5. Tag: `git tag vX.Y.Z`

---

## 3. Roadmap Status

### ✅ Completed (v0.1.0 → v0.6.4)
- **Core Features**: IMAP Sync, Email Parsing, JWT Authentication
- **AI Capabilities**: Summary, Sentiment, Classification, Contact Intelligence, Smart Reply
- **Optimization**: Spam filtering to reduce AI usage
- **User Interface**: Dashboard, Insights Graph, Account Settings
- **RAG & Search**: Vector embeddings, semantic search, search UI (Phase 5.2)

### 🚧 In Progress (Phase 5.3: v0.6.5 → v0.7.0)
**Plan**: [docs/sprints/week2_rag_polish/sprint-plan.md](docs/sprints/week2_rag_polish/sprint-plan.md)

- **Performance**: Search optimization (< 500ms), monitoring, benchmarking
- **Testing**: Integration tests, E2E tests, 80% coverage target
- **UX Improvements**: Search history, filters, error handling
- **Phase 6 Preparation**: Team Collaboration design, database schema planning

### 🔮 Future Roadmap (6-Month Plan)
- **Phase 6** (2026.01-02): Team Collaboration (Shared Labels, Organization)
- **Phase 7** (2026.03-04): Cross-Platform (Desktop, WeChat Voice AI)
- **Phase 8** (2026.05+): Commercialization (Stripe, SSO)

---

## 4. The Golden Rules (Non-Negotiable)

1. **Frequent Delivery**  
   Commit often. Don't hoard changes.

2. **Semantic Versioning**  
   `vMajor.Minor.Patch`. Update `Makefile` & `package.json` before tagging.

3. **Test-Driven Development**  
   Red → Green → Refactor. `make test` MUST pass before commit.

4. **Convention over Configuration**  
   Follow existing directory structure (`internal/`, `pkg/`) and naming conventions.

5. **Monorepo Discipline**  
   Respect the `backend/` vs `frontend/` boundary.

---

**Last Updated**: 2025-11-22  
**Project**: EchoMind v0.6.4 (Alpha)
