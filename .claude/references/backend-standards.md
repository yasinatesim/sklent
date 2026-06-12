# Backend Standards (Go 1.25)

Rejection criteria for `wtf-go`. Read before any Go work.

## Package layout

- Domain packages under `internal/`: `auth`, `cart`, `order`, `payment/iyzico`,
  `marketplace/{hb,ty}`, `llm`, `rag`, `product`, `promotion`, `reservation`, `invoice`, …
- Type definitions live in the module's `models/` subpackage (e.g. `internal/order/models`).
- Runtime values come from env. A package const is only an env default.

## Imports

- **No named imports.** `alias "pkg/path"` is forbidden. On a collision, rename the package at its
  declaration. Enforced by `post-edit-go.sh`.

## Tests — two-tier rule (Go toolchain constraint)

- **Black-box (exported API)** → `<module>/tests/<name>_test.go`, `package foo_test`.
- **White-box (unexported helpers)** → module root alongside `export_test.go`. Must stay in the
  same directory as `export_test.go`; moving it breaks `undefined` at build time. This is exactly
  how Go's own `fmt`, `bufio`, `bytes` packages work.
- Prefer table-driven tests.

## HTTP / middleware

- Middleware order: **Recover → RequestID → Logger → CORS**. Recover outermost.
- Every handler threads `c.Request.Context()` into DB and outbound calls.
- Background work that outlives the request uses `context.WithoutCancel(c.Request.Context())`.

## Errors

- Wrap with `%w`. Never leak `err.Error()` to a 500 response body — use an internal-error writer.
- GORM writes use explicit field lists (`Select`/`Updates(map)` with named fields). No blind
  struct updates → no mass assignment.

## Constants

- UPPER_SNAKE_CASE. Grep `internal/constants/` before adding. Extend, don't duplicate.
- Map dispatch over `switch` where it reads cleaner.

## Honesty

- Presence ≠ completeness. Inspect the code before claiming a feature works end to end. This is
  doubly true for marketplace sync, payment and invoice.
