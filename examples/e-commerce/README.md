<h3 align="center">
  <br />
  <a href="../../README.md"><img src="./assets/logo.svg" alt="Vela Commerce" width="160" /></a>
  <br />
  Vela Commerce
  <br />
</h3>

<p align="center">A brandless, full-stack e-commerce example that exercises the <a href="../../README.md">Sklent</a> agent ecosystem end to end.</p>

<hr />

## 📖 About

Vela Commerce is the worked example that proves Sklent's agents, skills, hooks and the BRAID
reasoning model on a real, non-trivial codebase. It is intentionally **unbranded** and runs end to
end with a single `docker compose up`.

> **Presence ≠ completeness.** Marketplace sync, payment and invoice are reference skeletons.
> Inspect the code before assuming any external push/import path works end to end. Audit it with
> the `intended-vs-implemented` skill.

## 🧐 What's inside

| Area | Choice |
|---|---|
| Backend | Go 1.25, Gin, GORM, PostgreSQL 16 |
| Frontend | Next.js, React 19, CSS Modules + SCSS, next-intl (TR + EN, JSON-driven) |
| Auth | JWT 15m access + rotating 7d refresh, httpOnly cookie, CSRF double-submit, rate limit |
| Cart | Guest (session cookie) + member; 15-minute stock reservation hold |
| Payment | Iyzico 3D Secure (sandbox) under `api/internal/payment/iyzico` |
| Promotions | Percent / fixed-TL, cart/product/category scope, coupon engine — pure `Evaluate` |
| Marketplace | Skeleton clients in `api/internal/marketplace/{hb,ty}` (category tree, attributes, publish) |
| LLM / RAG | Pluggable provider registry + ChromaDB-backed product copy (`api/internal/rag`) |
| Invoice | GIB e-Arşiv proxy concept (host-allowlisted) |
| Email | SMTP + templated order/shipping mails, sent off the request path |
| Deploy | `docker compose up`: Postgres + ChromaDB + Go API + Next.js web |

## 🚀 Getting Started

```bash
cd examples/e-commerce
cp .env.example .env
docker compose up --build               # Postgres + ChromaDB + API + web
docker compose run --rm api ./bin/seed  # 8 categories, 25 products, admin user
open http://localhost:3100/tr           # storefront (TR/EN, green theme, light/dark)
```

Default admin (seed): `admin@vela.test` / `admin12345`. Sign in at `/tr/giris`; admins land on
`/tr/admin` (dashboard, products with a full create form, campaigns, coupons, orders).

### Verify like an agent would

```bash
node e2e/verify.mjs   # real browser: admin login → create product → assert it shows in store
```

The verify script trusts the **screen the user sees**, not the API response.

### Local development (without Docker)

```bash
# Backend
cd api && go build ./... && go vet ./... && go test ./...

# Frontend
cd web && npm install && npm run dev
```

The visual design reference is [`index.html`](index.html), a single-file static prototype; the
Next.js app under [`web/`](web) is the real, componentised implementation.

## 🔑 License

MIT — see the [root LICENSE](../../LICENSE).
