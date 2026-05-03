
## 4. Color System

**Bukan tema. Tapi mood.**

```yaml
# ~/.waclaw/theme.yaml
# ubah apapun. rusak apapun. milik lu.

palette:
  # Background
  bg:           "#0A0A0B"     # hampir hitam, bukan pure — lebih soft
  bg_raised:    "#141416"     # elevasi halus
  bg_active:    "#1A1A1E"     # zona aktif

  # Teks
  text:         "#E8E8EC"     # primer — putih hangat
  text_muted:   "#6B6B76"     # sekunder — bisik, jangan teriak
  text_dim:     "#3D3D44"     # tersier — hampir ga keliatan, tapi bisa dibaca

  # Sinyal
  success:      "#34D399"     # hijau — bukan hijau agresif
  warning:      "#FBBF24"     # amber — perhatian, bukan alarm
  danger:       "#F87171"     # merah — jelas, bukan menakutkan
  accent:       "#818CF8"     # indigo — brand, aksi

  # Motion
  pulse:        "#818CF866"   # accent 40% opacity — elemen yang bernapas
  highlight:    "#FFFFFF22"   # putih 13% — zona hover/fokus

  # Celebration
  gold:         "#FFD700"     # jackpot & revenue — earned celebration
  celebration:  "#FFFFFF"     # full-screen flash — conversion only

  # Locale
  locale:       id            # "id" (default) atau "en" — bahasa UI
```

**Rules:**
- `danger` merah HANYA buat masalah teknis (wa putus, scraper error, config error). Rejection = `text_muted`.
- `accent` (indigo) cuma muncul di hal yang bisa lu interaksi. Kalau indigo, berarti bisa ditekan.
- `pulse` buat indikator hidup (connection status, scrape aktif). Bernapas = pulse.
- Rejection, skip, decline = netral. Bukan gagal. Cuma bukan jodoh.
- `gold` HANYA buat jackpot leads, revenue numbers, dan conversion celebration. Kalau gold, berarti UANG.
- `celebration` putih full-flash ONLY untuk conversion. Satu-satunya momen yang boleh secerah itu.
- `locale` mengatur bahasa TUI. `"id"` = casual Indonesian (default), `"en"` = casual English. Switchable runtime via Ctrl+K atau `l`.

---

## 5. Layout System — Vertical Borderless

**No `│ ┃ ║ ╼ ╽`. No `─ ━ ─`. No boxes. Cuma ruang.**

```
Spacing unit: 1 line = 8px equivalent
Section gap:  2 lines (16px)
Sub-section:  1 line (8px)
Item gap:     0 lines (touching = grouped)
Indent:       2 spaces per level
```

### Contoh: Lead Detail

```

  kopi nusantara
  cafe · jl. hasanuddin 23, kediri

  ⭐ 4.2  87 reviews
  no website  no instagram

  dikontak: 2 kali
  pesan terakhir: "boleh lihat desainnya?" — kemarin 14:23

  status: respond → nunggu offer

  1  kirim offer    2  balas custom    3  archive

```

**Yang lu ga liat:** No border. No label "Nama:" sebelum nama bisnis. No prefix "Status:". Layout-nya ITU hierarchy:
- Baris 1: Bold/bright = identitas
- Baris 2: Muted = konteks
- Baris 3-4: Detail = metadata
- Baris 5-6: Riwayat = narasi
- Baris 7: Status = kondisi sekarang
- Baris 9: Aksi = langkah lu

Baca = atas ke bawah. Keputusan = bawah ke atas. Dua arah lancar karena layout predictable setelah 2x pakai.

---
