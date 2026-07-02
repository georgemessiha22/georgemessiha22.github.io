package html

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
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
		PersonalProjects: model.PersonalProjects{
			Intro: "I love **automation** and the terminal.",
			Links: []model.ProjectLink{
				{Name: "GogoNvim", URL: "https://github.com/georgemessiha22/GogoNvim", Description: "My **Neovim** config.", Icon: "neovim"},
				{Name: "dotfiles", URL: "https://github.com/georgemessiha22/dotfiles", Description: "My terminal setup.", Icon: "terminal"},
			},
		},
		Experience: []model.Entry{{
			Role: "Senior Software Engineer", Org: "HungerStation",
			OrgURL: "https://hungerstation.com", Location: "Dubai, UAE", Mode: "Hybrid",
			Start: "Nov 2023", End: "Present", Intro: "Did things in **Go**.",
			Bullets: []string{"Built **Go** services."},
		}},
		Education:    []model.Entry{{Role: "BSc", Org: "GUC", Location: "Cairo, EG", Start: "2011", End: "2016", Note: "Thesis: Cloud"}},
		Skills:       []model.SkillGroup{{Category: "Languages", Items: []string{"Go", "Python"}}},
		Languages:    []model.Language{{Name: "Arabic", Level: "Native"}},
		Certificates: []model.Cert{{Title: "PM Foundation", Org: "Google", CertURL: "https://x/c.pdf", Start: "2024", End: "2024"}},
		ReleasesURL:  "https://github.com/example/example/releases/latest/download",
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
	// Download bar links to the release PDFs and Markdown files.
	for _, want := range []string{
		`class="downloads"`,
		"https://github.com/example/example/releases/latest/download/George_Messiha_Resume.pdf",
		"https://github.com/example/example/releases/latest/download/George_Messiha_detailed_resume_v2.pdf",
		"https://github.com/example/example/releases/latest/download/George_Messiha_Resume.md",
		"https://github.com/example/example/releases/latest/download/George_Messiha_detailed_resume.md",
		`class="projects"`,
		`class="project-list"`,
		`<svg`,
		"https://github.com/georgemessiha22/GogoNvim",
		"https://github.com/georgemessiha22/dotfiles",
		"<strong>Neovim</strong>",
	} {
		if !strings.Contains(string(index), want) {
			t.Fatalf("index.html missing %q", want)
		}
	}
}

func TestRenderMiniSite_NoDownloadsWhenNoReleasesURL(t *testing.T) {
	f := fixture()
	f.ReleasesURL = ""
	got, err := Renderer{}.Render(f, model.Detailed)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range got {
		if a.RelPath == "index.html" && strings.Contains(string(a.Bytes), "class=\"downloads\"") {
			t.Fatalf("download bar should be absent when ReleasesURL is empty")
		}
	}
}

func TestRender_NoBlogNavWhenNoBlogDir(t *testing.T) {
	got, err := Renderer{}.Render(fixture(), model.Detailed)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range got {
		if a.RelPath == "index.html" && strings.Contains(string(a.Bytes), `class="site-nav"`) {
			t.Fatalf("nav should be absent when there is no blog")
		}
		if strings.HasPrefix(a.RelPath, "blog/") {
			t.Fatalf("no blog artifacts expected, got %q", a.RelPath)
		}
	}
}

func TestRender_NoBlogNavWhenBlogDirEmpty(t *testing.T) {
	// A blog dir that exists but yields no posts (e.g. only a .gitkeep file, or
	// folders without an index.md) must behave like having no blog at all.
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".gitkeep"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "not-a-post"), 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := Renderer{BlogDir: dir}.Render(fixture(), model.Detailed)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range got {
		if a.RelPath == "index.html" && strings.Contains(string(a.Bytes), `class="site-nav"`) {
			t.Fatalf("nav should be absent when the blog dir has no posts")
		}
		if strings.HasPrefix(a.RelPath, "blog/") {
			t.Fatalf("no blog artifacts expected, got %q", a.RelPath)
		}
	}
}

func TestRender_WithBlog(t *testing.T) {
	dir := t.TempDir()
	postDir := filepath.Join(dir, "hello-world")
	if err := os.MkdirAll(postDir, 0o755); err != nil {
		t.Fatal(err)
	}
	md := "---\ntitle: Hello World\ndate: \"2026-06-30\"\nsummary: First post.\n---\n# Hi\n\n![cover](cover.png)\n"
	if err := os.WriteFile(filepath.Join(postDir, "index.md"), []byte(md), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(postDir, "cover.png"), []byte("PNG"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := Renderer{BlogDir: dir}.Render(fixture(), model.Detailed)
	if err != nil {
		t.Fatal(err)
	}

	files := map[string]string{}
	for _, a := range got {
		files[a.RelPath] = string(a.Bytes)
	}

	// Expected blog artifacts.
	for _, want := range []string{"blog/index.html", "blog/hello-world/index.html", "blog/hello-world/cover.png"} {
		if _, ok := files[want]; !ok {
			t.Fatalf("missing artifact %q; have %v", want, keys(files))
		}
	}
	// Asset content copied verbatim.
	if files["blog/hello-world/cover.png"] != "PNG" {
		t.Fatalf("asset not copied verbatim")
	}
	// Resume page shows the nav with a Blog tab.
	if !strings.Contains(files["index.html"], `class="site-nav"`) || !strings.Contains(files["index.html"], `href="blog/index.html"`) {
		t.Fatalf("resume page missing blog nav")
	}
	// Blog index lists the post.
	if !strings.Contains(files["blog/index.html"], "Hello World") || !strings.Contains(files["blog/index.html"], `href="hello-world/index.html"`) {
		t.Fatalf("blog index missing post listing")
	}
	// Post page renders the markdown (heading) and the image.
	post := files["blog/hello-world/index.html"]
	if !strings.Contains(post, "<h1") || !strings.Contains(post, `<img src="cover.png"`) {
		t.Fatalf("post page missing rendered content: %s", post)
	}
}

func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
