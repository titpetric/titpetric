# Diagrams

Diagrams are a common representation of:

- Execution flow (GitHub Actions)
- Dependency tracking (Go module graph)
- Software architecture (Go package graph)
- Coupling analyisis (Go code graph)
- Traffic monitoring (Flow analysis)
- Distributed request tracing

Some of these can be generated from a source of truth. For example, you
have a functioning database, migrations and schema - the schema is the
source of truth, and from that you can render relationship diagrams.

Guides for the typical options:

- [PlantUML](./plantuml.md)
- [MermaidJS](./mermaidjs.md)

For mechanical use cases, either being the source of truth, or generated
from the source of truth, PlantUML is the standard. It's good if the
plantuml could be rendered from markdown documentation directly, but in
case that's not possible certainly the preferred option is to render
`.svg` for any chart.

It's important to think if the documentation can go stale, and work
around it going stale with automated process whenever possible.

For workshopping diagrams, [draw.io](https://draw.io) is recommended.
It's a visual UML tool that allows some dragging and dropping to create
diagrams. It allows saving and loading those diagrams from `.drawio`
files, which can be maintained alongside documentation.

Since diagrams not generated from a source of truth will go stale, it's
always recommended to use mechanical means to generate them. With
planning-first efforts, UML being used as a source of truth means that
it will never go stale, and never be out of date.
