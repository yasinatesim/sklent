#!/usr/bin/env bash
# Block writing memory anywhere but .claude/memory/ (home ~/.claude/projects/.../memory or project-root memory/).
set -uo pipefail
source "$(dirname "$0")/_input.sh"

[[ -z "$HOOK_FILE" ]] && exit 0

case "$HOOK_FILE" in
  */.claude/memory/*) exit 0 ;;
  *MEMORY.md|*/memory/*)
    echo "BLOCK: memory yalnızca .claude/memory/ altına yazılır (feedback-memory-location)." >&2
    echo "  Reddedilen yol: $HOOK_FILE" >&2
    echo "  Doğru: <repo>/.claude/memory/<isim>.md  (index: .claude/memory/MEMORY.md)" >&2
    exit 2
    ;;
esac
exit 0
