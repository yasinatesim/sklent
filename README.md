# Vela Commerce

> An open-source, **AI-agent-driven** e-commerce reference platform.
> A brandless, generic distillation of a real production storefront — built to show *how* a
> single developer can grow and maintain a full commerce platform with a disciplined
> ecosystem of Claude Code agents, skills, hooks and a BRAID reasoning model.

Vela Commerce is the companion repository to the article
**"AI Agent'lar ile Sıfırdan Uçtan Uca E-Ticaret Projesi Geliştirme"**. It is intentionally
**unbranded** and runnable end-to-end with a single `docker compose up`.

The point of this repo is **not** the code. The point is the *agent ecosystem* that produces and
guards the code: constraints, facts, steps and checks, wired so that doing the wrong thing is
mechanically impossible.

---

## What's inside

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

> **Presence ≠ completeness.** Marketplace sync, payment and invoice are reference skeletons.
> Inspect the code before assuming any external push/import path works end to end.

---

## The agent ecosystem

Everything lives under [`.claude/`](.claude):

- **`agents/`** — `braid-solver`, `constants-guard`, `wtf-code-reviewer` (dispatcher) and the
  language/domain reviewers `wtf-go`, `wtf-js-react`, `wtf-security`, `wtf-ux-playwright`, plus
  `issue-auditor`.
- **`skills/`** — `braid-plan`, `spec-driven-development`, `coverage-gate`,
  `playwright-snapshot`, `ship-pr`, `issue-create`, `security-pentest`, and the reviewer skills.
- **`references/`** — language-agnostic coding standards, backend/frontend standards, security
  standards, git-flow, and the BRAID mental model.
- **`hooks/`** — shell scripts that make the rules impossible to skip (no aliased Go imports,
  no direct commits to protected branches, no agent-initiated merges, no long comments…).

[`CLAUDE.md`](CLAUDE.md) is loaded into every session and is the single source of project rules.

---

## Quick start

```bash
cp .env.example .env
docker compose up --build      # Postgres + ChromaDB + API + web
docker compose run --rm api ./bin/seed   # 8 categories, 25 products, admin user
open http://localhost:3100/tr  # storefront (TR/EN, green theme, light/dark)
```

Default admin (seed): `admin@vela.test` / `admin12345`. Sign in at `/tr/giris`; admins land on
`/tr/admin` (dashboard, products with a full create form, campaigns, coupons, orders).

### Storefront

Home (hero + category grid + best-sellers), category listing, product detail with quantity +
add-to-cart, cart with savings, a two-step checkout (address → card preview → pay), order
success, live search, auth, a gift-coupon modal and a theme toggle — all driven from the Go API
and translated through `next-intl` (TR/EN).

The pixel reference for the visual design is [`index.html`](index.html), a single-file static
prototype. The Next.js app under [`web/`](web) is the real, componentised implementation of it.

### Verify like an agent would

```bash
node e2e/verify.mjs            # drives a real browser: login, add product, see it in store
```

The verify script trusts the **screen the user sees**, not the API response. Green here means
the flow actually works.

---

## Local development (without Docker)

```bash
# Backend
cd api && go build ./... && go vet ./... && go test ./...

# Frontend
cd web && npm install && npm run dev
```

---

## Where the ideas come from

This ecosystem was not invented in a vacuum. It was adapted from public work:

- [Claude Code — Subagents](https://docs.claude.com/en/docs/claude-code/sub-agents)
- [Claude Code — Agent Skills](https://docs.claude.com/en/docs/claude-code/skills)
- [Claude Code — Hooks](https://docs.claude.com/en/docs/claude-code/hooks)
- [Anthropic Cookbook](https://github.com/anthropics/anthropic-cookbook)
- [mukul975/Anthropic-Cybersecurity-Skills](https://github.com/mukul975/Anthropic-Cybersecurity-Skills) — inspiration for the `security-pentest` skill family
- [BRAID: Bounded Reasoning for Autonomous Inference and Decisions (arXiv 2512.15959)](https://arxiv.org/abs/2512.15959)

Each agent, skill and reference file carries its own attribution where relevant. The goal is an
ecosystem you can read, not a black box.

---

## License

MIT — see [LICENSE](LICENSE).
