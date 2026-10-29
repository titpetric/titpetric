`atkins mdox:fmt` before `atkins mdox:lint`, always. The formatter re-pads every table cell, so a page with a hand-written table fails the check until it has been through it.

The vale half of `mdox:lint` walks the whole tree, not `docs/`. vale reads the comments and strings out of Go, PHP, JavaScript, CSS, SQL and the rest, and the sections in `.vale.ini` decide which files it holds to the prose rules. A page under `docs/` is often generated from a doc comment, so a scan that read only the markdown would report the copy and never the source.
