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

