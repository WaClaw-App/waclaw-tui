package infra

import (
        "fmt"
        "strings"

        "github.com/WaClaw-App/waclaw/internal/tui/component"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Shared helpers for the infra screen package
// ---------------------------------------------------------------------------
//
// These helpers consolidate patterns previously duplicated across
// workers.go, antiban.go, settings.go, and guardrail.go, following
// the DRY principle.

// kvItem is a key-value pair for label-value display sections.
// Used by both antiban.go (pattern/spam guard config) and settings.go
// (file paths, active config).
type kvItem struct {
        Label string
        Value string
}

// renderKVSection renders a section with a heading and key-value label rows.
// Shared by antiban.go and settings.go for consistent label-value layout.
func renderKVSection(b *strings.Builder, heading string, items []kvItem) {
        b.WriteString(style.SubHeadingStyle.Render(heading))
        b.WriteString(style.Section(style.SubSectionGap))
        for _, item := range items {
                b.WriteString(style.Indent(1))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Width(20).Render(item.Label))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(item.Value))
                b.WriteString("\n")
        }
}

// renderGutterError renders an error or warning line with context gutter
// and pointer underline. This pattern was duplicated in settings.go
// (viewReloadError) and guardrail.go (renderErrorDetail).
//
// Parameters:
//   - b: output builder
//   - linePrefix: the error/warning line text (e.g., "baris 23: parse error")
//   - context: surrounding code lines shown with │ gutter
//   - pointer: the ^^^^ underline showing exact position
//   - isError: true for danger-styled errors, false for warning-styled warnings
func renderGutterError(b *strings.Builder, linePrefix string, context []string, pointer string, isError bool) {
        prefixStyle := style.DangerStyle
        icon := "✗"
        if !isError {
                prefixStyle = style.WarningStyle
                icon = "⚠"
        }

        b.WriteString(style.Indent(1))
        b.WriteString(prefixStyle.Render(fmt.Sprintf("%s  %s", icon, linePrefix)))
        b.WriteString("\n")

        // Context lines with gutter
        for _, ctx := range context {
                b.WriteString(style.Indent(2))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render("│  "))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(ctx))
                b.WriteString("\n")
        }

        // Pointer with gutter
        if pointer != "" {
                b.WriteString(style.Indent(2))
                b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render("│  "))
                // Calculate padding to align pointer with the error position
                if len(context) > 0 {
                        padLen := len(context[0]) - len(strings.TrimLeft(context[0], " \t"))
                        if padLen > 0 {
                                b.WriteString(strings.Repeat(" ", padLen))
                        }
                }
                b.WriteString(prefixStyle.Render(pointer))
                b.WriteString("\n")
        }
}

// renderStatusIcon returns the appropriate icon and color for a status string.
// Shared by guardrail.go (validation status) and antiban.go (slot status)
// to ensure consistent icon/color mapping across screens.
func renderStatusIcon(status string) (icon string, color lipgloss.Color) {
        switch status {
        case "ok", "active":
                return "✓", style.Success
        case "error", "flagged":
                return "✗", style.Danger
        case "warning":
                return "⚠", style.Warning
        case "checking", "scanning":
                return "●", style.Text
        case "waiting", "idle", "cooldown":
                return "○", style.TextDim
        case "done":
                return "✓", style.TextMuted
        default:
                return "?", style.TextDim
        }
}

// renderHeadWithStatus renders a heading with a right-aligned status badge.
// This replaces the hardcoded space-alignment pattern (e.g., statusAlignGap)
// with a computed layout that works regardless of heading length.
func renderHeadWithStatus(heading, statusText string, statusColor lipgloss.Color) string {
        headingRendered := style.HeadingStyle.Render(heading)
        statusRendered := lipgloss.NewStyle().Foreground(statusColor).Render(statusText)
        // Place status after heading with enough spacing
        return headingRendered + "  " + statusRendered
}

// renderStageLine renders a single pipeline stage line (label + progress bar).
// Extracted from workers.go to DRY up the stage rendering in both
// viewOverview and viewDetail. W-DRY01.
func renderStageLine(label string, stage StageInfo, barWidth int) string {
        bar := component.NewProgressBar(barWidth)
        bar.Percent = stage.Progress
        bar.Label = stage.Detail
        bar.ShowPercent = !stage.Done

        var b strings.Builder
        b.WriteString(style.Indent(1))
        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Width(8).Render(label))
        b.WriteString(" ")
        b.WriteString(bar.View())
        return b.String()
}
