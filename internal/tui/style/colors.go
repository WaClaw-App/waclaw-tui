package style

import "github.com/charmbracelet/lipgloss"

// Color tokens from theme.yaml — the only source of truth for visual identity.
var (
        // Background
        Bg       = lipgloss.Color("#0A0A0B") // almost black, soft
        BgRaised = lipgloss.Color("#141416") // subtle elevation
        BgActive = lipgloss.Color("#1A1A1E") // active zone

        // Text
        Text      = lipgloss.Color("#E8E8EC") // primary — warm white
        TextMuted = lipgloss.Color("#6B6B76") // secondary — whisper, not shout
        TextDim   = lipgloss.Color("#3D3D44") // tertiary — barely visible, but readable

        // Signal
        Success = lipgloss.Color("#34D399") // green — not aggressive
        Warning = lipgloss.Color("#FBBF24") // amber — attention, not alarm
        Danger  = lipgloss.Color("#F87171") // red — clear, not scary

        // Brand
        Accent = lipgloss.Color("#818CF8") // indigo — brand, action

        // Motion
        Pulse     = lipgloss.Color("#818CF866") // accent 40% opacity — breathing elements
        Highlight = lipgloss.Color("#FFFFFF22") // white 13% — hover/focus zone

        // Dim signal variants
        DimBlue = lipgloss.Color("#5B7A9D") // cold/frozen — dim blue, not dead

        // Celebration
        Gold        = lipgloss.Color("#FFD700") // jackpot & revenue — earned celebration
        Celebration = lipgloss.Color("#FFFFFF") // full-screen flash — conversion only
)

// Pre-built styles for common text patterns.
// Badge-specific styles live in typography.go — not here.
var (
        PrimaryStyle  = lipgloss.NewStyle().Foreground(Text).Bold(true)
        MutedStyle    = lipgloss.NewStyle().Foreground(TextMuted)
        DimStyle      = lipgloss.NewStyle().Foreground(TextDim)
        AccentStyle   = lipgloss.NewStyle().Foreground(Accent)
        SuccessStyle  = lipgloss.NewStyle().Foreground(Success)
        WarningStyle  = lipgloss.NewStyle().Foreground(Warning)
        DangerStyle   = lipgloss.NewStyle().Foreground(Danger)
        GoldStyle     = lipgloss.NewStyle().Foreground(Gold).Bold(true)
        DimBlueStyle  = lipgloss.NewStyle().Foreground(DimBlue)
        BgStyle       = lipgloss.NewStyle().Background(Bg)
        BgRaisedStyle = lipgloss.NewStyle().Background(BgRaised)
        BgActiveStyle = lipgloss.NewStyle().Background(BgActive)
)

// RebuildStyles recreates all pre-built styles from the current color tokens.
// Called by the theme loader after overriding color values for hot-reload.
func RebuildStyles() {
        PrimaryStyle = lipgloss.NewStyle().Foreground(Text).Bold(true)
        MutedStyle = lipgloss.NewStyle().Foreground(TextMuted)
        DimStyle = lipgloss.NewStyle().Foreground(TextDim)
        AccentStyle = lipgloss.NewStyle().Foreground(Accent)
        SuccessStyle = lipgloss.NewStyle().Foreground(Success)
        WarningStyle = lipgloss.NewStyle().Foreground(Warning)
        DangerStyle = lipgloss.NewStyle().Foreground(Danger)
        GoldStyle = lipgloss.NewStyle().Foreground(Gold).Bold(true)
        DimBlueStyle = lipgloss.NewStyle().Foreground(DimBlue)
        BgStyle = lipgloss.NewStyle().Background(Bg)
        BgRaisedStyle = lipgloss.NewStyle().Background(BgRaised)
        BgActiveStyle = lipgloss.NewStyle().Background(BgActive)

        // Also rebuild typography styles since they depend on color tokens.
        rebuildTypography()
}
