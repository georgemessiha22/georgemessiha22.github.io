# Design: Parallel Typst Version of the Resume

Date: 2026-06-30
Status: Approved

## Goal

Add a buildable **Typst** version of both resumes alongside the existing LaTeX
sources. The existing LaTeX (`moderncv`), Docker build, and CI release pipeline
remain **untouched and fully functional**. Typst becomes a second, parallel way
to produce the two resume PDFs.

## Context

The project is a `moderncv`-based LaTeX CV (banking/grey theme) with two
top-level entry points assembled from `\subfile{}` fragments:

- `George_Messiha_Resume.tex` — **summary** resume
- `George_Messiha_detailed_resume.tex` — **detailed** resume (adds Huspy,
  Nexus Edge, internships, Activities, and the detailed certificates list)

Build today: `Makefile` → Docker (`pdflatex`) → `dist/`, released as PDFs by
`.github/workflows/release.yml`.

## Decisions (agreed with user)

1. **Parallel version** — LaTeX stays the source of truth; Typst is added
   side-by-side. Nothing in the LaTeX tree is modified or deleted.
2. **Template-based** — use the `modern-cv` Typst package (a port of
   Awesome-CV). The look is clean/modern, **not** a pixel match of moderncv's
   banking style. Accepted tradeoff.
3. **Fonts vendored** — Roboto, Source Sans Pro (Source Sans 3), and FontAwesome
   `.otf` files are committed under `typst/fonts/` so CLI + CI builds are
   reproducible without system font installs.
4. **Build + CI** — add Makefile targets and a CI step so the Typst PDFs are also
   attached to the GitHub release.
5. **Output naming** — Typst PDFs use a `_v2` suffix to avoid colliding with the
   LaTeX outputs:
   - `George_Messiha_Resume_v2.pdf` (summary)
   - `George_Messiha_detailed_resume_v2.pdf` (detailed)

## Package API (verified by probing `modern-cv:0.10.0`)

- `resume.with(author: (...), profile-picture: image(...), positions: ..., ...)`
- `author` supports: `firstname, lastname, email, homepage, phone, github,
  gitlab, linkedin, address, positions` — all present in `contact/main.tex`.
- `resume-entry(title:, location:, date:, description:, title-link:)` ↔ `\cventry`
- `resume-item[ - bullet ... ]` ↔ `\begin{itemize}\item`
- `resume-skill-item("Category", (..values..))` ↔ `\cvitem{label}{text}`
- `#include "sections/X.typ"` ↔ `\subfile{...}`

## File structure (all new; nothing existing is touched)

```
typst/
  resume.typ                    # summary entry point  -> _v2 summary PDF
  detailed_resume.typ           # detailed entry point -> _v2 detailed PDF
  sections/
    summary.typ
    experience_summary.typ      # entries used by the summary resume
    experience_detailed.typ     # entries used by the detailed resume (+huspy, nexus, internships)
    education.typ
    skills.typ
    languages.typ
    certificates_summary.typ
    certificates_detailed.typ
    activities.typ
  fonts/                        # vendored Roboto, Source Sans 3, FontAwesome .otf
  assets/
    profile.jpg                 # copy of pictures/61673.jpg
  README.md                     # Typst usage notes
```

The two entry points differ only in which section files they `#include`, mirroring
the two LaTeX entry points exactly.

## Content mapping (LaTeX → Typst)

| moderncv | modern-cv |
| --- | --- |
| contact macros + `\photo` | `resume.with(author, profile-picture)` |
| `\section{X}` | `= X` |
| `\cventry{date}{title}{org}{loc}{grade}{desc}` | `resume-entry(...)` + `resume-item[...]` |
| `\cvitem{label}{text}` (skills) | `resume-skill-item("label", (...))` |
| `\cvdoubleitem{a}{b}{c}{d}` (languages) | `resume-skill-item` rows |
| `\href{url}{txt}` / `\httplink[txt]{url}` | `#link("url")[txt]` |
| `\textbf{}` / `\textsc{}` | `*bold*` / `smallcaps()` |
| `\subfile{...}` | `#include "sections/...typ"` |

Notes:
- The `grade` field of `\cventry` (e.g. "Hybrid", "Remote", "Fulltime") has no
  direct `resume-entry` slot; it is appended to `location` (e.g. "Dubai, UAE ·
  Hybrid") so no information is lost.
- Commented-out LaTeX lines (`% \item Tools: ...`) are intentionally omitted, as
  they are in the current LaTeX output.
- `\href` targets missing a scheme (e.g. `github.com/...`) get `https://`
  prepended so links are valid in Typst.

## Build wiring

`Makefile` — append (do not alter existing targets):

```
typst-resume:
	typst compile --font-path typst/fonts typst/resume.typ dist/George_Messiha_Resume_v2.pdf

typst-detailed:
	typst compile --font-path typst/fonts typst/detailed_resume.typ dist/George_Messiha_detailed_resume_v2.pdf

typst: typst-resume typst-detailed
```

`all` is extended to also build the Typst PDFs.

CI (`release.yml`) — add a step that installs Typst and runs `make typst`, then
add the two `_v2.pdf` files to the release `files:` list.

## Verification

- `typst compile` succeeds for both entry points (Typst 0.15.0, local).
- Two `_v2.pdf` files are produced in `dist/`.
- Spot-check: every section, every `\cventry`, and all links from the LaTeX
  appear in the Typst output.
- `make detailed` / `make resume` (LaTeX) are unaffected (not run locally without
  Docker, but no LaTeX/Make existing targets are modified).

## Risks

- **modern-cv look differs from moderncv banking** — accepted per decision 2.
- **Font licensing in-repo** — Roboto/Source Sans/FontAwesome are OFL/Apache,
  redistributable. License notes added to `typst/fonts/`.
- **CI Typst install** — uses a maintained GitHub Action
  (`typst-community/setup-typst`) pinned to a version.
