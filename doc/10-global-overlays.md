### NERD STATS → VITALS (Toggle Overlay)

**Bukan screen. Overlay global. Tekan backtick (`) buat toggle.**

Ini bukan screen baru — ini overlay yang bisa muncul di atas screen APAPUN. Tujuannya: kasih lu info system vitals tanpa harus pindah screen. Default = hidden. Tekan ` sekali = muncul. Tekan ` lagi = ilang. Simple.

**Kenapa overlay, bukan screen?** Karena system vitals itu info yang lu butuh SESUDAH, bukan info yang lu cari. Lu lagi liat monitor, tiba-tiba penasaran "RAM lu berapa sekarang?" — tekan `, liat, tekan ` lagi, balik kerja. Nggak perlu pindah screen, nggak perlu kehilangan context.

**Dua mode:**

**Mode: minimal** (1-line footer, default saat toggle on)

```

  ── CPU 12% · RAM 134MB · Goroutines 23 · DB 2.4MB · Uptime 4j 12m ──

```

Minimal mode: cuma 1 baris di paling bawah layar. Warna `text_dim` (hampir ga keliatan tapi bisa dibaca). Nggak ganggu konten utama. Cuma info cepat.

**Mode: expanded** (5-line panel, tekan ` lagi dari minimal)

```

  ── nerd stats ──────────────────────────────────────

    CPU         12.3%  ████░░░░░░░░░░░░░░░░
    RAM         134MB  ██████░░░░░░░░░░░░░░  / 512MB
    Goroutines  23     ██░░░░░░░░░░░░░░░░░░  / 100
    DB Size     2.4MB  █░░░░░░░░░░░░░░░░░░░  / 50MB
    Uptime      4j 12m
    Version     v1.3.2

  ── logs ────────────────────────────────────────────

    14:23:01 scrape.web_dev  67 leads · 12 qualified
    14:22:48 send.slot_2     ✓ 8 terkirim · next 47m
    14:22:15 antiban         health 72 → 68 (warning)
    14:21:30 followup        3 due → auto-queue
    14:20:11 config          hot reload ✓ applied

  ────────────────────────────────────────────────────

```

Expanded mode: 5-line panel di bawah layar. Mini bar chart per metric + **live system logs stream**. Warna tetap `text_dim` kecuali kalau ada yang mendekati limit (RAM > 80% = `warning` amber, goroutine > 80 = `danger` merah).

**Live System Logs Stream:**
- 5 baris terakhir dari system event log, scroll otomatis (newest on top)
- Warna per log level: info = `text_dim`, warning = `warning` amber, error = `danger` merah, success = `success` green
- Format: `HH:MM:SS source          message`
- Source tags: `scrape.{niche}`, `send.slot_{n}`, `antiban`, `followup`, `config`, `license`, `worker.{niche}`, `wa.slot_{n}`
- Update real-time — log baru push yang lama ke bawah (tail -f style)
- Max 5 baris ditampilin. Sisa nya di full log file (`~/.waclaw/logs/`)
- Tekan `↵` saat di expanded mode buat buka full log di `$PAGER` (less)
- Logs stream pause kalau user scroll up (resume otomatis 5 detik setelah stop scroll)

**Metrics yang ditampilin:**
- **CPU%** — usage proses WaClaw saat ini
- **RAM** — memory dipakai / total
- **Goroutines** — jumlah goroutine aktif / max configured
- **DB Size** — ukuran database leads / max
- **Uptime** — berapa lama WaClaw jalan sejak startup
- **Version** — versi WaClaw yang lagi jalan (hanya di expanded mode)

**Toggle behavior:**
- Tekan ` pertama: hidden → minimal
- Tekan ` kedua: minimal → expanded
- Tekan ` ketiga: expanded → hidden
- Kalau 30 detik nggak ditekan: auto-collapse ke hidden (kembali ke mode default)
- Bisa di-toggle dari screen APAPUN, termasuk saat lagi di compose, review, atau monitor

**Data source:** Semua metrics diambil dari Go runtime (`runtime.ReadMemStats`, `runtime.NumGoroutine`) dan SQLite DB stats. Zero network call. Zero external dependency. Ini beneran vitals system lu, bukan data dari server.

Micro-interactions:
- Toggle on: minimal footer slide up dari bawah 150ms = nggak ngejut
- Toggle expanded: panel expand dari 1 baris ke 5 baris dengan height morph 200ms
- Toggle off: collapse ke hidden 150ms
- Bar chart fills gradient sweep 300ms = progress feel
- RAM > 80%: bar morph ke amber + subtle pulse = "keep an eye on this"
- Goroutine > 80: bar morph ke merah + double flash = "something might be wrong"
- Auto-collapse: gentle fade out 500ms = nggak abrupt
- Metrics update setiap 2 detik = live, bukan snapshot
- **Log entry new:** slide in dari kanan 100ms + `text_dim` → appropriate color = live feel
- **Log level color:** error entries glow `danger` red briefly = attention without interrupt
- **Log overflow:** oldest line fade out 150ms saat new entry push = clean transition

**The neuroscience:** Nerd stats itu kayak detak jantung. Lu nggak perlu liat terus — tapi kalau penasaran, 1 tombol. Overlay, bukan screen, karena vitals itu context-aware. Lu lagi liat leads, bukan lagi liat system. Tapi kalau system minta perhatian (RAM tinggi), warnanya yang ngomong, bukan notifikasi.

---

### CTRL+K → COMMAND PALETTE (Global Overlay)

**Bukan screen. Bukan menu. Ini teleport.**

Command palette itu cara tercepat buat ngapa-ngapain di WaClaw. Tekan `Ctrl+K` dari mana aja → ketik apa yang lu mau → ↵ → eksekusi. Nggak perlu hafal shortcut. Nggak perlu navigate 3 layer deep. Pikiran → ketik → jadi. 3 detik max.

**Kenapa command palette, bukan menu?** Menu itu hierarchy — lu harus tau dulu menu apa, terus submenu apa, terus item mana. Command palette itu FLAT — semua command di 1 level, lu cuma ketik yang lu mau. VS Code, Notion, Linear — semua pakai command palette karena itu cara tercepat buat power user. Dan WaClaw user itu power user.

**Kenapa Ctrl+K?** Karena `K` = "komando". Dan Ctrl+K udah jadi standard de facto buat command palette di terminal apps (kayak fzf, zoxide). Lagian, lu nggak butuh Ctrl+K buat apapun lain di WaClaw — semua action cukup 1 tombol tanpa modifier. Ctrl+K itu satu-satunya chord key di seluruh app = gampang ingat.

**State: cmd_closed** (default — palette nggak keliatan)

Palette hidden. Screen sekarang normal. Tekan Ctrl+K → buka palette.

**State: cmd_open** (palette terbuka, search aktif)

```

  ┌─────────────────────────────────────────────────────────────┐
  │  > scrape                                              │  ×  │
  ├─────────────────────────────────────────────────────────────┤
  │                                                             │
  │  ── recently used ──────────────────────────────────────    │
  │                                                             │
  │  ▸ Scrape semua niche sekarang            s  ·  scrape     │
  │    Force scrape satu worker               1  ·  scrape     │
  │                                                             │
  │  ── commands ──────────────────────────────────────────    │
  │                                                             │
  │    Pause semua worker                     p  ·  workers    │
  │    Pause satu worker                      2  ·  workers    │
  │    Lihat leads database                   1  ·  database  │
  │    Lihat dashboard                        d  ·  monitor   │
  │    Edit config                            e  ·  settings  │
  │    Validate semua config                  v  ·  guardrail │
  │    Lihat anti-ban shield                  5  ·  shield    │
  │    Cek update                             u  ·  update    │
  │    Lihat lisensi                          l  ·  license   │
  │    Lihat riwayat                          h  ·  history   │
  │    Explore niche baru                     n  ·  explorer  │
  │    Reload config                          r  ·  settings  │
  │    Kirim follow-up semua                  a  ·  followup  │
  │    Toggle nerd stats                      `  ·  overlay   │
  │                                                             │
  └─────────────────────────────────────────────────────────────┘

  ↑↓  pilih    ↵  eksekusi    esc  tutup    Ctrl+K  tutup

```

**Anatomy of command palette:**

1. **Search bar** (`> scrape`) — ketik apapun, filter langsung real-time. Debounce 50ms (instant). Match against command name + description + shortcut + category.

2. **Recently used** — 3 command terakhir yang lu pakai. Otomatis muncul di atas. Kalau lu sering "scrape semua", itu selalu di paling atas. Muscle memory + recency = fastest path.

3. **Filtered results** — semua command yang match search query. Diurutin berdasarkan: exact match > name match > description match > category match. Masing-masing baris punya: nama command, shortcut asli (kalau ada), dan kategori.

4. **Category tags** — setiap command punya category tag (`·  scrape`, `·  workers`, `·  database`). Ini buat filter juga — ketik "worker" muncul semua yang berkaitan. Ketik "shield" langsung ke anti-ban.

**State: cmd_executing** (command dipilih, lagi eksekusi)

```

  ┌─────────────────────────────────────────────────────────────┐
  │  > scrape                                              │  ×  │
  ├─────────────────────────────────────────────────────────────┤
  │                                                             │
  │  ▸ Scrape semua niche sekarang            s  ·  scrape     │
  │                                                             │
  │    ● executing...                                           │
  │                                                             │
  └─────────────────────────────────────────────────────────────┘

```

Execute animation: selected row pulse accent → palette collapses → command runs. Total transition 300ms. Nggak ada loading screen — palette tutup, screen baru muncul / action jalan.

**State: cmd_empty** (ketik sesuatu, nggak ada yang match)

```

  ┌─────────────────────────────────────────────────────────────┐
  │  > xyz                                                 │  ×  │
  ├─────────────────────────────────────────────────────────────┤
  │                                                             │
  │  nggak nemu command "xyz"                                   │
  │                                                             │
  │  coba: scrape, pause, leads, config,                       │
  │        shield, follow-up, update, license                   │
  │                                                             │
  └─────────────────────────────────────────────────────────────┘

```

Empty state: helpful suggestions, bukan cuma "no results". Lu tau command apa yang ada.

**Full Command Registry:**

```
── NAVIGATION ──────────────────────────────────────────────
  Dashboard / Monitor           d      ·  monitor
  Leads Database                1      ·  database
  Send Queue                    2      ·  send
  Workers                       3      ·  workers
  Templates                     4      ·  template
  Anti-Ban Shield               5      ·  shield
  Follow-Up                     6      ·  followup
  Settings                      7      ·  settings
  History                       h      ·  history
  Niche Explorer                n      ·  explorer
  License                       l      ·  license
  Update / Upgrade              u      ·  update

── ACTIONS ─────────────────────────────────────────────────
  Scrape semua niche sekarang   s      ·  scrape
  Scrape satu worker            —      ·  scrape
  Pause semua worker            p      ·  workers
  Resume semua worker           —      ·  workers
  Kirim follow-up semua         a      ·  followup
  Edit config                   e      ·  settings
  Reload config                 r      ·  settings
  Validate semua config         v      ·  guardrail
  Cek versi baru                u      ·  update
  Toggle nerd stats             `      ·  overlay
  Logout WA semua slot          —      ·  login
  Force scrape retry            —      ·  scrape

── OVERLAYS ────────────────────────────────────────────────
  Toggle nerd stats             `      ·  overlay
  Command palette               Ctrl+K ·  overlay
  Shortcut help                 ?      ·  overlay

── SPECIAL ─────────────────────────────────────────────────
  Compose custom reply          —      ·  compose
  Search leads                  /      ·  database
  Export leads CSV              —      ·  database
  Mark lead converted           —      ·  response
  Block lead                    —      ·  response
  Re-contact cold lead          —      ·  followup
  Switch language / Ganti bahasa —     ·  overlay
```

**Key behaviors:**

1. **Fuzzy search.** Ketik "scrp" → match "scrape". Ketik "pld" → match "pause all worker". Nggak perlu exact match — fuzzy finder yang cerdas. Algo: fzf-style scoring (exact > prefix > substring > fuzzy).

2. **Context-aware.** Command yang lagi relevant di screen sekarang muncul lebih tinggi. Kalau lu lagi di SCRAPE screen, "pause worker" dan "force retry" muncul di atas. Kalau di MONITOR, "scrape semua" dan "leads database" yang muncul duluan.

3. **Recently used priority.** 3 command terakhir selalu di bagian atas, di atas semua hasil search. Karena 80% waktu lu ngulang action yang sama.

4. **No dead ends.** Setiap command punya hasil. Navigasi = pindah screen. Action = execute. Toggle = toggle. Nggak ada command yang "coming soon" atau disabled.

5. **Esc / Ctrl+K = close.** Dua cara buat tutup. Sama kayak cara buka. Lu bisa Ctrl+K → liat → Ctrl+K → tutup. Atau Esc. Whichever muscle memory lu punya.

6. **Ctrl+C = hard cancel.** Kalau palette terbuka dan lu beneran mau keluar WaClaw, Ctrl+C dari palette = close palette aja (bukan exit app). Dari screen manapun tanpa palette, Ctrl+C = keluar.

**Variant: cmd_with_recent** (ada recently used commands)

```

  ┌─────────────────────────────────────────────────────────────┐
  │  >                                                    │  ×  │
  ├─────────────────────────────────────────────────────────────┤
  │                                                             │
  │  ── barusan ───────────────────────────────────────────     │
  │                                                             │
  │  ▸ Scrape semua niche sekarang            s  ·  scrape     │
  │    Lihat leads database                   1  ·  database  │
  │    Toggle nerd stats                      `  ·  overlay   │
  │                                                             │
  │  ── semua command ─────────────────────────────────────     │
  │                                                             │
  │    Dashboard / Monitor                    d  ·  monitor    │
  │    Pause semua worker                     p  ·  workers    │
  │    Edit config                            e  ·  settings   │
  │    Validate config                        v  ·  guardrail  │
  │    Lihat shield                           5  ·  shield     │
  │    ...                                                     │
  │                                                             │
  └─────────────────────────────────────────────────────────────┘

```

Kosong (belum ngetik) = recently used + semua command. Ngetik = filter. Always helpful.

**Variant: cmd_quick_action** (command yang langsung execute tanpa pindah screen)

Beberapa command itu ACTIONS, bukan navigasi. Contoh: "pause semua worker", "scrape sekarang", "reload config". Command ini langsung execute — palette tutup, toast notification muncul 2 detik konfirmasi, lu tetap di screen yang sama.

```

  ┌─────────────────────────────────────────────────────────────┐
  │  > pause                                              │  ×  │
  ├─────────────────────────────────────────────────────────────┤
  │                                                             │
  │  ▸ Pause semua worker                     p  ·  workers    │
  │    ⚡ langsung eksekusi, nggak pindah screen               │
  │                                                             │
  │    Pause satu worker                      2  ·  workers    │
  │    Resume semua worker                    —  ·  workers    │
  │                                                             │
  └─────────────────────────────────────────────────────────────┘

```

Quick action indicator: `⚡ langsung eksekusi` = lu tau ini command langsung jalan, bukan navigasi. Nggak ada surprise.

Micro-interactions:
- Palette open: slide down from top 150ms + backdrop dim = "focus here now"
- Search typing: results filter real-time zero lag = instant response
- Selection highlight: accent color slide 50ms = snappy, not sluggish
- Execute: selected row pulse accent → palette collapse → screen/action = 300ms total
- Recently used: badge `barusan` with subtle amber tint = these are YOUR commands
- Quick action: `⚡` icon with brief glow = "this runs immediately"
- Empty state: suggestion chips clickable = zero friction recovery
- Close: palette slides up 100ms + backdrop undim = smooth return to context
- Fuzzy match: matched characters get `accent` color highlight = lu tau kenapa ini match
- Category tags: `·  scrape` dimmed = metadata, bukan noise

**The neuroscience:** Command palette itu working memory externalizer. Otak lu tau mau ngapa — "mau pause workers" — tapi lu nggak perlu ingat tombolnya apa. Lu cuma ketik apa yang lu pikirin, dan WaClaw ngerti. Itu kenapa VS Code, Notion, dan semua app modern pakai command palette: karena **natural language > memorized shortcuts**. Tapi buat yang udah hafal shortcut, shortcut tetap jalan — command palette itu ALTERNATIF, bukan pengganti. Lu bisa pake `s` langsung, ATAU Ctrl+K → "scrape". Pilihan lu.

---

