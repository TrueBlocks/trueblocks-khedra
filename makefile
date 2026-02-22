.PHONY: app

SRC_GO := $(shell find . -name '*.go')

#-------------------------------------------------
bin=bin

#-------------------------------------------------
exec=khedra
dest=$(shell which khedra)

#-------------------------------------------------
all: $(SRC_GO)
	@make app

app:
	@mkdir -p $(bin)
	@rm -f $(dest)
	@go build -o $(dest) .
	@ls -l $(dest)

update:
	@go get "github.com/TrueBlocks/trueblocks-sdk/v6@latest"
	@go get github.com/TrueBlocks/trueblocks-chifra/v6@latest
	@go get -u ./...

install:
	@make app
	@mv $(dest) ~/go/bin

test: $(SRC_GO)
	@go test ./...

#-------------------------------------------------
clean:
	-@$(RM) -f $(dest)
