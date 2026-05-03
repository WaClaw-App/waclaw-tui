### SCREEN 7: MONITOR → COMMAND CENTER

**Ini screen utama. Home base. Tempat lu ngabupaten.**

**State: live_dashboard** (multi-niche command center)

```

  waclaw                          ● wa nyambung (3 nomor) · 3 niche aktif

  ░░ 3 7 1 4 8 2 9 ░░ 5 1 6 3 8 ░░ 2 7 4 ░░

  ────────────────────────────────────────────────────

  wa rotator

  📱 slot-1  0812-xxxx-3456   ● aktif   4/6 jam
  📱 slot-2  0813-xxxx-7890   ● aktif   3/6 jam
  📱 slot-3  0857-xxxx-2345   ○ cooldown  ready: 3m

  worker pool

  web_developer      ● scraping   24 antri   3 respond
  undangan_digital   ● sending    14 antri   1 respond
  social_media_mgr   ○ idle       0 antri    0 respond

  ────────────────────────────────────────────────────

  hari ini                         minggu ini

  lead nemu           250          total leads     1247
  pesan terkirim       43          pesan terkirim   312
  response              8          response          41
  convert               3          convert           14

  conversion rate: 5.1%         hari terbaik: selasa

  ────────────────────────────────────────────────────

  aktivitas terbaru (semua niche)

  14:23  [web_dev]   kopi nusantara    respond    "iya kak, boleh lihat"
  14:01  [undangan]  wedding bliss     terkirim   ──
  13:47  [web_dev]   gym fortress      sampai     ──
  12:31  [web_dev]   salon cantik      respond    "makasih tapi udah ada"
  11:58  [undangan]  venue gardenia    respond    "berapa harganya?"
  11:15  [web_dev]   toko makmur       terkirim   ──
  09:42  [web_dev]   bengkel jaya      terkirim   ──

  ────────────────────────────────────────────────────

  1  leads    2  pesan    3  worker    4  template    5  anti-ban    6  follow-up    7  pengaturan

  r  refresh    s  scrape semua    q  keluar    `  nerd stats

```

**Ambient data rain effect:** Baris `░░ 3 7 1 4 8 2 9 ░░ 5 1 6 3 8 ░░ 2 7 4 ░░` — angka-angka faint yang scroll perlahan ke bawah di background. Seperti Matrix rain tapi subtle, pake `text_dim` color (hampir ga keliatan tapi bisa dibaca kalau lu cari). Angka-angka ini represent lead IDs yang lagi di-proses. Efeknya: lu NGERAIN ada data yang mengalir, system hidup, tanpa bikin clutter. Update setiap 5 detik, angka baru masuk dari kiri, angka lalu hilang ke kanan.

**Breathing stat numbers:** Angka-angka stats (`250`, `43`, `8`, `3`) punya subtle pulse — opacity naik turun 0.9 → 1.0 → 0.9 dalam 4 detik cycle. Sinkron tapi nggak serentak (offset 200ms per number). Feel: angka-angka ini hidup, bukan screenshot.

**State: idle_background** (army lagi kerja, lu santai)

```

  waclaw                          ● auto-pilot · 3 niche

  ░░ 4 1 7 3 ░░ 6 2 8 ░░ 1 5 9 4 ░░

  ────────────────────────────────────────────────────

  army lu lagi kerja nih. semua auto.

  web_developer      ● scrape 5j 37m · kirim 11m 23s
  undangan_digital   ● scrape 2j 14m · kirim 8m 41s
  social_media_mgr   ● scrape 4j 02m · idle (malam)

  wa                 ● nyambung (3 nomor)
  antrian total      38 pesan
  response hari ini  8

  lu bisa minim window, waclaw tetap kerja.
  3 worker nyari leads + 3 nomor rotasi kirim.
  nanti waclaw notif kalau ada yang penting.

  ↵  liat detail    q  keluar

```

**State: dashboard_night** (di luar jam kerja, semua idle)

```

  waclaw                          ○ mode malam · 3 niche

  ────────────────────────────────────────────────────

  jam kerja: 09:00-17:00 wib
  sekarang: 22:15

  semua sender pause sampai besok pagi.
  scraper tetap jalan di background.
  3 worker auto gas lagi jam 09:00.

  ringkasan hari ini:
  43 terkirim · 8 response · 3 convert

  per niche:
  web_developer     28 terkirim · 5 respond · 2 convert
  undangan_digital  15 terkirim · 3 respond · 1 convert

  tips: selasa jam 10 itu waktu terbaik lu kirim pesan

   ↵  liat detail    q  keluar

```

**State: dashboard_error** (ada masalah)

```

  waclaw                          ✗ ada masalah · 3 niche

  ────────────────────────────────────────────────────

  ✗  wa putus (slot-1) — slot-2 & slot-3 tetap jalan
  ●  scraper tetap jalan (nggak butuh wa)
  ●  database ok

  2 nomor masih aktif. kirim tetap jalan.
  cuma slot-1 yang pending sampai nyambung.
  auto-reconnect nyala.

  1  login ulang    r  cek status    q  keluar

```

**State: dashboard_empty** (baru mulai, belum ada data)

```

  waclaw                                   ● wa nyambung

  ────────────────────────────────────────────────────

  belum ada data. mulai dari sini:

  1  scrape sekarang    2  atur niche    3  kirim pesan

  ────────────────────────────────────────────────────

  tips: mulai dari scrape dulu ya.
  nanti setelah ada lead, sisanya auto.

```

**Variant: dashboard_with_pending_responses** (ada response yang belum di-handle)

```

  waclaw                          ● wa nyambung · 3 niche

  ────────────────────────────────────────────────────

  ⚡ 5 response belum dibalas!

  [web_dev]   kopi nusantara     "iya kak, boleh lihat"
  [web_dev]   salon cantik       "berapa harganya?"
  [undangan]  wedding bliss      "boleh kirim contohnya?"
  [undangan]  venue gardenia     "tertarik, detailnya?"
  [web_dev]   bengkel jaya       "boleh info lebih lanjut?"

  ────────────────────────────────────────────────────

  1  leads    2  pesan    3  worker    4  template    5  anti-ban    6  follow-up    7  pengaturan

  ↵  balas satu-satu    1  auto-kirim offer ke semua
  2  auto per niche     (offer beda per niche)

  ────────────────────────────────────────────────────

  hari ini: 43 terkirim · 8 response · 3 convert

```

Micro-interactions:
- `● wa nyambung` pulses 3s cycle = alive
- `✗ ada masalah` only red thing on screen = auto-draw attention
- `⚡ 3 response belum dibalas!` amber flash = urgency tanpa alarm
- Numbers increment with highlight = movement catches eye
- "respond" rows get warm amber flash 2s = positive signal
- Rejection in neutral color = no negative visual weight
- "lu bisa minim window" = explicitly telling user they can leave
- "1 auto-kirim offer ke semua" = one-key bulk action, auto-pilot mindset
- Ambient data rain: faint scrolling numbers = system alive, data flowing, subliminal activity
- Breathing stats: subtle opacity pulse on numbers = living dashboard, bukan screenshot
- Data rain pauses when lu interact, resumes after 10s idle = focus mode vs ambient mode

---

### SCREEN 8: RESPONSE → REWARD

**Ini interrupt. WaClaw yang panggil lu, bukan sebaliknya.**

**State: response_positive** (jelas tertarik)

```

  💬 ada yang balas!

  ────────────────────────────────────────────────────

  kopi nusantara · cafe · kediri

  "iya kak, boleh lihat desainnya?"

  ────────────────────────────────────────────────────

  ↵  kirim offer    2  balas custom    3  nanti

```

**State: response_curious** (nanya-nanya, interested tapi ragu)

```

  💬 ada yang balas!

  ────────────────────────────────────────────────────

  salon cantik · salon · kediri

  "berapa harganya kak?"

  ────────────────────────────────────────────────────

  ↵  kirim info harga    2  balas custom    3  nanti

```

**State: response_negative** (tidak tertarik)

```

  💬 ada yang balas

  ────────────────────────────────────────────────────

  gym fortress · gym · kediri

  "makasih tapi udah ada yang handle"

  ────────────────────────────────────────────────────

  1  mark invalid    2  follow up nanti    ↵  skip aja

```

**State: response_stop_detected** (orang bilang berhenti — AUTO-ADD ke do_not_contact)

```

  🛑 ada yang minta berhenti

  ────────────────────────────────────────────────────

  bengkel jaya · bengkel · kediri

  "jangan hubungi lagi ya"

  ────────────────────────────────────────────────────

  ⚠  ini match closing_triggers.stop
  waclaw auto-tambah ke do_not_contact.yaml
  nomor ini nggak bakal dikontak lagi.

  ↵  setuju, block    2  block semua niche    3  batal, follow up aja

```

Micro-interactions:
- `🛑` red shield icon = berbeda dari response biasa, ini SERIUS
- Auto-action sudah taken: "waclaw auto-tambah ke do_not_contact" = lu cuma konfirmasi
- `2 block semua niche` = kalau nomor ini ada di multi-niche, block di semua
- `3 batal` = override, tapi warning: ini ngelawan anti-spam guard

**State: response_deal_detected** (CLOSING TRIGGER — auto-detected deal!)

```

  💬 ada yang balas — 🎯 DEAL TERDETEKSI!

  ────────────────────────────────────────────────────

  kopi nusantara · cafe · kediri

  "wah makasi dah transfer kak"

  ────────────────────────────────────────────────────

  ⚡ closing trigger: "dah transfer" → auto-deal

  waclaw deteksi ini sebagai deal.
  tekan ↵ buat konfirmasi, atau s buat koreksi.

  ↵  konfirmasi deal!    s  bukan deal    2  balas dulu

```

Micro-interactions:
- `🎯 DEAL TERDETEKSI!` amber flash + subtle gold shimmer = ini penting, tapi bukan sebesar manual conversion
- Auto-classification: WaClaw udah baca pattern-nya, lu cuma verify
- `s bukan deal` = false positive aman, bisa dikoreksi
- `2 balas dulu` = kadang mau bilang makasih dulu sebelum mark deal

**State: response_hot_lead_detected** (HOT LEAD — closing trigger hot_lead)

```

  💬 ada yang balas — 🔥 HOT LEAD!

  ────────────────────────────────────────────────────

  salon cantik · salon · kediri

  "berapa harganya kak?"

  ────────────────────────────────────────────────────

  ⚡ hot lead trigger: "berapa harga" → auto-prioritize

  waclaw naikin prioritas lead ini.
  tekan ↵ buat kirim offer langsung.

  ↵  kirim offer sekarang    2  balas custom    3  nanti

```

Micro-interactions:
- `🔥 HOT LEAD!` warm amber pulse = ini lead yang paling potensial
- Auto-prioritize: lead ini naik ke atas antrian
- Closing trigger bikin WaClaw proaktif, bukan cuma reactive

**State: response_maybe** (tidak jelas, bisa jadi tertarik)

```

  💬 ada yang balas

  ────────────────────────────────────────────────────

  toko makmur · toko bangunan · kediri

  "oh ya? bisa info lebih lanjut?"

  ────────────────────────────────────────────────────

  ↵  kirim offer    2  kirim info    3  balas custom

```

**State: response_auto_reply** (detected bot/auto-reply)

```

  💬 ada yang balas (auto)

  ────────────────────────────────────────────────────

  bengkel jaya · bengkel · kediri

  "Terima kasih telah menghubungi kami.
   Kami akan membalas pesan Anda secepatnya."

  ────────────────────────────────────────────────────

  ini auto-reply. skip aja ya.

  ↵  skip    1  tetap follow up    q  balik

```

**State: offer_preview** (sebelum kirim offer)

```

  kirim offer ke: kopi nusantara

  ────────────────────────────────────────────────────

  Halo Kak Kopi Nusantara! 👋

  Tadi aku iseng cari Kopi Nusantara di
  Google, ternyata belum ada website
  resminya ya?

  Kebetulan aku Haikal (programmer), udah
  buatin preview desain web yang clean &
  profesional khusus buat Kopi Nusantara.

  Boleh izin kirim fotonya ke sini buat
  diintip? Gratis kok kak lihat-lihat
  aja dulu 😁

  ────────────────────────────────────────────────────

  ↵  kirim    2  ganti template    3  edit dulu

```

**State: response_multi_queue** (multiple responses arrive simultaneously)

```

  💬 3 response masuk barengan!

  ────────────────────────────────────────────────────

  ▸ 01  kopi nusantara
        cafe · kediri
        "iya kak, boleh lihat desainnya?"
        ● positif

    02  wedding bliss WO
        WO · kediri
        "boleh kirim contohnya?"
        ○ penasaran

    03  bengkel jaya
        bengkel · kediri
        "Terima kasih telah menghubungi..."
        ○ auto-reply

  ────────────────────────────────────────────────────

  ↵  proses satu-satu    1  auto-offer semua yang positif
  2  auto per response type

```

Multi-queue micro-interactions:
- Stack appears with 3 items sliding in sequentially 100ms stagger = triage feel
- `● positif` / `○ penasaran` / `○ auto-reply` auto-classified badge = WaClaw udah sort buat lu
- Pre-highlighted `↵ proses satu-satu` = most safe default
- `1 auto-offer semua yang positif` = one-key bulk action, auto-pilot mindset
- Each response type has distinct tint: positive = warm amber, curious = neutral, auto-reply = dimmed
- After processing one, next slides up into focus = queue momentum

**State: conversion** (THE MOMENT — deal closed)

Ini. Ini dia. Puncak. The whole point. Ini screen yang bikin lu ketagihan.

```

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


   ★  ★  ★   D E A L !   ★  ★  ★


  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  kopi nusantara · cafe · kediri

  convert dari ice breaker → offer → deal
  waktu: 2 hari 4 jam

  ────────────────────────────────────────────────────

    ╔═══════════════════════════════════╗
    ║  ✦  ·  ★  ☆  ✧  ·  ★  ✦  ☆  ║
    ║  ☆  ✦  ·  ★  ✧  ☆  ·  ✦  ★  ║
    ║  ★  ☆  ✦  ·  ★  ✧  ☆  ✦  ·  ║
    ║  ✧  ★  ☆  ✦  ·  ★  ☆  ✧  ★  ║
    ╚═══════════════════════════════════╝

  ────────────────────────────────────────────────────

  🏆 conversion ke-3 minggu ini
  🏆 total revenue minggu ini: rp 7.5jt

   ↵  mark as converted

```

**CONVERSION — Full Drama Sequence:**

Phase 1: SHOCK (0-200ms)
- Seluruh screen flash PUTIH full-brightness — semua text hilang sesaat
- Terminal bell `\a` double-tap — dua kali bunyi, cepet
- Screen shock = "berhenti apapun yang lu lakuin"

Phase 2: REVEAL (200-800ms)
- White fades → `★ ★ ★ D E A L ! ★ ★ ★` scales dari 0 ke 1.3 (overshoot) → settle ke 1.0
- Spaced-out letters = monumental, bukan biasa
- Top/bottom `━━━` borders draw inward dari ujung, meet di tengah = framing the moment
- Particle cascade: `✦ · ★ ☆ ✧` scatter dari tengah ke luar, 40 particles, 600ms lifetime
- Particle colors cycle: gold → amber → white → fade
- Color wave: background pulses dari accent → success → normal dalam 400ms sweep

Phase 3: CONTEXT (800ms-1500ms)
- Business name + timeline fade in dari bawah
- `trophy 🏆` icon bounces in dari kanan dengan overshoot
- Revenue number glows gold 3x pulse = "ini uang beneran"
- "conversion ke-3 minggu ini" = streak context, bikin nagih

Phase 4: SETTLE (1500ms+)
- Particles dissolve
- Screen settles ke normal tones
- `↵ mark as converted` fades in = gentle return ke reality
- Hold 800ms sebelum keyboard accepted = lu HARUS ngerasain momen ini dulu

**Sound design suggestion:** Terminal bell `\a` saat shock phase. Kedua bunyi: BING BING. Kalau terminal support bell, ini jadi audio landmark. Kalau nggak, visual shock lebih dari cukup.

**Why this is the most important screen:** Ini dia. Ini satu-satunya alasan WaClaw ada. Semua scraping, template, rotator, anti-ban — semuanya cuma jalan ke screen ini. Kalau conversion screen ga bikin lu seneng, seluruh app gagal. Jadi kita ga hemat. Kita ALL OUT. Full flash. Full particles. Full drama. Karena MOMEN INI LAYAK DIPERINGATI.

**The neuroscience:** Different response types = different animations = different felt experiences. Your brain categorizes them automatically. You don't think "this is a positive response" — you FEEL it because the screen feels warm. Conversion? Conversion feels like WINNING. Karena emang menang.

---

