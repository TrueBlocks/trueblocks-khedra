SRC_GO := $(shell find . -name '*.go')

#-------------------------------------------------
bin=../bin

#-------------------------------------------------
exec=khedra
dest=$(bin)/$(exec)

#-------------------------------------------------
all: $(SRC_GO)
	@make app

every:
	@cd ../build ; make ; cd -
	@make app

app: $(SRC_GO)
	@mkdir -p $(bin)
	@go build -o $(dest) *.go

update:
	@go get "github.com/TrueBlocks/trueblocks-sdk/v4@latest"
	@go get github.com/TrueBlocks/trueblocks-core/src/apps/chifra@latest

install:
	@make build
	@mv khedra ~/go/bin

test: $(SRC_GO)
	@go test ./...

#-------------------------------------------------
clean:
	-@$(RM) -f $(dest)
