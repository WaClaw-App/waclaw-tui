// Package comms implements the Communication screens: Compose, History, and Follow-Up.
// These are screens 15–17 from the WaClaw TUI specification.
//
// Doc source: doc/09-screens-communicate.md
package comms

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
        tea "github.com/charmbracelet/bubbletea"
        "github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Snippet data
// ---------------------------------------------------------------------------

// Snippet represents a quick-reply snippet loaded from ~/.waclaw/snippets.md.
type Snippet struct {
        // Text is the snippet content.
        Text string

        // Category is the tone label shown to the user (e.g. "soft pitch", "free offer").
        Category string
}

// defaultSnippets provides an empty fallback when the backend hasn't
// pushed snippet data yet. Snippets should come from the backend
// (which reads ~/.waclaw/snippets.md), NOT be hardcoded in the TUI.
// The old hardcoded Indonesian strings were a 3G backend-frontend concern.
var defaultSnippets []Snippet

// snippetCategoryI18n maps snippet category keys to i18n keys for display-time lookup.
var snippetCategoryI18n = map[string]string{
        "soft_pitch":   i18n.KeyComposeCatSoftPitch,
        "free_offer":   i18n.KeyComposeCatFreeOffer,
        "move_to_call": i18n.KeyComposeCatMoveToCall,
        "direct_price": i18n.KeyComposeCatDirectPrice,
        "send_sample":  i18n.KeyComposeCatSendSample,
}

// ---------------------------------------------------------------------------
// Compose Model
// ---------------------------------------------------------------------------

// Compose is a MODAL OVERLAY screen for composing custom replies.
// It slides up from the bottom on top of the Response screen.
//
// States (from doc/09-screens-communicate.md):
//   - ComposeDraft: text input area for custom reply (double-enter send)
//   - ComposePreview: preview before sending (single enter)
//   - ComposeTemplatePick: quick-pick from template snippets
type Compose struct {
        screenBase
        state   protocol.StateID
        bus     *bus.Bus

        // Target is the business name the message will be sent to.
        Target string

        // Draft holds the current composed text.
        Draft string

        // CursorPos is the cursor position within Draft.
        CursorPos int

        // ConsecutiveEnters counts consecutive enter presses for double-enter send.
        ConsecutiveEnters int

        // PreviewEnteredAt tracks when the user entered preview (for hold delay).
        PreviewEnteredAt time.Time

        // Snippets holds the available quick-reply snippets.
        Snippets []Snippet

        // SnippetCursor is the selected snippet index.
        SnippetCursor int

        // CharCount tracks the character count for subtle awareness.
        CharCount int

        // MaxChars is the soft character limit.
        MaxChars int

        // Width and Height for layout.
        Width  int
        Height int

        // OpenedAt tracks when compose was opened (for slide-up animation).
        OpenedAt time.Time

        // TypeOut is the template preview for the compose preview state.
        TypeOut component.TemplatePreview
}

// NewCompose creates a Compose screen in draft state.
func NewCompose() *Compose {
        return &Compose{
                screenBase: screenBase{id: protocol.ScreenCompose},
                state:    protocol.ComposeDraft,
                MaxChars: composeDefaultMaxChars,
        }
}

// ---------------------------------------------------------------------------
// Screen interface
// ---------------------------------------------------------------------------

func (c *Compose) SetBus(b *bus.Bus) { c.bus = b }

func (c *Compose) Focus() {
        c.OpenedAt = time.Now()
        c.ConsecutiveEnters = 0
        if c.state == protocol.ComposeDraft {
                c.PreviewEnteredAt = time.Time{}
        }
}

func (c *Compose) Blur() {}

// ConsumesKey implements tui.KeyConsumer. Compose has sub-states (preview,
// template_pick) where "q" should navigate back locally to draft instead of
// popping the navigation stack.
func (c *Compose) ConsumesKey(msg tea.KeyMsg) bool {
        switch msg.String() {
        case "q":
                return c.state == protocol.ComposePreview || c.state == protocol.ComposeTemplatePick
        }
        return false
}

// HandleNavigate processes navigate commands from the backend.
func (c *Compose) HandleNavigate(params map[string]any) error {
        if t, ok := params[protocol.ParamTarget].(string); ok {
                c.Target = t
        }
        if state, ok := params[protocol.ParamState].(string); ok {
                c.state = protocol.StateID(state)
        }
        if snippets, ok := params[protocol.ParamSnippets].([]any); ok {
                c.Snippets = parseSnippets(snippets)
        }
        if mc, ok := params[protocol.ParamMaxChars]; ok {
                if f64, ok := mc.(float64); ok {
                        c.MaxChars = int(f64)
                }
        }
        return nil
}

// HandleUpdate processes update commands from the backend.
func (c *Compose) HandleUpdate(params map[string]any) error {
        if t, ok := params[protocol.ParamTarget].(string); ok {
                c.Target = t
        }
        if draft, ok := params[protocol.ParamDraft].(string); ok {
                c.Draft = draft
                c.CursorPos = len(draft)
                c.CharCount = len(draft)
        }
        if state, ok := params[protocol.ParamState].(string); ok {
                c.state = protocol.StateID(state)
        }
        if mc, ok := params[protocol.ParamMaxChars]; ok {
                if f64, ok := mc.(float64); ok {
                        c.MaxChars = int(f64)
                }
        }
        return nil
}

// Init implements tea.Model.
func (c *Compose) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (c *Compose) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch m := msg.(type) {
        case tea.WindowSizeMsg:
                c.Width = m.Width
                c.Height = m.Height
                return c, nil

        case tea.KeyMsg:
                return c.handleKey(m)
        }

        return c, nil
}

// handleKey routes key events to the current state handler.
func (c *Compose) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch c.state {
        case protocol.ComposeDraft:
                return c.handleDraftKey(msg)
        case protocol.ComposePreview:
                return c.handlePreviewKey(msg)
        case protocol.ComposeTemplatePick:
                return c.handlePickKey(msg)
        }
        return c, nil
}

// ---------------------------------------------------------------------------
// ComposeDraft key handling
// ---------------------------------------------------------------------------

func (c *Compose) handleDraftKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.Type {
        case tea.KeyEsc:
                // Close compose overlay.
                c.reset()
                return c, nil

        case tea.KeyTab:
                // Switch to template pick mode.
                c.state = protocol.ComposeTemplatePick
                c.SnippetCursor = 0
                return c, nil

        case tea.KeyEnter:
                c.ConsecutiveEnters++
                if c.ConsecutiveEnters >= composeSendEnterCount {
                        // Double-enter: move to preview.
                        c.transitionToPreview()
                        return c, nil
                }
                // Single enter: add newline.
                c.insertChar("\n")
                return c, nil

        case tea.KeyBackspace:
                c.backspace()
                c.ConsecutiveEnters = 0
                return c, nil

        case tea.KeyRunes:
                c.insertChar(msg.String())
                c.ConsecutiveEnters = 0
                return c, nil
        }

        c.ConsecutiveEnters = 0
        return c, nil
}

// ---------------------------------------------------------------------------
// ComposePreview key handling
// ---------------------------------------------------------------------------

func (c *Compose) handlePreviewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        // Enforce the 300ms hold before accepting enter.
        if !c.PreviewEnteredAt.IsZero() && time.Since(c.PreviewEnteredAt) < anim.ComposePreviewHold {
                return c, nil
        }

        switch msg.Type {
        case tea.KeyEsc:
                // Cancel: return to draft.
                c.state = protocol.ComposeDraft
                c.PreviewEnteredAt = time.Time{}
                return c, nil

        case tea.KeyEnter:
                // Send the message.
                c.sendAction()
                c.reset()
                return c, nil
        }

        // Handle "e" for edit.
        if msg.String() == "e" {
                c.state = protocol.ComposeDraft
                c.PreviewEnteredAt = time.Time{}
                return c, nil
        }

        // "q" → back to draft.
        if msg.String() == "q" {
                c.state = protocol.ComposeDraft
                c.PreviewEnteredAt = time.Time{}
                return c, nil
        }

        return c, nil
}

// ---------------------------------------------------------------------------
// ComposeTemplatePick key handling
// ---------------------------------------------------------------------------

func (c *Compose) handlePickKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.Type {
        case tea.KeyEsc:
                // Cancel: return to draft.
                c.state = protocol.ComposeDraft
                return c, nil

        case tea.KeyUp:
                if c.SnippetCursor > 0 {
                        c.SnippetCursor--
                }
                return c, nil

        case tea.KeyDown:
                if c.SnippetCursor < len(c.Snippets)-1 {
                        c.SnippetCursor++
                }
                return c, nil

        case tea.KeyEnter:
                // Use selected snippet.
                if c.SnippetCursor < len(c.Snippets) {
                        c.Draft = c.Snippets[c.SnippetCursor].Text
                        c.CursorPos = len(c.Draft)
                        c.CharCount = len(c.Draft)
                        c.transitionToPreview()
                }
                return c, nil
        }

        // "e" → insert snippet and open compose for editing.
        if msg.String() == "e" && c.SnippetCursor < len(c.Snippets) {
                c.Draft = c.Snippets[c.SnippetCursor].Text
                c.CursorPos = len(c.Draft)
                c.CharCount = len(c.Draft)
                c.state = protocol.ComposeDraft
                return c, nil
        }

        // "q" → back to draft.
        if msg.String() == "q" {
                c.state = protocol.ComposeDraft
                return c, nil
        }

        return c, nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

func (c *Compose) insertChar(ch string) {
        if c.MaxChars > 0 && c.CharCount >= c.MaxChars {
                return
        }
        before := c.Draft[:c.CursorPos]
        after := c.Draft[c.CursorPos:]
        c.Draft = before + ch + after
        c.CursorPos += len(ch) // byte offset — safe for ASCII; multi-byte runes increment by byte count
        c.CharCount = len(c.Draft)
}

func (c *Compose) backspace() {
        if c.CursorPos <= 0 {
                return
        }
        // Find the start of the last rune before the cursor to handle multi-byte correctly.
        runes := []rune(c.Draft[:c.CursorPos])
        if len(runes) == 0 {
                return
        }
        deleteLen := len(string(runes[len(runes)-1]))
        c.Draft = c.Draft[:c.CursorPos-deleteLen] + c.Draft[c.CursorPos:]
        c.CursorPos -= deleteLen
        c.CharCount = len(c.Draft)
}

func (c *Compose) transitionToPreview() {
        c.state = protocol.ComposePreview
        c.PreviewEnteredAt = time.Now()
        c.TypeOut = component.NewTemplatePreview(c.Draft)
        c.TypeOut.StartTypeOut()
}

func (c *Compose) sendAction() {
        if c.bus != nil {
                c.bus.Publish(bus.ActionMsg{
                        Action: string(protocol.ActionComposeSend),
                        Screen: protocol.ScreenCompose,
                        Params: map[string]any{
                                protocol.ParamTarget: c.Target,
                                protocol.ParamDraft:  c.Draft,
                        },
                })
        }
}

func (c *Compose) reset() {
        c.Draft = ""
        c.CursorPos = 0
        c.CharCount = 0
        c.ConsecutiveEnters = 0
        c.PreviewEnteredAt = time.Time{}
        c.state = protocol.ComposeDraft
}

// parseSnippets converts []any from backend params into []Snippet.
func parseSnippets(raw []any) []Snippet {
        var result []Snippet
        for _, item := range raw {
                if m, ok := item.(map[string]any); ok {
                        text, _ := m["text"].(string)
                        cat, _ := m["category"].(string)
                        if text != "" {
                                result = append(result, Snippet{Text: text, Category: cat})
                        }
                }
        }
        // Fallback to defaultSnippets (empty unless backend provides them).
        // The backend is the single source of truth; defaultSnippets is only
        // a safety net when the backend hasn't pushed data yet.
        if len(result) == 0 {
                return defaultSnippets
        }
        return result
}

// ---------------------------------------------------------------------------
// View rendering
// ---------------------------------------------------------------------------

// View renders the compose modal overlay.
func (c *Compose) View() string {
        switch c.state {
        case protocol.ComposeDraft:
                return c.viewDraft()
        case protocol.ComposePreview:
                return c.viewPreview()
        case protocol.ComposeTemplatePick:
                return c.viewPick()
        default:
                return c.viewDraft()
        }
}

// viewDraft renders the compose_draft state.
//
// P3 rule: No bordered box (┌──┐ │ │ └──┘). Hierarchy = brightness + size.
// The compose area uses BgRaised background + padding instead of box borders.
func (c *Compose) viewDraft() string {
        var b strings.Builder

        // Title line.
        title := i18n.T(i18n.KeyComposeDraftTitle)
        if c.Target != "" {
                title = fmt.Sprintf("%s: %s", title, c.Target)
        }
        b.WriteString(renderSectionTitle(title))
        b.WriteString("\n\n")

        // Hint.
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyComposeDraftHint)))
        b.WriteString("\n\n")

        // Text area — P3: raised background instead of bordered box.
        areaWidth := c.textAreaWidth()
        draftLines := strings.Split(c.Draft, "\n")
        if len(draftLines) == 0 {
                draftLines = []string{""}
        }

        var areaLines []string
        for i, line := range draftLines {
                if i == len(draftLines)-1 {
                        // Last line: show cursor per doc: _
                        cursorLine := style.BodyStyle.Render(line) + style.AccentStyle.Render("_")
                        areaLines = append(areaLines, cursorLine)
                } else {
                        areaLines = append(areaLines, style.BodyStyle.Render(line))
                }
        }

        // Add blank padding lines so area has minimum height.
        for len(areaLines) < composeMinBoxHeight {
                areaLines = append(areaLines, "")
        }

        areaContent := strings.Join(areaLines, "\n")
        innerStyle := lipgloss.NewStyle().
                Width(areaWidth).
                Padding(0, 1).
                Background(style.BgRaised)
        b.WriteString(innerStyle.Render(areaContent))

        // Char count in corner.
        var charCount string
        if c.MaxChars > 0 {
                charCount = style.DimStyle.Render(
                        fmt.Sprintf("%d/%d", c.CharCount, c.MaxChars),
                )
        } else {
                charCount = style.DimStyle.Render(
                        fmt.Sprintf("%d", c.CharCount),
                )
        }
        b.WriteString("  ")
        b.WriteString(charCount)
        b.WriteString("\n\n")

        // Actions.
        actions := []string{
                style.ActionStyle.Render(i18n.T(i18n.KeyComposeSend)),
                style.MutedStyle.Render(i18n.T(i18n.KeyComposeTab)),
                style.MutedStyle.Render(i18n.T(i18n.KeyComposeCancel)),
        }
        b.WriteString(strings.Join(actions, "    "))
        b.WriteString("\n\n")

        // Tip.
        b.WriteString(style.CaptionStyle.Render(i18n.T(i18n.KeyComposeTip)))

        return b.String()
}

// viewPreview renders the compose_preview state.
//
// Spec:
//
//      ── preview pesan ──
//      kirim ke: kopi nusantara
//      ────────────────────
//      (preview types out char-by-char)
//      ────────────────────
//      ↵ kirim    e edit lagi    esc batal
func (c *Compose) viewPreview() string {
        var b strings.Builder

        // Title.
        b.WriteString(renderSectionTitle(i18n.T(i18n.KeyComposePreviewTitle)))
        b.WriteString("\n\n")

        // Target.
        if c.Target != "" {
                b.WriteString(style.MutedStyle.Render(
                        fmt.Sprintf("%s: %s", i18n.T(i18n.KeyComposeSendTo), c.Target),
                ))
                b.WriteString("\n\n")
        }

        // Separator.
        writeSeparator(&b, c.Width)

        // Preview content (type-out animation).
        c.TypeOut.Tick(time.Now())
        b.WriteString(c.TypeOut.View())
        b.WriteString("\n\n")

        // Separator.
        writeSeparator(&b, c.Width)

        // Actions — preview uses single-enter, NOT double-enter.
        actions := []string{
                style.ActionStyle.Render(i18n.T(i18n.KeyComposeSendSingle)),
                style.MutedStyle.Render(i18n.T(i18n.KeyComposeEdit)),
                style.MutedStyle.Render(i18n.T(i18n.KeyComposeCancel)),
        }
        b.WriteString(strings.Join(actions, "    "))

        return b.String()
}

// viewPick renders the compose_template_pick state.
//
// Spec:
//
//      ── pilih snippet ──
//      snippet yang sering dipakai:
//      1  "boleh lihat dulu aja kak"           (soft pitch)
//      2  "aku kasih preview gratis ya"         (free offer)
//      ...
//      ↑↓ pilih    ↵ pake    e edit dulu    esc batal
//      snippet disimpan di: ~/.waclaw/snippets.md
func (c *Compose) viewPick() string {
        var b strings.Builder

        // Title.
        b.WriteString(renderSectionTitle(i18n.T(i18n.KeyComposePickTitle)))
        b.WriteString("\n\n")
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyComposeSnippetSubtitle)))
        b.WriteString("\n\n")

        // Snippet list.
        for i, snip := range c.Snippets {
                numStyle := style.CaptionStyle
                textStyle := style.BodyStyle
                catStyle := style.MutedStyle

                if i == c.SnippetCursor {
                        numStyle = style.AccentStyle
                        textStyle = style.SelectedBodyStyle
                        catStyle = style.AccentStyle
                }

                // Number.
                b.WriteString(numStyle.Render(fmt.Sprintf("%d", i+1)))
                b.WriteString("  ")

                // Snippet text — truncate if needed.
                text := snip.Text
                maxTextLen := c.textAreaWidth() - snippetTruncationMargin
                if maxTextLen < minSnippetTextWidth {
                        maxTextLen = minSnippetTextWidth
                }
                if len(text) > maxTextLen {
                        text = text[:maxTextLen-3] + "..."
                }
                b.WriteString(textStyle.Render(fmt.Sprintf("%q", text)))

                // Category label — i18n lookup at display time.
                if snip.Category != "" {
                        catLabel := snip.Category
                        if i18nKey, ok := snippetCategoryI18n[snip.Category]; ok {
                                catLabel = i18n.T(i18nKey)
                        }
                        b.WriteString("  ")
                        b.WriteString(catStyle.Render(fmt.Sprintf("(%s)", catLabel)))
                }
                b.WriteString("\n")
        }

        b.WriteString("\n")

        // Actions.
        actions := []string{
                style.MutedStyle.Render(i18n.T(i18n.KeyComposeSelect)),
                style.ActionStyle.Render(i18n.T(i18n.KeyComposeUse)),
                style.MutedStyle.Render(i18n.T(i18n.KeyComposeEditFirst)),
                style.MutedStyle.Render(i18n.T(i18n.KeyComposeCancel)),
        }
        b.WriteString(strings.Join(actions, "    "))
        b.WriteString("\n\n")

        // Snippet path hint.
        b.WriteString(style.CaptionStyle.Render(
                i18n.T(i18n.KeyComposeSnippetPath),
        ))

        return b.String()
}

// textAreaWidth returns a reasonable width for the compose text area.
func (c *Compose) textAreaWidth() int {
        w := c.Width
        if w <= 0 {
                w = defaultFallbackWidth
        }
        if w > maxTextAreaWidth {
                w = maxTextAreaWidth
        }
        return w - 4
}

// String returns a debug representation.
func (c *Compose) String() string {
        return fmt.Sprintf("Compose{state=%s, target=%s, chars=%d}", c.state, c.Target, c.CharCount)
}
