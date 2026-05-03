// Package license implements Screen 18: the license gate.
//
// This is a hard gate — WaClaw will not run without a valid license.
// The screen handles 7 states: input, validating, valid, invalid, expired,
// device_conflict, and server_error.
//
// Doc source: doc/13-screens-license.md
package license

import (
        "fmt"
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/internal/tui/component"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
        "github.com/charmbracelet/lipgloss"
        tea "github.com/charmbracelet/bubbletea"
)

// ---------------------------------------------------------------------------
// Internal messages (tea.Cmd results)
// ---------------------------------------------------------------------------

// holdElapsedMsg fires after the 800ms success hold period, enabling key input.
type holdElapsedMsg struct{}

// autoTransitionMsg fires after the 2s auto-transition delay on valid state.
type autoTransitionMsg struct{}

// validationTickMsg advances the sequential step animation during validation.
type validationTickMsg struct {
        step int // 0, 1, 2
}

// forceDisconnectConfirmedMsg is received when the user confirms the force
// disconnect action. It triggers automatic re-validation with the existing key.
type forceDisconnectConfirmedMsg struct{}

// License result types are now in pkg/protocol as protocol.LicenseResult.
// See LicenseResultValid, LicenseResultInvalid, etc.

// ---------------------------------------------------------------------------
// License screen model
// ---------------------------------------------------------------------------

// Model implements the license gate screen. It satisfies the Screen interface
// defined in the parent tui package via structural typing — no import of the
// parent package is required.
//
// Every color uses style.*. Every timing uses anim.*. No borders, no boxes.
// All display strings use i18n.T() for multi-language support with exact
// doc/13-screens-license.md parity in both locales.
type Model struct {
        // state is the current screen state within the license state machine.
        state protocol.StateID

        // bus holds the event bus reference, injected via SetBus.
        bus *bus.Bus

        // input is the raw user-entered license key characters (no formatting).
        input string

        // cursor is the cursor position within the raw input string.
        cursor int

        // Validation progress tracking.
        validationStep int     // 0 = connecting, 1 = checking, 2 = device check
        validationDone bool    // true once all three steps have ticked through
        validationPct  float64 // 0.0–1.0 progress for the progress bar

        // Result data populated by HandleNavigate / HandleUpdate.
        resultLicenseKey string
        resultDevice     string
        resultExpires    string
        resultExpiredAgo string
        resultOtherDev   string
        resultLastActive string
        graceHoursLeft   int

        // Animation timestamps for visual effects.
        greenPulseStart time.Time
        redGlowStart    time.Time

        // holdElapsed is true after the 800ms success hold period has passed,
        // allowing the user to press enter to continue.
        holdElapsed bool

        // needsHoldSchedule is set by HandleUpdate when a valid result arrives.
        // The next Update() call picks this up and schedules the hold +
        // auto-transition commands. This bridge is needed because HandleUpdate
        // returns error (not tea.Cmd), so async scheduling must go through
        // the bubbletea Update loop.
        needsHoldSchedule bool

        // Window dimensions (set via tea.WindowSizeMsg).
        width  int
        height int

        // focused is true when this screen is the active screen.
        focused bool

        // keyPrefix is the license key prefix from backend (default: protocol.LicenseKeyPrefixV1).
        keyPrefix string
}

// New creates a new license screen Model in the input state.
func New() Model {
        return Model{
                state:     protocol.LicenseInput,
                keyPrefix: protocol.LicenseKeyPrefixV1,
        }
}

// ---------------------------------------------------------------------------
// tea.Model interface
// ---------------------------------------------------------------------------

// Init returns the initial command.
func (m Model) Init() tea.Cmd {
        return nil
}

// Update handles all messages for the license screen.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch msg := msg.(type) {
        case tea.WindowSizeMsg:
                m.width = msg.Width
                m.height = msg.Height
                return m, nil

        case tea.KeyMsg:
                return m.handleKey(msg)

        case validationTickMsg:
                return m.handleValidationTick(msg.step)

        case holdElapsedMsg:
                m.holdElapsed = true
                return m, nil

        case autoTransitionMsg:
                // Auto-transition to next screen after 2s hold on valid state.
                m.holdElapsed = true
                return m, m.publishAction(string(protocol.ActionLicenseContinue), nil)

        case bus.ConfirmResultMsg:
                if msg.ConfirmType == protocol.ConfirmForceDisconnect && msg.Accepted {
                        // Doc spec: "Setelah putuskan: auto-revalidate → license_valid"
                        return m.startValidationFromServerError()
                }
        }

        // Bridge: if HandleUpdate (backend-driven) produced a valid result,
        // schedule the hold + auto-transition here where we CAN return tea.Cmd.
        if m.needsHoldSchedule {
                m.needsHoldSchedule = false
                return m, tea.Batch(
                        holdCmd(anim.SuccessHold),
                        autoTransitionCmd(anim.AutoTransitionDelay),
                )
        }

        return m, nil
}

// View renders the license screen based on the current state.
func (m Model) View() string {
        switch m.state {
        case protocol.LicenseInput:
                return m.viewInput()
        case protocol.LicenseValidating:
                return m.viewValidating()
        case protocol.LicenseValid:
                return m.viewValid()
        case protocol.LicenseInvalid:
                return m.viewInvalid()
        case protocol.LicenseExpired:
                return m.viewExpired()
        case protocol.LicenseDeviceConflict:
                return m.viewDeviceConflict()
        case protocol.LicenseServerError:
                return m.viewServerError()
        default:
                return style.DimStyle.Render(i18n.T(i18n.KeyLicenseUnknownState))
        }
}

// ---------------------------------------------------------------------------
// State accessor
// ---------------------------------------------------------------------------

// State returns the current license state ID. Useful for tests.
func (m Model) State() protocol.StateID { return m.state }

// ---------------------------------------------------------------------------
// Key handling per state
// ---------------------------------------------------------------------------

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch m.state {
        case protocol.LicenseInput, protocol.LicenseInvalid:
                return m.handleInputStateKey(msg)
        case protocol.LicenseValid:
                if m.holdElapsed && msg.String() == "enter" {
                        return m, m.publishAction(string(protocol.ActionLicenseContinue), nil)
                }
                if msg.String() == "q" {
                        return m, m.publishKeyPress("q")
                }
        case protocol.LicenseExpired:
                switch msg.String() {
                case "1":
                        m.state = protocol.LicenseInput
                        m.input = ""
                        m.cursor = 0
                        return m, nil
                case "2":
                        return m, m.publishAction(string(protocol.ActionBuyRenewal), nil)
                case "q":
                        return m, m.publishKeyPress("q")
                }
        case protocol.LicenseDeviceConflict:
                switch msg.String() {
                case "1":
                        m.state = protocol.LicenseInput
                        m.input = ""
                        m.cursor = 0
                        return m, nil
                case "2":
                        // Doc spec: confirmation overlay sebelum execute.
                        // Request the app-level confirmation overlay via bus, instead of
                        // directly publishing the destructive action.
                        return m, m.requestForceDisconnectConfirm()
                case "q":
                        return m, m.publishKeyPress("q")
                }
        case protocol.LicenseServerError:
                switch msg.String() {
                case "enter":
                        return m, m.publishAction(string(protocol.ActionLicenseOfflineContinue), nil)
                case "1":
                        // Doc spec: "1 coba lagi" = retry server connection.
                        // Re-validate with the existing key instead of making user retype.
                        return m.startValidationFromServerError()
                case "q":
                        return m, m.publishKeyPress("q")
                }
        }
        return m, nil
}

// handleInputStateKey handles key events for both LicenseInput and
// LicenseInvalid states. Action shortcuts ("1" for buy, "q" for exit)
// take priority over text input per the doc spec.
func (m Model) handleInputStateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // Action shortcuts take priority over text input
        switch msg.String() {
        case "1":
                return m, m.publishAction(string(protocol.ActionBuyLicense), nil)
        case "q":
                return m, m.publishKeyPress("q")
        }

        // Otherwise, handle as text input
        return m.handleInputKey(msg)
}

// handleInputKey processes key events when the user is typing a license key.
// Auto-formats input with WACL- prefix and hyphen every 4 characters.
// Only accepts alphanumeric characters, auto-uppercases.
func (m Model) handleInputKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.Type {
        case tea.KeyCtrlC:
                return m, tea.Quit

        case tea.KeyEnter:
                if len(m.formattedKey()) >= protocol.LicenseKeyFormattedLen {
                        return m.startValidation()
                }
                return m, nil

        case tea.KeyBackspace:
                if m.cursor > 0 {
                        m.cursor--
                        runes := []rune(m.input)
                        if m.cursor < len(runes) {
                                m.input = string(runes[:m.cursor]) + string(runes[m.cursor+1:])
                        } else {
                                m.input = string(runes[:m.cursor])
                        }
                }
                return m, nil

        case tea.KeyLeft:
                if m.cursor > 0 {
                        m.cursor--
                }
                return m, nil

        case tea.KeyRight:
                if m.cursor < len([]rune(m.input)) {
                        m.cursor++
                }
                return m, nil

        case tea.KeyRunes:
                ch := msg.String()
                if !isAlphanumeric(ch) {
                        return m, nil
                }
                ch = strings.ToUpper(ch)
                runes := []rune(m.input)
                insertPos := m.cursor
                newRunes := make([]rune, 0, len(runes)+1)
                newRunes = append(newRunes, runes[:insertPos]...)
                newRunes = append(newRunes, []rune(ch)...)
                newRunes = append(newRunes, runes[insertPos:]...)
                m.input = string(newRunes)
                m.cursor = insertPos + 1
                return m, nil
        }

        return m, nil
}

// formattedKey returns the license key with auto-formatting applied:
// WACL-XXXX-XXXX-XXXX-XXXX. The user types raw characters; this function
// adds the WACL- prefix and hyphens automatically.
func (m Model) formattedKey() string {
        raw := strings.ToUpper(strings.TrimSpace(m.input))
        raw = strings.ReplaceAll(raw, "-", "")

        // Strip key prefix if user typed it — we'll re-add it
        prefix := m.keyPrefix
        if strings.HasPrefix(raw, prefix) {
                raw = raw[len(prefix):]
        }

        var formatted strings.Builder
        formatted.WriteString(m.keyPrefix + "-")

        for i, ch := range raw {
                if i > 0 && i%protocol.LicenseKeyGroupSize == 0 && formatted.Len() < protocol.LicenseKeyFormattedLen {
                        formatted.WriteByte('-')
                }
                if formatted.Len() < protocol.LicenseKeyFormattedLen {
                        formatted.WriteRune(ch)
                }
        }

        return formatted.String()
}

// ---------------------------------------------------------------------------
// Validation flow — TUI sends action to backend, animates while waiting
// ---------------------------------------------------------------------------

// startValidation transitions to the validating state, starts the sequential
// step animation, AND sends a validate_license action to the backend via the
// bus. This is the correct backend-frontend concern split: the TUI animates
// progress while the backend performs actual validation. The backend responds
// via HandleUpdate with the result.
func (m Model) startValidation() (Model, tea.Cmd) {
        m.state = protocol.LicenseValidating
        m.validationStep = 0
        m.validationDone = false
        m.validationPct = 0
        m.resultLicenseKey = m.formattedKey()

        // Send validate_license action to backend so it can check the key.
        // The backend will respond via HandleUpdate with the result.
        validateCmd := m.publishAction(string(protocol.ActionValidateLicense), map[string]any{
                protocol.ParamLicenseKey: m.resultLicenseKey,
        })

        // Start the 3-step sequential animation (200ms per step, per doc spec).
        animCmd := m.validationStepCmd(0)

        return m, tea.Batch(validateCmd, animCmd)
}

// handleValidationTick advances the progress bar when a step completes.
func (m Model) handleValidationTick(step int) (tea.Model, tea.Cmd) {
        if m.state != protocol.LicenseValidating {
                return m, nil
        }
        m.validationPct = float64(step) / float64(protocol.ValidationStepCount)
        if step < m.validationStep {
                return m, nil
        }
        m.validationStep = step
        if step < protocol.ValidationStepCount {
                return m, m.validationStepCmd(step)
        }
        m.validationDone = true
        return m, nil
}

// applyResult applies a validation result to the model state, setting the
// correct protocol state and initializing animation timestamps.
// Shared by HandleUpdate (backend-driven) and HandleNavigate (direct navigation).
func (m *Model) applyResult(result protocol.LicenseResult) {
        m.validationDone = true
        switch result {
        case protocol.LicenseResultValid:
                m.state = protocol.LicenseValid
                m.greenPulseStart = time.Now()
                m.holdElapsed = false
                m.needsHoldSchedule = true // schedule hold + auto-transition in next Update()
        case protocol.LicenseResultInvalid:
                m.state = protocol.LicenseInvalid
                m.redGlowStart = time.Now()
        case protocol.LicenseResultExpired:
                m.state = protocol.LicenseExpired
                m.redGlowStart = time.Now()
        case protocol.LicenseResultDeviceConflict:
                m.state = protocol.LicenseDeviceConflict
                m.redGlowStart = time.Now()
        case protocol.LicenseResultServerError:
                m.state = protocol.LicenseServerError
        }
}

// validationStepCmd returns a command that fires a validation tick after a
// 200ms delay (matching the 200ms spinner smooth rotation from the doc).
func (m Model) validationStepCmd(currentStep int) tea.Cmd {
        nextStep := currentStep + 1
        return func() tea.Msg {
                <-time.After(anim.ValidationStep)
                return validationTickMsg{step: nextStep}
        }
}

// ---------------------------------------------------------------------------
// View implementations for each state
// ---------------------------------------------------------------------------

// viewHeader renders the common license screen header (title + gap).
// Extracted for DRY — used by all 7 view methods.
func (m Model) viewHeader() string {
        var b strings.Builder
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyLicenseTitle)))
        b.WriteString(style.Section(style.SubSectionGap))
        return b.String()
}

// renderDetail renders a single labeled detail line in caption style.
// Extracted for DRY — used by viewValid, viewExpired, viewDeviceConflict,
// and viewServerError.
func renderDetail(label, value string) string {
        if value == "" {
                return ""
        }
        return style.CaptionStyle.Render(fmt.Sprintf("  %s %s", label, value)) + "\n"
}

// viewInput renders the initial license key entry screen.
// Spec: input field pulse accent, auto-format WACL-XXXX-XXXX-XXXX-XXXX.
func (m Model) viewInput() string {
        var b strings.Builder
        key := m.formattedKey()

        b.WriteString(m.viewHeader())

        // Description — two lines from doc spec
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyLicenseNeedsLicense)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLicenseEnterKeyBelow)))
        b.WriteString(style.Section(style.SectionGap))

        displayKey := m.renderCursor(key, style.AccentStyle)
        b.WriteString(m.renderInputField(displayKey))
        b.WriteString(style.Section(style.SectionGap))

        // Action hints
        b.WriteString(m.renderActionHints(
                []actionHint{
                        {key: "\u21B5", label: i18n.T(i18n.KeyLicenseActionValidate), primary: true},
                        {key: "1", label: i18n.T(i18n.KeyLicenseActionBuyLicense)},
                        {key: "q", label: i18n.T(i18n.KeyLicenseActionExit)},
                },
        ))
        b.WriteString(style.Section(style.SubSectionGap))

        // Footer info — caption style
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyLicenseStoredAt)))
        b.WriteString("\n")
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyLicenseOneDevice)))

        return b.String()
}

// viewValidating renders the 3-step sequential validation animation.
// Spec: ● ○ ○ sequential animate, spinner 200ms, progress bar.
func (m Model) viewValidating() string {
        var b strings.Builder

        b.WriteString(m.viewHeader())

        // Three sequential steps with ● / ○ indicators
        steps := []string{
                i18n.T(i18n.KeyLicenseConnecting),
                i18n.T(i18n.KeyLicenseCheckValidity),
                i18n.T(i18n.KeyLicenseCheckDevice),
        }

        for i, step := range steps {
                if i < m.validationStep {
                        // Completed step — green checkmark
                        b.WriteString(style.SuccessStyle.Render(fmt.Sprintf("  \u25CF %s", step)))
                } else if i == m.validationStep && !m.validationDone {
                        // Current active step — bright text
                        b.WriteString(style.BodyStyle.Render(fmt.Sprintf("  \u25CF %s", step)))
                } else {
                        // Pending step — dim
                        b.WriteString(style.DimStyle.Render(fmt.Sprintf("  \u25CB %s", step)))
                }
                if i < len(steps)-1 {
                        b.WriteString("\n")
                }
        }
        b.WriteString(style.Section(style.SubSectionGap))

        // Progress bar using the shared component — no percentage per doc spec
        bar := component.NewProgressBar(style.LicenseProgressBarWidth).
                SetPercent(m.validationPct)
        bar.ShowPercent = false
        bar.Label = i18n.T(i18n.KeyLicenseValidating)
        b.WriteString(bar.View())

        return b.String()
}

// viewValid renders the success state.
// Spec: ✓ green pulse, hold 800ms, auto-transition 2s, ● ○ ○ all green.
func (m Model) viewValid() string {
        var b strings.Builder

        b.WriteString(m.viewHeader())

        // Green pulse: bold for the first SuccessPulse duration, then steady
        elapsed := time.Since(m.greenPulseStart)
        var titleStyle lipgloss.Style
        if elapsed < anim.SuccessPulse {
                titleStyle = style.SuccessStyle.Bold(true)
        } else {
                titleStyle = style.SuccessStyle
        }
        b.WriteString(titleStyle.Render(fmt.Sprintf("\u2713 %s", i18n.T(i18n.KeyLicenseValidShort))))
        b.WriteString(style.Section(style.SectionGap))

        // All three verification steps completed (green checkmarks).
        // Doc spec: ●  terhubung ke server / ●  lisensi valid / ●  device terdaftar
        // Note: KeyLicenseValid (bullet context) has NO exclamation mark,
        // unlike KeyLicenseValidShort (title context) which has one.
        for _, step := range []string{
                i18n.T(i18n.KeyLicenseConnected),
                i18n.T(i18n.KeyLicenseValid),
                i18n.T(i18n.KeyLicenseValidDevice),
        } {
                b.WriteString(style.SuccessStyle.Render(fmt.Sprintf("  \u25CF %s", step)))
                b.WriteString("\n")
        }
        b.WriteString(style.Section(style.SectionGap))

        // License details — caption style, indented. All labels use i18n for locale parity.
        b.WriteString(renderDetail(i18n.T(i18n.KeyLicenseLabelKey), m.resultLicenseKey))
        b.WriteString(renderDetail(i18n.T(i18n.KeyLicenseLabelDevice), m.resultDevice))
        b.WriteString(renderDetail(i18n.T(i18n.KeyLicenseLabelExpires), m.resultExpires))
        b.WriteString(style.Section(style.SectionGap))

        // Guidance text
        b.WriteString(style.CaptionStyle.Render(fmt.Sprintf("  %s", i18n.T(i18n.KeyLicenseKeySaved))))
        b.WriteString("\n")
        b.WriteString(style.CaptionStyle.Render(fmt.Sprintf("  %s", i18n.T(i18n.KeyLicenseReadyToRun))))
        b.WriteString(style.Section(style.SectionGap))

        // Action hints — only interactive after hold elapsed
        if m.holdElapsed {
                b.WriteString(m.renderActionHints(
                        []actionHint{
                                {key: "\u21B5", label: i18n.T(i18n.KeyLicenseActionContinue), primary: true},
                                {key: "q", label: i18n.T(i18n.KeyLicenseActionExit)},
                        },
                ))
        } else {
                // Dimmed placeholder during hold period
                b.WriteString(style.DimStyle.Render(fmt.Sprintf("  \u21B5  %s    q  %s",
                        i18n.T(i18n.KeyLicenseActionContinue),
                        i18n.T(i18n.KeyLicenseActionExit))))
        }

        return b.String()
}

// viewInvalid renders the invalid key error state.
// Spec: ✗ red flash, auto-focus input, red edge glow 800ms.
// IMPORTANT: No borders, no boxes (P3). The "red glow" is a danger-colored
// accent underline beneath the input using dots, not border line characters.
func (m Model) viewInvalid() string {
        var b strings.Builder

        b.WriteString(m.viewHeader())

        // Error message — danger color with timed bold glow
        b.WriteString(m.redGlowStyle().Render(fmt.Sprintf("\u2717 %s", i18n.T(i18n.KeyLicenseInvalidCheck))))
        b.WriteString(style.Section(style.SubSectionGap))

        // Explanation
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLicenseKeyNotMatch)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLicenseNoTypo)))
        b.WriteString(style.Section(style.SectionGap))

        key := m.formattedKey()
        displayKey := m.renderCursor(key, style.DangerStyle)
        b.WriteString(m.renderInputField(displayKey))

        // Red glow accent line under the input (P3-compliant: dots, not ─ border chars)
        elapsed := time.Since(m.redGlowStart)
        if elapsed < anim.ErrorGlow {
                glowWidth := style.LicenseGlowWidth
                glow := strings.Repeat("\u00B7", glowWidth)
                b.WriteString("\n")
                b.WriteString(style.DangerStyle.Render(glow))
        }
        b.WriteString(style.Section(style.SectionGap))

        // Action hints
        b.WriteString(m.renderActionHints(
                []actionHint{
                        {key: "\u21B5", label: i18n.T(i18n.KeyLicenseActionTryAgain), primary: true},
                        {key: "1", label: i18n.T(i18n.KeyLicenseActionBuyLicense)},
                        {key: "q", label: i18n.T(i18n.KeyLicenseActionExit)},
                },
        ))

        return b.String()
}

// viewExpired renders the expired license state.
// Spec: ✗ red, "data aman" reassurance, "beli perpanjangan" action.
func (m Model) viewExpired() string {
        var b strings.Builder

        b.WriteString(m.viewHeader())

        b.WriteString(m.redGlowStyle().Render(fmt.Sprintf("\u2717 %s", i18n.T(i18n.KeyLicenseExpiredLong))))
        b.WriteString(style.Section(style.SectionGap))

        // License details — all labels use i18n for locale parity
        b.WriteString(renderDetail(i18n.T(i18n.KeyLicenseLabelKey), m.resultLicenseKey))
        b.WriteString(renderDetail(i18n.T(i18n.KeyLicenseLabelDevice), m.resultDevice))

        expiredLine := "  " + i18n.T(i18n.KeyLicenseLabelExpired) + " "
        if m.resultExpires != "" {
                expiredLine += m.resultExpires
        }
        if m.resultExpiredAgo != "" {
                expiredLine += " (" + m.resultExpiredAgo + ")"
        }
        b.WriteString(style.CaptionStyle.Render(expiredLine))
        b.WriteString(style.Section(style.SectionGap))

        // Reassurance — data is safe, workers paused
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLicenseWorkersPaused)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLicenseRenewToContinue)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(m.renderActionHints(
                []actionHint{
                        {key: "1", label: i18n.T(i18n.KeyLicenseActionNewLicense), primary: true},
                        {key: "2", label: i18n.T(i18n.KeyLicenseActionBuyRenewal)},
                        {key: "q", label: i18n.T(i18n.KeyLicenseActionExit)},
                },
        ))

        return b.String()
}

// viewDeviceConflict renders the device conflict state.
// Spec: ✗ red flash, device info, "putuskan device lain" force transfer.
func (m Model) viewDeviceConflict() string {
        var b strings.Builder

        b.WriteString(m.viewHeader())

        b.WriteString(m.redGlowStyle().Render(fmt.Sprintf("\u2717 %s", i18n.T(i18n.KeyLicenseConflictLong))))
        b.WriteString(style.Section(style.SectionGap))

        // Device comparison info — all labels use i18n for locale parity
        b.WriteString(renderDetail(i18n.T(i18n.KeyLicenseLabelKey), m.resultLicenseKey))
        b.WriteString(renderDetail(i18n.T(i18n.KeyLicenseLabelThisDevice), m.resultDevice))
        b.WriteString(renderDetail(i18n.T(i18n.KeyLicenseLabelOtherDevice), m.resultOtherDev))
        b.WriteString(renderDetail(i18n.T(i18n.KeyLicenseLabelLastActive), m.resultLastActive))
        b.WriteString(style.Section(style.SectionGap))

        // Explanation — one license, one device
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLicenseOneDeviceExplain)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLicenseMoveDevice)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(m.renderActionHints(
                []actionHint{
                        {key: "1", label: i18n.T(i18n.KeyLicenseActionNewLicense), primary: true},
                        {key: "2", label: i18n.T(i18n.KeyLicenseActionDisconnect)},
                        {key: "q", label: i18n.T(i18n.KeyLicenseActionExit)},
                },
        ))
        b.WriteString(style.Section(style.SubSectionGap))

        // Warning footnote — explains what "putuskan" does (P3-compliant: em dashes, not ──)
        b.WriteString(style.CaptionStyle.Render(
                fmt.Sprintf("\u2014 %s \u2014", i18n.T(i18n.KeyLicenseForceDisconnectExplain)),
        ))

        return b.String()
}

// viewServerError renders the server connectivity error state.
// Spec: ⚠ amber (warning, not error), offline grace 72h, countdown.
// This is NOT a danger state — the user still has an offline path.
func (m Model) viewServerError() string {
        var b strings.Builder

        b.WriteString(m.viewHeader())

        // Amber warning — not red danger. User can still proceed offline.
        b.WriteString(style.WarningStyle.Render(fmt.Sprintf("\u26A0  %s", i18n.T(i18n.KeyLicenseServerFail))))
        b.WriteString(style.Section(style.SubSectionGap))

        // Status indicators: server reachable (●) vs connection failed (✗)
        b.WriteString(style.SuccessStyle.Render(fmt.Sprintf("  \u25CF  %s", i18n.T(i18n.KeyLicenseServerReachable))))
        b.WriteString("\n")
        b.WriteString(style.DangerStyle.Render(fmt.Sprintf("  \u2717  %s", i18n.T(i18n.KeyLicenseConnectionFailed))))
        b.WriteString(style.Section(style.SectionGap))

        // Offline grace explanation — reassuring, actionable
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLicenseHadValidBefore)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLicenseOfflineGrace)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLicenseMustOnline)))
        b.WriteString(style.Section(style.SectionGap))

        // License info with grace countdown — all labels use i18n for locale parity
        b.WriteString(renderDetail(i18n.T(i18n.KeyLicenseLabelKey), m.resultLicenseKey))
        graceText := fmt.Sprintf(i18n.T(i18n.KeyLicenseGraceRemaining), m.graceHoursLeft)
        b.WriteString(style.WarningStyle.Render(graceText))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(m.renderActionHints(
                []actionHint{
                        {key: "\u21B5", label: i18n.T(i18n.KeyLicenseActionOffline), primary: true},
                        {key: "1", label: i18n.T(i18n.KeyLicenseActionRetry)},
                        {key: "q", label: i18n.T(i18n.KeyLicenseActionExit)},
                },
        ))

        return b.String()
}

// ---------------------------------------------------------------------------
// Visual helpers (private) — DRY consolidated
// ---------------------------------------------------------------------------

// actionHint describes a single keyboard hint rendered in the footer area.
type actionHint struct {
        key     string
        label   string
        primary bool
}

// renderInputField renders the license key input field with consistent styling.
// Used by viewInput (accent cursor) and viewInvalid (danger cursor).
func (m Model) renderInputField(displayKey string) string {
        inputStyle := lipgloss.NewStyle().
                Foreground(style.Text).
                Background(style.BgActive).
                Padding(0, 1).
                Width(style.LicenseInputWidth)
        return lipgloss.NewStyle().Width(style.LicenseInputFrameWidth).Render(inputStyle.Render(displayKey))
}

// renderActionHints renders a row of keyboard action hints separated by spaces.
// Primary hints use accent bold, secondary hints use muted style.
func (m Model) renderActionHints(hints []actionHint) string {
        var parts []string
        for _, h := range hints {
                text := fmt.Sprintf("%s  %s", h.key, h.label)
                if h.primary {
                        parts = append(parts, style.AccentStyle.Bold(true).Render(text))
                } else {
                        parts = append(parts, style.MutedStyle.Render(text))
                }
        }
        return strings.Join(parts, strings.Repeat(" ", style.ActionHintSpacing))
}

// renderCursor renders the formatted key with a cursor in the given style.
// DRY consolidation of the former renderInputCursor and renderInputCursorDanger
// which differed only in cursor style (accent vs danger).
func (m Model) renderCursor(key string, cursorStyle lipgloss.Style) string {
        runes := []rune(key)
        cursorFormatted := m.mapCursorToFormatted()
        if cursorFormatted < len(runes) {
                return string(runes[:cursorFormatted]) +
                        cursorStyle.Render(string(runes[cursorFormatted])) +
                        string(runes[cursorFormatted+1:])
        }
        if len(key) < protocol.LicenseKeyFormattedLen {
                return key + cursorStyle.Render("_")
        }
        return key
}

// mapCursorToFormatted maps the raw input cursor position to the corresponding
// position in the formatted key string (which includes WACL- prefix and hyphens).
func (m Model) mapCursorToFormatted() int {
        rawPos := m.cursor
        prefixLen := len(m.keyPrefix) + 1 // prefix + "-"
        if rawPos <= protocol.LicenseKeyGroupSize {
                return prefixLen + rawPos // prefix- + up to group size chars, no hyphen yet
        }
        // After the first group, each group adds a hyphen
        extra := rawPos - protocol.LicenseKeyGroupSize
        hyphens := extra / protocol.LicenseKeyGroupSize
        return prefixLen + rawPos + hyphens
}

// redGlowStyle returns a danger-colored style with timed bold fade.
// Bold for the first ErrorGlow (800ms), then normal danger.
// Danger is only used for technical problems (P8), not rejection.
func (m Model) redGlowStyle() lipgloss.Style {
        if time.Since(m.redGlowStart) < anim.ErrorGlow {
                return style.DangerStyle.Bold(true)
        }
        return style.DangerStyle
}

// ---------------------------------------------------------------------------
// Shared param extraction (DRY)
// ---------------------------------------------------------------------------

// applyParams extracts common data fields from a params map into the model.
// Shared by HandleNavigate and HandleUpdate to eliminate duplication.
func (m *Model) applyParams(params map[string]any) {
        if key, ok := params[protocol.ParamLicenseKey].(string); ok {
                m.resultLicenseKey = key
        }
        if d, ok := params[protocol.ParamDevice].(string); ok {
                m.resultDevice = d
        }
        if e, ok := params[protocol.ParamExpires].(string); ok {
                m.resultExpires = e
        }
        if a, ok := params[protocol.ParamExpiredAgo].(string); ok {
                m.resultExpiredAgo = a
        }
        if o, ok := params[protocol.ParamOtherDevice].(string); ok {
                m.resultOtherDev = o
        }
        if l, ok := params[protocol.ParamLastActive].(string); ok {
                m.resultLastActive = l
        }
        m.graceHoursLeft = extractInt(params, protocol.ParamGraceHours)
        if kp, ok := params[protocol.ParamKeyPrefix].(string); ok && kp != "" {
                m.keyPrefix = kp
        }
}

// ---------------------------------------------------------------------------
// Bus / backend communication helpers
// ---------------------------------------------------------------------------

// publishKeyPress sends a key_press event to the backend via the bus.
func (m Model) publishKeyPress(key string) tea.Cmd {
        b := m.bus // capture bus reference for closure
        return func() tea.Msg {
                if b != nil {
                        b.Publish(bus.KeyPressMsg{
                                Key:    key,
                                Screen: protocol.ScreenLicense,
                        })
                }
                return nil
        }
}

// publishAction sends an action event to the backend via the bus.
func (m Model) publishAction(action string, params map[string]any) tea.Cmd {
        b := m.bus // capture bus reference for closure
        return func() tea.Msg {
                if b != nil {
                        b.Publish(bus.ActionMsg{
                                Action: action,
                                Screen: protocol.ScreenLicense,
                                Params: params,
                        })
                }
                return nil
        }
}

// requestForceDisconnectConfirm publishes a ShowConfirmMsg on the bus,
// requesting the app-level confirmation overlay before executing the
// destructive force_disconnect_device action.
// Doc spec: "ini bakal logout waclaw di device lain" confirmation overlay sebelum execute.
func (m Model) requestForceDisconnectConfirm() tea.Cmd {
        b := m.bus
        data := map[string]any{
                protocol.ParamLicenseKey: m.resultLicenseKey,
        }
        return func() tea.Msg {
                if b != nil {
                        b.Publish(bus.ShowConfirmMsg{
                                ConfirmType: protocol.ConfirmForceDisconnect,
                                Data:        data,
                        })
                }
                return nil
        }
}

// startValidationFromServerError re-triggers validation with the existing key
// when the user presses "1 coba lagi" from the server error state.
// Doc spec: "1 coba lagi" = retry connection, not re-enter key.
func (m Model) startValidationFromServerError() (Model, tea.Cmd) {
        m.state = protocol.LicenseValidating
        m.validationStep = 0
        m.validationDone = false
        m.validationPct = 0

        // Re-validate with the existing key
        validateCmd := m.publishAction(string(protocol.ActionValidateLicense), map[string]any{
                protocol.ParamLicenseKey: m.resultLicenseKey,
        })
        animCmd := m.validationStepCmd(0)
        return m, tea.Batch(validateCmd, animCmd)
}

// ---------------------------------------------------------------------------
// Screen interface compliance
//
// These methods satisfy tui.Screen via structural typing.
// The license package does NOT import the parent tui package (no circular dep).
// ---------------------------------------------------------------------------

// ID returns the screen identifier: ScreenLicense.
func (m Model) ID() protocol.ScreenID {
        return protocol.ScreenLicense
}

// SetBus injects the event bus reference for publishing key/action events
// to the backend.
func (m *Model) SetBus(b *bus.Bus) {
        m.bus = b
}

// HandleNavigate processes a "navigate" command from the backend.
// Supported params:
//
//      "state"         — force a specific license state (string)
//      "license_key"   — pre-fill the license key (string)
//      "device"        — device name (string)
//      "expires"       — expiration date display string (string)
//      "expired_ago"   — human-readable "how long ago expired" (string)
//      "other_device"  — conflicting device name (string)
//      "last_active"   — last active time on other device (string)
//      "grace_hours"   — offline grace hours remaining (int)
func (m *Model) HandleNavigate(params map[string]any) error {
        if s, ok := params[protocol.ParamState].(string); ok {
                m.state = protocol.StateID(s)
        }
        if key, ok := params[protocol.ParamLicenseKey].(string); ok {
                m.resultLicenseKey = key
                m.input = strings.ReplaceAll(strings.TrimPrefix(key, m.keyPrefix+"-"), "-", "")
        }
        // Apply remaining data fields via shared helper
        m.applyParams(params)

        // Initialize animation timestamps for states that use timed effects.
        // Without this, direct navigation skips the glow/pulse.
        switch m.state {
        case protocol.LicenseValid:
                m.greenPulseStart = time.Now()
                m.holdElapsed = false
                m.needsHoldSchedule = true // respect the 800ms hold, same as HandleUpdate
        case protocol.LicenseInvalid, protocol.LicenseExpired, protocol.LicenseDeviceConflict:
                m.redGlowStart = time.Now()
        }

        return nil
}

// HandleUpdate processes an "update" command from the backend.
// When the backend sends validation results while we're in the validating
// state, it transitions to the appropriate result state with correct
// animation scheduling (hold + auto-transition for valid results).
//
// Supported params:
//
//      "result"        — "valid", "invalid", "expired", "device_conflict", "server_error"
//      "license_key"   — the formatted license key (string)
//      "device"        — device name (string)
//      "expires"       — expiration date (string)
//      "expired_ago"   — how long ago expired (string)
//      "other_device"  — other device name (string)
//      "last_active"   — last active time (string)
//      "grace_hours"   — offline grace hours (int)
func (m *Model) HandleUpdate(params map[string]any) error {
        // Extract and apply all data fields via shared helper
        m.applyParams(params)

        // If we're validating, apply the result directly
        if m.state == protocol.LicenseValidating {
                if r, ok := params[protocol.ParamResult].(string); ok {
                        result := protocol.LicenseResult(r)
                        if protocol.IsValidLicenseResult(result) {
                                m.applyResult(result)
                        }
                }
        }

        return nil
}

// Focus is called when this screen becomes the active screen.
func (m *Model) Focus() {
        m.focused = true
}

// Blur is called when this screen is no longer the active screen.
func (m *Model) Blur() {
        m.focused = false
}

// ---------------------------------------------------------------------------
// Commands (tea.Cmd factories)
// ---------------------------------------------------------------------------

// holdCmd fires a holdElapsedMsg after the given duration.
func holdCmd(d time.Duration) tea.Cmd {
        return func() tea.Msg {
                <-time.After(d)
                return holdElapsedMsg{}
        }
}

// autoTransitionCmd fires an autoTransitionMsg after the given duration.
func autoTransitionCmd(d time.Duration) tea.Cmd {
        return func() tea.Msg {
                <-time.After(d)
                return autoTransitionMsg{}
        }
}

// ---------------------------------------------------------------------------
// Utility
// ---------------------------------------------------------------------------

// extractInt extracts an integer from a params map, handling both int and
// float64 (JSON numbers unmarshal as float64 in map[string]any).
func extractInt(params map[string]any, key string) int {
        if v, ok := params[key]; ok {
                switch n := v.(type) {
                case int:
                        return n
                case float64:
                        return int(n)
                }
        }
        return 0
}

// isAlphanumeric reports whether the string contains only ASCII letters or digits.
func isAlphanumeric(s string) bool {
        for _, r := range s {
                if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
                        return false
                }
        }
        return true
}
