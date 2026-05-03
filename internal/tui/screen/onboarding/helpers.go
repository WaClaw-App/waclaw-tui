// Package onboarding implements the Boot and Login screens — the first
// impression and WhatsApp authentication flow for the WaClaw TUI.
//
// This file (helpers.go) contains shared constants, types, and helper
// functions used by both boot.go and login.go, extracted for DRY.
package onboarding

import (
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/internal/tui/component"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/WaClaw-App/waclaw/internal/tui/util"
        "github.com/WaClaw-App/waclaw/pkg/protocol"

        "github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Shared layout constants — avoid magic numbers
// ---------------------------------------------------------------------------

const (
        // separatorWidth is the horizontal separator width for boot/login separators.
        separatorWidth = 50

        // animationTickInterval is the refresh rate for animation frames.
        animationTickInterval = 50 * time.Millisecond

        // armyMarchStagger is the per-row stagger delay for the army march animation.
        // Doc: 80ms stagger per row.
        armyMarchStagger = 80 * time.Millisecond

        // contentIndent is the left padding for content below the logo.
        contentIndent = 6

        // menuDescColumn is the column position after which menu descriptions start.
        // Used to column-align description text across menu items.
        menuDescColumn = 14
)

// ---------------------------------------------------------------------------
// Shared types
// ---------------------------------------------------------------------------

// tickMsg is an internal message for animation frame updates.
// Shared by both BootModel and LoginModel (previously animTickMsg / loginTickMsg).
type tickMsg time.Time

// ---------------------------------------------------------------------------
// Separator rendering
// ---------------------------------------------------------------------------

// renderSeparator renders a dim horizontal separator line.
func renderSeparator() string {
        return style.DimStyle.Render(strings.Repeat("─", separatorWidth))
}

// renderLabeledSeparator renders a separator with a centered label.
// Format: "── label ─────────────" (two em-dashes prefix, dashes fill).
// Uses lipgloss.Width for correct measurement of styled/ANSI strings.
func renderLabeledSeparator(label string) string {
        prefix := "── "
        suffix := " "
        availableWidth := separatorWidth - len(prefix) - lipgloss.Width(label) - len(suffix)
        if availableWidth < 0 {
                availableWidth = 0
        }
        return style.DimStyle.Render(prefix) +
                style.SubHeadingStyle.Render(label) +
                style.DimStyle.Render(suffix + strings.Repeat("─", availableWidth))
}

// ---------------------------------------------------------------------------
// Breathing dot rendering
// ---------------------------------------------------------------------------

// renderBreathingDot renders a ● indicator with breathing pulse.
func renderBreathingDot(opacities []float64, idx int, brightColor, dimColor lipgloss.Color) string {
        if idx < len(opacities) {
                return component.RenderBreathing("●", opacities[idx], brightColor, dimColor)
        }
        return component.RenderBreathing("●", 1.0, brightColor, dimColor)
}

// renderIndicatorLine renders a single status indicator line with breathing dot.
// Extracted to DRY up the repeated pattern across all returning/login state views.
func renderIndicatorLine(indent int, opacities []float64, idx int, brightColor, dimColor lipgloss.Color, text string) string {
        var b strings.Builder
        b.WriteString(style.Indent(indent))
        b.WriteString(renderBreathingDot(opacities, idx, brightColor, dimColor))
        b.WriteString(" ")
        b.WriteString(text)
        b.WriteString("\n")
        return b.String()
}

// renderDimIndicatorLine renders a dimmed indicator line (for license expired state).
func renderDimIndicatorLine(indent int, text string) string {
        var b strings.Builder
        b.WriteString(style.Indent(indent))
        b.WriteString(style.DimStyle.Render("●"))
        b.WriteString(" ")
        b.WriteString(text)
        b.WriteString("\n")
        return b.String()
}

// ---------------------------------------------------------------------------
// Attention flash rendering (amber 2x for "response baru!")
// ---------------------------------------------------------------------------

// renderAttentionFlash renders text with an amber double-flash effect.
// Flash pattern: bright → dim → bright → dim → settle at bright
// Total duration: anim.AttentionFlash (600ms)
func renderAttentionFlash(text string, start time.Time) string {
        elapsed := time.Since(start)
        flashDuration := anim.AttentionFlash

        if elapsed >= flashDuration {
                // Settled — bright amber
                return style.WarningStyle.Bold(true).Render(text)
        }

        // Double flash: 4 phases of 150ms each
        phase := int(elapsed / (flashDuration / 4))
        if phase%2 == 0 {
                return style.WarningStyle.Bold(true).Render(text)
        }
        return style.MutedStyle.Render(text)
}

// renderAttentionDot renders a ● with attention flash (amber).
func renderAttentionDot(opacities []float64, idx int, start time.Time) string {
        elapsed := time.Since(start)
        flashDuration := anim.AttentionFlash

        if elapsed >= flashDuration {
                // Settled — always bright
                return style.WarningStyle.Bold(true).Render("●")
        }

        // Double flash phase
        phase := int(elapsed / (flashDuration / 4))
        if phase%2 == 0 {
                return style.WarningStyle.Bold(true).Render("●")
        }
        return style.DimStyle.Render("●")
}

// ---------------------------------------------------------------------------
// Error cross ✗ rendering (red flash)
// ---------------------------------------------------------------------------

// renderErrorCross renders a ✗ with red flash effect.
// The ✗ is the only red element on the screen — auto-draws attention.
func renderErrorCross(start time.Time) string {
        elapsed := time.Since(start)
        flashDuration := anim.ErrorGlow

        if elapsed >= flashDuration {
                return style.DangerStyle.Render("✗")
        }

        // Flash: alternate between bold and normal red
        phase := int(elapsed / (flashDuration / 4))
        if phase%2 == 0 {
                return style.DangerStyle.Bold(true).Render("✗")
        }
        return style.DangerStyle.Render("✗")
}

// ---------------------------------------------------------------------------
// Menu stagger visibility
// ---------------------------------------------------------------------------

// isMenuStaggerVisible returns true if the menu item at index i should be
// visible based on the 120ms stagger animation (anim.MenuStagger).
func isMenuStaggerVisible(staggerStart time.Time, index int) bool {
        elapsed := time.Since(staggerStart)
        staggerDelay := time.Duration(index) * anim.MenuStagger
        return elapsed >= staggerDelay
}

// ---------------------------------------------------------------------------
// Backend communication helpers
// ---------------------------------------------------------------------------

// publishAction sends an action event to the backend via the bus.
func publishAction(b *bus.Bus, screen protocol.ScreenID, action string, params map[string]any) {
        if b != nil {
                b.Publish(bus.ActionMsg{
                        Action: action,
                        Screen: screen,
                        Params: params,
                })
        }
}

// ---------------------------------------------------------------------------
// Type conversion helpers — delegated to util package
// ---------------------------------------------------------------------------

// toInt is a deprecated wrapper around util.ToInt for backward compatibility.
// New code should use util.ToInt directly.
//
// Deprecated: Use util.ToInt(v, 0) instead.
func toInt(v any) int {
        return util.ToInt(v, 0)
}

// toStringSlice converts a []any to []string, skipping non-string elements.
// Wrapper around util.ToStringSlice for backward compatibility.
func toStringSlice(raw []any) []string {
        return util.ToStringSlice(raw)
}

// parseMarchingWorkers converts raw backend worker data into MarchingWorker slices.
func parseMarchingWorkers(raw []any) []MarchingWorker {
        workers := make([]MarchingWorker, 0, len(raw))
        for i, r := range raw {
                data, ok := r.(map[string]any)
                if !ok {
                        continue
                }
                name, _ := data[protocol.ParamName].(string)
                arrowCount := toInt(data[protocol.ParamArrowCount])
                if arrowCount == 0 {
                        // Default: decreasing arrow count based on position
                        arrowCount = 6 - i
                        if arrowCount < 3 {
                                arrowCount = 3
                        }
                }
                workers = append(workers, MarchingWorker{
                        Name:       name,
                        ArrowCount: arrowCount,
                })
        }
        return workers
}
