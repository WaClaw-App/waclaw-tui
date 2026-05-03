package component

import (
        "fmt"
        "strings"

        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/charmbracelet/lipgloss"
)

// ShieldLevel classifies the shield's health into three tiers.
// Used by ShieldArt to determine fill density, color, and label.
type ShieldLevel int

const (
        ShieldHealthy ShieldLevel = iota // 90-100: full fill, success green
        ShieldWarning                     // 50-89: partial fill, warning amber
        ShieldDanger                      // <50: minimal fill, danger red with crack
)

// ShieldArt renders a dynamic ASCII shield whose fill level reflects
// the aggregate health score of the anti-ban system.
//
// Visual spec (from doc/07-screens-workers-antiban.md):
//
//      Health 90-100 (HEALTHY):  solid fill, success green, label "HEALTHY"/"SEHAT"
//      Health 50-89  (WARNING):  partial fill, warning amber, label "WARNING"
//      Health <50    (DANGER):   minimal fill, danger red, ╳ crack, label "DANGER"/"BAHAYA"
//      Repair anim:  fill grows bottom→top at 50ms per health point
//
// The shield art is rendered as multi-line ASCII using block characters.
type ShieldArt struct {
        // Health is the aggregate health score [0, 100].
        Health int

        // Width controls the art width (default 12 for the inner shield).
        Width int
}

// NewShieldArt creates a ShieldArt with the given health score.
func NewShieldArt(health int) ShieldArt {
        if health < 0 {
                health = 0
        }
        if health > 100 {
                health = 100
        }
        return ShieldArt{Health: health, Width: 12}
}

// Level returns the shield level based on the current health score.
func (s ShieldArt) Level() ShieldLevel {
        switch {
        case s.Health >= 90:
                return ShieldHealthy
        case s.Health >= 50:
                return ShieldWarning
        default:
                return ShieldDanger
        }
}

// Label returns the i18n-aware label for the current shield level.
func (s ShieldArt) Label() string {
        switch s.Level() {
        case ShieldHealthy:
                return i18n.T(i18n.KeyShieldHealthy)
        case ShieldWarning:
                return i18n.T(i18n.KeyShieldWarning)
        default:
                return i18n.T(i18n.KeyShieldDanger)
        }
}

// Color returns the lipgloss color for the current shield level.
func (s ShieldArt) Color() lipgloss.Color {
        switch s.Level() {
        case ShieldHealthy:
                return style.Success
        case ShieldWarning:
                return style.Warning
        default:
                return style.Danger
        }
}

// View renders the shield ASCII art.
//
// The shield shape:
//
//             ╱╲
//            ╱  ╲
//           ╱ ░░ ╲       <- fill rows ( ░ fill, space = gap)
//          ╱ ░░░░ ╲
//         ╱ ░░░░░░ ╲
//        ╱──────────╲
//       ╱  ■ ■ ■ ■   ╲   <- slot row (only for healthy)
//      ╱──────────────╲
//      ────────────────────
func (s ShieldArt) View() string {
        w := s.Width
        if w < 8 {
                w = 8
        }

        fillColor := s.Color()
        fillChar := "░"
        gapChar := " "

        // Determine fill density based on health.
        fillRatio := float64(s.Health) / 100.0

        // Build shield line by line.
        var lines []string

        // Top point.
        lines = append(lines, lipgloss.NewStyle().Foreground(fillColor).Render(
                fmt.Sprintf("%s╱%s╲", indent(w/2+2), strings.Repeat(" ", w)),
        ))

        // Upper fill rows.
        fillRows := 3
        for i := 0; i < fillRows; i++ {
                rowFillCount := int(float64(w) * fillRatio)
                // Fill increases towards bottom.
                rowFillCount += i
                if rowFillCount > w {
                        rowFillCount = w
                }

                fillPart := strings.Repeat(fillChar, rowFillCount)
                gapPart := strings.Repeat(gapChar, w-rowFillCount)

                inner := fillPart + gapPart
                lines = append(lines, lipgloss.NewStyle().Foreground(fillColor).Render(
                        fmt.Sprintf("%s╱ %s ╲", indent(w/2+1-fillRows+i), inner),
                ))
        }

        // Separator.
        lines = append(lines, lipgloss.NewStyle().Foreground(fillColor).Render(
                fmt.Sprintf("%s╱%s╲", indent(2), strings.Repeat("─", w)),
        ))

        // Slot row (only for healthy).
        if s.Level() == ShieldHealthy {
                slotRow := "  ■ ■ ■ ■   "
                lines = append(lines, lipgloss.NewStyle().Foreground(fillColor).Render(
                        fmt.Sprintf("╱ %s ╲", slotRow[:min(w, len(slotRow))]),
                ))
        }

        // Bottom separator.
        lines = append(lines, lipgloss.NewStyle().Foreground(fillColor).Render(
                fmt.Sprintf("╱%s╲", strings.Repeat("─", w+2)),
        ))

        // Danger crack.
        if s.Level() == ShieldDanger {
                crackLine := fmt.Sprintf("%s───╳%s", indent(3), strings.Repeat("─", w-3))
                lines = append(lines, lipgloss.NewStyle().Foreground(style.Danger).Render(crackLine))
        } else {
                lines = append(lines, lipgloss.NewStyle().Foreground(fillColor).Render(
                        fmt.Sprintf("%s%s", indent(1), strings.Repeat("─", w+4)),
                ))
        }

        // Health score label below.
        healthLabel := lipgloss.NewStyle().Foreground(fillColor).Bold(true).Render(
                fmt.Sprintf("%d/100", s.Health),
        )
        levelLabel := lipgloss.NewStyle().Foreground(fillColor).Bold(true).Render(
                s.Label(),
        )
        lines = append(lines, fmt.Sprintf("  %s  %s", healthLabel, levelLabel))

        return strings.Join(lines, "\n")
}

// indent returns a left-margin indent string.
func indent(n int) string {
        if n <= 0 {
                return ""
        }
        return strings.Repeat(" ", n)
}
