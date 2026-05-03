package tui

import (
        "fmt"
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/overlay"
        "github.com/WaClaw-App/waclaw/internal/tui/rpc"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
        "github.com/charmbracelet/bubbles/key"
        tea "github.com/charmbracelet/bubbletea"
        "github.com/charmbracelet/lipgloss"
)

// App is the root bubbletea.Model that orchestrates the entire TUI.
//
// It holds the Router (screen stack), Bus (event queue), and manages the
// screen lifecycle, global key handling, transition animations, overlay
// composition, and locale injection. Every bubbletea message flows through
// App.Update before being delegated to the current screen.
type App struct {
        router *Router
        bus    *bus.Bus
        width  int
        height int
        ready  bool

        // validation holds the current validation result overlay, or nil.
        validation *bus.ValidateMsg

        // Overlays — Phase 4A global overlays composed on top of any screen.
        nerdStats    overlay.NerdStats
        cmdPalette   overlay.CmdPalette
        notifToast   overlay.NotificationToast
        confirmation overlay.ConfirmationOverlay
        shortcuts    overlay.ShortcutsOverlay

        // rpcClient is the JSON-RPC client for backend communication.
        // Set via SetRPCBuilders so the App can forward key events.
        rpcClient     *rpc.Client
        keyBuilder    rpc.KeyPressBuilder
        actionBuilder rpc.ActionBuilder

        // startTime tracks when the app was launched (for nerd stats uptime).
        startTime time.Time

        // sessionEnd tracks the double-quit state for session end flow.
        // First q from root screen shows session summary + Zeigarnik tip.
        // Second q actually quits. Any other key dismisses the summary.
        sessionEnd *sessionEndState
}

// NewApp creates a new TUI application with default configuration.
func NewApp() *App {
        b := bus.New()
        r := NewRouter(b)

        cp := overlay.NewCmdPalette()
        cp.OnExecute = func(cmd overlay.Command) {
                // Navigation commands → push to router.
                if cmd.Screen != protocol.ScreenID("") {
                        r.Push(cmd.Screen)
                }
                // Quick actions are handled by the backend via SendAction.
                // The overlay just triggers the action; actual execution
                // happens through the RPC channel.
        }

        confirmation := overlay.NewConfirmationOverlay()
        // confirmation.OnConfirm is wired after App construction because
        // it needs access to App.SendAction. See the line after NewApp().

        return &App{
                router:       r,
                bus:          b,
                nerdStats:    overlay.NewNerdStats(),
                cmdPalette:   cp,
                notifToast:   overlay.NewNotificationToast(),
                confirmation: confirmation,
                shortcuts:    overlay.NewShortcutsOverlay(),
                startTime:    time.Now(),
        }
}

// Bus returns the event bus for external access (e.g., RPC client).
func (a *App) Bus() *bus.Bus { return a.bus }

// Router returns the screen router for external access.
func (a *App) Router() *Router { return a.router }

// RegisterScreen adds a screen to the application and injects the bus.
func (a *App) RegisterScreen(s Screen) {
        s.SetBus(a.bus)
        a.router.Register(s)
}

// SetRPCClient sets the RPC client for backend communication.
// The client is used to forward key events and action events from the TUI
// to the backend over stdio.
func (a *App) SetRPCClient(client *rpc.Client) {
        a.rpcClient = client
}

// SetRPCBuilders configures the key and action builders that construct
// protocol events with the current screen context. These builders are
// used by screen code to send events to the backend without importing
// the rpc package directly.
func (a *App) SetRPCBuilders(keyBuilder rpc.KeyPressBuilder, actionBuilder rpc.ActionBuilder) {
        a.keyBuilder = keyBuilder
        a.actionBuilder = actionBuilder
}

// SendKeyPress forwards a key press event to the backend via the RPC client.
// Returns an error if the client is not running. This is a convenience method
// so screen code doesn't need to import the rpc package.
func (a *App) SendKeyPress(k string) error {
        if a.rpcClient == nil {
                return nil // standalone mode — no backend
        }
        return a.rpcClient.SendKeyPress(a.keyBuilder.Build(k))
}

// SendAction forwards an action event to the backend via the RPC client.
func (a *App) SendAction(action string, params map[string]any) error {
        if a.rpcClient == nil {
                return nil // standalone mode — no backend
        }
        return a.rpcClient.SendAction(a.actionBuilder.Build(action, params))
}

// SendRequest forwards a data request to the backend via the RPC client.
func (a *App) SendRequest(reqType string, params map[string]any) (<-chan *protocol.Response, error) {
        if a.rpcClient == nil {
                return nil, nil // standalone mode — no backend
        }
        builder := rpc.RequestBuilder{
                Screen: a.actionBuilder.Screen,
        }
        return a.rpcClient.SendRequest(builder.Build(reqType, params))
}

// WireConfirmation connects the confirmation overlay's OnConfirm callback
// to the App's SendAction method. Must be called after NewApp() because
// the callback needs the App reference.
func (a *App) WireConfirmation() {
        a.confirmation.OnConfirm = func(confirmType protocol.ConfirmationType, data map[string]any) {
                // Forward the confirmed action to the backend as an action event.
                // The action name encodes the confirmation type (e.g., "bulk_offer",
                // "bulk_delete", "bulk_archive", "force_device_disconnect").
                a.SendAction(string(confirmType), data)
        }
}

// Init implements tea.Model. It initialises the first screen on the stack
// (if any) and requests the initial window size.
func (a *App) Init() tea.Cmd {
        var cmds []tea.Cmd

        // If a screen is already on the stack (pushed before startup),
        // call its Init.
        if cur := a.router.Current(); cur != nil {
                if cmd := cur.Init(); cmd != nil {
                        cmds = append(cmds, cmd)
                }
        }

        return tea.Batch(cmds...)
}

// Update implements tea.Model.
// Routes key events, handles global keys, processes bus messages,
// and delegates to the current screen.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        var cmds []tea.Cmd

        // 1. Window resize — update dimensions and propagate to all screens.
        switch m := msg.(type) {
        case tea.WindowSizeMsg:
                a.handleWindowResize(m)
                return a, nil
        }

        // Tick overlays.
        a.tickOverlays()

        // 2. Key messages — handle global keys first, then overlay keys, then screen.
        if keyMsg, ok := msg.(tea.KeyMsg); ok {
                if cmd, consumed := a.handleGlobalKeys(keyMsg); consumed {
                        return a, cmd
                }
        }

        // 3. Process bus messages before delegating to the screen.
        a.processBusMessages()

        // 4. Handle internal transition messages.
        switch msg.(type) {
        case TransitionCompleteMsg:
                a.router.ClearTransition()
                return a, nil
        }

        // 5. Delegate to the current screen.
        if cur := a.router.Current(); cur != nil {
                updated, cmd := cur.Update(msg)
                if cmd != nil {
                        cmds = append(cmds, cmd)
                }
                // If the updated model still satisfies Screen and has the same ID,
                // replace the screen in the router's map so subsequent calls use it.
                if updated != nil {
                        if scr, ok := updated.(Screen); ok && scr.ID() == cur.ID() {
                                a.router.screens[cur.ID()] = scr
                        }
                }
        }

        return a, tea.Batch(cmds...)
}

// View implements tea.Model.
// Renders the current screen (or transition animation if active), then
// composes overlays on top in priority order.
func (a *App) View() string {
        // If no screen is registered, show a placeholder.
        if a.router.Current() == nil {
                return lipgloss.NewStyle().
                        Background(style.Bg).
                        Foreground(style.TextMuted).
                        Width(a.width).
                        Height(a.height).
                        Render(i18n.T("app.no_screen"))
        }

        // If a transition is active, render the transition animation.
        if t := a.router.Transition(); t != nil {
                return a.renderTransition(t)
        }

        // Render the current screen's view with the background style.
        content := a.router.Current().View()

        // Overlay: nerd stats (bottom of screen).
        if a.nerdStats.IsVisible() {
                content = a.renderNerdStatsOverlay(content)
        }

        // Overlay: notification toast (top of screen).
        if a.notifToast.IsVisible() {
                content = a.renderNotifToastOverlay(content)
        }

        // Overlay: validation results (bottom).
        if a.validation != nil {
                content = a.renderValidation(content)
        }

        // Overlay: confirmation dialog (center, full-screen dim).
        if a.confirmation.IsVisible() {
                content = a.renderConfirmationOverlay(content)
        }

        // Overlay: command palette (top, centered).
        if a.cmdPalette.IsOpen() {
                content = a.renderCmdPaletteOverlay(content)
        }

        // Overlay: shortcuts (centered panel).
        if a.shortcuts.IsVisible() {
                content = a.renderShortcutsOverlay(content)
        }

        // Overlay: session end (double-quit summary + Zeigarnik tip).
        if a.sessionEnd != nil {
                content = a.renderSessionEnd(content)
        }

        return a.applyBackground(content)
}

// handleGlobalKeys processes keys that work on every screen.
// Uses the key.Binding infrastructure from keymap.go for consistent,
// DRY key matching — no raw string comparisons.
// Returns a tea.Cmd and true if the key was consumed globally.
func (a *App) handleGlobalKeys(msg tea.KeyMsg) (tea.Cmd, bool) {
        // ── Command palette (highest priority overlay) ──
        // Ctrl+K toggles the command palette. If already open, closes it.
        if key.Matches(msg, KeyCmdPalette) {
                if a.cmdPalette.IsOpen() {
                        a.cmdPalette.Close()
                } else {
                        a.cmdPalette.SetCurrentScreen(a.router.CurrentID())
                        a.cmdPalette.Open()
                }
                return nil, true
        }

        // If command palette is open, all keys go to it.
        if a.cmdPalette.IsOpen() {
                consumed := a.cmdPalette.HandleKey(msg)
                return nil, consumed
        }

        // If confirmation overlay is active, all keys go to it.
        if a.confirmation.IsVisible() {
                consumed := a.confirmation.HandleKey(msg.String())
                return nil, consumed
        }

        // ── Nerd stats toggle ──
        // Backtick toggles: hidden → minimal → expanded → hidden.
        if key.Matches(msg, KeyNerdStats) {
                a.nerdStats.Toggle()
                return nil, true
        }

        // ── Shortcuts overlay ──
        // "?" toggles the shortcuts overlay. Any keypress while visible dismisses it.
        if key.Matches(msg, KeyHelp) {
                a.shortcuts.Toggle()
                return nil, true
        }
        // Dismiss shortcuts on any other keypress when visible.
        if a.shortcuts.IsVisible() {
                a.shortcuts.Hide()
                // Don't consume the key — let it fall through to the screen.
        }

        // Dismiss session end overlay on any key that isn't q.
        // (The q handler above already handles the second-quit case.)
        if a.sessionEnd != nil {
                // Any key other than q dismisses the session end overlay.
                // The q key is handled above (second q = quit).
                a.sessionEnd = nil
                // Don't consume the key — let it fall through to the screen.
        }

        // ── Escape → close overlays ──
        if key.Matches(msg, KeyEscape) {
                // Dismiss notification toast if active.
                if a.notifToast.IsVisible() {
                        a.notifToast.ForceDismiss()
                        return nil, true
                }
                // Dismiss validation overlay if active.
                if a.validation != nil {
                        a.validation = nil
                        return nil, true
                }
                // Dismiss nerd stats if visible.
                if a.nerdStats.IsVisible() {
                        a.nerdStats.Hide()
                        return nil, true
                }
                return nil, true
        }

        // ── Navigation shortcuts ──
        // "l" → navigate to license screen.
        if key.Matches(msg, KeyLicense) {
                a.router.Push(protocol.ScreenLicense)
                return nil, true
        }

        // "h" → navigate to history screen.
        if key.Matches(msg, KeyHistory) {
                a.router.Push(protocol.ScreenHistory)
                return nil, true
        }

        // "u" → check for updates.
        if key.Matches(msg, KeyUpdate) {
                a.router.Push(protocol.ScreenUpdate)
                return nil, true
        }

        // "v" → validate config.
        if key.Matches(msg, KeyValidate) {
                a.router.Push(protocol.ScreenGuardrail)
                return nil, true
        }

        // "q" → session end (double-quit) or go back.
        if key.Matches(msg, KeyBack) {
                if !a.router.IsRoot() {
                        a.router.Pop()
                        return nil, true
                }
                // At root screen — implement double-quit per doc/20 spec.
                if a.sessionEnd != nil {
                        // Second q → actually quit.
                        return tea.Quit, true
                }
                // First q → show session end summary.
                a.sessionEnd = &sessionEndState{
                        shownAt: time.Now(),
                        stats:   a.collectSessionStats(),
                }
                return nil, true
        }

        return nil, false
}

// ToggleLocale switches between the two supported locales.
// Called from the command palette (Ctrl+K) or programmatically.
func ToggleLocale() {
        current := i18n.GetLocale()
        if current == i18n.LocaleID {
                i18n.SetLocale(i18n.LocaleEN)
        } else {
                i18n.SetLocale(i18n.LocaleID)
        }
}

// tickOverlays advances all overlay animation states and timers.
func (a *App) tickOverlays() {
        now := time.Now()

        // Nerd stats: auto-collapse check + metric refresh.
        a.nerdStats.Tick(now)
        if a.nerdStats.IsVisible() && a.nerdStats.ShouldRefreshMetrics() {
                a.nerdStats.RefreshMetrics()
                a.nerdStats.Metrics.Uptime = now.Sub(a.startTime)
        }

        // Command palette animation.
        a.cmdPalette.Tick()

        // Notification toast auto-dismiss.
        a.notifToast.Tick()

        // Confirmation animation.
        a.confirmation.Tick()

        // Shortcuts animation.
        a.shortcuts.Tick()
}

// processBusMessages handles pending messages from the event bus.
func (a *App) processBusMessages() {
        for _, msg := range a.bus.Pending() {
                switch m := msg.(type) {
                case bus.NavigateMsg:
                        a.navigateToScreen(m)
                case bus.UpdateMsg:
                        if cur := a.router.Current(); cur != nil {
                                _ = cur.HandleUpdate(m.Params)
                        }
                case bus.NotifyMsg:
                        // Convert to notification type enum and enqueue in toast overlay.
                        notifType := protocol.NotificationType(m.Type)
                        a.notifToast.Enqueue(overlay.NotificationDataFromMsg(notifType, m.Severity, m.Data))
                case bus.ValidateMsg:
                        a.validation = &m
                }
        }
}

// handleWindowResize updates the app dimensions and propagates to all screens.
func (a *App) handleWindowResize(msg tea.WindowSizeMsg) {
        a.width = msg.Width
        a.height = msg.Height
        a.ready = true

        // Update overlay widths.
        a.nerdStats.Width = msg.Width
        a.cmdPalette.Width = msg.Width
        a.notifToast.Width = msg.Width
        a.confirmation.Width = msg.Width
        a.shortcuts.Width = msg.Width

        // Propagate the WindowSizeMsg to all registered screens via Update.
        // Each screen handles tea.WindowSizeMsg internally.
        for _, s := range a.router.screens {
                s.Update(msg)
        }
}

// navigateToScreen handles a NavigateMsg from the bus.
func (a *App) navigateToScreen(msg bus.NavigateMsg) {
        // Validate the target screen is registered.
        if a.router.Screen(msg.Screen) == nil {
                return
        }

        // Check params for navigation mode.
        if replace, _ := msg.Params["replace"].(bool); replace {
                a.router.Replace(msg.Screen)
        } else {
                a.router.Push(msg.Screen)
        }

        // If the screen has navigation params, forward them.
        if cur := a.router.Current(); cur != nil && len(msg.Params) > 0 {
                _ = cur.HandleNavigate(msg.Params)
        }
}

// ── Overlay rendering methods ──────────────────────────────────────────

// renderNerdStatsOverlay composes the nerd stats overlay at the bottom.
func (a *App) renderNerdStatsOverlay(content string) string {
        statsView := a.nerdStats.View()
        if statsView == "" {
                return content
        }

        lines := strings.Split(content, "\n")
        if len(lines) > 2 {
                // Insert before the last line.
                return strings.Join(lines[:len(lines)-1], "\n") + "\n" + statsView
        }
        return content + "\n" + statsView
}

// renderNotifToastOverlay composes the notification toast at the top.
func (a *App) renderNotifToastOverlay(content string) string {
        toastView := a.notifToast.View()
        if toastView == "" {
                return content
        }
        return toastView + "\n" + content
}

// renderConfirmationOverlay composes the confirmation dialog (centered).
func (a *App) renderConfirmationOverlay(content string) string {
        confirmView := a.confirmation.View()
        if confirmView == "" {
                return content
        }

        // Dim the background content.
        dimmed := a.dimContent(content)

        return a.overlayCentered(dimmed, confirmView)
}

// renderCmdPaletteOverlay composes the command palette at the top.
func (a *App) renderCmdPaletteOverlay(content string) string {
        paletteView := a.cmdPalette.View()
        if paletteView == "" {
                return content
        }

        // Dim the background content while palette is open.
        dimmed := a.dimContent(content)

        return paletteView + "\n" + dimmed
}

// renderShortcutsOverlay composes the shortcuts panel (centered).
func (a *App) renderShortcutsOverlay(content string) string {
        shortcutsView := a.shortcuts.View()
        if shortcutsView == "" {
                return content
        }

        return a.overlayCentered(content, shortcutsView)
}

// dimContent renders content in a dimmed style for overlay backgrounds.
func (a *App) dimContent(content string) string {
        return lipgloss.NewStyle().Foreground(style.TextDim).Render(content)
}

// overlayCentered inserts an overlay view at the vertical center of the content.
func (a *App) overlayCentered(content string, overlayView string) string {
        contentLines := strings.Split(content, "\n")
        overlayLines := strings.Split(overlayView, "\n")
        insertAt := len(contentLines) / 2
        if insertAt > len(contentLines) {
                insertAt = len(contentLines)
        }

        result := make([]string, 0, len(contentLines)+len(overlayLines))
        result = append(result, contentLines[:insertAt]...)
        result = append(result, overlayLines...)
        result = append(result, contentLines[insertAt:]...)

        return strings.Join(result, "\n")
}

// ── Transition rendering (unchanged from Phase 1) ──────────────────────

// renderTransition produces a visual transition animation between two screens.
func (a *App) renderTransition(t *TransitionState) string {
        cur := a.router.Current()
        if cur == nil {
                a.router.ClearTransition()
                return lipgloss.NewStyle().Background(style.Bg).Width(a.width).Height(a.height).Render("")
        }

        // Determine previous and next views.
        nextView := cur.View()
        if t.PreviousView == "" {
                t.PreviousView = nextView
        }
        if t.NextView == "" {
                t.NextView = nextView
        }

        // Advance the animation.
        if t.Tick() {
                // Transition complete — render the final screen.
                a.router.ClearTransition()
                return a.applyBackground(t.NextView)
        }

        // During transition, render a transition effect.
        var rendered string
        switch t.Direction {
        case TransitionForward:
                rendered = a.renderSlideTransition(t, 1)
        case TransitionBack:
                rendered = a.renderSlideTransition(t, -1)
        case TransitionFade:
                rendered = a.renderFadeTransition(t)
        default:
                rendered = t.NextView
        }

        return a.applyBackground(rendered)
}

// renderSlideTransition produces a slide animation for push/pop navigation.
// Uses easing to create a smooth visual effect: the incoming screen
// progressively replaces the outgoing one as progress advances.
// direction=1 slides right (forward), direction=-1 slides left (back).
func (a *App) renderSlideTransition(t *TransitionState, direction int) string {
        p := EaseOutCubic(t.Progress)

        // Render both views split at the progress point.
        // This creates a visual slide effect using line-based partitioning.
        prevLines := strings.Split(t.PreviousView, "\n")
        nextLines := strings.Split(t.NextView, "\n")

        maxLines := len(prevLines)
        if len(nextLines) > maxLines {
                maxLines = len(nextLines)
        }
        if maxLines == 0 {
                return t.NextView
        }

        // Calculate how many lines from each view to show.
        splitLine := int(float64(maxLines) * p)

        var result []string
        for i := 0; i < maxLines; i++ {
                if i < splitLine {
                        // Show next view lines (already slid in).
                        if i < len(nextLines) {
                                result = append(result, nextLines[i])
                        }
                } else {
                        // Show previous view lines (being pushed out).
                        if i < len(prevLines) {
                                result = append(result, prevLines[i])
                        }
                }
        }

        return strings.Join(result, "\n")
}

// renderFadeTransition produces a cross-fade animation for replace navigation.
// Uses eased progress for a smooth visual blend.
func (a *App) renderFadeTransition(t *TransitionState) string {
        p := EaseInOutCubic(t.Progress)

        prevLines := strings.Split(t.PreviousView, "\n")
        nextLines := strings.Split(t.NextView, "\n")

        maxLines := len(prevLines)
        if len(nextLines) > maxLines {
                maxLines = len(nextLines)
        }
        if maxLines == 0 {
                return t.NextView
        }

        // Cross-fade: blend lines from both views based on progress.
        splitLine := int(float64(maxLines) * p)

        var result []string
        for i := 0; i < maxLines; i++ {
                if i < splitLine {
                        if i < len(nextLines) {
                                result = append(result, nextLines[i])
                        }
                } else {
                        if i < len(prevLines) {
                                result = append(result, prevLines[i])
                        }
                }
        }

        return strings.Join(result, "\n")
}

// renderValidation overlays validation results on top of the current content.
func (a *App) renderValidation(content string) string {
        if a.validation == nil {
                return content
        }

        v := a.validation
        var lines []string

        if len(v.Errors) > 0 {
                errStyle := lipgloss.NewStyle().Foreground(style.Danger)
                lines = append(lines, errStyle.Render(fmt.Sprintf(i18n.T("validation.error_count"), len(v.Errors))))
                for _, e := range v.Errors {
                        lines = append(lines, errStyle.Render("  "+e))
                }
        }

        if len(v.Warnings) > 0 {
                warnStyle := lipgloss.NewStyle().Foreground(style.Warning)
                lines = append(lines, warnStyle.Render(fmt.Sprintf(i18n.T("validation.warning_count"), len(v.Warnings))))
                for _, w := range v.Warnings {
                        lines = append(lines, warnStyle.Render("  "+w))
                }
        }

        if len(v.Errors) == 0 && len(v.Warnings) == 0 {
                okStyle := lipgloss.NewStyle().Foreground(style.Success)
                lines = append(lines, okStyle.Render(i18n.T("validation.all_valid")))
        }

        valStyle := lipgloss.NewStyle().
                Background(style.BgRaised).
                Width(a.width).
                Padding(0, 2)

        valContent := valStyle.Render(strings.Join(lines, "\n"))
        return content + "\n" + valContent
}

// applyBackground wraps the content with the app's background color.
func (a *App) applyBackground(content string) string {
        // Remove trailing whitespace to avoid background color bleeding.
        content = strings.TrimRight(content, "\n ")

        bgStyle := lipgloss.NewStyle().
                Background(style.Bg).
                Width(a.width).
                Render(content)

        return bgStyle
}

// String returns a debug representation of the App.
func (a *App) String() string {
        return fmt.Sprintf("App{screen=%s, depth=%d, ready=%v, nerdStats=%d, cmdPalette=%d}",
                a.router.CurrentID(), a.router.Depth(), a.ready,
                a.nerdStats.Mode, a.cmdPalette.Mode)
}

// ── Session end (double-quit) ──────────────────────────────────────────

// sessionEndState tracks the double-quit flow per doc/20 spec:
//
//      First q (at root) → show session summary + Zeigarnik tip
//      Second q           → actually quit
//      Any other key      → dismiss, continue
//
// The summary shows:
//   - Session stats (sent, responses, deals)
//   - Best lead (fastest response)
//   - Zeigarnik tip (best send time from history insights)
type sessionEndState struct {
        shownAt time.Time
        stats   sessionStats
}

// sessionStats holds the session summary data displayed on quit.
type sessionStats struct {
        SentCount      int
        ResponseCount  int
        DealCount      int
        BestLeadName   string
        BestLeadTime   string
        ZeigarnikTip   string
}

// collectSessionStats gathers session statistics for the end summary.
// In demo mode, uses mock data. In production, these come from the backend.
func (a *App) collectSessionStats() sessionStats {
        // Request stats from backend if available.
        if a.rpcClient != nil {
                resp, err := a.SendRequest("get_stats", nil)
                if err == nil && resp != nil {
                        select {
                        case r := <-resp:
                                if r != nil && r.Result != nil {
                                        if m, ok := r.Result.(map[string]any); ok {
                                                return sessionStatsFromMap(m)
                                        }
                                }
                        default:
                        }
                }
        }

        // Fallback: sensible demo defaults.
        return sessionStats{
                SentCount:     43,
                ResponseCount: 7,
                DealCount:     2,
                BestLeadName:  "kopi nusantara",
                BestLeadTime:  "11 menit",
                ZeigarnikTip:  i18n.T("session.zeigarnik_tip"),
        }
}

// sessionStatsFromMap extracts session stats from a backend response map.
func sessionStatsFromMap(m map[string]any) sessionStats {
        s := sessionStats{}
        if v, ok := m["sent_count"].(float64); ok {
                s.SentCount = int(v)
        }
        if v, ok := m["response_count"].(float64); ok {
                s.ResponseCount = int(v)
        }
        if v, ok := m["deal_count"].(float64); ok {
                s.DealCount = int(v)
        }
        if v, ok := m["best_lead_name"].(string); ok {
                s.BestLeadName = v
        }
        if v, ok := m["best_lead_time"].(string); ok {
                s.BestLeadTime = v
        }
        if v, ok := m["zeigarnik_tip"].(string); ok {
                s.ZeigarnikTip = v
        } else {
                s.ZeigarnikTip = i18n.T("session.zeigarnik_tip")
        }
        return s
}

// renderSessionEnd renders the session end overlay per doc/20 spec.
//
//      sampai jumpa.
//
//      hari ini: 43 terkirim · 7 response · 2 deal
//      lead terbaik: kopi nusantara (respond 11 menit)
//
//      ── tips: selasa jam 10 itu waktu terbaik lu kirim pesan ──
//
//      tekan q lagi buat keluar, atau apa aja buat lanjut
func (a *App) renderSessionEnd(content string) string {
        if a.sessionEnd == nil {
                return content
        }

        s := a.sessionEnd.stats
        var lines []string

        // Farewell line.
        farewellStyle := lipgloss.NewStyle().
                Foreground(style.Text).
                Bold(true)
        lines = append(lines, farewellStyle.Render(i18n.T("session.farewell")))
        lines = append(lines, "")

        // Stats line.
        statsStyle := lipgloss.NewStyle().Foreground(style.TextMuted)
        statsText := fmt.Sprintf(i18n.T("session.stats"),
                s.SentCount, s.ResponseCount, s.DealCount)
        lines = append(lines, statsStyle.Render(statsText))

        // Best lead line.
        if s.BestLeadName != "" {
                bestStyle := lipgloss.NewStyle().Foreground(style.Accent)
                bestText := fmt.Sprintf(i18n.T("session.best_lead"),
                        s.BestLeadName, s.BestLeadTime)
                lines = append(lines, bestStyle.Render(bestText))
        }
        lines = append(lines, "")

        // Zeigarnik tip line.
        tipStyle := lipgloss.NewStyle().Foreground(style.Warning)
        if s.ZeigarnikTip != "" {
                lines = append(lines, tipStyle.Render(s.ZeigarnikTip))
        }
        lines = append(lines, "")

        // Action line.
        actionStyle := lipgloss.NewStyle().Foreground(style.TextDim)
        lines = append(lines, actionStyle.Render(i18n.T("session.quit_hint")))

        // Render the session end panel.
        panelStyle := lipgloss.NewStyle().
                Background(style.Bg).
                Width(a.width).
                Padding(1, 2)

        panelContent := panelStyle.Render(strings.Join(lines, "\n"))

        // Dim the background and overlay the panel at center.
        dimmed := a.dimContent(content)
        return a.overlayCentered(dimmed, panelContent)
}
