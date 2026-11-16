#!/usr/bin/env bash
# sweep.sh partitions a tree's vale findings into per-agent worklists and prints
# the prompt each agent runs. It writes no file of the repository and exits
# non-zero when a worklist comes out empty.
set -uo pipefail

root="$(git rev-parse --show-toplevel)"
cd "$root"

target="${1:-.}"
budget="${SWEEP_BUDGET:-40}"
out="${SWEEP_DIR:-/tmp/sweep-$(basename "$root")}"
exts='\.(md|go|php|tpl|css|sh|sql)$'

rm -rf "$out"
mkdir -p "$out"

# git decides the file list, because .gitignore already names the vendor trees,
# the agent worktrees and the generated documents. --cached --others
# --exclude-standard names the tracked files plus the untracked files no rule
# excludes.
git ls-files --cached --others --exclude-standard -- "$target" \
	| grep -E "$exts" > "$out/files.txt"

if [ ! -s "$out/files.txt" ]; then
	echo "sweep: no file vale reads under $target" >&2
	exit 2
fi

xargs -r -a "$out/files.txt" -d '\n' vale --output=line 2>/dev/null > "$out/findings.txt"

# One line per file, highest count first. The first worklist holds the densest
# documents and the last holds the tail.
awk -F: '{ print $1 }' "$out/findings.txt" | sort | uniq -c | sort -rn \
	| awk '{ print $1 "\t" $2 }' > "$out/counts.txt"

if [ ! -s "$out/counts.txt" ]; then
	echo "sweep: no findings under $target"
	exit 0
fi

# A worklist fills to SWEEP_BUDGET findings and then a new one opens. A file
# over the budget on its own takes a worklist to itself: one agent, one
# document is the shape, and the budget only groups the tail.
awk -v out="$out" -v budget="$budget" -F'\t' '
	BEGIN { n = 1; spent = 0 }
	{
		if (spent > 0 && spent + $1 > budget) { n++; spent = 0 }
		printf "%s\n", $2 >> sprintf("%s/wl-%02d.txt", out, n)
		spent += $1
		total[n] += $1
		files[n]++
	}
	END {
		printf "%-24s %6s %9s\n", "WORKLIST", "FILES", "FINDINGS"
		for (i = 1; i <= n; i++)
			printf "%s/wl-%02d.txt %6d %9d\n", out, i, files[i], total[i]
	}
' "$out/counts.txt"

cat <<'PROMPT'

Dispatch one agent per worklist, concurrently, with this prompt and the
worklist path substituted for WORKLIST:

  Repo: <root> (shared checkout; other agents are editing OTHER files).
  Your files are the list in WORKLIST. Read it first and open no other file.
  Goal: `vale --output=line $(cat WORKLIST)` prints nothing and exits 0.
  Loop: run it, fix the prose, run it again, until the run prints nothing.
  Read ~/.claude/output-styles/writing-gate.md before the first edit. For a
  finding whose intent is unclear, read .vale/styles/llm-slop/<Rule>.yml: the
  header comment argues the rule and names the fix.
  Rewrite every finding TRUE and SHORTER. Never drop a fact to silence a
  rule, never add a hedge, never invent behaviour.
  Touch prose only: no identifier, no command line, no flag, no string
  literal, nothing inside backticks or a fenced block. In a Go file only the
  comments, and a doc comment never grows.
  Use the Edit tool. Never sed, awk, a heredoc or a python script.
  Do not edit .vale.ini, anything under .vale/, or a baseline file. Do not add
  an ignore marker. Do not commit, stash or checkout.
  Report the before and after counts, and any finding you left with the
  file:line and why.
PROMPT
