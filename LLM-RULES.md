# LLM rules

Rules an assistant working in these repositories follows. Each numbered rule
holds the prompt that produced it, the settings that enforce it, and the
install steps for a Claude installation starting from an empty `~/.claude`.

`llm-rules/` holds the machine-readable half: the output style, the context
block, the installer and the score command. `.vale/styles/llm-slop/` holds the
45 rules the gate measures against, calibrated in
[CALIBRATION.md](.vale/styles/llm-slop/CALIBRATION.md).

Everything installs under `~/.claude`. No repository holds a copy.

## 1. Writing

### The prompt

Recorded verbatim. The frustration measures the gap. The rules below already
existed, and a vale run corrected the output after the fact on every artifact,
on every turn, for months.

```text
is it my bad that the docs:tone skill isn't loaded from CLAUDE.md? review what
skills you use and should use, reference mandatory usage in various
claude.md/agents.md files in every titpetric user scoped git repository, as well
as all the malta ones (git.cdn.si). All these extra consideration about how you
write code comments, the git commit style? a lot of it is context which is known
before you run the splint tool, before you make the edit. How can I beat the
machine-tone out of your voice for good? perform some goat sacrifices or
whatever to modify your settings, your agents context, or whatever else you need
to *author* code that for once passes any fucking rule defined in splint, for
comments, for composition, all of it. I am 10x more strict than the 1x strict I
expect, and you're fucking this about to be 20% effective for docs or comments
or ANY authored text in various contexts. like fuck, do some research, find some
way to learn, vale checks after violation are fine to correct your outputs, but
I want your outputs to be correct by default. function comments can mostly
follow a fucking template, with SRP and just coherent writing code you should
avoid spaghetti code and kind of respect all that in your planning and execution
work without explicit say so. Fucking fuck! Help a brother out, I'm not
expecting you to write in a british accent, does me asking for a brittish accent
change the way you write? do i have to anthropomorphize you to a personality or
SOUL.md, like I honestly want you to build a workflow, edit your settings, learn
all these things. save this EXACT prompt into titpetric/titpetric as
LLM-RULES.md, point 1: writing. Document all the settings to tweak your output
in comments, in docs, in issues, in PR descriptions, and replies to the user
giving you work, each has an expected output tone, style, and you pick
absolutely brutally idiotic word-sentences together that a regular expression is
a moving fucking goalpost to gate your fucking dumbassery of writing in some
style in the first place. Trained on so many works, and you picked up no style
rules that could match something like the go stdlib as some form of the ideal
expression? do i fucking have you measure vectors and sentence structure, why am
i building a grammar lookup database? i figure people corrected retards grammar
beforehand, and some of these choices on what kind of language are just that,
choices. Do include a checklist of choices on language, ask what kind of style
they prefer, what kind of brevity, what kind of style is undesired, make them
choose a form over some things like oxford comma yes no passive voice yes no,
and apply these settings to the fucking shell, I AM TIRED
```

### The diagnosis

`docs:tone` is a skill. A skill loads when an assistant invokes it. Tone
applies to every token an assistant emits, including the tokens of the sentence
deciding whether to invoke a skill. That layer loads last.

| Layer                         | Reaches                      | Covered before this rule |
|-------------------------------|------------------------------|--------------------------|
| Output style                  | The main loop, every session | No file present          |
| `~/.claude/CLAUDE.md`         | The main loop and subagents  | Replies only             |
| Skill (`/docs:tone`)          | On invocation                | The full specification   |
| `Stop` and `PreToolUse` hooks | After the text exists        | Replies, `gh pr` bodies  |

The specification sat in the two layers that load last. This rule moves it to
the two that load first.

One global ruleset, read in every tree. No repository holds a copy: a rule
that changes needs one edit, and a repository context file documents that
repository.

### Settings

Answers chosen 2026-11-06. An installation that wants different answers edits
`llm-rules/writing.md` and reruns the installer.

| Setting            | Answer                                            |
|--------------------|---------------------------------------------------|
| Oxford comma       | No. `a, b and c`                                  |
| Passive voice      | Allowed. The Go standard library writes it freely |
| Godoc form         | Verb-first from a closed list, no hedge           |
| Terminal reply     | 6 lines of prose, hard cap                        |
| Godoc length       | 1 line, 3 at most                                 |
| Commit message     | Subject only, 72 characters                       |
| Pull request body  | No section under 2 sentences                      |
| Em dash, smart quote, arrow | Never. ASCII only                        |

`Passive` stays allowed because the calibration measured it: 5,142 findings
across twelve standard library packages, 1,198 of them the BSD licence header.
Go's own contract idiom is passive (`ErrNegative is returned`), and a rule
against it bans Go.

### Per-context form

Each artifact has one form. An assistant fills the form.

| Artifact        | Length                     | Opens with                      | Never                                      |
|-----------------|----------------------------|---------------------------------|--------------------------------------------|
| Godoc           | 1 line, 3 at most          | The symbol name, then a verb    | A banned verb, a restatement of the body   |
| Inline comment  | 1 line                     | The constraint                  | What the next line already says            |
| Commit subject  | 72 characters              | `type(scope): ` then a fact     | A body, unless the diff cannot show it     |
| Pull request    | Sections of 2 sentences up | The behaviour after the change  | A preamble, a reviewer address             |
| Issue           | 3 paragraphs               | What breaks and where           | A proposed design before the report        |
| Markdown doc    | One line per paragraph     | The fact                        | A summary of the document in its own body  |
| Terminal reply  | 6 lines plus a list        | The result                      | A heading, a table, a victory lap          |

### The closed verb list

A godoc opens on one of these. The list exists because `VagueVerb` fired 195
times across one repository's `docs/`, more than every other rule combined.

```text
returns  writes   reports  reads    sets     refuses
holds    names    sends    builds   installs resolves
drops    keeps    counts   stops    starts   answers
parses   emits    loads    stores   owns     walks
```

Banned in the same position: `gives`, `provides`, `offers`, `handles`,
`carries`, `deals with`, `makes sure`, `takes care of`, `ends up`, `yields`,
`goes through`, `figures out`, `sorts out`, `grabs`, `spits out`.

### The joins

`RunOnClause` and `So` together account for 131 of 379 findings in a measured
`docs/` tree. Both read one fault: a second clause welded to a finished
sentence.

| Written                        | Fix                                                  |
|--------------------------------|------------------------------------------------------|
| `..., so the table keeps it.`  | Full stop. `The table keeps it.`                     |
| `..., and it returns nil.`     | Full stop. `It returns nil.`                         |
| `..., and the cache is cold.`  | Full stop. `The cache is cold.`                      |
| `..., and then it compiles.`   | Keep. A serial list under one subject.               |
| `..., because the key folds.`  | Keep. A subordinating clause earns its place.        |

Subordinators that earn a comma: `because`, `once`, `until`, `where`, `while`,
`after`, `before`, `unless`. Everything else takes a full stop.

### Install

`llm-rules/README.md` holds the steps. One command on a machine with a Claude
installation:

```bash
./llm-rules/install.sh
```

### Check

```bash
./llm-rules/score.sh            # findings per 1000 words over the working tree
./llm-rules/score.sh HEAD~1     # over a commit range
```

Zero findings is the target. A number above zero names the file, the line and
the rule.
