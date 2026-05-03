package monitor

import (
	"testing"
	"time"

	"github.com/WaClaw-App/waclaw/internal/tui/i18n"
	"github.com/WaClaw-App/waclaw/internal/tui/testutil"
	"github.com/WaClaw-App/waclaw/pkg/protocol"
)

// ---------------------------------------------------------------------------
// Dashboard screen tests
// ---------------------------------------------------------------------------

func TestDashboardNew(t *testing.T) {
	d := NewDashboard()
	if d.ID() != protocol.ScreenMonitor {
		t.Errorf("expected ScreenMonitor, got %s", d.ID())
	}
	if d.state != protocol.MonitorLiveDashboard {
		t.Errorf("expected MonitorLiveDashboard, got %s", d.state)
	}
}

func TestDashboardDefaultView(t *testing.T) {
	d := NewDashboard()
	helper := testutil.NewScreenHelper(d)
	view := helper.View()
	if view == "" {
		t.Error("expected non-empty view from dashboard")
	}
}

func TestDashboardHandleNavigate(t *testing.T) {
	d := NewDashboard()

	// Navigate to night mode
	err := d.HandleNavigate(map[string]any{"state": "monitor_night"})
	if err != nil {
		t.Fatalf("HandleNavigate failed: %v", err)
	}
	if d.state != protocol.MonitorNight {
		t.Errorf("expected MonitorNight, got %s", d.state)
	}

	// Navigate to error state
	err = d.HandleNavigate(map[string]any{"state": "monitor_error"})
	if err != nil {
		t.Fatalf("HandleNavigate failed: %v", err)
	}
	if d.state != protocol.MonitorError {
		t.Errorf("expected MonitorError, got %s", d.state)
	}

	// Navigate to empty state
	err = d.HandleNavigate(map[string]any{"state": "monitor_empty"})
	if err != nil {
		t.Fatalf("HandleNavigate failed: %v", err)
	}
	if d.state != protocol.MonitorEmpty {
		t.Errorf("expected MonitorEmpty, got %s", d.state)
	}

	// Navigate to pending responses
	err = d.HandleNavigate(map[string]any{"state": "monitor_pending_responses"})
	if err != nil {
		t.Fatalf("HandleNavigate failed: %v", err)
	}
	if d.state != protocol.MonitorPendingResponses {
		t.Errorf("expected MonitorPendingResponses, got %s", d.state)
	}
}

func TestDashboardViewNight(t *testing.T) {
	d := NewDashboard()
	d.HandleNavigate(map[string]any{"state": "monitor_night"})
	view := d.View()
	// Night mode view should contain night mode text
	if view == "" {
		t.Error("expected non-empty night view")
	}
}

func TestDashboardViewEmpty(t *testing.T) {
	d := NewDashboard()
	d.HandleNavigate(map[string]any{"state": "monitor_empty"})
	view := d.View()
	if view == "" {
		t.Error("expected non-empty empty view")
	}
}

func TestDashboardViewError(t *testing.T) {
	d := NewDashboard()
	d.HandleNavigate(map[string]any{
		"state":       "monitor_error",
		"error_slot":  "slot-1",
	})
	view := d.View()
	if view == "" {
		t.Error("expected non-empty error view")
	}
}

func TestDashboardViewPendingResponses(t *testing.T) {
	d := NewDashboard()
	d.HandleNavigate(map[string]any{
		"state": "monitor_pending_responses",
		"pending": []any{
			map[string]any{
				"niche":    "[web_dev]",
				"business": "kopi nusantara",
				"snippet":  "iya kak, boleh lihat",
			},
		},
	})
	view := d.View()
	if view == "" {
		t.Error("expected non-empty pending view")
	}
}

func TestDashboardPopulateData(t *testing.T) {
	d := NewDashboard()
	d.HandleNavigate(map[string]any{
		"state":        "monitor_live_dashboard",
		"niche_count":  float64(3),
		"wa_num_count": float64(2),
		"conv_rate":    "5.1%",
		"best_day":     "selasa",
	})
	if d.data.NicheCount != 3 {
		t.Errorf("expected 3 niches, got %d", d.data.NicheCount)
	}
	if d.data.WANumCount != 2 {
		t.Errorf("expected 2 WA numbers, got %d", d.data.WANumCount)
	}
	if d.data.ConvRate != "5.1%" {
		t.Errorf("expected 5.1%%, got %s", d.data.ConvRate)
	}
	if d.data.BestDay != "selasa" {
		t.Errorf("expected selasa, got %s", d.data.BestDay)
	}
}

func TestDashboardWASlotsFromParams(t *testing.T) {
	d := NewDashboard()
	d.HandleNavigate(map[string]any{
		"wa_slots": []any{
			map[string]any{
				"label":  "slot-1",
				"number": "0812-xxxx-3456",
				"active": true,
				"hours":  "4/6 jam",
			},
		},
	})
	if len(d.data.WASlots) != 1 {
		t.Fatalf("expected 1 WA slot, got %d", len(d.data.WASlots))
	}
	if d.data.WASlots[0].Label != "slot-1" {
		t.Errorf("expected slot-1, got %s", d.data.WASlots[0].Label)
	}
	if !d.data.WASlots[0].Active {
		t.Error("expected active slot")
	}
}

// ---------------------------------------------------------------------------
// Response screen tests
// ---------------------------------------------------------------------------

func TestResponseNew(t *testing.T) {
	r := NewResponse()
	if r.ID() != protocol.ScreenResponse {
		t.Errorf("expected ScreenResponse, got %s", r.ID())
	}
}

func TestResponseDefaultView(t *testing.T) {
	r := NewResponse()
	r.lead = LeadResponse{
		Business: "kopi nusantara",
		Category: "cafe",
		Area:     "kediri",
		Message:  "iya kak, boleh lihat desainnya?",
	}
	view := r.View()
	if view == "" {
		t.Error("expected non-empty view from response")
	}
}

func TestResponseHandleNavigate(t *testing.T) {
	r := NewResponse()

	// Navigate to positive response
	err := r.HandleNavigate(map[string]any{"state": "response_positive"})
	if err != nil {
		t.Fatalf("HandleNavigate failed: %v", err)
	}
	if r.state != protocol.ResponsePositive {
		t.Errorf("expected ResponsePositive, got %s", r.state)
	}

	// Navigate to stop detected
	err = r.HandleNavigate(map[string]any{"state": "response_stop_detected"})
	if err != nil {
		t.Fatalf("HandleNavigate failed: %v", err)
	}
	if r.state != protocol.ResponseStopDetected {
		t.Errorf("expected ResponseStopDetected, got %s", r.state)
	}

	// Navigate to deal detected
	err = r.HandleNavigate(map[string]any{"state": "response_deal_detected"})
	if err != nil {
		t.Fatalf("HandleNavigate failed: %v", err)
	}
	if r.state != protocol.ResponseDealDetected {
		t.Errorf("expected ResponseDealDetected, got %s", r.state)
	}

	// Navigate to hot lead
	err = r.HandleNavigate(map[string]any{"state": "response_hot_lead"})
	if err != nil {
		t.Fatalf("HandleNavigate failed: %v", err)
	}
	if r.state != protocol.ResponseHotLead {
		t.Errorf("expected ResponseHotLead, got %s", r.state)
	}
}

func TestResponseAllStates(t *testing.T) {
	r := NewResponse()
	r.lead = LeadResponse{
		Business: "kopi nusantara",
		Category: "cafe",
		Area:     "kediri",
		Message:  "iya kak",
	}

	states := []protocol.StateID{
		protocol.ResponsePositive,
		protocol.ResponseCurious,
		protocol.ResponseNegative,
		protocol.ResponseMaybe,
		protocol.ResponseAutoReply,
		protocol.ResponseStopDetected,
		protocol.ResponseDealDetected,
		protocol.ResponseHotLead,
		protocol.ResponseOfferPreview,
		protocol.ResponseMultiQueue,
	}

	for _, state := range states {
		r.state = state
		view := r.View()
		if view == "" {
			t.Errorf("expected non-empty view for state %s", state)
		}
	}
}

func TestResponseConversionDrama(t *testing.T) {
	r := NewResponse()
	r.width = 80
	r.height = 24
	r.conversion = ConversionData{
		Business:    "kopi nusantara",
		Pipeline:    "ice breaker → offer → deal",
		TimeTaken:   "2 hari 4 jam",
		TrophyCount: 3,
		Revenue:     "rp 7.5jt",
	}
	r.lead = LeadResponse{
		Business: "kopi nusantara",
		Category: "cafe",
		Area:     "kediri",
	}

	// Start conversion drama
	r.HandleNavigate(map[string]any{"state": "response_conversion"})

	if r.dramaPhase == dramaNone {
		t.Error("expected drama to start after navigating to conversion")
	}
	if r.dramaPhase != dramaShock {
		t.Errorf("expected dramaShock phase, got %d", r.dramaPhase)
	}

	// Verify shock view
	view := r.View()
	if view == "" {
		t.Error("expected non-empty shock view")
	}

	// Simulate time passing to reach settle phase
	r.dramaStart = time.Now().Add(-3 * time.Second)
	r.advanceDrama(time.Now())
	if r.dramaPhase != dramaSettle {
		t.Errorf("expected dramaSettle, got %d", r.dramaPhase)
	}

	// Verify key accepted after settle hold
	if !r.keyAccepted {
		t.Error("expected keyAccepted after settle hold")
	}
}

func TestResponseLeadDataFromParams(t *testing.T) {
	r := NewResponse()
	r.HandleNavigate(map[string]any{
		"state": "response_positive",
		"lead": map[string]any{
			"business": "kopi nusantara",
			"category": "cafe",
			"area":     "kediri",
			"message":  "iya kak, boleh lihat desainnya?",
		},
	})
	if r.lead.Business != "kopi nusantara" {
		t.Errorf("expected kopi nusantara, got %s", r.lead.Business)
	}
	if r.lead.Category != "cafe" {
		t.Errorf("expected cafe, got %s", r.lead.Category)
	}
	if r.lead.Area != "kediri" {
		t.Errorf("expected kediri, got %s", r.lead.Area)
	}
}

// ---------------------------------------------------------------------------
// i18n integration tests — verify key constants resolve
// ---------------------------------------------------------------------------

func TestMonitorI18nKeysResolve(t *testing.T) {
	// Ensure all monitor i18n keys resolve to non-empty strings
	keys := []string{
		i18n.KeyMonitorLive, i18n.KeyMonitorIdle,
		i18n.KeyMonitorNight, i18n.KeyMonitorError,
		i18n.KeyMonitorEmpty, i18n.KeyMonitorPending,
		i18n.KeyMonitorWAConnected, i18n.KeyMonitorWADisconnected,
		i18n.KeyMonitorNoData, i18n.KeyMonitorRecentActivity,
	}
	for _, k := range keys {
		val := i18n.T(k)
		if val == "" {
			t.Errorf("i18n key %q resolved to empty string", k)
		}
	}
}

func TestResponseI18nKeysResolve(t *testing.T) {
	keys := []string{
		i18n.KeyResponsePositive, i18n.KeyResponseCurious,
		i18n.KeyResponseNegative, i18n.KeyResponseMaybe,
		i18n.KeyResponseAuto, i18n.KeyResponseGotReply,
		i18n.KeyResponseStopDetected, i18n.KeyResponseDealDetected,
		i18n.KeyResponseHotLead, i18n.KeyResponseConversion,
	}
	for _, k := range keys {
		val := i18n.T(k)
		if val == "" {
			t.Errorf("i18n key %q resolved to empty string", k)
		}
	}
}
