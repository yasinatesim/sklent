# CLAUDE.md — Vela Commerce

Loaded into every session. Caveman-lite: short, imperative, no filler.

## Project

Generic, open-source e-commerce reference platform. Own order + payment + admin panel. Internal
stock + per-channel stock counters. 15-min checkout reservation. Two marketplace clients
(skeleton): a generic `hb` and `ty`. The point of the repo is the agent ecosystem, not the code.

## Stack

| Area | Choice |
|---|---|
| Backend | Go 1.25, Gin, GORM, PostgreSQL 16 |
| Frontend | Next.js, React 19, CSS Modules + SCSS, next-intl (TR + EN, JSON-driven) |
| Admin | `/admin` panel for catalog, promotions, coupons, orders |
| Auth | JWT 15m access + 7d rotating refresh, httpOnly cookie, CSRF, rate limit |
| Marketplace | Clients in `api/internal/marketplace/{hb,ty}` — partial; verify before claiming sync works |
| LLM | Pluggable provider registry under `api/internal/llm/`; AES-256-GCM encrypted API keys |
| RAG | `api/internal/rag` — ChromaDB retrieval + LLM generation; deterministic offline fallback |
| Payment | Iyzico 3D Secure under `api/internal/payment/iyzico`. Sandbox; verify before claiming prod |
| Invoice | GIB e-Arşiv proxy concept. Verify before referencing |
| Cart | Guest (session cookie) + member; Zustand; 15-min reservation hold |
| Promotions | Percent/fixed-TL, cart/product/category scope, coupon engine |
| Email | SMTP + React-Email-style templates, sent off the request path |
| Deploy | `docker compose up`: Postgres + ChromaDB + API + web |

## Non-negotiable Rules

1. **Git-flow.** `feature/*` → `development`. `hotfix/*` → `master`. No direct commits to either.
2. **Issue + milestone + label + PR.** Every change. Missing one = stop.
3. **AI agent never feels time pressure.** No shortcuts. Do the full job.
4. **Web changes → Playwright mandatory.** Start dev server, exercise, screenshot, attach.
5. **Unit tests mandatory.** Both `web/` and `api/`. PR body reports coverage delta.
6. **PR pre-check.** Lint + type-check + test + coverage (+ screenshot if UI). All PASS.
7. **Reviewer dispatch before PR.** `wtf-code-reviewer` routes by file pattern. Run in parallel.
8. **BRAID mental model.** Complex task = constraint → fact → step → check. On check fail, loop.
9. **Claude never merges PRs.** Open the PR. User reviews + merges.
10. **File limit.** Target 200–400 lines. Hard cap 800.
11. **PR body must contain `Closes #N`.**
12. **No filler questions — full autonomy.** Pick a sensible default, state it, finish the scope.

## Verification (before claiming done)

Backend: `cd api && go build ./... && go vet ./... && go test ./...`.
Frontend: `cd web && npm run lint && npm run type-check && npm test`.
UI: dev server up + Playwright snapshot saved.

## Frontend rules

- One component per file. `.tsx` exports exactly one default component.
- Arrow functions only. Never `function` declarations.
- Default export at file bottom: `const X = () => {...}; export default X;`.
- No inline styles. `*.module.scss` only. No Tailwind.
- No `React.*` namespace. Destructure hook imports.
- Status-based state for async: single `useState<State>` over a status enum. No parallel booleans.
- Render dispatch via object map. No `&&` chains / ternary trees / `switch` in JSX.
- Modals via central store (Zustand). Never React Context for state.
- SVG icons centralized under `web/src/components/icons/`.
- i18n: all strings in `web/src/i18n/messages/{tr,en}.json`. Never hardcode UI text.

## Backend rules

- No named imports in Go. No `alias "pkg/path"`. Rename at package declaration if conflict.
- Type definitions in `models/` subpackage.
- Tests two-tier: black-box (exported) → `<module>/tests/<name>_test.go`, `package foo_test`.
  White-box (unexported) → module root alongside `export_test.go`.
- Middleware order: Recover → RequestID → Logger → CORS. Recover outermost.
- Every handler uses `c.Request.Context()` for DB calls.
- UPPER_SNAKE_CASE for constants. Map dispatch over `switch`.
- Check existing constants first. Grep `api/internal/constants/` before adding. Extend, don't dup.
- Runtime values via env. Package const = env default only.

## Comments

- Default: no comments. Only when WHY is non-obvious. Never explain WHAT.
- Max 1 comment line. The `no-long-comments` hook blocks violations.

## Routes

- Admin EN: `/admin/promotions`, `/admin/coupons`, `/admin/orders`
- Public TR: `/sepet`, `/odeme`, `/siparis`, `/kategori`, `/urun`, `/arama`

## Security testing

- **Continuous:** `wtf-security` static review, auto-triggered on auth/payment/middleware/input/env diffs.
- **Periodic:** `security-pentest` black-box dynamic test (web + API + network) before releases.

## Hooks (`.claude/hooks/`)

- `post-edit-go.sh` — `go vet` + `gofmt -l` + named-import alias check on every `*.go` edit.
- `post-edit-ts.sh` — eslint on every web `*.{ts,tsx}` edit.
- `pre-commit-verify.sh` — CI-mirror verify lane; blocks direct commits to `master`/`development`.
- `enforce-branch-base.sh` — `feature/*` MUST target `development`; `hotfix/*` MUST target `master`.
- `block-pr-merge.sh` — Claude NEVER merges PRs.
- `no-long-comments.sh` — blocks 2+-line comment blocks.
- `constants-guard-trigger.sh` — enqueue `constants-guard` on new UPPER_SNAKE_CASE declarations.
