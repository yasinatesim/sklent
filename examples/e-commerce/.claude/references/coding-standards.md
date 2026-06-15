# Coding Standards (language-agnostic)

## Naming

- Descriptive function names. No bare `tick`, `run`, `handle`, `exec`, `fetch`, `load`.
  Use verb + object: `pollImportProgress`, `publishProduct`, `loadProducts`, `handleLoginSubmit`.
- Even effect-internal async functions get a name. Never `(async () => {})()`.
- UPPER_SNAKE_CASE for constants and immutable vars (Go + TS).

## File limits

- Target 200–400 lines per file. Hard cap 800. Split before you exceed it.

## Comments

- Default: no comments. Only when the WHY is non-obvious (hidden constraint, invariant, workaround).
- Never explain WHAT — the code already does that.
- No PR/task references in comments. Those belong in the commit/PR body.
- **Max 1 comment line.** If a WHY needs more, the code is too clever — simplify it. Enforced by
  the `no-long-comments` hook.

## Honesty

- Report outcomes faithfully. If tests fail, say so with the output. If a step was skipped, say it.
- "Done" means verified: it built, the test passed, the screen showed it.
