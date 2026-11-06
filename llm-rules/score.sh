#!/usr/bin/env bash
# score.sh counts prose findings per 1000 words over changed files and exits
# non-zero on any finding. With no argument it reads the working tree, and with
# a revision it reads what changed since that revision.
set -uo pipefail

image="${VALE_IMAGE:-registry.cdn.si/vale}"
root="$(git rev-parse --show-toplevel)"
cd "$root"

if [ $# -gt 0 ]; then
    files=$(git diff --name-only --diff-filter=d "$1")
else
    files=$(git status --porcelain | awk '$1 != "D" {print $NF}')
fi

files=$(echo "$files" | grep -E '\.(md|go|php|tpl)$' || true)
if [ -z "$files" ]; then
    echo "no prose changed"
    exit 0
fi

count=$(echo "$files" | wc -l | tr -d ' ')
words=$(echo "$files" | xargs wc -w 2>/dev/null | tail -1 | awk '{print $1}')
report=$(docker run --rm -v "${TASK_UI_HOST_PWD:-$PWD}:/src" "$image" \
    --output=line $files 2>&1 | grep -E ':[0-9]+:[0-9]+:' || true)

findings=0
[ -n "$report" ] && findings=$(echo "$report" | wc -l | tr -d ' ')

printf '%s\n' "$report"
printf '\n| Files | Words | Findings | Per 1000 words |\n'
printf -- '|------:|------:|---------:|---------------:|\n'
printf '| %d | %d | %d | %.1f |\n' \
    "$count" "$words" "$findings" "$(echo "$findings $words" | awk '{printf "%.1f", $2 ? $1 * 1000 / $2 : 0}')"

[ "$findings" -eq 0 ]
