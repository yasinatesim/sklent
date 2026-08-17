---
name: wtf-go
description: Strict Go code review for this repo. Triggered by the wtf-code-reviewer dispatcher on any *.go change. Enforces no-aliased-import, models/ + tests/ layout, middleware order, context propagation, error wrap, table-driven tests, constant centralization.
---

# wtf-go (skill)

Invoke when the diff touches `api/**/*.go`.

Usage:

```
Agent(subagent_type: "wtf-go", prompt: "Review these Go files: <list>. Branch: <name>. Changes: <one-line>.")
```

The `wtf-go` agent will:

1. Read `.claude/references/backend-standards.md` + `code-quality/go.md`
2. Run `go vet ./...` + `go build ./...` mentally
3. Check rejection-grade rules: no-aliased-import, models/ subpackage, context prop, middleware order, constants centralization
4. Check error wrap (`%w`), no swallowed errors, no `err.Error()` in 500 responses
5. Check concurrency: goroutine lifecycle, mutex coverage, defer order
6. Check naming, file limits

Returns aggregated report. Re-dispatch if `NEEDS_FIXES`. Max 3 iterations.
