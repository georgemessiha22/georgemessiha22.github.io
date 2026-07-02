package html

import (
	"bytes"
	"embed"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/georgemessiha22/georgemessiha22/internal/blog"
	"github.com/georgemessiha22/georgemessiha22/internal/model"
	"github.com/georgemessiha22/georgemessiha22/internal/render"
)

//go:embed templates/index.html.tmpl templates/blog_index.html.tmpl templates/blog_post.html.tmpl templates/style.css
var assetFS embed.FS

// Renderer renders the resume (and optional blog) to a self-contained mini-site.
type Renderer struct {
	// BlogDir is the directory holding blog/<slug>/index.md posts. When empty,
	// no blog pages are generated and the Blog nav tab is hidden.
	BlogDir string
}

// resumeView is the data passed to the resume page template.
type resumeView struct {
	model.Resume
	HasPhoto bool
	HasBlog  bool
}

// blogIndexView is the data passed to the blog listing template.
type blogIndexView struct {
	Resume  model.Resume
	Posts   []blog.Post
	HasBlog bool
}

// blogPostView is the data passed to a single-post template.
type blogPostView struct {
	Resume  model.Resume
	Post    blog.Post
	HasBlog bool
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

func newTemplate(name string) (*template.Template, error) {
	return template.New(name).Funcs(template.FuncMap{
		"emph": emphasize,
		// safeHTML embeds already-rendered HTML (goldmark post bodies).
		"safeHTML":    func(s string) template.HTML { return template.HTML(s) },
		"projectIcon": projectIcon,
	}).ParseFS(assetFS, "templates/"+name)
}

func (rd Renderer) Render(r model.Resume, v model.Variant) ([]render.Artifact, error) {
	rr := r.ForVariant(v)

	// Load blog posts (if a blog dir is configured and exists).
	var posts []blog.Post
	if rd.BlogDir != "" {
		p, err := blog.LoadPosts(rd.BlogDir)
		if err != nil {
			return nil, err
		}
		posts = p
	}
	hasBlog := len(posts) > 0

	var artifacts []render.Artifact

	// Profile photo.
	hasPhoto := false
	if rr.Contact.Photo != "" {
		if data, err := os.ReadFile(rr.Contact.Photo); err == nil {
			hasPhoto = true
			artifacts = append(artifacts, render.Artifact{RelPath: "profile.jpg", Bytes: data})
		}
	}

	// Resume home page.
	home, err := renderPage("index.html.tmpl", resumeView{Resume: rr, HasPhoto: hasPhoto, HasBlog: hasBlog})
	if err != nil {
		return nil, err
	}
	artifacts = append(artifacts, render.Artifact{RelPath: "index.html", Bytes: home})

	// Blog index + per-post pages.
	if hasBlog {
		idx, err := renderPage("blog_index.html.tmpl", blogIndexView{Resume: rr, Posts: posts, HasBlog: true})
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, render.Artifact{RelPath: path.Join("blog", "index.html"), Bytes: idx})

		for _, post := range posts {
			page, err := renderPage("blog_post.html.tmpl", blogPostView{Resume: rr, Post: post, HasBlog: true})
			if err != nil {
				return nil, err
			}
			base := path.Join("blog", post.Slug)
			artifacts = append(artifacts, render.Artifact{RelPath: path.Join(base, "index.html"), Bytes: page})
			for _, a := range post.Assets {
				artifacts = append(artifacts, render.Artifact{RelPath: path.Join(base, a.Name), Bytes: a.Bytes})
			}
		}
	}

	// Stylesheet.
	css, err := assetFS.ReadFile("templates/style.css")
	if err != nil {
		return nil, err
	}
	artifacts = append(artifacts, render.Artifact{RelPath: "style.css", Bytes: css})

	return artifacts, nil
}

func renderPage(name string, data any) ([]byte, error) {
	t, err := newTemplate(name)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
