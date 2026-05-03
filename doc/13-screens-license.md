### SCREEN 18: LICENSE → GATE

**Tanpa lisensi, army nggak jalan. Titik.**

Screen ini muncul kalau lisensi belum di-input, expired, atau konflik device. Ini hard gate — WaClaw NGGAK akan jalan tanpa lisensi valid. Bukan soft warning. Bukan trial mode. Full stop.

**Filosofi lisensi:**
- Satu lisensi = satu device. Fair untuk semua.
- Lisensi dicek tiap startup. Kalau konflik, WaClaw berhenti.
- Lisensi disimpan di `~/.waclaw/license.md` — bukan hardcoded, bukan registry.
- Lisensi expired = army pause, bukan army mati. Data aman, tinggal perpanjang.
- Versi beda = lisensi beda. v1 lisensi = v1 army. v2 lisensi = v2 army. Nggak ada forced upgrade.

**State: license_input** (pertama kali, belum ada lisensi)

```

  lisensi waclaw

  waclaw butuh lisensi buat jalan.
  masukin key lisensi lu di bawah.

  ┌───────────────────────────────────────────────────┐
  │  WACL-XXXX-XXXX-XXXX-XXXX                         │
  │_                                                   │
  └───────────────────────────────────────────────────┘

   ↵  validasi    1  beli lisensi    q  keluar

  lisensi disimpan di: ~/.waclaw/license.md
  satu lisensi cuma buat satu device.

```

Micro-interactions:
- Input field pulse accent = lu di tempat yang tepat, mulai ngetik
- Auto-format: ketik `WACL-` auto-capitalize, auto-insert hyphen tiap 4 karakter
- Key validation: tekan ↵ → loading spinner 500ms → result

**State: license_validating** (lagi cek ke server)

```

  lisensi waclaw

  ●  nyambung ke server lisensi...
  ○  cek validitas
  ○  cek device

   ━━━━━━━━━━━━░░░░░░░░░░  validating...

```

Micro-interactions:
- `● ○ ○` animate sequential = sama kayak login screen, progress feels alive
- Spinner smooth 200ms rotation

**State: license_valid** (lisensi ok!)

```

  lisensi waclaw

  ✓ lisensi valid!

   ●  terhubung ke server
   ●  lisensi valid
   ●  device terdaftar

  lisensi: WACL-XXXX-XXXX-XXXX-XXXX
  device:  LAPTOP-HOME
  expires: 30 juni 2025

  key disimpan ke ~/.waclaw/license.md
  waclaw siap jalan. tekan ↵ buat lanjut.

   ↵  lanjut    q  keluar

```

Micro-interactions:
- `✓ lisensi valid!` green pulse = the gate is open
- Hold 800ms sebelum bisa press ↵ = lu ngerasain momen "masuk" dulu
- Auto-transition ke next screen setelah 2 detik

**State: license_invalid** (key salah / nggak valid)

```

  lisensi waclaw

  ✗ lisensi nggak valid — cek lagi key nya

  key yang lu masukin nggak cocok.
  pastiin nggak ada salah ketik.

  ┌───────────────────────────────────────────────────┐
  │  WACL-XXXX-XXXX-XXXX-XXXX                         │
  │_                                                   │
  └───────────────────────────────────────────────────┘

   ↵  coba lagi    1  beli lisensi    q  keluar

```

Micro-interactions:
- `✗` red flash = instant feedback
- Input field auto-focus = langsung bisa benerin
- Field gets subtle red border glow 800ms = "fix this"

**State: license_expired** (lisensi expired)

```

  lisensi waclaw

  ✗ lisensi lu udah expired

  lisensi: WACL-XXXX-XXXX-XXXX-XXXX
  device:  LAPTOP-HOME
  expired: 15 april 2025 (17 hari lalu)

  semua worker di-pause. data aman.
  perpanjang lisensi buat lanjut.

  1  masukin lisensi baru    2  beli perpanjangan    q  keluar

```

Micro-interactions:
- `✗` red = hard stop, nggak bisa di-skip
- "data aman" = reassurance, data nggak hilang
- "2 beli perpanjangan" = langsung ada jalan keluar

**State: license_device_conflict** (lisensi aktif di device lain)

```

  lisensi waclaw

  ✗ lisensi lagi dipakai di device lain. waclaw berhenti.

  lisensi: WACL-XXXX-XXXX-XXXX-XXXX
  device ini:   LAPTOP-HOME
  device lain:  PC-KANTOR
  terakhir aktif: 12 menit lalu

  satu lisensi cuma buat satu device.
  kalau mau pindah device, putuskan dulu dari device lama.

  1  masukin lisensi baru    2  putuskan device lain    q  keluar

  ── "2 putuskan" = force logout device lain, ambil alih lisensi ──

```

Micro-interactions:
- `✗` red flash = hard stop
- Device info shown = lu tau siapa yang pakai, bukan mystery error
- "2 putuskan device lain" = force transfer, 1 tombol selesai
- Setelah putuskan: auto-revalidate → license_valid
- Warning: "ini bakal logout waclaw di device lain" confirmation overlay sebelum execute

**State: license_server_error** (gagal cek lisensi — network down)

```

  lisensi waclaw

  ⚠  gagal nyambung ke server lisensi

  ●  server lisensi
  ✗  gagal nyambung

  lu punya lisensi yang valid sebelumnya.
  waclaw bakal jalan pake lisensi offline selama 72 jam.
  setelah itu, harus online buat re-validate.

  lisensi: WACL-XXXX-XXXX-XXXX-XXXX
  offline grace: 71 jam tersisa

  ↵  lanjut offline    1  coba lagi    q  keluar

```

Micro-interactions:
- `⚠` amber = warning, bukan error. Lu masih bisa jalan.
- Offline grace period: 72 jam = lu punya waktu buat fix internet
- Countdown `71 jam tersisa` = lu tau batasnya
- Kalau grace period habis: hard stop sama kayak expired

**license.md — Lisensi Disimpan di File**

```markdown
# ~/.waclaw/license.md

# Lisensi WaClaw — jangan diubah manual
# Lisensi dicek otomatis tiap startup

key: WACL-XXXX-XXXX-XXXX-XXXX
device: LAPTOP-HOME
activated: 2025-01-15
expires: 2025-06-30
last_validated: 2025-05-02T14:23:00+07:00
```

Kenapa file, bukan database? Karena lu bisa backup, bisa pindah device (dengan re-activate), dan bisa liat sendiri kapan lisensi lu expired. Transparansi = trust.

---

