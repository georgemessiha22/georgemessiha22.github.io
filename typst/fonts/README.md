# Vendored fonts

These fonts are bundled so the Typst resume builds reproducibly (locally and in
CI) without relying on system-installed fonts.

| Font | Used for | License |
| --- | --- | --- |
| Roboto | section/header font | Apache License 2.0 |
| Source Sans 3 | body font | SIL Open Font License 1.1 |
| Font Awesome 7 Free (Regular/Solid) + Brands | contact / social icons | SIL Open Font License 1.1 (icons) |

Sources:
- Roboto: https://github.com/googlefonts/roboto-2
- Source Sans 3 (3.052R): https://github.com/adobe-fonts/source-sans
- Font Awesome 7 Free (7.2.0): https://github.com/FortAwesome/Font-Awesome

Compile with `--font-path typst/fonts` (see `typst/README.md` and the Makefile
`typst*` targets).
