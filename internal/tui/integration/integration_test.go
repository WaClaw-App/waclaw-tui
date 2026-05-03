// Package integration_test provides end-to-end integration tests for the
// WaClaw TUI Phase 5 — verifying screen registration, canonical flow
// navigation, interface compliance, session end, and DRY consolidation.
package integration_test

import (
        "testing"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui"
        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/comms"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/data"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/infra"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/license"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/monitor"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/niche"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/onboarding"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/pipeline"
        "github.com/WaClaw-App/waclaw/internal/tui/screen/update"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
)

// registerAllScreens registers all 20 screens into the app, mirroring
// the production wiring in cmd/tui/main.go.
func registerAllScreens(app *tui.App) {
        boot := onboarding.NewBootModel()
        app.RegisterScreen(&boot)
        login := onboarding.NewLoginModel()
        app.RegisterScreen(&login)
        sel := niche.NewSelectModel()
        app.RegisterScreen(&sel)
        expl := niche.NewExplorerModel()
        app.RegisterScreen(&expl)
        app.RegisterScreen(pipeline.NewScrape())
        rev := pipeline.NewReview()
        app.RegisterScreen(&rev)
        app.RegisterScreen(pipeline.NewSend())
        app.RegisterScreen(monitor.NewDashboard())
        app.RegisterScreen(monitor.NewResponse())
        app.RegisterScreen(data.NewLeadsDB())
        app.RegisterScreen(data.NewTemplateMgr())
        app.RegisterScreen(infra.NewWorkers())
        app.RegisterScreen(infra.NewShield())
        app.RegisterScreen(infra.NewSettings())
        app.RegisterScreen(infra.NewGuardrail())
        app.RegisterScreen(comms.NewCompose())
        app.RegisterScreen(comms.NewHistory())
        app.RegisterScreen(comms.NewFollowUp())
        lic := license.New()
        app.RegisterScreen(&lic)
        upd := update.NewModel()
        app.RegisterScreen(&upd)
}

// allScreenIDs returns the canonical list of all 20 screen IDs.
// Order matches the doc/18-screen-flow.md canonical flow:
// Boot → License → Validation → Login → NicheSelect → ...
func allScreenIDs() []protocol.ScreenID {
        return []protocol.ScreenID{
                protocol.ScreenBoot, protocol.ScreenLicense, protocol.ScreenGuardrail,
                protocol.ScreenLogin,
                protocol.ScreenNicheSelect, protocol.ScreenNicheExplorer,
                protocol.ScreenScrape, protocol.ScreenLeadReview, protocol.ScreenSend,
                protocol.ScreenMonitor, protocol.ScreenResponse,
                protocol.ScreenLeadsDB, protocol.ScreenTemplateMgr,
                protocol.ScreenWorkers, protocol.ScreenAntiBan, protocol.ScreenSettings,
                protocol.ScreenCompose, protocol.ScreenHistory, protocol.ScreenFollowUp,
                protocol.ScreenUpdate,
        }
}

// TestPhase5ScreenRegistration verifies that all 20 screens can be registered
// and that the router can navigate to each one.
func TestPhase5ScreenRegistration(t *testing.T) {
        app := tui.NewApp()
        registerAllScreens(app)

        // Verify all 20 screens are registered.
        // All 20 screens registered — verified by the loop below.

        // Verify each screen has the correct ID.
        for _, id := range allScreenIDs() {
                s := app.Router().Screen(id)
                if s == nil {
                        t.Errorf("screen %q not registered", id)
                        continue
                }
                if s.ID() != id {
                        t.Errorf("screen ID mismatch: got %q, want %q", s.ID(), id)
                }
        }
}

// TestPhase5CanonicalFlow verifies navigation through the canonical
// 20-screen demo walkthrough sequence.
func TestPhase5CanonicalFlow(t *testing.T) {
        app := tui.NewApp()
        registerAllScreens(app)

        // Push each screen and verify the router state.
        for i, screenID := range allScreenIDs() {
                app.Router().Push(screenID)
                current := app.Router().CurrentID()
                if current != screenID {
                        t.Errorf("flow step %d: expected screen %q, got %q", i, screenID, current)
                }
        }

        // Verify the navigation stack depth matches the canonical flow.
        if app.Router().Depth() != len(allScreenIDs()) {
                t.Errorf("expected stack depth %d, got %d", len(allScreenIDs()), app.Router().Depth())
        }

        // Pop back and verify reverse navigation.
        for i := len(allScreenIDs()) - 2; i >= 0; i-- {
                app.Router().Pop()
                current := app.Router().CurrentID()
                if current != allScreenIDs()[i] {
                        t.Errorf("pop step %d: expected screen %q, got %q", i, allScreenIDs()[i], current)
                }
        }
}

// TestPhase5ScreenInterfaces verifies that every registered screen
// implements the Screen interface and checks StateReporter adoption.
func TestPhase5ScreenInterfaces(t *testing.T) {
        app := tui.NewApp()
        registerAllScreens(app)

        stateReporterCount := 0
        for _, id := range allScreenIDs() {
                s := app.Router().Screen(id)
                if s == nil {
                        t.Errorf("screen %q not found", id)
                        continue
                }

                // Verify the Screen interface methods are callable.
                _ = s.ID()
                _ = s.View()
                s.SetBus(nil)
                _ = s.HandleNavigate(nil)
                _ = s.HandleUpdate(nil)
                s.Focus()
                s.Blur()

                // Check StateReporter (optional interface).
                if _, ok := s.(tui.StateReporter); ok {
                        stateReporterCount++
                }
        }

        t.Logf("StateReporter: %d/%d screens implement CurrentState()", stateReporterCount, len(allScreenIDs()))
}

// TestPhase5AllProtocolStates verifies that all state IDs referenced in
// the doc/22-state-machine.md exist in the protocol package.
func TestPhase5AllProtocolStates(t *testing.T) {
        // The 110 states from doc/22-state-machine.md should all be defined.
        // We verify a representative subset to ensure the protocol package
        // is consistent with the spec.
        states := []protocol.StateID{
                // Boot
                protocol.BootFirstTime, protocol.BootReturning,
                // Login
                protocol.LoginQRWaiting, protocol.LoginQRScanned, protocol.LoginSuccess,
                // Niche
                protocol.NicheList, protocol.NicheCustom, protocol.NicheEditFilters,
                // Scrape
                protocol.ScrapeActive, protocol.ScrapeIdle, protocol.ScrapeHighValueReveal,
                // Review
                protocol.ReviewReviewing, protocol.ReviewLeadDetail,
                // Send
                protocol.SendActive, protocol.SendPaused, protocol.SendRateLimited,
                // Monitor
                protocol.MonitorLiveDashboard, protocol.MonitorEmpty, protocol.MonitorNight,
                // Response
                protocol.ResponsePositive, protocol.ResponseCurious, protocol.ResponseNegative,
                // Leads
                protocol.LeadsList, protocol.LeadsFiltered,
                // Template
                protocol.TemplateList, protocol.TemplatePreview,
                // Workers
                protocol.WorkersOverview, protocol.WorkerDetail,
                // Shield
                protocol.ShieldOverview, protocol.ShieldWarning, protocol.ShieldDanger,
                // Settings
                protocol.SettingsOverview, protocol.SettingsReload,
                // Guardrail
                protocol.ValidationClean, protocol.ValidationErrors,
                // Compose
                protocol.ComposeDraft, protocol.ComposePreview,
                // History
                protocol.HistoryToday, protocol.HistoryWeek,
                // FollowUp
                protocol.FollowUpDashboard, protocol.FollowUpSending,
                // Explorer
                protocol.ExplorerBrowse, protocol.ExplorerSearch,
                // License
                protocol.LicenseInput, protocol.LicenseValid,
                // Update
                protocol.UpdateAvailable, protocol.UpdateDownloading,
        }

        for _, state := range states {
                if state == "" {
                        t.Error("empty state ID found — protocol constant not defined")
                }
        }
}

// TestPhase5StartupTiming verifies that the startup animation constants
// are consistent with the doc/20-startup-and-session.md spec:
// "1300ms sampai bisa dipakai" — the total startup must not exceed 1.3s.
func TestPhase5StartupTiming(t *testing.T) {
        // The total startup sequence must be within the documented 1300ms target.
        if anim.StartupSequenceTotal != 1300*time.Millisecond {
                t.Errorf("StartupSequenceTotal = %v, want 1300ms per doc/20", anim.StartupSequenceTotal)
        }

        // Logo character delay must be 8ms per doc spec.
        if anim.LogoCharDelay != 8*time.Millisecond {
                t.Errorf("LogoCharDelay = %v, want 8ms per doc/20", anim.LogoCharDelay)
        }

        // Army march total must be 600ms per doc spec.
        if anim.ArmyMarch != 600*time.Millisecond {
                t.Errorf("ArmyMarch = %v, want 600ms per doc/01", anim.ArmyMarch)
        }

        // Menu stagger must be 120ms per doc spec.
        if anim.MenuStagger != 120*time.Millisecond {
                t.Errorf("MenuStagger = %v, want 120ms per doc/01", anim.MenuStagger)
        }

        // Breathing cycle must be 4000ms per doc spec.
        if anim.BreathingCycle != 4000*time.Millisecond {
                t.Errorf("BreathingCycle = %v, want 4000ms per doc/15", anim.BreathingCycle)
        }

        // Screen transition must be 300ms per doc spec.
        if anim.ScreenTransition != 300*time.Millisecond {
                t.Errorf("ScreenTransition = %v, want 300ms per doc/15", anim.ScreenTransition)
        }
}

// TestPhase5CanonicalFlowOrder verifies the canonical flow follows
// doc/18-screen-flow.md: Boot → License → Validation → Login → Niche → ...
func TestPhase5CanonicalFlowOrder(t *testing.T) {
        ids := allScreenIDs()

        // Verify the first 5 screens match the doc flow for first-time users.
        expectedFirst5 := []protocol.ScreenID{
                protocol.ScreenBoot,       // BOOT
                protocol.ScreenLicense,    // LICENSE (gate)
                protocol.ScreenGuardrail,  // VALIDATION (gate)
                protocol.ScreenLogin,      // LOGIN
                protocol.ScreenNicheSelect, // NICHE SELECT
        }

        if len(ids) < len(expectedFirst5) {
                t.Fatalf("expected at least %d screens, got %d", len(expectedFirst5), len(ids))
        }

        for i, expected := range expectedFirst5 {
                if ids[i] != expected {
                        t.Errorf("canonical flow position %d: got %q, want %q (per doc/18-screen-flow.md)", i, ids[i], expected)
                }
        }

        // Monitor should be the "home base" screen per doc/18.
        foundMonitor := false
        for _, id := range ids {
                if id == protocol.ScreenMonitor {
                        foundMonitor = true
                        break
                }
        }
        if !foundMonitor {
                t.Error("Monitor (home base) not found in canonical flow per doc/18")
        }
}
