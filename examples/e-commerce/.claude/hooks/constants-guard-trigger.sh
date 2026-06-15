#!/usr/bin/env bash
# Runs on edit. If the diff adds a new UPPER_SNAKE_CASE const/var, enqueue constants-guard.
set -euo pipefail

input="$(cat)"
file="$(printf '%s' "$input" | grep -oE '"file_path"[^,]*' | head -1 | sed -E 's/.*:\s*"([^"]+)".*/\1/' || true)"

case "$file" in
  *.go|*.ts|*.tsx) ;;
  *) exit 0 ;;
esac
[ -f "$file" ] || exit 0

if grep -nE '(const|var|let)[[:space:]]+[A-Z][A-Z0-9_]+[[:space:]]*=' "$file" >/dev/null 2>&1; then
  echo "constants-guard: new UPPER_SNAKE_CASE declaration in $file — run constants-guard before commit." >&2
fi
exit 0
