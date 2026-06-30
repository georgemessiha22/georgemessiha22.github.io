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
