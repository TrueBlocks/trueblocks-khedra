#-------------------------------------------------
bin=../bin

#-------------------------------------------------
exec=khedra
dest=$(bin)/$(exec)

#-------------------------------------------------
all:
	@make app

every:
	@cd ../build ; make ; cd -
	@make app

app: **/*.go
	@echo Building khedra...
	@mkdir -p $(bin)
	@go build -o $(dest) *.go

update:
	@go get "github.com/TrueBlocks/trueblocks-sdk/v4@latest"
	@go get github.com/TrueBlocks/trueblocks-core/src/apps/chifra@latest

install:
	@make build
	@mv khedra ~/go/bin

test:
	@go test ./...

#-------------------------------------------------
clean:
	-@$(RM) -f $(dest)
