// Package overlay implements the five global overlays that appear on top of
// any screen: nerd stats, command palette, notification toast, confirmation
// dialog, and keyboard shortcuts.
//
// Overlays are NOT screens — they are lightweight state containers managed by
// the App's View composition layer. They share Bus, i18n, style, anim, and
// component packages, and follow the same DRY conventions (no hardcoded
// strings, no magic numbers).
package overlay

import (
        "fmt"
        "runtime"
        "strings"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/WaClaw-App/waclaw/internal/tui/component"
        "github.com/WaClaw-App/waclaw/internal/tui/i18n"
        "github.com/WaClaw-App/waclaw/internal/tui/style"
        "github.com/WaClaw-App/waclaw/pkg/protocol"
        "github.com/charmbracelet/lipgloss"
)

// NerdStatsMode represents the current visibility mode of the nerd stats overlay.
//
// Spec (doc/10-global-overlays.md):
//   - Hidden  → backtick once → Minimal
//   - Minimal → backtick again → Expanded
//   - Expanded → backtick again → Hidden
//   - Auto-collapse after 30s of inactivity
type NerdStatsMode int

const (
        NerdStatsHidden   NerdStatsMode = iota
        NerdStatsMinimal                 // 1-line footer
        NerdStatsExpanded                // 5-line panel + live logs
)

// Metric warning/danger thresholds (DRY: single source of truth).
const (
        // WarningThreshold is the percentage (0-1) at which metrics turn amber.
        WarningThreshold = 0.8

        // DangerThreshold is the percentage (0-1) at which bar charts turn red.
        DangerThreshold = 0.9
)

// NerdStats holds the state for the nerd stats overlay.
//
// Spec (doc/10-global-overlays.md):
//   - Minimal: 1-line footer with CPU, RAM, goroutines, DB, uptime
//   - Expanded: 5-line panel with mini bar charts per metric + live system logs
//   - Toggle: hidden → minimal → expanded → hidden (backtick key)
//   - Auto-collapse after 30s
//   - Metrics update every 2s
//   - RAM >80% = warning amber, goroutine >80 = danger red
type NerdStats struct {
        // Mode is the current visibility state.
        Mode NerdStatsMode

        // Width is the available terminal width for rendering.
        Width int

        // LastToggle tracks when the overlay was last toggled (for auto-collapse).
        LastToggle time.Time

        // LastMetricUpdate tracks when metrics were last refreshed.
        LastMetricUpdate time.Time

        // Logs is the live system log stream (shown in expanded mode).
        Logs component.LogStream

        // Metrics holds the latest runtime statistics.
        Metrics RuntimeMetrics

        // Anim tracks the slide/expand/collapse animation state.
        Anim anim.AnimationState
}

// RuntimeMetrics holds system vitals fetched from the Go runtime.
type RuntimeMetrics struct {
        // CPU percent (placeholder — real CPU requires /proc or cgo).
        CPUPercent float64

        // RAM allocated in MB.
        RAMAllocMB float64

        // RAM total system memory in MB (simplified: Sys from runtime).
        RAMTotalMB float64

        // Goroutines is the current number of goroutines.
        Goroutines int

        // GoroutinesMax is the configured max (default 100).
        GoroutinesMax int

        // DBSizeMB is the database file size in MB (placeholder).
        DBSizeMB float64

        // DBMaxMB is the max configured DB size (default 50).
        DBMaxMB float64

        // Uptime is how long the application has been running.
        Uptime time.Duration

        // Version is the application version string.
        Version string
}

// NewNerdStats creates a NerdStats in hidden mode with default settings.
func NewNerdStats() NerdStats {
        return NerdStats{
                Mode:           NerdStatsHidden,
                Logs:           component.NewLogStream(),
                Metrics: RuntimeMetrics{
                        GoroutinesMax: 100,
                        DBMaxMB:       50,
                        Version:       "v0.0.0",
                },
                LastToggle:       time.Now(),
                LastMetricUpdate: time.Now(),
        }
}

// Toggle advances the nerd stats mode: hidden→minimal→expanded→hidden.
// Returns a tea.Cmd for animation tick if the overlay becomes visible.
func (ns *NerdStats) Toggle() {
        ns.LastToggle = time.Now()
        switch ns.Mode {
        case NerdStatsHidden:
                ns.Mode = NerdStatsMinimal
                ns.Anim = anim.NewAnimationState(anim.AnimSlide, anim.NerdStatsSlideUp)
        case NerdStatsMinimal:
                ns.Mode = NerdStatsExpanded
                ns.Anim = anim.NewAnimationState(anim.AnimMorph, anim.NerdStatsExpand)
        case NerdStatsExpanded:
                ns.Mode = NerdStatsHidden
                ns.Anim = anim.NewAnimationState(anim.AnimFade, anim.NerdStatsCollapse)
        }
}

// Hide collapses to hidden state.
func (ns *NerdStats) Hide() {
        if ns.Mode != NerdStatsHidden {
                ns.Mode = NerdStatsHidden
                ns.Anim = anim.NewAnimationState(anim.AnimFade, anim.NerdStatsCollapse)
                ns.LastToggle = time.Now()
        }
}

// IsVisible returns true if the overlay is in minimal or expanded mode.
func (ns NerdStats) IsVisible() bool {
        return ns.Mode != NerdStatsHidden
}

// ShouldAutoCollapse returns true if 30s have elapsed since last toggle.
func (ns NerdStats) ShouldAutoCollapse() bool {
        if ns.Mode == NerdStatsHidden {
                return false
        }
        return time.Since(ns.LastToggle) >= anim.NerdStatsAutoCollapse
}

// ShouldRefreshMetrics returns true if the metric refresh interval has elapsed.
func (ns NerdStats) ShouldRefreshMetrics() bool {
        return time.Since(ns.LastMetricUpdate) >= anim.NerdStatsMetricRefresh
}

// RefreshMetrics updates runtime metrics from Go's runtime package.
func (ns *NerdStats) RefreshMetrics() {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)

        ns.Metrics.RAMAllocMB = float64(m.Alloc) / 1024 / 1024
        ns.Metrics.RAMTotalMB = float64(m.Sys) / 1024 / 1024
        ns.Metrics.Goroutines = runtime.NumGoroutine()
        ns.Metrics.CPUPercent = 0 // Placeholder: real CPU requires external measurement

        ns.LastMetricUpdate = time.Now()
}

// AddLog adds a log entry to the nerd stats log stream.
func (ns *NerdStats) AddLog(source, message string, level component.LogLevel) {
        ns.Logs.Add(component.LogEntry{
                Time:    time.Now(),
                Source:  source,
                Message: message,
                Level:   level,
        })
}

// Tick advances the animation and auto-collapse timer.
func (ns *NerdStats) Tick(now time.Time) {
        ns.Anim.UpdateProgress()
        ns.Logs.Tick(now)

        // Auto-collapse check.
        if ns.ShouldAutoCollapse() {
                ns.Hide()
        }
}

// View renders the nerd stats overlay according to its current mode.
func (ns NerdStats) View() string {
        switch ns.Mode {
        case NerdStatsMinimal:
                return ns.renderMinimal()
        case NerdStatsExpanded:
                return ns.renderExpanded()
        default:
                return ""
        }
}

// renderMinimal produces the 1-line footer.
//
// Spec: "── CPU 12% · RAM 134MB · Goroutines 23 · DB 2.4MB · Uptime 4j 12m ──"
// Color: text_dim (almost invisible but readable).
func (ns NerdStats) renderMinimal() string {
        m := ns.Metrics
        uptime := formatUptime(m.Uptime)

        text := fmt.Sprintf("── CPU %.0f%% · RAM %.0fMB · Goroutines %d · DB %.1fMB · Uptime %s ──",
                m.CPUPercent, m.RAMAllocMB, m.Goroutines, m.DBSizeMB, uptime)

        // Color warnings for thresholds.
        if m.RAMAllocMB/m.RAMTotalMB > WarningThreshold {
                ramText := fmt.Sprintf("RAM %.0fMB", m.RAMAllocMB)
                warningRam := lipgloss.NewStyle().Foreground(style.Warning).Render(ramText)
                text = strings.Replace(text, ramText, warningRam, 1)
        }
        if m.Goroutines > int(float64(m.GoroutinesMax)*WarningThreshold) {
                gorText := fmt.Sprintf("Goroutines %d", m.Goroutines)
                dangerGor := lipgloss.NewStyle().Foreground(style.Danger).Render(gorText)
                text = strings.Replace(text, gorText, dangerGor, 1)
        }

        return lipgloss.NewStyle().Foreground(style.TextDim).Render(text)
}

// renderExpanded produces the 5-line panel with mini bar charts + live logs.
//
// Spec (doc/10-global-overlays.md):
//   ── nerd stats ──────────────────────────────────────
//     CPU         12.3%  ████░░░░░░░░░░░░░░░░
//     RAM         134MB  ██████░░░░░░░░░░░░░░  / 512MB
//     Goroutines  23     ██░░░░░░░░░░░░░░░░░░  / 100
//     DB Size     2.4MB  █░░░░░░░░░░░░░░░░░░░  / 50MB
//     Uptime      4j 12m
//     Version     v1.3.2
//   ── logs ────────────────────────────────────────────
//     ...5 lines from LogStream...
//   ────────────────────────────────────────────────────
func (ns NerdStats) renderExpanded() string {
        m := ns.Metrics
        barWidth := 20

        var lines []string

        // Header.
        lines = append(lines, lipgloss.NewStyle().Foreground(style.TextDim).Render(
                fmt.Sprintf("── %s ──%s", i18n.T("nerd_stats.header"), strings.Repeat("─", max(0, 50-len(i18n.T("nerd_stats.header"))))),
        ))

        // CPU metric.
        cpuPercent := m.CPUPercent
        cpuBar := ns.renderMetricBar(cpuPercent/100, barWidth, false)
        lines = append(lines, fmt.Sprintf("  %-12s %5.1f%%  %s",
                "CPU", cpuPercent, cpuBar))

        // RAM metric.
        ramPercent := float64(0)
        if m.RAMTotalMB > 0 {
                ramPercent = m.RAMAllocMB / m.RAMTotalMB
        }
        ramAtLimit := ramPercent > WarningThreshold
        ramBar := ns.renderMetricBar(ramPercent, barWidth, ramAtLimit)
        lines = append(lines, fmt.Sprintf("  %-12s %5.0fMB  %s  / %.0fMB",
                "RAM", m.RAMAllocMB, ramBar, m.RAMTotalMB))

        // Goroutines metric.
        gorPercent := float64(m.Goroutines) / float64(m.GoroutinesMax)
        gorAtLimit := m.Goroutines > int(float64(m.GoroutinesMax)*WarningThreshold)
        gorBar := ns.renderMetricBar(gorPercent, barWidth, gorAtLimit)
        lines = append(lines, fmt.Sprintf("  %-12s %5d     %s  / %d",
                "Goroutines", m.Goroutines, gorBar, m.GoroutinesMax))

        // DB Size metric.
        dbPercent := float64(0)
        if m.DBMaxMB > 0 {
                dbPercent = m.DBSizeMB / m.DBMaxMB
        }
        dbBar := ns.renderMetricBar(dbPercent, barWidth, false)
        lines = append(lines, fmt.Sprintf("  %-12s %5.1fMB  %s  / %.0fMB",
                "DB Size", m.DBSizeMB, dbBar, m.DBMaxMB))

        // Uptime.
        lines = append(lines, fmt.Sprintf("  %-12s %s",
                "Uptime", formatUptime(m.Uptime)))

        // Version.
        lines = append(lines, fmt.Sprintf("  %-12s %s",
                "Version", m.Version))

        // Logs section.
        lines = append(lines, "")
        lines = append(lines, ns.Logs.ViewHeader())
        lines = append(lines, ns.Logs.View())

        // Footer.
        lines = append(lines, lipgloss.NewStyle().Foreground(style.TextDim).Render(
                "────────────────────────────────────────────────────────",
        ))

        return strings.Join(lines, "\n")
}

// renderMetricBar produces a mini bar chart for a single metric.
// Uses accent color normally, warning when atLimit, and danger when > 90%.
func (ns NerdStats) renderMetricBar(ratio float64, width int, atLimit bool) string {
        if ratio < 0 {
                ratio = 0
        }
        if ratio > 1 {
                ratio = 1
        }

        fillCount := int(float64(width) * ratio)
        emptyCount := width - fillCount

        fillColor := style.Accent
        if ratio > DangerThreshold {
                fillColor = style.Danger
        } else if atLimit {
                fillColor = style.Warning
        }

        filled := lipgloss.NewStyle().Foreground(fillColor).Render(
                strings.Repeat("█", fillCount),
        )
        empty := lipgloss.NewStyle().Foreground(style.TextDim).Render(
                strings.Repeat("░", emptyCount),
        )
        return filled + empty
}

// formatUptime formats a duration in a human-friendly way.
// E.g. "4j 12m" (Indonesian-style short) or "4d 12m" (English).
func formatUptime(d time.Duration) string {
        days := int(d.Hours()) / 24
        hours := int(d.Hours()) % 24
        minutes := int(d.Minutes()) % 60

        switch i18n.GetLocale() {
        case i18n.LocaleID:
                // Indonesian: j = jam (hours), m = menit (minutes), h = hari (days)
                if days > 0 {
                        return fmt.Sprintf("%dh %dj %dm", days, hours, minutes)
                }
                if hours > 0 {
                        return fmt.Sprintf("%dj %dm", hours, minutes)
                }
                return fmt.Sprintf("%dm", minutes)
        default:
                // English: d = days, h = hours, m = minutes
                if days > 0 {
                        return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
                }
                if hours > 0 {
                        return fmt.Sprintf("%dh %dm", hours, minutes)
                }
                return fmt.Sprintf("%dm", minutes)
        }
}

// severityDismissMap is the data-driven mapping from Severity to auto-dismiss duration.
// Single source of truth per doc/14-notification-system.md.
var severityDismissMap = map[protocol.Severity]time.Duration{
        protocol.SeverityCritical:    anim.NotifSeverityCritical,
        protocol.SeverityPositive:    anim.NotifSeverityPositive,
        protocol.SeverityNeutral:     anim.NotifSeverityNeutral,
        protocol.SeverityInformative: anim.NotifSeverityInformative,
}

// notifTypeDismissMap provides per-type overrides that deviate from the
// severity-based default. Doc specifies:
//   - UpdateAvailable (Positive): 15s instead of 10s
//   - UpgradeAvailable (Informative): 20s instead of 7s
//   - MultiResponse (Positive): 15s instead of 10s
var notifTypeDismissMap = map[protocol.NotificationType]time.Duration{
        protocol.NotifUpdateAvailable:  anim.NotifTypeUpdateAvailable,
        protocol.NotifUpgradeAvailable: anim.NotifTypeUpgradeAvailable,
        protocol.NotifMultiResponse:    anim.NotifTypeMultiResponse,
}

// SeverityAutoDismiss returns the auto-dismiss duration for a given severity.
// Centralises the mapping so both the notification overlay and App use the
// same source of truth (anim constants).
func SeverityAutoDismiss(sev protocol.Severity) time.Duration {
        if d, ok := severityDismissMap[sev]; ok {
                return d
        }
        return anim.NotifSeverityNeutral
}

// AutoDismissForType returns the auto-dismiss duration for a specific notification
// type, applying per-type overrides when the doc specifies a different duration
// than the severity default.
func AutoDismissForType(t protocol.NotificationType, sev protocol.Severity) time.Duration {
        if d, ok := notifTypeDismissMap[t]; ok {
                return d
        }
        return SeverityAutoDismiss(sev)
}
