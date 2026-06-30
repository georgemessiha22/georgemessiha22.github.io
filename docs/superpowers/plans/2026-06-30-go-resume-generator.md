# Go Resume Generator Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go CLI that reads one `resume.yaml` and generates Typst, LaTeX, and an HTML GitHub Pages mini-site, then wire the build/release pipeline to use it.

**Architecture:** Ports-and-adapters. A `model` package holds the domain data and variant filtering; `config` parses YAML into the model; a `render` package defines the `Renderer` interface plus shared helpers; each output engine lives in its own package (`typst`, `latex`, `html`) and implements `Renderer`; a `writer` persists artifacts; `cmd/resume` wires them via the interface. Templates are embedded with `//go:embed`. PDF compilation stays in the Makefile.

**Tech Stack:** Go 1.26 (stdlib `flag`, `text/template`, `html/template`, `embed`, `testing`), `gopkg.in/yaml.v3`. Existing Typst (`modern-cv`) + LaTeX (`moderncv`) toolchains compile the generated source.

---

## File Structure

| Path | Responsibility |
| --- | --- |
| `go.mod` / `go.sum` | Module `github.com/georgemessiha22/georgemessiha22`, dep on yaml.v3 |
| `cmd/resume/main.go` | CLI: parse subcommand + flags, load config, run renderer, write |
| `internal/model/resume.go` | Domain structs, `Variant`, `ForVariant` filtering, no format knowledge |
| `internal/model/resume_test.go` | Variant-filter unit tests |
| `internal/config/load.go` | YAML → model, defaults, validation (only YAML parser) |
| `internal/config/load_test.go` | Load + validation tests |
| `internal/config/testdata/*.yaml` | Good + malformed fixtures |
| `internal/render/render.go` | `Renderer` interface, `Artifact`, `Variant` alias, inline-emphasis + helpers |
| `internal/render/render_test.go` | Inline-emphasis parser tests |
| `internal/typst/typst.go` | `TypstRenderer` + embedded template + Typst escaping |
| `internal/typst/templates/resume.typ.tmpl` | Typst template |
| `internal/typst/typst_test.go` | Golden test |
| `internal/typst/testdata/*.golden` | Expected Typst output |
| `internal/latex/latex.go` | `LatexRenderer` + embedded template + LaTeX escaping |
| `internal/latex/templates/resume.tex.tmpl` | LaTeX (moderncv) template |
| `internal/latex/latex_test.go` | Golden test |
| `internal/latex/testdata/*.golden` | Expected LaTeX output |
| `internal/html/html.go` | `HTMLRenderer` mini-site + embedded templates + asset copy |
| `internal/html/templates/index.html.tmpl` | HTML timeline template |
| `internal/html/templates/style.css` | Mini-site stylesheet (static asset) |
| `internal/html/html_test.go` | Golden test (index.html) |
| `internal/html/testdata/*.golden` | Expected HTML output |
| `internal/writer/writer.go` | Write `[]Artifact` to disk (mkdir -p) |
| `internal/writer/writer_test.go` | Writer test (temp dir) |
| `resume.yaml` | Single source of truth |
| `Makefile` | Go-driven generate + compile targets |
| `.github/workflows/release.yml` | setup-go, build PDFs, attach 4 PDFs, deploy Pages |
| `GUIDE.md` | Fork/customization instructions |

---

## Task 1: Initialize Go module

**Files:**
- Create: `go.mod`
- Create: `.gitignore` (modify existing — add Go + generated dirs)

- [ ] **Step 1: Create the module**

Run:
```bash
cd /Users/george/projects/georgemessiha22
go mod init github.com/georgemessiha22/georgemessiha22
go get gopkg.in/yaml.v3@v3.0.1
```
Expected: creates `go.mod` and `go.sum` with `gopkg.in/yaml.v3 v3.0.1`.

- [ ] **Step 2: Add a smoke test package to prove the toolchain**

Create `internal/model/resume.go` with a package clause only (filled in Task 2):
```go
package model
```

Run: `go build ./...`
Expected: success, no output.

- [ ] **Step 3: Extend .gitignore**

Append to `.gitignore`:
```
# Go
/build/
/site/
*.test
*.out
```

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum internal/model/resume.go .gitignore
git commit -m "chore: initialize go module for resume generator"
```

---

## Task 2: Domain model + variant filtering

**Files:**
- Modify: `internal/model/resume.go`
- Test: `internal/model/resume_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/model/resume_test.go`:
```go
package model

import "testing"

func sampleResume() Resume {
	return Resume{
		Contact: Contact{Firstname: "George", Lastname: "Messiha"},
		Summary: "summary text",
		Experience: []Entry{
			{Role: "A", Variants: []Variant{Summary, Detailed}},
			{Role: "B", Variants: []Variant{Detailed}},
			{Role: "C"}, // untagged => both
		},
		Certificates: []Cert{
			{Title: "X", Variants: []Variant{Detailed}},
			{Title: "Y", Variants: []Variant{Summary, Detailed}},
		},
		Activities: []Entry{{Role: "Act", Variants: []Variant{Detailed}}},
	}
}

func TestForVariant_Summary(t *testing.T) {
	got := sampleResume().ForVariant(Summary)
	if len(got.Experience) != 2 {
		t.Fatalf("summary experience: want 2, got %d", len(got.Experience))
	}
	if got.Experience[0].Role != "A" || got.Experience[1].Role != "C" {
		t.Fatalf("unexpected roles: %+v", got.Experience)
	}
	if len(got.Certificates) != 1 || got.Certificates[0].Title != "Y" {
		t.Fatalf("summary certs wrong: %+v", got.Certificates)
	}
	if len(got.Activities) != 0 {
		t.Fatalf("summary activities: want 0, got %d", len(got.Activities))
	}
}

func TestForVariant_Detailed(t *testing.T) {
	got := sampleResume().ForVariant(Detailed)
	if len(got.Experience) != 3 {
		t.Fatalf("detailed experience: want 3, got %d", len(got.Experience))
	}
	if len(got.Certificates) != 2 {
		t.Fatalf("detailed certs: want 2, got %d", len(got.Certificates))
	}
	if len(got.Activities) != 1 {
		t.Fatalf("detailed activities: want 1, got %d", len(got.Activities))
	}
}
```

Run: `go test ./internal/model/ -run TestForVariant -v`
Expected: FAIL (compile error — types undefined).

- [ ] **Step 2: Implement the model**

Replace `internal/model/resume.go` with:
```go
package model

// Variant identifies which resume a piece of content belongs to.
type Variant string

const (
	Summary  Variant = "summary"
	Detailed Variant = "detailed"
)

// Contact is the header / personal information.
type Contact struct {
	Firstname string
	Lastname  string
	Title     string
	Email     string
	Phone     string
	Location  string
	Photo     string // path to image file
	Socials   Socials
}

// Socials holds usernames (not full URLs) per platform.
type Socials struct {
	GitHub   string
	GitLab   string
	LinkedIn string
}

// Entry is a timeline item used for experience, education, and activities.
type Entry struct {
	Role     string // job title / degree / activity role
	Org      string
	OrgURL   string
	URL      string // optional link on the title (e.g. certificate scan)
	Location string
	Mode     string // e.g. "Hybrid", "Remote", "Fulltime"
	Start    string
	End      string
	Note     string // e.g. thesis line
	Intro    string // paragraph before bullets
	Bullets  []string
	Variants []Variant
}

// SkillGroup is a labelled list of skills.
type SkillGroup struct {
	Category string
	Items    []string
}

// Language is a spoken language and proficiency.
type Language struct {
	Name  string
	Level string
}

// Cert is a certificate or award.
type Cert struct {
	Title    string
	Org      string
	OrgURL   string
	CertURL  string // link to the certificate scan
	Location string
	Start    string
	End      string
	Group    string // optional subsection, e.g. "Courses" / "Tracks"
	Variants []Variant
}

// Resume is the whole document.
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

func includesVariant(vs []Variant, v Variant) bool {
	if len(vs) == 0 {
		return true // untagged => appears in every variant
	}
	for _, x := range vs {
		if x == v {
			return true
		}
	}
	return false
}

func filterEntries(in []Entry, v Variant) []Entry {
	out := make([]Entry, 0, len(in))
	for _, e := range in {
		if includesVariant(e.Variants, v) {
			out = append(out, e)
		}
	}
	return out
}

func filterCerts(in []Cert, v Variant) []Cert {
	out := make([]Cert, 0, len(in))
	for _, c := range in {
		if includesVariant(c.Variants, v) {
			out = append(out, c)
		}
	}
	return out
}

// ForVariant returns a copy of the resume containing only the content that
// belongs to variant v. Skills and languages always pass through.
func (r Resume) ForVariant(v Variant) Resume {
	out := r
	out.Experience = filterEntries(r.Experience, v)
	out.Education = filterEntries(r.Education, v)
	out.Activities = filterEntries(r.Activities, v)
	out.Certificates = filterCerts(r.Certificates, v)
	return out
}
```

Run: `go test ./internal/model/ -run TestForVariant -v`
Expected: PASS (both tests).

- [ ] **Step 3: Commit**

```bash
git add internal/model/
git commit -m "feat(model): add resume domain model with variant filtering"
```

---

## Task 3: render core — Renderer interface + inline emphasis

**Files:**
- Create: `internal/render/render.go`
- Test: `internal/render/render_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/render/render_test.go`:
```go
package render

import "testing"

func TestParseSpans(t *testing.T) {
	got := ParseSpans("plain **bold** end")
	want := []Span{
		{Text: "plain ", Bold: false},
		{Text: "bold", Bold: true},
		{Text: " end", Bold: false},
	}
	if len(got) != len(want) {
		t.Fatalf("len: want %d got %d (%+v)", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("span %d: want %+v got %+v", i, want[i], got[i])
		}
	}
}

func TestParseSpans_NoBold(t *testing.T) {
	got := ParseSpans("just text")
	if len(got) != 1 || got[0].Bold || got[0].Text != "just text" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestParseSpans_Adjacent(t *testing.T) {
	got := ParseSpans("**a****b**")
	want := []Span{{Text: "a", Bold: true}, {Text: "b", Bold: true}}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("unexpected: %+v", got)
	}
}
```

Run: `go test ./internal/render/ -v`
Expected: FAIL (undefined `ParseSpans`, `Span`).

- [ ] **Step 2: Implement render core**

Create `internal/render/render.go`:
```go
package render

import (
	"strings"

	"github.com/georgemessiha22/georgemessiha22/internal/model"
)

// Variant is re-exported so engine packages depend only on render.
type Variant = model.Variant

// Artifact is one rendered file relative to the output directory.
type Artifact struct {
	RelPath string
	Bytes   []byte
}

// Renderer turns a resume + variant into one or more files.
type Renderer interface {
	Render(r model.Resume, v model.Variant) ([]Artifact, error)
}

// Span is a run of text that is either plain or bold.
type Span struct {
	Text string
	Bold bool
}

// ParseSpans splits text on **bold** markers into ordered spans. Empty spans
// (e.g. from adjacent markers) are dropped.
func ParseSpans(s string) []Span {
	var spans []Span
	bold := false
	for len(s) > 0 {
		idx := strings.Index(s, "**")
		if idx < 0 {
			spans = appendSpan(spans, s, bold)
			break
		}
		spans = appendSpan(spans, s[:idx], bold)
		bold = !bold
		s = s[idx+2:]
	}
	return spans
}

func appendSpan(spans []Span, text string, bold bool) []Span {
	if text == "" {
		return spans
	}
	return append(spans, Span{Text: text, Bold: bold})
}
```

Run: `go test ./internal/render/ -v`
Expected: PASS (all three tests).

- [ ] **Step 3: Commit**

```bash
git add internal/render/
git commit -m "feat(render): add Renderer interface and inline emphasis parser"
```

---

## Task 4: config — load + validate YAML

**Files:**
- Create: `internal/config/load.go`
- Test: `internal/config/load_test.go`
- Create: `internal/config/testdata/valid.yaml`, `internal/config/testdata/missing_name.yaml`

- [ ] **Step 1: Add fixtures**

Create `internal/config/testdata/valid.yaml`:
```yaml
contact:
  firstname: George
  lastname: Messiha
  title: Lead Software Engineer
  email: g@example.com
  phone: "(+971) 54 555 1032"
  location: "Dubai, UAE"
  photo: pictures/61673.jpg
  socials:
    github: georgemessiha22
    linkedin: georgemessiha22
summary: A short summary.
experience:
  - role: Senior Software Engineer
    org: HungerStation
    org_url: https://hungerstation.com
    location: Dubai, UAE
    mode: Hybrid
    start: Nov 2023
    end: Present
    variants: [summary, detailed]
    intro: Did things.
    bullets:
      - "Built **Go** services."
skills:
  - category: Languages
    items: [Go, Python]
languages:
  - name: Arabic
    level: Native
certificates:
  - title: Foundation of Project Management
    org: Google
    cert_url: https://example.com/cert.pdf
    start: "2024"
    end: "2024"
    variants: [summary, detailed]
```

Create `internal/config/testdata/missing_name.yaml`:
```yaml
contact:
  firstname: ""
  lastname: ""
summary: x
experience:
  - role: Dev
    org: Acme
```

- [ ] **Step 2: Write the failing test**

Create `internal/config/load_test.go`:
```go
package config

import (
	"strings"
	"testing"

	"github.com/georgemessiha22/georgemessiha22/internal/model"
)

func TestLoad_Valid(t *testing.T) {
	r, err := Load("testdata/valid.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Contact.Firstname != "George" || r.Contact.Lastname != "Messiha" {
		t.Fatalf("contact not parsed: %+v", r.Contact)
	}
	if len(r.Experience) != 1 || r.Experience[0].Org != "HungerStation" {
		t.Fatalf("experience not parsed: %+v", r.Experience)
	}
	if r.Experience[0].Variants[0] != model.Summary {
		t.Fatalf("variants not parsed: %+v", r.Experience[0].Variants)
	}
	if len(r.Skills) != 1 || r.Skills[0].Category != "Languages" {
		t.Fatalf("skills not parsed: %+v", r.Skills)
	}
}

func TestLoad_MissingName(t *testing.T) {
	_, err := Load("testdata/missing_name.yaml")
	if err == nil || !strings.Contains(err.Error(), "name") {
		t.Fatalf("expected name validation error, got %v", err)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("testdata/does_not_exist.yaml")
	if err == nil {
		t.Fatalf("expected error for missing file")
	}
}
```

Run: `go test ./internal/config/ -v`
Expected: FAIL (undefined `Load`).

- [ ] **Step 3: Implement loader with YAML tags**

Create `internal/config/load.go`:
```go
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/georgemessiha22/georgemessiha22/internal/model"
)

// yamlResume mirrors model.Resume with yaml tags. Keeping a separate DTO keeps
// yaml concerns out of the domain model.
type yamlResume struct {
	Contact struct {
		Firstname string `yaml:"firstname"`
		Lastname  string `yaml:"lastname"`
		Title     string `yaml:"title"`
		Email     string `yaml:"email"`
		Phone     string `yaml:"phone"`
		Location  string `yaml:"location"`
		Photo     string `yaml:"photo"`
		Socials   struct {
			GitHub   string `yaml:"github"`
			GitLab   string `yaml:"gitlab"`
			LinkedIn string `yaml:"linkedin"`
		} `yaml:"socials"`
	} `yaml:"contact"`
	Summary      string       `yaml:"summary"`
	Experience   []yamlEntry  `yaml:"experience"`
	Education    []yamlEntry  `yaml:"education"`
	Skills       []yamlSkill  `yaml:"skills"`
	Languages    []yamlLang   `yaml:"languages"`
	Certificates []yamlCert   `yaml:"certificates"`
	Activities   []yamlEntry  `yaml:"activities"`
}

type yamlEntry struct {
	Role     string   `yaml:"role"`
	Org      string   `yaml:"org"`
	OrgURL   string   `yaml:"org_url"`
	URL      string   `yaml:"url"`
	Location string   `yaml:"location"`
	Mode     string   `yaml:"mode"`
	Start    string   `yaml:"start"`
	End      string   `yaml:"end"`
	Note     string   `yaml:"note"`
	Intro    string   `yaml:"intro"`
	Bullets  []string `yaml:"bullets"`
	Variants []string `yaml:"variants"`
}

type yamlSkill struct {
	Category string   `yaml:"category"`
	Items    []string `yaml:"items"`
}

type yamlLang struct {
	Name  string `yaml:"name"`
	Level string `yaml:"level"`
}

type yamlCert struct {
	Title    string   `yaml:"title"`
	Org      string   `yaml:"org"`
	OrgURL   string   `yaml:"org_url"`
	CertURL  string   `yaml:"cert_url"`
	Location string   `yaml:"location"`
	Start    string   `yaml:"start"`
	End      string   `yaml:"end"`
	Group    string   `yaml:"group"`
	Variants []string `yaml:"variants"`
}

func toVariants(in []string) []model.Variant {
	if len(in) == 0 {
		return nil
	}
	out := make([]model.Variant, 0, len(in))
	for _, s := range in {
		out = append(out, model.Variant(s))
	}
	return out
}

func toEntry(y yamlEntry) model.Entry {
	return model.Entry{
		Role: y.Role, Org: y.Org, OrgURL: y.OrgURL, URL: y.URL,
		Location: y.Location, Mode: y.Mode, Start: y.Start, End: y.End,
		Note: y.Note, Intro: y.Intro, Bullets: y.Bullets,
		Variants: toVariants(y.Variants),
	}
}

func toEntries(in []yamlEntry) []model.Entry {
	out := make([]model.Entry, 0, len(in))
	for _, y := range in {
		out = append(out, toEntry(y))
	}
	return out
}

// Load reads a YAML file and returns a validated Resume.
func Load(path string) (model.Resume, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return model.Resume{}, fmt.Errorf("read resume file: %w", err)
	}
	var y yamlResume
	if err := yaml.Unmarshal(data, &y); err != nil {
		return model.Resume{}, fmt.Errorf("parse yaml: %w", err)
	}

	r := model.Resume{
		Contact: model.Contact{
			Firstname: y.Contact.Firstname,
			Lastname:  y.Contact.Lastname,
			Title:     y.Contact.Title,
			Email:     y.Contact.Email,
			Phone:     y.Contact.Phone,
			Location:  y.Contact.Location,
			Photo:     y.Contact.Photo,
			Socials: model.Socials{
				GitHub:   y.Contact.Socials.GitHub,
				GitLab:   y.Contact.Socials.GitLab,
				LinkedIn: y.Contact.Socials.LinkedIn,
			},
		},
		Summary:    y.Summary,
		Experience: toEntries(y.Experience),
		Education:  toEntries(y.Education),
		Activities: toEntries(y.Activities),
	}
	for _, s := range y.Skills {
		r.Skills = append(r.Skills, model.SkillGroup{Category: s.Category, Items: s.Items})
	}
	for _, l := range y.Languages {
		r.Languages = append(r.Languages, model.Language{Name: l.Name, Level: l.Level})
	}
	for _, c := range y.Certificates {
		r.Certificates = append(r.Certificates, model.Cert{
			Title: c.Title, Org: c.Org, OrgURL: c.OrgURL, CertURL: c.CertURL,
			Location: c.Location, Start: c.Start, End: c.End, Group: c.Group,
			Variants: toVariants(c.Variants),
		})
	}

	if err := validate(r); err != nil {
		return model.Resume{}, err
	}
	return r, nil
}

func validate(r model.Resume) error {
	if r.Contact.Firstname == "" && r.Contact.Lastname == "" {
		return fmt.Errorf("validation: contact name is required (firstname/lastname)")
	}
	if len(r.Experience) == 0 {
		return fmt.Errorf("validation: at least one experience entry is required")
	}
	return nil
}
```

Run: `go test ./internal/config/ -v`
Expected: PASS (all three tests).

- [ ] **Step 4: Commit**

```bash
git add internal/config/
git commit -m "feat(config): load and validate resume.yaml into the model"
```

---

## Task 5: Typst engine

**Files:**
- Create: `internal/typst/typst.go`
- Create: `internal/typst/templates/resume.typ.tmpl`
- Test: `internal/typst/typst_test.go`
- Create: `internal/typst/testdata/summary.typ.golden`

**Note on templates:** templates call funcs `esc` (escape Typst specials) and `emph` (escape + apply `*bold*`). The renderer registers these in a FuncMap.

- [ ] **Step 1: Write the template**

Create `internal/typst/templates/resume.typ.tmpl`:
```
#import "@preview/modern-cv:0.10.0": *

#show: resume.with(
  author: (
    firstname: "{{ esc .Contact.Firstname }}",
    lastname: "{{ esc .Contact.Lastname }}",
    email: "{{ esc .Contact.Email }}",
    phone: "{{ esc .Contact.Phone }}",
{{- if .Contact.Socials.GitHub }}
    github: "{{ esc .Contact.Socials.GitHub }}",
{{- end }}
{{- if .Contact.Socials.GitLab }}
    gitlab: "{{ esc .Contact.Socials.GitLab }}",
{{- end }}
{{- if .Contact.Socials.LinkedIn }}
    linkedin: "{{ esc .Contact.Socials.LinkedIn }}",
{{- end }}
    address: "{{ esc .Contact.Location }}",
    positions: ("{{ esc .Contact.Title }}",),
  ),
  profile-picture: image("assets/profile.jpg"),
  date: datetime.today().display(),
  language: "en",
  colored-headers: true,
  show-footer: false,
  paper-size: "a4",
)

= Professional Summary

#resume-item[
  {{ emph .Summary }}
]

= Experience

{{ range .Experience }}#resume-entry(
  title: {{ if .URL }}link("{{ esc .URL }}")[{{ esc .Role }}]{{ else }}"{{ esc .Role }}"{{ end }},
  location: "{{ esc .Location }}{{ if .Mode }} · {{ esc .Mode }}{{ end }}",
  date: "{{ esc .Start }}{{ if .End }} - {{ esc .End }}{{ end }}",
  description: {{ if .OrgURL }}smallcaps(link("{{ esc .OrgURL }}")[{{ esc .Org }}]){{ else }}smallcaps("{{ esc .Org }}"){{ end }},
)
{{ if or .Intro .Bullets }}#resume-item[
{{- if .Intro }}
  {{ emph .Intro }}
{{- end }}
{{- range .Bullets }}
  - {{ emph . }}
{{- end }}
]
{{ end }}
{{ end -}}
= Education

{{ range .Education }}#resume-entry(
  title: {{ if .URL }}link("{{ esc .URL }}")[{{ esc .Role }}]{{ else }}"{{ esc .Role }}"{{ end }},
  location: "{{ esc .Location }}",
  date: "{{ esc .Start }}{{ if .End }} - {{ esc .End }}{{ end }}",
  description: "{{ esc .Org }}",
)
{{ if .Note }}#resume-item[
  _{{ emph .Note }}_
]
{{ end }}
{{ end -}}
= Skills

{{ range .Skills }}#resume-skill-item("{{ esc .Category }}", ({{ range $i, $it := .Items }}{{ if $i }}, {{ end }}"{{ esc $it }}"{{ end }},))
{{ end }}
= Languages

{{ range .Languages }}#resume-skill-item("{{ esc .Name }}", (strong("{{ esc .Level }}"),))
{{ end }}
= Certificates & Awards

{{ range .Certificates }}#resume-entry(
  title: {{ if .CertURL }}link("{{ esc .CertURL }}")[{{ esc .Title }}]{{ else }}"{{ esc .Title }}"{{ end }},
  location: "{{ esc .Location }}",
  date: "{{ esc .Start }}{{ if .End }} - {{ esc .End }}{{ end }}",
  description: {{ if .OrgURL }}link("{{ esc .OrgURL }}")[{{ esc .Org }}]{{ else }}"{{ esc .Org }}"{{ end }},
)
{{ end -}}
{{ if .Activities }}= Activities

{{ range .Activities }}#resume-entry(
  title: "{{ esc .Role }}",
  location: "{{ esc .Location }}",
  date: "{{ esc .Start }}{{ if .End }} - {{ esc .End }}{{ end }}",
  description: "{{ esc .Org }}",
)
{{ if .Intro }}#resume-item[
  {{ emph .Intro }}
]
{{ end }}
{{ end }}{{ end -}}
```

- [ ] **Step 2: Implement the renderer**

Create `internal/typst/typst.go`:
```go
package typst

import (
	"bytes"
	"embed"
	"strings"
	"text/template"

	"github.com/georgemessiha22/georgemessiha22/internal/model"
	"github.com/georgemessiha22/georgemessiha22/internal/render"
)

//go:embed templates/resume.typ.tmpl
var tmplFS embed.FS

// Renderer renders the resume to a single Typst file.
type Renderer struct{}

// escape escapes Typst special characters in a plain string.
func escape(s string) string {
	r := strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
		`#`, `\#`,
		`$`, `\$`,
		`@`, `\@`,
		`*`, `\*`,
		`_`, `\_`,
	)
	return r.Replace(s)
}

// emphasize escapes text and converts **bold** spans to Typst *bold*.
func emphasize(s string) string {
	var b strings.Builder
	for _, sp := range render.ParseSpans(s) {
		if sp.Bold {
			b.WriteString("*" + escape(sp.Text) + "*")
		} else {
			b.WriteString(escape(sp.Text))
		}
	}
	return b.String()
}

func (Renderer) Render(r model.Resume, v model.Variant) ([]render.Artifact, error) {
	rr := r.ForVariant(v)
	t, err := template.New("resume.typ.tmpl").Funcs(template.FuncMap{
		"esc":  escape,
		"emph": emphasize,
	}).ParseFS(tmplFS, "templates/resume.typ.tmpl")
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, rr); err != nil {
		return nil, err
	}
	name := "resume.typ"
	if v == model.Detailed {
		name = "resume_detailed.typ"
	}
	return []render.Artifact{{RelPath: name, Bytes: buf.Bytes()}}, nil
}
```

- [ ] **Step 3: Write the golden test (auto-generate then lock)**

Create `internal/typst/typst_test.go`:
```go
package typst

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/georgemessiha22/georgemessiha22/internal/model"
)

var update = flag.Bool("update", false, "update golden files")

func fixture() model.Resume {
	return model.Resume{
		Contact: model.Contact{
			Firstname: "George", Lastname: "Messiha", Title: "Lead SWE",
			Email: "g@example.com", Phone: "(+971) 54 555 1032",
			Location: "Dubai, UAE",
			Socials:  model.Socials{GitHub: "georgemessiha22", LinkedIn: "georgemessiha22"},
		},
		Summary: "A summary with **bold**.",
		Experience: []model.Entry{{
			Role: "Senior Software Engineer", Org: "HungerStation",
			OrgURL: "https://hungerstation.com", Location: "Dubai, UAE", Mode: "Hybrid",
			Start: "Nov 2023", End: "Present", Intro: "Did things in **Go**.",
			Bullets: []string{"Built **Go** services.", "Cut build 60%."},
		}},
		Education: []model.Entry{{
			Role: "BSc", Org: "GUC", Location: "Cairo, EG",
			Start: "2011", End: "2016", Note: "Thesis: Cloud",
		}},
		Skills:    []model.SkillGroup{{Category: "Languages", Items: []string{"Go", "Python"}}},
		Languages: []model.Language{{Name: "Arabic", Level: "Native"}},
		Certificates: []model.Cert{{
			Title: "PM Foundation", Org: "Google", CertURL: "https://x/c.pdf",
			Start: "2024", End: "2024",
		}},
	}
}

func TestRenderSummaryGolden(t *testing.T) {
	got, err := Renderer{}.Render(fixture(), model.Summary)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].RelPath != "resume.typ" {
		t.Fatalf("unexpected artifacts: %+v", got)
	}
	golden := filepath.Join("testdata", "summary.typ.golden")
	if *update {
		if err := os.WriteFile(golden, got[0].Bytes, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if string(got[0].Bytes) != string(want) {
		t.Fatalf("output differs from golden; run with -update if intended")
	}
}
```

- [ ] **Step 4: Generate the golden file, then verify it locks**

Run:
```bash
mkdir -p internal/typst/testdata
go test ./internal/typst/ -run TestRenderSummaryGolden -update
go test ./internal/typst/ -run TestRenderSummaryGolden -v
```
Expected: first command writes `summary.typ.golden`; second PASSes.

- [ ] **Step 5: Sanity-check the golden compiles with Typst**

Run:
```bash
mkdir -p /tmp/typstcheck/assets
cp pictures/61673.jpg /tmp/typstcheck/assets/profile.jpg
cp internal/typst/testdata/summary.typ.golden /tmp/typstcheck/resume.typ
typst compile --font-path typst/fonts /tmp/typstcheck/resume.typ /tmp/typstcheck/out.pdf
```
Expected: produces `out.pdf` (font fallback warning is OK). If it errors, fix the template and re-run Step 4.

- [ ] **Step 6: Commit**

```bash
git add internal/typst/
git commit -m "feat(typst): add Typst renderer with golden test"
```

---

## Task 6: LaTeX engine

**Files:**
- Create: `internal/latex/latex.go`
- Create: `internal/latex/templates/resume.tex.tmpl`
- Test: `internal/latex/latex_test.go`
- Create: `internal/latex/testdata/summary.tex.golden`

- [ ] **Step 1: Write the template**

Create `internal/latex/templates/resume.tex.tmpl`:
```
\documentclass[10pt,a4paper,sans]{moderncv}
\usepackage{fontawesome5}
\usepackage[T1]{fontenc}
\usepackage{multicol}
\usepackage[scale=0.8,top=1.9cm, bottom=1.9cm]{geometry}
\usepackage[utf8]{inputenc}
\usepackage{lmodern}
\usepackage[unicode]{hyperref}
\definecolor{airforceblue}{rgb}{0.36, 0.54, 0.66}
\hypersetup{colorlinks=true, linkcolor=airforceblue, urlcolor=airforceblue}
\nopagenumbers{}
\moderncvtheme[grey]{banking}
\firstname{ {{- esc .Contact.Firstname -}} }
\familyname{ {{- esc .Contact.Lastname -}} }
\title{ {{- esc .Contact.Title -}} }
\address{ {{- esc .Contact.Location -}} }{}
\mobile{ {{- esc .Contact.Phone -}} }
\email{ {{- esc .Contact.Email -}} }
{{- if .Contact.Socials.LinkedIn }}
\social[linkedin]{ {{- esc .Contact.Socials.LinkedIn -}} }
{{- end }}
{{- if .Contact.Socials.GitHub }}
\social[github]{ {{- esc .Contact.Socials.GitHub -}} }
{{- end }}
{{- if .Contact.Socials.GitLab }}
\social[gitlab]{ {{- esc .Contact.Socials.GitLab -}} }
{{- end }}
\begin{document}
\makecvtitle

\section{\textbf{Professional Summary}}
\cvitem{}{ {{- emph .Summary -}} }

\section{Experience}
{{ range .Experience }}\cventry
  { {{- esc .Start }}{{ if .End }}--{{ esc .End }}{{ end -}} }
  {\textbf{ {{- esc .Role -}} }}
  {\textsc{ {{- if .OrgURL }}\href{ {{- esc .OrgURL -}} }{ {{- esc .Org -}} }{{ else }}{{ esc .Org }}{{ end -}} }}
  { {{- esc .Location -}} }
  { {{- esc .Mode -}} }
  {
{{- if .Intro }}
{{ emph .Intro }}
{{- end }}
{{- if .Bullets }}
\begin{itemize}
{{- range .Bullets }}
  \item {{ emph . }}
{{- end }}
\end{itemize}
{{- end }}
}
{{ end }}
\section{Education}
{{ range .Education }}\cventry
  { {{- esc .Start }}{{ if .End }}--{{ esc .End }}{{ end -}} }
  { {{- if .URL }}\href{ {{- esc .URL -}} }{ {{- esc .Role -}} }{{ else }}{{ esc .Role }}{{ end -}} }
  { {{- esc .Org -}} }
  { {{- esc .Location -}} }
  {}
  { {{- if .Note }}\textit{ {{- emph .Note -}} }{{ end -}} }
{{ end }}
\section{Skills}
{{ range .Skills }}\cvitem{ {{- esc .Category -}} }{ {{- range $i, $it := .Items }}{{ if $i }}, {{ end }}{{ esc $it }}{{ end -}} }
{{ end }}
\section{Languages}
{{ range .Languages }}\cvitem{ {{- esc .Name -}} }{ {{- esc .Level -}} }
{{ end }}
\section{Certificates \& Awards}
{{ range .Certificates }}\cventry
  { {{- esc .Start }}{{ if .End }}--{{ esc .End }}{{ end -}} }
  { {{- if .CertURL }}\href{ {{- esc .CertURL -}} }{ {{- esc .Title -}} }{{ else }}{{ esc .Title }}{{ end -}} }
  { {{- if .OrgURL }}\href{ {{- esc .OrgURL -}} }{ {{- esc .Org -}} }{{ else }}{{ esc .Org }}{{ end -}} }
  { {{- esc .Location -}} }
  {}
  {}
{{ end }}
{{- if .Activities }}\section{Activities}
{{ range .Activities }}\cventry
  { {{- esc .Start }}{{ if .End }}--{{ esc .End }}{{ end -}} }
  { {{- esc .Role -}} }
  { {{- esc .Org -}} }
  { {{- esc .Location -}} }
  {}
  { {{- emph .Intro -}} }
{{ end }}{{ end }}
\end{document}
```

- [ ] **Step 2: Implement the renderer**

Create `internal/latex/latex.go`:
```go
package latex

import (
	"bytes"
	"embed"
	"strings"
	"text/template"

	"github.com/georgemessiha22/georgemessiha22/internal/model"
	"github.com/georgemessiha22/georgemessiha22/internal/render"
)

//go:embed templates/resume.tex.tmpl
var tmplFS embed.FS

// Renderer renders the resume to a single LaTeX (moderncv) file.
type Renderer struct{}

// escape escapes LaTeX special characters in a plain string.
func escape(s string) string {
	r := strings.NewReplacer(
		`\`, `\textbackslash{}`,
		`&`, `\&`,
		`%`, `\%`,
		`$`, `\$`,
		`#`, `\#`,
		`_`, `\_`,
		`{`, `\{`,
		`}`, `\}`,
		`~`, `\textasciitilde{}`,
		`^`, `\textasciicircum{}`,
	)
	return r.Replace(s)
}

// emphasize escapes text and converts **bold** spans to \textbf{}.
func emphasize(s string) string {
	var b strings.Builder
	for _, sp := range render.ParseSpans(s) {
		if sp.Bold {
			b.WriteString(`\textbf{` + escape(sp.Text) + `}`)
		} else {
			b.WriteString(escape(sp.Text))
		}
	}
	return b.String()
}

func (Renderer) Render(r model.Resume, v model.Variant) ([]render.Artifact, error) {
	rr := r.ForVariant(v)
	t, err := template.New("resume.tex.tmpl").Funcs(template.FuncMap{
		"esc":  escape,
		"emph": emphasize,
	}).ParseFS(tmplFS, "templates/resume.tex.tmpl")
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, rr); err != nil {
		return nil, err
	}
	name := "resume.tex"
	if v == model.Detailed {
		name = "resume_detailed.tex"
	}
	return []render.Artifact{{RelPath: name, Bytes: buf.Bytes()}}, nil
}
```

- [ ] **Step 3: Write the golden test**

Create `internal/latex/latex_test.go` (reuse the same fixture shape as Typst):
```go
package latex

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/georgemessiha22/georgemessiha22/internal/model"
)

var update = flag.Bool("update", false, "update golden files")

func fixture() model.Resume {
	return model.Resume{
		Contact: model.Contact{
			Firstname: "George", Lastname: "Messiha", Title: "Lead SWE",
			Email: "g@example.com", Phone: "(+971) 54 555 1032",
			Location: "Dubai, UAE",
			Socials:  model.Socials{GitHub: "georgemessiha22", LinkedIn: "georgemessiha22"},
		},
		Summary: "A summary with **bold** & 50%.",
		Experience: []model.Entry{{
			Role: "Senior Software Engineer", Org: "HungerStation",
			OrgURL: "https://hungerstation.com", Location: "Dubai, UAE", Mode: "Hybrid",
			Start: "Nov 2023", End: "Present", Intro: "Did things in **Go**.",
			Bullets: []string{"Built **Go** services.", "Cut build 60%."},
		}},
		Education: []model.Entry{{
			Role: "BSc", Org: "GUC", Location: "Cairo, EG",
			Start: "2011", End: "2016", Note: "Thesis: Cloud",
		}},
		Skills:    []model.SkillGroup{{Category: "Languages", Items: []string{"Go", "Python"}}},
		Languages: []model.Language{{Name: "Arabic", Level: "Native"}},
		Certificates: []model.Cert{{
			Title: "PM Foundation", Org: "Google", CertURL: "https://x/c.pdf",
			Start: "2024", End: "2024",
		}},
	}
}

func TestRenderSummaryGolden(t *testing.T) {
	got, err := Renderer{}.Render(fixture(), model.Summary)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].RelPath != "resume.tex" {
		t.Fatalf("unexpected artifacts: %+v", got)
	}
	golden := filepath.Join("testdata", "summary.tex.golden")
	if *update {
		if err := os.WriteFile(golden, got[0].Bytes, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if string(got[0].Bytes) != string(want) {
		t.Fatalf("output differs from golden; run with -update if intended")
	}
}
```

- [ ] **Step 4: Generate golden + verify**

Run:
```bash
mkdir -p internal/latex/testdata
go test ./internal/latex/ -run TestRenderSummaryGolden -update
go test ./internal/latex/ -run TestRenderSummaryGolden -v
```
Expected: golden written; second run PASSes. Inspect the golden to confirm `&` and `%` are escaped (`\&`, `60\%`).

- [ ] **Step 5: Commit**

```bash
git add internal/latex/
git commit -m "feat(latex): add LaTeX/moderncv renderer with golden test"
```

---

## Task 7: HTML engine (GitHub Pages mini-site, timeline)

**Files:**
- Create: `internal/html/html.go`
- Create: `internal/html/templates/index.html.tmpl`
- Create: `internal/html/templates/style.css`
- Test: `internal/html/html_test.go`
- Create: `internal/html/testdata/index.html.golden`

**Design:** `Render` returns three artifacts: `index.html` (from template), `style.css` (embedded static), and `profile.jpg` (copied from `Contact.Photo` on disk; if the file is missing, the image artifact is skipped and the template omits the `<img>`). The HTML uses `html/template` for auto-escaping; bold spans are emitted as `template.HTML` via a `emph` func that escapes then wraps bold in `<strong>`.

- [ ] **Step 1: Write the stylesheet**

Create `internal/html/templates/style.css`:
```css
:root { --accent:#2f5fa6; --ink:#222; --muted:#666; --line:#dcdcdc; }
* { box-sizing: border-box; }
body { font-family: -apple-system, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  color: var(--ink); max-width: 820px; margin: 0 auto; padding: 2.5rem 1.25rem; line-height: 1.5; }
header.top { display: flex; gap: 1.5rem; align-items: center; border-bottom: 2px solid var(--accent); padding-bottom: 1rem; }
header.top img { width: 96px; height: 96px; border-radius: 50%; object-fit: cover; }
header.top h1 { margin: 0; font-size: 2rem; }
header.top h1 .last { color: var(--accent); }
.title { color: var(--accent); font-weight: 600; }
.contact { color: var(--muted); font-size: .9rem; margin-top: .25rem; }
.contact a { color: var(--accent); text-decoration: none; }
h2 { color: var(--accent); border-bottom: 1px solid var(--line); padding-bottom: .25rem; margin-top: 2rem; }
.timeline { position: relative; margin-left: 1rem; padding-left: 1.25rem; border-left: 2px solid var(--line); }
.tl-item { position: relative; margin-bottom: 1.25rem; }
.tl-item::before { content: ""; position: absolute; left: -1.66rem; top: .35rem; width: 10px; height: 10px;
  background: var(--accent); border-radius: 50%; }
.tl-head { display: flex; justify-content: space-between; gap: 1rem; flex-wrap: wrap; }
.tl-role { font-weight: 700; }
.tl-date { color: var(--muted); font-size: .85rem; white-space: nowrap; }
.tl-org { color: var(--muted); font-variant: small-caps; }
.tl-item ul { margin: .4rem 0 0; padding-left: 1.1rem; }
.skills div { margin: .2rem 0; }
.skills .cat { font-weight: 600; }
@media print { body { padding: 0; } a { color: var(--ink); } }
```

- [ ] **Step 2: Write the HTML template**

Create `internal/html/templates/index.html.tmpl`:
```
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{ .Contact.Firstname }} {{ .Contact.Lastname }} — Resume</title>
<link rel="stylesheet" href="style.css">
</head>
<body>
<header class="top">
{{- if .HasPhoto }}
  <img src="profile.jpg" alt="{{ .Contact.Firstname }} {{ .Contact.Lastname }}">
{{- end }}
  <div>
    <h1>{{ .Contact.Firstname }} <span class="last">{{ .Contact.Lastname }}</span></h1>
    <div class="title">{{ .Contact.Title }}</div>
    <div class="contact">
      {{ .Contact.Location }} &middot;
      <a href="mailto:{{ .Contact.Email }}">{{ .Contact.Email }}</a>
      {{- if .Contact.Socials.GitHub }} &middot; <a href="https://github.com/{{ .Contact.Socials.GitHub }}">GitHub</a>{{ end }}
      {{- if .Contact.Socials.GitLab }} &middot; <a href="https://gitlab.com/{{ .Contact.Socials.GitLab }}">GitLab</a>{{ end }}
      {{- if .Contact.Socials.LinkedIn }} &middot; <a href="https://linkedin.com/in/{{ .Contact.Socials.LinkedIn }}">LinkedIn</a>{{ end }}
    </div>
  </div>
</header>

<section>
  <h2>Professional Summary</h2>
  <p>{{ emph .Summary }}</p>
</section>

<section>
  <h2>Experience</h2>
  <div class="timeline">
  {{- range .Experience }}
    <div class="tl-item">
      <div class="tl-head">
        <span class="tl-role">{{ if .URL }}<a href="{{ .URL }}">{{ .Role }}</a>{{ else }}{{ .Role }}{{ end }}</span>
        <span class="tl-date">{{ .Start }}{{ if .End }} – {{ .End }}{{ end }}</span>
      </div>
      <div class="tl-org">{{ if .OrgURL }}<a href="{{ .OrgURL }}">{{ .Org }}</a>{{ else }}{{ .Org }}{{ end }}{{ if .Mode }} · {{ .Mode }}{{ end }}</div>
      {{- if .Intro }}<p>{{ emph .Intro }}</p>{{ end }}
      {{- if .Bullets }}<ul>{{ range .Bullets }}<li>{{ emph . }}</li>{{ end }}</ul>{{ end }}
    </div>
  {{- end }}
  </div>
</section>

<section>
  <h2>Education</h2>
  <div class="timeline">
  {{- range .Education }}
    <div class="tl-item">
      <div class="tl-head">
        <span class="tl-role">{{ if .URL }}<a href="{{ .URL }}">{{ .Role }}</a>{{ else }}{{ .Role }}{{ end }}</span>
        <span class="tl-date">{{ .Start }}{{ if .End }} – {{ .End }}{{ end }}</span>
      </div>
      <div class="tl-org">{{ .Org }} · {{ .Location }}</div>
      {{- if .Note }}<p><em>{{ emph .Note }}</em></p>{{ end }}
    </div>
  {{- end }}
  </div>
</section>

<section class="skills">
  <h2>Skills</h2>
  {{- range .Skills }}
  <div><span class="cat">{{ .Category }}:</span> {{ range $i, $it := .Items }}{{ if $i }}, {{ end }}{{ $it }}{{ end }}</div>
  {{- end }}
</section>

<section class="skills">
  <h2>Languages</h2>
  {{- range .Languages }}
  <div><span class="cat">{{ .Name }}:</span> {{ .Level }}</div>
  {{- end }}
</section>

<section>
  <h2>Certificates &amp; Awards</h2>
  <div class="timeline">
  {{- range .Certificates }}
    <div class="tl-item">
      <div class="tl-head">
        <span class="tl-role">{{ if .CertURL }}<a href="{{ .CertURL }}">{{ .Title }}</a>{{ else }}{{ .Title }}{{ end }}</span>
        <span class="tl-date">{{ .Start }}{{ if .End }} – {{ .End }}{{ end }}</span>
      </div>
      <div class="tl-org">{{ if .OrgURL }}<a href="{{ .OrgURL }}">{{ .Org }}</a>{{ else }}{{ .Org }}{{ end }}{{ if .Group }} · {{ .Group }}{{ end }}</div>
    </div>
  {{- end }}
  </div>
</section>

{{- if .Activities }}
<section>
  <h2>Activities</h2>
  <div class="timeline">
  {{- range .Activities }}
    <div class="tl-item">
      <div class="tl-head">
        <span class="tl-role">{{ .Role }}</span>
        <span class="tl-date">{{ .Start }}{{ if .End }} – {{ .End }}{{ end }}</span>
      </div>
      <div class="tl-org">{{ .Org }} · {{ .Location }}</div>
      {{- if .Intro }}<p>{{ emph .Intro }}</p>{{ end }}
    </div>
  {{- end }}
  </div>
</section>
{{- end }}
</body>
</html>
```

- [ ] **Step 3: Implement the renderer**

Create `internal/html/html.go`:
```go
package html

import (
	"bytes"
	"embed"
	"html/template"
	"os"
	"strings"

	"github.com/georgemessiha22/georgemessiha22/internal/model"
	"github.com/georgemessiha22/georgemessiha22/internal/render"
)

//go:embed templates/index.html.tmpl templates/style.css
var assetFS embed.FS

// Renderer renders the resume to a self-contained mini-site.
type Renderer struct{}

// viewModel adds derived fields the template needs.
type viewModel struct {
	model.Resume
	HasPhoto bool
}

// emphasize escapes text (html) and wraps **bold** spans in <strong>.
func emphasize(s string) template.HTML {
	var b strings.Builder
	for _, sp := range render.ParseSpans(s) {
		esc := template.HTMLEscapeString(sp.Text)
		if sp.Bold {
			b.WriteString("<strong>" + esc + "</strong>")
		} else {
			b.WriteString(esc)
		}
	}
	return template.HTML(b.String())
}

func (Renderer) Render(r model.Resume, v model.Variant) ([]render.Artifact, error) {
	rr := r.ForVariant(v)

	var artifacts []render.Artifact
	hasPhoto := false
	if rr.Contact.Photo != "" {
		if data, err := os.ReadFile(rr.Contact.Photo); err == nil {
			hasPhoto = true
			artifacts = append(artifacts, render.Artifact{RelPath: "profile.jpg", Bytes: data})
		}
	}

	t, err := template.New("index.html.tmpl").Funcs(template.FuncMap{
		"emph": emphasize,
	}).ParseFS(assetFS, "templates/index.html.tmpl")
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, viewModel{Resume: rr, HasPhoto: hasPhoto}); err != nil {
		return nil, err
	}
	artifacts = append(artifacts, render.Artifact{RelPath: "index.html", Bytes: buf.Bytes()})

	css, err := assetFS.ReadFile("templates/style.css")
	if err != nil {
		return nil, err
	}
	artifacts = append(artifacts, render.Artifact{RelPath: "style.css", Bytes: css})

	return artifacts, nil
}
```

- [ ] **Step 4: Write the golden test (index.html only; photo absent in test)**

Create `internal/html/html_test.go`:
```go
package html

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/georgemessiha22/georgemessiha22/internal/model"
)

var update = flag.Bool("update", false, "update golden files")

func fixture() model.Resume {
	return model.Resume{
		Contact: model.Contact{
			Firstname: "George", Lastname: "Messiha", Title: "Lead SWE",
			Email: "g@example.com", Location: "Dubai, UAE",
			Photo:   "", // no file => no image artifact, deterministic
			Socials: model.Socials{GitHub: "georgemessiha22", LinkedIn: "georgemessiha22"},
		},
		Summary: "A summary with **bold** & <tags>.",
		Experience: []model.Entry{{
			Role: "Senior Software Engineer", Org: "HungerStation",
			OrgURL: "https://hungerstation.com", Location: "Dubai, UAE", Mode: "Hybrid",
			Start: "Nov 2023", End: "Present", Intro: "Did things in **Go**.",
			Bullets: []string{"Built **Go** services."},
		}},
		Education: []model.Entry{{Role: "BSc", Org: "GUC", Location: "Cairo, EG", Start: "2011", End: "2016", Note: "Thesis: Cloud"}},
		Skills:    []model.SkillGroup{{Category: "Languages", Items: []string{"Go", "Python"}}},
		Languages: []model.Language{{Name: "Arabic", Level: "Native"}},
		Certificates: []model.Cert{{Title: "PM Foundation", Org: "Google", CertURL: "https://x/c.pdf", Start: "2024", End: "2024"}},
	}
}

func TestRenderMiniSite(t *testing.T) {
	got, err := Renderer{}.Render(fixture(), model.Detailed)
	if err != nil {
		t.Fatal(err)
	}
	// Expect index.html + style.css (no profile.jpg because Photo is empty).
	var index []byte
	names := map[string]bool{}
	for _, a := range got {
		names[a.RelPath] = true
		if a.RelPath == "index.html" {
			index = a.Bytes
		}
	}
	if !names["index.html"] || !names["style.css"] {
		t.Fatalf("missing artifacts: %v", names)
	}
	if names["profile.jpg"] {
		t.Fatalf("did not expect profile.jpg when Photo is empty")
	}
	golden := filepath.Join("testdata", "index.html.golden")
	if *update {
		if err := os.WriteFile(golden, index, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if string(index) != string(want) {
		t.Fatalf("index.html differs from golden; run with -update if intended")
	}
}
```

- [ ] **Step 5: Generate golden + verify escaping**

Run:
```bash
mkdir -p internal/html/testdata
go test ./internal/html/ -run TestRenderMiniSite -update
go test ./internal/html/ -run TestRenderMiniSite -v
```
Expected: golden written; second run PASSes. Open the golden and confirm `&amp;` and `&lt;tags&gt;` appear (auto-escaped) and `<strong>Go</strong>` is present.

- [ ] **Step 6: Commit**

```bash
git add internal/html/
git commit -m "feat(html): add timeline mini-site renderer with golden test"
```

---

## Task 8: writer — persist artifacts

**Files:**
- Create: `internal/writer/writer.go`
- Test: `internal/writer/writer_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/writer/writer_test.go`:
```go
package writer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/georgemessiha22/georgemessiha22/internal/render"
)

func TestWrite(t *testing.T) {
	dir := t.TempDir()
	arts := []render.Artifact{
		{RelPath: "resume.typ", Bytes: []byte("hello")},
		{RelPath: "sub/style.css", Bytes: []byte("body{}")},
	}
	if err := Write(dir, arts); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "resume.typ"))
	if err != nil || string(got) != "hello" {
		t.Fatalf("resume.typ wrong: %q err=%v", got, err)
	}
	got2, err := os.ReadFile(filepath.Join(dir, "sub", "style.css"))
	if err != nil || string(got2) != "body{}" {
		t.Fatalf("nested file wrong: %q err=%v", got2, err)
	}
}
```

Run: `go test ./internal/writer/ -v`
Expected: FAIL (undefined `Write`).

- [ ] **Step 2: Implement writer**

Create `internal/writer/writer.go`:
```go
package writer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/georgemessiha22/georgemessiha22/internal/render"
)

// Write persists artifacts under dir, creating subdirectories as needed.
func Write(dir string, artifacts []render.Artifact) error {
	for _, a := range artifacts {
		dest := filepath.Join(dir, a.RelPath)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("create dir for %s: %w", a.RelPath, err)
		}
		if err := os.WriteFile(dest, a.Bytes, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", a.RelPath, err)
		}
	}
	return nil
}
```

Run: `go test ./internal/writer/ -v`
Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add internal/writer/
git commit -m "feat(writer): persist rendered artifacts to disk"
```

---

## Task 9: CLI entrypoint

**Files:**
- Create: `cmd/resume/main.go`

- [ ] **Step 1: Implement the CLI**

Create `cmd/resume/main.go`:
```go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/georgemessiha22/georgemessiha22/internal/config"
	"github.com/georgemessiha22/georgemessiha22/internal/html"
	"github.com/georgemessiha22/georgemessiha22/internal/latex"
	"github.com/georgemessiha22/georgemessiha22/internal/model"
	"github.com/georgemessiha22/georgemessiha22/internal/render"
	"github.com/georgemessiha22/georgemessiha22/internal/typst"
	"github.com/georgemessiha22/georgemessiha22/internal/writer"
)

func usage() {
	fmt.Fprint(os.Stderr, `resume - generate resume sources from resume.yaml

Usage:
  resume <command> [flags]

Commands:
  typst   Generate Typst (.typ) source
  tex     Generate LaTeX (.tex) source
  html    Generate the HTML mini-site
  all     Generate typst+tex (both variants) and the html site

Flags:
  --input   path to YAML (default resume.yaml)
  --variant summary|detailed (default summary; html forces detailed)
  --out     output directory (default build/gen; site default site)
`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	input := fs.String("input", "resume.yaml", "path to YAML input")
	variant := fs.String("variant", "summary", "summary|detailed")
	out := fs.String("out", "", "output directory")
	_ = fs.Parse(os.Args[2:])

	if err := run(cmd, *input, *variant, *out); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(cmd, input, variant, out string) error {
	r, err := config.Load(input)
	if err != nil {
		return err
	}

	defGen := "build/gen"
	defSite := "site"

	switch cmd {
	case "typst":
		return generate(typst.Renderer{}, r, parseVariant(variant), pick(out, defGen))
	case "tex":
		return generate(latex.Renderer{}, r, parseVariant(variant), pick(out, defGen))
	case "html":
		return generate(html.Renderer{}, r, model.Detailed, pick(out, defSite))
	case "all":
		if err := generate(typst.Renderer{}, r, model.Summary, defGen); err != nil {
			return err
		}
		if err := generate(typst.Renderer{}, r, model.Detailed, defGen); err != nil {
			return err
		}
		if err := generate(latex.Renderer{}, r, model.Summary, defGen); err != nil {
			return err
		}
		if err := generate(latex.Renderer{}, r, model.Detailed, defGen); err != nil {
			return err
		}
		return generate(html.Renderer{}, r, model.Detailed, defSite)
	default:
		usage()
		return fmt.Errorf("unknown command %q", cmd)
	}
}

func generate(rd render.Renderer, r model.Resume, v model.Variant, dir string) error {
	arts, err := rd.Render(r, v)
	if err != nil {
		return err
	}
	return writer.Write(dir, arts)
}

func parseVariant(s string) model.Variant {
	if s == string(model.Detailed) {
		return model.Detailed
	}
	return model.Summary
}

func pick(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
```

- [ ] **Step 2: Build and smoke-test the whole pipeline**

Run:
```bash
go build ./...
go vet ./...
go test ./...
```
Expected: build OK, vet clean, all tests PASS.

- [ ] **Step 3: Commit**

```bash
git add cmd/
git commit -m "feat(cli): add resume CLI wiring renderers via Renderer interface"
```

---

## Task 10: Author `resume.yaml` (full content migration)

**Files:**
- Create: `resume.yaml`

This transcribes ALL existing content. Source of truth for the data is the
existing `.tex`/`.typ` files (HungerStation, Fincompare, Creatio, Huspy, Orange,
PayMob, Nexus Edge, BlueCloud, Aprcot, internships GUC/Orange/Xerox; education
GUC + Thanwya Amma; the full skills/languages; certs incl. the full DataCamp
course/track list; activities). Summary variant = the roles in
`George_Messiha_Resume.tex` (HungerStation, Fincompare, Creatio, Orange, PayMob,
BlueCloud, Aprcot) + one cert each; detailed = everything.

- [ ] **Step 1: Write resume.yaml**

Create `resume.yaml` using these rules:
- `contact`: George Messiha, "Lead Software Engineer", email `georgemessiha22@gmail.com`, phone `(+971) 54 555 1032`, location `Dubai, UAE`, photo `pictures/61673.jpg`, socials github/gitlab/linkedin = `georgemessiha22`.
- `summary`: the paragraph from `summary.tex`.
- `experience`: one entry per role (fields role/org/org_url/location/mode/start/end/intro/bullets), tagging the seven summary roles with `variants: [summary, detailed]` and the detailed-only roles (Huspy, Nexus Edge) and the three internships with `variants: [detailed]`. Internship entries set `url` to their certificate link and `mode: ""`.
- `education`: GUC bachelor (with `note: "Thesis: Security of Cloud Computing"`, `url` to the scan) and Thanwya Amma (with `url`). Both untagged (appear in both).
- `skills`: the six categories from `skills/main.tex`.
- `languages`: Arabic/Native, English/Fluent.
- `certificates`: PM Foundation (Google), Design Thinking (SkillSoft), PeopleCert DevOps, Deep Learning (Udacity), Data Analyst (Udacity), DataCamp Data Scientist track — all `variants: [summary, detailed]`. Then the Android cert and the full DataCamp course list (`group: Courses`) and track list (`group: Tracks`) tagged `variants: [detailed]`. Use the exact titles and `cert_url`s from `certificates_awards/**` (e.g. PeopleCert uses the shortened `https://candidate.peoplecert.org/`).
- `activities`: the four NGO/teaching entries from `Activities/activities.tex`, untagged-detailed via `variants: [detailed]`, mapping `role`/`org`/`location`/`start`/`end`/`intro`.

Reference the exact strings already present in:
`summary.tex`, `Experience/*.tex`, `Experience/internships/**`, `Education/**`,
`skills/*.tex`, `certificates_awards/**`, `Activities/activities.tex`.

- [ ] **Step 2: Generate everything and verify it parses**

Run:
```bash
go run ./cmd/resume all --input resume.yaml
ls build/gen site
```
Expected: `build/gen/{resume.typ,resume_detailed.typ,resume.tex,resume_detailed.tex}` and `site/{index.html,style.css,profile.jpg}` exist; no error.

- [ ] **Step 3: Compile the generated Typst to PDF (visual check)**

Run:
```bash
mkdir -p dist build/gen/assets
cp pictures/61673.jpg build/gen/assets/profile.jpg
typst compile --font-path typst/fonts build/gen/resume.typ dist/George_Messiha_Resume_v2.pdf
typst compile --font-path typst/fonts build/gen/resume_detailed.typ dist/George_Messiha_detailed_resume_v2.pdf
```
Expected: both PDFs build. (Note: the Typst template references `assets/profile.jpg` relative to the `.typ`; the Makefile in Task 11 handles copying the asset — here we copy manually to verify.)

- [ ] **Step 4: Open site/index.html and eyeball the timeline**

Run: `open site/index.html`
Expected: header with name/photo/contacts, timeline sections for experience/education/certs/activities, summary paragraph. Fix YAML/templates if anything is missing.

- [ ] **Step 5: Commit**

```bash
git add resume.yaml
git commit -m "feat: add resume.yaml as the single source of truth"
```

---

## Task 11: Rewire the Makefile

**Files:**
- Modify: `Makefile`

The generated Typst expects `assets/profile.jpg` next to the `.typ`; the generated
LaTeX uses `\photo{pictures/...}`? No — moderncv `\photo` needs a path. To keep
parity with the original, the LaTeX template does NOT set `\photo` (header without
photo) to avoid path issues; if a photo is desired later it can be added. The
Typst targets copy the profile image into `build/gen/assets/`.

- [ ] **Step 1: Replace the Makefile**

Replace `Makefile` with:
```make
VERSION_TAG=$(shell date +'%y.%m.%d')
PWD=$(shell pwd)
GEN=build/gen
TYPST_FONTS=typst/fonts

.PHONY: build bash gen resume detailed typst typst-resume typst-detailed site clean all

build:
	docker build -f build/Dockerfile -t tex:latest -t tex:$(VERSION_TAG) .

bash:
	docker run --rm -it -v $(PWD):/data tex:latest bash

# Generate all source files (typst + tex, both variants) and the html site.
gen:
	go run ./cmd/resume all --input resume.yaml
	mkdir -p $(GEN)/assets
	cp pictures/61673.jpg $(GEN)/assets/profile.jpg

# LaTeX PDFs (compiled in Docker) from generated .tex.
resume: gen
	docker run --rm -v $(PWD):/data tex:latest pdflatex -output-directory dist $(GEN)/resume.tex
	mv dist/resume.pdf dist/George_Messiha_Resume.pdf

detailed: gen
	docker run --rm -v $(PWD):/data tex:latest pdflatex -output-directory dist $(GEN)/resume_detailed.tex
	mv dist/resume_detailed.pdf dist/George_Messiha_detailed_resume.pdf

# Typst PDFs from generated .typ.
typst-resume: gen
	mkdir -p dist
	typst compile --font-path $(TYPST_FONTS) $(GEN)/resume.typ dist/George_Messiha_Resume_v2.pdf

typst-detailed: gen
	mkdir -p dist
	typst compile --font-path $(TYPST_FONTS) $(GEN)/resume_detailed.typ dist/George_Messiha_detailed_resume_v2.pdf

typst: typst-resume typst-detailed

# HTML mini-site (detailed) for GitHub Pages.
site: gen

clean:
	rm -f dist/*.aux dist/*.out dist/*.log

all: build resume detailed typst site clean
```

- [ ] **Step 2: Verify Typst + site targets (local, no Docker)**

Run:
```bash
make typst
make site
ls dist/*_v2.pdf site/index.html
```
Expected: both `_v2.pdf` files and `site/index.html` exist.

- [ ] **Step 3: Commit**

```bash
git add Makefile
git commit -m "build: drive Make targets from the Go generator"
```

---

## Task 12: Update CI — Go build, PDFs, and GitHub Pages deploy

**Files:**
- Modify: `.github/workflows/release.yml`

- [ ] **Step 1: Replace the workflow**

Replace `.github/workflows/release.yml` with:
```yaml
name: Release

on:
  push:
    branches:
      - master
      - releases

permissions:
  contents: write
  pages: write
  id-token: write

jobs:
  release:
    name: Release
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          submodules: recursive
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.26"
      - name: Setup Typst
        uses: typst-community/setup-typst@v3
        with:
          typst-version: "0.15.0"
      - name: Generate sources
        run: make gen
      - name: Build Typst PDFs
        run: make typst
      - name: Build LaTeX PDFs
        run: make build resume detailed
      - name: Build site
        run: make site
      - name: Bump version and push tag
        uses: laputansoft/github-tag-action@v4.6
        id: tags
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          create_annotated_tag: true
      - name: Release
        uses: softprops/action-gh-release@v2.2.1
        with:
          body: ${{ steps.tags.outputs.changelog }}
          tag_name: ${{ steps.tags.outputs.new_tag }}
          draft: false
          prerelease: false
          generate_release_notes: true
          files: |
            dist/George_Messiha_Resume.pdf
            dist/George_Messiha_detailed_resume.pdf
            dist/George_Messiha_Resume_v2.pdf
            dist/George_Messiha_detailed_resume_v2.pdf
      - name: Upload Pages artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: site

  deploy-pages:
    name: Deploy Pages
    needs: release
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

- [ ] **Step 2: Validate YAML**

Run: `ruby -ryaml -e "YAML.load_file('.github/workflows/release.yml'); puts 'valid'"`
Expected: `valid`.

- [ ] **Step 3: Commit**

```bash
git add .github/workflows/release.yml
git commit -m "ci: generate via Go, attach 4 PDFs, deploy site to Pages"
```

---

## Task 13: Remove obsolete hand-written source files

**Files:**
- Delete: all `*.tex` content/section files and the two top-level `.tex`
- Delete: `typst/sections/*.typ`, `typst/resume.typ`, `typst/detailed_resume.typ`, `typst/assets/`
- **Keep:** `typst/fonts/`, the `moderncv` submodule, `pictures/`, and **all certificate/document PDFs & images** under `Education/` and `certificates_awards/` (the resume links point to these files on GitHub — deleting them breaks the links).

> CRITICAL: directories like `certificates_awards/` and `Education/` contain the
> actual scanned certificate PDFs/images that the resume hyperlinks to. Only the
> `.tex` files are obsolete. Do NOT `git rm -r` these directories — remove `.tex`
> files only.

- [ ] **Step 1: Delete only the .tex files (preserve PDFs/images)**

Run:
```bash
# top-level tex
git rm George_Messiha_Resume.tex George_Messiha_detailed_resume.tex summary.tex
# contact/ contained only main.tex
git rm -r contact
# every other .tex anywhere, leaving PDFs/JPGs/PNGs intact
git rm $(git ls-files '*.tex')
```
Expected: only `.tex` files removed. Verify the certificate documents remain:
```bash
git ls-files 'certificates_awards/*' | head
git ls-files 'Education/*'
```
Expected: the `.pdf`/`.jpg`/`.png` files are still tracked.

- [ ] **Step 2: Delete hand-written Typst sources (keep fonts)**

Run:
```bash
git rm typst/resume.typ typst/detailed_resume.typ typst/README.md
git rm -r typst/sections typst/assets
```
(Note: `build/gen` now produces the `.typ` files; fonts remain under `typst/fonts/`.)

- [ ] **Step 3: Verify nothing references deleted files**

Run:
```bash
go test ./...
make typst && make site
grep -rn "subfile\|typst/sections" Makefile .github || echo "no stale refs"
git ls-files '*.tex' || true
```
Expected: tests pass; PDFs + site build; no stale references; **no `.tex` files remain** tracked.

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "chore: remove hand-written tex/typst sources superseded by resume.yaml"
```

---

## Task 14: Update GUIDE.md for forkers

**Files:**
- Modify: `GUIDE.md`

- [ ] **Step 1: Rewrite GUIDE.md**

Replace `GUIDE.md` with:
```markdown
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

1. **Fork & clone with submodules** (LaTeX needs the `moderncv` submodule):
   ```bash
   git clone --recurse-submodules https://github.com/<you>/<repo>.git
   ```
2. **Edit `resume.yaml`** — change `contact`, `experience`, `education`,
   `skills`, `languages`, `certificates`, `activities`. Tag entries with
   `variants: [summary, detailed]` (omit `variants` to include an entry in both
   resumes).
3. **Replace the photo** at `pictures/61673.jpg` (or change `contact.photo`).
4. **Generate locally** (see Usage). Commit and push to `master`; the release
   workflow builds the PDFs, attaches them to a GitHub Release, and deploys the
   site to Pages.
5. **Enable GitHub Pages**: repo Settings → Pages → Source = "GitHub Actions".

## Usage

### Prerequisites
- [Go](https://go.dev) 1.26+
- [Typst](https://github.com/typst/typst) 0.15+ (for the Typst PDFs)
- Docker (for the LaTeX PDFs)

### Generate sources only
```bash
go run ./cmd/resume all          # typst + tex (both variants) + html site
go run ./cmd/resume typst --variant detailed
go run ./cmd/resume html         # site/ (detailed)
```

### Build PDFs and site (Makefile)
```bash
make typst       # Typst PDFs  -> dist/*_v2.pdf
make resume      # LaTeX summary PDF (Docker) -> dist/George_Messiha_Resume.pdf
make detailed    # LaTeX detailed PDF (Docker)
make site        # HTML mini-site -> site/
make all         # everything
```

Outputs land in `dist/` (PDFs) and `site/` (HTML). Both are git-ignored.

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
```

- [ ] **Step 2: Commit**

```bash
git add GUIDE.md
git commit -m "docs: rewrite GUIDE.md for the YAML/Go workflow and forkers"
```

---

## Final verification

- [ ] **Step 1: Full clean build**

Run:
```bash
go test ./...
rm -rf build site dist
make typst && make site
ls dist/*_v2.pdf site/index.html
```
Expected: tests pass; both `_v2.pdf` files and `site/index.html` produced from scratch.

- [ ] **Step 2: Confirm no secrets and clean status**

Run:
```bash
git status
grep -rnE "[A-Za-z0-9_-]{40,}" resume.yaml && echo "CHECK long string" || echo "no long tokens"
```
Expected: clean tree (after commits), no long credential-like strings in `resume.yaml`.

