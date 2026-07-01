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

func newTemplate(name string) (*template.Template, error) {
	return template.New(name).Funcs(template.FuncMap{
		"emph": emphasize,
		// safeHTML embeds already-rendered HTML (goldmark post bodies).
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
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
