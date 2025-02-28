.PHONY: build test lint clean example file-listing

# Build the library
build:
	go build -v ./...

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/

# Run the simple example
example:
	go run examples/simple/main.go

# Run the file listing example
file-listing:
	go run examples/file_listing/main.go

# Install dependencies
deps:
	go get -v ./...
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Help command
help:
	@echo "Available commands:"
	@echo "  make build        - Build the library"
	@echo "  make test         - Run tests"
	@echo "  make lint         - Run linter"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make example      - Run the simple example"
	@echo "  make file-listing - Run the file listing example"
	@echo "  make deps         - Install dependencies" 