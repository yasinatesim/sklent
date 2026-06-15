#!/usr/bin/env bash
# Runs on every code edit. Block 2+ consecutive comment lines.
set -euo pipefail

input="$(cat)"
file="$(printf '%s' "$input" | grep -oE '"file_path"[^,]*' | head -1 | sed -E 's/.*:\s*"([^"]+)".*/\1/' || true)"

case "$file" in
  *.go|*.ts|*.tsx|*.js|*.mjs) ;;
  *) exit 0 ;;
esac
[ -f "$file" ] || exit 0

# count max run of consecutive lines whose first non-space char starts a // comment
awk '
  /^[[:space:]]*\/\// { run++; if (run > max) max = run; next }
  { run = 0 }
  END { exit (max >= 2 ? 1 : 0) }
' "$file" || {
  echo "no-long-comments: $file has a 2+-line comment block. Max 1 comment line." >&2
  exit 2
}
exit 0
