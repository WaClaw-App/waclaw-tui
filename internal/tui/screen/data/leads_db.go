// Package data implements the Leads Database (Screen 9) and Template Manager
// (Screen 10) screens for the WaClaw TUI. Both screens share the data package
// because they form the "Data & Armory" domain group per the plan.
package data

import (
        "fmt"
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/WaClaw-App/waclaw/internal/tui/component"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        tui "github.com/WaClaw-App/waclaw/internal/tui"
        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/pkg/protocol"

        "github.com/charmbracelet/bubbles/key"
        tea "github.com/charmbracelet/bubbletea"
        "github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Data types — lead model for display
// ---------------------------------------------------------------------------

// Lead holds the display data for a single lead in the database view.
// This is a TUI-side display model — it does NOT replace the backend's
// lead.Lead domain type. It is populated from HandleNavigate/HandleUpdate
// params and rendered in the View.
type Lead struct {
        Name      string
        Category  string
        Address   string
        City      string
        Rating    float64
        Reviews   int
        HasWeb    bool
        HasInsta  bool
        PhotoCount int
        Score     int // 0-10
        Phase     protocol.LeadPhase
        Niche     string

        // Timeline events for detail view.
        Timeline []LeadEvent

        // Contact tracking.
        ContactCount  int // total messages sent (ice breaker + follow-ups)
        FollowUpCount int // number of follow-ups sent
        ResponseText  string
        IceBreakerTime string // e.g. "kemarin 09:15" — when ice breaker was sent
        ResponseTime   string // e.g. "kemarin 14:23" — when last response came

        // Follow-up timestamps for cold/follow-up-due variants.
        FollowUpTimes []string // e.g. ["3 hari lalu", "1 hari lalu"]

        // FollowUpDueText is the backend-provided due-date text for the next follow-up.
        // e.g. "hari ini", "besok", "3 hari lagi". Empty if not applicable.
        FollowUpDueText string

        // Conversion info (only for converted leads).
        Duration     string // e.g. "2 hari 1 jam"
        TemplateName string
        WorkerName   string
        Revenue      string // e.g. "rp 2.5jt"
}

// LeadEvent is a single timeline entry for a lead.
type LeadEvent struct {
        Time    time.Time
        Action  string // e.g. "ice_breaker_sent", "response_received"
        Detail  string // e.g. response text
}

// FilterCategory is a single filter row in the leads list view.
type FilterCategory struct {
        Label string
        Count int
        Phase protocol.LeadPhase
        Note  string // e.g. "(2x follow-up, belum jawab)" — backend-provided context note
}

// ---------------------------------------------------------------------------
// LeadsDB screen
// ---------------------------------------------------------------------------

// LeadsDB implements Screen 9: Leads Database → Archive.
//
// States: LeadsList, LeadsFiltered, LeadsFullDetail,
// LeadsFollowUpDue, LeadsCold, LeadsNeverContacted, LeadsConverted.
//
// The screen follows the "vertical borderless" design: hierarchy is expressed
// through brightness, size, and motion — never through borders or boxes.
type LeadsDB struct {
        tui.ScreenBase

        // state is the current screen state.
        state protocol.StateID

        // width and height track the terminal dimensions.
        width  int
        height int

        // leads is the full dataset populated by HandleNavigate/HandleUpdate.
        leads []Lead

        // filtered holds the leads matching the current filter.
        filtered []Lead

        // categories holds the filter categories with counts.
        categories []FilterCategory

        // cursor is the currently selected index (list or detail).
        cursor int

        // search is the fuzzy search input for filtering.
        search component.SearchInput

        // timeline is the component for rendering lead detail timelines.
        timeline component.Timeline

        // currentLead is the lead being viewed in detail states.
        currentLead *Lead

        // totalLeads is the overall lead count (shown in header).
        totalLeads int

        // activeFilter is the currently selected filter phase (empty = all).
        activeFilter protocol.LeadPhase

        // animStart tracks when the current state was entered (for stagger animations).
        animStart time.Time
}

// NewLeadsDB creates a new LeadsDB screen with default state.
func NewLeadsDB() *LeadsDB {
        return &LeadsDB{
                ScreenBase: tui.NewScreenBase(protocol.ScreenLeadsDB),
                state:      protocol.LeadsList,
                search:     component.NewSearchInput(int(anim.SearchDebounce / time.Millisecond)),
                timeline:   component.NewTimeline(),
        }
}

// Focus is called when the screen becomes active.
func (s *LeadsDB) Focus() {
        s.animStart = time.Now()
}

// Blur is called when the screen becomes inactive.
func (s *LeadsDB) Blur() {}

// HandleNavigate processes a "navigate" command from the backend.
// The params map carries the initial screen state and data.
func (s *LeadsDB) HandleNavigate(params map[string]any) error {
        s.animStart = time.Now()
        s.parseLeadsParams(params)
        return nil
}

// HandleUpdate processes an "update" command from the backend.
func (s *LeadsDB) HandleUpdate(params map[string]any) error {
        s.parseLeadsParams(params)

        // Re-apply filter when leads data changes via update.
        if _, ok := params[protocol.ParamLeads].([]any); ok {
                s.applyFilter()
        }

        return nil
}

// parseLeadsParams extracts shared data from HandleNavigate/HandleUpdate params.
// DRY: both handlers parse the same fields — this avoids duplicating the logic.
func (s *LeadsDB) parseLeadsParams(params map[string]any) {
        if stateStr, ok := params[protocol.ParamState].(string); ok {
                oldState := s.state
                s.state = protocol.StateID(stateStr)
                if s.state != oldState {
                        s.animStart = time.Now()
                }
        }

        if total, ok := params[protocol.ParamTotal].(float64); ok {
                s.totalLeads = int(total)
        }

        if cats, ok := params[protocol.ParamCategories].([]any); ok {
                s.categories = parseFilterCategories(cats)
        }

        if leadsData, ok := params[protocol.ParamLeads].([]any); ok {
                s.leads = parseLeads(leadsData)
                s.totalLeads = len(s.leads)
        }

        if leadData, ok := params[protocol.ParamLead].(map[string]any); ok {
                lead := parseLead(leadData)
                s.currentLead = &lead
                buildTimelineForLead(&s.timeline, &lead)
        }
}

// Update handles bubbletea messages.
func (s *LeadsDB) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        now := time.Now()

        switch m := msg.(type) {
        case tea.WindowSizeMsg:
                s.width = m.Width
                s.height = m.Height
                return s, nil

        case tea.KeyMsg:
                return s.handleKey(m, now)
        }

        return s, nil
}

// handleKey dispatches keyboard input based on the current state.
func (s *LeadsDB) handleKey(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
        switch s.state {
        case protocol.LeadsList:
                return s.handleKeyList(msg, now)
        case protocol.LeadsFiltered:
                return s.handleKeyFiltered(msg, now)
        default:
                // All detail/variant states share the same key handling.
                return s.handleKeyDetail(msg, now)
        }
}

// handleKeyList handles keys in the leads_list state.
func (s *LeadsDB) handleKeyList(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
        switch {
        case key.Matches(msg, tui.KeyUp):
                if s.cursor > 0 {
                        s.cursor--
                }
        case key.Matches(msg, tui.KeyDown):
                maxIdx := len(s.categories) - 1
                if s.cursor < maxIdx {
                        s.cursor++
                }
        case key.Matches(msg, tui.KeyEnter):
                if s.cursor < len(s.categories) {
                        cat := s.categories[s.cursor]
                        s.activeFilter = cat.Phase
                        s.applyFilter()
                        s.state = protocol.LeadsFiltered
                        s.cursor = 0
                        s.animStart = now
                }
        case key.Matches(msg, tui.KeySearch):
                s.search.SetValue("")
                s.search.Focused = true
                s.state = protocol.LeadsFiltered
                s.activeFilter = ""
                s.applyFilter()
                s.animStart = now
        }
        return s, nil
}

// handleKeyFiltered handles keys in the leads_filtered state.
func (s *LeadsDB) handleKeyFiltered(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
        switch {
        case key.Matches(msg, tui.KeyUp):
                if s.cursor > 0 {
                        s.cursor--
                }
        case key.Matches(msg, tui.KeyDown):
                maxIdx := len(s.filtered) - 1
                if s.cursor < maxIdx {
                        s.cursor++
                }
        case key.Matches(msg, tui.KeyEnter):
                if s.cursor < len(s.filtered) {
                        s.currentLead = &s.filtered[s.cursor]
                        s.state = resolveDetailState(s.currentLead)
                        s.cursor = 0
                        s.animStart = now
                        buildTimelineForLead(&s.timeline, s.currentLead)
                }
        case key.Matches(msg, tui.Key1):
                // Bulk send offer to all filtered leads.
                if s.Bus() != nil {
                        s.Bus().Publish(bus.ActionMsg{
                                Action: protocol.ActionBulkOffer,
                                Screen: protocol.ScreenLeadsDB,
                                Params: map[string]any{protocol.ParamPhase: string(s.activeFilter)},
                        })
                }
        case key.Matches(msg, tui.KeySearch):
                s.search.Focused = true
        }

        // Handle search input when focused.
        if s.search.Focused {
                switch {
                case key.Matches(msg, tui.KeyEscape):
                        s.search.Focused = false
                case key.Matches(msg, tui.KeyBack):
                        s.search.Focused = false
                        if s.search.Value != "" {
                                s.search.Clear()
                                s.applyFilter()
                        } else {
                                s.state = protocol.LeadsList
                                s.cursor = 0
                                s.animStart = now
                        }
                default:
                        s.handleSearchInput(msg)
                        s.applyFilter()
                }
        }

        return s, nil
}

// handleKeyDetail handles keys in all detail/variant states.
func (s *LeadsDB) handleKeyDetail(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
        if s.currentLead == nil {
                return s, nil
        }

        // q back — navigate back to filtered/list view
        if key.Matches(msg, tui.KeyBack) {
                if s.activeFilter != "" {
                        s.state = protocol.LeadsFiltered
                } else {
                        s.state = protocol.LeadsList
                }
                s.currentLead = nil
                s.animStart = now
                return s, nil
        }

        switch s.state {
        case protocol.LeadsFullDetail:
                return s.handleKeyFullDetail(msg, now)
        case protocol.LeadsFollowUpDue:
                return s.handleKeyFollowUpDue(msg, now)
        case protocol.LeadsCold:
                return s.handleKeyCold(msg, now)
        case protocol.LeadsNeverContacted:
                return s.handleKeyNeverContacted(msg, now)
        case protocol.LeadsConverted:
                // Converted leads only have "q back".
                return s, nil
        }

        return s, nil
}

// handleKeyFullDetail handles keys in the lead_full_detail state.
func (s *LeadsDB) handleKeyFullDetail(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
        switch {
        case key.Matches(msg, tui.Key1):
                s.publishAction(protocol.ActionSendOffer)
        case key.Matches(msg, tui.Key2):
                s.publishAction(protocol.ActionCustomReply)
        case key.Matches(msg, tui.Key3):
                s.publishAction(protocol.ActionLater)
        case key.Matches(msg, tui.Key4):
                s.publishAction(protocol.ActionMarkConvert)
        case key.Matches(msg, tui.Key5):
                s.publishAction(protocol.ActionArchive)
        case key.Matches(msg, tui.Key6):
                s.publishAction(protocol.ActionBlock)
        }
        return s, nil
}

// handleKeyFollowUpDue handles keys in the lead_follow_up_due state.
func (s *LeadsDB) handleKeyFollowUpDue(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
        switch {
        case key.Matches(msg, tui.Key1):
                s.publishAction(protocol.ActionSendFollowUp)
        case key.Matches(msg, tui.Key2):
                s.publishAction(protocol.ActionSkip)
        case key.Matches(msg, tui.Key3):
                s.publishAction(protocol.ActionMarkCold)
        }
        return s, nil
}

// handleKeyCold handles keys in the leads_cold state.
func (s *LeadsDB) handleKeyCold(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
        switch {
        case key.Matches(msg, tui.Key1):
                s.publishAction(protocol.ActionLastFollowUp)
        case key.Matches(msg, tui.Key2):
                s.publishAction(protocol.ActionMarkCold)
        case key.Matches(msg, tui.Key3):
                s.publishAction(protocol.ActionArchive)
        }
        return s, nil
}

// handleKeyNeverContacted handles keys in the leads_never_contacted state.
func (s *LeadsDB) handleKeyNeverContacted(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
        switch {
        case key.Matches(msg, tui.Key1):
                s.publishAction(protocol.ActionSendIceBreaker)
        case key.Matches(msg, tui.Key2):
                s.publishAction(protocol.ActionSkip)
        }
        return s, nil
}

// publishAction sends an action event through the bus.
func (s *LeadsDB) publishAction(action string) {
        if s.Bus() == nil || s.currentLead == nil {
                return
        }
        s.Bus().Publish(bus.ActionMsg{
                Action: action,
                Screen: protocol.ScreenLeadsDB,
                Params: map[string]any{
                        protocol.ParamLeadName: s.currentLead.Name,
                        protocol.ParamPhase:     string(s.currentLead.Phase),
                },
        })
}

// handleSearchInput processes character input for the search field.
func (s *LeadsDB) handleSearchInput(msg tea.KeyMsg) {
        switch msg.Type {
        case tea.KeyRunes:
                s.search.AppendChar(msg.String())
        case tea.KeyBackspace:
                s.search.Backspace()
        }
}

// applyFilter sets the filtered leads based on activeFilter and search.
func (s *LeadsDB) applyFilter() {
        s.filtered = nil
        query := s.search.Value

        for i := range s.leads {
                lead := &s.leads[i]

                // Phase filter.
                if s.activeFilter != "" && lead.Phase != s.activeFilter {
                        continue
                }

                // Search filter.
                if query != "" {
                        result := component.FuzzyMatch(query, lead.Name)
                        if !result.Matched {
                                continue
                        }
                }

                s.filtered = append(s.filtered, *lead)
        }
}

// View renders the current state of the LeadsDB screen.
func (s *LeadsDB) View() string {
        now := time.Now()

        switch s.state {
        case protocol.LeadsList:
                return s.viewList(now)
        case protocol.LeadsFiltered:
                return s.viewFiltered(now)
        case protocol.LeadsFullDetail:
                return s.viewFullDetail(now)
        case protocol.LeadsFollowUpDue:
                return s.viewFollowUpDue(now)
        case protocol.LeadsCold:
                return s.viewCold(now)
        case protocol.LeadsNeverContacted:
                return s.viewNeverContacted(now)
        case protocol.LeadsConverted:
                return s.viewConverted(now)
        default:
                return s.viewList(now)
        }
}

// viewList renders the leads_list state.
//
// Visual spec from doc/06:
//
//      database leads                          total: 847
//      filter: semua    / buat cari...
//      ▸ baru          23
//      ▸ ice breaker    128
//      ...
//      terbaru:
//      kopi nusantara      respond    ⭐ 4.2
//      ...
//      ↑↓  pindah    ↵  liat detail    /  cari    q  balik
func (s *LeadsDB) viewList(now time.Time) string {
        var b strings.Builder

        // Title line.
        title := i18n.T(i18n.KeyDataLeadsTitle)
        totalLabel := i18n.T(i18n.KeyDataTotal)
        b.WriteString(style.HeadingStyle.Render(title))
        b.WriteString(strings.Repeat(" ", max(4, s.width-len(title)-len(totalLabel)-len(fmt.Sprint(s.totalLeads))-4)))
        b.WriteString(style.CaptionStyle.Render(fmt.Sprintf("%s: %d", totalLabel, s.totalLeads)))
        b.WriteString("\n")

        b.WriteString(style.Separator())
        b.WriteString("\n")

        // Filter line — doc spec: "filter: semua    / buat cari..."
        filterLabel := i18n.T(i18n.KeyDataFilterAll)
        searchHint := i18n.T(i18n.KeyDataFilterSearch)
        b.WriteString(fmt.Sprintf("%s %s    %s %s",
                i18n.T(i18n.KeyDataFilterLabel),
                style.MutedStyle.Render(filterLabel),
                style.DimStyle.Render("/"),
                style.DimStyle.Render(searchHint),
        ))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SubSectionGap))

        // Category rows.
        for i, cat := range s.categories {
                label := cat.Label
                count := cat.Count

                // Stagger fade-in animation.
                visibleAt := s.animStart.Add(time.Duration(i) * anim.MenuStagger)
                if now.Before(visibleAt) {
                        continue
                }

                var prefix string
                if i == s.cursor {
                        prefix = style.AccentStyle.Render("▸ ")
                } else {
                        prefix = style.DimStyle.Render("  ")
                }

                countStr := style.CaptionStyle.Render(fmt.Sprint(count))

                // Render the label with phase-appropriate badge color.
                labelStyle := phaseLabelStyle(cat.Phase)
                b.WriteString(prefix)
                b.WriteString(labelStyle.Render(label))
                b.WriteString(strings.Repeat(" ", max(2, categoryColWidth-len(label))))
                b.WriteString(countStr)

                // Backend-provided context note, e.g. "(2x follow-up, belum jawab)".
                if cat.Note != "" {
                        b.WriteString("   ")
                        b.WriteString(style.DimStyle.Render(cat.Note))
                }

                b.WriteString("\n")
        }

        b.WriteString(style.Separator())
        b.WriteString("\n")

        // Recent leads.
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyDataRecent)))
        b.WriteString("\n")

        recentCount := min(maxRecentLeads, len(s.leads))
        for i := 0; i < recentCount; i++ {
                lead := s.leads[i]
                rating := fmt.Sprintf("⭐ %.1f", lead.Rating)
                phaseLabel := phaseBadgeText(lead.Phase)
                phaseStyle := phaseBadgeStyle(lead.Phase)

                b.WriteString(style.Indent(1))
                b.WriteString(style.BodyStyle.Render(lead.Name))
                b.WriteString("    ")
                b.WriteString(phaseStyle.Render(phaseLabel))
                b.WriteString("  ")
                b.WriteString(style.CaptionStyle.Render(rating))
                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SectionGap))

        // Footer actions.
        b.WriteString(s.viewFooterList())

        return b.String()
}

// viewFiltered renders the leads_filtered state.
//
// Visual spec from doc/06:
//
//      database leads                      respond: 41
//      01  kopi nusantara
//          cafe · kediri · ⭐ 4.2
//          ice breaker: kemarin 14:23
//          response: "iya kak, boleh lihat"
//          → belum dikirim offer
func (s *LeadsDB) viewFiltered(now time.Time) string {
        var b strings.Builder

        // Title line with filter count.
        title := i18n.T(i18n.KeyDataLeadsTitle)
        filterName := phaseLabelDisplay(s.activeFilter)
        countStr := fmt.Sprint(len(s.filtered))
        b.WriteString(style.HeadingStyle.Render(title))
        b.WriteString(strings.Repeat(" ", max(4, s.width-len(title)-len(filterName)-len(countStr)-6)))
        b.WriteString(style.MutedStyle.Render(fmt.Sprintf("%s: %s", filterName, countStr)))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // Search bar when active.
        if s.search.Focused || s.search.Value != "" {
                b.WriteString(s.search.View())
                b.WriteString("\n")
                b.WriteString(style.Section(style.SubSectionGap))
        }

        // Lead entries.
        for i, lead := range s.filtered {
                // Stagger fade-in.
                visibleAt := s.animStart.Add(time.Duration(i) * anim.MenuStagger)
                if now.Before(visibleAt) {
                        continue
                }

                // Number prefix.
                numStyle := style.CaptionStyle
                if i == s.cursor {
                        numStyle = style.AccentStyle
                }
                b.WriteString(numStyle.Render(fmt.Sprintf("%02d  ", i+1)))

                // Lead name.
                nameStyle := style.BodyStyle
                if i == s.cursor {
                        nameStyle = style.PrimaryStyle
                }
                b.WriteString(nameStyle.Render(lead.Name))
                b.WriteString("\n")

                // Category · City · Rating line.
                b.WriteString(style.Indent(1))
                b.WriteString(style.CaptionStyle.Render(fmt.Sprintf("%s · %s · ⭐ %.1f",
                        lead.Category, lead.City, lead.Rating)))
                b.WriteString("\n")

                // Contextual info based on phase.
                b.WriteString(style.Indent(1))
                b.WriteString(viewLeadContext(&lead))
                b.WriteString("\n")

                // Response text line (e.g. response: "iya kak, boleh lihat").
                if respText := viewLeadResponseText(&lead); respText != "" {
                        b.WriteString(style.Indent(1))
                        b.WriteString(respText)
                        b.WriteString("\n")
                }

                // Actionable item (e.g. "belum dikirim offer").
                if actionable := viewLeadActionable(&lead); actionable != "" {
                        b.WriteString(style.Indent(1))
                        b.WriteString(style.WarningStyle.Render("→ " + actionable))
                        b.WriteString("\n")
                }

                if i < len(s.filtered)-1 {
                        b.WriteString(style.Section(style.ItemGap))
                }
        }

        b.WriteString(style.Section(style.SectionGap))

        // Footer actions.
        b.WriteString(s.viewFooterFiltered())

        return b.String()
}

// viewFullDetail renders the lead_full_detail state.
func (s *LeadsDB) viewFullDetail(now time.Time) string {
        if s.currentLead == nil {
                return ""
        }
        lead := s.currentLead

        var b strings.Builder

        // Lead name.
        b.WriteString(style.HeadingStyle.Render(lead.Name))
        b.WriteString("\n")

        b.WriteString(style.Separator())
        b.WriteString("\n")

        // Business info — uses DRY renderLeadInfoBlock with photos enabled.
        b.WriteString(renderLeadInfoBlock(lead,
                withPhotos(),
        ))

        b.WriteString(style.Section(style.SectionGap))

        // Timeline.
        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyDataTimeline)))
        b.WriteString("\n")

        s.timeline.Tick(now)
        b.WriteString(s.timeline.View())

        b.WriteString(style.Separator())
        b.WriteString("\n")

        // Status line.
        phaseLabel := phaseBadgeText(lead.Phase)
        phaseStyle := phaseBadgeStyle(lead.Phase)
        b.WriteString(style.Indent(1))
        b.WriteString(fmt.Sprintf("%s: %s → %s",
                i18n.T(i18n.KeyDataStatus),
                phaseStyle.Render(phaseLabel),
                style.MutedStyle.Render(i18n.T(i18n.KeyDataWaitingOffer)),
        ))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // Action keys — doc spec: "1 kirim offer  2 balas custom  3 nanti"
        b.WriteString(style.ActionStyle.Render("1 " + i18n.T(i18n.KeyDataSendOffer)))
        b.WriteString("    ")
        b.WriteString(style.ActionStyle.Render("2 " + i18n.T(i18n.KeyDataCustomReply)))
        b.WriteString("    ")
        b.WriteString(style.ActionStyle.Render("3 " + i18n.T(i18n.KeyDataLater)))
        b.WriteString("\n")
        b.WriteString(style.ActionStyle.Render("4 " + i18n.T(i18n.KeyDataMarkConvert)))
        b.WriteString("    ")
        b.WriteString(style.ActionStyle.Render("5 " + i18n.T(i18n.KeyDataArchive)))
        b.WriteString("    ")
        b.WriteString(style.ActionStyle.Render("6 " + i18n.T(i18n.KeyDataBlock)))
        b.WriteString("\n")

        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLabelBack)))

        return b.String()
}

// viewFollowUpDue renders the lead_follow_up_due variant.
func (s *LeadsDB) viewFollowUpDue(now time.Time) string {
        if s.currentLead == nil {
                return ""
        }
        lead := s.currentLead

        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(lead.Name))
        b.WriteString("\n")

        // Shared lead info block: location, rating, website, score.
        b.WriteString(renderLeadInfoBlock(lead,
                withoutInsta(),
        ))

        b.WriteString(style.Section(style.SubSectionGap))

        // Follow-up specific info.
        b.WriteString(style.Indent(1))
        b.WriteString(style.CaptionStyle.Render(
                fmt.Sprintf("%s %s (%s)", i18n.T(i18n.KeyDataIceBreakerColon), lead.IceBreakerTime,
                        fmt.Sprintf(i18n.T(i18n.KeyDataContactCount), lead.ContactCount))))
        b.WriteString("\n")

        b.WriteString(style.Indent(1))
        b.WriteString(fmt.Sprintf("%s: %s", i18n.T(i18n.KeyDataResponse),
                style.CaptionStyle.Render(i18n.T(i18n.KeyDataNoResponseYet))))
        b.WriteString("\n")

        // Follow-up due date — use backend-provided text, or i18n fallback.
        dueText := lead.FollowUpDueText
        if dueText == "" {
                dueText = i18n.T(i18n.KeyDataDueToday)
        }
        b.WriteString(style.Indent(1))
        b.WriteString(fmt.Sprintf("%s: %s (%s)",
                i18n.T(i18n.KeyDataFollowUpNext),
                style.WarningStyle.Render(dueText),
                style.MutedStyle.Render(i18n.T(i18n.KeyDataAutoSchedule))))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // Actions — doc spec: "1 kirim follow-up  2 skip  3 tandai dingin  q balik"
        b.WriteString(style.ActionStyle.Render("1 " + i18n.T(i18n.KeyDataSendFollowUp)))
        b.WriteString("    ")
        b.WriteString(style.ActionStyle.Render("2 " + i18n.T(i18n.KeyGeneralSkip)))
        b.WriteString("    ")
        b.WriteString(style.ActionStyle.Render("3 " + i18n.T(i18n.KeyDataMarkCold)))
        b.WriteString("    ")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLabelBack)))

        return b.String()
}

// viewCold renders the leads_cold variant.
func (s *LeadsDB) viewCold(now time.Time) string {
        if s.currentLead == nil {
                return ""
        }
        lead := s.currentLead

        var b strings.Builder

        // Name with COLD badge — dynamic right-alignment (mirrors viewConverted).
        nameRendered := style.HeadingStyle.Render(lead.Name)
        badgeRendered := style.BadgeColdStyle.Render(i18n.T(i18n.KeyBadgeCold))
        coldGap := s.width - len(lead.Name) - len(i18n.T(i18n.KeyBadgeCold)) - 4
        if coldGap < 4 {
                coldGap = 4
        }
        b.WriteString(nameRendered)
        b.WriteString(strings.Repeat(" ", coldGap))
        b.WriteString(badgeRendered)
        b.WriteString("\n")

        // Shared lead info block: location, rating, website, score.
        b.WriteString(renderLeadInfoBlock(lead,
                withoutInsta(),
        ))

        b.WriteString(style.Section(style.SubSectionGap))

        // Cold-specific info — dynamic data from backend, not hardcoded.
        b.WriteString(style.Indent(1))
        b.WriteString(style.CaptionStyle.Render(
                fmt.Sprintf("%s %s", i18n.T(i18n.KeyDataIceBreakerColon), lead.IceBreakerTime)))
        b.WriteString("\n")

        // Follow-up timestamps — rendered from lead.FollowUpTimes.
        for i, fuTime := range lead.FollowUpTimes {
                b.WriteString(style.Indent(1))
                b.WriteString(style.CaptionStyle.Render(
                        fmt.Sprintf(i18n.T(i18n.KeyDataFollowUpNum), i+1, fuTime)))
                b.WriteString("\n")
        }

        b.WriteString(style.Indent(1))
        b.WriteString(fmt.Sprintf("%s: %s (%s, %s)",
                i18n.T(i18n.KeyDataResponse),
                style.CaptionStyle.Render(i18n.T(i18n.KeyDataNoResponseYet)),
                fmt.Sprintf(i18n.T(i18n.KeyDataFollowUpCount), lead.FollowUpCount),
                style.DimStyle.Render(i18n.T(i18n.KeyDataStillQuiet))))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // Cold description.
        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyDataColdDesc)))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyDataColdOption)))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(style.DimStyle.Render(i18n.T(i18n.KeyDataColdAuto)))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // Actions — doc spec: "1 follow-up terakhir (ke-3)  2 tandai dingin  3 archive  q balik"
        b.WriteString(style.ActionStyle.Render("1 " + i18n.T(i18n.KeyDataLastFollowUp)))
        b.WriteString("    ")
        b.WriteString(style.ActionStyle.Render("2 " + i18n.T(i18n.KeyDataMarkCold)))
        b.WriteString("    ")
        b.WriteString(style.ActionStyle.Render("3 " + i18n.T(i18n.KeyDataArchive)))
        b.WriteString("    ")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLabelBack)))

        return b.String()
}

// viewNeverContacted renders the leads_never_contacted variant.
func (s *LeadsDB) viewNeverContacted(now time.Time) string {
        if s.currentLead == nil {
                return ""
        }
        lead := s.currentLead

        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(lead.Name))
        b.WriteString("\n")

        // Shared lead info block: location, rating, website, score.
        b.WriteString(renderLeadInfoBlock(lead,
                withoutInsta(),
        ))

        b.WriteString(style.Section(style.SubSectionGap))

        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyDataNotContacted)))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.ActionStyle.Render("1 " + i18n.T(i18n.KeyDataSendIceBreaker)))
        b.WriteString("    ")
        b.WriteString(style.ActionStyle.Render("2 " + i18n.T(i18n.KeyGeneralSkip)))
        b.WriteString("    ")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLabelBack)))

        return b.String()
}

// viewConverted renders the leads_converted variant.
//
// NOTE: This view intentionally deviates from renderLeadInfoBlock because the
// converted view has a different layout order: niche appears on its own line
// between location and rating, which doesn't match the renderLeadInfoBlock
// pattern of location → rating+web/insta → score. Keeping manual rendering
// here is the simplest DRY approach — the layout divergence makes a shared
// helper more complex than the duplication it would save.
func (s *LeadsDB) viewConverted(now time.Time) string {
        if s.currentLead == nil {
                return ""
        }
        lead := s.currentLead

        var b strings.Builder

        // Name with DEAL badge — dynamic right-alignment.
        nameRendered := style.HeadingStyle.Render(lead.Name)
        badgeRendered := style.BadgeConvertedStyle.Render(i18n.T(i18n.KeyBadgeConverted))
        gap := s.width - len(lead.Name) - len(i18n.T(i18n.KeyBadgeConverted)) - 4
        if gap < 4 {
                gap = 4
        }
        b.WriteString(nameRendered)
        b.WriteString(strings.Repeat(" ", gap))
        b.WriteString(badgeRendered)
        b.WriteString("\n")

        b.WriteString(style.Separator())
        b.WriteString("\n")

        b.WriteString(style.Indent(1))
        b.WriteString(style.CaptionStyle.Render(
                fmt.Sprintf("%s · %s, %s", lead.Category, lead.Address, lead.City),
        ))
        b.WriteString("\n")

        b.WriteString(style.Indent(1))
        b.WriteString(fmt.Sprintf("%s: %s", i18n.T(i18n.KeyDataNiche), lead.Niche))
        b.WriteString("\n")

        b.WriteString(style.Indent(1))
        b.WriteString(fmt.Sprintf("⭐ %.1f · %d %s", lead.Rating, lead.Reviews, i18n.T(i18n.KeyDataReviewsLabel)))
        b.WriteString("\n")

        b.WriteString(style.Separator())
        b.WriteString("\n")

        // Timeline.
        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyDataTimeline)))
        b.WriteString("\n")

        s.timeline.Tick(now)
        b.WriteString(s.timeline.View())

        b.WriteString(style.Separator())
        b.WriteString("\n")

        // Conversion details.
        b.WriteString(style.Indent(1))
        b.WriteString(fmt.Sprintf("%s: %s", i18n.T(i18n.KeyDataDuration), lead.Duration))
        b.WriteString("\n")

        b.WriteString(style.Indent(1))
        b.WriteString(fmt.Sprintf("%s: %s", i18n.T(i18n.KeyDataTemplate), lead.TemplateName))
        b.WriteString("\n")

        b.WriteString(style.Indent(1))
        b.WriteString(fmt.Sprintf("%s: %s", i18n.T(i18n.KeyDataWorker), lead.WorkerName))
        b.WriteString("\n")

        b.WriteString(style.Separator())
        b.WriteString("\n")

        // Revenue with gold shimmer.
        b.WriteString(style.Indent(1))
        b.WriteString(style.GoldStyle.Render(i18n.T(i18n.KeyDataTrophy) + " "))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyDataConversionRevenue)))
        b.WriteString(" ")
        b.WriteString(style.GoldStyle.Render(lead.Revenue))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyLabelBack)))

        return b.String()
}

// viewFooterList renders the action footer for the leads_list state.
func (s *LeadsDB) viewFooterList() string {
        parts := []string{
                i18n.T(i18n.KeyDataMove),
                i18n.T(i18n.KeyDataViewDetail),
                i18n.T(i18n.KeyLabelSearch),
                i18n.T(i18n.KeyLabelBack),
        }
        return style.DimStyle.Render(strings.Join(parts, "    "))
}

// viewFooterFiltered renders the action footer for the leads_filtered state.
func (s *LeadsDB) viewFooterFiltered() string {
        parts := []string{
                i18n.T(i18n.KeyDataMove),
                i18n.T(i18n.KeyDataView),
                i18n.T(i18n.KeyLabelSearch),
                i18n.T(i18n.KeyDataSendOfferAll),
                i18n.T(i18n.KeyLabelBack),
        }
        return style.DimStyle.Render(strings.Join(parts, "    "))
}

// ---------------------------------------------------------------------------
// Shared rendering helpers — DRY consolidations
// ---------------------------------------------------------------------------

// Layout constants for consistent spacing.
const (
        categoryColWidth = 18  // minimum column width for category labels
        maxRecentLeads   = 3   // number of recent leads shown in list view

        // timelineTimeLayout is the expected time format from the backend for
        // timeline event times. Kept as a named constant so the backend data
        // format assumption is explicit rather than buried in a literal.
        timelineTimeLayout = "15:04"
)

// renderLeadLocation returns the formatted "category · address, city" line.
func renderLeadLocation(l *Lead) string {
        return style.CaptionStyle.Render(
                fmt.Sprintf("%s · %s, %s", l.Category, l.Address, l.City))
}

// renderLeadRating returns the formatted "⭐ rating · N reviews" line.
func renderLeadRating(l *Lead) string {
        return fmt.Sprintf("⭐ %.1f · %d %s", l.Rating, l.Reviews, i18n.T(i18n.KeyDataReviewsLabel))
}

// renderWebStatus returns the formatted website status indicator.
func renderWebStatus(hasWeb bool) string {
        if hasWeb {
                return style.MutedStyle.Render(i18n.T(i18n.KeyDataHasWebsite))
        }
        return style.CaptionStyle.Render(i18n.T(i18n.KeyDataNoWebsite))
}

// renderInstaStatus returns the formatted instagram status indicator.
func renderInstaStatus(hasInsta bool) string {
        if hasInsta {
                return style.MutedStyle.Render(i18n.T(i18n.KeyDataHasInstagram))
        }
        return style.CaptionStyle.Render(i18n.T(i18n.KeyDataNoInstagram))
}

// renderLeadScore returns the formatted "score: N/10" line.
func renderLeadScore(l *Lead) string {
        return fmt.Sprintf("%s: %s", i18n.T(i18n.KeyDataScore),
                style.PrimaryStyle.Render(fmt.Sprintf("%d/10", l.Score)))
}

// renderLeadInfoBlock renders the common lead info block shared across
// all detail variants: location, rating, website/instagram, score.
// This is the single source of truth for the lead detail header pattern.
func renderLeadInfoBlock(l *Lead, opts ...renderOption) string {
        var b strings.Builder
        cfg := defaultRenderConfig()
        for _, opt := range opts {
                opt(&cfg)
        }

        b.WriteString(style.Indent(1))
        b.WriteString(renderLeadLocation(l))
        b.WriteString("\n")

        b.WriteString(style.Indent(1))
        b.WriteString(renderLeadRating(l))
        if cfg.showWeb {
                b.WriteString(fmt.Sprintf(" · %s", renderWebStatus(l.HasWeb)))
        }
        if cfg.showInsta {
                b.WriteString(fmt.Sprintf(" · %s", renderInstaStatus(l.HasInsta)))
        }
        if cfg.showPhotos && l.PhotoCount > 0 {
                b.WriteString(style.CaptionStyle.Render(
                        fmt.Sprintf(" · "+i18n.T(i18n.KeyDataPhotosCount), l.PhotoCount)))
        }
        b.WriteString("\n")

        if cfg.showScore {
                if style.SubSectionGap > 0 {
                        b.WriteString(style.Section(style.SubSectionGap))
                }
                b.WriteString(style.Indent(1))
                b.WriteString(renderLeadScore(l))
                b.WriteString("\n")
        }

        return b.String()
}

// renderConfig controls which elements appear in the info block.
type renderConfig struct {
        showWeb    bool
        showInsta  bool
        showPhotos bool
        showScore  bool
}

func defaultRenderConfig() renderConfig {
        return renderConfig{
                showWeb:    true,
                showInsta:  true,
                showPhotos: false,
                showScore:  true,
        }
}

// renderOption is a functional option for renderLeadInfoBlock.
type renderOption func(*renderConfig)

// withoutScore disables the score line in the info block.
func withoutScore() renderOption {
        return func(c *renderConfig) { c.showScore = false }
}

// withoutInsta disables the instagram indicator in the info block.
func withoutInsta() renderOption {
        return func(c *renderConfig) { c.showInsta = false }
}

// withPhotos enables the photo count line in the info block.
func withPhotos() renderOption {
        return func(c *renderConfig) { c.showPhotos = true }
}

// ---------------------------------------------------------------------------
// Helper functions — DRY shared logic
// ---------------------------------------------------------------------------

// resolveDetailState maps a lead's phase to the correct detail state.
func resolveDetailState(lead *Lead) protocol.StateID {
        switch lead.Phase {
        case protocol.LeadFollowUp1, protocol.LeadFollowUp2, protocol.LeadNoResponse:
                return protocol.LeadsFollowUpDue
        case protocol.LeadCold:
                return protocol.LeadsCold
        case protocol.LeadBaru:
                return protocol.LeadsNeverContacted
        case protocol.LeadConverted:
                return protocol.LeadsConverted
        default:
                return protocol.LeadsFullDetail
        }
}

// phaseBadgeText returns the display label for a lead phase.
func phaseBadgeText(phase protocol.LeadPhase) string {
        switch phase {
        case protocol.LeadBaru:
                return i18n.T(i18n.KeyBadgeNew)
        case protocol.LeadCold:
                return i18n.T(i18n.KeyBadgeCold)
        case protocol.LeadConverted:
                return i18n.T(i18n.KeyBadgeConverted)
        case protocol.LeadResponded, protocol.LeadOfferSent:
                return i18n.T(i18n.KeyBadgeResponded)
        case protocol.LeadFailed, protocol.LeadDead, protocol.LeadBlocked:
                return i18n.T(i18n.KeyBadgeFailed)
        case protocol.LeadIceBreakerSent:
                return i18n.T(i18n.KeyDataIceBreakerLabel)
        case protocol.LeadFollowUp1:
                return i18n.T(i18n.KeyDataFollowUp1Label)
        case protocol.LeadFollowUp2:
                return i18n.T(i18n.KeyDataFollowUp2Label)
        case protocol.LeadNegative:
                return i18n.T(i18n.KeyDataNegativeLabel)
        case protocol.LeadArchived:
                return i18n.T(i18n.KeyDataArchivedLabel)
        case protocol.LeadAutoReply:
                return i18n.T(i18n.KeyDataAutoReplyLabel)
        default:
                return string(phase)
        }
}

// phaseBadgeStyle returns the lipgloss style for a lead phase badge.
// Badge colors from doc/06: baru=accent, respond=amber, convert=success, failed=dimmed.
func phaseBadgeStyle(phase protocol.LeadPhase) lipgloss.Style {
        switch phase {
        case protocol.LeadBaru:
                return style.BadgeNewStyle
        case protocol.LeadCold:
                return style.BadgeColdStyle
        case protocol.LeadConverted:
                return style.BadgeConvertedStyle
        case protocol.LeadResponded, protocol.LeadOfferSent:
                return style.BadgeRespondedStyle
        case protocol.LeadFailed, protocol.LeadDead, protocol.LeadBlocked:
                return style.BadgeFailedStyle
        default:
                return style.MutedStyle
        }
}

// phaseLabelStyle returns the label style for a filter category row.
func phaseLabelStyle(phase protocol.LeadPhase) lipgloss.Style {
        switch phase {
        case protocol.LeadBaru:
                return style.AccentStyle
        case protocol.LeadCold:
                return style.DimStyle
        case protocol.LeadConverted:
                return style.SuccessStyle
        case protocol.LeadResponded, protocol.LeadOfferSent:
                return style.WarningStyle
        default:
                return style.MutedStyle
        }
}

// phaseLabelDisplay returns a human-readable label for a filter phase.
func phaseLabelDisplay(phase protocol.LeadPhase) string {
        switch phase {
        case protocol.LeadBaru:
                return i18n.T(i18n.KeyBadgeNew)
        case protocol.LeadIceBreakerSent:
                return i18n.T(i18n.KeyDataIceBreakerLabel)
        case protocol.LeadFollowUp1:
                return i18n.T(i18n.KeyDataFollowUp1Label)
        case protocol.LeadFollowUp2:
                return i18n.T(i18n.KeyDataFollowUp2Label)
        case protocol.LeadCold:
                return i18n.T(i18n.KeyDataColdLabel)
        case protocol.LeadResponded:
                return i18n.T(i18n.KeyBadgeResponded)
        case protocol.LeadOfferSent:
                return i18n.T(i18n.KeyDataOfferSentLabel)
        case protocol.LeadConverted:
                return i18n.T(i18n.KeyDataConvertLabel)
        case protocol.LeadFailed:
                return i18n.T(i18n.KeyDataFailedLabel)
        case protocol.LeadBlocked:
                return i18n.T(i18n.KeyDataBlockedLabel)
        default:
                return string(phase)
        }
}

// viewLeadContext renders the contextual line for a lead in the filtered list.
// Doc spec shows: "ice breaker: kemarin 14:23" then "response: \"iya kak, boleh lihat\""
func viewLeadContext(lead *Lead) string {
        switch lead.Phase {
        case protocol.LeadIceBreakerSent:
                return style.CaptionStyle.Render(
                        fmt.Sprintf(i18n.T(i18n.KeyDataIceBreakerTime), lead.IceBreakerTime))
        case protocol.LeadResponded:
                return style.CaptionStyle.Render(
                        fmt.Sprintf(i18n.T(i18n.KeyDataIceBreakerTime), lead.IceBreakerTime))
        case protocol.LeadOfferSent:
                return style.CaptionStyle.Render(
                        fmt.Sprintf(i18n.T(i18n.KeyDataOfferTime), lead.ResponseTime))
        default:
                return style.CaptionStyle.Render(
                        fmt.Sprintf("%s: %s", phaseBadgeText(lead.Phase), lead.ResponseTime))
        }
}

// viewLeadResponseText renders the response text line for a lead in the filtered list.
// Doc spec: response: "iya kak, boleh lihat"
func viewLeadResponseText(lead *Lead) string {
        if lead.ResponseText == "" {
                return ""
        }
        return style.CaptionStyle.Render(
                fmt.Sprintf("%s: \"%s\"", i18n.T(i18n.KeyDataResponse), lead.ResponseText))
}

// viewLeadActionable renders the actionable warning text for a lead, or empty
// string if no actionable item applies. "belum dikirim offer" in amber.
func viewLeadActionable(lead *Lead) string {
        if lead.Phase == protocol.LeadResponded && lead.ResponseText != "" {
                return i18n.T(i18n.KeyDataOfferNotSent)
        }
        if lead.Phase == protocol.LeadOfferSent && lead.ResponseText != "" {
                return i18n.T(i18n.KeyDataNotReplied)
        }
        return ""
}

// buildTimelineForLead populates the Timeline component from a Lead's events.
func buildTimelineForLead(tl *component.Timeline, lead *Lead) {
        tl.Events = nil
        tl.AnimStart = time.Now()

        for i, ev := range lead.Timeline {
                tl.Add(component.TimelineEvent{
                        Time:         ev.Time,
                        Title:        ev.Action,
                        Detail:       ev.Detail,
                        Highlight:    isHighlightAction(ev.Action),
                        StaggerIndex: i,
                        Visible:      false,
                })
        }
}

// isHighlightAction determines if a timeline action should be highlighted.
func isHighlightAction(action string) bool {
        highlightActions := map[string]bool{
                "response_received": true,
                "mark_convert":      true,
                "deal":              true,
        }
        return highlightActions[action]
}

// parseFilterCategories converts raw params into FilterCategory slices.
func parseFilterCategories(raw []any) []FilterCategory {
        var cats []FilterCategory
        for _, item := range raw {
                if m, ok := item.(map[string]any); ok {
                        cat := FilterCategory{}
                        if label, ok := m[protocol.ParamLabel].(string); ok {
                                cat.Label = label
                        }
                        if count, ok := m[protocol.ParamCount].(float64); ok {
                                cat.Count = int(count)
                        }
                        if phase, ok := m[protocol.ParamPhase].(string); ok {
                                cat.Phase = protocol.LeadPhase(phase)
                        }
                        if note, ok := m[protocol.ParamNote].(string); ok {
                                cat.Note = note
                        }
                        cats = append(cats, cat)
                }
        }
        return cats
}

// parseLeads converts raw params into Lead slices.
func parseLeads(raw []any) []Lead {
        var leads []Lead
        for _, item := range raw {
                if m, ok := item.(map[string]any); ok {
                        leads = append(leads, parseLead(m))
                }
        }
        return leads
}

// parseLead converts a single map into a Lead.
func parseLead(m map[string]any) Lead {
        lead := Lead{}
        if v, ok := m[protocol.ParamName].(string); ok {
                lead.Name = v
        }
        if v, ok := m[protocol.ParamCategory].(string); ok {
                lead.Category = v
        }
        if v, ok := m[protocol.ParamAddress].(string); ok {
                lead.Address = v
        }
        if v, ok := m[protocol.ParamCity].(string); ok {
                lead.City = v
        }
        if v, ok := m[protocol.ParamRating].(float64); ok {
                lead.Rating = v
        }
        if v, ok := m[protocol.ParamReviews].(float64); ok {
                lead.Reviews = int(v)
        }
        if v, ok := m[protocol.ParamHasWeb].(bool); ok {
                lead.HasWeb = v
        }
        if v, ok := m[protocol.ParamHasInsta].(bool); ok {
                lead.HasInsta = v
        }
        if v, ok := m[protocol.ParamPhotoCount].(float64); ok {
                lead.PhotoCount = int(v)
        }
        if v, ok := m[protocol.ParamScore].(float64); ok {
                lead.Score = int(v)
        }
        if v, ok := m[protocol.ParamPhase].(string); ok {
                lead.Phase = protocol.LeadPhase(v)
        }
        if v, ok := m[protocol.ParamNiche].(string); ok {
                lead.Niche = v
        }
        if v, ok := m[protocol.ParamResponseText].(string); ok {
                lead.ResponseText = v
        }
        if v, ok := m[protocol.ParamIceBreakerTime].(string); ok {
                lead.IceBreakerTime = v
        }
        if v, ok := m[protocol.ParamResponseTime].(string); ok {
                lead.ResponseTime = v
        }
        if v, ok := m[protocol.ParamContactCount].(float64); ok {
                lead.ContactCount = int(v)
        }
        if v, ok := m[protocol.ParamFollowupCount].(float64); ok {
                lead.FollowUpCount = int(v)
        }
        if times, ok := m[protocol.ParamFollowupTimes].([]any); ok {
                for _, t := range times {
                        if s, ok := t.(string); ok {
                                lead.FollowUpTimes = append(lead.FollowUpTimes, s)
                        }
                }
        }
        if v, ok := m[protocol.ParamFollowupDueText].(string); ok {
                lead.FollowUpDueText = v
        }
        if v, ok := m[protocol.ParamDuration].(string); ok {
                lead.Duration = v
        }
        if v, ok := m[protocol.ParamTemplateName].(string); ok {
                lead.TemplateName = v
        }
        if v, ok := m[protocol.ParamWorkerName].(string); ok {
                lead.WorkerName = v
        }
        if v, ok := m[protocol.ParamRevenue].(string); ok {
                lead.Revenue = v
        }

        // Parse timeline events.
        if events, ok := m[protocol.ParamTimeline].([]any); ok {
                for _, ev := range events {
                        if em, ok := ev.(map[string]any); ok {
                                event := LeadEvent{}
                                if t, ok := em[protocol.ParamTime].(string); ok {
                                        if parsed, err := time.Parse(timelineTimeLayout, t); err == nil {
                                                event.Time = parsed
                                        }
                                }
                                if a, ok := em[protocol.ParamAction].(string); ok {
                                        event.Action = a
                                }
                                if d, ok := em[protocol.ParamDetail].(string); ok {
                                        event.Detail = d
                                }
                                lead.Timeline = append(lead.Timeline, event)
                        }
                }
        }

        return lead
}

