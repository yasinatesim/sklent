# Git Flow (authoritative)

## Branches

- `feature/issue-N-<slug>` → **`development`**. Test environment.
- `hotfix/issue-N-<slug>` → **`master`**. Production.
- **No direct commits** to `development` or `master`. Enforced by `pre-commit-verify.sh`.
- Unknown branch prefixes are blocked.

## Per change

Every change needs: **issue + milestone + label + PR**. Missing one = stop.

## PR body

- Must contain `Closes #N` (or `Fixes #N`) so the issue auto-closes on merge.
- `Closes #N` only auto-closes on a **default-branch (`master`)** merge. For `feature →
  development`, either close manually or put `Closes #N` in the release PR (`development → master`).
- Must report: summary, verification results, coverage delta, screenshots (if UI), reviewer status.

## Merge

- **Claude never merges.** Open the PR; the user reviews and merges. Enforced by `block-pr-merge.sh`.
- `enforce-branch-base.sh` blocks a PR whose branch prefix and `--base` disagree.
