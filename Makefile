.PHONY: build test lint generate update-icons clean install bench

BINARY     := modern-ls
VERSION    ?= $(shell git describe --tags --dirty --always 2>/dev/null || echo "dev")
LDFLAGS    := -ldflags "-X github.com/the-mayankjha/modern-ls/internal/cli.Version=$(VERSION) -s -w"
BUILD_DIR  := ./dist

## build: Compile the binary for the current platform
build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/modern-ls/...

## install: Install the binary to $GOPATH/bin
install:
	go install $(LDFLAGS) ./cmd/modern-ls/...

## generate: Run the icon generator (requires internet on first run to clone vendor)
generate:
	go generate ./cmd/modern-ls/...

## update-icons: Fetch latest nvim-web-devicons and regenerate icons
update-icons: scripts/update-icons.sh
	bash scripts/update-icons.sh

## test: Run all tests with race detector
test:
	go test -race -count=1 ./...

## bench: Run benchmarks
bench:
	go test -bench=. -benchmem ./internal/icons/... ./internal/renderer/...

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## vet: Run go vet
vet:
	go vet ./...

## fmt: Format all Go source
fmt:
	gofmt -w .

## clean: Remove build artifacts
clean:
	rm -f $(BINARY)
	rm -rf $(BUILD_DIR)

## help: Print this help message
help:
	@echo "modern-ls — Available make targets:"
	@grep -E '^## ' Makefile | sed 's/## /  /'
