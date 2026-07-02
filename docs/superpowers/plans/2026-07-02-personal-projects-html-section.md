# Personal Projects Section (HTML-only) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a "Personal Projects" section (enthusiastic intro + two GitHub links) that renders only on the GitHub Pages HTML site, right after the Professional Summary.

**Architecture:** Add an optional `PersonalProjects` value to the domain model, populate it from a new `personal_projects` block in `resume.yaml` via the config loader, and render it exclusively in the HTML `index.html.tmpl` template. Each link carries an `icon` keyword that the HTML template maps to an inline SVG glyph. The Typst/LaTeX/Markdown renderers never reference the field, so it appears on GitHub Pages only.

**Tech Stack:** Go 1.26, `html/template`, `gopkg.in/yaml.v3`, golden-file tests, inline SVG.

---

### Task 1: Add PersonalProjects to the domain model

**Files:**
- Modify: `internal/model/resume.go`

- [ ] **Step 1: Add the new types and the field on `Resume`**

In `internal/model/resume.go`, add these types after the `Cert` type (before `type Resume struct`):

```go
// ProjectLink is a single personal/side project link (HTML site only).
type ProjectLink struct {
	Name        string
	URL         string
	Description string
	Icon        string // icon keyword: "github", "neovim", "terminal" (empty => "github")
}

// PersonalProjects is an optional block of side projects shown only on the
// HTML site (GitHub Pages). It is not rendered in the PDF or Markdown outputs.
type PersonalProjects struct {
	Intro string
	Links []ProjectLink
}
```

Then add the field to the `Resume` struct. Insert it immediately after the `SiteURL string` field, inside the struct:

```go
	// PersonalProjects is an optional block of side-project links rendered only
	// on the HTML site (GitHub Pages). Empty => the section is omitted. Optional.
	PersonalProjects PersonalProjects
```

- [ ] **Step 2: Verify it compiles**

Run: `go build ./...`
Expected: builds with no errors.

- [ ] **Step 3: Commit**

```bash
git add internal/model/resume.go
git commit -m "model: add optional PersonalProjects block (HTML-only)"
```

---

### Task 2: Parse `personal_projects` in the config loader

**Files:**
- Modify: `internal/config/load.go`
- Modify: `internal/config/testdata/valid.yaml`
- Test: `internal/config/load_test.go`

- [ ] **Step 1: Add the block to the test fixture YAML**

Append to the end of `internal/config/testdata/valid.yaml`:

```yaml
personal_projects:
  intro: I love **automation**.
  links:
    - name: GogoNvim
      url: https://github.com/georgemessiha22/GogoNvim
      description: My **Neovim** config.
      icon: neovim
    - name: dotfiles
      url: https://github.com/georgemessiha22/dotfiles
      description: My terminal setup.
      icon: terminal
```

- [ ] **Step 2: Write the failing test**

Add this test function to `internal/config/load_test.go`:

```go
func TestLoad_PersonalProjects(t *testing.T) {
	r, err := Load("testdata/valid.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.PersonalProjects.Intro != "I love **automation**." {
		t.Fatalf("intro not parsed: %q", r.PersonalProjects.Intro)
	}
	if len(r.PersonalProjects.Links) != 2 {
		t.Fatalf("expected 2 project links, got %d", len(r.PersonalProjects.Links))
	}
	l0 := r.PersonalProjects.Links[0]
	if l0.Name != "GogoNvim" || l0.URL != "https://github.com/georgemessiha22/GogoNvim" || l0.Description != "My **Neovim** config." {
		t.Fatalf("first link not parsed: %+v", l0)
	}
	if l0.Icon != "neovim" {
		t.Fatalf("first link icon not parsed: %q", l0.Icon)
	}
	if r.PersonalProjects.Links[1].Name != "dotfiles" {
		t.Fatalf("second link not parsed: %+v", r.PersonalProjects.Links[1])
	}
}
```

- [ ] **Step 3: Run the test to verify it fails**

Run: `go test ./internal/config/ -run TestLoad_PersonalProjects -v`
Expected: FAIL — intro is empty / 0 links parsed (field not yet mapped).

- [ ] **Step 4: Add the DTO to `yamlResume`**

In `internal/config/load.go`, add a field to the `yamlResume` struct immediately after the `SiteURL string \`yaml:"site_url"\`` line:

```go
	PersonalProjects yamlPersonalProjects `yaml:"personal_projects"`
```

- [ ] **Step 5: Add the DTO types**

In `internal/config/load.go`, add these types after the `yamlCert` struct definition:

```go
type yamlPersonalProjects struct {
	Intro string            `yaml:"intro"`
	Links []yamlProjectLink `yaml:"links"`
}

type yamlProjectLink struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	Description string `yaml:"description"`
	Icon        string `yaml:"icon"`
}
```

- [ ] **Step 6: Map the DTO into the model in `Load`**

In `internal/config/load.go`, inside `Load`, add the field to the `model.Resume{...}` composite literal, immediately after the `SiteURL: y.SiteURL,` line:

```go
		PersonalProjects: model.PersonalProjects{
			Intro: y.PersonalProjects.Intro,
			Links: toProjectLinks(y.PersonalProjects.Links),
		},
```

- [ ] **Step 7: Add the converter helper**

In `internal/config/load.go`, add this function after `toEntries`:

```go
func toProjectLinks(in []yamlProjectLink) []model.ProjectLink {
	if len(in) == 0 {
		return nil
	}
	out := make([]model.ProjectLink, 0, len(in))
	for _, l := range in {
		out = append(out, model.ProjectLink{Name: l.Name, URL: l.URL, Description: l.Description, Icon: l.Icon})
	}
	return out
}
```

- [ ] **Step 8: Run the test to verify it passes**

Run: `go test ./internal/config/ -run TestLoad_PersonalProjects -v`
Expected: PASS.

- [ ] **Step 9: Run the full config package tests (ensure nothing else broke)**

Run: `go test ./internal/config/`
Expected: ok.

- [ ] **Step 10: Commit**

```bash
git add internal/config/load.go internal/config/load_test.go internal/config/testdata/valid.yaml
git commit -m "config: parse personal_projects block into model"
```

---

### Task 3: Render the section (with SVG icons) in the HTML template + CSS

**Files:**
- Modify: `internal/html/html.go`
- Modify: `internal/html/templates/index.html.tmpl`
- Modify: `internal/html/templates/style.css`
- Modify: `internal/html/html_test.go`
- Modify (golden, regenerated): `internal/html/testdata/index.html.golden`

- [ ] **Step 1: Add PersonalProjects to the html test fixture**

In `internal/html/html_test.go`, inside the `fixture()` function's `model.Resume{...}` literal, add this field immediately after the `Summary:` line:

```go
		PersonalProjects: model.PersonalProjects{
			Intro: "I love **automation** and the terminal.",
			Links: []model.ProjectLink{
				{Name: "GogoNvim", URL: "https://github.com/georgemessiha22/GogoNvim", Description: "My **Neovim** config.", Icon: "neovim"},
				{Name: "dotfiles", URL: "https://github.com/georgemessiha22/dotfiles", Description: "My terminal setup.", Icon: "terminal"},
			},
		},
```

- [ ] **Step 2: Add assertions to `TestRenderMiniSite`**

In `internal/html/html_test.go`, in `TestRenderMiniSite`, extend the existing `for _, want := range []string{ ... }` slice (the one that already checks the download links) by adding these entries to it:

```go
		`class="projects"`,
		`class="project-list"`,
		`<svg`,
		"https://github.com/georgemessiha22/GogoNvim",
		"https://github.com/georgemessiha22/dotfiles",
		"<strong>Neovim</strong>",
```

- [ ] **Step 3: Run the test to verify it fails**

Run: `go test ./internal/html/ -run TestRenderMiniSite -v`
Expected: FAIL — index.html differs from golden AND/OR missing `class="projects"` (template not updated yet).

- [ ] **Step 4: Add the `projectIcon` template helper in `html.go`**

In `internal/html/html.go`, add this function immediately after the `emphasize` function (before `func newTemplate`):

```go
// projectIcon returns a small inline SVG glyph for a project link. The keyword
// selects the icon; an unknown or empty keyword falls back to the GitHub mark.
// All icons use currentColor so CSS controls their color. HTML-site use only.
func projectIcon(kind string) template.HTML {
	const github = `<svg class="pi" viewBox="0 0 16 16" width="16" height="16" aria-hidden="true" fill="currentColor"><path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0 0 16 8c0-4.42-3.58-8-8-8z"/></svg>`
	const neovim = `<svg class="pi" viewBox="0 0 16 16" width="16" height="16" aria-hidden="true" fill="currentColor"><path d="M2 2v12l3-3V5l6 8h3V2h-2v9L4 2H2z"/></svg>`
	const terminal = `<svg class="pi" viewBox="0 0 16 16" width="16" height="16" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="1.5" y="2.5" width="13" height="11" rx="1.5"/><path d="M4 6l2.5 2L4 10"/><path d="M8 10.5h4"/></svg>`
	switch kind {
	case "neovim":
		return template.HTML(neovim)
	case "terminal":
		return template.HTML(terminal)
	default:
		return template.HTML(github)
	}
}
```

- [ ] **Step 5: Register the helper in `newTemplate`**

In `internal/html/html.go`, in `newTemplate`, add `projectIcon` to the `template.FuncMap`. Change the FuncMap so it reads:

```go
	return template.New(name).Funcs(template.FuncMap{
		"emph": emphasize,
		// safeHTML embeds already-rendered HTML (goldmark post bodies).
		"safeHTML":    func(s string) template.HTML { return template.HTML(s) },
		"projectIcon": projectIcon,
	}).ParseFS(assetFS, "templates/"+name)
```

- [ ] **Step 6: Verify it still compiles**

Run: `go build ./...`
Expected: builds with no errors.

- [ ] **Step 7: Add the section to the template**

In `internal/html/templates/index.html.tmpl`, locate the Professional Summary section:

```html
<section>
  <h2>Professional Summary</h2>
  <p>{{ emph .Summary }}</p>
</section>
```

Immediately after that closing `</section>` (and before the `<section>` with `<h2>Experience</h2>`), insert:

```html

{{- if .PersonalProjects.Links }}
<section class="projects">
  <h2>Personal Projects</h2>
  {{- if .PersonalProjects.Intro }}<p>{{ emph .PersonalProjects.Intro }}</p>{{ end }}
  <ul class="project-list">
  {{- range .PersonalProjects.Links }}
    <li><span class="pi-wrap">{{ projectIcon .Icon }}</span><a href="{{ .URL }}">{{ .Name }}</a>{{ if .Description }} — {{ emph .Description }}{{ end }}</li>
  {{- end }}
  </ul>
</section>
{{- end }}
```

- [ ] **Step 8: Add CSS for the project list and icons**

In `internal/html/templates/style.css`, add these rules immediately after the line `.skills .cat { font-weight: 600; }`:

```css
.project-list { list-style: none; margin: .4rem 0 0; padding-left: 0; }
.project-list li { display: flex; align-items: baseline; gap: .5rem; margin: .35rem 0; }
.pi-wrap { flex: none; color: var(--accent); line-height: 1; transform: translateY(2px); }
.pi { display: block; }
```

- [ ] **Step 9: Regenerate the golden file**

Run: `go test ./internal/html/ -run TestRenderMiniSite -update`
Expected: PASS (golden file rewritten to include the new section).

- [ ] **Step 10: Verify the golden file contains the new section and an icon**

Run: `grep -c 'class="projects"' internal/html/testdata/index.html.golden && grep -c '<svg' internal/html/testdata/index.html.golden`
Expected: two lines — first `1`, second `2` (one SVG per link).

- [ ] **Step 11: Re-run html tests without -update to confirm stability**

Run: `go test ./internal/html/`
Expected: ok (golden matches, all assertions pass).

- [ ] **Step 12: Commit**

```bash
git add internal/html/html.go internal/html/templates/index.html.tmpl internal/html/templates/style.css internal/html/html_test.go internal/html/testdata/index.html.golden
git commit -m "html: render Personal Projects section with inline SVG icons"
```

---

### Task 4: Add real content to resume.yaml

**Files:**
- Modify: `resume.yaml`

- [ ] **Step 1: Add the `personal_projects` block**

In `resume.yaml`, add the following block immediately after the `summary:` block (after line ending the summary, before `experience:`). Keep the two-space indentation style used in the file:

```yaml
# Personal side projects — rendered ONLY on the HTML site (GitHub Pages).
# The PDF and Markdown outputs ignore this block. Each link's `icon` selects an
# inline SVG glyph: github | neovim | terminal (default: github).
personal_projects:
  intro: >-
    I'm a bit obsessed with **customization** and **automation** — I live in the
    terminal and love tuning every tool until it feels like an extension of my
    hands. When I'm not shipping backend systems, I'm sharpening my own developer
    environment.
  links:
    - name: GogoNvim
      url: https://github.com/georgemessiha22/GogoNvim
      description: My from-scratch **Neovim** configuration, tuned for fast, keyboard-driven development.
      icon: neovim
    - name: dotfiles
      url: https://github.com/georgemessiha22/dotfiles
      description: My portable terminal and system setup — shell, editor, and tooling, reproducible on any machine.
      icon: terminal
```

- [ ] **Step 2: Generate the site and confirm the section renders**

Run: `go run ./cmd/resume html`
Expected: exits 0, writes `site/`.

- [ ] **Step 3: Confirm the rendered index has the section and links**

Run: `grep -c 'Personal Projects' site/index.html && grep -c 'github.com/georgemessiha22/GogoNvim' site/index.html && grep -c 'github.com/georgemessiha22/dotfiles' site/index.html`
Expected: three lines, each `1`.

- [ ] **Step 4: Confirm the section does NOT leak into other formats**

Run: `go run ./cmd/resume all` then `grep -rl "GogoNvim" build/gen/ || echo "NOT FOUND (correct)"`
Expected: `NOT FOUND (correct)` — the Typst/LaTeX/Markdown generated sources must not contain the project links.

- [ ] **Step 5: Commit**

```bash
git add resume.yaml
git commit -m "content: add personal projects (GogoNvim, dotfiles) to HTML site"
```

---

### Task 5: Full verification

- [ ] **Step 1: Run the entire test suite**

Run: `go test ./...`
Expected: all packages `ok`.

- [ ] **Step 2: Vet**

Run: `go vet ./...`
Expected: no output (clean).

- [ ] **Step 3: Confirm working tree is clean**

Run: `git status --porcelain`
Expected: empty (everything committed). Note: `site/`, `build/gen/`, and `dist/` are git-ignored, so their generated files should not appear.

---

## Self-Review Notes

- **Spec coverage:** Model field + Icon (Task 1), YAML parsing (Task 2), HTML-only rendering with SVG icons + CSS (Task 3), real content with icons (Task 4), tests + isolation check (Tasks 2, 3, 5). All spec sections covered.
- **Type consistency:** `PersonalProjects{Intro, Links}` and `ProjectLink{Name, URL, Description, Icon}` are used identically across model, loader, template, and tests. YAML DTOs (`yamlPersonalProjects`, `yamlProjectLink`) map field-for-field including `icon`. The `projectIcon` helper is defined in `html.go` (Task 3 Step 4) and registered (Step 5) before the template uses it (Step 7).
- **HTML-only guarantee:** Verified explicitly in Task 4 Step 4 by grepping generated non-HTML sources for "GogoNvim". The `projectIcon` helper lives in the `html` package only.
- **Icons:** Inline SVG with `currentColor` + accent color via `.pi-wrap`; unknown/empty keyword falls back to the GitHub mark (no broken output).
- **No placeholders:** every code step includes the actual code and exact commands with expected output.
