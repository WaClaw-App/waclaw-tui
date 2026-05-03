package component

import (
        "fmt"
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/charmbracelet/lipgloss"
)

// QRDisplayState represents the current state of the QR display.
type QRDisplayState int

const (
        QRWaiting  QRDisplayState = iota // QR visible, waiting for scan
        QRScanned                        // QR being dissolved
        QRSuccess                        // Checkmark shown
        QRExpired                        // QR expired
        QRFailed                         // QR failed
)

// QRDisplay renders a QR code with pixel-dissolve animation on scan.
//
// Visual spec (from doc/01-screens-onboarding-boot-login.md):
//   - QR code displayed for WhatsApp login
//   - On scan: pixel-by-pixel dissolve (400ms)
//   - After dissolve: checkmark bounce overshoot
//   - ● ○ ○ sequential animate for multi-slot
//   - Hold 800ms success before auto-transition
type QRDisplay struct {
        // State is the current display state.
        State QRDisplayState

        // Data is the QR code data string (used to render the code).
        Data string

        // Size is the QR code size in characters (default 21 for v1 QR).
        Size int

        // DissolveProgress tracks the dissolve animation [0.0, 1.0].
        DissolveProgress float64

        // DissolveStart is when the dissolve animation began.
        DissolveStart time.Time

        // Slots shows the multi-slot indicator state (e.g. "● ○ ○").
        Slots string

        // ActiveSlot is the 0-based index of the active slot.
        ActiveSlot int

        // TotalSlots is the number of WA slots.
        TotalSlots int
}

// NewQRDisplay creates a QRDisplay with the given data.
func NewQRDisplay(data string) QRDisplay {
        return QRDisplay{
                State:      QRWaiting,
                Data:       data,
                Size:       21,
                Slots:      "●",
                ActiveSlot: 0,
                TotalSlots: 1,
        }
}

// Scan transitions to the dissolve animation.
func (q *QRDisplay) Scan() {
        q.State = QRScanned
        q.DissolveStart = time.Now()
        q.DissolveProgress = 0
}

// Tick advances the dissolve animation.
func (q *QRDisplay) Tick(now time.Time) {
        if q.State != QRScanned {
                return
        }

        elapsed := now.Sub(q.DissolveStart)
        q.DissolveProgress = float64(elapsed) / float64(anim.QRDissolve)
        if q.DissolveProgress >= 1.0 {
                q.DissolveProgress = 1.0
                q.State = QRSuccess
        }
}

// View renders the QR display based on the current state.
func (q QRDisplay) View() string {
        switch q.State {
        case QRWaiting:
                return q.renderQR()
        case QRScanned:
                return q.renderDissolve()
        case QRSuccess:
                return q.renderSuccess()
        case QRExpired:
                return q.renderExpired()
        case QRFailed:
                return q.renderFailed()
        default:
                return q.renderQR()
        }
}

// renderQR renders the full QR code using block characters.
// Since we don't have a real QR encoder, we render a placeholder pattern.
func (q QRDisplay) renderQR() string {
        return q.renderQRPattern(1.0)
}

// renderDissolve renders the QR code with pixels dissolving away.
func (q QRDisplay) renderDissolve() string {
        return q.renderQRPattern(1.0 - q.DissolveProgress)
}

// renderQRPattern renders a QR-like pattern with the given fill ratio.
// Pixels disappear as fillRatio decreases (for dissolve animation).
func (q QRDisplay) renderQRPattern(fillRatio float64) string {
        size := q.Size
        if size <= 0 {
                size = 21
        }

        // Generate a deterministic QR-like pattern from the data string.
        // In production, this would use a real QR encoder.
        pattern := q.generatePattern(size)

        var lines []string
        for y := 0; y < size; y++ {
                var b strings.Builder
                for x := 0; x < size; x++ {
                        idx := y*size + x
                        if idx < len(pattern) && pattern[idx] {
                                // Check if this pixel should be visible (for dissolve).
                                // Dissolve from center outward.
                                centerX, centerY := size/2, size/2
                                dist := float64((x-centerX)*(x-centerX) + (y-centerY)*(y-centerY))
                                maxDist := float64(centerX*centerX + centerY*centerY)
                                normalizedDist := dist / maxDist

                                if normalizedDist <= fillRatio {
                                        b.WriteString(lipgloss.NewStyle().Foreground(style.Text).Render("██"))
                                } else {
                                        b.WriteString("  ")
                                }
                        } else {
                                b.WriteString("  ")
                        }
                }
                lines = append(lines, b.String())
        }

        // Add slot indicators below the QR.
        slotLine := q.renderSlots()
        if slotLine != "" {
                lines = append(lines, "")
                lines = append(lines, slotLine)
        }

        return strings.Join(lines, "\n")
}

// generatePattern creates a deterministic QR-like boolean pattern.
func (q QRDisplay) generatePattern(size int) []bool {
        pattern := make([]bool, size*size)

        // Position detection patterns (corners) — always present in real QR.
        pdSize := 7
        // Top-left.
        q.placeFinderPattern(pattern, size, 0, 0, pdSize)
        // Top-right.
        q.placeFinderPattern(pattern, size, size-pdSize, 0, pdSize)
        // Bottom-left.
        q.placeFinderPattern(pattern, size, 0, size-pdSize, pdSize)

        // Fill remaining with a pseudo-random pattern based on data hash.
        seed := uint32(0)
        for _, c := range q.Data {
                seed = seed*31 + uint32(c)
        }
        for i := range pattern {
                if !pattern[i] {
                        seed = seed*1103515245 + 12345
                        pattern[i] = (seed>>16)%3 == 0
                }
        }

        return pattern
}

// placeFinderPattern draws a QR finder pattern at the given position.
func (q QRDisplay) placeFinderPattern(pattern []bool, size, ox, oy, s int) {
        for y := 0; y < s; y++ {
                for x := 0; x < s; x++ {
                        idx := (oy+y)*size + (ox + x)
                        if idx >= len(pattern) {
                                continue
                        }
                        // Outer border.
                        if x == 0 || x == s-1 || y == 0 || y == s-1 {
                                pattern[idx] = true
                        } else if x >= 2 && x <= s-3 && y >= 2 && y <= s-3 {
                                // Inner square.
                                pattern[idx] = true
                        } else {
                                pattern[idx] = false
                        }
                }
        }
}

// renderSlots renders the multi-slot WA indicator.
func (q QRDisplay) renderSlots() string {
        if q.TotalSlots <= 1 {
                return ""
        }

        var dots []string
        for i := 0; i < q.TotalSlots; i++ {
                if i == q.ActiveSlot {
                        dots = append(dots, lipgloss.NewStyle().Foreground(style.Success).Render("●"))
                } else {
                        dots = append(dots, lipgloss.NewStyle().Foreground(style.TextDim).Render("○"))
                }
        }
        return strings.Join(dots, " ")
}

// renderSuccess renders the checkmark after successful QR scan.
func (q QRDisplay) renderSuccess() string {
        checkmark := lipgloss.NewStyle().Foreground(style.Success).Bold(true).Render(
                "✓",
        )
        slotLine := q.renderSlots()

        lines := []string{
                "",
                "",
                "     " + checkmark,
                "",
        }
        if slotLine != "" {
                lines = append(lines, "     "+slotLine)
        }
        return strings.Join(lines, "\n")
}

// renderExpired renders the expired QR state.
func (q QRDisplay) renderExpired() string {
        return lipgloss.NewStyle().Foreground(style.Danger).Render(
                "  QR expired — press r to refresh",
        )
}

// renderFailed renders the failed QR state.
func (q QRDisplay) renderFailed() string {
        return lipgloss.NewStyle().Foreground(style.Danger).Render(
                "  ✗ QR failed — press r to retry",
        )
}

// String implements fmt.Stringer.
func (q QRDisplay) String() string {
        return fmt.Sprintf("QRDisplay{state=%d, size=%d, slots=%d}", q.State, q.Size, q.TotalSlots)
}
