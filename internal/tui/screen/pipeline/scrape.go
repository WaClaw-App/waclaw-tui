// Package pipeline implements Screen 4: SCRAPE — the data harvesting screen.
//
// This is the most complex screen in the TUI with 12 distinct states/variants
// that visualize scraping progress, WA validation, high-value lead reveals,
// batch completion cascades, and various edge cases (idle, empty, error, throttle).
package pipeline

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
        "github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Internal tea.Msg types for timed animations
// ---------------------------------------------------------------------------

// jackpotSettleMsg fires after anim.JackpotSettle to collapse the jackpot reveal.
type jackpotSettleMsg struct{}

// batchCascadeMsg fires for staggered batch-complete card reveals.
type batchCascadeMsg struct{ index int }

// bellMsg triggers a terminal bell character.
type bellMsg struct{}

// ---------------------------------------------------------------------------
// Scrape model
// ---------------------------------------------------------------------------

// Scrape is Screen 4 — the most complex screen in the TUI with 12 states.
// It visualizes Google Maps scraping progress, WA pre-validation, high-value
// lead reveals (slot machine effect), batch completion cascades, and various
// edge cases (idle, empty, error, throttle).
//
// The screen follows the "anticipation" principle: scraping feels like watching
// a slot machine. Numbers increment live, names appear one-by-one, and
// high-value leads get a jackpot reveal with particles + terminal bell.
type Scrape struct {
        tui.ScreenBase
        state  protocol.StateID
        width  int
        height int

        // Active scrape data (single niche).
        niche      string
        target     string
        area       string
        filter     string
        found      int64
        qualified  int64
        duplicates int64
        newLeads   int64
        leads      []LeadItem
        scanning   bool

        // Multi-niche data.
        workerCount int
        niches      []NicheScrapeData
        queueTotal  int64

        // Idle state.
        lastScrapeTime string
        nextScrapeTime string

        // Error state.
        errorMsg string
        retryIn  string

        // Throttle state.
        throttleIn string

        // WA validation (per-niche, shown in ScrapeWAValidation).
        waNiches []NicheScrapeData

        // WA validation detail (ScrapeWAValidationProgress).
        waNicheName string
        waTotal     int64
        waHas       int64
        waNot       int64
        waPending   int64
        waPercent   float64
        waEstimate  string

        // High-value reveal.
        jackpotLead   HighValueLead
        jackpotActive bool
        jackpotBell   bool

        // Batch complete.
        batchResults []NicheScrapeData
        batchNextIn  string

        // Animation: how many batch cards are visible (for cascade effect).
        batchVisible  int
        cascadeStarted bool

        // Pending jackpot command flag (set in HandleNavigate/HandleUpdate, consumed in Update).
        pendingJackpotCmd bool

        // Particle system for jackpot effect.
        particles component.ParticleSystem
}

// NewScrape creates a Scrape screen with default values.
func NewScrape() *Scrape {
        return &Scrape{
                ScreenBase:   tui.NewScreenBase(protocol.ScreenScrape),
                state:        protocol.ScrapeActive,
                niches:       make([]NicheScrapeData, 0),
                leads:        make([]LeadItem, 0),
                waNiches:     make([]NicheScrapeData, 0),
                batchResults: make([]NicheScrapeData, 0),
                particles:    component.NewParticleSystem(40, 6),
        }
}

// Init returns the initial command.
func (m *Scrape) Init() tea.Cmd {
        return nil
}

// Update handles all tea.Msg for the scrape screen.
func (m *Scrape) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        // Process pending jackpot animation commands.
        if m.pendingJackpotCmd {
                m.pendingJackpotCmd = false
                return m, tea.Batch(m.scheduleBell(), m.scheduleJackpotSettle())
        }

        switch msg := msg.(type) {
        case tea.WindowSizeMsg:
                m.width = msg.Width
                m.height = msg.Height
                m.particles.Width = msg.Width
                m.particles.Height = 6
                // Kick off batch cascade animation on first render after entering state.
                if m.state == protocol.ScrapeBatchComplete && !m.cascadeStarted && len(m.batchResults) > 0 {
                        m.cascadeStarted = true
                        return m, m.scheduleBatchCascade(1)
                }
                return m, nil

        case tea.KeyMsg:
                return m.handleKey(msg)

        case jackpotSettleMsg:
                m.jackpotActive = false
                m.jackpotBell = false
                return m, nil

        case batchCascadeMsg:
                if msg.index > m.batchVisible {
                        m.batchVisible = msg.index
                }
                if m.batchVisible < len(m.batchResults) {
                        return m, m.scheduleBatchCascade(m.batchVisible + 1)
                }
                return m, nil

        case bellMsg:
                m.jackpotBell = true
                return m, nil
        }

        return m, nil
}

// View renders the scrape screen based on the current state.
func (m *Scrape) View() string {
        switch m.state {
        case protocol.ScrapeActive:
                return m.viewActive()
        case protocol.ScrapeMultiActive:
                return m.viewMultiActive()
        case protocol.ScrapeMultiStaggered:
                return m.viewMultiStaggered()
        case protocol.ScrapeIdle:
                return m.viewIdle()
        case protocol.ScrapeEmpty:
                return m.viewEmpty()
        case protocol.ScrapeError:
                return m.viewError()
        case protocol.ScrapeGMapsLimited:
                return m.viewGMapsLimited()
        case protocol.ScrapeAutoApproved:
                return m.viewAutoApproved()
        case protocol.ScrapeWAValidation:
                return m.viewWAValidation()
        case protocol.ScrapeHighValueReveal:
                return m.viewHighValueReveal()
        case protocol.ScrapeBatchComplete:
                return m.viewBatchComplete()
        case protocol.ScrapeWAValidationProgress:
                return m.viewWAValidationProgress()
        default:
                return m.viewActive()
        }
}

// HandleNavigate switches the screen state based on backend navigation params.
func (m *Scrape) HandleNavigate(params map[string]any) error {
        if s, ok := params[protocol.ParamState].(string); ok {
                m.state = protocol.StateID(s)
        }

        // Reset batch cascade tracking when entering batch complete.
        if m.state == protocol.ScrapeBatchComplete {
                m.batchVisible = 0
                m.cascadeStarted = false
        }

        // When entering jackpot, fire the bell and schedule settle.
        if m.state == protocol.ScrapeHighValueReveal {
                m.jackpotActive = true
                m.particles.Burst()
                m.pendingJackpotCmd = true
        }

        return nil
}

// HandleUpdate updates scrape data without changing state.
func (m *Scrape) HandleUpdate(params map[string]any) error {
        if v, ok := params[protocol.ParamNiche].(string); ok {
                m.niche = v
        }
        if v, ok := params[protocol.ParamTarget].(string); ok {
                m.target = v
        }
        if v, ok := params[protocol.ParamArea].(string); ok {
                m.area = v
        }
        if v, ok := params[protocol.ParamFilter].(string); ok {
                m.filter = v
        }
        if v, ok := toInt64(params[protocol.ParamFound]); ok {
                m.found = v
        }
        if v, ok := toInt64(params[protocol.ParamQualified]); ok {
                m.qualified = v
        }
        if v, ok := toInt64(params[protocol.ParamDuplicates]); ok {
                m.duplicates = v
        }
        if v, ok := toInt64(params[protocol.ParamNewLeads]); ok {
                m.newLeads = v
        }
        if v, ok := params[protocol.ParamLeads]; ok {
                m.leads = parseLeadItems(v)
        }
        if v, ok := params[protocol.ParamScanning].(bool); ok {
                m.scanning = v
        }
        if v, ok := toInt64(params[protocol.ParamWorkerCount]); ok {
                m.workerCount = int(v)
        }
        if v, ok := params[protocol.ParamNiches]; ok {
                m.niches = parseNicheScrapeData(v)
        }
        if v, ok := params[protocol.ParamLastScrape].(string); ok {
                m.lastScrapeTime = v
        }
        if v, ok := params[protocol.ParamNextScrape].(string); ok {
                m.nextScrapeTime = v
        }
        if v, ok := params[protocol.ParamError].(string); ok {
                m.errorMsg = v
        }
        if v, ok := params[protocol.ParamRetryIn].(string); ok {
                m.retryIn = v
        }
        if v, ok := params[protocol.ParamThrottleIn].(string); ok {
                m.throttleIn = v
        }
        if v, ok := params[protocol.ParamWANiches]; ok {
                m.waNiches = parseNicheScrapeData(v)
        }
        if v, ok := params[protocol.ParamWANicheName].(string); ok {
                m.waNicheName = v
        }
        if v, ok := toInt64(params[protocol.ParamWATotal]); ok {
                m.waTotal = v
        }
        if v, ok := toInt64(params[protocol.ParamWAHas]); ok {
                m.waHas = v
        }
        if v, ok := toInt64(params[protocol.ParamWANot]); ok {
                m.waNot = v
        }
        if v, ok := toInt64(params[protocol.ParamWAPending]); ok {
                m.waPending = v
        }
        if v, ok := toFloat64(params[protocol.ParamWAPercent]); ok {
                m.waPercent = v
        }
        if v, ok := params[protocol.ParamWAEstimate].(string); ok {
                m.waEstimate = v
        }
        if v, ok := params[protocol.ParamJackpotLead]; ok {
                m.jackpotLead = parseHighValueLead(v)
                m.jackpotActive = true
                m.particles.Burst()
                m.jackpotBell = false
                m.pendingJackpotCmd = true
        }
        if v, ok := params[protocol.ParamBatchResults]; ok {
                m.batchResults = parseNicheScrapeData(v)
                m.batchVisible = 0
        }
        if v, ok := params[protocol.ParamBatchNextIn].(string); ok {
                m.batchNextIn = v
        }
        if v, ok := toInt64(params[protocol.ParamQueueTotal]); ok {
                m.queueTotal = v
        }
        return nil
}

// Focus is called when this screen becomes active.
func (m *Scrape) Focus() {
        if m.state == protocol.ScrapeBatchComplete && m.batchVisible == 0 && len(m.batchResults) > 0 {
                m.batchVisible = 0
        }
}

// Blur is called when this screen loses focus.
func (m *Scrape) Blur() {}

// ---------------------------------------------------------------------------
// Key handling
// ---------------------------------------------------------------------------

func (m *Scrape) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch {
        case key.Matches(msg, tui.KeyBack):
                m.publishKey("q")
                return m, nil
        case key.Matches(msg, tui.KeyEnter):
                m.publishKey("enter")
                return m, nil
        case key.Matches(msg, tui.KeyTab):
                m.publishKey("tab")
                return m, nil
        case key.Matches(msg, tui.KeySkip):
                m.publishKey("s")
                return m, nil
        case key.Matches(msg, tui.Key1):
                m.publishKey("1")
                return m, nil
        case key.Matches(msg, tui.Key2):
                m.publishKey("2")
                return m, nil
        case key.Matches(msg, tui.Key3):
                m.publishKey("3")
                return m, nil
        }
        return m, nil
}

// publishKey sends a KeyPressMsg on the bus for the app router to handle.
func (m *Scrape) publishKey(k string) {
        if m.Bus() != nil {
                m.Bus().Publish(bus.KeyPressMsg{
                        Key:    k,
                        Screen: protocol.ScreenScrape,
                })
        }
}

// ---------------------------------------------------------------------------
// State 1: ScrapeActive — single niche active scraping
// ---------------------------------------------------------------------------

func (m *Scrape) viewActive() string {
        var b strings.Builder
        barW := m.barWidth()

        // Title row: title left, niche right.
        title := style.HeadingStyle.Render(i18n.T(i18n.KeyScrapeActive))
        right := style.MutedStyle.Render(fmt.Sprintf("%s: %s", i18n.T(i18n.KeyScrapeNicheLabel), m.niche))
        b.WriteString(joinTitleRight(title, right, m.width))
        b.WriteString(style.Section(style.SectionGap))

        // Metadata: target, area, filter.
        if m.target != "" {
                b.WriteString(style.Indent(1))
                b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyScrapeTarget)))
                b.WriteString(": ")
                b.WriteString(style.BodyStyle.Render(m.target))
                b.WriteString("\n")
        }
        if m.area != "" {
                b.WriteString(style.Indent(1))
                b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyScrapeArea)))
                b.WriteString(": ")
                b.WriteString(style.BodyStyle.Render(m.area))
                b.WriteString("\n")
        }
        if m.filter != "" {
                b.WriteString(style.Indent(1))
                b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyScrapeFilter)))
                b.WriteString(": ")
                b.WriteString(style.BodyStyle.Render(m.filter))
                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SubSectionGap))

        // Found count progress bar.
        foundLabel := fmt.Sprintf("%s %d", i18n.T(i18n.KeyScrapeFound), m.found)
        bar := makeBar(barW, clampPercent(m.found, 500), foundLabel)
        b.WriteString(bar.View())
        b.WriteString("\n")

        b.WriteString(style.Section(style.SubSectionGap))

        // Lead list with check/skip marks.
        for _, lead := range m.leads {
                b.WriteString(style.Indent(1))
                b.WriteString(style.BodyStyle.Render(lead.Name))
                b.WriteString("    ")

                if lead.HasWebsite {
                        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyScrapeHasWebsite)))
                        b.WriteString(" ")
                        b.WriteString(style.MutedStyle.Render("✗"))
                        b.WriteString(" ")
                        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyScrapeSkipSuffix)))
                } else {
                        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyScrapeNoWebsite)))
                        b.WriteString(" ")
                        b.WriteString(style.SuccessStyle.Render("✓"))
                }
                b.WriteString("\n")
        }

        if m.scanning {
                b.WriteString(style.Indent(1))
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeScanning)))
                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SectionGap))

        // Summary: qualified, duplicates, new leads.
        b.WriteString(style.Indent(1))
        b.WriteString(fmt.Sprintf("%s: %d %s\n", i18n.T(i18n.KeyScrapeQualified), m.qualified, i18n.T(i18n.KeyScrapeLeadsUnit)))
        b.WriteString(style.Indent(1))
        b.WriteString(fmt.Sprintf("%s: %d\n", i18n.T(i18n.KeyScrapeDuplicates), m.duplicates))
        b.WriteString(style.Indent(1))
        b.WriteString(fmt.Sprintf("%s: %d\n", i18n.T(i18n.KeyScrapeNewLeads), m.newLeads))

        b.WriteString("\n")

        // New leads progress bar.
        newLabel := fmt.Sprintf("%d %s", m.newLeads, i18n.T(i18n.KeyScrapeNewLeads))
        newBar := makeBar(barW, clampPercent(m.newLeads, m.qualified), newLabel)
        b.WriteString(newBar.View())
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // Auto-review message.
        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeAutoReview)))
        b.WriteString("\n")
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeAutoReviewHint)))

        return b.String()
}

// ---------------------------------------------------------------------------
// State 2: ScrapeMultiActive — multi-niche parallel scraping
// ---------------------------------------------------------------------------

func (m *Scrape) viewMultiActive() string {
        var b strings.Builder
        barW := m.barWidth()

        // Title row.
        title := style.HeadingStyle.Render(i18n.T(i18n.KeyScrapeMulti))
        right := style.MutedStyle.Render(fmt.Sprintf("%d %s", m.workerCount, i18n.T(i18n.KeyScrapeWorkers)))
        b.WriteString(joinTitleRight(title, right, m.width))
        b.WriteString(style.Section(style.SectionGap))

        // Per-niche sections.
        for _, n := range m.niches {
                b.WriteString(style.CaptionStyle.Render(n.Name))
                b.WriteString(padRight(m.width - len(n.Name)))
                b.WriteString(style.SuccessStyle.Render(i18n.T(i18n.KeyStatusActive)))
                b.WriteString("\n")

                if n.Targets != "" {
                        b.WriteString(style.Indent(1))
                        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyScrapeTarget)))
                        b.WriteString(": ")
                        b.WriteString(style.BodyStyle.Render(n.Targets))
                        b.WriteString("\n")
                }
                if n.Area != "" {
                        b.WriteString(style.Indent(1))
                        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyScrapeArea)))
                        b.WriteString(": ")
                        b.WriteString(style.BodyStyle.Render(n.Area))
                        b.WriteString("\n")
                }

                // Niche found count bar.
                nicheLabel := fmt.Sprintf("%s %d", i18n.T(i18n.KeyScrapeFound), n.Found)
                nicheBar := makeBar(barW, clampPercent(n.Found, 500), nicheLabel)
                b.WriteString(nicheBar.View())
                b.WriteString("\n")

                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("%s: %d · %s: %d\n",
                        i18n.T(i18n.KeyScrapeQualified), n.Qualified,
                        i18n.T(i18n.KeyScrapeNewLeads), n.NewLeads))

                b.WriteString(style.Section(style.SubSectionGap))
        }

        // Blank line before total summary.

        // Total summary.
        var totalFound, totalQualified, totalNew int64
        for _, n := range m.niches {
                totalFound += n.Found
                totalQualified += n.Qualified
                totalNew += n.NewLeads
        }
        b.WriteString(fmt.Sprintf("%s: %d %s · %d %s · %d %s\n",
                i18n.T(i18n.KeyScrapeTotal), totalFound, i18n.T(i18n.KeyScrapeFound),
                totalQualified, i18n.T(i18n.KeyScrapeQualified),
                totalNew, i18n.T(i18n.KeyScrapeNewLeads)))

        b.WriteString(style.Section(style.SubSectionGap))

        // Auto-pilot message.
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeAutoPilotMsg)))
        b.WriteString("\n")
        b.WriteString(style.Section(style.SubSectionGap))

        // Key hints.
        b.WriteString(style.CaptionStyle.Render("tab"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyScrapeTabSwitch)))
        b.WriteString("    ")
        b.WriteString(style.CaptionStyle.Render("↵"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyScrapeEnterDetail)))
        b.WriteString("    ")
        b.WriteString(style.CaptionStyle.Render("q"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyLabelBack)))
        b.WriteString("\n")

        return b.String()
}

// ---------------------------------------------------------------------------
// State 3: ScrapeMultiStaggered — staggered multi-niche scraping
// ---------------------------------------------------------------------------

func (m *Scrape) viewMultiStaggered() string {
        var b strings.Builder
        barW := m.barWidth()

        // Title row.
        title := style.HeadingStyle.Render(i18n.T(i18n.KeyScrapeStaggered))
        right := style.MutedStyle.Render(fmt.Sprintf("%d %s", m.workerCount, i18n.T(i18n.KeyScrapeWorkers)))
        b.WriteString(joinTitleRight(title, right, m.width))
        b.WriteString(style.Section(style.SectionGap))

        // Per-niche sections with different phases.
        for _, n := range m.niches {
                b.WriteString(style.CaptionStyle.Render(n.Name))
                padding := 26
                if padding < len(n.Name) {
                        padding = 0
                } else {
                        padding -= len(n.Name)
                }
                b.WriteString(strings.Repeat(" ", padding))

                // Status with phase-specific display.
                switch n.Status {
                case WorkerStatusScraping:
                        b.WriteString(style.SuccessStyle.Render(i18n.T(i18n.KeyStatusActive)))
                case WorkerStatusIdle:
                        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyStatusIdle)))
                        b.WriteString(fmt.Sprintf(" (%s: %s)", i18n.T(i18n.KeyScrapeNextScrape), n.NextIn))
                case WorkerStatusStarting:
                        b.WriteString(style.WarningStyle.Render(i18n.T(i18n.KeyStatusStarting)))
                default:
                        b.WriteString(style.MutedStyle.Render(n.Status))
                }
                b.WriteString("\n")

                // Show last batch for idle niches.
                if n.Status == WorkerStatusIdle && n.LastBatch != "" {
                        b.WriteString(style.Indent(1))
                        b.WriteString(style.MutedStyle.Render(fmt.Sprintf("%s: %s", i18n.T(i18n.KeyScrapeLastScrape), n.LastBatch)))
                        b.WriteString("\n")
                }

                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("%s: %d · %s: %d · %s: %d\n",
                        i18n.T(i18n.KeyScrapeFound), n.Found,
                        i18n.T(i18n.KeyScrapeQualified), n.Qualified,
                        i18n.T(i18n.KeyScrapeNewLeads), n.NewLeads))

                if n.Targets != "" {
                        b.WriteString(style.Indent(1))
                        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyScrapeTarget)))
                        b.WriteString(": ")
                        b.WriteString(style.BodyStyle.Render(n.Targets))
                        b.WriteString("\n")
                }
                if n.Area != "" {
                        b.WriteString(style.Indent(1))
                        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyScrapeArea)))
                        b.WriteString(": ")
                        b.WriteString(style.BodyStyle.Render(n.Area))
                        b.WriteString("\n")
                }

                if n.Status == WorkerStatusScraping || n.Status == WorkerStatusStarting {
                        nicheLabel := fmt.Sprintf("%s: %d %s", i18n.T(i18n.KeyScrapeFound), n.Found, i18n.T(i18n.KeyScrapeScanning))
                        nicheBar := makeBar(barW, clampPercent(n.Found, 300), nicheLabel)
                        b.WriteString(nicheBar.View())
                        b.WriteString("\n")
                }

                b.WriteString(style.Section(style.SubSectionGap))
        }

        // Blank line before total summary.
        b.WriteString(style.Section(style.SectionGap))

        // Total summary.
        b.WriteString(fmt.Sprintf("%s: %d %s · %d %s\n",
                i18n.T(i18n.KeyScrapeTotalActive), len(m.niches), i18n.T(i18n.KeyScrapeWorkers),
                m.workerCount, i18n.T(i18n.KeyScrapeWorkers)))

        // Total queue.
        totalQueue := m.queueTotal
        if totalQueue == 0 {
                for _, n := range m.niches {
                        totalQueue += n.NewLeads
                }
        }
        b.WriteString(fmt.Sprintf("%s: %d\n", i18n.T(i18n.KeyScrapeQueue), totalQueue))

        b.WriteString(style.Section(style.SubSectionGap))

        // Auto-pilot message.
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeAutoPilotMsg)))
        b.WriteString("\n")

        return b.String()
}

// ---------------------------------------------------------------------------
// State 4: ScrapeIdle — waiting for next interval
// ---------------------------------------------------------------------------

func (m *Scrape) viewIdle() string {
        var b strings.Builder

        // Title.
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyScrapeIdle)))
        b.WriteString(style.Section(style.SectionGap))

        // Last scrape stats.
        if m.lastScrapeTime != "" {
                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("%s: %s\n", i18n.T(i18n.KeyScrapeLastScrape), m.lastScrapeTime))
        }
        if m.newLeads > 0 {
                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf(i18n.T(i18n.KeyScrapeIdleNewLeads)+"\n", m.newLeads))
        }

        // Next scrape countdown.
        if m.nextScrapeTime != "" {
                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("%s: %s\n", i18n.T(i18n.KeyScrapeNextScrape), m.nextScrapeTime))
        }

        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeIdleWaiting)))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // Hint to scrape now or go back.
        b.WriteString(style.Indent(1))
        b.WriteString(style.CaptionStyle.Render("↵"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyScrapeSkipWait)))
        b.WriteString("    ")
        b.WriteString(style.CaptionStyle.Render("q"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyLabelBack)))
        b.WriteString("\n")

        return b.String()
}

// ---------------------------------------------------------------------------
// State 5: ScrapeEmpty — zero results
// ---------------------------------------------------------------------------

func (m *Scrape) viewEmpty() string {
        var b strings.Builder
        barW := m.barWidth()

        // Title.
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyScrapeEmpty)))
        b.WriteString(style.Section(style.SectionGap))

        // Target and area info.
        if m.target != "" {
                b.WriteString(style.Indent(1))
                b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyScrapeTarget)))
                b.WriteString(": ")
                b.WriteString(style.BodyStyle.Render(m.target))
                b.WriteString("\n")
        }
        if m.area != "" {
                b.WriteString(style.Indent(1))
                b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyScrapeArea)))
                b.WriteString(": ")
                b.WriteString(style.BodyStyle.Render(m.area))
                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SubSectionGap))

        // Empty progress bar (0%).
        emptyBar := makeBar(barW, 0, "0")
        b.WriteString(emptyBar.View())
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // Empty hint message.
        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeEmptyHints)))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(style.CaptionStyle.Render("  - "))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeEmptyAreaHint)))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(style.CaptionStyle.Render("  - "))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeEmptyFilterHint)))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(style.CaptionStyle.Render("  - "))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeEmptyQueryHint)))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // Numbered action options.
        b.WriteString(style.Indent(1))
        b.WriteString(style.ActionStyle.Render("1"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyScrapeEditFilter)))
        b.WriteString("    ")
        b.WriteString(style.ActionStyle.Render("2"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyScrapeChangeArea)))
        b.WriteString("    ")
        b.WriteString(style.ActionStyle.Render("3"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyScrapeAddQuery)))
        b.WriteString("\n")

        return b.String()
}

// ---------------------------------------------------------------------------
// State 6: ScrapeError — scraper crash
// ---------------------------------------------------------------------------

func (m *Scrape) viewError() string {
        var b strings.Builder

        // Title.
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyScrapeError)))
        b.WriteString(style.Section(style.SectionGap))

        // Error message with ✗.
        b.WriteString(style.Indent(1))
        b.WriteString(style.DangerStyle.Render("✗ "))
        b.WriteString(style.DangerStyle.Render(i18n.T(i18n.KeyScrapeScraperError)))
        if m.errorMsg != "" {
                b.WriteString(" — ")
                b.WriteString(style.BodyStyle.Render(m.errorMsg))
        }
        b.WriteString("\n")

        b.WriteString(style.Section(style.SubSectionGap))

        // Auto-retry countdown.
        if m.retryIn != "" {
                b.WriteString(style.Indent(1))
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyGeneralRetry)))
                b.WriteString(" ")
                b.WriteString(style.BodyStyle.Render(m.retryIn))
                b.WriteString(".\n")
        }

        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeErrorReassure)))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SubSectionGap))

        // Retry option.
        b.WriteString(style.Indent(1))
        b.WriteString(style.ActionStyle.Render("1"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyScrapeRetry)))
        b.WriteString("    ")
        b.WriteString(style.CaptionStyle.Render("q"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyLabelBack)))
        b.WriteString("\n")

        return b.String()
}

// ---------------------------------------------------------------------------
// State 7: ScrapeGMapsLimited — rate limited
// ---------------------------------------------------------------------------

func (m *Scrape) viewGMapsLimited() string {
        var b strings.Builder

        // Title.
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyScrapeIdle)))
        b.WriteString(style.Section(style.SectionGap))

        // Throttle message.
        b.WriteString(style.Indent(1))
        b.WriteString(style.WarningStyle.Render("⏳ "))
        b.WriteString(style.WarningStyle.Render(i18n.T(i18n.KeyScrapeGMapsThrottle)))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SubSectionGap))

        // Auto-resume countdown.
        if m.throttleIn != "" {
                b.WriteString(style.Indent(1))
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeThrottleExplain)))
                b.WriteString("\n")
                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("%s %s.", i18n.T(i18n.KeyScrapeThrottleResume), m.throttleIn))
                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SubSectionGap))

        // Hint that other scrapers still run.
        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeThrottleGMapsOnly)))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SubSectionGap))

        // Retry option.
        b.WriteString(style.Indent(1))
        b.WriteString(style.ActionStyle.Render("1"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyScrapeRetry)))
        b.WriteString("    ")
        b.WriteString(style.CaptionStyle.Render("q"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyLabelBack)))
        b.WriteString("\n")

        return b.String()
}

// ---------------------------------------------------------------------------
// State 8: ScrapeAutoApproved — auto-pilot variant
// ---------------------------------------------------------------------------

func (m *Scrape) viewAutoApproved() string {
        var b strings.Builder
        barW := m.barWidth()

        // Title row.
        title := style.HeadingStyle.Render(i18n.T(i18n.KeyScrapeAutoApproved))
        right := style.MutedStyle.Render(fmt.Sprintf("%d %s", m.workerCount, i18n.T(i18n.KeyScrapeWorkers)))
        b.WriteString(joinTitleRight(title, right, m.width))
        b.WriteString(style.Section(style.SectionGap))

        var totalNew int64

        // Per-niche summary.
        for _, n := range m.niches {
                b.WriteString(style.CaptionStyle.Render(n.Name))
                b.WriteString("\n")
                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("%s: %d · %s: %d\n",
                        i18n.T(i18n.KeyScrapeFound), n.Found,
                        i18n.T(i18n.KeyScrapeQualified), n.Qualified))
                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("%s: %d · %s: %d\n",
                        i18n.T(i18n.KeyScrapeDuplicates), n.Duplicates,
                        i18n.T(i18n.KeyScrapeNewLeads), n.NewLeads))
                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("%s: %d %s, %d %s\n",
                        i18n.T(i18n.KeyScrapeAutoReview), n.NewLeads, i18n.T(i18n.KeyScrapeQueue),
                        n.Skipped, i18n.T(i18n.KeyScrapeSkipSuffix)))
                totalNew += n.NewLeads
                b.WriteString(style.Section(style.SubSectionGap))
        }

        // Blank line before total summary.
        b.WriteString(style.Section(style.SectionGap))

        // Total batch summary with progress bar.
        totalLabel := fmt.Sprintf("%d %s", totalNew, i18n.T(i18n.KeyScrapeNewLeadsWaiting))
        totalBar := makeBar(barW, 1.0, totalLabel)
        b.WriteString(totalBar.View())
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // Auto-pilot message.
        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeAutoPilotMsg)))
        b.WriteString("\n")

        return b.String()
}

// ---------------------------------------------------------------------------
// State 9: ScrapeWAValidation — WA validation running (per-niche overview)
// ---------------------------------------------------------------------------

func (m *Scrape) viewWAValidation() string {
        var b strings.Builder
        barW := m.barWidth()

        // Title row.
        title := style.HeadingStyle.Render(i18n.T(i18n.KeyScrapeWAValidation))
        right := style.MutedStyle.Render(fmt.Sprintf("%d %s", m.workerCount, i18n.T(i18n.KeyScrapeWorkers)))
        b.WriteString(joinTitleRight(title, right, m.width))
        b.WriteString(style.Section(style.SectionGap))

        var totalReady int64

        // Per-niche scrape + WA validation data.
        for _, n := range m.waNiches {
                b.WriteString(style.CaptionStyle.Render(n.Name))
                b.WriteString("\n")
                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("%s: %d · %s: %d\n",
                        i18n.T(i18n.KeyScrapeFound), n.Found,
                        i18n.T(i18n.KeyScrapeQualified), n.Qualified))
                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("%s: %d · %s: %d\n",
                        i18n.T(i18n.KeyScrapeDuplicates), n.Duplicates,
                        i18n.T(i18n.KeyScrapeNewLeads), n.NewLeads))

                // WA validation progress bar.
                b.WriteString("\n")
                b.WriteString(style.Indent(1))
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeWAValidation)))
                b.WriteString("\n")

                if n.WAChecked > 0 && n.NewLeads > 0 {
                        waPct := float64(n.WAChecked) / float64(n.NewLeads)
                        waLabel := fmt.Sprintf("%.0f%% (%d/%d %s)", waPct*100, n.WAChecked, n.NewLeads, i18n.T(i18n.KeyScrapeWACekUnit))
                        waBar := makeBar(barW, waPct, waLabel)
                        b.WriteString(waBar.View())
                        b.WriteString("\n")
                }

                // WA counts per niche.
                b.WriteString("\n")
                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("✓ %s:  %d   %s\n", i18n.T(i18n.KeyScrapeWAHas), n.WAHas, i18n.T(i18n.KeyScrapeWAQueueAnnotation)))
                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("✗ %s:  %d   %s\n", i18n.T(i18n.KeyScrapeWANot), n.WANot, i18n.T(i18n.KeyScrapeWAMarkedAnnotation)))
                if n.WAPending > 0 {
                        b.WriteString(style.Indent(1))
                        b.WriteString(fmt.Sprintf("⏳ %s: %d\n", i18n.T(i18n.KeyScrapeWAPending), n.WAPending))
                }

                totalReady += n.WAHas
                b.WriteString(style.Section(style.SectionGap))
        }

        // Total ready-to-send summary.
        b.WriteString(fmt.Sprintf("%s: %d %s\n", i18n.T(i18n.KeyScrapeTotalReady), totalReady, i18n.T(i18n.KeyScrapeLeadsUnit)))

        // Count non-WA leads.
        var totalNotWA int64
        for _, n := range m.waNiches {
                totalNotWA += n.WANot
        }
        if totalNotWA > 0 {
                b.WriteString(fmt.Sprintf("%d %s — %s\n", totalNotWA,
                        i18n.T(i18n.KeyScrapeWANot), i18n.T(i18n.KeyScrapeWASaving)))
        }

        b.WriteString(style.Section(style.SubSectionGap))

        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeWABackground)))
        b.WriteString("\n")

        return b.String()
}

// ---------------------------------------------------------------------------
// State 10: ScrapeHighValueReveal — slot machine jackpot!
// ---------------------------------------------------------------------------

func (m *Scrape) viewHighValueReveal() string {
        var b strings.Builder
        barW := m.barWidth()

        // Title row: scrape title, niche right.
        title := style.HeadingStyle.Render(i18n.T(i18n.KeyScrapeActive))
        right := style.MutedStyle.Render(m.niche)
        b.WriteString(joinTitleRight(title, right, m.width))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.Section(1))
        b.WriteString(style.Section(style.SubSectionGap))

        // ★ lead name with score (gold style).
        b.WriteString(style.Indent(1))
        b.WriteString(style.GoldStyle.Render(fmt.Sprintf("★  %s     %s: %.1f",
                m.jackpotLead.Name, i18n.T(i18n.KeyScrapeScore), m.jackpotLead.Score)))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SubSectionGap))

        // Business details: category, address.
        b.WriteString(style.Indent(1))
        b.WriteString(style.BodyStyle.Render(fmt.Sprintf("%s · %s",
                m.jackpotLead.Category, m.jackpotLead.Address)))
        b.WriteString("\n")

        // Rating, reviews, web, IG badges.
        var details []string
        if m.jackpotLead.Rating > 0 {
                details = append(details, fmt.Sprintf("⭐ %.1f", m.jackpotLead.Rating))
        }
        if m.jackpotLead.Reviews > 0 {
                details = append(details, fmt.Sprintf("%d %s", m.jackpotLead.Reviews, i18n.T(i18n.KeyWordReviews)))
        }
        if !m.jackpotLead.HasWeb {
                details = append(details, i18n.T(i18n.KeyScrapeNoWebsite))
        }
        if m.jackpotLead.HasIG {
                details = append(details, i18n.T(i18n.KeyWordActiveIG))
        }
        if len(details) > 0 {
                b.WriteString(style.Indent(1))
                b.WriteString(style.CaptionStyle.Render(strings.Join(details, " · ")))
                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SubSectionGap))

        b.WriteString(style.Section(style.SectionGap))

        // ★ ★ ★ JACKPOT LEAD! ★ ★ ★
        b.WriteString(style.Indent(1))
        jackpotText := fmt.Sprintf("★ ★ ★   %s   ★ ★ ★", i18n.T(i18n.KeyLabelJackpot))
        b.WriteString(style.GoldStyle.Render(jackpotText))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SubSectionGap))

        // Priority queue and template messages.
        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapePriorityQueue)))
        b.WriteString(".\n")
        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeBestTemplate)))
        b.WriteString(".\n")

        // Particle burst effect.
        if m.particles.Active {
                b.WriteString(style.Section(style.SubSectionGap))
                b.WriteString(m.particles.ViewCompact())
                b.WriteString("\n")
        }

        // Terminal bell character for audio feedback.
        if m.jackpotBell {
                b.WriteString("\a")
                m.jackpotBell = false
        }

        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.Section(style.SectionGap))

        // Progress counter below (found count continues during scraping).
        foundLabel := fmt.Sprintf("%s %d", i18n.T(i18n.KeyScrapeFound), m.found)
        foundBar := makeBar(barW, clampPercent(m.found, 500), foundLabel)
        b.WriteString(foundBar.View())
        b.WriteString("\n")

        if m.scanning {
                b.WriteString(style.Section(style.SubSectionGap))
                b.WriteString(style.Indent(1))
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeScanning)))
                b.WriteString("\n")
        }

        return b.String()
}

// ---------------------------------------------------------------------------
// State 11: ScrapeBatchComplete — batch finished with cascade effect
// ---------------------------------------------------------------------------

func (m *Scrape) viewBatchComplete() string {
        var b strings.Builder
        barW := m.barWidth()

        // Title row: scrape title, "batch selesai!" right.
        title := style.HeadingStyle.Render(i18n.T(i18n.KeyScrapeActive))
        right := style.SuccessStyle.Render(i18n.T(i18n.KeyScrapeBatchSelesai))
        b.WriteString(joinTitleRight(title, right, m.width))
        b.WriteString(style.Section(style.SectionGap))

        var totalNew int64

        // Per-niche completion cards (revealed by cascade animation).
        for idx, n := range m.batchResults {
                if idx >= m.batchVisible {
                        break
                }

                b.WriteString(style.Section(style.SubSectionGap))

                // ✓ niche name selesai.
                b.WriteString(style.SuccessStyle.Render(fmt.Sprintf("✓ batch %s %s", n.Name, i18n.T(i18n.KeyScrapeBatchDone))))
                b.WriteString("\n")
                b.WriteString(style.Section(style.SubSectionGap))

                // Niche stats.
                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("%d %s → %d %s → %d %s\n",
                        n.Found, i18n.T(i18n.KeyScrapeFound),
                        n.Qualified, i18n.T(i18n.KeyScrapeQualified),
                        n.NewLeads, i18n.T(i18n.KeyScrapeNewLeads)))
                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("%s: %d · %s: %d\n",
                        i18n.T(i18n.KeyScrapeDuplicates), n.Duplicates,
                        i18n.T(i18n.KeyScrapeSkipSuffix), n.Skipped))

                // Per-niche messages.
                b.WriteString(style.Indent(1))
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeBatchAutoQueue)))
                b.WriteString("\n")
                b.WriteString(style.Indent(1))
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeBatchWorkHours)))
                b.WriteString("\n")

                totalNew += n.NewLeads
        }

        // Cascade animation is driven by batchCascadeMsg in Update(),
        // started on the first WindowSizeMsg after entering batch complete state.
        // Cards up to batchVisible are rendered here.

        b.WriteString(style.Section(style.SectionGap))

        // Total summary with progress bar.
        totalLabel := fmt.Sprintf(i18n.T(i18n.KeyScrapeBatchTotalSummary), totalNew)
        totalBar := makeBar(barW, 1.0, totalLabel)
        b.WriteString(totalBar.View())
        b.WriteString("\n")

        // Next batch countdown.
        if m.batchNextIn != "" {
                b.WriteString(style.Section(style.SubSectionGap))
                b.WriteString(fmt.Sprintf("%s: %s\n", i18n.T(i18n.KeyScrapeNextScrape), m.batchNextIn))
        }

        return b.String()
}

// ---------------------------------------------------------------------------
// State 12: ScrapeWAValidationProgress — detailed WA validation
// ---------------------------------------------------------------------------

func (m *Scrape) viewWAValidationProgress() string {
        var b strings.Builder
        barW := m.barWidth()

        // Title row: WA validation title, niche name right.
        title := style.HeadingStyle.Render(i18n.T(i18n.KeyScrapeWAValidationProg))
        right := style.MutedStyle.Render(m.waNicheName)
        b.WriteString(joinTitleRight(title, right, m.width))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.Section(1))
        b.WriteString(style.Section(style.SubSectionGap))

        // Total nomor dicek.
        b.WriteString(style.Indent(1))
        b.WriteString(fmt.Sprintf("%s: %d\n", i18n.T(i18n.KeyScrapeWATotal), m.waTotal))

        // Large progress bar with percentage.
        b.WriteString(style.Section(style.SubSectionGap))
        // Normalize waPercent: backend may send 0-100 or 0-1.
        waPct := m.waPercent
        if waPct > 1 {
                waPct = waPct / 100
        }
        progressLabel := fmt.Sprintf("%.0f%%", waPct*100)
        progressBar := makeBar(barW, waPct, progressLabel)
        progressBar.ShowPercent = false // Label already contains percentage.
        b.WriteString(progressBar.View())
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // ✓/✗/⏳ counts with percentages.
        if m.waTotal > 0 {
                hasPct := float64(m.waHas) / float64(m.waTotal) * 100
                notPct := float64(m.waNot) / float64(m.waTotal) * 100
                pendPct := float64(m.waPending) / float64(m.waTotal) * 100

                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("✓ %-16s %4d  (%.1f%%)\n",
                        i18n.T(i18n.KeyScrapeWAHas), m.waHas, hasPct))

                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("✗ %-16s %4d  (%.1f%%)\n",
                        i18n.T(i18n.KeyScrapeWANot), m.waNot, notPct))

                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("⏳ %-16s %4d  (%.1f%%)\n",
                        i18n.T(i18n.KeyScrapeWAPending), m.waPending, pendPct))
        }

        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.Section(1))
        b.WriteString(style.Section(style.SubSectionGap))

        // Information messages about queue and retry.
        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeWAReady)))
        b.WriteString(".\n")
        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyScrapeWASkipped)))
        b.WriteString(".\n")

        // Estimated completion time.
        if m.waEstimate != "" {
                b.WriteString(style.Section(style.SubSectionGap))
                b.WriteString(style.Indent(1))
                b.WriteString(fmt.Sprintf("%s: %s\n", i18n.T(i18n.KeyScrapeEstimate), m.waEstimate))
        }

        b.WriteString(style.Section(style.SectionGap))

        // Key hints.
        b.WriteString(style.CaptionStyle.Render("↵"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyScrapeEnterDetail)))
        b.WriteString("    ")
        b.WriteString(style.CaptionStyle.Render("s"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyScrapeSkipWait)))
        b.WriteString("    ")
        b.WriteString(style.CaptionStyle.Render("q"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyLabelBack)))
        b.WriteString("\n")

        return b.String()
}

// ---------------------------------------------------------------------------
// Animation scheduling
// ---------------------------------------------------------------------------

// scheduleJackpotSettle returns a tea.Cmd that fires after anim.JackpotSettle
// to collapse the jackpot reveal back to normal scraping view.
func (m *Scrape) scheduleJackpotSettle() tea.Cmd {
        return tea.Tick(anim.JackpotSettle, func(_ time.Time) tea.Msg {
                return jackpotSettleMsg{}
        })
}

// scheduleBell returns a tea.Cmd that fires a terminal bell immediately.
func (m *Scrape) scheduleBell() tea.Cmd {
        return func() tea.Msg {
                return bellMsg{}
        }
}

// scheduleBatchCascade returns a tea.Cmd that fires after anim.BatchCascadeStagger
// to reveal the next batch completion card.
func (m *Scrape) scheduleBatchCascade(index int) tea.Cmd {
        return tea.Tick(anim.BatchCascadeStagger, func(_ time.Time) tea.Msg {
                return batchCascadeMsg{index: index}
        })
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// barWidth returns the default progress bar width based on terminal width.
func (m *Scrape) barWidth() int {
        if m.width > 20 {
                return m.width - 20
        }
        return 40
}

// makeBar creates a ProgressBar with the given width, percent, and label.
func makeBar(width int, percent float64, label string) component.ProgressBar {
        return component.ProgressBar{
                Width:       width,
                Percent:     percent,
                Label:       label,
                ShowPercent: true,
                FillChar:    "━",
                EmptyChar:   "░",
        }
}

// clampPercent returns a percentage clamped to [0.0, 1.0] based on value/max.
// Returns 0.0 if max is zero or negative.
func clampPercent(value, max int64) float64 {
        if max <= 0 {
                return 0
        }
        p := float64(value) / float64(max)
        if p > 1.0 {
                p = 1.0
        }
        if p < 0 {
                p = 0
        }
        return p
}

// joinTitleRight joins a title and right-aligned label within the given width.
// Uses lipgloss for proper rendering.
func joinTitleRight(left, right string, width int) string {
        if width <= 0 {
                return left
        }
        leftW := lipgloss.Width(left)
        rightW := lipgloss.Width(right)
        space := width - leftW - rightW
        if space < 1 {
                space = 1
        }
        return left + strings.Repeat(" ", space) + right + "\n"
}

// padRight returns a string of spaces with the given length (min 0).
func padRight(n int) string {
        if n < 0 {
                n = 0
        }
        return strings.Repeat(" ", n)
}

// ---------------------------------------------------------------------------
// Param parsing helpers
// ---------------------------------------------------------------------------

// toInt64 extracts an int64 from an any value.
// Supports int, int64, and float64 (from JSON numbers).
func toInt64(v any) (int64, bool) {
        switch n := v.(type) {
        case int64:
                return n, true
        case int:
                return int64(n), true
        case float64:
                return int64(n), true
        }
        return 0, false
}

// toFloat64 extracts a float64 from an any value.
func toFloat64(v any) (float64, bool) {
        switch n := v.(type) {
        case float64:
                return n, true
        case int:
                return float64(n), true
        case int64:
                return float64(n), true
        }
        return 0, false
}

// parseLeadItems converts a raw param value ([]any of maps) to a slice of LeadItem.
func parseLeadItems(v any) []LeadItem {
        items, ok := v.([]any)
        if !ok {
                return nil
        }
        result := make([]LeadItem, 0, len(items))
        for _, item := range items {
                m, ok := item.(map[string]any)
                if !ok {
                        continue
                }
                li := LeadItem{
                        Name:      castString(m["name"]),
                        HasWebsite: castBool(m["has_web"]),
                        IsNew:     castBool(m["is_new"]),
                        Qualified: castBool(m["qualified"]),
                }
                result = append(result, li)
        }
        return result
}

// parseNicheScrapeData converts a raw param value ([]any of maps) to a slice of NicheScrapeData.
func parseNicheScrapeData(v any) []NicheScrapeData {
        items, ok := v.([]any)
        if !ok {
                return nil
        }
        result := make([]NicheScrapeData, 0, len(items))
        for _, item := range items {
                m, ok := item.(map[string]any)
                if !ok {
                        continue
                }
                nsd := NicheScrapeData{
                        Name:       castString(m["name"]),
                        Status:     castString(m["status"]),
                        NextIn:     castString(m["next_in"]),
                        Targets:    castString(m[protocol.ParamTarget]),
                        Area:       castString(m["area"]),
                        Filter:     castString(m["filter"]),
                        Found:      castInt64(m["found"]),
                        Qualified:  castInt64(m["qualified"]),
                        Duplicates: castInt64(m["duplicates"]),
                        NewLeads:   castInt64(m["new_leads"]),
                        Skipped:    castInt64(m["skipped"]),
                        WAHas:      castInt64(m["wa_has"]),
                        WANot:      castInt64(m["wa_not"]),
                        WAPending:  castInt64(m["wa_pending"]),
                        WAChecked:  castInt64(m["wa_checked"]),
                        Leads:      parseLeadItems(m["leads"]),
                }
                result = append(result, nsd)
        }
        return result
}

// parseHighValueLead converts a raw param value (map) to a HighValueLead.
func parseHighValueLead(v any) HighValueLead {
        m, ok := v.(map[string]any)
        if !ok {
                return HighValueLead{}
        }
        return HighValueLead{
                Name:     castString(m["name"]),
                Category: castString(m["category"]),
                Address:  castString(m["address"]),
                Rating:   anyFloat64(m["rating"]),
                Reviews:  castInt64(m["reviews"]),
                HasWeb:   castBool(m["has_web"]),
                HasIG:    castBool(m["has_ig"]),
                Score:    anyFloat64(m["score"]),
        }
}

// anyString extracts a string from an any, returning "" if not a string.
func castString(v any) string {
        s, _ := v.(string)
        return s
}

// anyBool extracts a bool from an any, returning false if not a bool.
func castBool(v any) bool {
        b, _ := v.(bool)
        return b
}

// anyInt64 extracts an int64 from an any, returning 0 if conversion fails.
func castInt64(v any) int64 {
        n, _ := toInt64(v)
        return n
}

// anyFloat64 extracts a float64 from an any, returning 0 if conversion fails.
func anyFloat64(v any) float64 {
        f, _ := toFloat64(v)
        return f
}

// ---------------------------------------------------------------------------
// DRY rendering helpers — shared across pipeline package
// ---------------------------------------------------------------------------

// keyHint pairs a keyboard key with its i18n label.
type keyHint struct {
        Key   string // Raw key string (e.g. "↵", "q", "1")
        Label string // i18n key or literal label
}

// renderKeyHints renders a row of key hints in consistent format.
// Each hint is: <CaptionStyle>key</CaptionStyle> <BodyStyle>label</BodyStyle>
// separated by 4 spaces.
func renderKeyHints(hints []keyHint) string {
        var b strings.Builder
        for i, h := range hints {
                if i > 0 {
                        b.WriteString("    ")
                }
                b.WriteString(style.CaptionStyle.Render(h.Key))
                b.WriteString(" ")
                b.WriteString(style.BodyStyle.Render(h.Label))
        }
        b.WriteString("\n")
        return b.String()
}

// renderSeparator is deprecated — the design system is "vertical borderless".
// Use style.Section(style.SectionGap) instead for blank-line spacing.
// Kept only to satisfy any external references; do NOT use in new code.
func renderSeparator(width int) string {
        return style.Section(style.SectionGap)
}

// renderNicheStats renders a one-line stat row for a niche: "key1: val1 · key2: val2".
func renderNicheStats(pairs []string) string {
        if len(pairs)%2 != 0 || len(pairs) == 0 {
                return ""
        }
        var parts []string
        for i := 0; i < len(pairs); i += 2 {
                parts = append(parts, fmt.Sprintf("%s: %s", pairs[i], pairs[i+1]))
        }
        return strings.Join(parts, " · ") + "\n"
}
