# Calibration

A rule that fires on the Go standard library is measuring taste, not error.

Go's team writes Go for a living. A construction they write freely belongs to
the language. A rule against it bans Go. One they never write belongs to the
writer.

That is the test every rule here answers to. It costs one command: copy a
dozen standard library packages into a directory, point vale at them with this
style, and count.

```bash
mkdir /tmp/std && cd /tmp/std
for p in strings bytes io os errors sort bufio net/http encoding/json time fmt regexp; do
  mkdir -p $(dirname $p); cp -r $(go env GOROOT)/src/$p $(dirname $p)/
done
docker run --rm -v /tmp/std:/w -w /w registry.cdn.si/vale --output=line . | wc -l
```

## What the measurement found

Twelve packages, **10,361 findings**. Half of them came from one rule.

| Rule | Before | After | What the count was measuring |
|---|---|---|---|
| `Passive` | 5,142 | removed | 1,198 of them were the BSD licence header in every file: `is governed by`, `can be found in`. The rest was Go's contract idiom: `ErrNegative is returned`, `the file is closed`. Passive voice is an inverted tell, which a human technical writer writes freely and a model avoids. |
| `TooWordy` | 1,704 | 211 | 118 of its 216 tokens were single words. It flagged `multiple`, `equivalent`, `maximum`, `implement`, `subsequent`, `additional`, `However`: vocabulary, never wordiness. `it is` accounted for another 447, and `It is safe to call Close twice` is how Go states a contract. |
| `WasState` | 461 | 44 | A bare `was` matched `the first non-EOF error that was encountered`. A past participle in a relative clause is the only way to say a thing already happened. The rule now takes the auxiliary opening a clause, and the 44 survivors are real: `Was a bug:`, `Used to be issue 8305`. |
| `Assumption` | 516 | 222 | `should be` was 293 of it. `the buffer should be at least n bytes` states a requirement, never a state of mind. |
| `VagueVerb` | 349 | 159 | `gives` and `giving` were 74. Go documents a return that way. `gives back` stays, because nothing returns backwards. |
| `NegativeParallelism` | 205 | ~50 | `instead of` and `rather than` were 153 between them. Both are ordinary English. |
| `Anthropomorphism` | 230 | 123 | `wants to` was 93, nearly every one of them `if you want to`, which addresses a reader and animates no code. |

**10,361 to 2,565.** Every cut is one of those counts. `Passive` and
`TooWordy` were both verbatim write-good vendoring, calibrated against no
corpus at all.

## What the measurement protects

Nineteen of these rules fire **zero** times on the standard library. Those are
the well-aimed ones. They are the reason the style is worth running.

The three added most recently were measured before they were written, not
after:

| Rule | Rate against the standard library |
|---|---|
| `Posture` | `sits on` 240x, `stands on` 17x, `lives in` 5x |
| `OwnAppositive` | `of its own` 98x, `on its own` 17x |
| `QuantifierSubject` | `has none` 124x, `and none` 12x, `neither is` 11x |

Four candidates failed the same test and never went in: `none is` at 3.1x,
`both are` at 1.4x, bare `its own` at 11.9x over 305 standard library uses, and
`either is`, which this tree never writes.

## The page as its own subject

`PageNarration` is the rule a generated manual asked for. A reference page
that opens on "This page is a stub", names the generator that printed it and
labels an empty section "Not generated" has told the reader about the build
and not about the symbol. The reader runs no build.

Every token went through the same corpus first, and the rule fires **zero**
times on it. So do the additions to four other rules, measured as whole
phrases over the twelve packages:

| Rule | Added | Corpus |
|---|---|---|
| `PageNarration` | `this page is`, `is a stub`, `not generated`, `written by hand`, `the generator writes`, `the template prints`, `generated from the`, `every page is`, `has you covered` | 0 |
| `Register` | `powerful`, `versatile`, `effortlessly`, `unleash`, `full potential`, `rich set of`, `wealth of`, `cater to`, `tailored to`, `best practice`, `dive into`, `the world of`, `shed light on`, `gain insight`, `pave the way`, `a testament to` | 0 |
| `Importance` | `it is important to`, `important to note`, `plays a key role` | 0 |
| `Filler` | `Moreover`, `the journey` | 0 |
| `ReviewerAddress` | `take a look at`, `whether you are` | 0 |

Three candidates failed and stayed out: `exposes` and `exposed` at 11 between
them, `Additionally` at 13, `Furthermore` at 3. `on this page` measures zero
and stays out anyway, because it is the heading over a section list and names
no machinery.

## Two things the tokens cannot do

**A wrapped phrase needs `\s+`.** A literal space never matches a newline, so
`'\bmatters because\b'` passes a comment wrapped after "matters". 1,020 tokens
across 35 rules held a literal space, and every one of them passed a wrapped
phrase in silence. They are swept now, for 45 more findings across this
workspace and no more silent misses.

**A duplicate key drops the rule, in silence.** A `swap` map with the same
word on two lines parses in Python and fails vale's loader, and vale reports
neither the file nor the rule: `Register` stopped firing on `leverages` and
`utilizes` the moment a second `effortless` landed in it. The check is one
command, and it is worth running after any edit to a substitution rule:

```bash
python3 - .vale/styles/llm-slop/*.yml <<'EOF'
import sys, yaml
class Strict(yaml.SafeLoader): pass
def mapping(loader, node):
    seen = set()
    for key, _ in node.value:
        k = loader.construct_object(key)
        if k in seen:
            raise yaml.YAMLError("duplicate key %r" % k)
        seen.add(k)
    return yaml.SafeLoader.construct_mapping(loader, node)
Strict.add_constructor(yaml.resolver.BaseResolver.DEFAULT_MAPPING_TAG, mapping)
for f in sys.argv[1:]:
    try:
        yaml.load(open(f), Loader=Strict)
    except yaml.YAMLError as e:
        print(f, e)
EOF
```

A plain grep for a repeated line prefix reports eleven of these rules and
every one of the eleven is a token list, where the colon it splits on sits
inside a regex. The loader is the check; the grep is not.

**RE2 reads neither lookahead nor lookbehind.** `on its own line` is the bare
possessive and legal. `measured on its own.` is the appositive and not. With
no lookahead the rule matches the punctuation instead of asserting it. Where even
that fails, the gap is written into the rule file: see `WasState.yml`.

## What a regex cannot reach at all

Whether a comment is true, and whether it says anything.

`splint --linters=prose` answers the second from the other direction. It masks
the subject matter out of every sentence, counts what remains, and reports the
shapes a tree writes more often per word than the standard library does. The findings
come out of the measurement and not a list. That is how `beside` turned up at
95 uses here against **one** in two million words of Go comments.

`splint --linters=godoc` holds the rule a ban list cannot express. A comment
whose content words all come off the signature above it says nothing. The fix
is deletion instead of a rewrite. It fires on 0.36% of documented
functions in the standard library, and reading those found no false one.
