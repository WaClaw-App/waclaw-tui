// Package tui re-exports animation constants and types from the anim sub-package.
// The anim package is the canonical source of truth for animation timing values
// so that both the tui package and the component sub-package can import them
// without creating a circular dependency.
package tui

import (
        "time"

        tea "github.com/charmbracelet/bubbletea"

        "github.com/WaClaw-App/waclaw/internal/tui/anim"
)

// Re-export all animation constants from the anim package.
// These are kept here for backward compatibility — existing code in the
// tui package references them directly. New code should prefer importing
// anim directly, especially from sub-packages like component/.
const (
        ScreenTransition        = anim.ScreenTransition
        BackTransition          = anim.BackTransition
        TabSwitch               = anim.TabSwitch
        ScrollSmooth            = anim.ScrollSmooth
        NotifSlideIn            = anim.NotifSlideIn
        NotifSlideOut           = anim.NotifSlideOut
        NumberFlash             = anim.NumberFlash
        ItemSlideIn             = anim.ItemSlideIn
        ItemSlideOut            = anim.ItemSlideOut
        StatusMorph             = anim.StatusMorph
        BreathingCycle          = anim.BreathingCycle
        SuccessPulse            = anim.SuccessPulse
        AttentionFlash          = anim.AttentionFlash
        ErrorGlow               = anim.ErrorGlow
        CompletionBar           = anim.CompletionBar
        DealDrama               = anim.DealDrama
        ConfigErrorBlink        = anim.ConfigErrorBlink
        LogoCharDelay           = anim.LogoCharDelay
        MenuStagger             = anim.MenuStagger
        ArmyMarch               = anim.ArmyMarch
        QRDissolve              = anim.QRDissolve
        SuccessHold             = anim.SuccessHold
        NerdStatsSlideUp        = anim.NerdStatsSlideUp
        NerdStatsExpand         = anim.NerdStatsExpand
        NerdStatsCollapse       = anim.NerdStatsCollapse
        NerdStatsAutoCollapse   = anim.NerdStatsAutoCollapse
        CmdPaletteSlideDown     = anim.CmdPaletteSlideDown
        CmdPaletteSelect        = anim.CmdPaletteSelect
        CmdPaletteClose         = anim.CmdPaletteClose
        CmdPaletteDebounce      = anim.CmdPaletteDebounce
        SlotMachineScroll       = anim.SlotMachineScroll
        JackpotBounce           = anim.JackpotBounce
        ColorWavePulse          = anim.ColorWavePulse
        JackpotSettle           = anim.JackpotSettle
        BatchCascadeStagger     = anim.BatchCascadeStagger
        ShieldRepairPerPoint    = anim.ShieldRepairPerPoint
        BreathingMinOpacity     = anim.BreathingMinOpacity
        BreathingMaxOpacity     = anim.BreathingMaxOpacity
        BreathingOffsetPerItem  = anim.BreathingOffsetPerItem
        ScaleBumpFactor         = anim.ScaleBumpFactor
        SuccessPulseMin         = anim.SuccessPulseMin
        SuccessPulseMax         = anim.SuccessPulseMax
        JackpotOvershoot        = anim.JackpotOvershoot
        ParticleCount           = anim.ParticleCount
        ParticleLifetime        = anim.ParticleLifetime
        ConversionShockDuration = anim.ConversionShockDuration
        ConversionRevealDuration = anim.ConversionRevealDuration
        ConversionContextDuration = anim.ConversionContextDuration
        ConversionSettleHold    = anim.ConversionSettleHold
        NotifSeverityCritical   = anim.NotifSeverityCritical
        NotifSeverityPositive   = anim.NotifSeverityPositive
        NotifSeverityNeutral    = anim.NotifSeverityNeutral
        NotifSeverityInformative = anim.NotifSeverityInformative
        GuardrailCleanAutoDismiss = anim.GuardrailCleanAutoDismiss
        SearchDebounce          = anim.SearchDebounce
        DataRainUpdateInterval  = anim.DataRainUpdateInterval
        StartupSequenceTotal    = anim.StartupSequenceTotal
        LogStreamSlideIn        = anim.LogStreamSlideIn
        LogStreamOverflowFade   = anim.LogStreamOverflowFade
        ChangelogStagger        = anim.ChangelogStagger
        ItemFadeIn              = anim.ItemFadeIn
)

// Re-export types.
type AnimationType = anim.AnimationType

const (
        AnimSlide    = anim.AnimSlide
        AnimPulse    = anim.AnimPulse
        AnimMorph    = anim.AnimMorph
        AnimParticle = anim.AnimParticle
        AnimStagger  = anim.AnimStagger
        AnimFade     = anim.AnimFade
)

// Re-export functions.
var (
        EaseOutCubic    = anim.EaseOutCubic
        EaseInOutCubic  = anim.EaseInOutCubic
        EaseOutBack     = anim.EaseOutBack
)

// AnimationState is re-exported from anim for backward compatibility.
type AnimationState = anim.AnimationState

// NewAnimationState wraps anim.NewAnimationState for backward compatibility.
func NewAnimationState(animType AnimationType, duration time.Duration) AnimationState {
        return anim.NewAnimationState(animType, duration)
}

// UpdateProgress delegates to anim.AnimationState.UpdateProgress.
// This is needed because the re-exported type's methods are already available.

// TickCmd returns a bubbletea command that triggers a tick for animation updates.
// Uses the standard bubbletea tick pattern.
func TickCmd(interval time.Duration) tea.Cmd {
        return tea.Tick(interval, func(t time.Time) tea.Msg {
                return AnimationTickMsg(t)
        })
}

// AnimationTickMsg is the message type for animation frame updates.
type AnimationTickMsg time.Time
