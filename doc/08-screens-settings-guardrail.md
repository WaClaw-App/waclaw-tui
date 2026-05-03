### SCREEN 13: SETTINGS → CONFIG REFERENCE

**Bukan UI settings. Ini reference card. Semua config di file.**

**State: settings_overview**

```

  pengaturan

  semua konfigurasi disimpan di file.
  edit pakai editor favorit lu.

  ────────────────────────────────────────────────────

  config utama     ~/.waclaw/config.yaml
  tema warna       ~/.waclaw/theme.yaml
  query pencarian  ~/.waclaw/queries.md
  niche folder     ~/.waclaw/niches/

  ────────────────────────────────────────────────────

  config aktif:

  niche aktif      web_developer, undangan_digital, social_media_mgr
  wa slots         3 nomor aktif (rotator)
  worker pool      3 worker (1 per niche)
  area             multi (8 kota, 14 kecamatan)
  jam kerja        09:00-17:00 wib
  rate limit       6/jam per slot, 50/hari total
  rotator mode     round-robin + cooldown
  auto-pilot       aktif

  ────────────────────────────────────────────────────

  e  edit config    r  reload    q  balik

```

**State: settings_edit** (buka editor)

```

  edit config

  buka file ini di editor lu:

    ~/.waclaw/config.yaml

  setelah save, tekan r buat reload.

  r  reload    q  balik

```

**State: settings_reload** (setelah edit & reload)

```

  config di-reload ✓

  perubahan:
  - rate_limit: 6 → 8 per jam
  - template default: direct-curiosity → pattern-interrupt

  perubahan langsung berlaku.
  scraper & sender auto sesuaikan.

   ↵  balik ke dashboard

```

**State: settings_reload_error** (edit config, reload, tapi error)

```

  ✗ config error — reload gagal

  ────────────────────────────────────────────────────

  ~/.waclaw/config.yaml

  2 masalah:

  ✗  baris 23: parse error
     │  anti_ban:
     │    rate_limit_per_slot: "enam"
     │                          ^^^^^
     │  harus angka, bukan teks

  ✗  baris 41: field tidak dikenali
     │  teleport_mode: true
     │  ^^^^^^^^^^^^^
     │  field ini gak ada di schema

  ────────────────────────────────────────────────────

  config lama masih dipakai. nggak ada yang berubah.
  fix dulu baru reload.

  1  buka file lagi    2  revert ke backup    q  balik

  backup terakhir: ~/.waclaw/config.yaml.bak
  (auto-disimpan setiap reload sukses)

```

Micro-interactions:
- `settings_reload_error`: config yang error highlight dengan red line marker = langsung keliatan
- "config lama masih dipakai" = explicit reassurance, nothing broke
- `2 revert ke backup` = safety net, one-key rollback
- Auto-backup mention: "auto-disimpan setiap reload sukses" = trust builder
- Error line marker `^^^^^` underlines exact problem = zero ambiguity

---

### SCREEN 14: GUARDRAIL → CONFIG VALIDATION

**Config broken = army lumpuh, tapi user nggak tau kenapa.**

**Screen ini muncul otomatis kalau ada config error. Bisa juga diakses manual dari settings.**

**Filosofi:** Validate early, fail loudly. Broken config = paused army. Silent errors = invisible disaster. Setiap config error HARUS nongol ke permukaan SEKARANG, bukan pas runtime pas udah kirim 50 pesan ke orang salah.

**State: validation_clean** (semua config pass)

```

  config check                             ✓ semua bersih

  ────────────────────────────────────────────────────

  config.yaml          ✓  ok
  theme.yaml           ✓  ok
  queries.md           ✓  ok

  niches:
  web_developer        ✓  ok    3 template · 5 area
  undangan_digital     ✓  ok    2 template · 8 area
  social_media_mgr     ✓  ok    1 template · 3 area

  ────────────────────────────────────────────────────

  3 niche siap tempur. 6 template loaded.
  semua worker bisa jalan tanpa masalah.

  lu bisa gas. army lu siap.

   ↵  lanjut    q  balik

```

Micro-interactions:
- `✓ semua bersih` green pulse 1x = reassurance
- Setiap niche row brief green flash stagger 80ms = wave of approval
- Auto-dismiss setelah 3 detik kalau dari boot sequence = lu nggak perlu nunggu
- "lu bisa gas" = explicit green light, bukan cuma "tidak ada error"

**State: validation_errors** (ada config yang broken)

```

  config check                             ✗ ada error

  ────────────────────────────────────────────────────

  config.yaml          ✓  ok
  theme.yaml           ✓  ok
  queries.md           ✗  1 error

  niches:
  web_developer        ✓  ok    3 template · 5 area
  undangan_digital     ✗  2 error
  fotografer           ✗  3 error

  ────────────────────────────────────────────────────

  2 niche di-pause sampai diperbaiki:
  undangan_digital, fotografer

  web_developer tetap jalan.

  detail error:

  ── queries.md ──

  ✗  baris 8: format salah
     │  - "cafe di
     │        ^ kurang tanda kutip penutup

  ── undangan_digital/niche.yaml ──

  ✗  baris 14: parse error
     │  targets:
     │    - "wedding organizer
     │          ^ kurang tanda kutip penutup

  ✗  field wajib kosong: areas

  ── fotografer/niche.yaml ──

  ✗  baris 7: tipe salah
     │  rating_min: "rendah"
     │               ^^^^^^^^ harus angka

  ✗  field wajib kosong: areas

  ✗  fotografer/ice_breaker.md kosong

  ────────────────────────────────────────────────────

  1  buka file pertama yang error
  2  liat contoh config
  r  reload setelah fix    q  balik

  tekan 1 lagi buat buka file error berikutnya

```

Micro-interactions:
- `✗ ada error` red, tapi `✓ ok` rows tetap green = partial system masih jalan
- Error count per file: `✗ 2 error` = instant severity assessment
- "2 niche di-pause" = explicit consequence, bukan cuma "ada error"
- "web_developer tetap jalan" = reassurance, bukan total failure
- Error lines show `│` gutter + `^` pointer = exactly where to look
- `1 buka file pertama yang error` = zero friction, langsung ke masalah pertama
- `1 lagi` buat next error = tab-through semua error tanpa keluar screen
- On reload success: error rows collapse satu-satu ke `✓ ok` = satisfying fix feel

**State: validation_warnings** (config valid tapi ada warnings)

```

  config check                             ⚠️ ada warning

  ────────────────────────────────────────────────────

  config.yaml          ⚠️  1 warning
  theme.yaml           ✓  ok
  queries.md           ✓  ok

  niches:
  web_developer        ✓  ok    3 template · 5 area
  undangan_digital     ⚠️  1 warning
  social_media_mgr     ✓  ok    1 template · 3 area

  ────────────────────────────────────────────────────

  semua niche tetap jalan. warning bukan error.

  detail warning:

  ── config.yaml ──

  ⚠  field deprecated: send_delay
     │  send_delay: 12
     │  ^^^^^^^^^ gunakan min_delay_between
     │  masih works, tapi bakal dihapus versi depan

  ── undangan_digital/niche.yaml ──

  ⚠  rating_min (3.0) lebih rendah dari rekomendasi (3.5)
     │  lead dibawah 3.5 biasanya kurang potensial
     │  ini cuma saran, bukan wajib

  ────────────────────────────────────────────────────

  ↵  lanjut aja    1  buka file    r  reload    q  balik

```

Micro-interactions:
- `⚠️ ada warning` amber, bukan merah = not blocking, tapi worth knowing
- "semua niche tetap jalan. warning bukan error." = explicit difference
- Deprecated fields: show old + new field name = migration path jelas
- `↵ lanjut aja` = pre-highlighted default = warning = soft gate, bukan hard block

**State: validation_fix** (after user fixes and presses r to re-validate)

```

  config check                             ● re-validating...

  ────────────────────────────────────────────────────

  config.yaml          ✓  ok
  theme.yaml           ✓  ok
  queries.md           ✓  ok    ← fixed!

  niches:
  web_developer        ✓  ok    3 template · 5 area
  undangan_digital     ● checking...
  fotografer           ○ waiting

  ────────────────────────────────────────────────────

  sedang cek ulang...
  1 dari 3 sudah ok.

```

Validation runs sequentially per file. Tiap file yang pass gets `✓ ok` with brief green flash. Files with errors show errors again. When ALL pass:

```

  config check                             ✓ semua bersih!

  ────────────────────────────────────────────────────

  semua error udah diperbaiki. ✓

  config.yaml          ✓  ok
  queries.md           ✓  ok    ← fixed!
  undangan_digital     ✓  ok    ← fixed!
  fotografer           ✓  ok    ← fixed!

  ────────────────────────────────────────────────────

  3 niche yang tadi di-pause sekarang auto-resume.
  army lu udah siap lagi.

   ↵  gas

```

Micro-interactions:
- Re-validation: `● checking...` pulsing = processing, jangan panik
- Fixed items get `← fixed!` brief highlight = lu tau mana yang barusan lu perbaiki
- "3 niche yang tadi di-pause sekarang auto-resume" = explicit resolution
- `↵ gas` = back to action, no lingering

**Variant: validation_first_time** (first-time setup, more guidance)

```

  config check                             pertama kali ya?

  ────────────────────────────────────────────────────

  belom ada config. bikin dulu ya.

  yang perlu lu siapin:

  1  ~/.waclaw/config.yaml        (setting utama)
  2  ~/.waclaw/niches/.../         (folder per niche)

  tenang, waclaw bakal bikinin template
  yang lu cuma isi. tekan 1 buat mulai.

  ────────────────────────────────────────────────────

  1  bikin config otomatis    2  liat contoh    q  keluar

  template config bakal dibikin di:
  ~/.waclaw/config.yaml
  ~/.waclaw/niches/_contoh/niche.yaml

```

Micro-interactions:
- "pertama kali ya?" = friendly, bukan error
- "tenang, waclaw bakal bikinin template" = lowering anxiety
- `1 bikin config otomatis` = zero-to-running in one key
- After auto-generate: auto-switch ke `validation_clean` = instant validation loop

---

