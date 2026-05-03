package component

import (
	"fmt"
	"strings"
	"time"

	"github.com/WaClaw-App/waclaw/internal/tui/anim"
	"github.com/WaClaw-App/waclaw/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

// StatCard renders a dashboard stat with live increment and scale bump.
//
// Visual spec (from doc/05-screens-monitor-response.md, doc/15-micro-interactions.md):
//   - Live increment: flash bright + scale 1.05x for 200ms when value changes
//   - Breathing: subtle opacity pulse on numbers (0.9→1.0→0.9, 4s cycle)
//   - Layout: label on left, value on right (or stacked)
//   - Colors: label=TextMuted, value=Text (bright on increment)
type StatCard struct {
	// Label is the description text (e.g. "lead nemu", "pesan terkirim").
	Label string

	// Value is the current numeric value.
	Value int64

	// PreviousValue tracks the last value for increment detection.
	PreviousValue int64

	// IncrementAt tracks when the last increment occurred (for flash animation).
	IncrementAt time.Time

	// Width controls the card width for alignment.
	Width int
}

// NewStatCard creates a StatCard with the given label and initial value.
func NewStatCard(label string, value int64) StatCard {
	return StatCard{
		Label:  label,
		Value:  value,
		Width:  20,
	}
}

// SetValue updates the stat value and triggers the increment animation
// if the new value is different from the current value.
func (sc *StatCard) SetValue(v int64) {
	if v != sc.Value {
		sc.PreviousValue = sc.Value
		sc.Value = v
		sc.IncrementAt = time.Now()
	}
}

// IsFlashing returns whether the stat is currently in its increment flash
// animation window (200ms after value change, from anim.NumberFlash).
func (sc StatCard) IsFlashing(now time.Time) bool {
	return now.Sub(sc.IncrementAt) < anim.NumberFlash
}

// renderValue renders a numeric value with flash or normal styling.
func renderValue(text string, flashing bool) string {
	if flashing {
		return lipgloss.NewStyle().Foreground(style.Text).Bold(true).Render(text)
	}
	return lipgloss.NewStyle().Foreground(style.Text).Render(text)
}

// View renders the stat card as a single line.
func (sc StatCard) View() string {
	return sc.ViewAt(time.Now())
}

// ViewAt renders the stat card with a specific time for animation control.
func (sc StatCard) ViewAt(now time.Time) string {
	labelStyle := lipgloss.NewStyle().Foreground(style.TextMuted).Width(sc.Width / 2)
	valueText := fmt.Sprintf("%d", sc.Value)
	valueStr := renderValue(valueText, sc.IsFlashing(now))

	// Layout: label left-aligned, value right-aligned.
	return fmt.Sprintf("%s  %s", labelStyle.Render(sc.Label), valueStr)
}

// ViewStacked renders the stat card with label above and value below.
// Used in the "hari ini / minggu ini" stat grid layout.
func (sc StatCard) ViewStacked(now time.Time) string {
	labelText := lipgloss.NewStyle().Foreground(style.TextMuted).Render(sc.Label)
	valueText := fmt.Sprintf("%d", sc.Value)
	valueStr := renderValue(valueText, sc.IsFlashing(now))

	return fmt.Sprintf("%s\n%s", labelText, valueStr)
}

// StatGrid renders a 2-column grid of stat cards.
// Used for the monitor dashboard's "hari ini / minggu ini" layout.
type StatGrid struct {
	// LeftCards are the stats for the left column.
	LeftCards []StatCard

	// RightCards are the stats for the right column.
	RightCards []StatCard

	// ColumnWidth is the width of each column.
	ColumnWidth int
}

// NewStatGrid creates a 2-column stat grid.
func NewStatGrid(leftLabels, rightLabels []string, leftVals, rightVals []int64) StatGrid {
	left := make([]StatCard, len(leftLabels))
	for i, l := range leftLabels {
		var v int64
		if i < len(leftVals) {
			v = leftVals[i]
		}
		left[i] = NewStatCard(l, v)
	}

	right := make([]StatCard, len(rightLabels))
	for i, l := range rightLabels {
		var v int64
		if i < len(rightVals) {
			v = rightVals[i]
		}
		right[i] = NewStatCard(l, v)
	}

	return StatGrid{
		LeftCards:   left,
		RightCards:  right,
		ColumnWidth: 25,
	}
}

// View renders the stat grid as aligned columns.
func (sg StatGrid) View() string {
	return sg.ViewAt(time.Now())
}

// ViewAt renders the stat grid with a specific time.
func (sg StatGrid) ViewAt(now time.Time) string {
	maxRows := len(sg.LeftCards)
	if len(sg.RightCards) > maxRows {
		maxRows = len(sg.RightCards)
	}

	var lines []string
	for i := 0; i < maxRows; i++ {
		var left, right string
		if i < len(sg.LeftCards) {
			left = sg.LeftCards[i].ViewAt(now)
		}
		if i < len(sg.RightCards) {
			right = sg.RightCards[i].ViewAt(now)
		}

		// Pad left column to column width.
		leftPadded := left + strings.Repeat(" ", max(0, sg.ColumnWidth-len(left)))
		lines = append(lines, leftPadded+right)
	}

	return strings.Join(lines, "\n")
}
