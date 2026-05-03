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
Language   = Casual Indonesian (default) + Casual English — i18n, bukan hardcoded
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

**TUI support dua bahasa: casual Indonesian (default) dan casual English.**

Tone untuk kedua bahasa: santai, tapi nggak alay. Kayak ngobrol sama rekan kerja yang competent. Singkat, jelas, ga pake basa-basi. Error message juga santai — masalahnya serius, bahasanya nggak usah kaku.

Bahasa di-set di `~/.waclaw/config.yaml`:
```yaml
locale: id  # "id" (default) atau "en"
```

Bisa diganti runtime lewat Ctrl+K command palette → search "language" atau "bahasa". Tombol `l` buat navigasi ke License screen, bukan ganti bahasa.

Semua display string disimpen di `internal/tui/i18n/` — dua map: `id.go` (Indonesian) dan `en.go` (English). Kode tetap English semua. Yang Indonesia cuma yang user liat di layar.

Contoh perbandingan:
| Context | `id` (Indonesian) | `en` (English) |
|---------|-------------------|----------------|
| Proceed | `↵ gas` | `↵ go` |
| Cancel | `s batal` | `s cancel` |
| Back | `q kembali` | `q back` |
| Active | `● aktif` | `● active` |
| Sent | `✓ terkirim` | `✓ sent` |
| Delivered | `✓✓ sampai` | `✓✓ delivered` |
| Shield healthy | `SEHAT` | `HEALTHY` |
| Shield danger | `BAHAYA` | `DANGER` |
| Cold lead | `❄ DINGIN` | `❄ COLD` |
| Session end | `sampai jumpa.` | `see you.` |

---
