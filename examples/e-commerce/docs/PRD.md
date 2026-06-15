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
| ORD-03 | Admin order listing with status | P1 |
| ORD-04 | Order email notification (async goroutine) | P1 |

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
| PRO-07 | Admin coupon/promotion management in admin panel | P1 |

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

### 5.12 Admin Panel

| ID | Requirement | Priority |
|---|---|---|
| ADM-01 | Dashboard with stats, weekly orders chart, recent orders | P0 |
| ADM-02 | Product management: list + create form | P0 |
| ADM-03 | Campaign management | P1 |
| ADM-04 | Coupon management | P1 |
| ADM-05 | Order listing | P1 |

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

## 8. Out of Scope (v1)

- Image upload and CDN integration
- Product variants (size, color, etc.)
- Wishlist / favorites
- Product reviews and ratings
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
| POST | /admin/products | Admin+CSRF | Create product |

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
