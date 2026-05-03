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
)

// ---------------------------------------------------------------------------
// Shared data — template placeholder names (DRY: delegates to component package)
// ---------------------------------------------------------------------------

// templatePlaceholderNames lists all supported template placeholder names.
// Delegates to component.TemplatePlaceholderNames() — single source of truth.
// Cached at init time to avoid repeated function calls in render loops.
var templatePlaceholderNames = component.TemplatePlaceholderNames()

// Layout constants for consistent spacing.
const (
        // previewMaxLen is the maximum character length for the short preview
        // text in the template list view before truncation.
        previewMaxLen = 30

        // previewTruncateLen is the number of characters to keep before
        // adding "..." when truncating the preview text.
        previewTruncateLen = 27
)

// defaultSubstitutionValues provides fallback preview values when the backend
// has not yet supplied substitution data. These match doc/06 exactly.
// Used by NewTemplateMgr to initialize previewValues when no backend data
// is available yet.
var defaultSubstitutionValues = map[string]string{
        "Title":    "kopi nusantara",
        "Category": "cafe",
        "Address":  "jl. hasanuddin 23, kediri",
        "City":     "kediri",
        "Rating":   "4.2",
        "Reviews":  "87",
        "Area":     "kediri",
}

// ---------------------------------------------------------------------------
// Data types — template model for display
// ---------------------------------------------------------------------------

// TemplateEntry holds display data for a single template variant.
type TemplateEntry struct {
        Name        string // e.g. "direct-curiosity"
        Type        string // "ice_breaker" or "offer"
        Preview     string // short preview of the template text
        Recommended bool   // whether this is the best-performing variant
        FilePath    string // full file path for edit hint
        Body        string // full template text (for preview)
        HasError    bool   // whether this template has validation errors
        Errors      []TemplateError
}

// TemplateError describes a single validation error in a template.
type TemplateError struct {
        LineNumber int    // 0 if not line-specific
        Message    string // human-readable error description
        Severity   string // "error" or "warning"
        Code       string // backend error code: "empty_file", "unknown_placeholder", "encoding_error"
}

// TemplateGroup groups templates by type (ice_breaker vs offer).
type TemplateGroup struct {
        Type      string
        Label     string
        Templates []TemplateEntry
}

// ---------------------------------------------------------------------------
// TemplateMgr screen
// ---------------------------------------------------------------------------

// TemplateMgr implements Screen 10: Template Manager → Armory.
//
// States: TemplateList, TemplatePreview, TemplateEditHint, TemplateValidationError.
//
// Visual spec from doc/06-screens-database-templates.md:
//   - Template list grouped by type (ice breaker / offer)
//   - Preview fills placeholders with sample lead data
//   - Edit redirects to file editor (no in-app editor)
//   - Validation errors with ✗/✓ indicators and severity
type TemplateMgr struct {
        tui.ScreenBase

        // state is the current screen state.
        state protocol.StateID

        // width and height track the terminal dimensions.
        width  int
        height int

        // groups holds the template groups populated by HandleNavigate/HandleUpdate.
        groups []TemplateGroup

        // cursor is the currently selected template index (flattened across groups).
        cursor int

        // currentTemplate is the template being previewed or edited.
        currentTemplate *TemplateEntry

        // preview is the component for rendering template with placeholders.
        preview component.TemplatePreview

        // previewValues holds substitution values for template preview, provided
        // by the backend via the preview_values key in HandleNavigate/HandleUpdate params.
        previewValues map[string]string

        // nicheName is the current niche context (e.g. "web_developer").
        nicheName string

        // showPlaceholders toggles the placeholder reference view.
        showPlaceholders bool

        // animStart tracks when the current state was entered.
        animStart time.Time

        // errorBlinkCount tracks how many blinks remain for validation error.
        errorBlinkCount int
        errorBlinkStart time.Time

        // reloadHighlightStart tracks when a reload occurred for brief highlight
        // on changed items — doc spec: "On reload: template list refreshes with
        // brief highlight on changed items."
        reloadHighlightStart time.Time

        // previousGroups holds the groups before the last update for diffing
        // on reload — used to detect changed templates for highlight.
        previousGroups []TemplateGroup
}

// NewTemplateMgr creates a new TemplateMgr screen with default state.
func NewTemplateMgr() *TemplateMgr {
        return &TemplateMgr{
                ScreenBase:    tui.NewScreenBase(protocol.ScreenTemplateMgr),
                state:         protocol.TemplateList,
                preview:       component.NewTemplatePreview(""),
                previewValues: make(map[string]string),
        }
}

// Focus is called when the screen becomes active.
func (s *TemplateMgr) Focus() {
        s.animStart = time.Now()
        if s.state == protocol.TemplateValidationError {
                s.errorBlinkStart = time.Now()
                s.errorBlinkCount = 2
        }
}

// Blur is called when the screen becomes inactive.
func (s *TemplateMgr) Blur() {}

// HandleNavigate processes a "navigate" command from the backend.
func (s *TemplateMgr) HandleNavigate(params map[string]any) error {
        s.animStart = time.Now()
        s.parseTemplateParams(params)
        return nil
}

// HandleUpdate processes an "update" command from the backend.
func (s *TemplateMgr) HandleUpdate(params map[string]any) error {
        if stateStr, ok := params[protocol.ParamState].(string); ok {
                s.state = protocol.StateID(stateStr)
                s.animStart = time.Now()
                if s.state == protocol.TemplateValidationError {
                        s.errorBlinkStart = time.Now()
                        s.errorBlinkCount = 2
                }
        }

        // Save previous groups before parsing new data for reload diff.
        s.previousGroups = s.groups
        s.parseTemplateParams(params)

        // Detect reload highlight: if groups data changed, mark highlight.
        if _, ok := params[protocol.ParamGroups].([]any); ok {
                s.reloadHighlightStart = time.Now()
        }

        return nil
}

// parseTemplateParams extracts shared data from HandleNavigate/HandleUpdate params.
// DRY: both handlers parse the same fields — this avoids duplicating the logic.
func (s *TemplateMgr) parseTemplateParams(params map[string]any) {
        if niche, ok := params[protocol.ParamNiche].(string); ok {
                s.nicheName = niche
        }
        if groupsData, ok := params[protocol.ParamGroups].([]any); ok {
                s.groups = parseTemplateGroups(groupsData)
        }
        if rawPreview, ok := params[protocol.ParamPreviewValues].(map[string]any); ok {
                s.previewValues = make(map[string]string, len(rawPreview))
                for k, v := range rawPreview {
                        if vs, ok := v.(string); ok {
                                s.previewValues[k] = vs
                        }
                }
        }
        // Also accept "substitution_values" key from backend (alternative key name).
        // This provides compatibility with both naming conventions.
        if rawSub, ok := params[protocol.ParamSubstitutionValues].(map[string]any); ok && len(rawSub) > 0 {
                s.previewValues = make(map[string]string, len(rawSub))
                for k, v := range rawSub {
                        if vs, ok := v.(string); ok {
                                s.previewValues[k] = vs
                        }
                }
        }
        if tmplData, ok := params[protocol.ParamTemplate].(map[string]any); ok {
                tmpl := parseTemplateEntry(tmplData)
                s.currentTemplate = &tmpl
                s.preview = component.NewTemplatePreview(tmpl.Body)
                setSampleVars(&s.preview, s.effectivePreviewValues())
        }
}

// Update handles bubbletea messages.
func (s *TemplateMgr) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (s *TemplateMgr) handleKey(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
        switch s.state {
        case protocol.TemplateList:
                return s.handleKeyList(msg, now)
        case protocol.TemplatePreview:
                return s.handleKeyPreview(msg, now)
        case protocol.TemplateEditHint:
                return s.handleKeyEditHint(msg, now)
        case protocol.TemplateValidationError:
                return s.handleKeyValidationError(msg, now)
        }
        return s, nil
}

// handleKeyList handles keys in the template_list state.
func (s *TemplateMgr) handleKeyList(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
        totalTemplates := s.totalTemplateCount()

        switch {
        case key.Matches(msg, tui.KeyUp):
                if s.cursor > 0 {
                        s.cursor--
                }
        case key.Matches(msg, tui.KeyDown):
                if s.cursor < totalTemplates-1 {
                        s.cursor++
                }
        case key.Matches(msg, tui.KeyEnter):
                tmpl := s.templateAtCursor()
                if tmpl != nil {
                        s.currentTemplate = tmpl
                        s.preview = component.NewTemplatePreview(tmpl.Body)
                        setSampleVars(&s.preview, s.effectivePreviewValues())
                        s.state = protocol.TemplatePreview
                        s.animStart = now
                }
        case key.Matches(msg, tui.KeyEdit):
                // "e" for edit — go to edit hint.
                tmpl := s.templateAtCursor()
                if tmpl != nil {
                        s.currentTemplate = tmpl
                        s.state = protocol.TemplateEditHint
                        s.animStart = now
                }
        case key.Matches(msg, tui.KeyNew):
                // "n" for new — publish action to backend.
                s.publishAction(protocol.ActionNewTemplate, map[string]any{
                        protocol.ParamNiche: s.nicheName,
                })
        }
        return s, nil
}

// handleKeyPreview handles keys in the template_preview state.
func (s *TemplateMgr) handleKeyPreview(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
        switch {
        case key.Matches(msg, tui.KeyBack):
                s.state = protocol.TemplateList
                s.animStart = now
                return s, nil
        case key.Matches(msg, tui.KeyEnter):
                // Use this template.
                s.publishAction(protocol.ActionUseTemplate, map[string]any{
                        protocol.ParamTemplateName: s.currentTemplate.Name,
                        protocol.ParamTemplateType: s.currentTemplate.Type,
                })
        case key.Matches(msg, tui.KeyEdit):
                // "e" for edit in file — go to edit hint.
                s.state = protocol.TemplateEditHint
                s.animStart = now
        }
        return s, nil
}

// handleKeyEditHint handles keys in the template_edit_hint state.
func (s *TemplateMgr) handleKeyEditHint(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
        switch {
        case key.Matches(msg, tui.KeyBack):
                s.state = protocol.TemplateList
                s.animStart = now
                return s, nil
        case key.Matches(msg, tui.KeyRefresh):
                // Reload template from file.
                s.publishAction(protocol.ActionReloadTemplate, map[string]any{
                        protocol.ParamTemplateName: s.currentTemplate.Name,
                        protocol.ParamTemplateType: s.currentTemplate.Type,
                        protocol.ParamNiche:       s.nicheName,
                })
        }
        return s, nil
}

// handleKeyValidationError handles keys in the template_validation_error state.
func (s *TemplateMgr) handleKeyValidationError(msg tea.KeyMsg, now time.Time) (tea.Model, tea.Cmd) {
        switch {
        case key.Matches(msg, tui.KeyBack):
                s.state = protocol.TemplateList
                s.animStart = now
                return s, nil
        case key.Matches(msg, tui.Key1):
                // Open file.
                if s.currentTemplate != nil {
                        s.publishAction(protocol.ActionOpenFile, map[string]any{
                                protocol.ParamFilePath: s.currentTemplate.FilePath,
                        })
                }
        case key.Matches(msg, tui.Key2):
                // Toggle placeholder reference.
                s.showPlaceholders = !s.showPlaceholders
        case key.Matches(msg, tui.KeyRefresh):
                // Reload.
                s.publishAction(protocol.ActionReloadTemplate, map[string]any{
                        protocol.ParamNiche: s.nicheName,
                })
        }
        return s, nil
}

// publishAction sends an action event through the bus.
func (s *TemplateMgr) publishAction(action string, params map[string]any) {
        if s.Bus() == nil {
                return
        }
        s.Bus().Publish(bus.ActionMsg{
                Action: action,
                Screen: protocol.ScreenTemplateMgr,
                Params: params,
        })
}

// View renders the current state of the TemplateMgr screen.
func (s *TemplateMgr) View() string {
        now := time.Now()

        switch s.state {
        case protocol.TemplateList:
                return s.viewList(now)
        case protocol.TemplatePreview:
                return s.viewPreview(now)
        case protocol.TemplateEditHint:
                return s.viewEditHint(now)
        case protocol.TemplateValidationError:
                return s.viewValidationError(now)
        default:
                return s.viewList(now)
        }
}

// viewList renders the template_list state.
//
// Visual spec from doc/06:
//
//      template pesan
//      niche: web_developer
//      ice breaker:
//      ▸ default              "halo kak, apakah ini..."
//      offer:
//      ▸ direct-curiosity     "tadi aku iseng cari..."    ★ recommended
//      ▸ pattern-interrupt    "permisi kak! aku haikal..."
//      ▸ admin-bypass         "halo kak admin!..."
//      ↑↓  pilih    ↵  preview    n  baru    e  edit    q  balik
func (s *TemplateMgr) viewList(now time.Time) string {
        var b strings.Builder

        // Title.
        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyTemplateTitle)))
        b.WriteString("\n")

        b.WriteString(style.Separator())
        b.WriteString("\n")

        // Niche label.
        b.WriteString(style.Indent(1))
        b.WriteString(style.CaptionStyle.Render(
                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyTemplateNiche), s.nicheName)))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SubSectionGap))

        // Template groups.
        flatIdx := 0
        for _, group := range s.groups {
                // Group label.
                groupLabel := group.Label
                if groupLabel == "" {
                        groupLabel = group.Type
                }
                b.WriteString(style.MutedStyle.Render(groupLabel + ":"))
                b.WriteString("\n")

                // Template entries within group.
                for _, tmpl := range group.Templates {
                        visibleAt := s.animStart.Add(time.Duration(flatIdx) * anim.MenuStagger)
                        if now.Before(visibleAt) {
                                flatIdx++
                                continue
                        }

                        prefix := "  "
                        if flatIdx == s.cursor {
                                prefix = style.AccentStyle.Render("▸ ")
                        } else {
                                prefix = style.DimStyle.Render("  ")
                        }

                        b.WriteString(prefix)

                        // Template name.
                        nameStyle := style.BodyStyle
                        if flatIdx == s.cursor {
                                nameStyle = style.PrimaryStyle
                        }
                        // Reload highlight: brief accent flash on changed templates.
                        if s.isTemplateChanged(tmpl) {
                                nameStyle = style.WarningStyle
                        }
                        b.WriteString(nameStyle.Render(tmpl.Name))

                        // Preview text.
                        if tmpl.Preview != "" {
                                previewText := tmpl.Preview
                                if len(previewText) > previewMaxLen {
                                        previewText = previewText[:previewTruncateLen] + "..."
                                }
                                b.WriteString("  ")
                                b.WriteString(style.CaptionStyle.Render(
                                        fmt.Sprintf("\"%s\"", previewText)))
                        }

                        // Recommended badge.
                        if tmpl.Recommended {
                                b.WriteString("    ")
                                b.WriteString(style.WarningStyle.Render(
                                        i18n.T(i18n.KeyTemplateRecommended)))
                        }

                        // Error indicator.
                        if tmpl.HasError {
                                b.WriteString("    ")
                                b.WriteString(style.DangerStyle.Render("✗"))
                        }

                        b.WriteString("\n")
                        flatIdx++
                }
        }

        b.WriteString(style.Separator())
        b.WriteString("\n")

        // Footer actions.
        parts := []string{
                i18n.T(i18n.KeyTemplateSelect),
                i18n.T(i18n.KeyTemplatePreviewAction),
                i18n.T(i18n.KeyTemplateNew),
                i18n.T(i18n.KeyTemplateEdit),
                i18n.T(i18n.KeyLabelBack),
        }
        b.WriteString(style.DimStyle.Render(strings.Join(parts, "    ")))

        return b.String()
}

// viewPreview renders the template_preview state.
//
// Visual spec from doc/06:
//
//      preview: direct-curiosity                    ★ recommended
//      Halo Kak {{.Title}}! 👋
//      Tadi aku iseng cari {{.Title}} di Google,
//      ...
//      placeholder:
//      {{.Title}} → kopi nusantara
//      {{.Category}} → cafe
//      {{.Address}} → jl. hasanuddin 23, kediri
//      ↵  pake ini    e  edit di file    q  balik
func (s *TemplateMgr) viewPreview(now time.Time) string {
        if s.currentTemplate == nil {
                return ""
        }
        tmpl := s.currentTemplate

        var b strings.Builder

        // Title with recommended badge.
        b.WriteString(style.MutedStyle.Render(
                fmt.Sprintf("%s: ", i18n.T(i18n.KeyTemplatePreviewLabel))))
        b.WriteString(style.HeadingStyle.Render(tmpl.Name))

        if tmpl.Recommended {
                b.WriteString("                    ")
                b.WriteString(style.WarningStyle.Render(
                        i18n.T(i18n.KeyTemplateRecommended)))
        }
        b.WriteString("\n")

        b.WriteString(style.Separator())
        b.WriteString("\n")

        // Template body with placeholder highlighting.
        b.WriteString(s.preview.ViewWithHighlight())
        b.WriteString("\n")

        b.WriteString(style.Separator())
        b.WriteString("\n")

        // Placeholder mapping — only show placeholders that appear in the template.
        b.WriteString(style.MutedStyle.Render(
                i18n.T(i18n.KeyTemplatePlaceholderLabel) + ":"))
        b.WriteString("\n")

        for _, ph := range templatePlaceholderNames {
                placeholder := fmt.Sprintf("{{.%s}}", ph)
                if !strings.Contains(tmpl.Body, placeholder) {
                        continue
                }
                b.WriteString(style.Indent(1))
                b.WriteString(style.AccentStyle.Render(placeholder))
                b.WriteString(" → ")
                if val, ok := s.previewValues[ph]; ok {
                        b.WriteString(style.BodyStyle.Render(val))
                }
                b.WriteString("\n")
        }

        b.WriteString(style.Section(style.SectionGap))

        // Footer actions.
        parts := []string{
                i18n.T(i18n.KeyTemplateUseThis),
                i18n.T(i18n.KeyTemplateEditFile),
                i18n.T(i18n.KeyLabelBack),
        }
        b.WriteString(style.DimStyle.Render(strings.Join(parts, "    ")))

        return b.String()
}

// viewEditHint renders the template_edit_hint state.
//
// Visual spec from doc/06:
//
//      edit template
//      template disimpan sebagai file teks.
//      buka di editor favorit lu:
//        ~/.waclaw/niches/web_developer/offer_1.md
//      setelah save, tekan r buat reload.
//      r  reload    q  balik
func (s *TemplateMgr) viewEditHint(now time.Time) string {
        if s.currentTemplate == nil {
                return ""
        }
        tmpl := s.currentTemplate

        var b strings.Builder

        b.WriteString(style.HeadingStyle.Render(i18n.T(i18n.KeyTemplateEditTitle)))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyTemplateSavedAsFile)))
        b.WriteString("\n")

        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyTemplateOpenInEditor)))
        b.WriteString("\n")

        // File path with accent highlight.
        b.WriteString(style.Indent(2))
        b.WriteString(style.AccentStyle.Render(tmpl.FilePath))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SubSectionGap))

        b.WriteString(style.Indent(1))
        b.WriteString(style.MutedStyle.Render(i18n.T(i18n.KeyTemplateAfterSave)))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SectionGap))

        // Footer actions.
        parts := []string{
                i18n.T(i18n.KeyTemplateReload),
                i18n.T(i18n.KeyLabelBack),
        }
        b.WriteString(style.DimStyle.Render(strings.Join(parts, "    ")))

        return b.String()
}

// viewValidationError renders the template_validation_error state.
//
// Visual spec from doc/06:
//
//      ✗ template error
//      niche: fotografer
//      ice_breaker.md
//      ✗  file kosong — isi dulu pesan pembukanya
//         tapi harus ada minimal {{.Title}} placeholder
//      offer_2.md
//      ✗  placeholder gak dikenali: {{.NamaToko}}
//         placeholder yang tersedia: {{.Title}}, {{.Category}}, ...
//      ✗  baris 5: encoding error (bukan UTF-8)
//         kemungkinan file disave pake encoding salah
//      worker fotografer di-pause sampai template diperbaiki.
//      ice_breaker WAJIB ada. offer yang error cuma di-skip.
//      1  buka file    2  liat placeholder yang tersedia
//      r  reload    q  balik
func (s *TemplateMgr) viewValidationError(now time.Time) string {
        var b strings.Builder

        // Error title — blinks red 2x per doc spec.
        blinkVisible := true
        if !s.errorBlinkStart.IsZero() {
                elapsed := now.Sub(s.errorBlinkStart)
                blinkCycle := anim.ConfigErrorBlink / 4 // 4 phases in 2 blinks
                phase := int(elapsed / blinkCycle) % 4
                blinkVisible = phase < 2 // visible in first 2 phases, hidden in last 2
        }
        if blinkVisible {
                b.WriteString(style.DangerStyle.Render(i18n.T(i18n.KeyTemplateErrorTitle)))
        } else {
                b.WriteString(style.DimStyle.Render(i18n.T(i18n.KeyTemplateErrorTitle)))
        }
        b.WriteString("\n")

        b.WriteString(style.Separator())
        b.WriteString("\n")

        // Niche label.
        b.WriteString(style.Indent(1))
        b.WriteString(style.CaptionStyle.Render(
                fmt.Sprintf("%s: %s", i18n.T(i18n.KeyTemplateNiche), s.nicheName)))
        b.WriteString("\n")

        b.WriteString(style.Section(style.SubSectionGap))

        // Error entries per template.
        for _, group := range s.groups {
                for _, tmpl := range group.Templates {
                        if !tmpl.HasError || len(tmpl.Errors) == 0 {
                                continue
                        }

                        // Template file name.
                        b.WriteString(style.Indent(1))
                        b.WriteString(style.BodyStyle.Render(tmpl.Name))
                        b.WriteString("\n")

                        for _, err := range tmpl.Errors {
                                // Error marker with blink effect.
                                b.WriteString(style.Indent(1))

                                // Differentiate ice_breaker errors (WAJIB) from offer errors (di-skip)
                                isIceBreaker := tmpl.Type == "ice_breaker"
                                if isIceBreaker {
                                        b.WriteString(style.DangerStyle.Render(i18n.T(i18n.KeyDataErrorMark) + "  "))
                                } else {
                                        b.WriteString(style.WarningStyle.Render(i18n.T(i18n.KeyDataErrorMark) + "  "))
                                }

                                if err.LineNumber > 0 {
                                        b.WriteString(style.DangerStyle.Render(
                                                fmt.Sprintf("%s %d: ", i18n.T(i18n.KeyTemplateLine), err.LineNumber)))
                                }

                                b.WriteString(style.DangerStyle.Render(err.Message))
                                b.WriteString("\n")

                                // Hint for specific error types.
                                if hint := errorHint(err); hint != "" {
                                        b.WriteString(style.Indent(2))
                                        b.WriteString(style.CaptionStyle.Render(hint))
                                        b.WriteString("\n")
                                }
                        }
                }
        }

        b.WriteString(style.Separator())
        b.WriteString("\n")

        // Worker paused message.
        b.WriteString(style.Indent(1))
        b.WriteString(style.WarningStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyTemplateWorkerPausedFmt), s.nicheName)))
        b.WriteString("\n")

        b.WriteString(style.Indent(1))
        b.WriteString(style.CaptionStyle.Render(
                i18n.T(i18n.KeyTemplateIceBreakerReq)))
        b.WriteString("\n")

        // Optional: placeholder reference view.
        if s.showPlaceholders {
                b.WriteString(style.Section(style.SubSectionGap))
                b.WriteString(style.Indent(1))
                b.WriteString(style.MutedStyle.Render(
                        i18n.T(i18n.KeyTemplateAvailablePlaceholders) + ":"))
                b.WriteString("\n")

                for _, ph := range templatePlaceholderNames {
                        b.WriteString(style.Indent(2))
                        b.WriteString(style.AccentStyle.Render(
                                fmt.Sprintf("{{.%s}}", ph)))
                        b.WriteString("\n")
                }
        }

        b.WriteString(style.Section(style.SectionGap))

        // Footer actions.
        parts := []string{
                i18n.T(i18n.KeyTemplateOpenFile),
                i18n.T(i18n.KeyTemplateViewPlaceholders),
                i18n.T(i18n.KeyTemplateReload),
                i18n.T(i18n.KeyLabelBack),
        }
        b.WriteString(style.DimStyle.Render(strings.Join(parts, "    ")))

        return b.String()
}

// ---------------------------------------------------------------------------
// Helper functions — DRY shared logic
// ---------------------------------------------------------------------------

// isTemplateChanged checks if a template name existed in previousGroups
// with different content (preview, recommended, or error status).
func (s *TemplateMgr) isTemplateChanged(tmpl TemplateEntry) bool {
        if s.previousGroups == nil || s.reloadHighlightStart.IsZero() {
                return false
        }
        // Only highlight for 1.5 seconds after reload.
        if time.Since(s.reloadHighlightStart) > 1500*time.Millisecond {
                return false
        }
        for _, pg := range s.previousGroups {
                for _, pt := range pg.Templates {
                        if pt.Name == tmpl.Name && pt.Type == tmpl.Type {
                                return pt.Preview != tmpl.Preview ||
                                        pt.Recommended != tmpl.Recommended ||
                                        pt.HasError != tmpl.HasError
                        }
                }
        }
        // New template not in previous groups = changed (newly added).
        return true
}

// totalTemplateCount returns the total number of templates across all groups.
func (s *TemplateMgr) totalTemplateCount() int {
        count := 0
        for _, group := range s.groups {
                count += len(group.Templates)
        }
        return count
}

// templateAtCursor returns the template at the current cursor position,
// flattening across all groups.
func (s *TemplateMgr) templateAtCursor() *TemplateEntry {
        idx := 0
        for i := range s.groups {
                for j := range s.groups[i].Templates {
                        if idx == s.cursor {
                                return &s.groups[i].Templates[j]
                        }
                        idx++
                }
        }
        return nil
}

// effectivePreviewValues returns backend-provided preview values,
// falling back to defaultSubstitutionValues when the backend has not
// yet supplied any.
func (s *TemplateMgr) effectivePreviewValues() map[string]string {
        if len(s.previewValues) > 0 {
                return s.previewValues
        }
        return defaultSubstitutionValues
}

// setSampleVars populates the TemplatePreview with the given substitution values.
// The values are provided by the backend via the preview_values key in
// HandleNavigate/HandleUpdate params.
func setSampleVars(preview *component.TemplatePreview, values map[string]string) {
        for k, v := range values {
                preview.SetVar(k, v)
        }
}

// errorHint returns a contextual hint based on the error's structured code.
// Uses the backend-provided code field instead of string matching — proper
// separation of concerns between backend error codes and TUI display logic.
func errorHint(err TemplateError) string {
        code := err.Code
        if code == "" {
                code = err.Severity // Fallback for older backend that only sends severity
        }
        switch code {
        case "empty_file":
                return i18n.T(i18n.KeyTemplateMinPlaceholder)
        case "unknown_placeholder":
                var names []string
                for i, ph := range templatePlaceholderNames {
                        if i < 5 { // Show first 5 placeholders in hint
                                names = append(names, "{{."+ph+"}}")
                        }
                }
                return i18n.T(i18n.KeyTemplateAvailablePlaceholders) + ": " + strings.Join(names, ", ")
        case "encoding_error":
                return i18n.T(i18n.KeyTemplateEncodingHint)
        }
        return ""
}

// parseTemplateGroups converts raw params into TemplateGroup slices.
func parseTemplateGroups(raw []any) []TemplateGroup {
        var groups []TemplateGroup
        for _, item := range raw {
                if m, ok := item.(map[string]any); ok {
                        group := TemplateGroup{}
                        if typ, ok := m[protocol.ParamType].(string); ok {
                                group.Type = typ
                        }
                        if label, ok := m[protocol.ParamLabel].(string); ok {
                                group.Label = label
                        }
                        if templates, ok := m[protocol.ParamTemplates].([]any); ok {
                                for _, t := range templates {
                                        if tm, ok := t.(map[string]any); ok {
                                                group.Templates = append(group.Templates, parseTemplateEntry(tm))
                                        }
                                }
                        }
                        groups = append(groups, group)
                }
        }
        return groups
}

// parseTemplateEntry converts a single map into a TemplateEntry.
func parseTemplateEntry(m map[string]any) TemplateEntry {
        entry := TemplateEntry{}
        if v, ok := m[protocol.ParamName].(string); ok {
                entry.Name = v
        }
        if v, ok := m[protocol.ParamType].(string); ok {
                entry.Type = v
        }
        if v, ok := m[protocol.ParamPreview].(string); ok {
                entry.Preview = v
        }
        if v, ok := m[protocol.ParamRecommended].(bool); ok {
                entry.Recommended = v
        }
        if v, ok := m[protocol.ParamFilePath].(string); ok {
                entry.FilePath = v
        }
        if v, ok := m[protocol.ParamBody].(string); ok {
                entry.Body = v
        }
        if v, ok := m[protocol.ParamHasError].(bool); ok {
                entry.HasError = v
        }
        if errors, ok := m[protocol.ParamErrors].([]any); ok {
                for _, e := range errors {
                        if em, ok := e.(map[string]any); ok {
                                err := TemplateError{}
                                if ln, ok := em[protocol.ParamLine].(float64); ok {
                                        err.LineNumber = int(ln)
                                }
                                if msg, ok := em[protocol.ParamMessage].(string); ok {
                                        err.Message = msg
                                }
                                if sev, ok := em[protocol.ParamSeverity].(string); ok {
                                        err.Severity = sev
                                }
                                if code, ok := em[protocol.ParamCode].(string); ok {
                                        err.Code = code // Error code for display routing
                                }
                                entry.Errors = append(entry.Errors, err)
                        }
                }
        }
        return entry
}
