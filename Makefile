.PHONY: help build run test clean tidy build-all build-linux build-darwin build-windows

BINARY  := ttb
VERSION := 0.1.0
LDFLAGS := -s -w -X main.version=$(VERSION)

help: ## Show this help
	@echo "Tube Trend Buddy - make targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

tidy: ## go mod tidy
	go mod tidy

build: tidy ## Build for current OS
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) .

run: build ## Build + run --help
	./bin/$(BINARY) --help

test: ## go test ./...
	go test ./...

build-linux: tidy ## Cross-compile linux/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-linux-amd64 .

build-darwin: tidy ## Cross-compile darwin/arm64 (Apple Silicon)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-darwin-arm64 .

build-windows: tidy ## Cross-compile windows/amd64 -> .exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-windows-amd64.exe .

build-all: build-linux build-darwin build-windows ## Build all 3 platforms
	@ls -lh dist/

clean: ## Remove build artifacts
	rm -rf bin/ dist/
