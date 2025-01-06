SRC_MD := $(shell find src -name '*.md')
SRC_PNG := $(shell find src -name '*.png')

book: $(SRC_MD) $(SRC_PNG) makefile
	@echo building...
	@mdbook build

serve:
	@make clean
	@make book
	@mdbook serve --open --hostname 127.0.0.1

clean:
	@rm -fR book

deploy:
	@make clean && touch ${SRC_MD} && make book && cd book && scp -pr * "docs@khedra.trueblocks.io:/var/www/khedra.trueblocks.io/html/"

