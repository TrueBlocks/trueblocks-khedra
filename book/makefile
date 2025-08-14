SRC_MD := $(shell find src -name '*.md')
SRC_PNG := $(shell find src -name '*.png')

book: $(SRC_MD) $(SRC_PNG) makefile
	@echo building...
	@echo "injecting build info script..."
	@echo "// generated – do not commit\nwindow.MDBOOK_BUILD_DATE='$(shell git log -1 --format=%cI || date -u +%Y-%m-%dT%H:%M:%SZ)'; window.MDBOOK_GIT_COMMIT='$(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)';" > theme/build-info.js
	@mdbook build

serve:
	@make clean
	@make book
	@mdbook serve --open --hostname 127.0.0.1

clean:
	@rm -fR book

deploy:
	@make clean && touch ${SRC_MD} && make book && cd book && scp -pr * "docs@khedra.trueblocks.io:/var/www/khedra.trueblocks.io/html/"

