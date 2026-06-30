# Resume

A single-source-of-truth resume: edit one YAML file and generate a Typst PDF, a
LaTeX (moderncv) PDF, and an HTML timeline mini-site (published to GitHub Pages).

## How it works

```
resume.yaml ──> cmd/resume (Go) ──> Typst .typ ──> PDF
                                ──> LaTeX .tex ──> PDF
                                ──> site/ (index.html + css) ──> GitHub Pages
```

All resume content lives in `resume.yaml`. The Go app in `cmd/resume` renders it
into each format via independent engine packages (`internal/typst`,
`internal/latex`, `internal/html`).

## Forking this repo

1. **Fork & clone with submodules** (the LaTeX build references the `moderncv`
   submodule):
   ```bash
   git clone --recurse-submodules https://github.com/<you>/<repo>.git
   ```
2. **Edit `resume.yaml`** — change `contact`, `experience`, `education`,
   `skills`, `languages`, `certificates`, `activities`. Tag entries with
   `variants: [summary, detailed]` (omit `variants` to include an entry in both
   resumes; use `[detailed]` for detailed-only).
3. **Use `**bold**`** inside `summary`, `intro`, and `bullets` for emphasis — each
   renderer converts it (`*bold*` in Typst, `\textbf{}` in LaTeX, `<strong>` in
   HTML).
4. **Replace the photo** at `pictures/61673.jpg` (or change `contact.photo`).
5. **Generate locally** (see Usage). Commit and push to `master`; the release
   workflow builds the PDFs, attaches them to a GitHub Release, and deploys the
   site to Pages.
6. **Enable GitHub Pages**: repo Settings → Pages → Source = "GitHub Actions"
   (one-time step, required for the deploy job to publish the site).

## Usage

### Prerequisites
- [Go](https://go.dev) 1.26+
- [Typst](https://github.com/typst/typst) 0.15+ (for the Typst PDFs)
- Docker (for the LaTeX PDFs)

### Generate sources only (no PDF compilation)
```bash
go run ./cmd/resume all                      # typst + tex (both variants) + html site
go run ./cmd/resume typst --variant detailed # one format, one variant
go run ./cmd/resume html                     # site/ (detailed)
```
Flags: `--input` (default `resume.yaml`), `--variant` (`summary`|`detailed`),
`--out` (default `build/gen`; the site goes to `site/`).

### Build PDFs and site (Makefile)
```bash
make typst       # Typst PDFs  -> dist/George_Messiha_*_v2.pdf
make resume      # LaTeX summary PDF (Docker) -> dist/George_Messiha_Resume.pdf
make detailed    # LaTeX detailed PDF (Docker) -> dist/George_Messiha_detailed_resume.pdf
make site        # HTML mini-site -> site/
make all         # build everything (Docker image, all PDFs, site), then clean logs
```
Generated sources land in `build/gen/`, PDFs in `dist/`, and the site in `site/`
— all git-ignored.

### Run tests
```bash
go test ./...
```

## Project layout
- `resume.yaml` — your content (the only file you normally edit)
- `cmd/resume/` — CLI entrypoint
- `internal/model`, `internal/config` — data model + YAML loading
- `internal/render` — `Renderer` interface + shared helpers
- `internal/typst`, `internal/latex`, `internal/html` — output engines
- `typst/fonts/` — fonts vendored for reproducible Typst builds
- `moderncv/` — git submodule used by the LaTeX build
- Certificate/diploma scans under `certificates_awards/` and `Education/` are kept
  because the resume links to them.

## Thanks for checking my Resume
