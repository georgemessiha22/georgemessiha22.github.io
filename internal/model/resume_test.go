package model

import "testing"

func sampleResume() Resume {
	return Resume{
		Contact: Contact{Firstname: "George", Lastname: "Messiha"},
		Summary: "summary text",
		Experience: []Entry{
			{Role: "A", Variants: []Variant{Summary, Detailed}},
			{Role: "B", Variants: []Variant{Detailed}},
			{Role: "C"}, // untagged => both
		},
		Certificates: []Cert{
			{Title: "X", Variants: []Variant{Detailed}},
			{Title: "Y", Variants: []Variant{Summary, Detailed}},
		},
		Activities: []Entry{{Role: "Act", Variants: []Variant{Detailed}}},
	}
}

func TestForVariant_Summary(t *testing.T) {
	got := sampleResume().ForVariant(Summary)
	if len(got.Experience) != 2 {
		t.Fatalf("summary experience: want 2, got %d", len(got.Experience))
	}
	if got.Experience[0].Role != "A" || got.Experience[1].Role != "C" {
		t.Fatalf("unexpected roles: %+v", got.Experience)
	}
	if len(got.Certificates) != 1 || got.Certificates[0].Title != "Y" {
		t.Fatalf("summary certs wrong: %+v", got.Certificates)
	}
	if len(got.Activities) != 0 {
		t.Fatalf("summary activities: want 0, got %d", len(got.Activities))
	}
}

func TestForVariant_Detailed(t *testing.T) {
	got := sampleResume().ForVariant(Detailed)
	if len(got.Experience) != 3 {
		t.Fatalf("detailed experience: want 3, got %d", len(got.Experience))
	}
	if len(got.Certificates) != 2 {
		t.Fatalf("detailed certs: want 2, got %d", len(got.Certificates))
	}
	if len(got.Activities) != 1 {
		t.Fatalf("detailed activities: want 1, got %d", len(got.Activities))
	}
}
