package update

import (
        "fmt"
        "strings"
        "testing"

        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
        tea "github.com/charmbracelet/bubbletea"
)

// newTestModel creates a Model with test data pre-populated.
func newTestModel() Model {
        m := NewModel()
        m.currentVersion = "v1.3.2"
        m.newVersion = "v1.3.3"
        // isMajor is no longer a model field — state encodes the distinction.
        m.changelog = []string{
                "fix: scrape duplicate in certain areas",
                "fix: WA rotator cooldown calculation more accurate",
                "perf: database query 15% faster",
        }
        m.changelogStaggerIndex = len(m.changelog) // reveal all for test
        m.licenseKey = "WACL-XXXX-XXXX-XXXX-XXXX"
        m.licenseExpiry = "30 juni 2025"
        m.licenseStatus = string(protocol.LicenseResultValid)
        m.downloadSize = "5.1MB"
        m.downloadDownloaded = "4.2MB"
        m.downloadETA = "12 detik"
        m.downloadSource = "https://releases.waclaw.dev/v1.3.3/"
        m.checksumVerified = true
        m.backupPath = "~/.waclaw/waclaw-v1.3.2.bak"
        m.expiredDate = "15 april 2025"
        m.width = 80
        m.height = 24
        return m
}

// TestNewModel verifies the model initialises with correct defaults.
func TestNewModel(t *testing.T) {
        m := NewModel()
        if m.ID() != protocol.ScreenUpdate {
                t.Errorf("expected ID %s, got %s", protocol.ScreenUpdate, m.ID())
        }
        if m.state != protocol.UpdateAvailable {
                t.Errorf("expected initial state %s, got %s", protocol.UpdateAvailable, m.state)
        }
}

// TestUpdateAvailableState tests the update_available state rendering.
func TestUpdateAvailableState(t *testing.T) {
        m := newTestModel()
        m.state = protocol.UpdateAvailable

        view := m.View()

        if !strings.Contains(view, i18n.T(i18n.KeyUpdateAvailable)) {
                t.Errorf("expected view to contain update available title")
        }
        if !strings.Contains(view, m.currentVersion) {
                t.Errorf("expected view to contain current version %s", m.currentVersion)
        }
        if !strings.Contains(view, m.newVersion) {
                t.Errorf("expected view to contain new version %s", m.newVersion)
        }
        if !strings.Contains(view, i18n.T(i18n.KeyUpdateMinorFree)) {
                t.Errorf("expected view to contain free badge text")
        }
}

// TestDownloadingState tests the update_downloading state rendering.
func TestDownloadingState(t *testing.T) {
        m := newTestModel()
        m.state = protocol.UpdateDownloading
        m.downloadPercent = 0.82
        m.downloadBar = m.downloadBar.SetPercent(0.82)

        view := m.View()

        if !strings.Contains(view, m.newVersion) {
                t.Errorf("expected view to contain version %s", m.newVersion)
        }
        if !strings.Contains(view, m.downloadSource) {
                t.Errorf("expected view to contain download source")
        }
}

// TestReadyState tests the update_ready state rendering.
func TestReadyState(t *testing.T) {
        m := newTestModel()
        m.state = protocol.UpdateReady

        view := m.View()

        if !strings.Contains(view, m.newVersion) {
                t.Errorf("expected view to contain version %s in ready state", m.newVersion)
        }
        if !strings.Contains(view, m.downloadSize) {
                t.Errorf("expected view to contain download size %s", m.downloadSize)
        }
        if !strings.Contains(view, m.backupPath) {
                t.Errorf("expected view to contain backup path %s", m.backupPath)
        }
}

// TestUpgradeAvailableState tests the upgrade_available state rendering.
func TestUpgradeAvailableState(t *testing.T) {
        m := newTestModel()
        m.state = protocol.UpgradeAvailable
        m.newVersion = "v2.0.0"

        view := m.View()

        if !strings.Contains(view, i18n.T(i18n.KeyUpgradeAvailable)) {
                t.Errorf("expected view to contain upgrade available title")
        }
        if !strings.Contains(view, i18n.T(i18n.KeyUpgradeMajorWarning)) {
                t.Errorf("expected view to contain major warning text")
        }
        if !strings.Contains(view, i18n.T(i18n.KeyUpgradeNeedsLicense)) {
                t.Errorf("expected view to contain license requirement text")
        }
        if !strings.Contains(view, m.licenseKey) {
                t.Errorf("expected view to contain v1 license key")
        }
}

// TestLicenseInputState tests the upgrade_license_input state rendering.
func TestLicenseInputState(t *testing.T) {
        m := newTestModel()
        m.state = protocol.UpgradeLicenseInput
        m.licenseInputFocus()

        view := m.View()

        if !strings.Contains(view, i18n.T(i18n.KeyUpgradeV2Title)) {
                t.Errorf("expected view to contain v2 title")
        }
}

// TestExpiredWithUpgradeState tests the license_expired_with_upgrade state rendering.
func TestExpiredWithUpgradeState(t *testing.T) {
        m := newTestModel()
        m.state = protocol.LicenseExpiredWithUpgrade
        m.newVersion = "v2.0.0"

        view := m.View()

        if !strings.Contains(view, i18n.T(i18n.KeyLicenseExpiredTitle)) {
                t.Errorf("expected view to contain license expired title")
        }
        if !strings.Contains(view, i18n.T(i18n.KeyLicenseExpiredV2Features)) {
                t.Errorf("expected view to contain v2 features text")
        }
}

// TestHandleNavigate tests the HandleNavigate method.
func TestHandleNavigate(t *testing.T) {
        m := NewModel()
        err := m.HandleNavigate(map[string]any{
                protocol.ParamState:          string(protocol.UpdateDownloading),
                protocol.ParamCurrentVersion: "v1.3.2",
                protocol.ParamNewVersion:     "v1.3.3",
                protocol.ParamChangelog:      []any{"fix: bug", "perf: faster"},
                protocol.ParamLicenseKey:     "WACL-TEST",
                protocol.ParamLicenseExpiry:  "30 juni 2025",
                protocol.ParamLicenseStatus:  string(protocol.LicenseResultValid),
        })
        if err != nil {
                t.Fatalf("HandleNavigate returned error: %v", err)
        }

        if m.state != protocol.UpdateDownloading {
                t.Errorf("expected state %s, got %s", protocol.UpdateDownloading, m.state)
        }
        if m.currentVersion != "v1.3.2" {
                t.Errorf("expected current version v1.3.2, got %s", m.currentVersion)
        }
        if m.newVersion != "v1.3.3" {
                t.Errorf("expected new version v1.3.3, got %s", m.newVersion)
        }
        if len(m.changelog) != 2 {
                t.Errorf("expected 2 changelog items, got %d", len(m.changelog))
        }
}

// TestHandleNavigateNil tests HandleNavigate with nil params.
func TestHandleNavigateNil(t *testing.T) {
        m := NewModel()
        err := m.HandleNavigate(nil)
        if err != nil {
                t.Fatalf("HandleNavigate(nil) returned error: %v", err)
        }
}

// TestHandleUpdate tests the HandleUpdate method for download progress.
func TestHandleUpdate(t *testing.T) {
        m := NewModel()
        m.state = protocol.UpdateDownloading

        err := m.HandleUpdate(map[string]any{
                protocol.ParamPercent:    0.82,
                protocol.ParamSize:       "5.1MB",
                protocol.ParamDownloaded: "4.2MB",
                protocol.ParamSpeed:      "350KB/s",
                protocol.ParamETA:        "12 detik",
                protocol.ParamSource:     "https://releases.waclaw.dev/v1.3.3/",
        })
        if err != nil {
                t.Fatalf("HandleUpdate returned error: %v", err)
        }

        if m.downloadPercent != 0.82 {
                t.Errorf("expected download percent 0.82, got %f", m.downloadPercent)
        }
        if m.downloadSize != "5.1MB" {
                t.Errorf("expected download size 5.1MB, got %s", m.downloadSize)
        }
        if m.downloadETA != "12 detik" {
                t.Errorf("expected ETA '12 detik', got %s", m.downloadETA)
        }
}

// TestHandleUpdateNil tests HandleUpdate with nil params.
func TestHandleUpdateNil(t *testing.T) {
        m := NewModel()
        err := m.HandleUpdate(nil)
        if err != nil {
                t.Fatalf("HandleUpdate(nil) returned error: %v", err)
        }
}

// TestHandleUpdateStateTransition tests HandleUpdate for state changes.
func TestHandleUpdateStateTransition(t *testing.T) {
        m := NewModel()
        m.state = protocol.UpdateDownloading

        err := m.HandleUpdate(map[string]any{
                protocol.ParamState:            string(protocol.UpdateReady),
                protocol.ParamChecksumVerified: true,
                protocol.ParamBackupPath:       "~/.waclaw/waclaw-v1.3.2.bak",
        })
        if err != nil {
                t.Fatalf("HandleUpdate returned error: %v", err)
        }

        if m.state != protocol.UpdateReady {
                t.Errorf("expected state %s, got %s", protocol.UpdateReady, m.state)
        }
        if !m.checksumVerified {
                t.Error("expected checksum to be verified")
        }
}

// TestSendActionWithBus tests that sendAction publishes to the bus.
func TestSendActionWithBus(t *testing.T) {
        m := NewModel()
        b := bus.New()
        m.SetBus(b)

        // Subscribe to capture the action.
        var received bus.ActionMsg
        b.Subscribe(func(msg any) bool {
                if action, ok := msg.(bus.ActionMsg); ok {
                        received = action
                        return true
                }
                return false
        })

        m.state = protocol.UpdateAvailable
        _ = m.handleUpdateAvailableKeys(createKeyMsg("enter"))

        // Poll to deliver pending messages.
        b.Pending()

        if received.Action != string(protocol.ActionStartDownload) {
                t.Errorf("expected action %s, got %s", protocol.ActionStartDownload, received.Action)
        }
        if received.Screen != protocol.ScreenUpdate {
                t.Errorf("expected screen %s, got %s", protocol.ScreenUpdate, received.Screen)
        }
}

// TestKeyHandlerDownloadingCancel tests the cancel download key handler.
func TestKeyHandlerDownloadingCancel(t *testing.T) {
        m := NewModel()
        b := bus.New()
        m.SetBus(b)

        var received bus.ActionMsg
        b.Subscribe(func(msg any) bool {
                if action, ok := msg.(bus.ActionMsg); ok {
                        received = action
                        return true
                }
                return false
        })

        m.state = protocol.UpdateDownloading
        _ = m.handleDownloadingKeys(createKeyMsg("q"))
        b.Pending()

        if received.Action != string(protocol.ActionCancelDownload) {
                t.Errorf("expected action %s, got %s", protocol.ActionCancelDownload, received.Action)
        }
        // State transition is backend-driven: the TUI sends "cancel_download" action
        // and the backend confirms by pushing a HandleUpdate with the new state.
        // The TUI should NOT directly mutate state on key press.
        if m.state != protocol.UpdateDownloading {
                t.Errorf("expected state to remain UpdateDownloading until backend confirms, got %s", m.state)
        }
}

// TestI18NBothLocales verifies that all update screen keys resolve
// in both supported locales.
func TestI18NBothLocales(t *testing.T) {
        keys := []string{
                i18n.KeyUpdateAvailable, i18n.KeyUpdateDownloading, i18n.KeyUpdateReady,
                i18n.KeyUpgradeAvailable, i18n.KeyUpgradeLicense,
                i18n.KeyUpdateNow, i18n.KeyUpdateLater, i18n.KeyUpdateRestartNote,
                i18n.KeyUpdateMinorFree, i18n.KeyUpdateLicenseValid,
                i18n.KeyUpgradeMajorWarning, i18n.KeyUpgradeNeedsLicense,
                i18n.KeyUpgradeV1StillWorks, i18n.KeyUpgradeNoForced,
                i18n.KeyUpgradeBuyV2, i18n.KeyUpgradeStayV1,
                i18n.KeyLicenseExpiredTitle, i18n.KeyLicenseExpiredButV2,
                i18n.KeyUpgradeLicensePlaceholder, i18n.KeyUpgradeLicenseKey,
                i18n.KeyUpgradeLicenseExpires, i18n.KeyUpgradeLicenseStatus,
                i18n.KeyStartupCheckTitle, i18n.KeyStartupCheckLatest,
                i18n.KeyStartupCheckUpdate, i18n.KeyStartupCheckUpgrade,
                i18n.KeyUpdateSpeed,
                i18n.KeyUpdateChecksumPending,
        }

        for _, key := range keys {
                i18n.SetLocale(i18n.LocaleID)
                idVal := i18n.T(key)
                if idVal == key {
                        t.Errorf("key %q not found in Indonesian locale", key)
                }

                i18n.SetLocale(i18n.LocaleEN)
                enVal := i18n.T(key)
                if enVal == key {
                        t.Errorf("key %q not found in English locale", key)
                }
        }

        // Restore default locale.
        i18n.SetLocale(i18n.LocaleID)
}

// TestFocusBlur tests the Focus/Blur lifecycle.
func TestFocusBlur(t *testing.T) {
        m := NewModel()
        if m.focused {
                t.Error("expected model to start unfocused")
        }

        m.Focus()
        if !m.focused {
                t.Error("expected model to be focused after Focus()")
        }

        m.Blur()
        if m.focused {
                t.Error("expected model to be unfocused after Blur()")
        }
        if m.freeBadgePulsing {
                t.Error("expected free badge pulsing to stop on Blur()")
        }
}

// TestViewAllStates verifies that every state renders without panicking.
func TestViewAllStates(t *testing.T) {
        states := []protocol.StateID{
                protocol.UpdateAvailable,
                protocol.UpdateDownloading,
                protocol.UpdateReady,
                protocol.UpgradeAvailable,
                protocol.UpgradeLicenseInput,
                protocol.LicenseExpiredWithUpgrade,
                protocol.StartupCheck,
        }

        for _, s := range states {
                t.Run(fmt.Sprintf("state_%s", s), func(t *testing.T) {
                        m := newTestModel()
                        m.state = s
                        if s == protocol.UpgradeLicenseInput {
                                m.licenseInputFocus()
                        }
                        if s == protocol.StartupCheck {
                                m.startupCheckResult = string(protocol.CheckResultLatest)
                        }
                        if s == protocol.UpgradeAvailable {
                                m.newVersion = "v2.0.0"
                        }

                        view := m.View()
                        if view == "" {
                                t.Errorf("state %s rendered empty view", s)
                        }
                })
        }
}

// createKeyMsg creates a tea.KeyMsg for testing key handlers.
func createKeyMsg(key string) tea.KeyMsg {
        // Use the rune-based key message type for simple key presses.
        return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
}
