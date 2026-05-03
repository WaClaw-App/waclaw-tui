#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
BIN_DIR="$ROOT_DIR/bin"

# Ensure bin directory exists
mkdir -p "$BIN_DIR"

# Build both binaries (skip if already built and Makefile handled it)
if [ ! -x "$BIN_DIR/waclaw-backend" ] || [ ! -x "$BIN_DIR/waclaw-tui" ]; then
        echo "Building backend..."
        go build -o "$BIN_DIR/waclaw-backend" "$ROOT_DIR/cmd/backend/"

        echo "Building TUI..."
        go build -o "$BIN_DIR/waclaw-tui" "$ROOT_DIR/cmd/tui/"
fi

# Set default REST address (can be overridden)
WA_REST_ADDR="${WA_REST_ADDR:-:8080}"
export WA_REST_ADDR

# Create a temporary directory for the named pipes.
PIPE_DIR="$(mktemp -d)"
trap 'rm -rf "$PIPE_DIR"' EXIT

# Create two named pipes (FIFOs) for bidirectional stdio communication.
#   backend2tui: backend stdout → TUI stdin
#   tui2backend: TUI stdout → backend stdin
mkfifo "$PIPE_DIR/backend2tui"
mkfifo "$PIPE_DIR/tui2backend"

echo "Starting demo (REST API on $WA_REST_ADDR)..."
echo "  Backend → TUI: JSON-RPC 2.0 over stdio (bidirectional via named pipes)"
echo "  REST API: http://localhost${WA_REST_ADDR}/api/v1/state"
echo ""

# Start the backend reading from the TUI→backend pipe and writing to the
# backend→TUI pipe. The REST server also starts on $WA_REST_ADDR.
"$BIN_DIR/waclaw-backend" < "$PIPE_DIR/tui2backend" > "$PIPE_DIR/backend2tui" &
BACKEND_PID=$!

# Start the TUI reading from the backend→TUI pipe and writing to the
# TUI→backend pipe. WA_DEMO=1 tells the TUI to use /dev/tty for the
# terminal (bubbletea) while keeping stdin/stdout for RPC over pipes.
WA_DEMO=1 "$BIN_DIR/waclaw-tui" < "$PIPE_DIR/backend2tui" > "$PIPE_DIR/tui2backend" &
TUI_PID=$!

# Wait for either process to exit. When one dies, kill the other.
wait -n "$BACKEND_PID" "$TUI_PID" 2>/dev/null || true

# Clean up: kill any remaining process.
kill "$BACKEND_PID" 2>/dev/null || true
kill "$TUI_PID" 2>/dev/null || true
wait "$BACKEND_PID" 2>/dev/null || true
wait "$TUI_PID" 2>/dev/null || true

echo ""
echo "Demo ended."
