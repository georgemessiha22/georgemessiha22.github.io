package writer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/georgemessiha22/georgemessiha22/internal/render"
)

func TestWrite(t *testing.T) {
	dir := t.TempDir()
	arts := []render.Artifact{
		{RelPath: "resume.typ", Bytes: []byte("hello")},
		{RelPath: "sub/style.css", Bytes: []byte("body{}")},
	}
	if err := Write(dir, arts); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "resume.typ"))
	if err != nil || string(got) != "hello" {
		t.Fatalf("resume.typ wrong: %q err=%v", got, err)
	}
	got2, err := os.ReadFile(filepath.Join(dir, "sub", "style.css"))
	if err != nil || string(got2) != "body{}" {
		t.Fatalf("nested file wrong: %q err=%v", got2, err)
	}
}
