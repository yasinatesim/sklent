# Contributing to Vela Commerce

Thanks for your interest in contributing! This repository is as much about the **agent ecosystem** as it is about the code, so the workflow is deliberately disciplined. Please read this before opening a PR.

## Code of conduct

Be respectful, constructive and patient. Harassment of any kind is not tolerated.

## Git-flow

Branch policy is **non-negotiable**:

| Branch prefix | Targets | Purpose |
|---|---|---|
| `feature/*` | `development` | New features and enhancements |
| `hotfix/*` | `master` | Urgent production fixes |

- **No direct commits** to `development` or `master`.
- Branch names follow `feature/issue-<N>-<slug>` or `hotfix/issue-<N>-<slug>`.

## Before you start

Every change needs an **issue + milestone + label** first. No issue, no branch.

1. Open or claim an issue using the [bug report](.github/ISSUE_TEMPLATE/bug_report.md) or [feature request](.github/ISSUE_TEMPLATE/feature_request.md) template.
2. Create a branch from the correct base (see git-flow above).

## Development setup

```bash
cp .env.example .env
docker compose up --build
docker compose run --rm api ./bin/seed
```

Or run locally without Docker:

```bash
cd examples/e-commerce/api && go build ./... && go test ./...   # backend
cd examples/e-commerce/web && npm install && npm run dev        # frontend
```

## Coding standards

- **Go:** no named/aliased imports, types in `models/` subpackages, `UPPER_SNAKE_CASE` constants, map dispatch over `switch`, every handler uses `c.Request.Context()`.
- **Frontend:** one component per file, arrow functions only, default export at file bottom, no inline styles (`*.module.scss` only), no Tailwind, no `React.*` namespace, status-based async state, i18n strings in `web/src/i18n/messages/{tr,en}.json`.
- **Comments:** default to none. Only explain a non-obvious *why*. Max 1 comment line — the `no-long-comments` hook enforces this.
- **File limit:** target 200–400 lines, hard cap 800.

See [`CLAUDE.md`](CLAUDE.md) and `.claude/references/` for the full rule set.

## Verification (before opening a PR)

All lanes must pass:

```bash
# Backend
cd examples/e-commerce/api && go build ./... && go vet ./... && go test ./...

# Frontend
cd examples/e-commerce/web && npm run lint && npm run type-check && npm test
```

- **Unit tests are mandatory** for both `web/` and `api/`. The PR body reports the coverage delta.
- **Web changes require Playwright** — start the dev server, exercise the change, screenshot, attach.

## Pull requests

1. Use the [pull request template](.github/PULL_REQUEST_TEMPLATE.md) — fill in every section.
2. The PR body **must** contain `Closes #N` so the issue auto-closes on merge.
3. Ensure lint + type-check + test + coverage (+ screenshot for UI) all pass.
4. Maintainers review and merge. **Contributors do not merge their own PRs.**

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
