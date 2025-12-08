# 🧠 EchoMind Project Identity

**Vision**: Personal Neural Interface for Email Intelligence
**Version**: v1.2.1 (Desktop Excellence)
**Tech Stack**: Go (Gin/Asynq) | Next.js | Postgres (pgvector) | Redis | Docker

---

## 🗺️ Essence: Cognitive Map (认知地图)

> **What we are building**: A local-first, privacy-focused AI Chief of Staff. It ingests emails, generates knowledge graphs, and provides actionable intelligence via Chat & Dashboard.

*   **Architecture**: Modular Monolith (Clean Architecture).
    *   `backend/internal`: Domain Logic (Models, Services).
    *   `backend/pkg`: Infrastructure & Utilities.
    *   `frontend/src`: Feature-based Component Structure.
*   **Key Paths**:
    *   📖 **Docs**: `docs/` (Starting point: `docs/README.md`)
    *   🔌 **API**: `docs/api/openapi.yaml` (OpenAPI v3)
    *   🗄️ **Schema**: `backend/internal/model` & `backend/migrations`
    *   🚀 **Deploy**: `deploy/docker-compose.yaml`

---

## 📍 Context: Current Focus (当前焦点)

*   **Status**: **Post-Release (v1.2.1)**. Focusing on "Desktop Excellence".
*   **Priority**:
    1.  **Smart Release**: Stability of the new `/release` workflow.
    2.  **Documentation**: Intelligent archival and organization.
    3.  **Next Iteration**: Planning features for v1.2.0 (See `docs/product/product-roadmap.md`).

---

## ⚡ Operations: Smart Protocols (操作协议)

### 🤖 Intelligent Workflows (AI-first)
*   **`/release`**: **Full Release Cycle**. (Lint -> Test -> Bump -> Archive -> Tag).
*   **`/git`**: Smart Git Commit & Push (Auto-message).

### 🛠️ Engineering Defaults
*   **Command Hub**: Always use `Makefile`!
    *   `make dev`: Start everything.
    *   `make lint` / `make test`: QA suite.
    *   `make db-init`: Reset database.
*   **Verification**: NEVER assume. Always verify with `grep` or `ls` after changes.

---

## 🛡️ Heuristics: Prime Directives (核心法则)

1.  **Make-First Principle**: Do not run raw `go run` or `npm` commands if a `make` target exists.
2.  **Atomic Evolution**: Changes must be small, testable, and reversible.
3.  **Documentation Continuity**: If you change logic, you MUST update the corresponding doc in `docs/`.
4.  **Bilingual Cognition**: Code comments in English; User interaction/Docs in Chinese (unless specified). **ALL Reports must be in Simplified Chinese**.
5.  **Safety**: Never commit secrets. Validate all content before file writes.