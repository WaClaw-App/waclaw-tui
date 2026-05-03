package comms

import (
        "fmt"
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/internal/tui/component"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
        tea "github.com/charmbracelet/bubbletea"
        "github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// History data types
// ---------------------------------------------------------------------------

// HistoryEvent represents a single event in the history timeline.
type HistoryEvent struct {
        // Time is when the event occurred.
        Time time.Time

        // Icon is the emoji indicator: 💬 (response), 📤 (sent), ✅ (delivered),
        // ★ (jackpot), 🔍 (scrape), ⚡ (auto-pilot), 🎉 (converted).
        Icon string

        // BusinessName is the lead/business the event is about.
        BusinessName string

        // EventType is the event category label (e.g. "respond", "terkirim", "jackpot!").
        EventType string

        // Detail is the secondary info (e.g. response text, score).
        Detail string

        // Highlight indicates this event should stand out.
        Highlight bool

        // IsConversion indicates this is a conversion event (gold shimmer).
        IsConversion bool

        // Revenue is the deal revenue if converted.
        Revenue string
}

// DayStats holds aggregate statistics for a single day.
type DayStats struct {
        Sent     int
        Respond  int
        Convert  int
        NewLeads int
        Scrapes  int
}

// WeekData holds per-day data for the weekly chart view.
type WeekData struct {
        // Days contains 7 entries for each day of the week.
        Days []component.DayData

        // Messages is the daily send counts (parallel with Days).
        Messages []int64

        // Responses is the daily response counts (parallel with Days).
        Responses []int64

        // Converts is the daily conversion counts (parallel with Days).
        Converts []int64

        // BestDayIndex is the index of the best performing day.
        BestDayIndex int

        // BestDayLabel is the display name of the best day.
        BestDayLabel string

        // BestDayConvRate is the conversion rate string for the best day,
        // provided by the backend (e.g. "8.3%"). Empty if not available.
        // The TUI does not compute rates — that is a backend concern.
        BestDayConvRate string
}

// ---------------------------------------------------------------------------
// History Model
// ---------------------------------------------------------------------------

// History displays historical activity logs with timeline views and
// weekly mini charts.
//
// States (from doc/09-screens-communicate.md):
//   - HistoryToday: today's activity timeline
//   - HistoryWeek: weekly summary with mini bar charts
//   - HistoryDayDetail: specific past day detail
type History struct {
        screenBase
        state protocol.StateID
        bus   *bus.Bus

        // CurrentDate is the date being viewed (today or a selected past day).
        CurrentDate time.Time

        // TodayEvents holds today's timeline events.
        TodayEvents []HistoryEvent

        // TodayStats holds today's summary statistics.
        TodayStats DayStats

        // WeekData holds the weekly chart data.
        WeekData WeekData

        // WeekStats holds the aggregate weekly statistics.
        WeekStats DayStats

        // AvgResponseTime is the average response time string.
        AvgResponseTime string

        // DayDetailDate is the date shown in day_detail state.
        DayDetailDate time.Time

        // DayDetailEvents holds events for the selected day.
        DayDetailEvents []HistoryEvent

        // DayDetailStats holds stats for the selected day.
        DayDetailStats DayStats

        // Timeline is the reusable timeline component.
        Timeline component.Timeline

        // Width and Height for layout.
        Width  int
        Height int

        // CursorIndex tracks the selected item in list views.
        CursorIndex int

        // AnimStart tracks when the view animation started.
        AnimStart time.Time
}

// NewHistory creates a History screen in today state.
func NewHistory() *History {
        now := time.Now()
        return &History{
                screenBase:        screenBase{id: protocol.ScreenHistory},
                state:       protocol.HistoryToday,
                CurrentDate: now,
                Timeline:    component.NewTimeline(),
        }
}

// ---------------------------------------------------------------------------
// Screen interface
// ---------------------------------------------------------------------------

func (h *History) SetBus(b *bus.Bus) { h.bus = b }

func (h *History) Focus() {
        h.AnimStart = time.Now()
        h.CursorIndex = 0
}

func (h *History) Blur() {}

// ConsumesKey implements tui.KeyConsumer. History has sub-states (week,
// day_detail) where "q" should navigate back locally instead of popping the
// navigation stack.
func (h *History) ConsumesKey(msg tea.KeyMsg) bool {
        switch msg.String() {
        case "q":
                return h.state == protocol.HistoryWeek || h.state == protocol.HistoryDayDetail
        }
        return false
}

// HandleNavigate processes navigate commands from the backend.
func (h *History) HandleNavigate(params map[string]any) error {
        if state, ok := params[protocol.ParamState].(string); ok {
                h.state = protocol.StateID(state)
        }
        if dateStr, ok := params[protocol.ParamDate].(string); ok {
                if t, err := time.Parse("2006-01-02", dateStr); err == nil {
                        h.CurrentDate = t
                }
        }
        h.populateTimeline()
        return nil
}

// HandleUpdate processes update commands from the backend.
func (h *History) HandleUpdate(params map[string]any) error {
        // Parse events from update.
        if rawEvents, ok := params[protocol.ParamEvents].([]any); ok {
                h.TodayEvents = parseHistoryEvents(rawEvents)
                h.populateTimeline()
        }
        if rawStats, ok := params[protocol.ParamStats].(map[string]any); ok {
                h.TodayStats = parseDayStats(rawStats)
        }
        if rawWeek, ok := params[protocol.ParamWeek].(map[string]any); ok {
                h.WeekData = parseWeekData(rawWeek)
                h.WeekStats = parseWeekStats(rawWeek)
        }
        if avg, ok := params[protocol.ParamAvgResponseTime].(string); ok {
                h.AvgResponseTime = avg
        }
        return nil
}

// Init implements tea.Model.
func (h *History) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (h *History) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch m := msg.(type) {
        case tea.WindowSizeMsg:
                h.Width = m.Width
                h.Height = m.Height
                return h, nil

        case tea.KeyMsg:
                return h.handleKey(m)
        }

        // Advance timeline animation.
        h.Timeline.Tick(time.Now())
        return h, nil
}

// ---------------------------------------------------------------------------
// Key handling
// ---------------------------------------------------------------------------

func (h *History) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch h.state {
        case protocol.HistoryToday:
                return h.handleTodayKey(msg)
        case protocol.HistoryWeek:
                return h.handleWeekKey(msg)
        case protocol.HistoryDayDetail:
                return h.handleDayDetailKey(msg)
        }
        return h, nil
}

func (h *History) handleTodayKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // "w" → switch to week view.
        if msg.String() == "w" {
                h.state = protocol.HistoryWeek
                h.AnimStart = time.Now()
                h.CursorIndex = 0
                return h, nil
        }

        // "q" → back (let global handler pop screen).
        if msg.String() == "q" {
                return h, nil
        }

        switch msg.Type {
        case tea.KeyEnter:
                // View detail of selected event.
                if h.CursorIndex < len(h.TodayEvents) {
                        evt := h.TodayEvents[h.CursorIndex]
                        h.DayDetailDate = evt.Time
                        h.DayDetailEvents = []HistoryEvent{evt}
                        h.state = protocol.HistoryDayDetail
                        h.AnimStart = time.Now()
                }
                return h, nil

        case tea.KeyLeft:
                // Navigate to previous day.
                h.CurrentDate = h.CurrentDate.AddDate(0, 0, -1)
                h.AnimStart = time.Now()
                if h.bus != nil {
                        h.bus.Publish(bus.ActionMsg{
                                Action: string(protocol.ActionHistoryPrevDay),
                                Screen: protocol.ScreenHistory,
                                Params: map[string]any{protocol.ParamDate: h.CurrentDate.Format("2006-01-02")},
                        })
                }
                return h, nil

        case tea.KeyUp:
                if h.CursorIndex > 0 {
                        h.CursorIndex--
                }
                return h, nil

        case tea.KeyDown:
                if h.CursorIndex < len(h.TodayEvents)-1 {
                        h.CursorIndex++
                }
                return h, nil
        }
        return h, nil
}

func (h *History) handleWeekKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // "t" → switch back to today view.
        if msg.String() == "t" {
                h.state = protocol.HistoryToday
                h.AnimStart = time.Now()
                return h, nil
        }

        // "q" → back to today view.
        if msg.String() == "q" {
                h.state = protocol.HistoryToday
                h.AnimStart = time.Now()
                return h, nil
        }

        switch msg.Type {
        case tea.KeyUp:
                if h.CursorIndex > 0 {
                        h.CursorIndex--
                }
                return h, nil

        case tea.KeyDown:
                if h.CursorIndex < daysInWeek-1 {
                        h.CursorIndex++
                }
                return h, nil

        case tea.KeyEnter:
                // View detail of selected day.
                if h.CursorIndex < daysInWeek {
                        selectedDate := h.startOfWeek().AddDate(0, 0, h.CursorIndex)
                        h.DayDetailDate = selectedDate
                        h.state = protocol.HistoryDayDetail
                        h.AnimStart = time.Now()
                        if h.bus != nil {
                                h.bus.Publish(bus.ActionMsg{
                                        Action: string(protocol.ActionHistoryDayDetail),
                                        Screen: protocol.ScreenHistory,
                                        Params: map[string]any{protocol.ParamDate: selectedDate.Format("2006-01-02")},
                                })
                        }
                }
                return h, nil
        }
        return h, nil
}

func (h *History) handleDayDetailKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // "q" → back to today view.
        if msg.String() == "q" {
                h.state = protocol.HistoryToday
                h.AnimStart = time.Now()
                return h, nil
        }

        switch msg.Type {
        case tea.KeyLeft:
                // Previous day.
                h.DayDetailDate = h.DayDetailDate.AddDate(0, 0, -1)
                h.AnimStart = time.Now()
                return h, nil

        case tea.KeyRight:
                // Next day.
                h.DayDetailDate = h.DayDetailDate.AddDate(0, 0, 1)
                h.AnimStart = time.Now()
                return h, nil
        }
        return h, nil
}

// ---------------------------------------------------------------------------
// View rendering
// ---------------------------------------------------------------------------

// View renders the history screen.
func (h *History) View() string {
        switch h.state {
        case protocol.HistoryToday:
                return h.viewToday()
        case protocol.HistoryWeek:
                return h.viewWeek()
        case protocol.HistoryDayDetail:
                return h.viewDayDetail()
        default:
                return h.viewToday()
        }
}

// viewToday renders the history_today state.
//
// Spec:
//
//      hari ini                                14:23 wib
//      timeline:
//      14:23  💬  kopi nusantara      respond     "iya kak, boleh"
//      14:01  📤  wedding bliss        terkirim
//      ...
//      ringkasan:
//      terkirim  7    respond  4    convert  0
//      lead baru 67   scrape 1x
//      ↵ liat detail event    ← hari kemarin    q balik
func (h *History) viewToday() string {
        var b strings.Builder

        // Title with current time.
        title := i18n.T(i18n.KeyHistoryToday)
        timeStr := h.CurrentDate.Format("15:04")
        // Doc spec: ID locale shows "wib", EN locale shows the plain time.
        tzLabel := ""
        if i18n.GetLocale() == i18n.LocaleID {
                tzLabel = fmt.Sprintf("%s wib", timeStr)
        } else {
                tzLabel = timeStr
        }
        b.WriteString(renderTitleWithIndicator(title, tzLabel, style.CaptionStyle, h.Width))
        b.WriteString("\n\n")

        // Separator.
        writeSeparator(&b, h.Width)

        // Timeline label.
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyHistoryTimeline) + ":"))
        b.WriteString("\n\n")

        // Timeline events.
        if len(h.TodayEvents) > 0 {
                for i, evt := range h.TodayEvents {
                        h.renderEventLine(&b, evt, i == h.CursorIndex)
                }
        } else {
                b.WriteString(style.CaptionStyle.Render("—"))
        }

        b.WriteString("\n\n")

        // Separator.
        writeSeparator(&b, h.Width)

        // Summary.
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyHistorySummary) + ":"))
        b.WriteString("\n")
        h.renderSummary(&b, h.TodayStats)
        b.WriteString("\n\n")

        // Actions — doc: "↵  liat detail event    ←  hari kemarin    w  week view    q  balik"
        actions := []string{
                style.ActionStyle.Render(i18n.T(i18n.KeyHistoryViewEvent)),
                style.MutedStyle.Render(i18n.T(i18n.KeyHistoryPrevDay)),
                style.MutedStyle.Render(i18n.T(i18n.KeyHistorySwitchWeek)),
                style.MutedStyle.Render(i18n.T(i18n.KeyLabelBack)),
        }
        b.WriteString(strings.Join(actions, "    "))

        return b.String()
}

// viewWeek renders the history_week state.
//
// Spec:
//
//      minggu ini                              28 apr - 4 mei
//      pesan terkirim per hari:
//      senin   ████████████░░░░░░  12
//      selasa  ██████████████████  18  ★ terbaik
//      ...
//      response per hari:
//      ...
//      total minggu ini:
//      terkirim  68    respond  13    convert  3
//      ★ selasa = hari terbaik lu
func (h *History) viewWeek() string {
        var b strings.Builder

        // Title with week range.
        title := i18n.T(i18n.KeyHistoryWeek)
        startOfWeek := h.startOfWeek()
        endOfWeek := startOfWeek.AddDate(0, 0, daysInWeek-1)
        weekRange := fmt.Sprintf("%s - %s",
                startOfWeek.Format("2 Jan"),
                endOfWeek.Format("2 Jan"),
        )
        b.WriteString(renderTitleWithIndicator(title, weekRange, style.CaptionStyle, h.Width))
        b.WriteString("\n\n")

        // Separator.
        writeSeparator(&b, h.Width)

        // Messages per day chart.
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyHistoryMsgsPerDay) + ":"))
        b.WriteString("\n\n")
        h.renderInlineBarChart(&b, h.WeekData.Days, h.WeekData.Messages, h.WeekData.BestDayIndex)
        b.WriteString("\n\n")

        // Responses per day chart.
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyHistoryRespPerDay) + ":"))
        b.WriteString("\n\n")
        h.renderInlineBarChart(&b, h.WeekData.Days, h.WeekData.Responses, h.WeekData.BestDayIndex)
        b.WriteString("\n\n")

        // Separator.
        writeSeparator(&b, h.Width)

        // Week totals.
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyHistoryWeekTotal) + ":"))
        b.WriteString("\n")
        h.renderSummary(&b, h.WeekStats)
        if h.AvgResponseTime != "" {
                b.WriteString(style.MutedStyle.Render(
                        fmt.Sprintf("    %s: %s", i18n.T(i18n.KeyHistoryAvgRespTime), h.AvgResponseTime),
                ))
        }
        b.WriteString("\n\n")

        // Insight.
        if h.WeekData.BestDayLabel != "" {
                b.WriteString(style.GoldStyle.Render(
                        fmt.Sprintf("★ %s = %s", h.WeekData.BestDayLabel, i18n.T(i18n.KeyHistoryBestDay)),
                ))
                b.WriteString("\n")
                b.WriteString(style.CaptionStyle.Render(
                        fmt.Sprintf("  %s", i18n.T(i18n.KeyHistoryPrimeTime)),
                ))
                b.WriteString("\n")
                // Conversion rate for the best day — always backend-provided.
                // The TUI does NOT compute rates locally — the backend must provide
                // best_day_conv_rate. If the backend omits it, we simply don't show it.
                if h.WeekData.BestDayConvRate != "" {
                        b.WriteString(style.CaptionStyle.Render(
                                fmt.Sprintf("  %s %s: %s", i18n.T(i18n.KeyHistoryConvRate), h.WeekData.BestDayLabel, h.WeekData.BestDayConvRate),
                        ))
                }
                b.WriteString("\n\n")
        }

        // Separator before actions.
        writeSeparator(&b, h.Width)

        // Actions.
        actions := []string{
                style.MutedStyle.Render(i18n.T(i18n.KeyHistorySelectDay)),
                style.ActionStyle.Render(i18n.T(i18n.KeyHistoryViewDetail)),
                style.MutedStyle.Render(i18n.T(i18n.KeyHistorySwitchToday)),
                style.MutedStyle.Render(i18n.T(i18n.KeyLabelBack)),
        }
        b.WriteString(strings.Join(actions, "    "))

        return b.String()
}

// viewDayDetail renders the history_day_detail state.
//
// Spec:
//
//      selasa, 30 april 2024
//      ringkasan:
//      terkirim  18    respond  5    convert  1
//      timeline:
//      16:42  🎉  kopi nusantara      CONVERTED   rp 2.5jt
//      ...
//      → hari berikutnya    ← hari sebelumnya    q balik
func (h *History) viewDayDetail() string {
        var b strings.Builder

        // Date title — locale-aware formatting (doc: "selasa, 30 april 2024").
        dateTitle := i18n.FormatDate(h.DayDetailDate)
        b.WriteString(style.HeadingStyle.Render(dateTitle))
        b.WriteString("\n\n")

        // Separator.
        writeSeparator(&b, h.Width)

        // Summary.
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyHistorySummary) + ":"))
        b.WriteString("\n")
        h.renderSummary(&b, h.DayDetailStats)
        b.WriteString("\n\n")

        // Separator.
        writeSeparator(&b, h.Width)

        // Timeline.
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyHistoryTimeline) + ":"))
        b.WriteString("\n\n")

        for _, evt := range h.DayDetailEvents {
                h.renderEventLine(&b, evt, false)
        }

        b.WriteString("\n\n")

        // Separator.
        writeSeparator(&b, h.Width)

        // Actions.
        actions := []string{
                style.MutedStyle.Render(i18n.T(i18n.KeyHistoryNextDay)),
                style.MutedStyle.Render(i18n.T(i18n.KeyHistoryPrevDay)),
                style.MutedStyle.Render(i18n.T(i18n.KeyLabelBack)),
        }
        b.WriteString(strings.Join(actions, "    "))

        return b.String()
}

// ---------------------------------------------------------------------------
// Render helpers
// ---------------------------------------------------------------------------

// renderEventLine renders a single timeline event line.
func (h *History) renderEventLine(b *strings.Builder, evt HistoryEvent, selected bool) {
        // Time.
        if !evt.Time.IsZero() {
                b.WriteString(style.DimStyle.Render(evt.Time.Format("15:04")))
                b.WriteString("  ")
        }

        // Icon.
        if evt.Icon != "" {
                iconStyle := style.MutedStyle
                if evt.IsConversion {
                        iconStyle = style.GoldStyle
                } else if evt.Highlight {
                        iconStyle = style.WarningStyle
                }
                b.WriteString(iconStyle.Render(evt.Icon))
                b.WriteString("  ")
        }

        // Business name.
        nameStyle := style.BodyStyle
        if selected {
                nameStyle = style.SelectedBodyStyle
        }
        b.WriteString(nameStyle.Render(evt.BusinessName))

        // Event type.
        if evt.EventType != "" {
                b.WriteString("    ")
                typeStyle := style.MutedStyle
                if evt.IsConversion {
                        typeStyle = style.GoldStyle
                } else if evt.Highlight {
                        typeStyle = style.WarningStyle
                }
                b.WriteString(typeStyle.Render(evt.EventType))
        }

        // Detail.
        if evt.Detail != "" {
                b.WriteString("    ")
                b.WriteString(style.CaptionStyle.Render(evt.Detail))
        }

        // Revenue for conversions.
        if evt.Revenue != "" {
                b.WriteString("    ")
                b.WriteString(style.GoldStyle.Render(evt.Revenue))
        }

        b.WriteString("\n")
}

// renderSummary renders a day stats summary.
func (h *History) renderSummary(b *strings.Builder, stats DayStats) {
        b.WriteString(style.MutedStyle.Render(
                fmt.Sprintf("%s  %d    %s  %d    %s  %d",
                        i18n.T(i18n.KeyHistorySent), stats.Sent,
                        i18n.T(i18n.KeyHistoryRespond), stats.Respond,
                        i18n.T(i18n.KeyHistoryConvert), stats.Convert,
                ),
        ))
        b.WriteString("\n")
        b.WriteString(style.CaptionStyle.Render(
                fmt.Sprintf("%s %d   %s %dx",
                        i18n.T(i18n.KeyHistoryNewLeads), stats.NewLeads,
                        i18n.T(i18n.KeyHistoryScrape), stats.Scrapes,
                ),
        ))
        b.WriteString("\n")
}

// renderInlineBarChart renders the spec-style inline bar chart.
//
// Spec:
//
//      senin   ████████████░░░░░░  12
//      selasa  ██████████████████  18  ★ terbaik
func (h *History) renderInlineBarChart(b *strings.Builder, days []component.DayData, values []int64, bestIdx int) {
        maxVal := int64(1)
        for _, v := range values {
                if v > maxVal {
                        maxVal = v
                }
        }

        for i, day := range days {
                // Day label.
                label := day.Day
                if len(label) > dayLabelMaxRunes {
                        label = label[:dayLabelMaxRunes]
                }
                b.WriteString(style.BarLabelStyle.Render(label))

                // Bar.
                value := int64(0)
                if i < len(values) {
                        value = values[i]
                }
                fillWidth := int(float64(barMaxWidth) * float64(value) / float64(maxVal))
                if fillWidth > barMaxWidth {
                        fillWidth = barMaxWidth
                }
                emptyWidth := barMaxWidth - fillWidth

                barColor := style.TextMuted
                if i == bestIdx {
                        barColor = style.Gold
                }

                b.WriteString(lipgloss.NewStyle().Foreground(barColor).Render(
                        strings.Repeat("\u2588", fillWidth),
                ))
                b.WriteString(style.DimStyle.Render(
                        strings.Repeat("\u2591", emptyWidth),
                ))

                // Value.
                b.WriteString("  ")
                b.WriteString(style.BodyStyle.Render(fmt.Sprintf("%2d", value)))

                // Best day marker — i18n'd label, not hardcoded.
                if i == bestIdx {
                        b.WriteString("  ")
                        b.WriteString(style.WarningStyle.Render(fmt.Sprintf("★ %s", i18n.T(i18n.KeyHistoryBestLabel))))
                }

                // Weekend marker.
                weekday := h.startOfWeek().AddDate(0, 0, i).Weekday()
                if weekday == time.Saturday || weekday == time.Sunday {
                        b.WriteString("  ")
                        b.WriteString(style.CaptionStyle.Render(
                                fmt.Sprintf("(%s)", i18n.T(i18n.KeyHistoryWeekend)),
                        ))
                }

                b.WriteString("\n")
        }
}

// ---------------------------------------------------------------------------
// Data parsing helpers
// ---------------------------------------------------------------------------

func parseHistoryEvents(raw []any) []HistoryEvent {
        var events []HistoryEvent
        for _, item := range raw {
                if m, ok := item.(map[string]any); ok {
                        evt := HistoryEvent{
                                Icon:          strVal(m, protocol.ParamIcon),
                                BusinessName:  strVal(m, protocol.ParamBusinessName),
                                EventType:     strVal(m, protocol.ParamEventType),
                                Detail:        strVal(m, protocol.ParamDetail),
                                Revenue:       strVal(m, protocol.ParamRevenue),
                                Highlight:     boolVal(m, protocol.ParamHighlight),
                                IsConversion:  boolVal(m, protocol.ParamIsConversion),
                        }
                        if ts, ok := m[protocol.ParamTime].(string); ok {
                                if t, err := time.Parse("15:04", ts); err == nil {
                                        evt.Time = t
                                }
                        }
                        events = append(events, evt)
                }
        }
        return events
}

func parseDayStats(raw map[string]any) DayStats {
        return DayStats{
                Sent:     intVal(raw, protocol.ParamSent),
                Respond:  intVal(raw, protocol.ParamRespond),
                Convert:  intVal(raw, protocol.ParamConvert),
                NewLeads: intVal(raw, protocol.ParamNewLeads),
                Scrapes:  intVal(raw, protocol.ParamScrapes),
        }
}

func parseWeekData(raw map[string]any) WeekData {
        wd := WeekData{
                BestDayIndex:    intVal(raw, protocol.ParamBestDayIndex),
                BestDayLabel:    strVal(raw, protocol.ParamBestDayLabel),
                BestDayConvRate: strVal(raw, protocol.ParamBestDayConvRate),
        }

        // Use i18n-aware day labels — not hardcoded Indonesian.
        // The i18n system provides locale-correct day abbreviations.
        dayLabels := i18n.DayLabels()
        today := time.Now().Weekday()

        // Build day data.
        for i, label := range dayLabels {
                wd.Days = append(wd.Days, component.DayData{
                        Day:     label,
                        IsToday: int(today) == i+1, // Monday=1 in Go, Sunday=0
                })
        }

        // Parse message counts.
        if msgs, ok := raw[protocol.ParamMessages].([]any); ok {
                for _, v := range msgs {
                        if f, ok := v.(float64); ok {
                                wd.Messages = append(wd.Messages, int64(f))
                        }
                }
        }
        if resps, ok := raw[protocol.ParamResponses].([]any); ok {
                for _, v := range resps {
                        if f, ok := v.(float64); ok {
                                wd.Responses = append(wd.Responses, int64(f))
                        }
                }
        }
        if convs, ok := raw[protocol.ParamConverts].([]any); ok {
                for _, v := range convs {
                        if f, ok := v.(float64); ok {
                                wd.Converts = append(wd.Converts, int64(f))
                        }
                }
        }

        // Fill defaults if empty.
        for len(wd.Messages) < daysInWeek {
                wd.Messages = append(wd.Messages, 0)
        }
        for len(wd.Responses) < daysInWeek {
                wd.Responses = append(wd.Responses, 0)
        }
        for len(wd.Converts) < daysInWeek {
                wd.Converts = append(wd.Converts, 0)
        }

        return wd
}

func parseWeekStats(raw map[string]any) DayStats {
        return DayStats{
                Sent:     intVal(raw, protocol.ParamTotalSent),
                Respond:  intVal(raw, protocol.ParamTotalRespond),
                Convert:  intVal(raw, protocol.ParamTotalConvert),
                NewLeads: intVal(raw, protocol.ParamTotalNewLeads),
                Scrapes:  intVal(raw, protocol.ParamTotalScrapes),
        }
}

// populateTimeline syncs HistoryEvents into the Timeline component.
func (h *History) populateTimeline() {
        h.Timeline = component.NewTimeline()
        for i := len(h.TodayEvents) - 1; i >= 0; i-- {
                evt := h.TodayEvents[i]
                h.Timeline.Add(component.TimelineEvent{
                        Time:      evt.Time,
                        Icon:      evt.Icon,
                        Title:     evt.BusinessName,
                        Detail:    fmt.Sprintf("%s  %s", evt.EventType, evt.Detail),
                        Highlight: evt.Highlight,
                })
        }
}

// startOfWeek returns the Monday of the current date's week.
func (h *History) startOfWeek() time.Time {
        d := h.CurrentDate
        for d.Weekday() != time.Monday {
                d = d.AddDate(0, 0, -1)
        }
        return d
}

// String returns a debug representation.
func (h *History) String() string {
        return fmt.Sprintf("History{state=%s, date=%s, events=%d}", h.state, h.CurrentDate.Format("2006-01-02"), len(h.TodayEvents))
}
