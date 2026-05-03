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

// ---------------------------------------------------------------------------
// Data types for the Guardrail (Config Validation) screen
// ---------------------------------------------------------------------------

// ValidationResult holds the validation status of a single config file or niche.
type ValidationResult struct {
        Name    string
        Status  string // "ok", "error", "warning", "checking", "waiting"
        Details string // "3 template · 5 area"
        Errors  []ValidationError
        Warnings []ValidationWarning
        // FIX G-BF01: TemplateCount and AreaCount replace fragile parse logic
        TemplateCount int
        AreaCount     int
}

// ValidationError describes a single config parsing error.
type ValidationError struct {
        Line    int
        Message string
        Pointer string   // the ^^^^ underline
        Context []string // surrounding lines for gutter display
}

// ValidationWarning describes a non-blocking config issue.
// FIX G-05: add Context and Pointer fields for gutter rendering
type ValidationWarning struct {
        Line       int
        Message    string
        Suggestion string
        Context    []string // FIX G-05: context lines for gutter display
        Pointer    string   // FIX G-05: pointer underline
}

// ---------------------------------------------------------------------------
// Guardrail screen model
// ---------------------------------------------------------------------------

// Guardrail implements tui.Screen for Screen 14: Config Validation.
// It validates all config files on boot and whenever the user triggers a
// check, displaying errors with precise line/pointer diagnostics and
// warnings with migration guidance.
type Guardrail struct {
        tui.ScreenBase
        state         protocol.StateID
        results       []ValidationResult
        pausedNiches  []string
        selectedError int
        list          component.ListSelect
        width         int
        height        int
        focused       bool
        totalErrors   int
        totalWarnings int
        // FIX G-06/G-07: check progress stored in struct instead of hardcoded
        checkProgress int
}

// NewGuardrail creates a Guardrail screen with demo data.
// FIX 3: adds demo error and warning data
// FIX G-BF01: sets TemplateCount and AreaCount directly in demo data
func NewGuardrail() *Guardrail {
        g := &Guardrail{
                ScreenBase: tui.NewScreenBase(protocol.ScreenGuardrail),
                state:      protocol.ValidationClean,
                results: []ValidationResult{
                        {Name: "config.yaml", Status: "ok", TemplateCount: 0, AreaCount: 0},
                        {Name: "theme.yaml", Status: "ok", TemplateCount: 0, AreaCount: 0},
                        {Name: "queries.md", Status: "ok", TemplateCount: 0, AreaCount: 0},
                },
                // FIX G-06/G-07: check progress starts at zero — backend provides actual value
                checkProgress: 0,
                list: component.NewListSelect([]component.ListItem{
                        {Label: i18n.T(i18n.KeyGuardrailAutoGenerate), Description: ""},
                        {Label: i18n.T(i18n.KeyGuardrailSeeExample), Description: ""},
                }),
        }
        return g
}

// ---------------------------------------------------------------------------
// tea.Model interface
// ---------------------------------------------------------------------------

func (g *Guardrail) Init() tea.Cmd { return nil }

func (g *Guardrail) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch m := msg.(type) {
        case tea.WindowSizeMsg:
                g.width = m.Width
                g.height = m.Height
                return g, nil
        }

        if key, ok := msg.(tea.KeyMsg); ok {
                switch key.String() {
                case "q":
                        return g, nil
                case "1":
                        // Open first error file (or auto-generate for first_time)
                        if g.state == protocol.ValidationErrors {
                                g.selectedError++
                                if g.selectedError >= g.totalErrors {
                                        g.selectedError = 0
                                }
                                if g.Bus() != nil {
                                        g.Bus().Publish(bus.ActionMsg{Action: protocol.ActionOpenErrorFile, Screen: g.ID()})
                                }
                        } else if g.state == protocol.ValidationFirstTime {
                                // Auto-generate config
                                g.state = protocol.ValidationClean
                        }
                case "2":
                        if g.state == protocol.ValidationErrors {
                                // See example
                                if g.Bus() != nil {
                                        g.Bus().Publish(bus.ActionMsg{Action: protocol.ActionShowExample, Screen: g.ID()})
                                }
                        } else if g.state == protocol.ValidationReloadError {
                                // Revert backup
                                if g.Bus() != nil {
                                        g.Bus().Publish(bus.ActionMsg{Action: protocol.ActionRevertBackup, Screen: g.ID()})
                                }
                        }
                case "r":
                        // Reload after fix
                        if g.state == protocol.ValidationErrors || g.state == protocol.ValidationWarnings {
                                g.state = protocol.ValidationFix
                                // Simulate re-validation completing
                                g.results = []ValidationResult{
                                        {Name: "config.yaml", Status: "ok", TemplateCount: 0, AreaCount: 0},
                                        {Name: "theme.yaml", Status: "ok", TemplateCount: 0, AreaCount: 0},
                                        {Name: "queries.md", Status: "ok", Details: "← fixed!", TemplateCount: 0, AreaCount: 0},
                                        {Name: "web_developer", Status: "ok", Details: "3 template · 5 area", TemplateCount: 3, AreaCount: 5},
                                        {Name: "undangan_digital", Status: "ok", Details: "2 template · 8 area  ← fixed!", TemplateCount: 2, AreaCount: 8},
                                        {Name: "social_media_mgr", Status: "ok", Details: "1 template · 3 area", TemplateCount: 1, AreaCount: 3},
                                }
                        }
                case "enter":
                        switch g.state {
                        case protocol.ValidationClean:
                                g.state = protocol.ValidationClean // stay
                        case protocol.ValidationErrors:
                                g.state = protocol.ValidationErrors
                        case protocol.ValidationWarnings:
                                g.state = protocol.ValidationWarnings
                        }
                case "up", "k":
                        if g.state == protocol.ValidationFirstTime {
                                g.list.Up()
                        }
                case "down", "j":
                        if g.state == protocol.ValidationFirstTime {
                                g.list.Down()
                        }
                }
        }
        return g, nil
}

func (g *Guardrail) View() string {
        switch g.state {
        case protocol.ValidationClean:
                return g.viewClean()
        case protocol.ValidationErrors:
                return g.viewErrors()
        case protocol.ValidationWarnings:
                return g.viewWarnings()
        case protocol.ValidationFix:
                return g.viewFix()
        case protocol.ValidationFirstTime:
                return g.viewFirstTime()
        default:
                return g.viewClean()
        }
}

// ---------------------------------------------------------------------------
// Screen interface
// ---------------------------------------------------------------------------

func (g *Guardrail) HandleNavigate(params map[string]any) error {
        applyNavigateState(&g.state, params)
        return nil
}

func (g *Guardrail) HandleUpdate(params map[string]any) error {
        if st, ok := params[protocol.ParamState].(string); ok {
                g.state = protocol.StateID(st)
        }
        // Backend sends generic map data — convert to internal TUI types.
        // Do NOT assert TUI types from backend params (frontend/backend concern split).
        if raw, ok := params[protocol.ParamValidationResults]; ok {
                if list, ok := raw.([]map[string]any); ok {
                        var results []ValidationResult
                        for _, m := range list {
                                vr := ValidationResult{}
                                if v, ok := m[protocol.ParamName].(string); ok {
                                        vr.Name = v
                                }
                                if v, ok := m[protocol.ParamStatus].(string); ok {
                                        vr.Status = v
                                }
                                if v, ok := m[protocol.ParamDetails].(string); ok {
                                        vr.Details = v
                                }
                                // Parse errors array within each result
                                if rawErrs, ok := m[protocol.ParamErrors].([]any); ok {
                                        for _, re := range rawErrs {
                                                if em, ok := re.(map[string]any); ok {
                                                        ve := ValidationError{}
                                                        if v, ok := em[protocol.ParamLine].(int); ok {
                                                                ve.Line = v
                                                        }
                                                        if v, ok := em[protocol.ParamMessage].(string); ok {
                                                                ve.Message = v
                                                        }
                                                        if v, ok := em[protocol.ParamPointer].(string); ok {
                                                                ve.Pointer = v
                                                        }
                                                        if rawCtx, ok := em[protocol.ParamContext].([]any); ok {
                                                                for _, c := range rawCtx {
                                                                        if cs, ok := c.(string); ok {
                                                                                ve.Context = append(ve.Context, cs)
                                                                        }
                                                                }
                                                        }
                                                        vr.Errors = append(vr.Errors, ve)
                                                }
                                        }
                                }
                                // Parse warnings array within each result
                                if rawWarns, ok := m[protocol.ParamWarnings].([]any); ok {
                                        for _, rw := range rawWarns {
                                                if wm, ok := rw.(map[string]any); ok {
                                                        vw := ValidationWarning{}
                                                        if v, ok := wm[protocol.ParamLine].(int); ok {
                                                                vw.Line = v
                                                        }
                                                        if v, ok := wm[protocol.ParamMessage].(string); ok {
                                                                vw.Message = v
                                                        }
                                                        if v, ok := wm[protocol.ParamSuggestion].(string); ok {
                                                                vw.Suggestion = v
                                                        }
                                                        if rawCtx, ok := wm[protocol.ParamContext].([]any); ok {
                                                                for _, c := range rawCtx {
                                                                        if cs, ok := c.(string); ok {
                                                                                vw.Context = append(vw.Context, cs)
                                                                        }
                                                                }
                                                        }
                                                        if v, ok := wm[protocol.ParamPointer].(string); ok {
                                                                vw.Pointer = v
                                                        }
                                                        vr.Warnings = append(vr.Warnings, vw)
                                                }
                                        }
                                }
                                results = append(results, vr)
                        }
                        if len(results) > 0 {
                                g.results = results
                        }
                }
        }
        // FIX G-06/G-07: populate checkProgress from backend
        if v, ok := params[protocol.ParamCheckProgress].(int); ok {
                g.checkProgress = v
        }
        return nil
}

func (g *Guardrail) Focus() { g.focused = true }
func (g *Guardrail) Blur()  { g.focused = false }

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// FIX G-DRY02: replaced local statusIcon() with shared renderStatusIcon() from helpers.go

// renderResultRow renders a single validation result line.
// FIX G-04: shows error/warning counts per doc: "✗ 2 error", "⚠️ 1 warning"
func (g *Guardrail) renderResultRow(r ValidationResult, indent int) string {
        icon, color := renderStatusIcon(r.Status)
        nameStyle := lipgloss.NewStyle().Foreground(color)

        var b strings.Builder
        b.WriteString(style.Indent(indent))
        b.WriteString(nameStyle.Render(icon))
        b.WriteString("  ")
        b.WriteString(nameStyle.Render(r.Name))

        if r.Details != "" {
                b.WriteString("    ")
                // Check for "← fixed!" marker
                if strings.Contains(r.Details, "← fixed!") {
                        b.WriteString(style.SuccessStyle.Bold(true).Render(r.Details))
                } else {
                        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(r.Details))
                }
        }

        // FIX G-04: show error/warning counts per doc — i18n-aware pluralization
        if len(r.Errors) > 0 {
                b.WriteString("    ")
                errCount := len(r.Errors)
                errText := i18n.T(i18n.KeyGuardrailErrorCountOne)
                if errCount > 1 {
                        errText = i18n.T(i18n.KeyGuardrailErrorCountMany)
                }
                b.WriteString(style.DangerStyle.Render(fmt.Sprintf("✗ %d %s", errCount, errText)))
        }
        if len(r.Warnings) > 0 {
                b.WriteString("    ")
                warnCount := len(r.Warnings)
                warnText := i18n.T(i18n.KeyGuardrailWarningCountOne)
                if warnCount > 1 {
                        warnText = i18n.T(i18n.KeyGuardrailWarningCountMany)
                }
                b.WriteString(style.WarningStyle.Render(fmt.Sprintf("⚠ %d %s", warnCount, warnText)))
        }

        return b.String()
}

// renderErrorDetail renders an error with gutter and pointer.
// FIX G-DRY01: uses shared renderGutterError() from helpers.go
func (g *Guardrail) renderErrorDetail(err ValidationError, fileLabel string) string {
        var b strings.Builder

        // File section header — P3 compliance: use style.SectionLabel() instead of ── box-drawing chars
        if fileLabel != "" {
                b.WriteString(style.SectionLabel(fileLabel))
        }

        // FIX G-DRY01: use shared renderGutterError() instead of inline rendering
        linePrefix := fmt.Sprintf("%s %d: %s", i18n.T(i18n.KeyLabelLine), err.Line, err.Message)
        renderGutterError(&b, linePrefix, err.Context, err.Pointer, true)

        return b.String()
}

// renderWarningDetail renders a warning with suggestion.
// FIX G-05: renders Context and Pointer using renderGutterError()
func (g *Guardrail) renderWarningDetail(warn ValidationWarning, fileLabel string) string {
        var b strings.Builder

        // File section header — P3 compliance: use style.SectionLabel() instead of ── box-drawing chars
        if fileLabel != "" {
                b.WriteString(style.SectionLabel(fileLabel))
        }

        if len(warn.Context) > 0 || warn.Pointer != "" {
                // FIX G-05: render with gutter pattern using shared renderGutterError()
                linePrefix := fmt.Sprintf("%s %d: %s", i18n.T(i18n.KeyLabelLine), warn.Line, warn.Message)
                renderGutterError(&b, linePrefix, warn.Context, warn.Pointer, false)
        } else {
                // Simple warning without gutter
                b.WriteString(style.Indent(1))
                b.WriteString(style.WarningStyle.Render(
                        fmt.Sprintf("⚠  %s %d: %s", i18n.T(i18n.KeyLabelLine), warn.Line, warn.Message),
                ))
                b.WriteString("\n")
        }

        if warn.Suggestion != "" {
                b.WriteString(style.Indent(2))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render("│  "))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(warn.Suggestion))
                b.WriteString("\n")
        }

        return b.String()
}

// ---------------------------------------------------------------------------
// Views per state
// ---------------------------------------------------------------------------

// FIX 1: viewClean() replaces hardcoded strings with i18n keys
// FIX CC-05: uses renderHeadWithStatus() instead of statusAlignGap
func (g *Guardrail) viewClean() string {
        var b strings.Builder

        // Heading — FIX CC-05: use renderHeadWithStatus
        b.WriteString(renderHeadWithStatus(i18n.T(i18n.KeyGuardrailTitle), i18n.T(i18n.KeyGuardrailCleanEmoji), style.Success))
        b.WriteString(style.Section(style.SectionGap))

        // File results
        for _, r := range g.results[:3] { // config files first
                b.WriteString(g.renderResultRow(r, 0))
                b.WriteString("\n")
        }

        // Niche results
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyGuardrailNiches)))
        b.WriteString(style.Section(style.SubSectionGap))

        for _, r := range g.results[3:] {
                b.WriteString(g.renderResultRow(r, 0))
                b.WriteString("\n")
        }

        // Success message — FIX G-BF01: use TemplateCount/AreaCount directly instead of fragile parsing
        nicheCount := len(g.results) - 3
        templateCount := 0
        for _, r := range g.results[3:] {
                templateCount += r.TemplateCount
        }

        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(style.BodyStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyGuardrailNichesReady), nicheCount),
        ))
        b.WriteString("\n")
        b.WriteString(style.BodyStyle.Render(
                fmt.Sprintf("%d %s", templateCount, i18n.T(i18n.KeyGuardrailTemplatesLoaded)),
        ))
        b.WriteString("\n")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyGuardrailWorkersOK)))
        b.WriteString("\n")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyGuardrailCanGas)))

        // Key hints
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf(" ↵  %s    q  %s", i18n.T(i18n.KeyGuardrailGo), i18n.T(i18n.KeyLabelBack),
                )))

        return b.String()
}

// FIX 4: viewErrors() adds "niches:" subheading, paused niches count
// FIX G-01/G-02: iterates over g.results for actual error data instead of demoValidationErrors()
// FIX G-08: adds reassurance message for still-running niches
// FIX CC-05: uses renderHeadWithStatus()
func (g *Guardrail) viewErrors() string {
        var b strings.Builder

        // Heading — FIX CC-05: use renderHeadWithStatus
        b.WriteString(renderHeadWithStatus(i18n.T(i18n.KeyGuardrailTitle), i18n.T(i18n.KeyGuardrailErrorEmoji), style.Danger))
        b.WriteString(style.Section(style.SectionGap))

        // File/niche results with mixed statuses — split config files from niches with subheading
        pausedNiches := []string{}
        runningNiches := []string{} // FIX G-08: track still-running niches
        totalErrors := 0

        // Config file results first
        for _, r := range g.results {
                if !isConfigFile(r.Name) {
                        continue
                }
                b.WriteString(g.renderResultRow(r, 0))
                b.WriteString("\n")

                if r.Status == "error" {
                        totalErrors += len(r.Errors)
                }
        }

        // Niches subheading
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyGuardrailNiches)))
        b.WriteString(style.Section(style.SubSectionGap))

        for _, r := range g.results {
                if isConfigFile(r.Name) {
                        continue
                }
                b.WriteString(g.renderResultRow(r, 0))
                b.WriteString("\n")

                if r.Status == "error" {
                        totalErrors += len(r.Errors)
                        pausedNiches = append(pausedNiches, r.Name)
                } else if r.Status == "ok" || r.Status == "warning" {
                        // FIX G-08: niches that are still running
                        runningNiches = append(runningNiches, r.Name)
                }
        }

        // Error count summary
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(style.DangerStyle.Render(
                fmt.Sprintf(i18n.T(i18n.KeyGuardrailErrorCountFmt), totalErrors),
        ))

        // Paused niches with "niches:" subheading and count
        if len(pausedNiches) > 0 {
                b.WriteString("\n")
                b.WriteString(style.BodyStyle.Render(
                        fmt.Sprintf(i18n.T(i18n.KeyGuardrailPausedNichesFmt), len(pausedNiches)),
                ))
                b.WriteString(" ")
                b.WriteString(style.DangerStyle.Render(strings.Join(pausedNiches, ", ")))
        }

        // FIX G-08: reassurance message for still-running niches
        if len(runningNiches) > 0 {
                b.WriteString("\n")
                for _, name := range runningNiches {
                        b.WriteString(style.BodyStyle.Render(
                                fmt.Sprintf("%s %s", name, i18n.T(i18n.KeyGuardrailStillRunningMsg)),
                        ))
                        b.WriteString("\n")
                }
        }

        // Error details section — FIX G-01/G-02: iterate over g.results instead of demoValidationErrors()
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyGuardrailDetailErrors)))
        b.WriteString(style.Section(style.SubSectionGap))

        first := true
        for _, r := range g.results {
                if len(r.Errors) == 0 {
                        continue
                }
                for _, err := range r.Errors {
                        if !first {
                                b.WriteString(style.Section(style.SubSectionGap))
                        }
                        first = false
                        b.WriteString(g.renderErrorDetail(err, r.Name))
                }
        }

        // Key hints
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("1  %s    2  %s    r  %s    q  %s",
                        i18n.T(i18n.KeyGuardrailOpenFirst), i18n.T(i18n.KeyGuardrailSeeExample),
                        i18n.T(i18n.KeyLabelReload), i18n.T(i18n.KeyLabelBack),
                )))

        // Press-1-again hint
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                i18n.T(i18n.KeyGuardrailPress1Next),
        ))

        return b.String()
}

// isConfigFile returns true if the name looks like a config file (has extension).
func isConfigFile(name string) bool {
        return strings.Contains(name, ".")
}

// FIX 6: viewWarnings() adds detail warning section
// FIX G-03: iterates over g.results instead of demoValidationWarnings()
// FIX CC-05: uses renderHeadWithStatus()
func (g *Guardrail) viewWarnings() string {
        var b strings.Builder

        // Heading — FIX CC-05: use renderHeadWithStatus
        b.WriteString(renderHeadWithStatus(i18n.T(i18n.KeyGuardrailTitle), i18n.T(i18n.KeyGuardrailWarningEmoji), style.Warning))
        b.WriteString(style.Section(style.SectionGap))

        // Results with mixed statuses — split config files from niches with subheading
        for _, r := range g.results {
                if !isConfigFile(r.Name) {
                        continue
                }
                b.WriteString(g.renderResultRow(r, 0))
                b.WriteString("\n")
        }

        // Niches subheading
        b.WriteString(style.Section(style.SubSectionGap))
        b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyGuardrailNiches)))
        b.WriteString(style.Section(style.SubSectionGap))

        for _, r := range g.results {
                if isConfigFile(r.Name) {
                        continue
                }
                b.WriteString(g.renderResultRow(r, 0))
                b.WriteString("\n")
        }

        // Not blocking message
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyGuardrailNotBlocking)))

        // FIX G-03: detail warning section iterating over g.results
        hasWarnings := false
        for _, r := range g.results {
                if len(r.Warnings) > 0 {
                        hasWarnings = true
                        break
                }
        }
        if hasWarnings {
                b.WriteString(style.Section(style.SectionGap))
                b.WriteString(style.SubHeadingStyle.Render(i18n.T(i18n.KeyGuardrailDetailWarning)))
                b.WriteString(style.Section(style.SubSectionGap))

                first := true
                for _, r := range g.results {
                        if len(r.Warnings) == 0 {
                                continue
                        }
                        for _, w := range r.Warnings {
                                if !first {
                                        b.WriteString(style.Section(style.SubSectionGap))
                                }
                                first = false
                                b.WriteString(g.renderWarningDetail(w, r.Name))
                        }
                }
        }

        // Key hints
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("↵  %s    1  %s    r  %s    q  %s",
                        i18n.T(i18n.KeyGuardrailContinueAnyway), i18n.T(i18n.KeyGuardrailOpenFile),
                        i18n.T(i18n.KeyLabelReload), i18n.T(i18n.KeyLabelBack),
                )))

        return b.String()
}

// FIX 5: viewFix() shows sequential checking/waiting states instead of all ok immediately
// FIX G-06/G-07: uses g.checkProgress instead of hardcoded fixedCount
// FIX CC-05: uses renderHeadWithStatus()
func (g *Guardrail) viewFix() string {
        var b strings.Builder

        // Heading — FIX CC-05: use renderHeadWithStatus
        b.WriteString(renderHeadWithStatus(i18n.T(i18n.KeyGuardrailTitle), i18n.T(i18n.KeyGuardrailRevalidating), style.Text))
        b.WriteString(style.Section(style.SectionGap))

        // Show results with checking/waiting/done states derived from g.results
        for i, r := range g.results {
                if i < g.checkProgress {
                        b.WriteString(g.renderResultRow(ValidationResult{
                                Name:          r.Name,
                                Status:        "ok",
                                Details:       r.Details,
                                TemplateCount: r.TemplateCount,
                                AreaCount:     r.AreaCount,
                        }, 0))
                } else if i == g.checkProgress {
                        b.WriteString(g.renderResultRow(ValidationResult{
                                Name:   r.Name,
                                Status: "checking",
                        }, 0))
                } else {
                        b.WriteString(g.renderResultRow(ValidationResult{
                                Name:   r.Name,
                                Status: "waiting",
                        }, 0))
                }
                b.WriteString("\n")
        }

        // Progress message — FIX G-06/G-07: use g.checkProgress instead of hardcoded 4
        totalCount := len(g.results)
        fixedCount := g.checkProgress
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyGuardrailCheckingFile)))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(
                fmt.Sprintf("%d %s", fixedCount, i18n.T(i18n.KeyGuardrailOfDone)),
        ))
        b.WriteString("\n")
        b.WriteString(style.MutedStyle.Render(
                fmt.Sprintf("%d %s", totalCount-fixedCount, i18n.T(i18n.KeyGuardrailWaitingLabel)),
        ))

        return b.String()
}

// FIX 2: viewFirstTime() replaces hardcoded strings with i18n keys
// FIX CC-05: uses renderHeadWithStatus()
func (g *Guardrail) viewFirstTime() string {
        var b strings.Builder

        // Heading — FIX CC-05: use renderHeadWithStatus
        b.WriteString(renderHeadWithStatus(i18n.T(i18n.KeyGuardrailTitle), i18n.T(i18n.KeyGuardrailFirstTime), style.TextMuted))
        b.WriteString(style.Section(style.SectionGap))

        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyGuardrailNoConfig)))
        b.WriteString(style.Section(style.SectionGap))

        // FIX 2: use i18n keys
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyGuardrailNeedPrepare)))
        b.WriteString(style.Section(style.SubSectionGap))

        b.WriteString(style.Indent(1))
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyGuardrailConfigPathMain)))
        b.WriteString("\n")
        b.WriteString(style.Indent(1))
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyGuardrailConfigPathDesc)))
        b.WriteString(style.Section(style.SubSectionGap))

        // FIX 2: use i18n keys
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyGuardrailRelax)))
        b.WriteString("\n")
        b.WriteString(style.BodyStyle.Render(i18n.T(i18n.KeyGuardrailJustFill)))

        // ListSelect
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(g.list.ViewWithNumbers())
        b.WriteString(style.Section(style.SubSectionGap))

        // Template path
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                i18n.T(i18n.KeyGuardrailTemplatePathLabel),
        ))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                i18n.T(i18n.KeySettingsConfigPath),
        ))
        b.WriteString("\n")
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                i18n.T(i18n.KeyGuardrailNicheFolderPath),
        ))

        // Key hints
        b.WriteString(style.Section(style.SectionGap))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("1  %s    2  %s    q  %s",
                        i18n.T(i18n.KeyGuardrailAutoGenerate), i18n.T(i18n.KeyGuardrailSeeExample), i18n.T(i18n.KeyGuardrailQuitExit),
                )))

        return b.String()
}
