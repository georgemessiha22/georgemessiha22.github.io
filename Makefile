VERSION_TAG=$(shell date +'%y.%m.%d')
PWD=$(shell pwd)
GEN=build/gen
TYPST_FONTS=typst/fonts

.PHONY: build bash gen resume detailed typst typst-resume typst-detailed site clean all

build:
	docker build -f build/Dockerfile -t tex:latest -t tex:$(VERSION_TAG) .

bash:
	docker run --rm -it -v $(PWD):/data tex:latest bash

# Generate all source files (typst + tex, both variants) and the html site.
gen:
	go run ./cmd/resume all --input resume.yaml
	mkdir -p $(GEN)/assets
	cp pictures/61673.jpg $(GEN)/assets/profile.jpg

# LaTeX PDFs (compiled in Docker) from generated .tex.
resume: gen
	mkdir -p dist
	docker run --rm -v $(PWD):/data tex:latest pdflatex -output-directory dist $(GEN)/resume.tex
	mv dist/resume.pdf dist/George_Messiha_Resume.pdf

detailed: gen
	mkdir -p dist
	docker run --rm -v $(PWD):/data tex:latest pdflatex -output-directory dist $(GEN)/resume_detailed.tex
	mv dist/resume_detailed.pdf dist/George_Messiha_detailed_resume.pdf

# Typst PDFs from generated .typ.
typst-resume: gen
	mkdir -p dist
	typst compile --font-path $(TYPST_FONTS) $(GEN)/resume.typ dist/George_Messiha_Resume_v2.pdf

typst-detailed: gen
	mkdir -p dist
	typst compile --font-path $(TYPST_FONTS) $(GEN)/resume_detailed.typ dist/George_Messiha_detailed_resume_v2.pdf

typst: typst-resume typst-detailed

# HTML mini-site (detailed) for GitHub Pages.
site: gen

clean:
	rm -f dist/*.aux dist/*.out dist/*.log

all: build resume detailed typst site clean
