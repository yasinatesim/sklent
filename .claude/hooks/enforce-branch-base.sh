#!/usr/bin/env bash
# Runs on `gh pr create`. Block if branch prefix and --base disagree.
# feature/* MUST target development. hotfix/* MUST target master.
set -euo pipefail

input="$(cat)"
cmd="$(printf '%s' "$input" | grep -oE '"command"[^}]*' | head -1 || true)"

printf '%s' "$cmd" | grep -q 'gh pr create' || exit 0

branch="$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo unknown)"
base="$(printf '%s' "$cmd" | grep -oE '\-\-base[= ]+[a-zA-Z0-9_/-]+' | sed -E 's/.*[= ]+//' || true)"

case "$branch" in
  feature/*)
    if [ "$base" != "development" ]; then
      echo "feature/* must target development (got --base '$base')." >&2; exit 2
    fi ;;
  hotfix/*)
    if [ "$base" != "master" ]; then
      echo "hotfix/* must target master (got --base '$base')." >&2; exit 2
    fi ;;
esac
exit 0
