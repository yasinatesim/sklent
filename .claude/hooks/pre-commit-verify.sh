#!/usr/bin/env bash
# Runs on `git commit` / `gh pr create`. CI-mirror verify lane.
# Blocks direct commits to master/development and unknown branch prefixes.
set -euo pipefail

input="$(cat)"
cmd="$(printf '%s' "$input" | grep -oE '"command"[^}]*' | head -1 || true)"

# Only act on commit / PR-create commands. Everything else passes through.
case "$cmd" in
  *"git commit"*|*"gh pr create"*) ;;
  *) exit 0 ;;
esac

branch="$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo unknown)"

case "$branch" in
  master|development)
    echo "Direct commits to '$branch' are forbidden. Use feature/* or hotfix/*." >&2
    exit 2 ;;
  feature/*|hotfix/*) ;;
  *)
    echo "Unknown branch prefix '$branch'. Use feature/* or hotfix/*." >&2
    exit 2 ;;
esac

if true; then
  if [ -d api ]; then
    ( cd api && go build ./... && go vet ./... && go test ./... ) >&2 || { echo "api verify lane FAILED" >&2; exit 2; }
  fi
  if [ -d web ] && [ -f web/package.json ]; then
    ( cd web && npm run -s type-check && npm test --silent ) >&2 || { echo "web verify lane FAILED" >&2; exit 2; }
  fi
fi

echo "memory reminder: update .claude/memory before pushing if state changed." >&2
exit 0
