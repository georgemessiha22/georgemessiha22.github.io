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
