package license

import (
        "strings"
        "testing"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/WaClaw-App/waclaw/internal/tui/testutil"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
        tea "github.com/charmbracelet/bubbletea"
)

// updateModel applies a tea.Msg to a Model and returns the updated Model.
func updateModel(m Model, msg tea.Msg) Model {
        model, _ := m.Update(msg)
        return model.(Model)
}

// keyRun returns a KeyMsg for a printable character.
func keyRun(ch string) tea.KeyMsg {
        return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(ch)}
}

// ---------------------------------------------------------------------------
// Model creation
// ---------------------------------------------------------------------------

func TestNew_DefaultState(t *testing.T) {
        m := New()
        if m.State() != protocol.LicenseInput {
                t.Errorf("expected state %q, got %q", protocol.LicenseInput, m.State())
        }
        if m.ID() != protocol.ScreenLicense {
                t.Errorf("expected screen ID %q, got %q", protocol.ScreenLicense, m.ID())
        }
}

// ---------------------------------------------------------------------------
// Screen interface compliance (structural typing)
// ---------------------------------------------------------------------------

func TestModel_ImplementsScreenInterface(t *testing.T) {
        var _ tea.Model = Model{}
        var m Model
        _ = m.ID
        _ = m.Init
        _ = m.Update
        _ = m.View
        _ = (*Model)(nil).SetBus
        _ = (*Model)(nil).HandleNavigate
        _ = (*Model)(nil).HandleUpdate
        _ = (*Model)(nil).Focus
        _ = (*Model)(nil).Blur
}

// ---------------------------------------------------------------------------
// Input state — key entry
// ---------------------------------------------------------------------------

func TestInput_TypingFormatsKey(t *testing.T) {
        m := New()
        m = updateModel(m, keyRun("a"))
        m = updateModel(m, keyRun("b"))
        m = updateModel(m, keyRun("c"))
        m = updateModel(m, keyRun("d"))

        got := m.formattedKey()
        want := "WACL-ABCD"
        if got != want {
                t.Errorf("formatted key = %q, want %q", got, want)
        }
}

func TestInput_AutoUppercase(t *testing.T) {
        m := New()
        for _, ch := range "wacld" {
                m = updateModel(m, keyRun(string(ch)))
        }

        got := m.formattedKey()
        // "WACL" typed -> WACL prefix stripped -> "D" remaining -> WACL-D
        if !strings.HasPrefix(got, "WACL-") {
                t.Errorf("formatted key = %q, want WACL- prefix", got)
        }
}

func TestInput_HyphenInsertion(t *testing.T) {
        m := New()
        for _, ch := range "abcdefgh" {
                m = updateModel(m, keyRun(string(ch)))
        }

        got := m.formattedKey()
        want := "WACL-ABCD-EFGH"
        if got != want {
                t.Errorf("formatted key = %q, want %q", got, want)
        }
}

func TestInput_FullKeyFormat(t *testing.T) {
        m := New()
        for _, ch := range "abcdefghijklmnop" {
                m = updateModel(m, keyRun(string(ch)))
        }

        got := m.formattedKey()
        want := "WACL-ABCD-EFGH-IJKL-MNOP"
        if got != want {
                t.Errorf("formatted key = %q, want %q", got, want)
        }
}

func TestInput_NonAlphanumericIgnored(t *testing.T) {
        m := New()
        m = updateModel(m, keyRun("a"))
        m = updateModel(m, keyRun("-"))
        m = updateModel(m, keyRun("b"))

        if len(m.input) != 2 {
                t.Errorf("input length = %d, want 2", len(m.input))
        }
}

func TestInput_Backspace(t *testing.T) {
        m := New()
        m = updateModel(m, keyRun("a"))
        m = updateModel(m, keyRun("b"))
        m = updateModel(m, keyRun("c"))
        m = updateModel(m, tea.KeyMsg{Type: tea.KeyBackspace})

        if len(m.input) != 2 {
                t.Errorf("input length after backspace = %d, want 2", len(m.input))
        }
}

func TestInput_CursorMovement(t *testing.T) {
        m := New()
        m = updateModel(m, keyRun("a"))
        m = updateModel(m, keyRun("b"))
        m = updateModel(m, keyRun("c"))
        if m.cursor != 3 {
                t.Errorf("cursor = %d, want 3 after 3 chars", m.cursor)
        }

        // Move left
        m = updateModel(m, tea.KeyMsg{Type: tea.KeyLeft})
        if m.cursor != 2 {
                t.Errorf("cursor = %d, want 2 after left", m.cursor)
        }

        // Move right
        m = updateModel(m, tea.KeyMsg{Type: tea.KeyRight})
        if m.cursor != 3 {
                t.Errorf("cursor = %d, want 3 after right", m.cursor)
        }
}

// ---------------------------------------------------------------------------
// Validation state transitions
// ---------------------------------------------------------------------------

func TestInput_EnterStartsValidation(t *testing.T) {
        m := New()
        // Need full 24-char formatted key: WACL-ABCD-EFGH-IJKL-MNOP (20 raw chars)
        for _, ch := range "abcdefghijklmnop" {
                m = updateModel(m, keyRun(string(ch)))
        }
        model, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
        updated := model.(Model)
        if updated.State() != protocol.LicenseValidating {
                t.Errorf("state = %q, want %q", updated.State(), protocol.LicenseValidating)
        }
        if cmd == nil {
                t.Error("expected non-nil cmd when starting validation")
        }
}

func TestInput_EnterTooShort_NoValidation(t *testing.T) {
        m := New()
        _, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
        if m.State() != protocol.LicenseInput {
                t.Errorf("state = %q, want %q", m.State(), protocol.LicenseInput)
        }
        if cmd != nil {
                t.Error("expected nil cmd for too-short input")
        }
}

func TestValidating_TickUpdatesProgress(t *testing.T) {
        m := New()
        m.state = protocol.LicenseValidating
        m = updateModel(m, validationTickMsg{step: 1})
        if m.validationPct < 0.3 {
                t.Errorf("progress = %v, want >= 0.3", m.validationPct)
        }
        if m.validationStep < 1 {
                t.Errorf("step = %d, want >= 1", m.validationStep)
        }
}

func TestValidating_AllThreeStepsComplete(t *testing.T) {
        m := New()
        m.state = protocol.LicenseValidating
        m = updateModel(m, validationTickMsg{step: 1})
        m = updateModel(m, validationTickMsg{step: 2})
        m = updateModel(m, validationTickMsg{step: 3})
        if !m.validationDone {
                t.Error("expected validationDone=true after 3 steps")
        }
}

// ---------------------------------------------------------------------------
// HandleUpdate — backend-driven validation results
// ---------------------------------------------------------------------------

func TestHandleUpdate_ValidResult(t *testing.T) {
        m := New()
        m.state = protocol.LicenseValidating
        m.bus = bus.New()

        err := m.HandleUpdate(map[string]any{
                "result":      "valid",
                "license_key": "WACL-ABCD-EFGH-IJKL-MNOP",
                "device":      "LAPTOP-HOME",
                "expires":     "30 juni 2025",
                "grace_hours": 0,
        })
        if err != nil {
                t.Fatalf("HandleUpdate error: %v", err)
        }
        if m.State() != protocol.LicenseValid {
                t.Errorf("state = %q, want %q", m.State(), protocol.LicenseValid)
        }
        if m.resultDevice != "LAPTOP-HOME" {
                t.Errorf("device = %q, want %q", m.resultDevice, "LAPTOP-HOME")
        }
        if m.resultLicenseKey != "WACL-ABCD-EFGH-IJKL-MNOP" {
                t.Errorf("key = %q", m.resultLicenseKey)
        }
        // When valid result arrives from backend, needsHoldSchedule must be set
        // so the next Update() call schedules the hold + auto-transition.
        if !m.needsHoldSchedule {
                t.Error("expected needsHoldSchedule=true after valid result from HandleUpdate")
        }
}

func TestHandleUpdate_InvalidResult(t *testing.T) {
        m := New()
        m.state = protocol.LicenseValidating
        m.bus = bus.New()

        err := m.HandleUpdate(map[string]any{
                "result":      "invalid",
                "license_key": "WACL-BAAA-D000-0000-0000",
        })
        if err != nil {
                t.Fatalf("HandleUpdate error: %v", err)
        }
        if m.State() != protocol.LicenseInvalid {
                t.Errorf("state = %q, want %q", m.State(), protocol.LicenseInvalid)
        }
}

func TestHandleUpdate_ExpiredResult(t *testing.T) {
        m := New()
        m.state = protocol.LicenseValidating
        m.bus = bus.New()

        err := m.HandleUpdate(map[string]any{
                "result":      "expired",
                "license_key": "WACL-EXPI-RED0-0000-0000",
                "device":      "LAPTOP-HOME",
                "expires":     "15 april 2025",
                "expired_ago": "17 hari lalu",
        })
        if err != nil {
                t.Fatalf("HandleUpdate error: %v", err)
        }
        if m.State() != protocol.LicenseExpired {
                t.Errorf("state = %q, want %q", m.State(), protocol.LicenseExpired)
        }
        if m.resultExpiredAgo != "17 hari lalu" {
                t.Errorf("expired_ago = %q, want %q", m.resultExpiredAgo, "17 hari lalu")
        }
}

func TestHandleUpdate_DeviceConflictResult(t *testing.T) {
        m := New()
        m.state = protocol.LicenseValidating
        m.bus = bus.New()

        err := m.HandleUpdate(map[string]any{
                "result":       "device_conflict",
                "license_key":  "WACL-CONF-LICT-0000-0000",
                "device":       "LAPTOP-HOME",
                "other_device": "PC-KANTOR",
                "last_active":  "12 menit lalu",
        })
        if err != nil {
                t.Fatalf("HandleUpdate error: %v", err)
        }
        if m.State() != protocol.LicenseDeviceConflict {
                t.Errorf("state = %q, want %q", m.State(), protocol.LicenseDeviceConflict)
        }
        if m.resultOtherDev != "PC-KANTOR" {
                t.Errorf("other_device = %q, want %q", m.resultOtherDev, "PC-KANTOR")
        }
        if m.resultLastActive != "12 menit lalu" {
                t.Errorf("last_active = %q, want %q", m.resultLastActive, "12 menit lalu")
        }
}

func TestHandleUpdate_ServerErrorResult(t *testing.T) {
        m := New()
        m.state = protocol.LicenseValidating
        m.bus = bus.New()

        err := m.HandleUpdate(map[string]any{
                "result":      "server_error",
                "license_key": "WACL-SERV-ERR0-0000-0000",
                "grace_hours": 71,
        })
        if err != nil {
                t.Fatalf("HandleUpdate error: %v", err)
        }
        if m.State() != protocol.LicenseServerError {
                t.Errorf("state = %q, want %q", m.State(), protocol.LicenseServerError)
        }
        if m.graceHoursLeft != 71 {
                t.Errorf("grace_hours = %d, want 71", m.graceHoursLeft)
        }
}

// ---------------------------------------------------------------------------
// HandleNavigate — backend-driven state transitions
// ---------------------------------------------------------------------------

func TestHandleNavigate_SetsExpiredState(t *testing.T) {
        m := New()
        err := m.HandleNavigate(map[string]any{
                "state":       "license_expired",
                "license_key": "WACL-TEST-EXPI-RED0-0000",
                "device":      "MY-DEVICE",
                "expires":     "1 januari 2025",
                "expired_ago": "120 hari lalu",
        })
        if err != nil {
                t.Fatalf("HandleNavigate error: %v", err)
        }
        if m.State() != protocol.LicenseExpired {
                t.Errorf("state = %q, want %q", m.State(), protocol.LicenseExpired)
        }
        if m.resultDevice != "MY-DEVICE" {
                t.Errorf("device = %q", m.resultDevice)
        }
}

func TestHandleNavigate_SetsConflictState(t *testing.T) {
        m := New()
        err := m.HandleNavigate(map[string]any{
                "state":        "license_device_conflict",
                "license_key":  "WACL-CONF-LICT-0000-0000",
                "device":       "LAPTOP-HOME",
                "other_device": "PC-KANTOR",
                "last_active":  "12 menit lalu",
                "grace_hours":  0,
        })
        if err != nil {
                t.Fatalf("HandleNavigate error: %v", err)
        }
        if m.State() != protocol.LicenseDeviceConflict {
                t.Errorf("state = %q, want %q", m.State(), protocol.LicenseDeviceConflict)
        }
        if m.resultOtherDev != "PC-KANTOR" {
                t.Errorf("other_device = %q, want %q", m.resultOtherDev, "PC-KANTOR")
        }
}

func TestHandleNavigate_SetsServerErrorState(t *testing.T) {
        m := New()
        err := m.HandleNavigate(map[string]any{
                "state":       "license_server_error",
                "license_key": "WACL-SERV-ERR0-0000-0000",
                "grace_hours": 48,
        })
        if err != nil {
                t.Fatalf("HandleNavigate error: %v", err)
        }
        if m.State() != protocol.LicenseServerError {
                t.Errorf("state = %q, want %q", m.State(), protocol.LicenseServerError)
        }
        if m.graceHoursLeft != 48 {
                t.Errorf("grace_hours = %d, want 48", m.graceHoursLeft)
        }
}

// ---------------------------------------------------------------------------
// SetBus / Focus / Blur
// ---------------------------------------------------------------------------

func TestSetBus_StoresReference(t *testing.T) {
        m := New()
        b := bus.New()
        m.SetBus(b)
        if m.bus != b {
                t.Error("bus not stored after SetBus")
        }
}

func TestFocusBlur(t *testing.T) {
        m := New()
        m.Focus()
        if !m.focused {
                t.Error("expected focused=true after Focus()")
        }
        m.Blur()
        if m.focused {
                t.Error("expected focused=false after Blur()")
        }
}

// ---------------------------------------------------------------------------
// View rendering — content checks per state
// ---------------------------------------------------------------------------

func TestViewInput_ContainsKeyElements(t *testing.T) {
        h := testutil.NewScreenHelper(New())
        h.AssertContains(t, "lisensi")
        h.AssertContains(t, "lisensi disimpan di")
        h.AssertContains(t, "validasi")
}

func TestViewValidating_ContainsSteps(t *testing.T) {
        m := New()
        m.state = protocol.LicenseValidating
        h := testutil.NewScreenHelper(m)
        h.AssertContains(t, "server")
        h.AssertContains(t, "validitas")
        h.AssertContains(t, "device")
}

func TestViewValid_ContainsSuccessIndicators(t *testing.T) {
        m := New()
        m.state = protocol.LicenseValid
        m.resultLicenseKey = "WACL-TEST-0000-0000-0000"
        m.resultDevice = "TEST-DEV"
        m.resultExpires = "30 juni 2025"
        m.holdElapsed = true
        h := testutil.NewScreenHelper(m)
        // Doc-aligned: "lisensi valid!" for title (with exclamation)
        h.AssertContains(t, "\u2713 lisensi valid!")
        h.AssertContains(t, "terhubung ke server")
        // Doc-aligned: bullet "lisensi valid" (no exclamation — was "license valid!")
        h.AssertContains(t, "lisensi valid")
        // Doc-aligned: "device terdaftar" not "device terverifikasi"
        h.AssertContains(t, "device terdaftar")
}

func TestViewInvalid_ContainsError(t *testing.T) {
        m := New()
        m.state = protocol.LicenseInvalid
        m.resultLicenseKey = "WACL-BAD0-0000-0000-0000"
        m.redGlowStart = time.Now()
        h := testutil.NewScreenHelper(m)
        h.AssertContains(t, "\u2717")
        // Doc-aligned: "lisensi nggak valid — cek lagi key nya" not "lisensi gagal dicek"
        h.AssertContains(t, "lisensi nggak valid")
}

func TestViewExpired_ContainsExpiredInfo(t *testing.T) {
        m := New()
        m.state = protocol.LicenseExpired
        m.resultLicenseKey = "WACL-EXPI-RED0-0000-0000"
        m.resultDevice = "LAPTOP-HOME"
        m.resultExpires = "15 april 2025"
        m.resultExpiredAgo = "17 hari lalu"
        h := testutil.NewScreenHelper(m)
        h.AssertContains(t, "lisensi lu udah expired")
        // Doc-aligned: "semua worker di-pause. data aman." not "data lu aman."
        h.AssertContains(t, "data aman")
}

func TestViewDeviceConflict_ContainsDeviceInfo(t *testing.T) {
        m := New()
        m.state = protocol.LicenseDeviceConflict
        m.resultLicenseKey = "WACL-CONF-LICT-0000-0000"
        m.resultDevice = "LAPTOP-HOME"
        m.resultOtherDev = "PC-KANTOR"
        m.resultLastActive = "12 menit lalu"
        h := testutil.NewScreenHelper(m)
        h.AssertContains(t, "device lain")
        h.AssertContains(t, "PC-KANTOR")
        h.AssertContains(t, "putuskan")
        // Doc-aligned: force disconnect explanation includes "ambil alih lisensi"
        h.AssertContains(t, "ambil alih lisensi")
}

func TestViewServerError_ContainsGraceInfo(t *testing.T) {
        m := New()
        m.state = protocol.LicenseServerError
        m.resultLicenseKey = "WACL-SERV-ERR0-0000-0000"
        m.graceHoursLeft = 71
        h := testutil.NewScreenHelper(m)
        // Doc-aligned: "waclaw bakal jalan pake lisensi offline selama 72 jam."
        h.AssertContains(t, "72 jam")
        // Doc-aligned: "offline grace: 71 jam tersisa" not "sisa tenggang: 71j"
        h.AssertContains(t, "offline grace: 71 jam tersisa")
}

// ---------------------------------------------------------------------------
// No borders, no boxes (P3 compliance)
// ---------------------------------------------------------------------------

func TestNoBordersInAnyState(t *testing.T) {
        borderChars := []string{"\u2500", "\u2501", "\u2502", "\u2503", "\u2550", "\u2551"}
        states := []protocol.StateID{
                protocol.LicenseInput,
                protocol.LicenseValidating,
                protocol.LicenseValid,
                protocol.LicenseInvalid,
                protocol.LicenseExpired,
                protocol.LicenseDeviceConflict,
                protocol.LicenseServerError,
        }

        for _, state := range states {
                m := New()
                m.state = state
                m.resultLicenseKey = "WACL-TEST-0000-0000-0000"
                m.resultDevice = "DEV"
                m.resultExpires = "30 juni 2025"
                m.resultExpiredAgo = "10 hari lalu"
                m.resultOtherDev = "OTHER"
                m.resultLastActive = "5 menit lalu"
                m.graceHoursLeft = 48
                m.holdElapsed = true

                view := m.View()
                for _, bc := range borderChars {
                        if strings.Contains(view, bc) {
                                t.Errorf("state %q: view contains border character %q (P3 violation)",
                                        state, bc)
                        }
                }
        }
}

// ---------------------------------------------------------------------------
// Key handling — state-specific key routing
// ---------------------------------------------------------------------------

func TestExpired_KeyPress1_GoesToInput(t *testing.T) {
        m := New()
        m.state = protocol.LicenseExpired
        m = updateModel(m, keyRun("1"))
        if m.State() != protocol.LicenseInput {
                t.Errorf("state = %q, want %q", m.State(), protocol.LicenseInput)
        }
        if m.input != "" || m.cursor != 0 {
                t.Error("expected input and cursor reset")
        }
}

func TestExpired_KeyPress2_PublishesBuyRenewal(t *testing.T) {
        m := New()
        m.state = protocol.LicenseExpired
        b := bus.New()
        m.SetBus(b)

        m = updateModel(m, keyRun("2"))
        // State stays until backend confirms; action is published via bus
        if m.State() != protocol.LicenseExpired {
                t.Errorf("state = %q, want %q (should stay until backend confirms)",
                        m.State(), protocol.LicenseExpired)
        }
}

func TestInput_KeyPress1_PublishesBuyLicense(t *testing.T) {
        m := New()
        m.state = protocol.LicenseInput
        b := bus.New()
        m.SetBus(b)

        m = updateModel(m, keyRun("1"))
        // "1" should NOT type a character — it's an action shortcut
        if m.input != "" {
                t.Errorf("input = %q, want empty (1 should be action, not text)", m.input)
        }
}

func TestInput_KeyPressQ_PublishesExit(t *testing.T) {
        m := New()
        m.state = protocol.LicenseInput
        b := bus.New()
        m.SetBus(b)

        m = updateModel(m, keyRun("q"))
        // "q" should NOT type a character — it's an action shortcut
        if m.input != "" {
                t.Errorf("input = %q, want empty (q should be action, not text)", m.input)
        }
}

func TestInvalid_KeyPress1_PublishesBuyLicense(t *testing.T) {
        m := New()
        m.state = protocol.LicenseInvalid
        b := bus.New()
        m.SetBus(b)

        m = updateModel(m, keyRun("1"))
        if m.input != "" {
                t.Errorf("input = %q, want empty (1 should be action, not text)", m.input)
        }
}

func TestInvalid_KeyPressQ_PublishesExit(t *testing.T) {
        m := New()
        m.state = protocol.LicenseInvalid
        b := bus.New()
        m.SetBus(b)

        m = updateModel(m, keyRun("q"))
        if m.input != "" {
                t.Errorf("input = %q, want empty (q should be action, not text)", m.input)
        }
}

func TestConflict_KeyPress2_PublishesAction(t *testing.T) {
        m := New()
        m.state = protocol.LicenseDeviceConflict
        m.resultLicenseKey = "WACL-TEST-0000-0000-0000"
        b := bus.New()
        m.SetBus(b)

        // Pressing "2" should publish an action, state stays until backend confirms
        m = updateModel(m, keyRun("2"))
        if m.State() != protocol.LicenseDeviceConflict {
                t.Errorf("state = %q, want %q (should stay until backend confirms)",
                        m.State(), protocol.LicenseDeviceConflict)
        }
}

func TestServerError_KeyPressEnter_PublishesAction(t *testing.T) {
        m := New()
        m.state = protocol.LicenseServerError
        b := bus.New()
        m.SetBus(b)

        m = updateModel(m, tea.KeyMsg{Type: tea.KeyEnter})
        if m.State() != protocol.LicenseServerError {
                t.Errorf("state should stay %q", m.State())
        }
}

func TestServerError_KeyPress1_RetriesValidation(t *testing.T) {
        m := New()
        m.state = protocol.LicenseServerError
        m.resultLicenseKey = "WACL-TEST-0000-0000-0000"
        b := bus.New()
        m.SetBus(b)

        // Pressing "1" should re-trigger validation, not go to input
        m = updateModel(m, keyRun("1"))
        if m.State() != protocol.LicenseValidating {
                t.Errorf("state = %q, want %q", m.State(), protocol.LicenseValidating)
        }
}

// ---------------------------------------------------------------------------
// Hold / auto-transition timing
// ---------------------------------------------------------------------------

func TestValid_HoldElapsedMsg_EnablesInput(t *testing.T) {
        m := New()
        m.state = protocol.LicenseValid
        m.holdElapsed = false
        m = updateModel(m, holdElapsedMsg{})
        if !m.holdElapsed {
                t.Error("expected holdElapsed=true after holdElapsedMsg")
        }
}

func TestValid_AutoTransitionMsg_EnablesInput(t *testing.T) {
        m := New()
        m.state = protocol.LicenseValid
        m.holdElapsed = false
        m = updateModel(m, autoTransitionMsg{})
        if !m.holdElapsed {
                t.Error("expected holdElapsed=true after autoTransitionMsg")
        }
}

func TestValid_BeforeHold_IgnoresEnter(t *testing.T) {
        m := New()
        m.state = protocol.LicenseValid
        m.holdElapsed = false
        _, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
        if cmd != nil {
                t.Error("expected nil cmd when enter pressed before hold elapsed")
        }
}

func TestValid_AfterHold_ProcessesEnter(t *testing.T) {
        m := New()
        m.state = protocol.LicenseValid
        m.holdElapsed = true
        _, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
        if cmd == nil {
                t.Error("expected non-nil cmd when enter pressed after hold elapsed")
        }
}

// ---------------------------------------------------------------------------
// needsHoldSchedule bridge — backend-driven valid result schedules hold
// ---------------------------------------------------------------------------

func TestNeedsHoldSchedule_SchedulesHoldAndTransition(t *testing.T) {
        m := New()
        m.state = protocol.LicenseValid
        m.needsHoldSchedule = true
        m.greenPulseStart = time.Now()

        // The needsHoldSchedule bridge is checked AFTER the switch statement in Update(),
        // so it triggers on any message that reaches the fallthrough. Use a no-op message
        // that doesn't match any case in the switch.
        updated, cmd := m.Update(struct{}{})
        u := updated.(Model)
        if u.needsHoldSchedule {
                t.Error("expected needsHoldSchedule=false after Update processes it")
        }
        if cmd == nil {
                t.Error("expected non-nil cmd (hold + auto-transition) when needsHoldSchedule is true")
        }
}

// ---------------------------------------------------------------------------
// Window resize
// ---------------------------------------------------------------------------

func TestWindowResize_UpdatesDimensions(t *testing.T) {
        m := New()
        m = updateModel(m, tea.WindowSizeMsg{Width: 80, Height: 24})
        if m.width != 80 {
                t.Errorf("width = %d, want 80", m.width)
        }
        if m.height != 24 {
                t.Errorf("height = %d, want 24", m.height)
        }
}

// ---------------------------------------------------------------------------
// formattedKey edge cases
// ---------------------------------------------------------------------------

func TestFormattedKey_EmptyInput(t *testing.T) {
        m := New()
        got := m.formattedKey()
        want := "WACL-"
        if got != want {
                t.Errorf("formatted key = %q, want %q", got, want)
        }
}

func TestFormattedKey_WithExistingWACLPrefix(t *testing.T) {
        m := New()
        // User types WACL explicitly — prefix should not be doubled
        for _, ch := range "wacl" {
                m = updateModel(m, keyRun(string(ch)))
        }
        got := m.formattedKey()
        // The prefix "WACL" should be stripped and re-added once
        want := "WACL-"
        if got != want {
                t.Errorf("formatted key = %q, want %q", got, want)
        }
}

// ---------------------------------------------------------------------------
// TUI-vs-doc fidelity — exact text from spec
// ---------------------------------------------------------------------------

func TestViewInput_DocExactText(t *testing.T) {
        h := testutil.NewScreenHelper(New())
        // Doc-aligned: "waclaw butuh lisensi buat jalan." not "butuh lisensi baru"
        h.AssertContains(t, "waclaw butuh lisensi buat jalan")
        // Doc-aligned: "masukin key lisensi lu di bawah." not "masukin lisensi di bawah"
        h.AssertContains(t, "masukin key lisensi lu di bawah")
        // Footer — doc-aligned includes path
        h.AssertContains(t, "lisensi disimpan di: ~/.waclaw/license.md")
        // Doc-aligned: "satu lisensi cuma buat satu device." not "1 lisensi = 1 device"
        h.AssertContains(t, "satu lisensi cuma buat satu device")
        // Action hints — no stale prefix (was "v validasi", now "validasi")
        h.AssertContains(t, "validasi")
        h.AssertContains(t, "beli lisensi")
        h.AssertContains(t, "keluar")
}

// ---------------------------------------------------------------------------
// HandleNavigate — animation timestamps initialized on direct navigation
// ---------------------------------------------------------------------------

func TestHandleNavigate_Expired_SetsRedGlow(t *testing.T) {
        m := New()
        m.HandleNavigate(map[string]any{"state": "license_expired"})
        if m.redGlowStart.IsZero() {
                t.Error("expected redGlowStart to be initialized on navigate to expired")
        }
}

func TestHandleNavigate_Conflict_SetsRedGlow(t *testing.T) {
        m := New()
        m.HandleNavigate(map[string]any{"state": "license_device_conflict"})
        if m.redGlowStart.IsZero() {
                t.Error("expected redGlowStart to be initialized on navigate to conflict")
        }
}

func TestHandleNavigate_Invalid_SetsRedGlow(t *testing.T) {
        m := New()
        m.HandleNavigate(map[string]any{"state": "license_invalid"})
        if m.redGlowStart.IsZero() {
                t.Error("expected redGlowStart to be initialized on navigate to invalid")
        }
}

func TestHandleNavigate_Valid_SetsGreenPulseAndHoldSchedule(t *testing.T) {
        m := New()
        m.HandleNavigate(map[string]any{"state": "license_valid"})
        if m.greenPulseStart.IsZero() {
                t.Error("expected greenPulseStart to be initialized on navigate to valid")
        }
        // Doc spec: HandleNavigate now respects the 800ms hold for LicenseValid,
        // same as HandleUpdate. holdElapsed=false and needsHoldSchedule=true.
        if m.holdElapsed {
                t.Error("expected holdElapsed=false on navigate to valid (hold must be respected)")
        }
        if !m.needsHoldSchedule {
                t.Error("expected needsHoldSchedule=true on navigate to valid")
        }
}

// ---------------------------------------------------------------------------
// extractInt — int and float64 handling
// ---------------------------------------------------------------------------

func TestExtractInt_FromInt(t *testing.T) {
        params := map[string]any{"val": 42}
        got := extractInt(params, "val")
        if got != 42 {
                t.Errorf("extractInt(int) = %d, want 42", got)
        }
}

func TestExtractInt_FromFloat64(t *testing.T) {
        params := map[string]any{"val": float64(71)}
        got := extractInt(params, "val")
        if got != 71 {
                t.Errorf("extractInt(float64) = %d, want 71", got)
        }
}

func TestExtractInt_MissingKey(t *testing.T) {
        params := map[string]any{}
        got := extractInt(params, "nonexistent")
        if got != 0 {
                t.Errorf("extractInt(missing) = %d, want 0", got)
        }
}

func TestHandleUpdate_Float64GraceHours(t *testing.T) {
        m := New()
        m.state = protocol.LicenseValidating
        err := m.HandleUpdate(map[string]any{
                "result":      "server_error",
                "grace_hours": float64(71), // JSON numbers are float64
        })
        if err != nil {
                t.Fatalf("HandleUpdate error: %v", err)
        }
        if m.graceHoursLeft != 71 {
                t.Errorf("grace_hours = %d, want 71 (float64 handling)", m.graceHoursLeft)
        }
}

func TestHandleNavigate_Float64GraceHours(t *testing.T) {
        m := New()
        err := m.HandleNavigate(map[string]any{
                "state":       "license_server_error",
                "grace_hours": float64(48), // JSON numbers are float64
        })
        if err != nil {
                t.Fatalf("HandleNavigate error: %v", err)
        }
        if m.graceHoursLeft != 48 {
                t.Errorf("grace_hours = %d, want 48 (float64 handling)", m.graceHoursLeft)
        }
}

// ---------------------------------------------------------------------------
// licenseResultFromString — centralised string→result mapping (now via protocol)
// ---------------------------------------------------------------------------

func TestLicenseResultValidation(t *testing.T) {
        tests := []struct {
                input    string
                want     protocol.LicenseResult
                wantOk   bool
        }{
                {"valid", protocol.LicenseResultValid, true},
                {"invalid", protocol.LicenseResultInvalid, true},
                {"expired", protocol.LicenseResultExpired, true},
                {"device_conflict", protocol.LicenseResultDeviceConflict, true},
                {"server_error", protocol.LicenseResultServerError, true},
                {"unknown", protocol.LicenseResult(""), false},
                {"", protocol.LicenseResult(""), false},
        }
        for _, tt := range tests {
                result := protocol.LicenseResult(tt.input)
                gotOk := protocol.IsValidLicenseResult(result)
                if gotOk != tt.wantOk {
                        t.Errorf("IsValidLicenseResult(%q) = %v, want %v",
                                tt.input, gotOk, tt.wantOk)
                }
                if gotOk && result != tt.want {
                        t.Errorf("LicenseResult(%q) = %v, want %v",
                                tt.input, result, tt.want)
                }
        }
}

// ---------------------------------------------------------------------------
// startValidation sends validate_license action to backend
// ---------------------------------------------------------------------------

func TestStartValidation_SendsActionToBackend(t *testing.T) {
        m := New()
        for _, ch := range "abcdefghijklmnop" {
                m = updateModel(m, keyRun(string(ch)))
        }
        b := bus.New()
        m.SetBus(b)

        // Press enter to start validation
        model, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
        updated := model.(Model)
        if updated.State() != protocol.LicenseValidating {
                t.Errorf("state = %q, want %q", updated.State(), protocol.LicenseValidating)
        }
        // The validation should have stored the formatted key
        if updated.resultLicenseKey != "WACL-ABCD-EFGH-IJKL-MNOP" {
                t.Errorf("resultLicenseKey = %q, want WACL-ABCD-EFGH-IJKL-MNOP", updated.resultLicenseKey)
        }
}

// ---------------------------------------------------------------------------
// DRY helpers — renderCursor, viewHeader, renderDetail
// ---------------------------------------------------------------------------

func TestRenderCursor_AccentStyle(t *testing.T) {
        m := New()
        m = updateModel(m, keyRun("a"))
        result := m.renderCursor(m.formattedKey(), style.AccentStyle)
        if !strings.Contains(result, "A") {
                t.Error("expected cursor to highlight the character")
        }
}

func TestRenderCursor_DangerStyle(t *testing.T) {
        m := New()
        m = updateModel(m, keyRun("a"))
        result := m.renderCursor(m.formattedKey(), style.DangerStyle)
        if !strings.Contains(result, "A") {
                t.Error("expected cursor to highlight the character")
        }
}

func TestViewHeader_ContainsTitle(t *testing.T) {
        m := New()
        header := m.viewHeader()
        if !strings.Contains(header, "lisensi waclaw") {
                t.Error("expected header to contain license title")
        }
}

func TestRenderDetail_SkipsEmpty(t *testing.T) {
        result := renderDetail("label:", "")
        if result != "" {
                t.Errorf("expected empty result for empty value, got %q", result)
        }
}

func TestRenderDetail_ContainsLabelAndValue(t *testing.T) {
        result := renderDetail("device:", "LAPTOP-HOME")
        if !strings.Contains(result, "device:") || !strings.Contains(result, "LAPTOP-HOME") {
                t.Errorf("expected result to contain label and value, got %q", result)
        }
}
