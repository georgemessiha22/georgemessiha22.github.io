package render

import (
	"strings"

	"github.com/georgemessiha22/georgemessiha22/internal/model"
)

// Variant is re-exported so engine packages depend only on render.
type Variant = model.Variant

// Artifact is one rendered file relative to the output directory.
type Artifact struct {
	RelPath string
	Bytes   []byte
}

// Renderer turns a resume + variant into one or more files.
type Renderer interface {
	Render(r model.Resume, v model.Variant) ([]Artifact, error)
}

// Span is a run of text that is either plain or bold.
type Span struct {
	Text string
	Bold bool
}

// ParseSpans splits text on **bold** markers into ordered spans. Empty spans
// (e.g. from adjacent markers) are dropped.
func ParseSpans(s string) []Span {
	var spans []Span
	bold := false
	for len(s) > 0 {
		idx := strings.Index(s, "**")
		if idx < 0 {
			spans = appendSpan(spans, s, bold)
			break
		}
		spans = appendSpan(spans, s[:idx], bold)
		bold = !bold
		s = s[idx+2:]
	}
	return spans
}

func appendSpan(spans []Span, text string, bold bool) []Span {
	if text == "" {
		return spans
	}
	return append(spans, Span{Text: text, Bold: bold})
}
