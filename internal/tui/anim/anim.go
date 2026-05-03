// Package anim contains animation timing constants and easing functions
// shared between the tui root package and the component sub-package.
//
// This package exists to break the circular dependency:
//   - internal/tui depends on internal/tui/component
//   - internal/tui/component cannot import internal/tui (circular)
//   - Both need access to the same animation constants
//
// By placing constants here, both packages can import anim without cycles.
package anim

import (
        "time"
)

// ---------------------------------------------------------------------------
// Animation timing constants — exact values from doc/15-micro-interactions.md.
// All timing is in milliseconds unless otherwise noted.
// ---------------------------------------------------------------------------

const (
        // ScreenTransition is the horizontal slide duration for forward navigation.
        ScreenTransition = 300 * time.Millisecond

        // BackTransition is the reverse slide duration for back navigation.
        BackTransition = 300 * time.Millisecond

        // TabSwitch is the cross-fade + vertical shift duration.
        TabSwitch = 200 * time.Millisecond

        // ScrollSmooth is the smooth scroll duration per 2-line frame.
        ScrollSmooth = 150 * time.Millisecond

        // NotifSlideIn is the notification entry slide duration.
        NotifSlideIn = 250 * time.Millisecond

        // NotifSlideOut is the notification exit duration.
        NotifSlideOut = 200 * time.Millisecond

        // NumberFlash is the number increment flash + scale duration.
        NumberFlash = 200 * time.Millisecond

        // ItemSlideIn is the new item entry slide duration.
        ItemSlideIn = 250 * time.Millisecond

        // ItemSlideOut is the item removal fade + slide duration.
        ItemSlideOut = 300 * time.Millisecond

        // StatusMorph is the status change color morph duration.
        StatusMorph = 400 * time.Millisecond

        // BreathingCycle is the opacity pulse cycle (0.9→1.0→0.9).
        BreathingCycle = 4000 * time.Millisecond

        // SuccessPulse is the success green pulse duration.
        SuccessPulse = 500 * time.Millisecond

        // AttentionFlash is the amber double-flash duration.
        AttentionFlash = 600 * time.Millisecond

        // ErrorGlow is the red edge glow duration.
        ErrorGlow = 800 * time.Millisecond

        // CompletionBar is the full-width bar fill + hold + fade duration.
        CompletionBar = 1000 * time.Millisecond

        // DealDrama is the DEAL flash + particles + bell total duration.
        DealDrama = 2300 * time.Millisecond // Total: shock(200) + reveal(600) + context(700) + settle_hold(800)

        // ConfigErrorBlink is the red underline blink duration.
        ConfigErrorBlink = 600 * time.Millisecond

        // LogoCharDelay is the per-character render delay for the boot logo.
        LogoCharDelay = 8 * time.Millisecond

        // MenuStagger is the stagger delay per menu item.
        MenuStagger = 120 * time.Millisecond

        // TimelineStagger is the stagger delay per timeline entry.
        // Doc spec: "Timeline entries appear sequentially 100ms stagger".
        TimelineStagger = 100 * time.Millisecond

        // ArmyMarchStagger is the stagger delay for army march rows.
        // Doc spec: "army march ▸▸▸ 80ms stagger".
        ArmyMarchStagger = 80 * time.Millisecond

        // FrameTick is the interval for ~60fps animation frame updates.
        // Used by screens that need per-frame rendering (review, response).
        FrameTick = 16 * time.Millisecond

        // ArmyMarch is the total army march animation duration.
        ArmyMarch = 600 * time.Millisecond

        // QRDissolve is the QR code pixel-by-pixel dissolve duration.
        QRDissolve = 400 * time.Millisecond

        // SuccessHold is the hold time after success before auto-transition.
        SuccessHold = 800 * time.Millisecond

        // NerdStatsSlideUp is the nerd stats slide-up duration.
        NerdStatsSlideUp = 150 * time.Millisecond

        // NerdStatsExpand is the nerd stats panel expand duration.
        NerdStatsExpand = 200 * time.Millisecond

        // NerdStatsCollapse is the nerd stats collapse duration.
        NerdStatsCollapse = 150 * time.Millisecond

        // NerdStatsAutoCollapse is the auto-collapse timeout for nerd stats.
        NerdStatsAutoCollapse = 30 * time.Second

        // CmdPaletteSlideDown is the command palette slide-down duration.
        CmdPaletteSlideDown = 150 * time.Millisecond

        // CmdPaletteSelect is the command palette selection slide duration.
        CmdPaletteSelect = 50 * time.Millisecond

        // CmdPaletteClose is the command palette close duration.
        CmdPaletteClose = 100 * time.Millisecond

        // CmdPaletteDebounce is the search input debounce duration.
        CmdPaletteDebounce = 50 * time.Millisecond

        // SlotMachineScroll is the high-value lead name scroll duration.
        SlotMachineScroll = 400 * time.Millisecond

        // JackpotBounce is the jackpot label scale overshoot duration.
        JackpotBounce = 600 * time.Millisecond

        // ColorWavePulse is the gold/amber color wave pulse duration.
        ColorWavePulse = 300 * time.Millisecond

        // JackpotSettle is the auto-settle delay after jackpot.
        JackpotSettle = 2000 * time.Millisecond

        // BatchCascadeStagger is the batch complete cascade stagger delay.
        BatchCascadeStagger = 200 * time.Millisecond

        // ShieldRepairPerPoint is the shield repair animation rate per health point.
        ShieldRepairPerPoint = 50 * time.Millisecond

        // BreathingMinOpacity is the minimum opacity in the breathing cycle.
        BreathingMinOpacity = 0.9

        // BreathingMaxOpacity is the maximum opacity in the breathing cycle.
        BreathingMaxOpacity = 1.0

        // BreathingOffsetPerItem is the per-item offset in group breathing.
        BreathingOffsetPerItem = 200 * time.Millisecond

        // ScaleBumpFactor is the scale bump factor for number increments.
        ScaleBumpFactor = 1.05

        // SuccessPulseMin is the minimum scale for success pulse.
        SuccessPulseMin = 1.0

        // SuccessPulseMax is the maximum scale for success pulse.
        SuccessPulseMax = 1.2

        // JackpotOvershoot is the overshoot scale for jackpot bounce.
        JackpotOvershoot = 1.3

        // ParticleCount is the number of particles in conversion/jackpot effects.
        ParticleCount = 40

        // ParticleLifetime is the lifetime of a single particle.
        ParticleLifetime = 600 * time.Millisecond

        // ConversionShockDuration is the SHOCK phase of conversion drama.
        ConversionShockDuration = 200 * time.Millisecond

        // ConversionRevealDuration is the REVEAL phase of conversion drama.
        ConversionRevealDuration = 600 * time.Millisecond

        // ConversionContextDuration is the CONTEXT phase of conversion drama.
        ConversionContextDuration = 700 * time.Millisecond

        // ConversionSettleHold is the hold before keyboard accepted in SETTLE phase.
        ConversionSettleHold = 800 * time.Millisecond

        // NotifSeverityCritical is the auto-dismiss time for critical notifications.
        NotifSeverityCritical = 3 * time.Second

        // NotifSeverityPositive is the auto-dismiss time for positive notifications.
        NotifSeverityPositive = 10 * time.Second

        // NotifSeverityNeutral is the auto-dismiss time for neutral notifications.
        NotifSeverityNeutral = 5 * time.Second

        // NotifSeverityInformative is the auto-dismiss time for informative notifications.
        NotifSeverityInformative = 7 * time.Second

        // NotifTypeUpdateAvailable is the per-type override for UpdateAvailable (doc: 15s, not 10s).
        NotifTypeUpdateAvailable = 15 * time.Second

        // NotifTypeUpgradeAvailable is the per-type override for UpgradeAvailable (doc: 20s, not 7s).
        NotifTypeUpgradeAvailable = 20 * time.Second

        // NotifTypeMultiResponse is the per-type override for MultiResponse (doc: 15s, not 10s).
        NotifTypeMultiResponse = 15 * time.Second

        // GuardrailCleanAutoDismiss is the auto-dismiss for clean validation from boot.
        GuardrailCleanAutoDismiss = 3 * time.Second

        // SearchDebounce is the search input debounce duration (for niche explorer).
        SearchDebounce = 300 * time.Millisecond

        // DataRainUpdateInterval is the data rain background update interval.
        DataRainUpdateInterval = 5 * time.Second

        // DataRainPauseTimeout is the idle time before data rain resumes after pause.
        DataRainPauseTimeout = 10 * time.Second

        // StartupSequenceTotal is the total startup sequence duration.
        StartupSequenceTotal = 1300 * time.Millisecond

        // LogStreamSlideIn is the log entry slide-in duration.
        LogStreamSlideIn = 100 * time.Millisecond

        // LogStreamOverflowFade is the log overflow fade duration.
        LogStreamOverflowFade = 150 * time.Millisecond

        // LogStreamPauseTimeout is the idle time before log stream resumes after pause.
        LogStreamPauseTimeout = 5 * time.Second

        // TypeOutCharDelay is the per-character delay for template type-out preview.
        TypeOutCharDelay = 20 * time.Millisecond

        // ComposeSlideUp is the compose modal slide-up duration.
        ComposeSlideUp = 250 * time.Millisecond

        // ComposePreviewHold is the hold time before enter is accepted in preview.
        ComposePreviewHold = 300 * time.Millisecond

        // MiniChartBarStagger is the per-bar stagger delay for weekly mini charts.
        MiniChartBarStagger = 50 * time.Millisecond

        // FollowUpSendingProgress is the follow-up sending progress bar tick interval.
        FollowUpSendingProgress = 200 * time.Millisecond

        // ChangelogStagger is the per-item stagger delay for changelog reveal animation.
        // From doc/12-screens-update-upgrade.md: "changelog items slide in stagger 50ms".
        ChangelogStagger = 50 * time.Millisecond

        // ShortcutsFadeIn is the shortcuts overlay fade-in duration.
        ShortcutsFadeIn = 150 * time.Millisecond

        // ShortcutsFadeOut is the shortcuts overlay fade-out duration.
        ShortcutsFadeOut = 100 * time.Millisecond

        // NerdStatsMetricRefresh is the interval between metric data refreshes.
        NerdStatsMetricRefresh = 2 * time.Second

        // ItemFadeIn is the opacity ramp duration for newly revealed items.
        // Used for changelog items sliding in with a brief fade transition.
        ItemFadeIn = 200 * time.Millisecond

        // NicheStagger is the per-item stagger delay for niche/category list items.
        // Doc specifies 80ms per item for category/niche lists (vs 120ms for menus).
        NicheStagger = 80 * time.Millisecond

        // ConnectionPulseCycle is the 3s pulse cycle for "● wa nyambung" indicator.
        ConnectionPulseCycle = 3 * time.Second

        // ValidationStep is the per-step delay for the 3-step license validation animation.
        // Doc spec: "spinner smooth 200ms rotation".
        ValidationStep = 200 * time.Millisecond

        // AutoTransitionDelay is the delay before auto-transitioning to the next screen
        // after a successful license validation. Doc spec: "Auto-transition ke next screen setelah 2 detik".
        AutoTransitionDelay = 2 * time.Second
)

// AnimationType represents the category of animation being performed.
type AnimationType int

const (
        AnimSlide    AnimationType = iota // horizontal slide
        AnimPulse                         // opacity/scale pulse
        AnimMorph                         // color/shape morph
        AnimParticle                      // particle scatter
        AnimStagger                       // sequential stagger
        AnimFade                          // fade in/out
)

// AnimationState tracks the progress of an ongoing animation.
type AnimationState struct {
        Type      AnimationType
        StartTime time.Time
        Duration  time.Duration
        Progress  float64 // 0.0 to 1.0
}

// IsComplete returns whether the animation has finished.
func (a *AnimationState) IsComplete() bool { return a.Progress >= 1.0 }

// NewAnimationState creates a new animation state starting now.
func NewAnimationState(animType AnimationType, duration time.Duration) AnimationState {
        return AnimationState{
                Type:      animType,
                StartTime: time.Now(),
                Duration:  duration,
                Progress:  0,
        }
}

// UpdateProgress recalculates the progress based on elapsed time.
func (a *AnimationState) UpdateProgress() {
        elapsed := time.Since(a.StartTime)
        if a.Duration <= 0 {
                a.Progress = 1.0
                return
        }
        a.Progress = float64(elapsed) / float64(a.Duration)
        if a.Progress > 1.0 {
                a.Progress = 1.0
        }
}

// ---------------------------------------------------------------------------
// Easing functions
// ---------------------------------------------------------------------------

// EaseOutCubic applies an ease-out cubic curve to the progress value.
func EaseOutCubic(t float64) float64 { return 1 - (1-t)*(1-t)*(1-t) }

// EaseInOutCubic applies an ease-in-out cubic curve.
func EaseInOutCubic(t float64) float64 {
        if t < 0.5 {
                return 4 * t * t * t
        }
        return 1 - (-2*t+2)*(-2*t+2)*(-2*t+2)/2
}

// EaseOutBack applies an ease-out back curve (for overshoot bounces).
func EaseOutBack(t float64) float64 {
        c1 := 1.70158
        c3 := c1 + 1
        return 1 + c3*(t-1)*(t-1)*(t-1) + c1*(t-1)*(t-1)
}
