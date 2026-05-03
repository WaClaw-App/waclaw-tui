// Package pipeline implements the auto-pilot sending status screen
// (Screen 6: SEND) for the WaClaw TUI. It renders 8 distinct states
// covering active sending, pause, off-hours, rate/daily limits, failures,
// all-slots-down, and response interrupts.
package pipeline

import (
        "fmt"
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/internal/tui/component"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        tui "github.com/WaClaw-App/waclaw/internal/tui"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
        tea "github.com/charmbracelet/bubbletea"
        "github.com/charmbracelet/bubbles/key"
        "github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Send screen
// ---------------------------------------------------------------------------

// Send is the auto-pilot sending status screen. It implements tui.Screen
// and renders 8 states: active, paused, off-hours, rate-limited,
// daily-limit, failed, all-slots-down, and with-response.
type Send struct {
        tui.ScreenBase
        state      protocol.StateID
        width      int
        height     int
        slots      []SlotInfo
        queue      []SendQueueItem
        stats      SendStats
        failure    *SendFailure
        response   *ResponseInterrupt
        nicheTab   int      // Currently selected niche tab (0-based)
        nicheNames []string // Ordered niche names for tab switching
        workHours  string              // Work hours from backend config (e.g. "09:00-17:00")
        breathing  component.BreathingGroup
}

// NewSend creates a new Send screen with the default active state.
func NewSend() *Send {
        return &Send{
                ScreenBase: tui.NewScreenBase(protocol.ScreenSend),
                state:      protocol.SendActive,
                slots:      make([]SlotInfo, 0),
                queue:      make([]SendQueueItem, 0),
                breathing:  component.NewBreathingGroup(3),
        }
}

// Init returns the initial command for the send screen.
func (s *Send) Init() tea.Cmd {
        return nil
}

// Focus is called when this screen becomes the active screen.
func (s *Send) Focus() {}

// Blur is called when this screen is no longer the active screen.
func (s *Send) Blur() {}

// ConsumesKey implements tui.KeyConsumer. The Send screen claims "v" in the
// send_failed state so it can trigger validate-retry instead of navigating
// to the Guardrail screen.
func (s *Send) ConsumesKey(msg tea.KeyMsg) bool {
        switch msg.String() {
        case "v":
                return s.state == protocol.SendFailed
        }
        return false
}

// SetSize updates the terminal dimensions for layout calculations.
func (s *Send) SetSize(w, h int) {
        s.width = w
        s.height = h
}

// HandleNavigate processes navigate commands from the backend, switching
// the internal state machine.
func (s *Send) HandleNavigate(params map[string]any) error {
        if st, ok := params[protocol.ParamState].(string); ok {
                s.state = protocol.StateID(st)
        }
        if st, ok := params[protocol.ParamPause].(string); ok && st != "" {
                s.state = protocol.SendPaused
        }
        if st, ok := params[protocol.ParamOffHours].(string); ok && st != "" {
                s.state = protocol.SendOffHours
        }
        if wh, ok := params[protocol.ParamWorkHours].(string); ok && wh != "" {
                s.workHours = wh
        }
        // Rebuild breathing group when state changes.
        s.breathing = component.NewBreathingGroup(3)
        return nil
}

// HandleUpdate processes update commands from the backend, refreshing
// the screen's data model.
func (s *Send) HandleUpdate(params map[string]any) error {
        if slots, ok := params[protocol.ParamSlots].([]any); ok {
                s.slots = parseSlotInfos(slots)
        }
        if queue, ok := params[protocol.ParamQueue].([]any); ok {
                s.queue = parseSendQueueItems(queue)
        }
        if stats, ok := params[protocol.ParamStats].(map[string]any); ok {
                s.stats = parseSendStats(stats)
        }
        if failure, ok := params[protocol.ParamFailure].(map[string]any); ok {
                s.failure = parseSendFailure(failure)
        }
        if response, ok := params[protocol.ParamResponse].(map[string]any); ok {
                s.response = parseResponseInterrupt(response)
        }
        if names, ok := params[protocol.ParamNicheNames].([]any); ok {
                s.nicheNames = parseStringSlice(names)
                if s.nicheTab >= len(s.nicheNames) {
                        s.nicheTab = 0
                }
        }
        // Individual field updates for real-time counters.
        if v, ok := params[protocol.ParamRateHour].(int); ok {
                s.stats.RateHour = int64(v)
        }
        if v, ok := params[protocol.ParamDailySent].(int); ok {
                s.stats.DailySent = int64(v)
        }
        if v, ok := params[protocol.ParamQueueCount].(int); ok {
                s.stats.QueueTotal = int64(v)
        }
        if v, ok := params[protocol.ParamNextSendTime].(string); ok {
                s.stats.NextSendTime = v
        }
        if wh, ok := params[protocol.ParamWorkHours].(string); ok && wh != "" {
                s.workHours = wh
        }
        return nil
}

// Update handles keyboard input and routes bus messages.
func (s *Send) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch msg := msg.(type) {
        case tea.KeyMsg:
                return s.handleKey(msg)
        case tea.WindowSizeMsg:
                s.SetSize(msg.Width, msg.Height)
        case bus.NavigateMsg:
                if msg.Screen == protocol.ScreenSend {
                        _ = s.HandleNavigate(msg.Params)
                        return s, nil
                }
        case bus.UpdateMsg:
                if msg.Screen == protocol.ScreenSend {
                        _ = s.HandleUpdate(msg.Params)
                        return s, nil
                }
        }
        return s, nil
}

// handleKey dispatches keyboard events to context-sensitive actions.
func (s *Send) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch {
        case key.Matches(msg, tui.KeyBack):
                s.publishKeyPress("q")
                return s, nil

        case key.Matches(msg, tui.KeyPause):
                s.publishAction(string(protocol.ActionTogglePause))
                return s, nil

        case key.Matches(msg, tui.KeyTab):
                s.cycleNicheTab()
                return s, nil

        case key.Matches(msg, tui.KeyEnter):
                s.handleEnter()
                return s, nil

        case key.Matches(msg, tui.Key1):
                s.handleKey1()
                return s, nil

        case key.Matches(msg, tui.Key2):
                s.handleKey2()
                return s, nil

        case key.Matches(msg, tui.KeySkip):
                s.handleSkip()
                return s, nil

        case msg.String() == "v" && s.state == protocol.SendFailed:
                s.publishAction(string(protocol.ActionValidateRetry))
                return s, nil
        }
        return s, nil
}

// View renders the current state of the send screen.
func (s *Send) View() string {
        switch s.state {
        case protocol.SendActive:
                return s.viewActive()
        case protocol.SendPaused:
                return s.viewPaused()
        case protocol.SendOffHours:
                return s.viewOffHours()
        case protocol.SendRateLimited:
                return s.viewRateLimited()
        case protocol.SendDailyLimit:
                return s.viewDailyLimit()
        case protocol.SendFailed:
                return s.viewFailed()
        case protocol.SendAllSlotsDown:
                return s.viewAllSlotsDown()
        case protocol.SendWithResponse:
                return s.viewWithResponse()
        default:
                return s.viewActive()
        }
}

// ---------------------------------------------------------------------------
// State renderers
// ---------------------------------------------------------------------------

// viewActive renders the multi-niche batch sending view with WA rotator.
func (s *Send) viewActive() string {
        var b strings.Builder

        // Title row: left-aligned title, right-aligned summary.
        title := style.HeadingStyle.Render(i18n.T(i18n.KeySendActive))
        right := style.CaptionStyle.Render(s.formatSummaryHeader())
        b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, title, fillGap(s.width, len(title)+len(right)), right))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.Section(style.SectionGap))

        // WA Rotator section.
        b.WriteString(s.renderWARotator())
        b.WriteString(style.Section(style.SectionGap))

        // Per-niche queue groups.
        b.WriteString(s.renderNicheQueue())
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.Section(style.SectionGap))

        // Rate and daily stats.
        b.WriteString(s.renderRateStats())
        b.WriteString(style.Section(style.SubSectionGap))

        // Auto-pilot message.
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeySendAutoPilotMsg)))
        b.WriteString(style.Section(style.SubSectionGap))

        // Key hints.
        b.WriteString(s.renderActiveHints())

        return b.String()
}

// viewPaused renders the user-paused state.
func (s *Send) viewPaused() string {
        var b strings.Builder

        title := style.HeadingStyle.Render(i18n.T(i18n.KeySendActive))
        right := fmt.Sprintf("%d %s", s.stats.QueueTotal, i18n.T(i18n.KeySendQueueCount))
        b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, title, fillGap(s.width, len(title)+len(right)), style.CaptionStyle.Render(right)))
        b.WriteString(style.Section(style.SectionGap))

        // Pause banner.
        b.WriteString(style.WarningStyle.Render("⏸  " + i18n.T(i18n.KeyStatusPaused) + " — " + i18n.T(i18n.KeySendPauseReason)))
        b.WriteString(style.Section(style.SubSectionGap))

        // Queue items with statuses.
        b.WriteString(s.renderQueueItems())
        b.WriteString(style.Section(style.SectionGap))

        // Rate stats.
        b.WriteString(s.renderRateStats())
        b.WriteString(style.Section(style.SubSectionGap))

        // Resume hint.
        b.WriteString(s.renderPausedHints())

        return b.String()
}

// viewOffHours renders the outside-work-hours state.
func (s *Send) viewOffHours() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeySendActive)))
        b.WriteString(style.Section(style.SectionGap))

        workHours := s.workHours
        if workHours == "" {
                workHours = i18n.T(i18n.KeySendDefaultWorkHours)
        }
        b.WriteString(style.WarningStyle.Render("⏰  " + i18n.T(i18n.KeySendOffHours) + " (" + workHours + ")"))
        b.WriteString(style.Section(style.ItemGap))

        nowStr := s.stats.Now
        if nowStr == "" {
                nowStr = "--:-- wib"
        }
        b.WriteString(style.MutedStyle.Render(fmt.Sprintf("%s: %s", i18n.T(i18n.KeySendNow), nowStr)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeySendAutoContinue)))
        b.WriteString(style.Section(style.SectionGap))

        // Emergency send option.
        b.WriteString(style.CaptionStyle.Render("1  " + i18n.T(i18n.KeySendEmergency)))
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(s.renderBackHint())

        return b.String()
}

// viewRateLimited renders the hourly rate-limit state.
func (s *Send) viewRateLimited() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeySendActive)))
        b.WriteString(style.Section(style.SectionGap))

        rateLabel := fmt.Sprintf("⏳  %s (%d/%d)", i18n.T(i18n.KeySendRateLimit), s.stats.RateHour, s.stats.RateHourLimit)
        b.WriteString(style.WarningStyle.Render(rateLabel))
        b.WriteString(style.Section(style.ItemGap))

        // Progress bar showing hourly rate usage.
        ratio := float64(s.stats.RateHour) / float64(maxInt(int(s.stats.RateHourLimit), 1))
        bar := component.NewProgressBar(minInt(s.width-4, style.DefaultProgressBarWidth)).
                SetPercent(ratio)
        bar.AtLimit = ratio >= 0.8
        bar.ShowPercent = false
        bar.Label = fmt.Sprintf("%d/%d", s.stats.RateHour, s.stats.RateHourLimit)
        b.WriteString(bar.View())
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeySendIntervalHint)))
        b.WriteString(style.Section(style.SectionGap))

        // Two-line reassurance.
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeySendAutoManages)))
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeySendNoRefreshNeeded)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(s.renderBackHint())

        return b.String()
}

// viewDailyLimit renders the daily-limit-reached state.
func (s *Send) viewDailyLimit() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeySendActive)))
        b.WriteString(style.Section(style.SectionGap))

        // Daily limit message.
        limitMsg := fmt.Sprintf("📊  %s (%d/%d)", i18n.T(i18n.KeySendDailyLimit), s.stats.DailySent, s.stats.DailyLimit)
        b.WriteString(style.WarningStyle.Render(limitMsg))
        b.WriteString(style.Section(style.ItemGap))

        // "hari ini:" section label.
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeySendToday) + ":"))
        b.WriteString(style.Section(style.ItemGap))

        // Remaining queue per niche.
        if s.stats.QueueTotal > 0 {
                remainMsg := fmt.Sprintf("%s: %d → %s", i18n.T(i18n.KeySendRemaining), s.stats.QueueTotal, i18n.T(i18n.KeySendAutoContinue))
                b.WriteString(style.BodyStyle.Render(remainMsg))
                if len(s.nicheNames) > 0 {
                        b.WriteString(style.Section(style.ItemGap))
                        b.WriteString(s.renderRemainingByNiche())
                }
        }
        b.WriteString(style.Section(style.SectionGap))

        // Today's stats with breathing effect.
        b.WriteString(s.renderTodayStats())
        b.WriteString(style.Section(style.SectionGap))

        // Two-line reassurance.
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeySendGoodJob)))
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeySendContinueTomorrow)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(s.renderBackHint())

        return b.String()
}

// viewFailed renders the individual send failure state.
func (s *Send) viewFailed() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeySendActive)))
        b.WriteString(style.Section(style.SectionGap))

        if s.failure != nil {
                // Failure banner.
                b.WriteString(style.DangerStyle.Render(fmt.Sprintf("✗  %s — %s", s.failure.Name, i18n.T(i18n.KeySendFailed))))
                b.WriteString(style.Section(style.ItemGap))
                b.WriteString(style.MutedStyle.Render(s.failure.Reason))
                b.WriteString(style.Section(style.SectionGap))

                // Warning that this should be rare.
                b.WriteString(style.WarningStyle.Render("⚠  " + i18n.T(i18n.KeySendShouldNot)))
                b.WriteString(style.Section(style.ItemGap))
                b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeySendPreValMiss)))
                b.WriteString(style.Section(style.ItemGap))
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeySendAutoSkip)))
        }
        b.WriteString(style.Section(style.SectionGap))

        // Options.
        b.WriteString(style.CaptionStyle.Render("1  " + i18n.T(i18n.KeySendRetryManual)))
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(style.CaptionStyle.Render("s  " + i18n.T(i18n.KeySendSkipAja)))
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(style.CaptionStyle.Render("v  " + i18n.T(i18n.KeySendValidateRetry)))
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(s.renderBackHint())

        return b.String()
}

// viewAllSlotsDown renders the all-WA-disconnected state.
func (s *Send) viewAllSlotsDown() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeySendActive)))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.DangerStyle.Render("✗  " + i18n.T(i18n.KeySendAllSlotsDown) + "!"))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.Section(style.SectionGap))

        // Per-slot status (all down).
        for i, slot := range s.slots {
                slotLabel := slot.Number
                if slotLabel == "" {
                        slotLabel = fmt.Sprintf("slot-%d", i+1)
                }
                label := fmt.Sprintf("📱 %s  %s", slotLabel, renderWAStatus(slot.Status))
                b.WriteString(style.CaptionStyle.Render(label))
                if i < len(s.slots)-1 {
                        b.WriteString("\n")
                }
        }
        b.WriteString(style.Section(style.SectionGap))

        // Reassurance messages.
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeySendAllScraping) + "."))
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeySendLeadsSafe) + "."))
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeySendPendingOnly) + "."))
        b.WriteString(style.Section(style.SectionGap))

        // Options.
        b.WriteString(style.CaptionStyle.Render("1  " + i18n.T(i18n.KeySendLogin)))
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(style.CaptionStyle.Render("2  " + i18n.T(i18n.KeySendLoginOneByOne)))
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(s.renderBackHint())

        return b.String()
}

// viewWithResponse renders the normal sending view with a response
// interrupt overlay at the top.
func (s *Send) viewWithResponse() string {
        var b strings.Builder

        // Title with queue count.
        title := style.HeadingStyle.Render(i18n.T(i18n.KeySendActive))
        right := fmt.Sprintf("%d %s", s.stats.QueueTotal, i18n.T(i18n.KeySendQueueCount))
        b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, title, fillGap(s.width, len(title)+len(right)), style.CaptionStyle.Render(right)))
        b.WriteString(style.Section(style.SectionGap))

        // Response interrupt overlay.
        if s.response != nil {
                b.WriteString(s.renderResponseOverlay())
                b.WriteString(style.Section(style.SectionGap))
        }

        // Continue with active sending view below.
        b.WriteString(s.renderQueueItems())
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(s.renderBackHint())

        return b.String()
}

// ---------------------------------------------------------------------------
// Section renderers
// ---------------------------------------------------------------------------

// renderWARotator renders the WA rotator section showing all slots.
func (s *Send) renderWARotator() string {
        var b strings.Builder

        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeySendWARotator)))
        b.WriteString(style.Section(style.SubSectionGap))

        for i, slot := range s.slots {
                statusStr := renderWAStatus(slot.Status)
                rateStr := fmt.Sprintf("%d/%d %s %s", slot.SentThisHour, slot.HourLimit, i18n.T(i18n.KeySendRate), i18n.T(i18n.KeyWordThisHour))

                // Show cooldown or ready-in time.
                var timerStr string
                switch slot.Status {
                case SlotStatusActive:
                        timerStr = fmt.Sprintf("%s: %s", i18n.T(i18n.KeyStatusCooldown), slot.CooldownLeft)
                case SlotStatusCooldown:
                        timerStr = fmt.Sprintf("%s: %s", i18n.T(i18n.KeySendReady), slot.ReadyIn)
                case SlotStatusDown:
                        timerStr = ""
                default:
                        timerStr = ""
                }

                line := fmt.Sprintf("📱 %s   %s   %s", slot.Number, statusStr, rateStr)
                if timerStr != "" {
                        line += fmt.Sprintf("   %s", timerStr)
                }

                if i < len(s.slots)-1 {
                        b.WriteString(style.CaptionStyle.Render(line) + "\n")
                } else {
                        b.WriteString(style.CaptionStyle.Render(line))
                }
        }

        return b.String()
}

// renderNicheQueue renders per-niche queue groups with tab switching.
func (s *Send) renderNicheQueue() string {
        var b strings.Builder

        // Collect items for the current niche tab.
        activeNiche := ""
        if s.nicheTab < len(s.nicheNames) {
                activeNiche = s.nicheNames[s.nicheTab]
        }

        // Group queue items by niche.
        groups := s.groupQueueByNiche()

        for niche, items := range groups {
                isActive := niche == activeNiche

                // Niche header with queue count.
                indicator := "  "
                if isActive {
                        indicator = "▸ "
                }
                header := fmt.Sprintf("%s%s (%d %s, %s ✓)",
                        indicator, niche, len(items), i18n.T(i18n.KeySendQueueCount), i18n.T(i18n.KeySendWAValidated))

                if isActive {
                        b.WriteString(style.BodyStyle.Render(header))
                } else {
                        b.WriteString(style.MutedStyle.Render(header))
                }
                b.WriteString("\n")

                // Render items only for the active niche.
                if isActive {
                        for _, item := range items {
                                b.WriteString(s.renderQueueItem(item))
                        }
                }
        }

        return b.String()
}

// renderQueueItem renders a single queue item line.
func (s *Send) renderQueueItem(item SendQueueItem) string {
        var b strings.Builder

        // Index with active indicator.
        indexStr := fmt.Sprintf("%02d", item.Index)
        activeArrow := "  "
        if item.Status == QueueItemSending {
                activeArrow = "→ "
        }

        // Line 1: index + arrow + name + slot on same line (slot right-aligned).
        nameRendered := style.MutedStyle.Render(item.Name)
        if item.Status == QueueItemSending {
                nameRendered = style.BodyStyle.Render(item.Name)
        }
        slotRendered := style.CaptionStyle.Render("📱 " + item.Slot)

        line1Left := style.CaptionStyle.Render(indexStr+activeArrow) + nameRendered
        line1 := lipgloss.NewStyle().Width(s.width).Render(
                lipgloss.JoinHorizontal(lipgloss.Top, line1Left, slotRendered),
        )
        b.WriteString(line1)

        // Template + variant with rotation tag (indented on next line).
        templateStr := fmt.Sprintf("     %s: %s", item.TemplateType, item.TemplateVar)
        if item.Rotated {
                templateStr += "  " + style.WarningStyle.Render("← " + i18n.T(i18n.KeySendRotation) + "!")
        }
        b.WriteString("\n" + style.CaptionStyle.Render(templateStr))

        // Status line.
        switch item.Status {
        case QueueItemSending:
                b.WriteString("\n" + style.CaptionStyle.Render("     ━━━━━━━━━━━━━━━━━ " + i18n.T(i18n.KeySendSending)))
        case QueueItemSent:
                b.WriteString("\n" + style.CaptionStyle.Render("     ✓ " + i18n.T(i18n.KeyStatusSent)))
        case QueueItemWaiting:
                waitStr := fmt.Sprintf("     %s (%s: %s)", i18n.T(i18n.KeySendWaiting), i18n.T(i18n.KeySendNextLabel), item.NextAt)
                b.WriteString("\n" + style.CaptionStyle.Render(waitStr))
        }

        b.WriteString("\n")
        return b.String()
}

// renderQueueItems renders a flat list of queue items (used by paused/response views).
func (s *Send) renderQueueItems() string {
        var b strings.Builder

        items := s.queue
        if len(items) == 0 {
                return ""
        }

        for i, item := range items {
                if i > 0 {
                        b.WriteString("\n")
                }

                indexStr := fmt.Sprintf("%02d", item.Index)
                activeArrow := "  "
                if item.Status == QueueItemSending {
                        activeArrow = "→ "
                }

                line := fmt.Sprintf("%s%s  %s", indexStr, activeArrow, item.Name)

                switch item.Status {
                case QueueItemSending:
                        line += style.CaptionStyle.Render("  ━━━━━━━━━━━━━━━━━ " + i18n.T(i18n.KeySendSending))
                case QueueItemSent:
                        line += "  " + style.SuccessStyle.Render("✓ " + i18n.T(i18n.KeyStatusSent))
                case QueueItemWaiting:
                        line += "  " + style.CaptionStyle.Render(i18n.T(i18n.KeySendWaiting))
                }

                b.WriteString(line)
        }

        return b.String()
}

// renderRateStats renders the rate and daily statistics.
func (s *Send) renderRateStats() string {
        var b strings.Builder

        now := time.Now()
        opacities := s.breathing.Opacities(now)

        rateStr := fmt.Sprintf("%d/%d %s %s",
                s.stats.RateHour, s.stats.RateHourLimit, i18n.T(i18n.KeySendRate), i18n.T(i18n.KeyWordThisHour))

        slotCountStr := fmt.Sprintf("(%d %s)", s.stats.SlotCount, i18n.T(i18n.KeySendSlotCount))

        dailyStr := fmt.Sprintf("%s: %d/%d",
                i18n.T(i18n.KeySendToday), s.stats.DailySent, s.stats.DailyLimit)

        rateRendered := rateStr
        if len(opacities) > 0 {
                rateRendered = component.RenderBreathing(rateStr, opacities[0], style.Text, style.TextMuted)
        }

        b.WriteString(style.CaptionStyle.Render(rateRendered + " " + slotCountStr + " · " + dailyStr))

        // Next send time.
        if s.stats.NextSendTime != "" {
                nextStr := fmt.Sprintf("%s: %s (%s)", i18n.T(i18n.KeySendNextAt), s.stats.NextSendTime, s.stats.NextSendSlot)
                b.WriteString("\n" + style.CaptionStyle.Render(nextStr))
        }

        return b.String()
}

// renderTodayStats renders today's sent/respond/convert stats with breathing.
func (s *Send) renderTodayStats() string {
        var b strings.Builder

        now := time.Now()
        opacities := s.breathing.Opacities(now)

        sentStr := fmt.Sprintf("%s: %d", i18n.T(i18n.KeySendTodaySent), s.stats.DailySent)
        respStr := fmt.Sprintf("%s: %d", i18n.T(i18n.KeySendTodayResp), s.stats.DailyRespond)
        convStr := fmt.Sprintf("%s: %d", i18n.T(i18n.KeySendTodayConv), s.stats.DailyConvert)

        if len(opacities) > 0 {
                sentStr = component.RenderBreathing(sentStr, opacities[0], style.Text, style.TextMuted)
        }
        if len(opacities) > 1 {
                respStr = component.RenderBreathing(respStr, opacities[1], style.Warning, style.TextMuted)
        }
        if len(opacities) > 2 {
                convStr = component.RenderBreathing(convStr, opacities[2], style.Success, style.TextMuted)
        }

        b.WriteString(style.CaptionStyle.Render(sentStr))
        b.WriteString("\n" + style.CaptionStyle.Render(respStr))
        b.WriteString("\n" + style.CaptionStyle.Render(convStr))

        return b.String()
}

// renderRemainingByNiche shows remaining queue count per niche.
func (s *Send) renderRemainingByNiche() string {
        counts := s.countQueueByNiche()
        parts := make([]string, 0, len(counts))
        for niche, count := range counts {
                parts = append(parts, fmt.Sprintf("%d %s %s", count, i18n.T(i18n.KeyWordFrom), niche))
        }
        return style.CaptionStyle.Render("(" + strings.Join(parts, " · ") + ")")
}

// renderResponseOverlay renders the incoming response interrupt banner.
func (s *Send) renderResponseOverlay() string {
        if s.response == nil {
                return ""
        }

        var b strings.Builder

        b.WriteString(style.SuccessStyle.Render(fmt.Sprintf("💬  %s %s", i18n.T(i18n.KeySendResponseIn), s.response.Name)))
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(style.BodyStyle.Render(fmt.Sprintf("\"%s\"", s.response.Message)))
        b.WriteString(style.Section(style.SubSectionGap))

        // Actions.
        b.WriteString(style.ActionStyle.Render("↵  " + i18n.T(i18n.KeySendSendOffer)))
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(style.CaptionStyle.Render("2  " + i18n.T(i18n.KeySendReplyCustom)))
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(style.CaptionStyle.Render("s  " + i18n.T(i18n.KeyWordLater)))

        return b.String()
}

// ---------------------------------------------------------------------------
// Key-hint renderers
// ---------------------------------------------------------------------------

// renderActiveHints renders key hints for the active sending state.
func (s *Send) renderActiveHints() string {
        return style.CaptionStyle.Render(
                fmt.Sprintf("p  %s    ↵  %s    tab  %s    q  %s",
                        i18n.T(i18n.KeyLabelPause),
                        i18n.T(i18n.KeySendSkipWait),
                        i18n.T(i18n.KeySendSwitchNiche),
                        i18n.T(i18n.KeyLabelBack),
                ),
        )
}

// renderPausedHints renders key hints for the paused state.
func (s *Send) renderPausedHints() string {
        return style.ActionStyle.Render(
                fmt.Sprintf("↵  %s    q  %s",
                        i18n.T(i18n.KeySendResume),
                        i18n.T(i18n.KeyLabelBack),
                ),
        )
}

// renderBackHint renders a simple back hint.
func (s *Send) renderBackHint() string {
        return style.CaptionStyle.Render("q  " + i18n.T(i18n.KeyLabelBack))
}

// ---------------------------------------------------------------------------
// Key-action dispatchers
// ---------------------------------------------------------------------------

// handleEnter dispatches the enter key based on the current state.
func (s *Send) handleEnter() {
        switch s.state {
        case protocol.SendPaused:
                // Resume sending.
                s.publishAction(string(protocol.ActionTogglePause))
        case protocol.SendWithResponse:
                // Approve response — send offer.
                s.publishAction(string(protocol.ActionApproveResponse))
        default:
                // Skip wait on active sending.
                s.publishAction(string(protocol.ActionSkipWait))
        }
}

// handleKey1 dispatches the "1" key based on the current state.
func (s *Send) handleKey1() {
        switch s.state {
        case protocol.SendOffHours:
                s.publishAction(string(protocol.ActionEmergencySend))
        case protocol.SendFailed:
                s.publishAction(string(protocol.ActionRetryManual))
        case protocol.SendAllSlotsDown:
                s.publishAction(string(protocol.ActionReLogin))
        default:
                s.publishKeyPress("1")
        }
}

// handleKey2 dispatches the "2" key based on the current state.
func (s *Send) handleKey2() {
        switch s.state {
        case protocol.SendWithResponse:
                s.publishAction(string(protocol.ActionCustomReply))
        case protocol.SendAllSlotsDown:
                s.publishAction(string(protocol.ActionLoginOneByOne))
        default:
                s.publishKeyPress("2")
        }
}

// handleSkip dispatches the "s" key based on the current state.
func (s *Send) handleSkip() {
        switch s.state {
        case protocol.SendActive, protocol.SendWithResponse:
                s.publishAction(string(protocol.ActionSkipWait))
        case protocol.SendFailed:
                s.publishAction(string(protocol.ActionSkipFailure))
        default:
                s.publishKeyPress("s")
        }
}

// cycleNicheTab advances to the next niche tab.
func (s *Send) cycleNicheTab() {
        if len(s.nicheNames) == 0 {
                return
        }
        s.nicheTab = (s.nicheTab + 1) % len(s.nicheNames)
}

// ---------------------------------------------------------------------------
// Bus helpers
// ---------------------------------------------------------------------------

// publishKeyPress publishes a key_press event to the backend.
func (s *Send) publishKeyPress(k string) {
        if s.Bus() != nil {
                s.Bus().Publish(bus.KeyPressMsg{Key: k, Screen: s.ID()})
        }
}

// publishAction publishes an action event to the backend.
func (s *Send) publishAction(action string) {
        if s.Bus() != nil {
                s.Bus().Publish(bus.ActionMsg{Action: action, Screen: s.ID()})
        }
}

// ---------------------------------------------------------------------------
// Data helpers
// ---------------------------------------------------------------------------

// formatSummaryHeader formats the right-aligned header summary.
func (s *Send) formatSummaryHeader() string {
        slotCount := len(s.slots)
        nicheCount := len(s.nicheNames)
        queueCount := s.stats.QueueTotal

        return fmt.Sprintf("%d %s · %d %s · %d %s (%s ✓)",
                slotCount, i18n.T(i18n.KeySendSlotCount),
                nicheCount, i18n.T(i18n.KeySendNicheCount),
                queueCount, i18n.T(i18n.KeySendQueueCount),
                i18n.T(i18n.KeySendWAValidated),
        )
}

// groupQueueByNiche groups queue items by their niche field.
func (s *Send) groupQueueByNiche() map[string][]SendQueueItem {
        groups := make(map[string][]SendQueueItem)
        for _, item := range s.queue {
                niche := item.Niche
                if niche == "" {
                        niche = "default"
                }
                groups[niche] = append(groups[niche], item)
        }
        return groups
}

// countQueueByNiche returns the count of waiting items per niche.
func (s *Send) countQueueByNiche() map[string]int {
        counts := make(map[string]int)
        for _, item := range s.queue {
                niche := item.Niche
                if niche == "" {
                        niche = "default"
                }
                counts[niche]++
        }
        return counts
}

// ---------------------------------------------------------------------------
// Parsing helpers — convert any→domain types from backend updates.
// ---------------------------------------------------------------------------

func parseSlotInfos(raw []any) []SlotInfo {
        slots := make([]SlotInfo, 0, len(raw))
        for _, v := range raw {
                m, ok := v.(map[string]any)
                if !ok {
                        continue
                }
                slot := SlotInfo{
                        Number:   anyMapString(m, protocol.ParamSlotNumber),
                        Status:   anyMapString(m, protocol.ParamSlotStatus),
                        CooldownLeft: anyMapString(m, protocol.ParamSlotCooldown),
                        ReadyIn:  anyMapString(m, protocol.ParamSlotReadyInAlt),
                        SentThisHour: int64(anyMapInt(m, protocol.ParamRateHour)),
                        HourLimit:    int64(anyMapInt(m, protocol.ParamRateMax)),
                }
                slots = append(slots, slot)
        }
        return slots
}

func parseSendQueueItems(raw []any) []SendQueueItem {
        items := make([]SendQueueItem, 0, len(raw))
        for _, v := range raw {
                m, ok := v.(map[string]any)
                if !ok {
                        continue
                }
                item := SendQueueItem{
                        Index:        int64(anyMapInt(m, protocol.ParamIndex)),
                        Name:     anyMapString(m, protocol.ParamName),
                        Slot:         anyMapString(m, protocol.ParamSlot),
                        TemplateType: anyMapString(m, protocol.ParamTemplate),
                        TemplateVar:  anyMapString(m, protocol.ParamVariant),
                        Status:   anyMapString(m, protocol.ParamSlotStatus),
                        NextAt:       anyMapString(m, protocol.ParamNextIn),
                        Niche:    anyMapString(m, protocol.ParamNiche),
                        Rotated:  anyMapBool(m, protocol.ParamRotated),
                }
                items = append(items, item)
        }
        return items
}

func parseSendStats(raw map[string]any) SendStats {
        return SendStats{
                RateHour:      int64(anyMapInt(raw, protocol.ParamRateHourSent)),
                RateHourLimit: int64(anyMapInt(raw, protocol.ParamRateHourMax)),
                SlotCount:     int64(anyMapInt(raw, protocol.ParamRateSlotCount)),
                DailySent:     int64(anyMapInt(raw, protocol.ParamDailySent)),
                DailyRespond:  int64(anyMapInt(raw, protocol.ParamDailyResp)),
                DailyConvert:  int64(anyMapInt(raw, protocol.ParamDailyConv)),
                DailyLimit:    int64(anyMapInt(raw, protocol.ParamDailyMax)),
                QueueTotal:    int64(anyMapInt(raw, protocol.ParamQueueCount)),
                NextSendTime:  anyMapString(raw, protocol.ParamNextSendTime),
                NextSendSlot:  anyMapString(raw, protocol.ParamNextSendSlot),
                Now:           anyMapString(raw, protocol.ParamNow),
        }
}

func parseSendFailure(raw map[string]any) *SendFailure {
        if raw == nil {
                return nil
        }
        return &SendFailure{
                Name:   anyMapString(raw, protocol.ParamName),
                Reason: anyMapString(raw, protocol.ParamReason),
                Hint:   anyMapString(raw, protocol.ParamHint),
        }
}

func parseResponseInterrupt(raw map[string]any) *ResponseInterrupt {
        if raw == nil {
                return nil
        }
        return &ResponseInterrupt{
                Name:    anyMapString(raw, protocol.ParamName),
                Message: anyMapString(raw, protocol.ParamMessage),
        }
}

func parseStringSlice(raw []any) []string {
        out := make([]string, 0, len(raw))
        for _, v := range raw {
                if s, ok := v.(string); ok {
                        out = append(out, s)
                }
        }
        return out
}

// anyString is an alias for anyMapString for backward compatibility.
func anyString(m map[string]any, key string) string {
        return anyMapString(m, key)
}

// anyInt is an alias for anyMapInt for backward compatibility.
func anyInt(m map[string]any, key string) int {
        return anyMapInt(m, key)
}

// anyBool is an alias for anyMapBool for backward compatibility.
func anyBool(m map[string]any, key string) bool {
        return anyMapBool(m, key)
}

// ---------------------------------------------------------------------------
// WA status renderer
// ---------------------------------------------------------------------------

// renderWAStatus renders a WhatsApp slot status indicator with
// appropriate color styling.
func renderWAStatus(status string) string {
        switch status {
        case SlotStatusActive:
                return style.SuccessStyle.Render(i18n.T(i18n.KeyStatusActive))
        case SlotStatusCooldown:
                return style.WarningStyle.Render(i18n.T(i18n.KeyStatusCooldown))
        case SlotStatusDown:
                return style.DangerStyle.Render("✗ " + i18n.T(i18n.KeyWordDisconnected))
        default:
                return style.CaptionStyle.Render("○ " + status)
        }
}

// ---------------------------------------------------------------------------
// Utility helpers
// ---------------------------------------------------------------------------

// fillGap returns a string of spaces that fills the gap between left and
// right content in a horizontal join.
func fillGap(totalWidth, usedWidth int) string {
        gap := totalWidth - usedWidth
        if gap < 1 {
                gap = 1
        }
        return strings.Repeat(" ", gap)
}

// minInt returns the smaller of a and b.
func minInt(a, b int) int {
        if a < b {
                return a
        }
        return b
}

// maxInt returns the larger of a and b.
func maxInt(a, b int) int {
        if a > b {
                return a
        }
        return b
}

// ---------------------------------------------------------------------------
// Animation ticker — triggers periodic redraws for breathing effects.
// ---------------------------------------------------------------------------

// breathingTick is a tea.Msg that triggers a redraw to update breathing
// animations. Screens that embed BreathingGroup should return a
// BreathingTickCmd from Update when they need continuous animation.
type breathingTick struct{}

// BreathingTickCmd returns a command that sends a breathingTick message
// after the specified interval.
func BreathingTickCmd(interval time.Duration) tea.Cmd {
        return tea.Tick(interval, func(t time.Time) tea.Msg {
                return breathingTick{}
        })
}
