package webdata

import (
	"encoding/json"
	"testing"

	"github.com/georgemessiha22/georgemessiha22/internal/model"
)

func sampleResume() model.Resume {
	return model.Resume{
		Contact:     model.Contact{Firstname: "George", Lastname: "Messiha", Title: "Lead"},
		Summary:     "A **bold** summary",
		ReleasesURL: "https://example.com/dl",
		Experience: []model.Entry{
			{Role: "Both", Org: "A"},
			{Role: "Detailed only", Org: "B", Variants: []model.Variant{model.Detailed}},
		},
		Skills:    []model.SkillGroup{{Category: "Languages", Items: []string{"Go"}}},
		Languages: []model.Language{{Name: "English", Level: "Fluent"}},
	}
}

func decode(t *testing.T, b []byte) resumeJSON {
	t.Helper()
	var doc resumeJSON
	if err := json.Unmarshal(b, &doc); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	return doc
}

func TestRenderEmitsResumeJSON(t *testing.T) {
	arts, err := Renderer{}.Render(sampleResume(), model.Summary)
	if err != nil {
		t.Fatal(err)
	}
	var data []byte
	for _, a := range arts {
		if a.RelPath == "src/lib/data/resume.json" || a.RelPath == "src\\lib\\data\\resume.json" {
			data = a.Bytes
		}
	}
	if data == nil {
		t.Fatal("resume.json artifact not emitted")
	}
	doc := decode(t, data)

	if doc.Contact.Firstname != "George" {
		t.Errorf("firstname = %q", doc.Contact.Firstname)
	}
	// Both variants present.
	sum, ok := doc.Variants["summary"]
	if !ok {
		t.Fatal("summary variant missing")
	}
	det, ok := doc.Variants["detailed"]
	if !ok {
		t.Fatal("detailed variant missing")
	}
	if len(sum.Experience) != 1 {
		t.Errorf("summary experience = %d, want 1", len(sum.Experience))
	}
	if len(det.Experience) != 2 {
		t.Errorf("detailed experience = %d, want 2", len(det.Experience))
	}
	// Downloads built from ReleasesURL (6 entries).
	if len(doc.Downloads) != 6 {
		t.Errorf("downloads = %d, want 6", len(doc.Downloads))
	}
	if len(doc.Skills) != 1 || len(doc.Languages) != 1 {
		t.Errorf("skills/languages not carried through")
	}
}

func TestNoDownloadsWithoutReleasesURL(t *testing.T) {
	r := sampleResume()
	r.ReleasesURL = ""
	arts, err := Renderer{}.Render(r, model.Summary)
	if err != nil {
		t.Fatal(err)
	}
	var data []byte
	for _, a := range arts {
		if a.RelPath == "src/lib/data/resume.json" || a.RelPath == "src\\lib\\data\\resume.json" {
			data = a.Bytes
		}
	}
	doc := decode(t, data)
	if len(doc.Downloads) != 0 {
		t.Errorf("downloads = %d, want 0 when ReleasesURL empty", len(doc.Downloads))
	}
}
