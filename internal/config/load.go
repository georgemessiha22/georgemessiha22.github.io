package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/georgemessiha22/georgemessiha22/internal/model"
)

// yamlResume mirrors model.Resume with yaml tags. Keeping a separate DTO keeps
// yaml concerns out of the domain model.
type yamlResume struct {
	Contact struct {
		Firstname string `yaml:"firstname"`
		Lastname  string `yaml:"lastname"`
		Title     string `yaml:"title"`
		Email     string `yaml:"email"`
		Phone     string `yaml:"phone"`
		Location  string `yaml:"location"`
		Photo     string `yaml:"photo"`
		Socials   struct {
			GitHub   string `yaml:"github"`
			GitLab   string `yaml:"gitlab"`
			LinkedIn string `yaml:"linkedin"`
		} `yaml:"socials"`
	} `yaml:"contact"`
	Summary      string      `yaml:"summary"`
	Experience   []yamlEntry `yaml:"experience"`
	Education    []yamlEntry `yaml:"education"`
	Skills       []yamlSkill `yaml:"skills"`
	Languages    []yamlLang  `yaml:"languages"`
	Certificates []yamlCert  `yaml:"certificates"`
	Activities   []yamlEntry `yaml:"activities"`
	ReleasesURL  string      `yaml:"releases_url"`
}

type yamlEntry struct {
	Role     string   `yaml:"role"`
	Org      string   `yaml:"org"`
	OrgURL   string   `yaml:"org_url"`
	URL      string   `yaml:"url"`
	Location string   `yaml:"location"`
	Mode     string   `yaml:"mode"`
	Start    string   `yaml:"start"`
	End      string   `yaml:"end"`
	Note     string   `yaml:"note"`
	Intro    string   `yaml:"intro"`
	Bullets  []string `yaml:"bullets"`
	Variants []string `yaml:"variants"`
}

type yamlSkill struct {
	Category string   `yaml:"category"`
	Items    []string `yaml:"items"`
}

type yamlLang struct {
	Name  string `yaml:"name"`
	Level string `yaml:"level"`
}

type yamlCert struct {
	Title    string   `yaml:"title"`
	Org      string   `yaml:"org"`
	OrgURL   string   `yaml:"org_url"`
	CertURL  string   `yaml:"cert_url"`
	Location string   `yaml:"location"`
	Start    string   `yaml:"start"`
	End      string   `yaml:"end"`
	Group    string   `yaml:"group"`
	Variants []string `yaml:"variants"`
}

func toVariants(in []string) []model.Variant {
	if len(in) == 0 {
		return nil
	}
	out := make([]model.Variant, 0, len(in))
	for _, s := range in {
		out = append(out, model.Variant(s))
	}
	return out
}

func toEntry(y yamlEntry) model.Entry {
	return model.Entry{
		Role: y.Role, Org: y.Org, OrgURL: y.OrgURL, URL: y.URL,
		Location: y.Location, Mode: y.Mode, Start: y.Start, End: y.End,
		Note: y.Note, Intro: y.Intro, Bullets: y.Bullets,
		Variants: toVariants(y.Variants),
	}
}

func toEntries(in []yamlEntry) []model.Entry {
	out := make([]model.Entry, 0, len(in))
	for _, y := range in {
		out = append(out, toEntry(y))
	}
	return out
}

// Load reads a YAML file and returns a validated Resume.
func Load(path string) (model.Resume, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return model.Resume{}, fmt.Errorf("read resume file: %w", err)
	}
	var y yamlResume
	if err := yaml.Unmarshal(data, &y); err != nil {
		return model.Resume{}, fmt.Errorf("parse yaml: %w", err)
	}

	r := model.Resume{
		Contact: model.Contact{
			Firstname: y.Contact.Firstname,
			Lastname:  y.Contact.Lastname,
			Title:     y.Contact.Title,
			Email:     y.Contact.Email,
			Phone:     y.Contact.Phone,
			Location:  y.Contact.Location,
			Photo:     y.Contact.Photo,
			Socials: model.Socials{
				GitHub:   y.Contact.Socials.GitHub,
				GitLab:   y.Contact.Socials.GitLab,
				LinkedIn: y.Contact.Socials.LinkedIn,
			},
		},
		Summary:     y.Summary,
		Experience:  toEntries(y.Experience),
		Education:   toEntries(y.Education),
		Activities:  toEntries(y.Activities),
		ReleasesURL: y.ReleasesURL,
	}
	for _, s := range y.Skills {
		r.Skills = append(r.Skills, model.SkillGroup{Category: s.Category, Items: s.Items})
	}
	for _, l := range y.Languages {
		r.Languages = append(r.Languages, model.Language{Name: l.Name, Level: l.Level})
	}
	for _, c := range y.Certificates {
		r.Certificates = append(r.Certificates, model.Cert{
			Title: c.Title, Org: c.Org, OrgURL: c.OrgURL, CertURL: c.CertURL,
			Location: c.Location, Start: c.Start, End: c.End, Group: c.Group,
			Variants: toVariants(c.Variants),
		})
	}

	if err := validate(r); err != nil {
		return model.Resume{}, err
	}
	return r, nil
}

func validate(r model.Resume) error {
	if r.Contact.Firstname == "" && r.Contact.Lastname == "" {
		return fmt.Errorf("validation: contact name is required (firstname/lastname)")
	}
	if len(r.Experience) == 0 {
		return fmt.Errorf("validation: at least one experience entry is required")
	}
	return nil
}
