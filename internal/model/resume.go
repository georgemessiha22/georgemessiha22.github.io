package model

// Variant identifies which resume a piece of content belongs to.
type Variant string

const (
	Summary  Variant = "summary"
	Detailed Variant = "detailed"
)

// Contact is the header / personal information.
type Contact struct {
	Firstname string
	Lastname  string
	Title     string
	Email     string
	Phone     string
	Location  string
	Photo     string // path to image file
	Socials   Socials
}

// Socials holds usernames (not full URLs) per platform.
type Socials struct {
	GitHub   string
	GitLab   string
	LinkedIn string
}

// Entry is a timeline item used for experience, education, and activities.
type Entry struct {
	Role     string // job title / degree / activity role
	Org      string
	OrgURL   string
	URL      string // optional link on the title (e.g. certificate scan)
	Location string
	Mode     string // e.g. "Hybrid", "Remote", "Fulltime"
	Start    string
	End      string
	Note     string // e.g. thesis line
	Intro    string // paragraph before bullets
	Bullets  []string
	Variants []Variant
}

// SkillGroup is a labelled list of skills.
type SkillGroup struct {
	Category string
	Items    []string
}

// Language is a spoken language and proficiency.
type Language struct {
	Name  string
	Level string
}

// Cert is a certificate or award.
type Cert struct {
	Title    string
	Org      string
	OrgURL   string
	CertURL  string // link to the certificate scan
	Location string
	Start    string
	End      string
	Group    string // optional subsection, e.g. "Courses" / "Tracks"
	Variants []Variant
}

// Resume is the whole document.
type Resume struct {
	Contact      Contact
	Summary      string
	Experience   []Entry
	Education    []Entry
	Skills       []SkillGroup
	Languages    []Language
	Certificates []Cert
	Activities   []Entry
	// ReleasesURL is the base URL for downloadable release assets, e.g.
	// "https://github.com/owner/repo/releases/latest/download". When set, the
	// HTML site shows PDF download links. Optional.
	ReleasesURL string
}

func includesVariant(vs []Variant, v Variant) bool {
	if len(vs) == 0 {
		return true // untagged => appears in every variant
	}
	for _, x := range vs {
		if x == v {
			return true
		}
	}
	return false
}

func filterEntries(in []Entry, v Variant) []Entry {
	out := make([]Entry, 0, len(in))
	for _, e := range in {
		if includesVariant(e.Variants, v) {
			out = append(out, e)
		}
	}
	return out
}

func filterCerts(in []Cert, v Variant) []Cert {
	out := make([]Cert, 0, len(in))
	for _, c := range in {
		if includesVariant(c.Variants, v) {
			out = append(out, c)
		}
	}
	return out
}

// ForVariant returns a copy of the resume containing only the content that
// belongs to variant v. Skills and languages always pass through.
func (r Resume) ForVariant(v Variant) Resume {
	out := r
	out.Experience = filterEntries(r.Experience, v)
	out.Education = filterEntries(r.Education, v)
	out.Activities = filterEntries(r.Activities, v)
	out.Certificates = filterCerts(r.Certificates, v)
	return out
}
