VERSION_TAG=$(shell date +'%y.%m.%d')
PWD=$(shell pwd)

TYPST_FONTS=typst/fonts

.PHONY: build bash detailed resume clean all typst typst-resume typst-detailed

build:
	docker build -f build/Dockerfile -t tex:latest -t tex:$(VERSION_TAG) .

bash:
	docker run --rm -it -v $(PWD):/data tex:latest bash

detailed:
	docker run --rm -v $(PWD):/data tex:latest pdflatex -output-directory dist George_Messiha_detailed_resume.tex

resume:
	docker run --rm -v $(PWD):/data tex:latest pdflatex -output-directory dist George_Messiha_Resume.tex

clean:
	rm -f dist/*.aux dist/*.out dist/*.log

# ---------------------------------------------------------------------------
# Typst (parallel resume version) — uses the local `typst` binary, no Docker.
# Outputs use a _v2 suffix so they don't collide with the LaTeX PDFs.
# ---------------------------------------------------------------------------
typst-resume:
	mkdir -p dist
	typst compile --font-path $(TYPST_FONTS) typst/resume.typ dist/George_Messiha_Resume_v2.pdf

typst-detailed:
	mkdir -p dist
	typst compile --font-path $(TYPST_FONTS) typst/detailed_resume.typ dist/George_Messiha_detailed_resume_v2.pdf

typst: typst-resume typst-detailed

all: build detailed resume typst clean
