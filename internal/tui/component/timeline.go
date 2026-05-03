package component

import (
        "fmt"
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/charmbracelet/lipgloss"
)

// TimelineEvent represents a single event in a sequential timeline.
type TimelineEvent struct {
        // Time is when the event occurred.
        Time time.Time

        // Source is the origin tag (e.g. "[web_dev]", "scrape.web_dev").
        Source string

        // Title is the main event text (e.g. business name).
        Title string

        // Detail is the secondary detail (e.g. response text, status).
        Detail string

        // Icon is an optional emoji or symbol prefix.
        Icon string

        // Highlight indicates this event should stand out (e.g. response received).
        Highlight bool

        // StaggerIndex is the animation delay index for stagger fade-in.
        StaggerIndex int

        // Visible indicates whether this event has faded in.
        Visible bool
}

// Timeline renders a sequential event timeline with stagger fade-in animation.
//
// Visual spec (from doc/05-screens-monitor-response.md):
//   - Recent events listed chronologically
//   - Each event: HH:MM [source] title detail
//   - Highlight events (responses) get warm amber flash
//   - Stagger fade-in: 120ms per event when new events arrive
//   - Used in: monitor activity feed, history timeline
type Timeline struct {
        // Events holds the timeline entries (newest first).
        Events []TimelineEvent

        // MaxVisible limits the number of displayed events.
        MaxVisible int

        // ShowTime whether to display the time prefix.
        ShowTime bool

        // ShowSource whether to display the source tag.
        ShowSource bool

        // AnimStart tracks when the stagger animation began.
        AnimStart time.Time
}

// NewTimeline creates a Timeline with sensible defaults.
func NewTimeline() Timeline {
        return Timeline{
                MaxVisible: 8,
                ShowTime:   true,
                ShowSource: true,
        }
}

// Add appends a new event and triggers the stagger animation.
func (t *Timeline) Add(event TimelineEvent) {
        event.StaggerIndex = len(t.Events)
        event.Visible = false
        t.Events = append([]TimelineEvent{event}, t.Events...)
        t.AnimStart = time.Now()
}

// Tick advances the stagger fade-in animation.
// Returns true if any new events became visible.
func (t *Timeline) Tick(now time.Time) bool {
        changed := false
        for i := range t.Events {
                if t.Events[i].Visible {
                        continue
                }
                // Each event becomes visible after stagger delay from anim.TimelineStagger.
                visibleAt := t.AnimStart.Add(time.Duration(t.Events[i].StaggerIndex) * anim.TimelineStagger)
                if now.After(visibleAt) {
                        t.Events[i].Visible = true
                        changed = true
                }
        }
        return changed
}

// View renders the timeline as a list of events.
func (t Timeline) View() string {
        return t.ViewAt(time.Now())
}

// ViewAt renders the timeline with a specific time for animation control.
func (t Timeline) ViewAt(now time.Time) string {
        var lines []string
        count := min(len(t.Events), t.MaxVisible)

        for i := 0; i < count; i++ {
                event := t.Events[i]

                // Skip invisible events (still animating).
                if !event.Visible {
                        continue
                }

                var b strings.Builder

                // Time prefix.
                if t.ShowTime && !event.Time.IsZero() {
                        timeStr := event.Time.Format("15:04")
                        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(timeStr))
                        b.WriteString("  ")
                }

                // Source tag.
                if t.ShowSource && event.Source != "" {
                        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(event.Source))
                        b.WriteString("  ")
                }

                // Icon.
                if event.Icon != "" {
                        b.WriteString(event.Icon)
                        b.WriteString(" ")
                }

                // Title and detail with appropriate styling.
                if event.Highlight {
                        // Highlighted events (e.g. responses) get warm amber.
                        b.WriteString(lipgloss.NewStyle().Foreground(style.Warning).Bold(true).Render(event.Title))
                        if event.Detail != "" {
                                b.WriteString("    ")
                                b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(event.Detail))
                        }
                } else {
                        b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(event.Title))
                        if event.Detail != "" {
                                b.WriteString("    ")
                                b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(event.Detail))
                        }
                }

                lines = append(lines, b.String())
        }

        return strings.Join(lines, "\n")
}

// ViewCompact renders a compact single-line-per-event timeline
// suitable for the monitor dashboard's activity feed.
func (t Timeline) ViewCompact() string {
        var lines []string
        count := min(len(t.Events), 7) // Monitor shows ~7 recent events.

        for i := 0; i < count; i++ {
                event := t.Events[i]
                var b strings.Builder

                // Time.
                if !event.Time.IsZero() {
                        b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(event.Time.Format("15:04")))
                        b.WriteString("  ")
                }

                // Source.
                if event.Source != "" {
                        b.WriteString(lipgloss.NewStyle().Foreground(style.Accent).Render(
                                fmt.Sprintf("[%s]", event.Source),
                        ))
                        b.WriteString("   ")
                }

                // Title.
                if event.Highlight {
                        b.WriteString(lipgloss.NewStyle().Foreground(style.Warning).Render(event.Title))
                } else {
                        b.WriteString(lipgloss.NewStyle().Foreground(style.TextMuted).Render(event.Title))
                }

                // Detail.
                if event.Detail != "" {
                        b.WriteString("    ")
                        if event.Highlight {
                                b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render(event.Detail))
                        } else {
                                b.WriteString(lipgloss.NewStyle().Foreground(style.TextDim).Render(event.Detail))
                        }
                }

                lines = append(lines, b.String())
        }

        return strings.Join(lines, "\n")
}

// String implements fmt.Stringer.
func (t Timeline) String() string {
        return fmt.Sprintf("Timeline{events=%d, maxVisible=%d}", len(t.Events), t.MaxVisible)
}
