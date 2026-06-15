---
name: braid-plan
description: Generate a BRAID reasoning graph (Mermaid flowchart TD) for a complex multi-step task. Cache to .local-artifacts/braid/<task-slug>.mmd. Hand to the braid-solver agent for execution. Use before non-trivial features, multi-hypothesis debugging, architecture decisions.
---

# braid-plan

Produce a **BRAID graph** (Mermaid `flowchart TD`) for a task, then cache it.

## Steps

1. Read the task. Identify **Constraints** (from CLAUDE.md + references), **Facts** (real file
   paths / API shapes — verify them, do not recall), the ordered **Steps**, and the **Checks**.
2. Every Check has exactly two edges: `Pass` and `Fail`. The `Fail` edge points back to an earlier
   Step so the loop *is* the retry. Do not add a numeric retry cap.
3. Emit the Mermaid graph.
4. Cache it to `.local-artifacts/braid/<task-slug>.mmd`.
5. Hand the slug to the `braid-solver` agent.

## When to use

3+ file refactors, multi-hypothesis debugging, architecture decisions. Skip for one-liners — but
still apply the mental model. See `references/braid-mental-model.md`.
