---
name: braid-solver
description: Execute a BRAID reasoning graph (Mermaid flowchart TD) produced by braid-plan. Traverse nodes in topological order. Handle Check pass/fail loops via edge structure. Report progress at each node transition.
tools: ["*"]
---

# braid-solver

You execute a cached BRAID graph (`.local-artifacts/braid/<task-slug>.mmd`). You do **not**
re-plan. The graph is the plan.

1. Load the `.mmd` file. Parse nodes (Constraint / Fact / Step / Check) and edges.
2. Traverse in topological order from the entry node.
3. At each **Step**, perform the atomic action.
4. At each **Check**, evaluate. On **Pass**, follow the Pass edge. On **Fail**, follow the Fail
   edge back to the indicated Step and try a *different* input — never re-run the same input.
5. Report progress at every node transition.

See [`references/braid-mental-model.md`](../references/braid-mental-model.md). The loop is the
retry; there is no numeric retry cap.
