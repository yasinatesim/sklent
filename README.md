<h3 align="center">
  <br />
   <a href="https://github.com/yasinatesim/vela-commerce"><img src=".github/assets/logo.svg" alt="Vela Commerce" width="200" /></a>
  <br />
Vela Commerce
  <br />
</h3>

<hr />

<p align="center">An open-source, AI-agent-driven e-commerce reference platform — a brandless distillation of a real production storefront, built to show how a single developer can grow and maintain a full commerce platform with a disciplined ecosystem of Claude Code agents, skills, hooks and a BRAID reasoning model.</p>

<p align="center">
· <a href="./README.tr.md">🇹🇷 Türkçe</a>
· <a href="https://github.com/yasinatesim/vela-commerce/issues">Issues</a>
· <a href="./CONTRIBUTING.md">Contributing</a>
· <a href="https://github.com/yasinatesim/vela-commerce/blob/master/LICENSE">License</a>
</p>

## 📖 About

Vela Commerce is the companion repository to the article **"AI Agent'lar ile Sıfırdan Uçtan Uca E-Ticaret Projesi Geliştirme"**. It is intentionally **unbranded** and runnable end-to-end with a single `docker compose up`.

The point of this repo is **not** the code. The point is the *agent ecosystem* that produces and guards the code: constraints, facts, steps and checks, wired so that doing the wrong thing is mechanically impossible.

> **Presence ≠ completeness.** Marketplace sync, payment and invoice are reference skeletons. Inspect the code before assuming any external push/import path works end to end.

### 📚 Tech Stack

<table>
<tr>
  <td> <a href="https://go.dev">Go 1.25</a></td>
  <td>Backend API with Gin, GORM and PostgreSQL 16.</td>
</tr>
<tr>
  <td> <a href="https://nextjs.org">Next.js</a></td>
  <td>React 19 storefront + admin, CSS Modules + SCSS, no Tailwind.</td>
</tr>
<tr>
  <td> <a href="https://next-intl.dev">next-intl</a></td>
  <td>JSON-driven i18n, Turkish + English.</td>
</tr>
<tr>
  <td> <a href="https://www.postgresql.org">PostgreSQL 16</a></td>
  <td>Catalog, orders, stock and per-channel counters.</td>
</tr>
<tr>
  <td> <a href="https://www.trychroma.com">ChromaDB</a></td>
  <td>Vector retrieval for RAG-backed product copy generation.</td>
</tr>
<tr>
  <td> <a href="https://www.iyzico.com">Iyzico</a></td>
  <td>3D Secure payment integration (sandbox).</td>
</tr>
<tr>
  <td> <a href="https://www.docker.com">Docker Compose</a></td>
  <td>Postgres + ChromaDB + Go API + Next.js web, one command up.</td>
</tr>
<tr>
  <td> <a href="https://docs.claude.com/en/docs/claude-code">Claude Code</a></td>
  <td>Agents, skills, hooks and the BRAID reasoning model that build and guard the repo.</td>
</tr>
</table>

## 🧐 What's inside?

| Area | Choice |
|---|---|
| Backend | Go 1.25, Gin, GORM, PostgreSQL 16 |
| Frontend | Next.js, React 19, CSS Modules + SCSS, next-intl (TR + EN, JSON-driven) |
| Auth | JWT 15m access + rotating 7d refresh, httpOnly cookie, CSRF double-submit, rate limit |
| Cart | Guest (session cookie) + member; 15-minute stock reservation hold |
| Payment | Iyzico 3D Secure (sandbox) under `internal/payment/iyzico` |
| Promotions | Percent / fixed-TL, cart/product/category scope, coupon engine — pure `Evaluate` |
| Marketplace | Skeleton clients in `internal/marketplace/{hb,ty}` (category tree, attributes, publish) |
| LLM / RAG | Pluggable provider registry + ChromaDB-backed product copy generation (`internal/rag`) |
| Invoice | GIB e-Arşiv proxy concept (host-allowlisted) |
| Email | SMTP + templated order/shipping mails, sent off the request path |
| Deploy | `docker compose up`: Postgres + ChromaDB + Go API + Next.js web |

### The agent ecosystem

Everything lives under [`.claude/`](.claude):

- **`agents/`** — `braid-solver`, `constants-guard`, `wtf-code-reviewer` (dispatcher) and the language/domain reviewers `wtf-go`, `wtf-js-react`, `wtf-security`, `wtf-ux-playwright`, plus `issue-auditor`.
- **`skills/`** — `braid-plan`, `spec-driven-development`, `coverage-gate`, `playwright-snapshot`, `ship-pr`, `issue-create`, `security-pentest`, and the reviewer skills.
- **`references/`** — language-agnostic coding standards, backend/frontend standards, security standards, git-flow, and the BRAID mental model.
- **`hooks/`** — shell scripts that make the rules impossible to skip (no aliased Go imports, no direct commits to protected branches, no agent-initiated merges, no long comments…).

[`CLAUDE.md`](CLAUDE.md) is loaded into every session and is the single source of project rules.

## Getting Started

### 📦 Prerequisites

- [Docker](https://www.docker.com) + Docker Compose
- [Go 1.25+](https://go.dev) and [Node.js 20+](https://nodejs.org) (only for local development without Docker)

### ⚙️ How To Use

```bash
cp .env.example .env
docker compose up --build      # Postgres + ChromaDB + API + web
docker compose run --rm api ./bin/seed   # 8 categories, 25 products, admin user
open http://localhost:3100/tr  # storefront (TR/EN, green theme, light/dark)
```

Default admin (seed): `admin@vela.test` / `admin12345`. Sign in at `/tr/giris`; admins land on `/tr/admin` (dashboard, products with a full create form, campaigns, coupons, orders).

**Verify like an agent would:**

```bash
node e2e/verify.mjs            # drives a real browser: login, add product, see it in store
```

The verify script trusts the **screen the user sees**, not the API response. Green here means the flow actually works.

**Local development (without Docker):**

```bash
# Backend
cd api && go build ./... && go vet ./... && go test ./...

# Frontend
cd web && npm install && npm run dev
```

## 🔑 License

* Copyright © 2026 - MIT License.

See [LICENSE](https://github.com/yasinatesim/vela-commerce/blob/master/LICENSE) for more information.

---

_This README was generated with by [markdown-manager](https://github.com/yasinatesim/markdown-manager)_ 🥲
