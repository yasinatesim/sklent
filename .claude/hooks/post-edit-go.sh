#!/usr/bin/env bash
# Runs on every *.go edit. go vet + gofmt -l + named-import alias check. Block on failure.
set -euo pipefail

input="$(cat)"
file="$(printf '%s' "$input" | grep -oE '"file_path"[^,]*' | head -1 | sed -E 's/.*:\s*"([^"]+)".*/\1/' || true)"

case "$file" in
  *.go) ;;
  *) exit 0 ;;
esac

[ -f "$file" ] || exit 0

# gofmt
if [ -n "$(gofmt -l "$file")" ]; then
  echo "gofmt: $file is not formatted. Run gofmt -w." >&2
  exit 2
fi

# named import alias: an import line like `alias "pkg/path"` (excluding `_` and `.`)
if grep -nE '^[[:space:]]+[a-zA-Z_][a-zA-Z0-9_]*[[:space:]]+"[^"]+"' "$file" \
   | grep -vE '^\s*[0-9]+:\s*(_|\.)\s' >/dev/null 2>&1; then
  echo "named import alias detected in $file — rename at package declaration instead." >&2
  exit 2
fi

dir="$(dirname "$file")"
if ! go vet "./$dir" >/dev/null 2>&1; then
  go vet "./$dir" >&2 || true
  echo "go vet failed for $dir" >&2
  exit 2
fi

exit 0
