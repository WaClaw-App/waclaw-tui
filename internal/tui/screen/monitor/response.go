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
// Response data models
// ---------------------------------------------------------------------------

// LeadResponse holds data about a single lead response displayed on the
// Response screen.
type LeadResponse struct {
        Business string // e.g. "kopi nusantara"
        Category string // e.g. "cafe"
        Area     string // e.g. "kediri"
        Message   string // e.g. "iya kak, boleh lihat desainnya?"
        OfferText string // rendered offer message (from backend, not the lead's response message)
        Trigger   string // closing trigger match text (e.g. "dah transfer")
        Niche    string // e.g. "[web_dev]"
        Class    string // auto-classification: "positif", "penasaran", "auto-reply"
}

// MultiQueueItem represents a single response in the multi-queue view.
type MultiQueueItem struct {
        Index    int    // 1-based display index
        Business string
        Category string
        Area     string
        Message  string
        Class    string // "positif", "penasaran", "auto-reply"
}

// ConversionData holds data for the conversion drama screen.
type ConversionData struct {
        Business    string
        Pipeline    string // e.g. "ice breaker → offer → deal"
        TimeTaken   string // e.g. "2 hari 4 jam"
        TrophyCount int    // conversions this week
        Revenue     string // e.g. "rp 7.5jt"
}

// ---------------------------------------------------------------------------
// Response screen — Screen 8: Response → Reward
// ---------------------------------------------------------------------------

// Response is the screen that handles incoming lead responses, closing
// trigger detection, offer preview, multi-response triage, and the
// all-important conversion drama sequence.
//
// Visual spec from doc/05-screens-monitor-response.md.
type Response struct {
        base   protocol.ScreenID
        bus    *bus.Bus
        state  protocol.StateID
        width  int
        height int

        // Lead data
        lead     LeadResponse
        conversion ConversionData
        queue    []MultiQueueItem
        cursor   int // for multi-queue navigation

        // Conversion drama animation state
        dramaPhase  dramaPhase
        dramaStart  time.Time
        keyAccepted bool

        // Particle system for conversion/jackpot effects
        particles component.ParticleSystem
}

// dramaPhase tracks the 4-phase conversion drama animation.
type dramaPhase int

const (
        dramaNone    dramaPhase = iota
        dramaShock              // Phase 1: 0-200ms — full white flash + bell
        dramaReveal             // Phase 2: 200-800ms — DEAL text scales in + particles
        dramaContext            // Phase 3: 800-1500ms — business name + trophy
        dramaSettle             // Phase 4: 1500ms+ — settle + accept keyboard
)

// Backend classification constants for multi-queue badge rendering.
// These are aliases of the protocol package constants for convenience.
const (
        classPositif   = string(protocol.ClassPositif)
        classCurious   = string(protocol.ClassCurious)
        classAutoReply = string(protocol.ClassAutoReply)
)

// NewResponse creates a new Response screen.
func NewResponse() *Response {
        return &Response{
                base:  protocol.ScreenResponse,
                state: protocol.ResponsePositive,
        }
}

// ID returns the screen identifier.
func (r *Response) ID() protocol.ScreenID { return r.base }

// SetBus injects the event bus reference.
func (r *Response) SetBus(b *bus.Bus) { r.bus = b }

// HandleNavigate processes a "navigate" command from the backend.
func (r *Response) HandleNavigate(params map[string]any) error {
        if stateStr, ok := params[protocol.ParamState].(string); ok {
                newState := protocol.StateID(stateStr)
                r.state = newState
                // Start conversion drama if entering conversion state
                if newState == protocol.ResponseConversion {
                        r.startConversionDrama()
                }
        }
        r.populateData(params)
        return nil
}

// HandleUpdate processes an "update" command from the backend.
func (r *Response) HandleUpdate(params map[string]any) error {
        r.populateData(params)
        return nil
}

// Focus is called when this screen becomes active.
func (r *Response) Focus() {}

// Blur is called when this screen is no longer active.
func (r *Response) Blur() {}

// Init implements tea.Model.
func (r *Response) Init() tea.Cmd {
        return tickResponseCmd()
}

// Update implements tea.Model.
func (r *Response) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        now := time.Now()

        switch m := msg.(type) {
        case tea.WindowSizeMsg:
                r.width = m.Width
                r.height = m.Height
                r.particles = component.NewParticleSystem(r.width, r.height)
                return r, nil

        case tea.KeyMsg:
                // During conversion drama, don't accept keys until settle phase
                if r.state == protocol.ResponseConversion && !r.keyAccepted {
                        return r, nil
                }
                return r, r.handleKey(m)

        case responseTickMsg:
                r.advanceDrama(now)
                r.particles.Tick(now)
                if r.dramaPhase != dramaNone {
                        return r, tickResponseCmdFast()
                }
                return r, tickResponseCmd()
        }

        return r, nil
}

// handleKey processes key events for the response screen.
func (r *Response) handleKey(msg tea.KeyMsg) tea.Cmd {
        switch r.state {
        case protocol.ResponsePositive, protocol.ResponseHotLead:
                return r.handlePositiveKeys(msg)
        case protocol.ResponseCurious:
                return r.handleCuriousKeys(msg)
        case protocol.ResponseNegative:
                return r.handleNegativeKeys(msg)
        case protocol.ResponseMaybe:
                return r.handleMaybeKeys(msg)
        case protocol.ResponseAutoReply:
                return r.handleAutoReplyKeys(msg)
        case protocol.ResponseStopDetected:
                return r.handleStopKeys(msg)
        case protocol.ResponseDealDetected:
                return r.handleDealKeys(msg)
        case protocol.ResponseOfferPreview:
                return r.handleOfferKeys(msg)
        case protocol.ResponseMultiQueue:
                return r.handleMultiQueueKeys(msg)
        case protocol.ResponseConversion:
                return r.handleConversionKeys(msg)
        default:
                return r.handlePositiveKeys(msg)
        }
}

// handlePositiveKeys handles keys for positive/hot-lead response states.
// ↵ send offer, 2 reply custom, 3 later
func (r *Response) handlePositiveKeys(msg tea.KeyMsg) tea.Cmd {
        switch {
        case key.Matches(msg, KeyEnter):
                // Only publish action; backend will navigate via HandleNavigate
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionSendOffer,
                        Screen: protocol.ScreenResponse,
                        Params: map[string]any{protocol.ParamBusiness: r.lead.Business},
                })
                return nil
        case key.Matches(msg, Key2):
                // Navigate to compose — backend decides target screen
                publish(r.bus, bus.NavigateMsg{
                        Screen: protocol.ScreenCompose,
                        Params: map[string]any{protocol.ParamBusiness: r.lead.Business},
                })
                return nil
        case key.Matches(msg, Key3):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionLater,
                        Screen: protocol.ScreenResponse,
                })
                return nil
        }
        return nil
}

// handleCuriousKeys handles keys for curious response state.
// ↵ send info price, 2 reply custom, 3 later
func (r *Response) handleCuriousKeys(msg tea.KeyMsg) tea.Cmd {
        switch {
        case key.Matches(msg, KeyEnter):
                // Only publish action; backend will navigate via HandleNavigate
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionSendPricing,
                        Screen: protocol.ScreenResponse,
                        Params: map[string]any{protocol.ParamBusiness: r.lead.Business},
                })
                return nil
        case key.Matches(msg, Key2):
                // Navigate to compose — backend decides target screen
                publish(r.bus, bus.NavigateMsg{
                        Screen: protocol.ScreenCompose,
                        Params: map[string]any{protocol.ParamBusiness: r.lead.Business},
                })
                return nil
        case key.Matches(msg, Key3):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionLater,
                        Screen: protocol.ScreenResponse,
                })
                return nil
        }
        return nil
}

// handleNegativeKeys handles keys for negative response state.
// 1 mark invalid, 2 follow up later, ↵ skip
func (r *Response) handleNegativeKeys(msg tea.KeyMsg) tea.Cmd {
        switch {
        case key.Matches(msg, Key1):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionMarkInvalid,
                        Screen: protocol.ScreenResponse,
                        Params: map[string]any{protocol.ParamBusiness: r.lead.Business},
                })
                return nil
        case key.Matches(msg, Key2):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionFollowUpLater,
                        Screen: protocol.ScreenResponse,
                })
                return nil
        case key.Matches(msg, KeyEnter):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionSkip,
                        Screen: protocol.ScreenResponse,
                })
                return nil
        }
        return nil
}

// handleMaybeKeys handles keys for maybe response state.
// ↵ send offer, 2 send info, 3 reply custom
func (r *Response) handleMaybeKeys(msg tea.KeyMsg) tea.Cmd {
        switch {
        case key.Matches(msg, KeyEnter):
                // Only publish action; backend will navigate via HandleNavigate
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionSendOffer,
                        Screen: protocol.ScreenResponse,
                })
                return nil
        case key.Matches(msg, Key2):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionSendInfo,
                        Screen: protocol.ScreenResponse,
                })
                return nil
        case key.Matches(msg, Key3):
                // Navigate to compose — backend decides target screen
                publish(r.bus, bus.NavigateMsg{
                        Screen: protocol.ScreenCompose,
                        Params: map[string]any{protocol.ParamBusiness: r.lead.Business},
                })
                return nil
        }
        return nil
}

// handleAutoReplyKeys handles keys for auto-reply detected state.
// ↵ skip, 1 still follow up, q back
func (r *Response) handleAutoReplyKeys(msg tea.KeyMsg) tea.Cmd {
        switch {
        case key.Matches(msg, KeyEnter):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionSkip,
                        Screen: protocol.ScreenResponse,
                })
                return nil
        case key.Matches(msg, Key1):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionStillFollowUp,
                        Screen: protocol.ScreenResponse,
                })
                return nil
        case key.Matches(msg, KeyBack):
                // Navigate to monitor — backend decides target screen
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionNavMonitor,
                        Screen: protocol.ScreenResponse,
                })
                return nil
        }
        return nil
}

// handleStopKeys handles keys for stop/cease-contact detected state.
// ↵ agree block, 2 block all niches, 3 cancel
func (r *Response) handleStopKeys(msg tea.KeyMsg) tea.Cmd {
        switch {
        case key.Matches(msg, KeyEnter):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionBlockConfirm,
                        Screen: protocol.ScreenResponse,
                        Params: map[string]any{protocol.ParamBusiness: r.lead.Business},
                })
                return nil
        case key.Matches(msg, Key2):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionBlockAllNiches,
                        Screen: protocol.ScreenResponse,
                        Params: map[string]any{protocol.ParamBusiness: r.lead.Business},
                })
                return nil
        case key.Matches(msg, Key3):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionCancelBlock,
                        Screen: protocol.ScreenResponse,
                })
                return nil
        }
        return nil
}

// handleDealKeys handles keys for deal-detected state.
// ↵ confirm deal, s not deal, 2 reply first
func (r *Response) handleDealKeys(msg tea.KeyMsg) tea.Cmd {
        switch {
        case key.Matches(msg, KeyEnter):
                // Only publish confirm_deal action; backend will navigate to ResponseConversion
                // via HandleNavigate which triggers startConversionDrama()
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionConfirmDeal,
                        Screen: protocol.ScreenResponse,
                        Params: map[string]any{protocol.ParamBusiness: r.lead.Business},
                })
                return nil
        case key.Matches(msg, KeySkip):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionNotDeal,
                        Screen: protocol.ScreenResponse,
                        Params: map[string]any{protocol.ParamBusiness: r.lead.Business},
                })
                return nil
        case key.Matches(msg, Key2):
                // Navigate to compose — backend decides target screen
                publish(r.bus, bus.NavigateMsg{
                        Screen: protocol.ScreenCompose,
                        Params: map[string]any{protocol.ParamBusiness: r.lead.Business},
                })
                return nil
        }
        return nil
}

// handleOfferKeys handles keys for offer preview state.
// ↵ send, 2 change template, 3 edit
func (r *Response) handleOfferKeys(msg tea.KeyMsg) tea.Cmd {
        switch {
        case key.Matches(msg, KeyEnter):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionSendOffer,
                        Screen: protocol.ScreenResponse,
                        Params: map[string]any{protocol.ParamBusiness: r.lead.Business},
                })
                return nil
        case key.Matches(msg, Key2):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionChangeTemplate,
                        Screen: protocol.ScreenResponse,
                })
                return nil
        case key.Matches(msg, Key3):
                // Navigate to compose — backend decides target screen
                publish(r.bus, bus.NavigateMsg{
                        Screen: protocol.ScreenCompose,
                        Params: map[string]any{protocol.ParamBusiness: r.lead.Business},
                })
                return nil
        }
        return nil
}

// handleMultiQueueKeys handles keys for multi-response queue state.
// ↵ process one by one, 1 auto-offer positives, 2 auto per type
func (r *Response) handleMultiQueueKeys(msg tea.KeyMsg) tea.Cmd {
        switch {
        case key.Matches(msg, KeyEnter):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionProcessOne,
                        Screen: protocol.ScreenResponse,
                })
                return nil
        case key.Matches(msg, Key1):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionAutoOfferPos,
                        Screen: protocol.ScreenResponse,
                })
                return nil
        case key.Matches(msg, Key2):
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionAutoPerType,
                        Screen: protocol.ScreenResponse,
                })
                return nil
        }
        return nil
}

// handleConversionKeys handles keys for the conversion state.
// Only ↵ (mark as converted) is accepted after settle phase.
func (r *Response) handleConversionKeys(msg tea.KeyMsg) tea.Cmd {
        if key.Matches(msg, KeyEnter) {
                publish(r.bus, bus.ActionMsg{
                        Action: protocol.ActionMarkConverted,
                        Screen: protocol.ScreenResponse,
                        Params: map[string]any{protocol.ParamBusiness: r.conversion.Business},
                })
        }
        return nil
}

// ---------------------------------------------------------------------------
// Conversion drama — 4-phase animation sequence
// ---------------------------------------------------------------------------

// startConversionDrama initiates the full 4-phase conversion drama.
func (r *Response) startConversionDrama() {
        r.dramaPhase = dramaShock
        r.dramaStart = time.Now()
        r.keyAccepted = false
        // Terminal bell double-tap per doc spec
        fmt.Print("\a\a")
        if r.width > 0 && r.height > 0 {
                r.particles = component.NewParticleSystem(r.width, r.height)
                r.particles.Burst()
        }
}

// advanceDrama progresses through the 4 conversion drama phases based on
// elapsed time since dramaStart. The timing follows the exact spec from
// doc/05-screens-monitor-response.md:
//   - SHOCK: 0-200ms
//   - REVEAL: 200-800ms
//   - CONTEXT: 800-1500ms
//   - SETTLE: 1500ms+
func (r *Response) advanceDrama(now time.Time) {
        if r.dramaPhase == dramaNone {
                return
        }
        elapsed := now.Sub(r.dramaStart)
        shockEnd := anim.ConversionShockDuration
        revealEnd := shockEnd + anim.ConversionRevealDuration
        contextEnd := revealEnd + anim.ConversionContextDuration

        switch {
        case elapsed < shockEnd:
                r.dramaPhase = dramaShock
        case elapsed < revealEnd:
                if r.dramaPhase == dramaShock {
                        // Transition from shock to reveal: trigger particle burst
                        r.dramaPhase = dramaReveal
                        if r.particles.Width > 0 {
                                r.particles.Burst()
                        }
                }
        case elapsed < contextEnd:
                r.dramaPhase = dramaContext
        default:
                r.dramaPhase = dramaSettle
                // Accept keyboard after settle hold
                if elapsed > contextEnd+anim.ConversionSettleHold {
                        r.keyAccepted = true
                }
        }
}

// ---------------------------------------------------------------------------
// View — renders the current response state
// ---------------------------------------------------------------------------

// View implements tea.Model.
func (r *Response) View() string {
        switch r.state {
        case protocol.ResponsePositive:
                return r.viewPositive()
        case protocol.ResponseCurious:
                return r.viewCurious()
        case protocol.ResponseNegative:
                return r.viewNegative()
        case protocol.ResponseMaybe:
                return r.viewMaybe()
        case protocol.ResponseAutoReply:
                return r.viewAutoReply()
        case protocol.ResponseStopDetected:
                return r.viewStopDetected()
        case protocol.ResponseDealDetected:
                return r.viewDealDetected()
        case protocol.ResponseHotLead:
                return r.viewHotLead()
        case protocol.ResponseOfferPreview:
                return r.viewOfferPreview()
        case protocol.ResponseMultiQueue:
                return r.viewMultiQueue()
        case protocol.ResponseConversion:
                return r.viewConversion()
        default:
                return r.viewPositive()
        }
}

// viewPositive renders a positive response with offer options.
func (r *Response) viewPositive() string {
        return r.renderResponseHeader(i18n.T(i18n.KeyResponseGotReply), style.Success) +
                r.renderLeadInfo() +
                r.renderMessage() +
                r.renderSeparator() +
                r.renderActions([]actionItem{
                        {key: "↵", label: i18n.T(i18n.KeyResponseSendOffer), accent: true},
                        {key: "2", label: i18n.T(i18n.KeyResponseReplyCustom)},
                        {key: "3", label: i18n.T(i18n.KeyResponseLater)},
                })
}

// viewCurious renders a curious response with pricing info option.
func (r *Response) viewCurious() string {
        return r.renderResponseHeader(i18n.T(i18n.KeyResponseGotReply), style.Warning) +
                r.renderLeadInfo() +
                r.renderMessage() +
                r.renderSeparator() +
                r.renderActions([]actionItem{
                        {key: "↵", label: i18n.T(i18n.KeyResponseSendInfo), accent: true},
                        {key: "2", label: i18n.T(i18n.KeyResponseReplyCustom)},
                        {key: "3", label: i18n.T(i18n.KeyResponseLater)},
                })
}

// viewNegative renders a negative response with skip option.
func (r *Response) viewNegative() string {
        return r.renderResponseHeader(i18n.T(i18n.KeyResponseGotReplyPlain), style.TextMuted) +
                r.renderLeadInfo() +
                r.renderMessage() +
                r.renderSeparator() +
                r.renderActions([]actionItem{
                        {key: "1", label: i18n.T(i18n.KeyResponseMarkInvalid)},
                        {key: "2", label: i18n.T(i18n.KeyResponseFollowUpLater)},
                        {key: "↵", label: i18n.T(i18n.KeyResponseSkipIt), accent: true},
                })
}

// viewMaybe renders a maybe/ambiguous response.
func (r *Response) viewMaybe() string {
        return r.renderResponseHeader(i18n.T(i18n.KeyResponseGotReplyPlain), style.TextMuted) +
                r.renderLeadInfo() +
                r.renderMessage() +
                r.renderSeparator() +
                r.renderActions([]actionItem{
                        {key: "↵", label: i18n.T(i18n.KeyResponseSendOffer), accent: true},
                        {key: "2", label: i18n.T(i18n.KeyResponseSendInfo)},
                        {key: "3", label: i18n.T(i18n.KeyResponseReplyCustom)},
                })
}

// viewAutoReply renders a detected auto-reply response.
func (r *Response) viewAutoReply() string {
        return r.renderResponseHeader(i18n.T(i18n.KeyResponseGotReplyAuto), style.TextDim) +
                r.renderLeadInfo() +
                r.renderMessage() +
                r.renderNotice(i18n.T(i18n.KeyResponseAutoReplySkip), style.TextDim) +
                r.renderSeparator() +
                r.renderActions([]actionItem{
                        {key: "↵", label: i18n.T(i18n.KeyResponseSkipIt), accent: true},
                        {key: "1", label: i18n.T(i18n.KeyResponseStillFollowUp)},
                        {key: "q", label: i18n.T(i18n.KeyLabelBack)},
                })
}

// viewStopDetected renders a stop/cease-contact detected response.
func (r *Response) viewStopDetected() string {
        return r.renderResponseHeader(i18n.T(i18n.KeyResponseStopDetected), style.Danger) +
                r.renderLeadInfo() +
                r.renderMessage() +
                r.renderNotice(i18n.T(i18n.KeyResponseStopNotice1), style.Warning) +
                r.renderNotice(i18n.T(i18n.KeyResponseStopNotice2), style.TextMuted) +
                r.renderNotice(i18n.T(i18n.KeyResponseStopNotice3), style.TextDim) +
                r.renderSeparator() +
                r.renderActions([]actionItem{
                        {key: "↵", label: i18n.T(i18n.KeyResponseBlockConfirm), accent: true},
                        {key: "2", label: i18n.T(i18n.KeyResponseBlockAllNiche)},
                        {key: "3", label: i18n.T(i18n.KeyResponseCancelBlock)},
                })
}

// viewDealDetected renders a deal-detected (closing trigger) response.
func (r *Response) viewDealDetected() string {
        header := "💬 " + i18n.T(i18n.KeyResponseGotReply) + " — 🎯 " + i18n.T(i18n.KeyResponseDealDetected)
        return r.renderCustomHeader(header, style.Warning) +
                r.renderLeadInfo() +
                r.renderMessage() +
                r.renderTriggerNotice() +
                r.renderNotice(i18n.T(i18n.KeyResponseDealNotice1), style.Text) +
                r.renderNotice(i18n.T(i18n.KeyResponseDealNotice2), style.TextMuted) +
                r.renderSeparator() +
                r.renderActions([]actionItem{
                        {key: "↵", label: i18n.T(i18n.KeyResponseConfirmDeal), accent: true},
                        {key: "s", label: i18n.T(i18n.KeyResponseNotDeal)},
                        {key: "2", label: i18n.T(i18n.KeyResponseReplyFirst)},
                })
}

// viewHotLead renders a hot-lead detected response.
func (r *Response) viewHotLead() string {
        header := "💬 " + i18n.T(i18n.KeyResponseGotReply) + " — 🔥 " + i18n.T(i18n.KeyResponseHotLead)
        return r.renderCustomHeader(header, style.Warning) +
                r.renderLeadInfo() +
                r.renderMessage() +
                r.renderHotLeadNotices() +
                r.renderSeparator() +
                r.renderActions([]actionItem{
                        {key: "↵", label: i18n.T(i18n.KeyResponseOfferNow), accent: true},
                        {key: "2", label: i18n.T(i18n.KeyResponseReplyCustom)},
                        {key: "3", label: i18n.T(i18n.KeyResponseLater)},
                })
}

// viewOfferPreview renders the offer preview before sending.
func (r *Response) viewOfferPreview() string {
        var b strings.Builder

        // Title
        b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Bold(true).Render(
                fmt.Sprintf("%s %s", i18n.T(i18n.KeyResponseOfferPreview), r.lead.Business),
        ))
        b.WriteString("\n\n")

        // Separator
        b.WriteString(r.renderSeparator())

        // Offer text — rendered by backend with template substitution applied
        offerText := r.lead.OfferText
        if offerText == "" {
                offerText = r.lead.Message // Fallback for older backends
        }
        b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(offerText))
        b.WriteString("\n\n")

        // Separator
        b.WriteString(r.renderSeparator())

        // Actions
        b.WriteString(r.renderActions([]actionItem{
                {key: "↵", label: i18n.T(i18n.KeyResponseSend), accent: true},
                {key: "2", label: i18n.T(i18n.KeyResponseChangeTemplate)},
                {key: "3", label: i18n.T(i18n.KeyResponseEditFirst)},
        }))

        return b.String()
}

// viewMultiQueue renders multiple responses queued simultaneously.
func (r *Response) viewMultiQueue() string {
        var b strings.Builder

        // Header
        count := len(r.queue)
        b.WriteString(lipgloss.NewStyle().Foreground(style.Warning).Bold(true).Render(
                fmt.Sprintf("💬 %d %s", count, i18n.T(i18n.KeyResponseMultiQueue)),
        ))
        b.WriteString("\n\n")

        // Separator
        b.WriteString(r.renderSeparator())

        // Queue items
        for i, item := range r.queue {
                prefix := "  "
                focusMarker := "  "
                if i == r.cursor {
                        focusMarker = lipgloss.NewStyle().Foreground(style.Text).Render("▸")
                }
                indexLabel := lipgloss.NewStyle().Foreground(style.TextDim).Render(
                        fmt.Sprintf("%02d", item.Index),
                )
                bizName := lipgloss.NewStyle().Foreground(style.Text).Render(item.Business)
                category := lipgloss.NewStyle().Foreground(style.TextDim).Render(
                        fmt.Sprintf("%s · %s", item.Category, item.Area),
                )

                var classBadge string
                switch item.Class {
                case classPositif:
                        classBadge = lipgloss.NewStyle().Foreground(style.Warning).Render(
                                fmt.Sprintf("● %s", i18n.T(i18n.KeyResponseClassPositif)))
                case classCurious:
                        classBadge = lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                                fmt.Sprintf("○ %s", i18n.T(i18n.KeyResponseClassCurious)))
                case classAutoReply:
                        classBadge = lipgloss.NewStyle().Foreground(style.TextDim).Render(
                                fmt.Sprintf("○ %s", i18n.T(i18n.KeyResponseClassAutoReply)))
                default:
                        classBadge = lipgloss.NewStyle().Foreground(style.TextDim).Render(
                                fmt.Sprintf("○ %s", item.Class))
                }

                b.WriteString(fmt.Sprintf("%s%s %s\n", prefix, focusMarker, indexLabel))
                b.WriteString(fmt.Sprintf("%s    %s\n", prefix, bizName))
                b.WriteString(fmt.Sprintf("%s    %s\n", prefix, category))
                b.WriteString(fmt.Sprintf("%s    %s\n", prefix,
                        lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                                fmt.Sprintf("\"%s\"", item.Message),
                        ),
                ))
                b.WriteString(fmt.Sprintf("%s    %s\n", prefix, classBadge))
                b.WriteString("\n")
        }

        // Separator
        b.WriteString(r.renderSeparator())

        // Actions
        b.WriteString(r.renderActions([]actionItem{
                {key: "↵", label: i18n.T(i18n.KeyResponseProcessOne), accent: true},
                {key: "1", label: i18n.T(i18n.KeyResponseAutoOfferPos)},
                {key: "2", label: i18n.T(i18n.KeyResponseAutoPerType)},
        }))

        return b.String()
}

// viewConversion renders the conversion drama sequence.
// This is the most important screen — the moment a deal closes.
func (r *Response) viewConversion() string {
        now := time.Now()

        switch r.dramaPhase {
        case dramaShock:
                return r.viewConversionShock(now)
        case dramaReveal:
                return r.viewConversionReveal(now)
        case dramaContext:
                return r.viewConversionContext(now)
        case dramaSettle:
                return r.viewConversionSettle(now)
        default:
                return r.viewConversionSettle(now)
        }
}

// viewConversionShock renders Phase 1: SHOCK (0-200ms).
// Full white flash — entire screen flash PUTIH full-brightness.
func (r *Response) viewConversionShock(now time.Time) string {
        elapsed := now.Sub(r.dramaStart)
        // Flash intensity fades over the shock duration
        intensity := 1.0 - float64(elapsed)/float64(anim.ConversionShockDuration)
        if intensity < 0 {
                intensity = 0
        }

        // During shock, render a near-white screen
        flashStyle := lipgloss.NewStyle().
                Foreground(style.Celebration).
                Background(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x",
                        int(10+245*intensity), int(10+245*intensity), int(11+244*intensity)))).
                Width(r.width).
                Height(r.height)

        return flashStyle.Render(" ")
}

// viewConversionReveal renders Phase 2: REVEAL (200-800ms).
// DEAL text scales in, particles cascade, borders draw inward.
func (r *Response) viewConversionReveal(now time.Time) string {
        elapsed := now.Sub(r.dramaStart) - anim.ConversionShockDuration
        revealProgress := float64(elapsed) / float64(anim.ConversionRevealDuration)
        if revealProgress > 1.0 {
                revealProgress = 1.0
        }

        // Scale with overshoot: 0 → 1.3 → 1.0 (using EaseOutBack)
        scaledProgress := anim.EaseOutBack(revealProgress)
        if scaledProgress > anim.JackpotOvershoot {
                scaledProgress = anim.JackpotOvershoot
        }

        var b strings.Builder

        // Draw ━━ borders from edges inward
        borderWidth := max(r.width-4, minRenderWidth)
        drawnWidth := int(float64(borderWidth) * revealProgress)
        if drawnWidth < 2 {
                drawnWidth = 2
        }

        leftBorder := lipgloss.NewStyle().Foreground(style.Gold).Render(
                strings.Repeat("━", drawnWidth/2),
        )
        rightBorder := lipgloss.NewStyle().Foreground(style.Gold).Render(
                strings.Repeat("━", drawnWidth/2),
        )
        spaceBetween := borderWidth - drawnWidth
        if spaceBetween < 0 {
                spaceBetween = 0
        }

        b.WriteString(leftBorder)
        b.WriteString(strings.Repeat(" ", spaceBetween))
        b.WriteString(rightBorder)
        b.WriteString("\n")
        b.WriteString(leftBorder)
        b.WriteString(strings.Repeat(" ", spaceBetween))
        b.WriteString(rightBorder)
        b.WriteString("\n\n")

        // DEAL text with scale effect
        dealText := i18n.T(i18n.KeyResponseConversionDeal)
        if scaledProgress > 0.1 {
                // Color wave: accent → success → normal during reveal
                var dealColor lipgloss.Color
                if revealProgress < 0.4 {
                        dealColor = style.Gold
                } else if revealProgress < 0.7 {
                        dealColor = style.Success
                } else {
                        dealColor = style.Gold
                }
                b.WriteString(lipgloss.NewStyle().
                        Foreground(dealColor).
                        Bold(true).
                        Render(dealText))
        }
        b.WriteString("\n\n")

        // Particles
        if r.particles.Active {
                b.WriteString(r.particles.ViewCompact())
        }

        return b.String()
}

// viewConversionContext renders Phase 3: CONTEXT (800-1500ms).
// Business name + timeline fade in, trophy bounces from right.
func (r *Response) viewConversionContext(now time.Time) string {
        b := r.renderConversionBody(true)
        return b.String()
}

// viewConversionSettle renders Phase 4: SETTLE (1500ms+).
// Particles dissolve, screen settles, keyboard accepted after hold.
func (r *Response) viewConversionSettle(now time.Time) string {
        b := r.renderConversionBody(false)

        // Mark as converted — fades in only after settle hold
        if r.keyAccepted {
                b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render("↵"))
                b.WriteString("  ")
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                        i18n.T(i18n.KeyConversionMarkConverted),
                ))
        }

        return b.String()
}

// renderConversionBody is the shared renderer for the conversion screen body.
// Used by both CONTEXT and SETTLE phases to avoid DRY violation.
// showParticles controls whether particle box label says "dissolving".
func (r *Response) renderConversionBody(showParticles bool) strings.Builder {
        var b strings.Builder

        // Draw ━━ celebration borders (intentional design exception for conversion)
        borderWidth := max(r.width-4, minRenderWidth)
        borderLine := lipgloss.NewStyle().Foreground(style.Gold).Render(
                strings.Repeat("━", borderWidth),
        )
        b.WriteString(borderLine + "\n")
        b.WriteString(borderLine + "\n\n")

        // DEAL text
        b.WriteString(lipgloss.NewStyle().Foreground(style.Gold).Bold(true).Render(
                i18n.T(i18n.KeyResponseConversionDeal),
        ))
        b.WriteString("\n\n")

        // Business info
        b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(
                fmt.Sprintf("%s · %s · %s", r.lead.Business, r.lead.Category, r.lead.Area),
        ))
        b.WriteString("\n\n")

        // Pipeline timeline — Doc: "convert dari ice breaker → offer → deal" + "waktu: 2 hari 4 jam"
        pipelineText := r.conversion.Pipeline
        if pipelineText == "" {
                pipelineText = i18n.T(i18n.KeyConversionFromPipeline)
        }
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                pipelineText,
        ))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("%s %s", i18n.T(i18n.KeyConversionTimeTaken), r.conversion.TimeTaken),
        ))
        b.WriteString("\n\n")

        // Doc: separator uses ── lines (same visual rule as dashboard)
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                strings.Repeat("─", max(r.width-4, minRenderWidth)),
        ))
        b.WriteString("\n\n")

        // Particle box — use showParticles to control label
        particleBoxWidth := min(r.width-8, 35)
        if particleBoxWidth < 20 {
                particleBoxWidth = 20
        }
        particleBox := component.ParticleBox{Width: particleBoxWidth}
        if showParticles {
                particleBox.Label = i18n.T(i18n.KeyConversionCelebrating)
        } else {
                particleBox.Label = i18n.T(i18n.KeyConversionDissolving)
        }
        b.WriteString(particleBox.View())
        b.WriteString("\n\n")

        // Separator
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                strings.Repeat("─", max(r.width-4, minRenderWidth)),
        ))
        b.WriteString("\n\n")

        // Doc: 🏆 conversion ke-3 minggu ini / 🏆 total revenue minggu ini: rp 7.5jt
        ordinal := fmt.Sprintf(i18n.T(i18n.KeyConversionOrdinal), r.conversion.TrophyCount)
        b.WriteString(lipgloss.NewStyle().Foreground(style.Gold).Bold(true).Render(
                fmt.Sprintf("🏆 %s %s %s", ordinal, i18n.T(i18n.KeyConversionTrophyWeek), i18n.T(i18n.KeyMonitorThisWeek)),
        ))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Gold).Render(
                fmt.Sprintf("🏆 %s: %s", i18n.T(i18n.KeyConversionRevenueWeek), r.conversion.Revenue),
        ))
        b.WriteString("\n\n")

        return b
}

// ---------------------------------------------------------------------------
// Shared render helpers — DRY sub-renderers
// ---------------------------------------------------------------------------

// actionItem defines a key-label pair for the action bar.
type actionItem struct {
        key    string
        label  string
        accent bool
}

// renderResponseHeader renders the header line with emoji and title,
// followed by a separator line matching the doc spec (── after header).
// The emoji is determined by the current response state per doc spec:
//   💬 positive/curious, 🛑 stop, 🎯 deal, 🔥 hot lead
func (r *Response) renderResponseHeader(title string, color lipgloss.Color) string {
        var b strings.Builder
        // Prepend icon based on current state per doc spec
        var icon string
        switch r.state {
        case protocol.ResponseStopDetected:
                icon = "🛑 "
        case protocol.ResponseDealDetected:
                icon = "🎯 "
        case protocol.ResponseHotLead:
                icon = "🔥 "
        default:
                // Positive, curious, maybe, auto-reply all get 💬
                icon = "💬 "
        }
        b.WriteString(lipgloss.NewStyle().Foreground(color).Bold(true).Render(icon + title))
        b.WriteString("\n\n")
        // Doc: separator line after header (same ── rule as dashboard)
        b.WriteString(r.renderSeparator())
        return b.String()
}

// renderSeparator renders a ── line matching the doc spec.
func (r *Response) renderSeparator() string {
        return renderDimSeparator(r.width) + "\n"
}

// renderLeadInfo renders the business/category/area line.
func (r *Response) renderLeadInfo() string {
        return lipgloss.NewStyle().Foreground(style.Text).Render(
                fmt.Sprintf("%s · %s · %s", r.lead.Business, r.lead.Category, r.lead.Area),
        ) + "\n\n"
}

// renderMessage renders the response message as a quote.
func (r *Response) renderMessage() string {
        // Doc spec: message displayed in quotes, e.g. "iya kak, boleh lihat"
        // Do NOT use Go's %q which adds escaping; just wrap in literal quotes.
        return lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                fmt.Sprintf("\"%s\"", r.lead.Message),
        ) + "\n\n"
}

// renderNotice renders an informational notice block.
func (r *Response) renderNotice(text string, color lipgloss.Color) string {
        return lipgloss.NewStyle().Foreground(color).Render(text) + "\n\n"
}

// renderTriggerNotice renders the closing trigger match notice.
// Doc: ⚡ closing trigger: "dah transfer" → auto-deal
func (r *Response) renderTriggerNotice() string {
        if r.lead.Trigger == "" {
                return ""
        }
        var b strings.Builder
        b.WriteString(lipgloss.NewStyle().Foreground(style.Warning).Bold(true).Render(
                i18n.T(i18n.KeyResponseClosingTrigger),
        ))
        b.WriteString(" ")
        b.WriteString(lipgloss.NewStyle().Foreground(style.Warning).Render(
                fmt.Sprintf("%q → %s", r.lead.Trigger, i18n.T(i18n.KeyResponseAutoDealShort)),
        ))
        b.WriteString("\n\n")
        return b.String()
}

// renderCustomHeader renders a header with custom text and color.
// Used by hot lead and deal detected views that need a composite header.
func (r *Response) renderCustomHeader(text string, color lipgloss.Color) string {
        return lipgloss.NewStyle().Foreground(color).Bold(true).Render(text) + "\n\n"
}

// renderHotLeadNotices renders the trigger notice, auto-prioritize notice,
// and offer hint for hot lead responses.
func (r *Response) renderHotLeadNotices() string {
        var b strings.Builder
        triggerText := r.lead.Trigger
        if triggerText == "" {
                return ""
        }
        b.WriteString(lipgloss.NewStyle().Foreground(style.Warning).Render(
                fmt.Sprintf(i18n.T(i18n.KeyResponseHotLeadTrigger), triggerText)))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(
                i18n.T(i18n.KeyResponseHotLeadAutoPrior)))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                i18n.T(i18n.KeyResponseHotLeadOfferHint)))
        b.WriteString("\n")
        return b.String()
}

// renderActions renders the action key bar at the bottom of the screen.
// The i18n labels already include the key prefix (e.g. "2 balas custom"),
// so we highlight the key portion in accent and render the rest in muted.
func (r *Response) renderActions(actions []actionItem) string {
        var parts []string
        for _, a := range actions {
                keyStyle := lipgloss.NewStyle().Foreground(style.Accent)
                if a.accent {
                        keyStyle = lipgloss.NewStyle().Foreground(style.Accent).Bold(true)
                }
                // Label already contains the key prefix from i18n.
                // Highlight the key part in accent, rest in muted.
                if len(a.label) > 0 && len(a.key) > 0 && string(a.label[0]) == a.key {
                        parts = append(parts, fmt.Sprintf("%s%s",
                                keyStyle.Render(a.key),
                                lipgloss.NewStyle().Foreground(style.TextMuted).Render(a.label[1:])),
                        )
                } else {
                        parts = append(parts, fmt.Sprintf("%s  %s",
                                keyStyle.Render(a.key),
                                lipgloss.NewStyle().Foreground(style.TextMuted).Render(a.label)),
                        )
                }
        }
        return strings.Join(parts, "    ") + "\n"
}

// populateData extracts data from backend params.
func (r *Response) populateData(params map[string]any) {
        // Support both nested ("lead": {...}) and flat ("business": "...") param formats.
        // The backend may send lead fields as flat params for single-response states.
        if v, ok := params[protocol.ParamLead]; ok {
                r.lead = toLeadResponse(v)
        } else {
                // Extract flat lead params directly from the top-level params map.
                r.lead = toLeadResponse(params)
        }
        if v, ok := params[protocol.ParamConversion]; ok {
                r.conversion = toConversionData(v)
        } else {
                // Extract flat conversion params directly from the top-level params map.
                r.conversion = toConversionData(params)
        }
        if v, ok := params[protocol.ParamQueue]; ok {
                r.queue = toMultiQueue(v)
        }
        r.cursor = extractInt(params, protocol.ParamCursor)
}

func toLeadResponse(v any) LeadResponse {
        if lr, ok := v.(LeadResponse); ok {
                return lr
        }
        m, ok := v.(map[string]any)
        if !ok {
                return LeadResponse{}
        }
        lr := LeadResponse{
                Business:  extractString(m, protocol.ParamBusiness),
                Category:  extractString(m, protocol.ParamCategory),
                Area:      extractString(m, protocol.ParamArea),
                Message:   extractString(m, protocol.ParamMessage),
                OfferText: extractString(m, protocol.ParamOfferText),
                Trigger:   extractString(m, protocol.ParamTrigger),
                Niche:     extractString(m, protocol.ParamNiche),
                Class:     extractString(m, protocol.ParamClass),
        }
        return lr
}

func toConversionData(v any) ConversionData {
        if cd, ok := v.(ConversionData); ok {
                return cd
        }
        m, ok := v.(map[string]any)
        if !ok {
                return ConversionData{}
        }
        cd := ConversionData{
                Business:    extractString(m, protocol.ParamBusiness),
                Pipeline:    extractString(m, protocol.ParamPipeline),
                TimeTaken:   extractString(m, protocol.ParamTimeTaken),
                TrophyCount: extractInt(m, protocol.ParamTrophyCount),
                Revenue:     extractString(m, protocol.ParamRevenue),
        }
        return cd
}

func toMultiQueue(v any) []MultiQueueItem {
        if items, ok := v.([]MultiQueueItem); ok {
                return items
        }
        arr, ok := v.([]any)
        if !ok {
                return nil
        }
        var result []MultiQueueItem
        for i, item := range arr {
                m, ok := item.(map[string]any)
                if !ok {
                        continue
                }
                qi := MultiQueueItem{
                        Index:    i + 1,
                        Business: extractString(m, protocol.ParamBusiness),
                        Category: extractString(m, protocol.ParamCategory),
                        Area:     extractString(m, protocol.ParamArea),
                        Message:  extractString(m, protocol.ParamMessage),
                        Class:    extractString(m, protocol.ParamClass),
                }
                result = append(result, qi)
        }
        return result
}

// ---------------------------------------------------------------------------
// Tick messages
// ---------------------------------------------------------------------------

// responseTickMsg is the periodic tick for the response screen.
type responseTickMsg time.Time

// tickResponseCmd returns a tick at the standard data rain interval.
func tickResponseCmd() tea.Cmd {
        return tea.Tick(anim.DataRainUpdateInterval, func(t time.Time) tea.Msg {
                return responseTickMsg(t)
        })
}

// tickResponseCmdFast returns a faster tick for drama animation frames.
func tickResponseCmdFast() tea.Cmd {
        return tea.Tick(16*time.Millisecond, func(t time.Time) tea.Msg {
                return responseTickMsg(t)
        })
}
