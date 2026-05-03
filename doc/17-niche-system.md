
## 6. The Niche System — File-Based Power

**Setiap niche = folder. Setiap template = file. Zero UI config.**

```
~/.waclaw/
├── config.yaml              # setting utama + anti_ban + spam_guard
├── config.yaml.bak          # auto-backup (setiap reload sukses)
├── theme.yaml               # warna & feel
├── license.md               # lisensi waclaw (key + device + expires)
├── queries.md               # query pencarian
├── snippets.md              # template snippet buat custom reply
├── do_not_contact.yaml      # auto-block list (orang yang bilang stop)
├── wa_slots/
│   ├── slot_1.yaml          # nomor WA #1 (auto-generated)
│   ├── slot_2.yaml          # nomor WA #2
│   └── slot_3.yaml          # nomor WA #3
└── niches/
    ├── _contoh/
    │   ├── niche.yaml       # contoh config buat referensi
    │   ├── ice_breaker/
    │   │   ├── variant_1.md  # contoh ice breaker varian 1
    │   │   └── variant_2.md  # contoh ice breaker varian 2
    │   └── offer/
    │       ├── variant_1.md  # contoh offer varian 1
    │       └── variant_2.md  # contoh offer varian 2
    ├── web_developer/
    │   ├── niche.yaml       # filter & target + areas + closing_triggers
    │   ├── ice_breaker/     # ROTATABLE — 1 file = 1 varian
    │   │   ├── variant_1.md # "tadi iseng cari"
    │   │   ├── variant_2.md # "baru lihat websitenya"
    │   │   └── variant_3.md # "kebetulan lewat"
    │   ├── follow_up/      # ROTATABLE — follow-up templates
    │   │   ├── follow_up_1.md  # "halo kak, cuma ngingetin"
    │   │   ├── follow_up_2.md  # "kak, penawaran terbatas nih"
    │   │   └── follow_up_3.md  # "terakhir kak" (manual only)
    │   └── offer/           # ROTATABLE — 1 file = 1 varian
    │       ├── variant_1.md # direct-curiosity
    │       ├── variant_2.md # pattern-interrupt
    │       └── variant_3.md # admin-bypass
    ├── undangan_digital/
    │   ├── niche.yaml
    │   ├── ice_breaker/
    │   │   ├── variant_1.md
    │   │   └── variant_2.md
    │   └── offer/
    │       └── variant_1.md
    └── social_media_manager/
        ├── niche.yaml
        ├── ice_breaker/
        │   └── variant_1.md
        └── offer/
            └── variant_1.md
```

**BREAKING CHANGE dari spec sebelumnya:**
- `ice_breaker.md` → `ice_breaker/` folder dengan `variant_*.md` files
- `offer_*.md` → `offer/` folder dengan `variant_*.md` files
- Kenapa? Karena 1 file = 1 varian = WA bisa deteksi pola. Banyak varian = rotasi = aman.
- Validasi: minimal 1 varian per `ice_breaker/` dan per `offer/`. Kalau cuma 1, warning "tambah varian lagi biar aman dari ban".
- Rotasi: per kirim, WaClaw pick varian secara round-robin atau random (sesuai `template_rotation_mode` di config).

### niche.yaml — Full Filter Power

```yaml
# ~/.waclaw/niches/undangan_digital/niche.yaml

name: "Undangan Digital"
description: "target wedding organizer & venue"

targets:
  - "wedding organizer"
  - "venue pernikahan"
  - "gedung pertemuan"
  - "jasa dekorasi"

areas:
  - city: "kediri"
    radius: 15
    kecamatan:
      - "kota kediri"
      - "mojoroto"
      - "kotabaru"
      - "ngasem"
  - city: "nganjuk"
    radius: 10
    kecamatan:
      - "kota nganjuk"
      - "sukomoro"
  - city: "tulungagung"
    radius: 10
    kecamatan:
      - "kota tulungagung"
      - "boyolangu"
      - "kedungwaru"
  - city: "blitar"
    radius: 10
    kecamatan:
      - "kota blitar"
      - "kanigoro"
  - city: "madiun"
    radius: 10
    kecamatan:
      - "kota madiun"
  - city: "surabaya"
    radius: 10
    kecamatan:
      - "genteng"
      - "gubeng"
      - "rungkut"
  - city: "sidoarjo"
    radius: 10
    kecamatan:
      - "sidoarjo"
      - "waru"
  - city: "gresik"
    radius: 10
    kecamatan:
      - "gresik"
      - "cerme"

filters:
  has_website: any          # any | true | false
  has_instagram: true       # prefer yang PUNYA instagram
  rating_max: 4.8           # di bawah ini = masih butuh bantuan
  rating_min: 3.0           # terlalu rendah = mungkin nggak worth it
  review_count_max: 500     # bisnis kecil-menengah
  review_count_min: 5       # harus punya sedikit presence
  has_photos: true          # bisnis yang ada foto = aktif
  categories_exclude:
    - "gereja"
    - "masjid"

# CLOSING TRIGGERS — teks yang otomatis dianggap DEAL
# WaClaw scan response masuk, kalau match salah satu pattern = auto-flag sebagai closing
# Pattern: case-insensitive, substring match (tidak perlu exact)
# Gunakan regex kalau mau lebih spesifik
closing_triggers:
  deal:
    - "sudah transfer"
    - "dah transfer"
    - "sudah bayar"
    - "dah bayar"
    - "oke gas"
    - "ok gas"
    - "deal"
    - "saya mau"
    - "gas lah"
    - "kirim invoice"
    - "kirim rekening"
    - "mau pesan"
    - "pasang dong"
    - "oke saya ambil"
  hot_lead:
    - "berapa harga"
    - "harganya berapa"
    - "bisa diskon"
    - "paket apa aja"
    - "gimana cara pesan"
    - "minimal order berapa"
  stop:
    - "jangan hubungi"
    - "jangan dikontak"
    - "berhenti"
    - "tidak tertarik"
    - "sudah ada yang handle"
    - "jangan spam"
    - "unsubscribe"
    - "hapus nomor saya"

schedule:
  best_days: ["selasa", "rabu", "kamis"]
  best_hours: "10:00-15:00"    # owner WO biasa cek WA siang

scoring:
  has_instagram: +3           # aktif di sosmed = siap digital
  no_website: +5              # target sempurna
  rating_3_to_4: +2           # lumayan tapi masih bisa tumbuh
  review_count_under_50: +1   # bisnis yang lagi naik
```

**Kenapa area granular:** Satu kota = 50-200 leads. 8 kota = 400-1600 leads. Lu cuma tentuin radius per kota + kecamatan mana aja, WaClaw yang scan semua. Masing-masing area di-scrape paralel per worker.

### Template File — Pure Text, Dynamic Placeholders

```markdown
# ~/.waclaw/niches/undangan_digital/offer_1.md

Halo Kak {{.Title}}! 👋

Aku lihat {{.Title}} di {{.Address}} udah cukup ramai ya, tapi sayang belum ada cara digital buat client order undangannya.

Kebetulan aku bikinin template undangan digital yang bisa custom nama, tanggal, lokasi, dan bahkan embed video. Client tinggal klik link, isi data, langsung jadi.

Boleh aku kirim contohnya ke sini Kak? Gratis preview-nya, dilihat dulu aja 😁

— Haikal
```

**Available placeholders:**
```
{{.Title}}       → nama bisnis (contoh: "Kopi Nusantara")
{{.Category}}    → kategori (contoh: "cafe")
{{.Address}}     → alamat lengkap
{{.City}}        → kota saja
{{.Rating}}      → rating Google (contoh: "4.2")
{{.Reviews}}     → jumlah review (contoh: "87")
{{.Area}}        → area dari config (contoh: "kediri")
```

**Kenapa file, bukan UI:** Lu power user. Lu copy-paste template dari notes. Lu version-control pake git. Lu share ke tim. File nggak butuh settings screen. File ITU settings-nya.

### snippets.md — Quick Reply Snippets

```markdown
# ~/.waclaw/snippets.md

## soft pitch
boleh lihat dulu aja kak, gratis kok

## free offer
aku kasih preview gratis ya kak, tinggal liat aja

## move to call
bisa telepon bentar buat diskusi kak? lebih enak jelasinnya

## direct price
harganya mulai 500rb kak, tergantung fitur

## send sample
mau aku kirim contohnya ke sini kak?
```

### do_not_contact.yaml — Auto-Block List

```yaml
# ~/.waclaw/do_not_contact.yaml
# Orang yang minta berhenti dikontak. WaClaw nggak bakal kirim ke nomor ini.
# Auto-populated kalau closing_triggers.stop terdeteksi.
# Manual edit juga bisa — tambahin nomor yang lu tau nggak mau dikontak.

do_not_contact:
  - number: "0812-3456-7890"
    reason: "stop_detected"
    source: "web_developer"
    date: "2024-05-02"
    message: "jangan hubungi lagi ya"

  - number: "0813-9876-5432"
    reason: "manual_block"
    source: "undangan_digital"
    date: "2024-05-01"
    message: "sudah ada yang handle"

  - number: "0857-1111-2222"
    reason: "stop_detected"
    source: "all"
    date: "2024-04-28"
    message: "jangan spam"
```

**Kenapa file terpisah:** Ini bukan niche-specific. Kalau orang bilang "jangan hubungi" di satu niche, mereka nggak mau dikontak di niche lain juga. Satu list global = aman. Dan kalau lu tau sendiri nomor yang toxic, bisa manual tambahin.

---
