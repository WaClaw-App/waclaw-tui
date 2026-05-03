# WaClaw TUI Documentation

> Split from `tui.neuroscienced-customer-journey.md` (5280 lines) into contextually-grouped files.
> Zero regressions — every line from the original is preserved exactly in the split files.

---

## File Index

### Foundations

| File | Content |
|------|---------|
| [00-philosophy-and-design.md](00-philosophy-and-design.md) | Core philosophy (ARMY IN THE BACKGROUND), design language, and TUI language principles |

### Screens — Onboarding

| File | Content |
|------|---------|
| [01-screens-onboarding-boot-login.md](01-screens-onboarding-boot-login.md) | Screen 1: BOOT → First Impression, Screen 2: LOGIN → Trust |
| [13-screens-license.md](13-screens-license.md) | Screen 18: LICENSE → Gate (hard gate system) |

### Screens — Niche Setup

| File | Content |
|------|---------|
| [02-screens-niche-select.md](02-screens-niche-select.md) | Screen 3: NICHE SELECT → Identity |
| [11-screens-niche-explorer.md](11-screens-niche-explorer.md) | Screen 19: NICHE EXPLORER → Discovery |

### Screens — Lead Pipeline

| File | Content |
|------|---------|
| [03-screens-scrape.md](03-screens-scrape.md) | Screen 4: SCRAPE → Anticipation |
| [04-screens-lead-review-send.md](04-screens-lead-review-send.md) | Screen 5: LEAD REVIEW → Curated, Screen 6: SEND → Auto-Pilot |

### Screens — Monitor & Response

| File | Content |
|------|---------|
| [05-screens-monitor-response.md](05-screens-monitor-response.md) | Screen 7: MONITOR → Command Center, Screen 8: RESPONSE → Reward (conversion drama) |

### Screens — Data & Templates

| File | Content |
|------|---------|
| [06-screens-database-templates.md](06-screens-database-templates.md) | Screen 9: LEADS DATABASE → Archive, Screen 10: TEMPLATE MANAGER → Armory |

### Screens — Infrastructure & Safety

| File | Content |
|------|---------|
| [07-screens-workers-antiban.md](07-screens-workers-antiban.md) | Screen 11: WORKERS → Pipeline Visualizer, Screen 12: ANTI-BAN → Shield |
| [08-screens-settings-guardrail.md](08-screens-settings-guardrail.md) | Screen 13: SETTINGS → Config Reference, Screen 14: GUARDRAIL → Config Validation |

### Screens — Communication

| File | Content |
|------|---------|
| [09-screens-communicate.md](09-screens-communicate.md) | Screen 15: COMPOSE → Voice, Screen 16: HISTORY → Timeline, Screen 17: FOLLOW-UP → Persistence |

### Screens — Version Management

| File | Content |
|------|---------|
| [12-screens-update-upgrade.md](12-screens-update-upgrade.md) | Screen 20: UPDATE & UPGRADE → Renewal |

### Global Overlays

| File | Content |
|------|---------|
| [10-global-overlays.md](10-global-overlays.md) | NERD STATS → Vitals overlay (` key), CTRL+K → Command Palette overlay |

### System Design

| File | Content |
|------|---------|
| [14-notification-system.md](14-notification-system.md) | Notification system, confirmation overlays, and rules |
| [15-micro-interactions.md](15-micro-interactions.md) | Micro-interactions catalog: navigation, data, feedback, ambient, dramatic reveals |
| [16-design-system.md](16-design-system.md) | Color system (theme.yaml) and layout system (vertical borderless) |
| [17-niche-system.md](17-niche-system.md) | File-based niche system: directory structure, niche.yaml, templates, snippets, do_not_contact |

### Reference

| File | Content |
|------|---------|
| [18-screen-flow.md](18-screen-flow.md) | Complete screen flow diagram and lead lifecycle |
| [19-keyboard-grammar.md](19-keyboard-grammar.md) | Keyboard grammar and shortcut overlay |
| [20-startup-and-session.md](20-startup-and-session.md) | Startup sequence (4 seconds) and session end |
| [21-rules-and-tech-stack.md](21-rules-and-tech-stack.md) | Unwritten rules and tech stack (TUI layer) |
| [22-state-machine.md](22-state-machine.md) | State machine summary: lead states, worker states, pipeline, config validation, screen states |

### Stats

| File | Content |
|------|---------|
| [tui-screens.stats.md](tui-screens.stats.md) | Per-screen breakdown: states, variants, views count, and notification/confirmation overlay stats |

### Original

| File | Content |
|------|---------|
| [tui.neuroscienced-customer-journey.md](tui.neuroscienced-customer-journey.md) | The original monolithic file (5280 lines) — kept for reference |

---

## Quick Navigation by Screen Number

| Screen | Name | File |
|--------|------|------|
| 1 | BOOT → First Impression | [01-screens-onboarding-boot-login.md](01-screens-onboarding-boot-login.md) |
| 2 | LOGIN → Trust | [01-screens-onboarding-boot-login.md](01-screens-onboarding-boot-login.md) |
| 3 | NICHE SELECT → Identity | [02-screens-niche-select.md](02-screens-niche-select.md) |
| 4 | SCRAPE → Anticipation | [03-screens-scrape.md](03-screens-scrape.md) |
| 5 | LEAD REVIEW → Curated | [04-screens-lead-review-send.md](04-screens-lead-review-send.md) |
| 6 | SEND → Auto-Pilot | [04-screens-lead-review-send.md](04-screens-lead-review-send.md) |
| 7 | MONITOR → Command Center | [05-screens-monitor-response.md](05-screens-monitor-response.md) |
| 8 | RESPONSE → Reward | [05-screens-monitor-response.md](05-screens-monitor-response.md) |
| 9 | LEADS DATABASE → Archive | [06-screens-database-templates.md](06-screens-database-templates.md) |
| 10 | TEMPLATE MANAGER → Armory | [06-screens-database-templates.md](06-screens-database-templates.md) |
| 11 | WORKERS → Pipeline Visualizer | [07-screens-workers-antiban.md](07-screens-workers-antiban.md) |
| 12 | ANTI-BAN → Shield | [07-screens-workers-antiban.md](07-screens-workers-antiban.md) |
| 13 | SETTINGS → Config Reference | [08-screens-settings-guardrail.md](08-screens-settings-guardrail.md) |
| 14 | GUARDRAIL → Config Validation | [08-screens-settings-guardrail.md](08-screens-settings-guardrail.md) |
| 15 | COMPOSE → Voice | [09-screens-communicate.md](09-screens-communicate.md) |
| 16 | HISTORY → Timeline | [09-screens-communicate.md](09-screens-communicate.md) |
| 17 | FOLLOW-UP → Persistence | [09-screens-communicate.md](09-screens-communicate.md) |
| 18 | LICENSE → Gate | [13-screens-license.md](13-screens-license.md) |
| 19 | NICHE EXPLORER → Discovery | [11-screens-niche-explorer.md](11-screens-niche-explorer.md) |
| 20 | UPDATE & UPGRADE → Renewal | [12-screens-update-upgrade.md](12-screens-update-upgrade.md) |
| — | NERD STATS overlay | [10-global-overlays.md](10-global-overlays.md) |
| — | CTRL+K Command Palette overlay | [10-global-overlays.md](10-global-overlays.md) |
