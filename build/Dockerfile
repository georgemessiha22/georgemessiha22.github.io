# Full TeX Live 2025 (frozen). This image already bundles everything the
# generated resume.tex needs (moderncv, fontawesome5, latexmk, pdflatex), so no
# tlmgr install is required — which also avoids the "local TL older than remote
# mirror" failure that hits pinned historic images.
FROM registry.gitlab.com/islandoftex/images/texlive:TL2025-historic

LABEL maintainer="George Messiha (gmessiha)"

# python is used by texliveonfly / helper tooling if invoked; harmless otherwise.
RUN apt-get update && apt-get install -y --no-install-recommends python3 \
 && rm -rf /var/lib/apt/lists/* \
 && update-alternatives --install /usr/bin/python python /usr/bin/python3 1

WORKDIR /data

COPY build/entrypoint.sh /entrypoint
RUN chmod +x /entrypoint

ENTRYPOINT ["/entrypoint"]
