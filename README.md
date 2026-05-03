<div align="center">

```
                                                  
  ▄▄▄                 ▄   ▄▄▄▄ ▄▄                 
 █▀██  ██  ██▀▀       ▀██████▀  ██                
   ██  ██  ██           ██      ██                
   ██  ██  ██ ▄▀▀█▄     ██      ██ ▄▀▀█▄▀█▄ █▄ ██▀
   ██▄ ██▄ ██ ▄█▀██     ██      ██ ▄█▀██ ██▄██▄██ 
   ▀████▀███▀▄▀█▄██     ▀█████ ▄██▄▀█▄██  ▀██▀██▀ 
                                                  
```

### _Vertical-borderless. Micro-interactive. File-based._

**Lu cuma nonton. WaClaw yang kerja.**

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Charm](https://img.shields.io/badge/Charm.sh-Ecosystem-FF69B4?style=flat-square)](https://charm.sh/)
[![Screens](https://img.shields.io/badge/Screens-20-blueviolet?style=flat-square)]()
[![States](https://img.shields.io/badge/States-110-orange?style=flat-square)]()
[![License](https://img.shields.io/badge/License-Proprietary-red?style=flat-square)]()

</div>

---

## What is WaClaw?

WaClaw is a **multi-niche WhatsApp lead generation army** that runs in your terminal. It scrapes Google Maps for business leads, validates WhatsApp numbers, sends ice-breaker messages with template rotation, handles responses with closing-trigger detection, and follows up automatically — all from a single TUI dashboard.

One command. Multiple niches. Parallel workers. Auto-pilot by default.

```
Lu = jenderal. WaClaw = pasukan.
Lu tentuin strategi, mereka eksekusi.
Lu nggak perlu micromanage — tiap worker otonom.
```

> **This repo** contains the Go project structure and TUI frontend. The backend binary is closed-source and communicates with the TUI via **JSON-RPC over stdio**. The backend acts as a state scenario controller — it drives the TUI's screen states, transitions, and data. The `internal/` directory maps every documented screen, state, and system into its corresponding Go package.

---

## `waclaw demo`

Running `waclaw demo` launches the TUI with the **scenario demo backend** instead of the real working backend binary.

The demo backend is the same scenario engine that drives the TUI in production — it orchestrates screen navigation, injects mock data, and fires timed state transitions — but it does not connect to WhatsApp, scrape Google Maps, or send any real messages. Its sole purpose is to drive the TUI through the full screen flow with realistic-looking data so you can record marketing videos, take screenshots, or just see how the UI behaves end-to-end without needing a live WA session or a valid license.

In other words: `waclaw demo` = waclaw-tui + scenario demo backend binary. No real features. Pure showcase.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        WaClaw System                            │
│                                                                 │
│  ┌──────────────────────┐     JSON-RPC      ┌────────────────┐ │
│  │                      │     over stdio     │                │ │
│  │   Backend Binary     │◄──────────────────►│  TUI Frontend  │ │
│  │   (closed source)    │                    │  (this repo)   │ │
│  │                      │   ┌────────────┐   │                │ │
│  │  • State scenarios   │   │  Request:  │   │  • bubbletea   │ │
│  │  • Screen transitions│──►│  navigate   │──►│  • lipgloss    │ │
│  │  • Mock data feed    │   │  update     │   │  • bubbles     │ │
│  │  • Event simulation  │   │  notify     │   │  • glamour     │ │
│  │  • Timeline control  │   │  validate   │   │  • huh         │ │
│  │                      │   └────────────┘   │                │ │
│  │                      │   ┌────────────┐   │  20 Screens    │ │
│  │                      │◄──│  Response:  │◄──│  110 States    │ │
│  │                      │   │  key_press  │   │  17 Notifs     │ │
│  │                      │   │  action     │   │  2 Overlays    │ │
│  └──────────────────────┘   └────────────┘   └────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

**Communication Protocol:**

| Layer | Protocol | Direction | Package |
|-------|----------|-----------|---------|
| Transport (TUI) | JSON-RPC 2.0 over stdio | Bidirectional | `pkg/transport/stdio.go` |
| Transport (Web) | REST API + Swagger (OpenAPI 3.0) | Bidirectional | `pkg/transport/http.go` |
| Backend → TUI | `navigate`, `update`, `notify`, `validate` | Push | `internal/backend/rpc` → `internal/tui/bus` |
| TUI → Backend | `key_press`, `action`, `request` | Push | `internal/tui/rpc` → `internal/backend/rpc` |
| Backend → Web | Same RPC methods as REST endpoints | HTTP | `internal/backend/rest` |
| Swagger spec | Auto-generated from `pkg/protocol` types | — | `pkg/api/openapi.yaml` |

The backend is the **scenario engine** — it decides which screen to show, what data to display, and when transitions happen. The TUI renders whatever the backend tells it to. The REST API exposes the same scenario engine for future web frontends, with Swagger auto-generated from the shared `pkg/protocol` types.

**Multi-Language UI:** The TUI supports two locales — casual Indonesian (`id`, default) and casual English (`en`). All display strings live in `internal/tui/i18n/`. Locale is configured in `~/.waclaw/config.yaml` and switchable at runtime via Ctrl+K or `l` key.

**DRY Convention:** The codebase follows a strict Don't Repeat Yourself principle for types, interfaces, and constants:

1. **Every domain concept gets its own named Go type.** Severity levels, notification categories, worker phases, and lead lifecycle stages each have a distinct type in `pkg/protocol/types.go` — they never share `StateID` with screen states. This gives compile-time type safety and semantic clarity.

2. **Screen states belong in `state.go`, domain types belong in `types.go`.** `StateID` is reserved exclusively for screen-level UI states (which view is active). Business-domain concepts that have their own vocabulary (severity, notification type, worker phase, lead phase, confirmation type) live under their own named types in `types.go`.

3. **i18n keys are constants, not magic strings.** All lookup keys used with `i18n.T()` are defined as Go constants in `internal/tui/i18n/keys.go`. This catches typos at compile time and makes renaming a single-point change.

4. **No hardcoded protocol strings.** The JSON-RPC version (`"2.0"`) uses `protocol.Version`. Method names use `protocol.MethodNavigate` etc. Transport code delegates to `protocol.NewRequest()` / `protocol.NewNotification()` constructors instead of hand-building structs with hardcoded version strings.

5. **Theme field mapping is data-driven.** The `paletteFieldOrder` table in `theme.go` maps palette field names to style color pointers in a single loop — no repetitive per-field if-chains for `ApplyTheme` or `applyOverrides`.

---

## Project Structure

> **Two binaries. One module.** Each binary has its own `cmd/` entry point and `internal/` subtree. Shared protocol types live in `pkg/protocol/`. Shared transport I/O lives in `pkg/transport/` (stdio for JSON-RPC, HTTP for REST API). The `engine/` package breaks the scenario↔rpc circular dependency with interfaces. The `orchestrator/` package owns cross-domain workflows. The `schedule/` package owns shared work-hours logic. The `i18n/` package owns multi-language display strings.

```
waclaw/
├── cmd/
│   ├── backend/
│   │   └── main.go                        # Backend binary entry: scenario engine → RPC server + REST server → stdio + HTTP
│   └── tui/
│       └── main.go                        # TUI binary entry: RPC client → bubbletea bootstrap
│
├── pkg/
│   ├── protocol/                          # ── SHARED RPC TYPES (importable by both binaries) ───────
│   │   ├── types.go                       #   Centralized domain types: Severity, NotificationType, ConfirmationType, WorkerPhase, LeadPhase
│   │   ├── method.go                      #   Method name constants: navigate, update, notify, validate, key_press, action
│   │   ├── request.go                     #   JSON-RPC 2.0 Request type + Version constant
│   │   ├── response.go                    #   JSON-RPC 2.0 Response type + error codes
│   │   ├── notification.go               #   JSON-RPC 2.0 Notification type
│   │   ├── screen.go                      #   Screen ID enum: ScreenBoot, ScreenLogin, ..., ScreenUpdate
│   │   ├── state.go                       #   Screen state types: boot, login, niche, scrape, ..., overlay states
│   │   └── event.go                       #   TUI → Backend event types: key_press, action, request
│   │
│   └── transport/                         # ── SHARED I/O (split from protocol — no type↔I/O coupling) ──
│       ├── stdio.go                       #   Stdio transport: newline-delimited JSON read/write (TUI)
│       ├── http.go                        #   HTTP transport: REST API server (future web frontend)
│       └── handler.go                     #   Shared handler: protocol types → HTTP responses
│
│   └── api/                               # ── SWAGGER / OPENAPI SPEC ─────────────────────────────────
│       └── openapi.yaml                   #   Auto-generated from pkg/protocol types
│
├── internal/
│   │
│   ├── backend/                           # ── BACKEND BINARY (closed source) ───────────────────────
│   │   │
│   │   ├── engine/                        #   Interface layer (breaks scenario↔rpc circular dep)
│   │   │   └── interfaces.go              #     ScenarioEngine, RPCPusher, SlotPauser, LeadRepo interfaces
│   │   │
│   │   ├── orchestrator/                  #   Cross-domain workflow coordinator
│   │   │   ├── pipeline.go                #     Full pipeline: scrape → qualify → auto-review → send
│   │   │   ├── response.go                #     Response flow: incoming → trigger → classify → update lead → notify
│   │   │   └── shield.go                  #     Shield flow: flag → pause slot → redistribute → update score → notify
│   │   │
│   │   ├── scenario/                      #   State scenario controller (the brain)
│   │   │   ├── engine.go                  #     Scenario engine: implements engine.ScenarioEngine
│   │   │   ├── timeline.go               #     Timeline control: sequence events for marketing video
│   │   │   └── mock.go                   #     Mock data feed: generate realistic lead/worker/notification data
│   │   │
│   │   ├── rpc/                           #   JSON-RPC 2.0 server
│   │   │   ├── server.go                  #     RPC server: implements engine.RPCPusher, listen on stdio
│   │   │   ├── registry.go               #     Method registry: register handlers per method name
│   │   │   └── handler.go                #     Method handlers: depends on engine.ScenarioEngine (interface)
│   │   │
│   │   ├── rest/                          #   REST API server (future web frontend)
│   │   │   ├── server.go                  #     REST server: same scenario engine, HTTP transport
│   │   │   ├── handler.go                 #     REST handlers: protocol types → JSON responses
│   │   │   └── middleware.go             #     Middleware: CORS, logging, auth
│   │   │
│   │   ├── worker/                        #   Worker pool (army)
│   │   │   ├── pool.go                    #     Worker pool manager: spawn, pause, resume, stop per niche
│   │   │   └── worker.go                  #     Single niche worker: scrape → qualify → queue → send loop
│   │   │
│   │   ├── pipeline/                      #   Batch pipeline (split from worker/)
│   │   │   └── batch.go                   #     Batch pipeline: queries → scrape → qualify → auto-review → queue → batch send → wait → loop
│   │   │
│   │   ├── schedule/                      #   Shared scheduling (extracted from worker/ + sender/)
│   │   │   └── scheduler.go              #     Work hours guard + interval scheduling: 09:00-17:00 WIB, off-hours queue
│   │   │
│   │   ├── niche/                         #   Niche system (file-based)
│   │   │   ├── niche.go                   #     Niche model: name, targets, areas, filters, closing_triggers, scoring
│   │   │   ├── loader.go                  #     Load all niches from ~/.waclaw/niches/*/niche.yaml
│   │   │   ├── generator.go              #     Auto-generate niche.yaml + ice_breaker/ + offer/ from explorer
│   │   │   ├── explorer.go               #     Niche explorer: browse categories, live search (WA Biz Dir + GMaps)
│   │   │   └── area.go                    #     Granular area model: city, radius, kecamatan list
│   │   │
│   │   ├── lead/                          #   Lead lifecycle state machine
│   │   │   ├── lead.go                    #     Lead model: business data, scores, timeline, status
│   │   │   ├── state.go                   #     State machine: baru → wa_validated → ice_breaker_sent → responded → offer_sent → converted
│   │   │   ├── lifecycle.go              #     Lifecycle transitions: no_response → follow_up_1 → follow_up_2 → cold
│   │   │   ├── recontact.go              #     Re-contact engine: responded + dingin → 7 hari jeda → re_contact
│   │   │   ├── scorer.go                 #     Lead scoring: has_instagram, no_website, rating, review_count weights
│   │   │   └── storage.go                #     LeadRepo interface: consumer-defined contract, implemented by database/
│   │   │
│   │   ├── scrape/                        #   Google Maps scraper
│   │   │   ├── scraper.go                 #     Google Maps scraper: query → parse → deduplicate → qualify
│   │   │   ├── query.go                   #     Query builder: targets × areas → search queries
│   │   │   ├── qualifier.go              #     Filter engine: has_website, has_instagram, rating range, review count
│   │   │   ├── dedup.go                   #     Duplicate detector: cross-niche, cross-batch deduplication
│   │   │   └── throttle.go               #     Rate limit handler: GMaps throttle detection, auto-backoff, retry
│   │   │
│   │   ├── sender/                        #   Message sender + WA rotator
│   │   │   ├── sender.go                  #     Batch sender: queue → template pick → WA slot → send → confirm
│   │   │   ├── rotator.go                 #     WA number rotator: round-robin + cooldown, load balancing
│   │   │   ├── pauser.go                  #     SlotPauser interface: consumer-defined, lets antiban/ pause slots without importing sender/
│   │   │   └── tracker.go                #     Rate tracker: per slot/hour/day, cooldown countdowns
│   │   │
│   │   ├── template/                      #   Template rotation engine
│   │   │   ├── template.go               #     Template model: type (ice_breaker, follow_up, offer), variant files
│   │   │   ├── loader.go                 #     Load template variants from niche folders (variant_*.md)
│   │   │   ├── rotator.go                #     Rotation engine: round-robin or random per send
│   │   │   ├── renderer.go               #     Placeholder substitution: {{.Title}}, {{.Category}}, {{.Address}}, {{.City}}, {{.Rating}}, {{.Reviews}}, {{.Area}}
│   │   │   └── snippet.go                #     Snippet loader: ~/.waclaw/snippets.md quick replies
│   │   │
│   │   ├── followup/                      #   Follow-up persistence engine
│   │   │   ├── engine.go                  #     Auto follow-up: detect due leads, schedule, send with variant rotation
│   │   │   ├── cold.go                    #     Cold detection: 2x follow-up no response → auto-tandai dingin
│   │   │   └── limit.go                   #     Follow-up guard: max 3 messages lifetime, 24h gap, must use different variant
│   │   │
│   │   ├── antiban/                       #   Anti-ban shield system
│   │   │   ├── shield.go                  #     Health score calculator: aggregate per-slot metrics → 0-100 score
│   │   │   ├── spam_guard.go             #     Spam guard: per lead limits, do_not_contact, duplicate cross-niche, re-contact delay
│   │   │   ├── pattern_guard.go          #     Pattern detection: template rotation, time variance, emoji variation, paragraph shuffle
│   │   │   ├── flag_detector.go          #     WA flag detection: auto-pause flagged slot via sender.SlotPauser (interface)
│   │   │   └── donotcontact.go           #     Do-not-contact manager: ~/.waclaw/do_not_contact.yaml, auto-populate on stop trigger
│   │   │
│   │   ├── wa/                            #   WhatsApp connection layer
│   │   │   ├── connection.go             #     whatsmeow client: connect, QR generation, multi-slot management
│   │   │   ├── slot.go                    #     WA slot model: number, status, session data, health
│   │   │   ├── validator.go              #     WA pre-validation: check-registration / send-silent before queuing
│   │   │   ├── listener.go               #     Incoming message listener: response detection, auto-reply detection
│   │   │   └── session.go                #     Session persistence: ~/.waclaw/wa_slots/slot_*.yaml
│   │   │
│   │   ├── license/                       #   License gate system
│   │   │   ├── license.go                #     License model: key, device, activation date, expiration
│   │   │   ├── validator.go              #     License validation: server check, offline grace (72h), version binding
│   │   │   ├── device.go                 #     Device fingerprint: identify unique machine, conflict detection
│   │   │   └── store.go                  #     License file store: ~/.waclaw/license.md read/write
│   │   │
│   │   ├── config/                        #   Config management
│   │   │   ├── config.go                 #     Main config model: anti_ban, spam_guard, schedule, work_hours
│   │   │   ├── loader.go                 #     YAML loader: config.yaml, theme.yaml, niche.yaml
│   │   │   ├── backup.go                 #     Auto-backup: config.yaml.bak on every successful reload
│   │   │   ├── reload.go                 #     Hot reload: watch file changes, re-validate, apply without restart
│   │   │   └── migrator.go              #     Config migration: deprecated fields, version upgrades
│   │   │
│   │   ├── validator/                     #   Config validation engine
│   │   │   ├── validator.go              #     Validation orchestrator: check all config files sequentially
│   │   │   ├── schema.go                 #     Schema definitions: required fields, types, value ranges
│   │   │   ├── template_check.go        #     Template validation: required placeholders, encoding, non-empty
│   │   │   ├── niche_check.go           #     Niche validation: targets, areas, closing_triggers, scoring
│   │   │   └── report.go                 #     Validation report: errors, warnings, file/line pointers, fix suggestions
│   │   │
│   │   ├── notification/                  #   Notification system (backend dispatch)
│   │   │   ├── dispatcher.go             #     Notification dispatcher: queue, prioritize, push to TUI via RPC
│   │   │   ├── types.go                  #     17 notification types: response_masuk, scrape_selesai, wa_flag, etc.
│   │   │   └── classifier.go            #     Severity classifier: critical (3s hold), positive (10s), neutral (5s), informative (7s)
│   │   │
│   │   ├── trigger/                       #   Closing triggers engine
│   │   │   ├── trigger.go                #     Trigger model: deal, hot_lead, stop — per-niche config
│   │   │   ├── matcher.go                #     Pattern matcher: case-insensitive substring, regex support
│   │   │   └── action.go                 #     Auto-actions: deal → auto-flag, hot_lead → auto-prioritize, stop → auto-block
│   │   │
│   │   ├── update/                        #   Update & upgrade system
│   │   │   ├── checker.go                #     Version checker: background startup check, non-blocking
│   │   │   ├── downloader.go            #     Download manager: progress tracking, checksum verification
│   │   │   ├── installer.go             #     Installer: backup old binary, replace, restart
│   │   │   └── version.go               #     Version model: major.minor.patch, license version binding
│   │   │
│   │   ├── history/                       #   History & timeline
│   │   │   ├── timeline.go               #     Event timeline: per-day activity log
│   │   │   ├── stats.go                  #     Statistics: daily/weekly aggregates, conversion rates, best times
│   │   │   └── insight.go                #     Insight engine: "selasa jam 10 = waktu terbaik"
│   │   │
│   │   └── database/                      #   Lead database (SQLite)
│   │       ├── db.go                      #     SQLite connection, migrations, query helpers
│   │       ├── lead_repo.go              #     SQLiteLeadRepo: implements lead.LeadRepo interface
│   │       ├── timeline_repo.go          #     Timeline events: log every state transition
│   │       └── stats_repo.go             #     Aggregate queries: daily counts, conversion rates, niche performance
│   │
│   └── tui/                               # ── TUI BINARY (this repo, open source) ─────────────────
│       │
│       ├── app.go                         #   bubbletea.Model root: screen router, key dispatch
│       ├── router.go                      #   Screen navigation: push/pop/replace transitions
│       ├── theme.go                       #   theme.yaml loader → lipgloss color tokens
│       ├── keymap.go                      #   Global keyboard grammar (↑↓ ↵ s q p r / ? v l h ` u Ctrl+K esc)
│       ├── animation.go                   #   Micro-interaction engine: slide, pulse, morph, particle
│       ├── transition.go                  #   Screen transition: horizontal slide 300ms, cross-fade 200ms
│       │
│       ├── bus/                           #   Internal event bus (decouples rpc/ from screen/ packages)
│       │   └── bus.go                     #     Publish/Subscribe: rpc handler emits tea.Msg, app.go routes to screens
│       │
│       ├── rpc/                           #   JSON-RPC 2.0 client (connects to backend binary)
│       │   ├── client.go                  #     RPC client: connect to backend over stdio
│       │   └── handler.go                #     Translates RPC responses → tea.Msg values, publishes to bus
│       │
│       ├── i18n/                          #   Multi-language display strings (casual Indonesian + casual English)
│       │   ├── keys.go                    #     Centralized i18n key constants (compile-time safe)
│       │   ├── i18n.go                    #     T(key) → looks up current locale, runtime switchable
│       │   ├── en.go                      #     English locale map
│       │   └── id.go                      #     Indonesian locale map (default)
│       │
│       ├── screen/                        #   20 screen models in domain-grouped sub-packages
│       │   │
│       │   ├── onboarding/               #   Screens 1-2: first contact
│       │   │   ├── boot.go               #     Screen 1: BOOT — first_time, returning, +5 variants
│       │   │   └── login.go              #     Screen 2: LOGIN — QR scan, multi-slot WA rotator
│       │   │
│       │   ├── niche/                     #   Screens 3, 19: identity & discovery
│       │   │   ├── select.go             #     Screen 3: NICHE SELECT — multi-select, filters, config error
│       │   │   └── explorer.go           #     Screen 19: NICHE EXPLORER — browse, search, auto-generate
│       │   │
│       │   ├── pipeline/                  #   Screens 4-6: lead pipeline
│       │   │   ├── scrape.go             #     Screen 4: SCRAPE — multi-niche, high-value reveal, batch cascade
│       │   │   ├── review.go             #     Screen 5: LEAD REVIEW — optional manual override
│       │   │   └── send.go              #     Screen 6: SEND — auto-pilot, WA rotator, rate limits
│       │   │
│       │   ├── monitor/                   #   Screens 7-8: command center & reward
│       │   │   ├── dashboard.go          #     Screen 7: MONITOR — command center, ambient data rain
│       │   │   └── response.go           #     Screen 8: RESPONSE — closing triggers, conversion drama
│       │   │
│       │   ├── data/                      #   Screens 9-10: archive & armory
│       │   │   ├── leads_db.go           #     Screen 9: LEADS DATABASE — archive, filter, follow-up due
│       │   │   └── template_mgr.go       #     Screen 10: TEMPLATE MANAGER — variant preview, validation
│       │   │
│       │   ├── infra/                     #   Screens 11-14: infrastructure & safety
│       │   │   ├── workers.go            #     Screen 11: WORKERS — pipeline visualizer, add/pause
│       │   │   ├── antiban.go            #     Screen 12: ANTI-BAN — shield art, health score, spam guard
│       │   │   ├── settings.go           #     Screen 13: SETTINGS — config reference card
│       │   │   └── guardrail.go          #     Screen 14: GUARDRAIL — config validation, errors/warnings
│       │   │
│       │   ├── comms/                     #   Screens 15-17: communication
│       │   │   ├── compose.go            #     Screen 15: COMPOSE — custom reply modal overlay
│       │   │   ├── history.go            #     Screen 16: HISTORY — timeline, weekly mini charts
│       │   │   └── followup.go           #     Screen 17: FOLLOW-UP — persistence dashboard, cold list
│       │   │
│       │   ├── license/                   #   Screen 18: gate
│       │   │   └── license.go            #     Screen 18: LICENSE — hard gate, device conflict, offline grace
│       │   │
│       │   └── update/                    #   Screen 20: renewal
│       │       └── update.go             #     Screen 20: UPDATE & UPGRADE — minor free, major new license
│       │
│       ├── overlay/                       #   Global overlays (not screens)
│       │   ├── nerd_stats.go              #     `` toggle: hidden → minimal → expanded → hidden
│       │   ├── cmd_palette.go             #     Ctrl+K: fuzzy search, recently used, quick actions
│       │   ├── notification.go            #     Toast overlay: critical/positive/neutral/informative
│       │   ├── confirmation.go            #     Confirmation: bulk_offer, bulk_delete, bulk_archive, force_disconnect
│       │   └── shortcuts.go               #     `?` overlay: keyboard grammar cheat sheet
│       │
│       ├── component/                     #   Reusable TUI components
│       │   ├── progress_bar.go            #     Gradient sweep progress bars
│       │   ├── shield_art.go              #     Dynamic ASCII shield (health-based fill level)
│       │   ├── particle.go                #     Particle cascade system (conversion drama)
│       │   ├── data_rain.go               #     Ambient faint number scroll (monitor background)
│       │   ├── breathing.go               #     Opacity pulse engine (0.9→1.0→0.9, 4s cycle)
│       │   ├── list_select.go             #     Multi-select list with checkbox states
│       │   ├── stat_card.go               #     Dashboard stat with live increment + scale bump
│       │   ├── timeline.go                #     Sequential event timeline (stagger fade-in)
│       │   ├── mini_chart.go              #     Weekly bar charts (history screen)
│       │   ├── qr_display.go              #     QR code renderer with pixel dissolve
│       │   ├── search_input.go            #     Fuzzy search input with debounce
│       │   └── template_preview.go        #     Template text with placeholder substitution
│       │
│       ├── style/                         #   Lipgloss style definitions
│       │   ├── colors.go                  #     Color tokens: bg, text, text_muted, text_dim, success, warning, danger, accent, gold, celebration
│       │   ├── layout.go                  #     Vertical borderless layout: spacing units, section gaps
│       │   └── typography.go              #     Weight hierarchy: bold primary, muted secondary, dim tertiary
│       │
│       └── testutil/                      #   Test helpers
│           ├── fake_rpc.go               #     Fake RPC client: replay canned responses
│           └── screen_helper.go          #     Screen model test harness: init, key press, assert state
│
├── scripts/                              # Build & automation
│   ├── demo.sh                           #   Run TUI with demo backend
│   └── release.sh                        #   Build cross-platform binaries
│
├── doc/                                   # ── DOCUMENTATION ──────────────────────────
│   ├── README.md                          #   Documentation index & navigation
│   │
│   ├── # ── FOUNDATIONS ─────────────────────────────────────────
│   ├── 00-philosophy-and-design.md        #   Core philosophy: ARMY IN THE BACKGROUND
│   │
│   ├── # ── SCREENS — ONBOARDING ────────────────────────────────
│   ├── 01-screens-onboarding-boot-login.md       # Screen 1: BOOT, Screen 2: LOGIN
│   ├── 13-screens-license.md                     # Screen 18: LICENSE gate
│   │
│   ├── # ── SCREENS — NICHE SETUP ──────────────────────────────
│   ├── 02-screens-niche-select.md                # Screen 3: NICHE SELECT
│   ├── 11-screens-niche-explorer.md              # Screen 19: NICHE EXPLORER
│   │
│   ├── # ── SCREENS — LEAD PIPELINE ────────────────────────────
│   ├── 03-screens-scrape.md                      # Screen 4: SCRAPE
│   ├── 04-screens-lead-review-send.md            # Screen 5: LEAD REVIEW, Screen 6: SEND
│   │
│   ├── # ── SCREENS — MONITOR & RESPONSE ───────────────────────
│   ├── 05-screens-monitor-response.md            # Screen 7: MONITOR, Screen 8: RESPONSE
│   │
│   ├── # ── SCREENS — DATA & TEMPLATES ─────────────────────────
│   ├── 06-screens-database-templates.md          # Screen 9: LEADS DATABASE, Screen 10: TEMPLATES
│   │
│   ├── # ── SCREENS — INFRASTRUCTURE & SAFETY ──────────────────
│   ├── 07-screens-workers-antiban.md             # Screen 11: WORKERS, Screen 12: ANTI-BAN
│   ├── 08-screens-settings-guardrail.md          # Screen 13: SETTINGS, Screen 14: GUARDRAIL
│   │
│   ├── # ── SCREENS — COMMUNICATION ────────────────────────────
│   ├── 09-screens-communicate.md                 # Screen 15: COMPOSE, Screen 16: HISTORY, Screen 17: FOLLOW-UP
│   │
│   ├── # ── SCREENS — VERSION MANAGEMENT ───────────────────────
│   ├── 12-screens-update-upgrade.md              # Screen 20: UPDATE & UPGRADE
│   │
│   ├── # ── GLOBAL OVERLAYS ────────────────────────────────────
│   ├── 10-global-overlays.md                     # NERD STATS overlay, CTRL+K Command Palette
│   │
│   ├── # ── SYSTEM DESIGN ──────────────────────────────────────
│   ├── 14-notification-system.md                 # Notification types, confirmation overlays
│   ├── 15-micro-interactions.md                  # Animation catalog: nav, data, feedback, ambient
│   ├── 16-design-system.md                       # Color system (theme.yaml), layout system
│   ├── 17-niche-system.md                        # File-based niches: directory structure, templates
│   │
│   ├── # ── REFERENCE ──────────────────────────────────────────
│   ├── 18-screen-flow.md                         # Complete screen flow & lead lifecycle
│   ├── 19-keyboard-grammar.md                    # Keyboard shortcuts & overlay
│   ├── 20-startup-and-session.md                 # 4-second startup sequence & session end
│   ├── 21-rules-and-tech-stack.md                # Unwritten rules & Charm.sh tech stack
│   ├── 22-state-machine.md                       # State machine: lead, worker, screen, notif states
│   │
│   ├── # ── STATS ──────────────────────────────────────────────
│   ├── tui-screens.stats.md                      # Per-screen breakdown: 110 states, 20 variants
│   │
│   └── # ── ORIGINAL ──────────────────────────────────────────
│       └── tui.neuroscienced-customer-journey.md # Original monolithic spec (5280 lines)
│
├── go.mod                                 # Go module: github.com/WaClaw-App/waclaw
├── go.sum                                 # Dependency checksums
├── Makefile                               # Build targets: make backend, make tui, make test, make all
├── .golangci.yml                          # Linter config: errcheck, staticcheck, govet, revive
└── README.md                              # You are here
```

### Package Dependency Map

```
╔══════════════════════════════════════════════════════════════════╗
║  cmd/backend/main.go                                            ║
║  └── internal/backend                                           ║
║        ├── engine (interface layer — breaks circular deps)       ║
║        │     ├── ScenarioEngine interface                       ║
║        │     ├── RPCPusher interface                            ║
║        │     ├── SlotPauser interface                           ║
║        │     └── LeadRepo interface                             ║
║        │                                                        ║
║        ├── orchestrator (cross-domain workflow coordinator)      ║
║        │     ├── engine (ScenarioEngine interface)              ║
║        │     ├── lead (LeadRepo interface)                      ║
║        │     ├── sender (SlotPauser interface)                  ║
║        │     ├── trigger (trigger matching)                     ║
║        │     ├── notification (dispatch to TUI)                 ║
║        │     └── rpc (push results to TUI)                      ║
║        │                                                        ║
║        ├── scenario (scenario engine — the brain)               ║
║        │     ├── engine (implements ScenarioEngine)             ║
║        │     ├── rpc (uses RPCPusher interface — no circular)   ║
║        │     ├── notification (dispatch events to TUI)          ║
║        │     └── pkg/protocol (shared types)                    ║
║        │                                                        ║
║        ├── rpc/server                                           ║
║        │     ├── engine (ScenarioEngine interface — not concrete)║
║        │     ├── worker (worker control methods)                ║
║        │     ├── pkg/protocol (request/response types)          ║
║        │     └── pkg/transport (stdio read/write)               ║
║        │                                                        ║
║        ├── worker                                               ║
║        │     ├── pipeline (batch processing)                    ║
║        │     ├── schedule (work hours + intervals)              ║
║        │     ├── scrape (Google Maps scraper)                   ║
║        │     ├── sender (WA message sender)                     ║
║        │     ├── followup (auto follow-up)                      ║
║        │     ├── trigger (response classification)              ║
║        │     ├── niche (target/filter data)                     ║
║        │     ├── template (variant selection)                   ║
║        │     ├── antiban (rate limiting, flag detection)        ║
║        │     ├── wa (WhatsApp connection)                       ║
║        │     └── database (lead persistence via LeadRepo)       ║
║        │                                                        ║
║        ├── sender                                               ║
║        │     ├── wa (slot rotator)                              ║
║        │     ├── template (variant rendering)                   ║
║        │     ├── schedule (work hours guard)                    ║
║        │     ├── antiban (spam guard, via SlotPauser interface) ║
║        │     └── database (status updates via LeadRepo)         ║
║        │                                                        ║
║        ├── antiban                                              ║
║        │     ├── wa (read slot health)                          ║
║        │     └── sender (SlotPauser interface — no direct import)║
║        │                                                        ║
║        ├── scrape                                               ║
║        │     ├── niche (query generation)                       ║
║        │     ├── lead (scoring, qualification via LeadRepo)     ║
║        │     └── database (lead storage via LeadRepo)           ║
║        │                                                        ║
║        └── validator                                            ║
║              ├── config (schema validation)                     ║
║              ├── niche (niche.yaml validation)                  ║
║              └── template (placeholder validation)              ║
║                                                                 ║
╠══════════════════════════════════════════════════════════════════╣
║  cmd/tui/main.go                                                ║
║  └── internal/tui (app bootstrap)                               ║
║        ├── bus (event bus — decouples rpc from screens)         ║
║        │                                                        ║
║        ├── rpc/client                                           ║
║        │     ├── pkg/protocol (request/response types)          ║
║        │     ├── pkg/transport (stdio read/write)               ║
║        │     └── bus (publishes tea.Msg from RPC responses)     ║
║        │                                                        ║
║        ├── screen/* (domain-grouped sub-packages)               ║
║        │     ├── onboarding/ (boot, login)                      ║
║        │     ├── niche/ (select, explorer)                      ║
║        │     ├── pipeline/ (scrape, review, send)               ║
║        │     ├── monitor/ (dashboard, response)                 ║
║        │     ├── data/ (leads_db, template_mgr)                 ║
║        │     ├── infra/ (workers, antiban, settings, guardrail) ║
║        │     ├── comms/ (compose, history, followup)            ║
║        │     ├── license/ (license)                             ║
║        │     └── update/ (update)                               ║
║        │     └── bus (subscribes to events, no direct rpc import)║
║        │                                                        ║
║        ├── overlay/* (nerd stats, cmd palette, notifications)   ║
║        ├── component/* (reusable TUI widgets)                   ║
║        ├── style/* (lipgloss color tokens, layout, typography)  ║
║        └── testutil/* (fake_rpc, screen_helper)                 ║
║                                                                 ║
╠══════════════════════════════════════════════════════════════════╣
║  pkg/protocol (SHARED TYPES — imported by both binaries)        ║
║        ├── types.go         (domain types: Severity, NotifType, ║
║        │                     WorkerPhase, LeadPhase, ConfirmType)║
║        ├── method.go        (method name constants)             ║
║        ├── request.go       (JSON-RPC 2.0 Request + Version)   ║
║        ├── response.go      (JSON-RPC 2.0 Response + err codes)║
║        ├── notification.go  (JSON-RPC 2.0 Notification)        ║
║        ├── screen.go        (Screen ID enum)                   ║
║        ├── state.go         (screen state types only)           ║
║        └── event.go         (TUI → Backend event types)        ║
║                                                                 ║
║  pkg/transport (SHARED I/O — split from protocol)              ║
║        └── stdio.go         (newline-delimited JSON read/write) ║
╚══════════════════════════════════════════════════════════════════╝
```

### Binary Communication Flow

```
┌─────────────────────────┐                  ┌─────────────────────────┐
│   cmd/backend/main.go   │                  │     cmd/tui/main.go     │
│                         │                  │                         │
│  scenario engine        │  pkg/protocol +  │  bubbletea app          │
│  ├─ scenario/engine.go ─┤  pkg/transport   ├─ rpc/client.go          │
│  ├─ rpc/server.go       │  JSON-RPC 2.0    │  ├─ bus/bus.go          │
│  │  ├─ registry.go      │  over stdio      │  ├─ screen/* (9 groups) │
│  │  └─ handler.go       │                  │  ├─ overlay/* (5)       │
│  ├─ orchestrator/       │  ┌───────────┐   │  ├─ component/* (12)    │
│  │  ├─ pipeline.go      │  │ navigate  │   │  ├─ style/* (3)        │
│  │  ├─ response.go      │  │ update    │   │  └─ testutil/* (2)     │
│  │  └─ shield.go        │  │ notify    │──►│                         │
│  ├─ engine/interfaces.go │  │ validate  │   │  RPC → backend:         │
│  ├─ worker/pool.go      │  │ key_press │   │  • key_press            │
│  ├─ pipeline/batch.go   │  │ action    │◄──│  • action               │
│  ├─ schedule/scheduler  │  └───────────┘   │  • request              │
│  ├─ niche/loader.go     │                  │                         │
│  ├─ lead/state.go       │                  │                         │
│  ├─ scrape/scraper.go   │                  │                         │
│  ├─ sender/sender.go    │                  │                         │
│  ├─ template/rotator.go │                  │                         │
│  ├─ followup/engine.go  │                  │                         │
│  ├─ antiban/shield.go   │                  │                         │
│  ├─ wa/connection.go    │                  │                         │
│  ├─ license/validator.go│                  │                         │
│  ├─ config/loader.go    │                  │                         │
│  ├─ validator/          │                  │                         │
│  ├─ notification/       │                  │                         │
│  ├─ trigger/            │                  │                         │
│  ├─ update/             │                  │                         │
│  ├─ history/            │                  │                         │
│  └─ database/           │                  │                         │
└─────────────────────────┘                  └─────────────────────────┘
       CLOSED SOURCE                              THIS REPO
```

### Key Interfaces

```go
// ── BACKEND INTERFACE LAYER (internal/backend/engine) ──────────────────
// Breaks circular dependency: scenario and rpc both depend on interfaces, not each other.

// ScenarioEngine — implemented by scenario.Engine, consumed by rpc.Handler
type ScenarioEngine interface {
    Run(ctx context.Context) error
    Navigate(screen protocol.ScreenID, state protocol.State)
    PushNotification(n protocol.Notification)
    InjectData(screen protocol.ScreenID, data json.RawMessage)
}

// RPCPusher — implemented by rpc.Server, consumed by scenario.Engine
type RPCPusher interface {
    Push(method string, params json.RawMessage) error
}

// SlotPauser — implemented by sender.Sender, consumed by antiban.FlagDetector
// Lets antiban/ pause slots without importing sender/ (avoids coupling)
type SlotPauser interface {
    PauseSlot(slotID string) error
    RedistributeLoad(fromSlot string) error
}

// LeadRepo — consumer-defined interface in lead/, implemented by database.SQLiteLeadRepo
// Packages that need lead data depend on this interface, not on database/ directly
type LeadRepo interface {
    Insert(lead Lead) error
    UpdateStatus(id string, status State) error
    Filter(opts FilterOpts) ([]Lead, error)
}

// ── TUI BINARY (internal/tui) ──────────────────────────────────────

// Screen — every TUI screen implements bubbletea.Model + this
type Screen interface {
    tea.Model
    ID() protocol.ScreenID    // e.g. ScreenBoot, ScreenMonitor
    State() protocol.State    // current state string
    SetState(protocol.State)  // backend-driven state transition via bus event
}

// ── BACKEND BINARY (internal/backend) ──────────────────────────────

// Worker — one niche worker, runs independently in the army
type Worker interface {
    Start() error
    Pause()
    Resume()
    Stop()
    Status() WorkerStatus     // spawning|scraping|qualifying|queuing|sending|idle|paused|error
    Niche() *niche.Niche
}

// Slot — one WhatsApp number in the rotator
type Slot interface {
    Connect() error
    Disconnect()
    Send(to, message string) error
    Health() int              // 0-100 health score
    Status() SlotStatus       // aktif|cooldown|flagged|disconnected
}

// ── SHARED (pkg/protocol) ─────────────────────────────────────────

// RPCMethod — JSON-RPC 2.0 method handler (used by backend rpc/server)
type RPCMethod interface {
    Method() string
    Handle(params json.RawMessage) (interface{}, error)
}
```

---

## The 20 Screens

| # | Screen | Tagline | Key States | Views |
|---|--------|---------|------------|-------|
| 1 | **BOOT** | First Impression | `first_time`, `returning`, `+5 variants` | 7 |
| 2 | **LOGIN** | Trust | `qr_waiting` → `qr_scanned` → `login_success` | 5 |
| 3 | **NICHE SELECT** | Identity | `niche_list`, `multi_selected`, `custom`, `config_error` | 6 |
| 4 | **SCRAPE** | Anticipation | `scraping_active`, `multi_active`, `high_value_reveal` | 12 |
| 5 | **LEAD REVIEW** | Curated | `reviewing`, `lead_detail`, `template_preview` | 4 |
| 6 | **SEND** | Auto-Pilot | `sending_active`, `off_hours`, `rate_limited` | 8 |
| 7 | **MONITOR** | Command Center | `live_dashboard`, `idle_background`, `night` | 6 |
| 8 | **RESPONSE** | Reward | `positive`, `deal_detected`, `conversion` | 11 |
| 9 | **LEADS DATABASE** | Archive | `leads_list`, `filtered`, `lead_full_detail` | 7 |
| 10 | **TEMPLATE MANAGER** | Armory | `template_list`, `preview`, `validation_error` | 4 |
| 11 | **WORKERS** | Pipeline Visualizer | `workers_overview`, `worker_detail` | 4 |
| 12 | **ANTI-BAN** | Shield | `shield_overview`, `warning`, `danger` | 5 |
| 13 | **SETTINGS** | Config Reference | `settings_overview`, `reload`, `reload_error` | 4 |
| 14 | **GUARDRAIL** | Config Validation | `validation_clean`, `errors`, `warnings` | 5 |
| 15 | **COMPOSE** | Voice | `compose_draft`, `preview`, `template_pick` | 3 |
| 16 | **HISTORY** | Timeline | `history_today`, `history_week`, `day_detail` | 3 |
| 17 | **FOLLOW-UP** | Persistence | `followup_dashboard`, `sending`, `cold_list` | 6 |
| 18 | **LICENSE** | Gate | `license_input`, `valid`, `expired`, `device_conflict` | 7 |
| 19 | **NICHE EXPLORER** | Discovery | `explorer_browse`, `search`, `generated` | 6 |
| 20 | **UPDATE & UPGRADE** | Renewal | `update_available`, `upgrade_available` | 7 |

**Plus 2 Global Overlays:**
- **NERD STATS** — Toggle with `` ` `` — CPU, RAM, goroutines, DB size
- **CTRL+K Command Palette** — Fuzzy search, navigate, execute

**Total: 20 screens, 110 states, 20 variants = 130 views, 17 notification types, 4 confirmation overlays**

---

## TUI State Management

The backend owns **all** state. The TUI is a thin rendering layer that reflects whatever the backend pushes via RPC. The TUI never mutates domain state on its own — it dispatches user intents (`key_press`, `action`, `request`) and waits for a state update.

### State Ownership Model

```
┌──────────────────────────────────────────────────────────────────┐
│                      STATE OWNERSHIP                             │
│                                                                  │
│  BACKEND (source of truth)          TUI (projection)            │
│  ─────────────────────────          ───────────────             │
│  • Current screen ID                • Render current screen     │
│  • Current screen state             • Apply lipgloss styles     │
│  • Worker pool status               • Run animations            │
│  • Lead lifecycle states            • Dispatch key presses       │
│  • WA slot health/flags             • Show overlay on request   │
│  • Notification queue               • Display toasts            │
│  • Config validation results        • Render errors/warnings    │
│  • License status                   • Render license gate       │
│  • Update availability              • Show update prompt        │
│                                                                  │
│  Backend pushes ──► navigate/update/notify/validate              │
│  TUI pushes    ──► key_press/action/request                      │
│                                                                  │
│  TUI local state ONLY:                                            │
│  • Cursor position (↑↓)                                          │
│  • Scroll offset                                                  │
│  • Overlay toggle states (nerd stats, cmd palette)               │
│  • Animation frame counters                                       │
│  • Search input text (cmd palette, lead filter)                  │
│  • Confirmation overlay pending response                         │
└──────────────────────────────────────────────────────────────────┘
```

### State Machine Categories

Four independent state machines that the backend controls in parallel:

#### 1. Screen State Machine

Screen transitions are backend-driven. The TUI sends a user intent, the backend decides the next screen + state.

```
SCREEN TRANSITION FLOW:

  TUI                          Backend                          TUI
  ───                          ───────                          ───
  User presses ↵ on niche ──► key_press("↵") ──► Backend decides:
                                                         navigate(screen=SCRAPE,
                                                                  state=scraping_active)
                                                    ◄── JSON-RPC response
  Screen SCRAPE renders ◄── SetState(scraping_active)
  with new data injected
```

| Screen | States |
|--------|--------|
| **BOOT** | `first_time` \| `returning` \| `returning+response` \| `returning+error` \| `returning+config_error` \| `returning+license_expired` \| `returning+device_conflict` |
| **LOGIN** | `qr_waiting` \| `qr_scanned` \| `login_success` \| `login_expired` \| `login_failed` |
| **NICHE SELECT** | `niche_list` \| `niche_multi_selected` \| `niche_custom` \| `niche_edit_filters` \| `niche_config_error` |
| **SCRAPE** | `scraping_active` \| `scraping_multi_active` \| `scraping_multi_staggered` \| `scrape_idle` \| `scrape_empty` \| `scrape_error` \| `scrape_gmaps_limited` \| `scrape_auto_approved` \| `scrape_high_value_reveal` \| `scrape_batch_complete` |
| **LEAD REVIEW** | `reviewing` \| `lead_detail` \| `template_preview` \| `queue_complete` |
| **SEND** | `sending_active` \| `sending_paused` \| `sending_off_hours` \| `sending_rate_limited` \| `sending_daily_limit` \| `sending_failed` \| `sending_all_slots_down` \| `sending_with_response` |
| **MONITOR** | `live_dashboard` \| `idle_background` \| `dashboard_night` \| `dashboard_error` \| `dashboard_empty` \| `dashboard_with_pending_responses` |
| **RESPONSE** | `response_positive` \| `response_curious` \| `response_negative` \| `response_maybe` \| `response_auto_reply` \| `offer_preview` \| `response_multi_queue` \| `conversion` |
| **LEADS DB** | `leads_list` \| `leads_filtered` \| `lead_full_detail` \| `lead_follow_up_due` \| `lead_cold` \| `lead_never_contacted` \| `lead_converted` |
| **TEMPLATE** | `template_list` \| `template_preview` \| `template_edit_hint` \| `template_validation_error` |
| **WORKERS** | `workers_overview` \| `worker_detail` \| `worker_add_niche` \| `worker_paused` |
| **ANTI-BAN** | `shield_overview` \| `shield_warning` \| `shield_danger` \| `shield_slot_detail` \| `shield_settings` |
| **SETTINGS** | `settings_overview` \| `settings_edit` \| `settings_reload` \| `settings_reload_error` |
| **GUARDRAIL** | `validation_clean` \| `validation_errors` \| `validation_warnings` \| `validation_fix` \| `validation_first_time` |
| **COMPOSE** | `compose_draft` \| `compose_preview` \| `compose_template_pick` |
| **HISTORY** | `history_today` \| `history_week` \| `history_day_detail` |
| **FOLLOW-UP** | `followup_dashboard` \| `followup_niche_detail` \| `followup_sending` \| `followup_empty` \| `followup_cold_list` \| `followup_recontact` |
| **NICHE EXPLORER** | `explorer_browse` \| `explorer_search` \| `explorer_category_detail` \| `explorer_generating` \| `explorer_generated` |
| **UPDATE** | `update_available` \| `update_downloading` \| `update_ready` \| `upgrade_available` \| `upgrade_license_input` \| `license_expired_with_upgrade` |
| **LICENSE** | `license_input` \| `license_validating` \| `license_valid` \| `license_invalid` \| `license_expired` \| `license_device_conflict` \| `license_server_error` |

#### 2. Worker State Machine (per niche)

```
WORKER STATES (per niche):
  spawning → scraping → qualifying → queuing → sending → idle (loop)
    │            │           │           │          │
    │            └→ error → retry       │          └→ rate_limited
    │                                    └→ paused (manual)
    └→ config_error → paused → fixed → spawning
    └→ stopped (manual)

BATCH PIPELINE (per worker):
  queries.md → scrape → qualify → auto-review → queue → batch send → wait → loop
                                                                  └→ response → offer → deal
```

#### 3. Lead Lifecycle State Machine

```
LEAD STATES:
  baru → ice_breaker_sent → responded → offer_sent → converted
    │         │                 │            │
    │         │                 └→ negative → archived
    │         │                 └→ auto_reply → skipped
    │         └→ no_response → follow_up_1 → follow_up_2 → cold
    │         │                                    │            │
    │         │                                    └→ responded └→ offer_sent
    │         └→ failed → retry (max 3x) → dead
    └→ blocked (manual)

  RE-CONTACT: responded + dingin → 7 hari jeda → re_contact → responded
  FOLLOW-UP LIMIT: ice_breaker + follow_up_1 + follow_up_2 = max 3 per lead
  COLD: 2x follow-up tanpa response → auto-tandai dingin
```

#### 4. Config Validation State Machine

```
CONFIG VALIDATION STATES:
  unchecked → validating → clean | errors | warnings
    │                         │        │         │
    │                         │        └→ fix → revalidate → clean
    │                         └→ auto-dismiss (from boot)
    └→ first_time → generate template → revalidate → clean
```

### Notification State Machine

Notifications are queued and dispatched one at a time by the backend.

```
NOTIFICATION LIFECYCLE:

  Backend event ──► notification dispatched ──► TUI overlay renders
                          │                           │
                          │                     user responds? ──► action sent to backend
                          │                           │
                          │                     auto-dismiss timer?
                          │                           │
                     queue next ◄───────────── TUI reports dismissed
```

| Severity | Hold Time | Auto-Dismiss | Examples |
|----------|-----------|-------------|----------|
| **Critical** | 3s mandatory | 10s | `wa_disconnect`, `wa_flag`, `config_error`, `validation_error`, `license_expired`, `device_conflict` |
| **Positive** | — | 10s (15s for updates) | `response_masuk`, `multi_response`, `streak_milestone`, `update_available` |
| **Neutral** | — | 5s | `scrape_selesai`, `batch_selesai`, `followup_terjadwal` |
| **Informative** | — | 7s (20s for upgrades) | `lead_dingin`, `health_drop`, `limit_harian`, `upgrade_available` |

**17 Notification Types:**
`response_masuk` \| `multi_response` \| `scrape_selesai` \| `batch_selesai` \| `wa_disconnect` \| `wa_flag` \| `health_drop` \| `limit_harian` \| `streak_milestone` \| `config_error` \| `validation_error` \| `license_expired` \| `device_conflict` \| `followup_terjadwal` \| `lead_dingin` \| `update_available` \| `upgrade_available`

**4 Confirmation Overlays:**
`bulk_offer` \| `bulk_delete` \| `bulk_archive` \| `force_device_disconnect`

### Overlay State Machines

Global overlays maintain their own local toggle state in the TUI — not backend-driven.

#### Nerd Stats Toggle (` key)

```
hidden ──[` key]──► minimal ──[` key]──► expanded ──[` key]──► hidden
                        │                       │
                        └──[30s timeout]──► hidden
```

- **minimal**: 1-line footer (`CPU 12% · RAM 134MB · Goroutines 23 · DB 2.4MB · Uptime 4j 12m`)
- **expanded**: 3-line panel with mini bar charts, RAM/goroutine thresholds (80% = warning amber, >80% = danger red)
- **Data source**: Go runtime (`runtime.ReadMemStats`, `runtime.NumGoroutine`) + SQLite DB stats

#### Command Palette (Ctrl+K)

```
cmd_closed ──[Ctrl+K]──► cmd_open ──[select]──► cmd_executing ──► cmd_closed
                              │                                            │
                              ├──[esc/Ctrl+K]──► cmd_closed               │
                              ├──[no match]──► cmd_empty                  │
                              └──[has recent]──► cmd_with_recent          │
                                                                         │
                         cmd_quick_action ──[⚡ execute immediately]──► cmd_closed
```

- **Fuzzy search**: fzf-style scoring (exact > prefix > substring > fuzzy), 50ms debounce
- **Context-aware**: commands relevant to current screen ranked higher
- **Recently used**: last 3 commands always shown first
- **Quick actions**: `⚡` tagged commands execute without screen navigation

### RPC Method Mapping

| RPC Method | Direction | State Effect |
|------------|-----------|-------------|
| `navigate` | Backend → TUI | Change current screen + screen state |
| `update` | Backend → TUI | Inject data into current screen without navigation |
| `notify` | Backend → TUI | Queue notification overlay for display |
| `validate` | Backend → TUI | Push config validation results to GUARDRAIL screen |
| `key_press` | TUI → Backend | Report user key input for backend decision |
| `action` | TUI → Backend | Report user action (approve, skip, classify, etc.) |
| `request` | TUI → Backend | Request specific data (lead detail, stats, etc.) |

---

## Design Philosophy

### ARMY IN THE BACKGROUND

```
WaClaw itu bukan satu asisten.
WaClaw itu army — satu worker per niche, jalan paralel 24/7.

Lu = jenderal. WaClaw = pasukan.
Lu tentuin strategi, mereka eksekusi.
Lu nggak perlu micromanage — tiap worker otonom.
Lu cuma di-interupt kalau ada yang perlu keputusan bos.
```

### Core Principles

| Principle | Meaning |
|-----------|---------|
| **Auto-pilot = default** | If you don't touch the keyboard for 1 hour and 3 niches are still finding leads, WaClaw is working |
| **Notification-first** | WaClaw finds moments to ask you, not the other way around |
| **One-key decisions** | Every interrupt only needs 1 key. `↵` = agree, `s` = skip |
| **Validate early, fail loudly** | Broken config = paused army. Silent errors = invisible disaster |
| **Every message rotatable** | One template = many variants. Single variant = WA pattern detection |
| **Assume nothing about numbers** | Scraped Google Maps number ≠ has WhatsApp. Validate first |
| **Follow-up = persistence, not spam** | Max 3 messages lifetime, 24h gap, different variant each time |
| **One license, one device** | Shared = stopped. Fair for everyone |
| **Minor update = free. Major upgrade = new license** | v1 army = v1 license. v2 = different product |

### Design Language

```
Hierarchy  = Brightness + Size + Motion
Separation = Vertical rhythm, never lines
Navigation = Muscle memory, never menus
Feedback   = Felt, not read
Language   = Netizen indo, bukan bahasa buku
Validation = Early, visible, actionable
Rotation   = Every message, every time
```

**No borders. No boxes. Only space, weight, and motion.**

---

## Screen Flow

```
BOOT (first time)
  │
  ├─→ LICENSE ── belum ada key? ──→ LICENSE INPUT
  │     │
  │     ▼
  ├─→ VALIDATION ── config missing? ──→ CONFIG SETUP
  │     │
  │     ▼
  │   LOGIN ──→ NICHE SELECT ──→ VALIDATION ──→ SCRAPE (auto)
  │
BOOT (returning)
  │
  ▼
LICENSE CHECK ── VALIDATION ── MONITOR (home base)
                                  │
                                  ├── auto: SCRAPER runs on interval
                                  ├── auto: SENDER runs during work hours
                                  ├── auto: FOLLOW-UP for unanswered leads
                                  ├── auto: NOTIFICATION on events
                                  │         ├── closing_triggers.deal → auto-flag
                                  │         ├── closing_triggers.stop → auto-block
                                  │         └── closing_triggers.hot_lead → auto-prioritize
                                  ├── interrupt: RESPONSE → you approve
                                  └── manual: DATABASE, TEMPLATE, SETTINGS, HISTORY

LEAD LIFECYCLE:
baru → wa_validated → ice_breaker_sent → responded → offer_sent → converted
                          │                   │
                          │                   └→ negative → archived
                          └→ no_response → follow_up_1 → follow_up_2 → cold
```

---

## Tech Stack

### Go Runtime

| Requirement | Version |
|-------------|---------|
| Go | 1.22+ |

### Binary Architecture

| Binary | Entry Point | Source | Purpose |
|--------|-------------|--------|---------|
| `waclaw-backend` | `cmd/backend/main.go` | Closed source | Scenario engine: drive TUI states, mock data, timeline control |
| `waclaw-tui` | `cmd/tui/main.go` | This repo (open source) | TUI rendering: 20 screens, overlays, micro-interactions |

### Shared Protocol Layer

| Library | Purpose | Package |
|---------|---------|---------|
| Go stdlib `encoding/json` | JSON-RPC 2.0 message serialization | `pkg/protocol` |
| Go stdlib `os` | Stdio transport (stdin/stdout) | `pkg/protocol/transport` |

Both binaries import `pkg/protocol` for shared request/response/notification types, screen IDs, state enums, and the stdio transport.

### TUI Binary — Charm.sh Ecosystem

| Library | Purpose | Package |
|---------|---------|---------|
| [bubbletea](https://github.com/charmbracelet/bubbletea) | MVC framework | `internal/tui` |
| [lipgloss](https://github.com/charmbracelet/lipgloss) | Styling & layout | `internal/tui/style` |
| [bubbles](https://github.com/charmbracelet/bubbles) | Pre-built components | `internal/tui/component` |
| [glamour](https://github.com/charmbracelet/glamour) | Markdown rendering | `internal/tui/screen/template_mgr` |
| [huh](https://github.com/charmbracelet/huh) | Forms & prompts | `internal/tui/screen/license`, `compose` |

**Why Charm.sh?** Renders the same everywhere. Components designed by people who understand terminals. Borderless aesthetic built-in — we just push it further.

### Backend Binary — Communication Layer

| Component | Protocol | Package |
|-----------|----------|---------|
| RPC server | JSON-RPC 2.0 over stdio | `internal/backend/rpc` |
| Scenario engine | Drives TUI via RPC | `internal/backend/scenario` |
| Notification dispatch | Push to TUI via RPC | `internal/backend/notification` |

### Backend Binary — WhatsApp Layer

| Library | Purpose | Package |
|---------|---------|---------|
| [whatsmeow](https://github.com/tulir/whatsmeow) | WhatsApp Web client — QR login, send, listen | `internal/backend/wa` |

### Backend Binary — Data Layer

| Library | Purpose | Package |
|---------|---------|---------|
| [go-sqlite3](https://github.com/mattn/go-sqlite3) | Lead database — embedded, zero-config | `internal/backend/database` |
| [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3) | Config & niche YAML parsing | `internal/backend/config`, `internal/backend/niche` |

### Backend Binary — Scraping Layer

| Library | Purpose | Package |
|---------|---------|---------|
| [rod](https://github.com/go-rod/rod) | Google Maps headless browser scraper | `internal/backend/scrape` |

### TUI Binary — Search Layer

| Library | Purpose | Package |
|---------|---------|---------|
| [sahilm/fuzzy](https://github.com/sahilm/fuzzy) | Fuzzy search — command palette, lead filter | `internal/tui/overlay/cmd_palette` |

### Data Layer (file-based, zero UI config)

```
~/.waclaw/
├── config.yaml              # Main settings + anti_ban + spam_guard
├── config.yaml.bak          # Auto-backup on every successful reload
├── theme.yaml               # Colors & feel
├── license.md               # License key + device + expiration
├── queries.md               # Search queries
├── snippets.md              # Quick reply snippets for custom messages
├── do_not_contact.yaml      # Auto-block list (stop triggers + manual)
├── wa_slots/
│   ├── slot_1.yaml          # WA number #1 (auto-generated)
│   ├── slot_2.yaml          # WA number #2
│   └── slot_3.yaml          # WA number #3
└── niches/
    ├── _contoh/             # Example niche for reference
    ├── web_developer/
    │   ├── niche.yaml       # Filters, targets, areas, closing_triggers
    │   ├── ice_breaker/     # ROTATABLE — 1 file = 1 variant
    │   │   ├── variant_1.md
    │   │   ├── variant_2.md
    │   │   └── variant_3.md
    │   ├── follow_up/       # ROTATABLE — follow-up templates
    │   │   ├── follow_up_1.md
    │   │   ├── follow_up_2.md
    │   │   └── follow_up_3.md
    │   └── offer/           # ROTATABLE — 1 file = 1 variant
    │       ├── variant_1.md
    │       ├── variant_2.md
    │       └── variant_3.md
    ├── undangan_digital/
    └── social_media_manager/
```

---

## Keyboard Grammar

**Memorize once. Use forever. The same keys, every screen.**

| Key | Always Does |
|-----|-------------|
| `↑` / `↓` | Move between items |
| `↵` | Primary action (the most sensible one) |
| `1`-`9` | Secondary actions (varies per screen) |
| `s` | Skip / discard |
| `q` | Done / back / exit (context-dependent) |
| `p` | Pause whatever's running |
| `r` | Refresh / reload |
| `/` | Search / filter |
| `?` | Show shortcuts overlay |
| `v` | Validate all config |
| `l` | License — view / change |
| `h` | History — timeline |
| `` ` `` | Nerd stats — toggle RAM/CPU overlay |
| `u` | Update — check new version |
| `Ctrl+K` | Command palette — search & execute anything |
| `esc` | Cancel / close modal / exit compose |

> **Keyboard is a privilege, not an obligation.** If you don't touch the keyboard for 1 hour and 3 niches are still finding leads, that means WaClaw is working.

---

## Micro-Interactions

Every animation has purpose. Nothing is decoration.

| Category | Example | Duration | Feel |
|----------|---------|----------|------|
| **Navigation** | Screen transition horizontal slide | 300ms | Forward momentum |
| **Data** | Numbers increment with scale bump | 200ms | Tangible change |
| **Feedback** | Success green pulse (1.0→1.2→1.0) | 500ms | Achievement |
| **Dramatic** | Conversion: white flash + particles + bell | 1200ms | WINNING |
| **Ambient** | Data rain `░░ 3 7 1 4 ░░` on monitor | Continuous | System alive |
| **Ambient** | Breathing stats (opacity 0.9→1.0→0.9) | 4000ms cycle | Dashboard breathes |

### The Conversion — Full Drama Sequence

This is the most important screen in the entire app. Everything leads here.

| Phase | Time | What Happens |
|-------|------|-------------|
| **SHOCK** | 0-200ms | Full-screen white flash + double terminal bell `\a` |
| **REVEAL** | 200-800ms | `★ ★ ★ D E A L ! ★ ★ ★` scales from 0→1.3→1.0, particle cascade (40 particles), gold→amber→white color wave |
| **CONTEXT** | 800-1500ms | Business name + timeline fade in, trophy bounces from right, revenue glows gold 3x |
| **SETTLE** | 1500ms+ | Particles dissolve, screen settles, `↵ mark as converted` fades in |

---

## Color System

**Not a theme. A mood.**

| Token | Hex | Purpose |
|-------|-----|---------|
| `bg` | `#0A0A0B` | Almost black, softer than pure black |
| `text` | `#E8E8EC` | Primary — warm white |
| `text_muted` | `#6B6B76` | Secondary — whisper, don't shout |
| `text_dim` | `#3D3D44` | Tertiary — barely visible, still readable |
| `success` | `#34D399` | Green — not aggressive |
| `warning` | `#FBBF24` | Amber — attention, not alarm |
| `danger` | `#F87171` | Red — clear, not frightening |
| `accent` | `#818CF8` | Indigo — brand, action |
| `gold` | `#FFD700` | Jackpot & revenue — earned celebration |
| `celebration` | `#FFFFFF` | Full-screen flash — conversion only |

**Rules:** Red is ONLY for technical problems. Rejection = neutral. Gold = money. Celebration white = conversion only.

---

## Startup Sequence — 4 Seconds

```
$ waclaw

  t +0ms     Logo render per character
  t +80ms    Tagline fade in
  t +200ms   System check (WA, config, DB, license)
  t +300ms   License check (valid / expired / device conflict)
  t +400ms   Config validation (all niche.yaml, templates)
  t +700ms   Status report: ● ok ○ paused
  t +800ms   Auto-pilot: ON
  t +900ms   Army marching: workers → ● aktif
  t +1100ms  Dashboard fade in
  t +1300ms  Ready. Cursor blinks.
```

**1300ms until usable.** Every millisecond before that pulls attention. Every millisecond after that is your decision.

---

## Documentation

Full documentation lives in the [`doc/`](doc/) directory. Start with the [documentation index](doc/README.md).

### Quick Links

| Topic | File |
|-------|------|
| Philosophy & Design Language | [00-philosophy-and-design.md](doc/00-philosophy-and-design.md) |
| Screen Flow & Lead Lifecycle | [18-screen-flow.md](doc/18-screen-flow.md) |
| State Machine Reference | [22-state-machine.md](doc/22-state-machine.md) |
| Niche System & File Structure | [17-niche-system.md](doc/17-niche-system.md) |
| Color & Layout System | [16-design-system.md](doc/16-design-system.md) |
| Micro-Interactions Catalog | [15-micro-interactions.md](doc/15-micro-interactions.md) |
| Notification System | [14-notification-system.md](doc/14-notification-system.md) |
| Keyboard Grammar | [19-keyboard-grammar.md](doc/19-keyboard-grammar.md) |
| Startup & Session | [20-startup-and-session.md](doc/20-startup-and-session.md) |
| Tech Stack | [21-rules-and-tech-stack.md](doc/21-rules-and-tech-stack.md) |
| Screen Stats Breakdown | [tui-screens.stats.md](doc/tui-screens.stats.md) |

---

## Stats at a Glance

```
SCREEN              STATES  VARIANTS  TOTAL
────────────────────────────────────────────
1  BOOT                2       5        7
2  LOGIN               5       0        5
3  NICHE SELECT        5       1        6
4  SCRAPE              8       4       12
5  LEAD REVIEW         4       0        4
6  SEND                7       1        8
7  MONITOR             5       1        6
8  RESPONSE           11       0       11
9  LEADS DB            3       4        7
10 TEMPLATE            4       0        4
11 WORKERS             4       0        4
12 ANTI-BAN            5       0        5
13 SETTINGS            4       0        4
14 GUARDRAIL           4       1        5
15 COMPOSE             3       0        3
16 HISTORY             3       0        3
17 FOLLOW-UP           6       0        6
18 LICENSE             7       0        7
19 EXPLORER            5       1        6
20 UPDATE              6       1        7
────────────────────────────────────────────
TOTAL               110      20      130

+ 17 notification types
+ 4 confirmation overlays
+ 2 global overlays (nerd stats + command palette)
```

---

## The Unwritten Rules

1. Never show an empty state without a next action
2. Never use red for rejection — red = broken, rejection = neutral
3. Never animate for decoration — every animation = meaningful state change
4. Never hide rate limits — visible limits = trust, hidden = anxiety
5. Never ask for confirmation twice — `↵` = go, trust the user
6. Never break keyboard grammar — `q` = back/exit, always
7. Never show numbers without context — "4.6% conversion" not "46% selesai"
8. Never use formal language — netizen Indo, casual but clear
9. Never make the user wait without info — if waiting, say why + how long
10. Auto-pilot = default, manual = bonus
11. Config error = partial pause, never full stop
12. Celebration is earned, never given — conversion gets full drama because you genuinely won

---

<div align="center">

**WaClaw — army lu cerdas, keras, dan aman.**

_Lu cuma nonton. WaClaw yang kerja._  
_Tapi kalau lu mau intervene, 1 tombol cukup._

</div>
