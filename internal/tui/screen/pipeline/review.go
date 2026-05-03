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
// typeOutTickMsg drives the char-by-char template preview animation.
// ---------------------------------------------------------------------------

// typeOutTickMsg is an internal tea.Msg that ticks the template type-out
// animation forward at ~60fps.
type typeOutTickMsg struct {
        Time time.Time
}

// typeOutTickCmd returns a tea.Cmd that sends a typeOutTickMsg after a short
// delay to advance the template type-out animation.
func typeOutTickCmd() tea.Cmd {
        return tea.Tick(16*time.Millisecond, func(t time.Time) tea.Msg {
                return typeOutTickMsg{Time: t}
        })
}

// ---------------------------------------------------------------------------
// Review — Screen 5: LEAD REVIEW (optional manual override)
// ---------------------------------------------------------------------------

// Review implements the lead review screen with four states:
//   - ReviewReviewing: one-lead-at-a-time approval/skip flow
//   - ReviewLeadDetail: full lead detail view
//   - ReviewTemplatePreview: template char-by-char preview
//   - ReviewQueueComplete: all leads reviewed summary
//
// The screen is optional — WaClaw auto-reviews by default. It only appears
// when the user explicitly wants to inspect leads before sending.
type Review struct {
        tui.ScreenBase
        state  protocol.StateID
        width  int
        height int

        // Queue data.
        leads []LeadItem // The full review queue.
        cursor int       // Currently selected lead index.

        // Session counters.
        stats ReviewStats

        // Template preview state.
        previewVar  int    // Currently previewed variant index (1-based).
        templates   []string
        preview     component.TemplatePreview
        previewText string // The selected template text.

        // Animation.
        typeOutDone bool

        // prevState tracks the state before entering template_preview,
        // so that canceling preview returns to the correct state.
        prevState protocol.StateID
}

// NewReview creates a Review screen in the default ReviewReviewing state.
func NewReview() Review {
        return Review{
                ScreenBase: tui.NewScreenBase(protocol.ScreenLeadReview),
                state:      protocol.ReviewReviewing,
        }
}

// ---------------------------------------------------------------------------
// tea.Model interface
// ---------------------------------------------------------------------------

// Init returns the initial command. No async work needed for review.
func (r Review) Init() tea.Cmd { return nil }

// Update handles all tea.Msg values: key presses, window size changes,
// bus messages, and internal animation ticks.
func (r Review) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch m := msg.(type) {
        case tea.KeyMsg:
                return r.handleKey(m)
        case tea.WindowSizeMsg:
                r.width = m.Width
                r.height = m.Height
                r.preview.Width = r.width - (4 * style.IndentPerLevel)
                return r, nil
        case typeOutTickMsg:
                return r.handleTypeOutTick(m)
        case bus.UpdateMsg:
                if m.Screen == protocol.ScreenLeadReview {
                        _ = r.HandleUpdate(m.Params)
                }
                return r, nil
        case bus.NavigateMsg:
                if m.Screen == protocol.ScreenLeadReview {
                        _ = r.HandleNavigate(m.Params)
                }
                return r, nil
        }
        return r, nil
}

// View renders the current state of the review screen.
func (r Review) View() string {
        switch r.state {
        case protocol.ReviewReviewing:
                return r.viewReviewing()
        case protocol.ReviewLeadDetail:
                return r.viewLeadDetail()
        case protocol.ReviewTemplatePreview:
                return r.viewTemplatePreview()
        case protocol.ReviewQueueComplete:
                return r.viewQueueComplete()
        default:
                return r.viewReviewing()
        }
}

// ---------------------------------------------------------------------------
// Screen interface — HandleNavigate
// ---------------------------------------------------------------------------

// HandleNavigate processes backend navigation commands. Supported params:
//   - "state" (string): switch to the given StateID
//   - "leads" ([]LeadItem): initial lead data
//   - "templates" ([]string): template variant names
func (r *Review) HandleNavigate(params map[string]any) error {
        if s, ok := params[protocol.ParamState].(string); ok {
                r.state = protocol.StateID(s)
        }
        if leads, ok := params[protocol.ParamLeads].([]LeadItem); ok {
                r.leads = leads
                r.clampCursor()
        }
        if templates, ok := params[protocol.ParamTemplates].([]string); ok {
                r.templates = templates
        }
        return nil
}

// ---------------------------------------------------------------------------
// Screen interface — HandleUpdate
// ---------------------------------------------------------------------------

// HandleUpdate processes backend data updates. Supported params:
//   - "leads" ([]LeadItem): replace the review queue
//   - "stats" (ReviewStats): update session counters
//   - "templates" ([]string): template variant names
//   - "current" (int): set cursor position
//   - "template_text" (string): set template text for preview
//   - "lead" (LeadItem): update current lead detail
func (r *Review) HandleUpdate(params map[string]any) error {
        if leads, ok := params[protocol.ParamLeads].([]LeadItem); ok {
                r.leads = leads
                r.clampCursor()
        }
        if stats, ok := params[protocol.ParamStats].(ReviewStats); ok {
                r.stats = stats
        }
        if templates, ok := params[protocol.ParamTemplates].([]string); ok {
                r.templates = templates
        }
        if current, ok := params[protocol.ParamCurrent].(int); ok {
                r.cursor = current
                r.clampCursor()
        }
        if text, ok := params[protocol.ParamTemplateText].(string); ok {
                r.previewText = text
        }
        if lead, ok := params[protocol.ParamLead].(LeadItem); ok {
                if lead.Index >= 0 && int(lead.Index) < len(r.leads) {
                        r.leads[lead.Index] = lead
                }
        }
        return nil
}

// ---------------------------------------------------------------------------
// Focus / Blur
// ---------------------------------------------------------------------------

// Focus is called when this screen becomes the active screen.
func (r *Review) Focus() {}

// Blur is called when this screen is no longer the active screen.
func (r *Review) Blur() {}

// ---------------------------------------------------------------------------
// Key handling
// ---------------------------------------------------------------------------

// handleKey dispatches key events based on the current review state.
func (r Review) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // Navigation keys: cursor movement.
        if key.Matches(msg, tui.KeyUp) {
                if r.cursor > 0 {
                        r.cursor--
                        r.publishKeyPress("up")
                }
                return r, nil
        }
        if key.Matches(msg, tui.KeyDown) {
                if r.cursor < len(r.leads)-1 {
                        r.cursor++
                        r.publishKeyPress("down")
                }
                return r, nil
        }

        // State-specific key dispatch.
        switch r.state {
        case protocol.ReviewReviewing:
                return r.handleKeyReviewing(msg)
        case protocol.ReviewLeadDetail:
                return r.handleKeyLeadDetail(msg)
        case protocol.ReviewTemplatePreview:
                return r.handleKeyTemplatePreview(msg)
        case protocol.ReviewQueueComplete:
                return r.handleKeyQueueComplete(msg)
        }
        return r, nil
}

// handleKeyReviewing handles keys in the main reviewing state.
func (r Review) handleKeyReviewing(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // Enter: approve lead.
        if key.Matches(msg, tui.KeyEnter) {
                r.publishAction(protocol.ActionReviewApprove, nil)
                return r, nil
        }
        // s: skip.
        if key.Matches(msg, tui.KeySkip) {
                r.publishAction(protocol.ActionReviewSkip, nil)
                return r, nil
        }
        // x: skip & block.
        if msg.String() == "x" {
                r.publishAction(protocol.ActionReviewBlock, nil)
                return r, nil
        }
        // d: show detail.
        if msg.String() == "d" {
                r.state = protocol.ReviewLeadDetail
                return r, nil
        }
        // 1/2/3: show template variant preview.
        if v := variantFromKey(msg.String()); v > 0 {
                return r.switchToTemplatePreview(v)
        }
        // q: exit review.
        if key.Matches(msg, tui.KeyBack) {
                r.publishKeyPress("q")
                return r, nil
        }
        return r, nil
}

// handleKeyLeadDetail handles keys in the lead detail state.
func (r Review) handleKeyLeadDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // Enter: approve from detail.
        if key.Matches(msg, tui.KeyEnter) {
                r.publishAction(protocol.ActionReviewApprove, nil)
                return r, nil
        }
        // s: skip from detail.
        if key.Matches(msg, tui.KeySkip) {
                r.publishAction(protocol.ActionReviewSkip, nil)
                return r, nil
        }
        // 1/2/3: show template variant preview.
        if v := variantFromKey(msg.String()); v > 0 {
                return r.switchToTemplatePreview(v)
        }
        // q or esc: back to reviewing.
        if key.Matches(msg, tui.KeyBack) || key.Matches(msg, tui.KeyEscape) {
                r.state = protocol.ReviewReviewing
                return r, nil
        }
        return r, nil
}

// handleKeyTemplatePreview handles keys in the template preview state.
func (r Review) handleKeyTemplatePreview(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // Enter: use this template (approve with variant).
        if key.Matches(msg, tui.KeyEnter) {
                params := map[string]any{
                        protocol.ParamVariant: r.previewVar,
                }
                r.publishAction(protocol.ActionReviewApprove, params)
                return r, nil
        }
        // 1/2/3: switch variant.
        if v := variantFromKey(msg.String()); v > 0 && v != r.previewVar {
                return r.switchToTemplatePreview(v)
        }
        // q or esc: cancel preview, return to previous state.
        if key.Matches(msg, tui.KeyBack) || key.Matches(msg, tui.KeyEscape) {
                r.state = r.prevState
                r.typeOutDone = true
                return r, nil
        }
        return r, nil
}

// handleKeyQueueComplete handles keys when the queue is empty.
func (r Review) handleKeyQueueComplete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // Enter: go to dashboard.
        if key.Matches(msg, tui.KeyEnter) {
                r.publishKeyPress("enter")
                return r, nil
        }
        // q: exit.
        if key.Matches(msg, tui.KeyBack) {
                r.publishKeyPress("q")
                return r, nil
        }
        return r, nil
}

// ---------------------------------------------------------------------------
// State transitions
// ---------------------------------------------------------------------------

// switchToTemplatePreview transitions to the template preview state for the
// given variant (1-based). It sets up the component.TemplatePreview with
// variables substituted from the current lead.
func (r Review) switchToTemplatePreview(variant int) (tea.Model, tea.Cmd) {
        r.prevState = r.state
        r.state = protocol.ReviewTemplatePreview
        r.previewVar = variant
        r.typeOutDone = false

        tp := component.NewTemplatePreview(r.previewText)
        tp.Width = r.width - (4 * style.IndentPerLevel)

        lead := r.currentLead()
        tp.SetVar("Title", lead.Name)
        tp.SetVar("Category", lead.Category)
        tp.SetVar("Address", lead.Address)
        tp.SetVar("City", lead.City)
        tp.SetVar("Rating", fmt.Sprintf("%.1f", lead.Rating))
        tp.SetVar("Reviews", fmt.Sprintf("%d", lead.Reviews))
        tp.SetVar("Area", "")

        tp.StartTypeOut()
        r.preview = tp

        return r, typeOutTickCmd()
}

// ---------------------------------------------------------------------------
// Animation
// ---------------------------------------------------------------------------

// handleTypeOutTick advances the template type-out animation one frame.
func (r Review) handleTypeOutTick(msg typeOutTickMsg) (tea.Model, tea.Cmd) {
        if r.state != protocol.ReviewTemplatePreview || r.typeOutDone {
                return r, nil
        }

        r.preview.Tick(msg.Time)

        if !r.preview.TypeOutEnabled {
                // Animation finished.
                r.typeOutDone = true
                return r, nil
        }

        return r, typeOutTickCmd()
}

// ---------------------------------------------------------------------------
// View — ReviewReviewing (main review state)
// ---------------------------------------------------------------------------

// viewReviewing renders the main lead-by-lead review view.
func (r Review) viewReviewing() string {
        var b strings.Builder

        lead := r.currentLead()
        if lead.Name == "" {
                return style.CaptionStyle.Render(i18n.T(i18n.KeyReviewWaiting))
        }

        // Title row: "review leads" left, "{n} nunggu" right.
        titleLeft := style.HeadingStyle.Render(i18n.T(i18n.KeyReviewTitle))
        titleRight := style.MutedStyle.Render(fmt.Sprintf("%d %s", len(r.leads), i18n.T(i18n.KeyReviewWaiting)))
        title := lipgloss.NewStyle().Width(r.width).Render(
                lipgloss.JoinHorizontal(lipgloss.Top, titleLeft, titleRight),
        )
        b.WriteString(title)
        b.WriteString(style.Section(style.SectionGap))

        // Lead number (01-indexed, zero-padded).
        b.WriteString(style.CaptionStyle.Render(fmt.Sprintf("%02d", lead.Index+1)))
        b.WriteString(style.Section(style.ItemGap))

        // Lead name.
        b.WriteString(style.PrimaryStyle.Render(lead.Name))
        b.WriteString(style.Section(style.ItemGap))

        // Category · Address.
        b.WriteString(style.BodyStyle.Render(fmt.Sprintf("%s · %s", lead.Category, lead.Address)))
        b.WriteString(style.Section(style.ItemGap))

        // Rating · reviews · website status.
        ratingLine := fmt.Sprintf("⭐ %.1f (%d %s)", lead.Rating, lead.Reviews, i18n.T(i18n.KeyWordReviews))
        if !lead.HasWebsite {
                ratingLine += " · " + i18n.T(i18n.KeyReviewNoWeb)
        }
        b.WriteString(style.BodyStyle.Render(ratingLine))
        b.WriteString(style.Section(style.ItemGap))

        // Instagram status.
        if !lead.HasInstagram {
                b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyReviewNoIG)))
                b.WriteString(style.Section(style.ItemGap))
        }

        // WA registration status.
        if lead.HasWA {
                waLine := fmt.Sprintf("WA: %s %s", "✓", i18n.T(i18n.KeyReviewWAReg))
                b.WriteString(style.SuccessStyle.Render(waLine))
                b.WriteString(style.Section(style.ItemGap))
        }

        // Blank line before prompt.
        b.WriteString(style.Section(style.SubSectionGap))

        // Prompt.
        promptStyle := lipgloss.NewStyle().Foreground(style.TextDim).Bold(true)
        b.WriteString(promptStyle.Render(i18n.T(i18n.KeyReviewQA)))
        b.WriteString(style.Section(style.SubSectionGap))

        // Action keys (two-column layout).
        indent := style.Indent(1)
        col1 := style.ActionStyle.Render("↵  " + i18n.T(i18n.KeyReviewApproveAction))
        col2 := style.MutedStyle.Render("1-3  " + i18n.T(i18n.KeyReviewPickTemplate))
        b.WriteString(lipgloss.NewStyle().Width(r.width).Render(
                lipgloss.JoinHorizontal(lipgloss.Top, indent+col1, col2),
        ))
        b.WriteString(style.Section(style.ItemGap))

        col3 := style.MutedStyle.Render("s  " + i18n.T(i18n.KeyReviewSkip))
        col4 := style.MutedStyle.Render("d  " + i18n.T(i18n.KeyReviewDetail))
        b.WriteString(lipgloss.NewStyle().Width(r.width).Render(
                lipgloss.JoinHorizontal(lipgloss.Top, indent+col3, col4),
        ))
        b.WriteString(style.Section(style.ItemGap))

        col5 := style.MutedStyle.Render("x  " + i18n.T(i18n.KeyReviewSkipBlock))
        b.WriteString(indent + col5)

        // Bottom blank line.
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.Section(style.SubSectionGap))

        // Stats bar.
        statsLine := fmt.Sprintf(
                "%s: %d  %s: %d  %s: %d",
                i18n.T(i18n.KeyReviewQueued), r.stats.Queued,
                i18n.T(i18n.KeyReviewSkipped), r.stats.Skipped,
                i18n.T(i18n.KeyReviewRemaining), r.stats.Remaining,
        )
        b.WriteString(style.MutedStyle.Render(statsLine))
        b.WriteString(style.Section(style.SubSectionGap))

        // Navigation hints.
        hints := fmt.Sprintf(
                "↑↓  %s   ↵  %s   q  %s",
                i18n.T(i18n.KeyReviewMove),
                i18n.T(i18n.KeyReviewGasQueue),
                i18n.T(i18n.KeyReviewDone),
        )
        b.WriteString(style.CaptionStyle.Render(hints))

        return b.String()
}

// ---------------------------------------------------------------------------
// View — ReviewLeadDetail (full detail state)
// ---------------------------------------------------------------------------

// viewLeadDetail renders the full detail view for the currently selected lead.
func (r Review) viewLeadDetail() string {
        var b strings.Builder
        lead := r.currentLead()

        // Business name heading.
        b.WriteString(style.HeadingStyle.Render(lead.Name))
        b.WriteString(style.Section(style.ItemGap))

        // Category · Address.
        b.WriteString(style.BodyStyle.Render(fmt.Sprintf("%s · %s", lead.Category, lead.Address)))
        b.WriteString(style.Section(style.SubSectionGap))

        // Rating · reviews count.
        b.WriteString(fmt.Sprintf("⭐ %.1f · %d %s", lead.Rating, lead.Reviews, i18n.T(i18n.KeyWordReviews)))
        b.WriteString(style.Section(style.ItemGap))

        // Website / Instagram status.
        var meta []string
        if !lead.HasWebsite {
                meta = append(meta, i18n.T(i18n.KeyReviewNoWeb))
        }
        if !lead.HasInstagram {
                meta = append(meta, i18n.T(i18n.KeyReviewNoIG))
        }
        if len(meta) > 0 {
                b.WriteString(style.MutedStyle.Render(strings.Join(meta, " · ")))
                b.WriteString(style.Section(style.ItemGap))
        }

        // Photo count.
        if lead.PhotoCount > 0 {
                b.WriteString(style.BodyStyle.Render(fmt.Sprintf("%s: %d", i18n.T(i18n.KeyReviewPhotos), lead.PhotoCount)))
                b.WriteString(style.Section(style.ItemGap))
        }

        // Google search results.
        if len(lead.GoogleResults) > 0 {
                b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyReviewGoogleSrch) + ":"))
                b.WriteString(style.Section(style.ItemGap))
                indent := style.Indent(1)
                for _, result := range lead.GoogleResults {
                        b.WriteString(style.MutedStyle.Render(indent + "- " + result))
                        b.WriteString(style.Section(style.ItemGap))
                }
        }

        // Lead score.
        scoreLine := fmt.Sprintf("%s: %d/10", i18n.T(i18n.KeyReviewScore), lead.Score)
        if lead.Potential != "" {
                scoreLine += fmt.Sprintf(" (%s)", lead.Potential)
        }
        scoreColor := style.Success
        if lead.Score < 5 {
                scoreColor = style.Warning
        } else if lead.Score < 7 {
                scoreColor = style.Text
        }
        b.WriteString(lipgloss.NewStyle().Foreground(scoreColor).Render(scoreLine))
        b.WriteString(style.Section(style.SubSectionGap))

        // Contact history.
        if lead.ContactHistory != "" {
                historyLine := fmt.Sprintf("%s: %s", i18n.T(i18n.KeyReviewHistory), lead.ContactHistory)
                b.WriteString(style.BodyStyle.Render(historyLine))
                b.WriteString(style.Section(style.ItemGap))
        }

        // Follow-up status.
        if lead.FollowUpStatus != "" {
                b.WriteString(style.MutedStyle.Render(fmt.Sprintf("%s: %s", i18n.T(i18n.KeyReviewFollowUp), lead.FollowUpStatus)))
                b.WriteString(style.Section(style.ItemGap))
        }

        // Blank line before actions.
        b.WriteString(style.Section(style.SectionGap))

        // Actions.
        indent := style.Indent(0)
        act1 := style.ActionStyle.Render("↵  " + i18n.T(i18n.KeyReviewApproveShort))
        act2 := style.MutedStyle.Render("s  " + i18n.T(i18n.KeyReviewSkip))
        act3 := style.MutedStyle.Render("1-3  " + i18n.T(i18n.KeyReviewPickTemplate))
        act4 := style.MutedStyle.Render("q  " + i18n.T(i18n.KeyReviewBack))

        line1 := lipgloss.NewStyle().Width(r.width).Render(
                lipgloss.JoinHorizontal(lipgloss.Top, indent+act1, act2),
        )
        line2 := lipgloss.NewStyle().Width(r.width).Render(
                lipgloss.JoinHorizontal(lipgloss.Top, indent+act3, act4),
        )
        b.WriteString(line1)
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(line2)

        return b.String()
}

// ---------------------------------------------------------------------------
// View — ReviewTemplatePreview (template preview state)
// ---------------------------------------------------------------------------

// viewTemplatePreview renders the template preview with char-by-char animation.
func (r Review) viewTemplatePreview() string {
        var b strings.Builder

        // Header: "preview varian: {variant_name}".
        variantName := r.variantName(r.previewVar)
        header := fmt.Sprintf("%s: %s", i18n.T(i18n.KeyReviewPreviewVar), variantName)
        b.WriteString(style.SubHeadingStyle.Render(header))

        // Separator.
        b.WriteString(style.Section(style.SubSectionGap))

        // Template content via component.TemplatePreview.
        templateView := r.preview.View()
        b.WriteString(templateView)

        // Separator.
        b.WriteString(style.Section(style.SubSectionGap))

        // Other variant hints (show the other 2 variants).
        otherHints := r.renderOtherVariantHints()
        if otherHints != "" {
                b.WriteString(style.MutedStyle.Render(otherHints))
                b.WriteString(style.Section(style.ItemGap))
        }

        // Actions.
        act1 := style.ActionStyle.Render("↵  " + i18n.T(i18n.KeyReviewUseThis))
        act2 := style.MutedStyle.Render("1-3  " + i18n.T(i18n.KeyReviewChangeVar))
        act3 := style.MutedStyle.Render("q  " + i18n.T(i18n.KeyReviewCancel))

        line1 := lipgloss.NewStyle().Width(r.width).Render(
                lipgloss.JoinHorizontal(lipgloss.Top, act1, act2),
        )
        b.WriteString(line1)
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(act3)

        return b.String()
}

// ---------------------------------------------------------------------------
// View — ReviewQueueComplete (queue empty state)
// ---------------------------------------------------------------------------

// viewQueueComplete renders the summary shown after all leads are reviewed.
func (r Review) viewQueueComplete() string {
        var b strings.Builder

        // Heading.
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyReviewComplete)))
        b.WriteString(style.Section(style.SectionGap))

        // Stats.
        indent := style.Indent(0)
        b.WriteString(indent + style.BodyStyle.Render(fmt.Sprintf("%s: %d", i18n.T(i18n.KeyReviewQueued), r.stats.Queued)))
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(indent + style.BodyStyle.Render(fmt.Sprintf("%s: %d", i18n.T(i18n.KeyReviewSkipped), r.stats.Skipped)))
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(indent + style.BodyStyle.Render(fmt.Sprintf("%s: %d", i18n.T(i18n.KeyReviewBlocked), r.stats.Blocked)))

        // Auto-pilot messages.
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyReviewAutoSendTiming)))
        b.WriteString(style.Section(style.ItemGap))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyReviewAutoPilotRelax)))

        // Blank line before actions.
        b.WriteString(style.Section(style.SectionGap))

        // Actions.
        act1 := style.ActionStyle.Render("↵  " + i18n.T(i18n.KeyReviewDashLink))
        act2 := style.MutedStyle.Render("q  " + i18n.T(i18n.KeyReviewExit))

        line := lipgloss.NewStyle().Width(r.width).Render(
                lipgloss.JoinHorizontal(lipgloss.Top, act1, act2),
        )
        b.WriteString(line)

        return b.String()
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// currentLead returns the lead at the current cursor position, or an empty
// LeadItem if the cursor is out of bounds.
func (r Review) currentLead() LeadItem {
        if r.cursor >= 0 && r.cursor < len(r.leads) {
                return r.leads[r.cursor]
        }
        return LeadItem{}
}

// clampCursor ensures the cursor stays within the valid range [0, len(leads)-1].
func (r *Review) clampCursor() {
        if r.cursor >= len(r.leads) {
                r.cursor = len(r.leads) - 1
        }
        if r.cursor < 0 {
                r.cursor = 0
        }
}

// publishKeyPress sends a key_press event to the backend via the bus.
func (r Review) publishKeyPress(key string) {
        if r.Bus() != nil {
                r.Bus().Publish(bus.KeyPressMsg{
                        Key:    key,
                        Screen: protocol.ScreenLeadReview,
                })
        }
}

// publishAction sends an action event to the backend via the bus.
func (r Review) publishAction(action string, params map[string]any) {
        if r.Bus() != nil {
                r.Bus().Publish(bus.ActionMsg{
                        Action: action,
                        Screen: protocol.ScreenLeadReview,
                        Params: params,
                })
        }
}

// variantFromKey extracts a 1-based variant index from a key press string.
// Returns 0 if the key is not a valid variant selector (1, 2, or 3).
func variantFromKey(k string) int {
        if k == "1" {
                return 1
        }
        if k == "2" {
                return 2
        }
        if k == "3" {
                return 3
        }
        return 0
}

// variantName returns the display name for the given variant index.
// Returns an empty string if no template names are configured.
func (r Review) variantName(v int) string {
        if v >= 1 && v <= len(r.templates) {
                return r.templates[v-1]
        }
        return ""
}

// renderOtherVariantHints builds a string showing the other two variant
// selectors, e.g. "varian lain:  2 variant_2    3 variant_3".
func (r Review) renderOtherVariantHints() string {
        var others []string
        for i := 1; i <= 3; i++ {
                if i == r.previewVar {
                        continue
                }
                name := r.variantName(i)
                others = append(others, fmt.Sprintf("%d  %s", i, name))
        }
        if len(others) == 0 {
                return ""
        }
        prefix := i18n.T(i18n.KeyReviewOtherVar) + ":  "
        return prefix + strings.Join(others, "    ")
}
