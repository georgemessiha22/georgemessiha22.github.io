package markdown

import (
	"bytes"
	"embed"
	"strings"
	"text/template"

	"github.com/georgemessiha22/georgemessiha22/internal/model"
	"github.com/georgemessiha22/georgemessiha22/internal/render"
)

//go:embed templates/resume.md.tmpl
var tmplFS embed.FS

// Renderer renders the resume to a single GitHub-flavored Markdown file.
type Renderer struct{}

// escape escapes characters that would otherwise be interpreted as Markdown
// formatting when they appear in plain resume text.
func escape(s string) string {
	r := strings.NewReplacer(
		`\`, `\\`,
		"`", "\\`",
		`*`, `\*`,
		`_`, `\_`,
		`[`, `\[`,
		`]`, `\]`,
		`<`, `\<`,
		`>`, `\>`,
		`|`, `\|`,
		`#`, `\#`,
	)
	return r.Replace(s)
}

// emphasize escapes text and converts **bold** spans to Markdown **bold**.
func emphasize(s string) string {
	var b strings.Builder
	for _, sp := range render.ParseSpans(s) {
		if sp.Bold {
			b.WriteString("**" + escape(sp.Text) + "**")
		} else {
			b.WriteString(escape(sp.Text))
		}
	}
	return b.String()
}

// urlEncode encodes the characters in a URL that would break a Markdown
// inline link target `[text](url)` — spaces and parentheses. Existing percent
// escapes and other characters are left untouched.
func urlEncode(s string) string {
	r := strings.NewReplacer(
		" ", "%20",
		"(", "%28",
		")", "%29",
	)
	return r.Replace(s)
}

// contactLine builds a single line of contact details separated by " · ".
func contactLine(r model.Resume) string {
	c := r.Contact
	var parts []string
	if c.Location != "" {
		parts = append(parts, escape(c.Location))
	}
	if c.Email != "" {
		parts = append(parts, "["+escape(c.Email)+"](mailto:"+c.Email+")")
	}
	if c.Phone != "" {
		parts = append(parts, escape(c.Phone))
	}
	if c.Socials.GitHub != "" {
		parts = append(parts, "[GitHub](https://github.com/"+c.Socials.GitHub+")")
	}
	if c.Socials.GitLab != "" {
		parts = append(parts, "[GitLab](https://gitlab.com/"+c.Socials.GitLab+")")
	}
	if c.Socials.LinkedIn != "" {
		parts = append(parts, "[LinkedIn](https://linkedin.com/in/"+c.Socials.LinkedIn+")")
	}
	return strings.Join(parts, " · ")
}

func (Renderer) Render(r model.Resume, v model.Variant) ([]render.Artifact, error) {
	rr := r.ForVariant(v)
	t, err := template.New("resume.md.tmpl").Funcs(template.FuncMap{
		"md":          escape,
		"emph":        emphasize,
		"contactline": contactLine,
		"urlenc":      urlEncode,
	}).ParseFS(tmplFS, "templates/resume.md.tmpl")
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, rr); err != nil {
		return nil, err
	}
	name := "resume.md"
	if v == model.Detailed {
		name = "resume_detailed.md"
	}
	return []render.Artifact{{RelPath: name, Bytes: buf.Bytes()}}, nil
}
