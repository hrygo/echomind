# Smart Release Workflow

Follow these steps to release a new version of EchoMind.

## 1. Strategy & Preparation
### Determine Version
1.  **Identify Current Version**: Check `Makefile` or `GEMINI.md` to find the `<CURRENT_VERSION>` (e.g., `1.1.0`).
2.  **Determine New Version**: Analyze git history (`git log`) since the last tag.
    *   *Patch* (x.x.+1): Bug fixes, minor tweaks.
    *   *Minor* (x.+1.0): New features, non-breaking changes.
    *   *Major* (+1.0.0): Breaking changes.
    *   **Decision**: Set `<NEW_VERSION>`.

### Pre-flight Checks
Ensure you are on the `main` branch and synchronized.
```bash
git checkout main
git pull origin main
```

## 2. Smart Quality Assurance
Run tests and linters. **If linting fails, attempt to auto-fix using `make format` separate tools.**

```bash
# Try to auto-fix format issues if they exist
make format

# Run QA suite
// turbo
make lint
// turbo
make lint-fe
// turbo
make test
// turbo
make test-fe
// turbo
make build-check
```
> **Smart Fix**: If `make lint` still fails after formatting, read the error log and attempt to fix the specific code issues (e.g., missing error checks, unused variables) before proceeding.

## 3. Dynamic Version Bumping
Update the version number from `<CURRENT_VERSION>` to `<NEW_VERSION>` in the following files.

### Target Files
*   `Makefile`
*   `GEMINI.md`
*   `backend/pkg/logger/config.go`
*   `backend/configs/logger.yaml` (Watch out for strict parsing)
*   `frontend/package.json`
*   `docs/api/openapi.yaml`

### Verification & Smart Fix
Search for any leftovers of the old version.
```bash
grep -r "<CURRENT_VERSION>" . --exclude-dir=node_modules --exclude-dir=.git --exclude-dir=bin --exclude-dir=coverage --exclude-dir=tmp
```
> **Smart Fix**: If `grep` returns any files (other than huge lockfiles or historical logs), **automatically update them** to `<NEW_VERSION>`.

## 4. Documentation & Changelog
### Archival
Check `docs/` for stale reports (older than 1 month) and move them to `docs/archive/<YYYY>/` if needed.

### Changelog Update
Generate comments for the `CHANGELOG.md`.
```bash
# View commits since last release (replace v1.1.0 with actual previous tag)
git log --oneline v<CURRENT_VERSION>..HEAD
```
*Action*: Append a new entry to `CHANGELOG.md` with sections for **Features**, **Fixes**, and **Performance**.

## 5. Git Release
Stage all changed files (version bumps + docs).
```bash
git add .
```

Commit and Tag.
```bash
# Replace placeholders with actual values
git commit -m "chore(release): v<NEW_VERSION>"
git tag -a v<NEW_VERSION> -m "Release v<NEW_VERSION>"
```

## 6. Post-Release
Review upcoming tasks.
```bash
cat docs/product/product-roadmap.md
```
