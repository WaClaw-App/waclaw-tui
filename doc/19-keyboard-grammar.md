
## 8. Keyboard Grammar

**Hafal sekali. Pakai selamanya. Tombol yang sama, tiap screen.**

### Universal Keys (selalu ada di semua screen)

| Tombol | Selalu Ngejek |
|--------|---------------|
| `↑/↓` | Pindah item |
| `↵` | Aksi utama (yang paling masuk akal) |
| `1-9` | Aksi sekunder (beda tiap screen) |
| `s` | Skip / buang |
| `q` | Selesai / balik / keluar (context-dependent) |
| `p` | Pause yang lagi jalan |
| `r` | Refresh / reload |
| `/` | Cari / filter (kalau ada list) |
| `?` | Tampilin shortcut (overlay, bukan screen baru) |
| `v` | Validate — cek semua config |
| `l` | License — liat / ganti lisensi |
| `h` | History — liat timeline |
| `` ` `` | Nerd stats — toggle RAM/CPU overlay |
| `u` | Update — cek versi baru |
| `Ctrl+K` | Command palette — search & execute apapun |
| `esc` | Batal / tutup modal / keluar compose / tutup palette |
| `n` | Baru — bikin item baru (template, worker) |
| `e` | Edit — edit item / config |
| `tab` | Pindah tab / niche group |
| `space` | Centang / hapus pilihan (niche select) |

### Screen-Local Keys (cuma ada di screen tertentu)

| Tombol | Screen | Aksi |
|--------|--------|------|
| `x` | Lead Review | Skip & block lead |
| `d` | Lead Review | Liat detail lead |
| `a` | Follow-Up | Auto-approve semua |
| `+` | Login | Tambah slot WA |
| `←` | History | Hari sebelumnya |
| `→` | History | Hari berikutnya |
| `w` | History | Pindah ke tampilan minggu |
| `t` | History | Pindah ke tampilan hari ini |
| `space` | Niche Select | Centang/hapus niche |

### Perilaku `q` (Context-Dependent)

`q` selalu berarti "balik" tapi beda-beda tergantung context:

- **Di sub-state** (detail, settings, preview): balik ke state sebelumnya dalam screen yang sama
- **Di screen utama** (overview, dashboard): balik ke screen sebelumnya (pop navigation stack)
- **Di root screen**: pertama kali → tampilin session summary, kedua kali → keluar app

Contoh:
- Di Shield slot detail → `q` balik ke Shield overview
- Di Shield overview → `q` balik ke Monitor
- Di Monitor (root) → `q` tampilin session summary → `q` lagi keluar

### Perilaku `v` (Context-Dependent)

- **Di Send failed state**: `v` = validate & retry (aksi lokal)
- **Di screen lain**: `v` = navigasi ke Guardrail screen

**Overlay `?`:**

```

  ── tombol ─────────────

  ↑↓   pindah
  ↵    aksi utama
  1-9  pilih opsi
  s    skip
  q    balik/keluar
  p    pause
  n    baru
  e    edit
  /    cari
  tab  pindah tab
  esc  batal/tutup
  v    validasi config
  l    lisensi
  h    riwayat
  r    reload
  `    nerd stats
  u    cek update
  ^K   command palette

  tekan apa aja buat tutup

```

Fade in di atas screen sekarang. Fade out di keypress apapun. Nggak pernah ilang context.

**Tapi ingat:** kalau lu nggak ngetik apa-apa, itu berarti WaClaw kerja dengan bener. Keyboard itu hak istimewa, bukan kewajiban.

