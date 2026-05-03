package onboarding

import (
        "fmt"
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/internal/tui/component"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        tui "github.com/WaClaw-App/waclaw/internal/tui"
        "github.com/WaClaw-App/waclaw/pkg/protocol"

        tea "github.com/charmbracelet/bubbletea"
        "github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Login-specific layout constants
// ---------------------------------------------------------------------------

const (
        // loginQRIndent is the left padding for the QR code display area.
        loginQRIndent = 9

        // loginQRBoxWidth is the inner width of the QR code box frame.
        loginQRBoxWidth = 19

        // loginQRBoxHeight is the inner height of the QR code box frame (excluding borders).
        loginQRBoxHeight = 3

        // loginIndicatorStagger is the per-step stagger delay for the
        // ● ○ ○ sequential animation.
        loginIndicatorStagger = 120 * time.Millisecond

        // loginContactSyncInterval is the tick interval for the contact
        // sync counter animation.
        loginContactSyncInterval = 100 * time.Millisecond
)

// ---------------------------------------------------------------------------
// LoginModel — bubbletea.Model for the Login screen
// ---------------------------------------------------------------------------

// LoginModel is the bubbletea.Model for the Login screen.
//
// It implements all 5 login states from the doc:
//   - login_qr_waiting: Shows QR code, waiting for scan
//   - login_qr_scanned: QR scanned, syncing contacts
//   - login_success: Successfully connected, slot info
//   - login_expired: Session expired, re-scan needed
//   - login_failed: Network error, connection failed
//
// Micro-interactions:
//   - ● ○ ○ animate sequential with breathing component offset
//   - QR dissolve pixel-by-pixel when scan detected
//   - Contact sync counter live: numbers moving
//   - On success: hold 800ms "udah nyambung" → auto transition
//   - On failed: ✗ red, tapi pesan tetap santai
type LoginModel struct {
        base    tui.ScreenBase
        state   protocol.StateID
        width   int
        height  int
        focused bool

        // QR display component
        qr component.QRDisplay

        // Slot tracking
        filledSlots int
        totalSlots  int
        activeSlot  int

        // Contact sync
        contactCount  int
        syncAnimating bool
        syncStart     time.Time

        // Connected phone display
        phoneNumbers []string

        // Success auto-transition
        successStart   time.Time
        autoTransition bool

        // Expired session info
        expiredSlot    int
        activeSlots    int
        lastSessionAgo string

        // Stagger animation for indicators
        staggerStart time.Time

        // Breathing group for ● pulse indicators
        breathing component.BreathingGroup
}

// NewLoginModel creates a Login screen model with default values.
func NewLoginModel() LoginModel {
        base := tui.NewScreenBase(protocol.ScreenLogin)
        qr := component.NewQRDisplay("")
        qr.TotalSlots = 0
        qr.ActiveSlot = 0

        return LoginModel{
                base:           base,
                state:          protocol.LoginQRWaiting,
                qr:             qr,
                filledSlots:    0,
                totalSlots:     0,
                activeSlot:     0,
                phoneNumbers:   []string{},
                autoTransition: false,
                breathing:      component.NewBreathingGroup(3), // up to 3 indicators
        }
}

// ID returns the screen identifier.
func (m LoginModel) ID() protocol.ScreenID { return m.base.ID() }

// SetBus injects the event bus reference.
func (m *LoginModel) SetBus(b *bus.Bus) { m.base.SetBus(b) }

// Bus returns the event bus.
func (m *LoginModel) Bus() *bus.Bus { return m.base.Bus() }

// Focus is called when this screen becomes the active screen.
func (m *LoginModel) Focus() {
        m.focused = true
        now := time.Now()
        m.staggerStart = now
        m.syncStart = now
        m.successStart = now
        m.autoTransition = false
        m.breathing = component.NewBreathingGroup(m.indicatorCount())
}

// Blur is called when this screen is no longer the active screen.
func (m *LoginModel) Blur() { m.focused = false }

// applyLoginParams extracts common login data fields from a params map.
// Shared by HandleNavigate and HandleUpdate to avoid duplication.
func (m *LoginModel) applyLoginParams(params map[string]any) {
        m.filledSlots = toInt(params[protocol.ParamFilledSlots])
        m.totalSlots = toInt(params[protocol.ParamTotalSlots])
        m.activeSlot = toInt(params[protocol.ParamActiveSlot])
        m.contactCount = toInt(params[protocol.ParamContactCount])
        m.expiredSlot = toInt(params[protocol.ParamExpiredSlot])
        m.activeSlots = toInt(params[protocol.ParamActiveSlots])

        if v, ok := params[protocol.ParamLastSessionAgo].(string); ok {
                m.lastSessionAgo = v
        }
        if raw, ok := params[protocol.ParamPhoneNumbers].([]any); ok {
                m.phoneNumbers = toStringSlice(raw)

        }

        // Update QR display component
        m.qr.TotalSlots = m.totalSlots
        m.qr.ActiveSlot = m.activeSlot
}

// transitionQRState updates QR display state based on the current login state.
// Shared by HandleNavigate and HandleUpdate to avoid duplication.
func (m *LoginModel) transitionQRState() {
        switch m.state {
        case protocol.LoginQRScanned:
                if m.qr.State == component.QRWaiting {
                        m.qr.Scan()
                }
        case protocol.LoginSuccess:
                m.qr.State = component.QRSuccess
        case protocol.LoginExpired:
                m.qr.State = component.QRExpired
        case protocol.LoginFailed:
                m.qr.State = component.QRFailed
        }
}

// HandleNavigate processes a "navigate" command from the backend.
func (m *LoginModel) HandleNavigate(params map[string]any) error {
        if stateStr, ok := params[protocol.ParamState].(string); ok {
                m.state = protocol.StateID(stateStr)
        }

        m.applyLoginParams(params)
        m.transitionQRState()

        // Reset animation timers on navigate
        now := time.Now()
        m.staggerStart = now
        m.syncStart = now
        m.successStart = now
        m.autoTransition = false
        m.breathing = component.NewBreathingGroup(m.indicatorCount())

        return nil
}

// HandleUpdate processes an "update" command from the backend.
func (m *LoginModel) HandleUpdate(params map[string]any) error {
        if stateStr, ok := params[protocol.ParamState].(string); ok {
                m.state = protocol.StateID(stateStr)
        }

        m.applyLoginParams(params)
        m.transitionQRState()

        // Reset stagger on state change
        now := time.Now()
        m.staggerStart = now
        m.syncStart = now
        if m.state == protocol.LoginSuccess {
                m.successStart = now
                m.autoTransition = false
        }
        m.breathing = component.NewBreathingGroup(m.indicatorCount())

        return nil
}

// Init implements tea.Model. Starts the animation tick loop.
func (m LoginModel) Init() tea.Cmd {
        return tea.Tick(animationTickInterval, func(t time.Time) tea.Msg {
                return tickMsg(t)
        })
}

// Update implements tea.Model.
func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch msg := msg.(type) {
        case tea.WindowSizeMsg:
                m.width = msg.Width
                m.height = msg.Height
                return m, nil

        case tea.KeyMsg:
                return m.handleKey(msg)

        case tickMsg:
                // Continue the animation tick loop
                cmd := tea.Tick(animationTickInterval, func(t time.Time) tea.Msg {
                        return tickMsg(t)
                })

                now := time.Time(msg)

                // Advance QR dissolve animation
                m.qr.Tick(now)

                // Check auto-transition on success (800ms hold per anim.SuccessHold)
                if m.state == protocol.LoginSuccess && !m.autoTransition {
                        elapsed := time.Since(m.successStart)
                        if elapsed >= anim.SuccessHold {
                                m.autoTransition = true
                                publishAction(m.base.Bus(), protocol.ScreenLogin, protocol.ActionLoginSuccessContinue, nil)
                        }
                }

                return m, cmd
        }

        return m, nil
}

// handleKey routes key events based on the current state.
// F8 FIX: Forward key_press events to the backend so the backend stays
// in sync with user input, per plan requirement.
func (m LoginModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // Forward the key press to the backend for state tracking.
        if b := m.base.Bus(); b != nil {
                b.Publish(bus.KeyPressMsg{
                        Screen: protocol.ScreenLogin,
                        State:  string(m.state),
                        Key:    msg.String(),
                })
        }

        switch m.state {
        case protocol.LoginQRWaiting:
                return m.handleQRWaitingKey(msg)
        case protocol.LoginQRScanned:
                return m.handleQRScannedKey(msg)
        case protocol.LoginSuccess:
                return m.handleSuccessKey(msg)
        case protocol.LoginExpired:
                return m.handleExpiredKey(msg)
        case protocol.LoginFailed:
                return m.handleFailedKey(msg)
        default:
                return m, nil
        }
}

// handleQRWaitingKey handles key events in the qr_waiting state.
func (m LoginModel) handleQRWaitingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "enter":
                // Skip — navigate to niche_select or dashboard
                publishAction(m.base.Bus(), protocol.ScreenLogin, protocol.ActionLoginSkip, nil)
        case "+":
                // Add slot
                publishAction(m.base.Bus(), protocol.ScreenLogin, protocol.ActionLoginAddSlot, nil)
        }
        return m, nil
}

// handleQRScannedKey handles key events in the qr_scanned state.
// F9 FIX: Enter now navigates to niche select (not just an action) so the
// user is never stuck waiting for a backend response that won't come.
func (m LoginModel) handleQRScannedKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "enter":
                // "cukup" — continue to niche select
                publishAction(m.base.Bus(), protocol.ScreenLogin, protocol.ActionLoginEnough, nil)
        case "+":
                // Add another number — let the backend handle state transition
                publishAction(m.base.Bus(), protocol.ScreenLogin, protocol.ActionLoginAddAnother, nil)
        }
        return m, nil
}

// handleSuccessKey handles key events in the login_success state.
// F5 FIX: Both Enter ("cukup, gas") and the auto-transition now go to
// NicheSelect — consistent "continue after login" semantics. The "q later"
// path goes to the dashboard because the user explicitly defers setup.
func (m LoginModel) handleSuccessKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // Cancel auto-transition on any key press
        m.autoTransition = true

        switch msg.String() {
        case "enter":
                // "cukup, gas" — continue to niche setup (same as auto-transition)
                publishAction(m.base.Bus(), protocol.ScreenLogin, protocol.ActionLoginGas, nil)
        case "+":
                // Add another number — let the backend handle state transition
                publishAction(m.base.Bus(), protocol.ScreenLogin, protocol.ActionLoginAddNumber, nil)
        case "q":
                // "nanti" — skip setup for now, go directly to dashboard
                publishAction(m.base.Bus(), protocol.ScreenLogin, protocol.ActionLoginLater, nil)
        }
        return m, nil
}

// handleExpiredKey handles key events in the login_expired state.
// F12 FIX: Previously returned nil for all keys, trapping the user.
// Now supports 'q' to go back to boot and forwards other keys to backend.
func (m LoginModel) handleExpiredKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "q":
                publishAction(m.base.Bus(), protocol.ScreenLogin, protocol.ActionLoginBack, nil)
        default:
                // Forward other key presses to the backend for potential state changes
                publishAction(m.base.Bus(), protocol.ScreenLogin, protocol.ActionLoginExpiredKey, map[string]any{
                        protocol.ParamKey: msg.String(),
                })
        }
        return m, nil
}

// handleFailedKey handles key events in the login_failed state.
func (m LoginModel) handleFailedKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "1":
                // Try again — let the backend handle state transition
                publishAction(m.base.Bus(), protocol.ScreenLogin, protocol.ActionLoginRetry, nil)
        case "2":
                // Change slot
                publishAction(m.base.Bus(), protocol.ScreenLogin, protocol.ActionLoginChangeSlot, nil)
        case "q":
                // Go back
                publishAction(m.base.Bus(), protocol.ScreenLogin, protocol.ActionLoginBack, nil)
        }
        return m, nil
}

// ---------------------------------------------------------------------------
// View — render the current state
// ---------------------------------------------------------------------------

// View implements tea.Model.
func (m LoginModel) View() string {
        switch m.state {
        case protocol.LoginQRWaiting:
                return m.viewQRWaiting()
        case protocol.LoginQRScanned:
                return m.viewQRScanned()
        case protocol.LoginSuccess:
                return m.viewSuccess()
        case protocol.LoginExpired:
                return m.viewExpired()
        case protocol.LoginFailed:
                return m.viewFailed()
        default:
                return m.viewQRWaiting()
        }
}

// ---------------------------------------------------------------------------
// viewQRWaiting renders the login_qr_waiting state.
// ---------------------------------------------------------------------------

func (m LoginModel) viewQRWaiting() string {
        var b strings.Builder

        // Title
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyLoginTitle)))
        b.WriteString(style.Section(style.SectionGap))

        // Instructions
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLoginQRInstruction)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLoginMultiSlotHint)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLoginMoreNumbersSafer)))
        b.WriteString(style.Section(style.SectionGap))

        // QR code in box frame
        b.WriteString(m.renderQRBox())
        b.WriteString(style.Section(style.SectionGap))

        // Waiting text + slot indicator
        waitingText := style.BodyStyle.Render(i18n.T(i18n.KeyLoginWaitingScan))
        slotIndicator := style.DimStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyLoginSlotIndicator), m.activeSlot+1, m.totalSlots),
        )
        b.WriteString(style.Indent(loginQRIndent))
        b.WriteString(waitingText)
        b.WriteString("   ")
        b.WriteString(slotIndicator)
        b.WriteString(style.Section(style.SectionGap))

        // Progress steps with stagger animation
        now := time.Now()
        opacities := m.breathing.Opacities(now)
        idx := 0

        // Step 1: ●  nyambung ke server wa (filled)
        if isLoginStaggerVisible(m.staggerStart, 0) {
                b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(i18n.T(i18n.KeyLoginConnectServer))))
                idx++
        }

        // Step 2: ○  nunggu scan dari hp (empty)
        if isLoginStaggerVisible(m.staggerStart, 1) {
                b.WriteString(style.Indent(contentIndent))
                b.WriteString(style.DimStyle.Render("○"))
                b.WriteString("  ")
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLoginWaitingHPScan)))
                b.WriteString("\n")
                idx++
        }

        // Step 3: ○  sinkron kontak (empty)
        if isLoginStaggerVisible(m.staggerStart, 2) {
                b.WriteString(style.Indent(contentIndent))
                b.WriteString(style.DimStyle.Render("○"))
                b.WriteString("  ")
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLoginSyncContacts)))
                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SectionGap))

        // Footer actions
        footer := style.BodyStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyLoginSlotsFilled), m.filledSlots, m.totalSlots),
        ) + "   " +
                style.ActionStyle.Render("+") + "  " + i18n.T(i18n.KeyLoginAddSlot) +
                "   " + style.ActionStyle.Render("↵") + "  " + i18n.T(i18n.KeyLoginSkip)
        b.WriteString(style.CaptionStyle.Render(footer))

        return b.String()
}

// ---------------------------------------------------------------------------
// viewQRScanned renders the login_qr_scanned state.
// ---------------------------------------------------------------------------

func (m LoginModel) viewQRScanned() string {
        var b strings.Builder

        // Title
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyLoginTitle)))
        b.WriteString(style.Section(style.SectionGap))

        // Scan detected message
        b.WriteString(style.SuccessStyle.Bold(true).Render(i18n.T(i18n.KeyLoginScanDetected)))
        b.WriteString(style.Section(style.SectionGap))

        // Progress steps
        now := time.Now()
        opacities := m.breathing.Opacities(now)
        idx := 0

        // Step 1: ●  nyambung ke server wa (filled)
        b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(i18n.T(i18n.KeyLoginConnectServer))))
        idx++

        // Step 2: ●  scan berhasil (filled)
        b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(i18n.T(i18n.KeyLoginScanSuccess))))
        idx++

        // Step 3: ○  sinkron kontak... 847 (animating)
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.DimStyle.Render("○"))
        b.WriteString("  ")
        syncLabel := fmt.Sprintf(i18n.T(i18n.KeyLoginSyncContactsCount), m.contactCount)
        b.WriteString(style.BodyStyle.Render(syncLabel))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // Footer actions
        footer := style.BodyStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyLoginSlotsFilled), m.filledSlots, m.totalSlots),
        ) + "   " +
                style.MutedStyle.Render(i18n.T(i18n.KeyLoginAddAnother)) +
                "   " + style.ActionStyle.Render("+") + "  " + i18n.T(i18n.KeyLoginYesAdd) +
                "   " + style.ActionStyle.Render("↵") + "  " + i18n.T(i18n.KeyLoginEnoughContinue)
        b.WriteString(style.CaptionStyle.Render(footer))

        return b.String()
}

// ---------------------------------------------------------------------------
// viewSuccess renders the login_success state.
// ---------------------------------------------------------------------------

func (m LoginModel) viewSuccess() string {
        var b strings.Builder

        // Title
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyLoginTitle)))
        b.WriteString(style.Section(style.SectionGap))

        // Progress steps — all filled
        now := time.Now()
        opacities := m.breathing.Opacities(now)
        idx := 0

        // Step 1: ●  nyambung ke server wa
        b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(i18n.T(i18n.KeyLoginConnectServer))))
        idx++

        // Step 2: ●  scan berhasil
        b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(i18n.T(i18n.KeyLoginScanSuccess))))
        idx++

        // Step 3: ●  kontak sinkron (847)
        b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyLoginContactsSynced), m.contactCount),
        )))

        b.WriteString(style.Section(style.SectionGap))

        // Connected slot info — only show if phone numbers are available
        if len(m.phoneNumbers) > 0 {
                for i, phone := range m.phoneNumbers {
                        b.WriteString(style.Indent(contentIndent))
                        b.WriteString(style.SuccessStyle.Render(
                                fmt.Sprintf(i18n.T(i18n.KeyLoginSlotConnected), i+1, phone),
                        ))
                        b.WriteString("\n")
                }
        }

        b.WriteString(style.Section(style.SectionGap))

        // Prompt
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyLoginConnectedNow)))
        b.WriteString("\n")
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLoginMoreSafer)))

        b.WriteString(style.Section(style.SectionGap))

        // Footer actions
        footer := style.ActionStyle.Render("+") + "  " + i18n.T(i18n.KeyLoginAddNumber) +
                "   " + style.ActionStyle.Render("↵") + "  " + i18n.T(i18n.KeyLoginEnoughGas) +
                "   " + style.ActionStyle.Render("q") + "  " + i18n.T(i18n.KeyLoginLater)
        b.WriteString(style.CaptionStyle.Render(footer))

        return b.String()
}

// ---------------------------------------------------------------------------
// viewExpired renders the login_expired state.
// ---------------------------------------------------------------------------

func (m LoginModel) viewExpired() string {
        var b strings.Builder

        // Title
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyLoginTitle)))
        b.WriteString(style.Section(style.SectionGap))

        // Expired message
        b.WriteString(style.WarningStyle.Render(i18n.T(i18n.KeyLoginSessionExpired)))
        b.WriteString(style.Section(style.SectionGap))

        // QR code in box frame (for re-scan)
        b.WriteString(m.renderQRBox())
        b.WriteString(style.Section(style.SectionGap))

        // Progress steps
        now := time.Now()
        opacities := m.breathing.Opacities(now)
        idx := 0

        // Step 1: ●  nyambung ke server wa (filled)
        b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(i18n.T(i18n.KeyLoginConnectServer))))
        idx++

        // Step 2: ○  nunggu scan dari hp (empty)
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.DimStyle.Render("○"))
        b.WriteString("  ")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLoginWaitingHPScan)))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // Slot status
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.BodyStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyLoginSlotExpired), m.expiredSlot, m.activeSlots),
        ))
        b.WriteString("\n")
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLoginExpiredAutoPause)))

        // Last session separator
        if m.lastSessionAgo != "" {
                b.WriteString(style.Section(style.SectionGap))
                b.WriteString(renderLabeledSeparator(
                        fmt.Sprintf(i18n.T(i18n.KeyLoginLastSession), m.lastSessionAgo),
                ))
        }

        return b.String()
}

// ---------------------------------------------------------------------------
// viewFailed renders the login_failed state.
// ---------------------------------------------------------------------------

func (m LoginModel) viewFailed() string {
        var b strings.Builder

        // Title
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyLoginTitle)))
        b.WriteString(style.Section(style.SectionGap))

        // Progress steps
        now := time.Now()
        opacities := m.breathing.Opacities(now)
        idx := 0

        // Step 1: ●  nyambung ke server wa (filled)
        b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(i18n.T(i18n.KeyLoginConnectServer))))
        idx++

        // Step 2: ✗  gagal nyambung (failed, red — santai tone)
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.DangerStyle.Bold(true).Render("✗"))
        b.WriteString("  ")
        b.WriteString(style.DangerStyle.Render(i18n.T(i18n.KeyLoginConnectionFailed)))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // Explanation — santai tone (problem, bukan disaster)
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyLoginSlotFailedNote)))
        b.WriteString("\n")
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLoginWAServerIssue)))
        b.WriteString("\n")
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLoginTryAgainLater)))

        b.WriteString(style.Section(style.SectionGap))

        // Footer actions
        footer := style.ActionStyle.Render("1") + "  " + i18n.T(i18n.KeyLoginTryAgain) +
                "   " + style.ActionStyle.Render("2") + "  " + i18n.T(i18n.KeyLoginChangeSlot) +
                "   " + style.ActionStyle.Render("q") + "  " + i18n.T(i18n.KeyLoginBack)
        b.WriteString(style.CaptionStyle.Render(footer))

        return b.String()
}

// ---------------------------------------------------------------------------
// QR code box rendering
// ---------------------------------------------------------------------------

// renderQRBox renders the QR code inside a box frame as shown in the doc:
//
//      ┌─────────────────┐
//      │                 │
//      │    [QR CODE]    │
//      │                 │
//      └─────────────────┘
func (m LoginModel) renderQRBox() string {
        var b strings.Builder

        // Top border
        b.WriteString(style.Indent(loginQRIndent))
        b.WriteString(style.DimStyle.Render("┌" + strings.Repeat("─", loginQRBoxWidth) + "┐"))
        b.WriteString("\n")

        // Empty line above QR
        b.WriteString(style.Indent(loginQRIndent))
        b.WriteString(style.DimStyle.Render("│" + strings.Repeat(" ", loginQRBoxWidth) + "│"))
        b.WriteString("\n")

        // QR code content — centered within the box
        qrView := m.qr.View()
        qrLines := strings.Split(qrView, "\n")
        for _, line := range qrLines {
                if line == "" {
                        continue
                }
                // Strip ANSI codes to get visible width for centering
                visibleLen := lipgloss.Width(line)
                leftPad := (loginQRBoxWidth - visibleLen) / 2
                if leftPad < 0 {
                        leftPad = 0
                }
                rightPad := loginQRBoxWidth - visibleLen - leftPad
                if rightPad < 0 {
                        rightPad = 0
                }
                b.WriteString(style.Indent(loginQRIndent))
                b.WriteString(style.DimStyle.Render("│"))
                b.WriteString(strings.Repeat(" ", leftPad))
                b.WriteString(line)
                b.WriteString(strings.Repeat(" ", rightPad))
                b.WriteString(style.DimStyle.Render("│"))
                b.WriteString("\n")
        }

        // Empty line below QR
        b.WriteString(style.Indent(loginQRIndent))
        b.WriteString(style.DimStyle.Render("│" + strings.Repeat(" ", loginQRBoxWidth) + "│"))
        b.WriteString("\n")

        // Bottom border
        b.WriteString(style.Indent(loginQRIndent))
        b.WriteString(style.DimStyle.Render("└" + strings.Repeat("─", loginQRBoxWidth) + "┘"))

        return b.String()
}

// ---------------------------------------------------------------------------
// Helper methods
// ---------------------------------------------------------------------------

// indicatorCount returns the number of ● ○ indicator steps for the current state.
func (m LoginModel) indicatorCount() int {
        switch m.state {
        case protocol.LoginQRWaiting:
                return 1 // only 1 filled (connect server)
        case protocol.LoginQRScanned:
                return 2 // 2 filled (connect server + scan success)
        case protocol.LoginSuccess:
                return 3 // all 3 filled
        case protocol.LoginExpired:
                return 1 // only 1 filled
        case protocol.LoginFailed:
                return 1 // only 1 filled, second is ✗
        default:
                return 1
        }
}

// isLoginStaggerVisible returns whether the i-th indicator should be visible
// based on the stagger animation timing.
func isLoginStaggerVisible(start time.Time, index int) bool {
        elapsed := time.Since(start)
        threshold := time.Duration(index) * loginIndicatorStagger
        return elapsed >= threshold
}

// Compile-time interface checks.
var (
        _ tui.Screen        = (*LoginModel)(nil)
        _ protocol.ScreenID = protocol.ScreenLogin
)
