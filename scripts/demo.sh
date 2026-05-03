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

# Open the write ends of the named pipes first to prevent deadlock.
#
# Named pipes (FIFOs) block on open() until both the read and write ends
# are connected. If both processes open their read end first (which is
# what happens with the `< pipe > pipe` shell redirection order), deadlock
# occurs: each process blocks waiting for the other to open the write end.
#
# By opening file descriptors 3 and 4 for writing before launching the
# processes, the subsequent read opens in the child processes unblock
# immediately. The file descriptors are inherited and kept open for the
# lifetime of each child process.
exec 3>"$PIPE_DIR/backend2tui"   # write end: backend stdout → TUI stdin
exec 4>"$PIPE_DIR/tui2backend"   # write end: TUI stdout → backend stdin

# Start the backend reading from the TUI→backend pipe and writing to the
# backend→TUI pipe. The REST server also starts on $WA_REST_ADDR.
"$BIN_DIR/waclaw-backend" < "$PIPE_DIR/tui2backend" >&3 &
BACKEND_PID=$!

# Start the TUI reading from the backend→TUI pipe and writing to the
# TUI→backend pipe. WA_DEMO=1 tells the TUI to use /dev/tty for the
# terminal (bubbletea) while keeping stdin/stdout for RPC over pipes.
WA_DEMO=1 "$BIN_DIR/waclaw-tui" < "$PIPE_DIR/backend2tui" >&4 &
TUI_PID=$!

# Wait for either process to exit. When one dies, kill the other.
# Use a polling loop instead of `wait -n` for bash 3.2+ compatibility
# (macOS ships bash 3.2; `wait -n PID...` requires bash 5.1+).
while kill -0 "$BACKEND_PID" 2>/dev/null && kill -0 "$TUI_PID" 2>/dev/null; do
        sleep 0.2
done

# Clean up: kill any remaining process and close pipe file descriptors.
kill "$BACKEND_PID" 2>/dev/null || true
kill "$TUI_PID" 2>/dev/null || true
wait "$BACKEND_PID" 2>/dev/null || true
wait "$TUI_PID" 2>/dev/null || true
exec 3>&- 4>&-

echo ""
echo "Demo ended."
