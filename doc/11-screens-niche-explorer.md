### SCREEN 19: NICHE EXPLORER → DISCOVERY

**Ngetik niche dari nol itu kerja keras. Browse dulu, pilih nanti.**

Screen ini buat user yang bingung harus ngetik apa. Daripada bikin niche.yaml manual, WaClaw nyediain explorer — browse kategori bisnis, cari pakai live search, pilih yang pas, dan WaClaw auto-generate semua file yang dibutuhkan.

**Filosofi niche explorer:**
- User cuma pilih. WaClaw yang nulis.
- Kategori diambil dari 2 sumber: WhatsApp Business Directory (via whatsmeow) + Google Maps categories
- Setelah user pilih kategori, WaClaw auto-generate: niche.yaml, ice_breaker.md (3 varian), queries.md
- Area auto-detect dari config.yaml yang udah ada — nggak perlu input ulang
- Explorer = discovery tool, bukan replacement buat custom niche. User yang udah jago bisa tetap bikin manual.

**State: explorer_browse** (browse kategori)

```

  niche explorer

  bingung mulai dari mana? browse aja.
  pilih kategori, waclaw yang bikin config-nya.

  ── populer ────────────────────────────────

  1  🍜 kuliner             cafe, restoran, catering, bakery
  2  💇 kecantikan          salon, barbershop, spa, nail art
  3  💪 fitness             gym, yoga studio, personal trainer
  4  🏠 properti            agent, developer, interior, renovation
  5  🎓 pendidikan          bimbel, kursus, training, workshop
  6  🏥 kesehatan           klinik, apotek, dokter, therapy
  7  📸 fotografi           wedding, product, studio, event
  8  💻 teknologi           web dev, app dev, it support, digital
  9  🚗 otomotif            bengkel, dealer, rental, sparepart
  10 🎉 event               wedding org, decorator, catering, mc

  ────────────────────────────────────────────

  /  cari kategori    ↑↓  pilih    ↵  liat detail    q  balik

  mau bikin sendiri? tekan 5 di screen niche select.

```

**State: explorer_search** (live search kategori)

```

  niche explorer                               mencari...

  ┌───────────────────────────────────────────────────┐
  │  cafe                                                │
  └───────────────────────────────────────────────────┘

  ── hasil ────────────────────────────────

  1  cafe & coffee shop      12 sub-kategori · 3 area
  2  cafe modern             8 sub-kategori · 2 area
  3  cafe traditional        5 sub-kategori · 1 area
  4  cafe catering           3 sub-kategori · 1 area

  sumber: WhatsApp Business Directory + Google Maps

  ↑↓  pilih    ↵  liat detail    esc  batal

```

Live search: tiap karakter yang lu ketik langsung filter hasil dari WhatsApp Business Directory + Google Maps categories. Debounce 300ms. Nggak ada tombol "cari" — ketik = cari.

**State: explorer_category_detail** (detail kategori sebelum generate)

```

  niche explorer → cafe & coffee shop

  ────────────────────────────────────────────────────

  kategori: cafe & coffee shop
  sub-kategori:
    cafe, coffee shop, kopi nusantara, kafe modern,
    kafe tradisional, coffee bar, espresso bar, roaster

  sumber:
    ● WhatsApp Business Directory  —  247 bisnis terdaftar
    ● Google Maps Categories       —  189 listing aktif

  area (auto-detect dari config):
    kediri (15km), nganjuk (10km), tulungagung (10km),
    blitar (10km), madiun (10km)

  filter default:
    ✗  punya website      (skip yang udah punya)
    ✓  punya instagram    (lebih potensial)
    ⭐ 3.0 - 4.8
    📊 5 - 500 reviews

  template yang bakal di-generate:
    ✓ niche.yaml          (filter + target + area)
    ✓ ice_breaker.md      (3 varian)
    ✓ queries.md          (search queries per area)

  ────────────────────────────────────────────────────

  sudah pas?   ↵  generate config   2  edit dulu   q  balik

  waclaw bakal bikin folder: ~/.waclaw/niches/cafe_coffee_shop/

```

**State: explorer_generating** (lagi generate config files)

```

  niche explorer → cafe & coffee shop

  lagi generate config...

  ●  membuat folder ~/.waclaw/niches/cafe_coffee_shop/
  ●  menulis niche.yaml
  ○  menulis ice_breaker.md (3 varian)
  ○  menulis queries.md

   ━━━━━━━━━━━━━━━━━━░░░░░░  generating...

```

**State: explorer_generated** (config berhasil di-generate)

```

  niche explorer → cafe & coffee shop

  ✓ config berhasil di-generate!

  ────────────────────────────────────────────────────

  ~/.waclaw/niches/cafe_coffee_shop/
  ├── niche.yaml          ✓  5 targets · 5 area · 4 filter
  ├── ice_breaker.md      ✓  3 varian
  └── queries.md          ✓  8 query per area

  ────────────────────────────────────────────────────

  lu bisa edit file-nya langsung kalo mau ubah.
  tekan r buat reload setelah edit.

   ↵  gas pake ini   1  edit config   2  liat template   q  balik

  worker baru bakal jalan paralel sama yang udah ada.

```

**Variant: explorer dengan area auto-detect** (area dari config.yaml yang udah ada)

```

  niche explorer → cafe & coffee shop

  area auto-detect:

  lu udah punya 5 area di config.yaml:
    kediri (15km), nganjuk (10km), tulungagung (10km),
    blitar (10km), madiun (10km)

  niche baru bakal pake area yang sama.
  nggak perlu input ulang. waclaw ambil dari config.

  mau tambah area lain?   1  tambah area   ↵  pake yang ada aja

```

Micro-interactions:
- Category list fade in stagger 80ms = scan-able, kayak menu restoran
- `🍜 💇 💪` emoji per kategori = instant recognition, nggak perlu baca
- Live search: karakter muncul per ketukan = lu tau input lu kebaca
- Search debounce 300ms = nggak spam API, tapi tetap responsive
- Category detail: sub-kategori slide in dari kanan = konteks baru masuk
- Source indicators `● WA ● GMaps` = lu tau datanya dari mana
- Generating: progress bar fills real-time = lu tau ini kerja, bukan hang
- Generated: `✓` checkmarks pulse hijau per file = micro-reward per completion
- Auto-detect area: area list dari config fade in = "waclaw udah tau area lu"
- "worker baru bakal jalan paralel" = explicit promise bahwa niche lama nggak keganggu

**The neuroscience:** Explorer itu mengurangi cognitive load. Daripada lu mikir "target apa ya?" + "filter apa?" + "area mana?" + "template gimana?" — lu cuma pilih 1 kategori. Sisanya WaClaw yang isi. Live search dari WA + GMaps = data nyata, bukan kategori fiktif. Dan auto-generate = zero friction dari "kepengen" ke "jalan". User yang bingung jadi user yang action. User yang udah jago tetap bisa custom.

---

