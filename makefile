build:
	@rm -f khedra && go build -o khedra main.go

test:
	@go test ./...
