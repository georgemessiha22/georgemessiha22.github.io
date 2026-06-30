package render

import "testing"

func TestParseSpans(t *testing.T) {
	got := ParseSpans("plain **bold** end")
	want := []Span{
		{Text: "plain ", Bold: false},
		{Text: "bold", Bold: true},
		{Text: " end", Bold: false},
	}
	if len(got) != len(want) {
		t.Fatalf("len: want %d got %d (%+v)", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("span %d: want %+v got %+v", i, want[i], got[i])
		}
	}
}

func TestParseSpans_NoBold(t *testing.T) {
	got := ParseSpans("just text")
	if len(got) != 1 || got[0].Bold || got[0].Text != "just text" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestParseSpans_Adjacent(t *testing.T) {
	got := ParseSpans("**a****b**")
	want := []Span{{Text: "a", Bold: true}, {Text: "b", Bold: true}}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("unexpected: %+v", got)
	}
}
