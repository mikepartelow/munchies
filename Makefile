all: lint
	@go build -o . ./cmd/...

lint:
	golangci-lint run

clean:
	@rm -f ./munchies
