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

// ShortURL strips the scheme and trailing slash so a URL can be printed
// compactly as visible text (e.g. "https://example.com/" -> "example.com").
func ShortURL(s string) string {
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "http://")
	return strings.TrimRight(s, "/")
}

// ContactItems returns the header contact details as plain text, with full
// visible URLs so ATS parsers can read them (no hyperlink-only text).
func ContactItems(c model.Contact, siteURL string) []string {
	var out []string
	add := func(s string) {
		if s != "" {
			out = append(out, s)
		}
	}
	add(c.Email)
	add(c.Phone)
	add(siteURL)
	if c.Socials.GitHub != "" {
		add("https://github.com/" + c.Socials.GitHub)
	}
	if c.Socials.GitLab != "" {
		add("https://gitlab.com/" + c.Socials.GitLab)
	}
	if c.Socials.LinkedIn != "" {
		add("https://linkedin.com/in/" + strings.TrimPrefix(c.Socials.LinkedIn, "in/"))
	}
	return out
}
