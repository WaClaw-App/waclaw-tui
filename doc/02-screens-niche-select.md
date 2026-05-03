
### SCREEN 3: NICHE SELECT → IDENTITY

**Ini screen baru. User milih siapa dia dan siapa targetnya.**

**State: niche_list** (pilih niche — BISA MULTI)

```

  pilih niche lu

  niche = siapa lu + siapa target lu
  bisa pilih lebih dari satu! tiap niche = 1 worker.
  makin banyak niche, makin luas jaring lu.

  1  ☐ web developer        buat yang jual jasa bikin web
  2  ☐ undangan digital     buat yang jual undangan digital
  3  ☐ social media mgr     buat yang jasa kelola sosmed
  4  ☐ fotografer           buat yang jasa foto & portfolio
  5  ☐ custom               bikin niche sendiri dari file

  space  centang/hapus   ↵  gas dengan yang dicentang
  semua niche jalan paralel, masing-masing punya worker.

```

**State: niche_multi_selected** (sudah centang beberapa)

```

  pilih niche lu

  1  ☑ web developer        buat yang jual jasa bikin web
  2  ☑ undangan digital     buat yang jual undangan digital
  3  ☐ social media mgr     buat yang jasa kelola sosmed
  4  ☐ fotografer           buat yang jasa foto & portfolio
  5  ☐ custom               bikin niche sendiri dari file

  2 niche dipilih:
  ▸ web_developer — kediri, 15km — 3 template
  ▸ undangan_digital — kediri + surabaya — 2 template

  ↵  gas jalanin 2 niche    space  ubah    q  balik

  kedua niche bakal jalan paralel.
  masing-masing scrape sendiri, kirim sendiri.

```

**State: niche_custom** (milih custom)

```

  niche custom

  taruh file niche lu di:
  ~/.waclaw/niches/nama_niche/

  butuh minimal:
  - niche.yaml    (filter & target)
  - ice_breaker.md

  contoh niche.yaml udah ada di:
  ~/.waclaw/niches/_contoh/

  kalau udah siap, tekan r buat reload.

   r  reload    1-5  pilih yang ada    q  balik

```

**State: niche_edit_filters** (setelah pilih niche, preview filters)

```

  niche: web developer

  target:
    cafe, gym, salon, toko bangunan, bengkel

  filter:
    ✗  punya website      (skip yang udah punya)
    ✓  punya instagram    (lebih potensial)
    ⭐ 3.0 - 4.8          (rating range)
    📊 5 - 500 reviews    (ukuran bisnis)

  area (5 kota):
    kediri       15km  ──  4 kecamatan
    nganjuk      10km  ──  2 kecamatan
    tulungagung  10km  ──  3 kecamatan
    blitar       10km  ──  2 kecamatan
    madiun       10km  ──  1 kecamatan

  ──────────────────────────────────────

  sudah pas?   ↵  gas scrape   2  edit filter   q  balik

```

**Variant: niche dengan granular area config**

```

  niche: undangan digital

  target:
    wedding organizer, venue pernikahan, gedung, jasa dekorasi

  area (8 kota, 14 kecamatan):
    kediri         15km  ──  kb. kediri, mojoroto, kotabaru
    nganjuk        10km  ──  kb. nganjuk, sukomoro
    tulungagung    10km  ──  kb. tulungagung, boyolangu, kedungwaru
    blitar         10km  ──  kb. blitar, kanigoro
    madiun         10km  ──  kb. madiun
    surabaya       10km  ──  genteng, gubeng, rungkut
    sidoarjo       10km  ──  sidoarjo, waru
    gresik         10km  ──  gresik, cerme

  filter:
    ✓  punya instagram    (justru yang ada IG = potensial)
    ✗  punya website      (bebas, bukan wajib)
    ⭐ 3.0 - 4.8
    📊 5 - 300

  makin banyak area = makin banyak leads.
  tiap area di-scrape paralel per worker.

```

**State: niche_config_error** (config yaml broken atau missing fields)

```

  ✗ niche error: fotografer

  ────────────────────────────────────────────────────

  ~/.waclaw/niches/fotografer/niche.yaml

  3 masalah:

  ✗  baris 14: parse error
     │  targets:
     │    - "wedding
     │          ^ kurang tanda kutip penutup

  ✗  field wajib kosong: areas
     niche.yaml harus punya minimal 1 area

  ✗  field scoring: rating_min harus angka, bukan "rendah"

  ────────────────────────────────────────────────────

  worker fotografer di-pause sampai config diperbaiki.
  worker lain tetap jalan normal.

  1  buka file    2  liat contoh config    r  reload    q  balik

```

Micro-interactions:
- Niche list items fade in stagger = scan-able
- `↵ buat pilih yang direkomendasiin` pre-highlighted = path of least resistance
- Filter preview shows ✗/✓ symbols = instant comprehension
- On `↵ gas scrape`: smooth transition ke scrape screen = momentum preserved
- `niche_config_error`: error lines highlight dengan merah blink = langsung keliatan mana yang salah
- "1 buka file" shortcut = buka $EDITOR langsung ke baris error = zero friction fix
- "worker lain tetap jalan" = explicit message bahwa partial failure ok

---

