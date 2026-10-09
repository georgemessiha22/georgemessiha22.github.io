// Package webdata renders the resume into a single JSON file consumed by the
// SvelteKit site (in web/). Unlike the other renderers it ignores the variant
// argument and always emits both the summary and detailed variants so the site
// can toggle between them client-side. It also copies the contact photo into
// the SvelteKit static directory.
package webdata

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/georgemessiha22/georgemessiha22/internal/model"
	"github.com/georgemessiha22/georgemessiha22/internal/render"
)

// Renderer emits resume.json (+ profile.jpg) for the SvelteKit site.
type Renderer struct{}

// Artifact-relative paths inside the web/ output directory.
const (
	dataPath  = "src/lib/data/resume.json"
	photoPath = "static/profile.jpg"
	photoName = "profile.jpg"
)

type socialsJSON struct {
	GitHub   string `json:"github,omitempty"`
	GitLab   string `json:"gitlab,omitempty"`
	LinkedIn string `json:"linkedin,omitempty"`
}

type contactJSON struct {
	Firstname string      `json:"firstname"`
	Lastname  string      `json:"lastname"`
	Title     string      `json:"title,omitempty"`
	Email     string      `json:"email,omitempty"`
	Phone     string      `json:"phone,omitempty"`
	Location  string      `json:"location,omitempty"`
	Photo     string      `json:"photo,omitempty"`
	Socials   socialsJSON `json:"socials"`
}

type projectLinkJSON struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

type personalProjectsJSON struct {
	Intro string            `json:"intro,omitempty"`
	Links []projectLinkJSON `json:"links,omitempty"`
}

type entryJSON struct {
	Role     string   `json:"role"`
	Org      string   `json:"org,omitempty"`
	OrgURL   string   `json:"orgUrl,omitempty"`
	URL      string   `json:"url,omitempty"`
	Location string   `json:"location,omitempty"`
	Mode     string   `json:"mode,omitempty"`
	Start    string   `json:"start,omitempty"`
	End      string   `json:"end,omitempty"`
	Note     string   `json:"note,omitempty"`
	Intro    string   `json:"intro,omitempty"`
	Bullets  []string `json:"bullets,omitempty"`
}

type certJSON struct {
	Title    string `json:"title"`
	Org      string `json:"org,omitempty"`
	OrgURL   string `json:"orgUrl,omitempty"`
	CertURL  string `json:"certUrl,omitempty"`
	Location string `json:"location,omitempty"`
	Start    string `json:"start,omitempty"`
	End      string `json:"end,omitempty"`
	Group    string `json:"group,omitempty"`
}

type skillJSON struct {
	Category string   `json:"category"`
	Items    []string `json:"items"`
}

type languageJSON struct {
	Name  string `json:"name"`
	Level string `json:"level,omitempty"`
}

type variantJSON struct {
	Experience   []entryJSON `json:"experience"`
	Education    []entryJSON `json:"education"`
	Certificates []certJSON  `json:"certificates"`
	Activities   []entryJSON `json:"activities"`
}

type downloadJSON struct {
	Label string `json:"label"`
	Href  string `json:"href"`
	Kind  string `json:"kind"`  // "PDF" | "MD"
	Group string `json:"group"` // "Summary" | "Detailed"
}

type resumeJSON struct {
	Contact          contactJSON            `json:"contact"`
	Summary          string                 `json:"summary,omitempty"`
	PersonalProjects personalProjectsJSON   `json:"personalProjects"`
	ReleasesURL      string                 `json:"releasesUrl,omitempty"`
	SiteURL          string                 `json:"siteUrl,omitempty"`
	Downloads        []downloadJSON         `json:"downloads"`
	Skills           []skillJSON            `json:"skills"`
	Languages        []languageJSON         `json:"languages"`
	Variants         map[string]variantJSON `json:"variants"`
}

// Render builds resume.json (both variants) and, when the contact photo file
// exists, a profile.jpg artifact. The variant argument is ignored.
func (Renderer) Render(r model.Resume, _ model.Variant) ([]render.Artifact, error) {
	doc := build(r)

	// Copy the photo if present; reflect it in contact.photo.
	var arts []render.Artifact
	if r.Contact.Photo != "" {
		if b, err := os.ReadFile(r.Contact.Photo); err == nil {
			doc.Contact.Photo = photoName
			arts = append(arts, render.Artifact{RelPath: filepath.FromSlash(photoPath), Bytes: b})
		}
	}

	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	data = append(data, '\n')
	arts = append(arts, render.Artifact{RelPath: filepath.FromSlash(dataPath), Bytes: data})
	return arts, nil
}

func build(r model.Resume) resumeJSON {
	doc := resumeJSON{
		Contact: contactJSON{
			Firstname: r.Contact.Firstname,
			Lastname:  r.Contact.Lastname,
			Title:     r.Contact.Title,
			Email:     r.Contact.Email,
			Phone:     r.Contact.Phone,
			Location:  r.Contact.Location,
			Socials: socialsJSON{
				GitHub:   r.Contact.Socials.GitHub,
				GitLab:   r.Contact.Socials.GitLab,
				LinkedIn: r.Contact.Socials.LinkedIn,
			},
		},
		Summary: r.Summary,
		PersonalProjects: personalProjectsJSON{
			Intro: r.PersonalProjects.Intro,
			Links: projectLinks(r.PersonalProjects.Links),
		},
		ReleasesURL: r.ReleasesURL,
		SiteURL:     r.SiteURL,
		Downloads:   downloads(r.ReleasesURL),
		Skills:      skills(r.Skills),
		Languages:   languages(r.Languages),
		Variants: map[string]variantJSON{
			string(model.Summary):  variantFor(r, model.Summary),
			string(model.Detailed): variantFor(r, model.Detailed),
		},
	}
	return doc
}

func variantFor(r model.Resume, v model.Variant) variantJSON {
	fr := r.ForVariant(v)
	return variantJSON{
		Experience:   entries(fr.Experience),
		Education:    entries(fr.Education),
		Certificates: certs(fr.Certificates),
		Activities:   entries(fr.Activities),
	}
}

func entries(in []model.Entry) []entryJSON {
	out := make([]entryJSON, 0, len(in))
	for _, e := range in {
		out = append(out, entryJSON{
			Role: e.Role, Org: e.Org, OrgURL: e.OrgURL, URL: e.URL,
			Location: e.Location, Mode: e.Mode, Start: e.Start, End: e.End,
			Note: e.Note, Intro: e.Intro, Bullets: e.Bullets,
		})
	}
	return out
}

func certs(in []model.Cert) []certJSON {
	out := make([]certJSON, 0, len(in))
	for _, c := range in {
		out = append(out, certJSON{
			Title: c.Title, Org: c.Org, OrgURL: c.OrgURL, CertURL: c.CertURL,
			Location: c.Location, Start: c.Start, End: c.End, Group: c.Group,
		})
	}
	return out
}

func skills(in []model.SkillGroup) []skillJSON {
	out := make([]skillJSON, 0, len(in))
	for _, s := range in {
		out = append(out, skillJSON{Category: s.Category, Items: s.Items})
	}
	return out
}

func languages(in []model.Language) []languageJSON {
	out := make([]languageJSON, 0, len(in))
	for _, l := range in {
		out = append(out, languageJSON{Name: l.Name, Level: l.Level})
	}
	return out
}

func projectLinks(in []model.ProjectLink) []projectLinkJSON {
	if len(in) == 0 {
		return nil
	}
	out := make([]projectLinkJSON, 0, len(in))
	for _, l := range in {
		out = append(out, projectLinkJSON{Name: l.Name, URL: l.URL, Description: l.Description, Icon: l.Icon})
	}
	return out
}

// downloads builds the release download list from the base URL. The filenames
// mirror the public asset names produced by the Makefile / release workflow.
func downloads(base string) []downloadJSON {
	if base == "" {
		return nil
	}
	href := func(name string) string { return base + "/" + name }
	return []downloadJSON{
		{Label: "Résumé (PDF)", Href: href("George_Messiha_Resume.pdf"), Kind: "PDF", Group: "Summary"},
		{Label: "Résumé · v2 (Typst)", Href: href("George_Messiha_Resume_v2.pdf"), Kind: "PDF", Group: "Summary"},
		{Label: "Résumé (Markdown)", Href: href("George_Messiha_Resume.md"), Kind: "MD", Group: "Summary"},
		{Label: "Detailed résumé (PDF)", Href: href("George_Messiha_detailed_resume.pdf"), Kind: "PDF", Group: "Detailed"},
		{Label: "Detailed résumé · v2 (Typst)", Href: href("George_Messiha_detailed_resume_v2.pdf"), Kind: "PDF", Group: "Detailed"},
		{Label: "Detailed résumé (Markdown)", Href: href("George_Messiha_detailed_resume.md"), Kind: "MD", Group: "Detailed"},
	}
}
