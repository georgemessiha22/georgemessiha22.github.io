# Design: `resume` Go app — YAML-driven multi-format resume generator

Date: 2026-06-30
Status: Approved

## Goal

A Go CLI in the repository root that reads a single `resume.yaml` (the new
**source of truth**) and renders the resume into three formats:

- **Typst** (`.typ`) — compiled to PDF later (summary + detailed)
- **LaTeX** (`.tex`, `moderncv`) — compiled to PDF later (summary + detailed)
- **HTML** — a self-contained GitHub Pages mini-site with a timeline layout
  (detailed variant), deployed to Pages on release

PDF compilation stays in the Makefile/CI (Go only generates source). The existing
hand-written `.typ`/`.tex` files are replaced by generated output; the `moderncv`
git submodule and the vendored Typst fonts are kept because they are still needed
to compile.

## Decisions (agreed with user)

1. **YAML is the new source of truth.** Go generates `.typ` and `.tex` for both
   variants; obsolete hand-written source files are removed.
2. **One YAML, variant-tagged entries.** Each experience/education/cert entry may
   carry `variants: [summary, detailed]`; untagged entries appear in both. The
   generator filters by variant.
3. **Ports-and-adapters architecture.** A `Renderer` interface decouples the
   domain model from output formats. **Each render engine lives in its own
   package** (`internal/typst`, `internal/latex`, `internal/html`) for strong
   separation of concerns.
4. **CLI with subcommands** generating source files only; PDF compilation remains
   in the Makefile.
5. **HTML = GitHub Pages mini-site** (detailed variant), deployed to Pages on
   release; the PDFs still attach to the release.
6. **stdlib + `gopkg.in/yaml.v3`**, templates embedded with `//go:embed`. No
   heavy CLI/template frameworks.

## Architecture & layout

```
go.mod                          # module github.com/georgemessiha22/georgemessiha22
resume.yaml                     # single source of truth (variant-tagged)
cmd/resume/main.go              # CLI entrypoint, flag parsing, wiring only
internal/
  model/
    resume.go                   # domain structs + Variant + ForVariant(v) filter (no format knowledge)
  config/
    load.go                     # YAML -> model, validation, defaults (only YAML parser)
  render/
    render.go                   # Renderer interface, Format, Artifact, inline-emphasis helper
  typst/
    typst.go                    # TypstRenderer
    templates/resume.typ.tmpl
  latex/
    latex.go                    # LatexRenderer
    templates/resume.tex.tmpl
  html/
    html.go                     # HTMLRenderer (mini-site)
    templates/index.html.tmpl
    templates/style.css
  writer/
    writer.go                   # writes []Artifact to disk (mkdir -p, multi-file site)
```

Dependency rule: `typst`, `latex`, `html` each import `render` and `model` but
**never each other**. `cmd/resume` wires everything via `render.Renderer`.

### Renderer interface

```go
package render

type Variant = model.Variant

type Artifact struct {
    RelPath string // e.g. "resume.typ" or "index.html"
    Bytes   []byte
}

type Renderer interface {
    Render(r model.Resume, v model.Variant) ([]Artifact, error)
}
```

Returning `[]Artifact` lets the HTML engine emit multiple files
(index.html + style.css + profile image) while Typst/LaTeX emit one.

## Data model & YAML

```yaml
contact:
  firstname: George
  lastname: Messiha
  title: Lead Software Engineer
  email: georgemessiha22@gmail.com
  phone: "(+971) 54 555 1032"
  location: "Dubai, UAE"
  photo: pictures/61673.jpg
  socials:
    github: georgemessiha22
    gitlab: georgemessiha22
    linkedin: georgemessiha22

summary: >
  SWE with over 9 years of experience ...

experience:
  - role: Senior Software Engineer
    org: HungerStation DeliveryHero Gmbh
    org_url: https://hungerstation.com
    location: Dubai, UAE
    mode: Hybrid                    # rendered as "Dubai, UAE · Hybrid"
    start: Nov 2023
    end: Present
    variants: [summary, detailed]
    intro: Engineered and maintained high-throughput microservices ...
    bullets:
      - "Architected ... new core service in **Go** ..."
      - "Drove ... 60% reduction in build times ..."

education:    [ ... ]               # title, org, location, start, end, url, note, variants
skills:
  - { category: Languages, items: [Go, Python, Shell] }
languages:
  - { name: Arabic, level: Native }
  - { name: English, level: Fluent }
certificates:                       # title, org, org_url, cert_url, location, start, end, group, variants
  - { title: Foundation of Project Management, org: Google, variants: [summary, detailed] }
activities:   [ ... ]               # detailed-only by tag
```

Model (`internal/model`):

```go
type Variant string // "summary" | "detailed"

type Resume struct {
    Contact      Contact
    Summary      string
    Experience   []Entry
    Education    []Entry
    Skills       []SkillGroup
    Languages    []Language
    Certificates []Cert
    Activities   []Entry
}

func (r Resume) ForVariant(v Variant) Resume // returns a filtered copy
```

Filtering rule: an entry is included when its `variants` list contains `v`, or
when it has no `variants` (defaults to both).

### Inline emphasis

Bullets/intros use Markdown-ish `**bold**`. A small shared parser in `render`
splits text into plain/bold spans so every engine stays consistent:

- Typst → `*bold*`
- LaTeX → `\textbf{bold}`
- HTML → `<strong>bold</strong>`

Each engine is also responsible for escaping its own special characters (LaTeX
`% & _ # { }`, Typst `# @ \\ $`, HTML `< > &`). Links in titles are expressed as
explicit `url`/`org_url`/`cert_url` fields (never inline), matching the current
files.

## CLI, output locations, compilation

CLI (stdlib `flag`):

```
resume typst [--variant summary|detailed] [--input resume.yaml] [--out DIR]
resume tex   [--variant summary|detailed] [--input resume.yaml] [--out DIR]
resume html  [--variant detailed]          [--input resume.yaml] [--out DIR]
resume all   [--input resume.yaml]   # both variants for typst+tex, plus html site
```

Defaults: `--input resume.yaml`, `--out build/gen` (site goes to `site/`).
Invalid YAML/validation → message on stderr, exit 1.

Generated output:

```
build/gen/
  resume.typ            resume_detailed.typ
  resume.tex            resume_detailed.tex
site/
  index.html  style.css  profile.jpg
```

Generated Typst/TeX are **single self-contained files** (the generator inlines
all sections — no `#include`/`\subfile`). They still
`#import "@preview/modern-cv:0.10.0"` / `\documentclass{moderncv}` and rely on the
vendored fonts and the `moderncv` submodule already in the repo.

Compilation (Makefile, unchanged division of labor — Go generates, external tools
compile):

- `make resume`   → `go run ./cmd/resume tex --variant summary` → pdflatex (Docker) → `George_Messiha_Resume.pdf`
- `make detailed` → tex detailed → `George_Messiha_detailed_resume.pdf`
- `make typst`    → `go run ./cmd/resume typst` (both) → `typst compile` → `*_v2.pdf`
- `make site`     → `go run ./cmd/resume html` → `site/`

## Pipeline, cleanup, testing

**Removed** (superseded by YAML): `typst/sections/*.typ`, the two hand-written
`typst/*.typ`, and all LaTeX content/section `*.tex` files.
**Kept:** `moderncv` submodule, `typst/fonts/`, `pictures/`, the PDF output names.

**CI (`release.yml`):**
- Add `actions/setup-go`.
- Build PDFs through the Go-driven Make targets (Go generates source, then
  Typst/pdflatex compile).
- Attach the 4 PDFs (`_Resume`, `_detailed_resume`, `_Resume_v2`,
  `_detailed_resume_v2`) to the release.
- Add a **GitHub Pages deploy job** that runs `make site` and publishes `site/`
  via `actions/upload-pages-artifact` + `actions/deploy-pages` (needs `pages:
  write`, `id-token: write` permissions).

**Testing:**
- `internal/model`: variant-filter unit tests.
- `internal/config`: YAML load + validation tests (good + malformed fixtures).
- `internal/typst`, `internal/latex`, `internal/html`: golden-file tests
  comparing a fixture model's rendered output to checked-in
  `testdata/*.golden`.
- `go test ./...` runs locally and in CI.

**Error handling:** `config` validates required fields (name, ≥1 experience) and
returns wrapped errors; renderers return errors and never panic; `cmd` prints to
stderr and exits non-zero.

## Risks

- **Content fidelity during YAML migration** — all existing entries/links must be
  transcribed faithfully into `resume.yaml`. Mitigated by golden tests + visual
  PDF/HTML check after first generation.
- **Escaping bugs** across three target syntaxes — covered by golden tests with
  entries containing `% & _ # { } < >`.
- **Pages first-enable** — repo must have Pages enabled (source = GitHub Actions);
  noted in docs.
- **Two PDF toolchains remain** (Typst + LaTeX/Docker) — accepted; both were
  requested.
