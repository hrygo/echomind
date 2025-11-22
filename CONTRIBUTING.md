# Contributing to EchoMind

Thank you for your interest in contributing to EchoMind! We welcome contributions from everyone.

## 1. Getting Started

1.  Clone the repository.
2.  Install dependencies (`go mod tidy`, `pnpm install`).
3.  Copy `backend/configs/config.example.yaml` to `backend/configs/config.yaml` and configure it.
4.  Run `make docker-up` to start infrastructure.

## 2. Commit & PR Guidelines

*   **Conventional Commits**: We follow [Conventional Commits](https://www.conventionalcommits.org/).
    *   Format: `<type>(<scope>): <subject>`
    *   **Types**: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `perf`, `ci`.
*   **Atomic Commits**: Commit often. Each commit should do one thing and do it well. Avoid "mega-commits".
*   **Pull Requests**: Keep PRs small and focused. One logical feature or bug fix per PR.

## 3. Development Protocols (Strict)

### 🛡️ Quality Assurance
*   **Backend**: `make test` must pass. All DB models must be compatible with both Postgres (Prod) and SQLite (Test).
*   **Frontend**: `pnpm build` AND `pnpm type-check` must pass. Do not rely solely on IDE checks.
*   **CI/CD**: Our CI runs these checks automatically. PRs with failing checks will be blocked.

### 🏗️ Frontend Guidelines
*   **Components**: Before using a generic UI component (e.g., `Dialog`, `Button`), verify it exists in `src/components/ui`. If not, create it first (following shadcn/ui patterns) in a separate commit.
*   **State Management**: Use Zustand. Be careful with `persist` middleware typing in v5.
*   **Imports**: Avoid default exports for core libraries if possible to simplify refactoring.

### 🔧 Refactoring Safety
*   **Blast Radius**: When changing a core type (e.g., DB Model) or API (e.g., `apiClient`), assume it breaks everything. Run a project-wide search (`grep` or IDE) to find usages.
*   **Incremental**: Deprecate old APIs before removing them if the codebase is large.

## 4. Versioning & Release Strategy

*   We use **Semantic Versioning (SemVer)**: `vMajor.Minor.Patch`.
*   **Pre-Release Checklist**: Before creating a git tag, you MUST update the version number in:
    1.  `Makefile` (VERSION variable)
    2.  `backend/cmd/main.go` (Version constant)
    3.  `frontend/package.json` (version field)
*   **Process**:
    1.  Update files.
    2.  Commit: `chore: bump version to vX.Y.Z`.
    3.  Tag: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`.
    4.  Push tags: `git push origin vX.Y.Z`.

## 5. Coding Standards

*   **Go**: Follow `Effective Go`. Use `testify` for assertions.
*   **Frontend**: React Functional Components, Hooks, TypeScript (Strict).
*   **i18n**: All user-facing text must be localized (`t('key')`).