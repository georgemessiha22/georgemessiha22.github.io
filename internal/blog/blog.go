// Package blog loads blog posts authored as Markdown under a blog directory and
// renders them to HTML. Each post lives in its own folder as
// blog/<slug>/index.md, with any sibling images referenced relatively.
package blog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"

	"gopkg.in/yaml.v3"
)

// Asset is a non-Markdown file that ships alongside a post (e.g. an image).
type Asset struct {
	// Name is the file name relative to the post folder (e.g. "cover.png").
	Name string
	// Bytes is the file content.
	Bytes []byte
}

// Post is a single rendered blog article.
type Post struct {
	// Slug is the URL-safe folder name (e.g. "hello-world").
	Slug string
	// Title is the display title (front-matter "title", else the folder name).
	Title string
	// Date is the optional display date (front-matter "date", raw string).
	Date string
	// Summary is an optional short description (front-matter "summary").
	Summary string
	// HTML is the rendered article body (goldmark output), safe to embed.
	HTML string
	// Assets are sibling files (images, etc.) to copy next to the post.
	Assets []Asset
}

type frontMatter struct {
	Title   string `yaml:"title"`
	Date    string `yaml:"date"`
	Summary string `yaml:"summary"`
}

var md = goldmark.New(goldmark.WithExtensions(extension.GFM))

// LoadPosts reads every blog/<slug>/index.md under dir, renders the Markdown to
// HTML, and returns the posts sorted by date descending (newest first). If dir
// does not exist, it returns an empty slice with no error.
func LoadPosts(dir string) ([]Post, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read blog dir: %w", err)
	}

	var posts []Post
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		postDir := filepath.Join(dir, e.Name())
		indexPath := filepath.Join(postDir, "index.md")
		raw, err := os.ReadFile(indexPath)
		if err != nil {
			// A folder without index.md is not a post; skip it.
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("read %s: %w", indexPath, err)
		}

		fm, body := splitFrontMatter(raw)
		var meta frontMatter
		if len(fm) > 0 {
			if err := yaml.Unmarshal(fm, &meta); err != nil {
				return nil, fmt.Errorf("parse front-matter in %s: %w", indexPath, err)
			}
		}

		var htmlBuf bytes.Buffer
		if err := md.Convert(body, &htmlBuf); err != nil {
			return nil, fmt.Errorf("render markdown in %s: %w", indexPath, err)
		}

		title := meta.Title
		if title == "" {
			title = e.Name()
		}

		assets, err := loadAssets(postDir)
		if err != nil {
			return nil, err
		}

		posts = append(posts, Post{
			Slug:    e.Name(),
			Title:   title,
			Date:    meta.Date,
			Summary: meta.Summary,
			HTML:    htmlBuf.String(),
			Assets:  assets,
		})
	}

	sort.SliceStable(posts, func(i, j int) bool {
		// Newest first; posts without a date sort after dated ones.
		if posts[i].Date == posts[j].Date {
			return posts[i].Slug > posts[j].Slug
		}
		return posts[i].Date > posts[j].Date
	})
	return posts, nil
}

// loadAssets returns every non-index.md file in the post folder.
func loadAssets(postDir string) ([]Asset, error) {
	entries, err := os.ReadDir(postDir)
	if err != nil {
		return nil, fmt.Errorf("read post dir %s: %w", postDir, err)
	}
	var assets []Asset
	for _, e := range entries {
		if e.IsDir() || e.Name() == "index.md" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(postDir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("read asset %s: %w", e.Name(), err)
		}
		assets = append(assets, Asset{Name: e.Name(), Bytes: data})
	}
	return assets, nil
}

// splitFrontMatter separates a leading YAML front-matter block delimited by
// "---" lines from the Markdown body. If there is no front-matter, it returns
// nil front-matter and the original content as the body.
func splitFrontMatter(raw []byte) (fm, body []byte) {
	s := string(raw)
	if !strings.HasPrefix(s, "---\n") && !strings.HasPrefix(s, "---\r\n") {
		return nil, raw
	}
	// Find the closing delimiter after the opening one.
	rest := s[strings.IndexByte(s, '\n')+1:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return nil, raw
	}
	fmStr := rest[:idx]
	after := rest[idx+len("\n---"):]
	// Drop the remainder of the closing delimiter line.
	if nl := strings.IndexByte(after, '\n'); nl >= 0 {
		after = after[nl+1:]
	} else {
		after = ""
	}
	return []byte(fmStr), []byte(after)
}
