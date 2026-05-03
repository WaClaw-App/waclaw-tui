package component

import (
        "fmt"
        "math"
        "time"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
        "github.com/charmbracelet/lipgloss"
)

// BreathingState tracks the opacity phase of a breathing animation.
//
// Visual spec (from doc/15-micro-interactions.md):
//   - Opacity cycle: 0.9 → 1.0 → 0.9
//   - Duration: 4000ms per cycle
//   - Per-item offset: 200ms (so group stats breathe out of sync)
//
// Breathing is used for:
//   - Monitor dashboard stat numbers
//   - Worker rows (each breathes independently)
//   - Shield pulse
//   - Any element that should feel "alive"
type BreathingState struct {
        // PhaseOffset is the per-item offset (e.g. 200ms per stat number).
        PhaseOffset time.Duration

        // CycleDuration is the total cycle duration (default 4000ms from anim.BreathingCycle).
        CycleDuration time.Duration

        // MinOpacity is the minimum opacity in the cycle (default 0.9 from anim.BreathingMinOpacity).
        MinOpacity float64

        // MaxOpacity is the maximum opacity in the cycle (default 1.0 from anim.BreathingMaxOpacity).
        MaxOpacity float64
}

// NewBreathingState creates a BreathingState with the given offset.
// Uses the default 4000ms cycle and 0.9→1.0 opacity range from spec.
func NewBreathingState(offset time.Duration) BreathingState {
        return BreathingState{
                PhaseOffset:   offset,
                CycleDuration: anim.BreathingCycle,
                MinOpacity:    anim.BreathingMinOpacity,
                MaxOpacity:    anim.BreathingMaxOpacity,
        }
}

// Opacity computes the current opacity value [0.9, 1.0] based on
// the current time and this item's phase offset.
func (b BreathingState) Opacity(now time.Time) float64 {
        // Apply offset so each item starts at a different point in the cycle.
        adjusted := now.Add(-b.PhaseOffset)

        // Position in cycle [0.0, 1.0).
        cycleNanos := b.CycleDuration.Nanoseconds()
        if cycleNanos <= 0 {
                return b.MaxOpacity
        }
        progress := float64(adjusted.UnixNano()%cycleNanos) / float64(cycleNanos)

        // Sine wave for natural breathing: 0→0.5→1→0.5→0
        amplitude := (b.MaxOpacity - b.MinOpacity) / 2.0
        midpoint := (b.MaxOpacity + b.MinOpacity) / 2.0
        opacity := midpoint + amplitude*math.Sin(2.0*math.Pi*progress)

        if opacity < b.MinOpacity {
                opacity = b.MinOpacity
        }
        if opacity > b.MaxOpacity {
                opacity = b.MaxOpacity
        }
        return opacity
}

// BreathingGroup manages a set of breathing items with staggered offsets.
// Used when multiple elements breathe independently (e.g. stat numbers).
type BreathingGroup struct {
        // Items holds the breathing states for each item.
        Items []BreathingState

        // OffsetPerItem is the delay between consecutive items (default 200ms from anim).
        OffsetPerItem time.Duration
}

// NewBreathingGroup creates a breathing group for n items.
// Each item gets a 200ms offset from the previous one (from anim.BreathingOffsetPerItem).
func NewBreathingGroup(n int) BreathingGroup {
        items := make([]BreathingState, n)
        for i := 0; i < n; i++ {
                items[i] = NewBreathingState(time.Duration(i) * anim.BreathingOffsetPerItem)
        }
        return BreathingGroup{
                Items:         items,
                OffsetPerItem: anim.BreathingOffsetPerItem,
        }
}

// Opacities returns the current opacity for each item in the group.
func (bg BreathingGroup) Opacities(now time.Time) []float64 {
        opacities := make([]float64, len(bg.Items))
        for i, item := range bg.Items {
                opacities[i] = item.Opacity(now)
        }
        return opacities
}

// RenderBreathing applies the breathing opacity to a string by choosing between
// a "dim" and "bright" style based on the current opacity.
// This is a practical helper for lipgloss-based rendering where
// true opacity isn't available — we interpolate between two colors.
func RenderBreathing(text string, opacity float64, brightColor, dimColor lipgloss.Color) string {
        // Since terminals don't support true opacity, we blend between
        // the dim and bright colors based on the opacity threshold.
        // Above 0.95 → bright color; below → dim color.
        if opacity >= 0.95 {
                return lipgloss.NewStyle().Foreground(brightColor).Bold(true).Render(text)
        }
        return lipgloss.NewStyle().Foreground(dimColor).Render(text)
}

// String implements fmt.Stringer for BreathingState.
func (b BreathingState) String() string {
        return fmt.Sprintf("BreathingState{offset=%v, cycle=%v, range=[%.1f,%.1f]}",
                b.PhaseOffset, b.CycleDuration, b.MinOpacity, b.MaxOpacity)
}
