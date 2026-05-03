
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
