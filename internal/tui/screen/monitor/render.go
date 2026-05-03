package monitor

import (
        "fmt"
        "strings"

        "github.com/WaClaw-App/waclaw/internal/tui/bus"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/charmbracelet/lipgloss"
)

// minRenderWidth is the minimum width for rendering separators, borders, and headers.
// Shared by both Dashboard and Response screens to avoid duplicating the magic number 40.
const minRenderWidth = 40

// statColumnWidth is the column width for the stats grid.
// Used for both the grid data and header alignment padding.
const statColumnWidth = 30

// renderDimSeparator renders a ── line in TextDim color.
// Shared by both Dashboard and Response screens.
func renderDimSeparator(width int) string {
        return lipgloss.NewStyle().Foreground(style.TextDim).Render(
                strings.Repeat("─", max(width-4, minRenderWidth)),
        ) + "\n"
}

// renderStatusDot renders an active/inactive status indicator.
// Active renders ● in success color, inactive renders ○ in dim color.
func renderStatusDot(active bool) string {
        if active {
                return lipgloss.NewStyle().Foreground(style.Success).Render("●")
        }
        return lipgloss.NewStyle().Foreground(style.TextDim).Render("○")
}

// publish safely publishes a message to the event bus if non-nil.
// Eliminates the repeated `if bus != nil { bus.Publish(...) }` guard.
func publish(b *bus.Bus, msg any) {
        if b != nil {
                b.Publish(msg)
        }
}

// extractString extracts a string value from a map[string]any by key.
// Returns empty string if key is missing or value is not a string.
func extractString(m map[string]any, key string) string {
        if v, ok := m[key].(string); ok {
                return v
        }
        return ""
}

// extractInt extracts an int value from a map[string]any by key.
// Supports float64 (JSON), int64, and int types.
func extractInt(m map[string]any, key string) int {
        switch v := m[key].(type) {
        case float64:
                return int(v)
        case int64:
                return int(v)
        case int:
                return v
        }
        return 0
}

// extractFloat64 extracts a float64 value from a map[string]any by key.
func extractFloat64(m map[string]any, key string) float64 {
        if v, ok := m[key].(float64); ok {
                return v
        }
        return 0
}

// renderWAStatus renders the "● wa nyambung (N nomor)" status line.
// Extracted from the 3 dashboard views that duplicated this pattern.
func renderWAStatus(numCount int) string {
        statusText := lipgloss.NewStyle().Foreground(style.Success).Render("●") + " " + i18n.T(i18n.KeyMonitorWAConnected)
        if numCount > 0 {
                statusText = fmt.Sprintf("%s (%d %s)", statusText, numCount, i18n.T(i18n.KeyMonitorWARotatorNum))
        }
        return statusText
}

// renderDaySummary renders a "N msgs · N responses · N converts" line.
// Shared format used by both night and pending-responses dashboard states.
func renderDaySummary(stats [4]int64, prefix string) string {
        return fmt.Sprintf("%s%d %s · %d %s · %d %s",
                prefix,
                stats[1], i18n.T(i18n.KeyMonitorMsgsSent),
                stats[2], i18n.T(i18n.KeyMonitorResponses),
                stats[3], i18n.T(i18n.KeyMonitorConverts))
}
