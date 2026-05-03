// Package monitor implements the Monitor (dashboard) and Response screens.
// Agent 3D: Monitor Screens — Screen 7 (Dashboard) and Screen 8 (Response).
//
// Doc source: doc/05-screens-monitor-response.md
// States: ~17 (6 dashboard + 11 response)
package monitor

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
        "github.com/charmbracelet/bubbles/key"
        tea "github.com/charmbracelet/bubbletea"
        "github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Data models — populated by HandleNavigate / HandleUpdate from backend
// ---------------------------------------------------------------------------

// WASlot represents a WhatsApp number slot in the rotator.
type WASlot struct {
        Label    string // e.g. "slot-1"
        Number   string // e.g. "0812-xxxx-3456"
        Active   bool   // true = ● aktif, false = ○ cooldown
        Hours    string // e.g. "4/6 jam" or "ready: 3m"
}

// WorkerRow represents a single niche worker in the pool.
type WorkerRow struct {
        Name         string // e.g. "web_developer"
        Active       bool   // true = ●, false = ○
        Phase        string // e.g. "scraping", "sending", "idle"
        Queued       int    // antri count
        SentCount    int    // terkirim count (night view uses this instead of Queued)
        Responded    int    // respond count
        ConvertCount int    // convert count per niche (used in night view)
        Duration     string // e.g. "5j 37m" for idle view
        SendDur      string // e.g. "11m 23s"
}

// ActivityEvent represents a recent activity entry.
type ActivityEvent struct {
        Time     time.Time
        Niche    string // e.g. "[web_dev]"
        Business string // e.g. "kopi nusantara"
        Status   string // e.g. "respond", "terkirim"
        Detail   string // e.g. "iya kak, boleh lihat" or "──"
}

// PendingResponse represents an unhandled response in the pending view.
type PendingResponse struct {
        Niche    string // e.g. "[web_dev]"
        Business string // e.g. "kopi nusantara"
        Snippet  string // e.g. "iya kak, boleh lihat"
}

// DashboardData holds all data displayed by the Monitor dashboard.
type DashboardData struct {
        WASlots     []WASlot
        Workers     []WorkerRow
        Activities  []ActivityEvent
        Pending     []PendingResponse
        TodayStats  [4]int64 // leads_found, msgs_sent, responses, converts
        WeekStats   [4]int64 // same order
        ConvRate    string   // e.g. "5.1%"
        BestDay     string   // e.g. "selasa"
        BestTimeStr string   // e.g. "tuesday 10am" or "selasa jam 10" — backend is authoritative source
        NicheCount  int      // active niche count
        WANumCount  int      // connected WA numbers
        ActiveSlotCount int  // non-disconnected WA slots — backend is authoritative source
        ErrorSlot   string   // which slot disconnected (error state)
        WorkHours   string   // e.g. "09:00-17:00 wib"
        CurrentTime string   // e.g. "22:15" (night mode)
        AppName     string   // product name for dashboard headers (from backend)
}

// ---------------------------------------------------------------------------
// Dashboard screen — Screen 7: Monitor → Command Center
// ---------------------------------------------------------------------------

// Dashboard is the main monitor dashboard screen. It implements the Screen
// interface and renders 6 states: live_dashboard, idle_background,
// dashboard_night, dashboard_error, dashboard_empty, and
// dashboard_with_pending_responses.
//
// This is the home base screen — the command center where the user monitors
// all activity. Visual spec from doc/05-screens-monitor-response.md.
type Dashboard struct {
        base   protocol.ScreenID
        bus    *bus.Bus
        state  protocol.StateID
        data   DashboardData
        width  int
        height int

        // Animation components
        rain        component.DataRain

        // Tick management
        lastTick     time.Time
        rainPaused   bool      // true after user interaction
        rainResumeAt time.Time // when to resume data rain after idle
}

// NewDashboard creates a new Monitor Dashboard screen.
func NewDashboard() *Dashboard {
        return &Dashboard{
                base:     protocol.ScreenMonitor,
                state:    protocol.MonitorLiveDashboard,
                rain:     component.NewDataRain(60),
                lastTick: time.Now(),
        }
}

// ID returns the screen identifier.
func (d *Dashboard) ID() protocol.ScreenID { return d.base }

// SetBus injects the event bus reference.
func (d *Dashboard) SetBus(b *bus.Bus) { d.bus = b }

// HandleNavigate processes a "navigate" command from the backend.
// The backend sends the desired state via params[protocol.ParamState].
func (d *Dashboard) HandleNavigate(params map[string]any) error {
        if stateStr, ok := params[protocol.ParamState].(string); ok {
                d.state = protocol.StateID(stateStr)
        }
        d.populateData(params)
        return nil
}

// HandleUpdate processes an "update" command from the backend.
func (d *Dashboard) HandleUpdate(params map[string]any) error {
        d.populateData(params)
        d.rain.Pause() // Pause data rain on interaction
        d.rainPaused = true
        d.rainResumeAt = time.Now().Add(anim.DataRainPauseTimeout) // Resume after idle timeout
        return nil
}

// Focus is called when this screen becomes active.
func (d *Dashboard) Focus() {}

// Blur is called when this screen is no longer active.
func (d *Dashboard) Blur() {}

// Init implements tea.Model.
func (d *Dashboard) Init() tea.Cmd {
        return tickDashboardCmd()
}

// Update implements tea.Model.
func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        now := time.Now()

        switch m := msg.(type) {
        case tea.WindowSizeMsg:
                d.width = m.Width
                d.height = m.Height
                if d.width > 0 {
                        d.rain = component.NewDataRain(d.width)
                }
                return d, nil

        case tea.KeyMsg:
                // Pause data rain on any key press
                d.rain.Pause()
                d.rainPaused = true
                d.rainResumeAt = time.Now().Add(anim.DataRainPauseTimeout) // Resume after idle timeout
                return d, d.handleKey(m)

        case dashboardTickMsg:
                d.lastTick = now
                // Resume data rain after 10s idle
                if d.rainPaused && !d.rainResumeAt.IsZero() && now.After(d.rainResumeAt) {
                        d.rainPaused = false
                        d.rainResumeAt = time.Time{}
                }
                if !d.rainPaused {
                        d.rain.Tick(now)
                }
                return d, tickDashboardCmd()
        }

        return d, nil
}

// handleKey processes key events for the dashboard screen.
func (d *Dashboard) handleKey(msg tea.KeyMsg) tea.Cmd {
        switch {
        case key.Matches(msg, KeyEnter):
                // Navigate to response or detail view
                publish(d.bus, bus.ActionMsg{
                        Action: protocol.ActionViewDetail,
                        Screen: protocol.ScreenMonitor,
                })
                return nil

        case key.Matches(msg, KeyRefresh):
                publish(d.bus, bus.ActionMsg{
                        Action: protocol.ActionRefresh,
                        Screen: protocol.ScreenMonitor,
                })
                return nil

        case key.Matches(msg, Key1):
                // Navigate to leads — backend decides target screen
                publish(d.bus, bus.ActionMsg{
                        Action: protocol.ActionNavLeadsDB,
                        Screen: protocol.ScreenMonitor,
                })
                return nil

        case key.Matches(msg, Key2):
                // Navigate to response/messages — backend decides target screen
                publish(d.bus, bus.ActionMsg{
                        Action: protocol.ActionNavResponse,
                        Screen: protocol.ScreenMonitor,
                })
                return nil

        case key.Matches(msg, Key3):
                // Navigate to workers — backend decides target screen
                publish(d.bus, bus.ActionMsg{
                        Action: protocol.ActionNavWorkers,
                        Screen: protocol.ScreenMonitor,
                })
                return nil

        case key.Matches(msg, Key4):
                // Navigate to template manager — backend decides target screen
                publish(d.bus, bus.ActionMsg{
                        Action: protocol.ActionNavTemplate,
                        Screen: protocol.ScreenMonitor,
                })
                return nil

        case key.Matches(msg, Key5):
                // Navigate to anti-ban — backend decides target screen
                publish(d.bus, bus.ActionMsg{
                        Action: protocol.ActionNavAntiBan,
                        Screen: protocol.ScreenMonitor,
                })
                return nil

        case key.Matches(msg, Key6):
                // Navigate to follow-up — backend decides target screen
                publish(d.bus, bus.ActionMsg{
                        Action: protocol.ActionNavFollowUp,
                        Screen: protocol.ScreenMonitor,
                })
                return nil

        case key.Matches(msg, Key7):
                // Navigate to settings — backend decides target screen
                publish(d.bus, bus.ActionMsg{
                        Action: protocol.ActionNavSettings,
                        Screen: protocol.ScreenMonitor,
                })
                return nil

        case key.Matches(msg, KeySkip):
                publish(d.bus, bus.ActionMsg{
                        Action: protocol.ActionScrapeAll,
                        Screen: protocol.ScreenMonitor,
                })
                return nil

        case key.Matches(msg, KeyNerdStats):
                publish(d.bus, bus.ActionMsg{
                        Action: protocol.ActionToggleNerdStats,
                        Screen: protocol.ScreenMonitor,
                })
                return nil
        }

        return nil
}

// View implements tea.Model.
func (d *Dashboard) View() string {
        now := time.Now()

        switch d.state {
        case protocol.MonitorLiveDashboard:
                return d.viewLiveDashboard(now)
        case protocol.MonitorIdleBackground:
                return d.viewIdleBackground(now)
        case protocol.MonitorNight:
                return d.viewNight(now)
        case protocol.MonitorError:
                return d.viewError(now)
        case protocol.MonitorEmpty:
                return d.viewEmpty(now)
        case protocol.MonitorPendingResponses:
                return d.viewPendingResponses(now)
        default:
                return d.viewLiveDashboard(now)
        }
}

// ---------------------------------------------------------------------------
// State views — each renders one of the 6 dashboard states
// ---------------------------------------------------------------------------

// viewLiveDashboard renders the full command center with data rain, WA rotator,
// worker pool, stats, and activity timeline.
func (d *Dashboard) viewLiveDashboard(now time.Time) string {
        var b strings.Builder

        // Header line: "waclaw   ● wa nyambung (3 nomor) · 3 niche aktif"
        statusText := renderWAStatus(d.data.WANumCount)
        nicheText := fmt.Sprintf("%d %s", d.data.NicheCount, i18n.T(i18n.KeyMonitorNicheActive))
        headerRight := fmt.Sprintf("%s · %s", statusText, nicheText)
        d.renderHeader(&b, d.appName(), headerRight)
        b.WriteString("\n")

        // Data rain line
        if !d.rainPaused {
                b.WriteString(d.rain.ViewFormatted())
        }
        b.WriteString("\n\n")

        // Separator
        d.renderSeparator(&b)
        b.WriteString("\n")

        // WA Rotator section
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Bold(true).Render(
                i18n.T(i18n.KeyMonitorWARotator),
        ))
        b.WriteString("\n\n")
        for _, slot := range d.data.WASlots {
                d.renderWASlot(&b, slot)
        }
        b.WriteString("\n\n")

        // Worker Pool section
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Bold(true).Render(
                i18n.T(i18n.KeyMonitorWorkerPool),
        ))
        b.WriteString("\n\n")
        for _, w := range d.data.Workers {
                d.renderWorker(&b, w)
        }
        b.WriteString("\n\n")

        // Separator
        d.renderSeparator(&b)
        b.WriteString("\n")

        // Stats grid: hari ini / minggu ini
        d.renderStatsGrid(&b, now)
        b.WriteString("\n\n")

        // Separator
        d.renderSeparator(&b)
        b.WriteString("\n")

        // Activity timeline
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Bold(true).Render(
                i18n.T(i18n.KeyMonitorRecentActivityAll),
        ))
        b.WriteString("\n\n")
        d.renderActivityTimeline(&b)

        // Footer: nav bar + action keys
        b.WriteString("\n\n")
        d.renderSeparator(&b)
        b.WriteString("\n\n")
        d.renderNavBar(&b)
        b.WriteString("\n\n")
        d.renderActionKeys(&b)

        return b.String()
}

// viewIdleBackground renders the auto-pilot idle view.
func (d *Dashboard) viewIdleBackground(now time.Time) string {
        var b strings.Builder

        // Header
        statusText := lipgloss.NewStyle().Foreground(style.Success).Render("●") + " " + i18n.T(i18n.KeyMonitorAutoPilot)
        headerRight := fmt.Sprintf("%s · %d %s", statusText, d.data.NicheCount, i18n.T(i18n.KeyMonitorNicheActive))
        d.renderHeader(&b, d.appName(), headerRight)
        b.WriteString("\n")

        // Data rain
        if !d.rainPaused {
                b.WriteString(d.rain.ViewFormatted())
        }
        b.WriteString("\n\n")

        // Separator
        d.renderSeparator(&b)
        b.WriteString("\n")

        // Idle message
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                i18n.T(i18n.KeyMonitorIdleBody),
        ))
        b.WriteString("\n\n")

        // Worker summary
        for _, w := range d.data.Workers {
                phaseLabel := renderStatusDot(w.Active)
                line := fmt.Sprintf("%s %s   %s %s · %s %s",
                        style.Indent(style.IndentPerLevel),
                        lipgloss.NewStyle().Foreground(style.Text).Render(w.Name),
                        phaseLabel,
                        lipgloss.NewStyle().Foreground(style.TextMuted).Render(w.Phase),
                        lipgloss.NewStyle().Foreground(style.TextMuted).Render(w.Duration),
                        lipgloss.NewStyle().Foreground(style.TextDim).Render(w.SendDur),
                )
                b.WriteString(line)
                b.WriteString("\n")
        }
        b.WriteString("\n")

        // Summary line — doc: wa ● nyambung (3 nomor) / antrian total / response hari ini
        waStatus := renderWAStatus(d.data.WANumCount)
        b.WriteString(fmt.Sprintf("%s %s\n", style.Indent(style.IndentPerLevel), waStatus))
        b.WriteString(fmt.Sprintf("%s %s %d %s\n",
                style.Indent(style.IndentPerLevel),
                lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyMonitorQueuedTotal)),
                d.data.TodayStats[1],
                lipgloss.NewStyle().Foreground(style.TextDim).Render(i18n.T(i18n.KeyMonitorMsgsSent)),
        ))
        b.WriteString(fmt.Sprintf("%s %s %d\n",
                style.Indent(style.IndentPerLevel),
                lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyMonitorResponsesToday)),
                d.data.TodayStats[2],
        ))
        b.WriteString("\n")

        // Encouragement
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                i18n.T(i18n.KeyMonitorCanMinimize),
        ))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf(i18n.T(i18n.KeyMonitorWorkersDesc), d.data.NicheCount, d.data.WANumCount),
        ))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                i18n.T(i18n.KeyMonitorWillNotify),
        ))
        b.WriteString("\n\n")

        // Action keys
        b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render("↵"))
        b.WriteString("  ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyMonitorViewDetail)))
        b.WriteString("    ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render("q"))
        b.WriteString("  ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyLabelQuit)))

        return b.String()
}

// viewNight renders the night mode dashboard.
func (d *Dashboard) viewNight(now time.Time) string {
        var b strings.Builder

        // Header
        statusText := lipgloss.NewStyle().Foreground(style.TextDim).Render("○") + " " + lipgloss.NewStyle().Foreground(style.TextDim).Render(i18n.T(i18n.KeyMonitorNightMode))
        headerRight := fmt.Sprintf("%s · %d %s", statusText, d.data.NicheCount, i18n.T(i18n.KeyMonitorNicheActive))
        d.renderHeader(&b, d.appName(), headerRight)
        b.WriteString("\n\n")

        // Separator
        d.renderSeparator(&b)
        b.WriteString("\n")

        // Time info — use data field for work hours from backend.
        // No hardcoded fallback — if backend hasn't provided work hours yet,
        // show a generic placeholder. The backend is the authoritative source.
        workHours := d.data.WorkHours
        if workHours == "" {
                workHours = i18n.T(i18n.KeyMonitorDefaultWorkHours) // Generic placeholder — backend will override
        }
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("%s %s", i18n.T(i18n.KeyMonitorWorkHours), workHours),
        ))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(
                fmt.Sprintf("%s %s", i18n.T(i18n.KeyMonitorNow), d.data.CurrentTime),
        ))
        b.WriteString("\n\n")

        // Status
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyMonitorSenderPaused)))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyMonitorScraperRunning)))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf(i18n.T(i18n.KeyMonitorAutoResumeDetail), d.data.NicheCount, d.data.WorkHours),
        ))
        b.WriteString("\n\n")

        // Day summary
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Bold(true).Render(
                i18n.T(i18n.KeyMonitorDaySummary),
        ))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(
                renderDaySummary(d.data.TodayStats, ""),
        ))
        b.WriteString("\n\n")

        // Per niche
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Bold(true).Render(
                i18n.T(i18n.KeyMonitorPerNiche),
        ))
        b.WriteString("\n")
        // Doc: per niche: web_developer     28 terkirim · 5 respond · 2 convert
        for _, w := range d.data.Workers {
                sentCount := w.SentCount // Backend provides dedicated sent count for night view
                if sentCount == 0 && w.Queued > 0 {
                        sentCount = w.Queued // Fallback: older backend may still use Queued for both
                }
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                        fmt.Sprintf("%s   %d %s · %d %s · %d %s",
                                w.Name, sentCount, i18n.T(i18n.KeyMonitorMsgsSent),
                                w.Responded, i18n.T(i18n.KeyMonitorResponses), w.ConvertCount, i18n.T(i18n.KeyMonitorConverts)),
                ))
                b.WriteString("\n")
        }
        b.WriteString("\n")

        // Tips
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(fmt.Sprintf(i18n.T(i18n.KeyMonitorBestTime), d.data.BestTimeStr))) // Backend is authoritative source
        b.WriteString("\n\n")

        // Action keys
        b.WriteString(" ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render("↵"))
        b.WriteString("  ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyMonitorViewDetail)))
        b.WriteString("    ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render("q"))
        b.WriteString("  ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyLabelQuit)))

        return b.String()
}

// viewError renders the error state dashboard.
func (d *Dashboard) viewError(now time.Time) string {
        var b strings.Builder

        // Header
        statusText := lipgloss.NewStyle().Foreground(style.Danger).Render("✗") + " " + lipgloss.NewStyle().Foreground(style.Danger).Render(i18n.T(i18n.KeyMonitorProblemStatus))
        headerRight := fmt.Sprintf("%s · %d %s", statusText, d.data.NicheCount, i18n.T(i18n.KeyMonitorNicheActive))
        d.renderHeader(&b, d.appName(), headerRight)
        b.WriteString("\n\n")

        // Separator
        d.renderSeparator(&b)
        b.WriteString("\n\n")

        // Error details — Doc: ✗  wa putus (slot-1) — slot-2 & slot-3 tetap jalan
        b.WriteString(lipgloss.NewStyle().Foreground(style.Danger).Render("✗"))
        b.WriteString("  ")
        slotName := d.data.ErrorSlot
        // Backend must always provide error_slot in the error state.
        // If missing, the error view still renders but without a specific slot name.
        if slotName != "" {
                b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(
                        fmt.Sprintf(i18n.T(i18n.KeyMonitorSlotDisconnected), slotName),
                ))
        } else {
                b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(
                        i18n.T(i18n.KeyMonitorWADisconnected),
                ))
        }
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Success).Render("●"))
        b.WriteString("  ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyMonitorScraperOK)))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Success).Render("●"))
        b.WriteString("  ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyMonitorDatabaseOK)))
        b.WriteString("\n\n")

        // Reassurance
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(fmt.Sprintf(i18n.T(i18n.KeyMonitorSlotsActive), d.data.ActiveSlotCount))) // Backend is authoritative source
        b.WriteString("\n")
        // Doc: hanya slot-1 yang pending sampai nyambung.
        if slotName != "" {
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                        fmt.Sprintf(i18n.T(i18n.KeyMonitorSlotPending), slotName),
                ))
        }
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyMonitorAutoReconnect)))
        b.WriteString("\n\n")

        // Action keys
        b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render("1"))
        b.WriteString("  ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyMonitorRelogin)))
        b.WriteString("    ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render("r"))
        b.WriteString("  ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyMonitorCheckStatus)))
        b.WriteString("    ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render("q"))
        b.WriteString("  ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyLabelQuit)))

        return b.String()
}

// viewEmpty renders the empty/no-data state.
func (d *Dashboard) viewEmpty(now time.Time) string {
        var b strings.Builder

        // Header
        statusText := lipgloss.NewStyle().Foreground(style.Success).Render("● " + i18n.T(i18n.KeyMonitorWAConnected))
        d.renderHeader(&b, d.appName(), statusText)
        b.WriteString("\n\n")

        // Separator
        d.renderSeparator(&b)
        b.WriteString("\n\n")

        // Start options
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyMonitorNoData)))
        b.WriteString("\n\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render("1"))
        b.WriteString("  ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(i18n.T(i18n.KeyMonitorStartScrape)))
        b.WriteString("    ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render("2"))
        b.WriteString("  ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(i18n.T(i18n.KeyMonitorSetupNiche)))
        b.WriteString("    ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render("3"))
        b.WriteString("  ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(i18n.T(i18n.KeyMonitorStartSend)))
        b.WriteString("\n\n")

        // Separator
        d.renderSeparator(&b)
        b.WriteString("\n\n")

        // Tips
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(i18n.T(i18n.KeyMonitorTipStartScrape)))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(i18n.T(i18n.KeyMonitorTipAuto)))

        return b.String()
}

// viewPendingResponses renders the dashboard with pending responses overlay.
func (d *Dashboard) viewPendingResponses(now time.Time) string {
        var b strings.Builder

        // Header
        statusText := renderWAStatus(d.data.WANumCount)
        headerRight := fmt.Sprintf("%s · %d %s", statusText, d.data.NicheCount, i18n.T(i18n.KeyMonitorNicheActive))
        d.renderHeader(&b, d.appName(), headerRight)
        b.WriteString("\n\n")

        // Separator
        d.renderSeparator(&b)
        b.WriteString("\n\n")

        // Pending count with amber flash
        pendingCount := len(d.data.Pending)
        b.WriteString(lipgloss.NewStyle().Foreground(style.Warning).Bold(true).Render(
                fmt.Sprintf("⚡ %d %s", pendingCount, i18n.T(i18n.KeyMonitorPendingUnreplied)),
        ))
        b.WriteString("\n\n")

        // Pending responses list
        for _, p := range d.data.Pending {
                nicheTag := lipgloss.NewStyle().Foreground(style.Accent).Render(p.Niche)
                bizName := lipgloss.NewStyle().Foreground(style.Text).Render(p.Business)
                snippet := lipgloss.NewStyle().Foreground(style.TextMuted).Render(p.Snippet)
                b.WriteString(fmt.Sprintf("%s %s     %s\n", nicheTag, bizName, snippet))
        }
        b.WriteString("\n")

        // Separator
        d.renderSeparator(&b)
        b.WriteString("\n\n")

        // Nav bar
        d.renderNavBar(&b)
        b.WriteString("\n\n")

        // Action keys
        b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render("↵"))
        b.WriteString("  ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyResponseProcessOne)))
        b.WriteString("    ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render("1"))
        b.WriteString("  ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyMonitorAutoOfferAll)))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render("2"))
        b.WriteString("  ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(i18n.T(i18n.KeyMonitorAutoPerNiche)))
        b.WriteString("    ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(i18n.T(i18n.KeyMonitorOfferPerNiche)))
        b.WriteString("\n\n")

        // Separator
        d.renderSeparator(&b)
        b.WriteString("\n\n")

        // Day summary
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyMonitorToday),
                        renderDaySummary(d.data.TodayStats, "")),
        ))

        return b.String()
}

// ---------------------------------------------------------------------------
// Render helpers — DRY sub-renderers used by multiple states
// ---------------------------------------------------------------------------

// defaultAppName is the generic fallback when the backend has not yet provided
// an app name. It is deliberately not the product name — the backend is the
// single source of truth for the display name via protocol.ParamAppName.
const defaultAppName = "App"

// appName returns the product name for dashboard headers, sourced from the backend.
// Falls back to a generic default only if the backend has not yet provided a name.
func (d *Dashboard) appName() string {
        if d.data.AppName != "" {
                return d.data.AppName
        }
        return defaultAppName
}

// renderSeparator renders a ── line in TextDim color.
func (d *Dashboard) renderSeparator(b *strings.Builder) {
        b.WriteString(renderDimSeparator(d.width))
}

// renderHeader renders the "<appName>   <rightText>" pattern with dynamic padding.
func (d *Dashboard) renderHeader(b *strings.Builder, leftText, rightText string) {
        left := lipgloss.NewStyle().Foreground(style.Text).Bold(true).Render(leftText)
        right := lipgloss.NewStyle().Foreground(style.TextMuted).Render(rightText)
        totalWidth := d.width
        if totalWidth < minRenderWidth {
                totalWidth = minRenderWidth
        }
        padLen := totalWidth - len(left) - 2
        if padLen < 0 {
                padLen = 0
        }
        b.WriteString(lipgloss.NewStyle().Width(totalWidth).Render(
                left + strings.Repeat(" ", padLen) + right,
        ))
        b.WriteString("\n")
}

// renderKeyHint renders a single key-label pair. The i18n label already
// includes the key prefix (e.g. "q keluar"), so we highlight the first
// character in accent and render the rest in text_muted to match the doc
// spec: "q keluar" where q is accent-colored and "keluar" is muted.
func (d *Dashboard) renderKeyHint(b *strings.Builder, key, label string) {
        // Label already contains the key prefix from i18n (e.g. "q keluar").
        // Render key in accent, rest in muted.
        if len(label) > 0 && string(label[0]) == key {
                b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render(key))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(label[1:]))
        } else {
                // Fallback for labels without key prefix
                b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render(key))
                b.WriteString("  ")
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(label))
        }
}
// renderWASlot renders a single WA slot row per doc spec:
// 📱 slot-1  0812-xxxx-3456   ● aktif   4/6 jam
func (d *Dashboard) renderWASlot(b *strings.Builder, slot WASlot) {
        icon := "📱"
        var status string
        if slot.Active {
                status = lipgloss.NewStyle().Foreground(style.Success).Render(
                        i18n.T(i18n.KeyStatusActive))
        } else {
                status = lipgloss.NewStyle().Foreground(style.TextDim).Render(
                        fmt.Sprintf("○ %s", i18n.T(i18n.KeyMonitorCooldown)))
        }
        hoursLabel := lipgloss.NewStyle().Foreground(style.TextDim).Render(slot.Hours)
        label := lipgloss.NewStyle().Foreground(style.TextMuted).Render(slot.Label)
        number := lipgloss.NewStyle().Foreground(style.Text).Render(slot.Number)
        b.WriteString(fmt.Sprintf("  %s %s  %s   %s   %s\n", icon, label, number, status, hoursLabel))
}

// renderWorker renders a single worker row.
func (d *Dashboard) renderWorker(b *strings.Builder, w WorkerRow) {
        statusDot := renderStatusDot(w.Active)
        name := lipgloss.NewStyle().Foreground(style.Text).Render(w.Name)
        phase := lipgloss.NewStyle().Foreground(style.TextMuted).Render(w.Phase)
        queued := lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("%d %s", w.Queued, i18n.T(i18n.KeyMonitorQueued)),
        )
        responded := lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("%d %s", w.Responded, i18n.T(i18n.KeyMonitorResponded)),
        )
        // Doc: "web_developer      ● scraping   24 antri   3 respond"
        b.WriteString(fmt.Sprintf("  %s   %s %s   %s   %s\n",
                name, statusDot, phase, queued, responded))
}

// renderStatsGrid renders the hari ini / minggu ini stat grid.
// TODO(issue-22): Apply breathing animation to stat numbers (requires component/stat_card.go changes).
// TODO(issue-23): "● wa nyambung" 3s pulse using anim.ConnectionPulseCycle needs component integration.
// TODO(issue-24): "⚡ response belum dibalas!" amber flash using anim.AttentionFlash needs integration.
func (d *Dashboard) renderStatsGrid(b *strings.Builder, now time.Time) {
        leftLabels := []string{
                i18n.T(i18n.KeyMonitorLeadsFound),
                i18n.T(i18n.KeyMonitorMsgsSent),
                i18n.T(i18n.KeyMonitorResponses),
                i18n.T(i18n.KeyMonitorConverts),
        }
        rightLabels := []string{
                i18n.T(i18n.KeyMonitorTotalLeads),
                i18n.T(i18n.KeyMonitorMsgsSent),
                i18n.T(i18n.KeyMonitorResponses),
                i18n.T(i18n.KeyMonitorConverts),
        }

        grid := component.NewStatGrid(leftLabels, rightLabels, d.data.TodayStats[:], d.data.WeekStats[:])
        grid.ColumnWidth = statColumnWidth

        // Column headers
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Bold(true).Render(
                i18n.T(i18n.KeyMonitorToday),
        ))
        padLen := statColumnWidth
        if d.width > 60 {
                padLen = d.width/2 - 10
        }
        b.WriteString(strings.Repeat(" ", padLen))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Bold(true).Render(
                i18n.T(i18n.KeyMonitorThisWeek),
        ))
        b.WriteString("\n\n")

        b.WriteString(grid.ViewAt(now))
        b.WriteString("\n\n")

        // Conversion rate + best day
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("%s %s         %s %s",
                        i18n.T(i18n.KeyMonitorConversionRate), d.data.ConvRate,
                        i18n.T(i18n.KeyMonitorBestDay), d.data.BestDay),
        ))
}

// renderActivityTimeline renders recent activity events.
func (d *Dashboard) renderActivityTimeline(b *strings.Builder) {
        for _, ev := range d.data.Activities {
                timeStr := lipgloss.NewStyle().Foreground(style.TextDim).Render(ev.Time.Format("15:04"))
                nicheTag := lipgloss.NewStyle().Foreground(style.Accent).Render(ev.Niche)
                bizName := lipgloss.NewStyle().Foreground(style.Text).Render(ev.Business)

                var statusStr string
                if ev.Status == string(protocol.ActivityStatusRespond) || ev.Status == i18n.T(i18n.KeyMonitorResponses) {
                        statusStr = lipgloss.NewStyle().Foreground(style.Warning).Render(ev.Status)
                } else {
                        statusStr = lipgloss.NewStyle().Foreground(style.TextMuted).Render(ev.Status)
                }

                detailStr := lipgloss.NewStyle().Foreground(style.TextDim).Render(ev.Detail)

                b.WriteString(fmt.Sprintf("  %s  %s  %s    %s    %s\n",
                        timeStr, nicheTag, bizName, statusStr, detailStr))
        }
}

// renderNavBar renders the 1-7 navigation bar.
func (d *Dashboard) renderNavBar(b *strings.Builder) {
        navItems := []struct {
                key   string
                label string
        }{
                {"1", i18n.T(i18n.KeyMonitorNavLeads)},
                {"2", i18n.T(i18n.KeyMonitorNavMessages)},
                {"3", i18n.T(i18n.KeyMonitorNavWorkers)},
                {"4", i18n.T(i18n.KeyMonitorNavTemplate)},
                {"5", i18n.T(i18n.KeyMonitorNavAntiban)},
                {"6", i18n.T(i18n.KeyMonitorNavFollowup)},
                {"7", i18n.T(i18n.KeyMonitorNavSettings)},
        }
        for i, item := range navItems {
                if i > 0 {
                        b.WriteString("    ")
                }
                d.renderKeyHint(b, item.key, item.label)
        }
}

// renderActionKeys renders the bottom action key hints.
func (d *Dashboard) renderActionKeys(b *strings.Builder) {
        keys := []struct {
                key   string
                label string
        }{
                {"r", i18n.T(i18n.KeyLabelRefresh)},
                {"s", i18n.T(i18n.KeyMonitorScrapeAll)},
                {"q", i18n.T(i18n.KeyLabelQuit)},
                {"`", i18n.T(i18n.KeyLabelNerd)},
        }
        for i, k := range keys {
                if i > 0 {
                        b.WriteString("    ")
                }
                d.renderKeyHint(b, k.key, k.label)
        }
}

// populateData extracts data from backend params into DashboardData.
func (d *Dashboard) populateData(params map[string]any) {
        if v, ok := params[protocol.ParamWASlots]; ok {
                d.data.WASlots = toWASlots(v)
        }
        if v, ok := params[protocol.ParamWorkers]; ok {
                d.data.Workers = toWorkers(v)
        }
        if v, ok := params[protocol.ParamActivities]; ok {
                d.data.Activities = toActivities(v)
        }
        if v, ok := params[protocol.ParamPending]; ok {
                d.data.Pending = toPendingResponses(v)
        }
        if v, ok := params[protocol.ParamTodayStats]; ok {
                d.data.TodayStats = toStatArray(v)
        }
        if v, ok := params[protocol.ParamWeekStats]; ok {
                d.data.WeekStats = toStatArray(v)
        }
        d.data.ConvRate = extractString(params, protocol.ParamConvRate)
        d.data.BestDay = extractString(params, protocol.ParamBestDay)
        d.data.BestTimeStr = extractString(params, protocol.ParamBestTimeStr) // Backend is authoritative source
        d.data.NicheCount = extractInt(params, protocol.ParamNicheCount)
        d.data.WANumCount = extractInt(params, protocol.ParamWANumCount)
        d.data.ActiveSlotCount = extractInt(params, protocol.ParamActiveSlotCount) // Backend is authoritative source
        d.data.ErrorSlot = extractString(params, protocol.ParamErrorSlot)
        d.data.WorkHours = extractString(params, protocol.ParamWorkHours)
        d.data.CurrentTime = extractString(params, protocol.ParamCurrentTime)
        if v := extractString(params, protocol.ParamAppName); v != "" {
                d.data.AppName = v
        }
}

// ---------------------------------------------------------------------------
// Tick message for periodic dashboard updates
// ---------------------------------------------------------------------------

// dashboardTickMsg is the internal tick message for the dashboard.
type dashboardTickMsg time.Time

// tickDashboardCmd returns a tea.Cmd that produces a tick at the data rain
// update interval (5s from anim.DataRainUpdateInterval).
func tickDashboardCmd() tea.Cmd {
        return tea.Tick(anim.DataRainUpdateInterval, func(t time.Time) tea.Msg {
                return dashboardTickMsg(t)
        })
}

// ---------------------------------------------------------------------------
// Param conversion helpers — DRY shared param extraction
// ---------------------------------------------------------------------------

func toWASlots(v any) []WASlot {
        slots, ok := v.([]WASlot)
        if ok {
                return slots
        }
        // Try []map[string]any from JSON
        items, ok := v.([]any)
        if !ok {
                return nil
        }
        var result []WASlot
        for _, item := range items {
                m, ok := item.(map[string]any)
                if !ok {
                        continue
                }
                s := WASlot{}
                if v, ok := m["label"].(string); ok {
                        s.Label = v
                }
                if v, ok := m["number"].(string); ok {
                        s.Number = v
                }
                if v, ok := m["active"].(bool); ok {
                        s.Active = v
                }
                if v, ok := m["hours"].(string); ok {
                        s.Hours = v
                }
                result = append(result, s)
        }
        return result
}

func toWorkers(v any) []WorkerRow {
        workers, ok := v.([]WorkerRow)
        if ok {
                return workers
        }
        items, ok := v.([]any)
        if !ok {
                return nil
        }
        var result []WorkerRow
        for _, item := range items {
                m, ok := item.(map[string]any)
                if !ok {
                        continue
                }
                w := WorkerRow{}
                if v, ok := m["name"].(string); ok {
                        w.Name = v
                }
                if v, ok := m["active"].(bool); ok {
                        w.Active = v
                }
                if v, ok := m["phase"].(string); ok {
                        w.Phase = v
                }
                if v, ok := m["queued"].(float64); ok {
                        w.Queued = int(v)
                }
                if v, ok := m["responded"].(float64); ok {
                        w.Responded = int(v)
                }
                if v, ok := m["convert_count"].(float64); ok {
                        w.ConvertCount = int(v)
                }
                if v := extractInt(m, "sent"); v > 0 {
                        w.SentCount = v
                }
                if v, ok := m["duration"].(string); ok {
                        w.Duration = v
                }
                if v, ok := m["send_dur"].(string); ok {
                        w.SendDur = v
                }
                result = append(result, w)
        }
        return result
}

func toActivities(v any) []ActivityEvent {
        events, ok := v.([]ActivityEvent)
        if ok {
                return events
        }
        items, ok := v.([]any)
        if !ok {
                return nil
        }
        var result []ActivityEvent
        for _, item := range items {
                m, ok := item.(map[string]any)
                if !ok {
                        continue
                }
                e := ActivityEvent{}
                if v, ok := m["niche"].(string); ok {
                        e.Niche = v
                }
                if v, ok := m["business"].(string); ok {
                        e.Business = v
                }
                if v, ok := m["status"].(string); ok {
                        e.Status = v
                }
                if v, ok := m["detail"].(string); ok {
                        e.Detail = v
                }
                result = append(result, e)
        }
        return result
}

func toPendingResponses(v any) []PendingResponse {
        pending, ok := v.([]PendingResponse)
        if ok {
                return pending
        }
        items, ok := v.([]any)
        if !ok {
                return nil
        }
        var result []PendingResponse
        for _, item := range items {
                m, ok := item.(map[string]any)
                if !ok {
                        continue
                }
                p := PendingResponse{}
                if v, ok := m["niche"].(string); ok {
                        p.Niche = v
                }
                if v, ok := m["business"].(string); ok {
                        p.Business = v
                }
                if v, ok := m["snippet"].(string); ok {
                        p.Snippet = v
                }
                result = append(result, p)
        }
        return result
}

func toStatArray(v any) [4]int64 {
        arr, ok := v.([4]int64)
        if ok {
                return arr
        }
        items, ok := v.([]any)
        if !ok || len(items) < 4 {
                return [4]int64{}
        }
        var result [4]int64
        for i, item := range items {
                switch n := item.(type) {
                case float64:
                        result[i] = int64(n)
                case int64:
                        result[i] = n
                case int:
                        result[i] = int64(n)
                }
        }
        return result
}
