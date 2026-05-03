.PHONY: all tui backend demo test lint vet fmt clean

# Go parameters
GOROOT  ?= $(shell go env GOROOT)
GOPATH  ?= $(shell go env GOPATH)
GO      ?= go
MAINMOD := github.com/WaClaw-App/waclaw

# Build output
BUILDDIR := ./bin

# Linter
LINTER := golangci-lint

all: tui backend

## tui: Build the TUI binary
tui:
	$(GO) build -o $(BUILDDIR)/waclaw-tui ./cmd/tui/

## backend: Build the backend binary
backend:
	$(GO) build -o $(BUILDDIR)/waclaw-backend ./cmd/backend/

## demo: Run TUI with demo backend
demo: backend tui
	@bash scripts/demo.sh

## test: Run all tests
test:
	$(GO) test -race -count=1 ./...

## lint: Run golangci-lint
lint:
	$(LINTER) run ./...

## vet: Run go vet
vet:
	$(GO) vet ./...

## fmt: Format all Go files
fmt:
	gofmt -w -s .

## tidy: Run go mod tidy
tidy:
	$(GO) mod tidy

## clean: Remove build artifacts
clean:
	rm -rf $(BUILDDIR)

## proto: Generate OpenAPI spec from protocol types
proto:
	$(GO) generate ./pkg/api/...
