# 🧠 WaClaw TUI — Neuroscienced Customer Journey

> Vertical-borderless. Micro-interactive. File-based.
> Lu cuma nonton. WaClaw yang kerja.
> Every pixel earns its place. Every pause has purpose.
> Validate early, fail loudly. Broken config = paused army.
> Assume nothing. Verify everything. Scrapped ≠ has WhatsApp.
> Follow-up = persistent, bukan spam. Timing + varian = persistence yang aman.
> Satu lisensi, satu device. Shared = stopped.
> Niche explorer = jalan keluar buat yang bingung. Search dulu, pilih nanti.
> Nerd stats = bawah sadar system. Toggle kapan aja.
> Versi baru = upgrade lisensi. Beda versi = beli baru.
> Ctrl+K = command palette. Apapun bisa, dari mana aja, 3 detik.

---

## 0. Filosofi Utama: ARMY IN THE BACKGROUND

```
WaClaw itu bukan satu asisten.
WaClaw itu army — satu worker per niche, jalan paralel 24/7.

Lu = jenderal. WaClaw = pasukan.
Lu tentuin strategi, mereka eksekusi.
Lu nggak perlu micromanage — tiap worker otonom.
Lu cuma di-interupt kalau ada yang perlu keputusan bos.
```

**Prinsip:**
- **Multi-niche by default.** Bukan satu pipeline. Banyak pipeline jalan bareng. Worker per niche.
- **Auto-run by default.** `waclaw` tanpa argument = langsung jalan semua niche, auto-pilot
- **Notification-first.** WaClaw yang nyari momen buat nanya lu, bukan lu yang nyari WaClaw
- **One-key decisions.** Setiap interrupt cuma butuh 1 tombol. `↵` = setuju, `s` = skip
- **Background is the default.** Scrapers jalan paralel, senders jalan paralel, lu cuma liat dashboard
- **Batch mindset.** Bukan kirim 1-1 kayak manual. Antri → batch → kirim → report. Rinse repeat.
- **Interrupt only when it matters.** Response masuk? Interrupt. Scrape selesai? Silent update ke dashboard. Error? Interrupt. Batch selesai? Notification ringan.
- **Validate early, fail loudly.** Broken config = paused army. Silent errors = invisible disaster. Setiap config error harus nongol ke permukaan SEKARANG, bukan pas runtime.
- **Every message rotatable.** Satu template = banyak varian. Kalau cuma 1 varian, WA bisa deteksi pola. Varian = perisai.
- **Assume nothing about numbers.** Nomor yang di-scrape dari Google Maps BUKAN berarti punya WhatsApp. Cek dulu, baru antri. Jangan buang slot kirim ke nomor kosong.
- **Anti-spam = anti-ban.** Spam sama WA = alasan ban. Jangan kirim ke orang yang bilang stop. Jangan kirim ulang ke yang udah dijawab. Anti-spam itu keselamatan, bukan etika doang.
- **Deal triggers are data, bukan tebakan.** User yang bilang "dah transfer" = deal. User yang bilang "oke gas" = deal. Closing trigger di-config, bukan di-hardcode.
- **Follow-up itu persistence, bukan spam.** Lead yang nggak response BUKAN lead mati — cuma belum waktunya. Tapi follow-up punya batas: max 3 pesan lifetime, jeda minimum 24 jam, varian berbeda tiap follow-up. Tanpa follow-up = leads terbuang. Tanpa batas = ban.
- **Lisensi = satu kunci, satu device.** WaClaw ngcek lisensi tiap startup. Kalau lisensi aktif di device lain, WaClaw berhenti. Bukan karena pelit — karena fair. Satu lisensi = satu army. Dua army = dua lisensi.
- **Niche explorer = discovery, bukan typing.** User nggak harus ngetik niche dari nol. WaClaw nyediain explorer buat browse kategori bisnis, search pakai WhatsApp Business Directory atau Google Maps categories, lalu pilih yang pas. Semua otomatis masuk ke config. User cuma pilih, WaClaw yang nulis niche.yaml.
- **Nerd stats = system vitals on demand.** Tekan backtick (`) buat toggle overlay RAM, CPU, goroutine count, dan DB size di mana aja. Bukan screen baru — overlay transparan di bawah layar. Mau liat? Toggle on. Nggak mau? Nggak keliatan. Default = hidden.
- **Versi baru = lisensi baru.** WaClaw ngecek versi terbaru pas startup. Kalau ada update, user bisa update langsung. Tapi beda major version = beda lisensi. Lisensi v1 nggak berlaku buat v2. Ini bukan pelit — ini product yang beda. v1 army = v1 lisensi. v2 army = v2 lisensi.
- **Ctrl+K = instant command.** Command palette yang bisa ngapa-ngapain dari mana aja. Search action, pindah screen, execute command — semua dari 1 tempat. Lu lagi di monitor, mau ke leads? Ctrl+K → ketik "lead" → ↵. Mau pause semua worker? Ctrl+K → ketik "pause" → ↵. Nggak perlu hafal semua shortcut. Nggak perlu navigate menu. 3 detik dari pikiran ke eksekusi.

**Keyboard itu hak istimewa, bukan kewajiban.**
Kalau lu nggak sentuh keyboard sama sekali selama 1 jam dan 3 niche tetap nyari leads, itu berarti WaClaw kerja dengan bener.

---

## 0.1 Design Language

**No borders. No boxes. Only space, weight, and motion.**

```
Hierarchy  = Brightness + Size + Motion
Separation = Vertical rhythm, never lines
Navigation = Muscle memory, never menus
Feedback   = Felt, not read
Language   = Netizen indo, bukan bahasa buku
Validation = Early, visible, actionable
Rotation   = Every message, every time — varian = perisai
Verification = Scrapped number ≠ WhatsApp number — cek dulu
Follow-up   = Persistent bukan spam — timing + varian = persistence aman
License     = Satu kunci satu device — shared = stopped
Explorer   = Search dulu pilih nanti — nggak usah ngetik dari nol
Nerd Stats = Vitals on demand — toggle kapan aja, hidden default
Versioning  = Major version = major license — v1 ≠ v2
CmdPalette = Ctrl+K from anywhere — search, navigate, execute, 3 seconds
```

Color system lives in `~/.waclaw/theme.yaml`. Everything else in `~/.waclaw/config.yaml`. Zero UI settings screens. You own your configs. Nerd stats toggle lives on backtick (`) key — always available, never in the way. Command palette lives on Ctrl+K — the universal backdoor to everything.

---

## 0.2 Bahasa TUI

**Semua text di TUI pake bahasa netizen indo. Bukan formal. Bukan inggris.**

Tone: santai, tapi nggak alay. Kayak ngobrol sama rekan kerja yang competent. Singkat, jelas, ga pake basa-basi. Error message juga santai — masalahnya serius, bahasanya nggak usah kaku.

---

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

### SCREEN 4: SCRAPE → ANTICIPATION

**State: scraping_active** (single niche view)

```

  lagi nyari...                           niche: web_developer

  target: cafe, gym, salon, toko bangunan
  area: kediri (15km)
  filter: tanpa website

  nemu ━━━━━━━━━━━━━━━━━━━━━━  156

   kopi nusantara          ──  no web  ✓
   gym fortress pro        ──  no web  ✓
   salon cantik alami      ──  has web  ✗ skip
   toko makmur jaya        ──  no web  ✓
   ...scanning


   lolos: 89 leads
   duplikat dibuang: 12
   baru masuk db: 67

   ━━━━━━━━━━━━━━━━━━━━━━━━━━  67 baru

   67 lead baru nunggu. waclaw auto-review ya.
   tekan ↵ buat review sendiri, atau biarin aja.

```

**State: scraping_multi_active** (MULTI NICHE — ini yang bikin nagih)

```

  lagi nyari...                           2 worker aktif

  ────────────────────────────────────────────────────

  web_developer                          ● scraping
    target: cafe, gym, salon, toko
    area: kediri (15km)
    nemu ━━━━━━━━━━━━  156
    lolos: 89 · baru: 67

  undangan_digital                       ● scraping
    target: wedding organizer, venue
    area: kediri + surabaya
    nemu ━━━━━━━━━━━━  94
    lolos: 61 · baru: 48

  ────────────────────────────────────────────────────

  total: 250 nemu · 150 lolos · 115 baru
  semua auto-review & auto-antri.
  lu nggak perlu ngapa-ngapain 👌

  tab  pindah niche    ↵  liat detail    q  balik

```

**State: scraping_multi_staggered** (berbeda fase per niche)

```

  lagi nyari...                           3 worker aktif

  ────────────────────────────────────────────────────

  web_developer                          ● scraping
    nemu: 156 · lolos: 89 · baru: 67

  undangan_digital                       ○ idle (berikutnya: 2j 14m)
    nemu: 94 · lolos: 61 · baru: 48
    batch terakhir: 23 menit lalu

  social_media_mgr                       ● baru mulai
    target: umkm, kafe lokal
    area: malang (10km)
    nemu: 12 · scanning...

  ────────────────────────────────────────────────────

  total aktif: 3 niche · 3 worker
  antrian kirim: 115 pesan

  tiap niche punya timing sendiri.
  lu cuma liat, waclaw yang ngatur semua.

```

**State: scrape_idle** (menunggu interval berikutnya)

```

  nyari

  scrape terakhir: 23 menit lalu
  nemu 67 lead baru
  scrape berikutnya: 5j 37m

  lagi nunggu interval. semua auto.
   ↵  scrape sekarang    q  balik

```

**State: scrape_empty** (zero results)

```

  nyari

  target: cafe, gym, salon, toko bangunan
  area: kediri (15km)

  ━━━━━━━━━━━━━━━━━━━━━━━━━━  0

  kosong nih. bisa jadi:
  - area lu kecil — coba perbesar radius
  - filter terlalu ketat — coba longgarkan
  - query terlalu spesifik — tambah variasi

  1  edit filter    2  ganti area    3  tambah query

```

**State: scrape_error** (scraper crash, network, etc)

```

  nyari

  ✗ scraper error — gagal nyambung google maps

  auto-retry 3 menit lagi.
  kalau masih gagal, waclaw bakal kasih tau.

   1  coba sekarang    q  balik

```

**State: scrape_gmaps_limited** (Google Maps rate limit / throttled)

```

  nyari

  ⏳ google maps throttle — permintaan kebanyakan

  gmaps lagi ngasih jeda. bukan error, cuma antri.
  auto-resume 12 menit lagi.

  scraper yang lain tetap jalan.
  cuma query ke gmaps yang dikasih jeda.

   1  coba sekarang    q  balik

```

**Variant: scrape_auto_approved** (auto-pilot mode, WaClaw auto-qualifies)

```

  lagi nyari...                           2 worker aktif

  ── web_developer ──
  nemu: 156 · lolos filter: 89
  duplikat: 12 · baru: 67
  auto-review: 67 masuk antrian, 8 skip

  ── undangan_digital ──
  nemu: 94 · lolos filter: 61
  duplikat: 3 · baru: 48
  auto-review: 48 masuk antrian, 10 skip

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  total batch ini: 115 lead baru masuk antrian

  pesan bakal dikirim jam kerja nanti.
  tiap niche pakai template sendiri.
  lu nggak perlu ngapa-ngapain 👌

```

**Variant: scrape_with_wa_validation** (setelah scrape, cek nomor WA dulu)

Ini beda sama sebelumnya. Sebelum langsung masuk antrian kirim, WaClaw CEK DULU apakah nomor yang di-scrape punya WhatsApp. Karena realitanya: **banyak nomor Google Maps itu telepon kantor, bukan WA.**

```

  lagi nyari...                           2 worker aktif

  ── web_developer ──
  nemu: 156 · lolos filter: 89
  duplikat: 12 · baru: 67

  cek nomor WA...
  ██████████████░░░░  78% (52/67 cek)

  ✓ punya WA:     34   → masuk antrian kirim
  ✗ bukan WA:     15   → ditandai, nggak dikirim
  ⏳ belum dicek:  18   → masih nunggu

  ── undangan_digital ──
  nemu: 94 · lolos filter: 61
  duplikat: 3 · baru: 48

  cek nomor WA...
  ██████████████████  100% (48/48 cek)

  ✓ punya WA:     29   → masuk antrian kirim
  ✗ bukan WA:     19   → ditandai, nggak dikirim

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  total siap kirim: 63 lead (bukan 115!)
  34 nomor bukan WA — nggak bakal buang slot kirim.

  cek WA jalan di background. nggak ngaretin proses.
  yang udah lolos langsung masuk antrian.

```

**Kenapa WA pre-validation:**
- Google Maps ngasih nomor telepon bisnis, BUKAN nomor WhatsApp
- Realitanya, 30-50% nomor yang di-scrape NGGAK punya WA
- Tanpa pre-validation: antri 115 → cuma 63 yang bisa dikirim → 52 slot kirim terbuang
- Dengan pre-validation: antri 63 → semua bisa dikirim → nggak ada slot terbuang
- Plus: nggak bikin WA flag karena nyoba kirim ke nomor yang bukan WA (ini salah satu trigger ban)
- Method: `check-registration` (cek ke WA server apakah nomor terdaftar) atau `send-silent` (kirim pesan silent type yang nggak keliatan)
- Catatan: pre-validation bisa lambat kalau banyak nomor, jadi jalan di background, lead yang udah validated langsung masuk antrian

**State: scrape_wa_validation_progress** (detail progress validasi WA)

```

  cek nomor WA                             niche: web_developer

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  total nomor dicek: 67
  progress: ██████████████░░░░  78%

  ✓ punya WA           34  (50.7%)
  ✗ bukan WA           15  (22.4%)
  ⏳ belum dicek        18  (26.9%)

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  yang punya WA langsung masuk antrian.
  yang bukan WA ditandai di database.
  nanti kalau mau, bisa dicoba ulang.

  estimasi selesai: 3 menit

  ↵  liat yang punya WA    s  skip tunggu    q  balik

```

Micro-interactions:
- Progress bar fills secara real-time = lu tau ini kerja, bukan hang
- `✓ punya WA` count naik = leads nyata, bukan angka kosong
- `✗ bukan WA` count dimmed = buang, bukan gagal — ini filter yang bikin data lebih bersih
- Percentage shown = data-driven, lu tau hit rate lu

**Variant: scrape_high_value_reveal** (slot machine effect when high-value lead found)

Ketika scraper nemu lead dengan skor 9+ (rating sempurna, no website, aktif IG, banyak reviews):

```

  lagi nyari...                           niche: web_developer

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

   ★  GRAND PALACE HOTEL     skor: 9.5

      hotel · jl. dhoho 45, kediri
      ⭐ 4.9 · 312 reviews · no web · aktif IG

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

   ★ ★ ★   JACKPOT LEAD!   ★ ★ ★

   langsung masuk prioritas antrian.
   template terbaik auto-dipilih.

  ──────────────────────────────────────

   nemu ━━━━━━━━━━━━━━━━━━━━━━  157

   ...scanning

```

Slot machine reveal animation:
1. Business name scrolls like slot reel (3 random names flash sebelum final) = 400ms
2. `★ ★ ★ JACKPOT LEAD! ★ ★ ★` scale up dari 0 dengan triple bounce overshoot = 600ms
3. Gold/amber color wave pulsed dari tengah ke luar = 300ms
4. Terminal bell sound (`\a`) = audio micro-reward
5. Score `9.5` glows amber 2x = perhatian fokus ke angka
6. Auto-settle: setelah 2 detik, jackpot section collapse ke normal list = lanjut kerja

**Variant: scrape_batch_complete** (cascade effect when batch finishes)

```

  lagi nyari...                           batch selesai!

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ✓ batch web_developer selesai

  156 nemu → 89 lolos → 67 baru
  duplikat: 12 · skip: 8

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  67 lead baru masuk antrian auto.
  pesan dikirim jam kerja nanti.

  ✓ batch undangan_digital selesai

  94 nemu → 61 lolos → 48 baru
  duplikat: 3 · skip: 10

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  total batch: 115 lead baru siap kirim
  berikutnya: 5j 37m

```

Batch cascade animation: tiap niche result "falls in" dari atas 200ms stagger, seolah hasil turun satu-per-satu. `✓` checkmark pulse hijau per niche = micro-reward per completion.

Micro-interactions:
- Business names appear one-by-one 30ms delay = slot machine effect
- `✓` marks pulse green briefly = micro-reward per hit
- `✗ skip` dimmed + slide faster = negatives don't linger
- Counter `156` increments live with scale bump = numbers feel tangible
- "waclaw auto-review ya" = explicit statement bahwa lu NGGAK PERLU ngapa-ngapain
- Auto-approved variant: the whole thing feels like a status report, bukan task list
- High-value reveal: slot machine scroll + jackpot bounce + terminal bell = peak excitement in an otherwise calm process
- Batch complete: cascade fall-in + sequential checkmarks = satisfying batch closure feel

**Prinsip utama:** WaClaw bilang "aku yang review, lu santai aja." User bisa interject kalau mau, tapi default = auto. High-value leads dapet momen khusus — lu NGERAIN ada sesuatu yang penting, tanpa perlu baca angka.

---

### SCREEN 5: LEAD REVIEW → CURATED (Optional, Manual Override)

**Screen ini OPTIONAL. WaClaw auto-review by default. Ini cuma muncul kalau user explicitly mau liat.**

**State: reviewing**

```

  review leads                              23 nunggu

  ────────────────────────────────────────────────────

  01  kopi nusantara
      cafe · jl. hasanuddin 23, kediri
      ⭐ 4.2 (87 reviews) · no website
      no instagram
      WA: ✓ terdaftar

      ── kirim ke ini? ──

      ↵  iya, antrikan       1-3  pilih varian
      s  skip                 d  detail
      x  skip & block

  ────────────────────────────────────────────────────

  diantrikan: 12  di-skip: 8  sisa: 23

  ↑↓  pindah   ↵  gas antrikan   q  selesai

```

**State: lead_detail** (press `d` untuk detail)

```

  kopi nusantara
  cafe · jl. hasanuddin 23, kediri

  ⭐ 4.2 · 87 reviews
  no website · no instagram
  ada foto: 12

  google search:
  - tidak ada website resmi
  - tidak ada link di maps
  - tidak ada di sosmed

  skor lead: 8/10 (sangat potensial)

  riwayat: belum pernah dikontak
  follow-up: belum — bisa follow up maks 3x

  ──────────────────────────────────────

  ↵  antrikan    s  skip    1-3  pilih template    q  balik

```

**State: template_preview** (press `1-3` untuk pilih varian)

```

  preview varian: ice_breaker/variant_1

  ──────────────────────────────────────

  Halo Kak Kopi Nusantara! 👋

  Tadi aku iseng cari Kopi Nusantara di
  Google, ternyata belum ada website
  resminya ya?

  Kebetulan aku Haikal (programmer), udah
  buatin preview desain web yang clean &
  profesional khusus buat Kopi Nusantara.

  Boleh izin kirim fotonya ke sini buat
  diintip? Gratis kok kak lihat-lihat
  aja dulu 😁

  ──────────────────────────────────────

  varian lain:  2  variant_2    3  variant_3

  ↵  pake ini    1-3  ganti varian    q  batal

```

**State: queue_complete** (semua lead sudah di-review)

```

  review selesai

  diantrikan: 18
  di-skip: 5
  di-block: 0

  pesan bakal dikirim pas jam kerja.
  waclaw yang ngatur timing, lu santai 👌

  ↵  liat dashboard    q  keluar

```

Micro-interactions:
- `↵ iya, antrikan` = one-key decision, pre-highlighted
- Slide-left on queue = decision momentum
- Fade-down on skip = different meaning per animation
- Template preview types out char-by-char = you see YOUR words forming
- "waclaw yang ngatur timing" = explicitly stating auto-pilot

---

### SCREEN 6: SEND → AUTO-PILOT

**Ini BUKAN screen yang lu operasikan. Ini status report.**

**State: sending_active** (jam kerja, multi-niche batch sending, WA rotator)

```

  lagi kirim pesan                    3 WA · 2 niche · 38 diantrikan (wa validated ✓)

  ────────────────────────────────────────────────────

  wa rotator

  📱 0812-xxxx-3456   ● aktif   4/6 jam ini   cooldown: 8m 12s
  📱 0813-xxxx-7890   ● aktif   3/6 jam ini   cooldown: 14m 37s
  📱 0857-xxxx-2345   ○ cooldown  2/6 jam ini   ready: 3m 05s

  ────────────────────────────────────────────────────

  ▸ web_developer (24 antri, semua wa validated)

  01  →  kopi nusantara       📱 slot-1
       ice_breaker: variant_2  ← rotasi!
       ━━━━━━━━━━━━━━━━━ mengirim...

  02     gym fortress pro     📱 slot-2
       ice_breaker: variant_1  ← rotasi!
       nunggu (berikutnya: 11m 23s)

  03     toko makmur jaya     📱 slot-1
       ice_breaker: variant_3  ← rotasi!
       nunggu (berikutnya: 24m 07s)

  ▸ undangan_digital (14 antri, semua wa validated)

  04  →  wedding bliss WO     📱 slot-3
       offer: variant_1
       nunggu (berikutnya: 8m 41s)

  ────────────────────────────────────────────────────

  rate: 9/18 per jam (3 nomor) · hari ini: 12/50
  kirim berikutnya: 08:41 (slot-3)

  semua lead udah dicek WA-nya. nggak ada slot terbuang.
  pesan + varian dirotasi ke 3 nomor. makin banyak varian = makin aman.
  p  pause    ↵  skip tunggu    tab  pindah niche    q  balik

```

**State: sending_paused** (user pause)

```

  lagi kirim pesan                           12 diantrikan

  ────────────────────────────────────────────────────

  ⏸  PAUSE — lu yang stop ini

  01  kopi nusantara       ✓ terkirim
  02  gym fortress pro     nunggu
  03     toko makmur jaya  nunggu

  rate: 4/6 per jam · hari ini: 12/50

  ↵  lanjut lagi    q  balik ke dashboard

```

**State: sending_off_hours** (di luar jam kerja)

```

  lagi kirim pesan

  ⏰ di luar jam kerja (09:00-17:00)
  sekarang: 21:47 wib

  pesan diantri, bakal kirim besok jam 09:00.
  lu nggak perlu ngapa-ngapain.

  1  tetap kirim (emergency)    q  balik

```

**State: sending_rate_limited** (limit jam ini capai)

```

  lagi kirim pesan

  ⏳ limit jam ini capai (6/6)
  bisa kirim lagi: 47 menit lagi

  waclaw ngatur sendiri, lu tinggal tunggu.
  nggak perlu refresh atau ngapa-ngapain.

   ↵  liat dashboard    q  balik

```

**State: sending_daily_limit** (limit harian capai)

```

  lagi kirim pesan

  📊 limit hari ini capai (50/50)
  sisa antrian: 38 pesan → dikirim besok
  (24 dari web_developer · 14 dari undangan_digital)

  hari ini:
  terkirim: 50
  respond: 7
  convert: 2

  bagus! istirahat dulu ya.
  waclaw lanjutin besok pagi auto, semua niche.

   ↵  liat dashboard    q  keluar

```

**State: sending_failed** (pesan gagal, number invalid, dll — RARE karena WA pre-validation)

```

  lagi kirim pesan

  ✗  gym fortress pro — gagal kirim
     nomor tidak terdaftar wa

  ⚠  ini seharusnya nggak terjadi — pre-validation harusnya udah filter.
  kemungkinan: nomor baru saja unreg WA, atau validasi miss.
  waclaw otomatis skip & mark sebagai invalid.
  nggak bakal dikirim ulang.

  1  coba lagi manual    s  skip aja    v  validasi ulang    q  balik

```

**State: sending_all_slots_down** (semua nomor WA putus)

```

  lagi kirim pesan

  ✗  semua nomor wa putus!

  📱 slot-1  ✗ putus
  📱 slot-2  ✗ putus
  📱 slot-3  ✗ putus

  semua sender mati. scraper tetap jalan.
  leads yang udah di-scrape tetap masuk database.
  cuma pengiriman yang pending sampai lu login ulang.

  ──────────────────────────────────────

  1  login ulang    2  login satu per satu    q  balik

```

**Variant: sending_with_response** (ada response masuk saat lagi kirim)

```

  lagi kirim pesan                           12 diantrikan

  ────────────────────────────────────────────────────

  💬 response masuk! kopi nusantara

  "iya kak, boleh lihat desainnya?"

  ↵  kirim offer    2  balas custom    s  nanti

  ────────────────────────────────────────────────────

  01  →  gym fortress pro
       mengirim...

  02     toko makmur jaya
       nunggu (berikutnya: 11m 23s)

```

Micro-interactions:
- `→` arrow hanya di active message = focus indicator
- `✓ terkirim` green pulse = completion reward
- `✓✓ sampai` second pulse different rhythm = escalating reward
- Countdown `11:23` live tick = time feels real
- "semuanya auto. lu cuma perlu nonton." = the whole point
- Response interrupt: gentle 200ms fade-in, not jarring = important, not alarming
- Rate limit bar fills organically → amber near limit = scarcity signal

---

### SCREEN 7: MONITOR → COMMAND CENTER

**Ini screen utama. Home base. Tempat lu ngabupaten.**

**State: live_dashboard** (multi-niche command center)

```

  waclaw                          ● wa nyambung (3 nomor) · 3 niche aktif

  ░░ 3 7 1 4 8 2 9 ░░ 5 1 6 3 8 ░░ 2 7 4 ░░

  ────────────────────────────────────────────────────

  wa rotator

  📱 slot-1  0812-xxxx-3456   ● aktif   4/6 jam
  📱 slot-2  0813-xxxx-7890   ● aktif   3/6 jam
  📱 slot-3  0857-xxxx-2345   ○ cooldown  ready: 3m

  worker pool

  web_developer      ● scraping   24 antri   3 respond
  undangan_digital   ● sending    14 antri   1 respond
  social_media_mgr   ○ idle       0 antri    0 respond

  ────────────────────────────────────────────────────

  hari ini                         minggu ini

  lead nemu           250          total leads     1247
  pesan terkirim       43          pesan terkirim   312
  response              8          response          41
  convert               3          convert           14

  conversion rate: 5.1%         hari terbaik: selasa

  ────────────────────────────────────────────────────

  aktivitas terbaru (semua niche)

  14:23  [web_dev]   kopi nusantara    respond    "iya kak, boleh lihat"
  14:01  [undangan]  wedding bliss     terkirim   ──
  13:47  [web_dev]   gym fortress      sampai     ──
  12:31  [web_dev]   salon cantik      respond    "makasih tapi udah ada"
  11:58  [undangan]  venue gardenia    respond    "berapa harganya?"
  11:15  [web_dev]   toko makmur       terkirim   ──
  09:42  [web_dev]   bengkel jaya      terkirim   ──

  ────────────────────────────────────────────────────

  1  leads    2  pesan    3  worker    4  template    5  anti-ban    6  follow-up    7  pengaturan

  r  refresh    s  scrape semua    q  keluar    `  nerd stats

```

**Ambient data rain effect:** Baris `░░ 3 7 1 4 8 2 9 ░░ 5 1 6 3 8 ░░ 2 7 4 ░░` — angka-angka faint yang scroll perlahan ke bawah di background. Seperti Matrix rain tapi subtle, pake `text_dim` color (hampir ga keliatan tapi bisa dibaca kalau lu cari). Angka-angka ini represent lead IDs yang lagi di-proses. Efeknya: lu NGERAIN ada data yang mengalir, system hidup, tanpa bikin clutter. Update setiap 5 detik, angka baru masuk dari kiri, angka lalu hilang ke kanan.

**Breathing stat numbers:** Angka-angka stats (`250`, `43`, `8`, `3`) punya subtle pulse — opacity naik turun 0.9 → 1.0 → 0.9 dalam 4 detik cycle. Sinkron tapi nggak serentak (offset 200ms per number). Feel: angka-angka ini hidup, bukan screenshot.

**State: idle_background** (army lagi kerja, lu santai)

```

  waclaw                          ● auto-pilot · 3 niche

  ░░ 4 1 7 3 ░░ 6 2 8 ░░ 1 5 9 4 ░░

  ────────────────────────────────────────────────────

  army lu lagi kerja nih. semua auto.

  web_developer      ● scrape 5j 37m · kirim 11m 23s
  undangan_digital   ● scrape 2j 14m · kirim 8m 41s
  social_media_mgr   ● scrape 4j 02m · idle (malam)

  wa                 ● nyambung (3 nomor)
  antrian total      38 pesan
  response hari ini  8

  lu bisa minim window, waclaw tetap kerja.
  3 worker nyari leads + 3 nomor rotasi kirim.
  nanti waclaw notif kalau ada yang penting.

  ↵  liat detail    q  keluar

```

**State: dashboard_night** (di luar jam kerja, semua idle)

```

  waclaw                          ○ mode malam · 3 niche

  ────────────────────────────────────────────────────

  jam kerja: 09:00-17:00 wib
  sekarang: 22:15

  semua sender pause sampai besok pagi.
  scraper tetap jalan di background.
  3 worker auto gas lagi jam 09:00.

  ringkasan hari ini:
  43 terkirim · 8 response · 3 convert

  per niche:
  web_developer     28 terkirim · 5 respond · 2 convert
  undangan_digital  15 terkirim · 3 respond · 1 convert

  tips: selasa jam 10 itu waktu terbaik lu kirim pesan

   ↵  liat detail    q  keluar

```

**State: dashboard_error** (ada masalah)

```

  waclaw                          ✗ ada masalah · 3 niche

  ────────────────────────────────────────────────────

  ✗  wa putus (slot-1) — slot-2 & slot-3 tetap jalan
  ●  scraper tetap jalan (nggak butuh wa)
  ●  database ok

  2 nomor masih aktif. kirim tetap jalan.
  cuma slot-1 yang pending sampai nyambung.
  auto-reconnect nyala.

  1  login ulang    r  cek status    q  keluar

```

**State: dashboard_empty** (baru mulai, belum ada data)

```

  waclaw                                   ● wa nyambung

  ────────────────────────────────────────────────────

  belum ada data. mulai dari sini:

  1  scrape sekarang    2  atur niche    3  kirim pesan

  ────────────────────────────────────────────────────

  tips: mulai dari scrape dulu ya.
  nanti setelah ada lead, sisanya auto.

```

**Variant: dashboard_with_pending_responses** (ada response yang belum di-handle)

```

  waclaw                          ● wa nyambung · 3 niche

  ────────────────────────────────────────────────────

  ⚡ 5 response belum dibalas!

  [web_dev]   kopi nusantara     "iya kak, boleh lihat"
  [web_dev]   salon cantik       "berapa harganya?"
  [undangan]  wedding bliss      "boleh kirim contohnya?"
  [undangan]  venue gardenia     "tertarik, detailnya?"
  [web_dev]   bengkel jaya       "boleh info lebih lanjut?"

  ────────────────────────────────────────────────────

  1  leads    2  pesan    3  worker    4  template    5  anti-ban    6  follow-up    7  pengaturan

  ↵  balas satu-satu    1  auto-kirim offer ke semua
  2  auto per niche     (offer beda per niche)

  ────────────────────────────────────────────────────

  hari ini: 43 terkirim · 8 response · 3 convert

```

Micro-interactions:
- `● wa nyambung` pulses 3s cycle = alive
- `✗ ada masalah` only red thing on screen = auto-draw attention
- `⚡ 3 response belum dibalas!` amber flash = urgency tanpa alarm
- Numbers increment with highlight = movement catches eye
- "respond" rows get warm amber flash 2s = positive signal
- Rejection in neutral color = no negative visual weight
- "lu bisa minim window" = explicitly telling user they can leave
- "1 auto-kirim offer ke semua" = one-key bulk action, auto-pilot mindset
- Ambient data rain: faint scrolling numbers = system alive, data flowing, subliminal activity
- Breathing stats: subtle opacity pulse on numbers = living dashboard, bukan screenshot
- Data rain pauses when lu interact, resumes after 10s idle = focus mode vs ambient mode

---

### SCREEN 8: RESPONSE → REWARD

**Ini interrupt. WaClaw yang panggil lu, bukan sebaliknya.**

**State: response_positive** (jelas tertarik)

```

  💬 ada yang balas!

  ────────────────────────────────────────────────────

  kopi nusantara · cafe · kediri

  "iya kak, boleh lihat desainnya?"

  ────────────────────────────────────────────────────

  ↵  kirim offer    2  balas custom    3  nanti

```

**State: response_curious** (nanya-nanya, interested tapi ragu)

```

  💬 ada yang balas!

  ────────────────────────────────────────────────────

  salon cantik · salon · kediri

  "berapa harganya kak?"

  ────────────────────────────────────────────────────

  ↵  kirim info harga    2  balas custom    3  nanti

```

**State: response_negative** (tidak tertarik)

```

  💬 ada yang balas

  ────────────────────────────────────────────────────

  gym fortress · gym · kediri

  "makasih tapi udah ada yang handle"

  ────────────────────────────────────────────────────

  1  mark invalid    2  follow up nanti    ↵  skip aja

```

**State: response_stop_detected** (orang bilang berhenti — AUTO-ADD ke do_not_contact)

```

  🛑 ada yang minta berhenti

  ────────────────────────────────────────────────────

  bengkel jaya · bengkel · kediri

  "jangan hubungi lagi ya"

  ────────────────────────────────────────────────────

  ⚠  ini match closing_triggers.stop
  waclaw auto-tambah ke do_not_contact.yaml
  nomor ini nggak bakal dikontak lagi.

  ↵  setuju, block    2  block semua niche    3  batal, follow up aja

```

Micro-interactions:
- `🛑` red shield icon = berbeda dari response biasa, ini SERIUS
- Auto-action sudah taken: "waclaw auto-tambah ke do_not_contact" = lu cuma konfirmasi
- `2 block semua niche` = kalau nomor ini ada di multi-niche, block di semua
- `3 batal` = override, tapi warning: ini ngelawan anti-spam guard

**State: response_deal_detected** (CLOSING TRIGGER — auto-detected deal!)

```

  💬 ada yang balas — 🎯 DEAL TERDETEKSI!

  ────────────────────────────────────────────────────

  kopi nusantara · cafe · kediri

  "wah makasi dah transfer kak"

  ────────────────────────────────────────────────────

  ⚡ closing trigger: "dah transfer" → auto-deal

  waclaw deteksi ini sebagai deal.
  tekan ↵ buat konfirmasi, atau s buat koreksi.

  ↵  konfirmasi deal!    s  bukan deal    2  balas dulu

```

Micro-interactions:
- `🎯 DEAL TERDETEKSI!` amber flash + subtle gold shimmer = ini penting, tapi bukan sebesar manual conversion
- Auto-classification: WaClaw udah baca pattern-nya, lu cuma verify
- `s bukan deal` = false positive aman, bisa dikoreksi
- `2 balas dulu` = kadang mau bilang makasih dulu sebelum mark deal

**State: response_hot_lead_detected** (HOT LEAD — closing trigger hot_lead)

```

  💬 ada yang balas — 🔥 HOT LEAD!

  ────────────────────────────────────────────────────

  salon cantik · salon · kediri

  "berapa harganya kak?"

  ────────────────────────────────────────────────────

  ⚡ hot lead trigger: "berapa harga" → auto-prioritize

  waclaw naikin prioritas lead ini.
  tekan ↵ buat kirim offer langsung.

  ↵  kirim offer sekarang    2  balas custom    3  nanti

```

Micro-interactions:
- `🔥 HOT LEAD!` warm amber pulse = ini lead yang paling potensial
- Auto-prioritize: lead ini naik ke atas antrian
- Closing trigger bikin WaClaw proaktif, bukan cuma reactive

**State: response_maybe** (tidak jelas, bisa jadi tertarik)

```

  💬 ada yang balas

  ────────────────────────────────────────────────────

  toko makmur · toko bangunan · kediri

  "oh ya? bisa info lebih lanjut?"

  ────────────────────────────────────────────────────

  ↵  kirim offer    2  kirim info    3  balas custom

```

**State: response_auto_reply** (detected bot/auto-reply)

```

  💬 ada yang balas (auto)

  ────────────────────────────────────────────────────

  bengkel jaya · bengkel · kediri

  "Terima kasih telah menghubungi kami.
   Kami akan membalas pesan Anda secepatnya."

  ────────────────────────────────────────────────────

  ini auto-reply. skip aja ya.

  ↵  skip    1  tetap follow up    q  balik

```

**State: offer_preview** (sebelum kirim offer)

```

  kirim offer ke: kopi nusantara

  ────────────────────────────────────────────────────

  Halo Kak Kopi Nusantara! 👋

  Tadi aku iseng cari Kopi Nusantara di
  Google, ternyata belum ada website
  resminya ya?

  Kebetulan aku Haikal (programmer), udah
  buatin preview desain web yang clean &
  profesional khusus buat Kopi Nusantara.

  Boleh izin kirim fotonya ke sini buat
  diintip? Gratis kok kak lihat-lihat
  aja dulu 😁

  ────────────────────────────────────────────────────

  ↵  kirim    2  ganti template    3  edit dulu

```

**State: response_multi_queue** (multiple responses arrive simultaneously)

```

  💬 3 response masuk barengan!

  ────────────────────────────────────────────────────

  ▸ 01  kopi nusantara
        cafe · kediri
        "iya kak, boleh lihat desainnya?"
        ● positif

    02  wedding bliss WO
        WO · kediri
        "boleh kirim contohnya?"
        ○ penasaran

    03  bengkel jaya
        bengkel · kediri
        "Terima kasih telah menghubungi..."
        ○ auto-reply

  ────────────────────────────────────────────────────

  ↵  proses satu-satu    1  auto-offer semua yang positif
  2  auto per response type

```

Multi-queue micro-interactions:
- Stack appears with 3 items sliding in sequentially 100ms stagger = triage feel
- `● positif` / `○ penasaran` / `○ auto-reply` auto-classified badge = WaClaw udah sort buat lu
- Pre-highlighted `↵ proses satu-satu` = most safe default
- `1 auto-offer semua yang positif` = one-key bulk action, auto-pilot mindset
- Each response type has distinct tint: positive = warm amber, curious = neutral, auto-reply = dimmed
- After processing one, next slides up into focus = queue momentum

**State: conversion** (THE MOMENT — deal closed)

Ini. Ini dia. Puncak. The whole point. Ini screen yang bikin lu ketagihan.

```

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


   ★  ★  ★   D E A L !   ★  ★  ★


  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  kopi nusantara · cafe · kediri

  convert dari ice breaker → offer → deal
  waktu: 2 hari 4 jam

  ────────────────────────────────────────────────────

    ╔═══════════════════════════════════╗
    ║  ✦  ·  ★  ☆  ✧  ·  ★  ✦  ☆  ║
    ║  ☆  ✦  ·  ★  ✧  ☆  ·  ✦  ★  ║
    ║  ★  ☆  ✦  ·  ★  ✧  ☆  ✦  ·  ║
    ║  ✧  ★  ☆  ✦  ·  ★  ☆  ✧  ★  ║
    ╚═══════════════════════════════════╝

  ────────────────────────────────────────────────────

  🏆 conversion ke-3 minggu ini
  🏆 total revenue minggu ini: rp 7.5jt

   ↵  mark as converted

```

**CONVERSION — Full Drama Sequence:**

Phase 1: SHOCK (0-200ms)
- Seluruh screen flash PUTIH full-brightness — semua text hilang sesaat
- Terminal bell `\a` double-tap — dua kali bunyi, cepet
- Screen shock = "berhenti apapun yang lu lakuin"

Phase 2: REVEAL (200-800ms)
- White fades → `★ ★ ★ D E A L ! ★ ★ ★` scales dari 0 ke 1.3 (overshoot) → settle ke 1.0
- Spaced-out letters = monumental, bukan biasa
- Top/bottom `━━━` borders draw inward dari ujung, meet di tengah = framing the moment
- Particle cascade: `✦ · ★ ☆ ✧` scatter dari tengah ke luar, 40 particles, 600ms lifetime
- Particle colors cycle: gold → amber → white → fade
- Color wave: background pulses dari accent → success → normal dalam 400ms sweep

Phase 3: CONTEXT (800ms-1500ms)
- Business name + timeline fade in dari bawah
- `trophy 🏆` icon bounces in dari kanan dengan overshoot
- Revenue number glows gold 3x pulse = "ini uang beneran"
- "conversion ke-3 minggu ini" = streak context, bikin nagih

Phase 4: SETTLE (1500ms+)
- Particles dissolve
- Screen settles ke normal tones
- `↵ mark as converted` fades in = gentle return ke reality
- Hold 800ms sebelum keyboard accepted = lu HARUS ngerasain momen ini dulu

**Sound design suggestion:** Terminal bell `\a` saat shock phase. Kedua bunyi: BING BING. Kalau terminal support bell, ini jadi audio landmark. Kalau nggak, visual shock lebih dari cukup.

**Why this is the most important screen:** Ini dia. Ini satu-satunya alasan WaClaw ada. Semua scraping, template, rotator, anti-ban — semuanya cuma jalan ke screen ini. Kalau conversion screen ga bikin lu seneng, seluruh app gagal. Jadi kita ga hemat. Kita ALL OUT. Full flash. Full particles. Full drama. Karena MOMEN INI LAYAK DIPERINGATI.

**The neuroscience:** Different response types = different animations = different felt experiences. Your brain categorizes them automatically. You don't think "this is a positive response" — you FEEL it because the screen feels warm. Conversion? Conversion feels like WINNING. Karena emang menang.

---

### SCREEN 9: LEADS DATABASE → ARCHIVE

**State: leads_list** (semua lead)

```

  database leads                          total: 847

  ────────────────────────────────────────────────────

  filter: semua    / buat cari...

  ▸ baru          23
  ▸ ice breaker    128
  ▸ follow-up 1    34
  ▸ follow-up 2    12
  ▸ dingin          8   (2x follow-up, belum jawab)
  ▸ respond        41
  ▸ offer terkirim 89
  ▸ convert        14
  ▸ gagal           5
  ▸ di-block        3

  ────────────────────────────────────────────────────

  terbaru:

  kopi nusantara      respond    ⭐ 4.2
  gym fortress        ice breaker ⭐ 3.8
  salon cantik        respond    ⭐ 4.5

  ↑↓  pindah    ↵  liat detail    /  cari    q  balik

```

**State: leads_filtered** (filter by status)

```

  database leads                      respond: 41

  ────────────────────────────────────────────────────

  01  kopi nusantara
      cafe · kediri · ⭐ 4.2
      ice breaker: kemarin 14:23
      response: "iya kak, boleh lihat"
      → belum dikirim offer

  02  salon cantik
      salon · kediri · ⭐ 4.5
      ice breaker: 2 hari lalu
      response: "berapa harganya?"
      → belum dibalas

  ↑↓  pindah    ↵  liat    /  cari    1  kirim offer semua    q  balik

```

**State: lead_full_detail** (single lead complete view)

```

  kopi nusantara

  ────────────────────────────────────────────────────

  cafe · jl. hasanuddin 23, kediri
  ⭐ 4.2 · 87 reviews
  no website · no instagram
  ada foto: 12

  ────────────────────────────────────────────────────

  skor: 8/10

  timeline:

  kemarin 09:15   ice breaker terkirim
         14:23   response masuk
                 "iya kak, boleh lihat desainnya?"

  ────────────────────────────────────────────────────

  status: respond → menunggu offer

  1  kirim offer    2  balas custom    3  nanti
  4  mark convert   5  archive         6  block

  q  balik

```

**Variant: lead_follow_up_due** (sudah dikirimi ice breaker, belum jawab, waktunya follow-up)

```

  kopi nusantara

  cafe · jl. hasanuddin 23, kediri
  ⭐ 4.2 · 87 reviews · no website

  skor: 8/10

  ice breaker: kemarin 09:15 (1x dikontak)
  response: belum ada
  follow-up berikutnya: hari ini (jadwal otomatis)

  1  kirim follow-up    2  skip    3  tandai dingin    q  balik

```

**Variant: lead_cold** (2x follow-up belum jawab — lead dingin)

```

  kopi nusantara                                    ❄ DINGIN

  cafe · jl. hasanuddin 23, kediri
  ⭐ 4.2 · 87 reviews · no website

  skor: 8/10

  ice breaker: 5 hari lalu
  follow-up 1: 3 hari lalu
  follow-up 2: 1 hari lalu
  response: belum ada (2x follow-up, masih diam)

  lead ini udah 2x di-follow up, belum jawab.
  bisa coba 1x lagi (terakhir), atau tandai dingin.
  lead dingin nggak bakal dikontak lagi auto.

  1  follow-up terakhir (ke-3)    2  tandai dingin    3  archive    q  balik

```

  Micro-interactions:
  - `❄ DINGIN` badge dim blue = bukan mati, tapi dingin. Masih bisa di-revive manual.
  - Follow-up count (`1x dikontak`, `2x follow-up`) = lu tau seberapa persistent
  - "jadwal otomatis" = WaClaw yang ngatur timing, lu cuma approve
  - Lead dingin: dimmed tapi nggak hilang = masih ada, tapi nggak ngerepotin

**Variant: lead_never_contacted** (baru masuk, belum dikontak)

```

  kopi nusantara

  cafe · jl. hasanuddin 23, kediri
  ⭐ 4.2 · 87 reviews · no website

  skor: 8/10

  belum pernah dikontak

  1  kirim ice breaker    2  skip    q  balik

```

**Variant: lead_converted** (sudah deal)

```

  kopi nusantara                                          ✓ DEAL

  ────────────────────────────────────────────────────

  cafe · jl. hasanuddin 23, kediri
  niche: web_developer
  ⭐ 4.2 · 87 reviews

  timeline:

  28 apr 09:15   ice breaker terkirim
  28 apr 14:23   response: "iya kak, boleh lihat"
  28 apr 14:25   offer terkirim
  30 apr 10:12   response: "oke gas, kita mulai kapan?"
  30 apr 10:15   ✓ mark convert

  durasi: 2 hari 1 jam
  template: direct-curiosity
  worker: web_developer

  🏆 conversion ini bawa pulang rp 2.5jt

  q  balik

```

Micro-interactions:
- Status badges use color: baru=accent, respond=amber, convert=success, gagal=dimmed
- Timeline entries appear sequentially 100ms stagger = narrative feel
- "belum dikirim offer" in amber = actionable, not forgotten
- Converted lead: green `✓ DEAL` badge + gold shimmer on revenue number = earned celebration

---

### SCREEN 10: TEMPLATE MANAGER → ARMORY

**State: template_list**

```

  template pesan

  ────────────────────────────────────────────────────

  niche: web_developer

  ice breaker:
  ▸ default              "halo kak, apakah ini..."

  offer:
  ▸ direct-curiosity     "tadi aku iseng cari..."    ★ recomended
  ▸ pattern-interrupt    "permisi kak! aku haikal..."
  ▸ admin-bypass         "halo kak admin!..."

  ────────────────────────────────────────────────────

  ↑↓  pilih    ↵  preview    n  baru    e  edit    q  balik

```

**State: template_preview**

```

  preview: direct-curiosity                    ★ recommended

  ────────────────────────────────────────────────────

  Halo Kak {{.Title}}! 👋

  Tadi aku iseng cari {{.Title}} di Google,
  ternyata belum ada website resminya ya?

  Kebetulan aku Haikal (programmer), udah
  buatin preview desain web yang clean &
  profesional khusus buat {{.Title}}.

  Boleh izin kirim fotonya ke sini buat
  diintip? Gratis kok kak lihat-lihat
  aja dulu 😁

  ────────────────────────────────────────────────────

  placeholder:
  {{.Title}} → kopi nusantara
  {{.Category}} → cafe
  {{.Address}} → jl. hasanuddin 23, kediri

  ↵  pake ini    e  edit di file    q  balik

```

**State: template_edit_hint** (redirect ke file editor)

```

  edit template

  template disimpan sebagai file teks.
  buka di editor favorit lu:

    ~/.waclaw/niches/web_developer/offer_1.md

  setelah save, tekan r buat reload.

  r  reload    q  balik

```

**State: template_validation_error** (broken template detected)

```

  ✗ template error

  ────────────────────────────────────────────────────

  niche: fotografer

  ice_breaker.md
  ✗  file kosong — isi dulu pesan pembukanya
     tapi harus ada minimal {{.Title}} placeholder

  offer_2.md
  ✗  placeholder gak dikenali: {{.NamaToko}}
     placeholder yang tersedia: {{.Title}}, {{.Category}},
     {{.Address}}, {{.Rating}}, {{.Reviews}}
  ✗  baris 5: encoding error (bukan UTF-8)
     kemungkinan file disave pake encoding salah

  ────────────────────────────────────────────────────

  worker fotografer di-pause sampai template diperbaiki.
  ice_breaker WAJIB ada. offer yang error cuma di-skip.

  1  buka file    2  liat placeholder yang tersedia
  r  reload    q  balik

```

Micro-interactions:
- `★ recommended` badge on best-performing template = social proof in your own data
- Preview fills placeholders with sample lead data = you see what they'll see
- "edit di file" = explicit redirect, no in-app editor needed
- On reload: template list refreshes with brief highlight on changed items
- `template_validation_error`: broken template name blinks red 2x = instant visibility
- "1 buka file" shortcut = zero friction fix
- "ice_breaker WAJIB ada" = explicit severity level, different dari offer error

---

### SCREEN 11: WORKERS → PIPELINE VISUALIZER

**Ini jantungnya. Live view of all workers running in background.**

**State: workers_overview**

```

  worker pool                               3 aktif · 0 idle

  ────────────────────────────────────────────────────

  web_developer
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  scrape ━━━━━━━━━━━━░░░░░  67%  (2/3 query selesai)
  review ━━━━━━━━━━━━━━━━━━  done (89 lolos)
  antri  ━━━━━━━━━━━░░░░░░  24 pesan
  kirim  ━━━━░░░░░░░░░░░░░  3/6 jam ini
  area: kediri (15km)
  template: direct-curiosity ★

  undangan_digital
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  scrape ━━━━━━━━━━━━━━━━━━  done (61 lolos)
  review ━━━━━━━━━━━━━━━━━━  done (48 lolos)
  antri  ━━━━━━░░░░░░░░░░░  14 pesan
  kirim  ━━━━░░░░░░░░░░░░░  1/6 jam ini
  area: kediri + surabaya
  template: undangan-offer ★

  social_media_mgr
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  scrape ━━━━━━━━━━━━━━━━━━  done (23 lolos)
  review ━━━━━━━━━━━━━━━━━━  done (18 lolos)
  antri  ━━━━━━━━━━━━━━━━━━  18 pesan
  kirim  ○ idle (malam, mulai besok 09:00)
  area: malang (10km)
  template: smm-pitch ★

  ────────────────────────────────────────────────────

  total pipeline: 250 nemu → 150 lolos → 56 antri → 4 terkirim
  semua jalan paralel. lu cuma nonton.

  ↑↓  pilih worker    ↵  liat detail    n  tambah niche    q  balik

```

**State: worker_detail** (deep dive satu worker)

```

  worker: web_developer

  ────────────────────────────────────────────────────

  pipeline:

  scrape     ━━━━━━━━━━━━━━━░░░  82%  aktif
             query: "cafe di kediri"  ✗ done
             query: "gym di kediri"   ● scanning
             query: "salon di kediri" ○ waiting

  qualify    ━━━━━━━━━━━━━━━━━━  done
             156 nemu → 89 lolos (57%)
             67 duplikat dibuang
             5 skip (rating rendah)

  antri      ━━━━━━━━━━━░░░░░░  24 pesan
             3 ice breaker
             21 offer (auto setelah respond)

  kirim      ━━━━░░░░░░░░░░░░  3/6 jam ini
             berikutnya: 11m 23s
             hari ini: 12/50

  ────────────────────────────────────────────────────

  performa niche ini:
  response rate: 16%
  conversion rate: 4.6%
  avg waktu respond: 3.2 jam

  ────────────────────────────────────────────────────

  1  pause worker    2  force scrape    3  liat leads
  q  balik

```

**State: worker_add_niche** (tambah niche baru ke pool)

```

  tambah worker baru

  pilih niche buat ditambah ke pool:
  (worker baru langsung jalan setelah dipilih)

  1  ☐ fotografer           jasa foto & portfolio
  2  ☐ akuntan              jasa pajak & keuangan umkm
  3  ☐ custom               bikin niche sendiri

  ↵  tambah    q  batal

  makin banyak niche, makin banyak leads.
  tiap worker jalan independen, nggak saling ganggu.

```

**State: worker_paused** (worker yang di-pause manual)

```

  worker: social_media_mgr              ⏸ PAUSED

  ────────────────────────────────────────────────────

  lu yang pause ini. alasan:

  1  lanjutin
  2  hapus worker ini
  3  liat leads yang udah dikumpulin

  18 leads udah di database.
  kalau lu lanjutin, auto gas lagi.

```

Micro-interactions:
- Pipeline bars fill in real-time = lu ngerain prosesnya hidup
- Each worker row breathes independently = paralel, bukan sequential
- When a pipeline stage completes: bar fills → holds 400ms → section collapses = stage done, next stage gets focus
- "semua jalan paralel" = explicitly stating the army metaphor
- Worker add: instant spin-up animation = new worker born, starts working immediately
- Worker pause: section dims but stays visible = paused, not dead

---

### SCREEN 12: ANTI-BAN → SHIELD

**Ini perisai lu. Monitor semua yang bikin lu aman dari ban.**

**State: shield_overview** (semua aman)

```

  perisai anti-ban                          🛡️ semua aman

           ╱╲
          ╱  ╲
         ╱ ░░ ╲       health score
        ╱ ░░░░ ╲      ━━━━━━━━━━  92/100
       ╱ ░░░░░░ ╲
      ╱──────────╲
     ╱  ■ ■ ■ ■   ╲    3 slot aktif
    ╱──────────────╲
   ╱                  ╲
  ──────────────────────

  ────────────────────────────────────────────────────

  wa rotator                                3 nomor

  📱 slot-1  0812-xxxx-3456   ● aktif
     ━━━━━━━━━━━━━━━━━━━━━━  4/6 jam   cooldown: 8m 12s
     hari ini: 12 terkirim · 0 warning
     status: sehat ✓

  📱 slot-2  0813-xxxx-7890   ● aktif
     ━━━━━━━━━━━━━━━━━━━━━━  3/6 jam   cooldown: 14m 37s
     hari ini: 9 terkirim · 0 warning
     status: sehat ✓

  📱 slot-3  0857-xxxx-2345   ○ cooldown
     ━━━━━━━━━━━━░░░░░░░░░░  2/6 jam   ready: 3m 05s
     hari ini: 6 terkirim · 0 warning
     status: sehat ✓

  ────────────────────────────────────────────────────

  rate limiting

  per slot:  6/jam (aktif 3 slot = 18/jam total)
  per hari:  50 total (terpakai: 27)
  per nomor: jeda minimum 8 menit antar pesan
  per lead:  1 pesan per 24 jam (nggak spam)

  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  daily budget: 27/50 ━━━━━━━━━━━━░░░░░░░░  54%

  ────────────────────────────────────────────────────

  jam kerja guard

  zona: wib (asia/jakarta)
  jam kirim: 09:00-17:00 (8 jam)
  jam scrape: 24/7 (nggak keliatan sama wa)
  sekarang: 14:23 ✓ dalam jam kerja

  ────────────────────────────────────────────────────

  pattern guard

  template rotasi:    aktif (3 varian ice_breaker, 3 varian offer per niche)
  variasi waktu:      ±30% random delay
  variasi pesan:      placeholder dynamic per lead
  emoji variation:    aktif
  paragraf shuffle:   aktif

  ────────────────────────────────────────────────────

  spam guard

  per lead:          1 pesan per 24 jam (nggak spam)
  per lead lifetime: max 3 pesan (ice_breaker + 2x follow up)
  do-not-contact:    12 nomor di block list
  stop detection:    aktif (auto-add ke block list kalau match closing_triggers.stop)
  duplicate guard:   aktif (nggak kirim ke nomor yang udah di-niche lain)
  re-contact delay:  7 hari (setelah response tanpa deal, tunggu 7 hari sebelum follow up)

  ────────────────────────────────────────────────────

  ban risk score

  🟢  rendah — semua indikator aman

  indikator:
  ✓  pengiriman merata antar nomor
  ✓  cooldown cukup antar pesan
  ✓  template bervariasi (ice_breaker + offer keduanya rotate)
  ✓  nggak ada nomor yang kelebihan beban
  ✓  jam kerja dipatuhi
  ✓  spam guard aktif (nggak ada lead yang kelebihan kontak)
  ✓  do-not-contact list dihormati

  ↵  liat detail nomor    r  refresh    q  balik

```

**ASCII Shield Art — Dynamic Based on Health:**

Shield visual berubah berdasarkan aggregate health score:

```
  Health 90-100 (SEHAT):
           ╱╲
          ╱  ╲
         ╱ ░░ ╲        solid, fill penuh
        ╱ ░░░░ ╲       warna: success (hijau)
       ╱ ░░░░░░ ╲
      ╱──────────╲

  Health 50-89 (WARNING):
           ╱╲
          ╱    ╲
         ╱ ░░   ╲       ada celah, fill partial
        ╱ ░░░░   ╲      warna: warning (amber)
       ╱ ░░░░    ╲
      ╱──────────╲

  Health <50 (BAHAYA):
           ╱╲
          ╱    ╲
         ╱      ╲       retak, fill minimal
        ╱  ░░    ╲      warna: danger (merah)
       ╱  ░░     ╲      cracks: ╳ di shield
      ╱───╳──────╲

  Repair animation:
  Saat health score naik, shield fill bertambah
  dari bawah ke atas 50ms per poin = perbaikan
  visual yang lu bisa liat terjadi.
```

**State: shield_warning** (ada indikator warning)

```

  perisai anti-ban                          ⚠️ ada warning

           ╱╲
          ╱    ╲
         ╱ ░░   ╲       health score
        ╱ ░░░░   ╲      ━━━━━━━━━━  71/100
       ╱ ░░░░    ╲
      ╱──────────╲

  ────────────────────────────────────────────────────

  wa rotator                                3 nomor

  📱 slot-1  0812-xxxx-3456   ⚠️ warning
     ━━━━━━━━━━━━━━━━━━━━━━  5/6 jam   cooldown: 23m 41s
     hari ini: 14 terkirim · 1 warning
     ⚠  terlalu banyak jam ini (5/6)
     auto-reduce: slot-1 dikurangi beban

  📱 slot-2  0813-xxxx-7890   ● aktif
     ━━━━━━━━━━━━━━━━━━━━━━  3/6 jam   cooldown: 14m 37s
     hari ini: 9 terkirim · 0 warning
     status: sehat ✓

  📱 slot-3  0857-xxxx-2345   ● aktif
     ━━━━━━━━━━━━━━━━━━━━━━  2/6 jam   cooldown: 3m 05s
     hari ini: 6 terkirim · 0 warning
     status: sehat ✓

  ────────────────────────────────────────────────────

  ban risk score

  🟡  sedang — slot-1 kelebihan beban

  waclaw otomatis pindah beban ke slot-2 & slot-3.
  lu nggak perlu ngapa-ngapain. auto-adjust.

  ↵  liat detail nomor    r  refresh    q  balik

```

**State: shield_danger** (nomor kena flag / ban risk tinggi)

```

  perisai anti-ban                          ✗ BAHAYA

           ╱╲
          ╱    ╲
         ╱      ╲       health score
        ╱  ░░    ╲      ━━━━━━━━━━  38/100
       ╱  ░░     ╲
      ╱───╳──────╲     RETAK!

  ────────────────────────────────────────────────────

  wa rotator                                3 nomor

  📱 slot-1  0812-xxxx-3456   ✗ FLAGGED
     ⚠  nomor ini kena flag sama wa
     kemungkinan: pengiriman terlalu agresif
     action: auto-pause slot-1

  📱 slot-2  0813-xxxx-7890   ● aktif
     ━━━━━━━━━━━━━━━━━━━━━━  4/6 jam   cooldown: 8m 12s
     menggantikan beban slot-1

  📱 slot-3  0857-xxxx-2345   ● aktif
     ━━━━━━━━━━━━━━━━━━━━━━  3/6 jam   cooldown: 3m 05s
     menggantikan beban slot-1

  ────────────────────────────────────────────────────

  ban risk score

  🔴  tinggi — 1 nomor kena flag

  waclaw udah auto-pause slot-1.
  semua pesan dipindah ke slot-2 & slot-3.
  lu nggak perlu ngapa-ngapain.

  rekomendasi:
  1  biarin aja (auto-recover 24 jam)
  2  tambah nomor baru buat ganti slot-1
  3  pause semua kirim, cuma scrape

  ↵  biarin aja    2  tambah nomor    3  pause kirim

```

**State: shield_slot_detail** (detail satu nomor)

```

  nomor: 0812-xxxx-3456                     ● aktif

  ────────────────────────────────────────────────────

  statistik 7 hari:

  terkirim    84
  respond     12 (14%)
  gagal       2 (invalid number)
  warning     1 (kemarin, kelebihan rate)

  ────────────────────────────────────────────────────

  riwayat:

  02 mei 14:23   terkirim ke kopi nusantara
  02 mei 14:11   cooldown 8m selesai
  02 mei 13:47   terkirim ke gym fortress
  02 mei 09:15   ⚠ rate jam ini capai 5/6
  01 mei 16:30   terkirim ke salon cantik
  01 mei 09:12   terkirim ke toko makmur
  30 apr 14:23   ⚠ warning: 6/6 jam ini

  ────────────────────────────────────────────────────

  health score: 87/100

  ↓  kalau turun di bawah 50, auto-pause
  ↑  naik 5 poin per hari tanpa warning

  1  pause nomor ini    2  liat leads    q  balik

```

**State: shield_settings** (anti-ban config reference)

```

  pengaturan anti-ban

  semua config di file. edit pakai editor lu.

  ────────────────────────────────────────────────────

  config utama     ~/.waclaw/config.yaml
  (section: anti_ban + spam_guard)

  config aktif:

  ── anti_ban ──

  rate_limit_per_slot   6/jam
  rate_limit_daily      50/hari
  min_delay_between     8 menit
  max_delay_between     25 menit
  delay_variance        30%  (random ±)
  cooldown_after_limit  47 menit
  work_hours            09:00-17:00 wib
  pause_on_flag         auto
  flag_recovery         24 jam
  health_threshold      50/100  (auto-pause)
  rotator_mode          round-robin + cooldown
  template_rotation     aktif
  template_rotation_mode  round-robin  # round-robin | random
  emoji_variation       aktif
  paragraph_shuffle     aktif

  ── spam_guard ──

  max_messages_per_lead   3  (lifetime: ice_breaker + follow_up_1 + follow_up_2)
  message_interval_hours  24  (min jam antar pesan ke lead yang sama)
  follow_up_delay_days    2  (min hari antar follow-up ke lead yang sama)
  follow_up_require_new_variant  aktif  (wajib beda template per follow-up)
  cold_after_followups    2  (setelah 2x follow-up tanpa response → tandai dingin)
  recontact_delay_days    7  (setelah response tanpa deal)
  auto_block_on_stop      aktif  (tambah ke do_not_contact.yaml)
  duplicate_cross_niche   aktif  (nggak kirim ke nomor yang udah dikontak niche lain)
  wa_pre_validation       aktif  (cek nomor WA sebelum antri kirim)
  wa_validation_method    "check-registration"  # check-registration | send-silent

  ── closing_triggers ──

  config per niche     ~/.waclaw/niches/*/niche.yaml
  section              closing_triggers
  auto_mark_deal       aktif  (kalau response match closing_triggers.deal)
  auto_mark_hot        aktif  (kalau response match closing_triggers.hot_lead)
  auto_block_stop      aktif  (kalau response match closing_triggers.stop)
  manual_override      aktif  (user bisa override auto-mark)

  ────────────────────────────────────────────────────

  e  edit config    r  reload    q  balik

```

Micro-interactions:
- `🛡️ semua aman` green pulse = perisai aktif, lu aman
- `⚠️ ada warning` amber flash 2x = waspada, bukan panic
- `✗ BAHAYA` red + auto-interrupt screen = penting, lu harus tau
- Ban risk score bar morphs color real-time (green → amber → red) = lu ngerain tingkat bahaya
- `📱` icon per slot pulses independently = masing-masing hidup sendiri
- Health score animates per poin perubahan = tiap poin terasa
- Shield art degrades/repairs: fill level berubah smooth 50ms per health point = lu NGERAIN perisai makin kuat/lemah
- Flag detection: instant screen interrupt = ini yang paling penting di seluruh app
- Auto-adjust message: "waclaw otomatis pindah beban" = lu tau system merespon, lu nggak perlu ngapa-ngapain
- Shield crack (`╳`) appears when health <50 = visual metaphor yang langsung dipahami

**The neuroscience:** Shield screen itu security blanket. Lu buka → liat hijau → aman. Lu nggak perlu mikir. Tapi kalau merah, lu LANGSUNG tau dan LANGSUNG ada action default (`↵ biarin aja`). Rasa aman = lu bisa minim window. Rasa bahaya = lu tau persis apa yang salah dan apa yang udah waclaw lakuin. Shield art yang berubah = lu ngerain kondisi tanpa baca angka.

---

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

### SCREEN 15: COMPOSE → VOICE

**Ketika user pilih "2 balas custom" dari RESPONSE screen, mereka butuh tempat nulis.**

**Ini cuma muncul sebagai modal overlay di atas RESPONSE screen. Bukan screen penuh.**

**State: compose_draft** (text input area for custom reply)

```

  ── balas custom: kopi nusantara ──────────────────────

  ketik pesan lu di bawah. enter 2x buat kirim.

  ┌───────────────────────────────────────────────────┐
  │ Halo Kak! Bokeh lihat desain website yang         │
  │ udah aku bikin khusus buat Kopi Nusantara?        │
  │                                                    │
  │ Aku kasih preview gratis ya Kak, tinggal           │
  │ liat aja dulu ☺                                    │
  │                                                    │
  │_                                                   │
  └───────────────────────────────────────────────────┘

  ↵↵  kirim    tab  pindah mode    esc  batal

  tip: pakai bahasa santai, kayak ngobrol doang.
  jangan terlalu formal — orang malah curiga.

```

Micro-interactions:
- Compose area slide up dari bawah = modal feel, context (response) masih keliatan di belakang
- Cursor blinks = lu di tempat yang tepat
- Character count subtle di corner = awareness tanpa pressure
- `↵↵ kirim` = double-enter to send, prevents accident
- `esc batal` = always a way out
- Tip di bawah: contextual advice, bukan tutorial

**State: compose_preview** (preview before sending)

```

  ── preview pesan ─────────────────────────────────────

  kirim ke: kopi nusantara

  ────────────────────────────────────────────────────

  Halo Kak! Bokeh lihat desain website yang
  udah aku bikin khusus buat Kopi Nusantara?

  Aku kasih preview gratis ya Kak, tinggal
  liat aja dulu ☺

  ────────────────────────────────────────────────────

  ↵  kirim    e  edit lagi    esc  batal

```

Micro-interactions:
- Preview types out char-by-char = lu baca apa yang bakal mereka baca
- "kirim ke: kopi nusantara" = explicit target, prevent wrong-send
- `↵ kirim` = single enter now = lu udah preview, ga perlu double-confirm
- Brief hold 300ms sebelum bisa press ↵ = prevent double-tap send

**State: compose_template_pick** (quick-pick from template snippets)

```

  ── pilih snippet ─────────────────────────────────────

  snippet yang sering dipakai:

  1  "boleh lihat dulu aja kak"           (soft pitch)
  2  "aku kasih preview gratis ya"         (free offer)
  3  "bisa telepon bentar buat diskusi?"   (move to call)
  4  "harganya mulai 500rb kak"            (direct price)
  5  "mau aku kirim contohnya?"            (send sample)

  ↑↓  pilih    ↵  pake    e  edit dulu    esc  batal

  snippet disimpan di: ~/.waclaw/snippets.md

```

Micro-interactions:
- Snippets load dari `~/.waclaw/snippets.md` = editable, bukan hardcoded
- `↑↓ pilih` → `↵ pake` → inserts ke compose area = quick draft
- `e edit dulu` → insert + buka compose = customize before send
- Category label `(soft pitch)`, `(free offer)` = lu tau tone tiap snippet

---

### SCREEN 16: HISTORY → TIMELINE

**Day/session history view. Cara liat performa masa lalu tanpa harus tunggu mingguan.**

**Akses dari monitor: tekan `h` dari dashboard manapun.**

**State: history_today** (today's activity timeline)

```

  hari ini                                14:23 wib

  ────────────────────────────────────────────────────

  timeline:

  14:23  💬  kopi nusantara      respond     "iya kak, boleh"
  14:01  📤  wedding bliss        terkirim
  13:47  ✅  gym fortress         sampai
  12:31  💬  salon cantik         respond     "makasih tapi..."
  11:58  💬  venue gardenia       respond     "berapa harganya?"
  11:15  📤  toko makmur          terkirim
  10:42  ★  gym fortress         jackpot!    skor: 9.2
  09:42  📤  bengkel jaya         terkirim
  09:15  🔍  batch scrape selesai  67 baru
  09:00  ⚡  auto-pilot on         3 niche

  ────────────────────────────────────────────────────

  ringkasan:
  terkirim  7    respond  4    convert  0
  lead baru 67   scrape 1x

  ↵  liat detail event    ←  hari kemarin    q  balik

```

Micro-interactions:
- Timeline entries fade in sequentially = narrative feel, lu ngerain cerita hari ini
- 💬 respond rows amber tint = yang paling penting langsung keliatan
- ★ jackpot rows gold flash = highlight moments
- `← hari kemarin` = arrow key navigation between days
- Auto-scroll ke event terbaru = lu selalu di "now"

**State: history_week** (weekly summary with mini charts)

```

  minggu ini                              28 apr - 4 mei

  ────────────────────────────────────────────────────

  pesan terkirim per hari:

  senin   ████████████░░░░░░  12
  selasa  ██████████████████  18  ★ terbaik
  rabu    ██████████████░░░░  14
  kamis   ██████████░░░░░░░░  10
  jumat   ████████░░░░░░░░░░   8
  sabtu   ████░░░░░░░░░░░░░░   4  (weekend)
  minggu  ██░░░░░░░░░░░░░░░░   2  (weekend)

  response per hari:

  senin   ████░░░░░░░░░░░░░░   2
  selasa  ██████████░░░░░░░░   5  ★ terbaik
  rabu    ██████░░░░░░░░░░░░   3
  kamis   ████░░░░░░░░░░░░░░   2
  jumat   ██░░░░░░░░░░░░░░░░   1
  sabtu   ░░░░░░░░░░░░░░░░░░   0
  minggu  ░░░░░░░░░░░░░░░░░░   0

  ────────────────────────────────────────────────────

  total minggu ini:
  terkirim  68    respond  13    convert  3
  lead baru 247   avg response time: 3.2 jam

  ★ selasa = hari terbaik lu
    jam 10-12 = prime time
    conversion rate selasa: 8.3%

  ────────────────────────────────────────────────────

  ↑↓  pilih hari    ↵  liat detail    q  balik

```

Mini charts micro-interactions:
- Bar chart bars grow dari kiri ke kanan 50ms per bar = sequential build
- `★ terbaik` star pulse amber 1x = highlight the insight
- Stats animate in: numbers count up from 0 = data feels discovered, bukan displayed
- Insight ("selasa = hari terbaik lu") fade in last = the takeaway lands after the data

**State: history_day_detail** (specific past day)

```

  selasa, 30 april 2024

  ────────────────────────────────────────────────────

  ringkasan:
  terkirim  18    respond  5    convert  1
  lead baru 52    scrape 2x

  ────────────────────────────────────────────────────

  timeline:

  16:42  🎉  kopi nusantara      CONVERTED   rp 2.5jt
  15:11  💬  kopi nusantara      respond     "oke gas!"
  14:23  📤  salon cantik        terkirim
  13:47  💬  wedding bliss       respond     "boleh kirim?"
  12:31  📤  toko makmur         terkirim
  11:58  📤  gym fortress        terkirim
  10:15  💬  bengkel jaya        respond     "berapa?"
  09:42  📤  venue gardenia      terkirim
  09:15  📤  kopi nusantara      offer sent
  09:00  🔍  batch scrape selesai  52 baru

  ────────────────────────────────────────────────────

  →  hari berikutnya    ←  hari sebelumnya    q  balik

```

Micro-interactions:
- `🎉 CONVERTED` row gets gold shimmer = peak moment in past data
- Revenue `rp 2.5jt` gold pulse = uang beneran, bahkan di history
- Day navigation: `← →` smooth crossfade = flipping through days feels natural
- Converted leads in history still get mini celebration = reward remembering

---

### SCREEN 17: FOLLOW-UP → PERSISTENCE

**Lead yang nggak jawab BUKAN lead mati. Cuma belum waktunya.**

Screen ini jalan otomatis di background. WaClaw ngatur timing, varian, dan jeda. Lu cuma liat status dan approve kalo mau. Tapi default-nya: auto.

**Filosofi follow-up:**
- Ice breaker = ketukan pertama. Follow-up = ketukan lagi. Tapi pintunya nggak selalu dibuka pertama kali.
- Max 3 pesan per lead lifetime (ice_breaker + follow_up_1 + follow_up_2)
- Jeda minimum 2 hari antar follow-up
- Wajib beda template/varian per follow-up — WA bisa deteksi pola kalau pesan sama
- Setelah 2x follow-up tanpa response → auto-tandai dingin, 1 chance terakhir optional
- Lead yang pernah respond tapi dingin lagi → re-contact setelah 7 hari (beda kategori)

**State: followup_dashboard** (overview semua follow-up)

```

  follow-up                                 ● auto-jalan

  ░░ 2 8 4 1 ░░ 7 3 9 ░░ 5 6 ░░

  ────────────────────────────────────────────────────

  antrian follow-up hari ini                 14 pesan

  ▸ web_developer
    follow-up 1    8 lead (ice breaker 2 hari lalu, belum jawab)
    follow-up 2    3 lead (follow-up 1 kemarin, belum jawab)
    dingin          2 lead (2x follow-up, masih diam)

  ▸ undangan_digital
    follow-up 1    4 lead
    dingin          1 lead

  ▸ social_media_mgr
    follow-up 1    2 lead
    follow-up 2    1 lead
    dingin          1 lead

  ────────────────────────────────────────────────────

  total: 14 follow-up hari ini · 4 dingin
  semua auto-kirim pas jam kerja.
  varian beda tiap follow-up. nggak ada yang sama.

   ↵  liat detail    a  auto-semua    q  balik

  lu nggak perlu ngapa-ngapain. waclaw ngatur timing + varian.

```

Ambient effect: sama kayak monitor — faint number rain `░░ 2 8 4 1 ░░` di background. Follow-up itu army yang kerja pelan tapi persisten. Angka-angka itu represent leads yang lagi di-follow up.

**State: followup_niche_detail** (detail per niche)

```

  follow-up: web_developer

  ────────────────────────────────────────────────────

  follow-up 1 (8 lead)

  01  kopi nusantara      ice breaker: 2 hari lalu   → follow-up hari ini
  02  gym fortress pro    ice breaker: 2 hari lalu   → follow-up hari ini
  03  salon cantik        ice breaker: kemarin       → follow-up besok
  04  toko makmur jaya    ice breaker: 3 hari lalu   → follow-up hari ini
  ...4 lainnya

  follow-up 2 (3 lead)

  05  bengkel jaya        follow-up 1: kemarin       → follow-up 2 besok lusa
  06  cafe nusantara      follow-up 1: 2 hari lalu   → follow-up 2 hari ini
  07  salon indah         follow-up 1: 2 hari lalu   → follow-up 2 hari ini

  dingin (2 lead)

  08  apotek sehat        2x follow-up, belum jawab  ❄
  09  toko elektronik     2x follow-up, belum jawab  ❄

  ────────────────────────────────────────────────────

  varian follow-up:
  ▸ follow_up_1.md → "halo kak, cuma ngingetin aja"
  ▸ follow_up_2.md → "kak, penawaran terbatas nih"
  ▸ follow_up_3.md → "terakhir kak, kalo berkenan"  (hanya buat manual)

   ↑↓  pilih lead    ↵  liat detail    a  auto-semua    q  balik

```

**State: followup_sending** (follow-up lagi dikirim)

```

  follow-up                                 ● lagi kirim

  ────────────────────────────────────────────────────

  ▸ web_developer

  01  →  kopi nusantara       📱 slot-2
       follow_up_1: variant_2  ← beda dari ice breaker!
       ━━━━━━━━━━━━━━━━━ mengirim...

  02     gym fortress pro     📱 slot-1
       follow_up_1: variant_1  ← rotasi!
       nunggu (berikutnya: 14m 02s)

  ────────────────────────────────────────────────────

  rate: 9/18 per jam · follow-up hari ini: 2/14

  varian follow-up beda dari ice_breaker. WA nggak bakal deteksi pola.
  tiap follow-up pakai template + varian sendiri.

   p  pause    ↵  skip tunggu    tab  pindah niche    q  balik

```

**State: followup_empty** (nggak ada yang perlu follow-up hari ini)

```

  follow-up

  nggak ada yang perlu di-follow up hari ini.
  semua lead masih dalam jeda minimum (2 hari).

  lead yang belom jawab ice breaker: 8
  lead dingin: 4

   ↵  liat lead dingin    q  balik

```

**State: followup_cold_list** (daftar lead dingin)

```

  lead dingin                                4 lead

  ────────────────────────────────────────────────────

  ini lead yang udah 2x di-follow up tapi belum jawab.
  waclaw nggak bakal kirim lagi otomatis.
  tapi lu bisa coba 1x terakhir manual.

  01  apotek sehat          ice breaker 5 hari lalu
      follow-up 1: 3 hari lalu
      follow-up 2: 1 hari lalu
      ❄ dingin — 0 response

  02  toko elektronik       ice breaker 6 hari lalu
      follow-up 1: 4 hari lalu
      follow-up 2: 2 hari lalu
      ❄ dingin — 0 response

  ...2 lainnya

  ────────────────────────────────────────────────────

  ↑↓  pilih    ↵  kirim follow-up terakhir (ke-3)    a  archive semua    q  balik

  hati-hati: follow-up ke-3 itu peluru terakhir.
  kalau masih nggak jawab, archive aja.

```

**State: followup_recontact** (lead yang pernah respond tapi dingin lagi)

```

  follow-up: re-contact

  ────────────────────────────────────────────────────

  ini lead yang pernah respond tapi nggak jadi deal.
  setelah 7 hari jeda, waclaw bisa coba lagi dengan pendekatan baru.

  01  salon cantik          respond 8 hari lalu: "berapa harganya?"
      offer terkirim 7 hari lalu
      setelah itu: diam
      → bisa re-contact hari ini

  02  wedding bliss WO      respond 10 hari lalu: "boleh kirim contohnya?"
      offer terkirim 9 hari lalu
      setelah itu: diam
      → bisa re-contact hari ini

  ────────────────────────────────────────────────────

  re-contact pakai template berbeda dari offer pertama.
  tone-nya lebih santai, "hai lagi kak" vibes.

  ↑↓  pilih    ↵  kirim re-contact    a  auto-semua    q  balik

```

Micro-interactions:
- Follow-up dashboard: sama kayak monitor — ambient data rain + breathing stats = army yang persisten
- `❄ DINGIN` badge dim blue = dingin, bukan mati. Masih bisa di-revive.
- Follow-up sending: sama kayak SEND screen tapi badge beda = lu tau ini follow-up, bukan ice breaker baru
- Cold list: dimmed rows = bukan prioritas, tapi nggak hilang
- Re-contact: warm amber tint = pernah tertarik, mungkin waktunya lagi
- Auto-semua (`a`) = satu tombol buat approve semua follow-up = auto-pilot mindset

**The neuroscience:** Follow-up itu consistency. Lead nggak jawab bukan karena nggak tertarik — sering karena sibuk, lupa, atau belum butuh. Tapi mereka NGERAIN kalau ada yang konsisten ngontak (tanpa nyebelin). Follow-up varian beda = lu bukan bot yang ngulang pesan yang sama. Lu orang yang sabar nunggu. Dan kalau mereka jawab di follow-up ke-2? Reward-nya lebih manis karena effort-nya lebih besar.

---

### NERD STATS → VITALS (Toggle Overlay)

**Bukan screen. Overlay global. Tekan backtick (`) buat toggle.**

Ini bukan screen baru — ini overlay yang bisa muncul di atas screen APAPUN. Tujuannya: kasih lu info system vitals tanpa harus pindah screen. Default = hidden. Tekan ` sekali = muncul. Tekan ` lagi = ilang. Simple.

**Kenapa overlay, bukan screen?** Karena system vitals itu info yang lu butuh SESUDAH, bukan info yang lu cari. Lu lagi liat monitor, tiba-tiba penasaran "RAM lu berapa sekarang?" — tekan `, liat, tekan ` lagi, balik kerja. Nggak perlu pindah screen, nggak perlu kehilangan context.

**Dua mode:**

**Mode: minimal** (1-line footer, default saat toggle on)

```

  ── CPU 12% · RAM 134MB · Goroutines 23 · DB 2.4MB · Uptime 4j 12m ──

```

Minimal mode: cuma 1 baris di paling bawah layar. Warna `text_dim` (hampir ga keliatan tapi bisa dibaca). Nggak ganggu konten utama. Cuma info cepat.

**Mode: expanded** (3-line panel, tekan ` lagi dari minimal)

```

  ── nerd stats ──────────────────────────────────────

    CPU         12.3%  ████░░░░░░░░░░░░░░░░
    RAM         134MB  ██████░░░░░░░░░░░░░░  / 512MB
    Goroutines  23     ██░░░░░░░░░░░░░░░░░░  / 100
    DB Size     2.4MB  █░░░░░░░░░░░░░░░░░░░  / 50MB
    Uptime      4j 12m
    Version     v1.3.2

  ────────────────────────────────────────────────────

```

Expanded mode: 3 baris panel di bawah layar. Mini bar chart per metric. Warna tetap `text_dim` kecuali kalau ada yang mendekati limit (RAM > 80% = `warning` amber, goroutine > 80 = `danger` merah).

**Metrics yang ditampilin:**
- **CPU%** — usage proses WaClaw saat ini
- **RAM** — memory dipakai / total
- **Goroutines** — jumlah goroutine aktif / max configured
- **DB Size** — ukuran database leads / max
- **Uptime** — berapa lama WaClaw jalan sejak startup
- **Version** — versi WaClaw yang lagi jalan (hanya di expanded mode)

**Toggle behavior:**
- Tekan ` pertama: hidden → minimal
- Tekan ` kedua: minimal → expanded
- Tekan ` ketiga: expanded → hidden
- Kalau 30 detik nggak ditekan: auto-collapse ke hidden (kembali ke mode default)
- Bisa di-toggle dari screen APAPUN, termasuk saat lagi di compose, review, atau monitor

**Data source:** Semua metrics diambil dari Go runtime (`runtime.ReadMemStats`, `runtime.NumGoroutine`) dan SQLite DB stats. Zero network call. Zero external dependency. Ini beneran vitals system lu, bukan data dari server.

Micro-interactions:
- Toggle on: minimal footer slide up dari bawah 150ms = nggak ngejut
- Toggle expanded: panel expand dari 1 baris ke 3 baris dengan height morph 200ms
- Toggle off: collapse ke hidden 150ms
- Bar chart fills gradient sweep 300ms = progress feel
- RAM > 80%: bar morph ke amber + subtle pulse = "keep an eye on this"
- Goroutine > 80: bar morph ke merah + double flash = "something might be wrong"
- Auto-collapse: gentle fade out 500ms = nggak abrupt
- Metrics update setiap 2 detik = live, bukan snapshot

**The neuroscience:** Nerd stats itu kayak detak jantung. Lu nggak perlu liat terus — tapi kalau penasaran, 1 tombol. Overlay, bukan screen, karena vitals itu context-aware. Lu lagi liat leads, bukan lagi liat system. Tapi kalau system minta perhatian (RAM tinggi), warnanya yang ngomong, bukan notifikasi.

---

### CTRL+K → COMMAND PALETTE (Global Overlay)

**Bukan screen. Bukan menu. Ini teleport.**

Command palette itu cara tercepat buat ngapa-ngapain di WaClaw. Tekan `Ctrl+K` dari mana aja → ketik apa yang lu mau → ↵ → eksekusi. Nggak perlu hafal shortcut. Nggak perlu navigate 3 layer deep. Pikiran → ketik → jadi. 3 detik max.

**Kenapa command palette, bukan menu?** Menu itu hierarchy — lu harus tau dulu menu apa, terus submenu apa, terus item mana. Command palette itu FLAT — semua command di 1 level, lu cuma ketik yang lu mau. VS Code, Notion, Linear — semua pakai command palette karena itu cara tercepat buat power user. Dan WaClaw user itu power user.

**Kenapa Ctrl+K?** Karena `K` = "komando". Dan Ctrl+K udah jadi standard de facto buat command palette di terminal apps (kayak fzf, zoxide). Lagian, lu nggak butuh Ctrl+K buat apapun lain di WaClaw — semua action cukup 1 tombol tanpa modifier. Ctrl+K itu satu-satunya chord key di seluruh app = gampang ingat.

**State: cmd_closed** (default — palette nggak keliatan)

Palette hidden. Screen sekarang normal. Tekan Ctrl+K → buka palette.

**State: cmd_open** (palette terbuka, search aktif)

```

  ┌─────────────────────────────────────────────────────────────┐
  │  > scrape                                              │  ×  │
  ├─────────────────────────────────────────────────────────────┤
  │                                                             │
  │  ── recently used ──────────────────────────────────────    │
  │                                                             │
  │  ▸ Scrape semua niche sekarang            s  ·  scrape     │
  │    Force scrape satu worker               1  ·  scrape     │
  │                                                             │
  │  ── commands ──────────────────────────────────────────    │
  │                                                             │
  │    Pause semua worker                     p  ·  workers    │
  │    Pause satu worker                      2  ·  workers    │
  │    Lihat leads database                   1  ·  database  │
  │    Lihat dashboard                        d  ·  monitor   │
  │    Edit config                            e  ·  settings  │
  │    Validate semua config                  v  ·  guardrail │
  │    Lihat anti-ban shield                  5  ·  shield    │
  │    Cek update                             u  ·  update    │
  │    Lihat lisensi                          l  ·  license   │
  │    Lihat riwayat                          h  ·  history   │
  │    Explore niche baru                     n  ·  explorer  │
  │    Reload config                          r  ·  settings  │
  │    Kirim follow-up semua                  a  ·  followup  │
  │    Toggle nerd stats                      `  ·  overlay   │
  │                                                             │
  └─────────────────────────────────────────────────────────────┘

  ↑↓  pilih    ↵  eksekusi    esc  tutup    Ctrl+K  tutup

```

**Anatomy of command palette:**

1. **Search bar** (`> scrape`) — ketik apapun, filter langsung real-time. Debounce 50ms (instant). Match against command name + description + shortcut + category.

2. **Recently used** — 3 command terakhir yang lu pakai. Otomatis muncul di atas. Kalau lu sering "scrape semua", itu selalu di paling atas. Muscle memory + recency = fastest path.

3. **Filtered results** — semua command yang match search query. Diurutin berdasarkan: exact match > name match > description match > category match. Masing-masing baris punya: nama command, shortcut asli (kalau ada), dan kategori.

4. **Category tags** — setiap command punya category tag (`·  scrape`, `·  workers`, `·  database`). Ini buat filter juga — ketik "worker" muncul semua yang berkaitan. Ketik "shield" langsung ke anti-ban.

**State: cmd_executing** (command dipilih, lagi eksekusi)

```

  ┌─────────────────────────────────────────────────────────────┐
  │  > scrape                                              │  ×  │
  ├─────────────────────────────────────────────────────────────┤
  │                                                             │
  │  ▸ Scrape semua niche sekarang            s  ·  scrape     │
  │                                                             │
  │    ● executing...                                           │
  │                                                             │
  └─────────────────────────────────────────────────────────────┘

```

Execute animation: selected row pulse accent → palette collapses → command runs. Total transition 300ms. Nggak ada loading screen — palette tutup, screen baru muncul / action jalan.

**State: cmd_empty** (ketik sesuatu, nggak ada yang match)

```

  ┌─────────────────────────────────────────────────────────────┐
  │  > xyz                                                 │  ×  │
  ├─────────────────────────────────────────────────────────────┤
  │                                                             │
  │  nggak nemu command "xyz"                                   │
  │                                                             │
  │  coba: scrape, pause, leads, config,                       │
  │        shield, follow-up, update, license                   │
  │                                                             │
  └─────────────────────────────────────────────────────────────┘

```

Empty state: helpful suggestions, bukan cuma "no results". Lu tau command apa yang ada.

**Full Command Registry:**

```
── NAVIGATION ──────────────────────────────────────────────
  Dashboard / Monitor           d      ·  monitor
  Leads Database                1      ·  database
  Send Queue                    2      ·  send
  Workers                       3      ·  workers
  Templates                     4      ·  template
  Anti-Ban Shield               5      ·  shield
  Follow-Up                     6      ·  followup
  Settings                      7      ·  settings
  History                       h      ·  history
  Niche Explorer                n      ·  explorer
  License                       l      ·  license
  Update / Upgrade              u      ·  update

── ACTIONS ─────────────────────────────────────────────────
  Scrape semua niche sekarang   s      ·  scrape
  Scrape satu worker            —      ·  scrape
  Pause semua worker            p      ·  workers
  Resume semua worker           —      ·  workers
  Kirim follow-up semua         a      ·  followup
  Edit config                   e      ·  settings
  Reload config                 r      ·  settings
  Validate semua config         v      ·  guardrail
  Cek versi baru                u      ·  update
  Toggle nerd stats             `      ·  overlay
  Logout WA semua slot          —      ·  login
  Force scrape retry            —      ·  scrape

── OVERLAYS ────────────────────────────────────────────────
  Toggle nerd stats             `      ·  overlay
  Command palette               Ctrl+K ·  overlay
  Shortcut help                 ?      ·  overlay

── SPECIAL ─────────────────────────────────────────────────
  Compose custom reply          —      ·  compose
  Search leads                  /      ·  database
  Export leads CSV              —      ·  database
  Mark lead converted           —      ·  response
  Block lead                    —      ·  response
  Re-contact cold lead          —      ·  followup
```

**Key behaviors:**

1. **Fuzzy search.** Ketik "scrp" → match "scrape". Ketik "pld" → match "pause all worker". Nggak perlu exact match — fuzzy finder yang cerdas. Algo: fzf-style scoring (exact > prefix > substring > fuzzy).

2. **Context-aware.** Command yang lagi relevant di screen sekarang muncul lebih tinggi. Kalau lu lagi di SCRAPE screen, "pause worker" dan "force retry" muncul di atas. Kalau di MONITOR, "scrape semua" dan "leads database" yang muncul duluan.

3. **Recently used priority.** 3 command terakhir selalu di bagian atas, di atas semua hasil search. Karena 80% waktu lu ngulang action yang sama.

4. **No dead ends.** Setiap command punya hasil. Navigasi = pindah screen. Action = execute. Toggle = toggle. Nggak ada command yang "coming soon" atau disabled.

5. **Esc / Ctrl+K = close.** Dua cara buat tutup. Sama kayak cara buka. Lu bisa Ctrl+K → liat → Ctrl+K → tutup. Atau Esc. Whichever muscle memory lu punya.

6. **Ctrl+C = hard cancel.** Kalau palette terbuka dan lu beneran mau keluar WaClaw, Ctrl+C dari palette = close palette aja (bukan exit app). Dari screen manapun tanpa palette, Ctrl+C = keluar.

**Variant: cmd_with_recent** (ada recently used commands)

```

  ┌─────────────────────────────────────────────────────────────┐
  │  >                                                    │  ×  │
  ├─────────────────────────────────────────────────────────────┤
  │                                                             │
  │  ── barusan ───────────────────────────────────────────     │
  │                                                             │
  │  ▸ Scrape semua niche sekarang            s  ·  scrape     │
  │    Lihat leads database                   1  ·  database  │
  │    Toggle nerd stats                      `  ·  overlay   │
  │                                                             │
  │  ── semua command ─────────────────────────────────────     │
  │                                                             │
  │    Dashboard / Monitor                    d  ·  monitor    │
  │    Pause semua worker                     p  ·  workers    │
  │    Edit config                            e  ·  settings   │
  │    Validate config                        v  ·  guardrail  │
  │    Lihat shield                           5  ·  shield     │
  │    ...                                                     │
  │                                                             │
  └─────────────────────────────────────────────────────────────┘

```

Kosong (belum ngetik) = recently used + semua command. Ngetik = filter. Always helpful.

**Variant: cmd_quick_action** (command yang langsung execute tanpa pindah screen)

Beberapa command itu ACTIONS, bukan navigasi. Contoh: "pause semua worker", "scrape sekarang", "reload config". Command ini langsung execute — palette tutup, toast notification muncul 2 detik konfirmasi, lu tetap di screen yang sama.

```

  ┌─────────────────────────────────────────────────────────────┐
  │  > pause                                              │  ×  │
  ├─────────────────────────────────────────────────────────────┤
  │                                                             │
  │  ▸ Pause semua worker                     p  ·  workers    │
  │    ⚡ langsung eksekusi, nggak pindah screen               │
  │                                                             │
  │    Pause satu worker                      2  ·  workers    │
  │    Resume semua worker                    —  ·  workers    │
  │                                                             │
  └─────────────────────────────────────────────────────────────┘

```

Quick action indicator: `⚡ langsung eksekusi` = lu tau ini command langsung jalan, bukan navigasi. Nggak ada surprise.

Micro-interactions:
- Palette open: slide down from top 150ms + backdrop dim = "focus here now"
- Search typing: results filter real-time zero lag = instant response
- Selection highlight: accent color slide 50ms = snappy, not sluggish
- Execute: selected row pulse accent → palette collapse → screen/action = 300ms total
- Recently used: badge `barusan` with subtle amber tint = these are YOUR commands
- Quick action: `⚡` icon with brief glow = "this runs immediately"
- Empty state: suggestion chips clickable = zero friction recovery
- Close: palette slides up 100ms + backdrop undim = smooth return to context
- Fuzzy match: matched characters get `accent` color highlight = lu tau kenapa ini match
- Category tags: `·  scrape` dimmed = metadata, bukan noise

**The neuroscience:** Command palette itu working memory externalizer. Otak lu tau mau ngapa — "mau pause workers" — tapi lu nggak perlu ingat tombolnya apa. Lu cuma ketik apa yang lu pikirin, dan WaClaw ngerti. Itu kenapa VS Code, Notion, dan semua app modern pakai command palette: karena **natural language > memorized shortcuts**. Tapi buat yang udah hafal shortcut, shortcut tetap jalan — command palette itu ALTERNATIF, bukan pengganti. Lu bisa pake `s` langsung, ATAU Ctrl+K → "scrape". Pilihan lu.

---

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

## 2. NOTIFICATION SYSTEM — The "Poke"

**WaClaw yang nanya, bukan lu yang nyari.**

Notification = interrupt yang gentle. Muncul sebagai overlay di atas screen apapun. Auto-dismiss setelah 10 detik kalau nggak di-respon.

**Notification: response masuk**

```

  ┌─────────────────────────────────────────────────┐
  │  💬 kopi nusantara balas pesan lu               │
  │  [web_dev] "iya kak, boleh lihat desainnya?"    │
  │                                                  │
  │  ↵  balas    s  nanti    q  dismiss              │
  └─────────────────────────────────────────────────┘

```

**Notification: scrape selesai** (multi-niche)

```

  ┌─────────────────────────────────────────────────┐
  │  🔍 2 niche selesai scrape                       │
  │  web_dev: 67 baru · undangan: 48 baru            │
  │  total: 115 lead baru masuk antrian               │
  │                                                  │
  │  ↵  liat    s  biarin aja                        │
  └─────────────────────────────────────────────────┘

```

**Notification: batch kirim selesai**

```

  ┌─────────────────────────────────────────────────┐
  │  📦 batch selesai — 12 pesan terkirim            │
  │  web_dev: 8 · undangan: 4                         │
  │  batch berikutnya: 47 menit lagi                  │
  │                                                  │
  │  ↵  liat    s  ok                                 │
  └─────────────────────────────────────────────────┘

```

**Notification: wa disconnect**

```

  ┌─────────────────────────────────────────────────┐
  │  ✗ wa putus (slot-1)! auto-reconnect nyala...   │
  │  slot-2 & slot-3 tetap jalan                     │
  │                                                  │
  │  1  login ulang    s  dismiss                     │
  └─────────────────────────────────────────────────┘

```

**Notification: wa flag (ban warning)**

```

  ┌─────────────────────────────────────────────────┐
  │  🚨 nomor 0812-xxxx-3456 kena flag!              │
  │  auto-pause slot-1, beban pindah ke slot-2/3     │
  │                                                  │
  │  ↵  liat shield    s  biarin waclaw urus         │
  └─────────────────────────────────────────────────┘

```

**Notification: health score drop**

```

  ┌─────────────────────────────────────────────────┐
  │  📉 health score slot-1 turun ke 55/100          │
  │  masih aman, tapi mendekati threshold (50)       │
  │                                                  │
  │  ↵  liat shield    s  ok                         │
  └─────────────────────────────────────────────────┘

```

**Notification: limit harian**

```

  ┌─────────────────────────────────────────────────┐
  │  📊 limit hari ini capai (50/50)                 │
  │  sisa antrian besok aja ya                        │
  │                                                  │
  │  ↵  ok                                           │
  └─────────────────────────────────────────────────┘

```

**Notification: streak / milestone**

```

  ┌─────────────────────────────────────────────────┐
  │  🔥 10 response minggu ini! conversion naik 2%   │
  │                                                  │
  │  ↵  liat stats    s  nice 👌                      │
  └─────────────────────────────────────────────────┘

```

**Notification: config error** (critical)

```

  ┌─────────────────────────────────────────────────┐
  │  ✗ config error — 1 worker di-pause              │
  │  fotografer/niche.yaml: parse error baris 14     │
  │  worker lain tetap jalan                          │
  │                                                  │
  │  ↵  liat error    v  validasi semua    s  nanti   │
  └─────────────────────────────────────────────────┘

```

**Notification: multi response** (3+ responses at once)

```

  ┌─────────────────────────────────────────────────┐
  │  💬 3 response masuk barengan!                    │
  │  kopi nusantara · wedding bliss · bengkel jaya   │
  │  2 positif · 1 auto-reply                         │
  │                                                  │
  │  ↵  proses    1  auto-offer yang positif    s  nanti│
  └─────────────────────────────────────────────────┘

```

**Notification: validation error**

```

  ┌─────────────────────────────────────────────────┐
  │  ✗ validasi gagal — 3 config error ditemukan     │
  │  2 niche di-pause, 1 tetap jalan                 │
  │                                                  │
  │  ↵  liat error    v  validasi semua    s  nanti   │
  └─────────────────────────────────────────────────┘

```

**Notification: lisensi expired** (critical — army berhenti)

```

  ┌─────────────────────────────────────────────────┐
  │  ✗ lisensi lu udah expired — waclaw berhenti     │
  │  semua worker di-pause. data aman.               │
  │                                                  │
  │  ↵  masukin lisensi baru    q  keluar             │
  └─────────────────────────────────────────────────┘

```

**Notification: device conflict** (critical — army berhenti)

```

  ┌─────────────────────────────────────────────────┐
  │  ✗ lisensi lagi dipakai di PC-KANTOR!            │
  │  waclaw berhenti. 1 lisensi = 1 device.          │
  │                                                  │
  │  ↵  liat lisensi    2  putuskan device lain       │
  │  s  keluar                                        │
  └─────────────────────────────────────────────────┘

```

**Notification: follow-up terjadwal** (gentle reminder)

```

  ┌─────────────────────────────────────────────────┐
  │  📋 14 follow-up terjadwal hari ini               │
  │  waclaw auto-kirim pas jam kerja                  │
  │                                                  │
  │  ↵  liat    s  ok                                 │
  └─────────────────────────────────────────────────┘

```

**Notification: lead dingin** (informative)

```

  ┌─────────────────────────────────────────────────┐
  │  ❄ 4 lead dingin — 2x follow-up tanpa response   │
  │  waclaw nggak bakal kirim lagi otomatis           │
  │                                                  │
  │  ↵  liat    s  ok                                 │
  └─────────────────────────────────────────────────┘

```

**Notification: update available** (positive/neutral — minor update)

```

  ┌─────────────────────────────────────────────────┐
  │  🔄 update tersedia: v1.3.3                       │
  │  bug fix + perf improvement. gratis!              │
  │                                                  │
  │  ↵  update    u  nanti    s  skip                  │
  └─────────────────────────────────────────────────┘

```

**Notification: upgrade available** (informative — major upgrade)

```

  ┌─────────────────────────────────────────────────┐
  │  ⬆️ upgrade tersedia: v2.0.0                      │
  │  major version — butuh lisensi baru               │
  │  v1 lu tetap jalan selama lisensi valid           │
  │                                                  │
  │  ↵  liat info    u  upgrade    s  nanti            │
  └─────────────────────────────────────────────────┘

```

**Confirmation Overlay: bulk action** (sebelum aksi besar/destruktif)

```

  ┌─────────────────────────────────────────────────┐
  │  ⚠️ kirim offer ke 5 leads sekaligus?            │
  │  3 dari web_dev · 2 dari undangan                │
  │                                                  │
  │  ↵  gas kirim    s  batal                         │
  └─────────────────────────────────────────────────┘

```

**Rules:**
- Max 1 notification di layar. Queue the rest.
- Critical (wa disconnect, config error, validation error, lisensi expired, device conflict) = instant, cannot dismiss for 3 seconds
- Positive (response, milestone, update available) = gentle fade-in, auto-dismiss 10s (update available: 15s)
- Neutral (scrape done, follow-up terjadwal) = brief, auto-dismiss 5s
- Informative (lead dingin, upgrade available) = brief, auto-dismiss 7s (upgrade available: 20s)
- Multi-response = slightly longer display (15s) because more data to process
- Confirmation overlay = muncul sebelum bulk action (auto-offer semua, delete semua, force device disconnect). Always `↵ gas` + `s batal`. Nggak pernah lebih dari 2 opsi.
- Never stack. Never spam. One at a time.

---

## 3. Micro-Interactions Catalog

### Navigation

| Aksi | Animasi | Durasi | Feel |
|------|---------|--------|------|
| Screen transition | Horizontal slide | 300ms | Maju |
| Back navigation | Slide reverse | 300ms | Mundur |
| Tab switch | Cross-fade + vertical shift | 200ms | Ganti konteks |
| Scroll | Smooth scroll (2 line/frame) | 150ms | Ngalir |
| Notification masuk | Slide dari atas + fade in | 250ms | Ada yang penting |
| Notification dismiss | Fade out + slide ke atas | 200ms | Sudah ditangani |

### Data

| Event | Animasi | Durasi | Feel |
|-------|---------|--------|------|
| Angka naik | Flash bright + scale 1.05x | 200ms | Ada perubahan |
| Item baru | Slide dari kanan | 250ms | Baru datang |
| Item hilang | Fade + slide ke bawah | 300ms | Dibuang |
| Status berubah | Color morph (never instant) | 400ms | Evolusi |
| Progress fill | Gradient sweep | Variable | Membangun |
| Breathing stats | Opacity pulse 0.9→1.0→0.9 | 4000ms | Hidup |

### Feedback

| Event | Animasi | Durasi | Feel |
|-------|---------|--------|------|
| Sukses | Pulse hijau (1.0→1.2→1.0 opacity) | 500ms | Berhasil |
| Perhatian | Amber double-flash | 600ms | Cek ini |
| Error | Merah edge glow, 3px, fades | 800ms | Ada masalah |
| Selesai | Full-width bar fills, hold, fade | 1000ms | Done |
| Deal! | Flash putih + ✦ particle scatter + terminal bell | 1200ms | MENANG |
| Config error | Red underline blink 2x | 600ms | Fix ini |

### Ambient Effects

| Effect | Location | Description | Purpose |
|--------|----------|-------------|---------|
| Data rain | Monitor dashboard | Faint scrolling numbers `░░ 3 7 1 4 ░░` | System alive, data flowing |
| Breathing stats | Monitor stats | Subtle opacity pulse on numbers | Dashboard bukan screenshot |
| Shield pulse | Anti-ban shield | Shield fill level breathes | Perisai hidup |
| Army march | Boot returning | `▸▸▸ → ●` worker rows | Pasukan siap |
| Nerd stats | Global overlay | 1-line footer / 3-line panel vitals | System pulse on demand |
| Command palette | Global overlay | Search + filtered list + fuzzy match | Instant access to everything |

### Dramatic Reveals

| Event | Screen | Description | Purpose |
|-------|--------|-------------|---------|
| High-value lead | Scrape | Slot machine name scroll + jackpot bounce + bell | Peak excitement |
| Batch complete | Scrape | Cascade fall-in + sequential checkmarks | Satisfying closure |
| Conversion | Response | Full white flash + particles + sound + color wave | THE moment |
| Config fix | Validation | Error rows collapse to ✓ ok | Satisfying fix |
| Shield repair | Anti-ban | Shield fill grows bottom-to-top | Recovery visible |

---

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
```

**Rules:**
- `danger` merah HANYA buat masalah teknis (wa putus, scraper error, config error). Rejection = `text_muted`.
- `accent` (indigo) cuma muncul di hal yang bisa lu interaksi. Kalau indigo, berarti bisa ditekan.
- `pulse` buat indikator hidup (connection status, scrape aktif). Bernapas = pulse.
- Rejection, skip, decline = netral. Bukan gagal. Cuma bukan jodoh.
- `gold` HANYA buat jackpot leads, revenue numbers, dan conversion celebration. Kalau gold, berarti UANG.
- `celebration` putih full-flash ONLY untuk conversion. Satu-satunya momen yang boleh secerah itu.

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

## 7. Complete Screen Flow

```
  BOOT (first time)
    │
    ├─→ LICENSE ── belum ada key? ──→ LICENSE INPUT
    │     │                              │
    │     ▼                              ▼
    ├─→ VALIDATION ── config missing? ──→ CONFIG SETUP
    │     │                                 │
    │     ▼                                 ▼
    │   LOGIN ──── udah pernah? ───→ SKIP   │
    │     │                                 │
    │     ▼                                 │
    │   NICHE SELECT                        │
    │     │                                 │
    │     ▼                                 │
    │   VALIDATION ── niche config ok?      │
    │     │                                 │
    │     ▼                                 │
    │   SCRAPE (auto)                       │
    │     │                                 │
    │     ▼                                 │
    │   WA VALIDATION ── cek nomor punya WA? (background)
    │     │                                 │
    └─────┘─────────────────────────────────┘

  BOOT (returning)
    │
    ▼
  LICENSE CHECK ── expired? ──→ LICENSE EXPIRED (hard stop)
    │                                  │
    │                device conflict? ──→ DEVICE CONFLICT (hard stop)
    │
    ▼
  VALIDATION ── config error? ──→ SHOW ERRORS (partial pause)
    │
    ▼
  MONITOR (home base)
    │
    ├── auto: SCRAPER jalan tiap interval
    │     └── WA VALIDATION otomatis setelah scrape
    ├── auto: SENDER jalan pas jam kerja (hanya WA-validated leads)
    │     ├── VARIAN ROTASI: ice_breaker + offer keduanya rotate
    │     └── FOLLOW-UP auto-kirim buat lead yang belum jawab
    │           └── varian follow-up beda dari ice_breaker
    ├── auto: NOTIFICATION muncul kalau ada event
    │     ├── closing_triggers.deal → auto-flag deal (lu verify)
    │     ├── closing_triggers.hot_lead → auto-prioritize
    │     └── closing_triggers.stop → auto-add do_not_contact (lu verify)
    │
    ├── interrupt: RESPONSE masuk → lu approve
    │     ├── deal detected? → auto-mark deal
    │     ├── stop detected? → auto-block
    │     └── COMPOSE → custom reply
    ├── interrupt: ERROR → lu fix
    │     └── VALIDATION → fix config
    ├── interrupt: MILESTONE → lu seneng
    │
    ├── manual: lu bisa akses DATABASE, TEMPLATE, SETTINGS, FOLLOW-UP kapan aja
    ├── manual: l LICENSE → liat / ganti lisensi
    ├── manual: h HISTORY → liat performa masa lalu
    └── manual: v VALIDATION → force check semua config

  LEAD LIFECYCLE:
  baru → wa_validated → ice_breaker_sent → responded → offer_sent → converted
    │         │                 │            │                            ↑
    │         │                 │            └→ negative → archived         │
    │         │                 │            └→ auto_reply → skipped        │
    │         │                 └→ no response → follow_up_1 → follow_up_2│
    │         │                                    │               │      │
    │         │                                    │               └→ cold │
    │         │                                    └→ responded ────┘──────┘
    │         └→ wa_invalid (skip, nggak dikirim)
    │
    │  RE-CONTACT: responded + dingin → 7 hari jeda → re_contact → responded
    │
    │  FOLLOW-UP LIMIT: max 3 pesan lifetime per lead
    │  ice_breaker + follow_up_1 + follow_up_2 = 3 pesan
    │  setelah 2x follow-up tanpa response = dingin
    │
    │  LICENSE GATE: tanpa lisensi valid = full stop
    │  expired → hard stop, data aman
    │  device conflict → hard stop, bisa force transfer
    │
  RESPONSE CLASSIFICATION:                                              │
  incoming → closing_triggers.deal match? ──→ auto-deal ──────────────┘
           → closing_triggers.hot_lead match? ──→ auto-prioritize
           → closing_triggers.stop match? ──→ auto-add do_not_contact
           → positive/curious/negative/maybe/auto-reply ──→ manual classify
```

**Auto-pilot = default. Manual = optional.**
Tidak ada dead end. Setiap screen punya next step. Setiap loop makin efisien.
Config error = partial pause, bukan full stop. Worker yang ok tetap jalan.
WA validation = bukan optimis, tapi realistis. Nomor yang nggak punya WA = skip, bukan gagal.
Closing triggers = data-driven, bukan tebakan. User define pattern-nya, WaClaw eksekusi.

---

## 8. Keyboard Grammar

**Hafal sekali. Pakai selamanya. Tombol yang sama, tiap screen.**

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

**Overlay `?`:**

```

  ── tombol ─────────────

  ↑↓   pindah
  ↵    aksi utama
  1-3  pilih opsi
  s    skip
  q    balik/keluar
  /    cari
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

## 11. Aturan Yang Ga Tertulis

1. **Jangan pernah tampilin empty state tanpa next action.** "Belum ada lead" → "tekan s buat mulai scrape"
2. **Jangan pernah pakai merah buat rejection.** Merah = rusak. Rejection = netral.
3. **Jangan pernah animate cuma buat hiasan.** Tiap animasi = state change yang penting.
4. **Jangan pernah sembunyiin rate limit.** Limit kelihatan = trust. Limit tersembunyi = cemas.
5. **Jangan pernah minta konfirmasi 2x.** `↵` = gas. Percaya user.
6. **Jangan pernah break keyboard grammar.** `q` = balik/keluar. Selalu.
7. **Jangan pernah tampilin angka tanpa konteks.** "4.6% conversion" ← good. "46% selesai" ← selesai apaan?
8. **Jangan pernah pakai bahasa formal.** Netizen indo. Santai tapi jelas.
9. **Jangan pernah buat user nunggu tanpa info.** Kalau nunggu, kasih tau ngapain + berapa lama.
10. **Auto-pilot = default. Manual = bonus.** Kalau user nggak ngapa-ngapain 1 jam dan WaClaw tetep kerja, itu SUKSES.
11. **Config error = partial pause, bukan full stop.** Worker yang ok tetap jalan. Yang error di-pause. Jangan pernah shut down seluruh army cuma karena 1 niche.yaml broken.
12. **Validate early, validate often.** Setiap reload, setiap startup, setiap template change. Silent errors = invisible disaster.
13. **Error message = solusi, bukan cuma masalah.** "parse error baris 14" + "1 buka file" = fixable. "error" saja = frustrating.
14. **Auto-backup sebelum overwrite.** `config.yaml.bak` always. Kalau reload gagal, revert 1 tombol.
15. **Celebration itu earned, bukan diberi.** Conversion screen full drama karena lu GENUNE menang. Kalau semua screen drama, yang beneran penting jadi biasa.

---

## 12. Tech Stack (TUI Layer)

```
  Terminal UI

  bubbletea    ──  MVC framework
  lipgloss     ──  styling & layout
  bubbles      ──  pre-built components
  glamour      ──  markdown rendering
  huh          ──  forms & prompts

  Semua Charm.sh ecosystem. Satu estetika.
```

Kenapa Charm.sh? Render sama di mana aja. Komponen dirancang orang yang ngerti terminal. Estetikanya udah borderless — cuma perlu di-push lebih jauh.

---

## 13. State Machine Summary

```
LEAD STATES:
  baru → ice_breaker_sent → responded → offer_sent → converted
    │         │                 │            │
    │         │                 └→ negative → archived
    │         │                 └→ auto_reply → skipped
    │         └→ no_response → follow_up_1 → follow_up_2 → cold
    │         │                                    │            │
    │         │                                    └→ responded └→ offer_sent
    │         └→ failed → retry (max 3x) → dead
    └→ blocked (manual)

  RE-CONTACT: responded + dingin → 7 hari → re_contact → responded
  FOLLOW-UP LIMIT: ice_breaker + follow_up_1 + follow_up_2 = max 3 per lead
  COLD: 2x follow-up tanpa response → auto-tandai dingin

WORKER STATES (per niche):
  spawning → scraping → qualifying → queuing → sending → idle (loop)
    │            │           │           │          │
    │            └→ error → retry       │          └→ rate_limited
    │                                    └→ paused (manual)
    └→ config_error → paused → fixed → spawning
    └→ stopped (manual)

BATCH PIPELINE (per worker):
  queries.md → scrape → qualify → auto-review → queue → batch send → wait → loop
                                                                  └→ response → offer → deal

CONFIG VALIDATION STATES:
  unchecked → validating → clean | errors | warnings
    │                         │        │         │
    │                         │        └→ fix → revalidate → clean
    │                         └→ auto-dismiss (from boot)
    └→ first_time → generate template → revalidate → clean

SCREEN STATES:
  boot: first_time | returning | returning+response | returning+error | returning+config_error | returning+license_expired | returning+device_conflict
  login: qr_waiting | qr_scanned | login_success | login_expired | login_failed
  niche: niche_list | niche_multi_selected | niche_custom | niche_edit_filters | niche_config_error
  scrape: scraping_active | scraping_multi_active | scraping_multi_staggered | scrape_idle | scrape_empty | scrape_error | scrape_gmaps_limited | scrape_auto_approved | scrape_high_value_reveal | scrape_batch_complete
  review: reviewing | lead_detail | template_preview | queue_complete
  send: sending_active | sending_paused | sending_off_hours | sending_rate_limited | sending_daily_limit | sending_failed | sending_all_slots_down | sending_with_response
  monitor: live_dashboard | idle_background | dashboard_night | dashboard_error | dashboard_empty | dashboard_with_pending_responses
  response: response_positive | response_curious | response_negative | response_maybe | response_auto_reply | offer_preview | response_multi_queue | conversion
  leads: leads_list | leads_filtered | lead_full_detail | lead_follow_up_due | lead_cold | lead_never_contacted | lead_converted
  template: template_list | template_preview | template_edit_hint | template_validation_error
  workers: workers_overview | worker_detail | worker_add_niche | worker_paused
  shield: shield_overview | shield_warning | shield_danger | shield_slot_detail | shield_settings
  settings: settings_overview | settings_edit | settings_reload | settings_reload_error
  validation: validation_clean | validation_errors | validation_warnings | validation_fix | validation_first_time
  compose: compose_draft | compose_preview | compose_template_pick
  history: history_today | history_week | history_day_detail
  followup: followup_dashboard | followup_niche_detail | followup_sending | followup_empty | followup_cold_list | followup_recontact
  explorer: explorer_browse | explorer_search | explorer_category_detail | explorer_generating | explorer_generated
  update: update_available | update_downloading | update_ready | upgrade_available | upgrade_license_input | license_expired_with_upgrade
  license: license_input | license_validating | license_valid | license_invalid | license_expired | license_device_conflict | license_server_error
  nerd_stats: hidden | minimal | expanded (global overlay, not a screen)
  cmd_palette: cmd_closed | cmd_open | cmd_executing | cmd_empty | cmd_with_recent | cmd_quick_action (global overlay, not a screen)

NOTIFICATION TYPES:
  response_masuk | multi_response | scrape_selesai | batch_selesai | wa_disconnect |
  wa_flag | health_drop | limit_harian | streak_milestone | config_error | validation_error |
  license_expired | device_conflict | followup_terjadwal | lead_dingin | update_available | upgrade_available

CONFIRMATION OVERLAYS:
  bulk_offer | bulk_delete | bulk_archive | force_device_disconnect
```

---

*WaClaw — army lu cerdas, keras, dan aman.*
*Broken config? Ketahuan. Multiple responses? Ditangani. History? Ada.*
*Follow-up? Persistent. Lisensi? Satu kunci satu device. Niche bingung? Explorer.*
*Update kecil? Gratis. Upgrade besar? Pilihan lu. Nerd stats? Toggle kapan aja.*
*Lupa shortcut? Ctrl+K. Semua command, dari mana aja, 3 detik.*
*Lu cuma nonton. WaClaw yang kerja. Tapi kalau lu mau intervene, 1 tombol cukup.*
