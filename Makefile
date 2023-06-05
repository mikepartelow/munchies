all:
	@go build -o . ./cmd/...

clean:
	@rm -f ./munchies
