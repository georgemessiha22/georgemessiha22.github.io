# Resume

[![Release](https://github.com/georgemessiha22/georgemessiha22.github.io/actions/workflows/release.yml/badge.svg?branch=master)](https://github.com/georgemessiha22/georgemessiha22.github.io/actions/workflows/release.yml)

A single-source-of-truth resume: edit one YAML file and generate a Typst PDF, a
LaTeX (moderncv) PDF, a Markdown resume, and a modern **SvelteKit + shadcn-svelte**
site with a résumé home page, a Timeline page, and a blog (published to GitHub
Pages).

## How it works

```
resume.yaml ──> cmd/resume (Go) ──> Typst .typ ──> PDF
                                ──> LaTeX .tex ──> PDF
                                ──> Markdown .md
                                ──> resume.json ──┐
build/assets/blog/<slug>/index.md ──────────────────────────────┤
                                                    ▼
                                    build/assets/web (SvelteKit) ──> build/gen/site/ ──> GitHub Pages
```

Resume content lives in `resume.yaml`. The Go app in `cmd/resume` renders it into
the print formats (`internal/typst`, `internal/latex`, `internal/markdown`) and
emits `resume.json` (`internal/webdata`) — the data source for the web app. The
SvelteKit app in `build/assets/web/` reads that JSON, renders the résumé/timeline pages, and
turns `build/assets/blog/<slug>/index.md` posts into pages via **mdsvex**.

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
    Markdown, `<strong>` on the site).
4. **Replace the photo** at `build/assets/pictures/61673.jpg` (or change `contact.photo`).
5. **Set `releases_url`** (top-level in `resume.yaml`) to
   `https://github.com/<you>/<repo>/releases/latest/download` — the site's
   **Download Résumé** dropdown points at your latest release assets. Omit it to
   hide the dropdown.
6. **Set `site_url`** (top-level) to your published site, e.g.
   `https://<you>.github.io`. It appears as your website link in the PDFs.
7. **Add blog posts** (optional): create `build/assets/blog/<slug>/index.md` with optional
   front-matter and images alongside it (see "Blogging" below). A `Blog` tab is
   always present; the home page shows a card grid of the latest posts.
8. **Generate locally** (see Usage). Commit and push to `master`; the release
   workflow builds the PDFs, attaches them to a GitHub Release, and deploys the
   site to Pages.
9. **Enable GitHub Pages**: repo Settings → Pages → Source = "GitHub Actions"
   (one-time step, required for the deploy job to publish the site).

## Blogging

The site includes a blog. The résumé is the home page (`/`), a `Blog` tab lists
your posts, and the home page shows a card grid of the latest ones.

Add a post by creating a folder under `build/assets/blog/` with an `index.md`:

```
build/assets/blog/
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
summary: A one-line description shown on the blog cards.
---

# My First Post

![A diagram](diagram.png)

Regular **Markdown** — lists, code blocks, quotes, links — all work.
```

Posts are sorted newest-first by `date`. Front-matter is optional: if `title` is
omitted the first `#` heading (or the folder name) is used, and `summary` falls
back to the first paragraph. Start the body with a `# Heading` — the post page
renders the date and the Markdown body (which owns the title). Regenerate with
`make site` (or `npm --prefix web run build`).

## Usage

### Prerequisites
- [Go](https://go.dev) 1.26+
- [Node.js](https://nodejs.org) 20+ (for the SvelteKit site)
- [Typst](https://github.com/typst/typst) 0.15+ (for the Typst PDFs)
- Docker (for the LaTeX PDFs)

### Generate sources only (no PDF compilation)
```bash
go run ./cmd/resume all                      # typst + tex + md (both variants) + web JSON
go run ./cmd/resume typst --variant detailed # one format, one variant
go run ./cmd/resume md                        # resume.md (summary)
go run ./cmd/resume web                       # resume.json + profile.jpg into build/assets/web/
```
Flags: `--input` (default `resume.yaml`), `--variant` (`summary`|`detailed`),
`--out` (default `build/gen`; the web data goes to `build/assets/web/`).

### Run the site locally
```bash
cd build/assets/web
npm install
npm run dev      # dev server (regenerates resume.json + syncs blog first)
npm run build    # static build -> build/gen/site
npm run preview  # preview the production build
```

### Build PDFs, Markdown, and site (Makefile)
```bash
make typst       # Typst PDFs  -> dist/George_Messiha_*_v2.pdf
make resume      # LaTeX summary PDF (Docker) -> dist/George_Messiha_Resume.pdf
make detailed    # LaTeX detailed PDF (Docker) -> dist/George_Messiha_detailed_resume.pdf
make md          # Markdown     -> dist/George_Messiha_*.md
make site        # SvelteKit site -> build/gen/site
make all         # build everything (Docker image, all PDFs, Markdown, site), then clean logs
```
Generated sources land in `build/gen/`, PDFs/Markdown in `dist/`, and the site in
`build/gen/site/` — all git-ignored.

### Run tests
```bash
go test ./...
```

## Project layout
- `resume.yaml` — your content (the only file you normally edit)
- `resume.schema.json` — JSON Schema for `resume.yaml` (used by the yaml language server via the modeline on line 1; keep in sync with the Go model)
- `build/assets/blog/<slug>/index.md` — blog posts (+ their images)
- `cmd/resume/` — CLI entrypoint
- `internal/model`, `internal/config` — data model + YAML loading
- `internal/render` — `Renderer` interface + shared helpers
- `internal/typst`, `internal/latex`, `internal/markdown` — print output engines
- `internal/webdata` — emits `resume.json` for the web app
- `build/assets/web/` — SvelteKit + shadcn-svelte site (résumé, timeline, blog via mdsvex)
- `build/assets/typst/fonts/` — fonts vendored for reproducible Typst builds
- `moderncv/` — git submodule used by the LaTeX build
- Certificate/diploma scans under `build/assets/certificates_awards/` and `build/assets/experience/` are kept
  because the resume links to them.

## Thanks for checking my Resume
