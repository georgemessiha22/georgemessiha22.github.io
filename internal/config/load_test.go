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
