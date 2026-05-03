// Package comms implements the Communication screens: Compose, History, and Follow-Up.
// These are screens 15–17 from the WaClaw TUI specification.
//
// Doc source: doc/09-screens-communicate.md
package comms

import (
        "fmt"
        "strings"
        "unicode/utf8"

        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/WaClaw-App/waclaw/internal/tui/util"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
        "github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Shared constants (DRY — used by all three screen models)
// ---------------------------------------------------------------------------

const (
        // defaultFallbackWidth is the assumed terminal width when the actual width
        // is not yet known (e.g. before the first tea.WindowSizeMsg).
        defaultFallbackWidth = 60

        // maxTextAreaWidth is the maximum width for the compose text area.
        // Overridable by backend via HandleUpdate("max_text_area_width").
        maxTextAreaWidth = 80

        // barMaxWidth is the maximum character width for inline bar charts.
        barMaxWidth = 20

        // composeDefaultMaxChars is the default soft character limit for compose messages.
        // 0 means unlimited until the backend sends a max_chars limit via HandleNavigate.
        composeDefaultMaxChars = 0

        // composeMinBoxHeight is the minimum visible height of the compose text area.
        composeMinBoxHeight = 4

        // maxSendingLeadsVisible limits how many leads are shown in the sending view.
        // Overridable by backend via HandleUpdate("max_sending_visible").
        maxSendingLeadsVisible = 3

        // maxColdLeadsVisible limits how many cold leads are shown before truncating.
        // Overridable by backend via HandleUpdate("max_cold_visible").
        maxColdLeadsVisible = 6

        // dataRainParticleCount is the number of particles in the data rain ambient effect.
        dataRainParticleCount = 40

        // daysInWeek is the number of days in a week — used by history bar charts.
        daysInWeek = 7

        // composeSendEnterCount is the number of consecutive enters required to send.
        composeSendEnterCount = 2

        // snippetTruncationMargin is the character margin subtracted for snippet text.
        snippetTruncationMargin = 20

        // minSnippetTextWidth is the minimum width for snippet text display.
        minSnippetTextWidth = 20

        // dayLabelMaxRunes is the maximum rune count for day labels in bar charts.
        dayLabelMaxRunes = 6

        // progressBarPadding is the padding subtracted for progress bar width.
        progressBarPadding = 8

        // titleIndicatorPadding is the padding between title and indicator.
        titleIndicatorPadding = 4

        // minTitleIndicatorSpacing is the minimum spacing between title and indicator.
        minTitleIndicatorSpacing = 4

        // defaultDataRainWidth is the default width for the data rain component
        // when the terminal size is not yet known.
        defaultDataRainWidth = defaultFallbackWidth

        // arrowIndicator is the arrow glyph used for selected niche rows.
        arrowIndicator = " ▸"
)

// ---------------------------------------------------------------------------
// Shared screen base type (DRY — eliminates three identical base structs)
// ---------------------------------------------------------------------------

// screenBase provides the Screen.ID() boilerplate for all comms screens.
type screenBase struct {
        id protocol.ScreenID
}

func (b screenBase) ID() protocol.ScreenID { return b.id }

// ---------------------------------------------------------------------------
// Shared render helpers (DRY — eliminates duplicated rendering patterns)
// ---------------------------------------------------------------------------

// calcSepWidth returns the separator width given the screen's current width.
func calcSepWidth(w int) int {
        if w <= 0 {
                w = defaultFallbackWidth
        }
        return w - 2
}

// writeSeparator appends a horizontal separator line to the builder.
// P3 rule: No box-drawing characters (─ │ ┌ └ etc). Hierarchy = brightness.
// Uses dim dots for subtle visual separation.
func writeSeparator(b *strings.Builder, width int) {
        b.WriteString(style.DimStyle.Render(strings.Repeat("·", calcSepWidth(width))))
        b.WriteString("\n\n")
}

// renderSectionTitle renders a section title using bold heading style.
// P3 rule: No box-drawing characters (── │). Hierarchy = brightness + size.
func renderSectionTitle(title string) string {
        return style.HeadingStyle.Render(title)
}

// renderLineNumber renders a two-digit numbered prefix (e.g. "01  ").
// This pattern was duplicated 5× in followup.go.
func renderLineNumber(n int) string {
        return lipgloss.NewStyle().Foreground(style.TextDim).Render(fmt.Sprintf("%02d", n)) + "  "
}

// ---------------------------------------------------------------------------
// Shared map helpers (DRY — delegates to util package for single source of truth)
// ---------------------------------------------------------------------------

func strVal(m map[string]any, key string) string {
        return util.ToString(m[key], "")
}

func intVal(m map[string]any, key string) int {
        return util.ToInt(m[key], 0)
}

func boolVal(m map[string]any, key string) bool {
        return util.ToBool(m[key], false)
}

// ---------------------------------------------------------------------------
// Shared render helpers (DRY — 3G audit additions)
// ---------------------------------------------------------------------------

// renderTitleWithIndicator renders a heading with a right-aligned indicator.
// Consolidates the repeated pattern of HeadingStyle + spacing + styled indicator.
// Uses rune-count widths (not byte lengths) for correct alignment with multi-byte text.
func renderTitleWithIndicator(title, indicator string, indicatorStyle lipgloss.Style, width int) string {
        titleRendered := style.HeadingStyle.Render(title)
        titleWidth := utf8.RuneCountInString(title)
        indicatorWidth := utf8.RuneCountInString(indicator)
        spacing := width - titleWidth - indicatorWidth - titleIndicatorPadding
        if spacing < minTitleIndicatorSpacing {
                spacing = minTitleIndicatorSpacing
        }
        return titleRendered + strings.Repeat(" ", spacing) + indicatorStyle.Render(indicator)
}

// renderSelectAction renders the "↑↓ <noun>" action label.
// The ↑↓ arrows are included directly — callers must NOT add their own.
func renderSelectAction(noun string) string {
        if noun == "" {
                return style.MutedStyle.Render("↑↓")
        }
        return style.MutedStyle.Render(fmt.Sprintf("↑↓ %s", strings.ToLower(noun)))
}
