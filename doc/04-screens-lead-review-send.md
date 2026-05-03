### SCREEN 5: LEAD REVIEW → CURATED (Optional, Manual Override)

**Screen ini OPTIONAL. WaClaw auto-review by default. Ini cuma muncul kalau user explicitly mau liat.**

**State: reviewing**

```

  review leads                              23 nunggu

  ────────────────────────────────────────────────────

  01  kopi nusantara
      cafe · jl. hasanuddin 23, kediri
      ⭐ 4.2 (87 reviews) · no website
      no instagram
      WA: ✓ terdaftar

      ── kirim ke ini? ──

      ↵  iya, antrikan       1-3  pilih varian
      s  skip                 d  detail
      x  skip & block

  ────────────────────────────────────────────────────

  diantrikan: 12  di-skip: 8  sisa: 23

  ↑↓  pindah   ↵  gas antrikan   q  selesai

```

**State: lead_detail** (press `d` untuk detail)

```

  kopi nusantara
  cafe · jl. hasanuddin 23, kediri

  ⭐ 4.2 · 87 reviews
  no website · no instagram
  ada foto: 12

  google search:
  - tidak ada website resmi
  - tidak ada link di maps
  - tidak ada di sosmed

  skor lead: 8/10 (sangat potensial)

  riwayat: belum pernah dikontak
  follow-up: belum — bisa follow up maks 3x

  ──────────────────────────────────────

  ↵  antrikan    s  skip    1-3  pilih template    q  balik

```

**State: template_preview** (press `1-3` untuk pilih varian)

```

  preview varian: ice_breaker/variant_1

  ──────────────────────────────────────

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

  ──────────────────────────────────────

  varian lain:  2  variant_2    3  variant_3

  ↵  pake ini    1-3  ganti varian    q  batal

```

**State: queue_complete** (semua lead sudah di-review)

```

  review selesai

  diantrikan: 18
  di-skip: 5
  di-block: 0

  pesan bakal dikirim pas jam kerja.
  waclaw yang ngatur timing, lu santai 👌

  ↵  liat dashboard    q  keluar

```

Micro-interactions:
- `↵ iya, antrikan` = one-key decision, pre-highlighted
- Slide-left on queue = decision momentum
- Fade-down on skip = different meaning per animation
- Template preview types out char-by-char = you see YOUR words forming
- "waclaw yang ngatur timing" = explicitly stating auto-pilot

---

### SCREEN 6: SEND → AUTO-PILOT

**Ini BUKAN screen yang lu operasikan. Ini status report.**

**State: sending_active** (jam kerja, multi-niche batch sending, WA rotator)

```

  lagi kirim pesan                    3 WA · 2 niche · 38 diantrikan (wa validated ✓)

  ────────────────────────────────────────────────────

  wa rotator

  📱 0812-xxxx-3456   ● aktif   4/6 jam ini   cooldown: 8m 12s
  📱 0813-xxxx-7890   ● aktif   3/6 jam ini   cooldown: 14m 37s
  📱 0857-xxxx-2345   ○ cooldown  2/6 jam ini   ready: 3m 05s

  ────────────────────────────────────────────────────

  ▸ web_developer (24 antri, semua wa validated)

  01  →  kopi nusantara       📱 slot-1
       ice_breaker: variant_2  ← rotasi!
       ━━━━━━━━━━━━━━━━━ mengirim...

  02     gym fortress pro     📱 slot-2
       ice_breaker: variant_1  ← rotasi!
       nunggu (berikutnya: 11m 23s)

  03     toko makmur jaya     📱 slot-1
       ice_breaker: variant_3  ← rotasi!
       nunggu (berikutnya: 24m 07s)

  ▸ undangan_digital (14 antri, semua wa validated)

  04  →  wedding bliss WO     📱 slot-3
       offer: variant_1
       nunggu (berikutnya: 8m 41s)

  ────────────────────────────────────────────────────

  rate: 9/18 per jam (3 nomor) · hari ini: 12/50
  kirim berikutnya: 08:41 (slot-3)

  semua lead udah dicek WA-nya. nggak ada slot terbuang.
  pesan + varian dirotasi ke 3 nomor. makin banyak varian = makin aman.
  p  pause    ↵  skip tunggu    tab  pindah niche    q  balik

```

**State: sending_paused** (user pause)

```

  lagi kirim pesan                           12 diantrikan

  ────────────────────────────────────────────────────

  ⏸  PAUSE — lu yang stop ini

  01  kopi nusantara       ✓ terkirim
  02  gym fortress pro     nunggu
  03     toko makmur jaya  nunggu

  rate: 4/6 per jam · hari ini: 12/50

  ↵  lanjut lagi    q  balik ke dashboard

```

**State: sending_off_hours** (di luar jam kerja)

```

  lagi kirim pesan

  ⏰ di luar jam kerja (09:00-17:00)
  sekarang: 21:47 wib

  pesan diantri, bakal kirim besok jam 09:00.
  lu nggak perlu ngapa-ngapain.

  1  tetap kirim (emergency)    q  balik

```

**State: sending_rate_limited** (limit jam ini capai)

```

  lagi kirim pesan

  ⏳ limit jam ini capai (6/6)
  bisa kirim lagi: 47 menit lagi

  waclaw ngatur sendiri, lu tinggal tunggu.
  nggak perlu refresh atau ngapa-ngapain.

   ↵  liat dashboard    q  balik

```

**State: sending_daily_limit** (limit harian capai)

```

  lagi kirim pesan

  📊 limit hari ini capai (50/50)
  sisa antrian: 38 pesan → dikirim besok
  (24 dari web_developer · 14 dari undangan_digital)

  hari ini:
  terkirim: 50
  respond: 7
  convert: 2

  bagus! istirahat dulu ya.
  waclaw lanjutin besok pagi auto, semua niche.

   ↵  liat dashboard    q  keluar

```

**State: sending_failed** (pesan gagal, number invalid, dll — RARE karena WA pre-validation)

```

  lagi kirim pesan

  ✗  gym fortress pro — gagal kirim
     nomor tidak terdaftar wa

  ⚠  ini seharusnya nggak terjadi — pre-validation harusnya udah filter.
  kemungkinan: nomor baru saja unreg WA, atau validasi miss.
  waclaw otomatis skip & mark sebagai invalid.
  nggak bakal dikirim ulang.

  1  coba lagi manual    s  skip aja    v  validasi ulang    q  balik

```

**State: sending_all_slots_down** (semua nomor WA putus)

```

  lagi kirim pesan

  ✗  semua nomor wa putus!

  📱 slot-1  ✗ putus
  📱 slot-2  ✗ putus
  📱 slot-3  ✗ putus

  semua sender mati. scraper tetap jalan.
  leads yang udah di-scrape tetap masuk database.
  cuma pengiriman yang pending sampai lu login ulang.

  ──────────────────────────────────────

  1  login ulang    2  login satu per satu    q  balik

```

**Variant: sending_with_response** (ada response masuk saat lagi kirim)

```

  lagi kirim pesan                           12 diantrikan

  ────────────────────────────────────────────────────

  💬 response masuk! kopi nusantara

  "iya kak, boleh lihat desainnya?"

  ↵  kirim offer    2  balas custom    s  nanti

  ────────────────────────────────────────────────────

  01  →  gym fortress pro
       mengirim...

  02     toko makmur jaya
       nunggu (berikutnya: 11m 23s)

```

Micro-interactions:
- `→` arrow hanya di active message = focus indicator
- `✓ terkirim` green pulse = completion reward
- `✓✓ sampai` second pulse different rhythm = escalating reward
- Countdown `11:23` live tick = time feels real
- "semuanya auto. lu cuma perlu nonton." = the whole point
- Response interrupt: gentle 200ms fade-in, not jarring = important, not alarming
- Rate limit bar fills organically → amber near limit = scarcity signal

---

