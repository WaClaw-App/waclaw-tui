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

