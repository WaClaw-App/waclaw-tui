// Package onboarding implements the Boot and Login screens — the first
// impression and WhatsApp authentication flow for the WaClaw TUI.
//
// Boot screen states per doc/01-screens-onboarding-boot-login.md:
//   - boot_first_time:           First-time user sees the logo + 3-step menu
//   - boot_returning:            Returning user sees army report
//   - boot_returning_response:   Returning + new responses waiting
//   - boot_returning_error:      Returning + WA disconnected
//   - boot_returning_config_error: Returning + config error detected
//   - boot_returning_license_expired: Returning + license expired
//   - boot_returning_device_conflict: Returning + device conflict
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
// Boot-specific layout constants
// ---------------------------------------------------------------------------

const (
        // autoTransitionDelay is the wait time before auto-navigating to the
        // dashboard for returning users (3s per doc spec).
        autoTransitionDelay = 3 * time.Second
)

// ---------------------------------------------------------------------------
// Startup phase timing — doc/20-startup-and-session.md
//
// The returning-user boot screen reveals content progressively at specific
// millisecond offsets. Each phase adds a new element to the display,
// creating the visual effect of the system progressively coming online.
// Total target: 1300ms until "ready. cursor blinks."
// ---------------------------------------------------------------------------

const (
        // startupPhaseTagline is when the tagline fades in after the logo.
        startupPhaseTagline = 80 * time.Millisecond

        // startupPhaseSystemCheck is when the system check line appears.
        startupPhaseSystemCheck = 200 * time.Millisecond

        // startupPhaseLicenseCheck is when the license check result appears.
        startupPhaseLicenseCheck = 300 * time.Millisecond

        // startupPhaseConfigValidation is when the config validation line appears.
        startupPhaseConfigValidation = 400 * time.Millisecond

        // startupPhaseValidationResult is when the validation result appears.
        startupPhaseValidationResult = 500 * time.Millisecond

        // startupPhaseStatus is when the status indicators (●) appear.
        startupPhaseStatus = 700 * time.Millisecond

        // startupPhaseAutopilot is when the "auto-pilot: ON" line appears.
        startupPhaseAutopilot = 800 * time.Millisecond

        // startupPhaseArmyMarch is when the army marching animation starts.
        startupPhaseArmyMarch = 900 * time.Millisecond

        // startupPhaseDashboardReady is when the dashboard-ready indicator appears.
        startupPhaseDashboardReady = 1100 * time.Millisecond
)

// ---------------------------------------------------------------------------
// ASCII Logo — exact glyphs from doc/01-screens-onboarding-boot-login.md
// ---------------------------------------------------------------------------

// waclawLogo holds the WACLAW ASCII art lines in order.
var waclawLogo = []string{
        "  ▄▄▄                 ▄   ▄▄▄▄ ▄▄                 ",
        " █▀██  ██  ██▀▀       ▀██████▀  ██                ",
        "   ██  ██  ██           ██      ██                ",
        "   ██  ██  ██ ▄▀▀█▄     ██      ██ ▄▀▀█▄▀█▄ █▄ ██▀",
        "   ██▄ ██▄ ██ ▄█▀██     ██      ██ ▄█▀██ ██▄██▄██ ",
        "   ▀████▀███▀▄▀█▄██     ▀█████ ▄██▄▀█▄██  ▀██▀██▀ ",
}

// ---------------------------------------------------------------------------
// MarchingWorker — army march animation data
// ---------------------------------------------------------------------------

// MarchingWorker represents a single worker row in the army march animation.
type MarchingWorker struct {
        // Name is the niche identifier shown after the marching arrows.
        Name string

        // ArrowCount is the number of ▸ arrows for this row (based on position).
        ArrowCount int
}

// ---------------------------------------------------------------------------
// BootModel — bubbletea.Model for the Boot screen
// ---------------------------------------------------------------------------

// BootModel is the bubbletea.Model for the Boot screen.
//
// It implements all 7 boot states from the doc, including:
//   - Logo reveal animation (8ms/char)
//   - Menu stagger animation (120ms/item)
//   - Army marching animation (▸▸▸ → ● morph, 600ms total, 80ms stagger)
//   - Breathing pulse for ● indicators
//   - Attention flash for "response baru!" amber 2x
//   - Red ✗ flash for error variants
//   - Auto-transition to dashboard after 3s for returning users
type BootModel struct {
        base    tui.ScreenBase
        state   protocol.StateID
        width   int
        height  int
        focused bool

        // Animation timing
        logoRevealStart time.Time
        staggerStart    time.Time
        marchStart      time.Time

        // Breathing group for ● pulse indicators
        breathing component.BreathingGroup

        // Attention flash tracking (for "response baru!" amber 2x)
        attentionStart time.Time

        // Error flash tracking (for ✗ red flash)
        errorFlashStart time.Time

        // Auto-transition timer for returning users
        autoTransitionStart time.Time
        autoTransitionFired bool

        // Boot data — populated from backend via HandleNavigate/HandleUpdate
        waCount      int // number of WA connections
        nicheCount   int // active niches
        workerCount  int // running workers
        leadsCount   int // leads in database
        responseCount int // new responses (for response variant)
        okNicheCount int // niches that are OK (for config_error variant)
        errorNicheCount int // niches with errors (for config_error variant)
        errorNiche   string // niche name with config error
        nicheAlreadySet bool // whether niche config already exists (first_time variant)
        deviceName   string // other device name (for device_conflict variant)
        lastActive   string // other device last active time (for device_conflict variant)

        // Marching workers for army march animation
        workers []MarchingWorker
}

// NewBootModel creates a Boot screen model with default values.
func NewBootModel() BootModel {
        base := tui.NewScreenBase(protocol.ScreenBoot)
        return BootModel{
                base:           base,
                state:          protocol.BootFirstTime,
                breathing:      component.NewBreathingGroup(5), // up to 5 ● indicators
                autoTransitionFired: false,
        }
}

// ID returns the screen identifier.
func (m BootModel) ID() protocol.ScreenID { return m.base.ID() }

// SetBus injects the event bus reference.
func (m *BootModel) SetBus(b *bus.Bus) { m.base.SetBus(b) }

// Bus returns the event bus.
func (m *BootModel) Bus() *bus.Bus { return m.base.Bus() }

// Focus is called when this screen becomes the active screen.
func (m *BootModel) Focus() {
        m.focused = true
        now := time.Now()
        m.logoRevealStart = now
        m.staggerStart = now
        // For returning users, the army march starts at the phased offset
        // (t+900ms per doc/20). For first-time users, it starts immediately.
        if m.isReturningState() {
                m.marchStart = now.Add(startupPhaseArmyMarch)
        } else {
                m.marchStart = now
        }
        m.attentionStart = now
        m.errorFlashStart = now
        m.autoTransitionStart = now
        m.autoTransitionFired = false

        // Reset breathing group based on current indicator count
        indicatorCount := m.indicatorCount()
        if indicatorCount > 0 {
                m.breathing = component.NewBreathingGroup(indicatorCount)
        }
}

// Blur is called when this screen is no longer the active screen.
func (m *BootModel) Blur() { m.focused = false }

// applyBootParams extracts common boot data fields from a params map.
// Shared by HandleNavigate and HandleUpdate to avoid duplication.
func (m *BootModel) applyBootParams(params map[string]any) {
        m.waCount = toInt(params[protocol.ParamWACount])
        m.nicheCount = toInt(params[protocol.ParamNicheCount])
        m.workerCount = toInt(params[protocol.ParamWorkerCount])
        m.leadsCount = toInt(params[protocol.ParamLeadsCount])
        m.responseCount = toInt(params[protocol.ParamResponseCount])
        m.okNicheCount = toInt(params[protocol.ParamOkNicheCount])

        m.errorNicheCount = toInt(params[protocol.ParamErrorNicheCount])

        if v, ok := params[protocol.ParamErrorNiche].(string); ok {
                m.errorNiche = v
        }
        if v, ok := params[protocol.ParamDeviceName].(string); ok {
                m.deviceName = v
        }
        if v, ok := params[protocol.ParamLastActive].(string); ok {
                m.lastActive = v
        }
}

// HandleNavigate processes a "navigate" command from the backend.
func (m *BootModel) HandleNavigate(params map[string]any) error {
        if stateStr, ok := params[protocol.ParamState].(string); ok {
                m.state = protocol.StateID(stateStr)
        }

        m.applyBootParams(params)

        if v, ok := params[protocol.ParamNicheAlreadySet].(bool); ok {
                m.nicheAlreadySet = v
        }

        if raw, ok := params[protocol.ParamWorkers].([]any); ok {
                m.workers = parseMarchingWorkers(raw)
        }

        // Reset animation timers on navigate
        now := time.Now()
        m.logoRevealStart = now
        m.staggerStart = now
        // For returning users, the army march starts at the phased offset
        // (t+900ms per doc/20). For first-time users, it starts immediately.
        if m.isReturningState() {
                m.marchStart = now.Add(startupPhaseArmyMarch)
        } else {
                m.marchStart = now
        }
        m.autoTransitionStart = now
        m.autoTransitionFired = false

        indicatorCount := m.indicatorCount()
        if indicatorCount > 0 {
                m.breathing = component.NewBreathingGroup(indicatorCount)
        }

        return nil
}

// HandleUpdate processes an "update" command from the backend.
func (m *BootModel) HandleUpdate(params map[string]any) error {
        if stateStr, ok := params[protocol.ParamState].(string); ok {
                m.state = protocol.StateID(stateStr)
        }
        m.applyBootParams(params)

        return nil
}

// Init implements tea.Model. Starts the animation tick loop.
func (m BootModel) Init() tea.Cmd {
        return tea.Tick(animationTickInterval, func(t time.Time) tea.Msg {
                return tickMsg(t)
        })
}

// Update implements tea.Model.
func (m BootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

                // Check auto-transition for returning users
                if m.isReturningState() && !m.autoTransitionFired {
                        elapsed := time.Since(m.autoTransitionStart)
                        if elapsed >= autoTransitionDelay {
                                m.autoTransitionFired = true
                                // Navigate to monitor dashboard
                                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootDashboard, nil)
                        }
                }

                return m, cmd
        }

        return m, nil
}

// handleKey routes key events based on the current state.
// F8 FIX: Forward key_press events to the backend in addition to
// the action message. This ensures the backend stays in sync with
// user input per plan requirement: "Key events are forwarded to the
// backend via the App's convenience methods."
func (m BootModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // Cancel auto-transition on any key press for returning users
        if m.isReturningState() {
                m.autoTransitionFired = true
        }

        // Forward the key press to the backend for state tracking.
        if b := m.base.Bus(); b != nil {
                b.Publish(bus.KeyPressMsg{
                        Screen: protocol.ScreenBoot,
                        State:  string(m.state),
                        Key:    msg.String(),
                })
        }

        switch m.state {
        case protocol.BootFirstTime:
                return m.handleFirstTimeKey(msg)
        case protocol.BootReturning:
                return m.handleReturningKey(msg)
        case protocol.BootReturningResponse:
                return m.handleResponseKey(msg)
        case protocol.BootReturningError:
                return m.handleErrorKey(msg)
        case protocol.BootReturningConfigError:
                return m.handleConfigErrorKey(msg)
        case protocol.BootReturningLicenseExpired:
                return m.handleLicenseExpiredKey(msg)
        case protocol.BootReturningDeviceConflict:
                return m.handleDeviceConflictKey(msg)
        default:
                return m, nil
        }
}

// handleFirstTimeKey handles key events in the first_time state.
func (m BootModel) handleFirstTimeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "1":
                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootLogin, nil)
        case "2":
                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootNiche, nil)
        case "3":
                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootGas, nil)
        }
        return m, nil
}

// handleReturningKey handles key events in the returning state.
// Any key navigates to the dashboard.
func (m BootModel) handleReturningKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootDashboard, nil)
        return m, nil
}

// handleResponseKey handles key events in the returning + response state.
func (m BootModel) handleResponseKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "enter":
                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootViewResponses, nil)
        default:
                // Any other key → dashboard
                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootDashboard, nil)
        }
        return m, nil
}

// handleErrorKey handles key events in the returning + WA disconnect state.
func (m BootModel) handleErrorKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "1":
                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootRelogin, nil)
        }
        return m, nil
}

// handleConfigErrorKey handles key events in the returning + config error state.
func (m BootModel) handleConfigErrorKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "v":
                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootViewError, map[string]any{
                        protocol.ParamNiche: m.errorNiche,
                })
        case "enter":
                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootDashboard, nil)
        case "q":
                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootExit, nil)
        }
        return m, nil
}

// handleLicenseExpiredKey handles key events in the returning + license expired state.
func (m BootModel) handleLicenseExpiredKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "1":
                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootEnterLicense, nil)
        case "2":
                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootBuyLicense, nil)
        case "q":
                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootExit, nil)
        }
        return m, nil
}

// handleDeviceConflictKey handles key events in the returning + device conflict state.
func (m BootModel) handleDeviceConflictKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "1":
                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootEnterLicense, nil)
        case "2":
                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootDisconnectOther, nil)
        case "q":
                publishAction(m.base.Bus(), protocol.ScreenBoot, protocol.ActionBootExit, nil)
        }
        return m, nil
}

// ---------------------------------------------------------------------------
// View — render the current state
// ---------------------------------------------------------------------------

// View implements tea.Model.
func (m BootModel) View() string {
        switch m.state {
        case protocol.BootFirstTime:
                return m.viewFirstTime()
        case protocol.BootReturning:
                return m.viewReturning()
        case protocol.BootReturningResponse:
                return m.viewReturningResponse()
        case protocol.BootReturningError:
                return m.viewReturningError()
        case protocol.BootReturningConfigError:
                return m.viewReturningConfigError()
        case protocol.BootReturningLicenseExpired:
                return m.viewReturningLicenseExpired()
        case protocol.BootReturningDeviceConflict:
                return m.viewReturningDeviceConflict()
        default:
                return m.viewFirstTime()
        }
}

// ---------------------------------------------------------------------------
// viewFirstTime renders the boot_first_time state.
// ---------------------------------------------------------------------------

func (m BootModel) viewFirstTime() string {
        var b strings.Builder

        // Logo with character-by-character reveal
        b.WriteString(m.renderLogoReveal())
        b.WriteString(style.Section(style.SectionGap))

        // Tagline
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyBootFirstTimeTagline)))
        b.WriteString(style.Section(style.SectionGap))

        // Separator with label — use the separator-specific i18n key
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(renderLabeledSeparator(i18n.T(i18n.KeyBootFirstTimeSeparator)))
        b.WriteString(style.Section(style.SectionGap))

        // Menu items with stagger
        menuItems := []struct {
                key      string
                labelKey string
                descKey  string
        }{
                {"1", i18n.KeyBootMenuLogin, i18n.KeyBootMenuLoginDesc},
                {"2", i18n.KeyBootMenuNiche, i18n.KeyBootMenuNicheDesc},
                {"3", i18n.KeyBootMenuGas, i18n.KeyBootMenuGasDesc},
        }

        for i, item := range menuItems {
                if !isMenuStaggerVisible(m.staggerStart, i) {
                        continue
                }

                b.WriteString(style.Indent(contentIndent))

                // Key number
                b.WriteString(style.ActionStyle.Render(item.key))
                b.WriteString("  ")

                // Label
                label := i18n.T(item.labelKey)
                b.WriteString(style.PrimaryStyle.Render(label))

                // Description or "already set" indicator
                if i == 1 && m.nicheAlreadySet {
                        b.WriteString("  ")
                        b.WriteString(style.SuccessStyle.Render(i18n.T(i18n.KeyBootNicheAlreadyConfigured)))
                } else {
                        // Dynamic padding to column-align descriptions
                        labelWidth := lipgloss.Width(label)
                        padLen := menuDescColumn - len(item.key) - 2 - labelWidth
                        if padLen < 2 {
                                padLen = 2
                        }
                        b.WriteString(strings.Repeat(" ", padLen))
                        b.WriteString(style.MutedStyle.Render(i18n.T(item.descKey)))
                }

                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SectionGap))

        // Bottom separator
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(renderSeparator())
        b.WriteString(style.Section(style.SectionGap))

        // Steps summary
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyBootStepsSummary)))

        return b.String()
}

// ---------------------------------------------------------------------------
// viewReturning renders the boot_returning state with phased startup
// sequence per doc/20-startup-and-session.md.
//
// The returning-user boot screen progressively reveals content at specific
// millisecond offsets, creating the visual effect of the system coming
// online step by step:
//
//      t +0ms    logo render per karakter
//      t +80ms   tagline fade in
//      t +200ms  system check (wa, config, db, lisensi)
//      t +300ms  license check result
//      t +400ms  config validation
//      t +500ms  validation result
//      t +700ms  status indicators (●)
//      t +800ms  auto-pilot: ON
//      t +900ms  army marching
//      t +1100ms  dashboard fade in
//      t +1300ms  ready. cursor blinks.
//
// Before 1300ms the user sees a progressive reveal. After 1300ms, the
// full screen is shown and the auto-transition timer (3s) begins counting.
// ---------------------------------------------------------------------------
func (m BootModel) viewReturning() string {
        var b strings.Builder
        elapsed := time.Since(m.logoRevealStart)

        // Phase 0: Logo (t+0ms) — always visible, character-by-character reveal.
        b.WriteString(m.renderLogoReveal())

        // Phase 1: Tagline (t+80ms).
        if elapsed >= startupPhaseTagline {
                b.WriteString(style.Section(style.SectionGap))
                b.WriteString(style.Indent(contentIndent))
                b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyBootFirstTimeTagline)))
        }

        // Phase 2: System check (t+200ms) — "system check (wa, config, db, lisensi)".
        if elapsed >= startupPhaseSystemCheck {
                b.WriteString(style.Section(style.SectionGap))
                b.WriteString(style.Indent(contentIndent))
                b.WriteString(style.DimStyle.Render(i18n.T(i18n.KeyBootSystemCheck)))
        }

        // Phase 3: License check result (t+300ms).
        if elapsed >= startupPhaseLicenseCheck {
                b.WriteString(style.Indent(contentIndent))
                b.WriteString(style.SuccessStyle.Render("✓"))
                b.WriteString(" ")
                b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyBootLicenseValid)))
                b.WriteString("\n")
        }

        // Phase 4: Config validation (t+400ms).
        if elapsed >= startupPhaseConfigValidation {
                b.WriteString(style.Indent(contentIndent))
                b.WriteString(style.SuccessStyle.Render("✓"))
                b.WriteString(" ")
                b.WriteString(style.BodyStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootConfigValidCount), m.okNicheCount)))
                b.WriteString("\n")
        }

        // Phase 5: Validation result (t+500ms).
        if elapsed >= startupPhaseValidationResult {
                b.WriteString(style.Indent(contentIndent))
                b.WriteString(style.DimStyle.Render(i18n.T(i18n.KeyBootValidationResult)))
                b.WriteString("\n")
        }

        // Phase 6: Status indicators (t+700ms) — ● breathing dots.
        if elapsed >= startupPhaseStatus {
                b.WriteString(style.Section(style.SectionGap))
                now := time.Now()
                opacities := m.breathing.Opacities(now)
                idx := 0

                b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootWAConnectedCount), m.waCount))))
                idx++
                b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootNicheWorkerSummary), m.nicheCount, m.workerCount))))
                idx++
                b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootLeadsCount), m.leadsCount))))
        }

        // Phase 7: Auto-pilot ON (t+800ms).
        if elapsed >= startupPhaseAutopilot {
                b.WriteString(style.Section(style.SectionGap))
                b.WriteString(style.Indent(contentIndent))
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyBootAutopilotActive)))
                b.WriteString("\n")
                b.WriteString(style.Indent(contentIndent))
                b.WriteString(style.MutedStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootWARotating), m.waCount)))
        }

        // Phase 8: Army marching (t+900ms) — start the march animation.
        if elapsed >= startupPhaseArmyMarch {
                b.WriteString(style.Section(style.SectionGap))
                b.WriteString(style.Indent(0))
                b.WriteString(renderLabeledSeparator(i18n.T(i18n.KeyBootArmyReport)))
                b.WriteString(style.Section(style.SectionGap))
                b.WriteString(m.renderArmyMarch())
        }

        // Phase 9: Dashboard ready (t+1100ms).
        if elapsed >= startupPhaseDashboardReady {
                b.WriteString(style.Indent(contentIndent))
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyBootPressAnyDashboard)))
                b.WriteString("\n")
                b.WriteString(style.Indent(contentIndent))
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyBootArmyWorking)))

                // Bottom separator
                b.WriteString(style.Section(style.SectionGap))
                b.WriteString(style.Indent(0))
                b.WriteString(renderSeparator())
        }

        return b.String()
}

// ---------------------------------------------------------------------------
// viewReturningResponse renders the boot_returning_response variant.
// ---------------------------------------------------------------------------

func (m BootModel) viewReturningResponse() string {
        var b strings.Builder

        // Logo
        b.WriteString(m.renderLogoReveal())
        b.WriteString(style.Section(style.SectionGap))

        // Army report separator
        b.WriteString(style.Indent(0))
        b.WriteString(renderLabeledSeparator(i18n.T(i18n.KeyBootArmyReport)))
        b.WriteString(style.Section(style.SectionGap))

        // Army marching animation
        b.WriteString(m.renderArmyMarch())
        b.WriteString(style.Section(style.SectionGap))

        // Status indicators with breathing
        now := time.Now()
        opacities := m.breathing.Opacities(now)
        idx := 0

        b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootWAConnectedCount), m.waCount))))
        idx++
        b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootNicheWorkerSummary), m.nicheCount, m.workerCount))))
        idx++
        b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootLeadsCount), m.leadsCount))))
        idx++

        // Attention flash: "3 response baru!" in amber, flashes 2x
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(renderAttentionDot(opacities, idx, m.attentionStart))
        b.WriteString(" ")
        b.WriteString(renderAttentionFlash(
                fmt.Sprintf(i18n.T(i18n.KeyBootNewResponses), m.responseCount),
                m.attentionStart,
        ))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // Response prompt
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.WarningStyle.Render(i18n.T(i18n.KeyBootPressEnterResponses)))

        return b.String()
}

// ---------------------------------------------------------------------------
// viewReturningError renders the boot_returning_error variant (WA disconnected).
// Doc shows this variant WITHOUT the logo — just the ✗ error header.
// ---------------------------------------------------------------------------

func (m BootModel) viewReturningError() string {
        var b strings.Builder

        // Error header with red flash ✗ — no logo per doc spec
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(renderErrorCross(m.errorFlashStart))
        b.WriteString(" ")
        b.WriteString(style.DangerStyle.Bold(true).Render(i18n.T(i18n.KeyBootWADisconnected)))
        b.WriteString(style.Section(style.SectionGap))

        // Status indicators
        now := time.Now()
        opacities := m.breathing.Opacities(now)
        idx := 0

        b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Warning, style.TextDim, style.BodyStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootScraperStillRunning), m.nicheCount))))
        idx++
        b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootLeadsCount), m.leadsCount))))

        b.WriteString(style.Section(style.SectionGap))

        // Explanation + action
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyBootScraperOnlyNote)))
        b.WriteString("\n")
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyBootPress1LoginAgain)))

        return b.String()
}

// ---------------------------------------------------------------------------
// viewReturningConfigError renders the boot_returning_config_error variant.
// ---------------------------------------------------------------------------

func (m BootModel) viewReturningConfigError() string {
        var b strings.Builder

        // Logo
        b.WriteString(m.renderLogoReveal())
        b.WriteString(style.Section(style.SectionGap))

        // Error header with red flash ✗
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(renderErrorCross(m.errorFlashStart))
        b.WriteString(" ")
        b.WriteString(style.DangerStyle.Bold(true).Render(i18n.T(i18n.KeyBootConfigError)))
        b.WriteString(style.Section(style.SectionGap))

        // Status indicators
        now := time.Now()
        opacities := m.breathing.Opacities(now)
        idx := 0

        b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootWAConnectedCount), m.waCount))))
        idx++
        b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootSomeNichesOK), m.okNicheCount, m.errorNicheCount))))
        idx++
        b.WriteString(renderIndicatorLine(contentIndent, opacities, idx, style.Success, style.TextDim, style.BodyStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootLeadsCount), m.leadsCount))))

        b.WriteString(style.Section(style.SectionGap))

        // Explanation
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.BodyStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootConfigErrorNiche), m.errorNiche)))
        b.WriteString("\n")
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyBootOtherWorkersStill)))
        b.WriteString("\n")
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyBootPressVError)))

        b.WriteString(style.Section(style.SectionGap))

        // Footer actions
        footer := style.ActionStyle.Render("v") + "  " + i18n.T(i18n.KeyBootVViewError) +
                "    " + style.ActionStyle.Render("↵") + "  " + i18n.T(i18n.KeyBootDashboard) +
                "    " + style.ActionStyle.Render("q") + "  " + i18n.T(i18n.KeyBootQExit)
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.CaptionStyle.Render(footer))

        return b.String()
}

// ---------------------------------------------------------------------------
// viewReturningLicenseExpired renders the boot_returning_license_expired variant.
// ---------------------------------------------------------------------------

func (m BootModel) viewReturningLicenseExpired() string {
        var b strings.Builder

        // Logo
        b.WriteString(m.renderLogoReveal())
        b.WriteString(style.Section(style.SectionGap))

        // Error header with red flash ✗
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(renderErrorCross(m.errorFlashStart))
        b.WriteString(" ")
        b.WriteString(style.DangerStyle.Bold(true).Render(i18n.T(i18n.KeyBootLicenseExpired)))
        b.WriteString(style.Section(style.SectionGap))

        // Status indicators — all ● dim because total pause
        b.WriteString(renderDimIndicatorLine(contentIndent, style.MutedStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootWAConnectedCount), m.waCount))))
        b.WriteString(renderDimIndicatorLine(contentIndent, style.MutedStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootAllNichesPaused), m.nicheCount))))
        b.WriteString(renderDimIndicatorLine(contentIndent, style.MutedStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyBootLeadsCount), m.leadsCount))))

        b.WriteString(style.Section(style.SectionGap))

        // Explanation
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyBootLicenseExpiredMsg)))
        b.WriteString("\n")
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyBootRenewToContinue)))

        b.WriteString(style.Section(style.SectionGap))

        // Footer actions
        footer := style.ActionStyle.Render("1") + "  " + i18n.T(i18n.KeyBootEnterNewLicense) +
                "    " + style.ActionStyle.Render("2") + "  " + i18n.T(i18n.KeyBootBuyLicense) +
                "    " + style.ActionStyle.Render("q") + "  " + i18n.T(i18n.KeyBootQExit)
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.CaptionStyle.Render(footer))

        return b.String()
}

// ---------------------------------------------------------------------------
// viewReturningDeviceConflict renders the boot_returning_device_conflict variant.
// ---------------------------------------------------------------------------

func (m BootModel) viewReturningDeviceConflict() string {
        var b strings.Builder

        // Logo
        b.WriteString(m.renderLogoReveal())
        b.WriteString(style.Section(style.SectionGap))

        // Error header with red flash ✗
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(renderErrorCross(m.errorFlashStart))
        b.WriteString(" ")
        b.WriteString(style.DangerStyle.Bold(true).Render(i18n.T(i18n.KeyBootDeviceConflict)))
        b.WriteString(style.Section(style.SectionGap))

        // Explanation
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyBootDeviceConflictMsg)))
        b.WriteString("\n")
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyBootOneLicenseOneDevice)))
        b.WriteString("\n")
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyBootAllPausedUntilResolved)))

        b.WriteString(style.Section(style.SectionGap))

        // Footer actions
        footer := style.ActionStyle.Render("1") + "  " + i18n.T(i18n.KeyBootEnterNewLicense) +
                "    " + style.ActionStyle.Render("2") + "  " + i18n.T(i18n.KeyBootDisconnectOther) +
                "    " + style.ActionStyle.Render("q") + "  " + i18n.T(i18n.KeyBootQExit)
        b.WriteString(style.Indent(contentIndent))
        b.WriteString(style.CaptionStyle.Render(footer))

        // Device info separator
        if m.deviceName != "" {
                b.WriteString(style.Section(style.SectionGap))
                info := fmt.Sprintf(i18n.T(i18n.KeyBootOtherDeviceInfo), m.deviceName, m.lastActive)
                b.WriteString(style.Indent(0))
                b.WriteString(renderLabeledSeparator(info))
        }

        return b.String()
}

// ---------------------------------------------------------------------------
// Logo rendering — character-by-character reveal
// ---------------------------------------------------------------------------

// renderLogoReveal renders the WACLAW ASCII logo with per-character reveal.
// Characters appear one at a time at 8ms intervals (anim.LogoCharDelay).
func (m BootModel) renderLogoReveal() string {
        elapsed := time.Since(m.logoRevealStart)

        // Calculate total characters in the logo
        totalChars := 0
        for _, line := range waclawLogo {
                totalChars += len(line)
        }

        // Number of characters revealed so far
        charsToShow := int(elapsed / anim.LogoCharDelay)
        if charsToShow > totalChars {
                charsToShow = totalChars
        }

        // If fully revealed, render normally
        if charsToShow >= totalChars {
                return renderFullLogo()
        }

        // Partial reveal — render character by character
        var b strings.Builder
        shown := 0
        for _, line := range waclawLogo {
                remaining := charsToShow - shown
                if remaining <= 0 {
                        // Line not started yet — empty line
                        b.WriteString("\n")
                        continue
                }
                if remaining >= len(line) {
                        // Full line visible
                        b.WriteString(style.PrimaryStyle.Render(line))
                        b.WriteString("\n")
                } else {
                        // Partial line
                        visible := line[:remaining]
                        b.WriteString(style.PrimaryStyle.Render(visible))
                        b.WriteString("\n")
                }
                shown += len(line)
        }

        return b.String()
}

// renderFullLogo renders the complete WACLAW ASCII logo without animation.
func renderFullLogo() string {
        var b strings.Builder
        for _, line := range waclawLogo {
                b.WriteString(style.PrimaryStyle.Render(line))
                b.WriteString("\n")
        }
        // Empty line after logo (as in doc)
        b.WriteString("\n")
        return b.String()
}

// ---------------------------------------------------------------------------
// Army march animation — ▸▸▸ → ● morph with overshoot bounce
// ---------------------------------------------------------------------------

// renderArmyMarch renders the army marching animation.
// Each row slides in with 80ms stagger, then ▸▸▸ morphs to ● with
// overshoot bounce. Total duration: 600ms (anim.ArmyMarch).
//
// F10 FIX: After the march animation settles (elapsed >= anim.ArmyMarch),
// only the summary line "3 worker udah jalan. lu telat datang, mereka nggak."
// is shown. The individual worker rows are hidden to avoid visual overlap
// with the status indicators that follow. This matches the doc spec where
// the army march is a transient animation that settles into the report.
func (m BootModel) renderArmyMarch() string {
        // If no workers configured, skip the march animation
        if len(m.workers) == 0 {
                return ""
        }

        elapsed := time.Since(m.marchStart)

        // After the full march animation completes, only show the summary line.
        // The individual worker rows overlap with the status indicators, so
        // they are hidden once the animation settles.
        if elapsed >= anim.ArmyMarch {
                var b strings.Builder
                b.WriteString(style.Indent(contentIndent))
                b.WriteString(style.MutedStyle.Render(
                        fmt.Sprintf(i18n.T(i18n.KeyBootWorkerReady), len(m.workers)),
                ))
                b.WriteString("\n")
                return b.String()
        }

        var b strings.Builder

        for i, w := range m.workers {
                rowStart := time.Duration(i) * armyMarchStagger
                rowElapsed := elapsed - rowStart

                if rowElapsed < 0 {
                        // Row hasn't started yet
                        continue
                }

                // March phase: ▸▸▸ marching (first 60% of total march duration per row)
                // Settle phase: morph to ● with bounce (last 40%)
                marchDuration := anim.ArmyMarch - rowStart
                if marchDuration <= 0 {
                        marchDuration = anim.ArmyMarch
                }

                progress := float64(rowElapsed) / float64(marchDuration)
                if progress > 1.0 {
                        progress = 1.0
                }

                // Cascading indent: each row indents further for visual depth
                b.WriteString(style.Indent(contentIndent + i*2))

                if progress < 0.6 {
                        // Marching phase — show ▸ arrows with partial reveal
                        arrowsToShow := int(float64(w.ArrowCount) * anim.EaseOutCubic(progress/0.6))
                        if arrowsToShow < 1 {
                                arrowsToShow = 1
                        }
                        if arrowsToShow > w.ArrowCount {
                                arrowsToShow = w.ArrowCount
                        }
                        arrows := strings.Repeat("▸", arrowsToShow)
                        padding := strings.Repeat(" ", w.ArrowCount-arrowsToShow)
                        b.WriteString(style.AccentStyle.Render(arrows))
                        b.WriteString(padding)
                        b.WriteString("  ")
                        b.WriteString(style.MutedStyle.Render(w.Name))
                } else {
                        // Settle phase — morph to ● with overshoot bounce
                        settleProgress := (progress - 0.6) / 0.4
                        if settleProgress > 1.0 {
                                settleProgress = 1.0
                        }
                        bounceScale := anim.EaseOutBack(settleProgress)

                        // Render: padding + name + spaces + ● aktif
                        // Name comes BEFORE the ● aktif status per doc spec
                        b.WriteString(strings.Repeat(" ", w.ArrowCount))
                        b.WriteString("  ")
                        b.WriteString(style.MutedStyle.Render(w.Name))
                        // padding to align the ● aktif
                        b.WriteString("    ")
                        if bounceScale > 1.0 {
                                // Overshoot — bright pulse
                                b.WriteString(style.SuccessStyle.Bold(true).Render("●"))
                        } else {
                                b.WriteString(style.SuccessStyle.Render("●"))
                        }
                        b.WriteString(" ")
                        // F6 FIX: Use KeyStatusActive ("● aktif"/"● active") instead of the
// redundant KeyBootArmyMarching for the army march settle text. The
// army march "● aktif" status should use the same i18n key as all other
// ● active indicators, ensuring consistency and DRY.
                                        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyStatusActive)))
                }

                b.WriteString("\n")
        }

        // Worker summary line
        if elapsed >= anim.ArmyMarch {
                b.WriteString(style.Indent(contentIndent))
                b.WriteString(style.MutedStyle.Render(
                        fmt.Sprintf(i18n.T(i18n.KeyBootWorkerReady), len(m.workers)),
                ))
                b.WriteString("\n")
        }

        return b.String()
}

// ---------------------------------------------------------------------------
// Helper: state classification
// ---------------------------------------------------------------------------

// isReturningState returns true if the current state is any of the returning
// user variants (which should auto-transition to dashboard after 3s).
func (m BootModel) isReturningState() bool {
        switch m.state {
        case protocol.BootReturning,
                protocol.BootReturningResponse,
                protocol.BootReturningError,
                protocol.BootReturningConfigError:
                return true
        default:
                // License expired and device conflict do NOT auto-transition
                // (hard gate — requires user action)
                return false
        }
}

// indicatorCount returns the number of ● indicators for the current state,
// used to size the breathing group.
func (m BootModel) indicatorCount() int {
        switch m.state {
        case protocol.BootReturning:
                return 3 // wa, niche/worker, leads
        case protocol.BootReturningResponse:
                return 4 // wa, niche/worker, leads, responses
        case protocol.BootReturningError:
                return 2 // niche (scrape), leads
        case protocol.BootReturningConfigError:
                return 3 // wa, niches, leads
        case protocol.BootReturningLicenseExpired:
                return 3 // wa, niches, leads (all dim)
        default:
                return 0
        }
}

// Compile-time interface checks.
var (
        _ tui.Screen        = (*BootModel)(nil)
        _ protocol.ScreenID = protocol.ScreenBoot
)
