---
name: wtf-security
description: Security audit. Triggered by wtf-code-reviewer when the diff touches auth, payment, middleware, input handlers, docker compose, or env templates. Checks IDOR, XSS, SQL injection, secrets, CSP, CSRF, info disclosure, payment bypass, pessimistic locks, mass assignment (GORM), API security posture.
---

# wtf-security (skill)

Invoke whenever the diff touches:

- `api/internal/auth/`, `api/internal/payment/`, `api/internal/order/`, `api/internal/cart/`, `api/internal/promotion/`, `api/internal/admin/` handlers
- `web/src/middleware.ts`, RequireAdmin, login/register/checkout pages
- `docker/docker-compose*.yml`, `Dockerfile*`
- Any new `.env.example` key
- Any handler that binds JSON into a model struct or calls `db.Updates()`

Usage:

```
Agent(subagent_type: "wtf-security", prompt: "Audit these files: <list>. Focus: <auth | payment | input | secrets | csp | mass-assignment | api-posture>. Branch: <name>.")
```

The `wtf-security` agent will:

1. Read `.claude/references/security-standards.md`
2. Check AuthN/AuthZ: server-side guards, JWT validation, session fixation, IDOR
3. Check input → output: sanitization, parameterized SQL, encoding, untrusted LLM content
4. Check **mass assignment**: GORM model binding, explicit field selection on all updates
5. Check payment: no debug/mock bypass, server-side amount, 3DS callback verify, no verbatim provider errors
6. Check info disclosure: no `err.Error()` in 500, no SQL fragments, no PII in logs
7. Check headers: HSTS, CSP (with payment-provider domains + correct nonce coverage), X-Frame-Options, Referrer-Policy
8. Check secrets: no values in source/compose, hardened datastore containers
9. Check CSRF: httpOnly + SameSite + Secure cookies, token rotation
10. Check concurrency: pessimistic locks on stock/balance/quota
11. Check **API posture**: rate limiting on auth endpoints, no keys in query params, consistent error format, pagination bounds

Returns aggregated report with CVE-style severity. Re-dispatch if `NEEDS_FIXES`. Max 3 iterations.

## Full pentest mode

When invoked for a project-wide sweep (not just a diff):

1. Refresh threat model:
   - Pull latest OWASP Top 10 + CWE Top 25
   - `npx skills add https://github.com/vercel-labs/skills --skill find-skills` and look for new security skills
   - WebSearch for CVEs in current dependency versions
2. Walk every handler, middleware, compose file, env template
3. For each finding: dispatch `issue-create` to open an issue under **Phase S: Security Hardening** milestone with severity + file:line + fix proposal
4. CRITICAL findings → also open a `hotfix/issue-N-<slug>` branch from `master` immediately
5. End with a roll-up summary. Coverage is the goal — no class skipped.
