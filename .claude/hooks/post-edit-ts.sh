#!/usr/bin/env bash
# Runs on every web/src/**/*.{ts,tsx} edit. eslint warn-only (does not block).
set -euo pipefail

input="$(cat)"
file="$(printf '%s' "$input" | grep -oE '"file_path"[^,]*' | head -1 | sed -E 's/.*:\s*"([^"]+)".*/\1/' || true)"

case "$file" in
  *web/src/*.ts|*web/src/*.tsx) ;;
  *) exit 0 ;;
esac

[ -f "$file" ] || exit 0

WEB_DIR="examples/e-commerce/web"
( cd "$WEB_DIR" && npm exec eslint "${file#*"$WEB_DIR"/}" ) >&2 || echo "eslint warnings in $file (non-blocking)" >&2
exit 0
