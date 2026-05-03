package component

import (
	"fmt"
	"strings"

	"github.com/WaClaw-App/waclaw/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

// DayData represents a single day's data point for the mini chart.
type DayData struct {
	// Day is the day label (e.g. "Sen", "Sel", "Rab").
	Day string

	// Value is the data value for this day.
	Value int64

	// IsToday highlights the current day.
	IsToday bool
}

// MiniChart renders weekly bar charts for the history screen.
//
// Visual spec (from doc/09-screens-communicate.md):
//   - Weekly bar charts per day
//   - Bars are vertical, using block characters (█ ▇ ▆ ▅ ▄ ▃ ▂)
//   - Day labels below
//   - Today highlighted with accent color
//   - Used in: History screen weekly view
type MiniChart struct {
	// Days holds the 7 days of data.
	Days []DayData

	// MaxValue is the maximum value for scaling. Auto-computed if 0.
	MaxValue int64

	// BarWidth is the character width of each bar.
	BarWidth int

	// ChartHeight is the number of rows for the bar area.
	ChartHeight int
}

// NewMiniChart creates a MiniChart for 7 days of data.
func NewMiniChart(days []DayData) MiniChart {
	mc := MiniChart{
		Days:        days,
		BarWidth:    3,
		ChartHeight: 5,
	}
	mc.MaxValue = mc.computeMax()
	return mc
}

// computeMax returns the maximum value across all days.
func (mc MiniChart) computeMax() int64 {
	var maxVal int64
	for _, d := range mc.Days {
		if d.Value > maxVal {
			maxVal = d.Value
		}
	}
	if maxVal == 0 {
		maxVal = 1 // Avoid division by zero.
	}
	return maxVal
}

// Block characters for different fill levels.
var blockChars = []string{"▂", "▃", "▄", "▅", "▆", "▇", "█"}

// View renders the mini chart as a multi-line string.
func (mc MiniChart) View() string {
	if len(mc.Days) == 0 {
		return ""
	}

	maxVal := mc.MaxValue
	if maxVal <= 0 {
		maxVal = mc.computeMax()
	}

	// Build bars for each day.
	// Each bar is ChartHeight rows tall, filled from bottom to top.
	barChars := make([][]string, len(mc.Days))
	for i, day := range mc.Days {
		barChars[i] = mc.renderBar(day, maxVal)
	}

	// Compose rows: top to bottom.
	var lines []string
	for row := 0; row < mc.ChartHeight; row++ {
		var b strings.Builder
		for _, bar := range barChars {
			if row < len(bar) {
				b.WriteString(bar[row])
			} else {
				b.WriteString(strings.Repeat(" ", mc.BarWidth))
			}
			b.WriteString(" ") // gap between bars
		}
		lines = append(lines, b.String())
	}

	// Day labels below the bars.
	var labelLine strings.Builder
	for _, day := range mc.Days {
		label := day.Day
		if len(label) > mc.BarWidth {
			label = label[:mc.BarWidth]
		}
		// Pad label to bar width.
		padding := mc.BarWidth - len(label)
		padded := strings.Repeat(" ", padding/2) + label + strings.Repeat(" ", padding-padding/2)

		if day.IsToday {
			labelLine.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Bold(true).Render(padded))
		} else {
			labelLine.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(padded))
		}
		labelLine.WriteString(" ")
	}
	lines = append(lines, labelLine.String())

	return strings.Join(lines, "\n")
}

// renderBar creates the rows for a single bar.
func (mc MiniChart) renderBar(day DayData, maxVal int64) []string {
	// Calculate how many rows should be filled.
	fillRatio := float64(day.Value) / float64(maxVal)
	fillRows := int(fillRatio * float64(mc.ChartHeight))
	if fillRows > mc.ChartHeight {
		fillRows = mc.ChartHeight
	}

	// Choose bar color.
	barColor := style.TextMuted
	if day.IsToday {
		barColor = style.Accent
	}

	// Build from top (empty) to bottom (filled).
	rows := make([]string, mc.ChartHeight)
	for i := 0; i < mc.ChartHeight; i++ {
		rowFromBottom := mc.ChartHeight - 1 - i
		if rowFromBottom < fillRows {
			// Filled row.
			rows[i] = lipgloss.NewStyle().Foreground(barColor).Render(
				strings.Repeat("█", mc.BarWidth),
			)
		} else if rowFromBottom == fillRows && fillRows > 0 {
			// Top of fill — use partial block for smooth appearance.
			rows[i] = lipgloss.NewStyle().Foreground(barColor).Render(
				strings.Repeat("▇", mc.BarWidth),
			)
		} else {
			// Empty row.
			rows[i] = strings.Repeat(" ", mc.BarWidth)
		}
	}

	return rows
}

// ViewInline renders a compact inline bar chart (single line).
// Used in stat cards and compact dashboard views.
func (mc MiniChart) ViewInline() string {
	maxVal := mc.MaxValue
	if maxVal <= 0 {
		maxVal = mc.computeMax()
	}

	var b strings.Builder
	for i, day := range mc.Days {
		ratio := float64(day.Value) / float64(maxVal)
		blockIdx := int(ratio * float64(len(blockChars)-1))
		if blockIdx >= len(blockChars) {
			blockIdx = len(blockChars) - 1
		}
		if blockIdx < 0 {
			blockIdx = 0
		}

		char := blockChars[blockIdx]
		if day.IsToday {
			b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render(char))
		} else {
			b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(char))
		}

		if i < len(mc.Days)-1 {
			b.WriteString(" ")
		}
	}

	return b.String()
}

// String implements fmt.Stringer.
func (mc MiniChart) String() string {
	return fmt.Sprintf("MiniChart{days=%d, max=%d}", len(mc.Days), mc.MaxValue)
}
