# Security Standards

Canonical for `wtf-security`. Read on any auth / payment / input / middleware / env diff.

## Auth

- JWT: 15-minute access, 7-day **rotating** refresh. Refresh tokens stored hashed.
- Tokens in `httpOnly`, `SameSite`, `Secure` cookies. Never in localStorage.
- CSRF: double-submit. Every state-changing request carries `X-CSRF-Token`, compared to the cookie.
- Rate limit by route. Tighter on login/register and guest order tracking.

## Access control

- Guest resources are reached by an unguessable token; member resources by JWT.
- Check ownership on every member resource. **No IDOR** — `GET /orders/:id` must verify the order
  belongs to the caller.

## Payment

- Never mark an order paid on an unverified callback.
- Verify: required fields present → 3DS status success → order exists → **paid amount equals order
  amount to the cent**. Any mismatch → release reservation, redirect to error.
- Order of writes: `MarkPaid` then `CommitByOrder`. On any failure, release the reservation.

## Input / output

- Validate and bind with limits (`max=...`). Reject oversized payloads.
- Never reflect `err.Error()` into a 500 body. Use an internal-error writer.
- GORM: explicit field lists on write (no mass assignment).

## Proxy / outbound

- Any server-side proxy (e.g. GIB e-Arşiv) **allowlists hosts**. An open proxy is a vulnerability.

## Secrets

- Only from env. No literal fallback (`process.env.X ?? ""`, never a baked-in default).
- Provider API keys encrypted at rest with AES-256-GCM (`internal/crypto`).
