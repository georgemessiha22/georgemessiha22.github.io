package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/georgemessiha22/georgemessiha22/internal/config"
	"github.com/georgemessiha22/georgemessiha22/internal/html"
	"github.com/georgemessiha22/georgemessiha22/internal/latex"
	"github.com/georgemessiha22/georgemessiha22/internal/markdown"
	"github.com/georgemessiha22/georgemessiha22/internal/model"
	"github.com/georgemessiha22/georgemessiha22/internal/render"
	"github.com/georgemessiha22/georgemessiha22/internal/typst"
	"github.com/georgemessiha22/georgemessiha22/internal/writer"
)

func usage() {
	fmt.Fprint(os.Stderr, `resume - generate resume sources from resume.yaml

Usage:
  resume <command> [flags]

Commands:
  typst   Generate Typst (.typ) source
  tex     Generate LaTeX (.tex) source
  md      Generate Markdown (.md) source
  html    Generate the HTML site (résumé home + blog)
  all     Generate typst+tex+md (both variants) and the html site

Flags:
  --input   path to YAML (default resume.yaml)
  --variant summary|detailed (default summary; html always uses summary)
  --out     output directory (default build/gen; site default site)
  --blog    blog posts directory (default blog) for the html site
`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	input := fs.String("input", "resume.yaml", "path to YAML input")
	variant := fs.String("variant", "summary", "summary|detailed")
	out := fs.String("out", "", "output directory")
	blog := fs.String("blog", "blog", "directory of blog/<slug>/index.md posts (html/all)")
	_ = fs.Parse(os.Args[2:])

	if err := run(cmd, *input, *variant, *out, *blog); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(cmd, input, variant, out, blogDir string) error {
	r, err := config.Load(input)
	if err != nil {
		return err
	}

	defGen := "build/gen"
	defSite := "site"

	switch cmd {
	case "typst":
		return generate(typst.Renderer{}, r, parseVariant(variant), pick(out, defGen))
	case "tex":
		return generate(latex.Renderer{}, r, parseVariant(variant), pick(out, defGen))
	case "md":
		return generate(markdown.Renderer{}, r, parseVariant(variant), pick(out, defGen))
	case "html":
		return generate(html.Renderer{BlogDir: blogDir}, r, model.Summary, pick(out, defSite))
	case "all":
		if err := generate(typst.Renderer{}, r, model.Summary, defGen); err != nil {
			return err
		}
		if err := generate(typst.Renderer{}, r, model.Detailed, defGen); err != nil {
			return err
		}
		if err := generate(latex.Renderer{}, r, model.Summary, defGen); err != nil {
			return err
		}
		if err := generate(latex.Renderer{}, r, model.Detailed, defGen); err != nil {
			return err
		}
		if err := generate(markdown.Renderer{}, r, model.Summary, defGen); err != nil {
			return err
		}
		if err := generate(markdown.Renderer{}, r, model.Detailed, defGen); err != nil {
			return err
		}
		return generate(html.Renderer{BlogDir: blogDir}, r, model.Summary, defSite)
	default:
		usage()
		return fmt.Errorf("unknown command %q", cmd)
	}
}

func generate(rd render.Renderer, r model.Resume, v model.Variant, dir string) error {
	arts, err := rd.Render(r, v)
	if err != nil {
		return err
	}
	return writer.Write(dir, arts)
}

func parseVariant(s string) model.Variant {
	if s == string(model.Detailed) {
		return model.Detailed
	}
	return model.Summary
}

func pick(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
