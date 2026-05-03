package component

import (
        "fmt"
        "math/rand"
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/charmbracelet/lipgloss"
)

// DataRain renders the ambient faint number scroll effect used in the
// monitor dashboard background.
//
// Visual spec (from doc/05-screens-monitor-response.md):
//   - Pattern: `░░ 3 7 1 4 8 2 9 ░░ 5 1 6 3 8 ░░ 2 7 4 ░░`
//   - Color: text_dim (happily visible but not distracting)
//   - Update interval: 5 seconds
//   - Pauses on interaction, resumes after 10s idle
//   - Represents lead IDs being processed — subliminal activity signal
//   - Numbers scroll from left to right: new numbers enter left, old exit right
type DataRain struct {
        // Width is the character width of the rain line.
        Width int

        // Values holds the current set of display digits/runes.
        Values []rune

        // LastUpdate tracks when the rain was last refreshed.
        LastUpdate time.Time

        // Paused indicates whether the rain is currently paused (user is interacting).
        Paused bool

        // PauseSince tracks when the pause started.
        PauseSince time.Time

        // rng is a local random source.
        rng *rand.Rand
}

// NewDataRain creates a DataRain of the given width.
func NewDataRain(width int) DataRain {
        dr := DataRain{
                Width:      width,
                LastUpdate: time.Now(),
                rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
        }
        dr.Values = dr.generateInitial()
        return dr
}

// generateInitial creates the starting set of rain values.
func (dr *DataRain) generateInitial() []rune {
        // Pattern: ░░ N N N N N N N ░░ N N N N N ░░ N N N ░░
        // Segments of numbers separated by ░░ delimiters.
        totalLen := dr.Width
        if totalLen <= 0 {
                totalLen = 40
        }

        vals := make([]rune, totalLen)
        i := 0
        for i < totalLen {
                // Add ░░ separator.
                if i < totalLen {
                        vals[i] = '░'
                        i++
                }
                if i < totalLen {
                        vals[i] = '░'
                        i++
                }
                // Add a group of 5-8 random digits with spaces.
                groupLen := 5 + dr.rng.Intn(4)
                for j := 0; j < groupLen && i < totalLen; j++ {
                        if i < totalLen {
                                vals[i] = rune('0' + dr.rng.Intn(10))
                                i++
                        }
                        if i < totalLen && j < groupLen-1 {
                                vals[i] = ' '
                                i++
                        }
                }
        }
        return vals[:min(i, totalLen)]
}

// Tick updates the rain values if enough time has elapsed.
// Returns true if values were updated.
func (dr *DataRain) Tick(now time.Time) bool {
        // Check if paused and should resume.
        if dr.Paused {
                if now.Sub(dr.PauseSince) >= anim.DataRainPauseTimeout {
                        dr.Paused = false
                        dr.LastUpdate = now
                }
                return false
        }

        // Update at the spec interval.
        if now.Sub(dr.LastUpdate) < anim.DataRainUpdateInterval {
                return false
        }

        dr.LastUpdate = now

        // Scroll: shift left, add new values on the right.
        newVals := make([]rune, len(dr.Values))
        copy(newVals, dr.Values[1:])

        // Add new digit or separator on the right.
        if dr.rng.Float64() < 0.3 {
                // 30% chance of ░ separator.
                newVals[len(newVals)-1] = '░'
        } else if dr.rng.Float64() < 0.2 {
                // Space.
                newVals[len(newVals)-1] = ' '
        } else {
                // Random digit.
                newVals[len(newVals)-1] = rune('0' + dr.rng.Intn(10))
        }

        dr.Values = newVals
        return true
}

// Pause pauses the rain (called on user interaction).
func (dr *DataRain) Pause() {
        dr.Paused = true
        dr.PauseSince = time.Now()
}

// Resume unpauses the rain immediately.
func (dr *DataRain) Resume() {
        dr.Paused = false
        dr.LastUpdate = time.Now()
}

// View renders the data rain as a single line.
func (dr DataRain) View() string {
        if len(dr.Values) == 0 {
                return ""
        }

        var b strings.Builder
        for _, v := range dr.Values {
                b.WriteRune(v)
        }

        return lipgloss.NewStyle().Foreground(style.TextDim).Render(b.String())
}

// ViewFormatted renders the data rain with the spec's canonical spacing pattern.
// This produces the exact `░░ 3 7 1 4 8 2 9 ░░` style format.
func (dr DataRain) ViewFormatted() string {
        if len(dr.Values) == 0 {
                return ""
        }

        // Build pattern: ░░ digits ░░ digits ░░ digits ░░
        var segments [][]rune
        var current []rune
        inSep := false

        for _, v := range dr.Values {
                if v == '░' {
                        if !inSep && len(current) > 0 {
                                segments = append(segments, current)
                                current = nil
                        }
                        inSep = true
                        current = append(current, v)
                } else {
                        if inSep && len(current) > 0 {
                                segments = append(segments, current)
                                current = nil
                        }
                        inSep = false
                        current = append(current, v)
                }
        }
        if len(current) > 0 {
                segments = append(segments, current)
        }

        var b strings.Builder
        for _, seg := range segments {
                for _, r := range seg {
                        b.WriteRune(r)
                }
                b.WriteRune(' ')
        }

        return lipgloss.NewStyle().Foreground(style.TextDim).Render(strings.TrimSpace(b.String()))
}

// String implements fmt.Stringer.
func (dr DataRain) String() string {
        return fmt.Sprintf("DataRain{width=%d, paused=%v}", dr.Width, dr.Paused)
}
