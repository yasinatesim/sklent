#!/usr/bin/env bash
# Stop hook: remind to save memory once. Must no-op when stop_hook_active
# is true, or the harness re-invokes this on every re-triggered stop forever.
set -uo pipefail
_INPUT=$(cat 2>/dev/null || true)

STOP_HOOK_ACTIVE="false"
if [[ -n "$_INPUT" ]] && command -v jq >/dev/null 2>&1; then
  STOP_HOOK_ACTIVE=$(printf '%s' "$_INPUT" | jq -r '.stop_hook_active // false' 2>/dev/null || echo false)
fi

[[ "$STOP_HOOK_ACTIVE" == "true" ]] && exit 0

echo '{"hookSpecificOutput":{"hookEventName":"Stop","additionalContext":"ZORUNLU: Session bitmeden once kontrol et - bu konusmada ogrenilenler var mi? Varsa memory/ klasorune kaydet (Write tool ile memory/MEMORY.md indexini de guncelle). Kaydedilecek bir sey yoksa gec."}}'
exit 0
