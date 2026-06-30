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
