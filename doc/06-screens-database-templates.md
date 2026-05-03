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

