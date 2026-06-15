---
name: coverage-gate
description: Run unit-test coverage for the changed surface and block if below threshold. Backend (Go) via go test -cover. Frontend (web) via vitest --coverage. Threshold 70% now, target 80%. Reports delta vs base branch.
---

# coverage-gate

Run coverage on the **changed surface** and block below threshold.

- Backend: `cd api && go test -cover ./<changed-pkgs>...`
- Frontend: `cd web && npm test -- --coverage`

Threshold: **70%** now, target **80%**. Report the delta versus the base branch in the PR body.
Block the PR if the changed surface drops below threshold.
