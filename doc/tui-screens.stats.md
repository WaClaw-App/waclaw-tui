# WaClaw TUI — Screen & Variant Stats

> Sumber: `tui.neuroscienced-customer-journey.md`

---

## Ringkasan

| Metric | Total |
|--------|-------|
| Total Screens | 20 |
| Total States | 110 |
| Total Variants | 20 |
| Total Views | 130 |
| Notification Types | 17 |
| Confirmation Overlays | 4 |

---

## Per-Screen Breakdown

### SCREEN 1: BOOT → FIRST IMPRESSION
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `first_time` | User baru pertama kali buka |
| 2 | state | `returning` | User sudah login & configure (army marching anim) |
| 3 | variant | `returning + ada response baru` | Returning + ada response masuk |
| 4 | variant | `returning + wa disconnect` | Returning + WA putus (per slot) |
| 5 | variant | `returning + config error` | Returning + config broken, partial pause |
| 6 | variant | `returning + lisensi expired` | Lisensi expired, army berhenti |
| 7 | variant | `returning + device conflict` | Lisensi aktif di device lain, hard stop |

**Subtotal: 2 states + 5 variants = 7 views**

---

### SCREEN 2: LOGIN → TRUST (WA ROTATOR)
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `qr_waiting` | Nunggu scan QR, slot info |
| 2 | state | `qr_scanned` | Scan terdeteksi, tanya tambah slot |
| 3 | state | `login_success` | Terhubung, tanya tambah nomor |
| 4 | state | `login_expired` | Session expired, slot aktif tetap jalan |
| 5 | state | `login_failed` | Gagal nyambung, slot lain jalan |

**Subtotal: 5 states + 0 variants = 5 views**

---

### SCREEN 3: NICHE SELECT → IDENTITY
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `niche_list` | Daftar niche, bisa multi-select |
| 2 | state | `niche_multi_selected` | Sudah centang beberapa niche |
| 3 | state | `niche_custom` | Milih custom niche |
| 4 | state | `niche_edit_filters` | Preview filter + granular area |
| 5 | state | `niche_config_error` | YAML parse error, missing fields, line pointers |
| 6 | variant | `niche granular area` | Multi-kota + kecamatan per niche |

**Subtotal: 5 states + 1 variant = 6 views**

---

### SCREEN 4: SCRAPE → ANTICIPATION
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `scraping_active` | Single niche aktif scraping |
| 2 | state | `scraping_multi_active` | Multi niche paralel scraping |
| 3 | state | `scraping_multi_staggered` | Multi niche beda fase |
| 4 | state | `scrape_idle` | Nunggu interval berikutnya |
| 5 | state | `scrape_empty` | Zero results |
| 6 | state | `scrape_error` | Scraper crash / network error |
| 7 | state | `scrape_gmaps_limited` | Google Maps throttle / rate limit |
| 8 | state | `scrape_wa_validation_progress` | Progress cek nomor WA (pre-validation) |
| 9 | variant | `scrape_auto_approved` | Auto-pilot mode, auto-qualify |
| 10 | variant | `scrape_with_wa_validation` | Scrape + WA pre-validation stats |
| 11 | variant | `scrape_high_value_reveal` | Slot machine jackpot lead (skor 9+) |
| 12 | variant | `scrape_batch_complete` | Cascade batch completion |

**Subtotal: 8 states + 4 variants = 12 views**

---

### SCREEN 5: LEAD REVIEW → CURATED
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `reviewing` | Review lead satu-satu (with WA status) |
| 2 | state | `lead_detail` | Detail satu lead (press `d`) |
| 3 | state | `template_preview` | Preview varian (press `1-3`) |
| 4 | state | `queue_complete` | Semua lead sudah di-reviewed |

**Subtotal: 4 states + 0 variants = 4 views**

---

### SCREEN 6: SEND → AUTO-PILOT (WA ROTATOR)
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `sending_active` | Multi-niche + WA rotator + varian rotasi + WA validated |
| 2 | state | `sending_paused` | User pause manual |
| 3 | state | `sending_off_hours` | Di luar jam kerja |
| 4 | state | `sending_rate_limited` | Limit per jam capai |
| 5 | state | `sending_daily_limit` | Limit harian capai |
| 6 | state | `sending_failed` | Pesan gagal (RARE karena WA pre-validation) |
| 7 | state | `sending_all_slots_down` | Semua nomor WA putus |
| 8 | variant | `sending_with_response` | Response masuk saat kirim |

**Subtotal: 7 states + 1 variant = 8 views**

---

### SCREEN 7: MONITOR → COMMAND CENTER
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `live_dashboard` | Multi-niche + ambient data rain + breathing stats |
| 2 | state | `idle_background` | Army kerja, WA rotator info |
| 3 | state | `dashboard_night` | Di luar jam kerja, semua idle |
| 4 | state | `dashboard_error` | WA putus per slot, partial operation |
| 5 | state | `dashboard_empty` | Baru mulai, belum ada data |
| 6 | variant | `dashboard_with_pending_responses` | Ada response belum di-handle |

**Subtotal: 5 states + 1 variant = 6 views**

---

### SCREEN 8: RESPONSE → REWARD
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `response_positive` | Jelas tertarik |
| 2 | state | `response_curious` | Nanya-nanya, interested tapi ragu |
| 3 | state | `response_negative` | Tidak tertarik |
| 4 | state | `response_stop_detected` | Orang bilang stop — auto-add do_not_contact |
| 5 | state | `response_deal_detected` | Closing trigger match — auto-flag deal |
| 6 | state | `response_hot_lead_detected` | Hot lead trigger — auto-prioritize |
| 7 | state | `response_maybe` | Tidak jelas, bisa jadi tertarik |
| 8 | state | `response_auto_reply` | Detected bot/auto-reply |
| 9 | state | `offer_preview` | Preview sebelum kirim offer |
| 10 | state | `response_multi_queue` | 3+ response barengan, triage view |
| 11 | state | `conversion` | DEAL closed — 4-phase full drama |

**Subtotal: 11 states + 0 variants = 11 views**

---

### SCREEN 9: LEADS DATABASE → ARCHIVE
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `leads_list` | Semua lead, grouped by status (incl. follow-up & dingin) |
| 2 | state | `leads_filtered` | Filter by status tertentu |
| 3 | state | `lead_full_detail` | Single lead complete view |
| 4 | variant | `lead_follow_up_due` | Sudah dikontak, belum jawab, waktunya follow-up |
| 5 | variant | `lead_cold` | 2x follow-up belum jawab — lead dingin |
| 6 | variant | `lead_never_contacted` | Baru masuk, belum dikontak |
| 7 | variant | `lead_converted` | Sudah deal — full timeline |

**Subtotal: 3 states + 4 variants = 7 views**

---

### SCREEN 10: TEMPLATE MANAGER → ARMORY
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `template_list` | Daftar varian per niche (ice_breaker/ + follow_up/ + offer/) |
| 2 | state | `template_preview` | Preview varian + placeholder |
| 3 | state | `template_edit_hint` | Redirect ke file editor |
| 4 | state | `template_validation_error` | Broken placeholder / empty / encoding error |

**Subtotal: 4 states + 0 variants = 4 views**

---

### SCREEN 11: WORKERS → PIPELINE VISUALIZER
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `workers_overview` | Live view semua worker |
| 2 | state | `worker_detail` | Deep dive satu worker |
| 3 | state | `worker_add_niche` | Tambah niche baru ke pool |
| 4 | state | `worker_paused` | Worker yang di-pause manual |

**Subtotal: 4 states + 0 variants = 4 views**

---

### SCREEN 12: ANTI-BAN → SHIELD
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `shield_overview` | Semua aman, dynamic ASCII shield + spam guard |
| 2 | state | `shield_warning` | Ada warning, auto-adjust |
| 3 | state | `shield_danger` | Nomor kena flag, auto-pause |
| 4 | state | `shield_slot_detail` | Detail statistik per nomor |
| 5 | state | `shield_settings` | Anti-ban + spam_guard + closing_triggers + follow-up config |

**Subtotal: 5 states + 0 variants = 5 views**

---

### SCREEN 13: SETTINGS → CONFIG REFERENCE
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `settings_overview` | Reference card + WA rotator + granular area |
| 2 | state | `settings_edit` | Buka editor file config |
| 3 | state | `settings_reload` | Setelah edit & reload sukses |
| 4 | state | `settings_reload_error` | Reload gagal, revert ke backup |

**Subtotal: 4 states + 0 variants = 4 views**

---

### SCREEN 14: GUARDRAIL → CONFIG VALIDATION
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `validation_clean` | Semua config valid, green checkmark |
| 2 | state | `validation_errors` | Satu atau lebih config broken, detail per file |
| 3 | state | `validation_warnings` | Valid tapi ada warning (deprecated, unusual, only 1 variant) |
| 4 | state | `validation_fix` | Setelah fix + re-validate |
| 5 | variant | `validation_first_time` | First-time setup validation, more guidance |

**Subtotal: 4 states + 1 variant = 5 views**

---

### SCREEN 15: COMPOSE → VOICE
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `compose_draft` | Text input area buat custom reply (modal) |
| 2 | state | `compose_preview` | Preview pesan sebelum kirim |
| 3 | state | `compose_template_pick` | Quick-pick dari snippets.md |

**Subtotal: 3 states + 0 variants = 3 views**

---

### SCREEN 16: HISTORY → TIMELINE
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `history_today` | Timeline aktivitas hari ini |
| 2 | state | `history_week` | Weekly summary + mini bar charts |
| 3 | state | `history_day_detail` | Detail hari spesifik + timeline |

**Subtotal: 3 states + 0 variants = 3 views**

---

### SCREEN 17: FOLLOW-UP → PERSISTENCE
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `followup_dashboard` | Overview semua follow-up multi-niche |
| 2 | state | `followup_niche_detail` | Detail per niche + varian follow-up |
| 3 | state | `followup_sending` | Follow-up lagi dikirim |
| 4 | state | `followup_empty` | Nggak ada yang perlu follow-up hari ini |
| 5 | state | `followup_cold_list` | Daftar lead dingin (2x follow-up belum jawab) |
| 6 | state | `followup_recontact` | Lead yang pernah respond tapi dingin lagi |

**Subtotal: 6 states + 0 variants = 6 views**

---

### SCREEN 18: LICENSE → GATE
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `license_input` | Input key lisensi (pertama kali) |
| 2 | state | `license_validating` | Cek ke server lisensi |
| 3 | state | `license_valid` | Lisensi valid, gate open |
| 4 | state | `license_invalid` | Key salah / nggak valid |
| 5 | state | `license_expired` | Lisensi expired, hard stop |
| 6 | state | `license_device_conflict` | Lisensi aktif di device lain, hard stop |
| 7 | state | `license_server_error` | Gagal cek lisensi, offline grace 72 jam |

**Subtotal: 7 states + 0 variants = 7 views**

---

### SCREEN 19: NICHE EXPLORER → DISCOVERY
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `explorer_browse` | Browse kategori bisnis populer |
| 2 | state | `explorer_search` | Live search kategori (WA Biz Dir + GMaps) |
| 3 | state | `explorer_category_detail` | Detail kategori + preview config yang bakal di-generate |
| 4 | state | `explorer_generating` | Auto-generate niche.yaml + template files |
| 5 | state | `explorer_generated` | Config berhasil di-generate, siap gas |
| 6 | variant | `explorer dengan area auto-detect` | Area auto-detect dari config.yaml yang udah ada |

**Subtotal: 5 states + 1 variant = 6 views**

---

### SCREEN 20: UPDATE & UPGRADE → RENEWAL
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `update_available` | Ada minor update (gratis, lisensi tetap valid) |
| 2 | state | `update_downloading` | Download update progress |
| 3 | state | `update_ready` | Download selesai, siap restart |
| 4 | state | `upgrade_available` | Ada major upgrade (butuh lisensi baru) |
| 5 | state | `upgrade_license_input` | Input lisensi v2 |
| 6 | state | `license_expired_with_upgrade` | Lisensi v1 expired + ada v2 tersedia |
| 7 | variant | `startup_check` | Background update check saat boot (non-blocking) |

**Subtotal: 6 states + 1 variant = 7 views**

---

## Notification System (Overlay, bukan screen)

| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | notif | `response_masuk` | Ada yang balas pesan |
| 2 | notif | `multi_response` | 3+ responses at once, auto-classified |
| 3 | notif | `scrape_selesai` | Multi-niche batch selesai |
| 4 | notif | `batch_kirim_selesai` | Batch pesan terkirim |
| 5 | notif | `wa_disconnect` | WA putus per slot, partial info |
| 6 | notif | `wa_flag` | Nomor kena flag WA, auto-pause |
| 7 | notif | `health_score_drop` | Health score mendekati threshold |
| 8 | notif | `limit_harian` | Limit harian capai |
| 9 | notif | `streak_milestone` | Achievement unlocked |
| 10 | notif | `config_error` | Config broken, worker paused |
| 11 | notif | `validation_error` | Full validation check gagal |
| 12 | notif | `license_expired` | Lisensi expired, army berhenti |
| 13 | notif | `device_conflict` | Lisensi aktif di device lain, hard stop |
| 14 | notif | `followup_terjadwal` | Follow-up terjadwal hari ini |
| 15 | notif | `lead_dingin` | Lead dingin setelah 2x follow-up |
| 16 | notif | `update_available` | Minor update tersedia (gratis, auto-dismiss 15s) |
| 17 | notif | `upgrade_available` | Major upgrade tersedia (butuh lisensi baru, auto-dismiss 20s) |

**Subtotal: 17 notification types**

---

## Confirmation Overlays (sebelum bulk action)

| # | Nama | Trigger |
|---|------|---------|
| 1 | `bulk_offer` | "auto-kirim offer ke semua" |
| 2 | `bulk_delete` | Hapus banyak leads |
| 3 | `bulk_archive` | Archive banyak leads |
| 4 | `force_device_disconnect` | Putuskan waclaw di device lain |

**Subtotal: 4 confirmation overlays**

---

## Visual Summary

```
SCREEN              STATES  VARIANTS  TOTAL  RATIO
─────────────────────────────────────────────────────────
1  BOOT                2       5       7     ███████░░░
2  LOGIN               5       0       5     █████░░░░░
3  NICHE SELECT        5       1       6     ██████░░░░
4  SCRAPE              8       4      12     ████████████
5  LEAD REVIEW         4       0       4     ████░░░░░░
6  SEND                7       1       8     ████████░░
7  MONITOR             5       1       6     ██████░░░░
8  RESPONSE           11       0      11     ███████████
9  LEADS DB            3       4       7     ███████░░░
10 TEMPLATE            4       0       4     ████░░░░░░
11 WORKERS             4       0       4     ████░░░░░░
12 ANTI-BAN            5       0       5     █████░░░░░
13 SETTINGS            4       0       4     ████░░░░░░
14 GUARDRAIL           4       1       5     █████░░░░░
15 COMPOSE             3       0       3     ███░░░░░░░
16 HISTORY             3       0       3     ███░░░░░░░
17 FOLLOW-UP           6       0       6     ██████░░░░
18 LICENSE             7       0       7     ███████░░░
19 EXPLORER            5       1       6     ██████░░░░
20 UPDATE              6       1       7     ███████░░░
─────────────────────────────────────────────────────────
TOTAL               110      20      130
```

---

## Top 3 Paling Kompleks

1. **SCREEN 4: SCRAPE** — 8 states, 4 variants = 12 views
   - Multi-niche paralel + high-value reveal + batch cascade + gmaps throttle + WA pre-validation = state explosion terbesar
2. **SCREEN 8: RESPONSE** — 11 states, 0 variants = 11 views
   - Paling banyak state: 5 tipe response klasik + 3 closing trigger states (deal/hot_lead/stop) + multi-queue + offer preview + conversion drama
3. **SCREEN 18: LICENSE** — 7 states, 0 variants = 7 views
   - Hard gate system: input → validating → valid/invalid/expired/device_conflict/server_error = semua kondisi lisensi
4. **SCREEN 20: UPDATE & UPGRADE** — 6 states, 1 variant = 7 views
   - Minor update (gratis) + major upgrade (lisensi baru) + download progress + startup check = semua kondisi versi

---

## Apa Yang Baru Dari Versi Sebelumnya

| Fitur Baru | Screen | Impact |
|------------|--------|--------|
| Ctrl+K command palette | Global overlay (bukan screen) | 6 states (cmd_closed/cmd_open/cmd_executing/cmd_empty/cmd_with_recent/cmd_quick_action) buat instant command |
| Command palette in ambient effects | Micro-Interactions | Global overlay search + fuzzy match |
| Ctrl+K key | Keyboard Grammar | Open command palette from anywhere |
| Niche explorer | SCREEN 19 (baru) | 5 states baru buat browse & discovery kategori bisnis |
| Update & upgrade | SCREEN 20 (baru) | 6 states baru buat versi management |
| Nerd stats overlay | Global overlay (bukan screen) | 3 mode (hidden/minimal/expanded) buat system vitals |
| Explorer area auto-detect | SCREEN 19 (variant) | Area auto-detect dari config.yaml |
| Startup check | SCREEN 20 (variant) | Background update check saat boot |
| Update available notif | Notification | Positive/neutral notif baru (auto-dismiss 15s) |
| Upgrade available notif | Notification | Informative notif baru (auto-dismiss 20s) |
| Backtick (`) key | Keyboard Grammar | Toggle nerd stats overlay |
| `u` key | Keyboard Grammar | Cek versi baru |
| Update check in startup | Section 9 | Update check jadi bagian startup sequence (non-blocking) |
| Nerd stats in ambient effects | Micro-Interactions | Global overlay vitals effect |
| Follow-up persistence | SCREEN 17 (baru) | 6 states baru buat auto follow-up |
| License gate | SCREEN 18 (baru) | 7 states baru buat lisensi system |
| License expired variant | SCREEN 1 (boot) | 2 variants baru di boot |
| Device conflict variant | SCREEN 1 (boot) | Variant baru di boot |
| Lead follow-up due | SCREEN 9 (leads db) | Variant baru buat lead yang waktunya follow-up |
| Lead cold/dingin | SCREEN 9 (leads db) | Variant baru buat lead dingin setelah 2x follow-up |
| Lisensi expired notif | Notification | Critical notif baru |
| Device conflict notif | Notification | Critical notif baru |
| Follow-up terjadwal notif | Notification | Neutral notif baru |
| Lead dingin notif | Notification | Informative notif baru |
| Force device disconnect | Confirmation | Overlay baru buat putuskan device lain |
| follow_up/ templates | File system | Folder baru per niche buat follow-up templates |
| license.md | File system | File baru buat simpen lisensi |
| Follow-up in spam_guard | SCREEN 12 (anti-ban) | Config follow-up delay + cold threshold |
| License check in startup | Section 9 | Lisensi check jadi bagian startup sequence |

---

*Generated from tui.neuroscienced-customer-journey.md (WaClaw TUI)*

---

## Global Overlays (bukan screens)

### NERD STATS → VITALS (Toggle Overlay)
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `hidden` | Default — overlay nggak keliatan |
| 2 | state | `minimal` | 1-line footer (CPU/RAM/Goroutines/DB/Uptime) |
| 3 | state | `expanded` | 3-line panel with mini bar charts |

**Subtotal: 3 overlay states (toggle via backtick ` key)**

### CTRL+K → COMMAND PALETTE (Global Overlay)
| # | Tipe | Nama | Deskripsi |
|---|------|------|-----------|
| 1 | state | `cmd_closed` | Default — palette hidden |
| 2 | state | `cmd_open` | Palette terbuka, search aktif, fuzzy filter |
| 3 | state | `cmd_executing` | Command dipilih, lagi eksekusi |
| 4 | state | `cmd_empty` | Search nggak nemu match + suggestions |
| 5 | variant | `cmd_with_recent` | Recently used commands di atas |
| 6 | variant | `cmd_quick_action` | Action langsung execute tanpa pindah screen |

**Subtotal: 4 overlay states + 2 variants = 6 overlay views (toggle via Ctrl+K)**
