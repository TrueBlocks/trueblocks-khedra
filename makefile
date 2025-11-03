.PHONY: app

SRC_GO := $(shell find . -name '*.go')

#-------------------------------------------------
bin=../bin

#-------------------------------------------------
exec=khedra
dest=$(bin)/$(exec)

#-------------------------------------------------
all: $(SRC_GO)
	@make app

app:
	@cd ../build ; make khedra ; cd -
	@mkdir -p $(bin)
	@go build -o $(dest) *.go

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
