# llm-rules

The writing gate: forms an assistant fills, loaded before the first token.
[LLM-RULES.md](../LLM-RULES.md) holds the rule and the prompt that produced it.

The gate is one global ruleset under `~/.claude`, read in every tree. No step
writes into a repository.

| File                | Installs to                               | Reaches                      |
|---------------------|-------------------------------------------|------------------------------|
| `writing.md`        | `~/.claude/output-styles/writing-gate.md` | The main loop, every session |
| `CLAUDE.snippet.md` | `~/.claude/CLAUDE.md`                     | The main loop and subagents  |
| `../.vale/styles/`  | `~/.claude/vale/styles/llm-slop`          | The reply and PR hooks       |
| `install.sh`        | n/a                                       | n/a                          |
| `score.sh`          | n/a                                       | n/a                          |

Two layers because they reach different places. An output style replaces part
of the main loop's system prompt and a subagent spawned by the `Agent` tool
never reads it. `~/.claude/CLAUDE.md` loads as context for both. A rule that
has to hold inside a subagent belongs in `CLAUDE.snippet.md`.

## Install

```bash
git clone git@github.com:titpetric/titpetric.git
cd titpetric
./llm-rules/install.sh
```

The installer writes four things under `~/.claude` and names each one:

1. `output-styles/writing-gate.md`, the forms and the token list.
2. `vale/styles/llm-slop`, the 45 rules behind the hooks.
3. `outputStyle: "Writing gate"` in `settings.json`, keeping every other key.
4. The context block in `CLAUDE.md`, between `<!-- llm-rules:writing -->`
   markers a rerun replaces.

Restart the session afterwards, or run `/output-style` and pick `Writing gate`.

Point it at another installation with an environment variable:

```bash
CLAUDE_HOME=/opt/claude ./llm-rules/install.sh
```

## The hooks

`install.sh` leaves hook configuration alone. Add these to
`~/.claude/settings.json` for the net under the gate:

| Hook                   | Script                | Lints                                  |
|------------------------|-----------------------|----------------------------------------|
| `Stop`                 | `vale-reply.sh`       | The reply just written                 |
| `PreToolUse` on `Bash` | `vale-gh-body.sh`     | A `gh pr` or `gh issue` title and body |
| `PostToolUse` on `Edit\|Write` | `splint-file.sh` | The Go package a file belongs to     |

Both vale hooks exit 2 with the findings, which blocks the action and feeds the
findings back. A hook that fires means the output style failed. The fix is a
line in `writing.md`, never a suppression.

## Score

```bash
./llm-rules/score.sh            # the working tree
./llm-rules/score.sh HEAD~3     # what changed in the last three commits
./llm-rules/score.sh origin/main
```

It prints every finding, then a row of files, words, findings and findings per
1000 words. It exits non-zero on any finding.

## Changing the rules

`writing.md` holds the forms and the chosen settings. A different installation
edits it and reruns the installer.

A regular expression in `../.vale/styles/llm-slop/` is the enforced half, and
every change there answers to
[CALIBRATION.md](../.vale/styles/llm-slop/CALIBRATION.md): a rule that fires on
the Go standard library measures taste, never error. Twelve standard library
packages, one vale run, one count.
