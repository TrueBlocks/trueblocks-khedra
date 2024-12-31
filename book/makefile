SRC_MD := $(shell find src -name '*.md')

book: $(SRC_MD)
	@echo building...
	mdbook build

serve:
	@make clean
	@mdbook serve --open --hostname 127.0.0.1

clean:
	@rm -fR book

deploy:
	@make clean && touch ${SRC_MD} && make book && cd book && scp -pr * "docs@khedra.trueblocks.io:/var/www/khedra.trueblocks.io/html/"

