---
name: ship-pr
description: Push the current branch and open a PR with the standard body template. feature/* targets development; hotfix/* targets master. Body includes summary, verification results, coverage delta, screenshots, reviewer status, and Closes #N.
---

# ship-pr

Push and open a PR. Never merge.

## Preconditions (all must PASS)

- `go build ./... && go vet ./... && go test ./...` (api)
- `npm run lint && npm run type-check && npm test` (web)
- Playwright snapshot saved (if UI touched)
- `wtf-code-reviewer` returned `VERIFIED`

## Target

- `feature/*` → `development`
- `hotfix/*` → `master`

(The `enforce-branch-base.sh` hook rejects a mismatched `--base`.)

## Body template

```
## Summary
<what changed, why>

## Verification
- api: build / vet / test ... PASS
- web: lint / type-check / test ... PASS
- Playwright: <screenshot link>

## Coverage
- <pkg>: <old>% → <new>%

## Reviewer
- wtf-code-reviewer: VERIFIED

Closes #N
```

`Closes #N` is mandatory. The `block-pr-merge.sh` hook prevents merging — the user merges.
