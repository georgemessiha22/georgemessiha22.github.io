# Resume

[![Release](https://github.com/georgemessiha22/georgemessiha22.github.io/actions/workflows/release.yml/badge.svg?branch=master)](https://github.com/georgemessiha22/georgemessiha22.github.io/actions/workflows/release.yml)

A single-source-of-truth resume: edit one YAML file and generate a Typst PDF, a
LaTeX (moderncv) PDF, a Markdown resume, and an HTML timeline mini-site
(published to GitHub Pages).

## How it works

```
resume.yaml ──> cmd/resume (Go) ──> Typst .typ ──> PDF
                                ──> LaTeX .tex ──> PDF
                                ──> Markdown .md
                                ──> site/ (index.html + css) ──> GitHub Pages
```

All resume content lives in `resume.yaml`. The Go app in `cmd/resume` renders it
into each format via independent engine packages (`internal/typst`,
`internal/latex`, `internal/markdown`, `internal/html`).

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
    renderer converts it (`*bold*` in Typst, `\textbf{}` in LaTeX, `**bold**` in
    Markdown, `<strong>` in HTML).
4. **Replace the photo** at `pictures/61673.jpg` (or change `contact.photo`).
5. **Set `releases_url`** (top-level in `resume.yaml`) to
   `https://github.com/<you>/<repo>/releases/latest/download` — the HTML site
   shows PDF download buttons pointing at your latest release. Omit it to hide
   the buttons.
6. **Set `site_url`** (top-level) to your published site, e.g.
   `https://<you>.github.io`. It appears as your website link in the PDFs.
7. **Add blog posts** (optional): create `blog/<slug>/index.md` with optional
   front-matter and images alongside it (see "Blogging" below). A `Blog` tab
   appears on the site automatically when at least one post exists.
8. **Generate locally** (see Usage). Commit and push to `master`; the release
   workflow builds the PDFs, attaches them to a GitHub Release, and deploys the
   site to Pages.
9. **Enable GitHub Pages**: repo Settings → Pages → Source = "GitHub Actions"
   (one-time step, required for the deploy job to publish the site).

## Blogging

The HTML site is a small blog. The résumé is the home page (`index.html`) and a
`Blog` tab lists your posts.

Add a post by creating a folder under `blog/` with an `index.md`:

```
blog/
  my-first-post/
    index.md
    diagram.png      # images (or any assets) referenced relatively
```

`index.md` supports optional YAML front-matter and standard Markdown (GFM),
including images:

```markdown
---
title: My First Post
date: "2026-06-30"
summary: A one-line description shown on the blog index.
---

# My First Post

![A diagram](diagram.png)

Regular **Markdown** — lists, code blocks, quotes, links — all work.
```

Posts are sorted newest-first by `date`. If you omit front-matter, the folder
name is used as the title. Regenerate with `make site` (or `go run ./cmd/resume
html`); each post becomes `blog/<slug>/index.html` with its images copied
alongside.

## Usage

### Prerequisites
- [Go](https://go.dev) 1.26+
- [Typst](https://github.com/typst/typst) 0.15+ (for the Typst PDFs)
- Docker (for the LaTeX PDFs)

### Generate sources only (no PDF compilation)
```bash
go run ./cmd/resume all                      # typst + tex + md (both variants) + html site
go run ./cmd/resume typst --variant detailed # one format, one variant
go run ./cmd/resume md                        # resume.md (summary)
go run ./cmd/resume html                     # site/ (summary)
```
Flags: `--input` (default `resume.yaml`), `--variant` (`summary`|`detailed`),
`--out` (default `build/gen`; the site goes to `site/`).

### Build PDFs, Markdown, and site (Makefile)
```bash
make typst       # Typst PDFs  -> dist/George_Messiha_*_v2.pdf
make resume      # LaTeX summary PDF (Docker) -> dist/George_Messiha_Resume.pdf
make detailed    # LaTeX detailed PDF (Docker) -> dist/George_Messiha_detailed_resume.pdf
make md          # Markdown     -> dist/George_Messiha_*.md
make site        # HTML mini-site -> site/
make all         # build everything (Docker image, all PDFs, Markdown, site), then clean logs
```
Generated sources land in `build/gen/`, PDFs/Markdown in `dist/`, and the site in
`site/` — all git-ignored.

### Run tests
```bash
go test ./...
```

## Project layout
- `resume.yaml` — your content (the only file you normally edit)
- `blog/<slug>/index.md` — blog posts (+ their images)
- `cmd/resume/` — CLI entrypoint
- `internal/model`, `internal/config` — data model + YAML loading
- `internal/render` — `Renderer` interface + shared helpers
- `internal/blog` — loads and renders Markdown blog posts
- `internal/typst`, `internal/latex`, `internal/markdown`, `internal/html` — output engines
- `typst/fonts/` — fonts vendored for reproducible Typst builds
- `moderncv/` — git submodule used by the LaTeX build
- Certificate/diploma scans under `certificates_awards/` and `Education/` are kept
  because the resume links to them.

## Thanks for checking my Resume
