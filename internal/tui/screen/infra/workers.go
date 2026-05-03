package infra

import (
        "fmt"
        "strings"

        "github.com/WaClaw-App/waclaw/internal/tui"
        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/internal/tui/component"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
        tea "github.com/charmbracelet/bubbletea"
        "github.com/charmbracelet/lipgloss"
)

// applyNavigateState extracts the "state" parameter from navigate params.
func applyNavigateState(state *protocol.StateID, params map[string]any) {
        if st, ok := params[protocol.ParamState].(string); ok {
                *state = protocol.StateID(st)
        }
}

// ---------------------------------------------------------------------------
// Data types for the Workers screen
// ---------------------------------------------------------------------------

// StageInfo represents a single pipeline stage's progress.
type StageInfo struct {
        Label    string // plain identifier: "scrape", "review", "antri", "kirim"
        Progress float64 // 0.0–1.0
        Detail   string
        Done     bool
}

// PerformaInfo holds a niche's performance metrics.
type PerformaInfo struct {
        ResponseRate    string
        ConversionRate string
        AvgRespond     string
}

// WorkerInfo holds all data for a single worker / niche.
type WorkerInfo struct {
        Name      string
        Active    bool
        Area      string
        Template  string
        Scrape    StageInfo
        Review    StageInfo
        Queue     StageInfo
        Send      StageInfo
        Queries   []QueryInfo
        Performa  PerformaInfo
        LeadsCollected int
        DuplicateCount  int
        LowRatingCount  int
        IceBreakerCount int
        AutoOfferCount  int
        NextIn          string
        TodaySent       int
        TodayLimit      int
        // FIX W-03: FoundPassed summary line under qualify
        FoundPassed     string
}

// QueryInfo represents a single scrape query status.
type QueryInfo struct {
        Text   string
        Status string // "done", "scanning", "waiting"
}

// PipelineTotals holds aggregate pipeline counts.
// FIX W-BF01: replaces hardcoded pipeline totals with struct fields
// that can be populated from HandleUpdate.
type PipelineTotals struct {
        Found  int
        Passed int
        Queued int
        Sent   int
}

// ---------------------------------------------------------------------------
// Workers screen model
// ---------------------------------------------------------------------------

// Workers implements tui.Screen for Screen 11: Workers Pipeline Visualizer.
// It visualises the per-niche worker pool with real-time pipeline progress,
// individual worker details, niche addition, and manual pause/resume.
type Workers struct {
        tui.ScreenBase
        state         protocol.StateID
        cursor        int
        workers       []WorkerInfo
        selected      int
        list          component.ListSelect
        width         int
        height        int
        focused       bool
        // FIX W-BF01: pipeline totals from backend, not hardcoded
        totals        PipelineTotals
}

// NewWorkers creates a Workers screen with demo data.
func NewWorkers() *Workers {
        w := &Workers{
                ScreenBase: tui.NewScreenBase(protocol.ScreenWorkers),
                state:      protocol.WorkersOverview,
                cursor:     0,
                workers:    demoWorkers(),
                // FIX W-BF01: pipeline totals start at zero — backend provides actual values
                totals: PipelineTotals{},
        }
        return w
}

// stageDisplayLabel resolves a stage identifier to a localized display label.
func stageDisplayLabel(id string) string {
        switch id {
        case "scrape":
                return i18n.T(i18n.KeyWorkersScrape)
        case "review", "qualify":
                // Detail view uses "qualify" per doc; overview uses "review"
                if id == "qualify" {
                        return i18n.T(i18n.KeyWorkersQualify)
                }
                return i18n.T(i18n.KeyWorkersReview)
        case "antri":
                return i18n.T(i18n.KeyWorkersQueue)
        case "kirim":
                return i18n.T(i18n.KeyWorkersSend)
        default:
                return id
        }
}

// demoWorkers returns realistic demo data for the workers screen.
// FIX 1: StageInfo.Label uses plain string identifiers, not i18n.T() calls.
func demoWorkers() []WorkerInfo {
        return []WorkerInfo{
                {
                        Name: "web_developer", Active: true,
                        Area: "kediri (15km)", Template: "direct-curiosity",
                        Scrape:   StageInfo{Label: "scrape", Progress: 0.67, Detail: "2/3 query selesai"},
                        Review:   StageInfo{Label: "review", Progress: 1.0, Detail: "done (89 lolos)", Done: true},
                        Queue:    StageInfo{Label: "antri", Progress: 0.6, Detail: "24 pesan"},
                        Send:     StageInfo{Label: "kirim", Progress: 0.5, Detail: "3/6 jam ini"},
                        Queries: []QueryInfo{
                                {Text: "cafe di kediri", Status: "done"},
                                {Text: "gym di kediri", Status: "scanning"},
                                {Text: "salon di kediri", Status: "waiting"},
                        },
                        DuplicateCount: 12, LowRatingCount: 5,
                        IceBreakerCount: 24, AutoOfferCount: 8,
                        Performa: PerformaInfo{"16%", "4.6%", "3.2 jam"},
                        // FIX W-03: found→passed summary
                        FoundPassed: "156 nemu → 89 lolos (57%)",
                },
                {
                        Name: "undangan_digital", Active: true,
                        Area: "kediri + surabaya", Template: "undangan-offer",
                        Scrape:   StageInfo{Label: "scrape", Progress: 1.0, Detail: "done (61 lolos)", Done: true},
                        Review:   StageInfo{Label: "review", Progress: 1.0, Detail: "done (48 lolos)", Done: true},
                        Queue:    StageInfo{Label: "antri", Progress: 0.35, Detail: "14 pesan"},
                        Send:     StageInfo{Label: "kirim", Progress: 0.17, Detail: "1/6 jam ini"},
                        DuplicateCount: 8, LowRatingCount: 3,
                        IceBreakerCount: 14, AutoOfferCount: 5,
                        Performa: PerformaInfo{"22%", "6.1%", "2.8 jam"},
                        // FIX W-03: found→passed summary
                        FoundPassed: "61 nemu → 48 lolos (79%)",
                },
                {
                        Name: "social_media_mgr", Active: false,
                        Area: "malang (10km)", Template: "smm-pitch",
                        Scrape:   StageInfo{Label: "scrape", Progress: 1.0, Detail: "done (23 lolos)", Done: true},
                        Review:   StageInfo{Label: "review", Progress: 1.0, Detail: "done (18 lolos)", Done: true},
                        Queue:    StageInfo{Label: "antri", Progress: 0.45, Detail: "18 pesan"},
                        Send:     StageInfo{Label: "kirim", Progress: 0.0, Detail: "○ idle (malam, mulai besok 09:00)"},
                        LeadsCollected: 18,
                        DuplicateCount: 4, LowRatingCount: 2,
                        IceBreakerCount: 18, AutoOfferCount: 6,
                        Performa: PerformaInfo{"19%", "5.3%", "3.5 jam"},
                        // FIX W-03: found→passed summary
                        FoundPassed: "23 nemu → 18 lolos (78%)",
                },
        }
}

// ---------------------------------------------------------------------------
// tea.Model interface
// ---------------------------------------------------------------------------

// Init returns the initial command (refresh animation).
func (w *Workers) Init() tea.Cmd {
        return nil
}

// Update handles key events and delegates to state-specific handlers.
func (w *Workers) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch msg.(type) {
        case tea.WindowSizeMsg:
                if m, ok := msg.(tea.WindowSizeMsg); ok {
                        w.width = m.Width
                        w.height = m.Height
                }
                return w, nil
        }

        if key, ok := msg.(tea.KeyMsg); ok {
                return w.handleKey(key), nil
        }

        return w, nil
}

// View renders the workers screen based on the current sub-state.
func (w *Workers) View() string {
        switch w.state {
        case protocol.WorkersOverview:
                return w.viewOverview()
        case protocol.WorkerDetail:
                return w.viewDetail()
        case protocol.WorkerAddNiche:
                return w.viewAddNiche()
        case protocol.WorkersPaused:
                return w.viewPaused()
        default:
                return w.viewOverview()
        }
}

// ---------------------------------------------------------------------------
// Screen interface
// ---------------------------------------------------------------------------

// ID returns the screen identifier (embedded via ScreenBase).
// HandleNavigate processes backend navigate commands.
func (w *Workers) HandleNavigate(params map[string]any) error {
        applyNavigateState(&w.state, params)
        if idx, ok := params[protocol.ParamSelected].(int); ok && idx >= 0 && idx < len(w.workers) {
                w.selected = idx
        }
        return nil
}

// HandleUpdate processes backend data updates.
func (w *Workers) HandleUpdate(params map[string]any) error {
        // Backend sends generic map data — convert to internal TUI types.
        // Do NOT assert TUI types from backend params (frontend/backend concern split).
        if raw, ok := params[protocol.ParamWorkersList]; ok {
                if list, ok := raw.([]map[string]any); ok {
                        var workers []WorkerInfo
                        for _, m := range list {
                                wi := WorkerInfo{}
                                if v, ok := m[protocol.ParamName].(string); ok {
                                        wi.Name = v
                                }
                                if v, ok := m[protocol.ParamActive].(bool); ok {
                                        wi.Active = v
                                }
                                if v, ok := m[protocol.ParamArea].(string); ok {
                                        wi.Area = v
                                }
                                if v, ok := m[protocol.ParamTemplate].(string); ok {
                                        wi.Template = v
                                }
                                workers = append(workers, wi)
                        }
                        if len(workers) > 0 {
                                w.workers = workers
                        }
                }
        }
        // FIX W-BF01: populate pipeline totals from backend
        if raw, ok := params[protocol.ParamPipelineTotals]; ok {
                if t, ok := raw.(map[string]any); ok {
                        if v, ok := t[protocol.ParamFound].(int); ok {
                                w.totals.Found = v
                        }
                        if v, ok := t[protocol.ParamPassed].(int); ok {
                                w.totals.Passed = v
                        }
                        if v, ok := t[protocol.ParamQueued].(int); ok {
                                w.totals.Queued = v
                        }
                        if v, ok := t[protocol.ParamSent].(int); ok {
                                w.totals.Sent = v
                        }
                }
        }
        return nil
}

// Focus marks the screen as active.
func (w *Workers) Focus() { w.focused = true }

// Blur marks the screen as inactive.
func (w *Workers) Blur() { w.focused = false }

// ---------------------------------------------------------------------------
// Key handling
// ---------------------------------------------------------------------------

func (w *Workers) handleKey(key tea.KeyMsg) tea.Model {
        switch key.String() {
        case "up", "k":
                if w.state == protocol.WorkersOverview {
                        if w.cursor > 0 {
                                w.cursor--
                        }
                }
        case "down", "j":
                if w.state == protocol.WorkersOverview {
                        if w.cursor < len(w.workers)-1 {
                                w.cursor++
                        }
                }
        case "enter":
                if w.state == protocol.WorkersOverview && w.cursor < len(w.workers) {
                        if !w.workers[w.cursor].Active {
                                w.selected = w.cursor
                                w.state = protocol.WorkersPaused
                        } else {
                                w.selected = w.cursor
                                w.state = protocol.WorkerDetail
                        }
                }
        case "q":
                // Let App handle back navigation
                return w
        case "n":
                if w.state == protocol.WorkersOverview {
                        w.state = protocol.WorkerAddNiche
                        w.list = component.NewListSelect([]component.ListItem{
                                {Label: i18n.T(i18n.KeyWorkersNicheFotografer), Description: i18n.T(i18n.KeyWorkersNicheFotograferDesc)},
                                {Label: i18n.T(i18n.KeyWorkersNicheAkuntan), Description: i18n.T(i18n.KeyWorkersNicheAkuntanDesc)},
                                {Label: i18n.T(i18n.KeyWorkersNicheCustom), Description: i18n.T(i18n.KeyWorkersNicheCustomDesc)},
                        })
                }
        case "1":
                if w.state == protocol.WorkerDetail {
                        // Pause worker — publish to backend instead of mutating locally
                        if w.Bus() != nil {
                                w.Bus().Publish(bus.ActionMsg{
                                        Action: protocol.ActionPauseWorker,
                                        Screen: w.ID(),
                                        Params: map[string]any{protocol.ParamNiche: w.workers[w.selected].Name},
                                })
                        }
                        w.state = protocol.WorkersPaused
                } else if w.state == protocol.WorkersPaused {
                        // Resume
                        w.workers[w.selected].Active = true
                        w.state = protocol.WorkerDetail
                }
        case "2":
                if w.state == protocol.WorkerDetail {
                        // Force scrape — notify backend
                        if w.Bus() != nil {
                                w.Bus().Publish(bus.ActionMsg{
                                        Action: protocol.ActionForceScrape,
                                        Screen: w.ID(),
                                        Params: map[string]any{protocol.ParamNiche: w.workers[w.selected].Name},
                                })
                        }
                } else if w.state == protocol.WorkersPaused {
                        // Delete — notify backend
                        if w.Bus() != nil {
                                w.Bus().Publish(bus.ActionMsg{
                                        Action: protocol.ActionDeleteWorker,
                                        Screen: w.ID(),
                                        Params: map[string]any{protocol.ParamNiche: w.workers[w.selected].Name},
                                })
                        }
                        w.state = protocol.WorkersOverview
                }
        }
        return w
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// FIX 2: viewOverview() key hint uses KeyWorkersChoose instead of KeyLabelHelp
func (w *Workers) viewOverview() string {
        activeCount := 0
        idleCount := 0
        for _, wk := range w.workers {
                if wk.Active {
                        activeCount++
                } else {
                        idleCount++
                }
        }

        var b strings.Builder

        // Heading — FIX CC-05: use renderHeadWithStatus
        b.WriteString(renderHeadWithStatus(i18n.T(i18n.KeyWorkersTitle), fmt.Sprintf("%d %s · %d %s", activeCount, i18n.T(i18n.KeyWorkersActive), idleCount, i18n.T(i18n.KeyWorkersIdle)), style.TextMuted))
        // FIX W-01: add separator after heading per doc
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(style.Separator())
        b.WriteString(style.Section(style.SectionGap))

        // Worker list
        for i, wk := range w.workers {
                if i > 0 {
                        b.WriteString(style.Section(style.SubSectionGap))
                }

                // Worker name with breathing highlight
                nameColor := style.Text
                if i == w.cursor {
                        nameColor = style.Accent
                }
                b.WriteString(lipgloss.NewStyle().Foreground(nameColor).Bold(true).Render(wk.Name))
                b.WriteString("\n")

                // FIX W-02: full-width underline bar under each worker name per doc
                // P3 compliance: use style.Separator() instead of ━ box-drawing chars
                b.WriteString(style.Separator())
                b.WriteString("\n")

                // Pipeline stages — FIX W-DRY01: use shared renderStageLine()
                stages := []StageInfo{wk.Scrape, wk.Review, wk.Queue, wk.Send}
                for _, st := range stages {
                        b.WriteString(renderStageLine(stageDisplayLabel(st.Label), st, 40))
                        b.WriteString("\n")
                }

                // Area + template
                b.WriteString(style.Indent(1))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                        fmt.Sprintf("%s: %s", i18n.T(i18n.KeyWorkersArea), wk.Area),
                ))
                b.WriteString("\n")
                b.WriteString(style.Indent(1))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                        fmt.Sprintf("%s: %s ★", i18n.T(i18n.KeyWorkersTemplate), wk.Template),
                ))
        }

        // Total pipeline summary — FIX W-BF01: use w.totals struct instead of hardcoded values
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("%s: %d %s → %d %s → %d %s → %d %s",
                        i18n.T(i18n.KeyWorkersTotalPipeline), w.totals.Found, i18n.T(i18n.KeyWorkersFound), w.totals.Passed, i18n.T(i18n.KeyWorkersPassed), w.totals.Queued, i18n.T(i18n.KeyWorkersQueued), w.totals.Sent, i18n.T(i18n.KeyWorkersSentLabel),
                )))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(i18n.T(i18n.KeyWorkersParallel)))

        // Key hints — FIX 2: use KeyWorkersChoose for "pilih worker"
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("↑↓  %s    ↵  %s    n  %s    q  %s",
                        i18n.T(i18n.KeyWorkersChoose), i18n.T(i18n.KeyWorkersDetail),
                        i18n.T(i18n.KeyWorkersAddNiche), i18n.T(i18n.KeyLabelBack),
                )))

        return b.String()
}

// FIX 3: viewDetail() heading uses "worker: {name}" format
// FIX 8: viewDetail() has sub-detail lines under qualify, queue, send stages
func (w *Workers) viewDetail() string {
        if w.selected >= len(w.workers) {
                return w.viewOverview()
        }
        wk := w.workers[w.selected]

        var b strings.Builder

        // Heading — FIX 3: "worker: {name}"
        b.WriteString(fmt.Sprintf("%s: %s", i18n.T(i18n.KeyWorkersWorkerLabel), wk.Name))
        b.WriteString(style.Section(style.SectionGap))

        // Pipeline label
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyWorkersPipelineLabel)))
        b.WriteString(style.Section(style.SubSectionGap))

        barWidth := 36

        // Scrape — FIX W-DRY01: use renderStageLine()
        b.WriteString(renderStageLine(i18n.T(i18n.KeyWorkersScrape), wk.Scrape, barWidth))
        if wk.Scrape.Progress > 0 && !wk.Scrape.Done {
                b.WriteString("  ")
                b.WriteString(lipgloss.NewStyle().Foreground(style.Success).Render(i18n.T(i18n.KeyWorkersActiveLabel)))
        }
        b.WriteString("\n")

        // Queries
        for _, q := range wk.Queries {
                b.WriteString(style.Indent(3))
                icon := "○"
                color := style.TextDim
                switch q.Status {
                case "done":
                        icon = "✗"
                        color = style.TextMuted
                case "scanning":
                        icon = "●"
                        color = style.Success
                }
                b.WriteString(lipgloss.NewStyle().Foreground(color).Render(icon))
                b.WriteString(" ")
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(q.Text))
                b.WriteString("\n")
        }

        // Review (qualify in detail view per doc)
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(renderStageLine(i18n.T(i18n.KeyWorkersQualify), wk.Review, barWidth))
        b.WriteString("\n")

        // FIX W-03: found→passed summary line under qualify
        if wk.FoundPassed != "" {
                b.WriteString(style.Indent(2))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(wk.FoundPassed))
                b.WriteString("\n")
        }

        // FIX 8: sub-detail lines under review (duplicates, low rating)
        b.WriteString(style.Indent(2))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf(i18n.T(i18n.KeyWorkersDuplicateCount), wk.DuplicateCount),
        ))
        b.WriteString("\n")
        b.WriteString(style.Indent(2))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf(i18n.T(i18n.KeyWorkersLowRatingCount), wk.LowRatingCount),
        ))
        b.WriteString("\n")

        // Queue
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(renderStageLine(i18n.T(i18n.KeyWorkersQueue), wk.Queue, barWidth))
        b.WriteString("\n")

        // FIX 8: sub-detail lines under queue (ice breaker, offer counts)
        b.WriteString(style.Indent(2))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf(i18n.T(i18n.KeyWorkersIceBreakerCount), wk.IceBreakerCount),
        ))
        b.WriteString("\n")
        b.WriteString(style.Indent(2))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf(i18n.T(i18n.KeyWorkersAutoOfferCount), wk.AutoOfferCount),
        ))
        b.WriteString("\n")

        // Send
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(renderStageLine(i18n.T(i18n.KeyWorkersSend), wk.Send, barWidth))
        b.WriteString("\n")

        // FIX 8: sub-detail lines under send (next in, today count)
        if wk.NextIn != "" {
                b.WriteString(style.Indent(2))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                        fmt.Sprintf(i18n.T(i18n.KeyWorkersNextIn), wk.NextIn),
                ))
                b.WriteString("\n")
        }
        b.WriteString(style.Indent(2))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf(i18n.T(i18n.KeyWorkersTodayCount), wk.TodaySent, wk.TodayLimit),
        ))

        // Performa
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyWorkersPerforma)))
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyWorkersResponseRate), wk.Performa.ResponseRate),
        ))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyWorkersConversionRate), wk.Performa.ConversionRate),
        ))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyWorkersAvgRespond), wk.Performa.AvgRespond),
        ))

        // Key hints
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("1  %s    2  %s    3  %s    q  %s",
                        i18n.T(i18n.KeyWorkersPauseWorker), i18n.T(i18n.KeyWorkersForceScrape),
                        i18n.T(i18n.KeyWorkersViewLeads), i18n.T(i18n.KeyLabelBack),
                )))

        return b.String()
}

// FIX W-04: key hints rendered BEFORE description lines per doc
func (w *Workers) viewAddNiche() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyWorkersAddNiche)))
        b.WriteString(style.Section(style.SubSectionGap))

        // TUI-vs-doc: "pilih niche buat ditambah ke pool:" and "(worker baru langsung jalan setelah dipilih)"
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyWorkersPickPool)))
        b.WriteString("\n")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyWorkersAutoStart)))
        b.WriteString(style.Section(style.SubSectionGap))

        b.WriteString(w.list.ViewWithNumbers())
        b.WriteString(style.Section(style.SectionGap))

        // FIX W-04: key hints BEFORE description lines per doc
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("%s    %s", i18n.T(i18n.KeyWorkersAddLabel), i18n.T(i18n.KeyWorkersQuitCancel)),
        ))
        b.WriteString(style.Section(style.SectionGap))

        // Description lines after key hints
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyWorkersMoreNiches)))
        b.WriteString("\n")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyWorkersIndependent)))

        return b.String()
}

// FIX W-05: render paused options as stacked numbered list, then lead count, then auto-resume
// FIX W-06: uses KeyWorkersStatusPaused which shows "⏸ PAUSED" per doc
func (w *Workers) viewPaused() string {
        if w.selected >= len(w.workers) {
                return w.viewOverview()
        }
        wk := w.workers[w.selected]

        var b strings.Builder

        // Heading with PAUSED badge — use renderHeadWithStatus from helpers.go
        headText := fmt.Sprintf("%s: %s", i18n.T(i18n.KeyWorkersWorkerLabel), wk.Name)
        b.WriteString(renderHeadWithStatus(headText, i18n.T(i18n.KeyWorkersStatusPaused), style.TextMuted))
        b.WriteString(style.Section(style.SectionGap))

        // Context — use KeyWorkersYouPaused
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyWorkersYouPaused)))
        b.WriteString(style.Section(style.SubSectionGap))

        // FIX W-05: 3 options as stacked numbered list per doc
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("1  %s", i18n.T(i18n.KeyWorkersOptionsResume)),
        ))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("2  %s", i18n.T(i18n.KeyWorkersOptionsDelete)),
        ))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("3  %s", i18n.T(i18n.KeyWorkersOptionsViewLeads)),
        ))
        b.WriteString(style.Section(style.SectionGap))

        // Lead count — FIX W-05: use KeyWorkersLeadCount
        b.WriteString(style.BodyStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyWorkersLeadCount), wk.LeadsCollected),
        ))
        b.WriteString("\n")

        // Auto-resume message — FIX W-05: use KeyWorkersAutoResumeMsg
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyWorkersAutoResumeMsg)))

        // Key hints
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("1  %s    2  %s    3  %s    q  %s",
                        i18n.T(i18n.KeyWorkersOptionsResume), i18n.T(i18n.KeyWorkersOptionsDelete),
                        i18n.T(i18n.KeyWorkersOptionsViewLeads), i18n.T(i18n.KeyLabelBack),
                )))

        return b.String()
}
