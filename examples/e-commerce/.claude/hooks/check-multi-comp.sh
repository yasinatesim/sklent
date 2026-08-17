#!/bin/bash
# Stop hook: warns (does not block) if any .ts/.tsx file touched this session
# declares more than one React component (react/no-multi-comp).
set -euo pipefail

cd "${CLAUDE_PROJECT_DIR:-.}"

files=$(
  {
    git diff --name-only HEAD -- '*.tsx' '*.ts' 2>/dev/null
    git diff --name-only --cached -- '*.tsx' '*.ts' 2>/dev/null
    git ls-files --others --exclude-standard -- '*.tsx' '*.ts' 2>/dev/null
  } | sort -u
)

existing=()
for f in $files; do
  [ -f "$f" ] && existing+=("$f")
done

if [ "${#existing[@]}" -eq 0 ]; then
  echo '{}'
  exit 0
fi

output=$(npx eslint --no-eslintrc --parser @typescript-eslint/parser --plugin react \
  --rule '{"react/no-multi-comp":"warn"}' \
  --parser-options=ecmaFeatures:{jsx:true} \
  "${existing[@]}" 2>/dev/null || true)

if echo "$output" | grep -q "no-multi-comp"; then
  hits=$(echo "$output" | grep -E '^/' | sed 's|.*/||' | sort -u | paste -sd, -)
  node -e '
    const hits = process.argv[1];
    console.log(JSON.stringify({
      systemMessage: `⚠️ react/no-multi-comp: ${hits} define more than one component per file — extract each into its own file before finishing.`
    }));
  ' "$hits"
else
  echo '{}'
fi
