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
	// Download bar links to the release PDFs.
	for _, want := range []string{
		`class="downloads"`,
		"https://github.com/example/example/releases/latest/download/George_Messiha_Resume.pdf",
		"https://github.com/example/example/releases/latest/download/George_Messiha_detailed_resume_v2.pdf",
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
