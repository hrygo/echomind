# Contributing to EchoMind

Thank you for your interest in contributing to EchoMind! We welcome contributions from everyone.

## 1. Getting Started

1.  Clone the repository.
2.  Install dependencies (`go mod tidy`, `pnpm install`).
3.  Copy `backend/configs/config.example.yaml` to `backend/configs/config.yaml` and configure it.
4.  Run `make docker-up` to start infrastructure.

## 2. Commit & PR Guidelines

*   **Conventional Commits**: We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.
    *   `feat`: New feature
    *   `fix`: Bug fix
    *   `docs`: Documentation only
    *   `refactor`: Code change that neither fixes a bug nor adds a feature
    *   `test`: Adding missing tests
*   **Atomic Commits**: Commit often. Each commit should do one thing and do it well. Avoid "mega-commits".
*   **Independent Release**: Once a feature module (and its tests) is complete, it should be considered ready for release.
*   **Pull Requests**: Keep PRs small and focused. One feature per PR.

## 3. Versioning Strategy

*   We use **Semantic Versioning (SemVer)**: `vMajor.Minor.Patch`.
*   **Patch**: Backward-compatible bug fixes.
*   **Minor**: Backward-compatible new features.
*   **Major**: Breaking changes.
*   Tag your releases using git tags (e.g., `git tag v1.0.1`).

## 4. Coding Standards

*   **Go**: Follow standard Go conventions (effective go). Run `golangci-lint` before committing.
*   **Frontend**: Use functional components and Hooks. Run `pnpm lint`.
*   **Tests**: TDD is encouraged. Ensure all tests pass before submitting.
