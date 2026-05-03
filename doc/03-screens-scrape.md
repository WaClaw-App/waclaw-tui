### SCREEN 4: SCRAPE → ANTICIPATION

**State: scraping_active** (single niche view)

```

  lagi nyari...                           niche: web_developer

  target: cafe, gym, salon, toko bangunan
  area: kediri (15km)
  filter: tanpa website

  nemu ━━━━━━━━━━━━━━━━━━━━━━  156

   kopi nusantara          ──  no web  ✓
   gym fortress pro        ──  no web  ✓
   salon cantik alami      ──  has web  ✗ skip
   toko makmur jaya        ──  no web  ✓
   ...scanning


   lolos: 89 leads
   duplikat dibuang: 12
   baru masuk db: 67

   ━━━━━━━━━━━━━━━━━━━━━━━━━━  67 baru

   67 lead baru nunggu. waclaw auto-review ya.
   tekan ↵ buat review sendiri, atau biarin aja.

```

**State: scraping_multi_active** (MULTI NICHE — ini yang bikin nagih)

```

  lagi nyari...                           2 worker aktif

  ────────────────────────────────────────────────────

  web_developer                          ● scraping
    target: cafe, gym, salon, toko
    area: kediri (15km)
    nemu ━━━━━━━━━━━━  156
    lolos: 89 · baru: 67

  undangan_digital                       ● scraping
    target: wedding organizer, venue
    area: kediri + surabaya
    nemu ━━━━━━━━━━━━  94
    lolos: 61 · baru: 48

  ────────────────────────────────────────────────────

  total: 250 nemu · 150 lolos · 115 baru
  semua auto-review & auto-antri.
  lu nggak perlu ngapa-ngapain 👌

  tab  pindah niche    ↵  liat detail    q  balik

```

**State: scraping_multi_staggered** (berbeda fase per niche)

```

  lagi nyari...                           3 worker aktif

  ────────────────────────────────────────────────────

  web_developer                          ● scraping
    nemu: 156 · lolos: 89 · baru: 67

  undangan_digital                       ○ idle (berikutnya: 2j 14m)
    nemu: 94 · lolos: 61 · baru: 48
    batch terakhir: 23 menit lalu

  social_media_mgr                       ● baru mulai
    target: umkm, kafe lokal
    area: malang (10km)
    nemu: 12 · scanning...

  ────────────────────────────────────────────────────

  total aktif: 3 niche · 3 worker
  antrian kirim: 115 pesan

  tiap niche punya timing sendiri.
  lu cuma liat, waclaw yang ngatur semua.

```

**State: scrape_idle** (menunggu interval berikutnya)

```

  nyari

  scrape terakhir: 23 menit lalu
  nemu 67 lead baru
  scrape berikutnya: 5j 37m

  lagi nunggu interval. semua auto.
   ↵  scrape sekarang    q  balik

```

**State: scrape_empty** (zero results)

```

  nyari

  target: cafe, gym, salon, toko bangunan
  area: kediri (15km)

  ━━━━━━━━━━━━━━━━━━━━━━━━━━  0

  kosong nih. bisa jadi:
  - area lu kecil — coba perbesar radius
  - filter terlalu ketat — coba longgarkan
  - query terlalu spesifik — tambah variasi

  1  edit filter    2  ganti area    3  tambah query

```

**State: scrape_error** (scraper crash, network, etc)

```

  nyari

  ✗ scraper error — gagal nyambung google maps

  auto-retry 3 menit lagi.
  kalau masih gagal, waclaw bakal kasih tau.

   1  coba sekarang    q  balik

```

**State: scrape_gmaps_limited** (Google Maps rate limit / throttled)

```

  nyari

  ⏳ google maps throttle — permintaan kebanyakan

  gmaps lagi ngasih jeda. bukan error, cuma antri.
  auto-resume 12 menit lagi.

  scraper yang lain tetap jalan.
  cuma query ke gmaps yang dikasih jeda.

   1  coba sekarang    q  balik

```

**Variant: scrape_auto_approved** (auto-pilot mode, WaClaw auto-qualifies)

```

  lagi nyari...                           2 worker aktif

  ── web_developer ──
  nemu: 156 · lolos filter: 89
  duplikat: 12 · baru: 67
  auto-review: 67 masuk antrian, 8 skip

  ── undangan_digital ──
  nemu: 94 · lolos filter: 61
  duplikat: 3 · baru: 48
  auto-review: 48 masuk antrian, 10 skip

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  total batch ini: 115 lead baru masuk antrian

  pesan bakal dikirim jam kerja nanti.
  tiap niche pakai template sendiri.
  lu nggak perlu ngapa-ngapain 👌

```

**Variant: scrape_with_wa_validation** (setelah scrape, cek nomor WA dulu)

Ini beda sama sebelumnya. Sebelum langsung masuk antrian kirim, WaClaw CEK DULU apakah nomor yang di-scrape punya WhatsApp. Karena realitanya: **banyak nomor Google Maps itu telepon kantor, bukan WA.**

```

  lagi nyari...                           2 worker aktif

  ── web_developer ──
  nemu: 156 · lolos filter: 89
  duplikat: 12 · baru: 67

  cek nomor WA...
  ██████████████░░░░  78% (52/67 cek)

  ✓ punya WA:     34   → masuk antrian kirim
  ✗ bukan WA:     15   → ditandai, nggak dikirim
  ⏳ belum dicek:  18   → masih nunggu

  ── undangan_digital ──
  nemu: 94 · lolos filter: 61
  duplikat: 3 · baru: 48

  cek nomor WA...
  ██████████████████  100% (48/48 cek)

  ✓ punya WA:     29   → masuk antrian kirim
  ✗ bukan WA:     19   → ditandai, nggak dikirim

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  total siap kirim: 63 lead (bukan 115!)
  34 nomor bukan WA — nggak bakal buang slot kirim.

  cek WA jalan di background. nggak ngaretin proses.
  yang udah lolos langsung masuk antrian.

```

**Kenapa WA pre-validation:**
- Google Maps ngasih nomor telepon bisnis, BUKAN nomor WhatsApp
- Realitanya, 30-50% nomor yang di-scrape NGGAK punya WA
- Tanpa pre-validation: antri 115 → cuma 63 yang bisa dikirim → 52 slot kirim terbuang
- Dengan pre-validation: antri 63 → semua bisa dikirim → nggak ada slot terbuang
- Plus: nggak bikin WA flag karena nyoba kirim ke nomor yang bukan WA (ini salah satu trigger ban)
- Method: `check-registration` (cek ke WA server apakah nomor terdaftar) atau `send-silent` (kirim pesan silent type yang nggak keliatan)
- Catatan: pre-validation bisa lambat kalau banyak nomor, jadi jalan di background, lead yang udah validated langsung masuk antrian

**State: scrape_wa_validation_progress** (detail progress validasi WA)

```

  cek nomor WA                             niche: web_developer

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  total nomor dicek: 67
  progress: ██████████████░░░░  78%

  ✓ punya WA           34  (50.7%)
  ✗ bukan WA           15  (22.4%)
  ⏳ belum dicek        18  (26.9%)

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  yang punya WA langsung masuk antrian.
  yang bukan WA ditandai di database.
  nanti kalau mau, bisa dicoba ulang.

  estimasi selesai: 3 menit

  ↵  liat yang punya WA    s  skip tunggu    q  balik

```

Micro-interactions:
- Progress bar fills secara real-time = lu tau ini kerja, bukan hang
- `✓ punya WA` count naik = leads nyata, bukan angka kosong
- `✗ bukan WA` count dimmed = buang, bukan gagal — ini filter yang bikin data lebih bersih
- Percentage shown = data-driven, lu tau hit rate lu

**Variant: scrape_high_value_reveal** (slot machine effect when high-value lead found)

Ketika scraper nemu lead dengan skor 9+ (rating sempurna, no website, aktif IG, banyak reviews):

```

  lagi nyari...                           niche: web_developer

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

   ★  GRAND PALACE HOTEL     skor: 9.5

      hotel · jl. dhoho 45, kediri
      ⭐ 4.9 · 312 reviews · no web · aktif IG

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

   ★ ★ ★   JACKPOT LEAD!   ★ ★ ★

   langsung masuk prioritas antrian.
   template terbaik auto-dipilih.

  ──────────────────────────────────────

   nemu ━━━━━━━━━━━━━━━━━━━━━━  157

   ...scanning

```

Slot machine reveal animation:
1. Business name scrolls like slot reel (3 random names flash sebelum final) = 400ms
2. `★ ★ ★ JACKPOT LEAD! ★ ★ ★` scale up dari 0 dengan triple bounce overshoot = 600ms
3. Gold/amber color wave pulsed dari tengah ke luar = 300ms
4. Terminal bell sound (`\a`) = audio micro-reward
5. Score `9.5` glows amber 2x = perhatian fokus ke angka
6. Auto-settle: setelah 2 detik, jackpot section collapse ke normal list = lanjut kerja

**Variant: scrape_batch_complete** (cascade effect when batch finishes)

```

  lagi nyari...                           batch selesai!

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ✓ batch web_developer selesai

  156 nemu → 89 lolos → 67 baru
  duplikat: 12 · skip: 8

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  67 lead baru masuk antrian auto.
  pesan dikirim jam kerja nanti.

  ✓ batch undangan_digital selesai

  94 nemu → 61 lolos → 48 baru
  duplikat: 3 · skip: 10

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  total batch: 115 lead baru siap kirim
  berikutnya: 5j 37m

```

Batch cascade animation: tiap niche result "falls in" dari atas 200ms stagger, seolah hasil turun satu-per-satu. `✓` checkmark pulse hijau per niche = micro-reward per completion.

Micro-interactions:
- Business names appear one-by-one 30ms delay = slot machine effect
- `✓` marks pulse green briefly = micro-reward per hit
- `✗ skip` dimmed + slide faster = negatives don't linger
- Counter `156` increments live with scale bump = numbers feel tangible
- "waclaw auto-review ya" = explicit statement bahwa lu NGGAK PERLU ngapa-ngapain
- Auto-approved variant: the whole thing feels like a status report, bukan task list
- High-value reveal: slot machine scroll + jackpot bounce + terminal bell = peak excitement in an otherwise calm process
- Batch complete: cascade fall-in + sequential checkmarks = satisfying batch closure feel

**Prinsip utama:** WaClaw bilang "aku yang review, lu santai aja." User bisa interject kalau mau, tapi default = auto. High-value leads dapet momen khusus — lu NGERAIN ada sesuatu yang penting, tanpa perlu baca angka.

---

