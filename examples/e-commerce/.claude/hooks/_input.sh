#!/usr/bin/env bash
# Shared: parse the hook stdin JSON into HOOK_FILE + HOOK_CMD (with env/git fallbacks).
_HOOK_INPUT=$(cat 2>/dev/null || true)
HOOK_FILE=""
HOOK_CMD="${CLAUDE_HOOK_BASH_COMMAND:-}"
if [[ -n "$_HOOK_INPUT" ]] && command -v jq >/dev/null 2>&1; then
  HOOK_FILE=$(printf '%s' "$_HOOK_INPUT" | jq -r '.tool_input.file_path // .tool_input.filePath // empty' 2>/dev/null || true)
  [[ -z "$HOOK_CMD" ]] && HOOK_CMD=$(printf '%s' "$_HOOK_INPUT" | jq -r '.tool_input.command // empty' 2>/dev/null || true)
fi
[[ -z "$HOOK_FILE" ]] && HOOK_FILE="${CLAUDE_HOOK_FILE_PATHS:-}"
