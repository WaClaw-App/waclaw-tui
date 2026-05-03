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
| UI evaluation | Experience all 20 screens and 110 states end-to-end |
| Prospect demo | Try the product visually without setup or commitment |

## Screen Flow (Demo Order)

```
BOOT → LOGIN → NICHE SELECT → SCRAPE → LEAD REVIEW → SEND
  → MONITOR → RESPONSE → LEADS DB → TEMPLATES → WORKERS
  → ANTI-BAN → SETTINGS → GUARDRAIL → COMPOSE → HISTORY
  → FOLLOW-UP → LICENSE → NICHE EXPLORER → UPDATE
```

The demo backend orchestrates this flow on a **timeline** — screens auto-advance with timed transitions and injected mock data. Users can also navigate freely via keyboard.

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

- **Frontend**: Go 1.22+ / bubbletea / lipgloss / bubbles / glamour / huh (Charm.sh)
- **Protocol**: JSON-RPC 2.0 over stdio (`pkg/protocol`, `pkg/transport`)
- **Entry point**: `scripts/demo.sh` pipes demo backend binary into TUI binary
