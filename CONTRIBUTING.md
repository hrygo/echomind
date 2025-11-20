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
    *   **Types**:
        *   `feat`: New feature
        *   `fix`: Bug fix
        *   `docs`: Documentation changes
        *   `refactor`: Code changes without fixing bugs or adding features
        *   `test`: Adding missing tests
        *   `chore`: Maintainance tasks (build, ci, deps)
*   **Atomic Commits**: Commit often. Each commit should do one thing and do it well. Avoid "mega-commits".
*   **Independent Release**: Once a feature module (and its tests) is complete, it should be considered ready for release, adhering to Semantic Versioning.
*   **Pull Requests**: Keep PRs small and focused. One logical feature or bug fix per PR.

## 3. Versioning & Release Strategy

*   We use **Semantic Versioning (SemVer)**: `vMajor.Minor.Patch`.
*   **Pre-Release Checklist**: Before creating a git tag, you MUST update the version number in the following locations:
    1.  `frontend/package.json` (version field)
    2.  `Makefile` (if `VERSION` variable exists)
    3.  `docs/` (if there are specific version references)
    4.  `backend/cmd/main.go` (if version flag is implemented)
*   **Process**:
    1.  Update files.
    2.  Commit: `chore: bump version to vX.Y.Z`.
    3.  Tag: `git tag vX.Y.Z`.
    4.  Push: `git push origin vX.Y.Z`.

## 4. Coding Standards

*   **Go**: Follow standard Go conventions (effective go). Run `golangci-lint` before committing.
*   **Frontend**: Use functional components and Hooks. Run `pnpm lint`.
*   **Tests**: TDD is encouraged. Ensure all tests pass before submitting.
