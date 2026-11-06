#!/usr/bin/env bash
# install.sh installs the writing gate into a Claude installation and reports
# what it wrote. Rerunning it is safe: every step replaces what it owns.
set -euo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo="$(dirname "$here")"
claude="${CLAUDE_HOME:-$HOME/.claude}"

say() { printf '%-28s %s\n' "$1" "$2"; }

mkdir -p "$claude/output-styles" "$claude/vale/styles"

# 1. The output style. Claude Code loads it into the system prompt of every
# session once settings.json names it.
cp "$here/writing.md" "$claude/output-styles/writing-gate.md"
say "output style" "$claude/output-styles/writing-gate.md"

# 2. The 43 rules the gate measures against, so a repository without its own
# .vale.ini is still covered by the reply and pull request hooks.
rm -rf "$claude/vale/styles/llm-slop"
cp -r "$repo/.vale/styles/llm-slop" "$claude/vale/styles/llm-slop"
say "vale style" "$(ls "$claude/vale/styles/llm-slop" | wc -l) rules"

# 3. settings.json names the style and keeps every other key.
python3 - "$claude/settings.json" <<'PY'
import json, sys
path = sys.argv[1]
with open(path) as fh:
    settings = json.load(fh)
settings["outputStyle"] = "Writing gate"
with open(path, "w") as fh:
    json.dump(settings, fh, indent=2)
    fh.write("\n")
PY
say "settings.json" "outputStyle=Writing gate"

# 4. The always-loaded context block, between its own markers so a rerun
# replaces the block and leaves the rest of the file alone.
block="$here/CLAUDE.snippet.md"
python3 - "$claude/CLAUDE.md" "$block" <<'PY'
import pathlib, sys
target, source = pathlib.Path(sys.argv[1]), pathlib.Path(sys.argv[2])
body = target.read_text() if target.exists() else ""
start, end = "<!-- llm-rules:writing -->", "<!-- /llm-rules:writing -->"
block = source.read_text().strip()
if start in body and end in body:
    head = body[: body.index(start)]
    tail = body[body.index(end) + len(end) :]
    body = head + block + tail
else:
    body = body.rstrip() + "\n\n" + block + "\n"
target.write_text(body)
PY
say "global CLAUDE.md" "$claude/CLAUDE.md"

# 5. Every repository under the workspace roots reads the same block from a
# session started in its tree, without the global file.
#
# A repository already carrying a context file keeps every file it has. One
# carrying none gets a CLAUDE.md holding the block alone, which the repository
# owner commits or discards.
roots="${LLM_RULES_ROOTS:-$HOME/workspace/github $HOME/workspace/malta}"
count=0
created=""
for root in $roots; do
    [ -d "$root" ] || continue
    for dir in "$root"/*/; do
        [ -d "$dir/.git" ] || [ -f "$dir/CLAUDE.md" ] || [ -f "$dir/AGENTS.md" ] || continue
        names="CLAUDE.md AGENTS.md"
        if [ ! -f "$dir/CLAUDE.md" ] && [ ! -f "$dir/AGENTS.md" ]; then
            printf '# %s\n' "$(basename "$dir")" > "$dir/CLAUDE.md"
            created="$created $(basename "$dir")"
        fi
        for name in $names; do
            [ -f "$dir$name" ] || continue
            python3 - "$dir$name" "$block" <<'PY'
import pathlib, sys
target, source = pathlib.Path(sys.argv[1]), pathlib.Path(sys.argv[2])
body = target.read_text()
start, end = "<!-- llm-rules:writing -->", "<!-- /llm-rules:writing -->"
block = source.read_text().strip()
if start in body and end in body:
    head = body[: body.index(start)]
    tail = body[body.index(end) + len(end) :]
    body = head + block + tail
else:
    body = body.rstrip() + "\n\n" + block + "\n"
target.write_text(body)
PY
            count=$((count + 1))
        done
    done
done
say "repository files" "$count updated"
[ -n "$created" ] && say "CLAUDE.md created in" "$(echo $created | tr ' ' ',')"

printf '\nRestart the session, or run /output-style to pick "Writing gate".\n'
