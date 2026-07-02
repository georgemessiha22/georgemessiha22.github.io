# Personal Projects Section (HTML site only) — Design

Date: 2026-07-02

## Goal

Add a "Personal Projects" callout to the résumé, shown **only on the GitHub
Pages HTML site** (not in the Typst/LaTeX PDFs or the Markdown output). It
appears immediately after the Professional Summary section and contains:

1. A short, enthusiastic intro paragraph about a love of customization,
   automation, and living in the terminal.
2. Two project links, each with a one-line description:
   - GogoNvim — https://github.com/georgemessiha22/GogoNvim
   - dotfiles — https://github.com/georgemessiha22/dotfiles

## Why HTML-only

The résumé is a single source of truth (`resume.yaml`). Each output format has
its own renderer (`internal/typst`, `internal/latex`, `internal/markdown`,
`internal/html`). By adding the new content as an optional field that **only the
HTML renderer reads**, it automatically appears only on GitHub Pages. No
variant/flag plumbing is needed — the other renderers simply ignore the field.

## Data model

New optional top-level field in `resume.yaml`:

```yaml
personal_projects:
  intro: >-
    Enthusiastic one-paragraph blurb about customization / automation / terminal.
  links:
    - name: GogoNvim
      url: https://github.com/georgemessiha22/GogoNvim
      description: One-line description.
    - name: dotfiles
      url: https://github.com/georgemessiha22/dotfiles
      description: One-line description.
```

Mirrored into the domain model:

```go
// PersonalProjects is an optional, HTML-only block of side projects.
type PersonalProjects struct {
    Intro string
    Links []ProjectLink
}

type ProjectLink struct {
    Name        string
    URL         string
    Description string
}
```

Added to `model.Resume` as `PersonalProjects PersonalProjects`. Because it is a
value type, an empty (zero) value means "not present" — the HTML template guards
on `len(.PersonalProjects.Links) > 0` (or a non-empty intro) before rendering.

`config.yamlResume` gets a matching DTO with yaml tags, converted in `Load`.

## Rendering (HTML only)

In `internal/html/templates/index.html.tmpl`, after the Professional Summary
`<section>`, add:

```html
{{- if .PersonalProjects.Links }}
<section class="projects">
  <h2>Personal Projects</h2>
  {{- if .PersonalProjects.Intro }}<p>{{ emph .PersonalProjects.Intro }}</p>{{ end }}
  <ul class="project-list">
  {{- range .PersonalProjects.Links }}
    <li><a href="{{ .URL }}">{{ .Name }}</a>{{ if .Description }} — {{ emph .Description }}{{ end }}</li>
  {{- end }}
  </ul>
</section>
{{- end }}
```

The `emph` helper (existing) supports `**bold**` in the intro and descriptions,
consistent with the rest of the site.

## Styling

Reuse existing look. Add minimal CSS to `style.css`:

```css
.project-list { margin: .4rem 0 0; padding-left: 1.1rem; }
.project-list li { margin: .25rem 0; }
```

The `h2` already has the accent underline treatment.

## Content (initial copy)

Intro (enthusiastic):
> I'm a bit obsessed with customization and automation — I live in the terminal
> and love tuning every tool until it feels like an extension of my hands. When
> I'm not shipping backend systems, I'm sharpening my own developer environment.

Links:
- **GogoNvim** — My from-scratch Neovim configuration, tuned for fast,
  keyboard-driven development.
- **dotfiles** — My portable terminal and system setup: shell, editor, and
  tooling, reproducible on any machine.

(Exact wording finalized in `resume.yaml`; can be tweaked freely later.)

## Testing

- `go test ./...` must pass (config loader + html renderer tests).
- Update `internal/html/html_test.go` / golden testdata if it asserts on the
  rendered index page, so the new section is covered.
- Update `internal/config/load_test.go` if it round-trips the full model.
- Manually verify the Typst/LaTeX/Markdown outputs are unchanged (they don't
  reference the field).

## Files touched

- `resume.yaml` — new `personal_projects` block.
- `internal/model/resume.go` — `PersonalProjects` + `ProjectLink` types, field on `Resume`.
- `internal/config/load.go` — DTO + mapping in `Load`.
- `internal/html/templates/index.html.tmpl` — new section after Summary.
- `internal/html/templates/style.css` — small list styling.
- Tests as needed.
