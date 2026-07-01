package blog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writePost(t *testing.T, root, slug, content string, assets map[string]string) {
	t.Helper()
	dir := filepath.Join(root, slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	for name, data := range assets {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLoadPosts_MissingDir(t *testing.T) {
	posts, err := LoadPosts(filepath.Join(t.TempDir(), "does-not-exist"))
	if err != nil {
		t.Fatalf("missing dir should not error: %v", err)
	}
	if len(posts) != 0 {
		t.Fatalf("expected 0 posts, got %d", len(posts))
	}
}

func TestLoadPosts_FrontMatterAndRender(t *testing.T) {
	root := t.TempDir()
	writePost(t, root, "hello-world", `---
title: Hello World
date: "2026-06-30"
summary: A first post.
---
# Heading

Some **bold** text and an image:

![cover](cover.png)
`, map[string]string{"cover.png": "PNGDATA"})

	posts, err := LoadPosts(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 1 {
		t.Fatalf("want 1 post, got %d", len(posts))
	}
	p := posts[0]
	if p.Slug != "hello-world" || p.Title != "Hello World" || p.Date != "2026-06-30" {
		t.Fatalf("front-matter not parsed: %+v", p)
	}
	if p.Summary != "A first post." {
		t.Fatalf("summary not parsed: %q", p.Summary)
	}
	if !strings.Contains(p.HTML, "<h1") || !strings.Contains(p.HTML, "<strong>bold</strong>") {
		t.Fatalf("markdown not rendered: %s", p.HTML)
	}
	if !strings.Contains(p.HTML, `<img src="cover.png"`) {
		t.Fatalf("image not rendered: %s", p.HTML)
	}
	if len(p.Assets) != 1 || p.Assets[0].Name != "cover.png" || string(p.Assets[0].Bytes) != "PNGDATA" {
		t.Fatalf("assets not collected: %+v", p.Assets)
	}
}

func TestLoadPosts_NoFrontMatter_TitleFromSlug(t *testing.T) {
	root := t.TempDir()
	writePost(t, root, "untitled", "Just body text.\n", nil)
	posts, err := LoadPosts(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 1 || posts[0].Title != "untitled" {
		t.Fatalf("expected title from slug, got %+v", posts)
	}
	if !strings.Contains(posts[0].HTML, "Just body text.") {
		t.Fatalf("body not rendered: %s", posts[0].HTML)
	}
}

func TestLoadPosts_SortedNewestFirst(t *testing.T) {
	root := t.TempDir()
	writePost(t, root, "older", "---\ntitle: Older\ndate: \"2025-01-01\"\n---\nold\n", nil)
	writePost(t, root, "newer", "---\ntitle: Newer\ndate: \"2026-01-01\"\n---\nnew\n", nil)
	posts, err := LoadPosts(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 2 || posts[0].Slug != "newer" || posts[1].Slug != "older" {
		t.Fatalf("expected newest first, got %+v", []string{posts[0].Slug, posts[1].Slug})
	}
}

func TestLoadPosts_SkipsFoldersWithoutIndex(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "notapost"), 0o755); err != nil {
		t.Fatal(err)
	}
	writePost(t, root, "real", "---\ntitle: Real\n---\nx\n", nil)
	posts, err := LoadPosts(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 1 || posts[0].Slug != "real" {
		t.Fatalf("should skip folder without index.md, got %+v", posts)
	}
}
