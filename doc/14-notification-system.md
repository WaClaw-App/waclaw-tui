## 2. NOTIFICATION SYSTEM — The "Poke"

**WaClaw yang nanya, bukan lu yang nyari.**

Notification = interrupt yang gentle. Muncul sebagai overlay di atas screen apapun. Auto-dismiss setelah 10 detik kalau nggak di-respon.

**Notification: response masuk**

```

  ┌─────────────────────────────────────────────────┐
  │  💬 kopi nusantara balas pesan lu               │
  │  [web_dev] "iya kak, boleh lihat desainnya?"    │
  │                                                  │
  │  ↵  balas    s  nanti    q  dismiss              │
  └─────────────────────────────────────────────────┘

```

**Notification: scrape selesai** (multi-niche)

```

  ┌─────────────────────────────────────────────────┐
  │  🔍 2 niche selesai scrape                       │
  │  web_dev: 67 baru · undangan: 48 baru            │
  │  total: 115 lead baru masuk antrian               │
  │                                                  │
  │  ↵  liat    s  biarin aja                        │
  └─────────────────────────────────────────────────┘

```

**Notification: batch kirim selesai**

```

  ┌─────────────────────────────────────────────────┐
  │  📦 batch selesai — 12 pesan terkirim            │
  │  web_dev: 8 · undangan: 4                         │
  │  batch berikutnya: 47 menit lagi                  │
  │                                                  │
  │  ↵  liat    s  ok                                 │
  └─────────────────────────────────────────────────┘

```

**Notification: wa disconnect**

```

  ┌─────────────────────────────────────────────────┐
  │  ✗ wa putus (slot-1)! auto-reconnect nyala...   │
  │  slot-2 & slot-3 tetap jalan                     │
  │                                                  │
  │  1  login ulang    s  dismiss                     │
  └─────────────────────────────────────────────────┘

```

**Notification: wa flag (ban warning)**

```

  ┌─────────────────────────────────────────────────┐
  │  🚨 nomor 0812-xxxx-3456 kena flag!              │
  │  auto-pause slot-1, beban pindah ke slot-2/3     │
  │                                                  │
  │  ↵  liat shield    s  biarin waclaw urus         │
  └─────────────────────────────────────────────────┘

```

**Notification: health score drop**

```

  ┌─────────────────────────────────────────────────┐
  │  📉 health score slot-1 turun ke 55/100          │
  │  masih aman, tapi mendekati threshold (50)       │
  │                                                  │
  │  ↵  liat shield    s  ok                         │
  └─────────────────────────────────────────────────┘

```

**Notification: limit harian**

```

  ┌─────────────────────────────────────────────────┐
  │  📊 limit hari ini capai (50/50)                 │
  │  sisa antrian besok aja ya                        │
  │                                                  │
  │  ↵  ok                                           │
  └─────────────────────────────────────────────────┘

```

**Notification: streak / milestone**

```

  ┌─────────────────────────────────────────────────┐
  │  🔥 10 response minggu ini! conversion naik 2%   │
  │                                                  │
  │  ↵  liat stats    s  nice 👌                      │
  └─────────────────────────────────────────────────┘

```

**Notification: config error** (critical)

```

  ┌─────────────────────────────────────────────────┐
  │  ✗ config error — 1 worker di-pause              │
  │  fotografer/niche.yaml: parse error baris 14     │
  │  worker lain tetap jalan                          │
  │                                                  │
  │  ↵  liat error    v  validasi semua    s  nanti   │
  └─────────────────────────────────────────────────┘

```

**Notification: multi response** (3+ responses at once)

```

  ┌─────────────────────────────────────────────────┐
  │  💬 3 response masuk barengan!                    │
  │  kopi nusantara · wedding bliss · bengkel jaya   │
  │  2 positif · 1 auto-reply                         │
  │                                                  │
  │  ↵  proses    1  auto-offer yang positif    s  nanti│
  └─────────────────────────────────────────────────┘

```

**Notification: validation error**

```

  ┌─────────────────────────────────────────────────┐
  │  ✗ validasi gagal — 3 config error ditemukan     │
  │  2 niche di-pause, 1 tetap jalan                 │
  │                                                  │
  │  ↵  liat error    v  validasi semua    s  nanti   │
  └─────────────────────────────────────────────────┘

```

**Notification: lisensi expired** (critical — army berhenti)

```

  ┌─────────────────────────────────────────────────┐
  │  ✗ lisensi lu udah expired — waclaw berhenti     │
  │  semua worker di-pause. data aman.               │
  │                                                  │
  │  ↵  masukin lisensi baru    q  keluar             │
  └─────────────────────────────────────────────────┘

```

**Notification: device conflict** (critical — army berhenti)

```

  ┌─────────────────────────────────────────────────┐
  │  ✗ lisensi lagi dipakai di PC-KANTOR!            │
  │  waclaw berhenti. 1 lisensi = 1 device.          │
  │                                                  │
  │  ↵  liat lisensi    2  putuskan device lain       │
  │  s  keluar                                        │
  └─────────────────────────────────────────────────┘

```

**Notification: follow-up terjadwal** (gentle reminder)

```

  ┌─────────────────────────────────────────────────┐
  │  📋 14 follow-up terjadwal hari ini               │
  │  waclaw auto-kirim pas jam kerja                  │
  │                                                  │
  │  ↵  liat    s  ok                                 │
  └─────────────────────────────────────────────────┘

```

**Notification: lead dingin** (informative)

```

  ┌─────────────────────────────────────────────────┐
  │  ❄ 4 lead dingin — 2x follow-up tanpa response   │
  │  waclaw nggak bakal kirim lagi otomatis           │
  │                                                  │
  │  ↵  liat    s  ok                                 │
  └─────────────────────────────────────────────────┘

```

**Notification: update available** (positive/neutral — minor update)

```

  ┌─────────────────────────────────────────────────┐
  │  🔄 update tersedia: v1.3.3                       │
  │  bug fix + perf improvement. gratis!              │
  │                                                  │
  │  ↵  update    u  nanti    s  skip                  │
  └─────────────────────────────────────────────────┘

```

**Notification: upgrade available** (informative — major upgrade)

```

  ┌─────────────────────────────────────────────────┐
  │  ⬆️ upgrade tersedia: v2.0.0                      │
  │  major version — butuh lisensi baru               │
  │  v1 lu tetap jalan selama lisensi valid           │
  │                                                  │
  │  ↵  liat info    u  upgrade    s  nanti            │
  └─────────────────────────────────────────────────┘

```

**Confirmation Overlay: bulk action** (sebelum aksi besar/destruktif)

```

  ┌─────────────────────────────────────────────────┐
  │  ⚠️ kirim offer ke 5 leads sekaligus?            │
  │  3 dari web_dev · 2 dari undangan                │
  │                                                  │
  │  ↵  gas kirim    s  batal                         │
  └─────────────────────────────────────────────────┘

```

**Rules:**
- Max 1 notification di layar. Queue the rest.
- Critical (wa disconnect, config error, validation error, lisensi expired, device conflict) = instant, cannot dismiss for 3 seconds
- Positive (response, milestone, update available) = gentle fade-in, auto-dismiss 10s (update available: 15s)
- Neutral (scrape done, follow-up terjadwal) = brief, auto-dismiss 5s
- Informative (lead dingin, upgrade available) = brief, auto-dismiss 7s (upgrade available: 20s)
- Multi-response = slightly longer display (15s) because more data to process
- Confirmation overlay = muncul sebelum bulk action (auto-offer semua, delete semua, force device disconnect). Always `↵ gas` / `↵ go` + `s batal` / `s cancel`. Nggak pernah lebih dari 2 opsi.
- Never stack. Never spam. One at a time.

**Notification type IDs** (English, used in code and REST API):

| ID | Severity | Display (`id`) | Display (`en`) |
|----|----------|----------------|----------------|
| `ResponseReceived` | Positive | `💬 {name} balas pesan lu` | `💬 {name} replied to your message` |
| `MultiResponse` | Positive | `💬 {n} response masuk barengan!` | `💬 {n} responses came in together!` |
| `ScrapeComplete` | Neutral | `🔍 {n} niche selesai scrape` | `🔍 {n} niches finished scraping` |
| `BatchSendComplete` | Neutral | `📦 batch selesai — {n} pesan terkirim` | `📦 batch done — {n} messages sent` |
| `WADisconnect` | Critical | `✗ wa putus (slot-{n})!` | `✗ wa disconnected (slot-{n})!` |
| `WAFlag` | Critical | `🚨 nomor {n} kena flag!` | `🚨 number {n} flagged!` |
| `HealthScoreDrop` | Neutral | `📉 health score slot-{n} turun ke {s}/100` | `📉 slot-{n} health score dropped to {s}/100` |
| `DailyLimit` | Neutral | `📊 limit hari ini capai ({n}/{max})` | `📊 daily limit reached ({n}/{max})` |
| `StreakMilestone` | Positive | `🔥 {n} response minggu ini!` | `🔥 {n} responses this week!` |
| `ConfigError` | Critical | `✗ config error — {n} worker di-pause` | `✗ config error — {n} workers paused` |
| `ValidationError` | Critical | `✗ validasi gagal — {n} config error` | `✗ validation failed — {n} config errors` |
| `LicenseExpired` | Critical | `✗ lisensi expired — waclaw berhenti` | `✗ license expired — waclaw stopped` |
| `DeviceConflict` | Critical | `✗ lisensi lagi dipakai di {device}!` | `✗ license in use on {device}!` |
| `FollowUpScheduled` | Informative | `📋 {n} follow-up terjadwal hari ini` | `📋 {n} follow-ups scheduled today` |
| `LeadCold` | Informative | `❄ {n} lead dingin` | `❄ {n} cold leads` |
| `UpdateAvailable` | Positive | `🔄 update tersedia: {v}` | `🔄 update available: {v}` |
| `UpgradeAvailable` | Informative | `⬆️ upgrade tersedia: {v}` | `⬆️ upgrade available: {v}` |

---
