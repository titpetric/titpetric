---
name: Writing gate
description: Forms to fill for every artifact, a closed verb list, and the tokens that never appear. Loads before the first token.
---

# Writing

Every artifact has a form. Fill the form. Composing freely is what produces the
tokens below.

## Before emitting any prose

Scan what you wrote for these. Each one is a regular expression in
`.vale/styles/llm-slop/`, calibrated against twelve Go standard library
packages. A hit is the prose being wrong.

| Token                                          | Rule                 | Fix                                  |
|------------------------------------------------|----------------------|--------------------------------------|
| `, so `                                        | So                   | Full stop                            |
| `, and it/they/this/there/we/the X <verb>`     | RunOnClause          | Full stop                            |
| `which is what`, `which is why`, `the point is`, `in other words` | Restatement | End at the fact        |
| `gives`, `provides`, `offers`, `handles`, `carries`, `deals with` | VagueVerb | A verb from the closed list |
| `makes sure`, `takes care of`, `ends up`, `yields`, `goes through` | VagueVerb | Name the operation    |
| `it is`, `there is`, `in order to`             | TooWordy             | Cut                                  |
| `instead of`, `rather than`, `not X but Y`     | NegativeParallelism  | State what is true                   |
| `several`, `a few`, `nothing else`, `various`  | Weasel, Nothing      | The number, or the name              |
| `load-bearing`, `crucial`, `the crux`          | Importance           | Say what breaks                      |
| `at the seam`, `where X meets Y`, `the glue`   | MetaphoricalLocation | Name the file and the function       |
| `wants to`, `knows about`, `decides to`        | Anthropomorphism     | State the contract                   |
| `sits under`, `lives in`, `stands as`          | Posture, SubstituteIs| `is`, or the location                |
| `was`, `used to`, `previously`, `no longer`    | WasState, SupersededState | State what is true now          |
| `—`, `–`, `…`, `'`, `"`, `→`                   | Punctuation          | `-`, `...`, `'`, `"`, `->`           |

Allowed and often correct: passive voice, `should be`, `multiple`,
`equivalent`, `However`. The calibration measured each against the Go standard
library and found a human technical writer using it freely.

Lists of three or more take no Oxford comma: `functions, classes and
constants`.

## Joining clauses

One finished sentence takes a full stop. A comma plus a connective welds a
second sentence onto the first, which is the single most frequent fault
measured.

| Written                          | Fix                                           |
|----------------------------------|-----------------------------------------------|
| `X holds the key, so Y reads it` | `X holds the key. Y reads it.`                |
| `X parses it, and it returns`    | `X parses it. It returns nil on a bad key.`   |
| `X runs, and the cache is cold`  | `X runs. The cache is cold.`                  |
| `X parses, compiles and runs it` | Keep. Three verbs, one subject.               |
| `X refuses it, because Y folds`  | Keep. A subordinator earns the comma.         |

Subordinators that earn a comma: `because`, `once`, `until`, `where`, `while`,
`after`, `before`, `unless`.

## Go doc comments

One line. Open on the symbol name, then a verb from this list:

```text
returns  writes   reports  reads    sets     refuses
holds    names    sends    builds   installs resolves
drops    keeps    counts   stops    starts   answers
parses   emits    loads    stores   owns     walks
```

```go
// AreaName returns the name a class answers to under its area.
// ParseConfig reads path and refuses a key no field matches.
// Shutdown stops the server, draining for one second.
```

A second line is for one constraint a reader cannot read off the code: a
protocol rule, a refused input, a hard-won ordering. Three lines is the
ceiling, and `splint --linters godoc` reports a comment that runs past the
budget of one line per function plus one per block it opens.

Write no comment at all where the signature already says it. An unexported
helper named `inClause(n int) string` needs none.

## Inline comments

One line. Two is a defect. The line names what the code cannot state:

```go
// Exact, not folded: Runtime\PX\Datetime would resolve to the class DateTime.
```

Never a comment that paraphrases the statement under it.

## Commit messages

Subject only, 72 characters, `type(scope): ` then a fact about the behaviour
after the change. A body only where the diff cannot show the fact.

```text
feat(runner): a class answers under the area that installs it
fix(parser): a namespaced file accepts a const declaration
```

## Pull request bodies

One line per paragraph, no hard wrapping. Every section runs to two sentences
or folds into the one above it. Open on the behaviour after the change. A
measured number goes in a table or a pasted output block, never in a sentence.

Never: a preamble, a reviewer address (`Worth a second look`, `Open
questions`), a summary of the body at the end of the body.

## Issues

Three paragraphs. What breaks, where it breaks, what reproduces it. A proposed
design comes after the report or in a comment, never in place of it.

## Terminal replies

Six lines of prose, hard cap, then a next-steps list where each entry is a
problem and what to do about it. No `##` heading, no table, no
"What changed / How it works / Limitations" scaffolding, no restatement of the
checks you ran.

A finding is one of two things: a change, which you make, or a problem, which
you state with a proposal. `BUG` goes first in the body, naming what it breaks
and the fix.

## Code, same rules

Single responsibility per function. Split a function that needs a comment
explaining its second job. Match the surrounding file's naming, comment density
and idiom. A reader of the diff cannot tell which lines an assistant wrote.
