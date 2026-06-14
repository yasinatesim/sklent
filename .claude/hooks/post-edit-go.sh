#!/usr/bin/env bash
# Runs on every *.go edit. gofmt + named-import alias check + go vet on the package. Block on fail.
set -euo pipefail

API_DIR="examples/e-commerce/api"

input="$(cat)"
file="$(printf '%s' "$input" | grep -oE '"file_path"[^,]*' | head -1 | sed -E 's/.*:\s*"([^"]+)".*/\1/' || true)"

case "$file" in
  *.go) ;;
  *) exit 0 ;;
esac
[ -f "$file" ] || exit 0

if [ -n "$(gofmt -l "$file")" ]; then
  echo "gofmt: $file is not formatted. Run gofmt -w." >&2
  exit 2
fi

if grep -nE '^[[:space:]]+[a-zA-Z_][a-zA-Z0-9_]*[[:space:]]+"[^"]+"' "$file" \
   | grep -vE '^\s*[0-9]+:\s*(_|\.)\s' >/dev/null 2>&1; then
  echo "named import alias detected in $file — rename at package declaration instead." >&2
  exit 2
fi

pkg="./${file#*"$API_DIR"/}"
pkg="$(dirname "$pkg")"
if ! ( cd "$API_DIR" && go vet "$pkg" ) >/dev/null 2>&1; then
  ( cd "$API_DIR" && go vet "$pkg" ) >&2 || true
  echo "go vet failed for $pkg" >&2
  exit 2
fi
exit 0
