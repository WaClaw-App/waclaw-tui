#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

VERSION ?= "$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')"
LDFLAGS="-s -w -X main.version=$VERSION"

echo "Building release binaries (version: $VERSION)..."

mkdir -p "$ROOT_DIR/bin"

# TUI binary
echo "  → waclaw-tui"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$ROOT_DIR/bin/waclaw-tui-linux-amd64" "$ROOT_DIR/cmd/tui/"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$ROOT_DIR/bin/waclaw-tui-darwin-amd64" "$ROOT_DIR/cmd/tui/"
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o "$ROOT_DIR/bin/waclaw-tui-darwin-arm64" "$ROOT_DIR/cmd/tui/"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$ROOT_DIR/bin/waclaw-tui-windows-amd64.exe" "$ROOT_DIR/cmd/tui/"

# Backend binary
echo "  → waclaw-backend"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$ROOT_DIR/bin/waclaw-backend-linux-amd64" "$ROOT_DIR/cmd/backend/"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$ROOT_DIR/bin/waclaw-backend-darwin-amd64" "$ROOT_DIR/cmd/backend/"
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o "$ROOT_DIR/bin/waclaw-backend-darwin-arm64" "$ROOT_DIR/cmd/backend/"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$ROOT_DIR/bin/waclaw-backend-windows-amd64.exe" "$ROOT_DIR/cmd/backend/"

echo "Done. Binaries in $ROOT_DIR/bin/"
