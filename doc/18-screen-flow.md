
## 7. Complete Screen Flow

```
  BOOT (first time)
    │
    ├─→ LICENSE ── belum ada key? ──→ LICENSE INPUT
    │     │                              │
    │     ▼                              ▼
    ├─→ VALIDATION ── config missing? ──→ CONFIG SETUP
    │     │                                 │
    │     ▼                                 ▼
    │   LOGIN ──── udah pernah? ───→ SKIP   │
    │     │                                 │
    │     ▼                                 │
    │   NICHE SELECT                        │
    │     │                                 │
    │     ▼                                 │
    │   VALIDATION ── niche config ok?      │
    │     │                                 │
    │     ▼                                 │
    │   SCRAPE (auto)                       │
    │     │                                 │
    │     ▼                                 │
    │   WA VALIDATION ── cek nomor punya WA? (background)
    │     │                                 │
    └─────┘─────────────────────────────────┘

  BOOT (returning)
    │
    ▼
  LICENSE CHECK ── expired? ──→ LICENSE EXPIRED (hard stop)
    │                                  │
    │                device conflict? ──→ DEVICE CONFLICT (hard stop)
    │
    ▼
  VALIDATION ── config error? ──→ SHOW ERRORS (partial pause)
    │
    ▼
  MONITOR (home base)
    │
    ├── auto: SCRAPER jalan tiap interval
    │     └── WA VALIDATION otomatis setelah scrape
    ├── auto: SENDER jalan pas jam kerja (hanya WA-validated leads)
    │     ├── VARIAN ROTASI: ice_breaker + offer keduanya rotate
    │     └── FOLLOW-UP auto-kirim buat lead yang belum jawab
    │           └── varian follow-up beda dari ice_breaker
    ├── auto: NOTIFICATION muncul kalau ada event
    │     ├── closing_triggers.deal → auto-flag deal (lu verify)
    │     ├── closing_triggers.hot_lead → auto-prioritize
    │     └── closing_triggers.stop → auto-add do_not_contact (lu verify)
    │
    ├── interrupt: RESPONSE masuk → lu approve
    │     ├── deal detected? → auto-mark deal
    │     ├── stop detected? → auto-block
    │     └── COMPOSE → custom reply
    ├── interrupt: ERROR → lu fix
    │     └── VALIDATION → fix config
    ├── interrupt: MILESTONE → lu seneng
    │
    ├── manual: lu bisa akses DATABASE, TEMPLATE, SETTINGS, FOLLOW-UP kapan aja
    ├── manual: l LICENSE → liat / ganti lisensi
    ├── manual: h HISTORY → liat performa masa lalu
    └── manual: v VALIDATION → force check semua config

  LEAD LIFECYCLE:
  baru → wa_validated → ice_breaker_sent → responded → offer_sent → converted
    │         │                 │            │                            ↑
    │         │                 │            └→ negative → archived         │
    │         │                 │            └→ auto_reply → skipped        │
    │         │                 └→ no response → follow_up_1 → follow_up_2│
    │         │                                    │               │      │
    │         │                                    │               └→ cold │
    │         │                                    └→ responded ────┘──────┘
    │         └→ wa_invalid (skip, nggak dikirim)
    │
    │  RE-CONTACT: responded + dingin → 7 hari jeda → re_contact → responded
    │
    │  FOLLOW-UP LIMIT: max 3 pesan lifetime per lead
    │  ice_breaker + follow_up_1 + follow_up_2 = 3 pesan
    │  setelah 2x follow-up tanpa response = dingin
    │
    │  LICENSE GATE: tanpa lisensi valid = full stop
    │  expired → hard stop, data aman
    │  device conflict → hard stop, bisa force transfer
    │
  RESPONSE CLASSIFICATION:                                              │
  incoming → closing_triggers.deal match? ──→ auto-deal ──────────────┘
           → closing_triggers.hot_lead match? ──→ auto-prioritize
           → closing_triggers.stop match? ──→ auto-add do_not_contact
           → positive/curious/negative/maybe/auto-reply ──→ manual classify
```

**Auto-pilot = default. Manual = optional.**
Tidak ada dead end. Setiap screen punya next step. Setiap loop makin efisien.
Config error = partial pause, bukan full stop. Worker yang ok tetap jalan.
WA validation = bukan optimis, tapi realistis. Nomor yang nggak punya WA = skip, bukan gagal.
Closing triggers = data-driven, bukan tebakan. User define pattern-nya, WaClaw eksekusi.

---
