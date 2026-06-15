---
name: spec-driven-development
description: Write a SPEC.md before coding any non-trivial feature. Covers objectives, scope, success criteria, affected files, and test strategy.
---

# spec-driven-development

Before coding a non-trivial feature, write a spec at
`docs/superpowers/specs/YYYY-MM-DD-<feature>-design.md` with:

1. **Objective** — what problem, for whom.
2. **Scope / Non-scope** — explicit boundaries.
3. **Success criteria** — observable, testable.
4. **Affected files / packages.**
5. **Test strategy** — which tests, which tier (black-box vs white-box), what Playwright flow.

The spec is the Constraint/Fact layer the `braid-plan` graph builds on.
