package component

import (
	"fmt"
	"strings"
	"time"

	"github.com/WaClaw-App/waclaw/internal/tui/anim"
	"github.com/WaClaw-App/waclaw/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

// LogLevel represents the severity level of a log entry.
type LogLevel int

const (
	LogInfo LogLevel = iota
	LogWarning
	LogError
	LogSuccess
)

// LogEntry represents a single log line in the stream.
type LogEntry struct {
	// Time is when the log was generated.
	Time time.Time

	// Source is the origin tag (e.g. "scrape.web_dev", "send.slot_2", "antiban").
	Source string

	// Message is the log text.
	Message string

	// Level is the log severity.
	Level LogLevel

	// SlideInAt tracks when this entry began its slide-in animation.
	SlideInAt time.Time

	// FadingOut indicates this entry is being removed from the visible area.
	FadingOut bool

	// FadeOutAt tracks when the fade-out began.
	FadeOutAt time.Time
}

// levelColor returns the color for the given log level.
func levelColor(level LogLevel) lipgloss.Color {
	switch level {
	case LogInfo:
		return style.TextDim
	case LogWarning:
		return style.Warning
	case LogError:
		return style.Danger
	case LogSuccess:
		return style.Success
	default:
		return style.TextDim
	}
}

// levelIcon returns a small icon for the log level.
func levelIcon(level LogLevel) string {
	switch level {
	case LogWarning:
		return "⚠"
	case LogError:
		return "✗"
	case LogSuccess:
		return "✓"
	default:
		return ""
	}
}

// LogStream renders a live system logs stream with color-coded levels.
//
// Visual spec (from doc/10-global-overlays.md — nerd stats expanded view):
//   - 5-line tail -f style display
//   - Color-coded levels: info=text_dim, warning=warning amber, error=danger red, success=success green
//   - Format: HH:MM:SS source          message
//   - Source tags: scrape.{niche}, send.slot_{n}, antiban, followup, config, etc.
//   - New entries: slide in from right 100ms (anim.LogStreamSlideIn)
//   - Overflow: oldest line fade out 150ms (anim.LogStreamOverflowFade)
//   - Max 5 lines visible; rest goes to full log file
//   - Pause on scroll up, resume after 5s (anim.LogStreamPauseTimeout)
type LogStream struct {
	// Entries holds all log entries (newest first).
	Entries []LogEntry

	// MaxVisible is the number of lines to display (default 5).
	MaxVisible int

	// Paused indicates the stream is paused (user scrolled up).
	Paused bool

	// PauseSince tracks when the pause started.
	PauseSince time.Time

	// SourceWidth is the column width for the source tag.
	SourceWidth int
}

// NewLogStream creates a LogStream with default settings.
func NewLogStream() LogStream {
	return LogStream{
		MaxVisible:  5,
		SourceWidth: 18,
	}
}

// Add appends a new log entry and triggers slide-in animation.
// If at max capacity, the oldest entry starts fading out.
func (ls *LogStream) Add(entry LogEntry) {
	now := time.Now()
	entry.SlideInAt = now

	ls.Entries = append([]LogEntry{entry}, ls.Entries...)

	// Trigger fade-out on overflow.
	if len(ls.Entries) > ls.MaxVisible {
		oldest := &ls.Entries[len(ls.Entries)-1]
		if !oldest.FadingOut {
			oldest.FadingOut = true
			oldest.FadeOutAt = now
		}
	}
}

// Tick advances animations and removes fully-faded entries.
func (ls *LogStream) Tick(now time.Time) {
	// Remove entries that have completed their fade-out.
	var visible []LogEntry
	for _, e := range ls.Entries {
		if e.FadingOut && now.Sub(e.FadeOutAt) >= anim.LogStreamOverflowFade {
			continue // Remove fully faded entry.
		}
		visible = append(visible, e)
	}
	ls.Entries = visible

	// Auto-resume after pause timeout.
	if ls.Paused && now.Sub(ls.PauseSince) >= anim.LogStreamPauseTimeout {
		ls.Paused = false
	}
}

// Pause pauses the log stream.
func (ls *LogStream) Pause() {
	ls.Paused = true
	ls.PauseSince = time.Now()
}

// View renders the log stream.
func (ls LogStream) View() string {
	return ls.ViewAt(time.Now())
}

// ViewAt renders the log stream with a specific time for animation control.
func (ls LogStream) ViewAt(now time.Time) string {
	var lines []string
	count := min(len(ls.Entries), ls.MaxVisible)

	for i := 0; i < count; i++ {
		entry := ls.Entries[i]
		lines = append(lines, ls.renderEntry(entry, now))
	}

	if len(lines) == 0 {
		return lipgloss.NewStyle().Foreground(style.TextDim).Render("  no logs yet")
	}

	return strings.Join(lines, "\n")
}

// renderEntry renders a single log entry with appropriate styling.
func (ls LogStream) renderEntry(entry LogEntry, now time.Time) string {
	color := levelColor(entry.Level)
	isSliding := now.Sub(entry.SlideInAt) < anim.LogStreamSlideIn

	var b strings.Builder

	// Timestamp.
	b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(
		entry.Time.Format("15:04:05"),
	))
	b.WriteString(" ")

	// Source tag (padded to fixed width).
	source := entry.Source
	if len(source) > ls.SourceWidth {
		source = source[:ls.SourceWidth]
	}
	paddedSource := source + strings.Repeat(" ", max(0, ls.SourceWidth-len(source)))
	b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(paddedSource))
	b.WriteString(" ")

	// Level icon.
	if icon := levelIcon(entry.Level); icon != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(color).Render(icon))
		b.WriteString(" ")
	}

	// Message — accent highlight during slide-in, normal color after.
	if isSliding {
		b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render(entry.Message))
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(color).Render(entry.Message))
	}

	return b.String()
}

// ViewHeader renders the "── logs ──" header line for the nerd stats panel.
func (ls LogStream) ViewHeader() string {
	return lipgloss.NewStyle().Foreground(style.TextDim).Render(
		"── logs ────────────────────────────────────────────",
	)
}

// ViewWithHeader renders the complete log section with header.
func (ls LogStream) ViewWithHeader() string {
	return ls.ViewWithHeaderAt(time.Now())
}

// ViewWithHeaderAt renders the complete log section with header at a specific time.
func (ls LogStream) ViewWithHeaderAt(now time.Time) string {
	return ls.ViewHeader() + "\n\n" + ls.ViewAt(now)
}

// String implements fmt.Stringer.
func (ls LogStream) String() string {
	return fmt.Sprintf("LogStream{entries=%d, maxVisible=%d, paused=%v}",
		len(ls.Entries), ls.MaxVisible, ls.Paused)
}
