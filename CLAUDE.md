# titpetric

<!-- llm-rules:writing -->

## Writing

Every artifact fills a form. `titpetric/titpetric/llm-rules/writing.md` holds
them, and the output style of the same name loads it on every session. Read it
before writing a doc, a comment, a commit message or a pull request body.

The four tokens that fire most. Scan a reply for these before sending it:

| Token                                       | Fix                       |
|---------------------------------------------|---------------------------|
| `, so `                                     | Full stop                 |
| `, and it/this/the X <verb>`                | Full stop                 |
| `gives`, `provides`, `handles`, `carries`   | A verb that names the act |
| `which is what`, `which is why`             | End at the fact           |

A godoc opens on the symbol name and a verb from the closed list in
`writing.md`. One line. Lists take no Oxford comma. Punctuation is ASCII.

`./llm-rules/score.sh` counts findings per 1000 words over a diff. Run it
before reporting work.

<!-- /llm-rules:writing -->
