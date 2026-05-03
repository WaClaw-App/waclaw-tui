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

# Pre-open the write ends of the named pipes to prevent deadlock.
#
# Named pipes (FIFOs) block on open() until both the read and write ends
# are connected. If both processes open their read end first (which is
# what happens with `< pipe > pipe` shell redirection order), deadlock
# occurs: each process blocks waiting for the other to open the write end.
#
# Strategy: Use read-write open (<> ) which does NOT block, then
# immediately close the unwanted read side by redirecting it away.
# This gives us a write-only FD without creating a competing reader.
#
#   1. exec 3<> pipe   — opens read-write (non-blocking)
#   2. exec 3>&3       — no-op to verify
#   Then use >&3 as write target for backend process.
#
# The extra read-side on FD 3 and FD 4 is harmless because we NEVER read
# from them — the kernel only delivers data to readers that actually call
# read(). Since these FDs are used purely as write targets, the kernel
# will deliver all pipe data to the actual reading processes (TUI/backend).

# Open write-end FDs in read-write mode (non-blocking).
exec 3<>"$PIPE_DIR/backend2tui"
exec 4<>"$PIPE_DIR/tui2backend"

# Start the backend reading from the TUI→backend pipe and writing to the
# backend→TUI pipe. The REST server also starts on $WA_REST_ADDR.
# Redirect stderr to /dev/null to avoid mixing backend log messages with
# the TUI's terminal output.
WA_DEMO=1 "$BIN_DIR/waclaw-backend" < "$PIPE_DIR/tui2backend" >&3 2>/dev/null &
BACKEND_PID=$!

# Start the TUI reading from the backend→TUI pipe and writing to the
# TUI→backend pipe. WA_DEMO=1 tells the TUI to use /dev/tty for the
# terminal (bubbletea) while keeping stdin/stdout for RPC over pipes.
WA_DEMO=1 "$BIN_DIR/waclaw-tui" < "$PIPE_DIR/backend2tui" >&4 &
TUI_PID=$!

# Close the script's copy of the write FDs so that when the processes
# exit, the pipes properly signal EOF to the other side.
exec 3>&- 4>&-

# Wait for either process to exit. When one dies, kill the other.
# Use a polling loop instead of `wait -n` for bash 3.2+ compatibility
# (macOS ships bash 3.2; `wait -n PID...` requires bash 5.1+).
while kill -0 "$BACKEND_PID" 2>/dev/null && kill -0 "$TUI_PID" 2>/dev/null; do
        sleep 0.2
done

# Clean up: kill any remaining process.
kill "$BACKEND_PID" 2>/dev/null || true
kill "$TUI_PID" 2>/dev/null || true
wait "$BACKEND_PID" 2>/dev/null || true
wait "$TUI_PID" 2>/dev/null || true

echo ""
echo "Demo ended."
