---

## 9. Startup Sequence — 4 Detik

```bash
$ waclaw
```

```
  t +0ms    logo render per karakter
  t +80ms   tagline fade in: "leads lu pada nunggu"
  t +200ms  system check (wa, config, db, lisensi)
  t +250ms  update check (background, non-blocking)
  t +300ms  lisensi check:
            ✓ lisensi valid (device: LAPTOP-HOME, expires: 30 jun 2025)
            ✗ lisensi expired → hard stop
            ✗ device conflict → hard stop
  t +400ms  config validation (semua niche.yaml, template, config.yaml)
  t +500ms  validation result:
            ✓ config.yaml ok
            ✓ 3/4 niche ok
            ✗ fotografer: parse error
  t +700ms  status:
            ● lisensi ok
            ● wa nyambung
            ● niche: web_developer, undangan_digital, social_media_mgr
            ○ niche: fotografer (di-pause — tekan v)
            ● 847 leads di database
  t +800ms  auto-pilot: ON (3 niche, 1 paused)
  t +900ms  army marching: 3 workers → ● aktif
  t +1100ms dashboard fade in
  t +1300ms ready. cursor blinks.
```

**1300ms sampai bisa dipakai.** Setiap ms sebelum itu narik perhatian. Setiap ms setelah itu keputusan lu.

Lisensi check jadi bagian dari startup. Kalau expired atau device conflict, army NGGAK jalan — hard stop. Config validation jalan setelah lisensi ok. Kalau ada config error, worker yang bermasalah di-pause tapi yang ok tetap jalan. Lu nggak perlu fix semuanya dulu baru bisa mulai.

Kalau ada masalah:

```
  t +500ms  validation result:
            ✗ lisensi belum ada → tekan l buat masukin
            ● wa nyambung
            ○ 0 leads di database
```

`✗` cuma warna merah di screen. Ga bisa di-ignore. Ga perlu baca error message — warnanya ITU pesannya.

**Variant: first time vs returning:**
- First time → lisensi check → license input → validation → 3 step onboarding (login → niche → scrape)
- Returning → lisensi check → validation → dashboard, auto-pilot ON (error niches paused)
- Lisensi expired → hard stop, masukin lisensi baru
- Device conflict → hard stop, putuskan device lain atau ganti lisensi

---

## 10. Session End — Kesan Terakhir

```bash
  # user tekan q dari monitor

  sampai jumpa.

  hari ini: 43 terkirim · 7 response · 2 deal
  lead terbaik: kopi nusantara (respond 11 menit)

  ── tips: selasa jam 10 itu waktu terbaik lu kirim pesan ──

  tekan q lagi buat keluar, atau apa aja buat lanjut
```

**Kenapa ini penting:** Hal terakhir yang lu liat = cara lu ingat whole session. Bukan tabel — cerita. "43 terkirim" = kerja keras. "7 response" = validasi. "2 deal" = menang. Tips selasa? Alasan buat balik.

Double-quit (`q` lagi) = prevent accident, tapi juga bikin momen reconsideration. "Eh, kirim beberapa lagi aja kali ya?" Itu Zeigarnik effect: **yang belum selesai nempel di kepala.**

---
