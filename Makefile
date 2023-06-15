all: lint
	@go build -o . ./cmd/...

lint:
	@golangci-lint run

test: lint
	@go test -cover ./...

clean:
	@rm -f ./munchies
