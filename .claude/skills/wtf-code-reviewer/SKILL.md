---
name: wtf-code-reviewer
description: Dispatch the wtf-code-reviewer agent after any implementation. Inspects diff, fans out to language-specific reviewers (wtf-go, wtf-js-react, wtf-security, wtf-ux-playwright) in parallel, aggregates findings. Mandatory before any PR. Loops until VERIFIED (max 3 iterations).
---

# wtf-code-reviewer (skill)

Dispatch the `wtf-code-reviewer` agent. It routes the diff by file pattern to the language/domain
reviewers, runs them in parallel, and aggregates. Loop until `VERIFIED`, max 3 iterations. See the
agent definition in `.claude/agents/wtf-code-reviewer.md`.
