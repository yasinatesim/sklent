# Product Requirements Document — E-Commerce Platform

## 1. Purpose

A full-stack, unbranded e-commerce platform designed as a reference implementation that exercises an end-to-end production-like ecosystem. The platform proves that a real, non-trivial codebase can be built and maintained with agent-driven tooling.

## 2. Problem Statement

Building a production-grade e-commerce platform requires coordinating many subsystems: catalog management, authentication, cart and checkout, payment processing, order management, promotions, marketplace integrations, content generation, invoicing, and notifications. Most demo projects either oversimplify (todo-list CRUD) or vendor-lock. There is no canonical open-source example that ties all these pieces together in a single deployable stack with both guest and member flows, a full admin panel, and a realistic but sandboxed external integration surface.

## 3. Target Audience

| Persona | Description |
|---|---|
| Store Owner / Admin | Manages catalog, promotions, coupons, orders via admin panel |
| Registered Member | Browses, purchases, views order history |
| Guest Shopper | Browses and checks out without registration |
| Developer | Studies the codebase as a reference for multi-service e-commerce architecture |

## 4. Goals & Success Criteria

| Goal | Success Metric |
|---|---|
| Guest checkout without account | Guest can complete purchase end-to-end |
| Member auth with security | JWT 15m access + 7d rotating refresh, bcrypt, CSRF, rate limiting |
| Stock integrity | 15-minute reservation hold prevents overselling |
| Admin panel | Full CRUD for products, promotions, coupons; order listing |
| Promotion engine | Percent / fixed-TL, product/category/all scope, min cart threshold |
| i18n | Full Turkish + English UI, locale-based routing |
| External integration readiness | Payment, marketplace, invoice sandbox clients |
| One-command deploy | `docker compose up` starts Postgres + ChromaDB + API + web |

## 5. Functional Requirements

### 5.1 Catalog Management

| ID | Requirement | Priority |
|---|---|---|
| CAT-01 | List products with category filtering | P0 |
| CAT-02 | Single product detail by slug | P0 |
| CAT-03 | List all categories (sorted, hierarchical) | P0 |
| CAT-04 | Admin creates products with title, description, price, stock, category, badge, seller | P0 |
| CAT-05 | Client-side real-time product search (title + description) | P1 |
| CAT-06 | SEO-friendly slug generation with Turkish character normalization | P1 |

### 5.2 Authentication & Authorization

| ID | Requirement | Priority |
|---|---|---|
| AUTH-01 | Register with email + password | P0 |
| AUTH-02 | Login with email + password | P0 |
| AUTH-03 | 15-minute access token (JWT, HMAC-SHA256) | P0 |
| AUTH-04 | 7-day rotating refresh token (SHA-256 hashed in DB) | P0 |
| AUTH-05 | httpOnly cookie-based session | P0 |
| AUTH-06 | Role-based access: user / admin | P0 |
| AUTH-07 | CSRF double-submit cookie protection on state-mutating endpoints | P0 |
| AUTH-08 | Rate limiting per IP (token bucket) | P1 |
| AUTH-09 | Optional auth binding (guest still works) | P1 |

### 5.3 Cart & Checkout

| ID | Requirement | Priority |
|---|---|---|
| CRT-01 | Guest cart persisted in localStorage | P0 |
| CRT-02 | Add item, set quantity, remove item, clear cart | P0 |
| CRT-03 | Computed total, count, savings | P0 |
| CRT-04 | Two-step checkout: address → payment | P0 |
| CRT-05 | Address form: name, surname, email, address, district, city | P0 |
| CRT-06 | Card form with live formatting and card type detection | P0 |
| CRT-07 | Place order creates order + stock reservation + sends confirmation email | P0 |
| CRT-08 | Order success page with order ID and amount | P0 |

### 5.4 Payment Processing

| ID | Requirement | Priority |
|---|---|---|
| PAY-01 | Iyzico 3D Secure sandbox integration | P0 |
| PAY-02 | Callback handler verifies mdStatus and amount match | P0 |
| PAY-03 | On success: mark order paid, commit reservations | P0 |
| PAY-04 | On failure: release reservations, redirect to error page | P0 |

### 5.5 Order Management

| ID | Requirement | Priority |
|---|---|---|
| ORD-01 | Guest order tracking via unguessable token | P0 |
| ORD-02 | Member order history (list + detail) | P0 |
| ORD-03 | Admin order listing with status filter, pagination | P0 |
| ORD-04 | Order email notification (async goroutine) | P1 |
| ORD-05 | Admin order status update (paid/shipped/cancelled) with shipping tracking number | P0 |

### 5.6 Stock Reservations

| ID | Requirement | Priority |
|---|---|---|
| RES-01 | 15-minute reservation hold on order placement | P0 |
| RES-02 | Commit reservation on payment success | P0 |
| RES-03 | Release reservation on payment failure | P0 |

### 5.7 Promotions & Coupons

| ID | Requirement | Priority |
|---|---|---|
| PRO-01 | Percent discount type | P0 |
| PRO-02 | Fixed TL discount type (capped at line subtotal) | P0 |
| PRO-03 | Scope: all products, specific products, specific categories | P0 |
| PRO-04 | Minimum cart threshold | P0 |
| PRO-05 | Coupon engine with unique codes | P0 |
| PRO-06 | Pure deterministic `Evaluate()` function with no side effects | P0 |
| PRO-07 | Admin coupon/promotion management: full CRUD (create, list, toggle active, delete) via `/admin/promotions` and `/admin/coupons` | P0 |

### 5.8 Marketplace Integration

| ID | Requirement | Priority |
|---|---|---|
| MKT-01 | HB (Hepsiburada) skeleton client with category tree and attributes | P2 |
| MKT-02 | TY (Trendyol) skeleton client with category tree | P2 |
| MKT-03 | Product publish stubs for both marketplaces | P2 |

### 5.9 LLM / RAG Product Copy

| ID | Requirement | Priority |
|---|---|---|
| LLM-01 | Pluggable LLM provider registry (OpenRouter + offline fallback) | P1 |
| LLM-02 | Product copy enhancement via LLM (SEO title, description) | P1 |
| LLM-03 | ChromaDB-backed retrieval for similar products | P1 |
| LLM-04 | Deterministic offline fallback when no API key configured | P1 |
| LLM-05 | AES-256-GCM encrypted API key storage | P1 |

### 5.10 Invoice (GIB e-Arsiv)

| ID | Requirement | Priority |
|---|---|---|
| INV-01 | Host-allowlist validation for GIB target URLs | P2 |
| INV-02 | Client-side proxy route to avoid CORS on GIB endpoints | P2 |

### 5.11 Email Notifications

| ID | Requirement | Priority |
|---|---|---|
| EML-01 | Order confirmation email on placement (async) | P1 |
| EML-02 | SMTP configuration support | P1 |
| EML-03 | Log-only fallback when no SMTP configured | P1 |
| EML-04 | Password reset email with time-limited token (async) | P1 |
| EML-05 | Low-stock alert email to admin when product stock crosses its threshold (async) | P1 |

### 5.12 Admin Panel

| ID | Requirement | Priority |
|---|---|---|
| ADM-01 | Dashboard with stats, weekly orders chart, recent orders | P0 |
| ADM-02 | Product management: list + create form | P0 |
| ADM-03 | Campaign (promotion) management: full CRUD | P0 |
| ADM-04 | Coupon management: full CRUD | P0 |
| ADM-05 | Order listing with status update + tracking number | P0 |
| ADM-06 | Manual stock tracking: name + quantity ledger, independent of catalog stock | P0 |
| ADM-07 | Review moderation: approve/reject pending product reviews | P1 |

### 5.13 Internationalization

| ID | Requirement | Priority |
|---|---|---|
| I18-01 | Full Turkish UI (default) | P0 |
| I18-02 | Full English UI | P0 |
| I18-03 | Locale-based routing (`/[locale]/...`) | P0 |
| I18-04 | Auto-redirect based on Accept-Language | P1 |
| I18-05 | Locale-aware currency formatting (TRY) | P1 |

### 5.14 Theme

| ID | Requirement | Priority |
|---|---|---|
| THM-01 | Light mode | P0 |
| THM-02 | Dark mode | P0 |
| THM-03 | Theme persisted in localStorage | P0 |
| THM-04 | Theme toggle in header | P0 |

## 6. Non-Functional Requirements

> **Maintainability is a first-class NFR here.** The architecture rules in §9.1 and their
> enforcement in §9.2 are acceptance criteria, not style guidance: the project is considered
> regressed if any structure rule drops back to `warn` or its violation count rises above zero.


| ID | Requirement | Target |
|---|---|---|
| NFR-01 | Backend language | Go 1.25 |
| NFR-02 | Frontend framework | Next.js + React 19 |
| NFR-03 | Database | PostgreSQL 16 |
| NFR-04 | ORM | GORM |
| NFR-05 | Styling | CSS Modules + SCSS (no Tailwind) |
| NFR-06 | API style | RESTful with Gin |
| NFR-07 | One-component-per-file frontend | Enforced |
| NFR-08 | No inline styles | Enforced |
| NFR-09 | Status-based async state | Enforced (no parallel booleans) |
| NFR-10 | Map dispatch over switch/ternary in JSX | Enforced |
| NFR-11 | Centralized i18n JSON files | Enforced (no hardcoded UI text) |
| NFR-12 | Middleware order | Recover → RequestID → Logger → CORS |
| NFR-13 | File size limit | 200–400 lines, hard cap 800 |
| NFR-14 | Unit test coverage | 70% threshold on changed surface |
| NFR-15 | E2E tests | Playwright with real browser |
| NFR-16 | Deploy | `docker compose up` (single command) |
| NFR-17 | No comments by default | Only when WHY is non-obvious |

## 7. User Stories

### Guest Shopper
1. Browse products by category on homepage
2. View product details with price and stock
3. Search products by keyword
4. Add items to cart
5. View cart with quantity controls and subtotal
6. Proceed to checkout
7. Fill in shipping address and payment info
8. Complete order and see success page
9. Track order with provided token

### Registered Member
1. Register and login
2. All guest capabilities
3. View order history
4. Automatic member binding on checkout (no re-entry of info)

### Admin
1. Login with admin credentials → redirected to admin panel
2. View dashboard with stats and recent orders
3. Create new products with full details
4. Manage campaigns and coupons
5. View all orders

## 8. Out of Scope (v2)

- Image upload and CDN integration
- Product variants (size, color, etc.)
- Wishlist / favorites
- Product review star aggregation / storefront review display (v2 adds submit + admin moderation only, see ADM-07)
- Automatic reservation cleanup cron
- Marketplace publish wiring (end-to-end)
- Real invoice generation (server-side)
- Shipping integration with cargo APIs
- Payment installments
- Multi-currency support (TRY only)
- Mobile app
- PWA
- WebSocket / real-time notifications
- SSO / OAuth providers
- Customer management, content/banner management, analytics dashboard, audit log (present in the private reference project this repo is adapted from; not yet ported)

## 9. Technical Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      Browser                             │
│  Next.js 15 + React 19 + Zustand + CSS Modules + SCSS   │
│  next-intl (TR/EN)                                       │
└──────────────────────┬──────────────────────────────────┘
                       │ HTTP / REST
┌──────────────────────▼──────────────────────────────────┐
│                   Go API (Gin 1.12)                       │
│                                                           │
│  ┌─────────┐ ┌──────────┐ ┌──────────┐ ┌────────────┐  │
│  │ Auth    │ │ Catalog  │ │ Orders   │ │ Payment    │  │
│  │ (JWT +  │ │ (GORM)   │ │ (GORM)   │ │ (Iyzico)   │  │
│  │ bcrypt) │ │          │ │          │ │            │  │
│  └─────────┘ └──────────┘ └──────────┘ └────────────┘  │
│  ┌─────────┐ ┌──────────┐ ┌──────────┐ ┌────────────┐  │
│  │ Market  │ │ Promo/   │ │ RAG/LLM  │ │ GIB        │  │
│  │ Place   │ │ Coupon   │ │ (Chroma) │ │ Invoice    │  │
│  └─────────┘ └──────────┘ └──────────┘ └────────────┘  │
└──────┬────────────────────┬──────────────────┬──────────┘
       │                    │                  │
┌──────▼──────┐   ┌────────▼───────┐   ┌─────▼──────────┐
│ PostgreSQL  │   │   ChromaDB     │   │   SMTP / Log   │
│    16       │   │   0.5.20       │   │   (Email)      │
└─────────────┘   └────────────────┘   └────────────────┘
```

### 9.1 Frontend Module Architecture

The web app is organised as **feature modules behind a thin route shell**, not as a flat
`components/ + lib/` tree. This is a hard requirement, not a preference: it is what makes the
codebase safe for agents to extend without an operator reviewing every file placement.

```
web/src/
  app/                      route shell ONLY — page.tsx is a one-line re-export
  features/
    admin/<module>/         admin panel features
    web/<module>/           public site features
  shared/                   constants, types, helpers, hooks, stores, ui, layout, styles, test
```

Module anatomy — the folder set is **closed**; any other folder name at module root is a violation:

```
<module>/
  page.tsx  page.module.scss  constants.ts
  components/<Comp>/<Comp>.tsx + <Comp>.module.scss + index.ts + __tests__/
  views/<View>/             status-dispatch views only
  api/  hooks/  helpers/  store/  styles/  types/  assets/  __tests__/
  [id]/  new/               nested route = sub-module, same anatomy
```

Import direction is one-way — `shared → feature → app`. Sibling feature modules may only reach each
other through an edge **declared with a written reason**. The rule is not prohibition; it is that
every cross-module dependency was a decision somebody recorded.

### 9.2 Architecture Enforcement (Requirement, not documentation)

Prose rules decay. Every structural rule above must be **machine-checked in CI and locally**:

| Requirement | Mechanism |
|---|---|
| Folder/file layout, style placement, eponymous component folders | `project-structure/folder-structure` (`web/scripts/projectStructure.mjs`) |
| Import direction + declared cross-module edges | `project-structure/independent-modules` (`web/scripts/independentModules.mjs`) |
| Exported types live in `types/` | `no-restricted-syntax` |
| Import order mirrors the layers | `simple-import-sort` (auto-fix) |
| No unused imports | `unused-imports` (auto-fix) |
| Code smells, complexity, duplicated branches | `sonarjs` |

**Acceptance criterion:** each rule reaches zero violations and is then set to `error`. A rule
parked at `warn` does not satisfy this requirement — see `MODULE_MIGRATION.md` for current counters
and the phase plan.

### 9.3 Agent Operating Requirements

The repository ships an agent ecosystem under `.claude/`; these are product requirements because
they determine whether generated code is safe to merge:

- **Rules in context before the first edit.** `.claude/SESSION_RULES.md` is injected by a
  SessionStart hook. Recalling a rule during code review means the wrong code already exists.
- **Durable memory in-repo.** `.claude/memory/` is committed, so a fresh clone reproduces the same
  agent behaviour. A hook pins the write location.
- **Mechanical guardrails over instructions.** Hooks block direct commits to protected branches,
  block PR merges by the agent, block multi-line comment blocks, and run the auto-fixable lint lane
  on every edit.
- **Periodic dead-code sweep** with a documented never-delete list, because static analysers cannot
  see reflection-driven calls or string-path asset references.

### 9.4 Catalog Search, Filtering and Sorting

Search runs on the server, never by filtering a fully-downloaded catalog in the browser. The match
is Turkish diacritic-insensitive (`unaccent()`), so "gumus" finds "gümüş", and it covers both locale
titles, the description and the slug — a product surfaces even when the query matches something
other than the leading title.

Filters (category, price range, in-stock) and sort (`newest`, `price_asc`, `price_desc`) are part of
the same query. **The sort value is matched against a closed set** and never interpolated into the
ORDER BY. Page size is capped, so an unbounded `pageSize` cannot be used as a denial-of-service
lever, and a reversed price range is corrected rather than silently returning nothing.

`/products/facets` returns counts for the *same filtered set* the list describes, except that the
category facet ignores the selected category — otherwise choosing one category would hide every
other option.

### 9.5 Returns, Shipping and Invoicing

**Returns.** A buyer may open a return only once the order is `shipped`, only with a known reason,
and only when no other return is still open for that order. The flow is a single transition table
(`requested → approved → refunded`, `requested → rejected`; both leaves terminal) shared by the API
and mirrored in the admin UI, so the screen can never offer a move the server will reject.

**Shipping.** One flat charge with an optional free-shipping threshold, both read from env. A zero
threshold means free shipping is disabled — never that everything ships free.

**Invoicing (GIB e-Arşiv).** The invoice is a GIB `RG_BASITFATURA` payload, not a generic
document: per-year sequential numbering, VAT split out of the gross total, TCKN/VKN validation, and
the Turkish amount-in-words line GIB prints on the invoice. Issuing is refused for an order that is
not `paid` or `shipped`, and refused when the lines do not sum to the order total. The portal
session is held **in process memory only** (`InMemoryGIBSessionStore`) — a restart forcing a fresh
login is the intended behaviour, and persisting a live GIB credential is treated as a security
defect.

## 10. Data Model (Core Entities)

| Entity | Key Fields |
|---|---|
| User | id, email, password_hash, role, full_name, closed_at |
| RefreshToken | id, user_id, token_hash, expires_at, revoked_at |
| Category | id, slug, name_tr, name_en, icon, parent_id |
| Product | id, slug, title_tr, title_en, price_cents, stock, category_id, badge, seller, material, image_url, published |
| Order | id, user_id, guest_token, email, status, payment_method, total_cents |
| OrderItem | id, order_id, product_id, title_tr, unit_cents, quantity |
| Reservation | id, order_id, product_id, quantity, expires_at, committed_at, released_at |
| Promotion | id, name, discount_type, discount_value, scope_type, product_ids, category_ids, min_cart_cents, active |
| Coupon | id, code, discount_type, discount_value, scope_type, min_cart_cents, active |
| StockTrackingItem | id, product_name, quantity |
| Review | id, product_id, author_name, rating, comment, status (pending/approved/rejected) |
| PasswordResetToken | id, user_id, token_hash, expires_at, used_at |

All IDs are UUIDs. Timestamps: created_at, updated_at on every entity.

## 11. API Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | /healthz | None | Health check + CSRF cookie prime |
| POST | /auth/register | Rate-limited | Register |
| POST | /auth/login | Rate-limited | Login |
| POST | /auth/logout | Optional | Logout |
| POST | /auth/refresh | Cookie | Rotate refresh token |
| GET | /auth/me | Cookie/Bearer | Current user info |
| GET | /products | None | List published products |
| GET | /products/:slug | None | Product detail |
| GET | /categories | None | List categories |
| POST | /orders | CSRF | Place order |
| GET | /orders/track/:token | Rate-limited | Guest order tracking |
| GET | /orders | RequireAuth | Member order list |
| GET | /orders/:id | RequireAuth | Member order detail |
| POST | /payments/iyzico/callback | None | Iyzico 3DS callback |
| POST | /auth/password/forgot | Rate-limited | Request password reset email |
| POST | /auth/password/reset | Rate-limited | Reset password with token |
| GET | /products/:slug/reviews | None | List approved reviews for a product |
| POST | /products/:slug/reviews | CSRF | Submit a review (pending moderation) |
| POST | /admin/products | Admin+CSRF | Create product |
| GET | /admin/promotions | Admin+CSRF | List promotions |
| POST | /admin/promotions | Admin+CSRF | Create promotion |
| PATCH | /admin/promotions/:id | Admin+CSRF | Toggle/update promotion |
| DELETE | /admin/promotions/:id | Admin+CSRF | Delete promotion |
| GET | /admin/coupons | Admin+CSRF | List coupons |
| POST | /admin/coupons | Admin+CSRF | Create coupon |
| PATCH | /admin/coupons/:id | Admin+CSRF | Toggle/update coupon |
| DELETE | /admin/coupons/:id | Admin+CSRF | Delete coupon |
| GET | /admin/orders | Admin+CSRF | List orders (status filter, pagination) |
| PATCH | /admin/orders/:id/status | Admin+CSRF | Update order status + tracking number |
| GET | /admin/stock-tracking | Admin+CSRF | List manual stock tracking rows |
| POST | /admin/stock-tracking | Admin+CSRF | Create stock tracking row |
| PATCH | /admin/stock-tracking/:id | Admin+CSRF | Update stock tracking row |
| DELETE | /admin/stock-tracking/:id | Admin+CSRF | Delete stock tracking row |
| GET | /admin/reviews | Admin+CSRF | List reviews (all statuses) |
| PATCH | /admin/reviews/:id/status | Admin+CSRF | Approve/reject a review |

## 12. Dependencies

### Backend (Go)
- gin-gonic/gin — HTTP framework
- joho/godotenv — env loading
- golang-jwt/jwt/v5 — JWT
- google/uuid — UUID generation
- gorm.io/gorm + gorm.io/driver/postgres — ORM
- golang.org/x/crypto — bcrypt

### Frontend (Node.js)
- next — React framework
- react + react-dom — UI library
- next-intl — internationalization
- zustand — state management
- sass — CSS preprocessing
- vitest + @testing-library/react — testing

## 13. Risk & Mitigation

| Risk | Mitigation |
|---|---|
| Payment sandbox differs from production | Sandbox verifier + amount-match check; swap verifier for production |
| LLM API key compromise | AES-256-GCM encryption at rest; env-based config |
| Overselling due to race condition | 15-min reservation + DB-level stock check in transaction |
| Reservation leaks (no cleanup) | Expired reservation column indexed for future cron job |
| Marketplace API changes | Skeleton clients isolate integration points |
| CSRF bypass | Double-submit cookie pattern with per-request token |

## 14. Glossary

| Term | Definition |
|---|---|
| Guest | Unauthenticated shopper |
| Member | Registered and authenticated shopper |
| Reservation | Temporary stock hold with expiration |
| CSRF | Cross-Site Request Forgery |
| 3DS | 3D Secure payment authentication |
| GIB | Gelir Idaresi Baskanligi (Turkish Revenue Administration) |
| e-Arsiv | Turkish electronic archive invoice |
| RAG | Retrieval-Augmented Generation |
| LLM | Large Language Model |
