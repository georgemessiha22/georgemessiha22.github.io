package config

import (
	"strings"
	"testing"

	"github.com/georgemessiha22/georgemessiha22/internal/model"
)

func TestLoad_Valid(t *testing.T) {
	r, err := Load("testdata/valid.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Contact.Firstname != "George" || r.Contact.Lastname != "Messiha" {
		t.Fatalf("contact not parsed: %+v", r.Contact)
	}
	if len(r.Experience) != 1 || r.Experience[0].Org != "HungerStation" {
		t.Fatalf("experience not parsed: %+v", r.Experience)
	}
	if r.Experience[0].Variants[0] != model.Summary {
		t.Fatalf("variants not parsed: %+v", r.Experience[0].Variants)
	}
	if len(r.Skills) != 1 || r.Skills[0].Category != "Languages" {
		t.Fatalf("skills not parsed: %+v", r.Skills)
	}
}

func TestLoad_PersonalProjects(t *testing.T) {
	r, err := Load("testdata/valid.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.PersonalProjects.Intro != "I love **automation**." {
		t.Fatalf("intro not parsed: %q", r.PersonalProjects.Intro)
	}
	if len(r.PersonalProjects.Links) != 2 {
		t.Fatalf("expected 2 project links, got %d", len(r.PersonalProjects.Links))
	}
	l0 := r.PersonalProjects.Links[0]
	if l0.Name != "GogoNvim" || l0.URL != "https://github.com/georgemessiha22/GogoNvim" || l0.Description != "My **Neovim** config." {
		t.Fatalf("first link not parsed: %+v", l0)
	}
	if l0.Icon != "neovim" {
		t.Fatalf("first link icon not parsed: %q", l0.Icon)
	}
	if r.PersonalProjects.Links[1].Name != "dotfiles" {
		t.Fatalf("second link not parsed: %+v", r.PersonalProjects.Links[1])
	}
}

func TestLoad_MissingName(t *testing.T) {
	_, err := Load("testdata/missing_name.yaml")
	if err == nil || !strings.Contains(err.Error(), "name") {
		t.Fatalf("expected name validation error, got %v", err)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("testdata/does_not_exist.yaml")
	if err == nil {
		t.Fatalf("expected error for missing file")
	}
}
