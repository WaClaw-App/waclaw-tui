
## 1. SCREENS — Semua State & Variant

Setiap screen punya **states** (kondisi sekarang) dan **variants** (tampilan alternatif berdasarkan context).

---

### SCREEN 1: BOOT → FIRST IMPRESSION

**State: first_time**

```

                                                  
  ▄▄▄                 ▄   ▄▄▄▄ ▄▄                 
 █▀██  ██  ██▀▀       ▀██████▀  ██                
   ██  ██  ██           ██      ██                
   ██  ██  ██ ▄▀▀█▄     ██      ██ ▄▀▀█▄▀█▄ █▄ ██▀
   ██▄ ██▄ ██ ▄█▀██     ██      ██ ▄█▀██ ██▄██▄██ 
   ▀████▀███▀▄▀█▄██     ▀█████ ▄██▄▀█▄██  ▀██▀██▀ 
                                                  
                                                  

      leads lu pada nunggu. yuk mulai.


      ── pertama kali? ─────────────────────

      1  login        hubungin whatsapp lu
      2  atur niche   pilih target & filter
      3  gas          mulai cari leads

      ──────────────────────────────────────

      1 → 2 → 3. gitu doang.

```

Micro-interactions:
- Logo render per karakter (8ms/char) → anticipation build
- Menu fade in sequential 120ms stagger → guided attention
- Press `1` → pulse bright → putih → transition, bukan jump
- Kalau `config.yaml` udah ada, step 2: `✓ udah diatur (2 buat edit)`

**State: returning** (sudah pernah login + configure)

```

                                                  
  ▄▄▄                 ▄   ▄▄▄▄ ▄▄                 
 █▀██  ██  ██▀▀       ▀██████▀  ██                
   ██  ██  ██           ██      ██                
   ██  ██  ██ ▄▀▀█▄     ██      ██ ▄▀▀█▄▀█▄ █▄ ██▀
   ██▄ ██▄ ██ ▄█▀██     ██      ██ ▄█▀██ ██▄██▄██ 
   ▀████▀███▀▄▀█▄██     ▀█████ ▄██▄▀█▄██  ▀██▀██▀ 
                                                  

  ── army report ────────────────────────────

      ● wa terhubung (3 nomor)
      ● 3 niche aktif · 4 worker jalan
      ● 847 leads di database

      auto-pilot aktif. semua niche jalan paralel.
      3 nomor WA rotasi kirim, aman dari ban.
      tekan apa aja buat liat dashboard, atau
      biarin aja — army lu lagi kerja nih.

  ──────────────────────────────────────────

```

**Army marching animation (returning users only):**
Setelah logo render, 3 baris "soldier" marching masuk dari kiri — tiap baris = 1 niche worker. Icon `▸▸▸` marching step-by-step lalu settle jadi `● aktif`. Durasi total 600ms. Feel: pasukan lu udah siap, udah jalan, lu cuma datang buat inspeksi.

```
  ── army marching ──

  ▸▸▸▸▸▸  web_developer     ● aktif
    ▸▸▸▸▸  undangan_digital  ● aktif
      ▸▸▸  social_media_mgr  ● aktif

  3 worker udah jalan. lu telat datang, mereka nggak.
```

Animation: tiap worker row slide in dari kiri 80ms stagger, lalu `▸▸▸` morph jadi `●` dengan overshoot bounce. Seperti unit militer yang nge-snap ke attention.

**Variant: returning + ada response baru**

```

      ● wa terhubung (3 nomor)
      ● 3 niche aktif · 4 worker jalan
      ● 847 leads di database
      ● 3 response baru!

      ada yang balas! tekan ↵ buat liat.

```

**Variant: returning + wa disconnect**

```

      ✗ wa putus — semua worker pause
      ● 3 niche (scrape tetap jalan, kirim pause)
      ● 847 leads di database

      scraper tetap nyari, cuma kirim yang pause.
      tekan 1 buat login ulang.

```

**Variant: returning + config error detected**

```
                                                  
  ▄▄▄                 ▄   ▄▄▄▄ ▄▄                 
 █▀██  ██  ██▀▀       ▀██████▀  ██                
   ██  ██  ██           ██      ██                
   ██  ██  ██ ▄▀▀█▄     ██      ██ ▄▀▀█▄▀█▄ █▄ ██▀
   ██▄ ██▄ ██ ▄█▀██     ██      ██ ▄█▀██ ██▄██▄██ 
   ▀████▀███▀▄▀█▄██     ▀█████ ▄██▄▀█▄██  ▀██▀██▀ 
                                                  
                                                  

      ✗ config error — worker pause
      ● wa terhubung (3 nomor)
      ● 2 niche ok · 1 niche bermasalah
      ● 847 leads di database

      niche "fotografer" punya config error.
      worker lain tetap jalan, yang error di-pause.
      tekan v buat liat detail error.

       v  liat error    ↵  dashboard    q  keluar

```

Micro-interactions:
- `●` pulses gently = alive
- `✗` satu-satunya warna merah di screen = auto-draw attention
- "3 response baru!" flash amber 2x = urgency tanpa panic
- Army marching: `▸▸▸` → `●` morph with bounce = workers reporting for duty
- Config error variant: `✗` red flash, tapi `● 2 niche ok` = partial system masih jalan, bukan total failure
- Auto-transition ke dashboard setelah 3 detik kalau nggak ada input = hands-off default

**Variant: returning + lisensi expired**

```

                                                  
  ▄▄▄                 ▄   ▄▄▄▄ ▄▄                 
 █▀██  ██  ██▀▀       ▀██████▀  ██                
   ██  ██  ██           ██      ██                
   ██  ██  ██ ▄▀▀█▄     ██      ██ ▄▀▀█▄▀█▄ █▄ ██▀
   ██▄ ██▄ ██ ▄█▀██     ██      ██ ▄█▀██ ██▄██▄██ 
   ▀████▀███▀▄▀█▄██     ▀█████ ▄██▄▀█▄██  ▀██▀██▀ 
                                                  
                                                  

      ✗ lisensi expired — army berhenti
      ● wa terhubung (3 nomor)
      ● 3 niche (semua pause)
      ● 847 leads di database

      lisensi lu udah expired. semua worker di-pause.
      perpanjang lisensi buat lanjut.

       1  masukin lisensi baru    2  beli lisensi    q  keluar

```

**Variant: returning + device conflict**

```

                                                  
  ▄▄▄                 ▄   ▄▄▄▄ ▄▄                 
 █▀██  ██  ██▀▀       ▀██████▀  ██                
   ██  ██  ██           ██      ██                
   ██  ██  ██ ▄▀▀█▄     ██      ██ ▄▀▀█▄▀█▄ █▄ ██▀
   ██▄ ██▄ ██ ▄█▀██     ██      ██ ▄█▀██ ██▄██▄██ 
   ▀████▀███▀▄▀█▄██     ▀█████ ▄██▄▀█▄██  ▀██▀██▀ 
                                                  
                                                  

      ✗ lisensi aktif di device lain — waclaw berhenti

      waclaw deteksi lisensi lu lagi dipakai di device lain.
      satu lisensi cuma buat satu device.
      semua worker di-pause sampai masalah ini selesai.

       1  masukin lisensi baru    2  putuskan device lain    q  keluar

      ── device lain: PC-KANTOR · terakhir aktif 12 menit lalu ──

```

Micro-interactions:
- License expired: `✗` red flash + semua `●` dim = total pause, bukan partial
- Device conflict: `✗` red flash + info device lain = lu tau siapa yang pakai
- "2 putuskan device lain" = force logout device lain, ambil alih lisensi
- Both variants: WaClaw NGGAK jalan sama sekali tanpa lisensi valid = hard gate

---

### SCREEN 2: LOGIN → TRUST

**Bisa konek lebih dari 1 nomor WA. Tiap nomor = 1 sender slot.**

**State: qr_waiting**

```

  hubungin whatsapp

  scan pakai hp lu. pelan aja.
  bisa tambah lebih dari 1 nomor buat rotator.
  makin banyak nomor, makin aman dari ban.

         ┌─────────────────┐
         │                 │
         │    [QR CODE]    │
         │                 │
         └─────────────────┘

         nunggu scan...          [1/3 slot]

   ●  nyambung ke server wa
   ○  nunggu scan dari hp
   ○  sinkron kontak

   slot terisi: 0/3   +  tambah slot    ↵  skip

```

**State: qr_scanned** (detected scan, syncing)

```

  hubungin whatsapp

  ✓ scan terdeteksi!

   ●  nyambung ke server wa
   ●  scan berhasil
   ○  sinkron kontak... 847

   slot terisi: 1/3   tambah lagi?   +  ya   ↵  cukup

```

**State: login_success**

```

  hubungin whatsapp

   ●  nyambung ke server wa
   ●  scan berhasil
   ●  kontak sinkron (847)

   slot 1 ✓  0812-xxxx-3456  terhubung

   udah nyambung. mau tambah nomor lagi?
   makin banyak nomor = rotator makin aman.

   +  tambah nomor   ↵  cukup, gas   q  nanti

```

**State: login_expired** (session expired, need re-login)

```

  hubungin whatsapp

  sesi lu udah expired. scan ulang ya.

         ┌─────────────────┐
         │                 │
         │    [QR CODE]    │
         │                 │
         └─────────────────┘

   ●  nyambung ke server wa
   ○  nunggu scan dari hp

   slot expired: 1  slot aktif: 2
   yang expired auto-pause, sisanya tetap jalan.

   ── sesi terakhir: 3 hari lalu ──

```

**State: login_failed** (network error, ban, etc)

```

  hubungin whatsapp

   ●  nyambung ke server wa
   ✗  gagal nyambung

   slot ini gagal. slot lain tetap jalan.
   wa server lagi bermasalah.
   coba lagi beberapa menit ya.

   1  coba lagi    2  ganti slot    q  kembali

```

Micro-interactions:
- `● ○ ○` animate sequential = progress feels alive
- QR dissolve pixel-by-pixel saat scan detected → checkmark bounce overshoot
- Contact sync counter live: `sinkron kontak... 847` = numbers moving = things happening
- On success: hold 800ms "udah nyambung" → auto transition = pause creates memory
- On failed: `✗` red, tapi pesan tetap santai = problem, bukan disaster

---
