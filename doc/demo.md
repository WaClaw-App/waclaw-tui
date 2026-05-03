# `waclaw demo` — Showcase Mode

## What

`waclaw demo` launches the TUI with a **scenario demo backend** instead of the real backend. Same UI, same animations, same 20-screen flow — but no real WhatsApp connection, no scraping, no messaging, no license required.

## Command

```bash
waclaw demo
```

## Architecture

```
waclaw demo = waclaw-tui (real frontend) + demo backend (mock engine)
```

The demo backend is the **same scenario/state engine** used in production. It drives screen navigation, injects mock data, and fires timed transitions. Communication uses the standard **JSON-RPC 2.0 over stdio** protocol.

## What It Does NOT Do

- Connect to WhatsApp
- Scrape Google Maps
- Send real messages
- Require a valid license
- Write to the real database

## Purpose

| Use Case | Description |
|----------|-------------|
| Marketing videos | Record the full TUI flow for promo content |
| Screenshots | Capture terminal shots for landing pages or docs |
| UI evaluation | Experience all 20 screens and 119 states end-to-end |
| Prospect demo | Try the product visually without setup or commitment |

## Screen Flow (Demo Order)

```
BOOT → LICENSE → GUARDRAIL (Validation) → LOGIN → NICHE SELECT → SCRAPE
  → LEAD REVIEW → SEND → MONITOR → RESPONSE → LEADS DB → TEMPLATES
  → WORKERS → ANTI-BAN → SETTINGS → COMPOSE → HISTORY
  → FOLLOW-UP → NICHE EXPLORER → UPDATE
```

The first-time demo flow includes the **License gate** and **Validation gate** (Guardrail screen) between Boot and Login — matching the documented screen flow in `doc/18-screen-flow.md`. On subsequent cycles, the demo loops back to Boot with the returning-user variant.

The demo backend orchestrates this flow on a **timeline** — screens auto-advance with timed transitions and injected mock data. A full cycle takes approximately **64 seconds**. After the Update screen, the demo loops back to the Boot screen with the `boot_returning` state, and the cycle repeats indefinitely. Users can also navigate freely via keyboard at any time.

## Key Interactions

| Key | Action |
|-----|--------|
| `Ctrl+K` | Command palette — search actions, navigate, execute |
| `` ` `` | Toggle nerd stats overlay (RAM, CPU, goroutines, DB) |
| `?` | Keyboard shortcut cheat sheet |
| `q` | Quit (with session summary) |
| `s` | Skip / pause |
| `↵` | Confirm / approve |
| `h` | History screen |
| `v` | Force config validation |

## Demo Backend Behavior

- **Mock data**: Realistic lead names, businesses, ratings, WhatsApp numbers
- **Timed transitions**: Screen-to-screen navigation on a script
- **Notification simulation**: Toast popups for scrape results, incoming responses, deals
- **State cycling**: Workers go through idle → scraping → qualifying → sending → idle
- **Closing triggers**: Simulated response classification (deal, hot_lead, stop)

## Tech Stack (Shared with Production)

- **Frontend**: Go 1.24+ / bubbletea / lipgloss / bubbles / glamour / huh (Charm.sh)
- **Protocol**: JSON-RPC 2.0 over stdio (`pkg/protocol`, `pkg/transport`)
- **Entry point**: `make demo` builds both binaries and runs `scripts/demo.sh`, which pipes the demo backend binary into the TUI binary via named pipes (FIFOs)
- **REST API**: The demo backend also starts a REST server on `:8080` (configurable via `WA_REST_ADDR`) for the future web frontend

## How to Run

```bash
# Prerequisites: Go 1.24+

# Clone and build
git clone https://github.com/WaClaw-App/waclaw-tui.git
cd waclaw-tui
make demo
```

The `make demo` target:
1. Builds `waclaw-backend` and `waclaw-tui` binaries into `./bin/`
2. Runs `scripts/demo.sh`, which sets up bidirectional named pipes and launches both processes
3. The TUI opens in your terminal with the demo timeline running automatically
