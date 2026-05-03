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

// ParticleSymbols are the characters used for the particle cascade effect.
// Spec: ✦ · ★ ☆ ✧ — 40 particles, 600ms lifetime, gold→amber→white→fade.
var ParticleSymbols = []string{"✦", "·", "★", "☆", "✧"}

// ParticleColorStages defines the color progression for particle lifetime.
// Spec: gold → amber → white → fade (text_dim).
var ParticleColorStages = []lipgloss.Color{
        style.Gold,    // 0-25%: gold
        style.Warning, // 25-50%: amber
        style.Text,    // 50-75%: white/bright
        style.TextDim, // 75-100%: fading out
}

// Particle represents a single particle in the cascade system.
type Particle struct {
        Symbol  string
        X       int
        Y       int
        Stage   int    // index into ParticleColorStages
        Born    time.Time
        Dead    bool
}

// ParticleSystem renders a cascade of particles for conversion/jackpot effects.
//
// Visual spec (from doc/05-screens-monitor-response.md):
//   - Particles: ✦ · ★ ☆ ✧
//   - Count: anim.ParticleCount (40)
//   - Lifetime: anim.ParticleLifetime (600ms) per particle
//   - Color cycle: gold → amber → white → fade
//   - Scatter from center outward
//   - Used in: conversion drama, jackpot reveal
type ParticleSystem struct {
        // Particles holds the active particle set.
        Particles []Particle

        // Width and Height define the rendering area.
        Width  int
        Height int

        // Active indicates whether the system is currently emitting.
        Active bool

        // StartTime is when the particle burst began.
        StartTime time.Time

        // rng is a local random source (not global).
        rng *rand.Rand
}

// NewParticleSystem creates a ParticleSystem ready to burst.
func NewParticleSystem(width, height int) ParticleSystem {
        return ParticleSystem{
                Width:  width,
                Height: height,
                rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
        }
}

// Burst triggers a new particle cascade, resetting all particles.
// Generates anim.ParticleCount (40) particles scattered from center.
func (ps *ParticleSystem) Burst() {
        ps.Active = true
        ps.StartTime = time.Now()
        ps.Particles = make([]Particle, 0, anim.ParticleCount)

        centerX := ps.Width / 2
        centerY := ps.Height / 2

        for i := 0; i < anim.ParticleCount; i++ {
                symbol := ParticleSymbols[ps.rng.Intn(len(ParticleSymbols))]
                // Scatter from center with random offset.
                offsetX := ps.rng.Intn(ps.Width) - centerX
                offsetY := ps.rng.Intn(ps.Height) - centerY
                x := centerX + offsetX
                y := centerY + offsetY

                // Clamp to bounds.
                if x < 0 {
                        x = 0
                }
                if x >= ps.Width {
                        x = ps.Width - 1
                }
                if y < 0 {
                        y = 0
                }
                if y >= ps.Height {
                        y = ps.Height - 1
                }

                ps.Particles = append(ps.Particles, Particle{
                        Symbol: symbol,
                        X:      x,
                        Y:      y,
                        Stage:  0,
                        Born:   ps.StartTime,
                })
        }
}

// Tick advances particle ages and returns true if all particles are dead.
func (ps *ParticleSystem) Tick(now time.Time) bool {
        if !ps.Active {
                return true
        }

        allDead := true
        for i := range ps.Particles {
                p := &ps.Particles[i]
                elapsed := now.Sub(p.Born)
                lifetime := anim.ParticleLifetime

                if elapsed >= lifetime {
                        p.Dead = true
                        p.Stage = len(ParticleColorStages) - 1
                        continue
                }

                allDead = false
                // Calculate stage based on lifetime progress.
                progress := float64(elapsed) / float64(lifetime)
                stageIdx := int(progress * float64(len(ParticleColorStages)))
                if stageIdx >= len(ParticleColorStages) {
                        stageIdx = len(ParticleColorStages) - 1
                }
                p.Stage = stageIdx
        }

        if allDead {
                ps.Active = false
        }
        return allDead
}

// View renders the particle system as a multi-line string.
// Each particle is placed at its (X, Y) position with the appropriate color.
func (ps ParticleSystem) View() string {
        if !ps.Active || len(ps.Particles) == 0 {
                return ""
        }

        // Build a grid of characters.
        grid := make([][]string, ps.Height)
        for y := range grid {
                grid[y] = make([]string, ps.Width)
                for x := range grid[y] {
                        grid[y][x] = " "
                }
        }

        // Place particles on the grid.
        for _, p := range ps.Particles {
                if p.Dead || p.Y < 0 || p.Y >= ps.Height || p.X < 0 || p.X >= ps.Width {
                        continue
                }
                color := ParticleColorStages[p.Stage]
                grid[p.Y][p.X] = lipgloss.NewStyle().Foreground(color).Render(p.Symbol)
        }

        // Join grid into string.
        var lines []string
        for _, row := range grid {
                lines = append(lines, strings.Join(row, ""))
        }
        return strings.Join(lines, "\n")
}

// ViewCompact renders particles as a single-line summary (for inline use).
// Returns a string like "✦ · ★ ☆ ✧" with appropriate colors.
func (ps ParticleSystem) ViewCompact() string {
        if !ps.Active {
                return ""
        }

        var b strings.Builder
        for i, p := range ps.Particles {
                if p.Dead {
                        continue
                }
                if i > 0 && i%8 == 0 {
                        b.WriteString("  ")
                } else if i > 0 {
                        b.WriteString(" ")
                }
                color := ParticleColorStages[p.Stage]
                b.WriteString(lipgloss.NewStyle().Foreground(color).Render(p.Symbol))
        }
        return b.String()
}

// IsComplete returns whether the particle system has finished its animation.
func (ps *ParticleSystem) IsComplete() bool {
        return !ps.Active
}

// ParticleBox renders a decorative particle frame used in conversion screens.
// This is the ╔════╗ boxed variant from the conversion screen spec.
type ParticleBox struct {
        Width int
        Label string // Optional label text rendered below the box (e.g. "celebrating" / "dissolving")
}

// View renders the particle box frame.
func (pb ParticleBox) View() string {
        if pb.Width <= 0 {
                pb.Width = 35
        }

        // Top and bottom border.
        border := fmt.Sprintf("╔%s╗", strings.Repeat("═", pb.Width))
        borderBottom := fmt.Sprintf("╚%s╝", strings.Repeat("═", pb.Width))

        // Particle rows inside the box.
        rows := 3
        var lines []string
        lines = append(lines, lipgloss.NewStyle().Foreground(style.Gold).Render(border))

        for r := 0; r < rows; r++ {
                var b strings.Builder
                b.WriteString("║")
                for i := 0; i < pb.Width; i++ {
                        sym := ParticleSymbols[(r*pb.Width+i)%len(ParticleSymbols)]
                        // Alternate colors: gold, amber, white.
                        stageIdx := (r + i) % len(ParticleColorStages)
                        color := ParticleColorStages[stageIdx]
                        b.WriteString(lipgloss.NewStyle().Foreground(color).Render(sym))
                        if i < pb.Width-1 {
                                b.WriteString(" ")
                        }
                }
                b.WriteString("║")
                lines = append(lines, b.String())
        }

        lines = append(lines, lipgloss.NewStyle().Foreground(style.Gold).Render(borderBottom))

        // Optional label text below the box.
        if pb.Label != "" {
                lines = append(lines, lipgloss.NewStyle().Foreground(style.TextDim).Render(pb.Label))
        }

        return strings.Join(lines, "\n")
}
