#!/usr/bin/env bash
# SessionStart: inject SESSION_RULES.md so the hard rules land before the first
# line of code, not at code-review time.
set -uo pipefail

RULES="${CLAUDE_PROJECT_DIR:-.}/.claude/SESSION_RULES.md"
[[ -f "$RULES" ]] || { echo '{}'; exit 0; }

python3 - "$RULES" <<'PY'
import json, sys

with open(sys.argv[1], encoding="utf-8") as fh:
    body = fh.read()

print(json.dumps({
    "hookSpecificOutput": {
        "hookEventName": "SessionStart",
        "additionalContext": (
            "MANDATORY CODING RULES (.claude/SESSION_RULES.md). "
            "Obey before writing code; recalling them at code review is too late.\n\n" + body
        ),
    }
}, ensure_ascii=False))
PY
exit 0
