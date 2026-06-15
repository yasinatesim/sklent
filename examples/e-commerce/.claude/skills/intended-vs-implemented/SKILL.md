---
name: intended-vs-implemented
description: Audit the gap between what the code is documented to do and what it actually does. Run before claiming any feature works, before a release, and whenever a doc/README/CLAUDE.md says a subsystem is "done". Turns "Presence ≠ completeness" into a repeatable check. Especially for marketplace sync, payment, invoice.
---

# intended-vs-implemented

An agent is optimistic by nature — it will say "this works". This skill forces proof. It makes the
project's central rule — **Presence ≠ completeness; verify, don't assume** — a repeatable audit.

## When to run

- Before claiming any feature works.
- Before a release.
- Whenever a doc (README, `CLAUDE.md`, a comment) asserts a subsystem is "done".
- Always for the reference-skeleton areas: `marketplace/{hb,ty}`, `payment/iyzico`, `invoice`.

## Procedure

1. **Collect the INTENDED behaviour.** From the docs/spec/PR body/`CLAUDE.md`: what is this
   subsystem *claimed* to do? Write it as a short list of concrete capabilities.
2. **Trace the IMPLEMENTED behaviour.** Read the real code path end to end. For each claimed
   capability, find the line that delivers it — or note its absence. A stub, a `TODO`, an
   `errors.New("not wired")`, or a happy-path with no error handling counts as **not implemented**.
3. **Run it, don't read it.** Where possible, exercise the path: a unit test, `go test`, or the
   `e2e/verify.mjs` browser run. Trust the **observed result**, not the source's intent.
4. **Produce the gap report.** Three buckets:
   - ✅ **Implemented + verified** (claim matched by code *and* a run/test).
   - ⚠️ **Implemented, unverified** (code exists, no test/run proves it).
   - ❌ **Claimed, missing** (doc says yes, code says no).
5. **Reconcile the docs.** Every ❌ and ⚠️ either gets fixed in code or downgraded in the docs.
   Never leave a doc claiming more than the code delivers.

## Output

A short report: subsystem, each claim, its bucket, and the file:line evidence. CRITICAL gaps
(payment/auth/invoice claiming more than they do) escalate to a hotfix.

---

Adapted from `pm-ai-shipping/intended-vs-implemented` in
[phuryn/pm-skills](https://github.com/phuryn/pm-skills) (MIT License, © 2026 Pawel Huryn).
Rewritten for this repo's stack and caveman style.
