// Package update implements Screen 20: UPDATE & UPGRADE.
//
// Visual spec: doc/12-screens-update-upgrade.md
//
// States:
//   - UpdateAvailable: minor update available (free, seamless)
//   - UpdateDownloading: download in progress with progress bar
//   - UpdateReady: download complete, ready to restart
//   - UpgradeAvailable: major upgrade (needs new license)
//   - UpgradeLicenseInput: entering v2 license key
//   - LicenseExpiredWithUpgrade: v1 expired + v2 available
//   - StartupCheck: background update check at boot (non-blocking)
//
// Design philosophy:
//   - Minor update = free, seamless, non-blocking
//   - Major upgrade = conscious choice, new product
//   - No forced upgrades, no dark patterns
//   - v1 army keeps running as long as license is valid
package update

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
        "github.com/charmbracelet/bubbles/key"
        "github.com/charmbracelet/bubbles/textinput"
        "github.com/charmbracelet/lipgloss"
)

// Model represents the Update & Upgrade screen.
type Model struct {
        // base embeds the ScreenBase for default Screen interface methods.
        base tui.ScreenBase

        // state is the current screen state.
        state protocol.StateID

        // width and height are the terminal dimensions.
        width  int
        height int

        // focused indicates whether this screen is the active screen.
        focused bool

        // --- Update data (populated by HandleNavigate/HandleUpdate) ---

        // currentVersion is the installed version string (e.g. "v1.3.2").
        currentVersion string

        // newVersion is the available version string (e.g. "v1.3.3" or "v2.0.0").
        newVersion string

        // changelog lists the changes in the new version.
        changelog []string

        // licenseKey is the user's current v1 license key.
        licenseKey string

        // licenseExpiry is the current license expiry date string.
        licenseExpiry string

        // licenseStatus is the current license status text.
        licenseStatus string

        // --- Download state ---

        // downloadPercent is the download progress [0.0, 1.0].
        downloadPercent float64

        // downloadSize is the total download size string (e.g. "5.1MB").
        downloadSize string

        // downloadDownloaded is the amount downloaded so far (e.g. "4.2MB").
        downloadDownloaded string

        // downloadSpeed is the current download speed (e.g. "350KB/s").
        downloadSpeed string

        // downloadETA is the estimated time remaining string (e.g. "12 detik").
        downloadETA string

        // downloadSource is the download URL.
        downloadSource string

        // --- Ready state ---

        // checksumVerified indicates the download checksum was verified.
        checksumVerified bool

        // backupPath is the path to the backup of the old binary.
        backupPath string

        // --- License input state ---

        // licenseInput is the text input for the v2 license key.
        licenseInput textinput.Model

        // licenseInputFocused indicates whether the license input is active.
        licenseInputFocused bool

        // --- Expired state ---

        // expiredDate is the date the v1 license expired.
        expiredDate string

        // --- Animation state ---

        // changelogStaggerIndex tracks which changelog items have been revealed.
        changelogStaggerIndex int

        // changelogStaggerStart is when the stagger animation began.
        changelogStaggerStart time.Time

        // freeBadgePulseStart is when the "free" badge pulse began.
        freeBadgePulseStart time.Time

        // freeBadgePulsing indicates whether the green pulse is active.
        freeBadgePulsing bool

        // downloadBar is the progress bar component for downloading state.
        // Uses █ fill character per doc spec ("████████████████████░░░░  82%").
        downloadBar component.ProgressBar

        // startupCheckResult holds the result of the background startup check.
        // Valid values are protocol.CheckResult constants: CheckResultLatest,
        // CheckResultUpdate, CheckResultUpgrade. Stored as string because
        // HandleNavigate/HandleUpdate receive params as map[string]any.
        startupCheckResult string

        // startupCheckFired indicates whether the background startup check
        // has been scheduled (t+250ms per doc spec).
        startupCheckFired bool

        // licensePrefix is the v2 license key prefix from backend
        // (default: protocol.LicenseKeyPrefixV2 + "-").
        licensePrefix string
}

// NewModel creates a new Update screen model with default values.
func NewModel() Model {
        // Initialise the license key text input.
        ti := textinput.New()
        ti.Placeholder = i18n.T(i18n.KeyUpgradeLicensePlaceholder)
        ti.Focus()
        ti.CharLimit = protocol.LicenseKeyFormattedLen
        ti.Width = style.UpgradeInputWidth

        // Progress bar with doc-matching fill character: █ (full block, U+2588).
        bar := component.NewProgressBar(style.UpgradeProgressBarWidth)
        bar.FillChar = "█"
        bar.EmptyChar = "░"

        return Model{
                base:               tui.NewScreenBase(protocol.ScreenUpdate),
                state:              protocol.UpdateAvailable,
                licenseInput:       ti,
                downloadBar:        bar,
                startupCheckResult: "",
        }
}

// ID returns the screen identifier.
func (m Model) ID() protocol.ScreenID {
        return m.base.ID()
}

// SetBus injects the event bus reference.
func (m *Model) SetBus(b *bus.Bus) {
        m.base.SetBus(b)
}

// HandleNavigate processes a "navigate" command from the backend.
// The params map carries screen-specific navigation parameters.
func (m *Model) HandleNavigate(params map[string]any) error {
        if params == nil {
                return nil
        }

        // Extract state from params.
        if s, ok := params[protocol.ParamState].(string); ok {
                m.state = protocol.StateID(s)
        }

        // Extract version info.
        if v, ok := params[protocol.ParamCurrentVersion].(string); ok {
                m.currentVersion = v
        }
        if v, ok := params[protocol.ParamNewVersion].(string); ok {
                m.newVersion = v
        }

        // Extract changelog.
        if items, ok := params[protocol.ParamChangelog].([]any); ok {
                m.setChangelog(items)
        }

        // Extract license info.
        if k, ok := params[protocol.ParamLicenseKey].(string); ok {
                m.licenseKey = k
        }
        if e, ok := params[protocol.ParamLicenseExpiry].(string); ok {
                m.licenseExpiry = e
        }
        if s, ok := params[protocol.ParamLicenseStatus].(string); ok {
                m.licenseStatus = s
        }

        // Extract v2 license prefix from backend (DRY: no hardcoded "WACL2-").
        if lp, ok := params[protocol.ParamLicensePrefix].(string); ok && lp != "" {
                m.licensePrefix = lp
        }

        // Trigger animations when entering relevant states.
        m.onStateEnter(m.state)

        // Extract startup check result.
        if r, ok := params[protocol.ParamCheckResult].(string); ok {
                m.startupCheckResult = r
        }

        return nil
}

// HandleUpdate processes an "update" command from the backend.
// Handles download progress updates and state transitions.
func (m *Model) HandleUpdate(params map[string]any) error {
        if params == nil {
                return nil
        }

        // Handle state change.
        if s, ok := params[protocol.ParamState].(string); ok {
                newState := protocol.StateID(s)
                m.state = newState
                m.onStateEnter(newState)

                if newState == protocol.StartupCheck {
                        // Startup check is non-blocking — just record the result.
                        if r, ok := params[protocol.ParamCheckResult].(string); ok {
                                m.startupCheckResult = r
                        }
                }
        }

        // Update download progress.
        if pct, ok := params[protocol.ParamPercent].(float64); ok {
                m.downloadPercent = pct
                m.downloadBar = m.downloadBar.SetPercent(pct)
        }
        if sz, ok := params[protocol.ParamSize].(string); ok {
                m.downloadSize = sz
        }
        if dl, ok := params[protocol.ParamDownloaded].(string); ok {
                m.downloadDownloaded = dl
        }
        if spd, ok := params[protocol.ParamSpeed].(string); ok {
                m.downloadSpeed = spd
        }
        if eta, ok := params[protocol.ParamETA].(string); ok {
                m.downloadETA = eta
        }
        if src, ok := params[protocol.ParamSource].(string); ok {
                m.downloadSource = src
        }

        // Update ready state info.
        if verified, ok := params[protocol.ParamChecksumVerified].(bool); ok {
                m.checksumVerified = verified
        }
        if bp, ok := params[protocol.ParamBackupPath].(string); ok {
                m.backupPath = bp
        }

        // Update expired state info.
        if d, ok := params[protocol.ParamExpiredDate].(string); ok {
                m.expiredDate = d
        }

        // Update changelog (may be sent incrementally).
        if items, ok := params[protocol.ParamChangelog].([]any); ok {
                m.setChangelog(items)
        }

        return nil
}

// onStateEnter triggers animations and side effects when entering a state.
// DRY helper — single place for all state-entry logic, eliminating the
// duplication that previously existed between HandleNavigate and HandleUpdate.
func (m *Model) onStateEnter(state protocol.StateID) {
        if state == protocol.UpdateAvailable {
                m.startFreeBadgePulse()
        }
        if state == protocol.UpgradeLicenseInput {
                m.licenseInputFocus()
                // Auto-prefix for v2 license input (doc spec: "prefix auto-detect").
                // Uses the licensePrefix from backend params, falling back to protocol constant if not set.
                prefix := m.licensePrefix
                if prefix == "" {
                        prefix = protocol.LicenseKeyPrefixV2 + "-"
                }
                if val := m.licenseInput.Value(); val == "" && !strings.HasPrefix(val, prefix) {
                        m.licenseInput.SetValue(prefix)
                }
        }
}

// startFreeBadgePulse begins the green pulse animation for the "free" badge.
// DRY helper — was previously duplicated in HandleNavigate and HandleUpdate.
func (m *Model) startFreeBadgePulse() {
        m.freeBadgePulseStart = time.Now()
        m.freeBadgePulsing = true
}

// Focus is called when this screen becomes the active screen.
func (m *Model) Focus() {
        m.focused = true
        // Re-trigger the stagger animation so changelog items reveal on focus.
        if len(m.changelog) > 0 {
                m.changelogStaggerIndex = 0
                m.changelogStaggerStart = time.Now()
        }
}

// Blur is called when this screen is no longer the active screen.
func (m *Model) Blur() {
        m.focused = false
        m.freeBadgePulsing = false
}

// Init implements tea.Model.
// Starts the changelog stagger tick so items reveal one-by-one (doc: "stagger 50ms").
// Without this initial tick, the stagger animation never begins.
func (m Model) Init() tea.Cmd {
        if len(m.changelog) > 0 || m.freeBadgePulsing {
                return tea.Tick(anim.ChangelogStagger, func(t time.Time) tea.Msg {
                        return tickMsg(t)
                })
        }
        return nil
}

// tickMsg is an internal message for animation ticks.
type tickMsg time.Time

// startupCheckMsg is fired after the t+250ms background update check delay.
type startupCheckMsg time.Time

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        var cmds []tea.Cmd

        switch msg := msg.(type) {
        case tea.WindowSizeMsg:
                m.width = msg.Width
                m.height = msg.Height
                m.downloadBar.Width = min(msg.Width-style.DownloadBarPadding, style.DownloadBarMaxWidth)

        case tea.KeyMsg:
                cmds = append(cmds, m.handleKey(msg))

        case tickMsg:
                // Advance changelog stagger animation by exactly 1 item per tick.
                // This guarantees the one-at-a-time visual effect the doc intends,
                // rather than computing from elapsed time which can cause batch reveals.
                if m.changelogStaggerIndex < len(m.changelog) {
                        m.changelogStaggerIndex++
                }

                // Auto-clear the free badge pulse after SuccessPulse (500ms) expires.
                if m.freeBadgePulsing && time.Since(m.freeBadgePulseStart) >= anim.SuccessPulse {
                        m.freeBadgePulsing = false
                }

                // Continue ticking if animation is still in progress.
                if m.changelogStaggerIndex < len(m.changelog) || m.freeBadgePulsing {
                        cmds = append(cmds, tea.Tick(anim.ChangelogStagger, func(t time.Time) tea.Msg {
                                return tickMsg(t)
                        }))
                }

        case startupCheckMsg:
                // Background startup check fired (doc: "t +250ms update check").
                // This is non-blocking — results are delivered as notifications.
                m.startupCheckFired = true
        }

        // Update the license text input if in license input state.
        if m.state == protocol.UpgradeLicenseInput && m.licenseInputFocused {
                var cmd tea.Cmd
                m.licenseInput, cmd = m.licenseInput.Update(msg)
                if cmd != nil {
                        cmds = append(cmds, cmd)
                }
        }

        return m, tea.Batch(cmds...)
}

// handleKey processes key events for the current state.
func (m *Model) handleKey(msg tea.KeyMsg) tea.Cmd {
        switch m.state {
        case protocol.UpdateAvailable:
                return m.handleUpdateAvailableKeys(msg)
        case protocol.UpdateDownloading:
                return m.handleDownloadingKeys(msg)
        case protocol.UpdateReady:
                return m.handleReadyKeys(msg)
        case protocol.UpgradeAvailable:
                return m.handleUpgradeAvailableKeys(msg)
        case protocol.UpgradeLicenseInput:
                return m.handleLicenseInputKeys(msg)
        case protocol.LicenseExpiredWithUpgrade:
                return m.handleExpiredKeys(msg)
        case protocol.StartupCheck:
                // Startup check is non-blocking — no key handling needed.
                // Results are delivered as notifications after the dashboard is ready.
                return nil
        }
        return nil
}

// handleUpdateAvailableKeys handles key events in the update_available state.
// Doc: "↵  update sekarang   1  nanti aja   q  skip"
func (m *Model) handleUpdateAvailableKeys(msg tea.KeyMsg) tea.Cmd {
        if key.Matches(msg, tui.KeyEnter) {
                // Start downloading the update — state transition confirmed by backend.
                return m.sendAction(string(protocol.ActionStartDownload), map[string]any{
                        protocol.ParamVersion: m.newVersion,
                })
        }
        if key.Matches(msg, tui.Key1) {
                // Remind later — go back to previous screen.
                return m.sendAction(string(protocol.ActionRemindLater), map[string]any{
                        protocol.ParamVersion: m.newVersion,
                })
        }
        if key.Matches(msg, tui.KeyBack) {
                // Skip this update entirely.
                return m.sendAction(string(protocol.ActionSkipUpdate), map[string]any{
                        protocol.ParamVersion: m.newVersion,
                })
        }
        return nil
}

// handleDownloadingKeys handles key events in the update_downloading state.
// Doc: "q  batal download"
func (m *Model) handleDownloadingKeys(msg tea.KeyMsg) tea.Cmd {
        if key.Matches(msg, tui.KeyBack) {
                // Cancel the download — let backend confirm state change.
                return m.sendAction(string(protocol.ActionCancelDownload), nil)
        }
        return nil
}

// handleReadyKeys handles key events in the update_ready state.
// Doc: "↵  restart sekarang   1  nanti (update tetap ada)   q  skip"
func (m *Model) handleReadyKeys(msg tea.KeyMsg) tea.Cmd {
        if key.Matches(msg, tui.KeyEnter) {
                // Restart to apply the update.
                return m.sendAction(string(protocol.ActionRestartNow), map[string]any{
                        protocol.ParamVersion: m.newVersion,
                })
        }
        if key.Matches(msg, tui.Key1) {
                // Defer restart — keep the downloaded update for later.
                return m.sendAction(string(protocol.ActionRestartLater), map[string]any{
                        protocol.ParamVersion: m.newVersion,
                })
        }
        if key.Matches(msg, tui.KeyBack) {
                // Skip — update file remains but won't auto-apply.
                return m.sendAction(string(protocol.ActionSkipRestart), nil)
        }
        return nil
}

// handleUpgradeAvailableKeys handles key events in the upgrade_available state.
// Doc: "↵  beli lisensi v2   1  liat detail v2   q  tetap v1"
func (m *Model) handleUpgradeAvailableKeys(msg tea.KeyMsg) tea.Cmd {
        if key.Matches(msg, tui.KeyEnter) {
                // Go to license input for v2 — let backend drive the state transition.
                return m.sendAction(string(protocol.ActionBuyLicense), nil)
        }
        if key.Matches(msg, tui.Key1) {
                // View v2 details — backend sends the details.
                return m.sendAction(string(protocol.ActionViewUpgradeDetails), map[string]any{
                        protocol.ParamVersion: m.newVersion,
                })
        }
        if key.Matches(msg, tui.KeyBack) {
                // Stay on v1.
                return m.sendAction(string(protocol.ActionStayV1), nil)
        }
        return nil
}

// handleLicenseInputKeys handles key events in the upgrade_license_input state.
// Doc: "↵  validasi v2    q  batal"
func (m *Model) handleLicenseInputKeys(msg tea.KeyMsg) tea.Cmd {
        if key.Matches(msg, tui.KeyEnter) {
                // Validate the entered license key.
                return m.sendAction(string(protocol.ActionValidateLicense), map[string]any{
                        protocol.ParamKey: m.licenseInput.Value(),
                })
        }
        if key.Matches(msg, tui.KeyBack) {
                // Cancel — go back to upgrade available.
                m.licenseInputBlur()
                return m.sendAction(string(protocol.ActionCancelLicenseInput), nil)
        }
        // All other keys are handled by the text input.
        return nil
}

// handleExpiredKeys handles key events in the license_expired_with_upgrade state.
// Doc: "1  perpanjang v1   2  upgrade ke v2   3  masukin lisensi baru   q  keluar"
func (m *Model) handleExpiredKeys(msg tea.KeyMsg) tea.Cmd {
        if key.Matches(msg, tui.Key1) {
                return m.sendAction(string(protocol.ActionRenewV1), nil)
        }
        if key.Matches(msg, tui.Key2) {
                return m.sendAction(string(protocol.ActionUpgradeV2), nil)
        }
        if key.Matches(msg, tui.Key3) {
                return m.sendAction(string(protocol.ActionEnterNewLicense), nil)
        }
        if key.Matches(msg, tui.KeyBack) {
                return m.sendAction(string(protocol.ActionExitExpired), nil)
        }
        return nil
}

// licenseInputFocus activates the license key text input.
func (m *Model) licenseInputFocus() {
        m.licenseInputFocused = true
        m.licenseInput.Focus()
}

// licenseInputBlur deactivates the license key text input.
func (m *Model) licenseInputBlur() {
        m.licenseInputFocused = false
        m.licenseInput.Blur()
}

// sendAction publishes an action message to the bus.
func (m *Model) sendAction(action string, params map[string]any) tea.Cmd {
        if m.base.Bus() == nil {
                return nil
        }
        m.base.Bus().Publish(bus.ActionMsg{
                Action: action,
                Screen: protocol.ScreenUpdate,
                Params: params,
        })
        return nil
}

// setChangelog parses a raw changelog list and resets the stagger animation.
// DRY helper — used by both HandleNavigate and HandleUpdate.
func (m *Model) setChangelog(items []any) {
        m.changelog = make([]string, 0, len(items))
        for _, item := range items {
                if s, ok := item.(string); ok {
                        m.changelog = append(m.changelog, s)
                }
        }
        m.changelogStaggerIndex = 0
        m.changelogStaggerStart = time.Now()
}

// resetDownload resets the download progress state.
// DRY helper — used when starting or cancelling a download.
func (m *Model) resetDownload() {
        m.downloadPercent = 0
        m.downloadBar = m.downloadBar.SetPercent(0)
}

// View implements tea.Model. Renders the current screen state.
func (m Model) View() string {
        switch m.state {
        case protocol.UpdateAvailable:
                return m.viewUpdateAvailable()
        case protocol.UpdateDownloading:
                return m.viewDownloading()
        case protocol.UpdateReady:
                return m.viewReady()
        case protocol.UpgradeAvailable:
                return m.viewUpgradeAvailable()
        case protocol.UpgradeLicenseInput:
                return m.viewLicenseInput()
        case protocol.LicenseExpiredWithUpgrade:
                return m.viewExpiredWithUpgrade()
        case protocol.StartupCheck:
                return m.viewStartupCheck()
        default:
                return m.viewUpdateAvailable()
        }
}

// viewUpdateAvailable renders the "minor update available" state.
//
// Visual spec from doc/12-screens-update-upgrade.md:
//
//      update tersedia
//
//      ────────────────────────────────────────────────────
//
//      versi lu:     v1.3.2
//      versi baru:   v1.3.3
//
//      yang baru:
//        • fix: scrape terkadang duplicate di area tertentu
//        • fix: WA rotator cooldown calculation lebih akurat
//        • perf: database query 15% lebih cepat
//
//      ini update kecil (minor). gratis.
//      lisensi lu tetap valid. nggak perlu bayar.
//
//      ────────────────────────────────────────────────────
//
//      ↵  update sekarang   1  nanti aja   q  skip
//
//      update butuh restart. waclaw bakal nunggu sampe
//      semua worker idle sebelum restart.
func (m Model) viewUpdateAvailable() string {
        var b strings.Builder

        // Title — the screen heading.
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyUpdateAvailable)))
        b.WriteString(style.Section(style.SectionGap))

        // Separator line.
        b.WriteString(style.Separator())
        b.WriteString(style.Section(style.SubSectionGap))

        // Version info.
        b.WriteString(m.renderVersionInfo())
        b.WriteString(style.Section(style.SectionGap))

        // Changelog items with stagger animation.
        b.WriteString(m.renderChangelog())
        b.WriteString(style.Section(style.SectionGap))

        // "Free" badge with green pulse.
        b.WriteString(m.renderFreeBadge())
        b.WriteString(style.Section(style.SubSectionGap))

        // Separator.
        b.WriteString(style.Separator())
        b.WriteString(style.Section(style.SectionGap))

        // Action bar.
        // Doc: "↵  update sekarang   1  nanti aja   q  skip"
        b.WriteString(m.renderActionBar([]actionItem{
                {key: "↵", label: i18n.T(i18n.KeyUpdateNow)},
                {key: "1", label: i18n.T(i18n.KeyUpdateLater)},
                {key: "q", label: i18n.T(i18n.KeyGeneralSkip)},
        }))
        b.WriteString(style.Section(style.SectionGap))

        // Footnote about restart.
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyUpdateRestartNote)))

        return b.String()
}

// viewDownloading renders the "downloading update" state.
//
// Visual spec from doc:
//
//      update tersedia
//
//      lagi download v1.3.3...
//
//      ████████████████████░░░░  82% (4.2MB / 5.1MB)
//
//      estimasi: 12 detik
//
//      sumber: https://releases.waclaw.dev/v1.3.3/
//
//      worker tetap jalan selama download.
//      restart cuma setelah download selesai + lu approve.
//
//       q  batal download
func (m Model) viewDownloading() string {
        var b strings.Builder

        // Title.
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyUpdateAvailable)))
        b.WriteString(style.Section(style.SectionGap))

        // Downloading label with version.
        b.WriteString(style.BodyStyle.Render(
                fmt.Sprintf("%s %s...", i18n.T(i18n.KeyUpdateDownloadingLabel), m.newVersion),
        ))
        b.WriteString(style.Section(style.SectionGap))

        // Progress bar.
        // Doc spec: "████████████████████░░░░  82% (4.2MB / 5.1MB)"
        // Percentage BEFORE the parenthesized size label.
        percentText := fmt.Sprintf("%3.0f%%", m.downloadPercent*100)
        label := fmt.Sprintf("(%s / %s)", m.downloadDownloaded, m.downloadSize)
        m.downloadBar.Label = fmt.Sprintf("%s %s", percentText, label)
        m.downloadBar.ShowPercent = false // We include percent in the label for doc-parity.
        b.WriteString(m.downloadBar.View())
        b.WriteString(style.Section(style.SectionGap))

        // ETA.
        if m.downloadETA != "" {
                b.WriteString(style.MutedStyle.Render(
                        fmt.Sprintf("%s: %s", i18n.T(i18n.KeyUpdateETA), m.downloadETA),
                ))
                b.WriteString("\n")
        }

        // Source URL.
        if m.downloadSource != "" {
                b.WriteString(style.CaptionStyle.Render(
                        fmt.Sprintf("%s: %s", i18n.T(i18n.KeyUpdateSource), m.downloadSource),
                ))
                b.WriteString("\n")
        }

        // Download speed — doc spec: "bar fills real-time + speed shown = lu tau ini kerja".
        if m.downloadSpeed != "" {
                b.WriteString(style.CaptionStyle.Render(
                        fmt.Sprintf("%s: %s", i18n.T(i18n.KeyUpdateSpeed), m.downloadSpeed),
                ))
                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SectionGap))

        // Worker reassurance.
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyUpdateWorkerNote)))
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyUpdateRestartApproval)))
        b.WriteString(style.Section(style.SectionGap))

        // Cancel action.
        // Doc: "q  batal download"
        b.WriteString(m.renderActionBar([]actionItem{
                {key: "q", label: i18n.T(i18n.KeyUpdateCancelDownload)},
        }))

        return b.String()
}

// viewReady renders the "update ready to install" state.
//
// Visual spec from doc:
//
//      update tersedia
//
//      ✓ v1.3.3 siap di-install!
//
//      ────────────────────────────────────────────────────
//
//      download selesai: 5.1MB
//      checksum: ✓ verified
//      backup: ~/.waclaw/waclaw-v1.3.2.bak
//
//      ────────────────────────────────────────────────────
//
//      restart sekarang? semua worker bakal di-pause dulu.
//      proses restart ±3 detik. data aman.
//
//       ↵  restart sekarang   1  nanti (update tetap ada)   q  skip
//
//      kalau lu skip, waclaw bakal ngingetin lagi pas startup berikutnya.
func (m Model) viewReady() string {
        var b strings.Builder

        // Title.
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyUpdateAvailable)))
        b.WriteString(style.Section(style.SectionGap))

        // Success line with green pulse.
        successText := fmt.Sprintf("✓ %s %s!", m.newVersion, i18n.T(i18n.KeyUpdateReadyInstall))
        b.WriteString(style.SuccessStyle.Render(successText))
        b.WriteString(style.Section(style.SectionGap))

        // Separator.
        b.WriteString(style.Separator())
        b.WriteString(style.Section(style.SubSectionGap))

        // Download complete info.
        b.WriteString(style.BodyStyle.Render(
                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyUpdateDownloadComplete), m.downloadSize),
        ))
        b.WriteString("\n")

        // Checksum verification — render label and value separately to avoid
        // nested lipgloss style conflicts (inner SuccessStyle + outer BodyStyle).
        b.WriteString(style.BodyStyle.Render(
                fmt.Sprintf("%s: ", i18n.T(i18n.KeyUpdateChecksum)),
        ))
        if m.checksumVerified {
                b.WriteString(style.SuccessStyle.Render("✓ " + i18n.T(i18n.KeyUpdateVerified)))
        } else {
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyUpdateChecksumPending)))
        }
        b.WriteString("\n")

        // Backup path.
        if m.backupPath != "" {
                b.WriteString(style.BodyStyle.Render(
                        fmt.Sprintf("%s: %s", i18n.T(i18n.KeyUpdateBackup), m.backupPath),
                ))
                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SubSectionGap))

        // Separator.
        b.WriteString(style.Separator())
        b.WriteString(style.Section(style.SectionGap))

        // Restart prompt.
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyUpdateRestartPrompt)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyUpdateRestartDuration)))
        b.WriteString(style.Section(style.SectionGap))

        // Action bar.
        // Doc: "↵  restart sekarang   1  nanti (update tetap ada)   q  skip"
        b.WriteString(m.renderActionBar([]actionItem{
                {key: "↵", label: i18n.T(i18n.KeyUpdateRestartNow)},
                {key: "1", label: i18n.T(i18n.KeyUpdateLaterKeep)},
                {key: "q", label: i18n.T(i18n.KeyGeneralSkip)},
        }))
        b.WriteString(style.Section(style.SectionGap))

        // Skip reminder footnote.
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyUpdateSkipReminder)))

        return b.String()
}

// viewUpgradeAvailable renders the "major upgrade available" state.
//
// Visual spec from doc:
//
//      upgrade tersedia
//
//      ────────────────────────────────────────────────────
//
//      versi lu:     v1.3.2
//      versi baru:   v2.0.0
//
//      ⚠ ini upgrade besar (major version)
//
//      yang baru:
//        • arsitektur baru: multi-device support
//        • AI-powered message personalization
//        • dashboard real-time collaboration
//        • 20+ fitur baru lainnya
//
//      ────────────────────────────────────────────────────
//
//      upgrade ke v2 butuh lisensi baru.
//      lisensi v1 lu nggak berlaku buat v2.
//
//      ini bukan pelit — ini product yang beda.
//      v1 lu tetep bisa dipakai selama lisensi valid.
//      nggak ada forced upgrade.
//
//      lisensi v1 lu:
//        key: WACL-XXXX-XXXX-XXXX-XXXX
//        expires: 30 juni 2025
//        status: ✓ valid — v1 tetap jalan
//
//       ↵  beli lisensi v2   1  liat detail v2   q  tetap v1
func (m Model) viewUpgradeAvailable() string {
        var b strings.Builder

        // Title.
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyUpgradeAvailable)))
        b.WriteString(style.Section(style.SectionGap))

        // Separator.
        b.WriteString(style.Separator())
        b.WriteString(style.Section(style.SubSectionGap))

        // Version info.
        b.WriteString(m.renderVersionInfo())
        b.WriteString(style.Section(style.SubSectionGap))

        // Warning badge — amber, not red (doc: "⚠ amber = perhatian, bukan bahaya").
        b.WriteString(style.WarningStyle.Render(
                fmt.Sprintf("⚠ %s", i18n.T(i18n.KeyUpgradeMajorWarning)),
        ))
        b.WriteString(style.Section(style.SectionGap))

        // Changelog with stagger.
        b.WriteString(m.renderChangelog())
        b.WriteString(style.Section(style.SubSectionGap))

        // Separator.
        b.WriteString(style.Separator())
        b.WriteString(style.Section(style.SectionGap))

        // License requirement explanation.
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyUpgradeNeedsLicense)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyUpgradeV1NotValid)))
        b.WriteString(style.Section(style.SubSectionGap))

        // Reassurance — no dark patterns.
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyUpgradeNotGreedy)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyUpgradeV1StillWorks)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyUpgradeNoForced)))
        b.WriteString(style.Section(style.SectionGap))

        // Current v1 license info.
        if m.licenseKey != "" {
                b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyUpgradeV1License)))
                b.WriteString("\n")
                b.WriteString(style.CaptionStyle.Render(
                        fmt.Sprintf("  %s: %s", i18n.T(i18n.KeyUpgradeLicenseKey), m.licenseKey),
                ))
                b.WriteString("\n")
                if m.licenseExpiry != "" {
                        b.WriteString(style.CaptionStyle.Render(
                                fmt.Sprintf("  %s: %s", i18n.T(i18n.KeyUpgradeLicenseExpires), m.licenseExpiry),
                        ))
                        b.WriteString("\n")
                }
                if m.licenseStatus != "" {
                        // Render label and status separately to avoid nested style conflicts.
                        b.WriteString(style.CaptionStyle.Render(
                                fmt.Sprintf("  %s: ", i18n.T(i18n.KeyUpgradeLicenseStatus)),
                        ))
                        if protocol.LicenseResult(m.licenseStatus) == protocol.LicenseResultValid {
                                b.WriteString(style.SuccessStyle.Render(
                                        fmt.Sprintf("✓ %s", i18n.T(i18n.KeyUpgradeV1StillRuns)),
                                ))
                        } else {
                                b.WriteString(style.CaptionStyle.Render(m.licenseStatus))
                        }
                        b.WriteString("\n")
                }
        }

        b.WriteString(style.Section(style.SectionGap))

        // Action bar.
        // Doc: "↵  beli lisensi v2   1  liat detail v2   q  tetap v1"
        b.WriteString(m.renderActionBar([]actionItem{
                {key: "↵", label: i18n.T(i18n.KeyUpgradeBuyV2)},
                {key: "1", label: i18n.T(i18n.KeyUpgradeViewDetails)},
                {key: "q", label: i18n.T(i18n.KeyUpgradeStayV1)},
        }))

        return b.String()
}

// viewLicenseInput renders the "enter v2 license key" state.
//
// Visual spec from doc:
//
//      upgrade ke v2
//
//      masukin lisensi v2 lu.
//
//      WACL2-XXXX-XXXX-XXXX-XXXX   (borderless raised input)
//
//      lisensi v1 lu tetap aktif buat v1.
//      kalau v2 lisensi nggak valid, waclaw tetap jalan v1.
//
//       ↵  validasi v2    q  batal
func (m Model) viewLicenseInput() string {
        var b strings.Builder

        // Title.
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyUpgradeV2Title)))
        b.WriteString(style.Section(style.SectionGap))

        // Instruction.
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyUpgradeV2EnterLicense)))
        b.WriteString(style.Section(style.SubSectionGap))

        // License input field — vertical borderless per P3.
        // Doc shows bordered box but P3 (Charm.sh vertical borderless) takes precedence.
        // We use BgRaised to visually distinguish the input area.
        inputStyle := lipgloss.NewStyle().
                Foreground(style.Text).
                Background(style.BgRaised).
                Width(style.UpgradeInputWidth).
                Padding(0, 1)

        b.WriteString(inputStyle.Render(m.licenseInput.View()))
        b.WriteString(style.Section(style.SectionGap))

        // Reassurance about v1 license.
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyUpgradeV1StillActive)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyUpgradeV2InvalidFallback)))
        b.WriteString(style.Section(style.SectionGap))

        // Action bar.
        // Doc: "↵  validasi v2    q  batal"
        b.WriteString(m.renderActionBar([]actionItem{
                {key: "↵", label: i18n.T(i18n.KeyUpgradeValidateV2)},
                {key: "q", label: i18n.T(i18n.KeyLabelBack)},
        }))

        return b.String()
}

// viewExpiredWithUpgrade renders the "v1 license expired + v2 available" state.
//
// Visual spec from doc:
//
//      lisensi expired
//
//      ────────────────────────────────────────────────────
//
//      lisensi v1 lu udah expired (15 april 2025).
//      v1 army berhenti.
//
//      tapi... ada v2!
//
//      ✓ v2.0.0 tersedia
//      ✓ lisensi v2 = akses ke semua fitur baru
//      ✓ bisa pake data yang sama (auto-migrate)
//
//      ────────────────────────────────────────────────────
//
//      1  perpanjang v1 (lisensi lama, fitur lama)
//      2  upgrade ke v2 (lisensi baru, fitur baru)
//      3  masukin lisensi baru (v1 atau v2)
//
//      q  keluar
func (m Model) viewExpiredWithUpgrade() string {
        var b strings.Builder

        // Title.
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyLicenseExpiredTitle)))
        b.WriteString(style.Section(style.SectionGap))

        // Separator.
        b.WriteString(style.Separator())
        b.WriteString(style.Section(style.SubSectionGap))

        // Expiry notice.
        b.WriteString(style.BodyStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyLicenseExpiredV1), m.expiredDate),
        ))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLicenseExpiredV1Stopped)))
        b.WriteString(style.Section(style.SectionGap))

        // Positive pivot — v2 available.
        // Doc describes this as a "positive pivot" — use AccentStyle (indigo) not WarningStyle.
        b.WriteString(style.AccentStyle.Render(i18n.T(i18n.KeyLicenseExpiredButV2)))
        b.WriteString(style.Section(style.SubSectionGap))

        // v2 benefits list.
        b.WriteString(style.SuccessStyle.Render(
                fmt.Sprintf("✓ %s %s", m.newVersion, i18n.T(i18n.KeyLicenseExpiredV2Available)),
        ))
        b.WriteString("\n")
        b.WriteString(style.SuccessStyle.Render(i18n.T(i18n.KeyLicenseExpiredV2Features)))
        b.WriteString("\n")
        b.WriteString(style.SuccessStyle.Render(i18n.T(i18n.KeyLicenseExpiredV2Migrate)))
        b.WriteString(style.Section(style.SubSectionGap))

        // Separator.
        b.WriteString(style.Separator())
        b.WriteString(style.Section(style.SectionGap))

        // Options — each on its own line per doc wireframe (not a horizontal action bar).
        // Doc shows numbered options stacked vertically, not inline.
        options := []actionItem{
                {key: "1", label: i18n.T(i18n.KeyLicenseExpiredRenewV1)},
                {key: "2", label: i18n.T(i18n.KeyLicenseExpiredUpgradeV2)},
                {key: "3", label: i18n.T(i18n.KeyLicenseExpiredNewLicense)},
        }
        for _, opt := range options {
                b.WriteString(style.ActionStyle.Render(opt.key))
                b.WriteString("  ")
                b.WriteString(style.MutedStyle.Render(opt.label))
                b.WriteString("\n")
        }
        b.WriteString("\n")
        // Doc: "q  keluar"
        b.WriteString(style.CaptionStyle.Render(
                fmt.Sprintf("q  %s", i18n.T(i18n.KeyLicenseExpiredExit)),
        ))

        return b.String()
}

// ---------------------------------------------------------------------------
// Rendering helpers
// ---------------------------------------------------------------------------

// renderVersionInfo renders the current/new version comparison.
func (m Model) renderVersionInfo() string {
        var b strings.Builder

        b.WriteString(style.CaptionStyle.Render(
                fmt.Sprintf("%s:", i18n.T(i18n.KeyUpdateCurrentVersion)),
        ))
        b.WriteString("  ")
        b.WriteString(style.BodyStyle.Render(m.currentVersion))
        b.WriteString("\n")

        b.WriteString(style.CaptionStyle.Render(
                fmt.Sprintf("%s:", i18n.T(i18n.KeyUpdateNewVersion)),
        ))
        b.WriteString("  ")

        // New version gets accent color — it's the interactive element.
        b.WriteString(style.AccentStyle.Render(m.newVersion))

        return b.String()
}

// renderChangelog renders the changelog items with stagger fade-in animation.
// Each item appears 50ms after the previous one (doc spec: "stagger 50ms").
func (m Model) renderChangelog() string {
        if len(m.changelog) == 0 {
                return ""
        }

        var b strings.Builder

        b.WriteString(style.SubHeadingStyle.Render(
                fmt.Sprintf("%s:", i18n.T(i18n.KeyUpdateWhatsNew)),
        ))
        b.WriteString("\n")

        now := time.Now()
        for i, item := range m.changelog {
                // Only render items that have been revealed by the stagger animation.
                if i >= m.changelogStaggerIndex {
                        break
                }

                // Each item has a slight opacity increase during its first ItemFadeIn duration.
                itemAge := now.Sub(m.changelogStaggerStart) - time.Duration(i+1)*anim.ChangelogStagger
                if itemAge < 0 {
                        continue
                }

                // Use breathing opacity for recently revealed items.
                opacity := 1.0
                if itemAge < anim.ItemFadeIn {
                        opacity = 0.5 + 0.5*float64(itemAge)/float64(anim.ItemFadeIn)
                }

                rendered := style.Indent(1)
                rendered += "• "

                if opacity >= 0.95 {
                        rendered += style.BodyStyle.Render(item)
                } else {
                        rendered += style.MutedStyle.Render(item)
                }

                b.WriteString(rendered)
                b.WriteString("\n")
        }

        return b.String()
}

// renderFreeBadge renders the "free" / "gratis" badge with a green pulse.
// The pulse lasts for SuccessPulse (500ms) and signals an instant positive.
func (m Model) renderFreeBadge() string {
        elapsed := time.Since(m.freeBadgePulseStart)

        badgeText := i18n.T(i18n.KeyUpdateMinorFree)
        licenseText := i18n.T(i18n.KeyUpdateLicenseValid)

        // During the first 500ms, the badge pulses with success color + bold.
        if m.freeBadgePulsing && elapsed < anim.SuccessPulse {
                return style.SuccessStyle.Bold(true).Render(badgeText) + "\n" +
                        style.MutedStyle.Render(licenseText)
        }

        // After the pulse settles, use normal success style.
        return style.SuccessStyle.Render(badgeText) + "\n" +
                style.MutedStyle.Render(licenseText)
}

// actionItem represents a key-label pair for the action bar.
type actionItem struct {
        key   string
        label string
}

// renderActionBar renders a horizontal action bar with key-label pairs.
// Uses accent color for keys and muted style for labels.
// Vertical borderless — no boxes, just spacing.
// Doc shows "↵  label" format (just the arrow, no "go"/"gas" suffix).
func (m Model) renderActionBar(actions []actionItem) string {
        var parts []string
        for _, a := range actions {
                part := style.ActionStyle.Render(a.key) + "  " + style.MutedStyle.Render(a.label)
                parts = append(parts, part)
        }
        return strings.Join(parts, style.Indent(1))
}

// viewStartupCheck renders the background update check variant.
//
// Visual spec from doc/12-screens-update-upgrade.md:
//
//      t +250ms  update check (background, non-blocking)
//                ✓ v1.3.2 (latest) — nggak ada update
//                ○ v1.3.3 tersedia → notif: "update Available"
//                ○ v2.0.0 tersedia → notif: "upgrade Available"
//
// This variant runs during the startup sequence, is non-blocking, and
// delivers results as notifications after the dashboard is ready.
func (m Model) viewStartupCheck() string {
        var b strings.Builder

        // Title.
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyStartupCheckTitle)))
        b.WriteString(style.Section(style.SectionGap))

        // Show the check result based on what the backend reported.
        switch protocol.CheckResult(m.startupCheckResult) {
        case protocol.CheckResultLatest:
                b.WriteString(style.SuccessStyle.Render(
                        fmt.Sprintf(i18n.T(i18n.KeyStartupCheckLatest), m.currentVersion),
                ))
        case protocol.CheckResultUpdate:
                b.WriteString(style.MutedStyle.Render(
                        fmt.Sprintf(i18n.T(i18n.KeyStartupCheckUpdate), m.newVersion),
                ))
        case protocol.CheckResultUpgrade:
                b.WriteString(style.WarningStyle.Render(
                        fmt.Sprintf(i18n.T(i18n.KeyStartupCheckUpgrade), m.newVersion),
                ))
        default:
                // No result yet — still checking (non-blocking, so just show spinner).
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyStatusChecking)))
        }

        return b.String()
}
