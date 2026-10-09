package markdown

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
		}, {
			Title: "Importing Data (Part 1)", Org: "DataCamp",
			CertURL: "https://ex.com/Importing Data (Part 1).pdf",
			Group:   "Courses", Start: "2018", End: "2018",
		}},
	}
}

func TestRenderSummaryGolden(t *testing.T) {
	got, err := Renderer{}.Render(fixture(), model.Summary)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].RelPath != "resume.md" {
		t.Fatalf("unexpected artifacts: %+v", got)
	}
	golden := filepath.Join("testdata", "summary.md.golden")
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

func TestRenderDetailedFilename(t *testing.T) {
	got, err := Renderer{}.Render(fixture(), model.Detailed)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].RelPath != "resume_detailed.md" {
		t.Fatalf("expected resume_detailed.md, got %+v", got)
	}
}

func TestURLEncode(t *testing.T) {
	got := urlEncode("https://ex.com/Importing Data (Part 1).pdf")
	want := "https://ex.com/Importing%20Data%20%28Part%201%29.pdf"
	if got != want {
		t.Fatalf("urlEncode: want %q got %q", want, got)
	}
}
