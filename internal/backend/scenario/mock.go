package scenario

import (
        "fmt"
        "math/rand"
        "time"

        "github.com/WaClaw-App/waclaw/pkg/protocol"
)

// MockData generates realistic-looking demo data for the WaClaw TUI.
//
// All data is Indonesian-locale-appropriate (business names, addresses, cities)
// because the demo backend targets the same audience as the production app.
// The data is deterministic within a single run to avoid visual jitter, but
// varies across runs for realism.
type MockData struct {
        rng *rand.Rand
}

// NewMockData creates a MockData with a time-seeded RNG.
func NewMockData() *MockData {
        return &MockData{rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

// Tagline returns the boot screen tagline.
func (m *MockData) Tagline() string {
        return "Vertical-borderless. Micro-interactive. File-based."
}

// BootData returns the mock data payload for the Boot screen. The keys
// match what the TUI BootModel reads in applyBootParams.
func (m *MockData) BootData(returning bool) map[string]any {
        waCount := 3
        nicheCount := 3
        workerCount := 4
        leadsCount := 847

        params := map[string]any{
                "wa_count":       waCount,
                "niche_count":    nicheCount,
                "worker_count":   workerCount,
                "leads_count":    leadsCount,
                "response_count": 0,
                "niche_already_set": returning,
        }

        if returning {
                params["workers"] = m.MarchingWorkers()
        }

        return params
}

// MarchingWorkers returns the worker rows used by the army march animation.
// The arrow_count values decrease per row to create the cascading indent
// shown in the doc (▸▸▸▸▸▸, ▸▸▸▸▸, ▸▸▸).
func (m *MockData) MarchingWorkers() []map[string]any {
        return []map[string]any{
                {"name": "web_developer", "arrow_count": 6},
                {"name": "undangan_digital", "arrow_count": 5},
                {"name": "social_media_mgr", "arrow_count": 3},
        }
}

// BootConfigErrorData returns variant data for the config_error boot state.
func (m *MockData) BootConfigErrorData() map[string]any {
        return map[string]any{
                "wa_count":         3,
                "niche_count":      3,
                "worker_count":     4,
                "leads_count":      847,
                "ok_niche_count":   2,
                "error_niche_count": 1,
                "error_niche":      "fotografer",
        }
}

// BootLicenseExpiredData returns variant data for the license_expired boot state.
func (m *MockData) BootLicenseExpiredData() map[string]any {
        return map[string]any{
                "wa_count":    3,
                "niche_count": 3,
                "worker_count": 4,
                "leads_count": 847,
        }
}

// BootDeviceConflictData returns variant data for the device_conflict boot state.
func (m *MockData) BootDeviceConflictData() map[string]any {
        return map[string]any{
                "device_name": "PC-KANTOR",
                "last_active": "12 menit",
        }
}

// BootResponseData returns variant data for the returning + new responses boot state.
func (m *MockData) BootResponseData() map[string]any {
        return map[string]any{
                "wa_count":       3,
                "niche_count":    3,
                "worker_count":   4,
                "leads_count":    847,
                "response_count": 3,
        }
}

// BootDisconnectData returns variant data for the WA disconnected boot state.
func (m *MockData) BootDisconnectData() map[string]any {
        return map[string]any{
                "niche_count":  3,
                "worker_count": 4,
                "leads_count":  847,
        }
}

// LoginData returns the mock data payload for the Login screen. The keys
// match what the TUI LoginModel reads in applyLoginParams.
func (m *MockData) LoginData() map[string]any {
        return map[string]any{
                "filled_slots":   0,
                "total_slots":   3,
                "active_slot":   0,
                "contact_count": 0,
                "qr_data":       m.QRData(),
        }
}

// LoginScannedData returns data for the qr_scanned login state.
func (m *MockData) LoginScannedData() map[string]any {
        return map[string]any{
                "filled_slots":   1,
                "total_slots":   3,
                "active_slot":   1,
                "contact_count": 847,
                "phone_numbers": []any{"0812-xxxx-3456"},
        }
}

// LoginSuccessData returns data for the login_success state.
func (m *MockData) LoginSuccessData() map[string]any {
        return map[string]any{
                "filled_slots":   1,
                "total_slots":   3,
                "active_slot":   1,
                "contact_count": 847,
                "phone_numbers": []any{"0812-xxxx-3456"},
        }
}

// LoginExpiredData returns data for the login_expired state.
func (m *MockData) LoginExpiredData() map[string]any {
        return map[string]any{
                "filled_slots":   1,
                "total_slots":   3,
                "active_slot":   1,
                "expired_slot":  1,
                "active_slots":  2,
                "last_session_ago": "3 hari",
        }
}

// LoginFailedData returns data for the login_failed state.
func (m *MockData) LoginFailedData() map[string]any {
        return map[string]any{
                "filled_slots":   0,
                "total_slots":   3,
                "active_slot":   0,
        }
}

// QRData returns a mock QR code data string for the login screen.
func (m *MockData) QRData() string {
        return fmt.Sprintf("WACLAW-QR-%d", m.rng.Intn(999999))
}

// Niches returns the list of available niches for the niche select screen.
func (m *MockData) Niches() []map[string]any {
        return []map[string]any{
                {
                        "name":     "kuliner",
                        "label":    "Kuliner & Restoran",
                        "emoji":    "🍜",
                        "targets":  []string{"restaurant", "cafe", "catering"},
                        "selected": true,
                },
                {
                        "name":     "kecantikan",
                        "label":    "Kecantikan & Salon",
                        "emoji":    "💇",
                        "targets":  []string{"salon", "spa", "barbershop"},
                        "selected": false,
                },
                {
                        "name":     "otomotif",
                        "label":    "Otomotif & Bengkel",
                        "emoji":    "🔧",
                        "targets":  []string{"bengkel", "sparepart", "carwash"},
                        "selected": true,
                },
                {
                        "name":     "kesehatan",
                        "label":    "Kesehatan & Klinik",
                        "emoji":    "🏥",
                        "targets":  []string{"klinik", "apotek", "dokter"},
                        "selected": false,
                },
                {
                        "name":     "pendidikan",
                        "label":    "Pendidikan & Kursus",
                        "emoji":    "📚",
                        "targets":  []string{"bimbel", "kursus", "seminar"},
                        "selected": false,
                },
                {
                        "name":     "properti",
                        "label":    "Properti & Interior",
                        "emoji":    "🏠",
                        "targets":  []string{"agen", "interior", "renovasi"},
                        "selected": true,
                },
                {
                        "name":     "fashion",
                        "label":    "Fashion & Boutique",
                        "emoji":    "👗",
                        "targets":  []string{"boutique", "tailor", "hijab"},
                        "selected": false,
                },
                {
                        "name":     "teknologi",
                        "label":    "Teknologi & IT",
                        "emoji":    "💻",
                        "targets":  []string{"webdev", "appdev", "it_support"},
                        "selected": false,
                },
                {
                        "name":     "photography",
                        "label":    "Photography & Video",
                        "emoji":    "📷",
                        "targets":  []string{"photographer", "videographer", "drone"},
                        "selected": false,
                },
                {
                        "name":     "event",
                        "label":    "Event & Organizer",
                        "emoji":    "🎉",
                        "targets":  []string{"wedding", "corporate", "birthday"},
                        "selected": false,
                },
        }
}

// ScrapeProgress returns mock scrape progress data.
func (m *MockData) ScrapeProgress() map[string]any {
        total := 50 + m.rng.Intn(100)
        found := total/2 + m.rng.Intn(total/2)
        validated := found/2 + m.rng.Intn(found/2)

        return map[string]any{
                "total_queries":    total,
                "leads_found":      found,
                "leads_validated":  validated,
                "high_value_count": m.rng.Intn(5),
                "progress_pct":     float64(found) / float64(total) * 100,
                "niches_active":    []string{"kuliner", "otomotif", "properti"},
        }
}

// Stats returns mock dashboard statistics.
func (m *MockData) Stats() map[string]any {
        leads := 200 + m.rng.Intn(300)
        sent := leads / 2
        responded := sent / 4
        converted := responded / 5

        return map[string]any{
                "leads_total":      leads,
                "leads_new":        30 + m.rng.Intn(50),
                "messages_sent":    sent,
                "responses":        responded,
                "conversions":      converted,
                "conversion_rate":  float64(converted) / float64(responded) * 100,
                "active_workers":   3,
                "health_score":     85 + m.rng.Intn(15),
                "streak_days":      7 + m.rng.Intn(14),
                "best_time_str":    "selasa jam 10",
                "active_slot_count": 2,
        }
}

// Leads returns a slice of mock leads.
func (m *MockData) Leads(count int) []map[string]any {
        businessNames := []string{
                "Warung Padang Sederhana", "Bengkel Jaya Motor", "Salon Cantik Alami",
                "Klinik Sehat Sentosa", "Toko Bangun Makmur", "Kedai Kopi Nusantara",
                "Apotek Farma Utama", "Boutique Elegance", "CV. Teknologi Maju",
                "Restoran Sedap Malam", "Fotografer Abadi", "Event Organizer Prima",
                "Interior Design Kita", "Bimbel Cerdas Bangsa", "Car Wash Bersih Kilat",
        }
        addresses := []string{
                "Jl. Sudirman No. 45, Jakarta Selatan",
                "Jl. Gatot Subroto No. 12, Bandung",
                "Jl. Ahmad Yani No. 78, Surabaya",
                "Jl. Diponegoro No. 33, Yogyakarta",
                "Jl. Imam Bonjol No. 56, Semarang",
                "Jl. Pemuda No. 21, Medan",
                "Jl. Veteran No. 9, Makassar",
                "Jl. Merdeka No. 15, Denpasar",
        }
        cities := []string{
                "Jakarta", "Bandung", "Surabaya", "Yogyakarta",
                "Semarang", "Medan", "Makassar", "Denpasar",
        }
        phases := []protocol.LeadPhase{
                protocol.LeadBaru, protocol.LeadIceBreakerSent, protocol.LeadResponded,
                protocol.LeadOfferSent, protocol.LeadConverted, protocol.LeadCold,
                protocol.LeadNoResponse, protocol.LeadFollowUp1,
        }

        leads := make([]map[string]any, 0, count)
        for i := 0; i < count; i++ {
                phase := phases[m.rng.Intn(len(phases))]
                rating := 3.0 + m.rng.Float64()*2.0
                reviews := m.rng.Intn(500)

                lead := map[string]any{
                        "id":         fmt.Sprintf("lead-%04d", i+1),
                        "name":       businessNames[m.rng.Intn(len(businessNames))],
                        "address":    addresses[m.rng.Intn(len(addresses))],
                        "city":       cities[m.rng.Intn(len(cities))],
                        "phone":      fmt.Sprintf("+62812%08d", m.rng.Intn(99999999)),
                        "phase":      string(phase),
                        "rating":     rating,
                        "reviews":    reviews,
                        "has_wa":     m.rng.Float64() > 0.2,
                        "has_insta":  m.rng.Float64() > 0.4,
                        "has_web":    m.rng.Float64() > 0.5,
                        "score":      m.rng.Intn(100),
                        "niche":      []string{"kuliner", "otomotif", "properti"}[m.rng.Intn(3)],
                }
                leads = append(leads, lead)
        }

        return leads
}

// Workers returns mock worker status data.
func (m *MockData) Workers() []map[string]any {
        phases := []protocol.WorkerPhase{
                protocol.WorkerScraping, protocol.WorkerSending, protocol.WorkerIdle,
        }
        niches := []string{"kuliner", "otomotif", "properti"}

        workers := make([]map[string]any, 0, len(niches))
        for i, niche := range niches {
                phase := phases[i%len(phases)]
                workers = append(workers, map[string]any{
                        "id":       fmt.Sprintf("worker-%d", i+1),
                        "niche":    niche,
                        "phase":    string(phase),
                        "leads":    20 + m.rng.Intn(80),
                        "sent":     10 + m.rng.Intn(40),
                        "responses": m.rng.Intn(15),
                        "slot":     fmt.Sprintf("slot_%d", i+1),
                })
        }
        return workers
}

// Responses returns mock response data for the response screen.
func (m *MockData) Responses() []map[string]any {
        classifications := []string{"positive", "curious", "negative", "maybe", "auto_reply"}
        businessNames := []string{
                "Warung Padang Sederhana", "Bengkel Jaya Motor", "Salon Cantik Alami",
        }

        responses := make([]map[string]any, 0, 5)
        for i := 0; i < 5; i++ {
                responses = append(responses, map[string]any{
                        "lead_id":        fmt.Sprintf("lead-%04d", i+1),
                        "business_name":  businessNames[m.rng.Intn(len(businessNames))],
                        "classification": classifications[m.rng.Intn(len(classifications))],
                        "message":        "Terima kasih, saya tertarik dengan penawarannya",
                        "received_at":    time.Now().Add(-time.Duration(i*15) * time.Minute).Format(time.RFC3339),
                })
        }
        return responses
}

// Template returns a mock template for the template manager.
func (m *MockData) Template() map[string]any {
        return map[string]any{
                "type":     "ice_breaker",
                "variant":  1,
                "content":  "Halo {{.Title}}, kami punya penawaran menarik untuk {{.Category}} di {{.City}}!",
                "placeholders": []string{
                        "{{.Title}}", "{{.Category}}", "{{.Address}}",
                        "{{.City}}", "{{.Rating}}", "{{.Reviews}}", "{{.Area}}",
                },
                "preview": "Halo Warung Padang Sederhana, kami punya penawaran menarik untuk Kuliner & Restoran di Jakarta!",
                "preview_values": map[string]any{
                        "Title":    "kopi nusantara",
                        "Category": "cafe",
                        "Address":  "jl. hasanuddin 23, kediri",
                        "City":     "kediri",
                        "Rating":   "4.2",
                        "Reviews":  "87",
                        "Area":     "kediri",
                },
        }
}

// Notifications returns a sequence of mock notifications for the demo.
func (m *MockData) Notifications() []map[string]any {
        return []map[string]any{
                {
                        "type":     string(protocol.NotifScrapeComplete),
                        "severity": string(protocol.SeverityPositive),
                        "data": map[string]any{
                                "message":   "Scraping selesai — 47 lead ditemukan",
                                "niche":     "kuliner",
                                "leads":     47,
                                "high_value": 3,
                        },
                },
                {
                        "type":     string(protocol.NotifResponseReceived),
                        "severity": string(protocol.SeverityPositive),
                        "data": map[string]any{
                                "message": "Response masuk dari Warung Padang Sederhana",
                                "lead_id": "lead-0001",
                        },
                },
                {
                        "type":     string(protocol.NotifBatchSendComplete),
                        "severity": string(protocol.SeverityNeutral),
                        "data": map[string]any{
                                "message": "Batch selesai — 25 pesan terkirim",
                                "count":   25,
                        },
                },
                {
                        "type":     string(protocol.NotifStreakMilestone),
                        "severity": string(protocol.SeverityInformative),
                        "data": map[string]any{
                                "message": "7 hari streak! Terus gas!",
                                "days":    7,
                        },
                },
                {
                        "type":     string(protocol.NotifWAFlag),
                        "severity": string(protocol.SeverityCritical),
                        "data": map[string]any{
                                "message": "Slot 2 ditandai WhatsApp — auto-pause",
                                "slot_id": "slot_2",
                        },
                },
        }
}

// ShieldData returns mock anti-ban shield data.
func (m *MockData) ShieldData() map[string]any {
        healthScore := 85 + m.rng.Intn(15)
        return map[string]any{
                "health_score":    healthScore,
                "health_level":    shieldLevel(healthScore),
                "warning_health_score": 71,
                "danger_health_score":  38,
                "flagged_slot": map[string]any{
                        "number": "0812-xxxx-3456",
                },
                "slot_detail_7day_sent":       84,
                "slot_detail_7day_responded":  12,
                "slot_detail_7day_failed":     "1",
                "slot_detail_7day_warnings":   1,
                "selected_slot_health":        87,
                "slots": []map[string]any{
                        {
                                "id":           "slot_1",
                                "phone":        "+62812-3456-7890",
                                "health_score": 92,
                                "status":       "active",
                                "sent_today":   45,
                                "daily_limit":  100,
                        },
                        {
                                "id":           "slot_2",
                                "phone":        "+62812-9876-5432",
                                "health_score": 67,
                                "status":       "warning",
                                "sent_today":   78,
                                "daily_limit":  100,
                        },
                        {
                                "id":           "slot_3",
                                "phone":        "+62812-5555-1234",
                                "health_score": 95,
                                "status":       "active",
                                "sent_today":   32,
                                "daily_limit":  100,
                        },
                },
                "spam_guard": map[string]any{
                        "per_lead_limit":     3,
                        "daily_gap_hours":    24,
                        "variant_required":   true,
                        "dnc_count":          12,
                },
        }
}

// ConfigData returns mock settings data.
func (m *MockData) ConfigData() map[string]any {
        return map[string]any{
                "anti_ban": map[string]any{
                        "per_slot_hourly": 15,
                        "per_slot_daily":  100,
                        "cooldown_min":    30,
                },
                "spam_guard": map[string]any{
                        "max_per_lead":  3,
                        "min_gap_hours": 24,
                },
                "schedule": map[string]any{
                        "work_hours_start": "09:00",
                        "work_hours_end":   "17:00",
                        "timezone":         "Asia/Jakarta",
                },
                "work_hours": "09:00-17:00 wib",
                "locale":     "id",
        }
}

// ValidationData returns mock validation results.
func (m *MockData) ValidationData() map[string]any {
        return map[string]any{
                "errors":   []string{},
                "warnings": []string{
                        "niche kuliner: closing_triggers kosong — pakai default",
                        "slot_2: health score di bawah 70",
                },
                "valid": true,
        }
}

// HistoryData returns mock history data for the history screen.
func (m *MockData) HistoryData() map[string]any {
        days := make([]map[string]any, 7)

        for i := 6; i >= 0; i-- {
                days[6-i] = map[string]any{
                        "sent":        20 + m.rng.Intn(30),
                        "responses":   3 + m.rng.Intn(10),
                        "conversions": m.rng.Intn(3),
                        "events":      m.rng.Intn(5),
                }
        }

        return map[string]any{
                "days":    days,
                "insight": "Selasa jam 10 = waktu terbaik untuk kirim",
                "streak":  7 + m.rng.Intn(14),
                // Backend-computed values — the TUI must NOT compute these locally.
                "best_day_index":     1, // Tuesday
                "best_day_label":     "Selasa",
                "best_day_conv_rate": "8.3%",
                "avg_response_time":  "2 jam 14 menit",
        }
}

// FollowUpData returns mock follow-up data with structured niche groups.
func (m *MockData) FollowUpData() map[string]any {
        return map[string]any{
                "total_today":            14,
                "cold_total":             4,
                "ice_breaker_unanswered": 8,
                "sending_rate":           "9/18",
                "variant_names":          []any{"follow_up_1.md", "follow_up_2.md", "follow_up_3.md"},
                "variant_manual_only":    []any{float64(2)}, // Index 2 (3rd variant) is manual-only
                "max_sending_visible":    float64(3),
                "max_cold_visible":       float64(6),
                "niches": []any{
                        map[string]any{
                                "niche_name": "web_developer",
                                "fu1_count":  8,
                                "fu2_count":  3,
                                "cold_count": 2,
                                "leads": []any{
                                        map[string]any{
                                                "business_name":   "Kedai Kopi Nusantara",
                                                "phase":           string(protocol.FUPhase1),
                                                "previous_action": "ice breaker: 2 hari lalu",
                                                "next_action":     "follow-up hari ini",
                                                "slot_number":     1,
                                                "variant_name":    "follow_up_1.md",
                                                "is_sending":      false,
                                                "wait_time":       "14m 02s",
                                        },
                                        map[string]any{
                                                "business_name":   "Bengkel Jaya Motor",
                                                "phase":           string(protocol.FUPhase1),
                                                "previous_action": "ice breaker: 3 hari lalu",
                                                "next_action":     "follow-up hari ini",
                                                "slot_number":     2,
                                                "variant_name":    "follow_up_1.md",
                                                "is_sending":      true,
                                                "wait_time":       "02m 30s",
                                        },
                                },
                        },
                        map[string]any{
                                "niche_name": "undangan_digital",
                                "fu1_count":  3,
                                "fu2_count":  0,
                                "cold_count": 2,
                                "leads":      []any{},
                        },
                },
                "cold_leads": []any{
                        map[string]any{
                                "business_name":      "Salon Cantik Alami",
                                "phase":              string(protocol.FUPhaseCold),
                                "ice_breaker_action": "5 hari lalu",
                                "follow_up_1_action": "3 hari lalu",
                                "follow_up_2_action": "1 hari lalu",
                        },
                        map[string]any{
                                "business_name":      "Apotek Farma Utama",
                                "phase":              string(protocol.FUPhaseCold),
                                "ice_breaker_action": "7 hari lalu",
                                "follow_up_1_action": "4 hari lalu",
                                "follow_up_2_action": "2 hari lalu",
                        },
                },
                "recontact_leads": []any{
                        map[string]any{
                                "business_name":      "Warung Padang Sederhana",
                                "previous_response":  "Terima kasih, saya tertarik",
                                "days_since_response": 8,
                                "days_since_offer":   7,
                                "can_recontact":      true,
                        },
                },
        }
}

// ComposeData returns mock data for the compose screen.
// The backend provides snippets and config — the TUI never hardcodes them.
func (m *MockData) ComposeData() map[string]any {
        return map[string]any{
                "target":    "Kedai Kopi Nusantara",
                "max_chars": float64(500),
                "snippets": []any{
                        map[string]any{"text": "boleh lihat dulu aja kak", "category": "soft_pitch"},
                        map[string]any{"text": "aku kasih preview gratis ya", "category": "free_offer"},
                        map[string]any{"text": "bisa telepon bentar buat diskusi?", "category": "move_to_call"},
                        map[string]any{"text": "harganya mulai 500rb kak", "category": "direct_price"},
                        map[string]any{"text": "mau aku kirim contohnya?", "category": "send_sample"},
                },
        }
}

// UpdateData returns mock update/upgrade data for the Update screen.
// Keys match what the TUI UpdateModel reads in HandleNavigate/HandleUpdate.
func (m *MockData) UpdateData(isMajor bool) map[string]any {
        currentVersion := "v1.3.2"
        newVersion := "v1.3.3"
        changelog := []any{
                "fix: scrape duplicate in certain areas",
                "fix: WA rotator cooldown calculation more accurate",
                "perf: database query 15% faster",
        }
        if isMajor {
                newVersion = "v2.0.0"
                changelog = []any{
                        "new architecture: multi-device support",
                        "AI-powered message personalization",
                        "dashboard real-time collaboration",
                        "20+ new features",
                }
        }
        return map[string]any{
                "current_version": currentVersion,
                "new_version":     newVersion,
                "is_major":        isMajor,
                "changelog":       changelog,
                "license_key":     "WACL-XXXX-XXXX-XXXX-XXXX",
                "license_expiry":  "30 juni 2025",
                "license_status":  "valid",
                "license_prefix":  "WACL2-",
        }
}

// FollowUpDetailData returns structured follow-up data matching what the TUI
// parseNicheGroups/parseColdLeads/parseRecontactLeads functions expect.
// The backend is the authoritative source — the TUI never hardcodes niches
// or variant names.
func (m *MockData) FollowUpDetailData() map[string]any {
        return map[string]any{
                "total_today":             14,
                "cold_total":              4,
                "ice_breaker_unanswered":  8,
                "variant_names":           []string{"follow_up_1.md", "follow_up_2.md", "follow_up_3.md"},
                "variant_previews":        []string{
                        "halo kak, cuma ngingetin aja",
                        "kak, penawaran terbatas nih",
                        "terakhir kak, kalo berkenan",
                },
                "variant_manual_only_idx": 2,
                "sending_rate":            "9/18",
                "niches": []any{
                        map[string]any{
                                "niche_name": "web_developer",
                                "fu1_count":  8,
                                "fu2_count":  3,
                                "cold_count": 2,
                                "leads": []any{
                                        map[string]any{
                                                "business_name":  "kopi nusantara",
                                                "phase":          string(protocol.FUPhase1),
                                                "previous_action": "ice breaker: 2 hari lalu",
                                                "next_action":    "follow-up hari ini",
                                                "slot_number":    2,
                                                "variant_name":   "variant_2",
                                                "is_sending":     true,
                                                "wait_time":      "0m 00s",
                                        },
                                        map[string]any{
                                                "business_name":  "gym fortress pro",
                                                "phase":          string(protocol.FUPhase1),
                                                "previous_action": "ice breaker: 2 hari lalu",
                                                "next_action":    "follow-up hari ini",
                                                "slot_number":    1,
                                                "variant_name":   "variant_1",
                                                "is_sending":     false,
                                                "wait_time":      "14m 02s",
                                        },
                                },
                        },
                        map[string]any{
                                "niche_name": "undangan_digital",
                                "fu1_count":  4,
                                "fu2_count":  0,
                                "cold_count": 1,
                                "leads":      []any{},
                        },
                        map[string]any{
                                "niche_name": "social_media_mgr",
                                "fu1_count":  2,
                                "fu2_count":  1,
                                "cold_count": 1,
                                "leads":      []any{},
                        },
                },
                "cold_leads": []any{
                        map[string]any{
                                "business_name":     "apotek sehat",
                                "phase":             string(protocol.FUPhaseCold),
                                "ice_breaker_action": "ice breaker 5 hari lalu",
                                "follow_up_1_action": "3 hari lalu",
                                "follow_up_2_action": "1 hari lalu",
                        },
                        map[string]any{
                                "business_name":     "toko elektronik",
                                "phase":             string(protocol.FUPhaseCold),
                                "ice_breaker_action": "ice breaker 6 hari lalu",
                                "follow_up_1_action": "4 hari lalu",
                                "follow_up_2_action": "2 hari lalu",
                        },
                },
                "recontact_leads": []any{
                        map[string]any{
                                "business_name":    "salon cantik",
                                "previous_response": "berapa harganya?",
                                "days_since_response": 8,
                                "days_since_offer":   7,
                                "can_recontact":      true,
                        },
                        map[string]any{
                                "business_name":    "wedding bliss WO",
                                "previous_response": "boleh kirim contohnya?",
                                "days_since_response": 10,
                                "days_since_offer":   9,
                                "can_recontact":      true,
                        },
                },
        }
}

// HistoryDetailData returns structured history data matching what the TUI
// parseHistoryEvents/parseWeekData functions expect.
// The backend provides all computed values (conversion rate, best day) —
// the TUI never computes rates locally.
func (m *MockData) HistoryDetailData() map[string]any {
        dayLabels := []string{"Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu", "Minggu"}
        messages := []any{12, 18, 14, 10, 8, 4, 2}
        responses := []any{2, 5, 3, 2, 1, 0, 0}
        converts := []any{0, 2, 1, 0, 0, 0, 0}

        events := make([]any, 0, 10)
        businessNames := []string{
                "kopi nusantara", "wedding bliss", "gym fortress",
                "salon cantik", "venue gardenia", "toko makmur",
                "bengkel jaya", "apotek sehat",
        }
        icons := []string{"💬", "📤", "✅", "★", "🔍", "⚡"}
        eventTypes := []string{"respond", "terkirim", "sampai", "jackpot!", "batch scrape selesai", "auto-pilot on"}
        for i := 0; i < 8; i++ {
                hour := 9 + i
                evt := map[string]any{
                        "time":          fmt.Sprintf("%02d:%02d", hour, 23-i*7),
                        "icon":          icons[i%len(icons)],
                        "business_name": businessNames[i%len(businessNames)],
                        "event_type":    eventTypes[i%len(eventTypes)],
                }
                if i == 3 {
                        evt["highlight"] = true
                        evt["detail"] = "skor: 9.2"
                }
                if i == 0 {
                        evt["detail"] = "\"iya kak, boleh\""
                }
                events = append(events, evt)
        }

        return map[string]any{
                "today_stats": map[string]any{
                        "sent":      7,
                        "respond":   4,
                        "convert":   0,
                        "new_leads": 67,
                        "scrapes":   1,
                },
                "events": events,
                "week": map[string]any{
                        "day_labels":        dayLabels,
                        "messages":          messages,
                        "responses":         responses,
                        "converts":          converts,
                        "best_day_index":    1,
                        "best_day_label":    "Selasa",
                        "best_day_conv_rate": "16.7%",
                },
                "week_stats": map[string]any{
                        "total_sent":       68,
                        "total_respond":    13,
                        "total_convert":    3,
                        "total_new_leads":  247,
                        "total_scrapes":    3,
                },
                "avg_response_time": "3.2 jam",
        }
}

// UpdateDownloadProgress returns incremental download progress data.
// Simulates realistic download progress for the demo timeline.
func (m *MockData) UpdateDownloadProgress(percent float64) map[string]any {
        totalMB := 5.1
        downloadedMB := totalMB * percent
        return map[string]any{
                "percent":     percent,
                "size":        fmt.Sprintf("%.1fMB", totalMB),
                "downloaded":  fmt.Sprintf("%.1fMB", downloadedMB),
                "speed":       "350KB/s",
                "eta":         fmt.Sprintf("%.0f detik", (1-percent)*15),
                "source":      "https://releases.waclaw.dev/v1.3.3/",
        }
}

// UpdateReadyData returns data for the update_ready state.
func (m *MockData) UpdateReadyData() map[string]any {
        return map[string]any{
                "checksum_verified": true,
                "backup_path":      "~/.waclaw/waclaw-v1.3.2.bak",
        }
}

// UpdateExpiredData returns data for the license_expired_with_upgrade state.
func (m *MockData) UpdateExpiredData() map[string]any {
        return map[string]any{
                "expired_date": "15 april 2025",
        }
}

// ConfigSettingsData returns mock config data for the Settings screen.
// All values that were previously hardcoded in the TUI are now
// provided by the backend through HandleNavigate params.
func (m *MockData) ConfigSettingsData() map[string]any {
        return map[string]any{
                "active_niches":  "web_developer, undangan_digital, social_media_mgr",
                "wa_slots":       "3 nomor aktif (rotator)",
                "worker_pool":    "3 worker (1 per niche)",
                "area":           "multi (8 kota, 14 kecamatan)",
                "work_hours":     "09:00-17:00 wib",
                "rate_limit":     "6/jam per slot, 50/hari total",
                "rotator_mode":   "round-robin + cooldown",
                "autopilot":      "aktif",
        }
}

// ShieldConfigData returns anti-ban config data for the Shield settings view.
// All values that were previously hardcoded in the TUI are now
// provided by the backend through HandleNavigate params.
func (m *MockData) ShieldConfigData() map[string]any {
        return map[string]any{
                "per_slot_hourly":      6,
                "per_slot_daily":       50,
                "cooldown_min":         47,
                "min_delay_min":        8,
                "max_delay_min":        25,
                "delay_variance_pct":   30,
                "work_hours":           "09:00-17:00 wib",
                "auto_pause":           "auto",
                "health_threshold":      50,
                "rotator_mode":         "round-robin + cooldown",
                "template_rotation":    "aktif",
                "rotation_mode":        "round-robin",
                "emoji_variation":      "aktif",
                "paragraph_shuffle":    "aktif",
                "per_lead_lifetime":    3,
                "msg_interval_hours":   24,
                "followup_delay_days":  2,
                "followup_variant":     "aktif",
                "cold_after":           2,
                "recontact_delay_days": 7,
                "auto_block":           "aktif",
                "dup_cross_niche":      "aktif",
                "wa_pre_validation":    "aktif",
                "wa_validation_method":  "check-registration",
                "daily_budget_sent":    27,
                "daily_budget_total":   50,

                // Additional i18n format string params — backend is authoritative source
                "dnc_count":            12,
                "ice_breaker_variants": 3,
                "offer_variants":       3,
                "work_hours_duration":  8,
                "health_recovery_pts":  5,
                "timezone_short":       "wib",
                "timezone_full":        "asia/jakarta",
                "current_time":         "14:23",
                "failed_count":         1,

                // Ban risk assessment — backend is authoritative source
                "ban_risk_level":       "low",
                "ban_risk_emoji":       "🟢",
                "risk_indicators": map[string]any{
                        "even_distribution":  true,
                        "cooldown_ok":        true,
                        "template_varied":    true,
                        "no_overload":        true,
                        "work_hours_ok":      true,
                        "spam_guard_active":  true,
                        "dnc_respected":      true,
                },
        }
}

// ExplorerCategories returns the list of niche explorer categories per doc/11.
func (m *MockData) ExplorerCategories() []map[string]any {
        return []map[string]any{
                {"emoji": "🍜", "name": "kuliner", "sub_categories": []string{"cafe", "restoran", "catering", "bakery"}, "sub_count": 4},
                {"emoji": "💇", "name": "kecantikan", "sub_categories": []string{"salon", "barbershop", "spa", "nail art"}, "sub_count": 4},
                {"emoji": "💪", "name": "fitness", "sub_categories": []string{"gym", "yoga studio", "personal trainer"}, "sub_count": 3},
                {"emoji": "🏠", "name": "properti", "sub_categories": []string{"agent", "developer", "interior", "renovation"}, "sub_count": 4},
                {"emoji": "🎓", "name": "pendidikan", "sub_categories": []string{"bimbel", "kursus", "training", "workshop"}, "sub_count": 4},
                {"emoji": "🏥", "name": "kesehatan", "sub_categories": []string{"klinik", "apotek", "dokter", "therapy"}, "sub_count": 4},
                {"emoji": "📸", "name": "fotografi", "sub_categories": []string{"wedding", "product", "studio", "event"}, "sub_count": 4},
                {"emoji": "💻", "name": "teknologi", "sub_categories": []string{"web dev", "app dev", "it support", "digital"}, "sub_count": 4},
                {"emoji": "🚗", "name": "otomotif", "sub_categories": []string{"bengkel", "dealer", "rental", "sparepart"}, "sub_count": 4},
                {"emoji": "🎉", "name": "event", "sub_categories": []string{"wedding org", "decorator", "catering", "mc"}, "sub_count": 4},
        }
}

// NicheSelectItems returns the niche select items per doc/02.
// Includes emoji and selected fields for DRY consistency with Niches().
func (m *MockData) NicheSelectItems() []map[string]any {
        return []map[string]any{
                {"name": "web developer", "description": "buat yang jual jasa bikin web", "area": "kediri, 15km", "templates": 3, "emoji": "💻", "selected": false, "targets": []string{"webdev", "appdev", "it_support"}},
                {"name": "undangan digital", "description": "buat yang jual undangan digital", "area": "kediri + surabaya", "templates": 2, "emoji": "🎉", "selected": false, "targets": []string{"wedding", "corporate", "birthday"}},
                {"name": "social media mgr", "description": "buat yang jasa kelola sosmed", "area": "", "templates": 0, "emoji": "📱", "selected": false, "targets": []string{"instagram", "tiktok", "youtube"}},
                {"name": "fotografer", "description": "buat yang jasa foto & portfolio", "area": "", "templates": 0, "emoji": "📸", "selected": false, "targets": []string{"wedding", "product", "event"}},
                {"name": "custom", "description": "bikin niche sendiri dari file", "area": "", "templates": 0, "emoji": "📁", "selected": false, "targets": []string{}},
        }
}

// shieldLevel returns the health level string for a given score.
func shieldLevel(score int) string {
        switch {
        case score >= 90:
                return "healthy"
        case score >= 50:
                return "warning"
        default:
                return "danger"
        }
}

// ExplorerSources returns the data sources for the explorer category detail per doc/11.
func (m *MockData) ExplorerSources() []map[string]any {
        return []map[string]any{
                {"name": "WhatsApp Business Directory", "count": 247, "unit": "bisnis terdaftar"},
                {"name": "Google Maps Categories", "count": 189, "unit": "listing aktif"},
        }
}

// ExplorerAreas returns the auto-detect areas for the explorer category detail per doc/11.
func (m *MockData) ExplorerAreas() []map[string]any {
        return []map[string]any{
                {"city": "kediri", "radius": "15km"},
                {"city": "nganjuk", "radius": "10km"},
                {"city": "tulungagung", "radius": "10km"},
                {"city": "blitar", "radius": "10km"},
                {"city": "madiun", "radius": "10km"},
        }
}

// ExplorerFilters returns the default filters for the explorer category detail per doc/11.
func (m *MockData) ExplorerFilters() []map[string]any {
        return []map[string]any{
                {"symbol": "✗", "label": "punya website", "detail": "skip yang udah punya", "inverted": true},
                {"symbol": "✓", "label": "punya instagram", "detail": "lebih potensial"},
                {"symbol": "⭐", "label": "3.0 - 4.8", "detail": "rating range", "neutral": true},
                {"symbol": "📊", "label": "5 - 500 reviews", "detail": "ukuran bisnis", "neutral": true},
        }
}

// ExplorerTemplates returns the template list for the explorer category detail per doc/11.
func (m *MockData) ExplorerTemplates() []map[string]any {
        return []map[string]any{
                {"name": "niche.yaml", "detail": "filter + target + area", "completed": true},
                {"name": "ice_breaker.md", "detail": "3 varian", "completed": true},
                {"name": "queries.md", "detail": "search queries per area", "completed": true},
        }
}

// ExplorerGenFiles returns the file list for the generating state per doc/11.
// This provides the TUI with the list of files that will be generated so it
// can render the generating progress UI immediately, before individual files
// are completed.
func (m *MockData) ExplorerGenFiles() []map[string]any {
        return []map[string]any{
                {"name": "membuat folder ~/.waclaw/niches/kuliner/", "detail": "", "completed": false},
                {"name": "menulis niche.yaml", "detail": "", "completed": false},
                {"name": "menulis ice_breaker.md (3 varian)", "detail": "", "completed": false},
                {"name": "menulis queries.md", "detail": "", "completed": false},
        }
}

// MonitorData returns comprehensive dashboard data for the Monitor screen.
// All values that were previously hardcoded in the TUI are now provided
// by the backend through HandleNavigate params.
func (m *MockData) MonitorData() map[string]any {
        stats := m.Stats()
        leads := stats["leads_total"].(int)
        sent := stats["messages_sent"].(int)
        responded := stats["responses"].(int)
        converted := stats["conversions"].(int)

        return map[string]any{
                "app_name":    "waclaw",
                "niche_count": 3,
                "wa_num_count": 3,
                "error_slot":  "slot-1",
                "work_hours":  "09:00-17:00 wib",
                "current_time": time.Now().Format("15:04"),
                "conv_rate":   fmt.Sprintf("%.1f%%", float64(converted)/float64(responded)*100),
                "best_day":    "selasa",
                "today_stats": []int64{int64(leads / 2), int64(sent / 2), int64(responded / 2), int64(converted / 2)},
                "week_stats":  []int64{int64(leads), int64(sent), int64(responded), int64(converted)},
                "wa_slots": []map[string]any{
                        {"label": "slot-1", "number": "0812-xxxx-3456", "active": true, "hours": "4/6 jam"},
                        {"label": "slot-2", "number": "0812-9876-5432", "active": true, "hours": "2/6 jam"},
                        {"label": "slot-3", "number": "0812-5555-1234", "active": false, "hours": "ready: 14m"},
                },
                "workers": []map[string]any{
                        {"name": "web_developer", "active": true, "phase": "scraping", "queued": 24, "sent": 28, "responded": 3, "convert_count": 2, "duration": "5j 37m", "send_dur": "11m 23s"},
                        {"name": "undangan_digital", "active": true, "phase": "sending", "queued": 18, "sent": 15, "responded": 5, "convert_count": 1, "duration": "3j 12m", "send_dur": "8m 45s"},
                        {"name": "social_media_mgr", "active": false, "phase": "idle", "queued": 7, "sent": 10, "responded": 1, "convert_count": 0, "duration": "12j", "send_dur": ""},
                },
                "activities": []map[string]any{
                        {"time": time.Now().Add(-5 * time.Minute).Format(time.RFC3339), "niche": "[web_dev]", "business": "kopi nusantara", "status": "respond", "detail": "iya kak, boleh lihat"},
                        {"time": time.Now().Add(-12 * time.Minute).Format(time.RFC3339), "niche": "[undangan]", "business": "wedding bliss", "status": "sent", "detail": "──"},
                        {"time": time.Now().Add(-25 * time.Minute).Format(time.RFC3339), "niche": "[social]", "business": "gym fortress", "status": "converted", "detail": "skor: 9.2"},
                },
                "pending": []map[string]any{
                        {"niche": "[web_dev]", "business": "kopi nusantara", "snippet": "iya kak, boleh lihat"},
                        {"niche": "[undangan]", "business": "wedding bliss", "snippet": "berapa harganya?"},
                },
        }
}

// ResponseScreenData returns comprehensive response data for the Response screen.
// Includes offer_text (rendered by backend) and trigger (closing trigger match)
// so the TUI never has to guess or hardcode these values.
func (m *MockData) ResponseScreenData() map[string]any {
        return map[string]any{
                "business":  "kopi nusantara",
                "category":  "cafe",
                "area":      "kediri",
                "message":   "iya kak, boleh lihat desainnya?",
                "offer_text": "Halo Kopi Nusantara, kami punya penawaran menarik untuk Cafe di Kediri! Dapatkan website profesional mulai 500rb.",
                "trigger":   "boleh lihat",
                "niche":     "[web_dev]",
                "class":     string(protocol.ClassPositif),
                "conversion": map[string]any{
                        "business":     "kopi nusantara",
                        "pipeline":     "ice breaker → offer → deal",
                        "time_taken":   "2 hari 4 jam",
                        "trophy_count": 3,
                        "revenue":      "rp 7.5jt",
                },
                "queue": []map[string]any{
                        {"index": 1, "business": "kopi nusantara", "category": "cafe", "area": "kediri", "message": "iya kak, boleh lihat", "class": string(protocol.ClassPositif)},
                        {"index": 2, "business": "bengkel jaya", "category": "otomotif", "area": "surabaya", "message": "berapa harganya?", "class": string(protocol.ClassCurious)},
                        {"index": 3, "business": "salon cantik", "category": "kecantikan", "area": "bandung", "message": "terima kasih sudah menghubungi", "class": string(protocol.ClassAutoReply)},
                },
        }
}

// ExplorerGenFilesDone returns the completed file list for the ExplorerGenerated
// state. Centralized here instead of being hardcoded in engine.go so that mock
// data lives in one place and can be updated consistently.
func (m *MockData) ExplorerGenFilesDone() []map[string]any {
        return []map[string]any{
                {"name": "niche.yaml", "detail": "5 targets · 5 area · 4 filter", "completed": true},
                {"name": "ice_breaker.md", "detail": "3 varian", "completed": true},
                {"name": "queries.md", "detail": "8 query per area", "completed": true},
        }
}
