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
)

// ---------------------------------------------------------------------------
// Follow-up data types
// ---------------------------------------------------------------------------

// FollowUpLead represents a lead in the follow-up queue.
type FollowUpLead struct {
        // BusinessName is the lead's business name.
        BusinessName string

        // Phase is the follow-up phase using protocol.FollowUpPhase constants.
        Phase protocol.FollowUpPhase

        // IceBreakerAction describes when the ice breaker was sent (e.g. "5 hari lalu").
        IceBreakerAction string

        // FollowUp1Action describes when follow-up 1 was sent (e.g. "3 hari lalu").
        FollowUp1Action string

        // FollowUp2Action describes when follow-up 2 was sent (e.g. "1 hari lalu").
        FollowUp2Action string

        // PreviousAction describes when the previous step was taken (e.g. "ice breaker: 2 hari lalu").
        PreviousAction string

        // NextAction describes when the next step is due (e.g. "follow-up hari ini").
        NextAction string

        // SlotNumber is the WA slot being used (for sending state).
        SlotNumber int

        // VariantName is the template variant name (for sending state).
        VariantName string

        // IsSending indicates the lead is currently being sent a follow-up.
        IsSending bool

        // WaitTime is the remaining wait time string (e.g. "14m 02s").
        WaitTime string

        // PreviousResponse is the lead's previous response text (for recontact).
        PreviousResponse string

        // DaysSinceLastAction is days since the last action was taken.
        DaysSinceLastAction int

        // IsSelected indicates cursor is on this lead.
        IsSelected bool
}

// NicheGroup holds follow-up leads grouped by niche.
type NicheGroup struct {
        // NicheName is the niche identifier (e.g. "web_developer").
        NicheName string

        // Leads contains the follow-up leads for this niche.
        Leads []FollowUpLead

        // FU1Count is the count of follow-up 1 leads.
        FU1Count int

        // FU2Count is the count of follow-up 2 leads.
        FU2Count int

        // ColdCount is the count of cold leads.
        ColdCount int
}

// RecontactLead represents a lead eligible for re-contact.
type RecontactLead struct {
        // BusinessName is the lead's business name.
        BusinessName string

        // PreviousResponse is the last response from the lead.
        PreviousResponse string

        // DaysSinceResponse is how many days since the lead responded.
        DaysSinceResponse int

        // DaysSinceOffer is how many days since the offer was sent.
        DaysSinceOffer int

        // CanRecontact indicates the lead is eligible for re-contact today.
        CanRecontact bool
}

// ---------------------------------------------------------------------------
// FollowUp Model
// ---------------------------------------------------------------------------

// FollowUp manages follow-up scheduling and execution for cold leads.
//
// States (from doc/09-screens-communicate.md):
//   - FollowUpDashboard: overview of all follow-ups with ambient data rain
//   - FollowUpNicheDetail: per-niche follow-up detail
//   - FollowUpSending: currently sending follow-up messages
//   - FollowUpEmpty: no follow-ups needed today
//   - FollowUpColdList: cold leads that got 2x follow-ups
//   - FollowUpRecontact: leads that responded before but went quiet
type FollowUp struct {
        screenBase
        state    protocol.StateID
        bus      *bus.Bus
        focused  bool

        // receivedNavigate is set when HandleNavigate is called, to prevent
        // handleDataUpdate from overriding a backend-driven state transition.
        receivedNavigate bool

        // Dashboard data.
        Niches               []NicheGroup
        TotalToday           int
        ColdTotal            int
        IceBreakerUnanswered int

        // Niche detail data.
        SelectedNicheIndex int
        SelectedLeadIndex  int

        // Sending state data.
        SendingNicheIndex int
        SendingRate       string
        SendingDoneCount  int
        SendingTotalCount int

        // Cold list data.
        ColdLeads []FollowUpLead

        // Recontact data.
        RecontactLeads []RecontactLead

        // Variant names for display.
        VariantNames []string

        // VariantPreviews holds the preview text for each variant, provided by backend.
        VariantPreviews []string

        // variantManualOnly tracks which variant indices are manual-only.
        // Populated from backend data — the TUI does not hardcode this rule.
        variantManualOnly map[int]bool

        // VariantManualOnlyIdx is the single variant index that is manual-only,
        // provided by the backend. -1 means not set.
        VariantManualOnlyIdx int

        // maxSendingVisible is the backend-configurable limit for visible sending leads.
        maxSendingVisible int
        // maxColdVisible is the backend-configurable limit for visible cold leads.
        maxColdVisible int

        // Ambient effects.
        DataRain component.DataRain

        // Breathing removed per 3G audit — data rain is the ambient effect.

        // Layout.
        Width  int
        Height int

        // Animation timing.
        LastUpdate time.Time
        AnimStart  time.Time

        // Progress bar for sending state.
        ProgressBar component.ProgressBar
}

// NewFollowUp creates a FollowUp screen in dashboard state.
func NewFollowUp() *FollowUp {
        return &FollowUp{
                screenBase:              screenBase{id: protocol.ScreenFollowUp},
                state:             protocol.FollowUpDashboard,
                DataRain:          component.NewDataRain(defaultDataRainWidth),
                // VariantNames come from the backend via HandleUpdate.
                // No hardcoded defaults — the backend owns variant file names.
                VariantNames:        nil,
                variantManualOnly:   nil, // not set yet — backend provides this
                maxSendingVisible:   maxSendingLeadsVisible,
                maxColdVisible:      maxColdLeadsVisible,
        }
}

// ---------------------------------------------------------------------------
// Screen interface
// ---------------------------------------------------------------------------

func (f *FollowUp) SetBus(b *bus.Bus) { f.bus = b }

func (f *FollowUp) Focus() {
        f.focused = true
        f.AnimStart = time.Now()
        f.SelectedLeadIndex = 0
        f.SelectedNicheIndex = 0
}

func (f *FollowUp) Blur() { f.focused = false }

// HandleNavigate processes navigate commands from the backend.
func (f *FollowUp) HandleNavigate(params map[string]any) error {
        if state, ok := params[protocol.ParamState].(string); ok {
                f.state = protocol.StateID(state)
        }
        f.receivedNavigate = true
        return f.handleDataUpdate(params)
}

// HandleUpdate processes update commands from the backend.
func (f *FollowUp) HandleUpdate(params map[string]any) error {
        return f.handleDataUpdate(params)
}

// handleDataUpdate is the shared data ingestion path for both Navigate and Update.
func (f *FollowUp) handleDataUpdate(params map[string]any) error {
        if rawNiches, ok := params[protocol.ParamNiches].([]any); ok {
                f.Niches = parseNicheGroups(rawNiches)
                f.TotalToday = 0
                f.ColdTotal = 0
                for _, n := range f.Niches {
                        f.TotalToday += n.FU1Count + n.FU2Count
                        f.ColdTotal += n.ColdCount
                }
        }

        // Only override computed totals if backend explicitly provides them.
        // intVal returns 0 for missing keys, which would wipe out the niche-derived totals.
        if v, ok := params[protocol.ParamTotalToday]; ok {
                if f64, ok := v.(float64); ok {
                        f.TotalToday = int(f64)
                }
        }
        if v, ok := params[protocol.ParamColdTotal]; ok {
                if f64, ok := v.(float64); ok {
                        f.ColdTotal = int(f64)
                }
        }

        if v, ok := params[protocol.ParamIceBreakerUnanswered]; ok {
                if f64, ok := v.(float64); ok {
                        f.IceBreakerUnanswered = int(f64)
                }
        }

        if rawCold, ok := params[protocol.ParamColdLeads].([]any); ok {
                f.ColdLeads = parseColdLeads(rawCold)
        }

        if rawRecontact, ok := params[protocol.ParamRecontactLeads].([]any); ok {
                f.RecontactLeads = parseRecontactLeads(rawRecontact)
        }

        if rate, ok := params[protocol.ParamSendingRate].(string); ok {
                f.SendingRate = rate
        }
        f.SendingDoneCount = intVal(params, protocol.ParamSendingDone)
        f.SendingTotalCount = intVal(params, protocol.ParamSendingTotal)

        // Variant names from backend — replaces any previous values.
        if rawVariants, ok := params[protocol.ParamVariantNames].([]any); ok {
                f.VariantNames = make([]string, 0, len(rawVariants))
                for _, v := range rawVariants {
                        if s, ok := v.(string); ok {
                                f.VariantNames = append(f.VariantNames, s)
                        }
                }
        }
        // Variant previews from backend.
        if rawPreviews, ok := params[protocol.ParamVariantPreviews].([]any); ok {
                f.VariantPreviews = make([]string, 0, len(rawPreviews))
                for _, v := range rawPreviews {
                        if s, ok := v.(string); ok {
                                f.VariantPreviews = append(f.VariantPreviews, s)
                        }
                }
        }

        // Per-variant manual-only flags from backend.
        // The backend decides which variants are manual-only, not the TUI.
        if rawManualOnly, ok := params[protocol.ParamVariantManualOnly].([]any); ok {
                f.variantManualOnly = make(map[int]bool, len(rawManualOnly))
                for _, v := range rawManualOnly {
                        if f64, ok := v.(float64); ok {
                                f.variantManualOnly[int(f64)] = true
                        }
                }
        }

        // Backend-configurable visible limits — replaces hardcoded constants.
        if mv, ok := params[protocol.ParamMaxSendingVisible]; ok {
                if f64, ok := mv.(float64); ok {
                        f.maxSendingVisible = int(f64)
                }
        }
        if mv, ok := params[protocol.ParamMaxColdVisible]; ok {
                if f64, ok := mv.(float64); ok {
                        f.maxColdVisible = int(f64)
                }
        }

        // Fallback: suggest empty state when dashboard has no follow-up data.
        // The backend is the authoritative state machine and should send a
        // navigate command with state=followup_empty. This local hint ensures
        // the TUI shows the correct view even when the backend hasn't
        // explicitly navigated yet (e.g. during initial data load).
        if f.state == protocol.FollowUpDashboard && f.TotalToday == 0 && len(f.ColdLeads) == 0 && len(f.RecontactLeads) == 0 && !f.receivedNavigate {
                f.state = protocol.FollowUpEmpty
        }

        return nil
}

// Init implements tea.Model.
func (f *FollowUp) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (f *FollowUp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        now := time.Now()

        // Advance data rain.
        if f.DataRain.Tick(now) {
                f.LastUpdate = now
        }

        switch m := msg.(type) {
        case tea.WindowSizeMsg:
                f.Width = m.Width
                f.Height = m.Height
                f.DataRain = component.NewDataRain(m.Width)
                return f, nil

        case tea.KeyMsg:
                return f.handleKey(m)
        }

        return f, nil
}

// ---------------------------------------------------------------------------
// Key handling
// ---------------------------------------------------------------------------

func (f *FollowUp) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // Pause data rain on interaction.
        f.DataRain.Pause()

        switch f.state {
        case protocol.FollowUpDashboard:
                return f.handleDashboardKey(msg)
        case protocol.FollowUpNicheDetail:
                return f.handleNicheDetailKey(msg)
        case protocol.FollowUpSending:
                return f.handleSendingKey(msg)
        case protocol.FollowUpEmpty:
                return f.handleEmptyKey(msg)
        case protocol.FollowUpColdList:
                return f.handleColdListKey(msg)
        case protocol.FollowUpRecontact:
                return f.handleRecontactKey(msg)
        }
        return f, nil
}

func (f *FollowUp) handleDashboardKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // "q" → back (let global handler pop screen).
        if msg.String() == "q" {
                return f, nil
        }

        switch msg.Type {
        case tea.KeyEnter:
                // View detail of selected niche.
                if f.SelectedNicheIndex < len(f.Niches) {
                        f.state = protocol.FollowUpNicheDetail
                        f.SelectedLeadIndex = 0
                        f.AnimStart = time.Now()
                }
                return f, nil

        case tea.KeyUp:
                if f.SelectedNicheIndex > 0 {
                        f.SelectedNicheIndex--
                }
                return f, nil

        case tea.KeyDown:
                if f.SelectedNicheIndex < len(f.Niches)-1 {
                        f.SelectedNicheIndex++
                }
                return f, nil
        }

        // "a" → auto-approve all follow-ups.
        if msg.String() == "a" {
                if f.bus != nil {
                        f.bus.Publish(bus.ActionMsg{
                                Action: string(protocol.ActionFollowUpAutoAll),
                                Screen: protocol.ScreenFollowUp,
                        })
                }
                return f, nil
        }

        return f, nil
}

func (f *FollowUp) handleNicheDetailKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        currentNiche := f.currentNiche()
        if currentNiche == nil {
                f.state = protocol.FollowUpDashboard
                return f, nil
        }

        // "q" → back to dashboard.
        if msg.String() == "q" {
                f.state = protocol.FollowUpDashboard
                f.AnimStart = time.Now()
                return f, nil
        }

        switch msg.Type {
        case tea.KeyUp:
                if f.SelectedLeadIndex > 0 {
                        f.SelectedLeadIndex--
                }
                return f, nil

        case tea.KeyDown:
                if f.SelectedLeadIndex < len(currentNiche.Leads)-1 {
                        f.SelectedLeadIndex++
                }
                return f, nil

        case tea.KeyEnter:
                // View lead detail (for future expansion).
                return f, nil
        }

        // "a" → auto-approve all.
        if msg.String() == "a" {
                if f.bus != nil {
                        f.bus.Publish(bus.ActionMsg{
                                Action: string(protocol.ActionFollowUpAutoAll),
                                Screen: protocol.ScreenFollowUp,
                                Params: map[string]any{protocol.ParamNiche: currentNiche.NicheName},
                        })
                }
                return f, nil
        }

        return f, nil
}

func (f *FollowUp) handleSendingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // "q" → back to dashboard.
        if msg.String() == "q" {
                f.state = protocol.FollowUpDashboard
                f.AnimStart = time.Now()
                return f, nil
        }

        switch msg.Type {
        case tea.KeyTab:
                // Switch niche.
                if f.SendingNicheIndex < len(f.Niches)-1 {
                        f.SendingNicheIndex++
                } else {
                        f.SendingNicheIndex = 0
                }
                return f, nil

        case tea.KeyEnter:
                // Skip wait for current lead.
                if f.bus != nil {
                        f.bus.Publish(bus.ActionMsg{
                                Action: string(protocol.ActionFollowUpSkipWait),
                                Screen: protocol.ScreenFollowUp,
                        })
                }
                return f, nil
        }

        // "p" → pause sending.
        if msg.String() == "p" {
                if f.bus != nil {
                        f.bus.Publish(bus.ActionMsg{
                                Action: string(protocol.ActionFollowUpPause),
                                Screen: protocol.ScreenFollowUp,
                        })
                }
                return f, nil
        }

        return f, nil
}

func (f *FollowUp) handleEmptyKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // "q" → back (let global handler pop screen).
        if msg.String() == "q" {
                return f, nil
        }

        switch msg.Type {
        case tea.KeyEnter:
                // View cold leads.
                if len(f.ColdLeads) > 0 {
                        f.state = protocol.FollowUpColdList
                        f.AnimStart = time.Now()
                }
                return f, nil
        }
        return f, nil
}

func (f *FollowUp) handleColdListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // "q" → back to dashboard.
        if msg.String() == "q" {
                f.state = protocol.FollowUpDashboard
                f.AnimStart = time.Now()
                return f, nil
        }

        switch msg.Type {
        case tea.KeyUp:
                if f.SelectedLeadIndex > 0 {
                        f.SelectedLeadIndex--
                }
                return f, nil

        case tea.KeyDown:
                if f.SelectedLeadIndex < len(f.ColdLeads)-1 {
                        f.SelectedLeadIndex++
                }
                return f, nil

        case tea.KeyEnter:
                // Send final follow-up (3rd).
                if f.SelectedLeadIndex < len(f.ColdLeads) {
                        if f.bus != nil {
                                f.bus.Publish(bus.ActionMsg{
                                        Action: string(protocol.ActionFollowUpSendFinal),
                                        Screen: protocol.ScreenFollowUp,
                                        Params: map[string]any{
                                                protocol.ParamBusinessName: f.ColdLeads[f.SelectedLeadIndex].BusinessName,
                                        },
                                })
                        }
                }
                return f, nil
        }

        // "a" → archive all cold leads.
        if msg.String() == "a" {
                if f.bus != nil {
                        f.bus.Publish(bus.ActionMsg{
                                Action: string(protocol.ActionFollowUpArchiveCold),
                                Screen: protocol.ScreenFollowUp,
                        })
                }
                return f, nil
        }

        return f, nil
}

func (f *FollowUp) handleRecontactKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // "q" → back to dashboard.
        if msg.String() == "q" {
                f.state = protocol.FollowUpDashboard
                f.AnimStart = time.Now()
                return f, nil
        }

        switch msg.Type {
        case tea.KeyUp:
                if f.SelectedLeadIndex > 0 {
                        f.SelectedLeadIndex--
                }
                return f, nil

        case tea.KeyDown:
                if f.SelectedLeadIndex < len(f.RecontactLeads)-1 {
                        f.SelectedLeadIndex++
                }
                return f, nil

        case tea.KeyEnter:
                // Send re-contact.
                if f.SelectedLeadIndex < len(f.RecontactLeads) {
                        if f.bus != nil {
                                f.bus.Publish(bus.ActionMsg{
                                        Action: string(protocol.ActionFollowUpRecontact),
                                        Screen: protocol.ScreenFollowUp,
                                        Params: map[string]any{
                                                protocol.ParamBusinessName: f.RecontactLeads[f.SelectedLeadIndex].BusinessName,
                                        },
                                })
                                }
                                return f, nil
                        }

                        // "a" → auto re-contact all.
                        if msg.String() == "a" {
                                if f.bus != nil {
                                        f.bus.Publish(bus.ActionMsg{
                                                Action: string(protocol.ActionFollowUpRecontactAll),
                                                Screen: protocol.ScreenFollowUp,
                                        })
                                }
                                return f, nil
                        }
                }

        return f, nil
}

// ---------------------------------------------------------------------------
// View rendering
// ---------------------------------------------------------------------------

// View renders the follow-up screen.
func (f *FollowUp) View() string {
        switch f.state {
        case protocol.FollowUpDashboard:
                return f.viewDashboard()
        case protocol.FollowUpNicheDetail:
                return f.viewNicheDetail()
        case protocol.FollowUpSending:
                return f.viewSending()
        case protocol.FollowUpEmpty:
                return f.viewEmpty()
        case protocol.FollowUpColdList:
                return f.viewColdList()
        case protocol.FollowUpRecontact:
                return f.viewRecontact()
        default:
                return f.viewDashboard()
        }
}

// viewDashboard renders the followup_dashboard state.
//
// Spec:
//
//      follow-up                                 ● auto-jalan
//      ░░ 2 8 4 1 ░░ 7 3 9 ░░ 5 6 ░░
//      antrian follow-up hari ini                 14 pesan
//      ▸ web_developer
//        follow-up 1    8 lead (ice breaker 2 hari lalu, belum jawab)
//        follow-up 2    3 lead (follow-up 1 kemarin, belum jawab)
//        dingin          2 lead (2x follow-up, masih diam)
//      ▸ undangan_digital
//      ...
func (f *FollowUp) viewDashboard() string {
        var b strings.Builder

        // Title with auto-running indicator.
        b.WriteString(renderTitleWithIndicator(i18n.T(i18n.KeyFollowUpTitle), i18n.T(i18n.KeyFollowUpAutoRunning), style.SuccessStyle, f.Width))
        b.WriteString("\n\n")

        // Ambient data rain.
        b.WriteString(f.DataRain.ViewFormatted())
        b.WriteString("\n\n")

        // Separator.
        writeSeparator(&b, f.Width)

        // Queue title with count — doc: "antrian follow-up hari ini                 14 pesan"
        queueTitle := i18n.T(i18n.KeyFollowUpQueueToday)
        queueCount := fmt.Sprintf("%d %s", f.TotalToday, i18n.T(i18n.KeyFollowUpMessagesLabel))
        b.WriteString(renderTitleWithIndicator(queueTitle, queueCount, style.HeadingStyle, f.Width))
        b.WriteString("\n\n")

        // Niche groups.
        for i, niche := range f.Niches {
                selected := i == f.SelectedNicheIndex

                // Niche name with arrow.
                arrow := arrowIndicator
                nicheStyle := style.BodyStyle
                if selected {
                        nicheStyle = style.AccentStyle
                }
                b.WriteString(nicheStyle.Render(arrow + " " + niche.NicheName))
                b.WriteString("\n")

                // Follow-up 1 row — doc: "follow-up 1    8 lead (ice breaker 2 hari lalu, belum jawab)"
                if niche.FU1Count > 0 {
                        b.WriteString(style.Indent(2))
                        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyFollowUpFU1Label)))
                        b.WriteString("    ")
                        b.WriteString(style.BodyStyle.Render(fmt.Sprintf("%d %s", niche.FU1Count, i18n.T(i18n.KeyFollowUpLeadNoun))))
                        b.WriteString(" ")
                        b.WriteString(style.CaptionStyle.Render(fmt.Sprintf("(%s)", i18n.T(i18n.KeyFollowUpFU1Detail))))
                        b.WriteString("\n")
                }

                // Follow-up 2 row — doc: "follow-up 2    3 lead (follow-up 1 kemarin, belum jawab)"
                if niche.FU2Count > 0 {
                        b.WriteString(style.Indent(2))
                        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyFollowUpFU2Label)))
                        b.WriteString("    ")
                        b.WriteString(style.BodyStyle.Render(fmt.Sprintf("%d %s", niche.FU2Count, i18n.T(i18n.KeyFollowUpLeadNoun))))
                        b.WriteString(" ")
                        b.WriteString(style.CaptionStyle.Render(fmt.Sprintf("(%s)", i18n.T(i18n.KeyFollowUpFU2Detail))))
                        b.WriteString("\n")
                }

                // Cold row — doc: "dingin          2 lead (2x follow-up, masih diam)"
                if niche.ColdCount > 0 {
                        b.WriteString(style.Indent(2))
                        b.WriteString(style.BadgeColdStyle.Render(i18n.T(i18n.KeyBadgeCold)))
                        b.WriteString("          ")
                        b.WriteString(style.CaptionStyle.Render(fmt.Sprintf("%d %s", niche.ColdCount, i18n.T(i18n.KeyFollowUpLeadNoun))))
                        b.WriteString(" ")
                        b.WriteString(style.CaptionStyle.Render(fmt.Sprintf("(%s)", i18n.T(i18n.KeyFollowUpColdDashboard))))
                        b.WriteString("\n")
                }

                b.WriteString("\n")
        }

        // Separator.
        writeSeparator(&b, f.Width)

        // Total line — doc: "total: 14 follow-up hari ini · 4 dingin"
        b.WriteString(style.MutedStyle.Render(
                fmt.Sprintf("%s: %d %s %s · %d %s",
                        i18n.T(i18n.KeyFollowUpTotal), f.TotalToday, i18n.T(i18n.KeyFollowUpTitle),
                        strings.ToLower(i18n.T(i18n.KeyHistoryToday)),
                        f.ColdTotal, i18n.T(i18n.KeyFollowUpTotalCold),
                ),
        ))
        b.WriteString("\n")
        // Doc: "semua auto-kirim pas jam kerja. varian beda tiap follow-up. nggak ada yang sama."
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyFollowUpAutoSchedule)))
        b.WriteString(" ")
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyFollowUpVariantNote)))
        b.WriteString(" ")
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyFollowUpNoDuplicate)))
        b.WriteString("\n\n")

        // Actions — doc: "↵  liat detail    a  auto-semua    q  balik"
        actions := []string{
                style.ActionStyle.Render(i18n.T(i18n.KeyFollowUpViewDetail)),
                style.ActionStyle.Render(i18n.T(i18n.KeyFollowUpAutoAll)),
                style.MutedStyle.Render(i18n.T(i18n.KeyLabelBack)),
        }
        b.WriteString(strings.Join(actions, "    "))
        b.WriteString("\n\n")

        // Auto-note.
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyFollowUpAutoNote)))

        return b.String()
}

// viewNicheDetail renders the followup_niche_detail state.
func (f *FollowUp) viewNicheDetail() string {
        niche := f.currentNiche()
        if niche == nil {
                return f.viewDashboard()
        }

        var b strings.Builder

        // Title.
        b.WriteString(style.HeadingStyle.Render(
                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyFollowUpTitle), niche.NicheName),
        ))
        b.WriteString("\n\n")

        // Separator.
        writeSeparator(&b, f.Width)

        // Follow-up 1 leads.
        fu1Leads := filterLeadsByPhase(niche.Leads, protocol.FUPhase1)
        if len(fu1Leads) > 0 {
                b.WriteString(style.SubHeadingStyle.Render(fmt.Sprintf("%s (%d %s)", i18n.T(i18n.KeyFollowUpFU1Label), len(fu1Leads), i18n.T(i18n.KeyFollowUpLeadNoun))))
                b.WriteString("\n\n")
                for i, lead := range fu1Leads {
                        f.renderLeadLine(&b, lead, i, 1)
                }
                b.WriteString("\n")
        }

        // Follow-up 2 leads.
        fu2Leads := filterLeadsByPhase(niche.Leads, protocol.FUPhase2)
        if len(fu2Leads) > 0 {
                b.WriteString(style.SubHeadingStyle.Render(fmt.Sprintf("%s (%d %s)", i18n.T(i18n.KeyFollowUpFU2Label), len(fu2Leads), i18n.T(i18n.KeyFollowUpLeadNoun))))
                b.WriteString("\n\n")
                for i, lead := range fu2Leads {
                        f.renderLeadLine(&b, lead, i, 2)
                }
                b.WriteString("\n")
        }

        // Cold leads.
        coldLeads := filterLeadsByPhase(niche.Leads, protocol.FUPhaseCold)
        if len(coldLeads) > 0 {
                b.WriteString(style.MutedStyle.Render(fmt.Sprintf("%s (%d %s)", i18n.T(i18n.KeyFollowUpCold), len(coldLeads), i18n.T(i18n.KeyFollowUpLeadNoun))))
                b.WriteString("\n\n")
                for i, lead := range coldLeads {
                        b.WriteString(renderLineNumber(i + 1))
                        b.WriteString(style.CaptionStyle.Render(lead.BusinessName))
                        b.WriteString("    ")
                        b.WriteString(style.CaptionStyle.Render(lead.PreviousAction))
                        b.WriteString("  ")
                        b.WriteString(style.BadgeColdStyle.Render(i18n.T(i18n.KeyBadgeCold)))
                        b.WriteString("\n")
                }
                b.WriteString("\n")
        }

        // Separator.
        writeSeparator(&b, f.Width)

        // Variant names with preview text.
        // Doc: ▸ follow_up_1.md → "halo kak, cuma ngingetin aja"
        //      ▸ follow_up_2.md → "kak, penawaran terbatas nih"
        //      ▸ follow_up_3.md → "terakhir kak, kalo berkenan"  (hanya buat manual)
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyFollowUpVariantLabel) + ":"))
        b.WriteString("\n")
        for i, vn := range f.VariantNames {
                preview := ""
                if i < len(f.VariantPreviews) {
                        preview = f.VariantPreviews[i]
                }
                line := fmt.Sprintf("▸ %s → %q", vn, preview)
                // Manual-only flag is backend-driven, not hardcoded.
                if f.variantManualOnly[i] {
                        line += fmt.Sprintf("  %s", i18n.T(i18n.KeyFollowUpVariantManualOnly))
                }
                b.WriteString(style.CaptionStyle.Render(line))
                b.WriteString("\n")
        }
        b.WriteString("\n")

        // Actions — doc: "↑↓  pilih lead    ↵  liat detail    a  auto-semua    q  balik"
        actions := []string{
                renderSelectAction(i18n.T(i18n.KeyFollowUpLeadNoun)),
                style.ActionStyle.Render(i18n.T(i18n.KeyFollowUpViewDetail)),
                style.ActionStyle.Render(i18n.T(i18n.KeyFollowUpAutoAll)),
                style.MutedStyle.Render(i18n.T(i18n.KeyLabelBack)),
        }
        b.WriteString(strings.Join(actions, "    "))

        return b.String()
}

// viewSending renders the followup_sending state.
func (f *FollowUp) viewSending() string {
        var b strings.Builder

        // Title with sending indicator.
        b.WriteString(renderTitleWithIndicator(i18n.T(i18n.KeyFollowUpTitle), i18n.T(i18n.KeyFollowUpSending), style.SuccessStyle, f.Width))
        b.WriteString("\n\n")

        // Separator.
        writeSeparator(&b, f.Width)

        // Show the niche being processed.
        if f.SendingNicheIndex < len(f.Niches) {
                niche := f.Niches[f.SendingNicheIndex]
                b.WriteString(style.AccentStyle.Render(fmt.Sprintf("▸ %s", niche.NicheName)))
                b.WriteString("\n\n")

                // Show leads being sent.
                for i, lead := range niche.Leads {
                        if i >= f.maxSendingVisible {
                                b.WriteString(style.CaptionStyle.Render("..."))
                                break
                        }

                        b.WriteString(renderLineNumber(i + 1))

                        if lead.IsSending {
                                // Sending: show progress.
                                b.WriteString(style.ActionStyle.Render("→"))
                                b.WriteString("  ")
                                b.WriteString(style.BodyStyle.Render(lead.BusinessName))
                                b.WriteString("    ")
                                b.WriteString(style.CaptionStyle.Render(fmt.Sprintf("📱 %s-%d", i18n.T(i18n.KeyFollowUpSlot), lead.SlotNumber)))
                                b.WriteString("\n")
                                b.WriteString(style.Indent(4))
                                b.WriteString(style.MutedStyle.Render(fmt.Sprintf("%s: %s", lead.Phase, lead.VariantName)))
                                b.WriteString("  ")
                                b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyFollowUpDiffVariant)))
                                b.WriteString("\n")
                                b.WriteString(style.Indent(4))

                                // Progress bar.
                                f.ProgressBar = component.NewProgressBar(calcSepWidth(f.Width) - progressBarPadding)
                                if f.SendingTotalCount > 0 {
                                        f.ProgressBar.Percent = float64(f.SendingDoneCount) / float64(f.SendingTotalCount)
                                }
                                f.ProgressBar.FillChar = "\u2588"
                                f.ProgressBar.EmptyChar = "░"
                                b.WriteString(f.ProgressBar.View())
                                b.WriteString(" ")
                                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyFollowUpSendingNow) + "..."))
                                b.WriteString("\n\n")
                        } else {
                                // Waiting.
                                b.WriteString(style.CaptionStyle.Render("   "))
                                b.WriteString(style.BodyStyle.Render(lead.BusinessName))
                                b.WriteString("    ")
                                b.WriteString(style.CaptionStyle.Render(fmt.Sprintf("📱 %s-%d", i18n.T(i18n.KeyFollowUpSlot), lead.SlotNumber)))
                                b.WriteString("\n")
                                b.WriteString(style.Indent(4))
                                b.WriteString(style.MutedStyle.Render(fmt.Sprintf("%s: %s", lead.Phase, lead.VariantName)))
                                b.WriteString("  ")
                                b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyFollowUpRotation)))
                                b.WriteString("\n")
                                b.WriteString(style.Indent(4))
                                b.WriteString(style.CaptionStyle.Render(
                                        fmt.Sprintf("%s (%s: %s)", i18n.T(i18n.KeyFollowUpWaitNext), i18n.T(i18n.KeyFollowUpNextLabel), lead.WaitTime),
                                ))
                                b.WriteString("\n\n")
                        }
                }
        }

        // Separator.
        writeSeparator(&b, f.Width)

        // Rate info — doc: "rate: 9/18 per jam · follow-up hari ini: 2/14"
        b.WriteString(style.MutedStyle.Render(
                fmt.Sprintf("%s: %s %s · %s %s: %d/%d",
                        i18n.T(i18n.KeyFollowUpRateInfo), f.SendingRate,
                        i18n.T(i18n.KeyFollowUpRatePerHour),
                        i18n.T(i18n.KeyFollowUpTitle), strings.ToLower(i18n.T(i18n.KeyHistoryToday)),
                        f.SendingDoneCount, f.SendingTotalCount,
                ),
        ))
        b.WriteString("\n\n")
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyFollowUpVariantNote)))
        b.WriteString("\n")
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyFollowUpVariantPerFu)))
        b.WriteString("\n\n")

        // Actions — doc: "p  pause    ↵  skip tunggu    tab  pindah niche    q  balik"
        actions := []string{
                style.MutedStyle.Render(i18n.T(i18n.KeyLabelPause)),
                style.ActionStyle.Render(i18n.T(i18n.KeyFollowUpSkipWait)),
                style.MutedStyle.Render(i18n.T(i18n.KeyFollowUpTabNiche)),
                style.MutedStyle.Render(i18n.T(i18n.KeyLabelBack)),
        }
        b.WriteString(strings.Join(actions, "    "))

        return b.String()
}

// viewEmpty renders the followup_empty state.
func (f *FollowUp) viewEmpty() string {
        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyFollowUpTitle)))
        b.WriteString("\n\n")

        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyFollowUpEmpty)))
        b.WriteString("\n")
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyFollowUpEmptyDesc)))
        b.WriteString("\n\n")

        // Doc: "lead yang belom jawab ice breaker: 8\nlead dingin: 4"
        b.WriteString(style.MutedStyle.Render(
                fmt.Sprintf("%s: %d", i18n.T(i18n.KeyFollowUpEmptyUnansweredIB), f.IceBreakerUnanswered),
        ))
        b.WriteString("\n")
        if f.ColdTotal > 0 {
                b.WriteString(style.MutedStyle.Render(
                        fmt.Sprintf("%s: %d", i18n.T(i18n.KeyFollowUpColdCountLabel), f.ColdTotal),
                ))
                b.WriteString("\n\n")
        }

        // Actions — doc: "↵  liat lead dingin    q  balik"
        actions := []string{
                style.ActionStyle.Render(i18n.T(i18n.KeyFollowUpViewCold)),
                style.MutedStyle.Render(i18n.T(i18n.KeyLabelBack)),
        }
        b.WriteString(strings.Join(actions, "    "))

        return b.String()
}

// viewColdList renders the followup_cold_list state.
func (f *FollowUp) viewColdList() string {
        var b strings.Builder

        // Title with count — uses renderTitleWithIndicator for responsive layout.
        b.WriteString(renderTitleWithIndicator(i18n.T(i18n.KeyFollowUpColdList), fmt.Sprintf("%d %s", len(f.ColdLeads), i18n.T(i18n.KeyFollowUpLeadNoun)), style.CaptionStyle, f.Width))
        b.WriteString("\n\n")

        // Separator.
        writeSeparator(&b, f.Width)

        // Description.
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyFollowUpColdDesc)))
        b.WriteString("\n\n")

        // Cold lead rows — multi-line per-lead detail per doc.
        for i, lead := range f.ColdLeads {
                if i >= f.maxColdVisible {
                        b.WriteString(style.CaptionStyle.Render(fmt.Sprintf("...%d %s", len(f.ColdLeads)-f.maxColdVisible, i18n.T(i18n.KeyFollowUpOthers))))
                        break
                }

                b.WriteString(renderLineNumber(i + 1))
                b.WriteString(style.BodyStyle.Render(lead.BusinessName))

                // Line 1: ice breaker action.
                if lead.IceBreakerAction != "" {
                        b.WriteString("  ")
                        b.WriteString(style.CaptionStyle.Render(lead.IceBreakerAction))
                }
                b.WriteString("\n")

                // Line 2: follow-up 1 action.
                if lead.FollowUp1Action != "" {
                        b.WriteString(style.Indent(2))
                        b.WriteString(style.CaptionStyle.Render(
                                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyFollowUpFU1Label), lead.FollowUp1Action),
                        ))
                        b.WriteString("\n")
                }

                // Line 3: follow-up 2 action.
                if lead.FollowUp2Action != "" {
                        b.WriteString(style.Indent(2))
                        b.WriteString(style.CaptionStyle.Render(
                                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyFollowUpFU2Label), lead.FollowUp2Action),
                        ))
                        b.WriteString("\n")
                }

                // Line 4: cold status.
                b.WriteString(style.Indent(2))
                b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyFollowUpColdDetail)))
                b.WriteString("\n\n")
        }

        // Separator.
        writeSeparator(&b, f.Width)

        // Actions — doc: "↑↓  pilih    ↵  kirim follow-up terakhir (ke-3)    a  archive semua    q  balik"
        actions := []string{
                renderSelectAction(""),
                style.ActionStyle.Render(i18n.T(i18n.KeyFollowUpSendFinal)),
                style.MutedStyle.Render(i18n.T(i18n.KeyActionArchiveAll)),
                style.MutedStyle.Render(i18n.T(i18n.KeyLabelBack)),
        }
        b.WriteString(strings.Join(actions, "    "))
        b.WriteString("\n\n")

        // Warning.
        b.WriteString(style.WarningStyle.Render(i18n.T(i18n.KeyFollowUpColdWarning)))

        return b.String()
}

// viewRecontact renders the followup_recontact state.
func (f *FollowUp) viewRecontact() string {
        var b strings.Builder

        // Title.
        b.WriteString(style.HeadingStyle.Render(
                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyFollowUpTitle), i18n.T(i18n.KeyFollowUpRecontact)),
        ))
        b.WriteString("\n\n")

        // Separator.
        writeSeparator(&b, f.Width)

        // Description.
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyFollowUpRecontactDesc)))
        b.WriteString("\n\n")

        // Recontact lead rows.
        for i, lead := range f.RecontactLeads {
                selected := i == f.SelectedLeadIndex

                b.WriteString(renderLineNumber(i + 1))

                nameStyle := style.BodyStyle
                if selected {
                        nameStyle = style.SelectedBodyStyle
                }
                b.WriteString(nameStyle.Render(lead.BusinessName))

                b.WriteString("    ")
                b.WriteString(style.CaptionStyle.Render(
                        fmt.Sprintf("%s %d %s: %q",
                                i18n.T(i18n.KeyHistoryRespond), lead.DaysSinceResponse,
                                i18n.T(i18n.KeyFollowUpDaysAgo), lead.PreviousResponse,
                        ),
                ))
                b.WriteString("\n")

                b.WriteString(style.Indent(4))
                b.WriteString(style.CaptionStyle.Render(
                        fmt.Sprintf("%s %d %s", i18n.T(i18n.KeyFollowUpOfferSent), lead.DaysSinceOffer, i18n.T(i18n.KeyFollowUpDaysAgo)),
                ))
                b.WriteString("\n")

                b.WriteString(style.Indent(4))
                b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyFollowUpAfterThatSilent)))

                if lead.CanRecontact {
                        b.WriteString("\n")
                        b.WriteString(style.Indent(4))
                        b.WriteString(style.SuccessStyle.Render(i18n.T(i18n.KeyFollowUpCanRecontact)))
                }

                b.WriteString("\n\n")
        }

        // Separator.
        writeSeparator(&b, f.Width)

        // Note about templates — doc: "re-contact pakai template berbeda dari offer pertama.\ntone-nya lebih santai, \"hai lagi kak\" vibes."
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyFollowUpRecontactTemplate)))
        b.WriteString("\n")
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyFollowUpRecontactTone)))
        b.WriteString(" ")
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyFollowUpRecontactVibes)))
        b.WriteString("\n\n")

        // Actions — doc: "↑↓  pilih    ↵  kirim re-contact    a  auto-semua    q  balik"
        actions := []string{
                renderSelectAction(""),
                style.ActionStyle.Render(i18n.T(i18n.KeyFollowUpSendRecontact)),
                style.ActionStyle.Render(i18n.T(i18n.KeyFollowUpAutoAll)),
                style.MutedStyle.Render(i18n.T(i18n.KeyLabelBack)),
        }
        b.WriteString(strings.Join(actions, "    "))

        return b.String()
}

// ---------------------------------------------------------------------------
// Render helpers
// ---------------------------------------------------------------------------

// renderLeadLine renders a single follow-up lead line with number.
func (f *FollowUp) renderLeadLine(b *strings.Builder, lead FollowUpLead, index int, fuLevel int) {
        b.WriteString(renderLineNumber(index + 1))
        b.WriteString(style.BodyStyle.Render(lead.BusinessName))
        b.WriteString("    ")
        b.WriteString(style.CaptionStyle.Render(lead.PreviousAction))
        b.WriteString("   ")
        b.WriteString(style.MutedStyle.Render("→"))
        b.WriteString(" ")
        b.WriteString(style.BodyStyle.Render(lead.NextAction))
        b.WriteString("\n")
}

// currentNiche returns the currently selected niche group.
func (f *FollowUp) currentNiche() *NicheGroup {
        if f.SelectedNicheIndex < len(f.Niches) {
                return &f.Niches[f.SelectedNicheIndex]
        }
        return nil
}

// filterLeadsByPhase returns leads matching the given follow-up phase.
func filterLeadsByPhase(leads []FollowUpLead, phase protocol.FollowUpPhase) []FollowUpLead {
        var filtered []FollowUpLead
        for _, l := range leads {
                if l.Phase == phase {
                        filtered = append(filtered, l)
                }
        }
        return filtered
}

// ---------------------------------------------------------------------------
// Data parsing helpers
// ---------------------------------------------------------------------------

func parseNicheGroups(raw []any) []NicheGroup {
        var groups []NicheGroup
        for _, item := range raw {
                if m, ok := item.(map[string]any); ok {
                        // Support both "name" (backend standard) and "niche_name" (robustness).
                        nicheName := strVal(m, protocol.ParamName)
                        if nicheName == "" {
                                nicheName = strVal(m, protocol.ParamNicheName)
                        }
                        ng := NicheGroup{
                                NicheName: nicheName,
                                FU1Count:  intVal(m, protocol.ParamFU1Count),
                                FU2Count:  intVal(m, protocol.ParamFU2Count),
                                ColdCount: intVal(m, protocol.ParamColdCount),
                        }
                        if rawLeads, ok := m[protocol.ParamLeads].([]any); ok {
                                ng.Leads = parseFollowUpLeads(rawLeads)
                        }
                        groups = append(groups, ng)
                }
        }
        return groups
}

func parseFollowUpLeads(raw []any) []FollowUpLead {
        var leads []FollowUpLead
        for _, item := range raw {
                if m, ok := item.(map[string]any); ok {
                        leads = append(leads, FollowUpLead{
                                BusinessName:       strVal(m, protocol.ParamBusinessName),
                                Phase:              protocol.FollowUpPhase(strVal(m, protocol.ParamPhase)),
                                IceBreakerAction:   strVal(m, protocol.ParamIceBreakerAction),
                                FollowUp1Action:    strVal(m, protocol.ParamFollowUp1Action),
                                FollowUp2Action:    strVal(m, protocol.ParamFollowUp2Action),
                                PreviousAction:     strVal(m, protocol.ParamPreviousAction),
                                NextAction:         strVal(m, protocol.ParamNextAction),
                                SlotNumber:         intVal(m, protocol.ParamSlotNumberAlt),
                                VariantName:        strVal(m, protocol.ParamVariantName),
                                IsSending:          boolVal(m, protocol.ParamIsSending),
                                WaitTime:           strVal(m, protocol.ParamWaitTime),
                                DaysSinceLastAction: intVal(m, protocol.ParamDaysSinceLastAction),
                        })
                }
        }
        return leads
}

func parseColdLeads(raw []any) []FollowUpLead {
        return parseFollowUpLeads(raw)
}

func parseRecontactLeads(raw []any) []RecontactLead {
        var leads []RecontactLead
        for _, item := range raw {
                if m, ok := item.(map[string]any); ok {
                        leads = append(leads, RecontactLead{
                                BusinessName:     strVal(m, protocol.ParamBusinessName),
                                PreviousResponse: strVal(m, protocol.ParamPreviousResponse),
                                DaysSinceResponse: intVal(m, protocol.ParamDaysSinceResponse),
                                DaysSinceOffer:    intVal(m, protocol.ParamDaysSinceOffer),
                                CanRecontact:      boolVal(m, protocol.ParamCanRecontact),
                        })
                }
        }
        return leads
}

// String returns a debug representation.
func (f *FollowUp) String() string {
        return fmt.Sprintf("FollowUp{state=%s, niches=%d, total=%d}", f.state, len(f.Niches), f.TotalToday)
}
