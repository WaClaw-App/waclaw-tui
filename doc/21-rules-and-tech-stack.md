
## 11. Aturan Yang Ga Tertulis

1. **Jangan pernah tampilin empty state tanpa next action.** "Belum ada lead" → "tekan s buat mulai scrape"
2. **Jangan pernah pakai merah buat rejection.** Merah = rusak. Rejection = netral.
3. **Jangan pernah animate cuma buat hiasan.** Tiap animasi = state change yang penting.
4. **Jangan pernah sembunyiin rate limit.** Limit kelihatan = trust. Limit tersembunyi = cemas.
5. **Jangan pernah minta konfirmasi 2x.** `↵` = gas. Percaya user.
6. **Jangan pernah break keyboard grammar.** `q` = balik/keluar. Selalu.
7. **Jangan pernah tampilin angka tanpa konteks.** "4.6% conversion" ← good. "46% selesai" ← selesai apaan?
8. **Jangan pernah pakai bahasa formal.** Casual Indonesian (default) atau casual English. Santai tapi jelas. Keduanya harus available via i18n.
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
  i18n         ──  multi-language (casual Indonesian + casual English)

  Semua Charm.sh ecosystem. Satu estetika.
```

Kenapa Charm.sh? Render sama di mana aja. Komponen dirancang orang yang ngerti terminal. Estetikanya udah borderless — cuma perlu di-push lebih jauh.

---

## 13. Tech Stack (REST API — Future Web Frontend)

```
  REST API

  net/http     ──  standard library HTTP server
  swagger      ──  OpenAPI 3.0 spec auto-generated from pkg/protocol types
  json         ──  encoding/json (same protocol types as JSON-RPC)
  cors         ──  CORS middleware for web frontend
```

**Kenapa REST di atas JSON-RPC?**

Backend punya satu scenario engine. TUI berkomunikasi via JSON-RPC over stdio. Web frontend (masa depan) berkomunikasi via REST API over HTTP. Keduanya pakai types yang sama dari `pkg/protocol/`. Swagger spec di-auto-generate dari protocol structs — jadi web frontend selalu sinkron dengan TUI.

**REST Endpoint Mapping** (same as RPC methods):

| RPC Method | REST Endpoint | Method |
|------------|---------------|--------|
| `navigate` | `POST /api/v1/navigate` | Backend → Frontend (SSE/WebSocket) |
| `update` | `POST /api/v1/update` | Backend → Frontend (SSE/WebSocket) |
| `notify` | `POST /api/v1/notify` | Backend → Frontend (SSE/WebSocket) |
| `validate` | `POST /api/v1/validate` | Backend → Frontend (SSE/WebSocket) |
| `key_press` | `POST /api/v1/events/keypress` | Frontend → Backend |
| `action` | `POST /api/v1/events/action` | Frontend → Backend |
| `request` | `POST /api/v1/events/request` | Frontend → Backend |
| — | `GET /api/v1/state` | Current full state snapshot |
| — | `GET /api/v1/screens` | Screen list + metadata |
| — | `GET /api/v1/swagger.json` | Swagger spec |

Swagger spec auto-generated dari `pkg/protocol/` types via `go:generate` directive. File: `pkg/api/openapi.yaml`.

---
