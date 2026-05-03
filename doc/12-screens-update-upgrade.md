### SCREEN 20: UPDATE & UPGRADE → RENEWAL

**Update kecil = gratis. Upgrade besar = lisensi baru.**

Screen ini ngatur versi WaClaw. Ada dua jenis perubahan:
- **Update** (minor, v1.1→v1.2) = bug fix, performance, fitur kecil. Gratis. Lisensi tetap valid.
- **Upgrade** (major, v1.x→v2.x) = arsitektur baru, fitur besar, breaking change. Butuh lisensi baru.

**Filosofi versi:**
- v1 army = v1 lisensi. Selama lisensi lu valid, v1 tetap jalan. Nggak ada forced upgrade.
- Minor update = gratis, otomatis, nggak ganggu kerjaan lu.
- Major upgrade = product baru. Kayak beli iPhone baru — yang lama tetep bisa dipakai.
- WaClaw ngecek update pas startup, tapi NON-BLOCKING. Kalau ada update, notif aja. User yang tentuin kapan update.
- Nggak pernah auto-update tanpa persetujuan. Kalau lu lagi kirim pesan, update nunggu.

**State: update_available** (ada minor update)

```

  update tersedia

  ────────────────────────────────────────────────────

  versi lu:     v1.3.2
  versi baru:   v1.3.3

  yang baru:
    • fix: scrape terkadang duplicate di area tertentu
    • fix: WA rotator cooldown calculation lebih akurat
    • perf: database query 15% lebih cepat

  ini update kecil (minor). gratis.
  lisensi lu tetap valid. nggak perlu bayar.

  ────────────────────────────────────────────────────

  ↵  update sekarang   1  nanti aja   q  skip

  update butuh restart. waclaw bakal nunggu sampe
  semua worker idle sebelum restart.

```

**State: update_downloading** (lagi download update)

```

  update tersedia

  lagi download v1.3.3...

  ████████████████████░░░░  82% (4.2MB / 5.1MB)

  estimasi: 12 detik

  sumber: https://releases.waclaw.dev/v1.3.3/

  worker tetap jalan selama download.
  restart cuma setelah download selesai + lu approve.

   q  batal download

```

**State: update_ready** (download selesai, siap restart)

```

  update tersedia

  ✓ v1.3.3 siap di-install!

  ────────────────────────────────────────────────────

  download selesai: 5.1MB
  checksum: ✓ verified
  backup: ~/.waclaw/waclaw-v1.3.2.bak

  ────────────────────────────────────────────────────

  restart sekarang? semua worker bakal di-pause dulu.
  proses restart ±3 detik. data aman.

   ↵  restart sekarang   1  nanti (update tetap ada)   q  skip

  kalau lu skip, waclaw bakal ngingetin lagi pas startup berikutnya.

```

**State: upgrade_available** (ada major upgrade — butuh lisensi baru)

```

  upgrade tersedia

  ────────────────────────────────────────────────────

  versi lu:     v1.3.2
  versi baru:   v2.0.0

  ⚠ ini upgrade besar (major version)

  yang baru:
    • arsitektur baru: multi-device support
    • AI-powered message personalization
    • dashboard real-time collaboration
    • 20+ fitur baru lainnya

  ────────────────────────────────────────────────────

  upgrade ke v2 butuh lisensi baru.
  lisensi v1 lu nggak berlaku buat v2.

  ini bukan pelit — ini product yang beda.
  v1 lu tetep bisa dipakai selama lisensi valid.
  nggak ada forced upgrade.

  lisensi v1 lu:
    key: WACL-XXXX-XXXX-XXXX-XXXX
    expires: 30 juni 2025
    status: ✓ valid — v1 tetap jalan

   ↵  beli lisensi v2   1  liat detail v2   q  tetap v1

```

**State: upgrade_license_input** (input lisensi v2)

```

  upgrade ke v2

  masukin lisensi v2 lu.

  ┌───────────────────────────────────────────────────┐
  │  WACL2-XXXX-XXXX-XXXX-XXXX                         │
  │_                                                   │
  └───────────────────────────────────────────────────┘

  lisensi v1 lu tetap aktif buat v1.
  kalau v2 lisensi nggak valid, waclaw tetap jalan v1.

   ↵  validasi v2    q  batal

```

**State: license_expired_with_upgrade** (lisensi v1 expired + ada v2)

```

  lisensi expired

  ────────────────────────────────────────────────────

  lisensi v1 lu udah expired (15 april 2025).
  v1 army berhenti.

  tapi... ada v2!

  ✓ v2.0.0 tersedia
  ✓ lisensi v2 = akses ke semua fitur baru
  ✓ bisa pake data yang sama (auto-migrate)

  ────────────────────────────────────────────────────

  1  perpanjang v1 (lisensi lama, fitur lama)
  2  upgrade ke v2 (lisensi baru, fitur baru)
  3  masukin lisensi baru (v1 atau v2)

  q  keluar

```

**Variant: startup_check** (background update check saat boot)

```

  # di startup sequence, setelah lisensi check:

  t +250ms  update check (background, non-blocking)
            ✓ v1.3.2 (latest) — nggak ada update
            ○ v1.3.3 tersedia → notif: "update Available"
            ○ v2.0.0 tersedia → notif: "upgrade Available"

  # check jalan di background. NGGAK nungguin hasilnya.
  # kalau ada update, notification muncul setelah dashboard ready.
  # NGGAK pernah block startup.

```

Micro-interactions:
- Update available: versi number morph dari lama ke baru = perubahan visible
- Changelog items slide in stagger 50ms = scan-able, bukan wall of text
- "gratis" badge pulse hijau = instant positive signal
- Download progress: bar fills real-time + speed shown = lu tau ini kerja
- Checksum verified `✓` = trust, file nggak corrupt
- Backup info shown = rollback possible, nggak ada fear
- Upgrade available: `⚠` amber = perhatian, bukan bahaya
- "v1 tetep bisa dipakai" = reassurance, nggak ada FOMO
- License input v2: prefix `WACL2-` auto-detect = beda format, beda product
- license_expired_with_upgrade: 3 opsi = lu punya pilihan, bukan dead end
- Startup check: non-blocking = WaClaw NGGAK nunggu network buat mulai kerja
- Restart countdown: `±3 detik` = expectation set, nggak takut nunggu lama

**The neuroscience:** Versi itu commitment. Minor update = low commitment (gratis, cepat, nggak ngubah apa-apa). Major upgrade = high commitment (bayar, pindah, adaptasi). WaClaw bikin dua-duanya jelas: minor = seamless, major = conscious choice. Nggak ada dark pattern, nggak ada forced upgrade, nggak ada "versi lu udah deprecated". v1 army lu jalan selama lisensi valid. Kalau mau v2, itu pilihan lu, bukan keharusan.

---

