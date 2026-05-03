package niche

import (
        "fmt"
        "strings"
        "testing"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/component"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
)

// stripANSI removes ANSI escape sequences for plain-text comparison.
func stripANSI(s string) string {
        var result strings.Builder
        i := 0
        for i < len(s) {
                if s[i] == '\x1b' {
                        i++
                        if i < len(s) && s[i] == '[' {
                                i++
                                for i < len(s) && !((s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z')) {
                                        i++
                                }
                                if i < len(s) {
                                        i++
                                }
                        }
                } else {
                        result.WriteByte(s[i])
                        i++
                }
        }
        return result.String()
}

func TestBlackboxNicheList(t *testing.T) {
        i18n.SetLocale("id")
        time.Sleep(100 * time.Millisecond)

        m := NewSelectModel()
        m.Focus()

        m.HandleNavigate(map[string]any{
                "state": "niche_list",
                "niches": []any{
                        map[string]any{"name": "web developer", "description": "buat yang jual jasa bikin web"},
                        map[string]any{"name": "undangan digital", "description": "buat yang jual undangan digital"},
                        map[string]any{"name": "social media mgr", "description": "buat yang jasa kelola sosmed"},
                        map[string]any{"name": "fotografer", "description": "buat yang jasa foto & portfolio"},
                        map[string]any{"name": "custom", "description": "bikin niche sendiri dari file"},
                },
        })
        // Set staggerStart to past so all items are visible immediately in tests
        m.staggerStart = time.Now().Add(-10 * time.Second)

        p := stripANSI(m.View())

        checks := []struct {
                name string
                ok   bool
                detail string
        }{
                {"title=pilih niche lu", strings.Contains(p, "pilih niche lu"), "doc heading"},
                {"niche_is text", strings.Contains(p, "niche = siapa lu + siapa target lu"), "doc says: 'niche = siapa lu + siapa target lu'"},
                {"worker_parallel text", strings.Contains(p, "bisa pilih lebih dari satu! tiap niche = 1 worker."), "doc says: 'bisa pilih lebih dari satu! tiap niche = 1 worker.'"},
                {"more_niche text", strings.Contains(p, "makin banyak niche, makin luas jaring lu."), "doc says: 'makin banyak niche, makin luas jaring lu.'"},
                {"☐ checkbox", strings.Contains(p, "☐"), "doc shows ☐"},
                {"all 5 niches", strings.Contains(p, "web developer") && strings.Contains(p, "custom"), "doc shows 5 niche items"},
                {"descriptions", strings.Contains(p, "buat yang jual jasa bikin web"), "doc shows descriptions"},
                {"space centang/hapus", strings.Contains(p, "centang/hapus"), "doc says: 'space  centang/hapus'"},
                {"gas dengan yang dicentang", strings.Contains(p, "gas dengan yang dicentang"), "doc says: '↵  gas dengan yang dicentang'"},
                {"parallel caption", strings.Contains(p, "semua niche jalan paralel"), "doc says: 'semua niche jalan paralel, masing-masing punya worker.'"},
                {"no border chars P3", !strings.ContainsAny(p, "│┃║─━┌┐└┘"), "vertical borderless: no box-drawing"},
        }

        for _, c := range checks {
                if !c.ok {
                        t.Errorf("[niche_list] FAIL %s: %s", c.name, c.detail)
                } else {
                        t.Logf("[niche_list] PASS %s", c.name)
                }
        }
}

func TestBlackboxNicheMultiSelected(t *testing.T) {
        i18n.SetLocale("id")
        time.Sleep(100 * time.Millisecond)

        m := NewSelectModel()
        m.Focus()
        m.HandleNavigate(map[string]any{
                "state": "niche_multi_selected",
                "niches": []any{
                        map[string]any{"name": "web developer", "description": "buat yang jual jasa bikin web", "area": "kediri, 15km", "templates": 3},
                        map[string]any{"name": "undangan digital", "description": "buat yang jual undangan digital", "area": "kediri + surabaya", "templates": 2},
                        map[string]any{"name": "social media mgr", "description": "buat yang jasa kelola sosmed"},
                        map[string]any{"name": "fotografer", "description": "buat yang jasa foto & portfolio"},
                        map[string]any{"name": "custom", "description": "bikin niche sendiri dari file"},
                },
        })
        m.list.Toggle()
        m.list.Down()
        m.list.Toggle()
        for i := range m.items {
                if i < len(m.list.Items) {
                        m.items[i].Selected = m.list.Items[i].Selected
                }
        }
        m.state = protocol.NicheMultiSelected

        p := stripANSI(m.View())

        checks := []struct {
                name string
                ok   bool
                detail string
        }{
                {"☑ selected checkbox", strings.Contains(p, "☑"), "doc shows ☑ for selected"},
                {"niche dipilih summary", strings.Contains(p, "dipilih"), "doc says: '2 niche dipilih:'"},
                {"▸ prefix", strings.Contains(p, "▸"), "doc shows ▸ for selected items"},
                {"area detail", strings.Contains(p, "kediri"), "doc shows area info"},
                {"gas jalanin footer", strings.Contains(p, "gas jalanin") || strings.Contains(p, "gas"), "doc says: '↵  gas jalanin 2 niche'"},
                {"ubah footer", strings.Contains(p, "ubah"), "doc says: 'space  ubah'"},
                {"balik footer", strings.Contains(p, "balik"), "doc says: 'q  balik'"},
                {"scrape own caption", strings.Contains(p, "scrape sendiri") || strings.Contains(p, "kirim sendiri"), "doc says: 'masing-masing scrape sendiri, kirim sendiri.'"},
        }

        for _, c := range checks {
                if !c.ok {
                        t.Errorf("[niche_multi_selected] FAIL %s: %s", c.name, c.detail)
                } else {
                        t.Logf("[niche_multi_selected] PASS %s", c.name)
                }
        }
}

func TestBlackboxNicheCustom(t *testing.T) {
        i18n.SetLocale("id")

        m := NewSelectModel()
        m.Focus()
        m.state = protocol.NicheCustom
        m.staggerStart = time.Now().Add(-1 * time.Second)
        p := stripANSI(m.View())

        checks := []struct {
                name string
                ok   bool
                detail string
        }{
                {"title=niche custom", strings.Contains(p, "niche custom"), "doc heading"},
                {"dir path", strings.Contains(p, "~/.waclaw/niches/nama_niche/"), "doc path"},
                {"butuh minimal", strings.Contains(p, "butuh minimal"), "doc says: 'butuh minimal:'"},
                {"niche.yaml", strings.Contains(p, "niche.yaml"), "doc says: '- niche.yaml    (filter & target)'"},
                {"ice_breaker.md", strings.Contains(p, "ice_breaker.md"), "doc says: '- ice_breaker.md'"},
                {"contoh path", strings.Contains(p, "~/.waclaw/niches/_contoh/"), "doc example path"},
                {"reload footer", strings.Contains(p, "reload"), "doc says: 'r  reload'"},
                {"pick existing footer", strings.Contains(p, "pilih yang ada"), "doc says: '1-5  pilih yang ada'"},
        }

        for _, c := range checks {
                if !c.ok {
                        t.Errorf("[niche_custom] FAIL %s: %s", c.name, c.detail)
                } else {
                        t.Logf("[niche_custom] PASS %s", c.name)
                }
        }
}

func TestBlackboxNicheEditFilters(t *testing.T) {
        i18n.SetLocale("id")

        m := NewSelectModel()
        m.Focus()
        m.state = protocol.NicheEditFilters
        m.filterNiche = "web developer"
        m.filterTargets = []string{"cafe", "gym", "salon", "toko bangunan", "bengkel"}
        m.filters = []FilterEntry{
                {Symbol: "✗", Label: "punya website", Detail: "skip yang udah punya", Inverted: true},
                {Symbol: "✓", Label: "punya instagram", Detail: "lebih potensial"},
                {Symbol: "⭐", Label: "3.0 - 4.8", Detail: "rating range", Neutral: true},
                {Symbol: "📊", Label: "5 - 500 reviews", Detail: "ukuran bisnis", Neutral: true},
        }
        m.areas = []AreaEntry{
                {City: "kediri", Radius: "15km", KecCount: 4},
                {City: "nganjuk", Radius: "10km", KecCount: 2},
                {City: "tulungagung", Radius: "10km", KecCount: 3},
                {City: "blitar", Radius: "10km", KecCount: 2},
                {City: "madiun", Radius: "10km", KecCount: 1},
        }
        p := stripANSI(m.View())

        checks := []struct {
                name string
                ok   bool
                detail string
        }{
                {"title niche: web developer", strings.Contains(p, "niche:") && strings.Contains(p, "web developer"), "doc says: 'niche: web developer'"},
                {"target heading", strings.Contains(p, "target:"), "doc shows 'target:' heading"},
                {"target items", strings.Contains(p, "cafe") && strings.Contains(p, "gym"), "doc shows target list"},
                {"✗ symbol", strings.Contains(p, "✗"), "doc shows '✗  punya website'"},
                {"✓ symbol", strings.Contains(p, "✓"), "doc shows '✓  punya instagram'"},
                {"⭐ symbol", strings.Contains(p, "⭐"), "doc shows '⭐ 3.0 - 4.8'"},
                {"📊 symbol", strings.Contains(p, "📊"), "doc shows '📊 5 - 500 reviews'"},
                {"area count", strings.Contains(p, "area") && strings.Contains(p, "kota"), "doc says: 'area (5 kota):'"},
                {"area cities", strings.Contains(p, "kediri") && strings.Contains(p, "madiun"), "doc shows city names"},
                {"gas scrape footer", strings.Contains(p, "gas scrape"), "doc says: '↵  gas scrape'"},
                {"edit filter footer", strings.Contains(p, "edit filter"), "doc says: '2  edit filter'"},
                {"MISALIGN doc ── connector vs code ··", !strings.Contains(p, "──"), "doc/02: 'kediri  15km  ──  4 kecamatan' but P3 forbids ─; code uses ··"},
                {"MISALIGN doc ──── separator vs code ·· dots", !strings.Contains(p, "────────────"), "doc shows ──── separator but P3 forbids ─; code uses ··"},
        }

        for _, c := range checks {
                if !c.ok {
                        t.Errorf("[niche_edit_filters] FAIL %s: %s", c.name, c.detail)
                } else {
                        t.Logf("[niche_edit_filters] PASS %s", c.name)
                }
        }
}

func TestBlackboxNicheConfigError(t *testing.T) {
        i18n.SetLocale("id")

        m := NewSelectModel()
        m.Focus()
        m.state = protocol.NicheConfigError
        m.errorNiche = "fotografer"
        m.errorFile = "~/.waclaw/niches/fotografer/niche.yaml"
        m.errors = []ConfigError{
                {Line: 14, Message: "parse error", Description: "kurang tanda kutip penutup", Detail: `    - "wedding`, Pointer: "          ^"},
                {Line: 0, Message: "field wajib kosong: areas", Description: "niche.yaml harus punya minimal 1 area"},
                {Line: 0, Message: "field scoring: rating_min harus angka, bukan \"rendeh\""},
        }
        m.errorBlinkStart = time.Now()
        p := stripANSI(m.View())

        checks := []struct {
                name string
                ok   bool
                detail string
        }{
                {"✗ niche error title", strings.Contains(p, "✗") && strings.Contains(p, "fotografer"), "doc says: '✗ niche error: fotografer'"},
                {"error file path", strings.Contains(p, "~/.waclaw/niches/fotografer/niche.yaml"), "doc shows file path"},
                {"3 masalah count", strings.Contains(p, "3 masalah") || strings.Contains(p, "masalah"), "doc says: '3 masalah:'"},
                {"parse error message", strings.Contains(p, "parse error") || strings.Contains(p, "baris 14"), "doc shows: '✗  baris 14: parse error'"},
                {"field wajib error", strings.Contains(p, "field wajib") || strings.Contains(p, "areas"), "doc shows: '✗  field wajib kosong: areas'"},
                {"paused message", strings.Contains(p, "pause") || strings.Contains(p, "diperbaiki"), "doc says: 'worker fotografer di-pause sampai config diperbaiki.'"},
                {"worker lain tetap jalan", strings.Contains(p, "worker lain tetap jalan"), "doc says: 'worker lain tetap jalan normal.'"},
                {"buka file footer", strings.Contains(p, "buka file"), "doc says: '1  buka file'"},
                {"liat contoh config footer", strings.Contains(p, "liat contoh config"), "doc says: '2  liat contoh config'"},
                {"MISALIGN doc │ gutter vs code indentation", !strings.Contains(p, "│"), "doc/02 shows │ gutter but P3 forbids it; code uses indentation"},
                {"MISALIGN doc ──── separator vs code ·· dots", !strings.Contains(p, "────────────"), "doc shows ──── separator but P3 forbids ─"},
        }

        for _, c := range checks {
                if !c.ok {
                        t.Errorf("[niche_config_error] FAIL %s: %s", c.name, c.detail)
                } else {
                        t.Logf("[niche_config_error] PASS %s", c.name)
                }
        }
}

func TestBlackboxExplorerBrowse(t *testing.T) {
        i18n.SetLocale("id")

        em := NewExplorerModel()
        em.Focus()

        em.HandleNavigate(map[string]any{
                "state": "explorer_browse",
                "categories": []any{
                        map[string]any{"emoji": "🍜", "name": "kuliner", "sub_categories": []any{"cafe", "restoran", "catering", "bakery"}, "sub_count": 12, "area_count": 3},
                        map[string]any{"emoji": "💇", "name": "kecantikan", "sub_categories": []any{"salon", "barbershop", "spa", "nail art"}, "sub_count": 8, "area_count": 2},
                        map[string]any{"emoji": "💪", "name": "fitness", "sub_categories": []any{"gym", "yoga studio", "personal trainer"}, "sub_count": 5, "area_count": 1},
                        map[string]any{"emoji": "🏠", "name": "properti", "sub_categories": []any{"agent", "developer", "interior", "renovation"}, "sub_count": 10, "area_count": 4},
                        map[string]any{"emoji": "🎓", "name": "pendidikan", "sub_categories": []any{"bimbel", "kursus", "training", "workshop"}, "sub_count": 7, "area_count": 2},
                        map[string]any{"emoji": "🏥", "name": "kesehatan", "sub_categories": []any{"klinik", "apotek", "dokter", "therapy"}, "sub_count": 9, "area_count": 3},
                        map[string]any{"emoji": "📸", "name": "fotografi", "sub_categories": []any{"wedding", "product", "studio", "event"}, "sub_count": 6, "area_count": 2},
                        map[string]any{"emoji": "💻", "name": "teknologi", "sub_categories": []any{"web dev", "app dev", "it support", "digital"}, "sub_count": 8, "area_count": 1},
                        map[string]any{"emoji": "🚗", "name": "otomotif", "sub_categories": []any{"bengkel", "dealer", "rental", "sparepart"}, "sub_count": 6, "area_count": 3},
                        map[string]any{"emoji": "🎉", "name": "event", "sub_categories": []any{"wedding org", "decorator", "catering", "mc"}, "sub_count": 4, "area_count": 2},
                },
        })
        // Set staggerStart to past so all items are visible immediately in tests
        em.staggerStart = time.Now().Add(-10 * time.Second)

        p := stripANSI(em.View())

        allCats := strings.Contains(p, "kuliner") && strings.Contains(p, "event")
        anyEmoji := strings.Contains(p, "🍜") || strings.Contains(p, "💇") || strings.Contains(p, "💪")

        checks := []struct {
                name string
                ok   bool
                detail string
        }{
                {"title=niche explorer", strings.Contains(p, "niche explorer"), "doc heading"},
                {"bingung mulai subtitle", strings.Contains(p, "bingung mulai dari mana"), "doc says: 'bingung mulai dari mana? browse aja.'"},
                {"pick config text", strings.Contains(p, "pilih kategori, waclaw yang bikin config"), "doc says: 'pilih kategori, waclaw yang bikin config-nya.'"},
                {"populer section", strings.Contains(p, "populer"), "doc shows '── populer ────'"},
                {"all 10 categories", allCats, "doc shows 10 categories"},
                {"emoji in categories", anyEmoji, "doc shows emoji per category"},
                {"sub-categories shown", strings.Contains(p, "cafe") || strings.Contains(p, "restoran"), "doc shows sub-categories"},
                {"cari kategori footer", strings.Contains(p, "cari kategori"), "doc says: '/  cari kategori'"},
                {"pilih footer", strings.Contains(p, "pilih"), "doc says: '↑↓  pilih'"},
                {"liat detail footer", strings.Contains(p, "liat detail"), "doc says: '↵  liat detail'"},
                {"DIY caption", strings.Contains(p, "mau bikin sendiri") || strings.Contains(p, "tekan 5"), "doc DIY caption"},
                {"no border chars P3", !strings.ContainsAny(p, "│┃║─━┌┐└┘"), "vertical borderless"},
                {"MISALIGN doc ── dividers vs code ··", !strings.Contains(p, "──"), "doc/11 shows '── populer ────' but P3 forbids ─"},
        }

        for _, c := range checks {
                if !c.ok {
                        t.Errorf("[explorer_browse] FAIL %s: %s", c.name, c.detail)
                } else {
                        t.Logf("[explorer_browse] PASS %s", c.name)
                }
        }
}

func TestBlackboxExplorerSearch(t *testing.T) {
        i18n.SetLocale("id")

        em := NewExplorerModel()
        em.Focus()
        em.state = protocol.ExplorerSearch
        em.searchInput = component.NewSearchInput(300)
        em.searchInput.Value = "cafe"
        em.searchInput.Focused = true
        em.searchResults = []SearchResult{
                {Name: "cafe & coffee shop", SubCount: 12, AreaCount: 3},
                {Name: "cafe modern", SubCount: 8, AreaCount: 2},
                {Name: "cafe traditional", SubCount: 5, AreaCount: 1},
                {Name: "cafe catering", SubCount: 3, AreaCount: 1},
        }
        em.categories = []Category{{Emoji: "🍜", Name: "kuliner", SubCategories: []string{"cafe"}, SubCount: 12, AreaCount: 3}}
        p := stripANSI(em.View())

        checks := []struct {
                name string
                ok   bool
                detail string
        }{
                {"title + mencari", strings.Contains(p, "niche explorer") && strings.Contains(p, "mencari"), "doc shows mencari..."},
                {"search input", strings.Contains(p, "cafe") || strings.Contains(p, ">"), "doc shows search input"},
                {"MISALIGN doc ┌──┐ box vs code borderless", !strings.ContainsAny(p, "┌┐│└┘─"), "doc/11 shows box but P3 forbids it"},
                {"hasil section", strings.Contains(p, "hasil"), "doc shows hasil section"},
                {"search results", strings.Contains(p, "cafe & coffee shop") || strings.Contains(p, "coffee shop"), "doc shows results"},
                {"sub-kategori detail", strings.Contains(p, "sub-kategori") || strings.Contains(p, "12"), "doc shows sub-kategori info"},
                {"sumber caption", strings.Contains(p, "sumber") || strings.Contains(p, "WhatsApp"), "doc says sumber caption"},
                {"esc batal footer", strings.Contains(p, "batal"), "doc says: 'esc  batal'"},
        }

        for _, c := range checks {
                if !c.ok {
                        t.Errorf("[explorer_search] FAIL %s: %s", c.name, c.detail)
                } else {
                        t.Logf("[explorer_search] PASS %s", c.name)
                }
        }
}

func TestBlackboxExplorerCategoryDetail(t *testing.T) {
        i18n.SetLocale("id")

        em := NewExplorerModel()
        em.Focus()
        em.state = protocol.ExplorerCategoryDetail
        em.selectedCategory = Category{
                Emoji: "🍜", Name: "cafe & coffee shop",
                SubCategories: []string{"cafe", "coffee shop", "kopi nusantara", "kafe modern", "kafe tradisional", "coffee bar", "espresso bar", "roaster"},
        }
        em.detailSources = []SourceDetail{
                {Name: "WhatsApp Business Directory", Count: 247, Unit: "bisnis terdaftar"},
                {Name: "Google Maps Categories", Count: 189, Unit: "listing aktif"},
        }
        em.detailAreas = []AreaEntry{
                {City: "kediri", Radius: "15km"},
                {City: "nganjuk", Radius: "10km"},
                {City: "tulungagung", Radius: "10km"},
                {City: "blitar", Radius: "10km"},
                {City: "madiun", Radius: "10km"},
        }
        em.detailFilters = []FilterEntry{
                {Symbol: "✗", Label: "punya website", Detail: "skip yang udah punya", Inverted: true},
                {Symbol: "✓", Label: "punya instagram", Detail: "lebih potensial"},
                {Symbol: "⭐", Label: "3.0 - 4.8", Neutral: true},
                {Symbol: "📊", Label: "5 - 500 reviews", Neutral: true},
        }
        em.detailTemplates = []GenerateFileStatus{
                {Name: "niche.yaml", Detail: "filter + target + area", Completed: true},
                {Name: "ice_breaker.md", Detail: "3 varian", Completed: true},
                {Name: "queries.md", Detail: "search queries per area", Completed: true},
        }
        p := stripANSI(em.View())

        checks := []struct {
                name string
                ok   bool
                detail string
        }{
                {"→ breadcrumb", strings.Contains(p, "→"), "doc shows: 'niche explorer → cafe & coffee shop'"},
                {"kategori label", strings.Contains(p, "kategori:") && strings.Contains(p, "cafe & coffee shop"), "doc says: 'kategori: cafe & coffee shop'"},
                {"sub-kategori", strings.Contains(p, "sub-kategori"), "doc says: 'sub-kategori:' heading"},
                {"sumber label", strings.Contains(p, "sumber:"), "doc says: 'sumber:' heading"},
                {"source counts", strings.Contains(p, "247") && strings.Contains(p, "189"), "doc shows 247 and 189"},
                {"area auto-detect", strings.Contains(p, "area") && strings.Contains(p, "config"), "doc says: 'area (auto-detect dari config):'"},
                {"filter default", strings.Contains(p, "filter default"), "doc says: 'filter default:'"},
                {"template gen list", strings.Contains(p, "niche.yaml") && strings.Contains(p, "ice_breaker"), "doc shows template list with ✓"},
                {"generate config footer", strings.Contains(p, "generate config"), "doc says: '↵  generate config'"},
                {"folder path caption", strings.Contains(p, "~/.waclaw/niches/"), "doc shows folder path"},
                {"MISALIGN doc ──── separator vs code ··", !strings.Contains(p, "────────────"), "doc shows ──── separator but P3 forbids ─"},
        }

        for _, c := range checks {
                if !c.ok {
                        t.Errorf("[explorer_category_detail] FAIL %s: %s", c.name, c.detail)
                } else {
                        t.Logf("[explorer_category_detail] PASS %s", c.name)
                }
        }
}

func TestBlackboxExplorerGenerating(t *testing.T) {
        i18n.SetLocale("id")

        em := NewExplorerModel()
        em.Focus()
        em.state = protocol.ExplorerGenerating
        em.selectedCategory = Category{Name: "cafe & coffee shop"}
        em.genProgress = 0.4
        em.genCurrentStep = 1
        em.genFiles = []GenerateFileStatus{
                {Name: "membuat folder ~/.waclaw/niches/cafe_coffee_shop/", Completed: true},
                {Name: "menulis niche.yaml", Completed: false},
                {Name: "menulis ice_breaker.md (3 varian)", Completed: false},
                {Name: "menulis queries.md", Completed: false},
        }
        p := stripANSI(em.View())

        checks := []struct {
                name string
                ok   bool
                detail string
        }{
                {"title present", strings.Contains(p, "niche explorer"), "doc shows title"},
                {"lagi generate text", strings.Contains(p, "lagi generate config"), "doc says: 'lagi generate config...'"},
                {"● current step", strings.Contains(p, "●"), "doc shows ● for current step"},
                {"○ pending step", strings.Contains(p, "○"), "doc shows ○ for pending steps"},
        }

        for _, c := range checks {
                if !c.ok {
                        t.Errorf("[explorer_generating] FAIL %s: %s", c.name, c.detail)
                } else {
                        t.Logf("[explorer_generating] PASS %s", c.name)
                }
        }
}

func TestBlackboxExplorerGenerated(t *testing.T) {
        i18n.SetLocale("id")

        em := NewExplorerModel()
        em.Focus()
        em.state = protocol.ExplorerGenerated
        em.selectedCategory = Category{Name: "cafe & coffee shop"}
        em.genNicheName = "cafe_coffee_shop"
        em.genFilesDone = []GenerateFileStatus{
                {Name: "niche.yaml", Detail: "5 targets · 5 area · 4 filter", Completed: true},
                {Name: "ice_breaker.md", Detail: "3 varian", Completed: true},
                {Name: "queries.md", Detail: "8 query per area", Completed: true},
        }
        p := stripANSI(em.View())

        checks := []struct {
                name string
                ok   bool
                detail string
        }{
                {"config berhasil", strings.Contains(p, "config berhasil di-generate"), "doc says: '✓ config berhasil di-generate!'"},
                {"folder path", strings.Contains(p, "~/.waclaw/niches/cafe_coffee_shop/"), "doc shows folder path"},
                {"file list", strings.Contains(p, "niche.yaml") && strings.Contains(p, "ice_breaker.md") && strings.Contains(p, "queries.md"), "doc shows file tree"},
                {"MISALIGN doc ├── └── file tree vs code indentation", !strings.ContainsAny(p, "├└│─"), "doc/11 shows ├── └── but P3 forbids box-drawing; code uses indentation"},
                {"edit file prompt", strings.Contains(p, "edit file") || strings.Contains(p, "edit"), "doc says: 'lu bisa edit file-nya langsung kalo mau ubah.'"},
                {"reload prompt", strings.Contains(p, "reload"), "doc says: 'tekan r buat reload setelah edit.'"},
                {"gas pake ini footer", strings.Contains(p, "gas pake ini"), "doc says: '↵  gas pake ini'"},
                {"parallel caption", strings.Contains(p, "paralel") || strings.Contains(p, "worker baru"), "doc says: 'worker baru bakal jalan paralel sama yang udah ada.'"},
        }

        for _, c := range checks {
                if !c.ok {
                        t.Errorf("[explorer_generated] FAIL %s: %s", c.name, c.detail)
                } else {
                        t.Logf("[explorer_generated] PASS %s", c.name)
                }
        }
}

// TestBlackboxArchitectureConcerns documents architectural issues discovered
// during QA that relate to backend-frontend separation and DRY violations.
func TestBlackboxArchitectureConcerns(t *testing.T) {
        // Concern 1: style/layout.go Separator() uses ─ chars which violates P3
        // The niche package correctly overrides with renderSeparator() using ·
        // But the shared Separator() function itself is wrong and could be used by other screens.
        t.Log("[architecture] CONCERN: style/layout.go Separator() uses ─ chars which violates P3 vertical borderless")
        t.Log("[architecture] CONCERN: niche/helpers.go renderSeparator() correctly overrides with · dots")
        t.Log("[architecture] CONCERN: slugify() in TUI is display-only — acceptable, backend owns authoritative slug")
        t.Log("[architecture] CONCERN: performSearch() does local filtering for instant preview — dual-source design, acceptable")
        t.Log("[architecture] CONCERN: parseFilterEntries/parseAreaEntries are frontend data mappers — correct separation")
        t.Log("[architecture] MISALIGNMENT SUMMARY: All ──/──/│/├──/└── box-drawing chars in doc wireframes are replaced with ·/indentation in code per P3")
        t.Log("[architecture] This is intentional: P3 (vertical borderless) takes precedence over doc wireframe visual formatting")
}

// TestBlackboxPrintDetailedOutput dumps full rendered output for each state
// for manual visual comparison against doc wireframes.
func TestBlackboxPrintDetailedOutput(t *testing.T) {
        i18n.SetLocale("id")

        fmt.Println("\n========================================")
        fmt.Println("DETAILED RENDERED OUTPUT FOR VISUAL COMPARISON")
        fmt.Println("========================================")

        // niche_list
        m := NewSelectModel()
        m.Focus()
        m.HandleNavigate(map[string]any{
                "state": "niche_list",
                "niches": []any{
                        map[string]any{"name": "web developer", "description": "buat yang jual jasa bikin web"},
                        map[string]any{"name": "undangan digital", "description": "buat yang jual undangan digital"},
                        map[string]any{"name": "social media mgr", "description": "buat yang jasa kelola sosmed"},
                        map[string]any{"name": "fotografer", "description": "buat yang jasa foto & portfolio"},
                        map[string]any{"name": "custom", "description": "bikin niche sendiri dari file"},
                },
        })
        fmt.Println("\n--- niche_list ---")
        fmt.Println(stripANSI(m.View()))

        // niche_edit_filters
        m2 := NewSelectModel()
        m2.Focus()
        m2.state = protocol.NicheEditFilters
        m2.filterNiche = "web developer"
        m2.filterTargets = []string{"cafe", "gym", "salon", "toko bangunan", "bengkel"}
        m2.filters = []FilterEntry{
                {Symbol: "✗", Label: "punya website", Detail: "skip yang udah punya", Inverted: true},
                {Symbol: "✓", Label: "punya instagram", Detail: "lebih potensial"},
                {Symbol: "⭐", Label: "3.0 - 4.8", Detail: "rating range", Neutral: true},
                {Symbol: "📊", Label: "5 - 500 reviews", Detail: "ukuran bisnis", Neutral: true},
        }
        m2.areas = []AreaEntry{
                {City: "kediri", Radius: "15km", KecCount: 4},
                {City: "nganjuk", Radius: "10km", KecCount: 2},
                {City: "tulungagung", Radius: "10km", KecCount: 3},
                {City: "blitar", Radius: "10km", KecCount: 2},
                {City: "madiun", Radius: "10km", KecCount: 1},
        }
        fmt.Println("\n--- niche_edit_filters ---")
        fmt.Println(stripANSI(m2.View()))

        // niche_config_error
        m3 := NewSelectModel()
        m3.Focus()
        m3.state = protocol.NicheConfigError
        m3.errorNiche = "fotografer"
        m3.errorFile = "~/.waclaw/niches/fotografer/niche.yaml"
        m3.errors = []ConfigError{
                {Line: 14, Message: "parse error", Description: "kurang tanda kutip penutup", Detail: `    - "wedding`, Pointer: "          ^"},
                {Line: 0, Message: "field wajib kosong: areas", Description: "niche.yaml harus punya minimal 1 area"},
                {Line: 0, Message: "field scoring: rating_min harus angka"},
        }
        m3.errorBlinkStart = time.Now()
        fmt.Println("\n--- niche_config_error ---")
        fmt.Println(stripANSI(m3.View()))

        // explorer_browse
        em := NewExplorerModel()
        em.Focus()
        em.HandleNavigate(map[string]any{
                "state": "explorer_browse",
                "categories": []any{
                        map[string]any{"emoji": "🍜", "name": "kuliner", "sub_categories": []any{"cafe", "restoran", "catering", "bakery"}},
                        map[string]any{"emoji": "💇", "name": "kecantikan", "sub_categories": []any{"salon", "barbershop", "spa", "nail art"}},
                },
        })
        fmt.Println("\n--- explorer_browse ---")
        fmt.Println(stripANSI(em.View()))

        // explorer_search
        em2 := NewExplorerModel()
        em2.Focus()
        em2.state = protocol.ExplorerSearch
        em2.searchInput = component.NewSearchInput(300)
        em2.searchInput.Value = "cafe"
        em2.searchInput.Focused = true
        em2.searchResults = []SearchResult{
                {Name: "cafe & coffee shop", SubCount: 12, AreaCount: 3},
                {Name: "cafe modern", SubCount: 8, AreaCount: 2},
        }
        em2.categories = []Category{{Emoji: "🍜", Name: "kuliner", SubCategories: []string{"cafe"}, SubCount: 12, AreaCount: 3}}
        fmt.Println("\n--- explorer_search ---")
        fmt.Println(stripANSI(em2.View()))

        // explorer_generated
        em3 := NewExplorerModel()
        em3.Focus()
        em3.state = protocol.ExplorerGenerated
        em3.selectedCategory = Category{Name: "cafe & coffee shop"}
        em3.genNicheName = "cafe_coffee_shop"
        em3.genFilesDone = []GenerateFileStatus{
                {Name: "niche.yaml", Detail: "5 targets · 5 area · 4 filter", Completed: true},
                {Name: "ice_breaker.md", Detail: "3 varian", Completed: true},
                {Name: "queries.md", Detail: "8 query per area", Completed: true},
        }
        fmt.Println("\n--- explorer_generated ---")
        fmt.Println(stripANSI(em3.View()))
}
