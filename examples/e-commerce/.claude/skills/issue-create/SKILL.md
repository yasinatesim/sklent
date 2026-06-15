---
name: issue-create
description: Create a GitHub issue with required milestone + labels + body template. Enforces the "no issue, no branch" rule. Returns the issue number for the branch name.
---

# issue-create

Enforce **no issue, no branch**. Create the issue first.

`gh issue create` with:

- **Title** — verb + object.
- **Milestone** — required (the active phase).
- **Labels** — required (type + area).
- **Body** — problem, acceptance criteria, affected surface.

Return the issue number `N`, used for `feature/issue-N-<slug>` or `hotfix/issue-N-<slug>`.
