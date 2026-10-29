# atkins skills

Reusable pipeline files for [atkins](https://github.com/titpetric/atkins),
a YAML task runner. A skill is an ordinary pipeline whose jobs are merged
into every project's job namespace, prefixed with the filename stem: the
job `fmt` in `go.yml` is run as `atkins go:fmt`.

Each skill carries a `when.files:` gate. The skill is only offered when
one of the listed files or directories exists in the project, so a Go
repository never sees the Tailwind jobs and a repository without
migrations never sees the schema jobs.

## Install

Copy the files into `~/.atkins/skills/`, where they apply to every
project:

```
git clone https://github.com/titpetric/titpetric
cp titpetric/.atkins/skills/*.yml titpetric/.atkins/skills/*.md ~/.atkins/skills/
```

A project can pin its own copies under `<project>/.atkins/skills/`, which
is searched first and survives `atkins --jail`. `atkins --vendor --write`
writes the skills a repository uses into that directory.

Only `*.yml` directly inside a skills directory is discovered;
subdirectories are ignored. A skill may carry a markdown file of the same
name beside it, `go.md` beside `go.yml`, which atkins prints as the usage
guide for that skill and which travels with the skill when vendored.

Confirm what was picked up with `atkins -l`, which prints a section per
source file.

## The skills

| File           | Gate                                | What it does                                            |
|----------------|-------------------------------------|---------------------------------------------------------|
| `go.yml`       | `go.mod`                            | generate, fmt, lint, test with coverage, matrix build   |
| `docker.yml`   | `docker/Dockerfile`                 | buildx build, tag with the newest git tag and `latest`  |
| `compose.yml`  | `compose.yml`, `docker-compose.yml` | up, down, build, pull, ps, logs                         |
| `release.yml`  | `.github/`, `docker/Dockerfile`     | cut a GitHub release and upload the binaries in `bin/`  |
| `mdox.yml`     | `.vale.ini`, `docs/`, `README.md`   | format markdown, run vale over the markdown and the code |
| `schema.yml`   | `schema/*.up.sql`                   | apply migrations, regenerate schema docs and Go models  |
| `tailwind.yml` | `tailwind.config.js`                | generate `tailwind.css`, or watch and regenerate        |

## What they expect on PATH

The skills call tools rather than vendoring them. Per skill:
`go.yml` needs [splint](https://github.com/titpetric/splint) and
[gotestsum](https://github.com/gotestyourself/gotestsum);
`docker.yml`, `compose.yml` and the vale step of `mdox.yml` need `docker`
with buildx; `release.yml` needs an authenticated
[gh](https://cli.github.com); `mdox.yml` needs
[mdox](https://github.com/bwplotka/mdox), which `mdox:setup` installs;
`schema.yml` needs [mig](https://github.com/go-bridget/mig);
`tailwind.yml` needs the `tailwindcss` CLI.

## The prose rules

`mdox:lint` runs vale at the repository root and lets it walk the tree.
The `.vale.ini` and the `llm-slop` style at the root of this repository
are the copy to start from: the style holds thirty-nine rules and the
config holds them against the markdown and the code together, `.go`,
`.php`, `.js`, `.css`, `.sql` and the rest, because vale reads the
comments and strings out of each.

Scanning only `docs/` reads the wrong file. Documentation is generated
from doc comments as often as it is written by hand, so a markdown-only
scan reports the generated page while the comment it came from goes
unread. The config excludes what nobody here writes: `vendor/`,
`node_modules/`, `testdata/`, `.claude/` and generated Go
(`*_templ.go`, `*.mig.go`, `*.pb.go`).

Copy both into a repository to use them:

```
cp -r titpetric/.vale.ini titpetric/.vale <your-repo>/
```

## Forking

`go.yml`, `docker.yml` and `release.yml` set `image: titpetric/${{name}}`,
where `name` is the directory name. Change that var to your own registry
namespace before running `docker:push` or `release:publish`.
