package component

import (
	"fmt"
	"strings"

	"github.com/WaClaw-App/waclaw/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

// ProgressBar renders a gradient-sweep progress bar with variable fill rate.
//
// Visual spec (from doc/15-micro-interactions.md):
//   - Fill uses gradient sweep animation (variable rate)
//   - Fill character: `━` (filled), `░` (unfilled)
//   - Color transitions: accent→warning when >80% of rate limit
//   - No borders, no boxes — just the bar
//
// Usage: screens embed ProgressBar and call View() to render.
type ProgressBar struct {
	// Width is the total character width of the bar (including label space).
	Width int

	// Percent is the fill ratio [0.0, 1.0].
	Percent float64

	// Label is an optional right-aligned label (e.g. "67%", "3/6 jam").
	Label string

	// ShowPercent appends the percentage after the bar when true.
	ShowPercent bool

	// FillChar is the character used for the filled portion.
	FillChar string

	// EmptyChar is the character used for the unfilled portion.
	EmptyChar string

	// AtLimit changes the fill color to warning (amber) when true.
	// Used for rate-limit bars approaching their ceiling.
	AtLimit bool
}

// NewProgressBar creates a ProgressBar with sensible defaults.
func NewProgressBar(width int) ProgressBar {
	return ProgressBar{
		Width:       width,
		Percent:     0,
		FillChar:    "━",
		EmptyChar:   "░",
		ShowPercent: true,
	}
}

// View renders the progress bar as a string.
func (p ProgressBar) View() string {
	if p.Width <= 0 {
		p.Width = 20
	}

	// Reserve space for label and percentage text.
	labelWidth := 0
	if p.Label != "" {
		labelWidth = len(p.Label) + 1 // +1 for spacing
	}
	percentText := ""
	if p.ShowPercent {
		percentText = fmt.Sprintf(" %3.0f%%", p.Percent*100)
		labelWidth += len(percentText)
	}

	barWidth := p.Width - labelWidth
	if barWidth < 4 {
		barWidth = 4
	}

	// Calculate fill count.
	fillCount := int(float64(barWidth) * p.Percent)
	if fillCount > barWidth {
		fillCount = barWidth
	}
	if fillCount < 0 {
		fillCount = 0
	}
	emptyCount := barWidth - fillCount

	// Choose fill color based on limit state.
	fillColor := style.Accent
	if p.AtLimit {
		fillColor = style.Warning
	}

	filled := lipgloss.NewStyle().Foreground(fillColor).Render(
		strings.Repeat(p.FillChar, fillCount),
	)
	empty := lipgloss.NewStyle().Foreground(style.TextDim).Render(
		strings.Repeat(p.EmptyChar, emptyCount),
	)

	// Assemble bar + label + percent.
	var b strings.Builder
	b.WriteString(filled)
	b.WriteString(empty)
	if p.Label != "" {
		b.WriteString(" ")
		b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(p.Label))
	}
	if p.ShowPercent {
		b.WriteString(percentText)
	}

	return b.String()
}

// SetPercent updates the fill ratio and returns the updated bar.
func (p ProgressBar) SetPercent(ratio float64) ProgressBar {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	p.Percent = ratio
	return p
}
