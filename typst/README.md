# Typst resume (parallel version)

A second, parallel implementation of the resume written in
[Typst](https://typst.app) using the
[`modern-cv`](https://typst.app/universe/package/modern-cv) template. The
original LaTeX (`moderncv`) version is unchanged and remains the primary source;
this directory adds a Typst version that produces the same two resumes.

## Outputs

Built into `../dist/` with a `_v2` suffix so they never collide with the LaTeX
PDFs:

- `George_Messiha_Resume_v2.pdf` — summary resume
- `George_Messiha_detailed_resume_v2.pdf` — detailed resume

## Layout

```
typst/
  resume.typ              # summary entry point
  detailed_resume.typ     # detailed entry point
  sections/               # one file per CV section (shared between both)
  fonts/                  # vendored fonts (see fonts/README.md)
  assets/profile.jpg      # profile picture
```

The two entry points differ only in which `sections/*.typ` files they
`#include`, mirroring the two LaTeX entry points.

## Build

Requires the [Typst CLI](https://github.com/typst/typst) (tested with 0.15.0).
Fonts are vendored, so `--font-path typst/fonts` is the only thing needed — no
system font installation.

From the repository root:

```bash
make typst            # build both
make typst-resume     # summary only
make typst-detailed   # detailed only
```

Or directly:

```bash
typst compile --font-path typst/fonts typst/resume.typ dist/George_Messiha_Resume_v2.pdf
typst compile --font-path typst/fonts typst/detailed_resume.typ dist/George_Messiha_detailed_resume_v2.pdf
```

`make all` also builds these alongside the LaTeX PDFs, and CI attaches the
`_v2.pdf` files to each GitHub release.

## Editing

- Update content in `sections/*.typ`.
- Update name / contact / photo in the `author` block of `resume.typ` and
  `detailed_resume.typ`.
- A harmless `unknown font family: source sans pro` warning may appear; it comes
  from the template's fallback font list. `Source Sans 3` is the font actually
  used and resolves correctly.
