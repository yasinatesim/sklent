#!/usr/bin/env bash
# Runs on `gh pr merge`. Hard block — Claude NEVER merges PRs.
set -euo pipefail

input="$(cat)"
cmd="$(printf '%s' "$input" | grep -oE '"command"[^}]*' | head -1 || true)"

if printf '%s' "$cmd" | grep -q 'gh pr merge'; then
  echo "Claude never merges PRs. Open the PR; the user reviews and merges." >&2
  exit 2
fi
exit 0
