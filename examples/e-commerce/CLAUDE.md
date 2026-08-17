# CLAUDE.md — Sklent

Loaded into every session. Caveman-lite: short, imperative, no filler.

## Read first

- `.claude/SESSION_RULES.md` — the hard rules that would otherwise only bite at code-review time.
  A SessionStart hook injects it, so it is already in context; obey it, do not re-read it.
- `.claude/references/project-structure.md` — canonical `web/` layout. Read before creating any
  file under `web/src/`.
- `.claude/memory/MEMORY.md` — durable cross-session facts.

## What this repo is

**Sklent is an agent bootstrap** — a reusable Claude Code ecosystem (agents + skills + references
+ hooks + the BRAID reasoning model) that you drop into any project to build and guard production
code. The headline is [`.claude/`](.claude), not the app.

This is a generic, open-source e-commerce platform — own order + payment + admin panel, internal + per-channel stock, 15-min
checkout reservation, two marketplace skeleton clients (`hb`, `ty`), LLM/RAG product copy, Iyzico
3D Secure sandbox, GIB e-Arşiv invoice proxy. Paths below are relative to the project root.

## Stack

| Area | Choice |
|---|---|
| Backend | Go 1.25, Gin, GORM, PostgreSQL 16 |
| Frontend | Next.js, React 19, CSS Modules + SCSS, next-intl (TR + EN, JSON-driven) |
| Admin | `/admin` panel for catalog, promotions, coupons, orders (status + tracking), manual stock tracking, review moderation |
| Auth | JWT 15m access + 7d rotating refresh, httpOnly cookie, CSRF, rate limit |
| Marketplace | Clients in `api/internal/marketplace/{hb,ty}` — partial; verify before claiming sync works |
| Catalog search | Server-side: diacritic-insensitive `unaccent()` match across both locale titles, description and slug; price range, in-stock and category filters; `newest`/`price_asc`/`price_desc` sort validated against a closed set; paginated with `MAX_PAGE_SIZE`; `/products/facets` returns category counts and the price bounds of the *filtered* set. `ParseListQuery` drops every hostile value before the repository sees it |
| LLM | Pluggable provider registry under `api/internal/llm/`; AES-256-GCM encrypted API keys |
| RAG | `api/internal/rag` — ChromaDB retrieval + LLM generation; deterministic offline fallback |
| Payment | Iyzico 3D Secure under `api/internal/payment/iyzico`. Sandbox; verify before claiming prod |
| Invoice | GIB e-Arşiv in `api/internal/invoice` + `invoice/gib`: per-year document numbering, RG_BASITFATURA payload builder (VAT split out of the gross total), Turkish amount-in-words, host-allowlisted proxy. **The GIB session lives in `InMemoryGIBSessionStore` — in process memory only, never the DB; a restart forcing a fresh login is the desired behaviour** |
| Returns | `api/internal/returnreq` — one transition table drives the state machine (`requested → approved → refunded`, `requested → rejected`); admin screen at `/admin/returns` |
| Shipping | `api/internal/shipping` — flat rate with an optional free-shipping threshold, read from env |
| Cart | Guest (session cookie) + member; Zustand; 15-min reservation hold |
| Promotions | Percent/fixed-TL, cart/product/category scope, coupon engine |
| Email | SMTP + React-Email-style templates, sent off the request path. Order confirmation, password reset, low-stock alert |
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
13. **Reuse before you write.** Grep for an existing component/function/util that already covers the
    case and extend it. Two near-identical renders of the same concept is a reviewer-blocking
    defect, not a style nit. This also covers **sibling screens for the same real-world concept** —
    a feature built into one and never checked against the other is the same defect as a copy-pasted
    util.
14. **"Also make X work for Y" = WIDEN the existing X. Never build a second X.** When a feature
    exists for one entity and is asked for on another, make the existing endpoint / service method
    / table / component accept both — never a parallel `XOther` method set, a second migration, a
    second route tree, or a second component. **Before writing a line:** grep the existing
    implementation and write down (a) which exact endpoint/function/table already does this, (b) the
    single narrowest thing blocking it from serving the new case. Fix only (b).
15. **Minimum code is the default, not an opt-in.** Least code that solves the stated problem, reuse
    what exists, no speculative abstraction, no rewriting a screen that already works.
16. **New route / new migration / new top-level component = justify or don't.** One line on why the
    existing one cannot be widened. A migration whose purpose is "the same thing the other entity
    already has" is the loudest signal you took the wrong path.
17. **`web/src/app/` is a route shell, nothing else.** No components, no `*.module.scss`, no
    `types.ts`, no helpers. `page.tsx` is a one-line re-export from the feature module. Feature code
    lives in `web/src/features/{admin,web}/<module>/`.
18. **Structure is lint-enforced, not taste.** `project-structure/folder-structure` +
    `project-structure/independent-modules` (`web/scripts/*.mjs`) decide where a file goes and who
    may import it. Sibling feature modules may only import each other through an edge **declared
    with a reason** in `independentModules.mjs`. See `MODULE_MIGRATION.md` for the current counters
    and the phase plan; a rule flips to `error` the moment its counter hits zero.
19. **Never persist third-party auth tokens or session state in the database.** Provider tokens and
    OTP session state stay in-process. A restart forcing a fresh login is desired behaviour.

## Verification (before claiming done)

Backend: `cd api && go build ./... && go vet ./... && go test ./...`.
Frontend: `cd web && npm run lint && npm run stylelint && npm run type-check && npm test`.
`npm run lint` runs `eslint .` (not `next lint`), so the structure rules are part of the gate.
Coverage: `npm run test:coverage`.
Dead code: `cd web && npx fallow dead-code` (TS/JS only — no Go support) and
`cd api && go run golang.org/x/tools/cmd/deadcode@latest -test ./...`.

**Never-delete list** (permanent dead-code false positives — deleting these breaks the build or a
runtime path no static analyser can see):

| Reported as unused | Why it stays |
|---|---|
| The SCSS compiler package | Next.js compiles `.scss` without an explicit import |
| Platform-specific optional binaries | installed per-arch, never imported |
| `public/sw.js` and similar | registered by string path, not imported |
| `MarshalJSON` / `UnmarshalJSON` and ORM hooks | called by `encoding/json` via reflection |
UI: dev server up + Playwright snapshot saved.
Before claiming a subsystem "done": run the `intended-vs-implemented` skill (docs vs code gap).

## Frontend rules

- One component per file. `.tsx` exports exactly one default component.
- Arrow functions only. Never `function` declarations.
- Default export at file bottom: `const X = () => {...}; export default X;`.
- No inline styles. `*.module.scss` only. No Tailwind.
- No `React.*` namespace. Destructure hook imports.
- Status-based state for async: single `useState<RequestStatus>`. No parallel `isLoading`/`isError`.
- **Screens whose body swaps per state render through a `STATE_VIEWS` map**, declared at module
  scope, keyed by `REQUEST_STATUS`, holding **component references** — never JSX elements, which
  would build every branch on every render. Type the map as the intersection of every view's props
  (`ComponentProps<typeof AView> & ComponentProps<typeof BView>`). Shared views live in
  `shared/ui/{NullView,LoadingView,ErrorView,EmptyView}`.
- **Forms are the exception:** a form stays mounted while it submits, so it uses `status` for the
  disabled button and an inline message, not a view swap.
- **Never let a state fall through unrendered.** An `ERROR` branch that renders nothing produces a
  silently empty screen — the defect this pattern exists to prevent.
- Render dispatch via object map. No `&&` chains / ternary trees / `switch` in JSX.
- Modals via central store (Zustand). Never React Context for state.
- SVG icons centralized in one icons folder. No inline SVG.
- i18n: all strings in the locale JSON files. Never hardcode UI text.
- Each component folder carries an `index.ts` holding exactly `export { default } from './Name';` —
  import the folder, not the file. Aggregate barrels (several exports, or a hook beside a component)
  stay banned.
- Every exported type/interface lives in `<module>/types/<name>.types.ts`, or the shared types
  folder when 2+ modules use it. `api/*.ts` declares none. Non-exported local types stay in their
  file. `constants/` is exempt — a `typeof X[keyof typeof X]` alias belongs with its const.
- SCSS placement: single consumer → beside the component with the same name; 2+ consumers inside one
  module → `<module>/styles/`; global → the shared styles folder. `assets/` holds binaries only.
- Import order is enforced and auto-fixable by `simple-import-sort`, grouped to mirror the layering.
  Never hand-sort; `eslint --fix` does it.
- **Splitting a fat component: move the state into a hook first.** Extracting JSX while the state
  stays in the parent produces a 30+ prop signature — worse than the complexity you started with.

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

**URL segments are locale-neutral English.** The `[locale]` segment carries the language and
next-intl translates the labels — the path is never a translation surface. All destinations live in
`web/src/constants/routes.ts` (`ROUTES`, `ADMIN_ROUTES`, `localePath`); never hardcode a path.
The API mirrors the two paths it redirects to in `api/internal/constants/frontend_routes.go`, so a
rename cannot drift between the codebases.

- Public: `/cart`, `/checkout`, `/checkout/success`, `/checkout/error`, `/category/[slug]`,
  `/product/[slug]`, `/search`, `/login`, `/forgot-password`, `/reset-password`
- Admin: `/admin`, `/admin/products`, `/admin/orders`, `/admin/promotions`, `/admin/coupons`,
  `/admin/stock-tracking`, `/admin/reviews`

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
- `session-rules.sh` — SessionStart: injects `.claude/SESSION_RULES.md` so the hard rules land
  before the first line of code, not at review time.
- `stop-memory-reminder.sh` — Stop: asks whether this session produced a durable fact worth saving.
- `memory-location.sh` — PreToolUse(Write|Edit): memory is written only under `.claude/memory/`.
- `check-multi-comp.sh` — warns when an edited file declares more than one React component.
